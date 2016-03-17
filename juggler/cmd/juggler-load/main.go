// Command juggler-load is a juggler load generator. It runs a
// number of client connections to a juggler server, and for a
// given duration, makes calls and collects results.
package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/exp/juggler"
	"github.com/PuerkitoBio/exp/juggler/client"
	"github.com/PuerkitoBio/exp/juggler/msg"
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

var tpl = template.Must(template.New("output").Parse(`
Connections:     {{ .Conns }}
Rate:            {{ .Rate | printf "%s" }}
Timeout:         {{ .Timeout | printf "%s" }}
Duration:        {{ .Duration | printf "%s" }}
Actual Duration: {{ .ActualDuration | printf "%s" }}

Calls:   {{ .Calls }}
OK:      {{ .OK }}
Errors:  {{ .Err }}
Results: {{ .Res }}
Expired: {{ .Exp }}

`))

type runStats struct {
	Conns          int
	Rate           time.Duration
	Timeout        time.Duration
	Duration       time.Duration
	ActualDuration time.Duration

	Calls int64
	OK    int64
	Err   int64
	Res   int64
	Exp   int64
}

func main() {
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

	log.SetFlags(0)

	if *connFlag <= 0 {
		log.Fatalf("invalid -c value, must be greater than 0")
	}

	stats := &runStats{
		Conns:    *connFlag,
		Rate:     *callRateFlag,
		Timeout:  *callTimeoutFlag,
		Duration: *durationFlag,
	}

	// TODO : collect expvars before and after if flag is set

	clientStarted := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(stats.Conns)

	stop := make(chan struct{})
	for i := 0; i < stats.Conns; i++ {
		go runClient(stats, &wg, clientStarted, stop)
	}

	// start clients with some jitter, up to 10ms
	log.Printf("%d connections started...", stats.Conns)
	start := time.Now()
	for i := 0; i < stats.Conns; i++ {
		<-time.After(time.Duration(rand.Intn(int(10 * time.Millisecond))))
		<-clientStarted
	}

	// run for the requested duration and signal stop
	<-time.After(stats.Duration)
	close(stop)
	log.Printf("stopping...")

	// wait for completion
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(time.Second):
			log.Fatalf("failed to stop clients")
		}
	}()
	wg.Wait()
	close(done)
	end := time.Now()
	stats.ActualDuration = end.Sub(start)

	log.Printf("stopped.")
	if err := tpl.Execute(os.Stdout, stats); err != nil {
		log.Fatalf("template.Execute failed: %v", err)
	}
}

func runClient(stats *runStats, wg *sync.WaitGroup, started chan<- struct{}, stop <-chan struct{}) {
	defer wg.Done()

	var wgResults sync.WaitGroup

	cli, err := client.Dial(
		&websocket.Dialer{Subprotocols: []string{*subprotoFlag}},
		*addrFlag, nil,
		client.SetLogFunc(juggler.DiscardLog),
		client.SetHandler(client.HandlerFunc(func(ctx context.Context, c *client.Client, m msg.Msg) {
			switch m.Type() {
			case msg.ResMsg:
				atomic.AddInt64(&stats.Res, 1)
			case client.ExpMsg:
				atomic.AddInt64(&stats.Exp, 1)
			case msg.OKMsg:
				atomic.AddInt64(&stats.OK, 1)
				return
			case msg.ErrMsg:
				atomic.AddInt64(&stats.Err, 1)
				return
			default:
				log.Fatalf("unexpected message type %s", m.Type())
			}
			wgResults.Done()
		})))

	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}

	var after time.Duration
	started <- struct{}{}
loop:
	for {
		select {
		case <-stop:
			break loop
		case <-time.After(after):
		}

		wgResults.Add(1)
		atomic.AddInt64(&stats.Calls, 1)
		_, err := cli.Call(*uriFlag, *payloadFlag, stats.Timeout)
		if err != nil {
			log.Fatalf("Call failed: %v", err)
		}
		after = stats.Rate
	}
	// wait for sent calls to return or expire
	wgResults.Wait()

	if err := cli.Close(); err != nil {
		log.Fatalf("Close failed: %v", err)
	}
}
