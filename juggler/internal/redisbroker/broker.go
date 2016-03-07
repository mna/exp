package redisbroker

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/exp/juggler/internal/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
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
	// Pool is the redis pool to use to get connections.
	Pool Pool

	// BlockingTimeout is the time to wait for a value on calls to
	// BRPOP.
	BlockingTimeout time.Duration

	// LogFunc is the logging function to use. If nil, log.Printf
	// is used. It can be set to juggler.DiscardLog to disable logging.
	LogFunc func(string, ...interface{})
}

const (
	// if no Broker.BlockingTimeout is provided.
	defaultBlockingTimeout = 5 * time.Second

	// CALL: callee BRPOPs on callKey. On a new payload, it checks if
	// callTimeoutKey is still valid and for how long (PTTL). If it is
	// still valid, it processes the call, otherwise it drops it.
	// callTimeoutKey is deleted.
	callKey            = "juggler:calls:{%s}"            // 1: URI
	callTimeoutKey     = "juggler:calls:timeout:{%s}:%s" // 1: URI, 2: mUUID
	defaultCallTimeout = time.Minute

	// RES: callee stores the result of the call in resKey (LPUSH) and
	// sets resTimeoutKey with an expiration of callTimeoutKey PTTL minus
	// the time of the call invocation.
	//
	// Caller BRPOPs on resKey. On a new payload, it checks if resTimeoutKey
	// is still valid. If it is, it sends the result on the connection,
	// otherwise it drops it. resTimeoutKey is deleted.
	resKey        = "juggler:results:{%s}"            // 1: cUUID
	resTimeoutKey = "juggler:results:timeout:{%s}:%s" // 1: cUUID, 2: mUUID
)

// Call registers a call request in the broker.
func (b *Broker) Call(cp *msg.CallPayload, timeout time.Duration) error {
	p, err := json.Marshal(cp)
	if err != nil {
		return err
	}

	rc := b.Pool.Get()
	defer rc.Close()

	to := int(timeout / time.Millisecond)
	if to == 0 {
		to = int(defaultCallTimeout / time.Millisecond)
	}

	// TODO : use script instead
	if err := rc.Send("SET", fmt.Sprintf(callTimeoutKey, cp.URI, cp.MsgUUID), to, "PX", to); err != nil {
		return err
	}
	_, err = rc.Do("LPUSH", fmt.Sprintf(callKey, cp.URI), p)

	// TODO : support capping the list with LTRIM

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
	rc := b.Pool.Get()
	return newPubSubConn(rc, b.LogFunc), nil
}

// Result registers a call result in the broker.
func (b *Broker) Result(rp *msg.ResPayload, timeout time.Duration) error {
	// TODO : implement...
	return nil
}

var prng = rand.New(rand.NewSource(time.Now().UnixNano()))

func expJitterDelay(att int, base, max time.Duration) time.Duration {
	exp := math.Pow(2, float64(att))
	top := float64(base) * exp
	return time.Duration(
		prng.Int63n(int64(math.Min(float64(max), top))),
	)
}

// Calls returns a channel that returns a stream of call requests
// for the specified URI. When the stop channel signals a stop, the
// returned channel is closed and the goroutine that listens for call
// requests is properly terminated.
func (b *Broker) Calls(uri string, stop <-chan struct{}) <-chan *msg.CallPayload {
	ch := make(chan *msg.CallPayload)
	go func() {
		defer close(ch)

		// compute the key and blocking timeout
		k := fmt.Sprintf(callKey, uri)
		to := int(b.BlockingTimeout / time.Second)
		if to == 0 {
			to = int(defaultBlockingTimeout / time.Second)
		}

		var rc redis.Conn
		defer func() {
			if rc != nil {
				rc.Close()
			}
		}()

		var attempt int
		for {
			// check for the stop signal
			select {
			case <-stop:
				return
			default:
			}

			// grab a redis connection if we don't have any valid one.
			if rc == nil {
				rc = b.Pool.Get()
			}

			// block checking for a call request to process.
			vals, err := redis.Values(rc.Do("BRPOP", k, to))
			switch err {
			case redis.ErrNil:
				// no value available
				attempt = 0 // successful redis call
				continue

			case nil:
				// got a call payload, process it
				attempt = 0 // successful redis call

				var p []byte
				_, err := redis.Scan(vals, nil, p)
				if err != nil {
					logf(b.LogFunc, "ProcessCalls: BRPOP failed to scan redis value: %v", err)
					continue
				}

				var cp msg.CallPayload
				if err := json.Unmarshal(p, &cp); err != nil {
					logf(b.LogFunc, "ProcessCalls: BRPOP failed to unmarshal call payload: %v", err)
					continue
				}

				toKey := fmt.Sprintf(callTimeoutKey, uri, cp.MsgUUID)
				if err := rc.Send("PTTL", toKey); err != nil {
					logf(b.LogFunc, "ProcessCalls: PTTL send failed: %v", err)
					continue
				}
				res, err := redis.Values(rc.Do("DEL", toKey))
				if err != nil {
					logf(b.LogFunc, "ProcessCalls: PTTL/DEL failed: %v", err)
					continue
				}
				var pttl int
				if _, err := redis.Scan(res, &pttl); err != nil {
					logf(b.LogFunc, "ProcessCalls: PTTL/DEL failed to scan redis value: %v", err)
					continue
				}
				if pttl <= 0 {
					logf(b.LogFunc, "ProcessCalls: message %v expired, dropping call", cp.MsgUUID)
					continue
				}

				cp.ReadTimestamp = time.Now().UTC()
				cp.TTLAfterRead = time.Duration(pttl) * time.Millisecond
				ch <- &cp

			default:
				// error, try again with a different redis connection, in
				// case that node went down.
				rc.Close()
				rc = nil

				delay := expJitterDelay(attempt, time.Second, time.Minute)
				select {
				case <-stop:
					return
				case <-time.After(delay):
					// go on
					attempt++
				}
			}
		}
	}()

	return ch
}

func (b *Broker) Results() {

}

func logf(fn func(string, ...interface{}), f string, args ...interface{}) {
	if fn != nil {
		fn(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
