package juggler

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/PuerkitoBio/exp/juggler/msg"
)

// MsgHandler defines the method required to handle a send or receive
// of a Msg over a connection.
type MsgHandler interface {
	Handle(*Conn, msg.Msg)
}

// MsgHandlerFunc is a function signature that implements the MsgHandler
// interface.
type MsgHandlerFunc func(*Conn, msg.Msg)

// Handle implements MsgHandler for the MsgHandlerFunc by calling the
// function itself.
func (h MsgHandlerFunc) Handle(c *Conn, m msg.Msg) {
	h(c, m)
}

// Chain returns a MsgHandler that calls the provided handlers
// in order, one after the other.
func Chain(hs ...MsgHandler) MsgHandler {
	return MsgHandlerFunc(func(c *Conn, m msg.Msg) {
		for _, h := range hs {
			h.Handle(c, m)
		}
	})
}

// PanicRecover returns a MsgHandler that recovers from panics that
// may happen in h and logs the panic to LogFunc. If close is true,
// the connection is closed on a panic.
func PanicRecover(h MsgHandler, closeConn bool, printStack bool) MsgHandler {
	return MsgHandlerFunc(func(c *Conn, m msg.Msg) {
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

				logf(c.srv, "%v: recovered from panic %v; serving message %v %s", c.UUID, e, m.UUID(), m.Type())
				if printStack {
					b := make([]byte, 4096)
					n := runtime.Stack(b, false)
					logf(c.srv, string(b[:n]))
				}
			}
		}()
		h.Handle(c, m)
	})
}

// LogConn is a function compatible with the Server.ConnState field
// type that logs connections and disconnections to LogFunc.
func LogConn(c *Conn, state ConnState) {
	switch state {
	case Connected:
		logf(c.srv, "%v: connected from %v with subprotocol %q", c.UUID, c.WSConn.RemoteAddr(), c.WSConn.Subprotocol())
	case Closing:
		logf(c.srv, "%v: closing from %v with error %v", c.UUID, c.WSConn.RemoteAddr(), c.CloseErr)
	}
}

// LogMsg is a MsgHandlerFunc that logs messages received or sent on
// c to LogFunc.
func LogMsg(c *Conn, m msg.Msg) {
	if m.Type().IsRead() {
		logf(c.srv, "%v: received message %v %s", c.UUID, m.UUID(), m.Type())
	} else if m.Type().IsWrite() {
		logf(c.srv, "%v: sending message %v %s", c.UUID, m.UUID(), m.Type())
	}
}

// ProcessMsg implements the default message processing. For client messages,
// it calls the appropriate RPC, PUB-SUB or AUTH mechanisms. For server
// messages, it marshals the message and sends it to the client.
//
// When a custom ReadHandler and/or WriterHandler is set on the Server,
// it should at some point call ProcessMsg so the expected behaviour
// happens.
func ProcessMsg(c *Conn, m msg.Msg) {
	switch m := m.(type) {
	case *msg.Auth:
		// TODO : think about it some more...

	case *msg.Call:
		if err := c.srv.pushRedisCall(c.UUID, m); err != nil {
			e := msg.NewErr(m, 500, err) // TODO : use HTTP-like error codes?
			c.Send(e)
			return
		}
		ok := msg.NewOK(m)
		c.Send(ok)

	case *msg.Pub:
	case *msg.Sub:

	case *msg.OK, *msg.Err, *msg.Evnt, *msg.Res:
		writeMsg(c, m)

	default:
		logf(c.srv, "unknown message in ProcessMsg: %T", m)
	}
}

func writeMsg(c *Conn, m msg.Msg) {
	w := c.Writer(c.srv.AcquireWriteLockTimeout)
	defer w.Close()

	if err := json.NewEncoder(w).Encode(m); err != nil {
		if err == ErrLockWriterTimeout {
			c.Close(fmt.Errorf("writeMsg failed: %v; closing connection", err))
			return
		}
		logf(c.srv, "%v: writeMsg %v failed: %v", c.UUID, m.UUID(), err)
		return
	}
}
