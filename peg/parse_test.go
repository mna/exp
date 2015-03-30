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
	"\xfe",
	`a = '\"'`,
	`a = "\'"`,
	`a = '\xzz`,
	`a = '\091`,
	`a = "b`,
	"a = `b",
	"a = 'b",
	`a = '\uA'`,
	`a = '\UA012'`,
	`a = [b`,
	`a = b /`,
	`a = +`,
	`a = *`,
	`a = ?`,
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
	"file:1:1 (0): no match found",
	"file:1:8 (9): rule UnicodeClassEscape: invalid Unicode class escape",
	"file:1:1 (0): invalid encoding",
	"file:1:7 (6): invalid escape character",
	"file:1:7 (6): invalid escape character",
	"file:1:7 (6): invalid hexadecimal escape",
	"file:1:7 (6): invalid octal escape",
	"file:1:7 (6): string literal not terminated",
	"file:1:7 (6): string literal not terminated",
	"file:1:7 (6): character literal not terminated",
	"file:1:7 (6): invalid Unicode escape",
	"file:1:7 (6): invalid Unicode escape",
	"file:1:7 (6): character class not terminated",
	"file:1:7 (6): invalid choice expression",
	"file:1:1 (0): no match found",
	"file:1:1 (0): no match found",
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
