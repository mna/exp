package juggler

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler/internal/wstest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientClose(t *testing.T) {
	done := make(chan bool, 1)
	srv := wstest.StartRecordingServer(t, done, ioutil.Discard)
	defer srv.Close()

	h := ClientHandlerFunc(func(ctx context.Context, cli *Client, m msg.Msg) {})
	cli, err := Dial(&websocket.Dialer{}, srv.URL, nil, SetHandler(h), SetLogFunc((&debugLog{t: t}).Printf))
	require.NoError(t, err, "Dial")

	_, err = cli.Call("a", "b", 0)
	require.NoError(t, err, "Call")
	require.NoError(t, cli.Close(), "Close")
	if err := cli.Close(); assert.Error(t, err, "Close") {
		assert.Contains(t, err.Error(), "use of closed network connection", "2nd Close")

		if _, err := cli.Call("c", "d", 0); assert.Error(t, err, "Call after Close") {
			assert.Contains(t, err.Error(), "use of closed network connection", "2nd Close")
		}
	}
}

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
	h := ClientHandlerFunc(func(ctx context.Context, cli *Client, m msg.Msg) {
		defer wg.Done()

		mu.Lock()
		cnt++
		if assert.Equal(t, msg.ExpMsg, m.Type(), "Expects EXP message") {
			expForUUID = m.(*msg.Exp).Payload.For
		}
		mu.Unlock()
	})

	cli, err := Dial(&websocket.Dialer{}, srv.URL, nil, SetHandler(h), SetCallTimeout(time.Millisecond), SetLogFunc((&debugLog{t: t}).Printf))
	require.NoError(t, err, "Dial")

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
	<-cli.stop

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
