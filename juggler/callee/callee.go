package callee

import (
	"encoding/json"
	"log"
	"time"

	"github.com/PuerkitoBio/exp/juggler/broker"
	"github.com/PuerkitoBio/exp/juggler/msg"
)

// Thunk is the function signature for functions that handle calls
// to a URI. Generally, they should be used to decode the arguments
// to the type expected by the actual underlying function, and to
// transfer the results back in the generic empty interface.
type Thunk func(string, json.RawMessage) (interface{}, error)

// Callee is a peer that handles call requests for some URIs.
type Callee struct {
	// Broker is the callee broker to use to listen for call requests
	// and to store results.
	Broker broker.CalleeBroker

	// LogFunc is the logging function to use. If nil, log.Printf
	// is used. It can be set to juggler.DiscardLog to disable logging.
	LogFunc func(string, ...interface{})
}

// TODO : helper to spin as many Listen as required for the URIs in
// a redis cluster setting.

// Listen listens for call requests for the requested URIs and calls the
// corresponding Thunk to execute the request. The m parameter has
// URIs as keys, and the associated Thunk function as value.
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
		ttl := cp.TTLAfterRead
		start := time.Now()

		fn := m[cp.URI]
		v, err := fn(cp.URI, cp.Args)

		if remain := ttl - time.Now().Sub(start); remain > 0 {
			// register the result
			if err := c.storeResult(cp, v, err, remain); err != nil {
				logf(c.LogFunc, "storeResult failed for message %v: %v", cp.MsgUUID, err)
				continue
			}
		} else {
			logf(c.LogFunc, "TTL expired for message %v", cp.MsgUUID)
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
