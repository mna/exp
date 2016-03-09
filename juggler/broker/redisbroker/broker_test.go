package redisbroker

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrokerCall(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	pool := redistest.NewPool(t, ":"+port)
	broker := &Broker{
		Pool:    pool,
		LogFunc: logIfVerbose,
		CallCap: 2,
	}

	var uuids []uuid.UUID
	connUUID := uuid.NewRandom()
	key := fmt.Sprintf(callKey, "a")
	for i := 0; i < 3; i++ {
		msgUUID := uuid.NewRandom()
		cp := &msg.CallPayload{
			ConnUUID: connUUID,
			MsgUUID:  msgUUID,
			URI:      "a",
		}
		uuids = append(uuids, msgUUID)
		err := broker.Call(cp, time.Second)
		if i < 2 {
			assert.NoError(t, err, "Call %d", i)
		} else {
			assert.Error(t, err, "Call %d", i)
			assert.Contains(t, err.Error(), "list capacity exceeded", "error has expected message")
			t.Logf("%T %[1]v", err)
		}
	}

	// the first 2 msg uuids should be present, in inverted order (LPUSH)
	expectUUIDs(t, pool.Get(), key, uuids[1], uuids[0])

	// call on a different URI works fine
	cp := &msg.CallPayload{
		ConnUUID: connUUID,
		MsgUUID:  uuid.NewRandom(),
		URI:      "b",
	}
	err := broker.Call(cp, time.Second)
	assert.NoError(t, err, "Call on different URI")

	// popping a value should pop uuids[0]
	rc := pool.Get()
	defer rc.Close()
	_, err = rc.Do("RPOP", key)
	require.NoError(t, err, "RPOP")

	expectUUIDs(t, pool.Get(), key, uuids[1])

	// call should now work
	cp = &msg.CallPayload{
		ConnUUID: connUUID,
		MsgUUID:  uuid.NewRandom(),
		URI:      "a",
	}
	uuids = append(uuids, cp.MsgUUID)
	err = broker.Call(cp, time.Second)
	assert.NoError(t, err, "Call after RPOP")

	expectUUIDs(t, pool.Get(), key, uuids[3], uuids[1])
}

func expectUUIDs(t *testing.T, rc redis.Conn, key string, uuids ...uuid.UUID) {
	defer rc.Close()
	vals, err := redis.ByteSlices(rc.Do("LRANGE", key, 0, -1))
	require.NoError(t, err, "LRANGE")

	if assert.Equal(t, len(uuids), len(vals), "number of items") {
		for i, v := range vals {
			var cp msg.CallPayload
			require.NoError(t, json.Unmarshal(v, &cp), "unmarshal into CallPayload")
			assert.Equal(t, uuids[i], cp.MsgUUID, "expected MsgUUID at %d", i)
		}
	}
}

func logIfVerbose(s string, args ...interface{}) {
	if testing.Verbose() {
		log.Printf(s, args...)
	}
}
