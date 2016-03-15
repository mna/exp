package main

import (
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
