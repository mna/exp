// Package redisbroker implements a juggler broker using Redis
// as backend. RPC calls and results are stored in Redis lists
// and queried via the BRPOP command, while pub-sub events
// are handled using Redis' built-in pub-sub support.
//
// Call timeouts are handled by an expiring key associated
// with each call request, and in a similar way for results.
// Keys are named in such a way that the call request list
// and associated expiring keys are in the same hash slot,
// and the same is true for results and their expiring key,
// so that using a redis cluster is supported. The call
// requests are hashed on the call URI, and the results
// are hashed on the calling connection's UUID.
package redisbroker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
)

var (
	// static check that *Broker implements all the broker interfaces
	_ broker.CallerBroker = (*Broker)(nil)
	_ broker.CalleeBroker = (*Broker)(nil)
	_ broker.PubSubBroker = (*Broker)(nil)
)

// Pool defines the methods required for a redis pool that provides
// a method to get a connection and to release the pool's resources.
type Pool interface {
	// Get returns a redis connection.
	Get() redis.Conn

	// Close releases the resources used by the pool.
	Close() error
}

// Broker is a broker that provides the methods to
// interact with Redis using the juggler protocol.
type Broker struct {
	// Pool is the redis pool to use to get short-lived connections.
	Pool Pool

	// Dial is the function to call to get a non-pooled, long-lived
	// redis connection. Typically, it can be set to redis.Pool.Dial.
	Dial func() (redis.Conn, error)

	// BlockingTimeout is the time to wait for a value on calls to
	// BRPOP before trying again. The default of 0 means no timeout.
	BlockingTimeout time.Duration

	// LogFunc is the logging function to use. If nil, log.Printf
	// is used. It can be set to juggler.DiscardLog to disable logging.
	LogFunc func(string, ...interface{})

	// CallCap is the capacity of the CALL queue per URI. If it is
	// exceeded for a given URI, subsequent Broker.Call calls for that
	// URI will fail with an error.
	CallCap int

	// ResultCap is the capacity of the RES queue per connection UUID.
	// If it is exceeded for a given connection, Broker.Result calls
	// for that connection will fail with an error.
	ResultCap int
}

const (
	callOrResScript = `
		redis.call("SET", KEYS[1], ARGV[1], "PX", tonumber(ARGV[1]))
		local res = redis.call("LPUSH", KEYS[2], ARGV[2])
		local limit = tonumber(ARGV[3])
		if res > limit and limit > 0 then
			redis.call("LTRIM", KEYS[2], 1, limit + 1)
			return redis.error_reply("list capacity exceeded")
		end
		return res
	`

	// redis cluster-compliant keys, so that both keys are in the same slot
	callKey        = "juggler:calls:{%s}"            // 1: URI
	callTimeoutKey = "juggler:calls:timeout:{%s}:%s" // 1: URI, 2: mUUID

	// redis cluster-compliant keys, so that both keys are in the same slot
	resKey        = "juggler:results:{%s}"            // 1: cUUID
	resTimeoutKey = "juggler:results:timeout:{%s}:%s" // 1: cUUID, 2: mUUID
)

// Call registers a call request in the broker.
func (b *Broker) Call(cp *msg.CallPayload, timeout time.Duration) error {
	k1 := fmt.Sprintf(callTimeoutKey, cp.URI, cp.MsgUUID)
	k2 := fmt.Sprintf(callKey, cp.URI)
	return registerCallOrRes(b.Pool, cp, timeout, b.CallCap, k1, k2)
}

// Result registers a call result in the broker.
func (b *Broker) Result(rp *msg.ResPayload, timeout time.Duration) error {
	k1 := fmt.Sprintf(resTimeoutKey, rp.ConnUUID, rp.MsgUUID)
	k2 := fmt.Sprintf(resKey, rp.ConnUUID)
	return registerCallOrRes(b.Pool, rp, timeout, b.ResultCap, k1, k2)
}

func registerCallOrRes(pool Pool, pld interface{}, timeout time.Duration, cap int, k1, k2 string) error {
	p, err := json.Marshal(pld)
	if err != nil {
		return err
	}

	rc := pool.Get()
	defer rc.Close()

	to := int(timeout / time.Millisecond)
	if to == 0 {
		to = int(broker.DefaultCallTimeout / time.Millisecond)
	}

	_, err = rc.Do("EVAL",
		callOrResScript,
		2,   // the number of keys
		k1,  // key[1] : the SET key with expiration
		k2,  // key[2] : the LIST key
		to,  // argv[1] : the timeout in milliseconds
		p,   // argv[2] : the call payload
		cap, // argv[3] : the LIST capacity
	)
	return err
}

// Publish publishes an event to a channel.
func (b *Broker) Publish(channel string, pp *msg.PubPayload) error {
	p, err := json.Marshal(pp)
	if err != nil {
		return err
	}

	rc := b.Pool.Get()
	defer rc.Close()

	_, err = rc.Do("PUBLISH", channel, p)
	return err
}

// PubSub returns a pub-sub connection that can be used to subscribe and
// unsubscribe to channels, and to process incoming events.
func (b *Broker) PubSub() (broker.PubSubConn, error) {
	rc, err := b.Dial()
	if err != nil {
		return nil, err
	}
	return newPubSubConn(rc, b.LogFunc), nil
}

// Calls returns a calls connection that can be used to process the call
// requests for the specified URIs.
func (b *Broker) Calls(uris ...string) (broker.CallsConn, error) {
	rc, err := b.Dial()
	if err != nil {
		return nil, err
	}
	return newCallsConn(rc, uris, b.BlockingTimeout, b.LogFunc), nil
}

// Results returns a results connection that can be used to process the call
// results for the specified connection UUID.
func (b *Broker) Results(connUUID uuid.UUID) (broker.ResultsConn, error) {
	rc, err := b.Dial()
	if err != nil {
		return nil, err
	}
	return newResultsConn(rc, connUUID, b.BlockingTimeout, b.LogFunc), nil
}

func logf(fn func(string, ...interface{}), f string, args ...interface{}) {
	if fn != nil {
		fn(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
