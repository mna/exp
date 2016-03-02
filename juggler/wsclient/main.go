package main

import (
	"flag"
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

var (
	commands    map[string]*cmd
	connections []*websocket.Conn
	term        *terminal.Terminal
)

func init() {
	commands = map[string]*cmd{
		"?":          helpCmd,
		"help":       helpCmd,
		"connect":    connectCmd,
		"disconnect": disconnectCmd,
		"send":       sendCmd,
		"close":      closeCmd,
	}
}

var (
	defaultConnFlag = flag.String("addr", "ws://localhost:9000/ws", "default dial address to use in connect commands")
)

func main() {
	var exitCode int

	flag.Parse()

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

	printf(welcomeMessage)
	for {
		l, err := t.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			printErr("wsclient: failed to read line: %v", err)
			exitCode = 1
			return
		}

		args := strings.Fields(l)
		if len(args) != 0 {
			if cmd := commands[args[0]]; cmd != nil {
				cmd.Run(args[1:]...)
			} else {
				printErr("unknown command %q", args[0])
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

func printf(msg string, args ...interface{}) {
	fmt.Fprintf(term, msg+"\n", args...)
}

func printErr(msg string, args ...interface{}) {
	term.Write(term.Escape.Red)
	printf(msg, args...)
	term.Write(term.Escape.Reset)
}
