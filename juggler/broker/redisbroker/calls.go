package redisbroker

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

var _ broker.CallsConn = (*callsConn)(nil)

const (
	delAndPTTLScript = `
		local res = redis.call("PTTL", KEYS[1])
		redis.call("DEL", KEYS[1])
		return res
	`
)

type callsConn struct {
	c       redis.Conn
	uris    []string
	timeout time.Duration
	logFn   func(string, ...interface{})

	// once makes sure only the first call to Calls starts the goroutine.
	once sync.Once
	ch   chan *msg.CallPayload

	// errmu protects access to err.
	errmu sync.Mutex
	err   error
}

func newCallsConn(rc redis.Conn, uris []string, to time.Duration, logFn func(string, ...interface{})) *callsConn {
	return &callsConn{c: rc, uris: uris, timeout: to, logFn: logFn}
}

// Close closes the connection.
func (c *callsConn) Close() error {
	return c.c.Close()
}

// CallsErr returns the error that caused the Calls channel to close.
func (c *callsConn) CallsErr() error {
	c.errmu.Lock()
	err := c.err
	c.errmu.Unlock()
	return err
}

// Calls returns a stream of call requests for the URIs specified when
// creating the callsConn.
func (c *callsConn) Calls() <-chan *msg.CallPayload {
	c.once.Do(func() {
		c.ch = make(chan *msg.CallPayload)

		go func() {
			defer close(c.ch)

			// compute all keys and timeout
			keys := make([]string, len(c.uris))
			for i, uri := range c.uris {
				keys[i] = fmt.Sprintf(callKey, uri)
			}
			to := int(c.timeout / time.Second)
			args := redis.Args{}.AddFlat(keys).Add(to)

			for {
				// BRPOP returns array with [0]: key name, [1]: payload.
				v, err := redis.Values(c.c.Do("BRPOP", args...))
				if err != nil {
					if err == redis.ErrNil {
						// no available value
						continue
					}

					// possibly a closed connection, in any case stop
					// the loop.
					c.errmu.Lock()
					c.err = err
					c.errmu.Unlock()
					return
				}

				// unmarshal the payload
				var cp msg.CallPayload
				if err := unmarshalBRPOPValue(&cp, v); err != nil {
					logf(c.logFn, "Calls: BRPOP failed to unmarshal call payload: %v", err)
					continue
				}

				// check if call is expired
				k := fmt.Sprintf(callTimeoutKey, cp.URI, cp.MsgUUID)
				pttl, err := redis.Int(c.c.Do("EVAL", delAndPTTLScript, 1, k))
				if err != nil {
					logf(c.logFn, "Calls: DEL/PTTL failed: %v", err)
					continue
				}
				if pttl <= 0 {
					logf(c.logFn, "Calls: message %v expired, dropping call", cp.MsgUUID)
					continue
				}

				cp.ReadTimestamp = time.Now().UTC()
				cp.TTLAfterRead = time.Duration(pttl) * time.Millisecond
				c.ch <- &cp
			}
		}()
	})

	return c.ch
}

func unmarshalBRPOPValue(dst interface{}, src []interface{}) error {
	var p []byte
	if _, err := redis.Scan(src, nil, p); err != nil {
		return err
	}
	if err := json.Unmarshal(p, dst); err != nil {
		return err
	}
	return nil
}
