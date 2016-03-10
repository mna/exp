package juggler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/internal/wstest"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerServe(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	done := make(chan bool, 1)
	srv := wstest.StartRecordingServer(t, done, ioutil.Discard)
	defer srv.Close()

	dbgl := &debugLog{t: t}
	pool := redistest.NewPool(t, ":"+port)
	broker := &redisbroker.Broker{
		Pool:    pool,
		Dial:    pool.Dial,
		LogFunc: dbgl.Printf,
	}

	conn := wstest.Dial(t, srv.URL)
	defer conn.Close()

	state := make(chan ConnState)
	fn := func(c *Conn, cs ConnState) {
		select {
		case state <- cs:
		case <-time.After(100 * time.Millisecond):
			assert.Fail(t, "could not sent state %d", cs)
		}
	}
	server := &Server{ConnState: fn, CallerBroker: broker, PubSubBroker: broker, LogFunc: dbgl.Printf}

	go server.ServeConn(conn)

	var got ConnState
	select {
	case got = <-state:
		assert.Equal(t, Connected, got, "received connected connection state")
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "no connected state received")
	}

	// closing the underlying websocket connection causes the juggler connection
	// to close too.
	conn.Close()

	select {
	case got = <-state:
		assert.Equal(t, Closing, got, "received closing connection state")
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "no closing state received")
	}
}

func TestUpgrade(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	dbgl := &debugLog{t: t}
	pool := redistest.NewPool(t, ":"+port)
	broker := &redisbroker.Broker{
		Pool:    pool,
		Dial:    pool.Dial,
		LogFunc: dbgl.Printf,
	}

	server := &Server{CallerBroker: broker, PubSubBroker: broker, LogFunc: dbgl.Printf}
	upg := &websocket.Upgrader{Subprotocols: Subprotocols}
	srv := httptest.NewServer(Upgrade(upg, server))
	srv.URL = strings.Replace(srv.URL, "http:", "ws:", 1)
	defer srv.Close()

	h := MsgHandlerFunc(func(m msg.Msg) {})
	closed := make(chan bool)
	dbgClientClosed = func(c *Client) { closed <- true }
	defer func() { dbgClientClosed = nil }()

	// valid subprotocol - no protocol will be set to juggler automatically
	cli, err := Dial(&websocket.Dialer{}, srv.URL, nil, h)
	require.NoError(t, err, "Dial 1")
	cli.Close()
	select {
	case <-closed:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "no close signal received for Dial 1")
	}

	// invalid subprotocol, websocket connection will be closed
	cli, err = Dial(&websocket.Dialer{}, srv.URL, http.Header{"Sec-WebSocket-Protocol": {"test"}}, h)
	require.NoError(t, err, "Dial 2")
	// no need to call Close, Upgrade will refuse the connection
	select {
	case <-closed:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "no close signal received for Dial 2")
	}
	cli.Close()
}
