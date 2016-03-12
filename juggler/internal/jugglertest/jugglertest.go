package jugglertest

import (
	"log"
	"sync/atomic"
	"testing"
)

// DebugLog is a logger that counts the number of calls it receives,
// and logs using log.Printf if testing.Verbose is set.
type DebugLog struct {
	T *testing.T
	n int64
}

func (d *debugLog) Printf(s string, args ...interface{}) {
	atomic.AddInt64(&d.n, 1)
	if testing.Verbose() {
		log.Printf(s, args...)
	}
}

func (d *debugLog) Calls() int {
	return int(atomic.LoadInt64(&d.n))
}
