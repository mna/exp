package juggler

import (
	"encoding/json"
	"fmt"

	"github.com/pborman/uuid"
)

/*
The juggler.1 subprotocol defines the following messages for the client:

- AUTH : for authentication
- CALL : to call an RPC function
- SUB  : to subscribe to a pub/sub channel
- PUB  : to publish to a pub/sub channel

And the following messages for the server:

- ERR  : in response to an invalid AUTH, CALL, SUB or PUB
- OK   : successful AUTH, SUB or PUB
- RES  : the result of a successful CALL message
- EVNT : an event triggered on a channel that the client is subscribed to

Closing the communication is done via the standard websocket close
process.

All messages must be of type websocket.TextMessage. Failing to properly
speak the protocol terminates the connection without notice from the
peer. That includes sending binary messages and sending unknown message
types.
*/

type MessageType int

const (
	startRead MessageType = iota
	AuthMsg
	CallMsg
	PubMsg
	SubMsg
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
	ErrMsg:  "ERR",
	OKMsg:   "OK",
	ResMsg:  "RES",
	EvntMsg: "EVNT",
}

func (mt MessageType) String() string {
	if mt > 0 && int(mt) < len(lookupMessageType) {
		return lookupMessageType[mt]
	}
	return fmt.Sprintf("<unknown: %d>", mt)
}

type Msg interface {
	Type() MessageType
	UUID() uuid.UUID
	IsRead() bool
	IsWrite() bool
}

type meta struct {
	T MessageType `json:"type"`
	U uuid.UUID   `json:"uuid"`
}

type partialMsg struct {
	Meta    meta            `json:"meta"`
	Payload json.RawMessage `json:"payload"`
}

func (m meta) Type() MessageType {
	return m.T
}

func (m meta) UUID() uuid.UUID {
	return m.U
}

func (m meta) IsRead() bool {
	return startRead < m.T && m.T < endRead
}

func (m meta) IsWrite() bool {
	return startWrite < m.T && m.T < endWrite
}

type Auth struct {
	meta
	AuthType string `json:"auth_type"`
	Token    string `json:"token,omitempty"`
	ID       string `json:"id,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

type Call struct {
	meta
	URI  string          `json:"uri"`
	Args json.RawMessage `json:"args"`
}

type Sub struct {
	meta
	Channel string `json:"channel"`
	Pattern bool   `json:"pattern"`
}

type Pub struct {
	meta
	Channel string          `json:"channel"`
	Args    json.RawMessage `json:"args"`
}

type Err struct {
	meta
	For      uuid.UUID `json:"for"`
	ForType  int       `json:"for_type"`
	AuthType string    `json:"auth_type,omitempty"` // when in response to an AUTH
	URI      string    `json:"uri,omitempty"`       // when in response to a CALL
	Channel  string    `json:"channel,omitempty"`   // when in response to a PUB or SUB
	Code     int       `json:"code"`
	Message  string    `json:"message"`
}

type OK struct {
	meta
	For      uuid.UUID `json:"for"`
	ForType  int       `json:"for_type"`            // AUTH, PUB or SUB
	AuthType string    `json:"auth_type,omitempty"` // when in response to an AUTH
	Channel  string    `json:"channel,omitempty"`   // when in response to PUB or SUB
	ID       string    `json:"id,omitempty"`        // ID of authenticated user, when in response to an AUTH
}

type Res struct {
	meta
	For  uuid.UUID       `json:"for"`           // no ForType, because always CALL
	URI  string          `json:"uri,omitempty"` // URI of the CALL
	Args json.RawMessage `json:"args"`
}

type Evnt struct {
	meta
	For     uuid.UUID       `json:"for"` // no ForType, because always PUB
	Channel string          `json:"channel"`
	Pattern string          `json:"pattern,omitempty"` // if triggered because of a pattern-based subscription
	Args    json.RawMessage `json:"args"`
}
