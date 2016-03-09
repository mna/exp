package redisbroker

import (
	"encoding/json"
	"sync"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

var _ broker.PubSubConn = (*pubSubConn)(nil)

type pubSubConn struct {
	psc   redis.PubSubConn
	logFn func(string, ...interface{})

	// wmu controls writes (sub/unsub calls) to the connection.
	wmu sync.Mutex

	// once makes sure only the first call to Events starts the goroutine.
	once sync.Once
	evch chan *msg.EvntPayload

	// errmu protects access to err.
	errmu sync.Mutex
	err   error
}

func newPubSubConn(rc redis.Conn, logFn func(string, ...interface{})) *pubSubConn {
	return &pubSubConn{psc: redis.PubSubConn{Conn: rc}, logFn: logFn}
}

// Close closes the connection.
func (c *pubSubConn) Close() error {
	return c.psc.Close()
}

// Subscribe subscribes the redis connection to the channel, which may
// be a pattern.
func (c *pubSubConn) Subscribe(channel string, pattern bool) error {
	return c.subUnsub(channel, pattern, true)
}

// Unsubscribe unsubscribes the redis connection from the channel, which
// may be a pattern.
func (c *pubSubConn) Unsubscribe(channel string, pattern bool) error {
	return c.subUnsub(channel, pattern, false)
}

func (c *pubSubConn) subUnsub(ch string, pat bool, sub bool) error {
	var fn func(...interface{}) error
	switch {
	case pat && sub:
		fn = c.psc.PSubscribe
	case pat && !sub:
		fn = c.psc.PUnsubscribe
	case !pat && sub:
		fn = c.psc.Subscribe
	case !pat && !sub:
		fn = c.psc.Unsubscribe
	}

	c.wmu.Lock()
	err := fn(ch)
	c.wmu.Unlock()
	return err
}

// Events returns the stream of events from channels that the redis
// connection is subscribed to.
func (c *pubSubConn) Events() <-chan *msg.EvntPayload {
	c.once.Do(func() {
		c.evch = make(chan *msg.EvntPayload)

		go func() {
			defer close(c.evch)

			for {
				switch v := c.psc.Receive().(type) {
				case redis.Message:
					ep, err := newEvntPayload(v.Channel, "", v.Data)
					if err != nil {
						logf(c.logFn, "Events: failed to unmarshal event payload: %v", err)
						continue
					}
					c.evch <- ep

				case redis.PMessage:
					ep, err := newEvntPayload(v.Channel, v.Pattern, v.Data)
					if err != nil {
						logf(c.logFn, "Events: failed to unmarshal event payload: %v", err)
						continue
					}
					c.evch <- ep

				case error:
					// possibly because the pub-sub connection was closed, but
					// in any case, the pub-sub is now broken, terminate the
					// loop.
					c.errmu.Lock()
					c.err = v
					c.errmu.Unlock()
					return
				}
			}
		}()
	})

	return c.evch
}

func newEvntPayload(channel, pattern string, pld []byte) (*msg.EvntPayload, error) {
	var pp msg.PubPayload
	if err := json.Unmarshal(pld, &pp); err != nil {
		return nil, err
	}
	ep := &msg.EvntPayload{
		MsgUUID: pp.MsgUUID,
		Channel: channel,
		Pattern: pattern,
		Args:    pp.Args,
	}
	return ep, nil
}

// EventsErr returns the error that caused the events channel to close.
func (c *pubSubConn) EventsErr() error {
	c.errmu.Lock()
	err := c.err
	c.errmu.Unlock()
	return err
}
