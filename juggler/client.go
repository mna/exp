package juggler

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

// MsgHandler defines the method required to handle a message received
// from the server.
type MsgHandler interface {
	Handle(msg.Msg)
}

// MsgHandlerFunc is a function that implements the MsgHandler interface.
type MsgHandlerFunc func(msg.Msg)

// Handle implements MsgHandler for a MsgHandlerFunc. It calls fn with m.
func (fn MsgHandlerFunc) Handle(m msg.Msg) {
	fn(m)
}

// Client is a juggler client based on a websocket connection. It can
// be used to send and receive messages to and from a juggler server.
type Client struct {
	// ResponseHeader is the map of HTTP response headers returned
	// from the initial websocket handshake.
	ResponseHeader http.Header

	// CallTimeout is the time to wait for the result of a call request.
	// The zero value uses the default timeout of the server.
	CallTimeout time.Duration

	// Handler is the message handler that is called with each message
	// received from the server.
	Handler MsgHandler

	conn *websocket.Conn
}

// NewClient creates a juggler client using the provided websocket
// connection and response header.
func NewClient(conn *websocket.Conn, resHeader http.Header, h MsgHandler) *Client {
	return &Client{
		ResponseHeader: resHeader,
		Handler:        h,
		conn:           conn,
	}
}

// Dial is a helper function to create a Client connected to urlStr using
// the provided *websocket.Dialer and request headers. If the connection
// succeeds, it returns the initialized client, otherwise it returns an
// error. It does not allow handling redirections and such, for a better
// control over the connection, directly use the *websocket.Dialer and
// create the client once the connection is established, using NewClient.
func Dial(d *websocket.Dialer, urlStr string, reqHeader http.Header, h MsgHandler) (*Client, error) {
	conn, res, err := d.Dial(urlStr, reqHeader)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, res.Header, h), nil
}

// UnderlyingConn returns the underlying websocket connection used by the
// client. Great care should be taken when using the websocket connection
// directly, as it can cause issues such as data races with the normal
// behaviour of the client.
func (c *Client) UnderlyingConn() *websocket.Conn {
	return c.conn
}

// Call makes a call request to the server for the remote procedure
// identified by uri. The v value is marshaled as JSON and sent as
// the parameters to the remote procedure. If timeout is > 0, it is used
// as the call-specific timeout, otherwise Client.CallTimeout is used.
//
// It returns the UUID of the call message on success, or an error if
// the call request could not be sent to the server.
func (c *Client) Call(uri string, v interface{}, timeout time.Duration) (uuid.UUID, error) {
	if timeout == 0 {
		timeout = c.CallTimeout
	}
	m, err := msg.NewCall(uri, v, timeout)
	if err != nil {
		return nil, err
	}
	if err := c.conn.WriteJSON(m); err != nil {
		return nil, err
	}
	return m.UUID(), nil
}
