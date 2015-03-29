package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/exp/peg/ast"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 5, col: 1, offset: 18},
			expr: &actionExpr{
				pos: position{line: 5, col: 11, offset: 30},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 5, col: 11, offset: 30},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 5, col: 11, offset: 30},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 5, col: 14, offset: 33},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 5, col: 26, offset: 45},
								expr: &seqExpr{
									pos: position{line: 5, col: 28, offset: 47},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 28, offset: 47},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 40, offset: 59},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 5, col: 46, offset: 65},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 5, col: 52, offset: 71},
								expr: &seqExpr{
									pos: position{line: 5, col: 54, offset: 73},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 54, offset: 73},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 59, offset: 78},
											name: "__",
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
			pos:  position{line: 24, col: 1, offset: 521},
			expr: &actionExpr{
				pos: position{line: 24, col: 15, offset: 537},
				run: (*parser).callonInitializer1,
				expr: &seqExpr{
					pos: position{line: 24, col: 15, offset: 537},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 24, col: 15, offset: 537},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 20, offset: 542},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 30, offset: 552},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 28, col: 1, offset: 582},
			expr: &actionExpr{
				pos: position{line: 28, col: 8, offset: 591},
				run: (*parser).callonRule1,
				expr: &seqExpr{
					pos: position{line: 28, col: 8, offset: 591},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 8, offset: 591},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 13, offset: 596},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 28, offset: 611},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 31, offset: 614},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 28, col: 39, offset: 622},
								expr: &seqExpr{
									pos: position{line: 28, col: 41, offset: 624},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 28, col: 41, offset: 624},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 55, offset: 638},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 644},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 71, offset: 654},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 74, offset: 657},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 79, offset: 662},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 90, offset: 673},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 41, col: 1, offset: 957},
			expr: &ruleRefExpr{
				pos:  position{line: 41, col: 14, offset: 972},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 43, col: 1, offset: 984},
			expr: &actionExpr{
				pos: position{line: 43, col: 14, offset: 999},
				run: (*parser).callonChoiceExpr1,
				expr: &seqExpr{
					pos: position{line: 43, col: 14, offset: 999},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 43, col: 14, offset: 999},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 20, offset: 1005},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 43, col: 31, offset: 1016},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 36, offset: 1021},
								expr: &seqExpr{
									pos: position{line: 43, col: 38, offset: 1023},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 43, col: 38, offset: 1023},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 43, col: 41, offset: 1026},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 45, offset: 1030},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 48, offset: 1033},
											name: "ActionExpr",
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
			name: "ActionExpr",
			pos:  position{line: 58, col: 1, offset: 1438},
			expr: &actionExpr{
				pos: position{line: 58, col: 14, offset: 1453},
				run: (*parser).callonActionExpr1,
				expr: &seqExpr{
					pos: position{line: 58, col: 14, offset: 1453},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 58, col: 14, offset: 1453},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 19, offset: 1458},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 58, col: 27, offset: 1466},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 58, col: 32, offset: 1471},
								expr: &seqExpr{
									pos: position{line: 58, col: 34, offset: 1473},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 58, col: 34, offset: 1473},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 58, col: 37, offset: 1476},
											name: "CodeBlock",
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
			name: "SeqExpr",
			pos:  position{line: 72, col: 1, offset: 1742},
			expr: &actionExpr{
				pos: position{line: 72, col: 11, offset: 1754},
				run: (*parser).callonSeqExpr1,
				expr: &seqExpr{
					pos: position{line: 72, col: 11, offset: 1754},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 72, col: 11, offset: 1754},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 17, offset: 1760},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 72, col: 29, offset: 1772},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 72, col: 34, offset: 1777},
								expr: &seqExpr{
									pos: position{line: 72, col: 36, offset: 1779},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 72, col: 36, offset: 1779},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 72, col: 39, offset: 1782},
											name: "LabeledExpr",
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
			name: "LabeledExpr",
			pos:  position{line: 85, col: 1, offset: 2133},
			expr: &choiceExpr{
				pos: position{line: 85, col: 15, offset: 2149},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 85, col: 15, offset: 2149},
						run: (*parser).callonLabeledExpr2,
						expr: &seqExpr{
							pos: position{line: 85, col: 15, offset: 2149},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 85, col: 15, offset: 2149},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 21, offset: 2155},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 32, offset: 2166},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 85, col: 35, offset: 2169},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 39, offset: 2173},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 85, col: 42, offset: 2176},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 47, offset: 2181},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 5, offset: 2354},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 93, col: 1, offset: 2368},
			expr: &choiceExpr{
				pos: position{line: 93, col: 16, offset: 2385},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 93, col: 16, offset: 2385},
						run: (*parser).callonPrefixedExpr2,
						expr: &seqExpr{
							pos: position{line: 93, col: 16, offset: 2385},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 93, col: 16, offset: 2385},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 19, offset: 2388},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 93, col: 30, offset: 2399},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 93, col: 33, offset: 2402},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 38, offset: 2407},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 5, offset: 2689},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 106, col: 1, offset: 2703},
			expr: &actionExpr{
				pos: position{line: 106, col: 14, offset: 2718},
				run: (*parser).callonPrefixedOp1,
				expr: &choiceExpr{
					pos: position{line: 106, col: 16, offset: 2720},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 16, offset: 2720},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 22, offset: 2726},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 110, col: 1, offset: 2768},
			expr: &choiceExpr{
				pos: position{line: 110, col: 16, offset: 2785},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 110, col: 16, offset: 2785},
						run: (*parser).callonSuffixedExpr2,
						expr: &seqExpr{
							pos: position{line: 110, col: 16, offset: 2785},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 110, col: 16, offset: 2785},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 21, offset: 2790},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 33, offset: 2802},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 36, offset: 2805},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 39, offset: 2808},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 129, col: 5, offset: 3338},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 131, col: 1, offset: 3352},
			expr: &actionExpr{
				pos: position{line: 131, col: 14, offset: 3367},
				run: (*parser).callonSuffixedOp1,
				expr: &choiceExpr{
					pos: position{line: 131, col: 16, offset: 3369},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 131, col: 16, offset: 3369},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 22, offset: 3375},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 28, offset: 3381},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 135, col: 1, offset: 3423},
			expr: &choiceExpr{
				pos: position{line: 135, col: 15, offset: 3439},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 135, col: 15, offset: 3439},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 28, offset: 3452},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 47, offset: 3471},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 60, offset: 3484},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 74, offset: 3498},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 135, col: 93, offset: 3517},
						run: (*parser).callonPrimaryExpr7,
						expr: &seqExpr{
							pos: position{line: 135, col: 93, offset: 3517},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 135, col: 93, offset: 3517},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 97, offset: 3521},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 135, col: 100, offset: 3524},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 105, offset: 3529},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 116, offset: 3540},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 135, col: 119, offset: 3543},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 138, col: 1, offset: 3572},
			expr: &actionExpr{
				pos: position{line: 138, col: 15, offset: 3588},
				run: (*parser).callonRuleRefExpr1,
				expr: &seqExpr{
					pos: position{line: 138, col: 15, offset: 3588},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 138, col: 15, offset: 3588},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 138, col: 20, offset: 3593},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 138, col: 35, offset: 3608},
							expr: &seqExpr{
								pos: position{line: 138, col: 38, offset: 3611},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 138, col: 38, offset: 3611},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 138, col: 41, offset: 3614},
										expr: &seqExpr{
											pos: position{line: 138, col: 43, offset: 3616},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 138, col: 43, offset: 3616},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 138, col: 57, offset: 3630},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 63, offset: 3636},
										name: "RuleDefOp",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredExpr",
			pos:  position{line: 143, col: 1, offset: 3752},
			expr: &actionExpr{
				pos: position{line: 143, col: 20, offset: 3773},
				run: (*parser).callonSemanticPredExpr1,
				expr: &seqExpr{
					pos: position{line: 143, col: 20, offset: 3773},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 143, col: 20, offset: 3773},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 23, offset: 3776},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 143, col: 38, offset: 3791},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 143, col: 41, offset: 3794},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 46, offset: 3799},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 154, col: 1, offset: 4076},
			expr: &actionExpr{
				pos: position{line: 154, col: 18, offset: 4095},
				run: (*parser).callonSemanticPredOp1,
				expr: &choiceExpr{
					pos: position{line: 154, col: 20, offset: 4097},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 20, offset: 4097},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 154, col: 26, offset: 4103},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 158, col: 1, offset: 4145},
			expr: &choiceExpr{
				pos: position{line: 158, col: 13, offset: 4159},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 158, col: 13, offset: 4159},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 19, offset: 4165},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 26, offset: 4172},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 37, offset: 4183},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 160, col: 1, offset: 4193},
			expr: &anyMatcher{
				line: 160, col: 14, offset: 4208,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 161, col: 1, offset: 4210},
			expr: &choiceExpr{
				pos: position{line: 161, col: 11, offset: 4222},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 161, col: 11, offset: 4222},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 161, col: 30, offset: 4241},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 162, col: 1, offset: 4259},
			expr: &seqExpr{
				pos: position{line: 162, col: 20, offset: 4280},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 162, col: 20, offset: 4280},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 162, col: 25, offset: 4285},
						expr: &seqExpr{
							pos: position{line: 162, col: 27, offset: 4287},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 162, col: 27, offset: 4287},
									expr: &litMatcher{
										pos:        position{line: 162, col: 28, offset: 4288},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 162, col: 33, offset: 4293},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 162, col: 47, offset: 4307},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 163, col: 1, offset: 4312},
			expr: &seqExpr{
				pos: position{line: 163, col: 36, offset: 4349},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 163, col: 36, offset: 4349},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 163, col: 41, offset: 4354},
						expr: &seqExpr{
							pos: position{line: 163, col: 43, offset: 4356},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 163, col: 43, offset: 4356},
									expr: &choiceExpr{
										pos: position{line: 163, col: 46, offset: 4359},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 163, col: 46, offset: 4359},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 163, col: 53, offset: 4366},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 163, col: 59, offset: 4372},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 163, col: 73, offset: 4386},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 164, col: 1, offset: 4391},
			expr: &seqExpr{
				pos: position{line: 164, col: 21, offset: 4413},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 164, col: 21, offset: 4413},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 164, col: 26, offset: 4418},
						expr: &seqExpr{
							pos: position{line: 164, col: 28, offset: 4420},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 164, col: 28, offset: 4420},
									expr: &ruleRefExpr{
										pos:  position{line: 164, col: 29, offset: 4421},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 164, col: 33, offset: 4425},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 166, col: 1, offset: 4440},
			expr: &choiceExpr{
				pos: position{line: 166, col: 14, offset: 4455},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 166, col: 14, offset: 4455},
						run: (*parser).callonIdentifier2,
						expr: &seqExpr{
							pos: position{line: 166, col: 14, offset: 4455},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 166, col: 14, offset: 4455},
									expr: &ruleRefExpr{
										pos:  position{line: 166, col: 15, offset: 4456},
										name: "ReservedWord",
									},
								},
								&labeledExpr{
									pos:   position{line: 166, col: 28, offset: 4469},
									label: "ident",
									expr: &ruleRefExpr{
										pos:  position{line: 166, col: 34, offset: 4475},
										name: "IdentifierName",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 168, col: 5, offset: 4518},
						run: (*parser).callonIdentifier8,
						expr: &ruleRefExpr{
							pos:  position{line: 168, col: 5, offset: 4518},
							name: "ReservedWord",
						},
					},
				},
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 171, col: 1, offset: 4637},
			expr: &actionExpr{
				pos: position{line: 171, col: 18, offset: 4656},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 171, col: 18, offset: 4656},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 171, col: 18, offset: 4656},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 171, col: 34, offset: 4672},
							expr: &ruleRefExpr{
								pos:  position{line: 171, col: 34, offset: 4672},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 174, col: 1, offset: 4754},
			expr: &charClassMatcher{
				pos:        position{line: 174, col: 19, offset: 4774},
				val:        "[\\pL_]",
				chars:      []rune{'_'},
				classes:    []*unicode.RangeTable{rangeTable("L")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 175, col: 1, offset: 4781},
			expr: &choiceExpr{
				pos: position{line: 175, col: 18, offset: 4800},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 175, col: 18, offset: 4800},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 175, col: 36, offset: 4818},
						val:        "[\\p{Nd}]",
						classes:    []*unicode.RangeTable{rangeTable("Nd")},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "LitMatcher",
			pos:  position{line: 177, col: 1, offset: 4828},
			expr: &actionExpr{
				pos: position{line: 177, col: 14, offset: 4843},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 177, col: 14, offset: 4843},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 177, col: 14, offset: 4843},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 177, col: 18, offset: 4847},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 177, col: 32, offset: 4861},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 177, col: 39, offset: 4868},
								expr: &litMatcher{
									pos:        position{line: 177, col: 39, offset: 4868},
									val:        "i",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 187, col: 1, offset: 5094},
			expr: &actionExpr{
				pos: position{line: 187, col: 17, offset: 5112},
				run: (*parser).callonStringLiteral1,
				expr: &choiceExpr{
					pos: position{line: 187, col: 19, offset: 5114},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 187, col: 19, offset: 5114},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 19, offset: 5114},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 187, col: 23, offset: 5118},
									expr: &ruleRefExpr{
										pos:  position{line: 187, col: 23, offset: 5118},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 187, col: 41, offset: 5136},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 187, col: 47, offset: 5142},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 47, offset: 5142},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 187, col: 51, offset: 5146},
									name: "SingleStringChar",
								},
								&litMatcher{
									pos:        position{line: 187, col: 68, offset: 5163},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 187, col: 74, offset: 5169},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 74, offset: 5169},
									val:        "`",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 187, col: 78, offset: 5173},
									name: "RawStringChar",
								},
								&litMatcher{
									pos:        position{line: 187, col: 92, offset: 5187},
									val:        "`",
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
			pos:  position{line: 190, col: 1, offset: 5258},
			expr: &choiceExpr{
				pos: position{line: 190, col: 20, offset: 5279},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 190, col: 20, offset: 5279},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 190, col: 20, offset: 5279},
								expr: &choiceExpr{
									pos: position{line: 190, col: 23, offset: 5282},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 190, col: 23, offset: 5282},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 190, col: 29, offset: 5288},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 190, col: 36, offset: 5295},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 42, offset: 5301},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 190, col: 55, offset: 5314},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 190, col: 55, offset: 5314},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 60, offset: 5319},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 191, col: 1, offset: 5338},
			expr: &choiceExpr{
				pos: position{line: 191, col: 20, offset: 5359},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 191, col: 20, offset: 5359},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 191, col: 20, offset: 5359},
								expr: &choiceExpr{
									pos: position{line: 191, col: 23, offset: 5362},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 191, col: 23, offset: 5362},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 191, col: 29, offset: 5368},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 191, col: 36, offset: 5375},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 191, col: 42, offset: 5381},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 191, col: 55, offset: 5394},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 191, col: 55, offset: 5394},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 191, col: 60, offset: 5399},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 192, col: 1, offset: 5418},
			expr: &seqExpr{
				pos: position{line: 192, col: 17, offset: 5436},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 192, col: 17, offset: 5436},
						expr: &litMatcher{
							pos:        position{line: 192, col: 18, offset: 5437},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 192, col: 22, offset: 5441},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 194, col: 1, offset: 5453},
			expr: &choiceExpr{
				pos: position{line: 194, col: 22, offset: 5476},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 194, col: 22, offset: 5476},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 194, col: 28, offset: 5482},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 195, col: 1, offset: 5503},
			expr: &choiceExpr{
				pos: position{line: 195, col: 22, offset: 5526},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 195, col: 22, offset: 5526},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 195, col: 28, offset: 5532},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 197, col: 1, offset: 5554},
			expr: &choiceExpr{
				pos: position{line: 197, col: 24, offset: 5579},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 197, col: 24, offset: 5579},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 43, offset: 5598},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 57, offset: 5612},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 69, offset: 5624},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 89, offset: 5644},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 198, col: 1, offset: 5663},
			expr: &choiceExpr{
				pos: position{line: 198, col: 20, offset: 5684},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 198, col: 20, offset: 5684},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 26, offset: 5690},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 32, offset: 5696},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 38, offset: 5702},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 44, offset: 5708},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 50, offset: 5714},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 56, offset: 5720},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 62, offset: 5726},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 199, col: 1, offset: 5731},
			expr: &seqExpr{
				pos: position{line: 199, col: 15, offset: 5747},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 199, col: 15, offset: 5747},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 26, offset: 5758},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 37, offset: 5769},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 200, col: 1, offset: 5780},
			expr: &seqExpr{
				pos: position{line: 200, col: 13, offset: 5794},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 200, col: 13, offset: 5794},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 17, offset: 5798},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 26, offset: 5807},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 201, col: 1, offset: 5816},
			expr: &seqExpr{
				pos: position{line: 201, col: 21, offset: 5838},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 201, col: 21, offset: 5838},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 25, offset: 5842},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 34, offset: 5851},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 43, offset: 5860},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 52, offset: 5869},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 61, offset: 5878},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 70, offset: 5887},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 79, offset: 5896},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 88, offset: 5905},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 202, col: 1, offset: 5914},
			expr: &seqExpr{
				pos: position{line: 202, col: 22, offset: 5937},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 202, col: 22, offset: 5937},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 26, offset: 5941},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 35, offset: 5950},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 44, offset: 5959},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 53, offset: 5968},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 204, col: 1, offset: 5978},
			expr: &charClassMatcher{
				pos:        position{line: 204, col: 14, offset: 5993},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 205, col: 1, offset: 5999},
			expr: &charClassMatcher{
				pos:        position{line: 205, col: 16, offset: 6016},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 206, col: 1, offset: 6022},
			expr: &charClassMatcher{
				pos:        position{line: 206, col: 12, offset: 6035},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 208, col: 1, offset: 6046},
			expr: &actionExpr{
				pos: position{line: 208, col: 20, offset: 6067},
				run: (*parser).callonCharClassMatcher1,
				expr: &seqExpr{
					pos: position{line: 208, col: 20, offset: 6067},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 208, col: 20, offset: 6067},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 208, col: 24, offset: 6071},
							expr: &choiceExpr{
								pos: position{line: 208, col: 26, offset: 6073},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 208, col: 26, offset: 6073},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 208, col: 43, offset: 6090},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 208, col: 55, offset: 6102},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 208, col: 55, offset: 6102},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 208, col: 60, offset: 6107},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 208, col: 82, offset: 6129},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 208, col: 86, offset: 6133},
							expr: &litMatcher{
								pos:        position{line: 208, col: 86, offset: 6133},
								val:        "i",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 213, col: 1, offset: 6238},
			expr: &seqExpr{
				pos: position{line: 213, col: 18, offset: 6257},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 213, col: 18, offset: 6257},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 213, col: 28, offset: 6267},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 213, col: 32, offset: 6271},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 214, col: 1, offset: 6281},
			expr: &choiceExpr{
				pos: position{line: 214, col: 13, offset: 6295},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 214, col: 13, offset: 6295},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 214, col: 13, offset: 6295},
								expr: &choiceExpr{
									pos: position{line: 214, col: 16, offset: 6298},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 214, col: 16, offset: 6298},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 214, col: 22, offset: 6304},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 214, col: 29, offset: 6311},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 35, offset: 6317},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 214, col: 48, offset: 6330},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 214, col: 48, offset: 6330},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 53, offset: 6335},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 215, col: 1, offset: 6351},
			expr: &choiceExpr{
				pos: position{line: 215, col: 19, offset: 6371},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 215, col: 19, offset: 6371},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 215, col: 25, offset: 6377},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 217, col: 1, offset: 6399},
			expr: &seqExpr{
				pos: position{line: 217, col: 22, offset: 6422},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 217, col: 22, offset: 6422},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 218, col: 7, offset: 6435},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 218, col: 7, offset: 6435},
								name: "SingleCharUnicodeClass",
							},
							&seqExpr{
								pos: position{line: 219, col: 7, offset: 6464},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 219, col: 7, offset: 6464},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 219, col: 11, offset: 6468},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 219, col: 24, offset: 6481},
										val:        "}",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SingleCharUnicodeClass",
			pos:  position{line: 221, col: 1, offset: 6491},
			expr: &charClassMatcher{
				pos:        position{line: 221, col: 26, offset: 6518},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 224, col: 1, offset: 6591},
			expr: &choiceExpr{
				pos: position{line: 224, col: 16, offset: 6608},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 224, col: 16, offset: 6608},
						val:        "Other_Default_Ignorable_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 225, col: 7, offset: 6651},
						val:        "Noncharacter_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 226, col: 7, offset: 6683},
						val:        "Logical_Order_Exception",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 227, col: 7, offset: 6715},
						val:        "Inscriptional_Parthian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 7, offset: 6746},
						val:        "Other_Grapheme_Extend",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 229, col: 7, offset: 6776},
						val:        "Inscriptional_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 230, col: 7, offset: 6806},
						val:        "Terminal_Punctuation",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 231, col: 7, offset: 6835},
						val:        "Meroitic_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 232, col: 7, offset: 6864},
						val:        "IDS_Trinary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 233, col: 7, offset: 6893},
						val:        "Egyptian_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 234, col: 7, offset: 6922},
						val:        "Pattern_White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 7, offset: 6950},
						val:        "IDS_Binary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 236, col: 7, offset: 6978},
						val:        "Canadian_Aboriginal",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 237, col: 7, offset: 7006},
						val:        "Variation_Selector",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 238, col: 7, offset: 7033},
						val:        "Caucasian_Albanian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 239, col: 7, offset: 7060},
						val:        "Unified_Ideograph",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 240, col: 7, offset: 7086},
						val:        "Other_ID_Continue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 241, col: 7, offset: 7112},
						val:        "Old_South_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 242, col: 7, offset: 7138},
						val:        "Old_North_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 243, col: 7, offset: 7164},
						val:        "Other_Alphabetic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 244, col: 7, offset: 7189},
						val:        "Meroitic_Cursive",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 245, col: 7, offset: 7214},
						val:        "Imperial_Aramaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 246, col: 7, offset: 7239},
						val:        "Psalter_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 247, col: 7, offset: 7263},
						val:        "Other_Uppercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 248, col: 7, offset: 7287},
						val:        "Other_Lowercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 249, col: 7, offset: 7311},
						val:        "ASCII_Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 250, col: 7, offset: 7335},
						val:        "Quotation_Mark",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 251, col: 7, offset: 7358},
						val:        "Pattern_Syntax",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 252, col: 7, offset: 7381},
						val:        "Other_ID_Start",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 253, col: 7, offset: 7404},
						val:        "Mende_Kikakui",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 254, col: 7, offset: 7426},
						val:        "Syloti_Nagri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 255, col: 7, offset: 7447},
						val:        "Sora_Sompeng",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 7, offset: 7468},
						val:        "Pahawh_Hmong",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 257, col: 7, offset: 7489},
						val:        "Meetei_Mayek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 258, col: 7, offset: 7510},
						val:        "Join_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 259, col: 7, offset: 7531},
						val:        "Bidi_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 260, col: 7, offset: 7552},
						val:        "White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 261, col: 7, offset: 7572},
						val:        "Warang_Citi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 7, offset: 7592},
						val:        "Soft_Dotted",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 263, col: 7, offset: 7612},
						val:        "Pau_Cin_Hau",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 264, col: 7, offset: 7632},
						val:        "Old_Persian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 265, col: 7, offset: 7652},
						val:        "New_Tai_Lue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 266, col: 7, offset: 7672},
						val:        "Ideographic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 267, col: 7, offset: 7692},
						val:        "Saurashtra",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 268, col: 7, offset: 7711},
						val:        "Phoenician",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 269, col: 7, offset: 7730},
						val:        "Other_Math",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 270, col: 7, offset: 7749},
						val:        "Old_Turkic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 271, col: 7, offset: 7768},
						val:        "Old_Permic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 272, col: 7, offset: 7787},
						val:        "Old_Italic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 273, col: 7, offset: 7806},
						val:        "Manichaean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 274, col: 7, offset: 7825},
						val:        "Kharoshthi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 275, col: 7, offset: 7844},
						val:        "Glagolitic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 276, col: 7, offset: 7863},
						val:        "Devanagari",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 277, col: 7, offset: 7882},
						val:        "Deprecated",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 278, col: 7, offset: 7901},
						val:        "Sundanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 7, offset: 7919},
						val:        "Samaritan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 280, col: 7, offset: 7937},
						val:        "Palmyrene",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 281, col: 7, offset: 7955},
						val:        "Nabataean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 282, col: 7, offset: 7973},
						val:        "Mongolian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 283, col: 7, offset: 7991},
						val:        "Malayalam",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 284, col: 7, offset: 8009},
						val:        "Khudawadi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 7, offset: 8027},
						val:        "Inherited",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 286, col: 7, offset: 8045},
						val:        "Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 287, col: 7, offset: 8063},
						val:        "Diacritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 288, col: 7, offset: 8081},
						val:        "Cuneiform",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 289, col: 7, offset: 8099},
						val:        "Bassa_Vah",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 290, col: 7, offset: 8117},
						val:        "Ugaritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 291, col: 7, offset: 8134},
						val:        "Tifinagh",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 292, col: 7, offset: 8151},
						val:        "Tai_Viet",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 293, col: 7, offset: 8168},
						val:        "Tai_Tham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 294, col: 7, offset: 8185},
						val:        "Tagbanwa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 295, col: 7, offset: 8202},
						val:        "Phags_Pa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 296, col: 7, offset: 8219},
						val:        "Ol_Chiki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 297, col: 7, offset: 8236},
						val:        "Mahajani",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 298, col: 7, offset: 8253},
						val:        "Linear_B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 299, col: 7, offset: 8270},
						val:        "Linear_A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 300, col: 7, offset: 8287},
						val:        "Kayah_Li",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 301, col: 7, offset: 8304},
						val:        "Katakana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 302, col: 7, offset: 8321},
						val:        "Javanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 303, col: 7, offset: 8338},
						val:        "Hiragana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 304, col: 7, offset: 8355},
						val:        "Gurmukhi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 305, col: 7, offset: 8372},
						val:        "Gujarati",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 306, col: 7, offset: 8389},
						val:        "Georgian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 307, col: 7, offset: 8406},
						val:        "Extender",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 308, col: 7, offset: 8423},
						val:        "Ethiopic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 309, col: 7, offset: 8440},
						val:        "Duployan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 310, col: 7, offset: 8457},
						val:        "Cyrillic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 311, col: 7, offset: 8474},
						val:        "Cherokee",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 312, col: 7, offset: 8491},
						val:        "Buginese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 313, col: 7, offset: 8508},
						val:        "Bopomofo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 314, col: 7, offset: 8525},
						val:        "Balinese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 315, col: 7, offset: 8542},
						val:        "Armenian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 316, col: 7, offset: 8559},
						val:        "Tirhuta",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 317, col: 7, offset: 8575},
						val:        "Tibetan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 318, col: 7, offset: 8591},
						val:        "Tagalog",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 319, col: 7, offset: 8607},
						val:        "Sinhala",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 320, col: 7, offset: 8623},
						val:        "Siddham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 321, col: 7, offset: 8639},
						val:        "Shavian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 322, col: 7, offset: 8655},
						val:        "Sharada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 323, col: 7, offset: 8671},
						val:        "Radical",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 324, col: 7, offset: 8687},
						val:        "Osmanya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 325, col: 7, offset: 8703},
						val:        "Myanmar",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 326, col: 7, offset: 8719},
						val:        "Mandaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 327, col: 7, offset: 8735},
						val:        "Kannada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 328, col: 7, offset: 8751},
						val:        "Hanunoo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 329, col: 7, offset: 8767},
						val:        "Grantha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 330, col: 7, offset: 8783},
						val:        "Elbasan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 331, col: 7, offset: 8799},
						val:        "Deseret",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 332, col: 7, offset: 8815},
						val:        "Cypriot",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 333, col: 7, offset: 8831},
						val:        "Braille",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 334, col: 7, offset: 8847},
						val:        "Bengali",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 335, col: 7, offset: 8863},
						val:        "Avestan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 336, col: 7, offset: 8879},
						val:        "Thaana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 337, col: 7, offset: 8894},
						val:        "Telugu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 338, col: 7, offset: 8909},
						val:        "Tai_Le",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 339, col: 7, offset: 8924},
						val:        "Syriac",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 340, col: 7, offset: 8939},
						val:        "Rejang",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 341, col: 7, offset: 8954},
						val:        "Lydian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 342, col: 7, offset: 8969},
						val:        "Lycian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 343, col: 7, offset: 8984},
						val:        "Lepcha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 344, col: 7, offset: 8999},
						val:        "Khojki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 345, col: 7, offset: 9014},
						val:        "Kaithi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 346, col: 7, offset: 9029},
						val:        "Hyphen",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 347, col: 7, offset: 9044},
						val:        "Hebrew",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 348, col: 7, offset: 9059},
						val:        "Hangul",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 349, col: 7, offset: 9074},
						val:        "Gothic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 350, col: 7, offset: 9089},
						val:        "Coptic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 351, col: 7, offset: 9104},
						val:        "Common",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 352, col: 7, offset: 9119},
						val:        "Chakma",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 353, col: 7, offset: 9134},
						val:        "Carian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 354, col: 7, offset: 9149},
						val:        "Brahmi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 355, col: 7, offset: 9164},
						val:        "Arabic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 356, col: 7, offset: 9179},
						val:        "Tamil",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 357, col: 7, offset: 9193},
						val:        "Takri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 358, col: 7, offset: 9207},
						val:        "STerm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 359, col: 7, offset: 9221},
						val:        "Runic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 360, col: 7, offset: 9235},
						val:        "Oriya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 361, col: 7, offset: 9249},
						val:        "Ogham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 362, col: 7, offset: 9263},
						val:        "Limbu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 363, col: 7, offset: 9277},
						val:        "Latin",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 364, col: 7, offset: 9291},
						val:        "Khmer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 365, col: 7, offset: 9305},
						val:        "Greek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 366, col: 7, offset: 9319},
						val:        "Buhid",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 367, col: 7, offset: 9333},
						val:        "Batak",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 368, col: 7, offset: 9347},
						val:        "Bamum",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 369, col: 7, offset: 9361},
						val:        "Thai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 370, col: 7, offset: 9374},
						val:        "Modi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 371, col: 7, offset: 9387},
						val:        "Miao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 372, col: 7, offset: 9400},
						val:        "Lisu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 373, col: 7, offset: 9413},
						val:        "Dash",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 374, col: 7, offset: 9426},
						val:        "Cham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 375, col: 7, offset: 9439},
						val:        "Vai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 376, col: 7, offset: 9451},
						val:        "Nko",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 377, col: 7, offset: 9463},
						val:        "Mro",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 378, col: 7, offset: 9475},
						val:        "Lao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 379, col: 7, offset: 9487},
						val:        "Han",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 380, col: 7, offset: 9499},
						val:        "Zs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 381, col: 7, offset: 9510},
						val:        "Zp",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 382, col: 7, offset: 9521},
						val:        "Zl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 383, col: 7, offset: 9532},
						val:        "Yi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 384, col: 7, offset: 9543},
						val:        "So",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 385, col: 7, offset: 9554},
						val:        "Sm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 386, col: 7, offset: 9565},
						val:        "Sk",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 387, col: 7, offset: 9576},
						val:        "Sc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 388, col: 7, offset: 9587},
						val:        "Ps",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 389, col: 7, offset: 9598},
						val:        "Po",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 390, col: 7, offset: 9609},
						val:        "Pi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 7, offset: 9620},
						val:        "Pf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 392, col: 7, offset: 9631},
						val:        "Pe",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 393, col: 7, offset: 9642},
						val:        "Pd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 394, col: 7, offset: 9653},
						val:        "Pc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 395, col: 7, offset: 9664},
						val:        "No",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 396, col: 7, offset: 9675},
						val:        "Nl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 7, offset: 9686},
						val:        "Nd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 398, col: 7, offset: 9697},
						val:        "Mn",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 399, col: 7, offset: 9708},
						val:        "Me",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 400, col: 7, offset: 9719},
						val:        "Mc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 401, col: 7, offset: 9730},
						val:        "Lu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 402, col: 7, offset: 9741},
						val:        "Lt",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 403, col: 7, offset: 9752},
						val:        "Lo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 404, col: 7, offset: 9763},
						val:        "Lm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 405, col: 7, offset: 9774},
						val:        "Ll",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 406, col: 7, offset: 9785},
						val:        "Cs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 407, col: 7, offset: 9796},
						val:        "Co",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 408, col: 7, offset: 9807},
						val:        "Cf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 409, col: 7, offset: 9818},
						val:        "Cc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 410, col: 7, offset: 9829},
						val:        "Z",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 411, col: 7, offset: 9839},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 412, col: 7, offset: 9849},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 413, col: 7, offset: 9859},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 414, col: 7, offset: 9869},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 415, col: 7, offset: 9879},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 416, col: 7, offset: 9889},
						val:        "C",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "ReservedWord",
			pos:  position{line: 418, col: 1, offset: 9894},
			expr: &choiceExpr{
				pos: position{line: 421, col: 2, offset: 9965},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 421, col: 2, offset: 9965},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 421, col: 2, offset: 9965},
								val:        "break",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 421, col: 10, offset: 9973},
								expr: &ruleRefExpr{
									pos:  position{line: 421, col: 11, offset: 9974},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 422, col: 4, offset: 9992},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 422, col: 4, offset: 9992},
								val:        "case",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 422, col: 11, offset: 9999},
								expr: &ruleRefExpr{
									pos:  position{line: 422, col: 12, offset: 10000},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 423, col: 4, offset: 10018},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 423, col: 4, offset: 10018},
								val:        "chan",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 423, col: 11, offset: 10025},
								expr: &ruleRefExpr{
									pos:  position{line: 423, col: 12, offset: 10026},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 424, col: 4, offset: 10044},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 424, col: 4, offset: 10044},
								val:        "const",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 424, col: 12, offset: 10052},
								expr: &ruleRefExpr{
									pos:  position{line: 424, col: 13, offset: 10053},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 425, col: 4, offset: 10071},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 425, col: 4, offset: 10071},
								val:        "continue",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 425, col: 15, offset: 10082},
								expr: &ruleRefExpr{
									pos:  position{line: 425, col: 16, offset: 10083},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 426, col: 4, offset: 10101},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 426, col: 4, offset: 10101},
								val:        "default",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 426, col: 14, offset: 10111},
								expr: &ruleRefExpr{
									pos:  position{line: 426, col: 15, offset: 10112},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 427, col: 4, offset: 10130},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 427, col: 4, offset: 10130},
								val:        "defer",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 427, col: 12, offset: 10138},
								expr: &ruleRefExpr{
									pos:  position{line: 427, col: 13, offset: 10139},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 428, col: 4, offset: 10157},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 428, col: 4, offset: 10157},
								val:        "else",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 428, col: 11, offset: 10164},
								expr: &ruleRefExpr{
									pos:  position{line: 428, col: 12, offset: 10165},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 429, col: 4, offset: 10183},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 429, col: 4, offset: 10183},
								val:        "fallthrough",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 429, col: 18, offset: 10197},
								expr: &ruleRefExpr{
									pos:  position{line: 429, col: 19, offset: 10198},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 430, col: 4, offset: 10216},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 430, col: 4, offset: 10216},
								val:        "for",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 430, col: 10, offset: 10222},
								expr: &ruleRefExpr{
									pos:  position{line: 430, col: 11, offset: 10223},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 431, col: 4, offset: 10241},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 431, col: 4, offset: 10241},
								val:        "func",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 431, col: 11, offset: 10248},
								expr: &ruleRefExpr{
									pos:  position{line: 431, col: 12, offset: 10249},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 432, col: 4, offset: 10267},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 432, col: 4, offset: 10267},
								val:        "goto",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 432, col: 11, offset: 10274},
								expr: &ruleRefExpr{
									pos:  position{line: 432, col: 12, offset: 10275},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 433, col: 4, offset: 10293},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 433, col: 4, offset: 10293},
								val:        "go",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 433, col: 9, offset: 10298},
								expr: &ruleRefExpr{
									pos:  position{line: 433, col: 10, offset: 10299},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 434, col: 4, offset: 10317},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 434, col: 4, offset: 10317},
								val:        "if",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 434, col: 9, offset: 10322},
								expr: &ruleRefExpr{
									pos:  position{line: 434, col: 10, offset: 10323},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 435, col: 4, offset: 10341},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 435, col: 4, offset: 10341},
								val:        "import",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 435, col: 13, offset: 10350},
								expr: &ruleRefExpr{
									pos:  position{line: 435, col: 14, offset: 10351},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 436, col: 4, offset: 10369},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 436, col: 4, offset: 10369},
								val:        "interface",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 436, col: 16, offset: 10381},
								expr: &ruleRefExpr{
									pos:  position{line: 436, col: 17, offset: 10382},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 437, col: 4, offset: 10400},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 437, col: 4, offset: 10400},
								val:        "map",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 437, col: 10, offset: 10406},
								expr: &ruleRefExpr{
									pos:  position{line: 437, col: 11, offset: 10407},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 438, col: 4, offset: 10425},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 438, col: 4, offset: 10425},
								val:        "package",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 438, col: 14, offset: 10435},
								expr: &ruleRefExpr{
									pos:  position{line: 438, col: 15, offset: 10436},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 439, col: 4, offset: 10454},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 439, col: 4, offset: 10454},
								val:        "range",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 439, col: 12, offset: 10462},
								expr: &ruleRefExpr{
									pos:  position{line: 439, col: 13, offset: 10463},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 440, col: 4, offset: 10481},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 440, col: 4, offset: 10481},
								val:        "return",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 440, col: 13, offset: 10490},
								expr: &ruleRefExpr{
									pos:  position{line: 440, col: 14, offset: 10491},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 441, col: 4, offset: 10509},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 441, col: 4, offset: 10509},
								val:        "select",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 441, col: 13, offset: 10518},
								expr: &ruleRefExpr{
									pos:  position{line: 441, col: 14, offset: 10519},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 442, col: 4, offset: 10537},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 442, col: 4, offset: 10537},
								val:        "struct",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 442, col: 13, offset: 10546},
								expr: &ruleRefExpr{
									pos:  position{line: 442, col: 14, offset: 10547},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 443, col: 4, offset: 10565},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 443, col: 4, offset: 10565},
								val:        "switch",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 443, col: 13, offset: 10574},
								expr: &ruleRefExpr{
									pos:  position{line: 443, col: 14, offset: 10575},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 444, col: 4, offset: 10593},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 444, col: 4, offset: 10593},
								val:        "type",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 444, col: 11, offset: 10600},
								expr: &ruleRefExpr{
									pos:  position{line: 444, col: 12, offset: 10601},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 445, col: 4, offset: 10619},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 445, col: 4, offset: 10619},
								val:        "var",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 445, col: 10, offset: 10625},
								expr: &ruleRefExpr{
									pos:  position{line: 445, col: 11, offset: 10626},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 449, col: 4, offset: 10725},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 449, col: 4, offset: 10725},
								val:        "bool",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 449, col: 11, offset: 10732},
								expr: &ruleRefExpr{
									pos:  position{line: 449, col: 12, offset: 10733},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 450, col: 4, offset: 10751},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 450, col: 4, offset: 10751},
								val:        "byte",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 450, col: 11, offset: 10758},
								expr: &ruleRefExpr{
									pos:  position{line: 450, col: 12, offset: 10759},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 451, col: 4, offset: 10777},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 451, col: 4, offset: 10777},
								val:        "complex64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 451, col: 16, offset: 10789},
								expr: &ruleRefExpr{
									pos:  position{line: 451, col: 17, offset: 10790},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 452, col: 4, offset: 10808},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 452, col: 4, offset: 10808},
								val:        "complex128",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 452, col: 17, offset: 10821},
								expr: &ruleRefExpr{
									pos:  position{line: 452, col: 18, offset: 10822},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 453, col: 4, offset: 10840},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 453, col: 4, offset: 10840},
								val:        "error",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 453, col: 12, offset: 10848},
								expr: &ruleRefExpr{
									pos:  position{line: 453, col: 13, offset: 10849},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 454, col: 4, offset: 10867},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 454, col: 4, offset: 10867},
								val:        "float32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 454, col: 14, offset: 10877},
								expr: &ruleRefExpr{
									pos:  position{line: 454, col: 15, offset: 10878},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 455, col: 4, offset: 10896},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 455, col: 4, offset: 10896},
								val:        "float64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 455, col: 14, offset: 10906},
								expr: &ruleRefExpr{
									pos:  position{line: 455, col: 15, offset: 10907},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 456, col: 4, offset: 10925},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 456, col: 4, offset: 10925},
								val:        "int8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 456, col: 11, offset: 10932},
								expr: &ruleRefExpr{
									pos:  position{line: 456, col: 12, offset: 10933},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 457, col: 4, offset: 10951},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 457, col: 4, offset: 10951},
								val:        "int16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 457, col: 12, offset: 10959},
								expr: &ruleRefExpr{
									pos:  position{line: 457, col: 13, offset: 10960},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 458, col: 4, offset: 10978},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 458, col: 4, offset: 10978},
								val:        "int32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 458, col: 12, offset: 10986},
								expr: &ruleRefExpr{
									pos:  position{line: 458, col: 13, offset: 10987},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 459, col: 4, offset: 11005},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 459, col: 4, offset: 11005},
								val:        "int64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 459, col: 12, offset: 11013},
								expr: &ruleRefExpr{
									pos:  position{line: 459, col: 13, offset: 11014},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 460, col: 4, offset: 11032},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 460, col: 4, offset: 11032},
								val:        "int",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 460, col: 10, offset: 11038},
								expr: &ruleRefExpr{
									pos:  position{line: 460, col: 11, offset: 11039},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 461, col: 4, offset: 11057},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 461, col: 4, offset: 11057},
								val:        "rune",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 461, col: 11, offset: 11064},
								expr: &ruleRefExpr{
									pos:  position{line: 461, col: 12, offset: 11065},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 462, col: 4, offset: 11083},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 462, col: 4, offset: 11083},
								val:        "string",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 462, col: 13, offset: 11092},
								expr: &ruleRefExpr{
									pos:  position{line: 462, col: 14, offset: 11093},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 463, col: 4, offset: 11111},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 463, col: 4, offset: 11111},
								val:        "uint8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 463, col: 12, offset: 11119},
								expr: &ruleRefExpr{
									pos:  position{line: 463, col: 13, offset: 11120},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 464, col: 4, offset: 11138},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 464, col: 4, offset: 11138},
								val:        "uint16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 464, col: 13, offset: 11147},
								expr: &ruleRefExpr{
									pos:  position{line: 464, col: 14, offset: 11148},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 465, col: 4, offset: 11166},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 465, col: 4, offset: 11166},
								val:        "uint32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 465, col: 13, offset: 11175},
								expr: &ruleRefExpr{
									pos:  position{line: 465, col: 14, offset: 11176},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 466, col: 4, offset: 11194},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 466, col: 4, offset: 11194},
								val:        "uint64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 466, col: 13, offset: 11203},
								expr: &ruleRefExpr{
									pos:  position{line: 466, col: 14, offset: 11204},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 467, col: 4, offset: 11222},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 467, col: 4, offset: 11222},
								val:        "uintptr",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 467, col: 14, offset: 11232},
								expr: &ruleRefExpr{
									pos:  position{line: 467, col: 15, offset: 11233},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 468, col: 4, offset: 11251},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 468, col: 4, offset: 11251},
								val:        "uint",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 468, col: 11, offset: 11258},
								expr: &ruleRefExpr{
									pos:  position{line: 468, col: 12, offset: 11259},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 469, col: 4, offset: 11277},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 469, col: 4, offset: 11277},
								val:        "true",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 469, col: 11, offset: 11284},
								expr: &ruleRefExpr{
									pos:  position{line: 469, col: 12, offset: 11285},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 470, col: 4, offset: 11303},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 470, col: 4, offset: 11303},
								val:        "false",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 470, col: 12, offset: 11311},
								expr: &ruleRefExpr{
									pos:  position{line: 470, col: 13, offset: 11312},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 471, col: 4, offset: 11330},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 471, col: 4, offset: 11330},
								val:        "iota",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 471, col: 11, offset: 11337},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 12, offset: 11338},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 472, col: 4, offset: 11356},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 472, col: 4, offset: 11356},
								val:        "nil",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 472, col: 10, offset: 11362},
								expr: &ruleRefExpr{
									pos:  position{line: 472, col: 11, offset: 11363},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 473, col: 4, offset: 11381},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 473, col: 4, offset: 11381},
								val:        "append",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 473, col: 13, offset: 11390},
								expr: &ruleRefExpr{
									pos:  position{line: 473, col: 14, offset: 11391},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 474, col: 4, offset: 11409},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 474, col: 4, offset: 11409},
								val:        "cap",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 474, col: 10, offset: 11415},
								expr: &ruleRefExpr{
									pos:  position{line: 474, col: 11, offset: 11416},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 475, col: 4, offset: 11434},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 475, col: 4, offset: 11434},
								val:        "close",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 475, col: 12, offset: 11442},
								expr: &ruleRefExpr{
									pos:  position{line: 475, col: 13, offset: 11443},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 476, col: 4, offset: 11461},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 476, col: 4, offset: 11461},
								val:        "complex",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 476, col: 14, offset: 11471},
								expr: &ruleRefExpr{
									pos:  position{line: 476, col: 15, offset: 11472},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 477, col: 4, offset: 11490},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 477, col: 4, offset: 11490},
								val:        "copy",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 477, col: 11, offset: 11497},
								expr: &ruleRefExpr{
									pos:  position{line: 477, col: 12, offset: 11498},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 478, col: 4, offset: 11516},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 478, col: 4, offset: 11516},
								val:        "delete",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 478, col: 13, offset: 11525},
								expr: &ruleRefExpr{
									pos:  position{line: 478, col: 14, offset: 11526},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 479, col: 4, offset: 11544},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 479, col: 4, offset: 11544},
								val:        "imag",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 479, col: 11, offset: 11551},
								expr: &ruleRefExpr{
									pos:  position{line: 479, col: 12, offset: 11552},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 480, col: 4, offset: 11570},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 480, col: 4, offset: 11570},
								val:        "len",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 480, col: 10, offset: 11576},
								expr: &ruleRefExpr{
									pos:  position{line: 480, col: 11, offset: 11577},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 481, col: 4, offset: 11595},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 481, col: 4, offset: 11595},
								val:        "make",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 481, col: 11, offset: 11602},
								expr: &ruleRefExpr{
									pos:  position{line: 481, col: 12, offset: 11603},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 482, col: 4, offset: 11621},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 482, col: 4, offset: 11621},
								val:        "new",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 482, col: 10, offset: 11627},
								expr: &ruleRefExpr{
									pos:  position{line: 482, col: 11, offset: 11628},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 483, col: 4, offset: 11646},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 483, col: 4, offset: 11646},
								val:        "panic",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 483, col: 12, offset: 11654},
								expr: &ruleRefExpr{
									pos:  position{line: 483, col: 13, offset: 11655},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 484, col: 4, offset: 11673},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 484, col: 4, offset: 11673},
								val:        "println",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 484, col: 14, offset: 11683},
								expr: &ruleRefExpr{
									pos:  position{line: 484, col: 15, offset: 11684},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 485, col: 4, offset: 11702},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 485, col: 4, offset: 11702},
								val:        "print",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 485, col: 12, offset: 11710},
								expr: &ruleRefExpr{
									pos:  position{line: 485, col: 13, offset: 11711},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 486, col: 4, offset: 11729},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 486, col: 4, offset: 11729},
								val:        "real",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 486, col: 11, offset: 11736},
								expr: &ruleRefExpr{
									pos:  position{line: 486, col: 12, offset: 11737},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 487, col: 4, offset: 11755},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 487, col: 4, offset: 11755},
								val:        "recover",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 487, col: 14, offset: 11765},
								expr: &ruleRefExpr{
									pos:  position{line: 487, col: 15, offset: 11766},
									name: "IdentifierPart",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 489, col: 1, offset: 11782},
			expr: &actionExpr{
				pos: position{line: 489, col: 14, offset: 11797},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 489, col: 14, offset: 11797},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 494, col: 1, offset: 11872},
			expr: &actionExpr{
				pos: position{line: 494, col: 13, offset: 11886},
				run: (*parser).callonCodeBlock1,
				expr: &seqExpr{
					pos: position{line: 494, col: 13, offset: 11886},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 494, col: 13, offset: 11886},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 494, col: 17, offset: 11890},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 494, col: 22, offset: 11895},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 500, col: 1, offset: 11993},
			expr: &zeroOrMoreExpr{
				pos: position{line: 500, col: 8, offset: 12002},
				expr: &choiceExpr{
					pos: position{line: 500, col: 10, offset: 12004},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 500, col: 10, offset: 12004},
							expr: &seqExpr{
								pos: position{line: 500, col: 12, offset: 12006},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 500, col: 12, offset: 12006},
										expr: &charClassMatcher{
											pos:        position{line: 500, col: 13, offset: 12007},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 500, col: 18, offset: 12012},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 500, col: 34, offset: 12028},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 500, col: 34, offset: 12028},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 500, col: 38, offset: 12032},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 500, col: 43, offset: 12037},
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
			name: "__",
			pos:  position{line: 502, col: 1, offset: 12045},
			expr: &zeroOrMoreExpr{
				pos: position{line: 502, col: 6, offset: 12052},
				expr: &choiceExpr{
					pos: position{line: 502, col: 8, offset: 12054},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 502, col: 8, offset: 12054},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 502, col: 21, offset: 12067},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 502, col: 27, offset: 12073},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 503, col: 1, offset: 12084},
			expr: &zeroOrMoreExpr{
				pos: position{line: 503, col: 5, offset: 12090},
				expr: &choiceExpr{
					pos: position{line: 503, col: 7, offset: 12092},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 503, col: 7, offset: 12092},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 503, col: 20, offset: 12105},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 505, col: 1, offset: 12142},
			expr: &charClassMatcher{
				pos:        position{line: 505, col: 14, offset: 12157},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 506, col: 1, offset: 12165},
			expr: &litMatcher{
				pos:        position{line: 506, col: 7, offset: 12173},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 507, col: 1, offset: 12178},
			expr: &choiceExpr{
				pos: position{line: 507, col: 7, offset: 12186},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 507, col: 7, offset: 12186},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 507, col: 7, offset: 12186},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 507, col: 10, offset: 12189},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 507, col: 16, offset: 12195},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 507, col: 16, offset: 12195},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 507, col: 18, offset: 12197},
								expr: &ruleRefExpr{
									pos:  position{line: 507, col: 18, offset: 12197},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 507, col: 37, offset: 12216},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 507, col: 43, offset: 12222},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 507, col: 43, offset: 12222},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 507, col: 46, offset: 12225},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 509, col: 1, offset: 12230},
			expr: &notExpr{
				pos: position{line: 509, col: 7, offset: 12238},
				expr: &anyMatcher{
					line: 509, col: 8, offset: 12239,
				},
			},
		},
	},
}

func (c *current) onGrammar1(initializer, rules interface{}) (interface{}, error) {
	pos := c.astPos()

	// create the grammar, assign its initializer
	g := ast.NewGrammar(pos)
	initSlice := toIfaceSlice(initializer)
	if len(initSlice) > 0 {
		g.Init = initSlice[0].(*ast.CodeBlock)
	}

	rulesSlice := toIfaceSlice(rules)
	g.Rules = make([]*ast.Rule, len(rulesSlice))
	for i, duo := range rulesSlice {
		g.Rules[i] = duo.([]interface{})[0].(*ast.Rule)
	}

	return g, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["initializer"], stack["rules"])
}

func (c *current) onInitializer1(code interface{}) (interface{}, error) {
	return code, nil
}

func (p *parser) callonInitializer1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInitializer1(stack["code"])
}

func (c *current) onRule1(name, display, expr interface{}) (interface{}, error) {
	pos := c.astPos()

	rule := ast.NewRule(pos, name.(*ast.Identifier))
	displaySlice := toIfaceSlice(display)
	if len(displaySlice) > 0 {
		rule.DisplayName = displaySlice[0].(*ast.StringLit)
	}
	rule.Expr = expr.(ast.Expression)

	return rule, nil
}

func (p *parser) callonRule1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRule1(stack["name"], stack["display"], stack["expr"])
}

func (c *current) onChoiceExpr1(first, rest interface{}) (interface{}, error) {
	restSlice := toIfaceSlice(rest)
	if len(restSlice) == 0 {
		return first, nil
	}

	pos := c.astPos()
	choice := ast.NewChoiceExpr(pos)
	choice.Alternatives = []ast.Expression{first.(ast.Expression)}
	for _, sl := range restSlice {
		choice.Alternatives = append(choice.Alternatives, sl.([]interface{})[3].(ast.Expression))
	}
	return choice, nil
}

func (p *parser) callonChoiceExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onChoiceExpr1(stack["first"], stack["rest"])
}

func (c *current) onActionExpr1(expr, code interface{}) (interface{}, error) {
	if code == nil {
		return expr, nil
	}

	pos := c.astPos()
	act := ast.NewActionExpr(pos)
	act.Expr = expr.(ast.Expression)
	codeSlice := toIfaceSlice(code)
	act.Code = codeSlice[1].(*ast.CodeBlock)

	return act, nil
}

func (p *parser) callonActionExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onActionExpr1(stack["expr"], stack["code"])
}

func (c *current) onSeqExpr1(first, rest interface{}) (interface{}, error) {
	restSlice := toIfaceSlice(rest)
	if len(restSlice) == 0 {
		return first, nil
	}
	seq := ast.NewSeqExpr(c.astPos())
	seq.Exprs = []ast.Expression{first.(ast.Expression)}
	for _, sl := range restSlice {
		seq.Exprs = append(seq.Exprs, sl.([]interface{})[1].(ast.Expression))
	}
	return seq, nil
}

func (p *parser) callonSeqExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSeqExpr1(stack["first"], stack["rest"])
}

func (c *current) onLabeledExpr2(label, expr interface{}) (interface{}, error) {
	pos := c.astPos()
	lab := ast.NewLabeledExpr(pos)
	lab.Label = label.(*ast.Identifier)
	lab.Expr = expr.(ast.Expression)
	return lab, nil
}

func (p *parser) callonLabeledExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabeledExpr2(stack["label"], stack["expr"])
}

func (c *current) onPrefixedExpr2(op, expr interface{}) (interface{}, error) {
	pos := c.astPos()
	opStr := op.(string)
	if opStr == "&" {
		and := ast.NewAndExpr(pos)
		and.Expr = expr.(ast.Expression)
		return and, nil
	}
	not := ast.NewNotExpr(pos)
	not.Expr = expr.(ast.Expression)
	return not, nil
}

func (p *parser) callonPrefixedExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedExpr2(stack["op"], stack["expr"])
}

func (c *current) onPrefixedOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonPrefixedOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedOp1()
}

func (c *current) onSuffixedExpr2(expr, op interface{}) (interface{}, error) {
	pos := c.astPos()
	opStr := op.(string)
	switch opStr {
	case "?":
		zero := ast.NewZeroOrOneExpr(pos)
		zero.Expr = expr.(ast.Expression)
		return zero, nil
	case "*":
		zero := ast.NewZeroOrMoreExpr(pos)
		zero.Expr = expr.(ast.Expression)
		return zero, nil
	case "+":
		one := ast.NewOneOrMoreExpr(pos)
		one.Expr = expr.(ast.Expression)
		return one, nil
	default:
		return nil, errors.New("unknown operator: " + opStr)
	}
}

func (p *parser) callonSuffixedExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedExpr2(stack["expr"], stack["op"])
}

func (c *current) onSuffixedOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonSuffixedOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedOp1()
}

func (c *current) onPrimaryExpr7(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonPrimaryExpr7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimaryExpr7(stack["expr"])
}

func (c *current) onRuleRefExpr1(name interface{}) (interface{}, error) {
	ref := ast.NewRuleRefExpr(c.astPos())
	ref.Name = name.(*ast.Identifier)
	return ref, nil
}

func (p *parser) callonRuleRefExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRuleRefExpr1(stack["name"])
}

func (c *current) onSemanticPredExpr1(op, code interface{}) (interface{}, error) {
	opStr := op.(string)
	if opStr == "&" {
		and := ast.NewAndCodeExpr(c.astPos())
		and.Code = code.(*ast.CodeBlock)
		return and, nil
	}
	not := ast.NewNotCodeExpr(c.astPos())
	not.Code = code.(*ast.CodeBlock)
	return not, nil
}

func (p *parser) callonSemanticPredExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredExpr1(stack["op"], stack["code"])
}

func (c *current) onSemanticPredOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonSemanticPredOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredOp1()
}

func (c *current) onIdentifier2(ident interface{}) (interface{}, error) {
	return ident, nil
}

func (p *parser) callonIdentifier2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier2(stack["ident"])
}

func (c *current) onIdentifier8(ident interface{}) (interface{}, error) {
	return ast.NewIdentifier(c.astPos(), string(c.text)), errors.New("identifier is a reserved word")
}

func (p *parser) callonIdentifier8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier8(stack["ident"])
}

func (c *current) onIdentifierName1() (interface{}, error) {
	return ast.NewIdentifier(c.astPos(), string(c.text)), nil
}

func (p *parser) callonIdentifierName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierName1()
}

func (c *current) onLitMatcher1(lit, ignore interface{}) (interface{}, error) {
	rawStr := lit.(*ast.StringLit).Val
	s, err := strconv.Unquote(rawStr)
	if err != nil {
		return nil, err
	}
	m := ast.NewLitMatcher(c.astPos(), s)
	m.IgnoreCase = ignore != nil
	return m, nil
}

func (p *parser) callonLitMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLitMatcher1(stack["lit"], stack["ignore"])
}

func (c *current) onStringLiteral1() (interface{}, error) {
	return ast.NewStringLit(c.astPos(), string(c.text)), nil
}

func (p *parser) callonStringLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral1()
}

func (c *current) onCharClassMatcher1() (interface{}, error) {
	pos := c.astPos()
	cc := ast.NewCharClassMatcher(pos, string(c.text))
	return cc, nil
}

func (p *parser) callonCharClassMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher1()
}

func (c *current) onAnyMatcher1() (interface{}, error) {
	any := ast.NewAnyMatcher(c.astPos(), ".")
	return any, nil
}

func (p *parser) callonAnyMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnyMatcher1()
}

func (c *current) onCodeBlock1() (interface{}, error) {
	pos := c.astPos()
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonCodeBlock1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock1()
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

// ParserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type ParserError struct {
	Inner  error
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
