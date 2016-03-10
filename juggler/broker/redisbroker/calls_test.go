package redisbroker

import (
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalls(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	pool := redistest.NewPool(t, ":"+port)
	brk := &Broker{
		Pool:            pool,
		Dial:            pool.Dial,
		BlockingTimeout: time.Second,
		LogFunc:         logIfVerbose,
	}

	// list calls on URI "a"
	cc, err := brk.Calls("a")
	require.NoError(t, err, "get Calls connection")

	// keep track of received calls
	wg := sync.WaitGroup{}
	wg.Add(1)
	var uuids []uuid.UUID
	go func() {
		defer wg.Done()
		for cp := range cc.Calls() {
			uuids = append(uuids, cp.MsgUUID)
		}
	}()

	cases := []struct {
		cp      *msg.CallPayload
		timeout time.Duration
		exp     bool
	}{
		{&msg.CallPayload{ConnUUID: uuid.NewRandom(), MsgUUID: uuid.NewRandom(), URI: "a"}, time.Second, true},
		{&msg.CallPayload{ConnUUID: uuid.NewRandom(), MsgUUID: uuid.NewRandom(), URI: "b"}, time.Second, false},
		{&msg.CallPayload{ConnUUID: uuid.NewRandom(), MsgUUID: uuid.NewRandom(), URI: "a"}, time.Minute, true},
	}
	var expected []uuid.UUID
	for i, c := range cases {
		if c.exp {
			expected = append(expected, c.cp.MsgUUID)
		}
		require.NoError(t, brk.Call(c.cp, c.timeout), "Call %d", i)
	}

	time.Sleep(10 * time.Millisecond) // ensure time to pop the last message :(
	require.NoError(t, cc.Close(), "close calls connection")
	wg.Wait()
	assert.Error(t, cc.CallsErr(), "CallsErr returns the error")
	assert.Equal(t, expected, uuids, "got expected UUIDs")
}
