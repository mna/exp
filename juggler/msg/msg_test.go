package msg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
