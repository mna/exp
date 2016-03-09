package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

func main() {
	pool := newRedisPool(":6379")
	broker := &redisbroker.Broker{
		Pool:      pool,
		Dial:      pool.Dial,
		ResultCap: 100,
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
