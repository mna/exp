package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/exp/juggler"
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
	Help: fmt.Sprintf("usage: connect [URL [PROTO]]\n\tconnect to URL using subprotocol PROTO (defaults to %s)", *defaultSubprotoFlag),

	Run: func(args ...string) {
		var d websocket.Dialer

		addr := *defaultConnFlag
		if len(args) > 0 {
			addr = args[0]
		}

		head := http.Header{"Sec-WebSocket-Protocol": {*defaultSubprotoFlag}}
		if len(args) > 1 {
			head.Set("Sec-WebSocket-Protocol", args[1])
		}

		conn, err := juggler.Dial(&d, addr, head, connMsgLogger(len(connections)+1))
		if err != nil {
			printErr("error: %v", err)
			return
		}
		connections = append(connections, conn)
		printf("connected to %s [%d]", addr, len(connections))
	},
}

type connMsgLogger int

func (l connMsgLogger) Handle(m msg.Msg) {
	printf("[%d] %s %v", l, m.Type(), m.UUID())
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
	Help: "usage: close CONN_ID [STATUS_TEXT]\n\tcleanly close the connection identified by CONN_ID, sending a websocket Close message",

	Run: func(args ...string) {
		if len(args) < 1 {
			printErr("usage: close CONN_ID [STATUS_TEXT]")
			return
		}
		if c, ix := getConn(args[0]); c != nil {
			wsc := c.UnderlyingConn()
			st := "bye"
			if len(args) > 1 {
				st = args[1]
			}
			if err := wsc.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, st), time.Time{}); err != nil {
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
			wsc := c.UnderlyingConn()
			if err := wsc.WriteMessage(websocket.TextMessage, []byte(strings.Join(args[1:], " "))); err != nil {
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
		if c, ix := getConn(args[0]); c != nil {
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

			uuid, err := c.Call(args[1], v, to)
			if err != nil {
				printErr("failed to send CALL message: %v", err)
				return
			}
			printf("[%d] sent CALL message %v", ix, uuid)
		} else {
			printErr("invalid connection ID")
		}
	},
}

func getConn(arg string) (*juggler.Client, int) {
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
