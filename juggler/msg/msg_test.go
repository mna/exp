package msg

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
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

	exp := NewExp(call)
	b, err := json.Marshal(exp)
	require.NoError(t, err, "Marshal EXP")
	_, err = Unmarshal(bytes.NewReader(b))
	assert.Error(t, err, "Unmarshal not allowed for EXP")
	assert.Contains(t, err.Error(), "unknown message EXP", "Unmarshal for EXP returns expected error")
}

func TestUnmarshalOK(t *testing.T) {
	raw := `{"meta":{"type":8,"uuid":"845ca2a4-7d1e-44eb-b1b1-304bb222a2bf"},"payload":{"for":"9dc9b548-3a7b-40d0-82b0-ccf7e9828e78","for_type":1}}`
	m, err := Unmarshal(strings.NewReader(raw))
	assert.NoError(t, err, "Unmarshal")
	if assert.IsType(t, &OK{}, m, "Is *OK") {
		ok := m.(*OK)
		assert.Equal(t, OKMsg, ok.Type(), "Type")
		assert.Equal(t, "845ca2a4-7d1e-44eb-b1b1-304bb222a2bf", ok.UUID().String(), "UUID")
		assert.Equal(t, "9dc9b548-3a7b-40d0-82b0-ccf7e9828e78", ok.Payload.For.String(), "For")
		assert.Equal(t, CallMsg, ok.Payload.ForType, "ForType")
	}
}
