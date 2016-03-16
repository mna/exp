package main

import (
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckRedisConfig(t *testing.T) {
	cases := []struct {
		in  string
		err bool
	}{
		{"", false},
		{`redis:
    MaxActive: 2
    pubsub:
        addr: :1234
`, true},
		{`redis:
    MaxActive: 2
    pubsub:
        addr: :1234
    caller:
        idle_timeout: 1s
`, true},
		{`redis:
    pubsub:
        addr: :1234
    caller:
        idle_timeout: 1s
`, true},
		{`redis:
    pubsub:
        addr: :1234
    caller:
        addr: :1235
        idle_timeout: 1s
`, false},
		{`redis:
    addr: :9876
    pubsub:
        addr: :1234
    caller:
        addr: :1235
        idle_timeout: 1s
`, true},
	}
	for i, c := range cases {
		conf, err := getConfigFromReader(strings.NewReader(c.in))
		require.NoError(t, err, "%d", i)
		err = checkRedisConfig(conf.Redis)
		assert.Equal(t, c.err, err != nil, "%d: %v", err)
		t.Logf("%d %v", i, err)
	}
}

func TestConfig(t *testing.T) {
	cases := []struct {
		in  string
		out *Config
	}{
		{"", getDefaultConfig()},
		{
			`
redis:
    addr: localhost:1234
`, &Config{
				Redis:        &Redis{Addr: "localhost:1234"},
				Server:       &Server{Addr: ":9000", Paths: []string{"/ws"}},
				CallerBroker: &CallerBroker{},
			},
		},
		{
			`
redis:
    pubsub:
        addr: :6666

    caller:
        addr: :6667
        max_active: 123
`, &Config{
				Redis: &Redis{
					Addr: ":6379",
					PubSub: &Redis{
						Addr: ":6666",
					},
					Caller: &Redis{
						Addr:      ":6667",
						MaxActive: 123,
					},
				},
				Server:       &Server{Addr: ":9000", Paths: []string{"/ws"}},
				CallerBroker: &CallerBroker{},
			},
		},
		{
			`
redis:
    addr: localhost:1234
    max_active: 34
    max_idle: 5
    idle_timeout: 1s

caller_broker:
    blocking_timeout: 2s
    call_cap: 987

server:
    addr: :9876

    paths:
    - /ws
    - /

    max_header_bytes: 23
    read_buffer_size: 4
    write_buffer_size: 5
    handshake_timeout: 1m

    whitelisted_origins:
    - http://localhost:4444

    read_limit: 6
    write_limit: 7
    read_timeout: 1h
    write_timeout: 2h
    acquire_write_lock_timeout: 3h

    allow_empty_subprotocol: true
`, &Config{
				Redis: &Redis{Addr: "localhost:1234", MaxActive: 34, MaxIdle: 5, IdleTimeout: time.Second},
				Server: &Server{Addr: ":9876", Paths: []string{"/ws", "/"}, MaxHeaderBytes: 23, ReadBufferSize: 4,
					WriteBufferSize: 5, HandshakeTimeout: time.Minute, WhitelistedOrigins: []string{"http://localhost:4444"},
					ReadLimit: 6, WriteLimit: 7, ReadTimeout: time.Hour, WriteTimeout: 2 * time.Hour,
					AcquireWriteLockTimeout: 3 * time.Hour, AllowEmptySubprotocol: true},
				CallerBroker: &CallerBroker{BlockingTimeout: 2 * time.Second, CallCap: 987},
			},
		},
	}

	for i, c := range cases {
		got, err := getConfigFromReader(strings.NewReader(c.in))
		require.NoError(t, err, "%d", i)
		if !assert.Equal(t, c.out, got, "%d", i) {
			spew.Dump(got)
		}
	}
}
