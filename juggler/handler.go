package juggler

// MsgHandler defines the method required to handle a send or receive
// of a Msg over a connection.
type MsgHandler interface {
	Handle(*Conn, Msg)
}

// MsgHandlerFunc is a function signature that implements the MsgHandler
// interface.
type MsgHandlerFunc func(*Conn, Msg)

// Handle implements MsgHandler for the MsgHandlerFunc by calling the
// function itself.
func (h MsgHandlerFunc) Handle(c *Conn, msg Msg) {
	h(c, msg)
}

// LogConn is a function compatible with the Server.ConnState field
// type that logs connections and disconnections to LogFunc.
func LogConn(c *Conn, state ConnState) {
	switch state {
	case Connected:
		LogFunc("%v: connected from %v with subprotocol %q", c.UUID, c.WSConn.RemoteAddr(), c.WSConn.Subprotocol())
	case Closing:
		LogFunc("%v: closing from %v with error %v", c.UUID, c.WSConn.RemoteAddr(), c.CloseErr)
	}
}

// LogMsg is a MsgHandlerFunc that logs messages received or sent on
// c to LogFunc.
func LogMsg(c *Conn, msg Msg) {
	if msg.IsRead() {
		LogFunc("%v: received message %v %s", c.UUID, msg.UUID(), msg.Type())
	} else if msg.IsWrite() {
		LogFunc("%v: sending message %v %s", c.UUID, msg.UUID(), msg.Type())
	}
}

// ProcessMsg implements the default message processing. For client messages,
// it calls the appropriate RPC, PUB-SUB or AUTH mechanisms. For server
// messages, it marshals the message and sends it to the client.
//
// When a custom ReadHandler and/or WriterHandler is set on the Server,
// it should at some point call ProcessMsg so the expected behaviour
// happens.
func ProcessMsg(c *Conn, msg Msg) {
	// TODO : default handling based on the type of msg
}
