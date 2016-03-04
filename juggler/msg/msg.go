// Package msg defines the supported types of messages in the juggler
// protocol.
//
// The juggler.0 subprotocol defines the following messages for the client:
//
//     - AUTH : for authentication (like a CALL, but with additional
//              structure, so the result knows if it succeeded and
//              for how long the auth is valid?).
//     - CALL : to call an RPC function
//     - SUB  : to subscribe to a pub/sub channel
//     - UNSB : to unsubscrube from a pub/sub channel
//     - PUB  : to publish to a pub/sub channel
//
// And the following messages for the server:
//
//     - ERR  : failed AUTH, CALL, SUB, UNSB or PUB
//     - OK   : successful AUTH, CALL, SUB, UNSB or PUB - but no result yet
//     - RES  : the result of an AUTH or CALL message
//     - EVNT : an event triggered on a channel that the client is subscribed to
//
// Closing the communication is done via the standard websocket close
// process.
//
// All messages must be of type websocket.TextMessage. Failing to properly
// speak the protocol terminates the connection without notice from the
// peer. That includes sending binary messages and sending unknown (or
// invalid for the peer) message types.
//
package msg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

// MessageType indicates the type of a message.
type MessageType int

// The list of supported message types.
const (
	startRead MessageType = iota
	AuthMsg
	CallMsg
	PubMsg
	SubMsg
	UnsbMsg
	endRead

	startWrite
	ErrMsg
	OKMsg
	ResMsg
	EvntMsg
	endWrite
)

var lookupMessageType = []string{
	AuthMsg: "AUTH",
	CallMsg: "CALL",
	PubMsg:  "PUB",
	SubMsg:  "SUB",
	UnsbMsg: "UNSB",
	ErrMsg:  "ERR",
	OKMsg:   "OK",
	ResMsg:  "RES",
	EvntMsg: "EVNT",
}

// String returns the human-readable representation of message types.
func (mt MessageType) String() string {
	if mt >= 0 && int(mt) < len(lookupMessageType) {
		if s := lookupMessageType[mt]; s != "" {
			return lookupMessageType[mt]
		}
	}
	return fmt.Sprintf("<unknown: %d>", mt)
}

// IsRead returns true if the message type is a "read" from the
// point of view of the server (that is, if this is a message
// that was sent by a client).
func (mt MessageType) IsRead() bool {
	return startRead < mt && mt < endRead
}

// IsWrite returns true if the message type is a "write" from the
// point of view of the server (that is, if this is a message
// that is being sent by the server).
func (mt MessageType) IsWrite() bool {
	return startWrite < mt && mt < endWrite
}

// Msg defines the common methods implemented by all messages.
type Msg interface {
	// Type returns the message type.
	Type() MessageType

	// UUID is the unique identifier of the message.
	UUID() uuid.UUID
}

// Meta contains the metadata for a message.
type Meta struct {
	T MessageType `json:"type"`
	U uuid.UUID   `json:"uuid"`
}

func newMeta(t MessageType) Meta {
	return Meta{T: t, U: uuid.NewRandom()}
}

// PartialMsg is a message that decodes only the metadata, leaving
// the payload in raw JSON. Primarily used by the server to
// decode the minimal part required to process a message.
type PartialMsg struct {
	Meta    Meta            `json:"meta"`
	Payload json.RawMessage `json:"payload"`
}

// Type returns the message type.
func (m Meta) Type() MessageType {
	return m.T
}

// UUID returns the message's unique identifier.
func (m Meta) UUID() uuid.UUID {
	return m.U
}

// Auth is an authentication message.
type Auth struct {
	Meta    `json:"meta"`
	Payload struct {
		AuthType string        `json:"auth_type"`
		Token    string        `json:"token,omitempty"`
		ID       string        `json:"id,omitempty"`
		Secret   string        `json:"secret,omitempty"`
		Timeout  time.Duration `json:"timeout"`
	} `json:"payload"`
}

// Call is a message that triggers an RPC call to a callee
// listening on the specified URI. The Args opaque field
// is transferred as-is to the callee. If the result is not
// available and sent back to the caller before the specified
// timeout, it is dropped.
type Call struct {
	Meta    `json:"meta"`
	Payload struct {
		URI     string          `json:"uri"`
		Timeout time.Duration   `json:"timeout"`
		Args    json.RawMessage `json:"args"`
	} `json:"payload"`
}

// NewCall creates a Call message using the provided arguments.
func NewCall(uri string, timeout time.Duration, args interface{}) (*Call, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	c := &Call{
		Meta: newMeta(CallMsg),
	}
	c.Payload.URI = uri
	c.Payload.Timeout = timeout
	c.Payload.Args = json.RawMessage(b)
	return c, nil
}

// Sub is a subscription message. It subscribes the caller to the
// Channel, which is treated as a pattern if Pattern is true. The
// pattern behaviour is the same as that of Redis.
type Sub struct {
	Meta    `json:"meta"`
	Payload struct {
		Channel string `json:"channel"`
		Pattern bool   `json:"pattern"`
	} `json:"payload"`
}

// Unsb is an unsubscription message. It unsubscribes the caller from
// the Channel, which is treated as a pattern if Pattern is true. The
// pattern behaviour is the same as that of Redis.
type Unsb Sub

// Pub is a publish message. It publishes an event on the specified
// Channel. The Args opaque field is transferred as-is to subscribers
// of that channel.
type Pub struct {
	Meta    `json:"meta"`
	Payload struct {
		Channel string          `json:"channel"`
		Args    json.RawMessage `json:"args"`
	} `json:"payload"`
}

// Err is an error message. It indicates the source message that
// failed to be delivered in the For (and ForType) fields. An Err
// is sent only when something failed to execute properly - notably,
// it is not sent if the result of a call was processed by the callee
// but resulted in an error. This would be returned by a Res message.
type Err struct {
	Meta    `json:"meta"`
	Payload struct {
		For      uuid.UUID   `json:"for"`
		ForType  MessageType `json:"for_type"`
		AuthType string      `json:"auth_type,omitempty"` // when in response to an AUTH
		URI      string      `json:"uri,omitempty"`       // when in response to a CALL
		Channel  string      `json:"channel,omitempty"`   // when in response to a PUB, SUB or UNSB
		Code     int         `json:"code"`
		Message  string      `json:"message"` // defaults to Err.Error()
		Err      error       `json:"-"`       // useful in the handler to have access to the source error
	} `json:"payload"`
}

// NewErr creates a new Err message to notify a failure to process
// the from message.
func NewErr(from Msg, code int, e error) *Err {
	err := &Err{
		Meta: newMeta(ErrMsg),
	}
	err.Payload.For = from.UUID()
	err.Payload.ForType = from.Type()
	err.Payload.Code = code
	err.Payload.Err = e
	err.Payload.Message = e.Error()

	switch from := from.(type) {
	case *Auth:
		err.Payload.AuthType = from.Payload.AuthType
	case *Call:
		err.Payload.URI = from.Payload.URI
	case *Pub:
		err.Payload.Channel = from.Payload.Channel
	case *Sub:
		err.Payload.Channel = from.Payload.Channel
	case *Unsb:
		err.Payload.Channel = from.Payload.Channel

		// other cases can happen e.g. if the message is too large
		// instead of sending the "from" info from the never-sent
		// OK, Err, Evnt or Res message, send back the origin "from"
		// information.
	case *OK:
		err.Payload.For = from.Payload.For
		err.Payload.ForType = from.Payload.ForType
		err.Payload.AuthType = from.Payload.AuthType
		err.Payload.URI = from.Payload.URI
		err.Payload.Channel = from.Payload.Channel
	case *Err:
		err.Payload.For = from.Payload.For
		err.Payload.ForType = from.Payload.ForType
		err.Payload.AuthType = from.Payload.AuthType
		err.Payload.URI = from.Payload.URI
		err.Payload.Channel = from.Payload.Channel
	case *Evnt:
		err.Payload.For = from.Payload.For
		err.Payload.ForType = PubMsg
		err.Payload.Channel = from.Payload.Channel
	case *Res:
		err.Payload.For = from.Payload.For
		err.Payload.ForType = CallMsg
		err.Payload.URI = from.Payload.URI
	}
	return err
}

// OK is a success message, when the request of the caller was successfully
// registered. It doesn't mean that e.g. a CALL has succeeded - only that
// the CALL was properly registered for a callee to process, eventually.
type OK struct {
	Meta    `json:"meta"`
	Payload struct {
		For      uuid.UUID   `json:"for"`
		ForType  MessageType `json:"for_type"`
		AuthType string      `json:"auth_type,omitempty"` // when in response to an AUTH
		URI      string      `json:"uri,omitempty"`       // when in response to a CALL
		Channel  string      `json:"channel,omitempty"`   // when in response to a PUB, SUB or UNSB
	} `json:"payload"`
}

// NewOK creates a new OK message to notify the successful execution
// of the from message.
func NewOK(from Msg) *OK {
	ok := &OK{
		Meta: newMeta(OKMsg),
	}
	ok.Payload.For = from.UUID()
	ok.Payload.ForType = from.Type()

	switch from := from.(type) {
	case *Auth:
		ok.Payload.AuthType = from.Payload.AuthType
	case *Call:
		ok.Payload.URI = from.Payload.URI
	case *Pub:
		ok.Payload.Channel = from.Payload.Channel
	case *Sub:
		ok.Payload.Channel = from.Payload.Channel
	case *Unsb:
		ok.Payload.Channel = from.Payload.Channel
	}
	return ok
}

// Res is a result message. It returns the result of the invocation
// of a Call message.
type Res struct {
	Meta    `json:"meta"`
	Payload struct {
		For  uuid.UUID       `json:"for"`           // TODO : no ForType, because always CALL, or AUTH too?
		URI  string          `json:"uri,omitempty"` // URI of the CALL
		Args json.RawMessage `json:"args"`
	} `json:"payload"`
}

// Evnt is a published event. It is sent to all subscribers of the
// Channel.
type Evnt struct {
	Meta    `json:"meta"`
	Payload struct {
		For     uuid.UUID       `json:"for"` // no ForType, because always PUB
		Channel string          `json:"channel"`
		Pattern string          `json:"pattern,omitempty"` // if triggered because of a pattern-based subscription
		Args    json.RawMessage `json:"args"`
	} `json:"payload"`
}
