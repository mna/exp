package juggler

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

// ClientHandler defines the method required to handle a message received
// from the server.
type ClientHandler interface {
	Handle(context.Context, *Client, msg.Msg)
}

// ClientHandlerFunc is a function that implements the ClientHandler interface.
type ClientHandlerFunc func(context.Context, *Client, msg.Msg)

// Handle implements ClientHandler for a ClientHandlerFunc. It calls fn
// with the parameters.
func (fn ClientHandlerFunc) Handle(ctx context.Context, cli *Client, m msg.Msg) {
	fn(ctx, cli, m)
}

// ClientOption sets an option on the Client.
type ClientOption func(*Client)

// SetCallTimeout sets the time to wait for the result of a call request.
// The zero value uses the default timeout of the server. Per-call
// timeouts can also be specified, see Client.Call.
func SetCallTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.callTimeout = timeout
	}
}

// SetClientHandler sets the client handler that is called with each message
// received from the server. Each invocation runs in its own
// goroutine, so proper synchronization must be used when accessing
// shared data.
func SetClientHandler(h ClientHandler) ClientOption {
	return func(c *Client) {
		c.handler = h
	}
}

// SetLogFunc sets the function used to log errors that occur outside
// the handler calls, such as when a message fails to be unmarshaled.
// If nil, it logs using log.Printf. It can be set to juggler.DiscardLog
// to disable logging.
func SetLogFunc(fn func(string, ...interface{})) ClientOption {
	return func(c *Client) {
		c.logFunc = fn
	}
}

// Client is a juggler client based on a websocket connection. It can
// be used to send and receive messages to and from a juggler server.
type Client struct {
	// ResponseHeader is the map of HTTP response headers returned
	// from the initial websocket handshake.
	ResponseHeader http.Header

	callTimeout time.Duration
	handler     ClientHandler
	logFunc     func(string, ...interface{})

	wg      sync.WaitGroup // wait for handleMessages goroutine
	stop    chan struct{}  // stop signal for expiration goroutines
	conn    *websocket.Conn
	mu      sync.Mutex // lock access to results map
	results map[string]struct{}
}

// NewClient creates a juggler client using the provided websocket
// connection and response header. Received messages are sent to
// the handler set by the SetClientHandler option.
func NewClient(conn *websocket.Conn, resHeader http.Header, opts ...ClientOption) *Client {
	c := &Client{
		ResponseHeader: resHeader,
		conn:           conn,
		stop:           make(chan struct{}),
		results:        make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.wg.Add(1)
	go c.handleMessages()
	return c
}

func (c *Client) handleMessages() {
	defer func() {
		close(c.stop)
		c.wg.Done()
	}()

	for {
		_, r, err := c.conn.NextReader()
		if err != nil {
			logf(c.logFunc, "client: NextReader failed: %v; stopping read loop", err)
			return
		}

		m, err := msg.UnmarshalResponse(r)
		if err != nil {
			logf(c.logFunc, "client: UnmarshalResponse failed: %v; skipping message", err)
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

		go c.handler.Handle(context.Background(), c, m)
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
func Dial(d *websocket.Dialer, urlStr string, reqHeader http.Header, opts ...ClientOption) (*Client, error) {
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
	return NewClient(conn, res.Header, opts...), nil
}

// Close closes the connection. No more messages will be received.
func (c *Client) Close() error {
	err := c.conn.Close()
	c.wg.Wait()
	return err
}

// UnderlyingConn returns the underlying websocket connection used by the
// client. Care should be taken when using the websocket connection
// directly, as it may interfere with the normal behaviour of the client.
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
		timeout = c.callTimeout
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
		timeout = broker.DefaultCallTimeout
	}
	select {
	case <-c.stop:
		return
	case <-time.After(timeout):
	}

	// check if still waiting for a result
	k := m.UUID().String()
	c.mu.Lock()
	_, ok := c.results[k]
	delete(c.results, k)
	c.mu.Unlock()

	if ok {
		// if so, send an Exp message
		exp := msg.NewExp(m)
		go c.handler.Handle(context.Background(), c, exp)
	}
}

// Sub makes a subscription request to the server for the specified
// channel, which is treated as a pattern if pattern is true. It
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
// channel, which is treated as a pattern if pattern is true. It
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
