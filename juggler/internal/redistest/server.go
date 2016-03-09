package redistest

import (
	"io"
	"net"
	"os/exec"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/require"
)

// StartServer starts a redis-server instance on a free port.
// It returns the started *exec.Cmd and the port used. The caller
// should make sure to stop the command. If the redis-server
// command is not found in the PATH, the test is skipped.
func StartServer(t *testing.T, w io.Writer) (*exec.Cmd, string) {
	if _, err := exec.LookPath("redis-server"); err != nil {
		t.Skip("redis-server not found in $PATH")
	}

	port := getFreePort(t)
	c := exec.Command("redis-server", "--port", port)
	if w != nil {
		c.Stderr = w
		c.Stdout = w
	}
	require.NoError(t, c.Start(), "start redis-server")

	// wait a bit for the server to start listening... better way?
	time.Sleep(500 * time.Millisecond)
	t.Logf("redis-server started on port %s", port)
	return c, port
}

func getFreePort(t *testing.T) string {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err, "listen on port 0")
	defer l.Close()
	_, p, err := net.SplitHostPort(l.Addr().String())
	require.NoError(t, err, "parse host and port")
	return p
}

// NewPool creates a redis pool to return connections on the specified
// addr.
func NewPool(t *testing.T, addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     2,
		MaxActive:   10,
		IdleTimeout: time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
