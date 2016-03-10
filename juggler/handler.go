package juggler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler/msg"
)

// Handler defines the method required to handle a send or receive
// of a Msg over a connection.
type Handler interface {
	Handle(context.Context, *Conn, msg.Msg)
}

// HandlerFunc is a function signature that implements the Handler
// interface.
type HandlerFunc func(context.Context, *Conn, msg.Msg)

// Handle implements Handler for the HandlerFunc by calling the
// function itself.
func (h HandlerFunc) Handle(ctx context.Context, c *Conn, m msg.Msg) {
	h(ctx, c, m)
}

// Chain returns a Handler that calls the provided handlers
// in order, one after the other.
func Chain(hs ...Handler) Handler {
	return HandlerFunc(func(ctx context.Context, c *Conn, m msg.Msg) {
		for _, h := range hs {
			h.Handle(ctx, c, m)
		}
	})
}

// PanicRecover returns a Handler that recovers from panics that
// may happen in h and logs the panic to LogFunc. If close is true,
// the connection is closed on a panic.
func PanicRecover(h Handler, closeConn bool, printStack bool) Handler {
	return HandlerFunc(func(ctx context.Context, c *Conn, m msg.Msg) {
		defer func() {
			if e := recover(); e != nil {
				if closeConn {
					var err error
					switch e := e.(type) {
					case error:
						err = e
					default:
						err = fmt.Errorf("%v", e)
					}
					c.Close(err)
				}

				logf(c.srv.LogFunc, "%v: recovered from panic %v; serving message %v %s", c.UUID, e, m.UUID(), m.Type())
				if printStack {
					b := make([]byte, 4096)
					n := runtime.Stack(b, false)
					logf(c.srv.LogFunc, string(b[:n]))
				}
			}
		}()
		h.Handle(ctx, c, m)
	})
}

// LogConn is a function compatible with the Server.ConnState field
// type that logs connections and disconnections to LogFunc.
func LogConn(c *Conn, state ConnState) {
	switch state {
	case Connected:
		logf(c.srv.LogFunc, "%v: connected from %v with subprotocol %q", c.UUID, c.RemoteAddr(), c.Subprotocol())
	case Closing:
		logf(c.srv.LogFunc, "%v: closing from %v with error %v", c.UUID, c.RemoteAddr(), c.CloseErr)
	}
}

// LogMsg is a HandlerFunc that logs messages received or sent on
// c to LogFunc.
func LogMsg(ctx context.Context, c *Conn, m msg.Msg) {
	if m.Type().IsRead() {
		logf(c.srv.LogFunc, "%v: received message %v %s", c.UUID, m.UUID(), m.Type())
	} else if m.Type().IsWrite() {
		logf(c.srv.LogFunc, "%v: sending message %v %s", c.UUID, m.UUID(), m.Type())
	}
}

// ProcessMsg implements the default message processing. For client messages,
// it calls the appropriate RPC, PUB-SUB or AUTH mechanisms. For server
// messages, it marshals the message and sends it to the client.
//
// When a custom ReadHandler and/or WriterHandler is set on the Server,
// it should at some point call ProcessMsg so the expected behaviour
// happens.
func ProcessMsg(ctx context.Context, c *Conn, m msg.Msg) {
	switch m := m.(type) {
	case *msg.Call:
		cp := &msg.CallPayload{
			ConnUUID: c.UUID,
			MsgUUID:  m.UUID(),
			URI:      m.Payload.URI,
			Args:     m.Payload.Args,
		}
		if err := c.srv.CallerBroker.Call(cp, m.Payload.Timeout); err != nil {
			c.Send(msg.NewErr(m, 500, err))
			return
		}
		c.Send(msg.NewOK(m))

	case *msg.Pub:
		pp := &msg.PubPayload{
			MsgUUID: m.UUID(),
			Args:    m.Payload.Args,
		}
		if err := c.srv.PubSubBroker.Publish(m.Payload.Channel, pp); err != nil {
			c.Send(msg.NewErr(m, 500, err))
			return
		}
		c.Send(msg.NewOK(m))

	case *msg.Sub:
		if err := c.psc.Subscribe(m.Payload.Channel, m.Payload.Pattern); err != nil {
			c.Send(msg.NewErr(m, 500, err))
			return
		}
		c.Send(msg.NewOK(m))

	case *msg.Unsb:
		if err := c.psc.Unsubscribe(m.Payload.Channel, m.Payload.Pattern); err != nil {
			c.Send(msg.NewErr(m, 500, err))
			return
		}
		c.Send(msg.NewOK(m))

	case *msg.OK, *msg.Err, *msg.Evnt, *msg.Res:
		if err := writeMsg(c, m); err != nil {
			switch err {
			case ErrLockWriterTimeout:
				c.Close(fmt.Errorf("writeMsg failed: %v; closing connection", err))

			case errWriteLimitExceeded:
				logf(c.srv.LogFunc, "%v: writeMsg %v failed: %v", c.UUID, m.UUID(), err)
				// no good http code for this case
				if err := writeMsg(c, msg.NewErr(m, 599, err)); err != nil {
					if err == ErrLockWriterTimeout {
						c.Close(fmt.Errorf("writeMsg failed: %v; closing connection", err))
					} else {
						logf(c.srv.LogFunc, "%v: writeMsg %v for write limit exceeded notification failed: %v", c.UUID, m.UUID(), err)
					}
					return
				}

			default:
				logf(c.srv.LogFunc, "%v: writeMsg %v failed: %v", c.UUID, m.UUID(), err)
			}
		}

	default:
		logf(c.srv.LogFunc, "unknown message in ProcessMsg: %T", m)
	}
}

var errWriteLimitExceeded = errors.New("write limit exceeded")

type limitedWriter struct {
	w io.Writer
	n int64
}

const minWriteLimit = 4096

func limitWriter(w io.Writer, limit int64) io.Writer {
	if limit < minWriteLimit {
		limit = minWriteLimit
	}
	return &limitedWriter{w: w, n: limit}
}

func (w *limitedWriter) Write(p []byte) (int, error) {
	w.n -= int64(len(p))
	if w.n < 0 {
		return 0, errWriteLimitExceeded
	}
	return w.w.Write(p)
}

func writeMsg(c *Conn, m msg.Msg) error {
	w := c.Writer(c.srv.AcquireWriteLockTimeout)
	defer w.Close()

	lw := io.Writer(w)
	if l := c.srv.WriteLimit; l > 0 {
		lw = limitWriter(w, l)
	}
	if err := json.NewEncoder(lw).Encode(m); err != nil {
		return err
	}
	return nil
}
