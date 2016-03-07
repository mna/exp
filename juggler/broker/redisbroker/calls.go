package redisbroker

import (
	"sync"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

type callsConn struct {
	c     redis.Conn
	uris  []string
	logFn func(string, ...interface{})

	// once makes sure only the first call to Calls starts the goroutine.
	once sync.Once
	ch   chan *msg.CallPayload

	// errmu protects access to err.
	errmu sync.Mutex
	err   error
}

func newCallsConn(rc redis.Conn, logFn func(string, ...interface{}), uris ...string) *callsConn {
	return &callsConn{c: rc, uris: uris, logFn: logFn}
}

// Close closes the connection.
func (c *callsConn) Close() error {
	return c.c.Close()
}

/*
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
*/
