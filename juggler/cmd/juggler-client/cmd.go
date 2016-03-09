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

var (
	commands    map[string]*cmd
	connections []*juggler.Client
)

func init() {
	commands = map[string]*cmd{
		"?":          helpCmd,
		"help":       helpCmd,
		"exit":       exitCmd,
		"connect":    connectCmd,
		"disconnect": disconnectCmd,
		"send":       sendCmd,
		"close":      closeCmd,
		"call":       callCmd,
		"pub":        pubCmd,
		"sub":        subCmd,
		"psub":       psubCmd,
		"unsb":       unsbCmd,
		"punsb":      punsbCmd,
	}
}

type cmd struct {
	Usage   string
	MinArgs int
	Help    string
	Run     func(*cmd, ...string)
}

var helpCmd = &cmd{
	Usage:   "usage: ? or help",
	MinArgs: 0,
	Help:    "print this message",

	Run: func(_ *cmd, _ ...string) {
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			printf("- %s :\n\t%s\n\t%s\n", k, commands[k].Usage, commands[k].Help)
		}
	},
}

var exitCmd = &cmd{
	Usage:   "usage: exit or ctrl-D",
	MinArgs: 0,
	Help:    "quit the program",

	Run: func(_ *cmd, _ ...string) {
		// special-cased in the readline loop
	},
}

var connectCmd = &cmd{
	Usage:   "usage: connect [URL [PROTO]]",
	MinArgs: 0,
	Help:    fmt.Sprintf("connect to URL using subprotocol PROTO (defaults to %s)", *defaultSubprotoFlag),

	Run: func(_ *cmd, args ...string) {
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
			printErr("Dial failed: %v", err)
			return
		}
		conn.LogFunc = printErr

		connections = append(connections, conn)
		printf("[%d] connected to %s", len(connections), addr)
	},
}

type connMsgLogger int

func (l connMsgLogger) Handle(m msg.Msg) {
	var s string
	switch m := m.(type) {
	case *msg.Err:
		s = fmt.Sprintf("for %s %v (%s)", m.Payload.ForType, m.Payload.For, m.Payload.Message)
	case *msg.OK:
		s = fmt.Sprintf("for %s %v", m.Payload.ForType, m.Payload.For)
	case *msg.Res:
		n := len(m.Payload.Args)
		if n > 20 {
			n = 20
		}
		val := string(m.Payload.Args[:n])
		s = fmt.Sprintf("for %s %v (%s)", msg.CallMsg, m.Payload.For, val)
	case *msg.Exp:
		s = fmt.Sprintf("for %s %v", msg.CallMsg, m.Payload.For)
	case *msg.Evnt:
		n := len(m.Payload.Args)
		if n > 20 {
			n = 20
		}
		val := string(m.Payload.Args[:n])
		s = fmt.Sprintf("for %s %v (%s)", msg.PubMsg, m.Payload.For, val)
	}
	printf("[%d] <<< %-4s message: %v %s", l, m.Type(), m.UUID(), s)
}

var disconnectCmd = &cmd{
	Usage:   "usage: disconnect CONN_ID",
	MinArgs: 1,
	Help:    "disconnect the connection identified by CONN_ID",

	Run: func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			c.Close()
			connections[ix] = nil
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	},
}

var closeCmd = &cmd{
	Usage:   "usage: close CONN_ID [STATUS_TEXT]",
	MinArgs: 1,
	Help:    "cleanly close the connection identified by CONN_ID, sending a\n\twebsocket Close message",

	Run: func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			wsc := c.UnderlyingConn()
			st := "bye"
			if len(args) > 1 {
				st = args[1]
			}

			if err := wsc.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, st), time.Time{}); err != nil {

				printErr("[%d] WriteControl failed: %v", ix+1, err)
				return
			}

			c.Close()
			connections[ix] = nil
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	},
}

var sendCmd = &cmd{
	Usage:   "usage: send CONN_ID MSG",
	MinArgs: 2,
	Help:    "send raw MSG (sent as-is) to the connection identified by CONN_ID",

	Run: func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			wsc := c.UnderlyingConn()
			if err := wsc.WriteMessage(websocket.TextMessage, []byte(strings.Join(args[1:], " "))); err != nil {
				printErr("[%d] WriteMessage failed: %v", ix+1, err)
				return
			}
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	},
}

var callCmd = &cmd{
	Usage:   "usage: call CONN_ID URI [TIMEOUT_SEC [ARGS]]",
	MinArgs: 2,
	Help:    "send a CALL message to the connection identified by CONN_ID\n\tto URI with optional ARGS as JSON",

	Run: func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			var to time.Duration
			if len(args) > 2 {
				d, err := time.ParseDuration(args[2])
				if err != nil {
					printErr("[%d] invalid timeout: %v", ix+1, err)
					return
				}
				to = d
			}

			var v json.RawMessage
			if len(args) > 3 {
				v = json.RawMessage(strings.Join(args[3:], " "))
			}

			uuid, err := c.Call(args[1], v, to)
			if err != nil {
				printErr("[%d] Call failed: %v", ix+1, err)
				return
			}
			printf("[%d] >>> CALL message: %v", ix+1, uuid)
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	},
}

var pubCmd = &cmd{
	Usage:   "usage: pub CONN_ID CHANNEL [ARGS]",
	MinArgs: 2,
	Help:    "send a PUB message to the connection identified by CONN_ID\n\tto CHANNEL with optional ARGS as JSON",

	Run: func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			var v json.RawMessage
			if len(args) > 2 {
				v = json.RawMessage(strings.Join(args[2:], " "))
			}

			uuid, err := c.Pub(args[1], v)
			if err != nil {
				printErr("[%d] Pub failed: %v", ix+1, err)
				return
			}
			printf("[%d] >>> PUB  message: %v", ix+1, uuid)
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	},
}

var subCmd = &cmd{
	Usage:   "usage: sub CONN_ID CHANNEL",
	MinArgs: 2,
	Help:    "send a SUB message to the connection identified by CONN_ID\n\tto subscribe the connection to the CHANNEL",

	Run: getSubFunc(false),
}

var psubCmd = &cmd{
	Usage:   "usage: psub CONN_ID CHANNEL_PATTERN",
	MinArgs: 2,
	Help:    "send a SUB message to the connection identified by CONN_ID\n\tto subscribe the connection to the pattern CHANNEL_PATTERN",

	Run: getSubFunc(true),
}

func getSubFunc(pattern bool) func(*cmd, ...string) {
	return func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			uuid, err := c.Sub(args[1], pattern)
			if err != nil {
				printErr("[%d] Sub failed: %v", ix+1, err)
				return
			}
			printf("[%d] >>> SUB  message: %v", ix+1, uuid)
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	}
}

var unsbCmd = &cmd{
	Usage:   "usage: unsb CONN_ID CHANNEL",
	MinArgs: 2,
	Help:    "send an UNSB message to the connection identified by CONN_ID\n\tto unsubscribe the connection from the CHANNEL",

	Run: getUnsbFunc(false),
}

var punsbCmd = &cmd{
	Usage:   "usage: punsb CONN_ID CHANNEL_PATTERN",
	MinArgs: 2,
	Help:    "send a UNSB message to the connection identified by CONN_ID\n\tto unsubscribe the connection from the pattern CHANNEL_PATTERN",

	Run: getUnsbFunc(true),
}

func getUnsbFunc(pattern bool) func(*cmd, ...string) {
	return func(cmd *cmd, args ...string) {
		if c, ix := getConn(args[0]); c != nil {
			uuid, err := c.Unsb(args[1], pattern)
			if err != nil {
				printErr("[%d] Unsb failed: %v", ix+1, err)
				return
			}
			printf("[%d] >>> UNSB message: %v", ix+1, uuid)
		} else {
			printErr("invalid connection ID: %s", args[0])
		}
	}
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
