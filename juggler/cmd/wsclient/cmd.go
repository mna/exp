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
	Usage:   "usage: disconnect CONN_ID",
	MinArgs: 1,
	Help:    "disconnect the connection identified by CONN_ID",

	Run: func(cmd *cmd, args ...string) {
		if len(args) < cmd.MinArgs {
			printErr(cmd.Usage)
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
	Usage:   "usage: close CONN_ID [STATUS_TEXT]",
	MinArgs: 1,
	Help:    "cleanly close the connection identified by CONN_ID, sending a websocket Close message",

	Run: func(cmd *cmd, args ...string) {
		if len(args) < cmd.MinArgs {
			printErr(cmd.Usage)
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
	Usage:   "usage: send CONN_ID MSG",
	MinArgs: 2,
	Help:    "send raw MSG (sent as-is) to the connection identified by CONN_ID",

	Run: func(cmd *cmd, args ...string) {
		if len(args) < cmd.MinArgs {
			printErr(cmd.Usage)
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
	Usage:   "usage: call CONN_ID URI [TIMEOUT_SEC [ARGS]]",
	MinArgs: 2,
	Help:    "send a CALL message to the connection identified by CONN_ID\n\tto URI with optional ARGS as JSON",

	Run: func(cmd *cmd, args ...string) {
		if len(args) < cmd.MinArgs {
			printErr(cmd.Usage)
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
				v = json.RawMessage(strings.Join(args[3:], " "))
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

var pubCmd = &cmd{
	Usage:   "usage: pub CONN_ID CHANNEL [ARGS]",
	MinArgs: 2,
	Help:    "send a PUB message to the connection identified by CONN_ID\n\tto CHANNEL with optional ARGS as JSON",

	Run: func(cmd *cmd, args ...string) {
		if len(args) < cmd.MinArgs {
			printErr(cmd.Usage)
			return
		}
		if c, ix := getConn(args[0]); c != nil {
			var v json.RawMessage
			if len(args) > 2 {
				v = json.RawMessage(strings.Join(args[2:], " "))
			}

			uuid, err := c.Pub(args[1], v)
			if err != nil {
				printErr("failed to send PUB message: %v", err)
				return
			}
			printf("[%d] sent PUB message %v", ix, uuid)
		} else {
			printErr("invalid connection ID")
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
		if len(args) < cmd.MinArgs {
			printErr(cmd.Usage)
			return
		}
		if c, ix := getConn(args[0]); c != nil {
			uuid, err := c.Sub(args[1], pattern)
			if err != nil {
				printErr("failed to send SUB message: %v", err)
				return
			}
			printf("[%d] sent SUB message %v", ix, uuid)
		} else {
			printErr("invalid connection ID")
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
