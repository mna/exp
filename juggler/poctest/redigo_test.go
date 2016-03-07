package poctest

import (
	"io"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClosedConnErr(t *testing.T) {
	cmd, port := redistest.StartRedisServer(t, nil)
	defer cmd.Process.Kill()

	conn, err := redis.Dial("tcp", ":"+port)
	require.NoError(t, err, "dial")
	defer conn.Close()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cmd.Process.Kill()
	}()

	before := time.Now()
	_, err = conn.Do("BRPOP", "some-key", 2)
	after := time.Now()
	assert.Error(t, err, "redis-server stopped during BRPOP")

	// A closed redis connection returns io.EOF
	assert.Equal(t, io.EOF, err, "error is io.EOF")

	assert.WithinDuration(t, before.Add(100*time.Millisecond), after, 100*time.Millisecond, "BRPOP returned as soon as the server stopped")
	t.Logf("%T %[1]v", err)
}
