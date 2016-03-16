package juggler

import (
	"expvar"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/gorilla/websocket"
)

// DiscardLog is a helper no-op function that can be assigned to
// Server.LogFunc to disable logging.
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
// the connections can get served by this server by calling
// Server.ServeConn.
//
// The fields should not be updated once a server has started
// serving connections.
type Server struct {
	// ReadLimit defines the maximum size, in bytes, of incoming
	// messages. If a client sends a message that exceeds this limit,
	// the connection is closed. The default of 0 means no limit.
	ReadLimit int64

	// ReadTimeout is the timeout to read an incoming message. It is
	// set on the websocket connection with SetReadDeadline before
	// reading each message. The default of 0 means no timeout.
	ReadTimeout time.Duration

	// WriteLimit defines the maximum size, in bytes, of outgoing
	// messages. If a message exceeds this limit, it is dropped and
	// an ERR message is sent to the client instead. The default of 0
	// means no limit.
	WriteLimit int64

	// WriteTimeout is the timeout to write an outgoing message. It is
	// set on the websocket connection with SetWriteDeadline before
	// writing each message. The default of 0 means no timeout.
	WriteTimeout time.Duration

	// AcquireWriteLockTimeout is the time to wait for the exclusive
	// write lock for a connection. If the lock cannot be acquired
	// before the timeout, the connection is dropped. The default of
	// 0 means no timeout.
	AcquireWriteLockTimeout time.Duration

	// ConnState specifies an optional callback function that is called
	// when a connection changes state. If non-nil, it is called for
	// Connected and Closing states.
	ConnState func(*Conn, ConnState)

	// Handler is the handler that is called when a message is
	// processed. The ProcessMsg function is called if the default
	// nil value is set. If a custom handler is set, it is assumed
	// that it will call ProcessMsg at some point, or otherwise
	// manually process the messages.
	Handler Handler

	// LogFunc is the function called to log events. By default,
	// it logs using log.Printf. Logging can be disabled by setting
	// LogFunc to DiscardLog.
	LogFunc func(string, ...interface{}) // TODO : normalize calls so that order of args is somewhat predictable

	// PubSubBroker is the broker to use for pub-sub messages. It must be
	// set before the Server can be used.
	PubSubBroker broker.PubSubBroker

	// CallerBroker is the broker to use for caller messages. It must be
	// set before the server can be used.
	CallerBroker broker.CallerBroker

	// Vars can be set to an *expvar.Map to collect metrics about the
	// server. It should be set before starting to listen for
	// connections.
	Vars *expvar.Map
}

// ServeConn serves the websocket connection as a juggler connection. It
// blocks until the juggler connection is closed, leaving the websocket
// connection open.
func (srv *Server) ServeConn(conn *websocket.Conn) {
	if srv.Vars != nil {
		srv.Vars.Add("ActiveConns", 1)
		srv.Vars.Add("TotalConns", 1)
		defer srv.Vars.Add("ActiveConns", -1)
	}

	conn.SetReadLimit(srv.ReadLimit)
	c := newConn(conn, srv)
	resConn, err := srv.CallerBroker.Results(c.UUID)
	if err != nil {
		logf(srv.LogFunc, "failed to create results connection: %v; dropping connection", err)
		return
	}
	pubSubConn, err := srv.PubSubBroker.PubSub()
	if err != nil {
		logf(srv.LogFunc, "failed to create pubsub connection: %v; dropping connection", err)
		return
	}
	c.psc = pubSubConn
	c.resc = resConn

	if cs := srv.ConnState; cs != nil {
		defer func() {
			cs(c, Closing)
		}()
	}

	// start lifecycle of the connection
	if cs := srv.ConnState; cs != nil {
		cs(c, Connected)
	}

	// receive, results loop, pub/sub loop
	go c.pubSub()
	go c.results()
	go c.receive()

	kill := c.CloseNotify()
	<-kill
}

// Upgrade returns an http.Handler that upgrades connections to
// the websocket protocol using upgrader. The websocket connection
// must be upgraded to a supported juggler subprotocol otherwise
// the connection is dropped.
//
// Once connected, the websocket connection is served via srv.ServeConn.
// The websocket connection is closed when the juggler connection is closed.
func Upgrade(upgrader *websocket.Upgrader, srv *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// upgrade the HTTP connection to the websocket protocol
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer wsConn.Close()

		// the agreed-upon subprotocol must be one of the supported ones.
		if !isIn(Subprotocols, wsConn.Subprotocol()) {
			logf(srv.LogFunc, "juggler: no supported subprotocol, closing connection")
			return
		}

		// this call blocks until the juggler connection is closed
		srv.ServeConn(wsConn)
	})
}

func logf(fn func(string, ...interface{}), f string, args ...interface{}) {
	if fn != nil {
		fn(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
