// Command calculator is a small PEG-generated parser that computes
// simple math using integers, additions, subtractions and multiplications.
//
// Example usage: $ calculator "3 + (2 - 5 * 12)"
//
// Inspired by pegjs arithmetic example:
// https://github.com/pegjs/pegjs/blob/master/examples/arithmetics.pegjs
//
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

var ops = map[string]func(int, int) int{
	"+": func(l, r int) int {
		return l + r
	},
	"-": func(l, r int) int {
		return l - r
	},
	"*": func(l, r int) int {
		return l * r
	},
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: calculator 'EXPR'")
	}
	got, err := Parse("", strings.NewReader(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=", got)
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func eval(first, rest interface{}) int {
	l := first.(int)
	restSl := toIfaceSlice(rest)
	for _, v := range restSl {
		restExpr := toIfaceSlice(v)
		r := restExpr[3].(int)
		op := restExpr[1].(string)
		l = ops[op](l, r)
	}
	return l
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 55, col: 1, offset: 1182},
			expr: &actionExpr{
				pos: position{line: 55, col: 10, offset: 1191},
				run: (*parser).callonInput1,
				expr: &seqExpr{
					pos: position{line: 55, col: 10, offset: 1191},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 55, col: 10, offset: 1191},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 55, col: 15, offset: 1196},
								name: "Expr",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 20, offset: 1201},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 59, col: 1, offset: 1231},
			expr: &actionExpr{
				pos: position{line: 59, col: 9, offset: 1239},
				run: (*parser).callonExpr1,
				expr: &seqExpr{
					pos: position{line: 59, col: 9, offset: 1239},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 59, col: 9, offset: 1239},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 59, col: 11, offset: 1241},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 17, offset: 1247},
								name: "Term",
							},
						},
						&labeledExpr{
							pos:   position{line: 59, col: 22, offset: 1252},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 59, col: 27, offset: 1257},
								expr: &seqExpr{
									pos: position{line: 59, col: 29, offset: 1259},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 59, col: 29, offset: 1259},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 59, col: 31, offset: 1261},
											name: "AddOp",
										},
										&ruleRefExpr{
											pos:  position{line: 59, col: 37, offset: 1267},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 59, col: 39, offset: 1269},
											name: "Term",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 59, col: 47, offset: 1277},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Term",
			pos:  position{line: 63, col: 1, offset: 1318},
			expr: &actionExpr{
				pos: position{line: 63, col: 9, offset: 1326},
				run: (*parser).callonTerm1,
				expr: &seqExpr{
					pos: position{line: 63, col: 9, offset: 1326},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 63, col: 9, offset: 1326},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 63, col: 15, offset: 1332},
								name: "Factor",
							},
						},
						&labeledExpr{
							pos:   position{line: 63, col: 22, offset: 1339},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 63, col: 27, offset: 1344},
								expr: &seqExpr{
									pos: position{line: 63, col: 29, offset: 1346},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 63, col: 29, offset: 1346},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 63, col: 31, offset: 1348},
											name: "MulOp",
										},
										&ruleRefExpr{
											pos:  position{line: 63, col: 37, offset: 1354},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 63, col: 39, offset: 1356},
											name: "Factor",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Factor",
			pos:  position{line: 67, col: 1, offset: 1405},
			expr: &choiceExpr{
				pos: position{line: 67, col: 11, offset: 1415},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 67, col: 11, offset: 1415},
						run: (*parser).callonFactor2,
						expr: &seqExpr{
							pos: position{line: 67, col: 11, offset: 1415},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 67, col: 11, offset: 1415},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 67, col: 15, offset: 1419},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 67, col: 20, offset: 1424},
										name: "Expr",
									},
								},
								&litMatcher{
									pos:        position{line: 67, col: 25, offset: 1429},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 69, col: 5, offset: 1460},
						run: (*parser).callonFactor8,
						expr: &labeledExpr{
							pos:   position{line: 69, col: 5, offset: 1460},
							label: "integer",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 13, offset: 1468},
								name: "Integer",
							},
						},
					},
				},
			},
		},
		{
			name: "AddOp",
			pos:  position{line: 73, col: 1, offset: 1505},
			expr: &actionExpr{
				pos: position{line: 73, col: 10, offset: 1514},
				run: (*parser).callonAddOp1,
				expr: &choiceExpr{
					pos: position{line: 73, col: 12, offset: 1516},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 73, col: 12, offset: 1516},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 73, col: 18, offset: 1522},
							val:        "-",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MulOp",
			pos:  position{line: 77, col: 1, offset: 1564},
			expr: &actionExpr{
				pos: position{line: 77, col: 10, offset: 1573},
				run: (*parser).callonMulOp1,
				expr: &litMatcher{
					pos:        position{line: 77, col: 10, offset: 1573},
					val:        "*",
					ignoreCase: false,
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 81, col: 1, offset: 1613},
			expr: &actionExpr{
				pos: position{line: 81, col: 12, offset: 1624},
				run: (*parser).callonInteger1,
				expr: &oneOrMoreExpr{
					pos: position{line: 81, col: 12, offset: 1624},
					expr: &charClassMatcher{
						pos:        position{line: 81, col: 12, offset: 1624},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 85, col: 1, offset: 1676},
			expr: &zeroOrMoreExpr{
				pos: position{line: 85, col: 19, offset: 1694},
				expr: &charClassMatcher{
					pos:        position{line: 85, col: 19, offset: 1694},
					val:        "[ \\n\\t\\r]",
					chars:      []rune{' ', '\n', '\t', '\r'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 87, col: 1, offset: 1706},
			expr: &notExpr{
				pos: position{line: 87, col: 8, offset: 1713},
				expr: &anyMatcher{
					line: 87, col: 9, offset: 1714,
				},
			},
		},
	},
}

func (c *current) onInput1(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["expr"])
}

func (c *current) onExpr1(first, rest interface{}) (interface{}, error) {
	return eval(first, rest), nil
}

func (p *parser) callonExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpr1(stack["first"], stack["rest"])
}

func (c *current) onTerm1(first, rest interface{}) (interface{}, error) {
	return eval(first, rest), nil
}

func (p *parser) callonTerm1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTerm1(stack["first"], stack["rest"])
}

func (c *current) onFactor2(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonFactor2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFactor2(stack["expr"])
}

func (c *current) onFactor8(expr, integer interface{}) (interface{}, error) {
	return integer, nil
}

func (p *parser) callonFactor8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFactor8(stack["expr"], stack["integer"])
}

func (c *current) onAddOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonAddOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAddOp1()
}

func (c *current) onMulOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonMulOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMulOp1()
}

func (c *current) onInteger1() (interface{}, error) {
	return strconv.Atoi(string(c.text))
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

var (
	// ErrNoRule is returned when the grammar to parse has no rule.
	ErrNoRule = errors.New("grammar has no rule")

	// ErrInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	ErrInvalidEncoding = errors.New("invalid encoding")

	// ErrNoMatch is returned if no match could be found.
	ErrNoMatch = errors.New("no match found")
)

var debug = false

// ParseFile parses the file identified by filename.
func ParseFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(filename, f)
}

// Parse parses the data from r, using filename as information in the
// error messages.
func Parse(filename string, r io.Reader) (interface{}, error) {
	return parse(filename, r, g)
}

type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
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
	run  func(*parser) (interface{}, error)
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
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
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

// ParserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type ParserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *ParserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
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

	depth  int
	rules  map[string]*rule
	vstack []map[string]interface{}
	rstack []*rule
}

func (p *parser) print(prefix, s string) string {
	if !debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &ParserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
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
			p.addErr(ErrInvalidEncoding)
		}
	}
}

func (p *parser) save() savepoint {
	if debug {
		defer p.out(p.in("save"))
	}
	return p.pt
}

func (p *parser) restore(pt savepoint) {
	if debug {
		defer p.out(p.in("restore"))
	}
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
		p.addErr(ErrNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	// panic can be used in action code to stop parsing immediately
	// and return the panic as an error.
	defer func() {
		if e := recover(); e != nil {
			if debug {
				defer p.out(p.in("panic handler"))
			}
			val = nil
			switch e := e.(type) {
			case error:
				p.addErr(e)
			default:
				p.addErr(fmt.Errorf("%v", e))
			}
			err = p.errs.err()
		}
	}()

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(ErrNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	start := p.save()
	p.rstack = append(p.rstack, rule)
	val, ok := p.parseExpr(rule.expr)
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.slice(start.position, p.save().position)))
	}
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
	if debug {
		defer p.out(p.in("parseActionExpr"))
	}

	p.vstack = append(p.vstack, make(map[string]interface{}))
	start := p.save()
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.slice(start.position, p.save().position)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
	if ok {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.slice(start.position, p.save().position)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.save()
	_, ok := p.parseExpr(and.expr)
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		p.read()
		return string(p.pt.rn), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

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
	if debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		val, ok := p.parseExpr(alt)
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	val, ok := p.parseExpr(lab.expr)
	if ok && lab.label != "" && len(p.vstack) > 0 {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseLitMatcher"))
	}

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
	if debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.save()
	_, ok := p.parseExpr(not.expr)
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

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
	if debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if debug {
		defer p.out(p.in("parseSeqExpr"))
	}

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
	if debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

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
	if debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

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

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
