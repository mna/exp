package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gorilla/websocket"
)

const welcomeMessage = `
Welcome to the websocket client. Enter ? or help for the available
commands. Press ^D (ctrl-D) to exit.
`

var commands map[string]*cmd

var connections []*websocket.Conn

func init() {
	commands = map[string]*cmd{
		"?":    helpCmd,
		"help": helpCmd,
		/*
			"connect":    connect,
			"disconnect": disconnect,
			"send":       send,
		*/
	}
}

var term *terminal.Terminal

func main() {
	var exitCode int

	// call os.Exit in a defer, otherwise defer to reset the terminal
	// will not be run.
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	// setup and restore the terminal
	t, fn := setupTerminal()
	defer fn()
	term = t

	print(welcomeMessage)
	for {
		l, err := t.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("wsclient: failed to read line: %v", err)
			exitCode = 1
			return
		}

		args := strings.Fields(l)
		if len(args) != 0 {
			if cmd := commands[args[0]]; cmd != nil {
				cmd.Run(args[1:]...)
			} else {
				print("unknown command %q", args[0])
			}
		}
	}
}

func setupTerminal() (*terminal.Terminal, func()) {
	// setup terminal
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalf("wsclient: failed to initialize the terminal: %v", err)
	}
	cleanUp := func() { terminal.Restore(0, oldState) }

	var screen = struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := terminal.NewTerminal(screen, "ws> ")
	return t, cleanUp
}

func print(msg string, args ...interface{}) {
	fmt.Fprintf(term, msg+"\n", args...)
}

/*
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
*/
