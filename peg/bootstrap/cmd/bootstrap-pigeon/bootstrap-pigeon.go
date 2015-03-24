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
			pos:  position{line: 56, col: 1, offset: 1163},
			expr: &actionExpr{
				pos: position{line: 56, col: 15, offset: 1177},
				run: (*parser).callonInitializer_1,
				expr: &ruleRefExpr{
					pos:  position{line: 56, col: 15, offset: 1177},
					name: "CodeBlock",
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 62, col: 1, offset: 1331},
			expr: &seqExpr{
				pos: position{line: 62, col: 8, offset: 1338},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 62, col: 8, offset: 1338},
						name: "IdentifierName",
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 23, offset: 1353},
						name: "sp",
					},
					&zeroOrOneExpr{
						pos: position{line: 62, col: 28, offset: 1358},
						expr: &ruleRefExpr{
							pos:  position{line: 62, col: 28, offset: 1358},
							name: "StringLiteral",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 45, offset: 1375},
						name: "RuleDefOp",
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 55, offset: 1385},
						name: "Expression",
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 66, offset: 1396},
						name: "EndOfRule",
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 64, col: 1, offset: 1407},
			expr: &ruleRefExpr{
				pos:  position{line: 64, col: 14, offset: 1420},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 66, col: 1, offset: 1432},
			expr: &seqExpr{
				pos: position{line: 66, col: 14, offset: 1445},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 66, col: 14, offset: 1445},
						name: "ActionExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 66, col: 27, offset: 1458},
						expr: &seqExpr{
							pos: position{line: 66, col: 27, offset: 1458},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 66, col: 27, offset: 1458},
									val:        "/",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 66, col: 31, offset: 1462},
									name: "sp",
								},
								&ruleRefExpr{
									pos:  position{line: 66, col: 34, offset: 1465},
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
			pos:  position{line: 68, col: 1, offset: 1480},
			expr: &seqExpr{
				pos: position{line: 68, col: 14, offset: 1493},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 68, col: 14, offset: 1493},
						name: "SeqExpr",
					},
					&zeroOrOneExpr{
						pos: position{line: 68, col: 24, offset: 1503},
						expr: &ruleRefExpr{
							pos:  position{line: 68, col: 24, offset: 1503},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 70, col: 1, offset: 1517},
			expr: &seqExpr{
				pos: position{line: 70, col: 11, offset: 1527},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 11, offset: 1527},
						name: "LabeledExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 70, col: 25, offset: 1541},
						expr: &ruleRefExpr{
							pos:  position{line: 70, col: 25, offset: 1541},
							name: "LabeledExpr",
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 72, col: 1, offset: 1557},
			expr: &choiceExpr{
				pos: position{line: 72, col: 15, offset: 1571},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 72, col: 15, offset: 1571},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 72, col: 15, offset: 1571},
								name: "Identifier",
							},
							&litMatcher{
								pos:        position{line: 72, col: 26, offset: 1582},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 72, col: 30, offset: 1586},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 72, col: 33, offset: 1589},
								name: "PrefixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 48, offset: 1604},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 74, col: 1, offset: 1618},
			expr: &choiceExpr{
				pos: position{line: 74, col: 16, offset: 1633},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 74, col: 16, offset: 1633},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 74, col: 16, offset: 1633},
								name: "PrefixedOp",
							},
							&ruleRefExpr{
								pos:  position{line: 74, col: 27, offset: 1644},
								name: "SuffixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 74, col: 42, offset: 1659},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 76, col: 1, offset: 1673},
			expr: &seqExpr{
				pos: position{line: 76, col: 14, offset: 1686},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 76, col: 16, offset: 1688},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 76, col: 16, offset: 1688},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 76, col: 22, offset: 1694},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 28, offset: 1700},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 78, col: 1, offset: 1704},
			expr: &choiceExpr{
				pos: position{line: 78, col: 16, offset: 1719},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 78, col: 16, offset: 1719},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 78, col: 16, offset: 1719},
								name: "PrimaryExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 78, col: 28, offset: 1731},
								name: "SuffixedOp",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 41, offset: 1744},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 80, col: 1, offset: 1758},
			expr: &seqExpr{
				pos: position{line: 80, col: 14, offset: 1771},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 80, col: 16, offset: 1773},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 80, col: 16, offset: 1773},
								val:        "?",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 80, col: 22, offset: 1779},
								val:        "*",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 80, col: 28, offset: 1785},
								val:        "+",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 34, offset: 1791},
						name: "sp",
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 82, col: 1, offset: 1795},
			expr: &choiceExpr{
				pos: position{line: 82, col: 15, offset: 1809},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 82, col: 15, offset: 1809},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 82, col: 28, offset: 1822},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 82, col: 47, offset: 1841},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 82, col: 60, offset: 1854},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 82, col: 74, offset: 1868},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 82, col: 93, offset: 1887},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 82, col: 93, offset: 1887},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 82, col: 97, offset: 1891},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 82, col: 100, offset: 1894},
								name: "Expression",
							},
							&litMatcher{
								pos:        position{line: 82, col: 111, offset: 1905},
								val:        ")",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 82, col: 115, offset: 1909},
								name: "sp",
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 84, col: 1, offset: 1913},
			expr: &seqExpr{
				pos: position{line: 84, col: 15, offset: 1927},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 84, col: 15, offset: 1927},
						name: "IdentifierName",
					},
					&notExpr{
						pos: position{line: 84, col: 30, offset: 1942},
						expr: &seqExpr{
							pos: position{line: 84, col: 33, offset: 1945},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 84, col: 35, offset: 1947},
									expr: &ruleRefExpr{
										pos:  position{line: 84, col: 35, offset: 1947},
										name: "StringLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 84, col: 52, offset: 1964},
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
			pos:  position{line: 86, col: 1, offset: 1971},
			expr: &seqExpr{
				pos: position{line: 86, col: 20, offset: 1990},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 86, col: 20, offset: 1990},
						name: "SemanticPredOp",
					},
					&ruleRefExpr{
						pos:  position{line: 86, col: 35, offset: 2005},
						name: "CodeBlock",
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 88, col: 1, offset: 2016},
			expr: &seqExpr{
				pos: position{line: 88, col: 18, offset: 2033},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 88, col: 20, offset: 2035},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 88, col: 20, offset: 2035},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 88, col: 26, offset: 2041},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 88, col: 32, offset: 2047},
						name: "sp",
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 90, col: 1, offset: 2051},
			expr: &seqExpr{
				pos: position{line: 90, col: 13, offset: 2063},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 90, col: 15, offset: 2065},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 90, col: 15, offset: 2065},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 90, col: 21, offset: 2071},
								val:        "<-",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 90, col: 28, offset: 2078},
								val:        "←",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 90, col: 39, offset: 2089},
								val:        "⟵",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 50, offset: 2100},
						name: "sp",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 92, col: 1, offset: 2104},
			expr: &seqExpr{
				pos: position{line: 92, col: 20, offset: 2123},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 92, col: 20, offset: 2123},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 92, col: 27, offset: 2130},
						expr: &seqExpr{
							pos: position{line: 92, col: 27, offset: 2130},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 92, col: 27, offset: 2130},
									expr: &litMatcher{
										pos:        position{line: 92, col: 28, offset: 2131},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 92, col: 33, offset: 2136,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 92, col: 38, offset: 2141},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 92, col: 43, offset: 2146},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 94, col: 1, offset: 2150},
			expr: &seqExpr{
				pos: position{line: 94, col: 21, offset: 2170},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 94, col: 21, offset: 2170},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 94, col: 28, offset: 2177},
						expr: &seqExpr{
							pos: position{line: 94, col: 28, offset: 2177},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 94, col: 28, offset: 2177},
									val:        "[^\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   true,
								},
								&anyMatcher{
									line: 94, col: 34, offset: 2183,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 96, col: 1, offset: 2189},
			expr: &ruleRefExpr{
				pos:  position{line: 96, col: 14, offset: 2202},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 98, col: 1, offset: 2218},
			expr: &actionExpr{
				pos: position{line: 98, col: 18, offset: 2235},
				run: (*parser).callonIdentifierName_1,
				expr: &seqExpr{
					pos: position{line: 98, col: 18, offset: 2235},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 98, col: 18, offset: 2235},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 98, col: 34, offset: 2251},
							expr: &ruleRefExpr{
								pos:  position{line: 98, col: 34, offset: 2251},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 102, col: 1, offset: 2303},
			expr: &charClassMatcher{
				pos:        position{line: 102, col: 19, offset: 2321},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 104, col: 1, offset: 2330},
			expr: &choiceExpr{
				pos: position{line: 104, col: 18, offset: 2347},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 104, col: 18, offset: 2347},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 104, col: 36, offset: 2365},
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
			pos:  position{line: 106, col: 1, offset: 2372},
			expr: &seqExpr{
				pos: position{line: 106, col: 14, offset: 2385},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 106, col: 14, offset: 2385},
						name: "StringLiteral",
					},
					&zeroOrOneExpr{
						pos: position{line: 106, col: 28, offset: 2399},
						expr: &litMatcher{
							pos:        position{line: 106, col: 28, offset: 2399},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 106, col: 33, offset: 2404},
						name: "sp",
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 108, col: 1, offset: 2408},
			expr: &choiceExpr{
				pos: position{line: 108, col: 17, offset: 2424},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 108, col: 17, offset: 2424},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 108, col: 17, offset: 2424},
								val:        "\"",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 108, col: 21, offset: 2428},
								expr: &ruleRefExpr{
									pos:  position{line: 108, col: 21, offset: 2428},
									name: "DoubleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 108, col: 39, offset: 2446},
								val:        "\"",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 108, col: 45, offset: 2452},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 108, col: 45, offset: 2452},
								val:        "'",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 108, col: 49, offset: 2456},
								expr: &ruleRefExpr{
									pos:  position{line: 108, col: 49, offset: 2456},
									name: "SingleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 108, col: 67, offset: 2474},
								val:        "'",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleStringChar",
			pos:  position{line: 110, col: 1, offset: 2479},
			expr: &seqExpr{
				pos: position{line: 110, col: 20, offset: 2498},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 110, col: 20, offset: 2498},
						expr: &choiceExpr{
							pos: position{line: 110, col: 23, offset: 2501},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 110, col: 23, offset: 2501},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 110, col: 29, offset: 2507},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 110, col: 36, offset: 2514},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 110, col: 43, offset: 2521,
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 112, col: 1, offset: 2525},
			expr: &seqExpr{
				pos: position{line: 112, col: 20, offset: 2544},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 112, col: 20, offset: 2544},
						expr: &choiceExpr{
							pos: position{line: 112, col: 23, offset: 2547},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 112, col: 23, offset: 2547},
									val:        "'",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 112, col: 29, offset: 2553},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 112, col: 36, offset: 2560},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 112, col: 43, offset: 2567,
					},
				},
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 114, col: 1, offset: 2570},
			expr: &seqExpr{
				pos: position{line: 114, col: 20, offset: 2589},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 114, col: 20, offset: 2589},
						val:        "[",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 114, col: 24, offset: 2593},
						expr: &litMatcher{
							pos:        position{line: 114, col: 24, offset: 2593},
							val:        "^",
							ignoreCase: false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 114, col: 31, offset: 2600},
						expr: &choiceExpr{
							pos: position{line: 114, col: 31, offset: 2600},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 114, col: 31, offset: 2600},
									name: "ClassCharRange",
								},
								&ruleRefExpr{
									pos:  position{line: 114, col: 48, offset: 2617},
									name: "ClassChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 114, col: 61, offset: 2630},
						val:        "]",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 114, col: 65, offset: 2634},
						expr: &litMatcher{
							pos:        position{line: 114, col: 65, offset: 2634},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 114, col: 70, offset: 2639},
						name: "sp",
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 116, col: 1, offset: 2643},
			expr: &seqExpr{
				pos: position{line: 116, col: 18, offset: 2660},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 116, col: 18, offset: 2660},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 116, col: 28, offset: 2670},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 116, col: 32, offset: 2674},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 118, col: 1, offset: 2685},
			expr: &seqExpr{
				pos: position{line: 118, col: 13, offset: 2697},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 118, col: 13, offset: 2697},
						expr: &choiceExpr{
							pos: position{line: 118, col: 16, offset: 2700},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 118, col: 16, offset: 2700},
									val:        "]",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 118, col: 22, offset: 2706},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 118, col: 29, offset: 2713},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 118, col: 36, offset: 2720,
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 120, col: 1, offset: 2724},
			expr: &litMatcher{
				pos:        position{line: 120, col: 14, offset: 2737},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 122, col: 1, offset: 2742},
			expr: &seqExpr{
				pos: position{line: 122, col: 13, offset: 2754},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 122, col: 13, offset: 2754},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 122, col: 17, offset: 2758},
						name: "sp",
					},
					&ruleRefExpr{
						pos:  position{line: 122, col: 20, offset: 2761},
						name: "Code",
					},
					&litMatcher{
						pos:        position{line: 122, col: 25, offset: 2766},
						val:        "}",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 122, col: 29, offset: 2770},
						name: "sp",
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 124, col: 1, offset: 2774},
			expr: &zeroOrMoreExpr{
				pos: position{line: 124, col: 10, offset: 2783},
				expr: &choiceExpr{
					pos: position{line: 124, col: 10, offset: 2783},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 124, col: 12, offset: 2785},
							expr: &seqExpr{
								pos: position{line: 124, col: 12, offset: 2785},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 124, col: 12, offset: 2785},
										val:        "[^{}]",
										chars:      []rune{'{', '}'},
										ignoreCase: false,
										inverted:   true,
									},
									&anyMatcher{
										line: 124, col: 18, offset: 2791,
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 124, col: 25, offset: 2798},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 124, col: 25, offset: 2798},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 124, col: 29, offset: 2802},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 124, col: 34, offset: 2807},
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
			pos:  position{line: 126, col: 1, offset: 2815},
			expr: &zeroOrMoreExpr{
				pos: position{line: 126, col: 6, offset: 2820},
				expr: &charClassMatcher{
					pos:        position{line: 126, col: 6, offset: 2820},
					val:        "[ \\n\\r\\t]",
					chars:      []rune{' ', '\n', '\r', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EndOfRule",
			pos:  position{line: 128, col: 1, offset: 2832},
			expr: &choiceExpr{
				pos: position{line: 128, col: 13, offset: 2844},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 128, col: 13, offset: 2844},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 128, col: 13, offset: 2844},
								val:        ";",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 128, col: 17, offset: 2848},
								name: "sp",
							},
						},
					},
					&seqExpr{
						pos: position{line: 128, col: 22, offset: 2853},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 128, col: 22, offset: 2853},
								expr: &ruleRefExpr{
									pos:  position{line: 128, col: 22, offset: 2853},
									name: "SingleLineComment",
								},
							},
							&charClassMatcher{
								pos:        position{line: 128, col: 41, offset: 2872},
								val:        "[\\n]",
								chars:      []rune{'\n'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 128, col: 48, offset: 2879},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "EndOfFile",
			pos:  position{line: 130, col: 1, offset: 2890},
			expr: &notExpr{
				pos: position{line: 130, col: 13, offset: 2902},
				expr: &anyMatcher{
					line: 130, col: 14, offset: 2903,
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

func (c *current) onIdentifierName_1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdentifierName_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierName_1()
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
