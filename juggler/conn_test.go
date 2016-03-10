package juggler

import (
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestConnClose(t *testing.T) {
	srv := &Server{}
	conn := newConn(&websocket.Conn{}, srv)
	conn.psc, conn.resc = fakePubSubConn{}, fakeResultsConn{}

	kill := conn.CloseNotify()
	select {
	case <-kill:
		assert.Fail(t, "close channel should block until call to Close")
	default:
	}

	conn.Close(errors.New("a"))
	select {
	case <-kill:
	default:
		assert.Fail(t, "close channel should be unblocked after call to Close")
	}

	conn.Close(errors.New("b"))
	select {
	case <-kill:
	default:
		assert.Fail(t, "close channel should still be unblocked after subsequent call to Close")
	}

	assert.Equal(t, errors.New("a"), conn.CloseErr, "got expected close error")
}
