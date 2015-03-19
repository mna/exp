package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/exp/peg/bootstrap"
)

var calcGrammar = `package test

start = additive eof
additive = left:multiplicative "+" space right:additive / multiplicative
multiplicative = left:primary "*" space right:multiplicative {code} / primary
primary = integer / "(" space additive:additive ")" space
integer = digits:[0123456789]+ space {integer}
space = ' '*
eof = !.
`

var src = `9 + 5 * (1+2)`

func TestCalcGrammar(t *testing.T) {
	// parse the grammar
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(calcGrammar))
	if err != nil {
		t.Fatal(err)
	}

	for _, rule := range g.Rules {
		fmt.Printf("%s\n", rule)
	}

	// res, err := parseUsingAST("", strings.NewReader(src), g)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("%#v\n", res)
}
