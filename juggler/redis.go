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

type resPayload callPayload

// called by the server to push a call request in the call list key.
func pushRedisCall(c *Conn, m *msg.Call) error {
	const defaultCallTimeout = time.Minute

	pld := &callPayload{
		ConnUUID: c.UUID,
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

	rc := c.srv.CallPool.Get()
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

// called by the callee to push a result in the result list key.
// TODO : move to a callee package or something.
func pushRedisRes() error {
	return nil
}

func pullRedisRes(c *Conn) {
	const minTimeoutSecs = 1

	rc := c.srv.CallPool.Get()
	defer rc.Close()

	for {
		// check for stop signal
		select {
		case <-c.kill:
			return
		default:
		}

		// get the next call result
		toSecs := int(c.srv.ResBrpopTimeout / time.Second)
		if toSecs <= minTimeoutSecs {
			toSecs = minTimeoutSecs
		}
		b, err := redis.Bytes(rc.Do("BRPOP", fmt.Sprintf(resKey, c.UUID), toSecs))
		if err != nil {
			// TODO : do not return
		}

		var m resPayload
		if err := json.Unmarshal(b, &m); err != nil {
			// TODO
		}

		// check if it is still expected (not timed-out)
		cnt, err := rc.Do("DEL", fmt.Sprintf(resTimeoutKey, c.UUID, m.MsgUUID))
		if err != nil {
			// TODO
		}

		if cnt == 1 {
			res := msg.NewRes(m)
			c.Send(res)
		}
	}
}
