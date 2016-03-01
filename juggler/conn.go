package juggler

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

type ConnState int

const (
	Unknown ConnState = iota
	Connected
	Closing
)

type ConnHandler interface {
	Handle(*Conn)
}

type Conn struct {
	UUID uuid.UUID

	// TODO : hide/show only as needed
	WSConn *websocket.Conn
	// TODO : some connection state (authenticated, etc.)

	mu       sync.RWMutex
	state    ConnState
	closeErr error
}

func (c *Conn) State() (state ConnState, closeErr error) {
	c.mu.RLock()
	st, err := c.state, c.closeErr
	c.mu.RUnlock()
	return st, err
}

func (c *Conn) setState(state ConnState, err error) {
	c.mu.Lock()
	c.state = state
	c.closeErr = err
	c.mu.Unlock()
}
