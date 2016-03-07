package msg

import (
	"encoding/json"
	"time"

	"github.com/pborman/uuid"
)

// CallPayload is the payload stored in the connector for a Call
// request.
type CallPayload struct {
	ConnUUID uuid.UUID       `json:"conn_uuid"`
	MsgUUID  uuid.UUID       `json:"msg_uuid"`
	Args     json.RawMessage `json:"args,omitempty"`

	// TTLAfterRead is the time-to-live remaining for the call request
	// once it has been extracted from the connector and just before it
	// is sent for processing to the callee.
	TTLAfterRead time.Duration `json:"-"`

	// ReadTimestamp is the timestamp in UTC of the call request once it
	// has been extracted from the connector and just before it is sent
	// for processing to the callee. It should be treated as informational,
	// as clocks may vary (sometimes wildly) between computers.
	ReadTimestamp time.Time `json:"-"`
}

// ResPayload is the payload stored in the connector for a result
// of a call request.
type ResPayload struct {
	ConnUUID uuid.UUID       `json:"conn_uuid"`
	MsgUUID  uuid.UUID       `json:"msg_uuid"`
	Args     json.RawMessage `json:"args,omitempty"`
}

// PubPayload is the payload to publish an event.
type PubPayload struct {
	MsgUUID uuid.UUID       `json:"msg_uuid"`
	Args    json.RawMessage `json:"args,omitempty"`
}

// EvntPayload is the payload of an event received by a subscriber.
type EvntPayload struct {
	MsgUUID uuid.UUID       `json:"msg_uuid"`
	Channel string          `json:"channel"`           // channel on which the event was sent
	Pattern string          `json:"pattern,omitempty"` // if received because of a pattern-based subscription
	Args    json.RawMessage `json:"args,omitempty"`
}
