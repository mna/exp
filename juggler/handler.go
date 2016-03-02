package juggler

type MsgHandlerFunc func(*Conn, Msg)

func (h MsgHandlerFunc) Handle(c *Conn, msg Msg) {
	h(c, msg)
}

func LogConn(c *Conn) {
	st, err := c.State()
	switch st {
	case Connected:
		LogFunc("%v: connected from %v with subprotocol %q", c.UUID, c.WSConn.RemoteAddr(), c.WSConn.Subprotocol())
	case Closing:
		LogFunc("%v: closing from %v with error %v", c.UUID, c.WSConn.RemoteAddr(), err)
	}
}

func LogMsg(c *Conn, msg Msg) {
	LogFunc("%v: received message %v %s", c.UUID, msg.UUID(), msg.Type())
}

func ProcessMsg(c *Conn, msg Msg) {
	// TODO : default handling based on the type of msg
}
