// Command juggler-client can be used to send commands to a juggler
// server.
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
	defaultConnFlag     = flag.String("addr", "ws://localhost:9000/ws", "Default server `address` used in connect command.")
	defaultSubprotoFlag = flag.String("proto", "juggler.0", "Default `subprotocol` used in connect command.")
	timestampFmtFlag    = flag.String("timestamp", time.StampMilli, "Timestamp `format`, using Go time format syntax.")
	helpFlag            = flag.Bool("help", false, "Show help.")
)

func main() {
	var exitCode int

	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

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

	printfTs(welcomeMessage, "")
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

func printfTs(msg, ts string, args ...interface{}) {
	if ts != "" {
		t := time.Now().Format(ts)
		msg = t + " | " + msg
	}
	fmt.Fprintf(term, msg+"\n", args...)
}

func printf(msg string, args ...interface{}) {
	printfTs(msg, *timestampFmtFlag, args...)
}

func printErr(msg string, args ...interface{}) {
	term.Write(term.Escape.Red)
	printf(msg, args...)
	term.Write(term.Escape.Reset)
}
