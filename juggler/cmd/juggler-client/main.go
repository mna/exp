package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const welcomeMessage = `
Welcome to the juggler client. Enter ? or help for the available
commands. Press ^D (ctrl-D) to exit.
`

var (
	term *terminal.Terminal
)

var (
	defaultConnFlag     = flag.String("addr", "ws://localhost:9000/ws", "default dial address to use in connect commands")
	defaultSubprotoFlag = flag.String("subprotocol", "juggler.0", "default subprotocol to request in the websocket handshake")
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
			printErr("failed to read line: %v", err)
			exitCode = 1
			return
		}

		args := strings.Fields(l)
		if len(args) != 0 {
			if cmd := commands[args[0]]; cmd != nil {
				args = args[1:]
				if len(args) < cmd.MinArgs {
					printErr(cmd.Usage)
					continue
				}
				if cmd == exitCmd {
					return
				}
				cmd.Run(cmd, args...)
			} else {
				printErr("unknown command: %q", args[0])
			}
		}
	}
}

func setupTerminal() (*terminal.Terminal, func()) {
	// setup terminal
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalf("failed to initialize the terminal: %v", err)
	}
	cleanUp := func() { terminal.Restore(0, oldState) }

	var screen = struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := terminal.NewTerminal(screen, "juggler> ")
	return t, cleanUp
}

func printf(msg string, args ...interface{}) {
	t := time.Now().Format(time.RFC3339)
	fmt.Fprintf(term, t+" | "+msg+"\n", args...)
}

func printErr(msg string, args ...interface{}) {
	term.Write(term.Escape.Red)
	printf(msg, args...)
	term.Write(term.Escape.Reset)
}
