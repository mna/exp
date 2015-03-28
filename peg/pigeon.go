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
			pos:  position{line: 167, col: 1, offset: 4468},
			expr: &ruleRefExpr{
				pos:  position{line: 167, col: 14, offset: 4483},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 168, col: 1, offset: 4498},
			expr: &actionExpr{
				pos: position{line: 168, col: 18, offset: 4517},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 168, col: 18, offset: 4517},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 168, col: 18, offset: 4517},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 168, col: 34, offset: 4533},
							expr: &ruleRefExpr{
								pos:  position{line: 168, col: 34, offset: 4533},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 171, col: 1, offset: 4615},
			expr: &charClassMatcher{
				pos:        position{line: 171, col: 19, offset: 4635},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 172, col: 1, offset: 4643},
			expr: &choiceExpr{
				pos: position{line: 172, col: 18, offset: 4662},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 172, col: 18, offset: 4662},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 172, col: 36, offset: 4680},
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
			pos:  position{line: 174, col: 1, offset: 4687},
			expr: &actionExpr{
				pos: position{line: 174, col: 14, offset: 4702},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 174, col: 14, offset: 4702},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 174, col: 14, offset: 4702},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 174, col: 18, offset: 4706},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 174, col: 32, offset: 4720},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 174, col: 39, offset: 4727},
								expr: &litMatcher{
									pos:        position{line: 174, col: 39, offset: 4727},
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
			pos:  position{line: 184, col: 1, offset: 4953},
			expr: &actionExpr{
				pos: position{line: 184, col: 17, offset: 4971},
				run: (*parser).callonStringLiteral1,
				expr: &choiceExpr{
					pos: position{line: 184, col: 19, offset: 4973},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 184, col: 19, offset: 4973},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 184, col: 19, offset: 4973},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 184, col: 23, offset: 4977},
									expr: &ruleRefExpr{
										pos:  position{line: 184, col: 23, offset: 4977},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 184, col: 41, offset: 4995},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 184, col: 47, offset: 5001},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 184, col: 47, offset: 5001},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 184, col: 51, offset: 5005},
									name: "SingleStringChar",
								},
								&litMatcher{
									pos:        position{line: 184, col: 68, offset: 5022},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 184, col: 74, offset: 5028},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 184, col: 74, offset: 5028},
									val:        "`",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 184, col: 78, offset: 5032},
									name: "RawStringChar",
								},
								&litMatcher{
									pos:        position{line: 184, col: 92, offset: 5046},
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
			pos:  position{line: 187, col: 1, offset: 5117},
			expr: &choiceExpr{
				pos: position{line: 187, col: 20, offset: 5138},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 187, col: 20, offset: 5138},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 187, col: 20, offset: 5138},
								expr: &choiceExpr{
									pos: position{line: 187, col: 23, offset: 5141},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 187, col: 23, offset: 5141},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 187, col: 29, offset: 5147},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 187, col: 36, offset: 5154},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 187, col: 42, offset: 5160},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 187, col: 55, offset: 5173},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 187, col: 55, offset: 5173},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 187, col: 60, offset: 5178},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 188, col: 1, offset: 5197},
			expr: &choiceExpr{
				pos: position{line: 188, col: 20, offset: 5218},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 188, col: 20, offset: 5218},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 188, col: 20, offset: 5218},
								expr: &choiceExpr{
									pos: position{line: 188, col: 23, offset: 5221},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 188, col: 23, offset: 5221},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 188, col: 29, offset: 5227},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 36, offset: 5234},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 188, col: 42, offset: 5240},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 188, col: 55, offset: 5253},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 188, col: 55, offset: 5253},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 188, col: 60, offset: 5258},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 189, col: 1, offset: 5277},
			expr: &seqExpr{
				pos: position{line: 189, col: 17, offset: 5295},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 189, col: 17, offset: 5295},
						expr: &litMatcher{
							pos:        position{line: 189, col: 18, offset: 5296},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 189, col: 22, offset: 5300},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 191, col: 1, offset: 5312},
			expr: &choiceExpr{
				pos: position{line: 191, col: 22, offset: 5335},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 191, col: 22, offset: 5335},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 191, col: 28, offset: 5341},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 192, col: 1, offset: 5362},
			expr: &choiceExpr{
				pos: position{line: 192, col: 22, offset: 5385},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 192, col: 22, offset: 5385},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 192, col: 28, offset: 5391},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 194, col: 1, offset: 5413},
			expr: &choiceExpr{
				pos: position{line: 194, col: 24, offset: 5438},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 194, col: 24, offset: 5438},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 194, col: 43, offset: 5457},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 194, col: 57, offset: 5471},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 194, col: 69, offset: 5483},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 194, col: 89, offset: 5503},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 195, col: 1, offset: 5522},
			expr: &choiceExpr{
				pos: position{line: 195, col: 20, offset: 5543},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 195, col: 20, offset: 5543},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 26, offset: 5549},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 32, offset: 5555},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 38, offset: 5561},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 44, offset: 5567},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 50, offset: 5573},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 56, offset: 5579},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 62, offset: 5585},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 196, col: 1, offset: 5590},
			expr: &seqExpr{
				pos: position{line: 196, col: 15, offset: 5606},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 196, col: 15, offset: 5606},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 196, col: 26, offset: 5617},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 196, col: 37, offset: 5628},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 197, col: 1, offset: 5639},
			expr: &seqExpr{
				pos: position{line: 197, col: 13, offset: 5653},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 197, col: 13, offset: 5653},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 17, offset: 5657},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 26, offset: 5666},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 198, col: 1, offset: 5675},
			expr: &seqExpr{
				pos: position{line: 198, col: 21, offset: 5697},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 198, col: 21, offset: 5697},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 25, offset: 5701},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 34, offset: 5710},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 43, offset: 5719},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 52, offset: 5728},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 61, offset: 5737},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 70, offset: 5746},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 79, offset: 5755},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 88, offset: 5764},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 199, col: 1, offset: 5773},
			expr: &seqExpr{
				pos: position{line: 199, col: 22, offset: 5796},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 199, col: 22, offset: 5796},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 26, offset: 5800},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 35, offset: 5809},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 44, offset: 5818},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 53, offset: 5827},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 201, col: 1, offset: 5837},
			expr: &charClassMatcher{
				pos:        position{line: 201, col: 14, offset: 5852},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 202, col: 1, offset: 5858},
			expr: &charClassMatcher{
				pos:        position{line: 202, col: 16, offset: 5875},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 203, col: 1, offset: 5881},
			expr: &charClassMatcher{
				pos:        position{line: 203, col: 12, offset: 5894},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 205, col: 1, offset: 5905},
			expr: &actionExpr{
				pos: position{line: 205, col: 20, offset: 5926},
				run: (*parser).callonCharClassMatcher1,
				expr: &seqExpr{
					pos: position{line: 205, col: 20, offset: 5926},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 205, col: 20, offset: 5926},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 205, col: 24, offset: 5930},
							expr: &choiceExpr{
								pos: position{line: 205, col: 26, offset: 5932},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 205, col: 26, offset: 5932},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 205, col: 43, offset: 5949},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 205, col: 55, offset: 5961},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 205, col: 55, offset: 5961},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 205, col: 60, offset: 5966},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 205, col: 82, offset: 5988},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 205, col: 86, offset: 5992},
							expr: &litMatcher{
								pos:        position{line: 205, col: 86, offset: 5992},
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
			pos:  position{line: 210, col: 1, offset: 6097},
			expr: &seqExpr{
				pos: position{line: 210, col: 18, offset: 6116},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 210, col: 18, offset: 6116},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 210, col: 28, offset: 6126},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 210, col: 32, offset: 6130},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 211, col: 1, offset: 6140},
			expr: &choiceExpr{
				pos: position{line: 211, col: 13, offset: 6154},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 211, col: 13, offset: 6154},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 211, col: 13, offset: 6154},
								expr: &choiceExpr{
									pos: position{line: 211, col: 16, offset: 6157},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 211, col: 16, offset: 6157},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 211, col: 22, offset: 6163},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 211, col: 29, offset: 6170},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 35, offset: 6176},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 211, col: 48, offset: 6189},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 211, col: 48, offset: 6189},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 53, offset: 6194},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 212, col: 1, offset: 6210},
			expr: &choiceExpr{
				pos: position{line: 212, col: 19, offset: 6230},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 212, col: 19, offset: 6230},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 25, offset: 6236},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 214, col: 1, offset: 6258},
			expr: &seqExpr{
				pos: position{line: 214, col: 22, offset: 6281},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 214, col: 22, offset: 6281},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 214, col: 28, offset: 6287},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 214, col: 28, offset: 6287},
								name: "SingleCharUnicodeClass",
							},
							&seqExpr{
								pos: position{line: 214, col: 53, offset: 6312},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 214, col: 53, offset: 6312},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 214, col: 57, offset: 6316},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 214, col: 70, offset: 6329},
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
			pos:  position{line: 215, col: 1, offset: 6335},
			expr: &charClassMatcher{
				pos:        position{line: 215, col: 26, offset: 6362},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 216, col: 1, offset: 6372},
			expr: &oneOrMoreExpr{
				pos: position{line: 216, col: 16, offset: 6389},
				expr: &charClassMatcher{
					pos:        position{line: 216, col: 16, offset: 6389},
					val:        "[a-z_]i",
					chars:      []rune{'_'},
					ranges:     []rune{'a', 'z'},
					ignoreCase: true,
					inverted:   false,
				},
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 218, col: 1, offset: 6433},
			expr: &actionExpr{
				pos: position{line: 218, col: 14, offset: 6448},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 218, col: 14, offset: 6448},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 223, col: 1, offset: 6523},
			expr: &actionExpr{
				pos: position{line: 223, col: 13, offset: 6537},
				run: (*parser).callonCodeBlock1,
				expr: &seqExpr{
					pos: position{line: 223, col: 13, offset: 6537},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 223, col: 13, offset: 6537},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 17, offset: 6541},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 223, col: 22, offset: 6546},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 229, col: 1, offset: 6644},
			expr: &zeroOrMoreExpr{
				pos: position{line: 229, col: 8, offset: 6653},
				expr: &choiceExpr{
					pos: position{line: 229, col: 10, offset: 6655},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 229, col: 10, offset: 6655},
							expr: &seqExpr{
								pos: position{line: 229, col: 12, offset: 6657},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 229, col: 12, offset: 6657},
										expr: &charClassMatcher{
											pos:        position{line: 229, col: 13, offset: 6658},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 229, col: 18, offset: 6663},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 229, col: 34, offset: 6679},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 229, col: 34, offset: 6679},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 229, col: 38, offset: 6683},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 229, col: 43, offset: 6688},
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
			pos:  position{line: 231, col: 1, offset: 6696},
			expr: &zeroOrMoreExpr{
				pos: position{line: 231, col: 6, offset: 6703},
				expr: &choiceExpr{
					pos: position{line: 231, col: 8, offset: 6705},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 231, col: 8, offset: 6705},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 231, col: 21, offset: 6718},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 231, col: 27, offset: 6724},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 232, col: 1, offset: 6735},
			expr: &zeroOrMoreExpr{
				pos: position{line: 232, col: 5, offset: 6741},
				expr: &choiceExpr{
					pos: position{line: 232, col: 7, offset: 6743},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 232, col: 7, offset: 6743},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 20, offset: 6756},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 234, col: 1, offset: 6793},
			expr: &charClassMatcher{
				pos:        position{line: 234, col: 14, offset: 6808},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 235, col: 1, offset: 6816},
			expr: &litMatcher{
				pos:        position{line: 235, col: 7, offset: 6824},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 236, col: 1, offset: 6829},
			expr: &choiceExpr{
				pos: position{line: 236, col: 7, offset: 6837},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 236, col: 7, offset: 6837},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 236, col: 7, offset: 6837},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 236, col: 10, offset: 6840},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 236, col: 16, offset: 6846},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 236, col: 16, offset: 6846},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 236, col: 18, offset: 6848},
								expr: &ruleRefExpr{
									pos:  position{line: 236, col: 18, offset: 6848},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 236, col: 37, offset: 6867},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 236, col: 43, offset: 6873},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 236, col: 43, offset: 6873},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 236, col: 46, offset: 6876},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 238, col: 1, offset: 6881},
			expr: &notExpr{
				pos: position{line: 238, col: 7, offset: 6889},
				expr: &anyMatcher{
					line: 238, col: 8, offset: 6890,
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
	if _, ok := err.(*ParserError); ok {
		p.errs.add(err)
		return
	}
	if len(p.rstack) == 0 {
		p.errs.add(err)
		return
	}
	rule := p.rstack[len(p.rstack)-1]
	pe := &ParserError{Inner: err, prefix: rule.name}
	if rule.displayName != "" {
		pe.prefix = rule.displayName
	}
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
		return nil, ErrNoRule
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
			p.addErr(err)
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

	// TODO : should be caught at the scan/parse step
	return &unicode.RangeTable{} // empty range
}
