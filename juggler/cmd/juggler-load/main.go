// Command juggler-load is a juggler load generator. It runs a
// number of client connections to a juggler server, and for a
// given duration, makes calls and collects results.
package main

import (
	"flag"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/PuerkitoBio/exp/juggler/client"
	"github.com/gorilla/websocket"
)

var (
	addrFlag        = flag.String("addr", "ws://localhost:9000/ws", "Server `address`.")
	subprotoFlag    = flag.String("proto", "juggler.0", "Websocket `subprotocol`.")
	connFlag        = flag.Int("c", 100, "Number of `connections`.")
	durationFlag    = flag.Duration("d", 10*time.Second, "Run `duration`.")
	callRateFlag    = flag.Duration("r", 100*time.Millisecond, "Call `rate` per connection.")
	callTimeoutFlag = flag.Duration("t", time.Second, "Call `timeout`.")
	uriFlag         = flag.String("u", "test.delay", "Call `URI`.")
	payloadFlag     = flag.String("p", "100", "Call `payload`.")
	collectVarsFlag = flag.Bool("vars", false, "Collect expvars before and after execution.")
	helpFlag        = flag.Bool("help", false, "Show help.")
)

func main() {
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

	if *connFlag <= 0 {
		log.Fatalf("invalid -c value, must be greater than 0")
	}

	clientStarted := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(*connFlag)

	stop := make(chan struct{})
	for i := 0; i < *connFlag; i++ {
		go func(i int) {
			defer wg.Done()

			var wgResults sync.WaitGroup

			cli, err := client.Dial(
				&websocket.Dialer{Subprotocols: []string{*subprotoFlag}},
				*addrFlag, nil,
				client.SetHandler(clientHandler(&clientStats)))

			if err != nil {
				log.Fatalf("Dial failed: %v", err)
			}

			after := time.Duration(0)
			clientStarted <- struct{}{}
		loop:
			for {
				select {
				case <-stop:
					break loop
				case <-time.After(after):
				}

				wgResults.Add(1)
				_, err := cli.Call(*uriFlag, *payloadFlag, *callTimeoutFlag)
				if err != nil {
					log.Fatalf("Call failed: %v", err)
				}
				after = *callRateFlag
			}
			// wait for sent calls to return or expire
			wgResults.Wait()

			if err := cli.Close(); err != nil {
				log.Fatalf("Close failed: %v", err)
			}
		}(i)
	}

	// start clients with some jitter, up to -r flag value
	start := time.Now()
	for i := 0; i < *connFlag; i++ {
		<-time.After(time.Duration(rand.Intn(int(*callRateFlag))))
		<-clientStarted
	}
	// wait for completion
	wg.Wait()
	end := time.Now()
}
