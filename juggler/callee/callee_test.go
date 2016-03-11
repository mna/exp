package callee

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCalleeBroker struct {
	cps []*msg.CallPayload
	err error
	rps []*msg.ResPayload
}

func (b *mockCalleeBroker) Result(rp *msg.ResPayload, timeout time.Duration) error {
	b.rps = append(b.rps, rp)
	return nil
}

func (b *mockCalleeBroker) Calls(uris ...string) (broker.CallsConn, error) {
	return &mockCallsConn{cps: b.cps, err: b.err}, nil
}

type mockCallsConn struct {
	cps []*msg.CallPayload
	err error
}

func (c *mockCallsConn) Calls() <-chan *msg.CallPayload {
	ch := make(chan *msg.CallPayload)
	go func() {
		for _, cp := range c.cps {
			ch <- cp
		}
		close(ch)
	}()
	return ch
}

func (c *mockCallsConn) CallsErr() error { return c.err }
func (c *mockCallsConn) Close() error    { return nil }

func okThunk(cp *msg.CallPayload) (interface{}, error) {
	time.Sleep(time.Millisecond)
	return "ok", nil
}

func errThunk(cp *msg.CallPayload) (interface{}, error) {
	time.Sleep(time.Millisecond)
	return nil, io.ErrUnexpectedEOF
}

func TestCallee(t *testing.T) {
	cuid := uuid.NewRandom()
	brk := &mockCalleeBroker{
		cps: []*msg.CallPayload{
			{ConnUUID: cuid, MsgUUID: uuid.NewRandom(), URI: "ok", TTLAfterRead: time.Second},
			{ConnUUID: cuid, MsgUUID: uuid.NewRandom(), URI: "err", TTLAfterRead: time.Second},
			{ConnUUID: cuid, MsgUUID: uuid.NewRandom(), URI: "ok", TTLAfterRead: time.Millisecond}, // result will be dropped
			{ConnUUID: cuid, MsgUUID: uuid.NewRandom(), URI: "err", TTLAfterRead: time.Second},
		},
		err: io.EOF,
	}

	var er errResult
	er.Error.Message = io.ErrUnexpectedEOF.Error()
	b, err := json.Marshal(er)
	require.NoError(t, err, "Marshal errResult")

	exp := []*msg.ResPayload{
		{ConnUUID: cuid, MsgUUID: brk.cps[0].MsgUUID, URI: "ok", Args: json.RawMessage(`"ok"`)},
		{ConnUUID: cuid, MsgUUID: brk.cps[1].MsgUUID, URI: "err", Args: b},
		{ConnUUID: cuid, MsgUUID: brk.cps[3].MsgUUID, URI: "err", Args: b},
	}

	cle := &Callee{Broker: brk, LogFunc: juggler.DiscardLog}
	err = cle.Listen(map[string]Thunk{
		"ok":  okThunk,
		"err": errThunk,
	})

	assert.Equal(t, io.EOF, err, "Listen returns expected error")
	assert.Equal(t, exp, brk.rps, "got expected results")
}
