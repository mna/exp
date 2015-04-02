package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	in := os.Stdin
	nm := "stdin"
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
		nm = os.Args[1]
	}

	got, err := Parse(nm, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "JSON",
			pos:  position{line: 25, col: 1, offset: 367},
			expr: &seqExpr{
				pos: position{line: 25, col: 8, offset: 376},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 25, col: 8, offset: 376},
						name: "_",
					},
					&oneOrMoreExpr{
						pos: position{line: 25, col: 10, offset: 378},
						expr: &ruleRefExpr{
							pos:  position{line: 25, col: 10, offset: 378},
							name: "Value",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 17, offset: 385},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 27, col: 1, offset: 390},
			expr: &seqExpr{
				pos: position{line: 27, col: 9, offset: 400},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 27, col: 11, offset: 402},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 27, col: 11, offset: 402},
								name: "Object",
							},
							&ruleRefExpr{
								pos:  position{line: 27, col: 20, offset: 411},
								name: "Array",
							},
							&ruleRefExpr{
								pos:  position{line: 27, col: 28, offset: 419},
								name: "Number",
							},
							&ruleRefExpr{
								pos:  position{line: 27, col: 37, offset: 428},
								name: "String",
							},
							&ruleRefExpr{
								pos:  position{line: 27, col: 46, offset: 437},
								name: "Bool",
							},
							&ruleRefExpr{
								pos:  position{line: 27, col: 53, offset: 444},
								name: "Null",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 27, col: 60, offset: 451},
						name: "_",
					},
				},
			},
		},
		{
			name: "Object",
			pos:  position{line: 29, col: 1, offset: 454},
			expr: &seqExpr{
				pos: position{line: 29, col: 10, offset: 465},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 29, col: 10, offset: 465},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 29, col: 14, offset: 469},
						name: "_",
					},
					&zeroOrOneExpr{
						pos: position{line: 29, col: 16, offset: 471},
						expr: &seqExpr{
							pos: position{line: 29, col: 18, offset: 473},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 29, col: 18, offset: 473},
									name: "String",
								},
								&ruleRefExpr{
									pos:  position{line: 29, col: 25, offset: 480},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 29, col: 27, offset: 482},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 29, col: 31, offset: 486},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 29, col: 33, offset: 488},
									name: "Value",
								},
								&zeroOrMoreExpr{
									pos: position{line: 29, col: 39, offset: 494},
									expr: &seqExpr{
										pos: position{line: 29, col: 41, offset: 496},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 29, col: 41, offset: 496},
												val:        ",",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 29, col: 45, offset: 500},
												name: "_",
											},
											&ruleRefExpr{
												pos:  position{line: 29, col: 47, offset: 502},
												name: "String",
											},
											&ruleRefExpr{
												pos:  position{line: 29, col: 54, offset: 509},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 29, col: 56, offset: 511},
												val:        ":",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 29, col: 60, offset: 515},
												name: "_",
											},
											&ruleRefExpr{
												pos:  position{line: 29, col: 62, offset: 517},
												name: "Value",
											},
										},
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 29, col: 74, offset: 529},
						val:        "}",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Array",
			pos:  position{line: 31, col: 1, offset: 534},
			expr: &seqExpr{
				pos: position{line: 31, col: 9, offset: 544},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 31, col: 9, offset: 544},
						val:        "[",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 31, col: 13, offset: 548},
						name: "_",
					},
					&zeroOrOneExpr{
						pos: position{line: 31, col: 15, offset: 550},
						expr: &seqExpr{
							pos: position{line: 31, col: 17, offset: 552},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 31, col: 17, offset: 552},
									name: "Value",
								},
								&zeroOrMoreExpr{
									pos: position{line: 31, col: 23, offset: 558},
									expr: &seqExpr{
										pos: position{line: 31, col: 25, offset: 560},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 31, col: 25, offset: 560},
												val:        ",",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 31, col: 29, offset: 564},
												name: "_",
											},
											&ruleRefExpr{
												pos:  position{line: 31, col: 31, offset: 566},
												name: "Value",
											},
										},
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 31, col: 43, offset: 578},
						val:        "]",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 33, col: 1, offset: 583},
			expr: &seqExpr{
				pos: position{line: 33, col: 10, offset: 594},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 33, col: 10, offset: 594},
						expr: &litMatcher{
							pos:        position{line: 33, col: 10, offset: 594},
							val:        "-",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 15, offset: 599},
						name: "Integer",
					},
					&zeroOrOneExpr{
						pos: position{line: 33, col: 23, offset: 607},
						expr: &seqExpr{
							pos: position{line: 33, col: 25, offset: 609},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 33, col: 25, offset: 609},
									val:        ".",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 33, col: 29, offset: 613},
									expr: &ruleRefExpr{
										pos:  position{line: 33, col: 29, offset: 613},
										name: "DecimalDigit",
									},
								},
							},
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 33, col: 46, offset: 630},
						expr: &ruleRefExpr{
							pos:  position{line: 33, col: 46, offset: 630},
							name: "Exponent",
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 35, col: 1, offset: 641},
			expr: &choiceExpr{
				pos: position{line: 35, col: 11, offset: 653},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 35, col: 11, offset: 653},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 35, col: 17, offset: 659},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 35, col: 17, offset: 659},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 35, col: 37, offset: 679},
								expr: &ruleRefExpr{
									pos:  position{line: 35, col: 37, offset: 679},
									name: "DecimalDigit",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Exponent",
			pos:  position{line: 37, col: 1, offset: 694},
			expr: &seqExpr{
				pos: position{line: 37, col: 12, offset: 707},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 37, col: 12, offset: 707},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 37, col: 17, offset: 712},
						expr: &charClassMatcher{
							pos:        position{line: 37, col: 17, offset: 712},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 37, col: 23, offset: 718},
						expr: &ruleRefExpr{
							pos:  position{line: 37, col: 23, offset: 718},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 39, col: 1, offset: 733},
			expr: &seqExpr{
				pos: position{line: 39, col: 10, offset: 744},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 39, col: 10, offset: 744},
						val:        "\"",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 39, col: 14, offset: 748},
						expr: &choiceExpr{
							pos: position{line: 39, col: 16, offset: 750},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 39, col: 16, offset: 750},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 39, col: 16, offset: 750},
											expr: &ruleRefExpr{
												pos:  position{line: 39, col: 17, offset: 751},
												name: "EscapedChar",
											},
										},
										&anyMatcher{
											line: 39, col: 29, offset: 763,
										},
									},
								},
								&seqExpr{
									pos: position{line: 39, col: 33, offset: 767},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 39, col: 33, offset: 767},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 39, col: 38, offset: 772},
											name: "EscapeSequence",
										},
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 39, col: 56, offset: 790},
						val:        "\"",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 41, col: 1, offset: 795},
			expr: &charClassMatcher{
				pos:        position{line: 41, col: 15, offset: 811},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'\x00', '0', '0', 'f', '"', '\\'},
				ranges:     []rune{'\x00', '1'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 43, col: 1, offset: 827},
			expr: &choiceExpr{
				pos: position{line: 43, col: 18, offset: 846},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 43, col: 18, offset: 846},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 37, offset: 865},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 45, col: 1, offset: 880},
			expr: &charClassMatcher{
				pos:        position{line: 45, col: 20, offset: 901},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 47, col: 1, offset: 914},
			expr: &seqExpr{
				pos: position{line: 47, col: 17, offset: 932},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 47, col: 17, offset: 932},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 21, offset: 936},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 30, offset: 945},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 39, offset: 954},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 48, offset: 963},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 49, col: 1, offset: 973},
			expr: &charClassMatcher{
				pos:        position{line: 49, col: 16, offset: 990},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 51, col: 1, offset: 997},
			expr: &charClassMatcher{
				pos:        position{line: 51, col: 23, offset: 1021},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 53, col: 1, offset: 1028},
			expr: &charClassMatcher{
				pos:        position{line: 53, col: 12, offset: 1041},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Bool",
			pos:  position{line: 55, col: 1, offset: 1052},
			expr: &choiceExpr{
				pos: position{line: 55, col: 8, offset: 1061},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 55, col: 8, offset: 1061},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 55, col: 8, offset: 1061},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 55, col: 38, offset: 1091},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 55, col: 38, offset: 1091},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 57, col: 1, offset: 1122},
			expr: &actionExpr{
				pos: position{line: 57, col: 8, offset: 1131},
				run: (*parser).callonNull1,
				expr: &litMatcher{
					pos:        position{line: 57, col: 8, offset: 1131},
					val:        "null",
					ignoreCase: false,
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 59, col: 1, offset: 1159},
			expr: &zeroOrMoreExpr{
				pos: position{line: 59, col: 18, offset: 1178},
				expr: &charClassMatcher{
					pos:        position{line: 59, col: 18, offset: 1178},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 61, col: 1, offset: 1190},
			expr: &notExpr{
				pos: position{line: 61, col: 7, offset: 1198},
				expr: &anyMatcher{
					line: 61, col: 8, offset: 1199,
				},
			},
		},
	},
}

func (c *current) onBool2() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBool2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool2()
}

func (c *current) onBool4() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBool4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool4()
}

func (c *current) onNull1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
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
	p.vstack = append(p.vstack, make(map[string]interface{}))
	val, ok := p.parseExpr(rule.expr)
	p.vstack = p.vstack[:len(p.vstack)-1]
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && debug {
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
