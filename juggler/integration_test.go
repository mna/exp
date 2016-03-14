package juggler_test

import (
	"flag"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
	"github.com/PuerkitoBio/exp/juggler/internal/jugglertest"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

var (
	// redis pool
	redisPoolMaxActiveFlag   = flag.Int("intg.redis-pool-max-active", 250, "redis: maximum active connections in the pool")
	redisPoolMaxIdleFlag     = flag.Int("intg.redis-pool-max-idle", 10, "redis: maximum idle connections in the pool")
	redisPoolIdleTimeoutFlag = flag.Duration("intg.redis-pool-idle-timeout", time.Minute, "redis: idle connection timeout")

	// broker configuration
	brokerBlockingTimeoutFlag = flag.Duration("intg.broker-blocking-timeout", 0, "broker: blocking timeout")
	brokerCallCapFlag         = flag.Int("intg.broker-call-cap", 0, "broker: call requests queue capacity")
	brokerResultCapFlag       = flag.Int("intg.broker-result-cap", 0, "broker: results queue capacity")

	// server configuration
	serverReadLimitFlag               = flag.Int64("intg.server-read-limit", 0, "server: read limit in bytes")
	serverReadTimeoutFlag             = flag.Duration("intg.server-read-timeout", 0, "server: read timeout")
	serverWriteLimitFlag              = flag.Int64("intg.server-write-limit", 0, "server: write limit in bytes")
	serverWriteTimeoutFlag            = flag.Duration("intg.server-write-timeout", 0, "server: write timeout")
	serverAcquireWriteLockTimeoutFlag = flag.Duration("intg.server-acquire-write-lock-timeout", time.Second, "server: acquire write lock timeout")

	nCalleesFlag          = flag.Int("intg.ncallees", 10, "number of registered callees")
	nWorkersPerCalleeFlag = flag.Int("intg.nworkers-per-callee", 10, "number of workers per callee")
	nClientsFlag          = flag.Int("intg.nclients", 10, "number of clients")

	nURIsFlag          = flag.Int("intg.nuris", 10, "number of URIs")
	nChannelsFlag      = flag.Int("intg.nchannels", 10, "number of pub-sub channels")
	durationFlag       = flag.Duration("intg.duration", 10*time.Second, "duration of the integration test")
	clientMsgRateFlag  = flag.Duration("intg.client-msg-rate", time.Second, "rate to send client message")
	serverPubRateFlag  = flag.Duration("intg.server-pub-rate", time.Second, "rate to send server publish message")
	thunkDelayFlag     = flag.Duration("intg.thunk-delay", 100*time.Millisecond, "delay of a call")
	callExpireOddsFlag = flag.Int("intg.call-expire-odds", 10, "one chance out of that value for a call to expire")
)

func getIntgConfig() *IntgConfig {
	return &IntgConfig{
		RedisPoolMaxActive:   *redisPoolMaxActiveFlag,
		RedisPoolMaxIdle:     *redisPoolMaxIdleFlag,
		RedisPoolIdleTimeout: *redisPoolIdleTimeoutFlag,

		BrokerBlockingTimeout: *brokerBlockingTimeoutFlag,
		BrokerCallCap:         *brokerCallCapFlag,
		BrokerResultCap:       *brokerResultCapFlag,

		ServerReadLimit:               *serverReadLimitFlag,
		ServerReadTimeout:             *serverReadTimeoutFlag,
		ServerWriteLimit:              *serverWriteLimitFlag,
		ServerWriteTimeout:            *serverWriteTimeoutFlag,
		ServerAcquireWriteLockTimeout: *serverAcquireWriteLockTimeoutFlag,

		NCallees:          *nCalleesFlag,
		NWorkersPerCallee: *nWorkersPerCalleeFlag,
		NClients:          *nClientsFlag,

		NURIs:          *nURIsFlag,
		NChannels:      *nChannelsFlag,
		Duration:       *durationFlag,
		ClientMsgRate:  *clientMsgRateFlag,
		ServerPubRate:  *serverPubRateFlag,
		ThunkDelay:     *thunkDelayFlag,
		CallExpireOdds: *callExpireOddsFlag,
	}
}

// IntgConfig holds the configuration of an integration test execution.
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

	NURIs          int           // number of different URIs to pick from at random
	NChannels      int           // number of different channels to pick from at random
	Duration       time.Duration // duration of the test
	ClientMsgRate  time.Duration // send a message at this rate
	ServerPubRate  time.Duration // publish a server-side event at this rate
	ThunkDelay     time.Duration // artificial delay of the call
	CallExpireOdds int           // one chance out of N that a call will expire (timeout < ThunkDelay)
}

func (conf *IntgConfig) URIs() []string {
	uris := make([]string, conf.NURIs)
	for i := 0; i < conf.NURIs; i++ {
		uris[i] = strconv.Itoa(i)
	}
	return uris
}

type runStats struct {
	Call    int64
	Pub     int64
	Sub     int64
	Unsb    int64
	Exp     int64
	OK      int64
	Err     int64
	Res     int64
	Evnt    int64
	Unknown int64
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping the integration test when the -short flag is set")
	}
	runIntegrationTest(t, getIntgConfig())
}

func incStats(stats *runStats, m msg.Msg) {
	switch m.Type() {
	case msg.CallMsg:
		atomic.AddInt64(&stats.Call, 1)
	case msg.PubMsg:
		atomic.AddInt64(&stats.Pub, 1)
	case msg.SubMsg:
		atomic.AddInt64(&stats.Sub, 1)
	case msg.UnsbMsg:
		atomic.AddInt64(&stats.Unsb, 1)
	case msg.ExpMsg:
		atomic.AddInt64(&stats.Exp, 1)
	case msg.OKMsg:
		atomic.AddInt64(&stats.OK, 1)
	case msg.ErrMsg:
		atomic.AddInt64(&stats.Err, 1)
	case msg.ResMsg:
		atomic.AddInt64(&stats.Res, 1)
	case msg.EvntMsg:
		atomic.AddInt64(&stats.Evnt, 1)
	default:
		atomic.AddInt64(&stats.Unknown, 1)
	}
}

func serverHandler(stats *runStats) juggler.Handler {
	return juggler.HandlerFunc(func(ctx context.Context, c *juggler.Conn, m msg.Msg) {
		incStats(stats, m)
		juggler.ProcessMsg(ctx, c, m)
	})
}

func clientHandler(stats *runStats) juggler.ClientHandler {
	return juggler.ClientHandlerFunc(func(ctx context.Context, c *juggler.Client, m msg.Msg) {
		incStats(stats, m)
	})
}

type runConfig struct {
	conf *IntgConfig

	randSeed      int64
	msgsPerClient [][]msg.Msg
	serverPubs    []msg.Msg

	start, end time.Time
}

func prepareExec(t *testing.T, conf *IntgConfig) *runConfig {
	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))

	nCliMsgs := int(conf.Duration / conf.ClientMsgRate)
	mpc := make([][]msg.Msg, conf.NClients)
	for i := 0; i < conf.NClients; i++ {
		mpc[i] = make([]msg.Msg, nCliMsgs)
		for j := 0; j < nCliMsgs; j++ {
			mpc[i][j] = newClientMsg(t, conf, rnd)
		}
	}

	nSrvMsgs := int(conf.Duration / conf.ServerPubRate)
	sps := make([]msg.Msg, nSrvMsgs)
	for i := 0; i < nSrvMsgs; i++ {
		sps[i] = newServerMsg(t, conf, rnd)
	}

	return &runConfig{
		conf:          conf,
		randSeed:      seed,
		msgsPerClient: mpc,
		serverPubs:    sps,
	}
}

func newServerMsg(t *testing.T, conf *IntgConfig, rnd *rand.Rand) msg.Msg {
	ch := strconv.Itoa(rnd.Intn(conf.NChannels))
	pub, err := msg.NewPub(ch, "server event")
	require.NoError(t, err, "NewPub failed")
	return pub
}

func newClientMsg(t *testing.T, conf *IntgConfig, rnd *rand.Rand) msg.Msg {
	var m msg.Msg

	mt := msg.MessageType(rnd.Intn(int(msg.UnsbMsg)) + 1)
	switch mt {
	case msg.CallMsg:
		uri := strconv.Itoa(rnd.Intn(conf.NURIs))
		exp := rnd.Intn(conf.CallExpireOdds)
		to := 10 * conf.ThunkDelay
		if exp == 0 {
			to = time.Millisecond
		}
		call, err := msg.NewCall(uri, "client call", to)
		require.NoError(t, err, "NewCall failed")
		m = call

	case msg.SubMsg:
		ch := strconv.Itoa(rnd.Intn(conf.NChannels))
		m = msg.NewSub(ch, false)

	case msg.UnsbMsg:
		ch := strconv.Itoa(rnd.Intn(conf.NChannels))
		m = msg.NewUnsb(ch, false)

	case msg.PubMsg:
		ch := strconv.Itoa(rnd.Intn(conf.NChannels))
		pub, err := msg.NewPub(ch, "client pub")
		require.NoError(t, err, "NewPub failed")
		m = pub
	}
	return m
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
	var srvStats runStats
	srv := &juggler.Server{
		CallerBroker: brk,
		PubSubBroker: brk,
		LogFunc:      dbgl.Printf,
		Handler:      serverHandler(&srvStats),

		ReadLimit:               conf.ServerReadLimit,
		ReadTimeout:             conf.ServerReadTimeout,
		WriteLimit:              conf.ServerWriteLimit,
		WriteTimeout:            conf.ServerWriteTimeout,
		AcquireWriteLockTimeout: conf.ServerAcquireWriteLockTimeout,
	}
	upg := &websocket.Upgrader{Subprotocols: juggler.Subprotocols}
	httpsrv := httptest.NewServer(juggler.Upgrade(upg, srv))
	defer httpsrv.Close()

	uris := conf.URIs()
	rc := prepareExec(t, conf)
	thunk := func(cp *msg.CallPayload) (interface{}, error) {
		time.Sleep(conf.ThunkDelay)
		return "ok", nil
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
	var clientStats runStats
	clientStarted := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(conf.NClients)

	for i := 0; i < conf.NClients; i++ {
		go func(i int) {
			defer wg.Done()

			cli, err := juggler.Dial(&websocket.Dialer{}, strings.Replace(httpsrv.URL, "http:", "ws:", 1), nil,
				juggler.SetHandler(clientHandler(&clientStats)),
				juggler.SetLogFunc(dbgl.Printf))
			if err != nil {
				t.Fatalf("Dial failed: %v", err)
			}

			clientStarted <- struct{}{}
			dbgl.Printf("client %d started: %d messages, %s delay", i, len(rc.msgsPerClient[i]), conf.ClientMsgRate)
			for _, m := range rc.msgsPerClient[i] {
				switch m := m.(type) {
				case *msg.Call:
					_, err := cli.Call(m.Payload.URI, m.Payload.Args, m.Payload.Timeout)
					require.NoError(t, err, "Call")
				case *msg.Sub:
					_, err := cli.Sub(m.Payload.Channel, m.Payload.Pattern)
					require.NoError(t, err, "Sub")
				case *msg.Unsb:
					_, err := cli.Unsb(m.Payload.Channel, m.Payload.Pattern)
					require.NoError(t, err, "Unsb")
				case *msg.Pub:
					_, err := cli.Pub(m.Payload.Channel, m.Payload.Args)
					require.NoError(t, err, "Pub")
				}
				<-time.After(conf.ClientMsgRate)
			}
			require.NoError(t, cli.Close(), "Close client %d", i)
			dbgl.Printf("client %d closed", i)
		}(i)
	}

	// wait for callees to come online
	for i, cnt := 0, conf.NCallees*conf.NWorkersPerCallee; i < cnt; i++ {
		<-calleeStarted
	}
	// start clients
	rc.start = time.Now()
	for i := 0; i < conf.NClients; i++ {
		<-clientStarted
	}
	// wait for completion
	wg.Wait()
	rc.end = time.Now()
	dbgl.Printf("done after %s", rc.end.Sub(rc.start))

	checkAndPrintResults(t, rc, &srvStats, &clientStats)
}

func checkAndPrintResults(t *testing.T, rc *runConfig, srv, cli *runStats) {

}
