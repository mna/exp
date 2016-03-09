package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
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
		"test.echo":    echoThunk,
		"test.reverse": reverseThunk,
	}); err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
}

func reverseThunk(raw json.RawMessage) (interface{}, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
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

func echoThunk(raw json.RawMessage) (interface{}, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
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
