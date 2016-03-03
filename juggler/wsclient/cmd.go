package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/gorilla/websocket"
)

type cmd struct {
	Help string
	Run  func(...string)
}

var helpCmd = &cmd{
	Help: "usage: ? or help\n\tprint this message",

	Run: func(_ ...string) {
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			printf("- %s :\n\t%s\n", k, commands[k].Help)
		}
	},
}

var connectCmd = &cmd{
	Help: "usage: connect [URL [PROTO]]\n\tconnect to URL using subprotocol PROTO (defaults to juggler.0)",

	Run: func(args ...string) {
		var d websocket.Dialer

		addr := *defaultConnFlag
		if len(args) > 0 {
			addr = args[0]
		}

		h := http.Header{"Sec-WebSocket-Protocol": {"juggler.0"}}
		if len(args) > 1 {
			h.Set("Sec-WebSocket-Protocol", args[1])
		}

		conn, _, err := d.Dial(addr, h)
		if err != nil {
			printErr("error: %v", err)
			return
		}
		connections = append(connections, conn)
		printf("connected to %s [%d]", addr, len(connections))
		go read(len(connections), conn)
	},
}

var disconnectCmd = &cmd{
	Help: "usage: disconnect CONN_ID\n\tdisconnect the connection identified by CONN_ID",

	Run: func(args ...string) {
		if len(args) != 1 {
			printErr("usage: disconnect CONN_ID")
			return
		}
		if c, ix := getConn(args[0]); c != nil {
			c.Close()
			connections[ix] = nil
		} else {
			printErr("invalid connection ID")
		}
	},
}

var closeCmd = &cmd{
	Help: "usage: close CONN_ID\n\tcleanly close the connection identified by CONN_ID",

	Run: func(args ...string) {
		if len(args) != 1 {
			printErr("usage: close CONN_ID")
			return
		}
		if c, ix := getConn(args[0]); c != nil {
			if err := c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "bye"), time.Time{}); err != nil {
				printErr("failed to send close message: %v", err)
				return
			}
			c.Close()
			connections[ix] = nil
		} else {
			printErr("invalid connection ID")
		}
	},
}

var sendCmd = &cmd{
	Help: "usage: send CONN_ID MSG\n\tsend free-form MSG to the connection identified by CONN_ID",

	Run: func(args ...string) {
		if len(args) < 2 {
			printErr("usage: send CONN_ID MSG")
			return
		}
		if c, _ := getConn(args[0]); c != nil {
			if err := c.WriteMessage(websocket.TextMessage, []byte(strings.Join(args[1:], " "))); err != nil {
				printErr("WriteMessage failed: %v", err)
				return
			}
		} else {
			printErr("invalid connection ID")
		}
	},
}

var callCmd = &cmd{
	Help: "usage: call CONN_ID URI [TIMEOUT_SEC [ARGS]]\n\tsend a CALL message to the connection identified by CONN_ID\n\tto URI with optional ARGS as JSON",

	Run: func(args ...string) {
		if len(args) < 2 {
			printErr("usage: call CONN_ID URI [ARGS]")
			return
		}
		if c, _ := getConn(args[0]); c != nil {
			var to time.Duration
			if len(args) > 2 {
				d, err := time.ParseDuration(args[2])
				if err != nil {
					printErr("invalid timeout: %v", err)
					return
				}
				to = d
			}

			var v json.RawMessage
			if len(args) > 3 {
				v = json.RawMessage(args[3])
			}

			call, err := msg.NewCall(args[1], to, v)
			if err != nil {
				printErr("failed to create CALL message: %v", err)
				return
			}
			if err := c.WriteJSON(call); err != nil {
				printErr("WriteJSON failed: %v", err)
				return
			}
		} else {
			printErr("invalid connection ID")
		}
	},
}

func read(ix int, c *websocket.Conn) {
	for {
		_, b, err := c.ReadMessage()
		if err != nil {
			printErr("[%d] NextReader failed: %v; closing connection", ix, err)
			c.Close()
			return
		}
		printf("[%d] %v", ix, string(b))
	}
}

func getConn(arg string) (*websocket.Conn, int) {
	ix, err := strconv.Atoi(arg)
	if err != nil {
		printErr("argument error: %v", err)
		return nil, 0
	}
	if ix > 0 && ix <= len(connections) {
		if c := connections[ix-1]; c != nil {
			return c, ix - 1
		}
	}
	return nil, 0
}
