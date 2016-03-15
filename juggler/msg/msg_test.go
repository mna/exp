package msg

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	call, err := NewCall("a", map[string]interface{}{"x": 3}, time.Second)
	require.NoError(t, err, "NewCall")
	pub, err := NewPub("d", map[string]interface{}{"y": "ok"})
	require.NoError(t, err, "NewPub")
	rp := &ResPayload{
		ConnUUID: uuid.NewRandom(),
		MsgUUID:  uuid.NewRandom(),
		URI:      "g",
		Args:     json.RawMessage("null"),
	}
	ep := &EvntPayload{
		MsgUUID: uuid.NewRandom(),
		Channel: "h",
		Pattern: "h*",
		Args:    json.RawMessage(`"string"`),
	}

	cases := []Msg{
		call,
		NewSub("b", false),
		NewUnsb("c", true),
		pub,
		NewErr(call, 500, io.EOF),
		NewOK(pub),
		NewRes(rp),
		NewEvnt(ep),
	}
	for i, m := range cases {
		b, err := json.Marshal(m)
		require.NoError(t, err, "Marshal %d", i)

		mm, err := Unmarshal(bytes.NewReader(b))
		require.NoError(t, err, "Unmarshal %d", i)

		// for ErrMsg, the Err is not marshaled, so zero it before the comparison
		if m.Type() == ErrMsg {
			m.(*Err).Payload.Err = nil
		}

		assert.True(t, reflect.DeepEqual(m, mm), "DeepEqual %d", i)

		_, err = UnmarshalRequest(bytes.NewReader(b))
		assert.Equal(t, m.Type().IsRead(), err == nil, "UnmarshalRequest for %d", i)

		_, err = UnmarshalResponse(bytes.NewReader(b))
		assert.Equal(t, m.Type().IsWrite(), err == nil, "UnmarshalResponse for %d", i)
	}
}

func TestNewErrFromOK(t *testing.T) {
	t.Parallel()

	pub, err := NewPub("d", map[string]interface{}{"y": "ok"})
	require.NoError(t, err, "NewPub")
	ok := NewOK(pub)
	e := NewErr(ok, 500, io.EOF)

	// should keep the "from" information of OK
	assert.Equal(t, e.Payload.For, ok.Payload.For, "For")
	assert.Equal(t, e.Payload.ForType, ok.Payload.ForType, "ForType")
	assert.Equal(t, e.Payload.URI, ok.Payload.URI, "URI")
	assert.Equal(t, e.Payload.Channel, ok.Payload.Channel, "Channel")
}
