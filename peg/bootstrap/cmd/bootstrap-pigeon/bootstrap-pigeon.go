package main

import (
	"bufio"
	"bytes"
	"errors"
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
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "USAGE: <cmd> FILE")
		os.Exit(1)
	}

	var in io.Reader

	nm := "stdin"
	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer f.Close()
		in = f
		nm = os.Args[1]
	} else {
		in = bufio.NewReader(os.Stdin)
	}

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
			pos:  position{line: 34, col: 1, offset: 457},
			expr: &actionExpr{
				pos: position{line: 34, col: 11, offset: 467},
				run: (*parser).callonGrammar_1,
				expr: &seqExpr{
					pos: position{line: 34, col: 11, offset: 467},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 34, col: 11, offset: 467},
							name: "sp",
						},
						&litMatcher{
							pos:        position{line: 34, col: 14, offset: 470},
							val:        "package",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 34, col: 24, offset: 480},
							name: "sp",
						},
						&labeledExpr{
							pos:   position{line: 34, col: 27, offset: 483},
							label: "pkg",
							expr: &ruleRefExpr{
								pos:  position{line: 34, col: 31, offset: 487},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 34, col: 46, offset: 502},
							name: "sp",
						},
						&labeledExpr{
							pos:   position{line: 34, col: 49, offset: 505},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 34, col: 62, offset: 518},
								expr: &seqExpr{
									pos: position{line: 34, col: 62, offset: 518},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 34, col: 62, offset: 518},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 34, col: 74, offset: 530},
											name: "sp",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 34, col: 79, offset: 535},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 34, col: 86, offset: 542},
								expr: &seqExpr{
									pos: position{line: 34, col: 86, offset: 542},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 34, col: 86, offset: 542},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 34, col: 91, offset: 547},
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
			pos:  position{line: 43, col: 1, offset: 780},
			expr: &ruleRefExpr{
				pos:  position{line: 43, col: 15, offset: 794},
				name: "CodeBlock",
			},
		},
		{
			name: "Rule",
			pos:  position{line: 45, col: 1, offset: 805},
			expr: &seqExpr{
				pos: position{line: 45, col: 8, offset: 812},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 45, col: 8, offset: 812},
						name: "IdentifierName",
					},
					&zeroOrOneExpr{
						pos: position{line: 45, col: 25, offset: 829},
						expr: &ruleRefExpr{
							pos:  position{line: 45, col: 25, offset: 829},
							name: "StringLiteral",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 45, col: 42, offset: 846},
						name: "RuleDefOp",
					},
					&ruleRefExpr{
						pos:  position{line: 45, col: 52, offset: 856},
						name: "Expression",
					},
					&ruleRefExpr{
						pos:  position{line: 45, col: 63, offset: 867},
						name: "EndOfRule",
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 47, col: 1, offset: 878},
			expr: &ruleRefExpr{
				pos:  position{line: 47, col: 14, offset: 891},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 49, col: 1, offset: 903},
			expr: &seqExpr{
				pos: position{line: 49, col: 14, offset: 916},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 49, col: 14, offset: 916},
						name: "ActionExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 49, col: 27, offset: 929},
						expr: &seqExpr{
							pos: position{line: 49, col: 27, offset: 929},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 49, col: 27, offset: 929},
									val:        "/",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 49, col: 31, offset: 933},
									name: "sp",
								},
								&ruleRefExpr{
									pos:  position{line: 49, col: 34, offset: 936},
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
			pos:  position{line: 51, col: 1, offset: 951},
			expr: &seqExpr{
				pos: position{line: 51, col: 14, offset: 964},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 51, col: 14, offset: 964},
						name: "SeqExpr",
					},
					&zeroOrOneExpr{
						pos: position{line: 51, col: 24, offset: 974},
						expr: &ruleRefExpr{
							pos:  position{line: 51, col: 24, offset: 974},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 53, col: 1, offset: 988},
			expr: &seqExpr{
				pos: position{line: 53, col: 11, offset: 998},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 53, col: 11, offset: 998},
						name: "LabeledExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 53, col: 25, offset: 1012},
						expr: &ruleRefExpr{
							pos:  position{line: 53, col: 25, offset: 1012},
							name: "LabeledExpr",
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 55, col: 1, offset: 1028},
			expr: &choiceExpr{
				pos: position{line: 55, col: 15, offset: 1042},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 55, col: 15, offset: 1042},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 55, col: 15, offset: 1042},
								name: "Identifier",
							},
							&litMatcher{
								pos:        position{line: 55, col: 26, offset: 1053},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 55, col: 30, offset: 1057},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 55, col: 33, offset: 1060},
								name: "PrefixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 55, col: 48, offset: 1075},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 57, col: 1, offset: 1089},
			expr: &choiceExpr{
				pos: position{line: 57, col: 16, offset: 1104},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 57, col: 16, offset: 1104},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 57, col: 16, offset: 1104},
								name: "PrefixedOp",
							},
							&ruleRefExpr{
								pos:  position{line: 57, col: 27, offset: 1115},
								name: "SuffixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 42, offset: 1130},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 59, col: 1, offset: 1144},
			expr: &seqExpr{
				pos: position{line: 59, col: 14, offset: 1157},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 59, col: 16, offset: 1159},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 59, col: 16, offset: 1159},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 59, col: 22, offset: 1165},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 59, col: 28, offset: 1171},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 61, col: 1, offset: 1175},
			expr: &choiceExpr{
				pos: position{line: 61, col: 16, offset: 1190},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 61, col: 16, offset: 1190},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 61, col: 16, offset: 1190},
								name: "PrimaryExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 61, col: 28, offset: 1202},
								name: "SuffixedOp",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 61, col: 41, offset: 1215},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 63, col: 1, offset: 1229},
			expr: &seqExpr{
				pos: position{line: 63, col: 14, offset: 1242},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 63, col: 16, offset: 1244},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 63, col: 16, offset: 1244},
								val:        "?",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 63, col: 22, offset: 1250},
								val:        "*",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 63, col: 28, offset: 1256},
								val:        "+",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 34, offset: 1262},
						name: "sp",
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 65, col: 1, offset: 1266},
			expr: &choiceExpr{
				pos: position{line: 65, col: 15, offset: 1280},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 65, col: 15, offset: 1280},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 65, col: 28, offset: 1293},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 65, col: 47, offset: 1312},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 65, col: 60, offset: 1325},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 65, col: 74, offset: 1339},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 65, col: 93, offset: 1358},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 65, col: 93, offset: 1358},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 65, col: 97, offset: 1362},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 65, col: 100, offset: 1365},
								name: "Expression",
							},
							&litMatcher{
								pos:        position{line: 65, col: 111, offset: 1376},
								val:        ")",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 65, col: 115, offset: 1380},
								name: "sp",
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 67, col: 1, offset: 1384},
			expr: &seqExpr{
				pos: position{line: 67, col: 15, offset: 1398},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 67, col: 15, offset: 1398},
						name: "IdentifierName",
					},
					&notExpr{
						pos: position{line: 67, col: 30, offset: 1413},
						expr: &seqExpr{
							pos: position{line: 67, col: 33, offset: 1416},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 67, col: 35, offset: 1418},
									expr: &ruleRefExpr{
										pos:  position{line: 67, col: 35, offset: 1418},
										name: "StringLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 67, col: 52, offset: 1435},
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
			pos:  position{line: 69, col: 1, offset: 1442},
			expr: &seqExpr{
				pos: position{line: 69, col: 20, offset: 1461},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 69, col: 20, offset: 1461},
						name: "SemanticPredOp",
					},
					&ruleRefExpr{
						pos:  position{line: 69, col: 35, offset: 1476},
						name: "CodeBlock",
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 71, col: 1, offset: 1487},
			expr: &seqExpr{
				pos: position{line: 71, col: 18, offset: 1504},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 71, col: 20, offset: 1506},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 71, col: 20, offset: 1506},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 71, col: 26, offset: 1512},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 32, offset: 1518},
						name: "sp",
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 73, col: 1, offset: 1522},
			expr: &seqExpr{
				pos: position{line: 73, col: 13, offset: 1534},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 73, col: 15, offset: 1536},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 73, col: 15, offset: 1536},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 73, col: 21, offset: 1542},
								val:        "<-",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 73, col: 28, offset: 1549},
								val:        "←",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 73, col: 39, offset: 1560},
								val:        "⟵",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 50, offset: 1571},
						name: "sp",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 75, col: 1, offset: 1575},
			expr: &seqExpr{
				pos: position{line: 75, col: 20, offset: 1594},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 75, col: 20, offset: 1594},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 75, col: 27, offset: 1601},
						expr: &seqExpr{
							pos: position{line: 75, col: 27, offset: 1601},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 75, col: 27, offset: 1601},
									expr: &litMatcher{
										pos:        position{line: 75, col: 28, offset: 1602},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 75, col: 33, offset: 1607,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 75, col: 38, offset: 1612},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 43, offset: 1617},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 77, col: 1, offset: 1621},
			expr: &seqExpr{
				pos: position{line: 77, col: 21, offset: 1641},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 77, col: 21, offset: 1641},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 77, col: 28, offset: 1648},
						expr: &seqExpr{
							pos: position{line: 77, col: 28, offset: 1648},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 77, col: 28, offset: 1648},
									val:        "[^\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   true,
								},
								&anyMatcher{
									line: 77, col: 34, offset: 1654,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 79, col: 1, offset: 1660},
			expr: &ruleRefExpr{
				pos:  position{line: 79, col: 14, offset: 1673},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 81, col: 1, offset: 1689},
			expr: &actionExpr{
				pos: position{line: 81, col: 18, offset: 1706},
				run: (*parser).callonIdentifierName_1,
				expr: &seqExpr{
					pos: position{line: 81, col: 18, offset: 1706},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 81, col: 18, offset: 1706},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 81, col: 34, offset: 1722},
							expr: &ruleRefExpr{
								pos:  position{line: 81, col: 34, offset: 1722},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 85, col: 1, offset: 1766},
			expr: &charClassMatcher{
				pos:        position{line: 85, col: 19, offset: 1784},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 87, col: 1, offset: 1793},
			expr: &choiceExpr{
				pos: position{line: 87, col: 18, offset: 1810},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 87, col: 18, offset: 1810},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 87, col: 36, offset: 1828},
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
			pos:  position{line: 89, col: 1, offset: 1835},
			expr: &seqExpr{
				pos: position{line: 89, col: 14, offset: 1848},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 89, col: 14, offset: 1848},
						name: "StringLiteral",
					},
					&zeroOrOneExpr{
						pos: position{line: 89, col: 28, offset: 1862},
						expr: &litMatcher{
							pos:        position{line: 89, col: 28, offset: 1862},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 89, col: 33, offset: 1867},
						name: "sp",
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 91, col: 1, offset: 1871},
			expr: &choiceExpr{
				pos: position{line: 91, col: 17, offset: 1887},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 91, col: 17, offset: 1887},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 91, col: 17, offset: 1887},
								val:        "\"",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 91, col: 21, offset: 1891},
								expr: &ruleRefExpr{
									pos:  position{line: 91, col: 21, offset: 1891},
									name: "DoubleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 91, col: 39, offset: 1909},
								val:        "\"",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 91, col: 45, offset: 1915},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 91, col: 45, offset: 1915},
								val:        "'",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 91, col: 49, offset: 1919},
								expr: &ruleRefExpr{
									pos:  position{line: 91, col: 49, offset: 1919},
									name: "SingleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 91, col: 67, offset: 1937},
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
			pos:  position{line: 93, col: 1, offset: 1942},
			expr: &seqExpr{
				pos: position{line: 93, col: 20, offset: 1961},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 93, col: 20, offset: 1961},
						expr: &choiceExpr{
							pos: position{line: 93, col: 23, offset: 1964},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 93, col: 23, offset: 1964},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 93, col: 29, offset: 1970},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 93, col: 36, offset: 1977},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 93, col: 43, offset: 1984,
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 95, col: 1, offset: 1988},
			expr: &seqExpr{
				pos: position{line: 95, col: 20, offset: 2007},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 95, col: 20, offset: 2007},
						expr: &choiceExpr{
							pos: position{line: 95, col: 23, offset: 2010},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 95, col: 23, offset: 2010},
									val:        "'",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 95, col: 29, offset: 2016},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 95, col: 36, offset: 2023},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 95, col: 43, offset: 2030,
					},
				},
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 97, col: 1, offset: 2033},
			expr: &seqExpr{
				pos: position{line: 97, col: 20, offset: 2052},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 97, col: 20, offset: 2052},
						val:        "[",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 97, col: 24, offset: 2056},
						expr: &litMatcher{
							pos:        position{line: 97, col: 24, offset: 2056},
							val:        "^",
							ignoreCase: false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 97, col: 31, offset: 2063},
						expr: &choiceExpr{
							pos: position{line: 97, col: 31, offset: 2063},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 97, col: 31, offset: 2063},
									name: "ClassCharRange",
								},
								&ruleRefExpr{
									pos:  position{line: 97, col: 48, offset: 2080},
									name: "ClassChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 97, col: 61, offset: 2093},
						val:        "]",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 97, col: 65, offset: 2097},
						expr: &litMatcher{
							pos:        position{line: 97, col: 65, offset: 2097},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 97, col: 70, offset: 2102},
						name: "sp",
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 99, col: 1, offset: 2106},
			expr: &seqExpr{
				pos: position{line: 99, col: 18, offset: 2123},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 99, col: 18, offset: 2123},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 99, col: 28, offset: 2133},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 32, offset: 2137},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 101, col: 1, offset: 2148},
			expr: &seqExpr{
				pos: position{line: 101, col: 13, offset: 2160},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 101, col: 13, offset: 2160},
						expr: &choiceExpr{
							pos: position{line: 101, col: 16, offset: 2163},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 101, col: 16, offset: 2163},
									val:        "]",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 101, col: 22, offset: 2169},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 101, col: 29, offset: 2176},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 101, col: 36, offset: 2183,
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 103, col: 1, offset: 2187},
			expr: &litMatcher{
				pos:        position{line: 103, col: 14, offset: 2200},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 105, col: 1, offset: 2205},
			expr: &seqExpr{
				pos: position{line: 105, col: 13, offset: 2217},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 105, col: 13, offset: 2217},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 105, col: 17, offset: 2221},
						name: "sp",
					},
					&ruleRefExpr{
						pos:  position{line: 105, col: 20, offset: 2224},
						name: "Code",
					},
					&litMatcher{
						pos:        position{line: 105, col: 25, offset: 2229},
						val:        "}",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 105, col: 29, offset: 2233},
						name: "sp",
					},
				},
			},
		},
		{
			name: "sp",
			pos:  position{line: 107, col: 1, offset: 2237},
			expr: &zeroOrMoreExpr{
				pos: position{line: 107, col: 6, offset: 2242},
				expr: &charClassMatcher{
					pos:        position{line: 107, col: 6, offset: 2242},
					val:        "[ \\n\\r\\t]",
					chars:      []rune{' ', '\n', '\r', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EndOfRule",
			pos:  position{line: 109, col: 1, offset: 2254},
			expr: &choiceExpr{
				pos: position{line: 109, col: 13, offset: 2266},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 109, col: 13, offset: 2266},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 109, col: 13, offset: 2266},
								val:        ";",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 109, col: 17, offset: 2270},
								name: "sp",
							},
						},
					},
					&seqExpr{
						pos: position{line: 109, col: 22, offset: 2275},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 109, col: 22, offset: 2275},
								expr: &ruleRefExpr{
									pos:  position{line: 109, col: 22, offset: 2275},
									name: "SingleLineComment",
								},
							},
							&charClassMatcher{
								pos:        position{line: 109, col: 41, offset: 2294},
								val:        "[\\n]",
								chars:      []rune{'\n'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 109, col: 48, offset: 2301},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "EndOfFile",
			pos:  position{line: 111, col: 1, offset: 2312},
			expr: &notExpr{
				pos: position{line: 111, col: 13, offset: 2324},
				expr: &anyMatcher{
					line: 111, col: 14, offset: 2325,
				},
			},
		},
	},
}

func (c *current) onGrammar_1(pkg, initializer, rules interface{}) (interface{}, error) {
	pos := ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}
	nm := ast.NewIdentifier(pos, pkg.(string))
	pack := ast.NewPackage(pos)
	pack.Name = nm
	g := ast.NewGrammar(pos, pack)
	return g, nil
}

func (p *parser) callonGrammar_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar_1(stack["pkg"], stack["initializer"], stack["rules"])
}

func (c *current) onIdentifierName_1() (interface{}, error) {
	return c.text, nil
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
