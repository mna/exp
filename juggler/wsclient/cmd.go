package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/websocket"
)

type cmd struct {
	Help string
	Run  func(...string)
}

var helpCmd = &cmd{
	Help: "print this message",

	Run: func(_ ...string) {
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			print("%s\n\t%s\n", k, commands[k].Help)
		}
	},
}

var connectCmd = &cmd{
	Help: "usage: connect [URL [PROTO]]\n\tconnect to URL using subprotocol PROTO (defaults to juggler.1)",

	Run: func(args ...string) {
		var d websocket.Dialer

		addr := *defaultConnFlag
		if len(args) > 0 {
			addr = args[0]
		}

		h := http.Header{"Sec-WebSocket-Protocol": {"juggler.1"}}
		if len(args) > 1 {
			h.Set("Sec-WebSocket-Protocol", args[1])
		}

		conn, _, err := d.Dial(addr, h)
		if err != nil {
			printErr("error: %v", err)
			return
		}
		connections = append(connections, conn)
		print("connected to %s [%d]", addr, len(connections))
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
		print("[%d] %v", ix, string(b))
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
