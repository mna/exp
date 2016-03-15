// Command juggler-server implements a juggler server that listens for
// connections and serves the requests.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
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
}

// Config defines the configuration options of the server.
type Config struct {
	Redis *Redis `yaml:"redis"`
}

func getConfigFromReader(r io.Reader) (*Config, error) {
	// set default values
	conf := &Config{
		Redis: &Redis{
			Addr:      *redisAddrFlag,
			MaxActive: 100,
			MaxIdle:   10,
		},
	}

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

	// wrap LogMsg and ProcessMsg in a PanicRecover handler
	h := juggler.PanicRecover(
		juggler.Chain(
			juggler.HandlerFunc(juggler.LogMsg),
			juggler.HandlerFunc(juggler.ProcessMsg),
		), true, true)

	pool := newRedisPool(conf.Redis)
	broker := &redisbroker.Broker{
		Pool:      pool,
		Dial:      pool.Dial,
		CallCap:   100,
		ResultCap: 100,
	}
	// test the connection so that it fails fast if redis is not available
	c := pool.Get()
	if _, err := c.Do("PING"); err != nil {
		log.Fatalf("redis PING failed: %v", err)
	}

	if *allowEmptyProtoFlag {
		juggler.Subprotocols = append(juggler.Subprotocols, "")
	}
	upg := &websocket.Upgrader{Subprotocols: juggler.Subprotocols}
	srv := &juggler.Server{
		ReadLimit:               4096,
		WriteLimit:              4096,
		AcquireWriteLockTimeout: 200 * time.Millisecond,
		ConnState:               juggler.LogConn,
		Handler:                 h,
		PubSubBroker:            broker,
		CallerBroker:            broker,
	}
	http.Handle("/ws", juggler.Upgrade(upg, srv))

	log.Printf("listening on port %d", *portFlag)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil); err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}

func newRedisPool(conf *Redis) *redis.Pool {
	addr := conf.Addr
	return &redis.Pool{
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
}
