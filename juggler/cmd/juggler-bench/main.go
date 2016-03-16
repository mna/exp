package main

import (
	"flag"
	"time"
)

var (
	addrFlag        = flag.String("addr", "ws://localhost:9000/ws", "Server `address`.")
	subprotoFlag    = flag.String("proto", "juggler.0", "Websocket `subprotocol`.")
	connFlag        = flag.Int("c", 100, "Number of `connections`.")
	durationFlag    = flag.Duration("d", 10*time.Second, "Run `duration`.")
	callRateFlag    = flag.Duration("r", 100*time.Millisecond, "Call `rate` per connection.")
	collectVarsFlag = flag.Bool("vars", false, "Collect expvars before and after execution.")
	helpFlag        = flag.Bool("help", false, "Show help.")
)

func main() {
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}
}
