package poctest

import (
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler/internal/redistest"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClosedConnErr(t *testing.T) {
	const stopAfter = 100 * time.Millisecond

	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	ch := make(chan struct{})
	ch2 := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go runRedisCmdExpectClose(t, wg, ch, stopAfter, port, "BRPOP", "some-key", 2)
	go runRedisCmdExpectClose(t, wg, ch, stopAfter, port, "LISTEN")
	go runRedisCmdExpectClose(t, wg, ch2, 0, port, "GET", "some-key") // runs AFTER server is closed

	go func() {
		close(ch)
		time.Sleep(stopAfter)
		cmd.Process.Kill()
		close(ch2)
	}()

	wg.Wait()
}

func runRedisCmdExpectClose(t *testing.T, wg *sync.WaitGroup, start <-chan struct{}, expectDelay time.Duration, port, cmd string, args ...interface{}) {
	defer wg.Done()

	conn, err := redis.Dial("tcp", ":"+port)
	require.NoError(t, err, "dial for %s", cmd)
	defer conn.Close()

	// wait to run the command
	<-start

	before := time.Now()

	if cmd == "LISTEN" {
		// special case, listen for pub-sub events
		psc := redis.PubSubConn{Conn: conn}
		got := psc.Receive()
		require.IsType(t, errors.New(""), got, "receive pub-sub returns an error")
		err = got.(error)
	} else {
		_, err = conn.Do(cmd, args...)
	}

	after := time.Now()
	assert.Error(t, err, "stopped redis-server caused error")

	// A closed redis connection returns io.EOF
	assert.Equal(t, io.EOF, err, "error is io.EOF")

	assert.WithinDuration(t, before.Add(expectDelay), after, 100*time.Millisecond, "returned as soon as the server stopped")
	t.Logf("%T %[1]v", err)
}

func TestCloseConnWhileBlocked(t *testing.T) {
	cmd, port := redistest.StartServer(t, nil)
	defer cmd.Process.Kill()

	conn, err := redis.Dial("tcp", ":"+port)
	require.NoError(t, err, "dial for %s", cmd)
	defer conn.Close()

	ch := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		ch <- struct{}{}
		_, err := conn.Do("BRPOP", "key", 5)
		assert.Contains(t, err.Error(), "use of closed network connection", "blocked conn returned with error")
		t.Logf("%T %[1]v", err)
	}()

	<-ch
	time.Sleep(100 * time.Millisecond)

	before := time.Now()
	err = conn.Close()
	require.NoError(t, err, "close without error")
	wg.Wait()
	after := time.Now()
	assert.WithinDuration(t, before, after, 100*time.Millisecond, "returned as soon as the connection closed")
}
