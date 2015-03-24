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
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 3, col: 1, offset: 14},
			expr: &seqExpr{
				pos: position{line: 3, col: 11, offset: 24},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 3, col: 11, offset: 24},
						name: "sp",
					},
					&oneOrMoreExpr{
						pos: position{line: 3, col: 14, offset: 27},
						expr: &ruleRefExpr{
							pos:  position{line: 3, col: 14, offset: 27},
							name: "Rule",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 3, col: 20, offset: 33},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 5, col: 1, offset: 44},
			expr: &seqExpr{
				pos: position{line: 5, col: 8, offset: 51},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 5, col: 8, offset: 51},
						name: "IdentifierName",
					},
					&zeroOrOneExpr{
						pos: position{line: 5, col: 25, offset: 68},
						expr: &ruleRefExpr{
							pos:  position{line: 5, col: 25, offset: 68},
							name: "StringLiteral",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 5, col: 42, offset: 85},
						name: "RuleDefOp",
					},
					&ruleRefExpr{
						pos:  position{line: 5, col: 52, offset: 95},
						name: "Expression",
					},
					&ruleRefExpr{
						pos:  position{line: 5, col: 63, offset: 106},
						name: "EndOfRule",
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 7, col: 1, offset: 117},
			expr: &ruleRefExpr{
				pos:  position{line: 7, col: 14, offset: 130},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 9, col: 1, offset: 142},
			expr: &seqExpr{
				pos: position{line: 9, col: 14, offset: 155},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 9, col: 14, offset: 155},
						name: "ActionExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 9, col: 27, offset: 168},
						expr: &seqExpr{
							pos: position{line: 9, col: 27, offset: 168},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 9, col: 27, offset: 168},
									val:        "/",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 9, col: 31, offset: 172},
									name: "sp",
								},
								&ruleRefExpr{
									pos:  position{line: 9, col: 34, offset: 175},
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
			pos:  position{line: 11, col: 1, offset: 190},
			expr: &seqExpr{
				pos: position{line: 11, col: 14, offset: 203},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 11, col: 14, offset: 203},
						name: "SeqExpr",
					},
					&zeroOrOneExpr{
						pos: position{line: 11, col: 24, offset: 213},
						expr: &ruleRefExpr{
							pos:  position{line: 11, col: 24, offset: 213},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 13, col: 1, offset: 227},
			expr: &seqExpr{
				pos: position{line: 13, col: 11, offset: 237},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 13, col: 11, offset: 237},
						name: "LabeledExpr",
					},
					&zeroOrMoreExpr{
						pos: position{line: 13, col: 25, offset: 251},
						expr: &ruleRefExpr{
							pos:  position{line: 13, col: 25, offset: 251},
							name: "LabeledExpr",
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 15, col: 1, offset: 267},
			expr: &choiceExpr{
				pos: position{line: 15, col: 15, offset: 281},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 15, col: 15, offset: 281},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 15, col: 15, offset: 281},
								name: "Identifier",
							},
							&litMatcher{
								pos:        position{line: 15, col: 26, offset: 292},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 15, col: 30, offset: 296},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 15, col: 33, offset: 299},
								name: "PrefixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 15, col: 48, offset: 314},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 17, col: 1, offset: 328},
			expr: &choiceExpr{
				pos: position{line: 17, col: 16, offset: 343},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 17, col: 16, offset: 343},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 17, col: 16, offset: 343},
								name: "PrefixedOp",
							},
							&ruleRefExpr{
								pos:  position{line: 17, col: 27, offset: 354},
								name: "SuffixedExpr",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 17, col: 42, offset: 369},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 19, col: 1, offset: 383},
			expr: &seqExpr{
				pos: position{line: 19, col: 14, offset: 396},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 19, col: 16, offset: 398},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 19, col: 16, offset: 398},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 19, col: 22, offset: 404},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 19, col: 28, offset: 410},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 21, col: 1, offset: 414},
			expr: &choiceExpr{
				pos: position{line: 21, col: 16, offset: 429},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 21, col: 16, offset: 429},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 21, col: 16, offset: 429},
								name: "PrimaryExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 21, col: 28, offset: 441},
								name: "SuffixedOp",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 21, col: 41, offset: 454},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 23, col: 1, offset: 468},
			expr: &seqExpr{
				pos: position{line: 23, col: 14, offset: 481},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 23, col: 16, offset: 483},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 23, col: 16, offset: 483},
								val:        "?",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 23, col: 22, offset: 489},
								val:        "*",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 23, col: 28, offset: 495},
								val:        "+",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 23, col: 34, offset: 501},
						name: "sp",
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 25, col: 1, offset: 505},
			expr: &choiceExpr{
				pos: position{line: 25, col: 15, offset: 519},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 25, col: 15, offset: 519},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 28, offset: 532},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 47, offset: 551},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 60, offset: 564},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 74, offset: 578},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 25, col: 93, offset: 597},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 25, col: 93, offset: 597},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 25, col: 97, offset: 601},
								name: "sp",
							},
							&ruleRefExpr{
								pos:  position{line: 25, col: 100, offset: 604},
								name: "Expression",
							},
							&litMatcher{
								pos:        position{line: 25, col: 111, offset: 615},
								val:        ")",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 25, col: 115, offset: 619},
								name: "sp",
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 27, col: 1, offset: 623},
			expr: &seqExpr{
				pos: position{line: 27, col: 15, offset: 637},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 27, col: 15, offset: 637},
						name: "IdentifierName",
					},
					&notExpr{
						pos: position{line: 27, col: 30, offset: 652},
						expr: &seqExpr{
							pos: position{line: 27, col: 33, offset: 655},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 27, col: 35, offset: 657},
									expr: &ruleRefExpr{
										pos:  position{line: 27, col: 35, offset: 657},
										name: "StringLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 27, col: 52, offset: 674},
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
			pos:  position{line: 29, col: 1, offset: 681},
			expr: &seqExpr{
				pos: position{line: 29, col: 20, offset: 700},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 29, col: 20, offset: 700},
						name: "SemanticPredOp",
					},
					&ruleRefExpr{
						pos:  position{line: 29, col: 35, offset: 715},
						name: "CodeBlock",
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 31, col: 1, offset: 726},
			expr: &seqExpr{
				pos: position{line: 31, col: 18, offset: 743},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 31, col: 20, offset: 745},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 31, col: 20, offset: 745},
								val:        "&",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 31, col: 26, offset: 751},
								val:        "!",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 31, col: 32, offset: 757},
						name: "sp",
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 33, col: 1, offset: 761},
			expr: &seqExpr{
				pos: position{line: 33, col: 13, offset: 773},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 33, col: 15, offset: 775},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 33, col: 15, offset: 775},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 33, col: 21, offset: 781},
								val:        "<-",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 33, col: 28, offset: 788},
								val:        "←",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 33, col: 39, offset: 799},
								val:        "⟵",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 50, offset: 810},
						name: "sp",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 35, col: 1, offset: 814},
			expr: &seqExpr{
				pos: position{line: 35, col: 20, offset: 833},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 35, col: 20, offset: 833},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 35, col: 27, offset: 840},
						expr: &seqExpr{
							pos: position{line: 35, col: 27, offset: 840},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 35, col: 27, offset: 840},
									expr: &litMatcher{
										pos:        position{line: 35, col: 28, offset: 841},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 35, col: 33, offset: 846,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 35, col: 38, offset: 851},
						val:        "*/",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 35, col: 43, offset: 856},
						name: "sp",
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 37, col: 1, offset: 860},
			expr: &seqExpr{
				pos: position{line: 37, col: 21, offset: 880},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 37, col: 21, offset: 880},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 37, col: 28, offset: 887},
						expr: &seqExpr{
							pos: position{line: 37, col: 28, offset: 887},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 37, col: 28, offset: 887},
									val:        "[^\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   true,
								},
								&anyMatcher{
									line: 37, col: 34, offset: 893,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 39, col: 1, offset: 899},
			expr: &ruleRefExpr{
				pos:  position{line: 39, col: 14, offset: 912},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 41, col: 1, offset: 928},
			expr: &seqExpr{
				pos: position{line: 41, col: 18, offset: 945},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 41, col: 18, offset: 945},
						name: "IdentifierStart",
					},
					&zeroOrMoreExpr{
						pos: position{line: 41, col: 34, offset: 961},
						expr: &ruleRefExpr{
							pos:  position{line: 41, col: 34, offset: 961},
							name: "IdentifierPart",
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 43, col: 1, offset: 978},
			expr: &charClassMatcher{
				pos:        position{line: 43, col: 19, offset: 996},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 45, col: 1, offset: 1005},
			expr: &choiceExpr{
				pos: position{line: 45, col: 18, offset: 1022},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 45, col: 18, offset: 1022},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 45, col: 36, offset: 1040},
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
			pos:  position{line: 47, col: 1, offset: 1047},
			expr: &seqExpr{
				pos: position{line: 47, col: 14, offset: 1060},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 47, col: 14, offset: 1060},
						name: "StringLiteral",
					},
					&zeroOrOneExpr{
						pos: position{line: 47, col: 28, offset: 1074},
						expr: &litMatcher{
							pos:        position{line: 47, col: 28, offset: 1074},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 33, offset: 1079},
						name: "sp",
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 49, col: 1, offset: 1083},
			expr: &choiceExpr{
				pos: position{line: 49, col: 17, offset: 1099},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 49, col: 17, offset: 1099},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 49, col: 17, offset: 1099},
								val:        "\"",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 49, col: 21, offset: 1103},
								expr: &ruleRefExpr{
									pos:  position{line: 49, col: 21, offset: 1103},
									name: "DoubleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 49, col: 39, offset: 1121},
								val:        "\"",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 49, col: 45, offset: 1127},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 49, col: 45, offset: 1127},
								val:        "'",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 49, col: 49, offset: 1131},
								expr: &ruleRefExpr{
									pos:  position{line: 49, col: 49, offset: 1131},
									name: "SingleStringChar",
								},
							},
							&litMatcher{
								pos:        position{line: 49, col: 67, offset: 1149},
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
			pos:  position{line: 51, col: 1, offset: 1154},
			expr: &seqExpr{
				pos: position{line: 51, col: 20, offset: 1173},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 51, col: 20, offset: 1173},
						expr: &choiceExpr{
							pos: position{line: 51, col: 23, offset: 1176},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 51, col: 23, offset: 1176},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 51, col: 29, offset: 1182},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 51, col: 36, offset: 1189},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 51, col: 43, offset: 1196,
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 53, col: 1, offset: 1200},
			expr: &seqExpr{
				pos: position{line: 53, col: 20, offset: 1219},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 53, col: 20, offset: 1219},
						expr: &choiceExpr{
							pos: position{line: 53, col: 23, offset: 1222},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 53, col: 23, offset: 1222},
									val:        "'",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 53, col: 29, offset: 1228},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 53, col: 36, offset: 1235},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 53, col: 43, offset: 1242,
					},
				},
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 55, col: 1, offset: 1245},
			expr: &seqExpr{
				pos: position{line: 55, col: 20, offset: 1264},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 55, col: 20, offset: 1264},
						val:        "[",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 55, col: 24, offset: 1268},
						expr: &litMatcher{
							pos:        position{line: 55, col: 24, offset: 1268},
							val:        "^",
							ignoreCase: false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 55, col: 31, offset: 1275},
						expr: &choiceExpr{
							pos: position{line: 55, col: 31, offset: 1275},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 55, col: 31, offset: 1275},
									name: "ClassCharRange",
								},
								&ruleRefExpr{
									pos:  position{line: 55, col: 48, offset: 1292},
									name: "ClassChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 55, col: 61, offset: 1305},
						val:        "]",
						ignoreCase: false,
					},
					&zeroOrOneExpr{
						pos: position{line: 55, col: 65, offset: 1309},
						expr: &litMatcher{
							pos:        position{line: 55, col: 65, offset: 1309},
							val:        "i",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 55, col: 70, offset: 1314},
						name: "sp",
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 57, col: 1, offset: 1318},
			expr: &seqExpr{
				pos: position{line: 57, col: 18, offset: 1335},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 57, col: 18, offset: 1335},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 57, col: 28, offset: 1345},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 32, offset: 1349},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 59, col: 1, offset: 1360},
			expr: &seqExpr{
				pos: position{line: 59, col: 13, offset: 1372},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 59, col: 13, offset: 1372},
						expr: &choiceExpr{
							pos: position{line: 59, col: 16, offset: 1375},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 59, col: 16, offset: 1375},
									val:        "]",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 59, col: 22, offset: 1381},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 59, col: 29, offset: 1388},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 59, col: 36, offset: 1395,
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 61, col: 1, offset: 1399},
			expr: &litMatcher{
				pos:        position{line: 61, col: 14, offset: 1412},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 63, col: 1, offset: 1417},
			expr: &seqExpr{
				pos: position{line: 63, col: 13, offset: 1429},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 63, col: 13, offset: 1429},
						val:        "{",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 17, offset: 1433},
						name: "sp",
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 20, offset: 1436},
						name: "Code",
					},
					&litMatcher{
						pos:        position{line: 63, col: 25, offset: 1441},
						val:        "}",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 29, offset: 1445},
						name: "sp",
					},
				},
			},
		},
		{
			name: "sp",
			pos:  position{line: 65, col: 1, offset: 1449},
			expr: &zeroOrMoreExpr{
				pos: position{line: 65, col: 6, offset: 1454},
				expr: &charClassMatcher{
					pos:        position{line: 65, col: 6, offset: 1454},
					val:        "[ \\n\\r\\t]",
					chars:      []rune{' ', '\n', '\r', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EndOfRule",
			pos:  position{line: 67, col: 1, offset: 1466},
			expr: &choiceExpr{
				pos: position{line: 67, col: 13, offset: 1478},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 67, col: 13, offset: 1478},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 67, col: 13, offset: 1478},
								val:        ";",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 67, col: 17, offset: 1482},
								name: "sp",
							},
						},
					},
					&seqExpr{
						pos: position{line: 67, col: 22, offset: 1487},
						exprs: []interface{}{
							&zeroOrOneExpr{
								pos: position{line: 67, col: 22, offset: 1487},
								expr: &ruleRefExpr{
									pos:  position{line: 67, col: 22, offset: 1487},
									name: "SingleLineComment",
								},
							},
							&charClassMatcher{
								pos:        position{line: 67, col: 41, offset: 1506},
								val:        "[\\n]",
								chars:      []rune{'\n'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 67, col: 48, offset: 1513},
						name: "EndOfFile",
					},
				},
			},
		},
		{
			name: "EndOfFile",
			pos:  position{line: 69, col: 1, offset: 1524},
			expr: &notExpr{
				pos: position{line: 69, col: 13, offset: 1536},
				expr: &anyMatcher{
					line: 69, col: 14, offset: 1537,
				},
			},
		},
	},
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
