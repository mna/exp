package juggler_test

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"text/tabwriter"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
	"github.com/PuerkitoBio/exp/juggler/client"
	"github.com/PuerkitoBio/exp/juggler/internal/jugglertest"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
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
	Call       int64
	PubSrvSent int64
	PubCliSent int64
	Sub        int64
	Unsb       int64
	Exp        int64
	OK         int64
	Err        int64
	Res        int64
	Evnt       int64
	Unknown    int64
}

func (s *runStats) clone() *runStats {
	return &runStats{
		Call:       atomic.LoadInt64(&s.Call),
		PubSrvSent: atomic.LoadInt64(&s.PubSrvSent),
		PubCliSent: atomic.LoadInt64(&s.PubCliSent),
		Sub:        atomic.LoadInt64(&s.Sub),
		Unsb:       atomic.LoadInt64(&s.Unsb),
		Exp:        atomic.LoadInt64(&s.Exp),
		OK:         atomic.LoadInt64(&s.OK),
		Err:        atomic.LoadInt64(&s.Err),
		Res:        atomic.LoadInt64(&s.Res),
		Evnt:       atomic.LoadInt64(&s.Evnt),
		Unknown:    atomic.LoadInt64(&s.Unknown),
	}
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping the integration test when the -short flag is set")
	}
	runIntegrationTest(t, getIntgConfig())
}

func incStats(stats *runStats, m msg.Msg, fromSrv bool) {
	switch m.Type() {
	case msg.CallMsg:
		atomic.AddInt64(&stats.Call, 1)
	case msg.PubMsg:
		if fromSrv {
			atomic.AddInt64(&stats.PubSrvSent, 1)
		} else {
			atomic.AddInt64(&stats.PubCliSent, 1)
		}
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

func serverHandler(t *testing.T, brk broker.PubSubBroker, rc *runConfig, stats *runStats) juggler.Handler {
	var once sync.Once
	return juggler.HandlerFunc(func(ctx context.Context, c *juggler.Conn, m msg.Msg) {
		incStats(stats, m, m.Type().IsWrite())
		juggler.ProcessMsg(ctx, c, m)

		// start sending PUB messages at the first received message
		once.Do(func() {
			go func() {
				for _, m := range rc.serverPubs {
					err := brk.Publish(m.Payload.Channel, &msg.PubPayload{
						MsgUUID: uuid.NewRandom(),
						Args:    m.Payload.Args,
					})
					require.NoError(t, err, "Publish failed")
					incStats(stats, m, true)
				}
			}()
		})
	})
}

func clientHandler(stats *runStats) client.Handler {
	return client.HandlerFunc(func(ctx context.Context, c *client.Client, m msg.Msg) {
		incStats(stats, m, false)
	})
}

type runConfig struct {
	conf *IntgConfig

	rnd           *rand.Rand
	randSeed      int64
	msgsPerClient [][]msg.Msg
	serverPubs    []*msg.Pub
	expectedExp   int

	start, end time.Time
}

func prepareExec(t *testing.T, conf *IntgConfig) *runConfig {
	var expectedExp int
	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))

	nCliMsgs := int(conf.Duration / conf.ClientMsgRate)
	mpc := make([][]msg.Msg, conf.NClients)
	for i := 0; i < conf.NClients; i++ {
		mpc[i] = make([]msg.Msg, nCliMsgs)
		for j := 0; j < nCliMsgs; j++ {
			m, expires := newClientMsg(t, conf, rnd)
			mpc[i][j] = m
			if expires {
				expectedExp++
			}
		}
	}

	nSrvMsgs := int(conf.Duration / conf.ServerPubRate)
	sps := make([]*msg.Pub, nSrvMsgs)
	for i := 0; i < nSrvMsgs; i++ {
		sps[i] = newServerMsg(t, conf, rnd)
	}

	return &runConfig{
		rnd:           rnd,
		conf:          conf,
		randSeed:      seed,
		msgsPerClient: mpc,
		serverPubs:    sps,
		expectedExp:   expectedExp,
	}
}

func newServerMsg(t *testing.T, conf *IntgConfig, rnd *rand.Rand) *msg.Pub {
	ch := strconv.Itoa(rnd.Intn(conf.NChannels))
	pub, err := msg.NewPub(ch, "server event")
	require.NoError(t, err, "NewPub failed")
	return pub
}

func newClientMsg(t *testing.T, conf *IntgConfig, rnd *rand.Rand) (m msg.Msg, expires bool) {
	mt := msg.MessageType(rnd.Intn(int(msg.UnsbMsg)) + 1)
	switch mt {
	case msg.CallMsg:
		uri := strconv.Itoa(rnd.Intn(conf.NURIs))
		exp := rnd.Intn(conf.CallExpireOdds)
		to := 10 * conf.ThunkDelay
		if exp == 0 {
			expires = true
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
	return m, expires
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
	rc := prepareExec(t, conf)

	var srvStats runStats
	srv := &juggler.Server{
		CallerBroker: brk,
		PubSubBroker: brk,
		LogFunc:      dbgl.Printf,
		Handler:      serverHandler(t, brk, rc, &srvStats),

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
	thunk := func(cp *msg.CallPayload) (interface{}, error) {
		time.Sleep(conf.ThunkDelay)
		return "ok", nil
	}

	// 4. start m callees
	calleeStarted := make(chan struct{})
	for i := 0; i < conf.NCallees; i++ {
		go func(i int) {
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

			wg := sync.WaitGroup{}
			wg.Add(conf.NWorkersPerCallee)
			dbgl.Printf("starting callee %d", i)
			for j := 0; j < conf.NWorkersPerCallee; j++ {
				go func() {
					defer wg.Done()

					calleeStarted <- struct{}{}
					for cp := range ch {
						if err := cle.InvokeAndStoreResult(cp, thunk); err != nil && err != callee.ErrCallExpired {
							t.Fatalf("InvokeAndStoreResult failed: %v", err)
						}
					}
				}()
			}
			wg.Wait()
			dbgl.Printf("stopping callee %d", i)
		}(i)
	}

	// 5. start n clients
	var clientStats runStats
	clientStarted := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(conf.NClients)

	for i := 0; i < conf.NClients; i++ {
		go func(i int) {
			defer wg.Done()

			cli, err := client.Dial(&websocket.Dialer{}, strings.Replace(httpsrv.URL, "http:", "ws:", 1), nil,
				client.SetHandler(clientHandler(&clientStats)),
				client.SetLogFunc(dbgl.Printf))

			clientStarted <- struct{}{}
			if err != nil {
				t.Errorf("Dial failed: %v", err)
				return
			}

			dbgl.Printf("client %d started: %d messages, %s delay", i, len(rc.msgsPerClient[i]), conf.ClientMsgRate)
			for _, m := range rc.msgsPerClient[i] {
				incStats(&clientStats, m, false)
				switch m := m.(type) {
				case *msg.Call:
					_, err := cli.Call(m.Payload.URI, m.Payload.Args, m.Payload.Timeout)
					if !assert.NoError(t, err, "Call") {
						return
					}
				case *msg.Sub:
					_, err := cli.Sub(m.Payload.Channel, m.Payload.Pattern)
					if !assert.NoError(t, err, "Sub") {
						return
					}
				case *msg.Unsb:
					_, err := cli.Unsb(m.Payload.Channel, m.Payload.Pattern)
					if !assert.NoError(t, err, "Unsb") {
						return
					}
				case *msg.Pub:
					_, err := cli.Pub(m.Payload.Channel, m.Payload.Args)
					if !assert.NoError(t, err, "Pub") {
						return
					}
				}
				<-time.After(conf.ClientMsgRate)
			}
			// wait some time for pending responses and potential EVNT messages
			<-time.After(conf.Duration / 10)

			require.NoError(t, cli.Close(), "Close client %d", i)
			dbgl.Printf("client %d closed", i)
		}(i)
	}

	// wait for callees to come online
	for i, cnt := 0, conf.NCallees*conf.NWorkersPerCallee; i < cnt; i++ {
		<-calleeStarted
	}
	// start clients with some jitter
	rc.start = time.Now()
	for i := 0; i < conf.NClients; i++ {
		<-time.After(time.Duration(rc.rnd.Intn(100)) * time.Millisecond)
		<-clientStarted
	}
	// wait for completion
	wg.Wait()
	rc.end = time.Now()

	checkAndPrintResults(t, rc, srvStats.clone(), clientStats.clone())
}

func checkAndPrintResults(t *testing.T, rc *runConfig, srv, cli *runStats) {
	runTime := rc.end.Sub(rc.start)
	assert.True(t, runTime > rc.conf.Duration, "Duration")

	assert.Equal(t, cli.Call, srv.Call, "Call")
	assert.Equal(t, cli.Sub, srv.Sub, "Sub")
	assert.Equal(t, cli.Unsb, srv.Unsb, "Unsb")
	assert.Equal(t, cli.PubCliSent, srv.PubCliSent, "Pub")
	assert.Equal(t, len(rc.serverPubs), int(srv.PubSrvSent), "Pub (server)")

	assert.Equal(t, rc.expectedExp, int(cli.Exp), "Exp (clients)")
	assert.Equal(t, 0, int(srv.Exp), "Exp (server)")
	expRes := int(cli.Call) - rc.expectedExp
	assert.Equal(t, expRes, int(cli.Res), "Res (clients)")
	assert.Equal(t, expRes, int(srv.Res), "Res (server)")

	var cntMsgs int
	for _, mpc := range rc.msgsPerClient {
		cntMsgs += len(mpc)
	}
	assert.Equal(t, cntMsgs, int(cli.OK), "OK (clients)")
	assert.Equal(t, cntMsgs, int(srv.OK), "OK (server)")

	assert.Equal(t, 0, int(cli.Err), "Err (clients)")
	assert.Equal(t, 0, int(srv.Err), "Err (server)")
	assert.Equal(t, cli.Evnt, srv.Evnt, "Evnt")

	assert.Equal(t, 0, int(cli.Unknown), "Unknown (clients)")
	assert.Equal(t, 0, int(srv.Unknown), "Unknown (server)")

	if testing.Verbose() {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "--- RESULTS")
		fmt.Fprintln(os.Stdout)

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)
		fmt.Fprintf(w, "• Duration\t%s\n", runTime)
		fmt.Fprintf(w, "• Random seed\t%d\n", rc.randSeed)
		fmt.Fprintf(w, "• Callees\t%d x %d\n", rc.conf.NCallees, rc.conf.NWorkersPerCallee)
		fmt.Fprintf(w, "• URIs\t%d\n", rc.conf.NURIs)
		fmt.Fprintf(w, "• Channels\t%d\n", rc.conf.NChannels)
		fmt.Fprintf(w, "• Clients\t%d\n", rc.conf.NClients)

		mpc := 0
		if len(rc.msgsPerClient) > 0 {
			mpc = len(rc.msgsPerClient[0])
		}
		fmt.Fprintf(w, "• Expected messages\t%d x %d\n", rc.conf.NClients, mpc)
		fmt.Fprintf(w, "• Expected expired calls\t%d\n", rc.expectedExp)
		w.Flush()

		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(w, "Stats\tServer\tClients\n")
		fmt.Fprintf(w, "• Calls\t%d\t%d\n", srv.Call, cli.Call)
		fmt.Fprintf(w, "• Sub\t%d\t%d\n", srv.Sub, cli.Sub)
		fmt.Fprintf(w, "• Unsb\t%d\t%d\n", srv.Unsb, cli.Unsb)
		fmt.Fprintf(w, "• Pub (from server)\t%d\t%d\n", srv.PubSrvSent, cli.PubSrvSent)
		fmt.Fprintf(w, "• Pub (from client)\t%d\t%d\n", srv.PubCliSent, cli.PubCliSent)
		fmt.Fprintf(w, "• Exp\t%d\t%d\n", srv.Exp, cli.Exp)
		fmt.Fprintf(w, "• Res\t%d\t%d\n", srv.Res, cli.Res)
		fmt.Fprintf(w, "• OK\t%d\t%d\n", srv.OK, cli.OK)
		fmt.Fprintf(w, "• Err\t%d\t%d\n", srv.Err, cli.Err)
		fmt.Fprintf(w, "• Evnt\t%d\t%d\n", srv.Evnt, cli.Evnt)
		fmt.Fprintf(w, "• Unknown\t%d\t%d\n", srv.Unknown, cli.Unknown)
		w.Flush()
		fmt.Fprintln(os.Stdout)
	}
}
