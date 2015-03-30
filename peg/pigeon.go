package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
						&ruleRefExpr{
							pos:  position{line: 5, col: 65, offset: 84},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Initializer",
			pos:  position{line: 24, col: 1, offset: 525},
			expr: &actionExpr{
				pos: position{line: 24, col: 15, offset: 541},
				run: (*parser).callonInitializer1,
				expr: &seqExpr{
					pos: position{line: 24, col: 15, offset: 541},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 24, col: 15, offset: 541},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 20, offset: 546},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 30, offset: 556},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 28, col: 1, offset: 586},
			expr: &actionExpr{
				pos: position{line: 28, col: 8, offset: 595},
				run: (*parser).callonRule1,
				expr: &seqExpr{
					pos: position{line: 28, col: 8, offset: 595},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 8, offset: 595},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 13, offset: 600},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 28, offset: 615},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 31, offset: 618},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 28, col: 39, offset: 626},
								expr: &seqExpr{
									pos: position{line: 28, col: 41, offset: 628},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 28, col: 41, offset: 628},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 55, offset: 642},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 648},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 71, offset: 658},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 74, offset: 661},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 79, offset: 666},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 90, offset: 677},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 41, col: 1, offset: 961},
			expr: &ruleRefExpr{
				pos:  position{line: 41, col: 14, offset: 976},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 43, col: 1, offset: 988},
			expr: &actionExpr{
				pos: position{line: 43, col: 14, offset: 1003},
				run: (*parser).callonChoiceExpr1,
				expr: &seqExpr{
					pos: position{line: 43, col: 14, offset: 1003},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 43, col: 14, offset: 1003},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 20, offset: 1009},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 43, col: 31, offset: 1020},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 36, offset: 1025},
								expr: &seqExpr{
									pos: position{line: 43, col: 38, offset: 1027},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 43, col: 38, offset: 1027},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 43, col: 41, offset: 1030},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 45, offset: 1034},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 48, offset: 1037},
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
			pos:  position{line: 58, col: 1, offset: 1442},
			expr: &actionExpr{
				pos: position{line: 58, col: 14, offset: 1457},
				run: (*parser).callonActionExpr1,
				expr: &seqExpr{
					pos: position{line: 58, col: 14, offset: 1457},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 58, col: 14, offset: 1457},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 19, offset: 1462},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 58, col: 27, offset: 1470},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 58, col: 32, offset: 1475},
								expr: &seqExpr{
									pos: position{line: 58, col: 34, offset: 1477},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 58, col: 34, offset: 1477},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 58, col: 37, offset: 1480},
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
			pos:  position{line: 72, col: 1, offset: 1746},
			expr: &actionExpr{
				pos: position{line: 72, col: 11, offset: 1758},
				run: (*parser).callonSeqExpr1,
				expr: &seqExpr{
					pos: position{line: 72, col: 11, offset: 1758},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 72, col: 11, offset: 1758},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 17, offset: 1764},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 72, col: 29, offset: 1776},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 72, col: 34, offset: 1781},
								expr: &seqExpr{
									pos: position{line: 72, col: 36, offset: 1783},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 72, col: 36, offset: 1783},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 72, col: 39, offset: 1786},
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
			pos:  position{line: 85, col: 1, offset: 2137},
			expr: &choiceExpr{
				pos: position{line: 85, col: 15, offset: 2153},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 85, col: 15, offset: 2153},
						run: (*parser).callonLabeledExpr2,
						expr: &seqExpr{
							pos: position{line: 85, col: 15, offset: 2153},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 85, col: 15, offset: 2153},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 21, offset: 2159},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 32, offset: 2170},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 85, col: 35, offset: 2173},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 39, offset: 2177},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 85, col: 42, offset: 2180},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 47, offset: 2185},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 5, offset: 2358},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 93, col: 1, offset: 2372},
			expr: &choiceExpr{
				pos: position{line: 93, col: 16, offset: 2389},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 93, col: 16, offset: 2389},
						run: (*parser).callonPrefixedExpr2,
						expr: &seqExpr{
							pos: position{line: 93, col: 16, offset: 2389},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 93, col: 16, offset: 2389},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 19, offset: 2392},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 93, col: 30, offset: 2403},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 93, col: 33, offset: 2406},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 38, offset: 2411},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 5, offset: 2693},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 106, col: 1, offset: 2707},
			expr: &actionExpr{
				pos: position{line: 106, col: 14, offset: 2722},
				run: (*parser).callonPrefixedOp1,
				expr: &choiceExpr{
					pos: position{line: 106, col: 16, offset: 2724},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 16, offset: 2724},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 22, offset: 2730},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 110, col: 1, offset: 2772},
			expr: &choiceExpr{
				pos: position{line: 110, col: 16, offset: 2789},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 110, col: 16, offset: 2789},
						run: (*parser).callonSuffixedExpr2,
						expr: &seqExpr{
							pos: position{line: 110, col: 16, offset: 2789},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 110, col: 16, offset: 2789},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 21, offset: 2794},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 33, offset: 2806},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 36, offset: 2809},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 39, offset: 2812},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 129, col: 5, offset: 3342},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 131, col: 1, offset: 3356},
			expr: &actionExpr{
				pos: position{line: 131, col: 14, offset: 3371},
				run: (*parser).callonSuffixedOp1,
				expr: &choiceExpr{
					pos: position{line: 131, col: 16, offset: 3373},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 131, col: 16, offset: 3373},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 22, offset: 3379},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 28, offset: 3385},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 135, col: 1, offset: 3427},
			expr: &choiceExpr{
				pos: position{line: 135, col: 15, offset: 3443},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 135, col: 15, offset: 3443},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 28, offset: 3456},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 47, offset: 3475},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 60, offset: 3488},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 74, offset: 3502},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 135, col: 93, offset: 3521},
						run: (*parser).callonPrimaryExpr7,
						expr: &seqExpr{
							pos: position{line: 135, col: 93, offset: 3521},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 135, col: 93, offset: 3521},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 97, offset: 3525},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 135, col: 100, offset: 3528},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 105, offset: 3533},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 116, offset: 3544},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 135, col: 119, offset: 3547},
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
			pos:  position{line: 138, col: 1, offset: 3576},
			expr: &actionExpr{
				pos: position{line: 138, col: 15, offset: 3592},
				run: (*parser).callonRuleRefExpr1,
				expr: &seqExpr{
					pos: position{line: 138, col: 15, offset: 3592},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 138, col: 15, offset: 3592},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 138, col: 20, offset: 3597},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 138, col: 35, offset: 3612},
							expr: &seqExpr{
								pos: position{line: 138, col: 38, offset: 3615},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 138, col: 38, offset: 3615},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 138, col: 41, offset: 3618},
										expr: &seqExpr{
											pos: position{line: 138, col: 43, offset: 3620},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 138, col: 43, offset: 3620},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 138, col: 57, offset: 3634},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 63, offset: 3640},
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
			pos:  position{line: 143, col: 1, offset: 3756},
			expr: &actionExpr{
				pos: position{line: 143, col: 20, offset: 3777},
				run: (*parser).callonSemanticPredExpr1,
				expr: &seqExpr{
					pos: position{line: 143, col: 20, offset: 3777},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 143, col: 20, offset: 3777},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 23, offset: 3780},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 143, col: 38, offset: 3795},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 143, col: 41, offset: 3798},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 46, offset: 3803},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 154, col: 1, offset: 4080},
			expr: &actionExpr{
				pos: position{line: 154, col: 18, offset: 4099},
				run: (*parser).callonSemanticPredOp1,
				expr: &choiceExpr{
					pos: position{line: 154, col: 20, offset: 4101},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 20, offset: 4101},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 154, col: 26, offset: 4107},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 158, col: 1, offset: 4149},
			expr: &choiceExpr{
				pos: position{line: 158, col: 13, offset: 4163},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 158, col: 13, offset: 4163},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 19, offset: 4169},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 26, offset: 4176},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 37, offset: 4187},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 160, col: 1, offset: 4197},
			expr: &anyMatcher{
				line: 160, col: 14, offset: 4212,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 161, col: 1, offset: 4214},
			expr: &choiceExpr{
				pos: position{line: 161, col: 11, offset: 4226},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 161, col: 11, offset: 4226},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 161, col: 30, offset: 4245},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 162, col: 1, offset: 4263},
			expr: &seqExpr{
				pos: position{line: 162, col: 20, offset: 4284},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 162, col: 20, offset: 4284},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 162, col: 25, offset: 4289},
						expr: &seqExpr{
							pos: position{line: 162, col: 27, offset: 4291},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 162, col: 27, offset: 4291},
									expr: &litMatcher{
										pos:        position{line: 162, col: 28, offset: 4292},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 162, col: 33, offset: 4297},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 162, col: 47, offset: 4311},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 163, col: 1, offset: 4316},
			expr: &seqExpr{
				pos: position{line: 163, col: 36, offset: 4353},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 163, col: 36, offset: 4353},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 163, col: 41, offset: 4358},
						expr: &seqExpr{
							pos: position{line: 163, col: 43, offset: 4360},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 163, col: 43, offset: 4360},
									expr: &choiceExpr{
										pos: position{line: 163, col: 46, offset: 4363},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 163, col: 46, offset: 4363},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 163, col: 53, offset: 4370},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 163, col: 59, offset: 4376},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 163, col: 73, offset: 4390},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 164, col: 1, offset: 4395},
			expr: &seqExpr{
				pos: position{line: 164, col: 21, offset: 4417},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 164, col: 21, offset: 4417},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 164, col: 26, offset: 4422},
						expr: &seqExpr{
							pos: position{line: 164, col: 28, offset: 4424},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 164, col: 28, offset: 4424},
									expr: &ruleRefExpr{
										pos:  position{line: 164, col: 29, offset: 4425},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 164, col: 33, offset: 4429},
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
			pos:  position{line: 166, col: 1, offset: 4444},
			expr: &choiceExpr{
				pos: position{line: 166, col: 14, offset: 4459},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 166, col: 14, offset: 4459},
						run: (*parser).callonIdentifier2,
						expr: &seqExpr{
							pos: position{line: 166, col: 14, offset: 4459},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 166, col: 14, offset: 4459},
									expr: &ruleRefExpr{
										pos:  position{line: 166, col: 15, offset: 4460},
										name: "ReservedWord",
									},
								},
								&labeledExpr{
									pos:   position{line: 166, col: 28, offset: 4473},
									label: "ident",
									expr: &ruleRefExpr{
										pos:  position{line: 166, col: 34, offset: 4479},
										name: "IdentifierName",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 168, col: 5, offset: 4522},
						run: (*parser).callonIdentifier8,
						expr: &ruleRefExpr{
							pos:  position{line: 168, col: 5, offset: 4522},
							name: "ReservedWord",
						},
					},
				},
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 171, col: 1, offset: 4641},
			expr: &actionExpr{
				pos: position{line: 171, col: 18, offset: 4660},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 171, col: 18, offset: 4660},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 171, col: 18, offset: 4660},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 171, col: 34, offset: 4676},
							expr: &ruleRefExpr{
								pos:  position{line: 171, col: 34, offset: 4676},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 174, col: 1, offset: 4758},
			expr: &charClassMatcher{
				pos:        position{line: 174, col: 19, offset: 4778},
				val:        "[\\pL_]",
				chars:      []rune{'_'},
				classes:    []*unicode.RangeTable{rangeTable("L")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 175, col: 1, offset: 4785},
			expr: &choiceExpr{
				pos: position{line: 175, col: 18, offset: 4804},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 175, col: 18, offset: 4804},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 175, col: 36, offset: 4822},
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
			pos:  position{line: 177, col: 1, offset: 4832},
			expr: &actionExpr{
				pos: position{line: 177, col: 14, offset: 4847},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 177, col: 14, offset: 4847},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 177, col: 14, offset: 4847},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 177, col: 18, offset: 4851},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 177, col: 32, offset: 4865},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 177, col: 39, offset: 4872},
								expr: &litMatcher{
									pos:        position{line: 177, col: 39, offset: 4872},
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
			pos:  position{line: 187, col: 1, offset: 5098},
			expr: &choiceExpr{
				pos: position{line: 187, col: 17, offset: 5116},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 187, col: 17, offset: 5116},
						run: (*parser).callonStringLiteral2,
						expr: &choiceExpr{
							pos: position{line: 187, col: 19, offset: 5118},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 187, col: 19, offset: 5118},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 187, col: 19, offset: 5118},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 187, col: 23, offset: 5122},
											expr: &ruleRefExpr{
												pos:  position{line: 187, col: 23, offset: 5122},
												name: "DoubleStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 187, col: 41, offset: 5140},
											val:        "\"",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 187, col: 47, offset: 5146},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 187, col: 47, offset: 5146},
											val:        "'",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 187, col: 51, offset: 5150},
											name: "SingleStringChar",
										},
										&litMatcher{
											pos:        position{line: 187, col: 68, offset: 5167},
											val:        "'",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 187, col: 74, offset: 5173},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 187, col: 74, offset: 5173},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 187, col: 78, offset: 5177},
											expr: &ruleRefExpr{
												pos:  position{line: 187, col: 78, offset: 5177},
												name: "RawStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 187, col: 93, offset: 5192},
											val:        "`",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 189, col: 5, offset: 5265},
						run: (*parser).callonStringLiteral18,
						expr: &choiceExpr{
							pos: position{line: 189, col: 7, offset: 5267},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 189, col: 9, offset: 5269},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 189, col: 9, offset: 5269},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 189, col: 13, offset: 5273},
											expr: &ruleRefExpr{
												pos:  position{line: 189, col: 13, offset: 5273},
												name: "DoubleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 189, col: 33, offset: 5293},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 189, col: 33, offset: 5293},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 189, col: 39, offset: 5299},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 189, col: 51, offset: 5311},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 189, col: 51, offset: 5311},
											val:        "'",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 189, col: 55, offset: 5315},
											name: "SingleStringChar",
										},
										&choiceExpr{
											pos: position{line: 189, col: 74, offset: 5334},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 189, col: 74, offset: 5334},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 189, col: 80, offset: 5340},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 189, col: 90, offset: 5350},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 189, col: 90, offset: 5350},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 189, col: 94, offset: 5354},
											expr: &ruleRefExpr{
												pos:  position{line: 189, col: 94, offset: 5354},
												name: "RawStringChar",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 189, col: 109, offset: 5369},
											name: "EOF",
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
			name: "DoubleStringChar",
			pos:  position{line: 193, col: 1, offset: 5471},
			expr: &choiceExpr{
				pos: position{line: 193, col: 20, offset: 5492},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 193, col: 20, offset: 5492},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 193, col: 20, offset: 5492},
								expr: &choiceExpr{
									pos: position{line: 193, col: 23, offset: 5495},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 23, offset: 5495},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 193, col: 29, offset: 5501},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 193, col: 36, offset: 5508},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 193, col: 42, offset: 5514},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 193, col: 55, offset: 5527},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 193, col: 55, offset: 5527},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 193, col: 60, offset: 5532},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 194, col: 1, offset: 5551},
			expr: &choiceExpr{
				pos: position{line: 194, col: 20, offset: 5572},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 194, col: 20, offset: 5572},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 194, col: 20, offset: 5572},
								expr: &choiceExpr{
									pos: position{line: 194, col: 23, offset: 5575},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 194, col: 23, offset: 5575},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 194, col: 29, offset: 5581},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 36, offset: 5588},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 194, col: 42, offset: 5594},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 194, col: 55, offset: 5607},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 194, col: 55, offset: 5607},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 194, col: 60, offset: 5612},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 195, col: 1, offset: 5631},
			expr: &seqExpr{
				pos: position{line: 195, col: 17, offset: 5649},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 195, col: 17, offset: 5649},
						expr: &litMatcher{
							pos:        position{line: 195, col: 18, offset: 5650},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 195, col: 22, offset: 5654},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 197, col: 1, offset: 5666},
			expr: &choiceExpr{
				pos: position{line: 197, col: 22, offset: 5689},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 197, col: 22, offset: 5689},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 28, offset: 5695},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 198, col: 1, offset: 5716},
			expr: &choiceExpr{
				pos: position{line: 198, col: 22, offset: 5739},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 198, col: 22, offset: 5739},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 28, offset: 5745},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 200, col: 1, offset: 5767},
			expr: &choiceExpr{
				pos: position{line: 200, col: 24, offset: 5792},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 200, col: 24, offset: 5792},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 43, offset: 5811},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 57, offset: 5825},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 69, offset: 5837},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 89, offset: 5857},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 201, col: 1, offset: 5876},
			expr: &choiceExpr{
				pos: position{line: 201, col: 20, offset: 5897},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 201, col: 20, offset: 5897},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 26, offset: 5903},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 32, offset: 5909},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 38, offset: 5915},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 44, offset: 5921},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 50, offset: 5927},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 56, offset: 5933},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 62, offset: 5939},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 202, col: 1, offset: 5944},
			expr: &seqExpr{
				pos: position{line: 202, col: 15, offset: 5960},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 202, col: 15, offset: 5960},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 26, offset: 5971},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 202, col: 37, offset: 5982},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 203, col: 1, offset: 5993},
			expr: &seqExpr{
				pos: position{line: 203, col: 13, offset: 6007},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 203, col: 13, offset: 6007},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 203, col: 17, offset: 6011},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 203, col: 26, offset: 6020},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 204, col: 1, offset: 6029},
			expr: &seqExpr{
				pos: position{line: 204, col: 21, offset: 6051},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 204, col: 21, offset: 6051},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 25, offset: 6055},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 34, offset: 6064},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 43, offset: 6073},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 52, offset: 6082},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 61, offset: 6091},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 70, offset: 6100},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 79, offset: 6109},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 88, offset: 6118},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 205, col: 1, offset: 6127},
			expr: &seqExpr{
				pos: position{line: 205, col: 22, offset: 6150},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 205, col: 22, offset: 6150},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 205, col: 26, offset: 6154},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 205, col: 35, offset: 6163},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 205, col: 44, offset: 6172},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 205, col: 53, offset: 6181},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 207, col: 1, offset: 6191},
			expr: &charClassMatcher{
				pos:        position{line: 207, col: 14, offset: 6206},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 208, col: 1, offset: 6212},
			expr: &charClassMatcher{
				pos:        position{line: 208, col: 16, offset: 6229},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 209, col: 1, offset: 6235},
			expr: &charClassMatcher{
				pos:        position{line: 209, col: 12, offset: 6248},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 211, col: 1, offset: 6259},
			expr: &actionExpr{
				pos: position{line: 211, col: 20, offset: 6280},
				run: (*parser).callonCharClassMatcher1,
				expr: &seqExpr{
					pos: position{line: 211, col: 20, offset: 6280},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 211, col: 20, offset: 6280},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 211, col: 24, offset: 6284},
							expr: &choiceExpr{
								pos: position{line: 211, col: 26, offset: 6286},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 211, col: 26, offset: 6286},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 211, col: 43, offset: 6303},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 211, col: 55, offset: 6315},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 211, col: 55, offset: 6315},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 211, col: 60, offset: 6320},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 211, col: 82, offset: 6342},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 211, col: 86, offset: 6346},
							expr: &litMatcher{
								pos:        position{line: 211, col: 86, offset: 6346},
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
			pos:  position{line: 216, col: 1, offset: 6451},
			expr: &seqExpr{
				pos: position{line: 216, col: 18, offset: 6470},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 216, col: 18, offset: 6470},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 216, col: 28, offset: 6480},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 216, col: 32, offset: 6484},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 217, col: 1, offset: 6494},
			expr: &choiceExpr{
				pos: position{line: 217, col: 13, offset: 6508},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 217, col: 13, offset: 6508},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 217, col: 13, offset: 6508},
								expr: &choiceExpr{
									pos: position{line: 217, col: 16, offset: 6511},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 217, col: 16, offset: 6511},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 217, col: 22, offset: 6517},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 217, col: 29, offset: 6524},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 217, col: 35, offset: 6530},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 217, col: 48, offset: 6543},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 217, col: 48, offset: 6543},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 217, col: 53, offset: 6548},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 218, col: 1, offset: 6564},
			expr: &choiceExpr{
				pos: position{line: 218, col: 19, offset: 6584},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 218, col: 19, offset: 6584},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 218, col: 25, offset: 6590},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 220, col: 1, offset: 6612},
			expr: &seqExpr{
				pos: position{line: 220, col: 22, offset: 6635},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 220, col: 22, offset: 6635},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 221, col: 7, offset: 6648},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 221, col: 7, offset: 6648},
								name: "SingleCharUnicodeClass",
							},
							&actionExpr{
								pos: position{line: 222, col: 7, offset: 6677},
								run: (*parser).callonUnicodeClassEscape5,
								expr: &seqExpr{
									pos: position{line: 222, col: 7, offset: 6677},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 222, col: 7, offset: 6677},
											expr: &litMatcher{
												pos:        position{line: 222, col: 8, offset: 6678},
												val:        "{",
												ignoreCase: false,
											},
										},
										&ruleRefExpr{
											pos:  position{line: 222, col: 12, offset: 6682},
											name: "SourceChar",
										},
									},
								},
							},
							&seqExpr{
								pos: position{line: 223, col: 7, offset: 6758},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 223, col: 7, offset: 6758},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 223, col: 11, offset: 6762},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 223, col: 24, offset: 6775},
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
			pos:  position{line: 225, col: 1, offset: 6785},
			expr: &charClassMatcher{
				pos:        position{line: 225, col: 26, offset: 6812},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 228, col: 1, offset: 6885},
			expr: &choiceExpr{
				pos: position{line: 228, col: 16, offset: 6902},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 228, col: 16, offset: 6902},
						val:        "Other_Default_Ignorable_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 229, col: 7, offset: 6945},
						val:        "Noncharacter_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 230, col: 7, offset: 6977},
						val:        "Logical_Order_Exception",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 231, col: 7, offset: 7009},
						val:        "Inscriptional_Parthian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 232, col: 7, offset: 7040},
						val:        "Other_Grapheme_Extend",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 233, col: 7, offset: 7070},
						val:        "Inscriptional_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 234, col: 7, offset: 7100},
						val:        "Terminal_Punctuation",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 7, offset: 7129},
						val:        "Meroitic_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 236, col: 7, offset: 7158},
						val:        "IDS_Trinary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 237, col: 7, offset: 7187},
						val:        "Egyptian_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 238, col: 7, offset: 7216},
						val:        "Pattern_White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 239, col: 7, offset: 7244},
						val:        "IDS_Binary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 240, col: 7, offset: 7272},
						val:        "Canadian_Aboriginal",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 241, col: 7, offset: 7300},
						val:        "Variation_Selector",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 242, col: 7, offset: 7327},
						val:        "Caucasian_Albanian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 243, col: 7, offset: 7354},
						val:        "Unified_Ideograph",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 244, col: 7, offset: 7380},
						val:        "Other_ID_Continue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 245, col: 7, offset: 7406},
						val:        "Old_South_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 246, col: 7, offset: 7432},
						val:        "Old_North_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 247, col: 7, offset: 7458},
						val:        "Other_Alphabetic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 248, col: 7, offset: 7483},
						val:        "Meroitic_Cursive",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 249, col: 7, offset: 7508},
						val:        "Imperial_Aramaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 250, col: 7, offset: 7533},
						val:        "Psalter_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 251, col: 7, offset: 7557},
						val:        "Other_Uppercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 252, col: 7, offset: 7581},
						val:        "Other_Lowercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 253, col: 7, offset: 7605},
						val:        "ASCII_Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 254, col: 7, offset: 7629},
						val:        "Quotation_Mark",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 255, col: 7, offset: 7652},
						val:        "Pattern_Syntax",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 7, offset: 7675},
						val:        "Other_ID_Start",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 257, col: 7, offset: 7698},
						val:        "Mende_Kikakui",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 258, col: 7, offset: 7720},
						val:        "Syloti_Nagri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 259, col: 7, offset: 7741},
						val:        "Sora_Sompeng",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 260, col: 7, offset: 7762},
						val:        "Pahawh_Hmong",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 261, col: 7, offset: 7783},
						val:        "Meetei_Mayek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 7, offset: 7804},
						val:        "Join_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 263, col: 7, offset: 7825},
						val:        "Bidi_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 264, col: 7, offset: 7846},
						val:        "White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 265, col: 7, offset: 7866},
						val:        "Warang_Citi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 266, col: 7, offset: 7886},
						val:        "Soft_Dotted",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 267, col: 7, offset: 7906},
						val:        "Pau_Cin_Hau",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 268, col: 7, offset: 7926},
						val:        "Old_Persian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 269, col: 7, offset: 7946},
						val:        "New_Tai_Lue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 270, col: 7, offset: 7966},
						val:        "Ideographic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 271, col: 7, offset: 7986},
						val:        "Saurashtra",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 272, col: 7, offset: 8005},
						val:        "Phoenician",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 273, col: 7, offset: 8024},
						val:        "Other_Math",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 274, col: 7, offset: 8043},
						val:        "Old_Turkic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 275, col: 7, offset: 8062},
						val:        "Old_Permic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 276, col: 7, offset: 8081},
						val:        "Old_Italic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 277, col: 7, offset: 8100},
						val:        "Manichaean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 278, col: 7, offset: 8119},
						val:        "Kharoshthi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 7, offset: 8138},
						val:        "Glagolitic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 280, col: 7, offset: 8157},
						val:        "Devanagari",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 281, col: 7, offset: 8176},
						val:        "Deprecated",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 282, col: 7, offset: 8195},
						val:        "Sundanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 283, col: 7, offset: 8213},
						val:        "Samaritan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 284, col: 7, offset: 8231},
						val:        "Palmyrene",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 7, offset: 8249},
						val:        "Nabataean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 286, col: 7, offset: 8267},
						val:        "Mongolian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 287, col: 7, offset: 8285},
						val:        "Malayalam",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 288, col: 7, offset: 8303},
						val:        "Khudawadi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 289, col: 7, offset: 8321},
						val:        "Inherited",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 290, col: 7, offset: 8339},
						val:        "Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 291, col: 7, offset: 8357},
						val:        "Diacritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 292, col: 7, offset: 8375},
						val:        "Cuneiform",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 293, col: 7, offset: 8393},
						val:        "Bassa_Vah",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 294, col: 7, offset: 8411},
						val:        "Ugaritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 295, col: 7, offset: 8428},
						val:        "Tifinagh",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 296, col: 7, offset: 8445},
						val:        "Tai_Viet",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 297, col: 7, offset: 8462},
						val:        "Tai_Tham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 298, col: 7, offset: 8479},
						val:        "Tagbanwa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 299, col: 7, offset: 8496},
						val:        "Phags_Pa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 300, col: 7, offset: 8513},
						val:        "Ol_Chiki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 301, col: 7, offset: 8530},
						val:        "Mahajani",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 302, col: 7, offset: 8547},
						val:        "Linear_B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 303, col: 7, offset: 8564},
						val:        "Linear_A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 304, col: 7, offset: 8581},
						val:        "Kayah_Li",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 305, col: 7, offset: 8598},
						val:        "Katakana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 306, col: 7, offset: 8615},
						val:        "Javanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 307, col: 7, offset: 8632},
						val:        "Hiragana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 308, col: 7, offset: 8649},
						val:        "Gurmukhi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 309, col: 7, offset: 8666},
						val:        "Gujarati",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 310, col: 7, offset: 8683},
						val:        "Georgian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 311, col: 7, offset: 8700},
						val:        "Extender",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 312, col: 7, offset: 8717},
						val:        "Ethiopic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 313, col: 7, offset: 8734},
						val:        "Duployan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 314, col: 7, offset: 8751},
						val:        "Cyrillic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 315, col: 7, offset: 8768},
						val:        "Cherokee",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 316, col: 7, offset: 8785},
						val:        "Buginese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 317, col: 7, offset: 8802},
						val:        "Bopomofo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 318, col: 7, offset: 8819},
						val:        "Balinese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 319, col: 7, offset: 8836},
						val:        "Armenian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 320, col: 7, offset: 8853},
						val:        "Tirhuta",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 321, col: 7, offset: 8869},
						val:        "Tibetan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 322, col: 7, offset: 8885},
						val:        "Tagalog",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 323, col: 7, offset: 8901},
						val:        "Sinhala",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 324, col: 7, offset: 8917},
						val:        "Siddham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 325, col: 7, offset: 8933},
						val:        "Shavian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 326, col: 7, offset: 8949},
						val:        "Sharada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 327, col: 7, offset: 8965},
						val:        "Radical",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 328, col: 7, offset: 8981},
						val:        "Osmanya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 329, col: 7, offset: 8997},
						val:        "Myanmar",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 330, col: 7, offset: 9013},
						val:        "Mandaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 331, col: 7, offset: 9029},
						val:        "Kannada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 332, col: 7, offset: 9045},
						val:        "Hanunoo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 333, col: 7, offset: 9061},
						val:        "Grantha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 334, col: 7, offset: 9077},
						val:        "Elbasan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 335, col: 7, offset: 9093},
						val:        "Deseret",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 336, col: 7, offset: 9109},
						val:        "Cypriot",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 337, col: 7, offset: 9125},
						val:        "Braille",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 338, col: 7, offset: 9141},
						val:        "Bengali",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 339, col: 7, offset: 9157},
						val:        "Avestan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 340, col: 7, offset: 9173},
						val:        "Thaana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 341, col: 7, offset: 9188},
						val:        "Telugu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 342, col: 7, offset: 9203},
						val:        "Tai_Le",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 343, col: 7, offset: 9218},
						val:        "Syriac",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 344, col: 7, offset: 9233},
						val:        "Rejang",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 345, col: 7, offset: 9248},
						val:        "Lydian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 346, col: 7, offset: 9263},
						val:        "Lycian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 347, col: 7, offset: 9278},
						val:        "Lepcha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 348, col: 7, offset: 9293},
						val:        "Khojki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 349, col: 7, offset: 9308},
						val:        "Kaithi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 350, col: 7, offset: 9323},
						val:        "Hyphen",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 351, col: 7, offset: 9338},
						val:        "Hebrew",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 352, col: 7, offset: 9353},
						val:        "Hangul",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 353, col: 7, offset: 9368},
						val:        "Gothic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 354, col: 7, offset: 9383},
						val:        "Coptic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 355, col: 7, offset: 9398},
						val:        "Common",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 356, col: 7, offset: 9413},
						val:        "Chakma",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 357, col: 7, offset: 9428},
						val:        "Carian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 358, col: 7, offset: 9443},
						val:        "Brahmi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 359, col: 7, offset: 9458},
						val:        "Arabic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 360, col: 7, offset: 9473},
						val:        "Tamil",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 361, col: 7, offset: 9487},
						val:        "Takri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 362, col: 7, offset: 9501},
						val:        "STerm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 363, col: 7, offset: 9515},
						val:        "Runic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 364, col: 7, offset: 9529},
						val:        "Oriya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 365, col: 7, offset: 9543},
						val:        "Ogham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 366, col: 7, offset: 9557},
						val:        "Limbu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 367, col: 7, offset: 9571},
						val:        "Latin",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 368, col: 7, offset: 9585},
						val:        "Khmer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 369, col: 7, offset: 9599},
						val:        "Greek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 370, col: 7, offset: 9613},
						val:        "Buhid",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 371, col: 7, offset: 9627},
						val:        "Batak",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 372, col: 7, offset: 9641},
						val:        "Bamum",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 373, col: 7, offset: 9655},
						val:        "Thai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 374, col: 7, offset: 9668},
						val:        "Modi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 375, col: 7, offset: 9681},
						val:        "Miao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 376, col: 7, offset: 9694},
						val:        "Lisu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 377, col: 7, offset: 9707},
						val:        "Dash",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 378, col: 7, offset: 9720},
						val:        "Cham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 379, col: 7, offset: 9733},
						val:        "Vai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 380, col: 7, offset: 9745},
						val:        "Nko",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 381, col: 7, offset: 9757},
						val:        "Mro",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 382, col: 7, offset: 9769},
						val:        "Lao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 383, col: 7, offset: 9781},
						val:        "Han",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 384, col: 7, offset: 9793},
						val:        "Zs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 385, col: 7, offset: 9804},
						val:        "Zp",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 386, col: 7, offset: 9815},
						val:        "Zl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 387, col: 7, offset: 9826},
						val:        "Yi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 388, col: 7, offset: 9837},
						val:        "So",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 389, col: 7, offset: 9848},
						val:        "Sm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 390, col: 7, offset: 9859},
						val:        "Sk",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 7, offset: 9870},
						val:        "Sc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 392, col: 7, offset: 9881},
						val:        "Ps",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 393, col: 7, offset: 9892},
						val:        "Po",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 394, col: 7, offset: 9903},
						val:        "Pi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 395, col: 7, offset: 9914},
						val:        "Pf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 396, col: 7, offset: 9925},
						val:        "Pe",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 7, offset: 9936},
						val:        "Pd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 398, col: 7, offset: 9947},
						val:        "Pc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 399, col: 7, offset: 9958},
						val:        "No",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 400, col: 7, offset: 9969},
						val:        "Nl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 401, col: 7, offset: 9980},
						val:        "Nd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 402, col: 7, offset: 9991},
						val:        "Mn",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 403, col: 7, offset: 10002},
						val:        "Me",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 404, col: 7, offset: 10013},
						val:        "Mc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 405, col: 7, offset: 10024},
						val:        "Lu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 406, col: 7, offset: 10035},
						val:        "Lt",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 407, col: 7, offset: 10046},
						val:        "Lo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 408, col: 7, offset: 10057},
						val:        "Lm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 409, col: 7, offset: 10068},
						val:        "Ll",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 410, col: 7, offset: 10079},
						val:        "Cs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 411, col: 7, offset: 10090},
						val:        "Co",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 412, col: 7, offset: 10101},
						val:        "Cf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 413, col: 7, offset: 10112},
						val:        "Cc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 414, col: 7, offset: 10123},
						val:        "Z",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 415, col: 7, offset: 10133},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 416, col: 7, offset: 10143},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 417, col: 7, offset: 10153},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 418, col: 7, offset: 10163},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 419, col: 7, offset: 10173},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 420, col: 7, offset: 10183},
						val:        "C",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "ReservedWord",
			pos:  position{line: 422, col: 1, offset: 10188},
			expr: &choiceExpr{
				pos: position{line: 425, col: 2, offset: 10259},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 425, col: 2, offset: 10259},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 425, col: 2, offset: 10259},
								val:        "break",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 425, col: 10, offset: 10267},
								expr: &ruleRefExpr{
									pos:  position{line: 425, col: 11, offset: 10268},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 426, col: 4, offset: 10286},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 426, col: 4, offset: 10286},
								val:        "case",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 426, col: 11, offset: 10293},
								expr: &ruleRefExpr{
									pos:  position{line: 426, col: 12, offset: 10294},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 427, col: 4, offset: 10312},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 427, col: 4, offset: 10312},
								val:        "chan",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 427, col: 11, offset: 10319},
								expr: &ruleRefExpr{
									pos:  position{line: 427, col: 12, offset: 10320},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 428, col: 4, offset: 10338},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 428, col: 4, offset: 10338},
								val:        "const",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 428, col: 12, offset: 10346},
								expr: &ruleRefExpr{
									pos:  position{line: 428, col: 13, offset: 10347},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 429, col: 4, offset: 10365},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 429, col: 4, offset: 10365},
								val:        "continue",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 429, col: 15, offset: 10376},
								expr: &ruleRefExpr{
									pos:  position{line: 429, col: 16, offset: 10377},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 430, col: 4, offset: 10395},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 430, col: 4, offset: 10395},
								val:        "default",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 430, col: 14, offset: 10405},
								expr: &ruleRefExpr{
									pos:  position{line: 430, col: 15, offset: 10406},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 431, col: 4, offset: 10424},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 431, col: 4, offset: 10424},
								val:        "defer",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 431, col: 12, offset: 10432},
								expr: &ruleRefExpr{
									pos:  position{line: 431, col: 13, offset: 10433},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 432, col: 4, offset: 10451},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 432, col: 4, offset: 10451},
								val:        "else",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 432, col: 11, offset: 10458},
								expr: &ruleRefExpr{
									pos:  position{line: 432, col: 12, offset: 10459},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 433, col: 4, offset: 10477},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 433, col: 4, offset: 10477},
								val:        "fallthrough",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 433, col: 18, offset: 10491},
								expr: &ruleRefExpr{
									pos:  position{line: 433, col: 19, offset: 10492},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 434, col: 4, offset: 10510},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 434, col: 4, offset: 10510},
								val:        "for",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 434, col: 10, offset: 10516},
								expr: &ruleRefExpr{
									pos:  position{line: 434, col: 11, offset: 10517},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 435, col: 4, offset: 10535},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 435, col: 4, offset: 10535},
								val:        "func",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 435, col: 11, offset: 10542},
								expr: &ruleRefExpr{
									pos:  position{line: 435, col: 12, offset: 10543},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 436, col: 4, offset: 10561},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 436, col: 4, offset: 10561},
								val:        "goto",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 436, col: 11, offset: 10568},
								expr: &ruleRefExpr{
									pos:  position{line: 436, col: 12, offset: 10569},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 437, col: 4, offset: 10587},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 437, col: 4, offset: 10587},
								val:        "go",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 437, col: 9, offset: 10592},
								expr: &ruleRefExpr{
									pos:  position{line: 437, col: 10, offset: 10593},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 438, col: 4, offset: 10611},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 438, col: 4, offset: 10611},
								val:        "if",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 438, col: 9, offset: 10616},
								expr: &ruleRefExpr{
									pos:  position{line: 438, col: 10, offset: 10617},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 439, col: 4, offset: 10635},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 439, col: 4, offset: 10635},
								val:        "import",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 439, col: 13, offset: 10644},
								expr: &ruleRefExpr{
									pos:  position{line: 439, col: 14, offset: 10645},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 440, col: 4, offset: 10663},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 440, col: 4, offset: 10663},
								val:        "interface",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 440, col: 16, offset: 10675},
								expr: &ruleRefExpr{
									pos:  position{line: 440, col: 17, offset: 10676},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 441, col: 4, offset: 10694},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 441, col: 4, offset: 10694},
								val:        "map",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 441, col: 10, offset: 10700},
								expr: &ruleRefExpr{
									pos:  position{line: 441, col: 11, offset: 10701},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 442, col: 4, offset: 10719},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 442, col: 4, offset: 10719},
								val:        "package",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 442, col: 14, offset: 10729},
								expr: &ruleRefExpr{
									pos:  position{line: 442, col: 15, offset: 10730},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 443, col: 4, offset: 10748},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 443, col: 4, offset: 10748},
								val:        "range",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 443, col: 12, offset: 10756},
								expr: &ruleRefExpr{
									pos:  position{line: 443, col: 13, offset: 10757},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 444, col: 4, offset: 10775},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 444, col: 4, offset: 10775},
								val:        "return",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 444, col: 13, offset: 10784},
								expr: &ruleRefExpr{
									pos:  position{line: 444, col: 14, offset: 10785},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 445, col: 4, offset: 10803},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 445, col: 4, offset: 10803},
								val:        "select",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 445, col: 13, offset: 10812},
								expr: &ruleRefExpr{
									pos:  position{line: 445, col: 14, offset: 10813},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 446, col: 4, offset: 10831},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 446, col: 4, offset: 10831},
								val:        "struct",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 446, col: 13, offset: 10840},
								expr: &ruleRefExpr{
									pos:  position{line: 446, col: 14, offset: 10841},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 447, col: 4, offset: 10859},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 447, col: 4, offset: 10859},
								val:        "switch",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 447, col: 13, offset: 10868},
								expr: &ruleRefExpr{
									pos:  position{line: 447, col: 14, offset: 10869},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 448, col: 4, offset: 10887},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 448, col: 4, offset: 10887},
								val:        "type",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 448, col: 11, offset: 10894},
								expr: &ruleRefExpr{
									pos:  position{line: 448, col: 12, offset: 10895},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 449, col: 4, offset: 10913},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 449, col: 4, offset: 10913},
								val:        "var",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 449, col: 10, offset: 10919},
								expr: &ruleRefExpr{
									pos:  position{line: 449, col: 11, offset: 10920},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 453, col: 4, offset: 11019},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 453, col: 4, offset: 11019},
								val:        "bool",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 453, col: 11, offset: 11026},
								expr: &ruleRefExpr{
									pos:  position{line: 453, col: 12, offset: 11027},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 454, col: 4, offset: 11045},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 454, col: 4, offset: 11045},
								val:        "byte",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 454, col: 11, offset: 11052},
								expr: &ruleRefExpr{
									pos:  position{line: 454, col: 12, offset: 11053},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 455, col: 4, offset: 11071},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 455, col: 4, offset: 11071},
								val:        "complex64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 455, col: 16, offset: 11083},
								expr: &ruleRefExpr{
									pos:  position{line: 455, col: 17, offset: 11084},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 456, col: 4, offset: 11102},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 456, col: 4, offset: 11102},
								val:        "complex128",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 456, col: 17, offset: 11115},
								expr: &ruleRefExpr{
									pos:  position{line: 456, col: 18, offset: 11116},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 457, col: 4, offset: 11134},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 457, col: 4, offset: 11134},
								val:        "error",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 457, col: 12, offset: 11142},
								expr: &ruleRefExpr{
									pos:  position{line: 457, col: 13, offset: 11143},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 458, col: 4, offset: 11161},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 458, col: 4, offset: 11161},
								val:        "float32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 458, col: 14, offset: 11171},
								expr: &ruleRefExpr{
									pos:  position{line: 458, col: 15, offset: 11172},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 459, col: 4, offset: 11190},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 459, col: 4, offset: 11190},
								val:        "float64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 459, col: 14, offset: 11200},
								expr: &ruleRefExpr{
									pos:  position{line: 459, col: 15, offset: 11201},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 460, col: 4, offset: 11219},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 460, col: 4, offset: 11219},
								val:        "int8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 460, col: 11, offset: 11226},
								expr: &ruleRefExpr{
									pos:  position{line: 460, col: 12, offset: 11227},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 461, col: 4, offset: 11245},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 461, col: 4, offset: 11245},
								val:        "int16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 461, col: 12, offset: 11253},
								expr: &ruleRefExpr{
									pos:  position{line: 461, col: 13, offset: 11254},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 462, col: 4, offset: 11272},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 462, col: 4, offset: 11272},
								val:        "int32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 462, col: 12, offset: 11280},
								expr: &ruleRefExpr{
									pos:  position{line: 462, col: 13, offset: 11281},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 463, col: 4, offset: 11299},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 463, col: 4, offset: 11299},
								val:        "int64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 463, col: 12, offset: 11307},
								expr: &ruleRefExpr{
									pos:  position{line: 463, col: 13, offset: 11308},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 464, col: 4, offset: 11326},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 464, col: 4, offset: 11326},
								val:        "int",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 464, col: 10, offset: 11332},
								expr: &ruleRefExpr{
									pos:  position{line: 464, col: 11, offset: 11333},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 465, col: 4, offset: 11351},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 465, col: 4, offset: 11351},
								val:        "rune",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 465, col: 11, offset: 11358},
								expr: &ruleRefExpr{
									pos:  position{line: 465, col: 12, offset: 11359},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 466, col: 4, offset: 11377},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 466, col: 4, offset: 11377},
								val:        "string",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 466, col: 13, offset: 11386},
								expr: &ruleRefExpr{
									pos:  position{line: 466, col: 14, offset: 11387},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 467, col: 4, offset: 11405},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 467, col: 4, offset: 11405},
								val:        "uint8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 467, col: 12, offset: 11413},
								expr: &ruleRefExpr{
									pos:  position{line: 467, col: 13, offset: 11414},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 468, col: 4, offset: 11432},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 468, col: 4, offset: 11432},
								val:        "uint16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 468, col: 13, offset: 11441},
								expr: &ruleRefExpr{
									pos:  position{line: 468, col: 14, offset: 11442},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 469, col: 4, offset: 11460},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 469, col: 4, offset: 11460},
								val:        "uint32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 469, col: 13, offset: 11469},
								expr: &ruleRefExpr{
									pos:  position{line: 469, col: 14, offset: 11470},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 470, col: 4, offset: 11488},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 470, col: 4, offset: 11488},
								val:        "uint64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 470, col: 13, offset: 11497},
								expr: &ruleRefExpr{
									pos:  position{line: 470, col: 14, offset: 11498},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 471, col: 4, offset: 11516},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 471, col: 4, offset: 11516},
								val:        "uintptr",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 471, col: 14, offset: 11526},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 15, offset: 11527},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 472, col: 4, offset: 11545},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 472, col: 4, offset: 11545},
								val:        "uint",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 472, col: 11, offset: 11552},
								expr: &ruleRefExpr{
									pos:  position{line: 472, col: 12, offset: 11553},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 473, col: 4, offset: 11571},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 473, col: 4, offset: 11571},
								val:        "true",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 473, col: 11, offset: 11578},
								expr: &ruleRefExpr{
									pos:  position{line: 473, col: 12, offset: 11579},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 474, col: 4, offset: 11597},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 474, col: 4, offset: 11597},
								val:        "false",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 474, col: 12, offset: 11605},
								expr: &ruleRefExpr{
									pos:  position{line: 474, col: 13, offset: 11606},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 475, col: 4, offset: 11624},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 475, col: 4, offset: 11624},
								val:        "iota",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 475, col: 11, offset: 11631},
								expr: &ruleRefExpr{
									pos:  position{line: 475, col: 12, offset: 11632},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 476, col: 4, offset: 11650},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 476, col: 4, offset: 11650},
								val:        "nil",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 476, col: 10, offset: 11656},
								expr: &ruleRefExpr{
									pos:  position{line: 476, col: 11, offset: 11657},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 477, col: 4, offset: 11675},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 477, col: 4, offset: 11675},
								val:        "append",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 477, col: 13, offset: 11684},
								expr: &ruleRefExpr{
									pos:  position{line: 477, col: 14, offset: 11685},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 478, col: 4, offset: 11703},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 478, col: 4, offset: 11703},
								val:        "cap",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 478, col: 10, offset: 11709},
								expr: &ruleRefExpr{
									pos:  position{line: 478, col: 11, offset: 11710},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 479, col: 4, offset: 11728},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 479, col: 4, offset: 11728},
								val:        "close",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 479, col: 12, offset: 11736},
								expr: &ruleRefExpr{
									pos:  position{line: 479, col: 13, offset: 11737},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 480, col: 4, offset: 11755},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 480, col: 4, offset: 11755},
								val:        "complex",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 480, col: 14, offset: 11765},
								expr: &ruleRefExpr{
									pos:  position{line: 480, col: 15, offset: 11766},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 481, col: 4, offset: 11784},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 481, col: 4, offset: 11784},
								val:        "copy",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 481, col: 11, offset: 11791},
								expr: &ruleRefExpr{
									pos:  position{line: 481, col: 12, offset: 11792},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 482, col: 4, offset: 11810},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 482, col: 4, offset: 11810},
								val:        "delete",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 482, col: 13, offset: 11819},
								expr: &ruleRefExpr{
									pos:  position{line: 482, col: 14, offset: 11820},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 483, col: 4, offset: 11838},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 483, col: 4, offset: 11838},
								val:        "imag",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 483, col: 11, offset: 11845},
								expr: &ruleRefExpr{
									pos:  position{line: 483, col: 12, offset: 11846},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 484, col: 4, offset: 11864},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 484, col: 4, offset: 11864},
								val:        "len",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 484, col: 10, offset: 11870},
								expr: &ruleRefExpr{
									pos:  position{line: 484, col: 11, offset: 11871},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 485, col: 4, offset: 11889},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 485, col: 4, offset: 11889},
								val:        "make",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 485, col: 11, offset: 11896},
								expr: &ruleRefExpr{
									pos:  position{line: 485, col: 12, offset: 11897},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 486, col: 4, offset: 11915},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 486, col: 4, offset: 11915},
								val:        "new",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 486, col: 10, offset: 11921},
								expr: &ruleRefExpr{
									pos:  position{line: 486, col: 11, offset: 11922},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 487, col: 4, offset: 11940},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 487, col: 4, offset: 11940},
								val:        "panic",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 487, col: 12, offset: 11948},
								expr: &ruleRefExpr{
									pos:  position{line: 487, col: 13, offset: 11949},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 488, col: 4, offset: 11967},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 488, col: 4, offset: 11967},
								val:        "println",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 488, col: 14, offset: 11977},
								expr: &ruleRefExpr{
									pos:  position{line: 488, col: 15, offset: 11978},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 489, col: 4, offset: 11996},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 489, col: 4, offset: 11996},
								val:        "print",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 489, col: 12, offset: 12004},
								expr: &ruleRefExpr{
									pos:  position{line: 489, col: 13, offset: 12005},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 490, col: 4, offset: 12023},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 490, col: 4, offset: 12023},
								val:        "real",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 490, col: 11, offset: 12030},
								expr: &ruleRefExpr{
									pos:  position{line: 490, col: 12, offset: 12031},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 491, col: 4, offset: 12049},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 491, col: 4, offset: 12049},
								val:        "recover",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 491, col: 14, offset: 12059},
								expr: &ruleRefExpr{
									pos:  position{line: 491, col: 15, offset: 12060},
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
			pos:  position{line: 493, col: 1, offset: 12076},
			expr: &actionExpr{
				pos: position{line: 493, col: 14, offset: 12091},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 493, col: 14, offset: 12091},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 498, col: 1, offset: 12166},
			expr: &actionExpr{
				pos: position{line: 498, col: 13, offset: 12180},
				run: (*parser).callonCodeBlock1,
				expr: &seqExpr{
					pos: position{line: 498, col: 13, offset: 12180},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 498, col: 13, offset: 12180},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 498, col: 17, offset: 12184},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 498, col: 22, offset: 12189},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 504, col: 1, offset: 12287},
			expr: &zeroOrMoreExpr{
				pos: position{line: 504, col: 8, offset: 12296},
				expr: &choiceExpr{
					pos: position{line: 504, col: 10, offset: 12298},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 504, col: 10, offset: 12298},
							expr: &seqExpr{
								pos: position{line: 504, col: 12, offset: 12300},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 504, col: 12, offset: 12300},
										expr: &charClassMatcher{
											pos:        position{line: 504, col: 13, offset: 12301},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 504, col: 18, offset: 12306},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 504, col: 34, offset: 12322},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 504, col: 34, offset: 12322},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 504, col: 38, offset: 12326},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 504, col: 43, offset: 12331},
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
			pos:  position{line: 506, col: 1, offset: 12339},
			expr: &zeroOrMoreExpr{
				pos: position{line: 506, col: 6, offset: 12346},
				expr: &choiceExpr{
					pos: position{line: 506, col: 8, offset: 12348},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 506, col: 8, offset: 12348},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 506, col: 21, offset: 12361},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 506, col: 27, offset: 12367},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 507, col: 1, offset: 12378},
			expr: &zeroOrMoreExpr{
				pos: position{line: 507, col: 5, offset: 12384},
				expr: &choiceExpr{
					pos: position{line: 507, col: 7, offset: 12386},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 507, col: 7, offset: 12386},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 507, col: 20, offset: 12399},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 509, col: 1, offset: 12436},
			expr: &charClassMatcher{
				pos:        position{line: 509, col: 14, offset: 12451},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 510, col: 1, offset: 12459},
			expr: &litMatcher{
				pos:        position{line: 510, col: 7, offset: 12467},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 511, col: 1, offset: 12472},
			expr: &choiceExpr{
				pos: position{line: 511, col: 7, offset: 12480},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 511, col: 7, offset: 12480},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 511, col: 7, offset: 12480},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 511, col: 10, offset: 12483},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 511, col: 16, offset: 12489},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 511, col: 16, offset: 12489},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 511, col: 18, offset: 12491},
								expr: &ruleRefExpr{
									pos:  position{line: 511, col: 18, offset: 12491},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 511, col: 37, offset: 12510},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 511, col: 43, offset: 12516},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 511, col: 43, offset: 12516},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 511, col: 46, offset: 12519},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 513, col: 1, offset: 12524},
			expr: &notExpr{
				pos: position{line: 513, col: 7, offset: 12532},
				expr: &anyMatcher{
					line: 513, col: 8, offset: 12533,
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

func (c *current) onStringLiteral2() (interface{}, error) {
	return ast.NewStringLit(c.astPos(), string(c.text)), nil
}

func (p *parser) callonStringLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral2()
}

func (c *current) onStringLiteral18() (interface{}, error) {
	return ast.NewStringLit(c.astPos(), "``"), errors.New("string literal not terminated")
}

func (p *parser) callonStringLiteral18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral18()
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

func (c *current) onUnicodeClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid Unicode class escape")
}

func (p *parser) callonUnicodeClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape5()
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

	depth  int
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
	val, ok := p.parseExpr(rule.expr)
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok {
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
