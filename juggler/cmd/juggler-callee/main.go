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
	"sync"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/callee"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

var (
	redisAddrFlag             = flag.String("redis", ":6379", "Redis `address`.")
	redisPoolMaxActiveFlag    = flag.Int("redis-max-active", 100, "Maximum active redis `connections`.")
	redisPoolMaxIdleFlag      = flag.Int("redis-max-idle", 10, "Maximum idle redis `connections`.")
	redisPoolIdleTimeoutFlag  = flag.Duration("redis-idle-timeout", time.Minute, "Redis idle connection `timeout`.")
	brokerResultCapFlag       = flag.Int("broker-result-cap", 100, "Capacity of the `results` queue.")
	brokerBlockingTimeoutFlag = flag.Duration("broker-blocking-timeout", 0, "Blocking `timeout` when polling for call requests.")
	workersFlag               = flag.Int("workers", 1, "Number of concurrent `workers` processing call requests.")
	helpFlag                  = flag.Bool("help", false, "Show help.")
)

var uris = map[string]callee.Thunk{
	"test.echo":    echoThunk,
	"test.reverse": reverseThunk,
	"test.delay":   delayThunk,
}

func main() {
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}
	if *workersFlag <= 0 {
		*workersFlag = 1
	}

	pool := newRedisPool(*redisAddrFlag)
	c := &callee.Callee{Broker: newBroker(pool)}

	log.Printf("listening for call requests on %s with %d workers", *redisAddrFlag, *workersFlag)

	keys := make([]string, 0, len(uris))
	for k := range uris {
		keys = append(keys, k)
	}

	cc, err := c.Broker.Calls(keys...)
	if err != nil {
		log.Fatalf("Calls failed: %v", err)
	}
	defer cc.Close()

	wg := sync.WaitGroup{}
	wg.Add(*workersFlag)
	for i := 0; i < *workersFlag; i++ {
		go func() {
			defer wg.Done()

			ch := cc.Calls()
			for cp := range ch {
				log.Printf("received request %v %s", cp.MsgUUID, cp.URI)
				if err := c.InvokeAndStoreResult(cp, uris[cp.URI]); err != nil {
					if err != callee.ErrCallExpired {
						log.Printf("InvokeAndStoreResult failed: %v", err)
						continue
					}
					log.Printf("expired request %v %s", cp.MsgUUID, cp.URI)
					continue
				}
				log.Printf("sent result %v %s", cp.MsgUUID, cp.URI)
			}
		}()
	}
	wg.Wait()
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

func newBroker(pool *redis.Pool) broker.CalleeBroker {
	return &redisbroker.Broker{
		Pool:            pool,
		Dial:            pool.Dial,
		BlockingTimeout: *brokerBlockingTimeoutFlag,
		ResultCap:       *brokerResultCapFlag,
	}
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
