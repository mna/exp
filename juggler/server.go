package juggler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var LogFunc = log.Printf

var Subprotocols = []string{
	"juggler.1",
}

type Server struct {
	ReadLimit    int64
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	ConnHandler  ConnHandler
	ReadHandler  MsgHandler
	WriteHandler MsgHandler
}

func Upgrade(upgrader *websocket.Upgrader, srv *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// upgrade the HTTP connection to the websocket protocol
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer wsConn.Close()
		if wsConn.Subprotocol() == "" {
			LogFunc("juggler: no supported subprotocol, closing connection")
			return
		}

		// configure the websocket connection
		wsConn.SetReadLimit(srv.ReadLimit)
		c := newConn(wsConn)
		defer func() {
			if srv.ConnHandler != nil {
				srv.ConnHandler.Handle(c)
			}
		}()

		// start lifecycle of the connection
		if srv.ConnHandler != nil {
			srv.ConnHandler.Handle(c)
		}

		if err := srv.read(c); err != nil {
			c.setState(Closing, err)
			LogFunc("juggler: read failed: %v; closing connection", err)
			return
		}

		/*
			if err := c.WSConn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
				c.setState(Closing, err)
				LogFunc("juggler: WriteMessage failed: %v; closing connection", err)
				return
			}
			for {
				c.WSConn.SetWriteDeadline(time.Time{})

				mt, r, err := c.WSConn.NextReader()
				if err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: NextReader failed: %v; closing connection", err)
					return
				}
				c.WSConn.SetReadDeadline(time.Now().Add(srv.ReadTimeout))

				w, err := c.WSConn.NextWriter(mt)
				if err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: NextWriter failed: %v; closing connection", err)
					return
				}
				c.WSConn.SetWriteDeadline(time.Now().Add(srv.WriteTimeout))

				if _, err := io.Copy(w, r); err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: Copy failed: %v; closing connection", err)
					return
				}
				if err := w.Close(); err != nil {
					c.setState(Closing, err)
					LogFunc("juggler: Close failed: %v; closing connection", err)
					return
				}
			}
		*/
	})
}

func (s *Server) read(c *Conn) error {
	for {
		c.WSConn.SetReadDeadline(time.Time{})

		mt, r, err := c.WSConn.NextReader()
		if err != nil {
			return err
		}
		if mt != websocket.TextMessage {
			return fmt.Errorf("invalid websocket message type: %d", mt)
		}
		c.WSConn.SetReadDeadline(time.Now().Add(s.ReadTimeout))

		msg, err := unmarshalMessage(r)
		if err != nil {
			return err
		}

		if s.ReadHandler != nil {
			s.ReadHandler.Handle(c, msg)
		} else {
			ProcessMsg(c, msg)
		}
	}
}

func unmarshalMessage(r io.Reader) (Msg, error) {
	var pm partialMsg
	if err := json.NewDecoder(r).Decode(&pm); err != nil {
		return nil, fmt.Errorf("invalid JSON message: %v", err)
	}

	genericUnmarshal := func(v interface{}, metaDst *meta) error {
		if err := json.Unmarshal(pm.Payload, v); err != nil {
			return fmt.Errorf("invalid %s message: %v", pm.Meta.T, err)
		}
		*metaDst = pm.Meta
		return nil
	}

	var msg Msg
	switch pm.Meta.T {
	case AuthMsg:
		var auth Auth
		if err := genericUnmarshal(&auth, &auth.meta); err != nil {
			return nil, err
		}
		msg = &auth

	case CallMsg:
		var call Call
		if err := genericUnmarshal(&call, &call.meta); err != nil {
			return nil, err
		}
		msg = &call

	case SubMsg:
		var sub Sub
		if err := genericUnmarshal(&sub, &sub.meta); err != nil {
			return nil, err
		}
		msg = &sub

	case PubMsg:
		var pub Pub
		if err := genericUnmarshal(&pub, &pub.meta); err != nil {
			return nil, err
		}
		msg = &pub

	case ErrMsg, OKMsg, ResMsg, EvntMsg:
		return nil, fmt.Errorf("invalid message %s for client peer", pm.Meta.T)
	default:
		return nil, fmt.Errorf("unknown message %s", pm.Meta.T)
	}

	return msg, nil
}
