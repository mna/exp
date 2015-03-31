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

	// non-terminated, empty, EOF "quoted" tokens
	"{":     "file:1:1 (0): rule CodeBlock: code block not terminated",
	"\n{":   "file:2:1 (1): rule CodeBlock: code block not terminated",
	`a = "`: "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = `": "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = '": "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = [`: "file:1:5 (4): rule CharClassMatcher: character class not terminated",
	`a = [\p{]`: `file:1:9 (8): rule UnicodeClass: invalid Unicode class escape
file:1:8 (7): rule UnicodeClassEscape: Unicode class not terminated
file:1:5 (4): rule CharClassMatcher: character class not terminated`,

	// non-terminated, empty, EOL "quoted" tokens
	"{\n":          "file:1:1 (0): rule CodeBlock: code block not terminated",
	"\n{\n":        "file:2:1 (1): rule CodeBlock: code block not terminated",
	"a = \"\n":     "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = `\n":      "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = '\n":      "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = [\n":      "file:1:5 (4): rule CharClassMatcher: character class not terminated",
	"a = [\\p{\n]": `file:1:5 (4): rule CharClassMatcher: character class not terminated`,

	// non-terminated quoted tokens with escaped closing char
	`a = "\"`: "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = '\'`: "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = [\]`: "file:1:5 (4): rule CharClassMatcher: character class not terminated",

	// non-terminated, non-empty, EOF "quoted" tokens
	"{a":     "file:1:1 (0): rule CodeBlock: code block not terminated",
	"\n{{}":  "file:2:1 (1): rule CodeBlock: code block not terminated",
	`a = "b`: "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = `b": "file:1:5 (4): rule StringLiteral: string literal not terminated",
	"a = 'b": "file:1:5 (4): rule StringLiteral: string literal not terminated",
	`a = [b`: "file:1:5 (4): rule CharClassMatcher: character class not terminated",
	`a = [\p{W]`: `file:1:9 (8): rule UnicodeClass: invalid Unicode class escape
file:1:8 (7): rule UnicodeClassEscape: Unicode class not terminated
file:1:5 (4): rule CharClassMatcher: character class not terminated`,

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

	// escapes followed by newline
	"a = '\\\n": `file:2:0 (6): rule SingleStringEscape: invalid escape character
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\x\n": `file:1:7 (6): rule HexEscape: invalid hexadecimal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\0\n": `file:1:7 (6): rule OctalEscape: invalid octal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\u\n": `file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\U\n": `file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\\n": `file:2:0 (6): rule DoubleStringEscape: invalid escape character
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\x\n": `file:1:7 (6): rule HexEscape: invalid hexadecimal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\0\n": `file:1:7 (6): rule OctalEscape: invalid octal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\u\n": `file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\U\n": `file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = [\\\n": `file:2:0 (6): rule CharClassEscape: invalid escape character
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\x\n": `file:1:7 (6): rule HexEscape: invalid hexadecimal escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\0\n": `file:1:7 (6): rule OctalEscape: invalid octal escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\u\n": `file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\U\n": `file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\p\n": `file:2:0 (7): rule UnicodeClassEscape: invalid Unicode class escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\p{\n": `file:1:5 (4): rule CharClassMatcher: character class not terminated`,

	// escapes followed by EOF
	"a = '\\": `file:1:7 (6): rule SingleStringEscape: invalid escape character
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\x": `file:1:7 (6): rule HexEscape: invalid hexadecimal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\0": `file:1:7 (6): rule OctalEscape: invalid octal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\u": `file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = '\\U": `file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\": `file:1:7 (6): rule DoubleStringEscape: invalid escape character
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\x": `file:1:7 (6): rule HexEscape: invalid hexadecimal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\0": `file:1:7 (6): rule OctalEscape: invalid octal escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\u": `file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = \"\\U": `file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule StringLiteral: string literal not terminated`,
	"a = [\\": `file:1:7 (6): rule CharClassEscape: invalid escape character
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\x": `file:1:7 (6): rule HexEscape: invalid hexadecimal escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\0": `file:1:7 (6): rule OctalEscape: invalid octal escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\u": `file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\U": `file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\p": `file:1:8 (7): rule UnicodeClassEscape: invalid Unicode class escape
file:1:5 (4): rule CharClassMatcher: character class not terminated`,
	"a = [\\p{": `file:1:5 (4): rule CharClassMatcher: character class not terminated`,

	// multi-char escapes, fail after 2 chars
	`a = '\x0z'`: "file:1:7 (6): rule HexEscape: invalid hexadecimal escape",
	`a = '\00z'`: "file:1:7 (6): rule OctalEscape: invalid octal escape",
	`a = '\u0z'`: "file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape",
	`a = '\U0z'`: "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",
	// multi-char escapes, fail after 3 chars
	`a = '\u00z'`: "file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape",
	`a = '\U00z'`: "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",
	// multi-char escapes, fail after 4 chars
	`a = '\u000z'`: "file:1:7 (6): rule ShortUnicodeEscape: invalid Unicode escape",
	`a = '\U000z'`: "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",
	// multi-char escapes, fail after 5 chars
	`a = '\U0000z'`: "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",
	// multi-char escapes, fail after 6 chars
	`a = '\U00000z'`: "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",
	// multi-char escapes, fail after 7 chars
	`a = '\U000000z'`: "file:1:7 (6): rule LongUnicodeEscape: invalid Unicode escape",

	// combine escape errors
	`a = "\a\b\c\t\n\r\xab\xz\ux"`: `file:1:11 (10): rule DoubleStringEscape: invalid escape character
file:1:23 (22): rule HexEscape: invalid hexadecimal escape
file:1:26 (25): rule ShortUnicodeEscape: invalid Unicode escape`,
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
