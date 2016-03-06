package redconn

import (
	"encoding/json"
	"fmt"
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
}

const (
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

type callPayload struct {
	ConnUUID uuid.UUID       `json:"conn_uuid"`
	MsgUUID  uuid.UUID       `json:"msg_uuid"`
	Args     json.RawMessage `json:"args,omitempty"`
}

func (c *Connector) Call(connUUID uuid.UUID, m *msg.Call) error {
	pld := &callPayload{
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
	return err
}

func (c *Connector) ProcessCalls() {

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
	rc := c.Pool.Get()
	defer rc.Close()

	cmd := subUnsubCmds[struct{ pat, sub bool }{pat, sub}]
	_, err := rc.Do(cmd, ch)
	return err
}

func (c *Connector) ProcessEvents() {

}
