package juggler

import "github.com/garyburd/redigo/redis"

// RedisPool defines the methods required for a Redis pool. The
// pool provides connections via Get and must be closed to release
// its resources.
type RedisPool interface {
	// Get returns a redis connection.
	Get() redis.Conn

	// Close releases the resources used by the pool.
	Close() error
}
