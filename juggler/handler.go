package juggler

import (
	"encoding/json"
	"fmt"
)

// MsgHandler defines the method required to handle a send or receive
// of a Msg over a connection.
type MsgHandler interface {
	Handle(*Conn, Msg)
}

// MsgHandlerFunc is a function signature that implements the MsgHandler
// interface.
type MsgHandlerFunc func(*Conn, Msg)

// Handle implements MsgHandler for the MsgHandlerFunc by calling the
// function itself.
func (h MsgHandlerFunc) Handle(c *Conn, msg Msg) {
	h(c, msg)
}

// Chain returns a MsgHandler that calls the provided handlers
// in order, one after the other.
func Chain(hs ...MsgHandler) MsgHandler {
	return MsgHandlerFunc(func(c *Conn, msg Msg) {
		for _, h := range hs {
			h.Handle(c, msg)
		}
	})
}

// PanicRecover returns a MsgHandler that recovers from panics that
// may happen in h and logs the panic to LogFunc. If close is true,
// the connection is closed on a panic.
func PanicRecover(h MsgHandler, close bool) MsgHandler {
	return MsgHandlerFunc(func(c *Conn, msg Msg) {
		defer func() {
			if e := recover(); e != nil {
				if close {
					var err error
					switch e := e.(type) {
					case error:
						err = e
					default:
						err = fmt.Errorf("%v", e)
					}
					c.Close(err)
				}

				logf(c.srv, "%v: recovered from panic %v; serving message %v %s", c.UUID, e, msg.UUID(), msg.Type())
			}
		}()
		h.Handle(c, msg)
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
func LogMsg(c *Conn, msg Msg) {
	if msg.IsRead() {
		logf(c.srv, "%v: received message %v %s", c.UUID, msg.UUID(), msg.Type())
	} else if msg.IsWrite() {
		logf(c.srv, "%v: sending message %v %s", c.UUID, msg.UUID(), msg.Type())
	}
}

// ProcessMsg implements the default message processing. For client messages,
// it calls the appropriate RPC, PUB-SUB or AUTH mechanisms. For server
// messages, it marshals the message and sends it to the client.
//
// When a custom ReadHandler and/or WriterHandler is set on the Server,
// it should at some point call ProcessMsg so the expected behaviour
// happens.
func ProcessMsg(c *Conn, msg Msg) {
	switch msg := msg.(type) {
	case *Auth:
	case *Call:
		if err := c.srv.redisCall(msg); err != nil {
			e := newErr(msg, 500, err.Error()) // TODO : use HTTP-like error codes?
			c.Send(e)
		}
		ok := newOK(msg)
		c.Send(ok)

	case *Pub:
	case *Sub:

	case *OK, *Err, *Evnt, *Res:
		writeMsg(c, msg)

	default:
		logf(c.srv, "unknown message in ProcessMsg: %T", msg)
	}
}

func writeMsg(c *Conn, msg Msg) {
	w := c.Writer(c.srv.AcquireWriteLockTimeout)
	defer w.Close()

	if err := json.NewEncoder(w).Encode(msg); err != nil {
		if err == ErrLockWriterTimeout {
			c.Close(fmt.Errorf("writeMsg failed: %v; closing connection", err))
			return
		}
		logf(c.srv, "%v: writeMsg %v failed: %v", c.UUID, msg.UUID(), err)
		return
	}
}
