package juggler

import (
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

var dbgClientClosed func(*Client)

// MsgHandler defines the method required to handle a message received
// from the server.
type MsgHandler interface {
	Handle(msg.Msg)
}

// TODO : harmonize MsgHandler with (server) Handler? Or rename ClientHandler?

// MsgHandlerFunc is a function that implements the MsgHandler interface.
type MsgHandlerFunc func(msg.Msg)

// Handle implements MsgHandler for a MsgHandlerFunc. It calls fn with m.
func (fn MsgHandlerFunc) Handle(m msg.Msg) {
	fn(m)
}

// Client is a juggler client based on a websocket connection. It can
// be used to send and receive messages to and from a juggler server.
// It is not safe to call Client methods concurrently.
type Client struct {
	// ResponseHeader is the map of HTTP response headers returned
	// from the initial websocket handshake.
	ResponseHeader http.Header

	// CallTimeout is the time to wait for the result of a call request.
	// The zero value uses the default timeout of the server. Per-call
	// timeouts can also be specified, see Client.Call.
	CallTimeout time.Duration

	// Handler is the message handler that is called with each message
	// received from the server. Each invocation runs in its own
	// goroutine, so proper synchronization must be used when accessing
	// shared data.
	Handler MsgHandler

	// LogFunc is used to log errors that occur outside the handler calls,
	// such as when a message fails to be unmarshaled. If nil, it logs
	// using log.Printf. It can be set to juggler.DiscardLog to disable
	// logging.
	LogFunc func(string, ...interface{})

	wg      sync.WaitGroup
	conn    *websocket.Conn
	mu      sync.Mutex
	results map[string]struct{}
}

// NewClient creates a juggler client using the provided websocket
// connection and response header. Received messages are sent to
// the MsgHandler h.
func NewClient(conn *websocket.Conn, resHeader http.Header, h MsgHandler) *Client {
	c := &Client{
		ResponseHeader: resHeader,
		Handler:        h,
		conn:           conn,
		results:        make(map[string]struct{}),
	}
	c.wg.Add(1)
	go c.handleMessages()
	return c
}

func (c *Client) handleMessages() {
	defer func() {
		c.wg.Done()
		if dbgClientClosed != nil {
			dbgClientClosed(c)
		}
	}()

	for {
		_, r, err := c.conn.NextReader()
		if err != nil {
			logf(c.LogFunc, "client: NextReader failed: %v; stopping read loop", err)
			return
		}

		m, err := msg.UnmarshalResponse(r)
		if err != nil {
			logf(c.LogFunc, "client: UnmarshalResponse failed: %v; skipping message", err)
			continue
		}

		switch m := m.(type) {
		case *msg.Res:
			// got the result, do not trigger an expired message
			k := m.Payload.For.String()
			c.mu.Lock()
			_, ok := c.results[k]
			delete(c.results, k)
			c.mu.Unlock()

			// if an expired message got here first, then drop the
			// result, client treated this call as expired already.
			if !ok {
				continue
			}

		case *msg.Err:
			if m.Payload.ForType == msg.CallMsg {
				// won't get any result for this call
				k := m.Payload.For.String()
				c.mu.Lock()
				_, ok := c.results[k]
				delete(c.results, k)
				c.mu.Unlock()

				// if an expired message got here first, then drop the
				// result, client treated this call as expired already.
				if !ok {
					continue
				}
			}
		}

		go c.Handler.Handle(m)
	}
}

// Dial is a helper function to create a Client connected to urlStr using
// the provided *websocket.Dialer and request headers. If the connection
// succeeds, it returns the initialized client, otherwise it returns an
// error. It does not allow handling redirections and such, for a better
// control over the connection, directly use the *websocket.Dialer and
// create the client once the connection is established, using NewClient.
//
// If the request header doesn't have a Sec-WebSocket-Protocol header,
// it is set to the last entry of juggler.Subprotocols.
func Dial(d *websocket.Dialer, urlStr string, reqHeader http.Header, h MsgHandler) (*Client, error) {
	if reqHeader == nil {
		reqHeader = http.Header{}
	}
	if v := reqHeader["Sec-WebSocket-Protocol"]; (len(v) == 0 || v[0] == "") && len(Subprotocols) > 0 {
		reqHeader["Sec-WebSocket-Protocol"] = []string{Subprotocols[len(Subprotocols)-1]}
	}
	conn, res, err := d.Dial(urlStr, reqHeader)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, res.Header, h), nil
}

// Close closes the connection.
func (c *Client) Close() error {
	err := c.conn.Close()
	c.wg.Wait()
	return err
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

	// add the expected result
	c.mu.Lock()
	c.results[m.UUID().String()] = struct{}{}
	c.mu.Unlock()

	go c.handleExpiredCall(m, timeout)

	return m.UUID(), nil
}

func (c *Client) handleExpiredCall(m *msg.Call, timeout time.Duration) {
	// wait for the timeout
	if timeout <= 0 {
		timeout = time.Minute // TODO : make that available as const somewhere?
	}
	<-time.After(timeout)

	// check if still waiting for a result
	k := m.UUID().String()
	c.mu.Lock()
	_, ok := c.results[k]
	delete(c.results, k)
	c.mu.Unlock()

	if ok {
		// if so, send an Exp message
		exp := msg.NewExp(m)
		go c.Handler.Handle(exp)
	}
}

// Sub makes a subscription request to the server for the specified
// channel, which is treated as a pattern if patter is true. It
// returns the UUID of the sub message on success, or an error if
// the request could not be sent to the server.
func (c *Client) Sub(channel string, pattern bool) (uuid.UUID, error) {
	m := msg.NewSub(channel, pattern)
	if err := c.conn.WriteJSON(m); err != nil {
		return nil, err
	}
	return m.UUID(), nil
}

// Unsb makes an unsubscription request to the server for the specified
// channel, which is treated as a pattern if patter is true. It
// returns the UUID of the unsb message on success, or an error if
// the request could not be sent to the server.
func (c *Client) Unsb(channel string, pattern bool) (uuid.UUID, error) {
	m := msg.NewUnsb(channel, pattern)
	if err := c.conn.WriteJSON(m); err != nil {
		return nil, err
	}
	return m.UUID(), nil
}

// Pub makes a publish request to the server on the specified channel.
// The v value is marshaled as JSON and sent as event payload. It returns
// the UUID of the pub message on success, or an error if the request could
// not be sent to the server.
func (c *Client) Pub(channel string, v interface{}) (uuid.UUID, error) {
	m, err := msg.NewPub(channel, v)
	if err != nil {
		return nil, err
	}
	if err := c.conn.WriteJSON(m); err != nil {
		return nil, err
	}
	return m.UUID(), nil
}
