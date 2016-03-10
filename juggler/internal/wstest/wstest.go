package wstest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StartServer starts a websocket server that calls fn with the created
// websocket connection for each request it receives. It sends true on the
// done channel when the connection is terminated. The server should
// be closed by the caller.
func StartServer(t *testing.T, done chan<- bool, fn func(*websocket.Conn)) *httptest.Server {
	upg := &websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upg.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() {
			wsConn.Close()
			done <- true
		}()

		fn(wsConn)
	}))
	return srv
}

// StartRecordingServer starts a websocket server that listens for
// messages from the websocket connections and writes them to
// w. It sends true on the done channel when the connection is
// terminated. Control messages are ignored. The server should
// be closed by the caller.
func StartRecordingServer(t *testing.T, done chan<- bool, w io.Writer) *httptest.Server {
	srv := StartServer(t, done, func(c *websocket.Conn) {
		for {
			_, r, err := c.NextReader()
			if err != nil {
				break
			}
			_, err = io.Copy(w, r)
			if !assert.NoError(t, err, "record message") {
				break
			}
		}
	})
	return srv
}

// Dial starts a new connection to urlStr and returns the created
// websocket connection. If urlStr uses an http: scheme, it is replaced
// by ws:. The connection should be closed by the caller.
func Dial(t *testing.T, urlStr string) *websocket.Conn {
	var d websocket.Dialer
	c, res, err := d.Dial(strings.Replace(urlStr, "http:", "ws:", 1), nil)
	require.NoError(t, err, "Dial")
	assert.Equal(t, 101, res.StatusCode, "status code")
	return c
}
