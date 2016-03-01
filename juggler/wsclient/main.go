package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/kless/term/readline"
	"github.com/mitchellh/go-homedir"
)

var commands map[string]func(...string)

var connections []*websocket.Conn

func init() {
	commands = map[string]func(...string){
		"?":          showHelp,
		"help":       showHelp,
		"connect":    connect,
		"disconnect": disconnect,
		"send":       send,
	}
}

func main() {
	fmt.Println(`
Welcome to the websocket client. Enter ? or help for the available
commands. Press ^D (ctrl-D) to exit.
`)

	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("wsclient: failed to read home directory: %v", err)
	}
	hist, err := readline.NewHistory(filepath.Join(home, ".wsclienthist"))
	if err != nil {
		log.Fatalf("wsclient: failed to create history file: %v", err)
	}

	line, err := readline.NewDefaultLine(hist)
	if err != nil {
		log.Fatalf("wsclient: failed to create line: %v", err)
	}

	go func() {
		for {
			select {
			case <-readline.ChanCtrlC:
			case <-readline.ChanCtrlD:
				line.Restore()
				for _, conn := range connections {
					if conn != nil {
						conn.Close()
					}
				}
				hist.Save()
				os.Exit(0)
			}
		}
	}()

	for {
		l, err := line.Read()
		if err != nil {
			line.Restore()
			log.Fatalf("wsclient: failed to read line: %v", err)
		}
		args := strings.Fields(l)
		if len(args) != 0 {
			if cmd := commands[args[0]]; cmd != nil {
				cmd(args[1:]...)
			} else {
				fmt.Printf("unknown command %q\r\n", args[0])
			}
		}
	}
}

func showHelp(_ ...string) {
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s\r\n", k)
	}
}

func connect(args ...string) {
	var d websocket.Dialer

	if len(args) < 1 {
		fmt.Print("usage: connect URL [PROTO]\r\n")
		return
	}

	var h http.Header
	if len(args) == 2 {
		h = http.Header{"Sec-WebSocket-Protocol": {args[1]}}
	}
	conn, _, err := d.Dial(args[0], h)
	if err != nil {
		fmt.Printf("error: %v\r\n", err)
		return
	}
	connections = append(connections, conn)
	fmt.Printf("connected to %s [%d]\r\n", args[0], len(connections))
	go read(len(connections), conn)
}

func read(ix int, c *websocket.Conn) {
	for {
		_, b, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("[%d] NextReader failed: %v; closing connection\r\n", ix, err)
			c.Close()
			return
		}
		fmt.Printf("[%d] %v\r\n", ix, string(b))
	}
}

func getConn(arg string) (*websocket.Conn, int) {
	ix, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Printf("argument error: %v\r\n", err)
		return nil, 0
	}
	if ix > 0 && ix <= len(connections) {
		if c := connections[ix-1]; c != nil {
			return c, ix - 1
		}
	}
	return nil, 0
}

func disconnect(args ...string) {
	if len(args) != 1 {
		fmt.Print("usage: disconnect CONN_ID\r\n")
		return
	}
	if c, ix := getConn(args[0]); c != nil {
		c.Close()
		connections[ix] = nil
	}
}

func send(args ...string) {
	if len(args) < 2 {
		fmt.Print("usage: send CONN_ID MSG\r\n")
		return
	}
	if c, _ := getConn(args[0]); c != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(strings.Join(args[1:], " "))); err != nil {
			fmt.Printf("WriteMessage failed: %v\r\n", err)
			return
		}
	}
}
