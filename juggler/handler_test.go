package juggler

import (
	"io/ioutil"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestLimitedWriter(t *testing.T) {
	// use int8/uint8 to keep size reasonable
	checker := func(limit int16, n uint8) bool {
		// create a limited writer with the specified limit
		w := limitWriter(ioutil.Discard, int64(limit))
		// create the payload for each write
		p := make([]byte, n)

		var cnt, tot int
		var err error
		for {
			cnt, err = w.Write(p)
			tot += cnt
			if err != nil {
				break
			}
		}

		// property 1: the total number of bytes written cannot be > limit
		// except if limit < minLimit (4096).
		if limit < minWriteLimit {
			limit = minWriteLimit
		}
		if tot > int(limit) {
			return false
		}
		// property 2: by writing repeatedly, it necessarily terminates with
		// an errWriteLimitExceeded
		return err == errWriteLimitExceeded
	}
	assert.NoError(t, quick.Check(checker, nil))
}
