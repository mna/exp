package juggler_test

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
	"github.com/PuerkitoBio/exp/juggler/internal/jugglertest"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
)

type IntgConfig struct {
	RedisPoolMaxActive   int
	RedisPoolMaxIdle     int
	RedisPoolIdleTimeout time.Duration

	BrokerBlockingTimeout time.Duration
	BrokerCallCap         int
	BrokerResultCap       int

	ServerReadLimit               int64
	ServerReadTimeout             time.Duration
	ServerWriteLimit              int64
	ServerWriteTimeout            time.Duration
	ServerAcquireWriteLockTimeout time.Duration

	NCallees          int // number of callees listening on the URIs
	NWorkersPerCallee int // number of workers per callee
	NClients          int // number of clients generating calls

	NURIs         int           // number of different URIs to pick from at random
	NChannels     int           // number of different channels to pick from at random
	Duration      time.Duration // duration of the test
	ClientMsgRate time.Duration // send a message at this rate
	ServerPubRate time.Duration // publish a server-side event at this rate
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests don't run with the -short flag")
	}
	runIntegrationTest(t, &IntgConfig{}) // TODO : parse flags into IntgConfig
}

func runIntegrationTest(t *testing.T, conf *IntgConfig) {
	dbgl := &jugglertest.DebugLog{T: t}

	// start/create:
	// 1. redis-server
	// 2. redis pool and broker
	// 3. juggler server
	// 4. m callees
	// 5. n clients

	// 1. redis-server
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	// 2. create the redis pool and broker
	pool := redistest.NewPool(t, ":"+port)
	pool.MaxActive = conf.RedisPoolMaxActive
	pool.MaxIdle = conf.RedisPoolMaxIdle
	pool.IdleTimeout = conf.RedisPoolIdleTimeout
	brk := &redisbroker.Broker{
		Pool:    pool,
		Dial:    pool.Dial,
		LogFunc: dbgl.Printf,

		BlockingTimeout: conf.BrokerBlockingTimeout,
		CallCap:         conf.BrokerCallCap,
		ResultCap:       conf.BrokerResultCap,
	}

	// 3. create the juggler server
	srv := &juggler.Server{
		CallerBroker: brk,
		PubSubBroker: brk,
		LogFunc:      dbgl.Printf,

		// TODO : set those to something that can keep track of metrics/correctness
		ReadHandler:  nil,
		WriteHandler: nil,

		ReadLimit:               conf.ServerReadLimit,
		ReadTimeout:             conf.ServerReadTimeout,
		WriteLimit:              conf.ServerWriteLimit,
		WriteTimeout:            conf.ServerWriteTimeout,
		AcquireWriteLockTimeout: conf.ServerAcquireWriteLockTimeout,
	}
	upg := &websocket.Upgrader{Subprotocols: juggler.Subprotocols}
	httpsrv := httptest.NewServer(juggler.Upgrade(upg, srv))
	defer httpsrv.Close()

	// TODO : get URIs
	uris := []string{}
	thunk := func(cp *msg.CallPayload) (interface{}, error) {
		return nil, nil
	}

	// 4. start m callees
	calleeStarted := make(chan struct{})
	for i := 0; i < conf.NCallees; i++ {
		go func() {
			cle := callee.Callee{
				Broker:  brk,
				LogFunc: dbgl.Printf,
			}

			conn, err := brk.Calls(uris...)
			if err != nil {
				t.Fatalf("failed to get CallsConn: %v", err)
			}
			defer conn.Close()
			ch := conn.Calls()

			for j := 0; j < conf.NWorkersPerCallee; j++ {
				go func() {
					calleeStarted <- struct{}{}
					for cp := range ch {
						if err := cle.InvokeAndStoreResult(cp, thunk); err != nil {
							t.Fatalf("InvokeAndStoreResult failed: %v", err)
						}
					}
				}()
			}
		}()
	}

	// 5. start n clients
	clientStarted := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(conf.NClients)
	for i := 0; i < conf.NClients; i++ {
		go func() {
			defer wg.Done()

			cli, err := juggler.Dial(&websocket.Dialer{}, strings.Replace(httpsrv.URL, "http:", "ws:", 1), nil)
			if err != nil {
				t.Fatalf("Dial failed: %v", err)
			}
			// juggler.SetHandler()) // TODO : set to something that keeps track of metrics/correctness
			_ = cli

			// TODO : run predetermined requests...
			clientStarted <- struct{}{}
		}()
	}

	// wait for callees to come online
	for i, cnt := 0, conf.NCallees*conf.NWorkersPerCallee; i < cnt; i++ {
		<-calleeStarted
	}
	// start clients
	for i := 0; i < conf.NClients; i++ {
		<-clientStarted
	}
	// wait for completion
	wg.Wait()
}
