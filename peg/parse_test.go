package main

import (
	"strings"
	"testing"
)

var invalidParseCases = map[string]string{
	"":             "file:1:1 (0): no match found",
	"a":            "file:1:1 (0): no match found",
	"abc":          "file:1:1 (0): no match found",
	" ":            "file:1:1 (0): no match found",
	"a ←":          "file:1:1 (0): no match found",
	"{":            "file:1:1 (0): no match found",
	"{{}":          "file:1:1 (0): no match found",
	"a ← nil:b":    "file:1:5 (6): rule Identifier: identifier is a reserved word",
	"a ← b\nb ←":   "file:1:1 (0): no match found",
	`a ← [\pA]`:    "file:1:8 (9): rule UnicodeClassEscape: invalid Unicode class escape",
	"\xfe":         "file:1:1 (0): invalid encoding",
	`a = '\"'`:     "file:1:7 (6): invalid escape character",
	`a = "\'"`:     "file:1:7 (6): invalid escape character",
	`a = '\xzz`:    "file:1:7 (6): invalid hexadecimal escape",
	`a = '\091`:    "file:1:7 (6): invalid octal escape",
	`a = "b`:       "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = `b":       "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = 'b":       "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = '\uA'`:    "file:1:7 (6): invalid Unicode escape",
	`a = '\UA012'`: "file:1:7 (6): invalid Unicode escape",
	`a = [b`:       "file:1:5 (4): rule CharClassMatcher: character class not terminated",
	`a = b /`:      "file:1:7 (6): invalid choice expression",
	`a = +`:        "file:1:1 (0): no match found",
	`a = *`:        "file:1:1 (0): no match found",
	`a = ?`:        "file:1:1 (0): no match found",
}

func TestInvalidParseCases(t *testing.T) {
	for tc, exp := range invalidParseCases {
		_, err := Parse("file", strings.NewReader(tc))
		if err == nil {
			t.Errorf("%q: want error, got none", tc)
			continue
		}
		if err.Error() != exp {
			t.Errorf("%q: want \n%s\n, got \n%s\n", tc, exp, err)
		}
	}
}
