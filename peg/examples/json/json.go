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

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "JSON",
			pos:  position{line: 32, col: 1, offset: 491},
			expr: &actionExpr{
				pos: position{line: 32, col: 8, offset: 500},
				run: (*parser).callonJSON1,
				expr: &seqExpr{
					pos: position{line: 32, col: 8, offset: 500},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 32, col: 8, offset: 500},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 32, col: 10, offset: 502},
							label: "vals",
							expr: &oneOrMoreExpr{
								pos: position{line: 32, col: 15, offset: 507},
								expr: &ruleRefExpr{
									pos:  position{line: 32, col: 15, offset: 507},
									name: "Value",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 32, col: 22, offset: 514},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 44, col: 1, offset: 705},
			expr: &actionExpr{
				pos: position{line: 44, col: 9, offset: 715},
				run: (*parser).callonValue1,
				expr: &seqExpr{
					pos: position{line: 44, col: 9, offset: 715},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 44, col: 9, offset: 715},
							label: "val",
							expr: &choiceExpr{
								pos: position{line: 44, col: 15, offset: 721},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 44, col: 15, offset: 721},
										name: "Object",
									},
									&ruleRefExpr{
										pos:  position{line: 44, col: 24, offset: 730},
										name: "Array",
									},
									&ruleRefExpr{
										pos:  position{line: 44, col: 32, offset: 738},
										name: "Number",
									},
									&ruleRefExpr{
										pos:  position{line: 44, col: 41, offset: 747},
										name: "String",
									},
									&ruleRefExpr{
										pos:  position{line: 44, col: 50, offset: 756},
										name: "Bool",
									},
									&ruleRefExpr{
										pos:  position{line: 44, col: 57, offset: 763},
										name: "Null",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 44, col: 64, offset: 770},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Object",
			pos:  position{line: 48, col: 1, offset: 797},
			expr: &actionExpr{
				pos: position{line: 48, col: 10, offset: 808},
				run: (*parser).callonObject1,
				expr: &seqExpr{
					pos: position{line: 48, col: 10, offset: 808},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 48, col: 10, offset: 808},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 48, col: 14, offset: 812},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 48, col: 16, offset: 814},
							label: "vals",
							expr: &zeroOrOneExpr{
								pos: position{line: 48, col: 21, offset: 819},
								expr: &seqExpr{
									pos: position{line: 48, col: 23, offset: 821},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 48, col: 23, offset: 821},
											name: "String",
										},
										&ruleRefExpr{
											pos:  position{line: 48, col: 30, offset: 828},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 48, col: 32, offset: 830},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 48, col: 36, offset: 834},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 48, col: 38, offset: 836},
											name: "Value",
										},
										&zeroOrMoreExpr{
											pos: position{line: 48, col: 44, offset: 842},
											expr: &seqExpr{
												pos: position{line: 48, col: 46, offset: 844},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 48, col: 46, offset: 844},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 48, col: 50, offset: 848},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 48, col: 52, offset: 850},
														name: "String",
													},
													&ruleRefExpr{
														pos:  position{line: 48, col: 59, offset: 857},
														name: "_",
													},
													&litMatcher{
														pos:        position{line: 48, col: 61, offset: 859},
														val:        ":",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 48, col: 65, offset: 863},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 48, col: 67, offset: 865},
														name: "Value",
													},
												},
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 48, col: 79, offset: 877},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Array",
			pos:  position{line: 63, col: 1, offset: 1219},
			expr: &actionExpr{
				pos: position{line: 63, col: 9, offset: 1229},
				run: (*parser).callonArray1,
				expr: &seqExpr{
					pos: position{line: 63, col: 9, offset: 1229},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 63, col: 9, offset: 1229},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 63, col: 13, offset: 1233},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 63, col: 15, offset: 1235},
							label: "vals",
							expr: &zeroOrOneExpr{
								pos: position{line: 63, col: 20, offset: 1240},
								expr: &seqExpr{
									pos: position{line: 63, col: 22, offset: 1242},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 63, col: 22, offset: 1242},
											name: "Value",
										},
										&zeroOrMoreExpr{
											pos: position{line: 63, col: 28, offset: 1248},
											expr: &seqExpr{
												pos: position{line: 63, col: 30, offset: 1250},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 63, col: 30, offset: 1250},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 63, col: 34, offset: 1254},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 63, col: 36, offset: 1256},
														name: "Value",
													},
												},
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 63, col: 48, offset: 1268},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 77, col: 1, offset: 1574},
			expr: &actionExpr{
				pos: position{line: 77, col: 10, offset: 1585},
				run: (*parser).callonNumber1,
				expr: &seqExpr{
					pos: position{line: 77, col: 10, offset: 1585},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 77, col: 10, offset: 1585},
							expr: &litMatcher{
								pos:        position{line: 77, col: 10, offset: 1585},
								val:        "-",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 77, col: 15, offset: 1590},
							name: "Integer",
						},
						&zeroOrOneExpr{
							pos: position{line: 77, col: 23, offset: 1598},
							expr: &seqExpr{
								pos: position{line: 77, col: 25, offset: 1600},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 77, col: 25, offset: 1600},
										val:        ".",
										ignoreCase: false,
									},
									&oneOrMoreExpr{
										pos: position{line: 77, col: 29, offset: 1604},
										expr: &ruleRefExpr{
											pos:  position{line: 77, col: 29, offset: 1604},
											name: "DecimalDigit",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 77, col: 46, offset: 1621},
							expr: &ruleRefExpr{
								pos:  position{line: 77, col: 46, offset: 1621},
								name: "Exponent",
							},
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 83, col: 1, offset: 1776},
			expr: &choiceExpr{
				pos: position{line: 83, col: 11, offset: 1788},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 83, col: 11, offset: 1788},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 83, col: 17, offset: 1794},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 83, col: 17, offset: 1794},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 83, col: 37, offset: 1814},
								expr: &ruleRefExpr{
									pos:  position{line: 83, col: 37, offset: 1814},
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
			pos:  position{line: 85, col: 1, offset: 1829},
			expr: &seqExpr{
				pos: position{line: 85, col: 12, offset: 1842},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 85, col: 12, offset: 1842},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 85, col: 17, offset: 1847},
						expr: &charClassMatcher{
							pos:        position{line: 85, col: 17, offset: 1847},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 85, col: 23, offset: 1853},
						expr: &ruleRefExpr{
							pos:  position{line: 85, col: 23, offset: 1853},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 87, col: 1, offset: 1868},
			expr: &actionExpr{
				pos: position{line: 87, col: 10, offset: 1879},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 87, col: 10, offset: 1879},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 87, col: 10, offset: 1879},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 87, col: 14, offset: 1883},
							expr: &choiceExpr{
								pos: position{line: 87, col: 16, offset: 1885},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 87, col: 16, offset: 1885},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 87, col: 16, offset: 1885},
												expr: &ruleRefExpr{
													pos:  position{line: 87, col: 17, offset: 1886},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 87, col: 29, offset: 1898,
											},
										},
									},
									&seqExpr{
										pos: position{line: 87, col: 33, offset: 1902},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 87, col: 33, offset: 1902},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 87, col: 38, offset: 1907},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 87, col: 56, offset: 1925},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 93, col: 1, offset: 2097},
			expr: &charClassMatcher{
				pos:        position{line: 93, col: 15, offset: 2113},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 95, col: 1, offset: 2129},
			expr: &choiceExpr{
				pos: position{line: 95, col: 18, offset: 2148},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 95, col: 18, offset: 2148},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 37, offset: 2167},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 97, col: 1, offset: 2182},
			expr: &charClassMatcher{
				pos:        position{line: 97, col: 20, offset: 2203},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 99, col: 1, offset: 2216},
			expr: &seqExpr{
				pos: position{line: 99, col: 17, offset: 2234},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 99, col: 17, offset: 2234},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 21, offset: 2238},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 30, offset: 2247},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 39, offset: 2256},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 48, offset: 2265},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 101, col: 1, offset: 2275},
			expr: &charClassMatcher{
				pos:        position{line: 101, col: 16, offset: 2292},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 103, col: 1, offset: 2299},
			expr: &charClassMatcher{
				pos:        position{line: 103, col: 23, offset: 2323},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 105, col: 1, offset: 2330},
			expr: &charClassMatcher{
				pos:        position{line: 105, col: 12, offset: 2343},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Bool",
			pos:  position{line: 107, col: 1, offset: 2354},
			expr: &choiceExpr{
				pos: position{line: 107, col: 8, offset: 2363},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 107, col: 8, offset: 2363},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 107, col: 8, offset: 2363},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 107, col: 38, offset: 2393},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 107, col: 38, offset: 2393},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 109, col: 1, offset: 2424},
			expr: &actionExpr{
				pos: position{line: 109, col: 8, offset: 2433},
				run: (*parser).callonNull1,
				expr: &litMatcher{
					pos:        position{line: 109, col: 8, offset: 2433},
					val:        "null",
					ignoreCase: false,
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 111, col: 1, offset: 2461},
			expr: &zeroOrMoreExpr{
				pos: position{line: 111, col: 18, offset: 2480},
				expr: &charClassMatcher{
					pos:        position{line: 111, col: 18, offset: 2480},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 113, col: 1, offset: 2492},
			expr: &notExpr{
				pos: position{line: 113, col: 7, offset: 2500},
				expr: &anyMatcher{
					line: 113, col: 8, offset: 2501,
				},
			},
		},
	},
}

func (c *current) onJSON1(vals interface{}) (interface{}, error) {
	valsSl := toIfaceSlice(vals)
	switch len(valsSl) {
	case 0:
		return nil, nil
	case 1:
		return valsSl[0], nil
	default:
		return valsSl, nil
	}
}

func (p *parser) callonJSON1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onJSON1: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onJSON1(stack["vals"])
}

func (c *current) onValue1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onValue1: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onValue1(stack["val"])
}

func (c *current) onObject1(vals interface{}) (interface{}, error) {
	res := make(map[string]interface{})
	valsSl := toIfaceSlice(vals)
	if len(valsSl) == 0 {
		return res, nil
	}
	res[valsSl[0].(string)] = valsSl[4]
	restSl := toIfaceSlice(valsSl[5])
	for _, v := range restSl {
		vSl := toIfaceSlice(v)
		res[vSl[2].(string)] = vSl[6]
	}
	return res, nil
}

func (p *parser) callonObject1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onObject1: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onObject1(stack["vals"])
}

func (c *current) onArray1(vals interface{}) (interface{}, error) {
	valsSl := toIfaceSlice(vals)
	if len(valsSl) == 0 {
		return []interface{}{}, nil
	}
	res := []interface{}{valsSl[0]}
	restSl := toIfaceSlice(valsSl[1])
	for _, v := range restSl {
		vSl := toIfaceSlice(v)
		res = append(res, vSl[2])
	}
	return res, nil
}

func (p *parser) callonArray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onArray1: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onArray1(stack["vals"])
}

func (c *current) onNumber1() (interface{}, error) {
	// JSON numbers have the same syntax as Go's, and are parseable using
	// strconv.
	return strconv.ParseFloat(string(c.text), 64)
}

func (p *parser) callonNumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onNumber1: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onNumber1()
}

func (c *current) onString1() (interface{}, error) {
	// TODO : the forward slash (solidus) is not a valid escape in Go, it will
	// fail if there's one in the string
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onString1: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onString1()
}

func (c *current) onBool2() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBool2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onBool2: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onBool2()
}

func (c *current) onBool4() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBool4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onBool4: stack %d: %v\n", len(p.vstack), stack)
	return p.cur.onBool4()
}

func (c *current) onNull1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	//fmt.Printf("CALL onNull1: stack %d: %v\n", len(p.vstack), stack)
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

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
	}
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

func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
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
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
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
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
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

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
		//fmt.Printf("LABEL: set %q = %T (%s) to stack %d\n", lab.label, val, val, len(p.vstack))
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
