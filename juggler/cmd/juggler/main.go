package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/juggler"
	"github.com/gorilla/websocket"
)

var (
	portFlag      = flag.Int("port", 9000, "port to listen on")
	readLimitFlag = flag.Int("read-limit", 4096, "read message size limit")
	readTOFlag    = flag.Duration("read-timeout", 10*time.Second, "read deadline duration")
	writeTOFlag   = flag.Duration("write-timeout", 10*time.Second, "write deadline duration")
)

func main() {
	flag.Parse()

	upg := &websocket.Upgrader{Subprotocols: juggler.Subprotocols}
	srv := &juggler.Server{
		ReadLimit:    int64(*readLimitFlag),
		ReadTimeout:  *readTOFlag,
		WriteTimeout: *writeTOFlag,
		ConnHandler:  juggler.ConnHandlerFunc(juggler.LogConn),
		ReadHandler:  juggler.MsgHandlerFunc(juggler.LogMsg),
	}
	http.Handle("/ws", juggler.Upgrade(upg, srv))

	log.Printf("juggler: listening on port %d", *portFlag)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil); err != nil {
		log.Fatalf("juggler: ListenAndServe failed: %v", err)
	}
}
