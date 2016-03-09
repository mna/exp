package redisbroker

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cap = 2

func testBrokerCallOrRes(t *testing.T, keyFmt string, run func(*Broker, uuid.UUID) (uuid.UUID, error)) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	pool := redistest.NewPool(t, ":"+port)
	broker := &Broker{
		Pool:      pool,
		LogFunc:   logIfVerbose,
		CallCap:   cap,
		ResultCap: cap,
	}

	var uuids []uuid.UUID
	// run all on same key
	keyUUID := uuid.NewRandom()
	for i := 0; i <= cap; i++ {
		uid, err := run(broker, keyUUID)
		uuids = append(uuids, uid)
		if i < cap {
			assert.NoError(t, err, "Call %d", i)
		} else {
			assert.Error(t, err, "Call %d", i)
			assert.Contains(t, err.Error(), "list capacity exceeded", "error has expected message")
		}
	}

	// the first 2 msg uuids should be present, in inverted order (LPUSH)
	key := fmt.Sprintf(keyFmt, keyUUID)
	expectUUIDs(t, pool.Get(), key, uuids[1], uuids[0])

	// call on a different URI works fine
	diffKeyUUID := uuid.NewRandom()
	_, err := run(broker, diffKeyUUID)
	assert.NoError(t, err, "Call on different key")

	// popping a value should pop uuids[0]
	rc := pool.Get()
	defer rc.Close()
	_, err = rc.Do("RPOP", key)
	require.NoError(t, err, "RPOP")

	expectUUIDs(t, pool.Get(), key, uuids[1])

	// call should now work on original key
	uid, err := run(broker, keyUUID)
	uuids = append(uuids, uid)
	assert.NoError(t, err, "Call after RPOP")

	expectUUIDs(t, pool.Get(), key, uuids[3], uuids[1])
}

func TestBrokerCall(t *testing.T) {
	connUUID := uuid.NewRandom()
	testBrokerCallOrRes(t, callKey, func(b *Broker, keyParm uuid.UUID) (uuid.UUID, error) {
		cp := &msg.CallPayload{
			ConnUUID: connUUID,
			MsgUUID:  uuid.NewRandom(),
			URI:      keyParm.String(),
		}
		err := b.Call(cp, time.Second)
		return cp.MsgUUID, err
	})
}

func TestBrokerResult(t *testing.T) {
	testBrokerCallOrRes(t, resKey, func(b *Broker, keyParm uuid.UUID) (uuid.UUID, error) {
		rp := &msg.ResPayload{
			ConnUUID: keyParm,
			MsgUUID:  uuid.NewRandom(),
			URI:      "z",
		}
		err := b.Result(rp, time.Second)
		return rp.MsgUUID, err
	})
}

func TestPublish(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	pool := redistest.NewPool(t, ":"+port)
	brk := broker.PubSubBroker(&Broker{
		Pool:    pool,
		Dial:    pool.Dial,
		LogFunc: logIfVerbose,
	})

	psc, err := brk.PubSub()
	require.NoError(t, err, "get PubSubConn")

	// subscribe to channel "a"
	require.NoError(t, psc.Subscribe("a", false), "Subscribe")

	// listen to events on "a"
	var cnt int
	expPlds := []string{`"abc"`, `{"v":3}`}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for ev := range psc.Events() {
			var want string

			if cnt < len(expPlds) {
				want = expPlds[cnt]
			}
			assert.Equal(t, "a", ev.Channel, "event is from the subscribed channel")
			assert.Equal(t, want, string(ev.Args), "event payload")
			cnt++
		}
	}()

	cases := []struct {
		v  interface{}
		ch string
	}{
		{"abc", "a"},
		{"def", "b"},
		{map[string]interface{}{"v": 3}, "a"},
		{5, "c"},
	}
	for i, c := range cases {
		b, err := json.Marshal(c.v)
		require.NoError(t, err, "marshal case %d", i)
		pp := &msg.PubPayload{MsgUUID: uuid.NewRandom(), Args: b}
		require.NoError(t, brk.Publish(c.ch, pp), "Publish event %d", i)
	}

	require.NoError(t, psc.Close(), "close subscribed connection")
	wg.Wait()
	assert.Equal(t, 2, cnt, "number of events received")
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
