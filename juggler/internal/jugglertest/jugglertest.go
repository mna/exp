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

// Printf implements the LogFunc-compatible signature. It logs to log.Printf
// if testing.Verbose is true.
func (d *DebugLog) Printf(s string, args ...interface{}) {
	atomic.AddInt64(&d.n, 1)
	if testing.Verbose() {
		log.Printf(s, args...)
	}
}

// Calls returns the number of calls received by the logger.
func (d *DebugLog) Calls() int {
	return int(atomic.LoadInt64(&d.n))
}
