package juggler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// LogFunc is the function called to log events. It should never be
// set to nil, use DiscardLog instead to disable logging. By default,
// it logs using log.Printf.
var LogFunc = log.Printf

// DiscardLog is a helper no-op function that can be assigned to LogFunc
// to disable logging.
func DiscardLog(f string, args ...interface{}) {}

// Subprotocols is the list of juggler protocol versions supported by this
// package. It should be set as-is on the websocket.Upgrader Subprotocols
// field.
var Subprotocols = []string{
	"juggler.1",
}

// Server is a juggler server. Once a websocket handshake has been
// established with a juggler subprotocol over a standard HTTP server,
// the connections get served by this server.
type Server struct {
	// ReadLimit defines the maximum size, in bytes, of incoming
	// messages. If a client sends a message that exceeds this limit,
	// the connection is closed. The default of 0 means no limit.
	ReadLimit int64

	// ReadTimeout is the timeout to read an incoming message. It is
	// set on the websocket connection with SetReadDeadline before
	// reading each message. The default of 0 means no timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the timeout to write an outgoing message. It is
	// set on the websocket connection with SetWriteDeadline before
	// writing each message. The default of 0 means no timeout.
	WriteTimeout time.Duration

	// ConnState specifies an optional callback function that is called
	// when a connection changes state. It is called for Connected and
	// Closing states.
	ConnState func(*Conn, ConnState)

	// ReadHandler is the handler that is called when an incoming
	// message is processed. The ProcessMsg function is called
	// if the default nil value is set. If a custom handler is set,
	// it is assumed that it will call ProcessMsg at some point,
	// or otherwise manually process the messages.
	ReadHandler MsgHandler

	// WriteHandler is the handler that is called when an outgoing
	// message is processed. The ProcessMsg function is called
	// if the default nil value is set. If a custom handler is set,
	// it is assumed that it will call ProcessMsg at some point,
	// or otherwise manually process the messages.
	WriteHandler MsgHandler
}

func Upgrade(upgrader *websocket.Upgrader, srv *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// upgrade the HTTP connection to the websocket protocol
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer wsConn.Close()
		if wsConn.Subprotocol() == "" {
			LogFunc("juggler: no supported subprotocol, closing connection")
			return
		}

		// configure the websocket connection
		wsConn.SetReadLimit(srv.ReadLimit)
		c := newConn(wsConn)
		defer func() {
			if srv.ConnState != nil {
				srv.ConnState(c, Closing)
			}
		}()

		// start lifecycle of the connection
		if srv.ConnState != nil {
			srv.ConnState(c, Connected)
		}

		if err := srv.read(c); err != nil {
			c.setState(Closing, err)
			LogFunc("juggler: read failed: %v; closing connection", err)
			return
		}

		/*
			if err := c.WSConn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
				c.setState(Closing, err)
				LogFunc("juggler: WriteMessage failed: %v; closing connection", err)
				return
			}
			for {
				c.WSConn.SetWriteDeadline(time.Time{})

				mt, r, err := c.WSConn.NextReader()
				if err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: NextReader failed: %v; closing connection", err)
					return
				}
				c.WSConn.SetReadDeadline(time.Now().Add(srv.ReadTimeout))

				w, err := c.WSConn.NextWriter(mt)
				if err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: NextWriter failed: %v; closing connection", err)
					return
				}
				c.WSConn.SetWriteDeadline(time.Now().Add(srv.WriteTimeout))

				if _, err := io.Copy(w, r); err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: Copy failed: %v; closing connection", err)
					return
				}
				if err := w.Close(); err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: Close failed: %v; closing connection", err)
					return
				}
			}
		*/
	})
}

func (s *Server) read(c *Conn) error {
	for {
		c.WSConn.SetReadDeadline(time.Time{})

		mt, r, err := c.WSConn.NextReader()
		if err != nil {
			return err
		}
		if mt != websocket.TextMessage {
			return fmt.Errorf("invalid websocket message type: %d", mt)
		}
		if s.ReadTimeout > 0 {
			c.WSConn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		}

		msg, err := unmarshalMessage(r)
		if err != nil {
			return err
		}

		if s.ReadHandler != nil {
			s.ReadHandler.Handle(c, msg)
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
