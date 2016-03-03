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
	callTimeoutKey     = "juggler:calls:{%s}"
)

type callPayload struct {
	UUID uuid.UUID       `json:"uuid"`
	Args json.RawMessage `json:"args,omitempty"`
}

func (s *Server) redisCall(m *msg.Call) error {
	c := s.CallPool.Get()
	defer c.Close()

	pld := &callPayload{UUID: m.UUID(), Args: m.Payload.Args}
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

	// TODO : use {} to ensure both keys are on the same node/slot when
	// using redis-cluster. (e.g. timeout key is juggler:call:{uri}:uuid).

	to := int(m.Payload.Timeout / time.Millisecond)
	if to == 0 {
		to = int(defaultCallTimeout / time.Millisecond)
	}
	if err := c.Send("SET", fmt.Sprintf(callTimeoutKey, m.UUID()), to, "PX", to); err != nil {
		return err
	}
	_, err = c.Do("LPUSH", m.Payload.URI, b)
	return err
}
