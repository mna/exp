package juggler

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler/broker"
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

	// CloseErr is the error, if any, that caused the connection
	// to close. Must only be accessed after the close notification
	// has been received (i.e. after a <-conn.CloseNotify()).
	CloseErr error

	// TODO : some connection state (authenticated, etc.)?

	// the underlying websocket connection.
	wsConn *websocket.Conn

	wmu  chan struct{} // write lock
	srv  *Server
	psc  broker.PubSubConn  // single pub-sub-dedicated broker connection
	resc broker.ResultsConn // single results-dedicated broker connection

	// ensure the kill channel can only be closed once
	closeOnce sync.Once
	kill      chan struct{}
}

func newConn(c *websocket.Conn, srv *Server) *Conn {
	// wmu is the write lock, used as a semaphore of 1, so start with
	// an available slot (initialize with a sent value).
	wmu := make(chan struct{}, 1)
	wmu <- struct{}{}

	return &Conn{
		UUID:   uuid.NewRandom(),
		wsConn: c,
		wmu:    wmu,
		srv:    srv,
		kill:   make(chan struct{}),
	}
}

// UnderlyingConn returns the underlying websocket connection. Great
// care should be taken when using the websocket connection directly,
// as it may interfere and create data races with the juggler connection.
func (c *Conn) UnderlyingConn() *websocket.Conn {
	return c.wsConn
}

// CloseNotify returns a signal channel that is closed when the
// Conn is closed.
func (c *Conn) CloseNotify() <-chan struct{} {
	return c.kill
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.wsConn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.wsConn.RemoteAddr()
}

// Subprotocol returns the negotiated protocol for the connection.
func (c *Conn) Subprotocol() string {
	return c.wsConn.Subprotocol()
}

// Close closes the connection, setting err as CloseErr to identify
// the reason of the close. It does not send a websocket close message,
// nor does it close the underlying websocket connection.
// As with all Conn methods, it is safe to call concurrently, but
// only the first call will set the CloseErr field to err.
func (c *Conn) Close(err error) {
	c.closeOnce.Do(func() {
		c.CloseErr = err
		c.psc.Close()
		c.resc.Close()
		close(c.kill)
	})
}

// writer that acquires the connection's write lock prior to writing.
type exclusiveWriter struct {
	w            io.WriteCloser
	init         bool
	writeLock    chan struct{}
	lockTimeout  time.Duration
	writeTimeout time.Duration
	wsConn       *websocket.Conn
}

func (w *exclusiveWriter) Write(p []byte) (int, error) {
	if !w.init {
		var wait <-chan time.Time
		if to := w.lockTimeout; to > 0 {
			wait = time.After(to)
		}

		// try to acquire the write lock before the timeout
		select {
		case <-wait:
			return 0, ErrLockWriterTimeout

		case <-w.writeLock:
			// lock acquired, get next writer from the websocket connection
			w.init = true
			wc, err := w.wsConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return 0, err
			}
			w.w = wc
			if to := w.writeTimeout; to > 0 {
				w.wsConn.SetWriteDeadline(time.Now().Add(to))
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

	var err error
	if w.w != nil {
		// if w.init is true, then NextWriter was called and that writer
		// must be properly closed.
		err = w.w.Close()
		w.wsConn.SetWriteDeadline(time.Time{})
	}

	// release the write lock
	w.writeLock <- struct{}{}
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
		writeLock:    c.wmu,
		lockTimeout:  timeout,
		writeTimeout: c.srv.WriteTimeout,
		wsConn:       c.wsConn,
	}
}

// Send sends the msg to the client. It calls the Server's
// WriteHandler if any, or ProcessMsg if nil.
func (c *Conn) Send(m msg.Msg) {
	if wh := c.srv.WriteHandler; wh != nil {
		wh.Handle(context.Background(), c, m)
	} else {
		ProcessMsg(context.Background(), c, m)
	}
}

// results is the loop that looks for call results.
func (c *Conn) results() {
	ch := c.resc.Results()
	for res := range ch {
		c.Send(msg.NewRes(res))
	}

	// results loop was stopped, the connection should be closed if it
	// isn't already.
	c.Close(c.resc.ResultsErr())
}

// pubSub is the loop that receives events that the connection is subscribed
// to.
func (c *Conn) pubSub() {
	ch := c.psc.Events()
	for ev := range ch {
		c.Send(msg.NewEvnt(ev))
	}

	// pubsub loop was stopped, the connection should be closed if it
	// isn't already.
	c.Close(c.psc.EventsErr())
}

// receive is the read loop, started in its own goroutine.
func (c *Conn) receive() {
	for {
		c.wsConn.SetReadDeadline(time.Time{})

		// NextReader returns with an error once a connection is closed,
		// so this loop doesn't need to check the c.kill channel.
		mt, r, err := c.wsConn.NextReader()
		if err != nil {
			c.Close(err)
			return
		}
		if mt != websocket.TextMessage {
			c.Close(fmt.Errorf("invalid websocket message type: %d", mt))
			return
		}
		if to := c.srv.ReadTimeout; to > 0 {
			c.wsConn.SetReadDeadline(time.Now().Add(to))
		}

		m, err := msg.UnmarshalRequest(r)
		if err != nil {
			c.Close(err)
			return
		}

		if rh := c.srv.ReadHandler; rh != nil {
			rh.Handle(context.Background(), c, m)
		} else {
			ProcessMsg(context.Background(), c, m)
		}
	}
}
