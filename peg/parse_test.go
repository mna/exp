package main

import (
	"strings"
	"testing"
)

var invalidParseCases = []string{
	"",
	"a",
}

var expInvalidParseErrs = []string{
	"file: no match found",
	"file:1:1 (0): no match found",
}

func TestInvalidParseCases(t *testing.T) {
	for i, tc := range invalidParseCases {
		_, err := Parse("file", strings.NewReader(tc))
		if err == nil {
			t.Errorf("[%d] %q: want error, got none", i, tc)
			continue
		}
		if err.Error() != expInvalidParseErrs[i] {
			t.Errorf("[%d] %q: want \n%s\n, got \n%s\n", i, tc, expInvalidParseErrs[i], err)
		}
	}
}
