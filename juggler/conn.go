package juggler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

// ErrLockWriterTimeout is returned when a call to Write fails
// because the write lock of the connection cannot be acquired before
// the timeout.
var ErrLockWriterTimeout = errors.New("juggler: timed out waiting for writer lock")

// ConnState represents the possible states of a connection.
type ConnState int

// List of possible connection states.
const (
	Unknown ConnState = iota
	Connected
	Closing
)

// Conn is a juggler connection. Each connection is identified by
// a UUID and has an underlying websocket connection. It is safe to
// call methods on a Conn concurrently, but the fields should be
// treated as read-only.
type Conn struct {
	// UUID is the unique identifier of the connection.
	UUID uuid.UUID

	// WSConn is the underlying websocket connection.
	WSConn *websocket.Conn // TODO : hide/show only as needed

	// CloseErr is the error, if any, that caused the connection
	// to close.
	CloseErr error

	// TODO : some connection state (authenticated, etc.)?
	wmu       chan struct{} // write lock
	srv       *Server
	kill      chan struct{}
	closeOnce sync.Once
}

func newConn(c *websocket.Conn, srv *Server) *Conn {
	// wmu is the write lock, used as a semaphore of 1, so start with
	// an available slot (initialize with a sent value).
	wmu := make(chan struct{}, 1)
	wmu <- struct{}{}

	return &Conn{
		UUID:   uuid.NewRandom(),
		WSConn: c,
		wmu:    wmu,
		srv:    srv,
		kill:   make(chan struct{}),
	}
}

// CloseNotify returns a signal channel that is closed when the
// Conn is closed.
func (c *Conn) CloseNotify() <-chan struct{} {
	return c.kill
}

// Close closes the connection, setting err as CloseErr to identify
// the reason of the close. It does not send a websocket close message.
// As with all Conn methods, it is safe to call concurrently, but
// only the first call will set the CloseErr field to err.
func (c *Conn) Close(err error) {
	c.closeOnce.Do(func() {
		c.CloseErr = err
		close(c.kill)
	})
}

// writer that acquires the connection's write lock prior to writing.
type exclusiveWriter struct {
	w       io.WriteCloser
	c       *Conn
	timeout time.Duration
	init    bool
}

func (w *exclusiveWriter) Write(p []byte) (int, error) {
	if !w.init {
		var wait <-chan time.Time
		if w.timeout > 0 {
			wait = time.After(w.timeout)
		}

		// try to acquire the write lock before the timeout
		select {
		case <-wait:
			return 0, ErrLockWriterTimeout

		case <-w.c.wmu:
			// lock acquired, get next writer from the websocket connection
			w.init = true
			wc, err := w.c.WSConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return 0, err
			}
			w.w = wc
			if w.c.srv.WriteTimeout > 0 {
				w.c.WSConn.SetWriteDeadline(time.Now().Add(w.c.srv.WriteTimeout))
			}
		}
	}

	return w.w.Write(p)
}

func (w *exclusiveWriter) Close() error {
	if !w.init {
		// no write, Close is a no-op
		return nil
	}

	// if w.init is true, then NextWriter was called and that writer
	// must be properly closed.
	err := w.w.Close()
	w.c.WSConn.SetWriteDeadline(time.Time{})

	// release the write lock
	w.c.wmu <- struct{}{}
	return err
}

func (c *Conn) Writer(timeout time.Duration) io.WriteCloser {
	return &exclusiveWriter{
		c:       c,
		timeout: timeout,
	}
}

func (c *Conn) Send(msg Msg) {
	if c.srv.WriteHandler != nil {
		c.srv.WriteHandler.Handle(c, msg)
	} else {
		ProcessMsg(c, msg)
	}
}

func (c *Conn) receive() {
	for {
		c.WSConn.SetReadDeadline(time.Time{})

		mt, r, err := c.WSConn.NextReader()
		if err != nil {
			c.Close(err)
			return
		}
		if mt != websocket.TextMessage {
			c.Close(fmt.Errorf("invalid websocket message type: %d", mt))
			return
		}
		if c.srv.ReadTimeout > 0 {
			c.WSConn.SetReadDeadline(time.Now().Add(c.srv.ReadTimeout))
		}

		msg, err := unmarshalMessage(r)
		if err != nil {
			c.Close(err)
			return
		}

		if c.srv.ReadHandler != nil {
			c.srv.ReadHandler.Handle(c, msg)
		} else {
			ProcessMsg(c, msg)
		}
	}
}

func unmarshalMessage(r io.Reader) (Msg, error) {
	var pm partialMsg
	if err := json.NewDecoder(r).Decode(&pm); err != nil {
		return nil, fmt.Errorf("invalid JSON message: %v", err)
	}

	genericUnmarshal := func(v interface{}, metaDst *meta) error {
		if err := json.Unmarshal(pm.Payload, v); err != nil {
			return fmt.Errorf("invalid %s message: %v", pm.Meta.T, err)
		}
		*metaDst = pm.Meta
		return nil
	}

	var msg Msg
	switch pm.Meta.T {
	case AuthMsg:
		var auth Auth
		if err := genericUnmarshal(&auth, &auth.meta); err != nil {
			return nil, err
		}
		msg = &auth

	case CallMsg:
		var call Call
		if err := genericUnmarshal(&call, &call.meta); err != nil {
			return nil, err
		}
		msg = &call

	case SubMsg:
		var sub Sub
		if err := genericUnmarshal(&sub, &sub.meta); err != nil {
			return nil, err
		}
		msg = &sub

	case PubMsg:
		var pub Pub
		if err := genericUnmarshal(&pub, &pub.meta); err != nil {
			return nil, err
		}
		msg = &pub

	case ErrMsg, OKMsg, ResMsg, EvntMsg:
		return nil, fmt.Errorf("invalid message %s for client peer", pm.Meta.T)
	default:
		return nil, fmt.Errorf("unknown message %s", pm.Meta.T)
	}

	return msg, nil
}
