package parser

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/exp/peg/bootstrap"
)

var grammar = `package test

start = additive
additive = left:multiplicative "+" space right:additive / multiplicative
multiplicative = left:primary "*" space right:multiplicative / primary
primary = integer / "(" space additive:additive ")" space
integer "integer" = digits:[0123456789]+ space
space = ' '*`

var src = `9 + 5 * (1+2)`

func TestCalcGrammar(t *testing.T) {
	// parse the grammar
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	if err != nil {
		t.Fatal(err)
	}

	res, err := parseUsingAST("", strings.NewReader(src), g)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", res)
}
