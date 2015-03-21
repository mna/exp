package bootstrap

import (
	"strings"
	"testing"
)

var parseValidCases = []string{
	"package a\n",
	"package a\n{code}",
	"package a\nR <- 'c'",
	"package a\n;\n\nR <- 'c'\n\n",
	`package a

A = ident:B / C+ / D?;`,
	`package a

{ code }

R "name" <- "abc"i
R2 = 'd'i
R3 = ( R2+ ![;] )`,
}

var parseExpRes = []string{
	`1:1 (0): *ast.Grammar{Package: 1:1 (0): *ast.Package{Name: 1:9 (8): *ast.Identifier{Val: "a"}}, Init: <nil>, Rules: [
]}`,
	`1:1 (0): *ast.Grammar{Package: 1:1 (0): *ast.Package{Name: 1:9 (8): *ast.Identifier{Val: "a"}}, Init: 2:1 (10): *ast.CodeBlock{Val: "{code}"}, Rules: [
]}`,
	`1:1 (0): *ast.Grammar{Package: 1:1 (0): *ast.Package{Name: 1:9 (8): *ast.Identifier{Val: "a"}}, Init: <nil>, Rules: [
2:1 (10): *ast.Rule{Name: 2:1 (10): *ast.Identifier{Val: "R"}, DisplayName: <nil>, Expr: 2:6 (15): *ast.LitMatcher{Val: "c", IgnoreCase: false}},
]}`,
	`1:1 (0): *ast.Grammar{Package: 1:1 (0): *ast.Package{Name: 1:9 (8): *ast.Identifier{Val: "a"}}, Init: <nil>, Rules: [
4:1 (13): *ast.Rule{Name: 4:1 (13): *ast.Identifier{Val: "R"}, DisplayName: <nil>, Expr: 4:6 (18): *ast.LitMatcher{Val: "c", IgnoreCase: false}},
]}`,
	`1:1 (0): *ast.Grammar{Package: 1:1 (0): *ast.Package{Name: 1:9 (8): *ast.Identifier{Val: "a"}}, Init: <nil>, Rules: [
3:1 (11): *ast.Rule{Name: 3:1 (11): *ast.Identifier{Val: "A"}, DisplayName: <nil>, Expr: 3:5 (15): *ast.ChoiceExpr{Alternatives: [
3:5 (15): *ast.LabeledExpr{Label: 3:5 (15): *ast.Identifier{Val: "ident"}, Expr: 3:11 (21): *ast.RuleRefExpr{Name: 3:11 (21): *ast.Identifier{Val: "B"}}},
3:15 (25): *ast.OneOrMoreExpr{Expr: 3:15 (25): *ast.RuleRefExpr{Name: 3:15 (25): *ast.Identifier{Val: "C"}}},
3:20 (30): *ast.ZeroOrOneExpr{Expr: 3:20 (30): *ast.RuleRefExpr{Name: 3:20 (30): *ast.Identifier{Val: "D"}}},
]}},
]}`,
	`1:1 (0): *ast.Grammar{Package: 1:1 (0): *ast.Package{Name: 1:9 (8): *ast.Identifier{Val: "a"}}, Init: 3:1 (11): *ast.CodeBlock{Val: "{ code }"}, Rules: [
5:1 (21): *ast.Rule{Name: 5:1 (21): *ast.Identifier{Val: "R"}, DisplayName: 5:3 (23): *ast.StringLit{Val: "name"}, Expr: 5:13 (33): *ast.LitMatcher{Val: "abc", IgnoreCase: true}},
6:1 (40): *ast.Rule{Name: 6:1 (40): *ast.Identifier{Val: "R2"}, DisplayName: <nil>, Expr: 6:6 (45): *ast.LitMatcher{Val: "d", IgnoreCase: true}},
7:1 (50): *ast.Rule{Name: 7:1 (50): *ast.Identifier{Val: "R3"}, DisplayName: <nil>, Expr: 7:8 (57): *ast.SeqExpr{Exprs: [
7:8 (57): *ast.OneOrMoreExpr{Expr: 7:8 (57): *ast.RuleRefExpr{Name: 7:8 (57): *ast.Identifier{Val: "R2"}}},
7:12 (61): *ast.NotExpr{Expr: 7:13 (62): *ast.CharClassMatcher{Val: "[;]", IgnoreCase: false, Inverted: false}},
]}},
]}`,
}

func TestParseValid(t *testing.T) {
	p := NewParser()
	for i, c := range parseValidCases {
		g, err := p.Parse("", strings.NewReader(c))
		if err != nil {
			t.Errorf("%d: got error %v", i, err)
			continue
		}

		want := parseExpRes[i]
		got := g.String()
		if want != got {
			t.Errorf("%d: want \n%s\n, got \n%s\n", i, want, got)
		}
	}
}

var parseInvalidCases = []string{
	"",
	"a",
	"package",
	"package a",
	`package a
	Rule
`,
	`package a
R "a"i`,
	`package a
R = )`,
}

var parseExpErrs = [][]string{
	{"1:0 (0): expected keyword, got eof", "1:0 (0): no grammar"},
	{"1:1 (0): expected keyword, got ident", "1:1 (0): no grammar"},
	{"1:7 (6): expected ident, got eof", "1:7 (6): no grammar"},
	{"1:9 (8): expected any of [eol semicolon], got eof", "1:9 (8): no grammar"},
	{"3:0 (15): expected ruledef, got eol"},
	{"2:3 (12): invalid suffix 'i'"},
	{"2:5 (14): no expression in sequence", "2:5 (14): no expression in choice", "2:5 (14): missing expression"},
}

func TestParseInvalid(t *testing.T) {
	p := NewParser()
	for i, c := range parseInvalidCases {
		_, err := p.Parse("", strings.NewReader(c))
		el := *(err.(*errList))
		if len(el) != len(parseExpErrs[i]) {
			t.Errorf("%d: want %d errors, got %d", i, len(parseExpErrs[i]), len(el))
			continue
		}
		for j, err := range el {
			want := parseExpErrs[i][j]
			got := err.Error()
			if want != got {
				t.Errorf("%d: error %d: want %q, got %q", i, j, want, got)
			}
		}
	}
}
