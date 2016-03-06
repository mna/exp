package redconn

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
)

// Pool defines the methods required for a redis pool that provides
// a method to get a connection and to release the pool's resources.
type Pool interface {
	Get() redis.Conn
	Close() error
}

// Connector is a redis connector that provides the methods to
// interact with Redis using the juggler protocol.
type Connector struct {
	Pool            Pool
	BlockingTimeout time.Duration
	LogFunc         func(string, ...interface{})
}

const (
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

// TODO: not redis-specific, should go elsewhere...
type CallPayload struct {
	ConnUUID      uuid.UUID       `json:"conn_uuid"`
	MsgUUID       uuid.UUID       `json:"msg_uuid"`
	Args          json.RawMessage `json:"args,omitempty"`
	TTLAfterRead  time.Duration   `json:"-"`
	ReadTimestamp time.Time       `json:"-"`
}

func (c *Connector) Call(connUUID uuid.UUID, m *msg.Call) error {
	pld := &CallPayload{
		ConnUUID: connUUID,
		MsgUUID:  m.UUID(),
		Args:     m.Payload.Args,
	}
	b, err := json.Marshal(pld)
	if err != nil {
		return err
	}

	// a call generates two redis key values:
	// - SET that expires after timeout
	// - LPUSH that adds the call payload to the list of calls under URI
	//
	// A callee will read with BRPOP on the list, and will check the
	// expiring key to see if it still exists. If it doesn't, the call is
	// dropped, unprocessed, as the client is not waiting for the response
	// anymore.
	//
	// If it is still there, the callee gets its PTTL and deletes it, and
	// it processes the call and stores the response payload under a new
	// key with an expiration of PTTL.

	rc := c.Pool.Get()
	defer rc.Close()

	to := int(m.Payload.Timeout / time.Millisecond)
	if to == 0 {
		to = int(defaultCallTimeout / time.Millisecond)
	}
	if err := rc.Send("SET", fmt.Sprintf(callTimeoutKey, m.Payload.URI, m.UUID()), to, "PX", to); err != nil {
		return err
	}
	_, err = rc.Do("LPUSH", fmt.Sprintf(callKey, m.Payload.URI), b)

	// TODO : support capping the list with LTRIM

	return err
}

var prng = rand.New(rand.NewSource(time.Now().UnixNano()))

func expJitterDelay(att int, base, max time.Duration) time.Duration {
	exp := math.Pow(2, float64(att))
	top := float64(base) * exp
	return time.Duration(
		prng.Int63n(int64(math.Min(float64(max), top))),
	)
}

// ProcessCalls returns a channel that returns a stream of call requests
// for the specified URI. When the stop channel signals a stop, the
// returned channel is closed and the goroutine that listens for call
// requests is properly terminated.
func (c *Connector) ProcessCalls(uri string, stop <-chan struct{}) <-chan *CallPayload {
	ch := make(chan *CallPayload)
	go func() {
		defer close(ch)

		// compute the key and blocking timeout
		k := fmt.Sprintf(callKey, uri)
		to := int(c.BlockingTimeout / time.Second)
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
				rc = c.Pool.Get()
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

				var b []byte
				_, err := redis.Scan(vals, nil, b)
				if err != nil {
					logf(c, "ProcessCalls: BRPOP failed to scan redis value: %v", err)
					continue
				}

				var cp CallPayload
				if err := json.Unmarshal(b, &cp); err != nil {
					logf(c, "ProcessCalls: BRPOP failed to unmarshal call payload: %v", err)
					continue
				}

				toKey := fmt.Sprintf(callTimeoutKey, uri, cp.MsgUUID)
				if err := rc.Send("PTTL", toKey); err != nil {
					logf(c, "ProcessCalls: PTTL send failed: %v", err)
					continue
				}
				res, err := redis.Values(rc.Do("DEL", toKey))
				if err != nil {
					logf(c, "ProcessCalls: PTTL/DEL failed: %v", err)
					continue
				}
				var pttl int
				if _, err := redis.Scan(res, &pttl); err != nil {
					logf(c, "ProcessCalls: PTTL/DEL failed to scan redis value: %v", err)
					continue
				}
				if pttl <= 0 {
					logf(c, "ProcessCalls: message %v expired, dropping call", cp.MsgUUID)
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

func (c *Connector) ProcessResults() {

}

func (c *Connector) Publish(m *msg.Pub) error {
	rc := c.Pool.Get()
	defer rc.Close()

	_, err := rc.Do("PUBLISH", m.Payload.Channel, m.Payload.Args)
	return err
}

func (c *Connector) Subscribe(m *msg.Sub) error {
	return c.subUnsub(m.Payload.Channel, m.Payload.Pattern, true)
}

func (c *Connector) Unsubscribe(m *msg.Unsb) error {
	return c.subUnsub(m.Payload.Channel, m.Payload.Pattern, false)
}

var subUnsubCmds = map[struct{ pat, sub bool }]string{
	{true, true}:   "PSUBSCRIBE",
	{true, false}:  "PUNSUBSCRIBE",
	{false, true}:  "SUBSCRIBE",
	{false, false}: "UNSUBSCRIBE",
}

func (c *Connector) subUnsub(ch string, pat bool, sub bool) error {
	// TODO : no, must be on the same connection always...
	rc := c.Pool.Get()
	defer rc.Close()

	cmd := subUnsubCmds[struct{ pat, sub bool }{pat, sub}]
	_, err := rc.Do(cmd, ch)
	return err
}

func (c *Connector) ProcessEvents() {
	// TODO : must be on the same connection as the sub
}

func logf(c *Connector, f string, args ...interface{}) {
	if c.LogFunc != nil {
		c.LogFunc(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
