package juggler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Client is a juggler client based on a websocket connection. It can
// be used to send and receive messages to and from a juggler server.
type Client struct {
	// ResponseHeader is the map of HTTP response headers returned
	// from the initial websocket handshake.
	ResponseHeader http.Header

	conn *websocket.Conn
}

// NewClient creates a juggler client using the provided websocket
// connection and response header.
func NewClient(conn *websocket.Conn, resHeader http.Header) *Client {
	return &Client{
		ResponseHeader: resHeader,
	}
}

// Dial is a helper function to create a Client connected to urlStr using
// the provided *websocket.Dialer and request headers. If the connection
// succeeds, it returns the initialized client, otherwise it returns an
// error. It does not allow handling redirections and such, for a better
// control over the connection, directly use the *websocket.Dialer and
// create the client once the connection is established, using NewClient.
func Dial(d *websocket.Dialer, urlStr string, reqHeader http.Header) (*Client, error) {
	conn, res, err := d.Dial(urlStr, reqHeader)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, res.Header), nil
}

// UnderlyingConn returns the underlying websocket connection used by the
// client. Great care should be taken when using the websocket connection
// directly, as it can cause issues such as data races with the normal
// behaviour of the client.
func (c *Client) UnderlyingConn() *websocket.Conn {
	return c.conn
}
