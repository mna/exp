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
			pos:  position{line: 171, col: 1, offset: 4595},
			expr: &actionExpr{
				pos: position{line: 171, col: 18, offset: 4614},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 171, col: 18, offset: 4614},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 171, col: 18, offset: 4614},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 171, col: 34, offset: 4630},
							expr: &ruleRefExpr{
								pos:  position{line: 171, col: 34, offset: 4630},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 174, col: 1, offset: 4712},
			expr: &charClassMatcher{
				pos:        position{line: 174, col: 19, offset: 4732},
				val:        "[\\pL_]",
				chars:      []rune{'_'},
				classes:    []*unicode.RangeTable{rangeTable("L")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 175, col: 1, offset: 4739},
			expr: &choiceExpr{
				pos: position{line: 175, col: 18, offset: 4758},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 175, col: 18, offset: 4758},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 175, col: 36, offset: 4776},
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
			pos:  position{line: 177, col: 1, offset: 4786},
			expr: &actionExpr{
				pos: position{line: 177, col: 14, offset: 4801},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 177, col: 14, offset: 4801},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 177, col: 14, offset: 4801},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 177, col: 18, offset: 4805},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 177, col: 32, offset: 4819},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 177, col: 39, offset: 4826},
								expr: &litMatcher{
									pos:        position{line: 177, col: 39, offset: 4826},
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
			pos:  position{line: 187, col: 1, offset: 5052},
			expr: &actionExpr{
				pos: position{line: 187, col: 17, offset: 5070},
				run: (*parser).callonStringLiteral1,
				expr: &choiceExpr{
					pos: position{line: 187, col: 19, offset: 5072},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 187, col: 19, offset: 5072},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 19, offset: 5072},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 187, col: 23, offset: 5076},
									expr: &ruleRefExpr{
										pos:  position{line: 187, col: 23, offset: 5076},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 187, col: 41, offset: 5094},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 187, col: 47, offset: 5100},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 47, offset: 5100},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 187, col: 51, offset: 5104},
									name: "SingleStringChar",
								},
								&litMatcher{
									pos:        position{line: 187, col: 68, offset: 5121},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 187, col: 74, offset: 5127},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 187, col: 74, offset: 5127},
									val:        "`",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 187, col: 78, offset: 5131},
									name: "RawStringChar",
								},
								&litMatcher{
									pos:        position{line: 187, col: 92, offset: 5145},
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
			pos:  position{line: 190, col: 1, offset: 5216},
			expr: &choiceExpr{
				pos: position{line: 190, col: 20, offset: 5237},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 190, col: 20, offset: 5237},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 190, col: 20, offset: 5237},
								expr: &choiceExpr{
									pos: position{line: 190, col: 23, offset: 5240},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 190, col: 23, offset: 5240},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 190, col: 29, offset: 5246},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 190, col: 36, offset: 5253},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 42, offset: 5259},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 190, col: 55, offset: 5272},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 190, col: 55, offset: 5272},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 190, col: 60, offset: 5277},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 191, col: 1, offset: 5296},
			expr: &choiceExpr{
				pos: position{line: 191, col: 20, offset: 5317},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 191, col: 20, offset: 5317},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 191, col: 20, offset: 5317},
								expr: &choiceExpr{
									pos: position{line: 191, col: 23, offset: 5320},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 191, col: 23, offset: 5320},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 191, col: 29, offset: 5326},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 191, col: 36, offset: 5333},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 191, col: 42, offset: 5339},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 191, col: 55, offset: 5352},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 191, col: 55, offset: 5352},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 191, col: 60, offset: 5357},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 192, col: 1, offset: 5376},
			expr: &seqExpr{
				pos: position{line: 192, col: 17, offset: 5394},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 192, col: 17, offset: 5394},
						expr: &litMatcher{
							pos:        position{line: 192, col: 18, offset: 5395},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 192, col: 22, offset: 5399},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 194, col: 1, offset: 5411},
			expr: &choiceExpr{
				pos: position{line: 194, col: 22, offset: 5434},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 194, col: 22, offset: 5434},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 194, col: 28, offset: 5440},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 195, col: 1, offset: 5461},
			expr: &choiceExpr{
				pos: position{line: 195, col: 22, offset: 5484},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 195, col: 22, offset: 5484},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 195, col: 28, offset: 5490},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 197, col: 1, offset: 5512},
			expr: &choiceExpr{
				pos: position{line: 197, col: 24, offset: 5537},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 197, col: 24, offset: 5537},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 43, offset: 5556},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 57, offset: 5570},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 69, offset: 5582},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 89, offset: 5602},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 198, col: 1, offset: 5621},
			expr: &choiceExpr{
				pos: position{line: 198, col: 20, offset: 5642},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 198, col: 20, offset: 5642},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 26, offset: 5648},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 32, offset: 5654},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 38, offset: 5660},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 44, offset: 5666},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 50, offset: 5672},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 56, offset: 5678},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 62, offset: 5684},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 199, col: 1, offset: 5689},
			expr: &seqExpr{
				pos: position{line: 199, col: 15, offset: 5705},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 199, col: 15, offset: 5705},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 26, offset: 5716},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 199, col: 37, offset: 5727},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 200, col: 1, offset: 5738},
			expr: &seqExpr{
				pos: position{line: 200, col: 13, offset: 5752},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 200, col: 13, offset: 5752},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 17, offset: 5756},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 26, offset: 5765},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 201, col: 1, offset: 5774},
			expr: &seqExpr{
				pos: position{line: 201, col: 21, offset: 5796},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 201, col: 21, offset: 5796},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 25, offset: 5800},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 34, offset: 5809},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 43, offset: 5818},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 52, offset: 5827},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 61, offset: 5836},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 70, offset: 5845},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 79, offset: 5854},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 88, offset: 5863},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 202, col: 1, offset: 5872},
			expr: &seqExpr{
				pos: position{line: 202, col: 22, offset: 5895},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 202, col: 22, offset: 5895},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 26, offset: 5899},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 35, offset: 5908},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 44, offset: 5917},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 53, offset: 5926},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 204, col: 1, offset: 5936},
			expr: &charClassMatcher{
				pos:        position{line: 204, col: 14, offset: 5951},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 205, col: 1, offset: 5957},
			expr: &charClassMatcher{
				pos:        position{line: 205, col: 16, offset: 5974},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 206, col: 1, offset: 5980},
			expr: &charClassMatcher{
				pos:        position{line: 206, col: 12, offset: 5993},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 208, col: 1, offset: 6004},
			expr: &actionExpr{
				pos: position{line: 208, col: 20, offset: 6025},
				run: (*parser).callonCharClassMatcher1,
				expr: &seqExpr{
					pos: position{line: 208, col: 20, offset: 6025},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 208, col: 20, offset: 6025},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 208, col: 24, offset: 6029},
							expr: &choiceExpr{
								pos: position{line: 208, col: 26, offset: 6031},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 208, col: 26, offset: 6031},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 208, col: 43, offset: 6048},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 208, col: 55, offset: 6060},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 208, col: 55, offset: 6060},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 208, col: 60, offset: 6065},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 208, col: 82, offset: 6087},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 208, col: 86, offset: 6091},
							expr: &litMatcher{
								pos:        position{line: 208, col: 86, offset: 6091},
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
			pos:  position{line: 213, col: 1, offset: 6196},
			expr: &seqExpr{
				pos: position{line: 213, col: 18, offset: 6215},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 213, col: 18, offset: 6215},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 213, col: 28, offset: 6225},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 213, col: 32, offset: 6229},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 214, col: 1, offset: 6239},
			expr: &choiceExpr{
				pos: position{line: 214, col: 13, offset: 6253},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 214, col: 13, offset: 6253},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 214, col: 13, offset: 6253},
								expr: &choiceExpr{
									pos: position{line: 214, col: 16, offset: 6256},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 214, col: 16, offset: 6256},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 214, col: 22, offset: 6262},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 214, col: 29, offset: 6269},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 35, offset: 6275},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 214, col: 48, offset: 6288},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 214, col: 48, offset: 6288},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 53, offset: 6293},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 215, col: 1, offset: 6309},
			expr: &choiceExpr{
				pos: position{line: 215, col: 19, offset: 6329},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 215, col: 19, offset: 6329},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 215, col: 25, offset: 6335},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 217, col: 1, offset: 6357},
			expr: &seqExpr{
				pos: position{line: 217, col: 22, offset: 6380},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 217, col: 22, offset: 6380},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 217, col: 28, offset: 6386},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 217, col: 28, offset: 6386},
								name: "SingleCharUnicodeClass",
							},
							&seqExpr{
								pos: position{line: 217, col: 53, offset: 6411},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 217, col: 53, offset: 6411},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 217, col: 57, offset: 6415},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 217, col: 70, offset: 6428},
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
			pos:  position{line: 218, col: 1, offset: 6434},
			expr: &charClassMatcher{
				pos:        position{line: 218, col: 26, offset: 6461},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 221, col: 1, offset: 6534},
			expr: &choiceExpr{
				pos: position{line: 221, col: 16, offset: 6551},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 221, col: 16, offset: 6551},
						val:        "Other_Default_Ignorable_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 222, col: 7, offset: 6594},
						val:        "Noncharacter_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 223, col: 7, offset: 6626},
						val:        "Logical_Order_Exception",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 224, col: 7, offset: 6658},
						val:        "Inscriptional_Parthian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 225, col: 7, offset: 6689},
						val:        "Other_Grapheme_Extend",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 226, col: 7, offset: 6719},
						val:        "Inscriptional_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 227, col: 7, offset: 6749},
						val:        "Terminal_Punctuation",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 7, offset: 6778},
						val:        "Meroitic_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 229, col: 7, offset: 6807},
						val:        "IDS_Trinary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 230, col: 7, offset: 6836},
						val:        "Egyptian_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 231, col: 7, offset: 6865},
						val:        "Pattern_White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 232, col: 7, offset: 6893},
						val:        "IDS_Binary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 233, col: 7, offset: 6921},
						val:        "Canadian_Aboriginal",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 234, col: 7, offset: 6949},
						val:        "Variation_Selector",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 7, offset: 6976},
						val:        "Caucasian_Albanian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 236, col: 7, offset: 7003},
						val:        "Unified_Ideograph",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 237, col: 7, offset: 7029},
						val:        "Other_ID_Continue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 238, col: 7, offset: 7055},
						val:        "Old_South_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 239, col: 7, offset: 7081},
						val:        "Old_North_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 240, col: 7, offset: 7107},
						val:        "Other_Alphabetic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 241, col: 7, offset: 7132},
						val:        "Meroitic_Cursive",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 242, col: 7, offset: 7157},
						val:        "Imperial_Aramaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 243, col: 7, offset: 7182},
						val:        "Psalter_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 244, col: 7, offset: 7206},
						val:        "Other_Uppercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 245, col: 7, offset: 7230},
						val:        "Other_Lowercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 246, col: 7, offset: 7254},
						val:        "ASCII_Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 247, col: 7, offset: 7278},
						val:        "Quotation_Mark",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 248, col: 7, offset: 7301},
						val:        "Pattern_Syntax",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 249, col: 7, offset: 7324},
						val:        "Other_ID_Start",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 250, col: 7, offset: 7347},
						val:        "Mende_Kikakui",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 251, col: 7, offset: 7369},
						val:        "Syloti_Nagri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 252, col: 7, offset: 7390},
						val:        "Sora_Sompeng",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 253, col: 7, offset: 7411},
						val:        "Pahawh_Hmong",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 254, col: 7, offset: 7432},
						val:        "Meetei_Mayek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 255, col: 7, offset: 7453},
						val:        "Join_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 7, offset: 7474},
						val:        "Bidi_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 257, col: 7, offset: 7495},
						val:        "White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 258, col: 7, offset: 7515},
						val:        "Warang_Citi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 259, col: 7, offset: 7535},
						val:        "Soft_Dotted",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 260, col: 7, offset: 7555},
						val:        "Pau_Cin_Hau",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 261, col: 7, offset: 7575},
						val:        "Old_Persian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 7, offset: 7595},
						val:        "New_Tai_Lue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 263, col: 7, offset: 7615},
						val:        "Ideographic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 264, col: 7, offset: 7635},
						val:        "Saurashtra",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 265, col: 7, offset: 7654},
						val:        "Phoenician",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 266, col: 7, offset: 7673},
						val:        "Other_Math",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 267, col: 7, offset: 7692},
						val:        "Old_Turkic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 268, col: 7, offset: 7711},
						val:        "Old_Permic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 269, col: 7, offset: 7730},
						val:        "Old_Italic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 270, col: 7, offset: 7749},
						val:        "Manichaean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 271, col: 7, offset: 7768},
						val:        "Kharoshthi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 272, col: 7, offset: 7787},
						val:        "Glagolitic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 273, col: 7, offset: 7806},
						val:        "Devanagari",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 274, col: 7, offset: 7825},
						val:        "Deprecated",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 275, col: 7, offset: 7844},
						val:        "Sundanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 276, col: 7, offset: 7862},
						val:        "Samaritan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 277, col: 7, offset: 7880},
						val:        "Palmyrene",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 278, col: 7, offset: 7898},
						val:        "Nabataean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 7, offset: 7916},
						val:        "Mongolian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 280, col: 7, offset: 7934},
						val:        "Malayalam",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 281, col: 7, offset: 7952},
						val:        "Khudawadi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 282, col: 7, offset: 7970},
						val:        "Inherited",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 283, col: 7, offset: 7988},
						val:        "Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 284, col: 7, offset: 8006},
						val:        "Diacritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 7, offset: 8024},
						val:        "Cuneiform",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 286, col: 7, offset: 8042},
						val:        "Bassa_Vah",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 287, col: 7, offset: 8060},
						val:        "Ugaritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 288, col: 7, offset: 8077},
						val:        "Tifinagh",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 289, col: 7, offset: 8094},
						val:        "Tai_Viet",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 290, col: 7, offset: 8111},
						val:        "Tai_Tham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 291, col: 7, offset: 8128},
						val:        "Tagbanwa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 292, col: 7, offset: 8145},
						val:        "Phags_Pa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 293, col: 7, offset: 8162},
						val:        "Ol_Chiki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 294, col: 7, offset: 8179},
						val:        "Mahajani",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 295, col: 7, offset: 8196},
						val:        "Linear_B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 296, col: 7, offset: 8213},
						val:        "Linear_A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 297, col: 7, offset: 8230},
						val:        "Kayah_Li",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 298, col: 7, offset: 8247},
						val:        "Katakana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 299, col: 7, offset: 8264},
						val:        "Javanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 300, col: 7, offset: 8281},
						val:        "Hiragana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 301, col: 7, offset: 8298},
						val:        "Gurmukhi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 302, col: 7, offset: 8315},
						val:        "Gujarati",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 303, col: 7, offset: 8332},
						val:        "Georgian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 304, col: 7, offset: 8349},
						val:        "Extender",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 305, col: 7, offset: 8366},
						val:        "Ethiopic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 306, col: 7, offset: 8383},
						val:        "Duployan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 307, col: 7, offset: 8400},
						val:        "Cyrillic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 308, col: 7, offset: 8417},
						val:        "Cherokee",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 309, col: 7, offset: 8434},
						val:        "Buginese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 310, col: 7, offset: 8451},
						val:        "Bopomofo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 311, col: 7, offset: 8468},
						val:        "Balinese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 312, col: 7, offset: 8485},
						val:        "Armenian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 313, col: 7, offset: 8502},
						val:        "Tirhuta",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 314, col: 7, offset: 8518},
						val:        "Tibetan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 315, col: 7, offset: 8534},
						val:        "Tagalog",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 316, col: 7, offset: 8550},
						val:        "Sinhala",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 317, col: 7, offset: 8566},
						val:        "Siddham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 318, col: 7, offset: 8582},
						val:        "Shavian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 319, col: 7, offset: 8598},
						val:        "Sharada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 320, col: 7, offset: 8614},
						val:        "Radical",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 321, col: 7, offset: 8630},
						val:        "Osmanya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 322, col: 7, offset: 8646},
						val:        "Myanmar",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 323, col: 7, offset: 8662},
						val:        "Mandaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 324, col: 7, offset: 8678},
						val:        "Kannada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 325, col: 7, offset: 8694},
						val:        "Hanunoo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 326, col: 7, offset: 8710},
						val:        "Grantha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 327, col: 7, offset: 8726},
						val:        "Elbasan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 328, col: 7, offset: 8742},
						val:        "Deseret",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 329, col: 7, offset: 8758},
						val:        "Cypriot",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 330, col: 7, offset: 8774},
						val:        "Braille",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 331, col: 7, offset: 8790},
						val:        "Bengali",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 332, col: 7, offset: 8806},
						val:        "Avestan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 333, col: 7, offset: 8822},
						val:        "Thaana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 334, col: 7, offset: 8837},
						val:        "Telugu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 335, col: 7, offset: 8852},
						val:        "Tai_Le",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 336, col: 7, offset: 8867},
						val:        "Syriac",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 337, col: 7, offset: 8882},
						val:        "Rejang",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 338, col: 7, offset: 8897},
						val:        "Lydian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 339, col: 7, offset: 8912},
						val:        "Lycian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 340, col: 7, offset: 8927},
						val:        "Lepcha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 341, col: 7, offset: 8942},
						val:        "Khojki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 342, col: 7, offset: 8957},
						val:        "Kaithi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 343, col: 7, offset: 8972},
						val:        "Hyphen",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 344, col: 7, offset: 8987},
						val:        "Hebrew",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 345, col: 7, offset: 9002},
						val:        "Hangul",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 346, col: 7, offset: 9017},
						val:        "Gothic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 347, col: 7, offset: 9032},
						val:        "Coptic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 348, col: 7, offset: 9047},
						val:        "Common",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 349, col: 7, offset: 9062},
						val:        "Chakma",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 350, col: 7, offset: 9077},
						val:        "Carian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 351, col: 7, offset: 9092},
						val:        "Brahmi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 352, col: 7, offset: 9107},
						val:        "Arabic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 353, col: 7, offset: 9122},
						val:        "Tamil",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 354, col: 7, offset: 9136},
						val:        "Takri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 355, col: 7, offset: 9150},
						val:        "STerm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 356, col: 7, offset: 9164},
						val:        "Runic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 357, col: 7, offset: 9178},
						val:        "Oriya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 358, col: 7, offset: 9192},
						val:        "Ogham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 359, col: 7, offset: 9206},
						val:        "Limbu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 360, col: 7, offset: 9220},
						val:        "Latin",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 361, col: 7, offset: 9234},
						val:        "Khmer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 362, col: 7, offset: 9248},
						val:        "Greek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 363, col: 7, offset: 9262},
						val:        "Buhid",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 364, col: 7, offset: 9276},
						val:        "Batak",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 365, col: 7, offset: 9290},
						val:        "Bamum",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 366, col: 7, offset: 9304},
						val:        "Thai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 367, col: 7, offset: 9317},
						val:        "Modi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 368, col: 7, offset: 9330},
						val:        "Miao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 369, col: 7, offset: 9343},
						val:        "Lisu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 370, col: 7, offset: 9356},
						val:        "Dash",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 371, col: 7, offset: 9369},
						val:        "Cham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 372, col: 7, offset: 9382},
						val:        "Vai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 373, col: 7, offset: 9394},
						val:        "Nko",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 374, col: 7, offset: 9406},
						val:        "Mro",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 375, col: 7, offset: 9418},
						val:        "Lao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 376, col: 7, offset: 9430},
						val:        "Han",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 377, col: 7, offset: 9442},
						val:        "Zs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 378, col: 7, offset: 9453},
						val:        "Zp",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 379, col: 7, offset: 9464},
						val:        "Zl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 380, col: 7, offset: 9475},
						val:        "Yi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 381, col: 7, offset: 9486},
						val:        "So",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 382, col: 7, offset: 9497},
						val:        "Sm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 383, col: 7, offset: 9508},
						val:        "Sk",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 384, col: 7, offset: 9519},
						val:        "Sc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 385, col: 7, offset: 9530},
						val:        "Ps",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 386, col: 7, offset: 9541},
						val:        "Po",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 387, col: 7, offset: 9552},
						val:        "Pi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 388, col: 7, offset: 9563},
						val:        "Pf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 389, col: 7, offset: 9574},
						val:        "Pe",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 390, col: 7, offset: 9585},
						val:        "Pd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 7, offset: 9596},
						val:        "Pc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 392, col: 7, offset: 9607},
						val:        "No",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 393, col: 7, offset: 9618},
						val:        "Nl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 394, col: 7, offset: 9629},
						val:        "Nd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 395, col: 7, offset: 9640},
						val:        "Mn",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 396, col: 7, offset: 9651},
						val:        "Me",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 7, offset: 9662},
						val:        "Mc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 398, col: 7, offset: 9673},
						val:        "Lu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 399, col: 7, offset: 9684},
						val:        "Lt",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 400, col: 7, offset: 9695},
						val:        "Lo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 401, col: 7, offset: 9706},
						val:        "Lm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 402, col: 7, offset: 9717},
						val:        "Ll",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 403, col: 7, offset: 9728},
						val:        "Cs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 404, col: 7, offset: 9739},
						val:        "Co",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 405, col: 7, offset: 9750},
						val:        "Cf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 406, col: 7, offset: 9761},
						val:        "Cc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 407, col: 7, offset: 9772},
						val:        "Z",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 408, col: 7, offset: 9782},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 409, col: 7, offset: 9792},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 410, col: 7, offset: 9802},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 411, col: 7, offset: 9812},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 412, col: 7, offset: 9822},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 413, col: 7, offset: 9832},
						val:        "C",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "ReservedWord",
			pos:  position{line: 415, col: 1, offset: 9837},
			expr: &choiceExpr{
				pos: position{line: 418, col: 2, offset: 9908},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 418, col: 2, offset: 9908},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 418, col: 2, offset: 9908},
								val:        "break",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 418, col: 10, offset: 9916},
								expr: &ruleRefExpr{
									pos:  position{line: 418, col: 11, offset: 9917},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 419, col: 4, offset: 9935},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 419, col: 4, offset: 9935},
								val:        "case",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 419, col: 11, offset: 9942},
								expr: &ruleRefExpr{
									pos:  position{line: 419, col: 12, offset: 9943},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 420, col: 4, offset: 9961},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 420, col: 4, offset: 9961},
								val:        "chan",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 420, col: 11, offset: 9968},
								expr: &ruleRefExpr{
									pos:  position{line: 420, col: 12, offset: 9969},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 421, col: 4, offset: 9987},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 421, col: 4, offset: 9987},
								val:        "const",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 421, col: 12, offset: 9995},
								expr: &ruleRefExpr{
									pos:  position{line: 421, col: 13, offset: 9996},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 422, col: 4, offset: 10014},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 422, col: 4, offset: 10014},
								val:        "continue",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 422, col: 15, offset: 10025},
								expr: &ruleRefExpr{
									pos:  position{line: 422, col: 16, offset: 10026},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 423, col: 4, offset: 10044},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 423, col: 4, offset: 10044},
								val:        "default",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 423, col: 14, offset: 10054},
								expr: &ruleRefExpr{
									pos:  position{line: 423, col: 15, offset: 10055},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 424, col: 4, offset: 10073},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 424, col: 4, offset: 10073},
								val:        "defer",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 424, col: 12, offset: 10081},
								expr: &ruleRefExpr{
									pos:  position{line: 424, col: 13, offset: 10082},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 425, col: 4, offset: 10100},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 425, col: 4, offset: 10100},
								val:        "else",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 425, col: 11, offset: 10107},
								expr: &ruleRefExpr{
									pos:  position{line: 425, col: 12, offset: 10108},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 426, col: 4, offset: 10126},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 426, col: 4, offset: 10126},
								val:        "fallthrough",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 426, col: 18, offset: 10140},
								expr: &ruleRefExpr{
									pos:  position{line: 426, col: 19, offset: 10141},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 427, col: 4, offset: 10159},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 427, col: 4, offset: 10159},
								val:        "for",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 427, col: 10, offset: 10165},
								expr: &ruleRefExpr{
									pos:  position{line: 427, col: 11, offset: 10166},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 428, col: 4, offset: 10184},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 428, col: 4, offset: 10184},
								val:        "func",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 428, col: 11, offset: 10191},
								expr: &ruleRefExpr{
									pos:  position{line: 428, col: 12, offset: 10192},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 429, col: 4, offset: 10210},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 429, col: 4, offset: 10210},
								val:        "goto",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 429, col: 11, offset: 10217},
								expr: &ruleRefExpr{
									pos:  position{line: 429, col: 12, offset: 10218},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 430, col: 4, offset: 10236},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 430, col: 4, offset: 10236},
								val:        "go",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 430, col: 9, offset: 10241},
								expr: &ruleRefExpr{
									pos:  position{line: 430, col: 10, offset: 10242},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 431, col: 4, offset: 10260},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 431, col: 4, offset: 10260},
								val:        "if",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 431, col: 9, offset: 10265},
								expr: &ruleRefExpr{
									pos:  position{line: 431, col: 10, offset: 10266},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 432, col: 4, offset: 10284},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 432, col: 4, offset: 10284},
								val:        "import",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 432, col: 13, offset: 10293},
								expr: &ruleRefExpr{
									pos:  position{line: 432, col: 14, offset: 10294},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 433, col: 4, offset: 10312},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 433, col: 4, offset: 10312},
								val:        "interface",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 433, col: 16, offset: 10324},
								expr: &ruleRefExpr{
									pos:  position{line: 433, col: 17, offset: 10325},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 434, col: 4, offset: 10343},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 434, col: 4, offset: 10343},
								val:        "map",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 434, col: 10, offset: 10349},
								expr: &ruleRefExpr{
									pos:  position{line: 434, col: 11, offset: 10350},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 435, col: 4, offset: 10368},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 435, col: 4, offset: 10368},
								val:        "package",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 435, col: 14, offset: 10378},
								expr: &ruleRefExpr{
									pos:  position{line: 435, col: 15, offset: 10379},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 436, col: 4, offset: 10397},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 436, col: 4, offset: 10397},
								val:        "range",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 436, col: 12, offset: 10405},
								expr: &ruleRefExpr{
									pos:  position{line: 436, col: 13, offset: 10406},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 437, col: 4, offset: 10424},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 437, col: 4, offset: 10424},
								val:        "return",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 437, col: 13, offset: 10433},
								expr: &ruleRefExpr{
									pos:  position{line: 437, col: 14, offset: 10434},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 438, col: 4, offset: 10452},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 438, col: 4, offset: 10452},
								val:        "select",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 438, col: 13, offset: 10461},
								expr: &ruleRefExpr{
									pos:  position{line: 438, col: 14, offset: 10462},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 439, col: 4, offset: 10480},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 439, col: 4, offset: 10480},
								val:        "struct",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 439, col: 13, offset: 10489},
								expr: &ruleRefExpr{
									pos:  position{line: 439, col: 14, offset: 10490},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 440, col: 4, offset: 10508},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 440, col: 4, offset: 10508},
								val:        "switch",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 440, col: 13, offset: 10517},
								expr: &ruleRefExpr{
									pos:  position{line: 440, col: 14, offset: 10518},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 441, col: 4, offset: 10536},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 441, col: 4, offset: 10536},
								val:        "type",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 441, col: 11, offset: 10543},
								expr: &ruleRefExpr{
									pos:  position{line: 441, col: 12, offset: 10544},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 442, col: 4, offset: 10562},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 442, col: 4, offset: 10562},
								val:        "var",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 442, col: 10, offset: 10568},
								expr: &ruleRefExpr{
									pos:  position{line: 442, col: 11, offset: 10569},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 446, col: 4, offset: 10668},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 446, col: 4, offset: 10668},
								val:        "bool",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 446, col: 11, offset: 10675},
								expr: &ruleRefExpr{
									pos:  position{line: 446, col: 12, offset: 10676},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 447, col: 4, offset: 10694},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 447, col: 4, offset: 10694},
								val:        "byte",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 447, col: 11, offset: 10701},
								expr: &ruleRefExpr{
									pos:  position{line: 447, col: 12, offset: 10702},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 448, col: 4, offset: 10720},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 448, col: 4, offset: 10720},
								val:        "complex64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 448, col: 16, offset: 10732},
								expr: &ruleRefExpr{
									pos:  position{line: 448, col: 17, offset: 10733},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 449, col: 4, offset: 10751},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 449, col: 4, offset: 10751},
								val:        "complex128",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 449, col: 17, offset: 10764},
								expr: &ruleRefExpr{
									pos:  position{line: 449, col: 18, offset: 10765},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 450, col: 4, offset: 10783},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 450, col: 4, offset: 10783},
								val:        "error",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 450, col: 12, offset: 10791},
								expr: &ruleRefExpr{
									pos:  position{line: 450, col: 13, offset: 10792},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 451, col: 4, offset: 10810},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 451, col: 4, offset: 10810},
								val:        "float32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 451, col: 14, offset: 10820},
								expr: &ruleRefExpr{
									pos:  position{line: 451, col: 15, offset: 10821},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 452, col: 4, offset: 10839},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 452, col: 4, offset: 10839},
								val:        "float64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 452, col: 14, offset: 10849},
								expr: &ruleRefExpr{
									pos:  position{line: 452, col: 15, offset: 10850},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 453, col: 4, offset: 10868},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 453, col: 4, offset: 10868},
								val:        "int8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 453, col: 11, offset: 10875},
								expr: &ruleRefExpr{
									pos:  position{line: 453, col: 12, offset: 10876},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 454, col: 4, offset: 10894},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 454, col: 4, offset: 10894},
								val:        "int16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 454, col: 12, offset: 10902},
								expr: &ruleRefExpr{
									pos:  position{line: 454, col: 13, offset: 10903},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 455, col: 4, offset: 10921},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 455, col: 4, offset: 10921},
								val:        "int32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 455, col: 12, offset: 10929},
								expr: &ruleRefExpr{
									pos:  position{line: 455, col: 13, offset: 10930},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 456, col: 4, offset: 10948},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 456, col: 4, offset: 10948},
								val:        "int64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 456, col: 12, offset: 10956},
								expr: &ruleRefExpr{
									pos:  position{line: 456, col: 13, offset: 10957},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 457, col: 4, offset: 10975},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 457, col: 4, offset: 10975},
								val:        "int",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 457, col: 10, offset: 10981},
								expr: &ruleRefExpr{
									pos:  position{line: 457, col: 11, offset: 10982},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 458, col: 4, offset: 11000},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 458, col: 4, offset: 11000},
								val:        "rune",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 458, col: 11, offset: 11007},
								expr: &ruleRefExpr{
									pos:  position{line: 458, col: 12, offset: 11008},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 459, col: 4, offset: 11026},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 459, col: 4, offset: 11026},
								val:        "string",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 459, col: 13, offset: 11035},
								expr: &ruleRefExpr{
									pos:  position{line: 459, col: 14, offset: 11036},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 460, col: 4, offset: 11054},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 460, col: 4, offset: 11054},
								val:        "uint8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 460, col: 12, offset: 11062},
								expr: &ruleRefExpr{
									pos:  position{line: 460, col: 13, offset: 11063},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 461, col: 4, offset: 11081},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 461, col: 4, offset: 11081},
								val:        "uint16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 461, col: 13, offset: 11090},
								expr: &ruleRefExpr{
									pos:  position{line: 461, col: 14, offset: 11091},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 462, col: 4, offset: 11109},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 462, col: 4, offset: 11109},
								val:        "uint32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 462, col: 13, offset: 11118},
								expr: &ruleRefExpr{
									pos:  position{line: 462, col: 14, offset: 11119},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 463, col: 4, offset: 11137},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 463, col: 4, offset: 11137},
								val:        "uint64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 463, col: 13, offset: 11146},
								expr: &ruleRefExpr{
									pos:  position{line: 463, col: 14, offset: 11147},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 464, col: 4, offset: 11165},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 464, col: 4, offset: 11165},
								val:        "uintptr",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 464, col: 14, offset: 11175},
								expr: &ruleRefExpr{
									pos:  position{line: 464, col: 15, offset: 11176},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 465, col: 4, offset: 11194},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 465, col: 4, offset: 11194},
								val:        "uint",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 465, col: 11, offset: 11201},
								expr: &ruleRefExpr{
									pos:  position{line: 465, col: 12, offset: 11202},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 466, col: 4, offset: 11220},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 466, col: 4, offset: 11220},
								val:        "true",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 466, col: 11, offset: 11227},
								expr: &ruleRefExpr{
									pos:  position{line: 466, col: 12, offset: 11228},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 467, col: 4, offset: 11246},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 467, col: 4, offset: 11246},
								val:        "false",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 467, col: 12, offset: 11254},
								expr: &ruleRefExpr{
									pos:  position{line: 467, col: 13, offset: 11255},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 468, col: 4, offset: 11273},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 468, col: 4, offset: 11273},
								val:        "iota",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 468, col: 11, offset: 11280},
								expr: &ruleRefExpr{
									pos:  position{line: 468, col: 12, offset: 11281},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 469, col: 4, offset: 11299},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 469, col: 4, offset: 11299},
								val:        "nil",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 469, col: 10, offset: 11305},
								expr: &ruleRefExpr{
									pos:  position{line: 469, col: 11, offset: 11306},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 470, col: 4, offset: 11324},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 470, col: 4, offset: 11324},
								val:        "append",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 470, col: 13, offset: 11333},
								expr: &ruleRefExpr{
									pos:  position{line: 470, col: 14, offset: 11334},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 471, col: 4, offset: 11352},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 471, col: 4, offset: 11352},
								val:        "cap",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 471, col: 10, offset: 11358},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 11, offset: 11359},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 472, col: 4, offset: 11377},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 472, col: 4, offset: 11377},
								val:        "close",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 472, col: 12, offset: 11385},
								expr: &ruleRefExpr{
									pos:  position{line: 472, col: 13, offset: 11386},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 473, col: 4, offset: 11404},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 473, col: 4, offset: 11404},
								val:        "complex",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 473, col: 14, offset: 11414},
								expr: &ruleRefExpr{
									pos:  position{line: 473, col: 15, offset: 11415},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 474, col: 4, offset: 11433},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 474, col: 4, offset: 11433},
								val:        "copy",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 474, col: 11, offset: 11440},
								expr: &ruleRefExpr{
									pos:  position{line: 474, col: 12, offset: 11441},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 475, col: 4, offset: 11459},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 475, col: 4, offset: 11459},
								val:        "delete",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 475, col: 13, offset: 11468},
								expr: &ruleRefExpr{
									pos:  position{line: 475, col: 14, offset: 11469},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 476, col: 4, offset: 11487},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 476, col: 4, offset: 11487},
								val:        "imag",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 476, col: 11, offset: 11494},
								expr: &ruleRefExpr{
									pos:  position{line: 476, col: 12, offset: 11495},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 477, col: 4, offset: 11513},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 477, col: 4, offset: 11513},
								val:        "len",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 477, col: 10, offset: 11519},
								expr: &ruleRefExpr{
									pos:  position{line: 477, col: 11, offset: 11520},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 478, col: 4, offset: 11538},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 478, col: 4, offset: 11538},
								val:        "make",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 478, col: 11, offset: 11545},
								expr: &ruleRefExpr{
									pos:  position{line: 478, col: 12, offset: 11546},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 479, col: 4, offset: 11564},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 479, col: 4, offset: 11564},
								val:        "new",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 479, col: 10, offset: 11570},
								expr: &ruleRefExpr{
									pos:  position{line: 479, col: 11, offset: 11571},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 480, col: 4, offset: 11589},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 480, col: 4, offset: 11589},
								val:        "panic",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 480, col: 12, offset: 11597},
								expr: &ruleRefExpr{
									pos:  position{line: 480, col: 13, offset: 11598},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 481, col: 4, offset: 11616},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 481, col: 4, offset: 11616},
								val:        "println",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 481, col: 14, offset: 11626},
								expr: &ruleRefExpr{
									pos:  position{line: 481, col: 15, offset: 11627},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 482, col: 4, offset: 11645},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 482, col: 4, offset: 11645},
								val:        "print",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 482, col: 12, offset: 11653},
								expr: &ruleRefExpr{
									pos:  position{line: 482, col: 13, offset: 11654},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 483, col: 4, offset: 11672},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 483, col: 4, offset: 11672},
								val:        "real",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 483, col: 11, offset: 11679},
								expr: &ruleRefExpr{
									pos:  position{line: 483, col: 12, offset: 11680},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 484, col: 4, offset: 11698},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 484, col: 4, offset: 11698},
								val:        "recover",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 484, col: 14, offset: 11708},
								expr: &ruleRefExpr{
									pos:  position{line: 484, col: 15, offset: 11709},
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
			pos:  position{line: 486, col: 1, offset: 11725},
			expr: &actionExpr{
				pos: position{line: 486, col: 14, offset: 11740},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 486, col: 14, offset: 11740},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 491, col: 1, offset: 11815},
			expr: &actionExpr{
				pos: position{line: 491, col: 13, offset: 11829},
				run: (*parser).callonCodeBlock1,
				expr: &seqExpr{
					pos: position{line: 491, col: 13, offset: 11829},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 491, col: 13, offset: 11829},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 491, col: 17, offset: 11833},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 491, col: 22, offset: 11838},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 497, col: 1, offset: 11936},
			expr: &zeroOrMoreExpr{
				pos: position{line: 497, col: 8, offset: 11945},
				expr: &choiceExpr{
					pos: position{line: 497, col: 10, offset: 11947},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 497, col: 10, offset: 11947},
							expr: &seqExpr{
								pos: position{line: 497, col: 12, offset: 11949},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 497, col: 12, offset: 11949},
										expr: &charClassMatcher{
											pos:        position{line: 497, col: 13, offset: 11950},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 497, col: 18, offset: 11955},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 497, col: 34, offset: 11971},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 497, col: 34, offset: 11971},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 497, col: 38, offset: 11975},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 497, col: 43, offset: 11980},
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
			pos:  position{line: 499, col: 1, offset: 11988},
			expr: &zeroOrMoreExpr{
				pos: position{line: 499, col: 6, offset: 11995},
				expr: &choiceExpr{
					pos: position{line: 499, col: 8, offset: 11997},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 499, col: 8, offset: 11997},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 499, col: 21, offset: 12010},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 499, col: 27, offset: 12016},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 500, col: 1, offset: 12027},
			expr: &zeroOrMoreExpr{
				pos: position{line: 500, col: 5, offset: 12033},
				expr: &choiceExpr{
					pos: position{line: 500, col: 7, offset: 12035},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 500, col: 7, offset: 12035},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 500, col: 20, offset: 12048},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 502, col: 1, offset: 12085},
			expr: &charClassMatcher{
				pos:        position{line: 502, col: 14, offset: 12100},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 503, col: 1, offset: 12108},
			expr: &litMatcher{
				pos:        position{line: 503, col: 7, offset: 12116},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 504, col: 1, offset: 12121},
			expr: &choiceExpr{
				pos: position{line: 504, col: 7, offset: 12129},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 504, col: 7, offset: 12129},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 504, col: 7, offset: 12129},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 504, col: 10, offset: 12132},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 504, col: 16, offset: 12138},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 504, col: 16, offset: 12138},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 504, col: 18, offset: 12140},
								expr: &ruleRefExpr{
									pos:  position{line: 504, col: 18, offset: 12140},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 504, col: 37, offset: 12159},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 504, col: 43, offset: 12165},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 504, col: 43, offset: 12165},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 504, col: 46, offset: 12168},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 506, col: 1, offset: 12173},
			expr: &notExpr{
				pos: position{line: 506, col: 7, offset: 12181},
				expr: &anyMatcher{
					line: 506, col: 8, offset: 12182,
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
	return nil, errors.New("identifier is a reserved word")
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
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if p.pt.w > 0 || p.pt.offset > 0 {
		if buf.Len() > 0 {
			buf.WriteString(":")
		}
		// parsing has started, so add the position of the error
		buf.WriteString(fmt.Sprintf("%d:%d (%d)", p.pt.line, p.pt.col, p.pt.offset))
	}
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

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
