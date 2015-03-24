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
			name: "Rule",
			pos:  position{line: 43, col: 1, offset: 780},
			expr: &seqExpr{
				pos: position{line: 43, col: 8, offset: 787},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 43, col: 8, offset: 787},
						name: "IdentifierName",
					},
					&zeroOrOneExpr{
						pos: position{line: 43, col: 25, offset: 804},
						expr: &ruleRefExpr{
							pos:  position{line: 43, col: 25, offset: 804},
							name: "StringLiteral",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 42, offset: 821},
						name: "RuleDefOp",
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 52, offset: 831},
						name: "Expression",
					},
					&ruleRefExpr{
						pos:  position{line: 43, col: 63, offset: 842},
						name: "EndOfRule",
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 45, col: 1, offset: 853},
			expr: &ruleRefExpr{
				pos:  position{line: 45, col: 14, offset: 866},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 47, col: 1, offset: 878},
			expr: &seqExpr{
				pos: position{line: 47, col: 14, offset: 891},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 47, col: 14, offset: 891},
						name: "ActionExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 47, col: 27, offset: 904},
						expr: &seqExpr{
							pos: position{line: 47, col: 27, offset: 904},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 47, col: 27, offset: 904},
									val:        "/",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 47, col: 31, offset: 908},
									name: "sp",
								},
								&ruleRefExpr{
									pos:  position{line: 47, col: 34, offset: 911},
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
			pos:  position{line: 49, col: 1, offset: 926},
			expr: &seqExpr{
				pos: position{line: 49, col: 14, offset: 939},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 49, col: 14, offset: 939},
						name: "SeqExpr",
					},
					&zeroOrOneExpr{
						pos: position{line: 49, col: 24, offset: 949},
						expr: &ruleRefExpr{
							pos:  position{line: 49, col: 24, offset: 949},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 51, col: 1, offset: 963},
			expr: &seqExpr{
				pos: position{line: 51, col: 11, offset: 973},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 51, col: 11, offset: 973},
						name: "LabeledExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 51, col: 25, offset: 987},
						expr: &ruleRefExpr{
							pos:  position{line: 51, col: 25, offset: 987},
							name: "LabeledExpr",
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 53, col: 1, offset: 1003},
			expr: &choiceExpr{
				pos: position{line: 53, col: 15, offset: 1017},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 53, col: 15, offset: 1017},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 53, col: 15, offset: 1017},
								name: "Identifier",
							},
							&litMatcher{
								pos:        position{line: 53, col: 26, offset: 1028},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 53, col: 30, offset: 1032},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 53, col: 33, offset: 1035},
								name: "PrefixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 53, col: 48, offset: 1050},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 55, col: 1, offset: 1064},
			expr: &choiceExpr{
				pos: position{line: 55, col: 16, offset: 1079},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 55, col: 16, offset: 1079},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 55, col: 16, offset: 1079},
								name: "PrefixedOp",
							},
							&ruleRefExpr{
								pos:  position{line: 55, col: 27, offset: 1090},
								name: "SuffixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 55, col: 42, offset: 1105},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 57, col: 1, offset: 1119},
			expr: &seqExpr{
				pos: position{line: 57, col: 14, offset: 1132},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 57, col: 16, offset: 1134},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 57, col: 16, offset: 1134},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 57, col: 22, offset: 1140},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 28, offset: 1146},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 59, col: 1, offset: 1150},
			expr: &choiceExpr{
				pos: position{line: 59, col: 16, offset: 1165},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 59, col: 16, offset: 1165},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 59, col: 16, offset: 1165},
								name: "PrimaryExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 59, col: 28, offset: 1177},
								name: "SuffixedOp",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 59, col: 41, offset: 1190},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 61, col: 1, offset: 1204},
			expr: &seqExpr{
				pos: position{line: 61, col: 14, offset: 1217},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 61, col: 16, offset: 1219},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 61, col: 16, offset: 1219},
								val:        "?",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 61, col: 22, offset: 1225},
								val:        "*",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 61, col: 28, offset: 1231},
								val:        "+",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 61, col: 34, offset: 1237},
						name: "sp",
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 63, col: 1, offset: 1241},
			expr: &choiceExpr{
				pos: position{line: 63, col: 15, offset: 1255},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 63, col: 15, offset: 1255},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 28, offset: 1268},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 47, offset: 1287},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 60, offset: 1300},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 74, offset: 1314},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 63, col: 93, offset: 1333},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 63, col: 93, offset: 1333},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 63, col: 97, offset: 1337},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 63, col: 100, offset: 1340},
								name: "Expression",
							},
							&litMatcher{
								pos:        position{line: 63, col: 111, offset: 1351},
								val:        ")",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 63, col: 115, offset: 1355},
								name: "sp",
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 65, col: 1, offset: 1359},
			expr: &seqExpr{
				pos: position{line: 65, col: 15, offset: 1373},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 65, col: 15, offset: 1373},
						name: "IdentifierName",
					},
					&notExpr{
						pos: position{line: 65, col: 30, offset: 1388},
						expr: &seqExpr{
							pos: position{line: 65, col: 33, offset: 1391},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 65, col: 35, offset: 1393},
									expr: &ruleRefExpr{
										pos:  position{line: 65, col: 35, offset: 1393},
										name: "StringLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 65, col: 52, offset: 1410},
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
			pos:  position{line: 67, col: 1, offset: 1417},
			expr: &seqExpr{
				pos: position{line: 67, col: 20, offset: 1436},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 67, col: 20, offset: 1436},
						name: "SemanticPredOp",
					},
					&ruleRefExpr{
						pos:  position{line: 67, col: 35, offset: 1451},
						name: "CodeBlock",
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 69, col: 1, offset: 1462},
			expr: &seqExpr{
				pos: position{line: 69, col: 18, offset: 1479},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 69, col: 20, offset: 1481},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 69, col: 20, offset: 1481},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 69, col: 26, offset: 1487},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 69, col: 32, offset: 1493},
						name: "sp",
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 71, col: 1, offset: 1497},
			expr: &seqExpr{
				pos: position{line: 71, col: 13, offset: 1509},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 71, col: 15, offset: 1511},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 71, col: 15, offset: 1511},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 71, col: 21, offset: 1517},
								val:        "<-",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 71, col: 28, offset: 1524},
								val:        "←",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 71, col: 39, offset: 1535},
								val:        "⟵",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 50, offset: 1546},
						name: "sp",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 73, col: 1, offset: 1550},
			expr: &seqExpr{
				pos: position{line: 73, col: 20, offset: 1569},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 73, col: 20, offset: 1569},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 73, col: 27, offset: 1576},
						expr: &seqExpr{
							pos: position{line: 73, col: 27, offset: 1576},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 73, col: 27, offset: 1576},
									expr: &litMatcher{
										pos:        position{line: 73, col: 28, offset: 1577},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 73, col: 33, offset: 1582,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 73, col: 38, offset: 1587},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 43, offset: 1592},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 75, col: 1, offset: 1596},
			expr: &seqExpr{
				pos: position{line: 75, col: 21, offset: 1616},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 75, col: 21, offset: 1616},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 75, col: 28, offset: 1623},
						expr: &seqExpr{
							pos: position{line: 75, col: 28, offset: 1623},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 75, col: 28, offset: 1623},
									val:        "[^\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   true,
								},
								&anyMatcher{
									line: 75, col: 34, offset: 1629,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 77, col: 1, offset: 1635},
			expr: &ruleRefExpr{
				pos:  position{line: 77, col: 14, offset: 1648},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 79, col: 1, offset: 1664},
			expr: &actionExpr{
				pos: position{line: 79, col: 18, offset: 1681},
				run: (*parser).callonIdentifierName_1,
				expr: &seqExpr{
					pos: position{line: 79, col: 18, offset: 1681},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 79, col: 18, offset: 1681},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 79, col: 34, offset: 1697},
							expr: &ruleRefExpr{
								pos:  position{line: 79, col: 34, offset: 1697},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 83, col: 1, offset: 1741},
			expr: &charClassMatcher{
				pos:        position{line: 83, col: 19, offset: 1759},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 85, col: 1, offset: 1768},
			expr: &choiceExpr{
				pos: position{line: 85, col: 18, offset: 1785},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 85, col: 18, offset: 1785},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 85, col: 36, offset: 1803},
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
			pos:  position{line: 87, col: 1, offset: 1810},
			expr: &seqExpr{
				pos: position{line: 87, col: 14, offset: 1823},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 87, col: 14, offset: 1823},
						name: "StringLiteral",
					},
					&zeroOrOneExpr{
						pos: position{line: 87, col: 28, offset: 1837},
						expr: &litMatcher{
							pos:        position{line: 87, col: 28, offset: 1837},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 87, col: 33, offset: 1842},
						name: "sp",
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 89, col: 1, offset: 1846},
			expr: &choiceExpr{
				pos: position{line: 89, col: 17, offset: 1862},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 89, col: 17, offset: 1862},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 89, col: 17, offset: 1862},
								val:        "\"",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 89, col: 21, offset: 1866},
								expr: &ruleRefExpr{
									pos:  position{line: 89, col: 21, offset: 1866},
									name: "DoubleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 89, col: 39, offset: 1884},
								val:        "\"",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 89, col: 45, offset: 1890},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 89, col: 45, offset: 1890},
								val:        "'",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 89, col: 49, offset: 1894},
								expr: &ruleRefExpr{
									pos:  position{line: 89, col: 49, offset: 1894},
									name: "SingleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 89, col: 67, offset: 1912},
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
			pos:  position{line: 91, col: 1, offset: 1917},
			expr: &seqExpr{
				pos: position{line: 91, col: 20, offset: 1936},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 91, col: 20, offset: 1936},
						expr: &choiceExpr{
							pos: position{line: 91, col: 23, offset: 1939},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 91, col: 23, offset: 1939},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 91, col: 29, offset: 1945},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 91, col: 36, offset: 1952},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 91, col: 43, offset: 1959,
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 93, col: 1, offset: 1963},
			expr: &seqExpr{
				pos: position{line: 93, col: 20, offset: 1982},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 93, col: 20, offset: 1982},
						expr: &choiceExpr{
							pos: position{line: 93, col: 23, offset: 1985},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 93, col: 23, offset: 1985},
									val:        "'",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 93, col: 29, offset: 1991},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 93, col: 36, offset: 1998},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 93, col: 43, offset: 2005,
					},
				},
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 95, col: 1, offset: 2008},
			expr: &seqExpr{
				pos: position{line: 95, col: 20, offset: 2027},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 95, col: 20, offset: 2027},
						val:        "[",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 95, col: 24, offset: 2031},
						expr: &litMatcher{
							pos:        position{line: 95, col: 24, offset: 2031},
							val:        "^",
							ignoreCase: false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 95, col: 31, offset: 2038},
						expr: &choiceExpr{
							pos: position{line: 95, col: 31, offset: 2038},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 95, col: 31, offset: 2038},
									name: "ClassCharRange",
								},
								&ruleRefExpr{
									pos:  position{line: 95, col: 48, offset: 2055},
									name: "ClassChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 95, col: 61, offset: 2068},
						val:        "]",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 95, col: 65, offset: 2072},
						expr: &litMatcher{
							pos:        position{line: 95, col: 65, offset: 2072},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 70, offset: 2077},
						name: "sp",
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 97, col: 1, offset: 2081},
			expr: &seqExpr{
				pos: position{line: 97, col: 18, offset: 2098},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 97, col: 18, offset: 2098},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 97, col: 28, offset: 2108},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 97, col: 32, offset: 2112},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 99, col: 1, offset: 2123},
			expr: &seqExpr{
				pos: position{line: 99, col: 13, offset: 2135},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 99, col: 13, offset: 2135},
						expr: &choiceExpr{
							pos: position{line: 99, col: 16, offset: 2138},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 99, col: 16, offset: 2138},
									val:        "]",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 99, col: 22, offset: 2144},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 99, col: 29, offset: 2151},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 99, col: 36, offset: 2158,
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 101, col: 1, offset: 2162},
			expr: &litMatcher{
				pos:        position{line: 101, col: 14, offset: 2175},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 103, col: 1, offset: 2180},
			expr: &seqExpr{
				pos: position{line: 103, col: 13, offset: 2192},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 103, col: 13, offset: 2192},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 103, col: 17, offset: 2196},
						name: "sp",
					},
					&ruleRefExpr{
						pos:  position{line: 103, col: 20, offset: 2199},
						name: "Code",
					},
					&litMatcher{
						pos:        position{line: 103, col: 25, offset: 2204},
						val:        "}",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 103, col: 29, offset: 2208},
						name: "sp",
					},
				},
			},
		},
		{
			name: "sp",
			pos:  position{line: 105, col: 1, offset: 2212},
			expr: &zeroOrMoreExpr{
				pos: position{line: 105, col: 6, offset: 2217},
				expr: &charClassMatcher{
					pos:        position{line: 105, col: 6, offset: 2217},
					val:        "[ \\n\\r\\t]",
					chars:      []rune{' ', '\n', '\r', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EndOfRule",
			pos:  position{line: 107, col: 1, offset: 2229},
			expr: &choiceExpr{
				pos: position{line: 107, col: 13, offset: 2241},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 107, col: 13, offset: 2241},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 107, col: 13, offset: 2241},
								val:        ";",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 107, col: 17, offset: 2245},
								name: "sp",
							},
						},
					},
					&seqExpr{
						pos: position{line: 107, col: 22, offset: 2250},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 107, col: 22, offset: 2250},
								expr: &ruleRefExpr{
									pos:  position{line: 107, col: 22, offset: 2250},
									name: "SingleLineComment",
								},
							},
							&charClassMatcher{
								pos:        position{line: 107, col: 41, offset: 2269},
								val:        "[\\n]",
								chars:      []rune{'\n'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 107, col: 48, offset: 2276},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "EndOfFile",
			pos:  position{line: 109, col: 1, offset: 2287},
			expr: &notExpr{
				pos: position{line: 109, col: 13, offset: 2299},
				expr: &anyMatcher{
					line: 109, col: 14, offset: 2300,
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
