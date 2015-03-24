package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/exp/peg/ast"
)

func main() {
	fmt.Println(g)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 9, col: 1, offset: 54},
			expr: &actionExpr{
				pos: position{line: 9, col: 11, offset: 64},
				run: (*parser).callonGrammar_1,
				expr: &seqExpr{
					pos: position{line: 9, col: 11, offset: 64},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 9, col: 11, offset: 64},
							name: "sp",
						},
						&litMatcher{
							pos:        position{line: 9, col: 14, offset: 67},
							val:        "package",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 9, col: 24, offset: 77},
							label: "pkg",
							expr: &ruleRefExpr{
								pos:  position{line: 9, col: 28, offset: 81},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 9, col: 43, offset: 96},
							name: "sp",
						},
						&labeledExpr{
							pos:   position{line: 9, col: 46, offset: 99},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 9, col: 59, offset: 112},
								expr: &seqExpr{
									pos: position{line: 9, col: 59, offset: 112},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 9, col: 59, offset: 112},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 9, col: 71, offset: 124},
											name: "sp",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 9, col: 76, offset: 129},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 9, col: 83, offset: 136},
								expr: &seqExpr{
									pos: position{line: 9, col: 83, offset: 136},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 9, col: 83, offset: 136},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 9, col: 88, offset: 141},
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
			pos:  position{line: 18, col: 1, offset: 374},
			expr: &seqExpr{
				pos: position{line: 18, col: 8, offset: 381},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 18, col: 8, offset: 381},
						name: "IdentifierName",
					},
					&zeroOrOneExpr{
						pos: position{line: 18, col: 25, offset: 398},
						expr: &ruleRefExpr{
							pos:  position{line: 18, col: 25, offset: 398},
							name: "StringLiteral",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 18, col: 42, offset: 415},
						name: "RuleDefOp",
					},
					&ruleRefExpr{
						pos:  position{line: 18, col: 52, offset: 425},
						name: "Expression",
					},
					&ruleRefExpr{
						pos:  position{line: 18, col: 63, offset: 436},
						name: "EndOfRule",
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 20, col: 1, offset: 447},
			expr: &ruleRefExpr{
				pos:  position{line: 20, col: 14, offset: 460},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 22, col: 1, offset: 472},
			expr: &seqExpr{
				pos: position{line: 22, col: 14, offset: 485},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 22, col: 14, offset: 485},
						name: "ActionExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 22, col: 27, offset: 498},
						expr: &seqExpr{
							pos: position{line: 22, col: 27, offset: 498},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 22, col: 27, offset: 498},
									val:        "/",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 22, col: 31, offset: 502},
									name: "sp",
								},
								&ruleRefExpr{
									pos:  position{line: 22, col: 34, offset: 505},
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
			pos:  position{line: 24, col: 1, offset: 520},
			expr: &seqExpr{
				pos: position{line: 24, col: 14, offset: 533},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 24, col: 14, offset: 533},
						name: "SeqExpr",
					},
					&zeroOrOneExpr{
						pos: position{line: 24, col: 24, offset: 543},
						expr: &ruleRefExpr{
							pos:  position{line: 24, col: 24, offset: 543},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 26, col: 1, offset: 557},
			expr: &seqExpr{
				pos: position{line: 26, col: 11, offset: 567},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 26, col: 11, offset: 567},
						name: "LabeledExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 26, col: 25, offset: 581},
						expr: &ruleRefExpr{
							pos:  position{line: 26, col: 25, offset: 581},
							name: "LabeledExpr",
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 28, col: 1, offset: 597},
			expr: &choiceExpr{
				pos: position{line: 28, col: 15, offset: 611},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 28, col: 15, offset: 611},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 28, col: 15, offset: 611},
								name: "Identifier",
							},
							&litMatcher{
								pos:        position{line: 28, col: 26, offset: 622},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 28, col: 30, offset: 626},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 28, col: 33, offset: 629},
								name: "PrefixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 28, col: 48, offset: 644},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 30, col: 1, offset: 658},
			expr: &choiceExpr{
				pos: position{line: 30, col: 16, offset: 673},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 30, col: 16, offset: 673},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 30, col: 16, offset: 673},
								name: "PrefixedOp",
							},
							&ruleRefExpr{
								pos:  position{line: 30, col: 27, offset: 684},
								name: "SuffixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 30, col: 42, offset: 699},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 32, col: 1, offset: 713},
			expr: &seqExpr{
				pos: position{line: 32, col: 14, offset: 726},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 32, col: 16, offset: 728},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 32, col: 16, offset: 728},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 32, col: 22, offset: 734},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 32, col: 28, offset: 740},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 34, col: 1, offset: 744},
			expr: &choiceExpr{
				pos: position{line: 34, col: 16, offset: 759},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 34, col: 16, offset: 759},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 34, col: 16, offset: 759},
								name: "PrimaryExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 34, col: 28, offset: 771},
								name: "SuffixedOp",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 34, col: 41, offset: 784},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 36, col: 1, offset: 798},
			expr: &seqExpr{
				pos: position{line: 36, col: 14, offset: 811},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 36, col: 16, offset: 813},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 36, col: 16, offset: 813},
								val:        "?",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 36, col: 22, offset: 819},
								val:        "*",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 36, col: 28, offset: 825},
								val:        "+",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 36, col: 34, offset: 831},
						name: "sp",
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 38, col: 1, offset: 835},
			expr: &choiceExpr{
				pos: position{line: 38, col: 15, offset: 849},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 38, col: 15, offset: 849},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 38, col: 28, offset: 862},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 38, col: 47, offset: 881},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 38, col: 60, offset: 894},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 38, col: 74, offset: 908},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 38, col: 93, offset: 927},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 38, col: 93, offset: 927},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 38, col: 97, offset: 931},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 38, col: 100, offset: 934},
								name: "Expression",
							},
							&litMatcher{
								pos:        position{line: 38, col: 111, offset: 945},
								val:        ")",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 38, col: 115, offset: 949},
								name: "sp",
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 40, col: 1, offset: 953},
			expr: &seqExpr{
				pos: position{line: 40, col: 15, offset: 967},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 40, col: 15, offset: 967},
						name: "IdentifierName",
					},
					&notExpr{
						pos: position{line: 40, col: 30, offset: 982},
						expr: &seqExpr{
							pos: position{line: 40, col: 33, offset: 985},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 40, col: 35, offset: 987},
									expr: &ruleRefExpr{
										pos:  position{line: 40, col: 35, offset: 987},
										name: "StringLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 40, col: 52, offset: 1004},
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
			pos:  position{line: 42, col: 1, offset: 1011},
			expr: &seqExpr{
				pos: position{line: 42, col: 20, offset: 1030},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 42, col: 20, offset: 1030},
						name: "SemanticPredOp",
					},
					&ruleRefExpr{
						pos:  position{line: 42, col: 35, offset: 1045},
						name: "CodeBlock",
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 44, col: 1, offset: 1056},
			expr: &seqExpr{
				pos: position{line: 44, col: 18, offset: 1073},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 44, col: 20, offset: 1075},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 44, col: 20, offset: 1075},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 44, col: 26, offset: 1081},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 44, col: 32, offset: 1087},
						name: "sp",
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 46, col: 1, offset: 1091},
			expr: &seqExpr{
				pos: position{line: 46, col: 13, offset: 1103},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 46, col: 15, offset: 1105},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 46, col: 15, offset: 1105},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 46, col: 21, offset: 1111},
								val:        "<-",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 46, col: 28, offset: 1118},
								val:        "←",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 46, col: 39, offset: 1129},
								val:        "⟵",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 46, col: 50, offset: 1140},
						name: "sp",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 48, col: 1, offset: 1144},
			expr: &seqExpr{
				pos: position{line: 48, col: 20, offset: 1163},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 48, col: 20, offset: 1163},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 48, col: 27, offset: 1170},
						expr: &seqExpr{
							pos: position{line: 48, col: 27, offset: 1170},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 48, col: 27, offset: 1170},
									expr: &litMatcher{
										pos:        position{line: 48, col: 28, offset: 1171},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 48, col: 33, offset: 1176,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 48, col: 38, offset: 1181},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 43, offset: 1186},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 50, col: 1, offset: 1190},
			expr: &seqExpr{
				pos: position{line: 50, col: 21, offset: 1210},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 50, col: 21, offset: 1210},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 50, col: 28, offset: 1217},
						expr: &seqExpr{
							pos: position{line: 50, col: 28, offset: 1217},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 50, col: 28, offset: 1217},
									val:        "[^\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   true,
								},
								&anyMatcher{
									line: 50, col: 34, offset: 1223,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 52, col: 1, offset: 1229},
			expr: &ruleRefExpr{
				pos:  position{line: 52, col: 14, offset: 1242},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 54, col: 1, offset: 1258},
			expr: &seqExpr{
				pos: position{line: 54, col: 18, offset: 1275},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 54, col: 18, offset: 1275},
						name: "IdentifierStart",
					},
					&zeroOrMoreExpr{
						pos: position{line: 54, col: 34, offset: 1291},
						expr: &ruleRefExpr{
							pos:  position{line: 54, col: 34, offset: 1291},
							name: "IdentifierPart",
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 56, col: 1, offset: 1308},
			expr: &charClassMatcher{
				pos:        position{line: 56, col: 19, offset: 1326},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 58, col: 1, offset: 1335},
			expr: &choiceExpr{
				pos: position{line: 58, col: 18, offset: 1352},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 58, col: 18, offset: 1352},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 58, col: 36, offset: 1370},
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
			pos:  position{line: 60, col: 1, offset: 1377},
			expr: &seqExpr{
				pos: position{line: 60, col: 14, offset: 1390},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 60, col: 14, offset: 1390},
						name: "StringLiteral",
					},
					&zeroOrOneExpr{
						pos: position{line: 60, col: 28, offset: 1404},
						expr: &litMatcher{
							pos:        position{line: 60, col: 28, offset: 1404},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 60, col: 33, offset: 1409},
						name: "sp",
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 62, col: 1, offset: 1413},
			expr: &choiceExpr{
				pos: position{line: 62, col: 17, offset: 1429},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 62, col: 17, offset: 1429},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 62, col: 17, offset: 1429},
								val:        "\"",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 62, col: 21, offset: 1433},
								expr: &ruleRefExpr{
									pos:  position{line: 62, col: 21, offset: 1433},
									name: "DoubleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 62, col: 39, offset: 1451},
								val:        "\"",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 62, col: 45, offset: 1457},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 62, col: 45, offset: 1457},
								val:        "'",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 62, col: 49, offset: 1461},
								expr: &ruleRefExpr{
									pos:  position{line: 62, col: 49, offset: 1461},
									name: "SingleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 62, col: 67, offset: 1479},
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
			pos:  position{line: 64, col: 1, offset: 1484},
			expr: &seqExpr{
				pos: position{line: 64, col: 20, offset: 1503},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 64, col: 20, offset: 1503},
						expr: &choiceExpr{
							pos: position{line: 64, col: 23, offset: 1506},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 64, col: 23, offset: 1506},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 64, col: 29, offset: 1512},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 64, col: 36, offset: 1519},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 64, col: 43, offset: 1526,
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 66, col: 1, offset: 1530},
			expr: &seqExpr{
				pos: position{line: 66, col: 20, offset: 1549},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 66, col: 20, offset: 1549},
						expr: &choiceExpr{
							pos: position{line: 66, col: 23, offset: 1552},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 66, col: 23, offset: 1552},
									val:        "'",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 66, col: 29, offset: 1558},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 66, col: 36, offset: 1565},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 66, col: 43, offset: 1572,
					},
				},
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 68, col: 1, offset: 1575},
			expr: &seqExpr{
				pos: position{line: 68, col: 20, offset: 1594},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 68, col: 20, offset: 1594},
						val:        "[",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 68, col: 24, offset: 1598},
						expr: &litMatcher{
							pos:        position{line: 68, col: 24, offset: 1598},
							val:        "^",
							ignoreCase: false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 68, col: 31, offset: 1605},
						expr: &choiceExpr{
							pos: position{line: 68, col: 31, offset: 1605},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 68, col: 31, offset: 1605},
									name: "ClassCharRange",
								},
								&ruleRefExpr{
									pos:  position{line: 68, col: 48, offset: 1622},
									name: "ClassChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 68, col: 61, offset: 1635},
						val:        "]",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 68, col: 65, offset: 1639},
						expr: &litMatcher{
							pos:        position{line: 68, col: 65, offset: 1639},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 70, offset: 1644},
						name: "sp",
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 70, col: 1, offset: 1648},
			expr: &seqExpr{
				pos: position{line: 70, col: 18, offset: 1665},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 18, offset: 1665},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 70, col: 28, offset: 1675},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 70, col: 32, offset: 1679},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 72, col: 1, offset: 1690},
			expr: &seqExpr{
				pos: position{line: 72, col: 13, offset: 1702},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 72, col: 13, offset: 1702},
						expr: &choiceExpr{
							pos: position{line: 72, col: 16, offset: 1705},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 72, col: 16, offset: 1705},
									val:        "]",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 72, col: 22, offset: 1711},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 72, col: 29, offset: 1718},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 72, col: 36, offset: 1725,
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 74, col: 1, offset: 1729},
			expr: &litMatcher{
				pos:        position{line: 74, col: 14, offset: 1742},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 76, col: 1, offset: 1747},
			expr: &seqExpr{
				pos: position{line: 76, col: 13, offset: 1759},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 76, col: 13, offset: 1759},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 17, offset: 1763},
						name: "sp",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 20, offset: 1766},
						name: "Code",
					},
					&litMatcher{
						pos:        position{line: 76, col: 25, offset: 1771},
						val:        "}",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 29, offset: 1775},
						name: "sp",
					},
				},
			},
		},
		{
			name: "sp",
			pos:  position{line: 78, col: 1, offset: 1779},
			expr: &zeroOrMoreExpr{
				pos: position{line: 78, col: 6, offset: 1784},
				expr: &charClassMatcher{
					pos:        position{line: 78, col: 6, offset: 1784},
					val:        "[ \\n\\r\\t]",
					chars:      []rune{' ', '\n', '\r', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EndOfRule",
			pos:  position{line: 80, col: 1, offset: 1796},
			expr: &choiceExpr{
				pos: position{line: 80, col: 13, offset: 1808},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 80, col: 13, offset: 1808},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 80, col: 13, offset: 1808},
								val:        ";",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 80, col: 17, offset: 1812},
								name: "sp",
							},
						},
					},
					&seqExpr{
						pos: position{line: 80, col: 22, offset: 1817},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 80, col: 22, offset: 1817},
								expr: &ruleRefExpr{
									pos:  position{line: 80, col: 22, offset: 1817},
									name: "SingleLineComment",
								},
							},
							&charClassMatcher{
								pos:        position{line: 80, col: 41, offset: 1836},
								val:        "[\\n]",
								chars:      []rune{'\n'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 48, offset: 1843},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "EndOfFile",
			pos:  position{line: 82, col: 1, offset: 1854},
			expr: &notExpr{
				pos: position{line: 82, col: 13, offset: 1866},
				expr: &anyMatcher{
					line: 82, col: 14, offset: 1867,
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
	return p.cur.onGrammar_1(stack["pkg"], stack["initializer"], stack["rules"])
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
				err = fmt.Errorf("%!v(MISSING)", e)
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
		panic(fmt.Sprintf("unknown expression type %!T(MISSING)", expr))
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
		panic(fmt.Sprintf("%!s(MISSING): invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.errs.add(fmt.Errorf("undefined rule: %!s(MISSING)", ref.name))
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
