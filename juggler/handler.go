package juggler

type ConnHandlerFunc func(*Conn)

func (h ConnHandlerFunc) Handle(c *Conn) {
	h(c)
}

type MsgHandlerFunc func(*Conn, Msg, Direction)

func (h MsgHandlerFunc) Handle(c *Conn, msg Msg, dir Direction) {
	h(c, msg, dir)
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

func ProcessMsg(c *Conn, msg Msg, dir Direction) {
	// TODO : default handling based on the type of msg
}
