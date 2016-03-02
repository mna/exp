package juggler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	UUID   uuid.UUID
	WSConn *websocket.Conn // TODO : hide/show only as needed

	// TODO : some connection state (authenticated, etc.)
	wmu chan struct{} // write lock
	srv *Server
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
			// TODO : set write deadline
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
	// TODO : reset write deadline
	w.c.wmu <- struct{}{}
	return err
}

func (c *Conn) Writer(timeout time.Duration) io.WriteCloser {
	return &exclusiveWriter{
		c:       c,
		timeout: timeout,
	}
}

func (c *Conn) receive() error {
	for {
		c.WSConn.SetReadDeadline(time.Time{})

		mt, r, err := c.WSConn.NextReader()
		if err != nil {
			return err
		}
		if mt != websocket.TextMessage {
			return fmt.Errorf("invalid websocket message type: %d", mt)
		}
		if c.srv.ReadTimeout > 0 {
			c.WSConn.SetReadDeadline(time.Now().Add(c.srv.ReadTimeout))
		}

		msg, err := unmarshalMessage(r)
		if err != nil {
			return err
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
