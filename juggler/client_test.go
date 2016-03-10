package juggler

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/internal/wstest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	var buf bytes.Buffer
	done := make(chan bool, 1)
	srv := wstest.StartRecordingServer(t, done, &buf)
	defer srv.Close()

	// the only received message should be EXP
	var (
		mu         sync.Mutex
		cnt        int
		expForUUID uuid.UUID
		wg         sync.WaitGroup
	)
	h := MsgHandlerFunc(func(m msg.Msg) {
		defer wg.Done()

		mu.Lock()
		cnt++
		if assert.Equal(t, msg.ExpMsg, m.Type(), "Expects EXP message") {
			expForUUID = m.(*msg.Exp).Payload.For
		}
		mu.Unlock()
	})

	closed := make(chan bool)
	dbgClientClosed = func(c *Client) { closed <- true }
	defer func() { dbgClientClosed = nil }()

	cli, err := Dial(&websocket.Dialer{}, srv.URL, nil, h)
	require.NoError(t, err, "Dial")
	cli.CallTimeout = time.Millisecond
	cli.LogFunc = (&debugLog{t: t}).Printf

	// call
	wg.Add(1)
	type expected struct {
		uid uuid.UUID
		mt  msg.MessageType
	}
	var expectedResults []expected
	callUUID, err := cli.Call("a", "call", 0)
	require.NoError(t, err, "Call")
	expectedResults = append(expectedResults, expected{callUUID, msg.CallMsg})

	uid, err := cli.Pub("b", "pub")
	require.NoError(t, err, "Pub")
	expectedResults = append(expectedResults, expected{uid, msg.PubMsg})

	uid, err = cli.Sub("c", false)
	require.NoError(t, err, "Sub")
	expectedResults = append(expectedResults, expected{uid, msg.SubMsg})

	uid, err = cli.Unsb("d", true)
	require.NoError(t, err, "Unsb")
	expectedResults = append(expectedResults, expected{uid, msg.UnsbMsg})

	// wait for any pending handlers
	wg.Wait()

	cli.Close()
	<-done
	<-closed

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 1, cnt, "Expected calls to Handler")
	assert.Equal(t, callUUID, expForUUID, "Expired message should be for the Call message")

	// read the messages received by the server
	var p json.RawMessage
	r := bytes.NewReader(buf.Bytes())
	dec := json.NewDecoder(r)
	for i, exp := range expectedResults {
		require.NoError(t, dec.Decode(&p), "Decode %d", i)
		m, err := msg.Unmarshal(bytes.NewReader(p))
		require.NoError(t, err, "Unmarshal %d", i)
		assert.Equal(t, exp.uid, m.UUID(), "%d: uuid", i)
		assert.Equal(t, exp.mt, m.Type(), "%d: type", i)
	}

	// no superfluous bytes
	finalErr := dec.Decode(&p)
	if assert.Error(t, finalErr, "Decode after expected results") {
		assert.Equal(t, io.EOF, finalErr, "EOF")
	}
}
