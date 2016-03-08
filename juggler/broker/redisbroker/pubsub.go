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
	c     redis.Conn
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
	return &pubSubConn{c: rc, logFn: logFn}
}

// Close closes the connection.
func (c *pubSubConn) Close() error {
	return c.c.Close()
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

var subUnsubCmds = map[struct{ pat, sub bool }]string{
	{true, true}:   "PSUBSCRIBE",
	{true, false}:  "PUNSUBSCRIBE",
	{false, true}:  "SUBSCRIBE",
	{false, false}: "UNSUBSCRIBE",
}

func (c *pubSubConn) subUnsub(ch string, pat bool, sub bool) error {
	cmd := subUnsubCmds[struct{ pat, sub bool }{pat, sub}]

	c.wmu.Lock()
	_, err := c.c.Do(cmd, ch)
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

			psc := redis.PubSubConn{Conn: c.c}
			for {
				switch v := psc.Receive().(type) {
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
