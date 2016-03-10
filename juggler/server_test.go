package juggler

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker/redisbroker"
	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/PuerkitoBio/exp/juggler/internal/wstest"
	"github.com/stretchr/testify/assert"
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
