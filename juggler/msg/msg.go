// Package msg defines the supported types of messages in the juggler
// protocol.
//
// The juggler.0 subprotocol defines the following messages for the client:
//
//     - CALL : to call an RPC function
//     - SUB  : to subscribe to a pub/sub channel
//     - UNSB : to unsubscrube from a pub/sub channel
//     - PUB  : to publish to a pub/sub channel
//
// And the following messages for the server:
//
//     - ERR  : failed CALL, SUB, UNSB or PUB
//     - OK   : successful CALL, SUB, UNSB or PUB - but no result yet
//     - RES  : the result of a CALL message
//     - EVNT : an event triggered on a channel that the client is subscribed to
//
// There's another message that is not sent by either end, but can be
// triggered by the client for itself:
//
//     - EXP  : expired CALL, meaning that no RES will be received for this call.
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
	"io"
	"time"

	"github.com/pborman/uuid"
)

// MessageType indicates the type of a message.
type MessageType int

// The list of supported message types.
const (
	startRead MessageType = iota
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

	startMeta
	ExpMsg
	endMeta
)

var lookupMessageType = []string{
	CallMsg: "CALL",
	PubMsg:  "PUB",
	SubMsg:  "SUB",
	UnsbMsg: "UNSB",
	ErrMsg:  "ERR",
	OKMsg:   "OK",
	ResMsg:  "RES",
	EvntMsg: "EVNT",
	ExpMsg:  "EXP",
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

// partialMsg is a message that decodes only the metadata, leaving
// the payload in raw JSON.
type partialMsg struct {
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

// NewCall creates a Call message using the provided arguments. The uri
// identifies the function to call. The args value is marshaled to JSON
// and used as the parameters to the call. If the result is not available
// before the timeout, it is dropped.
func NewCall(uri string, args interface{}, timeout time.Duration) (*Call, error) {
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

// NewSub creates a Sub message using the provided arguments. The
// channel indicates the pub-sub channel to subscribe to. It is
// treated as a pattern if pattern is true.
func NewSub(channel string, pattern bool) *Sub {
	sub := &Sub{
		Meta: newMeta(SubMsg),
	}
	sub.Payload.Channel = channel
	sub.Payload.Pattern = pattern
	return sub
}

// Unsb is an unsubscription message. It unsubscribes the caller from
// the Channel, which is treated as a pattern if Pattern is true. The
// pattern behaviour is the same as that of Redis.
type Unsb Sub

// NewUnsb creates an Unsb message using the provided arguments. The
// channel indicates the pub-sub channel to unsubscribe from. It is
// treated as a pattern if pattern is true.
func NewUnsb(channel string, pattern bool) *Unsb {
	un := &Unsb{
		Meta: newMeta(UnsbMsg),
	}
	un.Payload.Channel = channel
	un.Payload.Pattern = pattern
	return un
}

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

// NewPub creates a Pub message using the provided arguments. The channel
// identifies the channel on which this event is published. The args value
// is marshaled to JSON and used as the payload of the event.
func NewPub(channel string, args interface{}) (*Pub, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	p := &Pub{
		Meta: newMeta(PubMsg),
	}
	p.Payload.Channel = channel
	p.Payload.Args = json.RawMessage(b)
	return p, nil
}

// Err is an error message. It indicates the source message that
// failed to be delivered in the For (and ForType) fields. An Err
// is sent only when something failed to execute properly - notably,
// it is not sent if the result of a call was processed by the callee
// but resulted in an error. This would be returned by a Res message.
type Err struct {
	Meta    `json:"meta"`
	Payload struct {
		For     uuid.UUID   `json:"for"`
		ForType MessageType `json:"for_type"`
		URI     string      `json:"uri,omitempty"`     // when in response to a CALL
		Channel string      `json:"channel,omitempty"` // when in response to a PUB, SUB or UNSB
		Code    int         `json:"code"`
		Message string      `json:"message"` // defaults to Err.Error()
		Err     error       `json:"-"`       // useful in the handler to have access to the source error
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
		err.Payload.URI = from.Payload.URI
		err.Payload.Channel = from.Payload.Channel
	case *Err:
		err.Payload.For = from.Payload.For
		err.Payload.ForType = from.Payload.ForType
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
		For     uuid.UUID   `json:"for"`
		ForType MessageType `json:"for_type"`
		URI     string      `json:"uri,omitempty"`     // when in response to a CALL
		Channel string      `json:"channel,omitempty"` // when in response to a PUB, SUB or UNSB
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
		For  uuid.UUID       `json:"for"`           // no ForType, because always CALL
		URI  string          `json:"uri,omitempty"` // URI of the CALL
		Args json.RawMessage `json:"args"`
	} `json:"payload"`
}

// NewRes creates a new Res message corresponding to a call result.
func NewRes(pld *ResPayload) *Res {
	res := &Res{
		Meta: newMeta(ResMsg),
	}
	res.Payload.For = pld.MsgUUID
	res.Payload.URI = pld.URI
	res.Payload.Args = pld.Args
	return res
}

// Evnt is a published event. It is sent to all subscribers of the
// Channel.
type Evnt struct {
	Meta    `json:"meta"`
	Payload struct {
		For     uuid.UUID       `json:"for"` // no ForType, because always PUB
		Channel string          `json:"channel,omitempty"`
		Pattern string          `json:"pattern,omitempty"` // if triggered because of a pattern-based subscription
		Args    json.RawMessage `json:"args"`
	} `json:"payload"`
}

// NewEvnt creates a new Evnt message corresponding to an event that
// occurred on a subscribed channel.
func NewEvnt(pld *EvntPayload) *Evnt {
	ev := &Evnt{
		Meta: newMeta(EvntMsg),
	}
	ev.Payload.Channel = pld.Channel
	ev.Payload.Pattern = pld.Pattern
	ev.Payload.For = pld.MsgUUID
	ev.Payload.Args = pld.Args
	return ev
}

// Exp is an expired call message. It is never sent over the network, but
// it can be raised by a client, for itself, when the timeout for a call
// result has expired. As such, the ExpMsg message type returns false for
// both IsRead and IsWrite.
type Exp struct {
	Meta    `json:"meta"`
	Payload struct {
		For  uuid.UUID       `json:"for"`           // no ForType, because always CALL
		URI  string          `json:"uri,omitempty"` // URI of the CALL
		Args json.RawMessage `json:"args"`
	} `json:"payload"`
}

// NewExp creates a new expired message for the provided call message.
func NewExp(m *Call) *Exp {
	exp := &Exp{
		Meta: newMeta(ExpMsg),
	}
	exp.Payload.For = m.UUID()
	exp.Payload.URI = m.Payload.URI
	exp.Payload.Args = m.Payload.Args
	return exp
}

// UnmarshalRequest unmarshals a JSON-encoded message from r into the
// correct concrete message type. It returns an error if the message
// type is invalid for a request (client -> server).
func UnmarshalRequest(r io.Reader) (Msg, error) {
	return unmarshalIf(r, CallMsg, SubMsg, UnsbMsg, PubMsg)
}

// UnmarshalResponse unmarshals a JSON-encoded message from r into the
// correct concrete message type. It returns an error if the message
// type is invalid for a response (client <- server).
func UnmarshalResponse(r io.Reader) (Msg, error) {
	return unmarshalIf(r, ErrMsg, OKMsg, EvntMsg, ResMsg)
}

// Unmarshal unmarshals a JSON-encoded message from r into the correct
// concrete message type.
func Unmarshal(r io.Reader) (Msg, error) {
	return unmarshalIf(r)
}

func isIn(list []MessageType, v MessageType) bool {
	for _, vv := range list {
		if v == vv {
			return true
		}
	}
	return false
}

func unmarshalIf(r io.Reader, allowed ...MessageType) (Msg, error) {
	var pm partialMsg
	if err := json.NewDecoder(r).Decode(&pm); err != nil {
		return nil, fmt.Errorf("invalid JSON message: %v", err)
	}

	if len(allowed) > 0 && !isIn(allowed, pm.Meta.T) {
		return nil, fmt.Errorf("invalid message %s for this peer", pm.Meta.T)
	}

	genericUnmarshal := func(v interface{}, metaDst *Meta) error {
		var b []byte
		b = append(b, `{"payload":`...)
		b = append(b, pm.Payload...)
		b = append(b, '}')
		if err := json.Unmarshal(b, v); err != nil {
			return fmt.Errorf("invalid %s message: %v", pm.Meta.T, err)
		}
		*metaDst = pm.Meta
		return nil
	}

	var m Msg
	switch pm.Meta.T {
	case CallMsg:
		var call Call
		if err := genericUnmarshal(&call, &call.Meta); err != nil {
			return nil, err
		}
		m = &call

	case SubMsg:
		var sub Sub
		if err := genericUnmarshal(&sub, &sub.Meta); err != nil {
			return nil, err
		}
		m = &sub

	case UnsbMsg:
		var uns Unsb
		if err := genericUnmarshal(&uns, &uns.Meta); err != nil {
			return nil, err
		}
		m = &uns

	case PubMsg:
		var pub Pub
		if err := genericUnmarshal(&pub, &pub.Meta); err != nil {
			return nil, err
		}
		m = &pub

	case ErrMsg:
		var e Err
		if err := genericUnmarshal(&e, &e.Meta); err != nil {
			return nil, err
		}
		m = &e

	case OKMsg:
		var ok OK
		if err := genericUnmarshal(&ok, &ok.Meta); err != nil {
			return nil, err
		}
		m = &ok

	case ResMsg:
		var res Res
		if err := genericUnmarshal(&res, &res.Meta); err != nil {
			return nil, err
		}
		m = &res

	case EvntMsg:
		var ev Evnt
		if err := genericUnmarshal(&ev, &ev.Meta); err != nil {
			return nil, err
		}
		m = &ev

	default:
		return nil, fmt.Errorf("unknown message %s", pm.Meta.T)
	}

	return m, nil
}
