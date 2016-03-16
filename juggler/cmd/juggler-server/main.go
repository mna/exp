// Command juggler-server implements a juggler server that listens for
// connections and serves the requests.
package main

import (
	"errors"
	_ "expvar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"gopkg.in/yaml.v2"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

var (
	redisAddrFlag       = flag.String("redis", ":6379", "Redis `address`.")
	allowEmptyProtoFlag = flag.Bool("allow-empty-subprotocol", false, "Allow empty subprotocol during handshake.")
	portFlag            = flag.Int("port", 9000, "Server `port`.")
	configFlag          = flag.String("config", "", "Path of the configuration `file`.")
	helpFlag            = flag.Bool("help", false, "Show help.")
)

// Redis defines the redis-specific configuration options.
type Redis struct {
	Addr        string        `yaml:"addr"`
	MaxActive   int           `yaml:"max_active"`
	MaxIdle     int           `yaml:"max_idle"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	PubSub      *Redis        `yaml:"pubsub"`
	Caller      *Redis        `yaml:"caller"`
}

// CallerBroker defines the configuration options for the caller broker.
type CallerBroker struct {
	BlockingTimeout time.Duration `yaml:"blocking_timeout"`
	CallCap         int           `yaml:"call_cap"`
}

// Server defines the juggler server configuration options.
type Server struct {
	// HTTP server configuration for the websocket handshake/upgrade
	Addr               string        `yaml:"addr"`
	Paths              []string      `yaml:"paths"`
	MaxHeaderBytes     int           `yaml:"max_header_bytes"`
	ReadBufferSize     int           `yaml:"read_buffer_size"`
	WriteBufferSize    int           `yaml:"write_buffer_size"`
	HandshakeTimeout   time.Duration `yaml:"handshake_timeout"`
	WhitelistedOrigins []string      `yaml:"whitelisted_origins"`

	// websocket/juggler configuration
	ReadLimit               int64         `yaml:"read_limit"`
	ReadTimeout             time.Duration `yaml:"read_timeout"`
	WriteLimit              int64         `yaml:"write_limit"`
	WriteTimeout            time.Duration `yaml:"write_timeout"`
	AcquireWriteLockTimeout time.Duration `yaml:"acquire_write_lock_timeout"`
	AllowEmptySubprotocol   bool          `yaml:"allow_empty_subprotocol"`

	// handler options
	CloseURI string `yaml:"close_uri"`
}

// Config defines the configuration options of the server.
type Config struct {
	Redis        *Redis        `yaml:"redis"`
	CallerBroker *CallerBroker `yaml:"caller_broker"`
	Server       *Server       `yaml:"server"`
}

func getDefaultConfig() *Config {
	return &Config{
		Redis: &Redis{
			Addr:        *redisAddrFlag,
			MaxActive:   0,
			MaxIdle:     0,
			IdleTimeout: 0,
		},
		CallerBroker: &CallerBroker{
			BlockingTimeout: 0,
			CallCap:         0,
		},
		Server: &Server{
			Addr:                    ":" + strconv.Itoa(*portFlag),
			Paths:                   []string{"/ws"},
			ReadLimit:               0,
			ReadTimeout:             0,
			WriteLimit:              0,
			WriteTimeout:            0,
			AcquireWriteLockTimeout: 0,
			AllowEmptySubprotocol:   *allowEmptyProtoFlag,
			CloseURI:                "",
		},
	}
}

func getConfigFromReader(r io.Reader) (*Config, error) {
	conf := getDefaultConfig()

	// set default values
	if r != nil {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, conf); err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func getConfigFromFile(file string) (*Config, error) {
	var r io.Reader
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r = f
	}
	return getConfigFromReader(r)
}

var zeroRedis = Redis{}

func isZeroRedis(rc *Redis) bool {
	if rc == nil {
		return true
	}

	// nil the pubsub and caller
	copy := *rc
	copy.PubSub = nil
	copy.Caller = nil
	return copy == zeroRedis
}

// check redis configuration: use Config.Redis to use the same pool
// for pubsub and caller, or use Config.Redis.PubSub and Config.Redis.Caller.
// No other combination is accepted.
func checkRedisConfig(conf *Redis) error {
	// if either PubSub or Caller is set, then both must be set
	if !isZeroRedis(conf.PubSub) || !isZeroRedis(conf.Caller) {
		if (conf.PubSub == nil || conf.PubSub.Addr == "") || (conf.Caller == nil || conf.Caller.Addr == "") {
			return errors.New("both redis.pubsub and redis.caller sections must be configured")
		}

		// and the generic redis must not be set
		if conf.Addr == *redisAddrFlag {
			conf.Addr = ""
		}
		if !isZeroRedis(conf) {
			return errors.New("redis must not be configured if redis.pubsub and redis.caller are configured")
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

	conf, err := getConfigFromFile(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration file: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if err := checkRedisConfig(conf.Redis); err != nil {
		fmt.Fprintf(os.Stderr, "invalid redis configuration: %v\n", err)
		flag.Usage()
		os.Exit(3)
	}

	// create pool, brokers, server, upgrader, HTTP server
	var poolp, poolc *redis.Pool
	if conf.Redis.Addr != "" {
		pool := newRedisPool(conf.Redis)
		poolp, poolc = pool, pool
		log.Printf("redis pool configured on %s", conf.Redis.Addr)
	} else {
		poolp = newRedisPool(conf.Redis.PubSub)
		poolc = newRedisPool(conf.Redis.Caller)
		log.Printf("redis pool configured on %s (pubsub) and %s (caller)", conf.Redis.PubSub.Addr, conf.Redis.Caller.Addr)
	}

	psb := newPubSubBroker(poolp)
	cb := newCallerBroker(conf.CallerBroker, poolc)

	srv := newServer(conf.Server, psb, cb)
	srv.Handler = newHandler(conf.Server)

	upg := newUpgrader(conf.Server) // must be after newServer, for Subprotocols

	upgh := juggler.Upgrade(upg, srv)
	for _, p := range conf.Server.Paths {
		http.Handle(p, upgh)
	}

	httpSrv := newHTTPServer(conf.Server)

	log.Printf("listening for connections on %s", conf.Server.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}

func newHandler(conf *Server) juggler.Handler {
	closeURI := conf.CloseURI
	writeTimeout := conf.WriteTimeout

	process := juggler.HandlerFunc(func(ctx context.Context, c *juggler.Conn, m msg.Msg) {
		if call, ok := m.(*msg.Call); ok {
			if call.Payload.URI == closeURI {
				wsc := c.UnderlyingConn()

				deadline := time.Now().Add(writeTimeout)
				if writeTimeout == 0 {
					deadline = time.Time{}
				}

				if err := wsc.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"),
					deadline); err != nil {

					log.Printf("WriteControl failed: %v", err)
				}
				return
			}
		}
		juggler.ProcessMsg(ctx, c, m)
	})

	return juggler.PanicRecover(
		juggler.Chain(
			juggler.HandlerFunc(juggler.LogMsg),
			process,
		), true, true)
}

func newPubSubBroker(pool *redis.Pool) broker.PubSubBroker {
	return &redisbroker.Broker{
		Pool: pool,
		Dial: pool.Dial,
	}
}

func newCallerBroker(conf *CallerBroker, pool *redis.Pool) broker.CallerBroker {
	return &redisbroker.Broker{
		Pool:            pool,
		Dial:            pool.Dial,
		BlockingTimeout: conf.BlockingTimeout,
		CallCap:         conf.CallCap,
	}
}

func isIn(list []string, v string) bool {
	for _, vv := range list {
		if v == vv {
			return true
		}
	}
	return false
}

func newUpgrader(conf *Server) *websocket.Upgrader {
	upg := &websocket.Upgrader{
		HandshakeTimeout: conf.HandshakeTimeout,
		ReadBufferSize:   conf.ReadBufferSize,
		WriteBufferSize:  conf.WriteBufferSize,
		Subprotocols:     juggler.Subprotocols,
	}

	if len(conf.WhitelistedOrigins) > 0 {
		oris := conf.WhitelistedOrigins
		upg.CheckOrigin = func(r *http.Request) bool {
			o := r.Header.Get("Origin")
			return isIn(oris, o)
		}
	}
	return upg
}

func newHTTPServer(conf *Server) *http.Server {
	return &http.Server{
		Addr:           conf.Addr,
		ReadTimeout:    conf.ReadTimeout,
		WriteTimeout:   conf.WriteTimeout,
		MaxHeaderBytes: conf.MaxHeaderBytes,
	}
}

func newServer(conf *Server, pubSub broker.PubSubBroker, caller broker.CallerBroker) *juggler.Server {
	if conf.AllowEmptySubprotocol {
		juggler.Subprotocols = append(juggler.Subprotocols, "")
	}

	return &juggler.Server{
		ReadLimit:               conf.ReadLimit,
		ReadTimeout:             conf.ReadTimeout,
		WriteLimit:              conf.WriteLimit,
		WriteTimeout:            conf.WriteTimeout,
		AcquireWriteLockTimeout: conf.AcquireWriteLockTimeout,
		ConnState:               juggler.LogConn,
		PubSubBroker:            pubSub,
		CallerBroker:            caller,
	}
}

func newRedisPool(conf *Redis) *redis.Pool {
	addr := conf.Addr
	p := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// test the connection so that it fails fast if redis is not available
	c := p.Get()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		log.Fatalf("redis PING failed: %v", err)
	}

	return p
}
