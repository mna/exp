package juggler

import (
	"errors"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"testing/quick"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

func TestChain(t *testing.T) {
	var b []byte

	genHandler := func(char byte) HandlerFunc {
		return HandlerFunc(func(ctx context.Context, c *Conn, m msg.Msg) {
			b = append(b, char)
		})
	}
	ch := Chain(genHandler('a'), genHandler('b'), genHandler('c'))
	ch.Handle(context.Background(), &Conn{}, &msg.OK{})

	assert.Equal(t, "abc", string(b))
}

type fakePubSubConn struct{}

func (f fakePubSubConn) Subscribe(channel string, pattern bool) error   { return nil }
func (f fakePubSubConn) Unsubscribe(channel string, pattern bool) error { return nil }
func (f fakePubSubConn) Events() <-chan *msg.EvntPayload                { return nil }
func (f fakePubSubConn) EventsErr() error                               { return nil }
func (f fakePubSubConn) Close() error                                   { return nil }

type fakeResultsConn struct{}

func (f fakeResultsConn) Results() <-chan *msg.ResPayload { return nil }
func (f fakeResultsConn) ResultsErr() error               { return nil }
func (f fakeResultsConn) Close() error                    { return nil }

type debugLog struct {
	t *testing.T
	n int64
}

func (d *debugLog) Printf(s string, args ...interface{}) {
	atomic.AddInt64(&d.n, 1)
	if testing.Verbose() {
		log.Printf(s, args...)
	}
}

func (d *debugLog) Calls() int {
	return int(atomic.LoadInt64(&d.n))
}

func TestPanicRecover(t *testing.T) {
	defer func() {
		require.Nil(t, recover(), "panic escaped the PanicRecover handler")
	}()

	panicer := HandlerFunc(func(ctx context.Context, c *Conn, m msg.Msg) {
		panic("a")
	})
	ph := PanicRecover(panicer, true, true)

	dbgl := &debugLog{t: t}
	srv := &Server{LogFunc: dbgl.Printf}
	conn := newConn(&websocket.Conn{}, srv)
	conn.psc, conn.resc = fakePubSubConn{}, fakeResultsConn{}
	ph.Handle(context.Background(), conn, &msg.OK{})

	err := conn.CloseErr
	if assert.NotNil(t, err, "connection has been closed") {
		assert.Equal(t, errors.New("a"), err, "error is as expected")
	}
	// with the stack, PanicRecover calls the log twice
	assert.Equal(t, 2, dbgl.Calls(), "log calls")
}

func TestLimitedWriter(t *testing.T) {
	// use int8/uint8 to keep size reasonable
	checker := func(limit int16, n uint8) bool {
		// create a limited writer with the specified limit
		w := limitWriter(ioutil.Discard, int64(limit))
		// create the payload for each write
		p := make([]byte, n)

		var cnt, tot int
		var err error
		for {
			cnt, err = w.Write(p)
			tot += cnt
			if err != nil {
				break
			}
		}

		// property 1: the total number of bytes written cannot be > limit
		// except if limit < minLimit (4096).
		if limit < minWriteLimit {
			limit = minWriteLimit
		}
		if tot > int(limit) {
			return false
		}
		// property 2: by writing repeatedly, it necessarily terminates with
		// an errWriteLimitExceeded
		return err == errWriteLimitExceeded
	}
	assert.NoError(t, quick.Check(checker, nil))
}
