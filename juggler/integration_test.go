package juggler_test

import (
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/internal/jugglertest"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
)

type IntgConfig struct {
	BrokerBlockingTimeout time.Duration
	BrokerCallCap         int
	BrokerResultCap       int

	ServerReadLimit               int
	ServerReadTimeout             time.Duration
	ServerWriteLimit              int
	ServerWriteTimeout            time.Duration
	ServerAcquireWriteLockTimeout time.Duration

	NCallees int
	NClients int
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests don't run with the -short flag")
	}
	runIntegrationTests(t, &IntgConfig{}) // TODO : parse flags into IntgConfig
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
	_ = srv

	// 4. start m callees
	stopCallees := make(chan struct{})
	wg := sync.WaitGroup{}
}
