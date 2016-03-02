package juggler

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// LogFunc is the function called to log events. It must never be
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
	"juggler.0",
}

func isIn(list []string, v string) bool {
	for _, vv := range list {
		if vv == v {
			return true
		}
	}
	return false
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
	ConnState func(*Conn, ConnState, error)

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

// ListenAndServe starts a default HTTP server by calling
// http.ListenAndServe on the provided address, and upgrades requests
// made to path to a websocket connection for a supported juggler
// subprotocol. The provided read and write MsgHandlers are used
// to process messages.
func ListenAndServe(addr, path string, read, write MsgHandler) error {
	upg := &websocket.Upgrader{Subprotocols: Subprotocols}
	srv := &Server{ReadHandler: read, WriteHandler: write}
	mux := http.NewServeMux()
	mux.Handle(path, Upgrade(upg, srv))
	return http.ListenAndServe(addr, mux)
}

// Upgrade returns an http.Handler that upgrades connections to
// the websocket protocol using upgrader. The websocket connection
// must be upgraded to a supported juggler subprotocol otherwise
// the connection is dropped.
//
// Once connected, the websocket connection is served via srv.
func Upgrade(upgrader *websocket.Upgrader, srv *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// upgrade the HTTP connection to the websocket protocol
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer wsConn.Close()
		if wsConn.Subprotocol() == "" || !isIn(Subprotocols, wsConn.Subprotocol()) {
			LogFunc("juggler: no supported subprotocol, closing connection")
			return
		}

		var closeErr error
		wsConn.SetReadLimit(srv.ReadLimit)
		c := newConn(wsConn, srv)
		defer func() {
			if srv.ConnState != nil {
				srv.ConnState(c, Closing, closeErr)
			}
		}()

		// start lifecycle of the connection
		if srv.ConnState != nil {
			srv.ConnState(c, Connected, nil)
		}

		if err := c.receive(); err != nil {
			closeErr = err
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
