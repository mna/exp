// TODO : document package.
package callee

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
)

// ErrCallExpired is returned when a call is processed but the
// call timeout is exceeded, meaning that the client is no longer
// expecting the result. The result is dropped and this error is
// returned from InvokeAndStoreResult.
var ErrCallExpired = errors.New("call expired")

// Thunk is the function signature for functions that handle calls
// to a URI. Generally, they should be used to decode the arguments
// to the type expected by the actual underlying function, and to
// transfer the results back in the generic empty interface.
type Thunk func(*msg.CallPayload) (interface{}, error)

// Callee is a peer that handles call requests for some URIs.
type Callee struct {
	// Broker is the callee broker to use to listen for call requests
	// and to store results.
	Broker broker.CalleeBroker

	// LogFunc is the logging function to use. If nil, log.Printf
	// is used. It can be set to juggler.DiscardLog to disable logging.
	LogFunc func(string, ...interface{})
}

// SplitByHashSlot takes a list of URIs and splits them into groups
// of URIs belonging to the same redis cluster hash slot. URIs in
// the same hash slot can be listened to using the same broker.CallsConn,
// optimizing the number of redis connections.
//
// See the redis cluster documentation for details:
// http://redis.io/topics/cluster-tutorial
func SplitByHashSlot(uris []string) [][]string {
	// TODO : implement when redis-cluster package is ready
	return nil
}

// InvokeAndStoreResult processes the provided call payload by calling
// fn with the payload's arguments, and storing the result so that
// it can be sent back to the caller. If the call timeout is exceeded,
// the result is dropped and ErrCallExpired is returned.
func (c *Callee) InvokeAndStoreResult(cp *msg.CallPayload, fn Thunk) error {
	ttl := cp.TTLAfterRead
	start := time.Now()

	v, err := fn(cp)
	if remain := ttl - time.Now().Sub(start); remain > 0 {
		// register the result
		if err := c.storeResult(cp, v, err, remain); err != nil {
			return err
		}
	}
	return ErrCallExpired
}

// Listen is a helper method that listens for call requests for the
// requested URIs and calls the corresponding Thunk to execute the
// request. The m map has URIs as keys, and the associated Thunk
// function as value. If a redis cluster is used, all URIs in m
// must belong to the same hash slot (see SplitByHashSlot).
//
// The method implements a single-producer, single-consumer helper,
// where a single redis connection is used to listen for call requests
// on the URIs, and for each request, a single goroutine executes
// the calls and stores the results. More powerful concurrency
// patterns can be implemented using Callee.Broker.Calls directly,
// and starting multiple consumer goroutines reading from the same calls
// channel and calling InvokeAndStoreResult to process each call request.
//
// The function blocks until the call request loop exits. It returns
// the error that caused the loop to stop, or the error to initiate
// the connection to the broker.
func (c *Callee) Listen(m map[string]Thunk) error {
	if len(m) == 0 {
		return nil
	}

	uris := make([]string, 0, len(m))
	for k := range m {
		uris = append(uris, k)
	}
	conn, err := c.Broker.Calls(uris...)
	if err != nil {
		return err
	}
	defer conn.Close()

	for cp := range conn.Calls() {
		if err := c.InvokeAndStoreResult(cp, m[cp.URI]); err != nil {
			if err == ErrCallExpired {
				logf(c.LogFunc, "dropping expired message %v", cp.MsgUUID)
				continue
			}
			logf(c.LogFunc, "storeResult failed for message %v: %v", cp.MsgUUID, err)
			continue
		}
	}
	return conn.CallsErr()
}

type errResult struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Callee) storeResult(cp *msg.CallPayload, v interface{}, e error, timeout time.Duration) error {
	// if there's an error, that's what gets stored
	if e != nil {
		var er errResult
		er.Error.Message = e.Error()
		v = er
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	rp := &msg.ResPayload{
		ConnUUID: cp.ConnUUID,
		MsgUUID:  cp.MsgUUID,
		URI:      cp.URI,
		Args:     b,
	}
	return c.Broker.Result(rp, timeout)
}

func logf(fn func(string, ...interface{}), s string, args ...interface{}) {
	if fn != nil {
		fn(s, args...)
	} else {
		log.Printf(s, args...)
	}
}
