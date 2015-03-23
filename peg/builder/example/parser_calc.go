package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
This is a hand-generated example of what the parser generator would do
for the following grammar. Positions have been left off, but the rest is
as close as possible to what would/should be generated.

// Options: -imports=strconv,fmt,strings -node-type=int

package main

{
func init() {
	fmt.Println("this is in the grammar's Init code.")
}

func main() {
	res, err := Parse("", strings.NewReader("2 + 3 * (5 +1)"))
	fmt.Println("got ", res, err)
}
}

start = result:additive eof {
	fmt.Println("result: ", result)
	return result, nil
}
additive = left:multiplicative "+" space right:additive { return left + right, nil }
	/ multiplicative
multiplicative = left:primary "*" space right:multiplicative { return left * right, nil }
	/ primary
primary = integer / "(" space additive:additive ")" space { return additive, nil }
integer = digits:[0-9]+ space {
	return strconv.Atoi(digits)
}
space = ' '*
eof = !.
*/

func init() {
	fmt.Println("this is in the grammar's Init code.")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: <cmd> EXPR")
	}
	res, err := Parse("", strings.NewReader(os.Args[1]))
	fmt.Println("got ", res, err)
}

var Grammar = &grammar{
	rules: []*rule{
		{
			name: "start",
			expr: &actionExpr{
				run: (*parser).callonstart_0,
				expr: &seqExpr{
					exprs: []interface{}{
						&labeledExpr{
							label: "result",
							expr: &ruleRefExpr{
								name: "additive",
							},
						},
						&ruleRefExpr{
							name: "eof",
						},
					},
				},
			},
		},
		{
			name: "additive",
			expr: &choiceExpr{
				alternatives: []interface{}{
					&actionExpr{
						run: (*parser).callonadditive_1,
						expr: &seqExpr{
							exprs: []interface{}{
								&labeledExpr{
									label: "left",
									expr: &ruleRefExpr{
										name: "multiplicative",
									},
								},
								&litMatcher{
									val: "+",
								},
								&ruleRefExpr{
									name: "space",
								},
								&labeledExpr{
									label: "right",
									expr: &ruleRefExpr{
										name: "additive",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						name: "multiplicative",
					},
				},
			},
		},
		{
			name: "multiplicative",
			expr: &choiceExpr{
				alternatives: []interface{}{
					&actionExpr{
						run: (*parser).callonmultiplicative_1,
						expr: &seqExpr{
							exprs: []interface{}{
								&labeledExpr{
									label: "left",
									expr: &ruleRefExpr{
										name: "primary",
									},
								},
								&litMatcher{
									val: "*",
								},
								&ruleRefExpr{
									name: "space",
								},
								&labeledExpr{
									label: "right",
									expr: &ruleRefExpr{
										name: "multiplicative",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						name: "primary",
					},
				},
			},
		},
		{
			name: "primary",
			expr: &choiceExpr{
				alternatives: []interface{}{
					&ruleRefExpr{
						name: "integer",
					},
					&actionExpr{
						run: (*parser).callonprimary_2,
						expr: &seqExpr{
							exprs: []interface{}{
								&litMatcher{
									val: "(",
								},
								&ruleRefExpr{
									name: "space",
								},
								&labeledExpr{
									label: "additive",
									expr: &ruleRefExpr{
										name: "additive",
									},
								},
								&litMatcher{
									val: ")",
								},
								&ruleRefExpr{
									name: "space",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "integer",
			expr: &actionExpr{
				run: (*parser).calloninteger_0,
				expr: &seqExpr{
					exprs: []interface{}{
						&labeledExpr{
							expr: &oneOrMoreExpr{
								expr: &charClassMatcher{
									val:    "[0-9]",
									ranges: []rune{'0', '9'},
								},
							},
							label: "digits",
						},
						&ruleRefExpr{
							name: "space",
						},
					},
				},
			},
		},
		{
			name: "space",
			expr: &zeroOrMoreExpr{
				expr: &litMatcher{
					val: " ",
				},
			},
		},
		{
			name: "eof",
			expr: &notExpr{
				expr: &anyMatcher{},
			},
		},
	},
}

func (c *current) onstart_0(result int) (int, error) {
	fmt.Println("onstart_0: ", result)
	return result, nil
}

func (c *current) onadditive_1(left, right int) (int, error) {
	fmt.Println("onadditive_1: ", left, right)
	return left + right, nil
}

func (c *current) onmultiplicative_1(left, right int) (int, error) {
	fmt.Println("onmultiplicative_1: ", left, right)
	return left * right, nil
}

func (c *current) onprimary_2(additive int) (int, error) {
	fmt.Println("onprimary_2: ", additive)
	return additive, nil
}

// type inferred to string since the label is on a litMatcher
func (c *current) oninteger_0(digits string) (int, error) {
	fmt.Println("oninteger_0: ", digits)
	return strconv.Atoi(digits)
}

func (p *parser) callonstart_0() (int, error) {
	stack := p.vstack[len(p.vstack)-1]
	return p.cur.onstart_0(stack["result"].(int))
}

func (p *parser) callonadditive_1() (int, error) {
	stack := p.vstack[len(p.vstack)-1]
	return p.cur.onadditive_1(stack["left"].(int), stack["right"].(int))
}

func (p *parser) callonmultiplicative_1() (int, error) {
	stack := p.vstack[len(p.vstack)-1]
	return p.cur.onmultiplicative_1(stack["left"].(int), stack["right"].(int))
}

func (p *parser) callonprimary_2() (int, error) {
	stack := p.vstack[len(p.vstack)-1]
	return p.cur.onprimary_2(stack["additive"].(int))
}

func (p *parser) calloninteger_0() (int, error) {
	stack := p.vstack[len(p.vstack)-1]
	val := stack["digits"].([]interface{})
	var buf bytes.Buffer
	for _, v := range val {
		buf.WriteString(v.(string))
	}
	return p.cur.oninteger_0(buf.String())
}

// --------------------
// Start of static code
// --------------------

var (
	ErrNoRule          = errors.New("grammar has no rule")
	ErrInvalidEncoding = errors.New("invalid encoding")
	ErrNoMatch         = errors.New("no match found")
)

func ParseFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(filename, f)
}

func Parse(filename string, r io.Reader) (interface{}, error) {
	return parse(filename, r, Grammar)
}

type position struct {
	line, col, offset int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (int, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e *errList) err() error {
	if len(*e) == 0 {
		return nil
	}
	return e
}

func (e *errList) Error() string {
	switch len(*e) {
	case 0:
		return ""
	case 1:
		return (*e)[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range *e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

func parse(filename string, r io.Reader, g *grammar) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	p := &parser{filename: filename, errs: new(errList), data: b, pt: savepoint{position: position{line: 1}}}
	return p.parse(g)
}

type savepoint struct {
	position
	rn rune
	w  int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	rules  map[string]*rule
	vstack []map[string]interface{}
	rstack []*rule
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.errs.add(ErrInvalidEncoding)
		}
	}
}

func (p *parser) save() savepoint {
	return p.pt
}

func (p *parser) restore(pt savepoint) {
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

func (p *parser) slice(start, end position) []byte {
	return p.data[start.offset:end.offset]
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		return nil, ErrNoRule
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	// panic can be used in action code to stop parsing immediately
	// and return the panic as an error.
	defer func() {
		if e := recover(); e != nil {
			val = nil
			switch e := e.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		err := p.errs.err()
		if err == nil {
			// make sure this doesn't go out silently
			err = ErrNoMatch
		}
		return nil, err
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	// TODO : build error messages with references to the current rule
	p.rstack = append(p.rstack, rule)
	val, ok := p.parseExpr(rule.expr)
	p.rstack = p.rstack[:len(p.rstack)-1]
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	switch expr := expr.(type) {
	case *actionExpr:
		return p.parseActionExpr(expr)
	case *andCodeExpr:
		return p.parseAndCodeExpr(expr)
	case *andExpr:
		return p.parseAndExpr(expr)
	case *anyMatcher:
		return p.parseAnyMatcher(expr)
	case *charClassMatcher:
		return p.parseCharClassMatcher(expr)
	case *choiceExpr:
		return p.parseChoiceExpr(expr)
	case *labeledExpr:
		return p.parseLabeledExpr(expr)
	case *litMatcher:
		return p.parseLitMatcher(expr)
	case *notCodeExpr:
		return p.parseNotCodeExpr(expr)
	case *notExpr:
		return p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		return p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		return p.parseRuleRefExpr(expr)
	case *seqExpr:
		return p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		return p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		return p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	p.vstack = append(p.vstack, make(map[string]interface{}))
	start := p.save()
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.slice(start.position, p.save().position)
		actVal, err := act.run(p)
		if err != nil {
			p.errs.add(err) // TODO : transform, or use directly?
		}
		val = actVal
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	ok, err := and.run(p)
	if err != nil {
		p.errs.add(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	pt := p.save()
	_, ok := p.parseExpr(and.expr)
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.pt.rn != utf8.RuneError {
		p.read()
		return string(p.pt.rn), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	cur := p.pt.rn
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return string(cur), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return string(cur), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return string(cur), true
		}
	}

	if chr.inverted {
		p.read()
		return string(cur), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	for _, alt := range ch.alternatives {
		val, ok := p.parseExpr(alt)
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	val, ok := p.parseExpr(lab.expr)
	if ok && lab.label != "" && len(p.vstack) > 0 {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	var buf bytes.Buffer
	pt := p.save()
	for _, want := range lit.val {
		cur := p.pt.rn
		buf.WriteRune(cur)
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(pt)
			return nil, false
		}
		p.read()
	}
	return buf.String(), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	ok, err := not.run(p)
	if err != nil {
		p.errs.add(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	pt := p.save()
	_, ok := p.parseExpr(not.expr)
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	var vals []interface{}

	for {
		val, ok := p.parseExpr(expr.expr)
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.errs.add(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	var vals []interface{}

	pt := p.save()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	var vals []interface{}

	for {
		val, ok := p.parseExpr(expr.expr)
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	val, _ := p.parseExpr(expr.expr)
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// TODO : should be caught at the scan/parse step
	return &unicode.RangeTable{} // empty range
}
