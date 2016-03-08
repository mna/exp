package redisbroker

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
)

var _ broker.ResultsConn = (*resultsConn)(nil)

type resultsConn struct {
	c        redis.Conn
	connUUID uuid.UUID
	timeout  time.Duration
	logFn    func(string, ...interface{})

	// once makes sure only the first call to Results starts the goroutine.
	once sync.Once
	ch   chan *msg.ResPayload

	// errmu protects access to err.
	errmu sync.Mutex
	err   error
}

func newResultsConn(rc redis.Conn, connUUID uuid.UUID, to time.Duration, logFn func(string, ...interface{})) *resultsConn {
	return &resultsConn{c: rc, connUUID: connUUID, timeout: to, logFn: logFn}
}

// Close closes the connection.
func (c *resultsConn) Close() error {
	return c.c.Close()
}

// ResultsErr returns the error that caused the Results channel to close.
func (c *resultsConn) ResultsErr() error {
	c.errmu.Lock()
	err := c.err
	c.errmu.Unlock()
	return err
}

// Results returns a stream of call results for the connUUID specified when
// creating the resultsConn.
func (c *resultsConn) Results() <-chan *msg.ResPayload {
	c.once.Do(func() {
		c.ch = make(chan *msg.ResPayload)

		go func() {
			defer close(c.ch)

			// compute key and timeout
			key := fmt.Sprintf(resKey, c.connUUID)
			to := int(c.timeout / time.Second)
			for {
				// BRPOP returns array with [0]: key name, [1]: payload.
				v, err := redis.Values(c.c.Do("BRPOP", key, to))
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
				var p []byte
				if _, err = redis.Scan(v, nil, p); err != nil {
					logf(c.logFn, "Results: BRPOP failed to scan redis value: %v", err)
					continue
				}
				var rp msg.ResPayload
				if err := json.Unmarshal(p, &rp); err != nil {
					logf(c.logFn, "Results: BRPOP failed to unmarshal result payload: %v", err)
					continue
				}

				// check if call is expired
				k := fmt.Sprintf(resTimeoutKey, rp.ConnUUID, rp.MsgUUID)
				pttl, err := redis.Int(c.c.Do("EVAL", delAndPTTLScript, 1, k))
				if err != nil {
					logf(c.logFn, "Results: DEL/PTTL failed: %v", err)
					continue
				}
				if pttl <= 0 {
					logf(c.logFn, "Results: message %v expired, dropping call", rp.MsgUUID)
					continue
				}

				c.ch <- &rp
			}
		}()
	})

	return c.ch
}
