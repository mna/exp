package main

import (
	"strings"
	"testing"
)

var invalidParseCases = map[string]string{
	"":           "file:1:1 (0): no match found",
	"a":          "file:1:1 (0): no match found",
	"abc":        "file:1:1 (0): no match found",
	" ":          "file:1:1 (0): no match found",
	`a = +`:      "file:1:1 (0): no match found",
	`a = *`:      "file:1:1 (0): no match found",
	`a = ?`:      "file:1:1 (0): no match found",
	"a ←":        "file:1:1 (0): no match found",
	"a ← b\nb ←": "file:1:1 (0): no match found",
	"a ← nil:b":  "file:1:5 (6): rule Identifier: identifier is a reserved word",
	"\xfe":       "file:1:1 (0): invalid encoding",
	// non-terminated "quoted" tokens
	"{":       "file:1:1 (0): rule CodeBlock: code block not terminated",
	"\n{{}":   "file:2:1 (1): rule CodeBlock: code block not terminated",
	`a = "b`:  "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = `b":  "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = 'b":  "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = [b`:  "file:1:5 (4): rule CharClassMatcher: character class not terminated",
	`a = "\"`: "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = '\'`: "file:1:5 (4): rule StringLiteral: string literal not terminated",
	// invalid escapes
	`a ← [\pA]`:    "file:1:8 (9): rule UnicodeClassEscape: invalid Unicode class escape",
	`a ← [\p{WW}]`: "file:1:9 (10): rule UnicodeClass: invalid Unicode class escape",
	`a = '\"'`:     "file:1:7 (6): rule SingleStringEscape: invalid escape character",
	`a = "\'"`:     "file:1:7 (6): rule DoubleStringEscape: invalid escape character",
	`a = [\']`:     "file:1:7 (6): rule CharClassEscape: invalid escape character",
	`a = '\xz'`:    "file:1:7 (6): rule HexEscape: invalid hexadecimal escape",
	`a = '\0z'`:    "file:1:7 (6): rule OctalEscape: invalid octal escape",
	`a = '\uz'`:    "file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape",
	`a = '\Uz'`:    "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",
	/*
		// escapes followed by newline
		"a = '\\\n":   "",
		"a = '\\x\n":  "",
		"a = '\\0\n":  "",
		"a = '\\u\n":  "",
		"a = '\\U\n":  "",
		"a = \"\\\n":  "",
		"a = \"\\x\n": "",
		"a = \"\\0\n": "",
		"a = \"\\u\n": "",
		"a = \"\\U\n": "",
		"a = [\\\n":   "",
		"a = [\\x\n":  "",
		"a = [\\0\n":  "",
		"a = [\\u\n":  "",
		"a = [\\U\n":  "",
		"a = [\\p\n":  "",
		"a = [\\p{\n": "",
		// escapes followed by EOF
		"a = '\\":   "",
		"a = '\\x":  "",
		"a = '\\0":  "",
		"a = '\\u":  "",
		"a = '\\U":  "",
		"a = \"\\":  "",
		"a = \"\\x": "",
		"a = \"\\0": "",
		"a = \"\\u": "",
		"a = \"\\U": "",
		"a = [\\":   "",
		"a = [\\x":  "",
		"a = [\\0":  "",
		"a = [\\u":  "",
		"a = [\\U":  "",
		"a = [\\p":  "",
		"a = [\\p{": "",
	*/
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
