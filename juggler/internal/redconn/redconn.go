package redconn

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// Pool defines the methods required for a redis pool that provides
// a method to get a connection and to release the pool's resources.
type Pool interface {
	Get() redis.Conn
	Close() error
}

// Connector is a redis connector that provides the methods to
// interact with Redis using the juggler protocol.
type Connector struct {
	Pool            Pool
	BlockingTimeout time.Duration
}

func (c *Connector) Call() error {

}

func (c *Connector) ProcessCalls() {

}

func (c *Connector) ProcessResults() {

}

func (c *Connector) Publish() error {

}

func (c *Connector) Subscribe() error {

}

func (c *Connector) Unsubscribe() error {

}

func (c *Connector) ProcessEvents() {

}
