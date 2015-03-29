package main

import (
	"strings"
	"testing"
)

var invalidParseCases = []string{
	"",
	"a",
	"abc",
	" ",
	"a ←",
	"{",
	"{{}",
	"a ← nil:b", // U#2190 is 3 bytes long
	"a ← b\nb ←",
	`a ← [\pA]`,
}

var expInvalidParseErrs = []string{
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
	"file:1:5 (6): rule Identifier: identifier is a reserved word",
	"file:2:1 (8): no match found",
	"file:1:8 (9): rule UnicodeClassEscape: invalid Unicode class escape",
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
