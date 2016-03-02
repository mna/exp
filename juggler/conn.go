package juggler

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

var (
	ErrLockWriterTimeout = errors.New("juggler: timed out waiting for writer lock")
)

type ConnState int

const (
	Unknown ConnState = iota
	Connected
	Closing
)

type ConnHandler interface {
	Handle(*Conn)
}

type Conn struct {
	UUID uuid.UUID

	// TODO : hide/show only as needed
	WSConn *websocket.Conn
	// TODO : some connection state (authenticated, etc.)

	wmu chan struct{}

	mu       sync.RWMutex
	state    ConnState
	closeErr error
}

func newConn(c *websocket.Conn) *Conn {
	// wmu is the write lock, used as a semaphore of 1, so start with
	// an available slot (initialize with a sent value).
	wmu := make(chan struct{}, 1)
	wmu <- struct{}{}

	return &Conn{
		UUID:   uuid.NewRandom(),
		WSConn: c,
		wmu:    wmu,
		state:  Connected,
	}
}

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
		select {
		case <-wait:
			return 0, ErrLockWriterTimeout
		case <-w.c.wmu:
			w.init = true
			wc, err := w.c.WSConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return 0, err
			}
			w.w = wc
		}
	}

	return w.w.Write(p)
}

func (w *exclusiveWriter) Close() error {
	if !w.init {
		// no write, Close is a no-op
		return nil
	}
	err := w.w.Close()
	w.c.wmu <- struct{}{}
	return err
}

func (c *Conn) Writer(timeout time.Duration) io.WriteCloser {
	return &exclusiveWriter{
		c:       c,
		timeout: timeout,
	}
}

func (c *Conn) State() (state ConnState, closeErr error) {
	c.mu.RLock()
	st, err := c.state, c.closeErr
	c.mu.RUnlock()
	return st, err
}

func (c *Conn) setState(state ConnState, err error) {
	c.mu.Lock()
	c.state = state
	c.closeErr = err
	c.mu.Unlock()
}
