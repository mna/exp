package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/exp/peg/ast"
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
	"{}{}":       "file:1:1 (0): no match found",

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

var validParseCases = map[string]*ast.Grammar{
	"a = b": &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
		},
	},
	"a ← b\nc=d \n e <- f \ng\u27f5h": &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
			{
				Name: ast.NewIdentifier(ast.Pos{}, "c"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "d")},
			},
			{
				Name: ast.NewIdentifier(ast.Pos{}, "e"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "f")},
			},
			{
				Name: ast.NewIdentifier(ast.Pos{}, "g"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "h")},
			},
		},
	},
	`a "A"← b`: &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name:        ast.NewIdentifier(ast.Pos{}, "a"),
				DisplayName: ast.NewStringLit(ast.Pos{}, `"A"`),
				Expr:        &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
		},
	},
	"{ init \n}\na 'A'← b": &ast.Grammar{
		Init: ast.NewCodeBlock(ast.Pos{}, "{ init \n}"),
		Rules: []*ast.Rule{
			{
				Name:        ast.NewIdentifier(ast.Pos{}, "a"),
				DisplayName: ast.NewStringLit(ast.Pos{}, `'A'`),
				Expr:        &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
		},
	},
	"a\n<-\nb": &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
		},
	},
	"a\n<-\nb\nc": &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.SeqExpr{
					Exprs: []ast.Expression{
						&ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
						&ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "c")},
					},
				},
			},
		},
	},
	"a\n<-\nb\nc\n=\nd": &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
			{
				Name: ast.NewIdentifier(ast.Pos{}, "c"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "d")},
			},
		},
	},
	"a\n<-\nb\nc\n'C'\n=\nd": &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "b")},
			},
			{
				Name:        ast.NewIdentifier(ast.Pos{}, "c"),
				DisplayName: ast.NewStringLit(ast.Pos{}, `'C'`),
				Expr:        &ast.RuleRefExpr{Name: ast.NewIdentifier(ast.Pos{}, "d")},
			},
		},
	},
	`a = [a-def]`: &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.CharClassMatcher{
					Chars:  []rune{'e', 'f'},
					Ranges: []rune{'a', 'd'},
				},
			},
		},
	},
	`a = [abc-f]`: &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.CharClassMatcher{
					Chars:  []rune{'a', 'b'},
					Ranges: []rune{'c', 'f'},
				},
			},
		},
	},
	`a = [abc-fg]`: &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.CharClassMatcher{
					Chars:  []rune{'a', 'b', 'g'},
					Ranges: []rune{'c', 'f'},
				},
			},
		},
	},
	`a = [abc-fgh-l]`: &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.CharClassMatcher{
					Chars:  []rune{'a', 'b', 'g'},
					Ranges: []rune{'c', 'f', 'h', 'l'},
				},
			},
		},
	},
	`a = [\x00-\xabc]`: &ast.Grammar{
		Rules: []*ast.Rule{
			{
				Name: ast.NewIdentifier(ast.Pos{}, "a"),
				Expr: &ast.CharClassMatcher{
					Chars:  []rune{'c'},
					Ranges: []rune{'\x00', '\xab'},
				},
			},
		},
	},
}

func TestValidParseCases(t *testing.T) {
	for tc, exp := range validParseCases {
		got, err := Parse("", strings.NewReader(tc))
		if err != nil {
			t.Errorf("%q: got error %v", tc, err)
			continue
		}
		gotg, ok := got.(*ast.Grammar)
		if !ok {
			t.Errorf("%q: want grammar type %T, got %T", tc, exp, got)
			continue
		}
		compareGrammars(t, tc, exp, gotg)
	}
}

func compareGrammars(t *testing.T, src string, exp, got *ast.Grammar) bool {
	if (exp.Init != nil) != (got.Init != nil) {
		t.Errorf("%q: want Init? %t, got %t", src, exp.Init != nil, got.Init != nil)
		return false
	}
	if exp.Init != nil {
		if exp.Init.Val != got.Init.Val {
			t.Errorf("%q: want Init %q, got %q", src, exp.Init.Val, got.Init.Val)
			return false
		}
	}

	rn, rm := len(exp.Rules), len(got.Rules)
	if rn != rm {
		t.Errorf("%q: want %d rules, got %d", src, rn, rm)
		return false
	}

	for i, r := range got.Rules {
		if !compareRule(t, src+": "+exp.Rules[i].Name.Val, exp.Rules[i], r) {
			return false
		}
	}

	return true
}

func compareRule(t *testing.T, prefix string, exp, got *ast.Rule) bool {
	if exp.Name.Val != got.Name.Val {
		t.Errorf("%q: want rule name %q, got %q", prefix, exp.Name.Val, got.Name.Val)
		return false
	}
	if (exp.DisplayName != nil) != (got.DisplayName != nil) {
		t.Errorf("%q: want DisplayName? %t, got %t", prefix, exp.DisplayName != nil, got.DisplayName != nil)
		return false
	}
	if exp.DisplayName != nil {
		if exp.DisplayName.Val != got.DisplayName.Val {
			t.Errorf("%q: want DisplayName %q, got %q", prefix, exp.DisplayName.Val, got.DisplayName.Val)
			return false
		}
	}
	return compareExpr(t, prefix, 0, exp.Expr, got.Expr)
}

func compareExpr(t *testing.T, prefix string, ix int, exp, got ast.Expression) bool {
	ixPrefix := prefix + " (" + strconv.Itoa(ix) + ")"

	switch exp := exp.(type) {
	case *ast.ActionExpr:
		got, ok := got.(*ast.ActionExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if (exp.Code != nil) != (got.Code != nil) {
			t.Errorf("%q: want Code?: %t, got %t", ixPrefix, exp.Code != nil, got.Code != nil)
			return false
		}
		if exp.Code != nil {
			if exp.Code.Val != got.Code.Val {
				t.Errorf("%q: want code %q, got %q", ixPrefix, exp.Code.Val, got.Code.Val)
				return false
			}
		}
		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	case *ast.AndCodeExpr:
		got, ok := got.(*ast.AndCodeExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if (exp.Code != nil) != (got.Code != nil) {
			t.Errorf("%q: want Code?: %t, got %t", ixPrefix, exp.Code != nil, got.Code != nil)
			return false
		}
		if exp.Code != nil {
			if exp.Code.Val != got.Code.Val {
				t.Errorf("%q: want code %q, got %q", ixPrefix, exp.Code.Val, got.Code.Val)
				return false
			}
		}

	case *ast.AndExpr:
		got, ok := got.(*ast.AndExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	case *ast.AnyMatcher:
		got, ok := got.(*ast.AnyMatcher)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		// for completion's sake...
		if exp.Val != got.Val {
			t.Errorf("%q: want value %q, got %q", ixPrefix, exp.Val, got.Val)
		}

	case *ast.CharClassMatcher:
		got, ok := got.(*ast.CharClassMatcher)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if exp.IgnoreCase != got.IgnoreCase {
			t.Errorf("%q: want IgnoreCase %t, got %t", ixPrefix, exp.IgnoreCase, got.IgnoreCase)
			return false
		}
		if exp.Inverted != got.Inverted {
			t.Errorf("%q: want Inverted %t, got %t", ixPrefix, exp.Inverted, got.Inverted)
			return false
		}

		ne, ng := len(exp.Chars), len(got.Chars)
		if ne != ng {
			t.Errorf("%q: want %d Chars, got %d (%v)", ixPrefix, ne, ng, got.Chars)
			return false
		}
		for i, r := range exp.Chars {
			if r != got.Chars[i] {
				t.Errorf("%q: want Chars[%d] %#U, got %#U", ixPrefix, i, r, got.Chars[i])
				return false
			}
		}

		ne, ng = len(exp.Ranges), len(got.Ranges)
		if ne != ng {
			t.Errorf("%q: want %d Ranges, got %d", ixPrefix, ne, ng)
			return false
		}
		for i, r := range exp.Ranges {
			if r != got.Ranges[i] {
				t.Errorf("%q: want Ranges[%d] %#U, got %#U", ixPrefix, i, r, got.Ranges[i])
				return false
			}
		}

		ne, ng = len(exp.UnicodeClasses), len(got.UnicodeClasses)
		if ne != ng {
			t.Errorf("%q: want %d UnicodeClasses, got %d", ixPrefix, ne, ng)
			return false
		}
		for i, s := range exp.UnicodeClasses {
			if s != got.UnicodeClasses[i] {
				t.Errorf("%q: want UnicodeClasses[%d] %q, got %q", ixPrefix, i, s, got.UnicodeClasses[i])
				return false
			}
		}

	case *ast.ChoiceExpr:
		got, ok := got.(*ast.ChoiceExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		ne, ng := len(exp.Alternatives), len(got.Alternatives)
		if ne != ng {
			t.Errorf("%q: want %d Alternatives, got %d", ixPrefix, ne, ng)
			return false
		}

		for i, alt := range exp.Alternatives {
			if !compareExpr(t, prefix, ix+1, alt, got.Alternatives[i]) {
				return false
			}
		}

	case *ast.LabeledExpr:
		got, ok := got.(*ast.LabeledExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if (exp.Label != nil) != (got.Label != nil) {
			t.Errorf("%q: want Label?: %t, got %t", ixPrefix, exp.Label != nil, got.Label != nil)
			return false
		}
		if exp.Label != nil {
			if exp.Label.Val != got.Label.Val {
				t.Errorf("%q: want label %q, got %q", ixPrefix, exp.Label.Val, got.Label.Val)
				return false
			}
		}

		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	case *ast.LitMatcher:
		got, ok := got.(*ast.LitMatcher)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if exp.IgnoreCase != got.IgnoreCase {
			t.Errorf("%q: want IgnoreCase %t, got %t", ixPrefix, exp.IgnoreCase, got.IgnoreCase)
			return false
		}
		if exp.Val != got.Val {
			t.Errorf("%q: want value %q, got %q", ixPrefix, exp.Val, got.Val)
			return false
		}

	case *ast.NotCodeExpr:
		got, ok := got.(*ast.NotCodeExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if (exp.Code != nil) != (got.Code != nil) {
			t.Errorf("%q: want Code?: %t, got %t", ixPrefix, exp.Code != nil, got.Code != nil)
			return false
		}
		if exp.Code != nil {
			if exp.Code.Val != got.Code.Val {
				t.Errorf("%q: want code %q, got %q", ixPrefix, exp.Code.Val, got.Code.Val)
				return false
			}
		}

	case *ast.NotExpr:
		got, ok := got.(*ast.NotExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	case *ast.OneOrMoreExpr:
		got, ok := got.(*ast.OneOrMoreExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	case *ast.RuleRefExpr:
		got, ok := got.(*ast.RuleRefExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		if (exp.Name != nil) != (got.Name != nil) {
			t.Errorf("%q: want Name?: %t, got %t", ixPrefix, exp.Name != nil, got.Name != nil)
			return false
		}
		if exp.Name != nil {
			if exp.Name.Val != got.Name.Val {
				t.Errorf("%q: want name %q, got %q", ixPrefix, exp.Name.Val, got.Name.Val)
				return false
			}
		}

	case *ast.SeqExpr:
		got, ok := got.(*ast.SeqExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		ne, ng := len(exp.Exprs), len(got.Exprs)
		if ne != ng {
			t.Errorf("%q: want %d Exprs, got %d", ixPrefix, ne, ng)
			return false
		}

		for i, expr := range exp.Exprs {
			if !compareExpr(t, prefix, ix+1, expr, got.Exprs[i]) {
				return false
			}
		}

	case *ast.ZeroOrMoreExpr:
		got, ok := got.(*ast.ZeroOrMoreExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	case *ast.ZeroOrOneExpr:
		got, ok := got.(*ast.ZeroOrOneExpr)
		if !ok {
			t.Errorf("%q: want expression type %T, got %T", ixPrefix, exp, got)
			return false
		}
		return compareExpr(t, prefix, ix+1, exp.Expr, got.Expr)

	default:
		t.Fatalf("unexpected expression type %T", exp)
	}
	return true
}
