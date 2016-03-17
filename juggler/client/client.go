// Package client implements a juggler client.
package client

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

// Client is a juggler client based on a websocket connection. It can
// be used to send and receive messages to and from a juggler server.
type Client struct {
	// ResponseHeader is the map of HTTP response headers returned
	// from the initial websocket handshake.
	ResponseHeader http.Header

	callTimeout time.Duration
	handler     Handler
	logFunc     func(string, ...interface{})

	wg      sync.WaitGroup // wait for handleMessages goroutine
	stop    chan struct{}  // stop signal for expiration goroutines
	conn    *websocket.Conn
	mu      sync.Mutex // lock access to results map
	results map[string]struct{}
}

// NewClient creates a juggler client using the provided websocket
// connection and response header. Received messages are sent to
// the handler set by the SetHandler option.
func NewClient(conn *websocket.Conn, resHeader http.Header, opts ...Option) *Client {
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
			if ok := c.deletePending(m.Payload.For.String()); !ok {
				// if an expired message got here first, then drop the
				// result, client treated this call as expired already.
				continue
			}

		case *msg.Err:
			if m.Payload.ForType == msg.CallMsg {
				// won't get any result for this call
				if ok := c.deletePending(m.Payload.For.String()); !ok {
					// if an expired message got here first, then drop the
					// result, client treated this call as expired already.
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
// The Dialer's Subprotocols field should be set to one of (or any/all of)
// juggler.Subprotocol.
func Dial(d *websocket.Dialer, urlStr string, reqHeader http.Header, opts ...Option) (*Client, error) {
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

// CloseNotify returns a channel that is closed when the client is
// closed.
func (c *Client) CloseNotify() <-chan struct{} {
	return c.stop
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
	c.addPending(m.UUID().String())

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
	if ok := c.deletePending(m.UUID().String()); ok {
		// if so, send an Exp message
		exp := newExp(m)
		go c.handler.Handle(context.Background(), c, exp)
	}
}

// add a pending call.
func (c *Client) addPending(key string) {
	c.mu.Lock()
	c.results[key] = struct{}{}
	c.mu.Unlock()
}

// delete the pending call, returning true if it was still pending.
func (c *Client) deletePending(key string) bool {
	c.mu.Lock()
	_, ok := c.results[key]
	delete(c.results, key)
	c.mu.Unlock()

	return ok
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

// Handler defines the method required to handle a message received
// from the server.
type Handler interface {
	Handle(context.Context, *Client, msg.Msg)
}

// HandlerFunc is a function that implements the Handler interface.
type HandlerFunc func(context.Context, *Client, msg.Msg)

// Handle implements Handler for a HandlerFunc. It calls fn
// with the parameters.
func (fn HandlerFunc) Handle(ctx context.Context, cli *Client, m msg.Msg) {
	fn(ctx, cli, m)
}

// Option sets an option on the Client.
type Option func(*Client)

// SetCallTimeout sets the time to wait for the result of a call request.
// The zero value uses the default timeout of the server. Per-call
// timeouts can also be specified, see Client.Call.
func SetCallTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.callTimeout = timeout
	}
}

// SetHandler sets the handler that is called with each message
// received from the server. Each invocation runs in its own
// goroutine, so proper synchronization must be used when accessing
// shared data.
func SetHandler(h Handler) Option {
	return func(c *Client) {
		c.handler = h
	}
}

// SetLogFunc sets the function used to log errors that occur outside
// the handler calls, such as when a message fails to be unmarshaled.
// If nil, it logs using log.Printf. It can be set to juggler.DiscardLog
// to disable logging.
func SetLogFunc(fn func(string, ...interface{})) Option {
	return func(c *Client) {
		c.logFunc = fn
	}
}

// Exp is an expired call message. It is never sent over the network, but
// it is raised by the client for itself, when the timeout for a call
// result has expired. As such, its message type returns false for
// both IsRead and IsWrite.
type Exp struct {
	msg.Meta `json:"meta"`
	Payload  struct {
		For  uuid.UUID       `json:"for"`           // no ForType, because always CALL
		URI  string          `json:"uri,omitempty"` // URI of the CALL
		Args json.RawMessage `json:"args"`
	} `json:"payload"`
}

// ExpMsg is the message type of the call expiration message.
const ExpMsg msg.MessageType = msg.CustomMsg + iota

// newExp creates a new expired message for the provided call message.
func newExp(m *msg.Call) *Exp {
	exp := &Exp{
		Meta: msg.NewMeta(ExpMsg),
	}
	exp.Payload.For = m.UUID()
	exp.Payload.URI = m.Payload.URI
	exp.Payload.Args = m.Payload.Args
	return exp
}

func logf(fn func(string, ...interface{}), f string, args ...interface{}) {
	if fn != nil {
		fn(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
