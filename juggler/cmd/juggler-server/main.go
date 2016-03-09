package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

var (
	// TODO : work out a config file for all server and broker options
	portFlag      = flag.Int("port", 9000, "port to listen on")
	readLimitFlag = flag.Int("read-limit", 4096, "read message size limit")
	readTOFlag    = flag.Duration("read-timeout", 10*time.Second, "read deadline duration")
	writeTOFlag   = flag.Duration("write-timeout", 10*time.Second, "write deadline duration")
)

func main() {
	flag.Parse()

	// wrap LogMsg and ProcessMsg in a PanicRecover handler
	h := juggler.PanicRecover(
		juggler.Chain(
			juggler.HandlerFunc(juggler.LogMsg),
			juggler.HandlerFunc(juggler.ProcessMsg),
		), true, true)

	pool := newRedisPool(":6379")
	broker := &redisbroker.Broker{
		Pool:      pool,
		Dial:      pool.Dial,
		CallCap:   100,
		ResultCap: 100,
	}

	upg := &websocket.Upgrader{Subprotocols: juggler.Subprotocols}
	srv := &juggler.Server{
		ReadLimit:               int64(*readLimitFlag),
		ReadTimeout:             *readTOFlag,
		WriteLimit:              4096,
		WriteTimeout:            *writeTOFlag,
		AcquireWriteLockTimeout: 200 * time.Millisecond,
		ConnState:               juggler.LogConn,
		ReadHandler:             h,
		WriteHandler:            h,
		PubSubBroker:            broker,
		CallerBroker:            broker,
	}
	http.Handle("/ws", juggler.Upgrade(upg, srv))

	log.Printf("juggler: listening on port %d", *portFlag)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil); err != nil {
		log.Fatalf("juggler: ListenAndServe failed: %v", err)
	}
}

func newRedisPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10000,
		IdleTimeout: 2 * time.Minute,
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
