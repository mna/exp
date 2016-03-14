// Command juggler-callee implements a testing callee that provides
// simple URI functions.
//
//     - test.echo (string) : returns the received string
//     - test.reverse (string) : reverses each rune in the received string
//     - test.delay (string) : sleeps for the duration received as string, converted to number (in ms)
//
package main

import (
	"encoding/json"
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

var (
	redisAddrFlag             = flag.String("redis-addr", ":6379", "redis address to connect to")
	redisPoolMaxActiveFlag    = flag.Int("redis-pool-max-active", 100, "redis pool max active connections")
	redisPoolMaxIdleFlag      = flag.Int("redis-pool-max-idle", 10, "redis pool max idle connections")
	redisPoolIdleTimeoutFlag  = flag.Duration("redis-pool-idle-timeout", time.Minute, "redis pool idle connection timeout")
	brokerResultCapFlag       = flag.Int("broker-result-cap", 100, "broker result queue capacity")
	brokerBlockingTimeoutFlag = flag.Duration("broker-blocking-timeout", 0, "broker blocking timeout")
)

func main() {
	flag.Parse()

	pool := newRedisPool(*redisAddrFlag)
	broker := &redisbroker.Broker{
		Pool:            pool,
		Dial:            pool.Dial,
		BlockingTimeout: *brokerBlockingTimeoutFlag,
		ResultCap:       *brokerResultCapFlag,
	}

	c := &callee.Callee{Broker: broker}

	log.Printf("listening for call requests")
	if err := c.Listen(map[string]callee.Thunk{
		"test.echo":    logWrapThunk(echoThunk),
		"test.reverse": logWrapThunk(reverseThunk),
		"test.delay":   logWrapThunk(delayThunk),
	}); err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
}

func logWrapThunk(t callee.Thunk) callee.Thunk {
	return func(cp *msg.CallPayload) (interface{}, error) {
		log.Printf("received call for %s from %v", cp.URI, cp.MsgUUID)
		v, err := t(cp)
		log.Printf("sending result for %s from %v", cp.URI, cp.MsgUUID)
		return v, err
	}
}

func delayThunk(cp *msg.CallPayload) (interface{}, error) {
	var s string
	if err := json.Unmarshal(cp.Args, &s); err != nil {
		return nil, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return delay(i), nil
}

func delay(i int) int {
	time.Sleep(time.Duration(i) * time.Millisecond)
	return i
}

func reverseThunk(cp *msg.CallPayload) (interface{}, error) {
	var s string
	if err := json.Unmarshal(cp.Args, &s); err != nil {
		return nil, err
	}
	return reverse(s), nil
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func echoThunk(cp *msg.CallPayload) (interface{}, error) {
	var s string
	if err := json.Unmarshal(cp.Args, &s); err != nil {
		return nil, err
	}
	return echo(s), nil
}

func echo(s string) string {
	return s
}

func newRedisPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     *redisPoolMaxIdleFlag,
		MaxActive:   *redisPoolMaxActiveFlag,
		IdleTimeout: *redisPoolIdleTimeoutFlag,
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
