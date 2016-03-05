package juggler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
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
	kill      chan struct{} // signal channel, closed when Close is called
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
	if !w.init || w.w == nil {
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

// Writer returns an io.WriteCloser that can be used to send a
// message on the connection. Only one writer can be active at
// any moment for a given connection, so the returned writer
// will acquire a lock on the first call to Write, and will
// release it only when Close is called. The timeout controls
// the time to wait to acquire the lock on the first call to
// Write. If the lock cannot be acquired within that time,
// ErrLockWriterTimeout is returned and no write is performed.
//
// It is possible to enter a deadlock state if Writer is called
// with no timeout, an initial Write is executed, and Writer is
// called again from the same goroutine, without a timeout.
// To avoid this, make sure each goroutine closes the Writer
// before asking for another one, and ideally always use a timeout.
//
// The returned writer itself is not safe for concurrent use, but
// as all Conn methods, Writer can be called concurrently.
func (c *Conn) Writer(timeout time.Duration) io.WriteCloser {
	return &exclusiveWriter{
		c:       c,
		timeout: timeout,
	}
}

// Send sends the msg to the client. It calls the Server's
// WriteHandler if any, or ProcessMsg if nil.
func (c *Conn) Send(m msg.Msg) {
	if c.srv.WriteHandler != nil {
		c.srv.WriteHandler.Handle(c, m)
	} else {
		ProcessMsg(c, m)
	}
}

func (c *Conn) receive() {
	for {
		c.WSConn.SetReadDeadline(time.Time{})

		// NextReader returns with an error once a connection is closed,
		// so this loop doesn't need to check the c.kill channel.
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

func unmarshalMessage(r io.Reader) (msg.Msg, error) {
	var pm msg.PartialMsg
	if err := json.NewDecoder(r).Decode(&pm); err != nil {
		return nil, fmt.Errorf("invalid JSON message: %v", err)
	}

	genericUnmarshal := func(v interface{}, metaDst *msg.Meta) error {
		if err := json.Unmarshal(pm.Payload, v); err != nil {
			return fmt.Errorf("invalid %s message: %v", pm.Meta.T, err)
		}
		*metaDst = pm.Meta
		return nil
	}

	var m msg.Msg
	switch pm.Meta.T {
	case msg.AuthMsg:
		var auth msg.Auth
		if err := genericUnmarshal(&auth, &auth.Meta); err != nil {
			return nil, err
		}
		m = &auth

	case msg.CallMsg:
		var call msg.Call
		if err := genericUnmarshal(&call, &call.Meta); err != nil {
			return nil, err
		}
		m = &call

	case msg.SubMsg:
		var sub msg.Sub
		if err := genericUnmarshal(&sub, &sub.Meta); err != nil {
			return nil, err
		}
		m = &sub

	case msg.UnsbMsg:
		var uns msg.Unsb
		if err := genericUnmarshal(&uns, &uns.Meta); err != nil {
			return nil, err
		}
		m = &uns

	case msg.PubMsg:
		var pub msg.Pub
		if err := genericUnmarshal(&pub, &pub.Meta); err != nil {
			return nil, err
		}
		m = &pub

	case msg.ErrMsg, msg.OKMsg, msg.ResMsg, msg.EvntMsg:
		return nil, fmt.Errorf("invalid message %s for client peer", pm.Meta.T)
	default:
		return nil, fmt.Errorf("unknown message %s", pm.Meta.T)
	}

	return m, nil
}
