package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/exp/peg/ast"
)

func main() {
	dbgFlag := flag.Bool("debug", false, "set debug mode")
	//noBuildFlag := flag.Bool("x", false, "do not build, only parse")
	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "USAGE: <cmd> FILE")
		os.Exit(1)
	}

	var in io.Reader

	nm := "stdin"
	if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer f.Close()
		in = f
		nm = flag.Arg(0)
	} else {
		in = bufio.NewReader(os.Stdin)
	}

	debug = *dbgFlag
	res, err := Parse(nm, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 39, col: 1, offset: 626},
			expr: &actionExpr{
				pos: position{line: 39, col: 11, offset: 636},
				run: (*parser).callonGrammar_1,
				expr: &seqExpr{
					pos: position{line: 39, col: 11, offset: 636},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 39, col: 11, offset: 636},
							name: "sp",
						},
						&litMatcher{
							pos:        position{line: 39, col: 14, offset: 639},
							val:        "package",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 39, col: 24, offset: 649},
							name: "sp",
						},
						&labeledExpr{
							pos:   position{line: 39, col: 27, offset: 652},
							label: "pkg",
							expr: &ruleRefExpr{
								pos:  position{line: 39, col: 31, offset: 656},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 39, col: 46, offset: 671},
							name: "sp",
						},
						&labeledExpr{
							pos:   position{line: 39, col: 49, offset: 674},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 39, col: 62, offset: 687},
								expr: &seqExpr{
									pos: position{line: 39, col: 62, offset: 687},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 39, col: 62, offset: 687},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 39, col: 74, offset: 699},
											name: "sp",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 39, col: 79, offset: 704},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 39, col: 86, offset: 711},
								expr: &seqExpr{
									pos: position{line: 39, col: 86, offset: 711},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 39, col: 86, offset: 711},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 39, col: 91, offset: 716},
											name: "sp",
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
			name: "Initializer",
			pos:  position{line: 63, col: 1, offset: 1406},
			expr: &actionExpr{
				pos: position{line: 63, col: 15, offset: 1420},
				run: (*parser).callonInitializer_1,
				expr: &ruleRefExpr{
					pos:  position{line: 63, col: 15, offset: 1420},
					name: "CodeBlock",
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 69, col: 1, offset: 1574},
			expr: &actionExpr{
				pos: position{line: 69, col: 8, offset: 1581},
				run: (*parser).callonRule_1,
				expr: &seqExpr{
					pos: position{line: 69, col: 8, offset: 1581},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 69, col: 8, offset: 1581},
							label: "ident",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 14, offset: 1587},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 29, offset: 1602},
							name: "sp",
						},
						&labeledExpr{
							pos:   position{line: 69, col: 32, offset: 1605},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 69, col: 42, offset: 1615},
								expr: &seqExpr{
									pos: position{line: 69, col: 42, offset: 1615},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 69, col: 42, offset: 1615},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 69, col: 56, offset: 1629},
											name: "sp",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 62, offset: 1635},
							name: "RuleDefOp",
						},
						&labeledExpr{
							pos:   position{line: 69, col: 72, offset: 1645},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 77, offset: 1650},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 88, offset: 1661},
							name: "EndOfRule",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 83, col: 1, offset: 2007},
			expr: &ruleRefExpr{
				pos:  position{line: 83, col: 14, offset: 2020},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 85, col: 1, offset: 2032},
			expr: &seqExpr{
				pos: position{line: 85, col: 14, offset: 2045},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 85, col: 14, offset: 2045},
						name: "ActionExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 85, col: 27, offset: 2058},
						expr: &seqExpr{
							pos: position{line: 85, col: 27, offset: 2058},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 85, col: 27, offset: 2058},
									val:        "/",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 31, offset: 2062},
									name: "sp",
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 34, offset: 2065},
									name: "ActionExpr",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ActionExpr",
			pos:  position{line: 87, col: 1, offset: 2080},
			expr: &seqExpr{
				pos: position{line: 87, col: 14, offset: 2093},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 87, col: 14, offset: 2093},
						name: "SeqExpr",
					},
					&zeroOrOneExpr{
						pos: position{line: 87, col: 24, offset: 2103},
						expr: &ruleRefExpr{
							pos:  position{line: 87, col: 24, offset: 2103},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 89, col: 1, offset: 2117},
			expr: &seqExpr{
				pos: position{line: 89, col: 11, offset: 2127},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 89, col: 11, offset: 2127},
						name: "LabeledExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 89, col: 25, offset: 2141},
						expr: &ruleRefExpr{
							pos:  position{line: 89, col: 25, offset: 2141},
							name: "LabeledExpr",
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 91, col: 1, offset: 2157},
			expr: &choiceExpr{
				pos: position{line: 91, col: 15, offset: 2171},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 91, col: 15, offset: 2171},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 91, col: 15, offset: 2171},
								name: "Identifier",
							},
							&litMatcher{
								pos:        position{line: 91, col: 26, offset: 2182},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 91, col: 30, offset: 2186},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 91, col: 33, offset: 2189},
								name: "PrefixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 48, offset: 2204},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 93, col: 1, offset: 2218},
			expr: &choiceExpr{
				pos: position{line: 93, col: 16, offset: 2233},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 93, col: 16, offset: 2233},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 93, col: 16, offset: 2233},
								name: "PrefixedOp",
							},
							&ruleRefExpr{
								pos:  position{line: 93, col: 27, offset: 2244},
								name: "SuffixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 93, col: 42, offset: 2259},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 95, col: 1, offset: 2273},
			expr: &seqExpr{
				pos: position{line: 95, col: 14, offset: 2286},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 95, col: 16, offset: 2288},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 95, col: 16, offset: 2288},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 95, col: 22, offset: 2294},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 28, offset: 2300},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 97, col: 1, offset: 2304},
			expr: &choiceExpr{
				pos: position{line: 97, col: 16, offset: 2319},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 97, col: 16, offset: 2319},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 97, col: 16, offset: 2319},
								name: "PrimaryExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 97, col: 28, offset: 2331},
								name: "SuffixedOp",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 97, col: 41, offset: 2344},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 99, col: 1, offset: 2358},
			expr: &seqExpr{
				pos: position{line: 99, col: 14, offset: 2371},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 99, col: 16, offset: 2373},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 99, col: 16, offset: 2373},
								val:        "?",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 99, col: 22, offset: 2379},
								val:        "*",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 99, col: 28, offset: 2385},
								val:        "+",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 34, offset: 2391},
						name: "sp",
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 101, col: 1, offset: 2395},
			expr: &choiceExpr{
				pos: position{line: 101, col: 15, offset: 2409},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 101, col: 15, offset: 2409},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 101, col: 28, offset: 2422},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 101, col: 47, offset: 2441},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 101, col: 60, offset: 2454},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 101, col: 74, offset: 2468},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 101, col: 93, offset: 2487},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 101, col: 93, offset: 2487},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 101, col: 97, offset: 2491},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 101, col: 100, offset: 2494},
								name: "Expression",
							},
							&litMatcher{
								pos:        position{line: 101, col: 111, offset: 2505},
								val:        ")",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 101, col: 115, offset: 2509},
								name: "sp",
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 103, col: 1, offset: 2513},
			expr: &seqExpr{
				pos: position{line: 103, col: 15, offset: 2527},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 103, col: 15, offset: 2527},
						name: "IdentifierName",
					},
					&notExpr{
						pos: position{line: 103, col: 30, offset: 2542},
						expr: &seqExpr{
							pos: position{line: 103, col: 33, offset: 2545},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 103, col: 35, offset: 2547},
									expr: &ruleRefExpr{
										pos:  position{line: 103, col: 35, offset: 2547},
										name: "StringLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 103, col: 52, offset: 2564},
									val:        "=",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredExpr",
			pos:  position{line: 105, col: 1, offset: 2571},
			expr: &seqExpr{
				pos: position{line: 105, col: 20, offset: 2590},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 105, col: 20, offset: 2590},
						name: "SemanticPredOp",
					},
					&ruleRefExpr{
						pos:  position{line: 105, col: 35, offset: 2605},
						name: "CodeBlock",
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 107, col: 1, offset: 2616},
			expr: &seqExpr{
				pos: position{line: 107, col: 18, offset: 2633},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 107, col: 20, offset: 2635},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 107, col: 20, offset: 2635},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 107, col: 26, offset: 2641},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 107, col: 32, offset: 2647},
						name: "sp",
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 109, col: 1, offset: 2651},
			expr: &seqExpr{
				pos: position{line: 109, col: 13, offset: 2663},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 109, col: 15, offset: 2665},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 109, col: 15, offset: 2665},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 109, col: 21, offset: 2671},
								val:        "<-",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 109, col: 28, offset: 2678},
								val:        "←",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 109, col: 39, offset: 2689},
								val:        "⟵",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 109, col: 50, offset: 2700},
						name: "sp",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 111, col: 1, offset: 2704},
			expr: &seqExpr{
				pos: position{line: 111, col: 20, offset: 2723},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 111, col: 20, offset: 2723},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 111, col: 27, offset: 2730},
						expr: &seqExpr{
							pos: position{line: 111, col: 27, offset: 2730},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 111, col: 27, offset: 2730},
									expr: &litMatcher{
										pos:        position{line: 111, col: 28, offset: 2731},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 111, col: 33, offset: 2736,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 111, col: 38, offset: 2741},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 111, col: 43, offset: 2746},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 113, col: 1, offset: 2750},
			expr: &seqExpr{
				pos: position{line: 113, col: 21, offset: 2770},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 113, col: 21, offset: 2770},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 113, col: 28, offset: 2777},
						expr: &seqExpr{
							pos: position{line: 113, col: 28, offset: 2777},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 113, col: 28, offset: 2777},
									val:        "[^\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   true,
								},
								&anyMatcher{
									line: 113, col: 34, offset: 2783,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 115, col: 1, offset: 2789},
			expr: &ruleRefExpr{
				pos:  position{line: 115, col: 14, offset: 2802},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 117, col: 1, offset: 2818},
			expr: &actionExpr{
				pos: position{line: 117, col: 18, offset: 2835},
				run: (*parser).callonIdentifierName_1,
				expr: &seqExpr{
					pos: position{line: 117, col: 18, offset: 2835},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 117, col: 18, offset: 2835},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 117, col: 34, offset: 2851},
							expr: &ruleRefExpr{
								pos:  position{line: 117, col: 34, offset: 2851},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 121, col: 1, offset: 2903},
			expr: &charClassMatcher{
				pos:        position{line: 121, col: 19, offset: 2921},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 123, col: 1, offset: 2930},
			expr: &choiceExpr{
				pos: position{line: 123, col: 18, offset: 2947},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 123, col: 18, offset: 2947},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 123, col: 36, offset: 2965},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "LitMatcher",
			pos:  position{line: 125, col: 1, offset: 2972},
			expr: &seqExpr{
				pos: position{line: 125, col: 14, offset: 2985},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 125, col: 14, offset: 2985},
						name: "StringLiteral",
					},
					&zeroOrOneExpr{
						pos: position{line: 125, col: 28, offset: 2999},
						expr: &litMatcher{
							pos:        position{line: 125, col: 28, offset: 2999},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 125, col: 33, offset: 3004},
						name: "sp",
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 127, col: 1, offset: 3008},
			expr: &actionExpr{
				pos: position{line: 127, col: 17, offset: 3024},
				run: (*parser).callonStringLiteral_1,
				expr: &choiceExpr{
					pos: position{line: 127, col: 19, offset: 3026},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 127, col: 19, offset: 3026},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 127, col: 19, offset: 3026},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 127, col: 23, offset: 3030},
									expr: &ruleRefExpr{
										pos:  position{line: 127, col: 23, offset: 3030},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 127, col: 41, offset: 3048},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 127, col: 47, offset: 3054},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 127, col: 47, offset: 3054},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 127, col: 51, offset: 3058},
									expr: &ruleRefExpr{
										pos:  position{line: 127, col: 51, offset: 3058},
										name: "SingleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 127, col: 69, offset: 3076},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleStringChar",
			pos:  position{line: 134, col: 1, offset: 3279},
			expr: &seqExpr{
				pos: position{line: 134, col: 20, offset: 3298},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 134, col: 20, offset: 3298},
						expr: &choiceExpr{
							pos: position{line: 134, col: 23, offset: 3301},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 134, col: 23, offset: 3301},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 134, col: 29, offset: 3307},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 134, col: 36, offset: 3314},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 134, col: 43, offset: 3321,
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 136, col: 1, offset: 3325},
			expr: &seqExpr{
				pos: position{line: 136, col: 20, offset: 3344},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 136, col: 20, offset: 3344},
						expr: &choiceExpr{
							pos: position{line: 136, col: 23, offset: 3347},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 136, col: 23, offset: 3347},
									val:        "'",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 136, col: 29, offset: 3353},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 136, col: 36, offset: 3360},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 136, col: 43, offset: 3367,
					},
				},
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 138, col: 1, offset: 3370},
			expr: &seqExpr{
				pos: position{line: 138, col: 20, offset: 3389},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 138, col: 20, offset: 3389},
						val:        "[",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 138, col: 24, offset: 3393},
						expr: &litMatcher{
							pos:        position{line: 138, col: 24, offset: 3393},
							val:        "^",
							ignoreCase: false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 138, col: 31, offset: 3400},
						expr: &choiceExpr{
							pos: position{line: 138, col: 31, offset: 3400},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 138, col: 31, offset: 3400},
									name: "ClassCharRange",
								},
								&ruleRefExpr{
									pos:  position{line: 138, col: 48, offset: 3417},
									name: "ClassChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 138, col: 61, offset: 3430},
						val:        "]",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 138, col: 65, offset: 3434},
						expr: &litMatcher{
							pos:        position{line: 138, col: 65, offset: 3434},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 138, col: 70, offset: 3439},
						name: "sp",
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 140, col: 1, offset: 3443},
			expr: &seqExpr{
				pos: position{line: 140, col: 18, offset: 3460},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 140, col: 18, offset: 3460},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 140, col: 28, offset: 3470},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 140, col: 32, offset: 3474},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 142, col: 1, offset: 3485},
			expr: &seqExpr{
				pos: position{line: 142, col: 13, offset: 3497},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 142, col: 13, offset: 3497},
						expr: &choiceExpr{
							pos: position{line: 142, col: 16, offset: 3500},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 142, col: 16, offset: 3500},
									val:        "]",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 142, col: 22, offset: 3506},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 142, col: 29, offset: 3513},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 142, col: 36, offset: 3520,
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 144, col: 1, offset: 3524},
			expr: &litMatcher{
				pos:        position{line: 144, col: 14, offset: 3537},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 146, col: 1, offset: 3542},
			expr: &seqExpr{
				pos: position{line: 146, col: 13, offset: 3554},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 146, col: 13, offset: 3554},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 146, col: 17, offset: 3558},
						name: "sp",
					},
					&ruleRefExpr{
						pos:  position{line: 146, col: 20, offset: 3561},
						name: "Code",
					},
					&litMatcher{
						pos:        position{line: 146, col: 25, offset: 3566},
						val:        "}",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 146, col: 29, offset: 3570},
						name: "sp",
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 148, col: 1, offset: 3574},
			expr: &zeroOrMoreExpr{
				pos: position{line: 148, col: 10, offset: 3583},
				expr: &choiceExpr{
					pos: position{line: 148, col: 10, offset: 3583},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 148, col: 12, offset: 3585},
							expr: &seqExpr{
								pos: position{line: 148, col: 12, offset: 3585},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 148, col: 12, offset: 3585},
										val:        "[^{}]",
										chars:      []rune{'{', '}'},
										ignoreCase: false,
										inverted:   true,
									},
									&anyMatcher{
										line: 148, col: 18, offset: 3591,
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 148, col: 25, offset: 3598},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 148, col: 25, offset: 3598},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 148, col: 29, offset: 3602},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 148, col: 34, offset: 3607},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "sp",
			pos:  position{line: 150, col: 1, offset: 3615},
			expr: &zeroOrMoreExpr{
				pos: position{line: 150, col: 6, offset: 3620},
				expr: &charClassMatcher{
					pos:        position{line: 150, col: 6, offset: 3620},
					val:        "[ \\n\\r\\t]",
					chars:      []rune{' ', '\n', '\r', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EndOfRule",
			pos:  position{line: 152, col: 1, offset: 3632},
			expr: &choiceExpr{
				pos: position{line: 152, col: 13, offset: 3644},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 152, col: 13, offset: 3644},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 152, col: 13, offset: 3644},
								val:        ";",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 152, col: 17, offset: 3648},
								name: "sp",
							},
						},
					},
					&seqExpr{
						pos: position{line: 152, col: 22, offset: 3653},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 152, col: 22, offset: 3653},
								expr: &ruleRefExpr{
									pos:  position{line: 152, col: 22, offset: 3653},
									name: "SingleLineComment",
								},
							},
							&charClassMatcher{
								pos:        position{line: 152, col: 41, offset: 3672},
								val:        "[\\n]",
								chars:      []rune{'\n'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 152, col: 48, offset: 3679},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "EndOfFile",
			pos:  position{line: 154, col: 1, offset: 3690},
			expr: &notExpr{
				pos: position{line: 154, col: 13, offset: 3702},
				expr: &anyMatcher{
					line: 154, col: 14, offset: 3703,
				},
			},
		},
	},
}

func (c *current) onGrammar_1(pkg, initializer, rules interface{}) (interface{}, error) {
	pos := ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}

	// create the package identifier
	nm := ast.NewIdentifier(pos, pkg.(string))
	// create the package
	pack := ast.NewPackage(pos)
	pack.Name = nm
	// create the grammar, assign its initializer
	g := ast.NewGrammar(pos, pack)
	if initializer != nil {
		g.Init = initializer.([]interface{})[0].(*ast.CodeBlock)
	}

	rulesTuples := rules.([]interface{})
	g.Rules = make([]*ast.Rule, len(rulesTuples))
	for i, duo := range rulesTuples {
		g.Rules[i] = duo.([]interface{})[0].(*ast.Rule)
		fmt.Println("rule: ", g.Rules[i].Name.Val)
	}

	return g, nil
}

func (p *parser) callonGrammar_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar_1(stack["pkg"], stack["initializer"], stack["rules"])
}

func (c *current) onInitializer_1() (interface{}, error) {
	pos := ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonInitializer_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInitializer_1()
}

func (c *current) onRule_1(ident, display, expr interface{}) (interface{}, error) {
	pos := ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}

	// create the rule identifier
	nm := ast.NewIdentifier(pos, ident.(string))
	rule := ast.NewRule(pos, nm)
	if display != nil {
		rule.DisplayName = display.([]interface{})[0].(*ast.StringLit)
	}
	// TODO : expr

	return rule, nil
}

func (p *parser) callonRule_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRule_1(stack["ident"], stack["display"], stack["expr"])
}

func (c *current) onIdentifierName_1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdentifierName_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierName_1()
}

func (c *current) onStringLiteral_1() (interface{}, error) {
	pos := ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}
	fmt.Println("StringLit match: ", string(c.text))
	sl := ast.NewStringLit(pos, string(c.text))
	return sl, nil
}

func (p *parser) callonStringLiteral_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral_1()
}

var (
	ErrNoRule          = errors.New("grammar has no rule")
	ErrInvalidEncoding = errors.New("invalid encoding")
	ErrNoMatch         = errors.New("no match found")
)

var debug = false

func ParseFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(filename, f)
}

func Parse(filename string, r io.Reader) (interface{}, error) {
	return parse(filename, r, g)
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

func (p *parser) print(prefix, s string) string {
	if !debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	return p.print(">", s)
}

func (p *parser) out(s string) string {
	return p.print("<", s)
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
	defer p.out(p.in("save"))
	return p.pt
}

func (p *parser) restore(pt savepoint) {
	defer p.out(p.in("restore"))
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
			defer p.out(p.in("panic handler"))
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
	defer p.out(p.in("parseRule " + rule.name))

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
	defer p.out(p.in("parseActionExpr"))

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
	defer p.out(p.in("parseAndCodeExpr"))

	ok, err := and.run(p)
	if err != nil {
		p.errs.add(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	defer p.out(p.in("parseAndExpr"))

	pt := p.save()
	_, ok := p.parseExpr(and.expr)
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	defer p.out(p.in("parseAnyMatcher"))

	if p.pt.rn != utf8.RuneError {
		p.read()
		return string(p.pt.rn), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	defer p.out(p.in("parseCharClassMatcher"))

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
	defer p.out(p.in("parseChoiceExpr"))

	for _, alt := range ch.alternatives {
		val, ok := p.parseExpr(alt)
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	defer p.out(p.in("parseLabeledExpr"))

	val, ok := p.parseExpr(lab.expr)
	if ok && lab.label != "" && len(p.vstack) > 0 {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	defer p.out(p.in("parseLitMatcher"))

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
	defer p.out(p.in("parseNotCodeExpr"))

	ok, err := not.run(p)
	if err != nil {
		p.errs.add(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	defer p.out(p.in("parseNotExpr"))

	pt := p.save()
	_, ok := p.parseExpr(not.expr)
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	defer p.out(p.in("parseOneOrMoreExpr"))

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
	defer p.out(p.in("parseRuleRefExpr " + ref.name))

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
	defer p.out(p.in("parseSeqExpr"))

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
	defer p.out(p.in("parseZeroOrMoreExpr"))

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
	defer p.out(p.in("parseZeroOrOneExpr"))

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
