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

func TestResults(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	pool := redistest.NewPool(t, ":"+port)
	brk := &Broker{
		Pool:            pool,
		Dial:            pool.Dial,
		BlockingTimeout: time.Second,
		LogFunc:         logIfVerbose,
	}

	// list results on this conn UUID
	connUUID := uuid.NewRandom()
	rc, err := brk.Results(connUUID)
	require.NoError(t, err, "get Results connection")

	// keep track of received results
	wg := sync.WaitGroup{}
	wg.Add(1)
	var uuids []uuid.UUID
	go func() {
		defer wg.Done()
		for rp := range rc.Results() {
			uuids = append(uuids, rp.MsgUUID)
		}
	}()

	// wait 1s to test the ErrNil case
	time.Sleep(1100 * time.Millisecond)

	cases := []struct {
		rp      *msg.ResPayload
		timeout time.Duration
		exp     bool
	}{
		{&msg.ResPayload{ConnUUID: connUUID, MsgUUID: uuid.NewRandom(), URI: "a"}, time.Second, true},
		{&msg.ResPayload{ConnUUID: uuid.NewRandom(), MsgUUID: uuid.NewRandom(), URI: "b"}, time.Second, false},
		{&msg.ResPayload{ConnUUID: connUUID, MsgUUID: uuid.NewRandom(), URI: "c"}, 0, true},
	}
	var expected []uuid.UUID
	for i, c := range cases {
		if c.exp {
			expected = append(expected, c.rp.MsgUUID)
		}
		require.NoError(t, brk.Result(c.rp, c.timeout), "Result %d", i)
	}

	time.Sleep(10 * time.Millisecond) // ensure time to pop the last message :(
	require.NoError(t, rc.Close(), "close results connection")
	wg.Wait()
	if assert.Error(t, rc.ResultsErr(), "ResultsErr returns the error") {
		assert.Contains(t, rc.ResultsErr().Error(), "use of closed network connection", "ResultsErr is the expected error")
	}
	assert.Equal(t, expected, uuids, "got expected UUIDs")
}
