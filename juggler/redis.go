package juggler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
)

// RedisPool defines the methods required for a Redis pool. The
// pool provides connections via Get and must be closed to release
// its resources.
type RedisPool interface {
	// Get returns a redis connection.
	Get() redis.Conn

	// Close releases the resources used by the pool.
	Close() error
}

const (
	defaultCallTimeout = time.Minute

	// CALL: callee BRPOPs on callKey. On a new payload, it checks if
	// callTimeoutKey is still valid and for how long (PTTL). If it is
	// still valid, it processes the call, otherwise it drops it.
	// callTimeoutKey is deleted.
	callKey        = "juggler:calls:{%s}"            // 1: URI
	callTimeoutKey = "juggler:calls:timeout:{%s}:%s" // 1: URI, 2: mUUID

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

func (s *Server) pushRedisCall(connUUID uuid.UUID, m *msg.Call) error {
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

	c := s.CallPool.Get()
	defer c.Close()

	to := int(m.Payload.Timeout / time.Millisecond)
	if to == 0 {
		to = int(defaultCallTimeout / time.Millisecond)
	}
	if err := c.Send("SET", fmt.Sprintf(callTimeoutKey, m.Payload.URI, m.UUID()), to, "PX", to); err != nil {
		return err
	}
	_, err = c.Do("LPUSH", fmt.Sprintf(callKey, m.Payload.URI), b)
	return err
}
