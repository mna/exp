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
			pos:  position{line: 190, col: 1, offset: 5271},
			expr: &choiceExpr{
				pos: position{line: 190, col: 17, offset: 5289},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 190, col: 17, offset: 5289},
						run: (*parser).callonStringLiteral2,
						expr: &choiceExpr{
							pos: position{line: 190, col: 19, offset: 5291},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 190, col: 19, offset: 5291},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 190, col: 19, offset: 5291},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 190, col: 23, offset: 5295},
											expr: &ruleRefExpr{
												pos:  position{line: 190, col: 23, offset: 5295},
												name: "DoubleStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 190, col: 41, offset: 5313},
											val:        "\"",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 190, col: 47, offset: 5319},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 190, col: 47, offset: 5319},
											val:        "'",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 190, col: 51, offset: 5323},
											name: "SingleStringChar",
										},
										&litMatcher{
											pos:        position{line: 190, col: 68, offset: 5340},
											val:        "'",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 190, col: 74, offset: 5346},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 190, col: 74, offset: 5346},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 190, col: 78, offset: 5350},
											expr: &ruleRefExpr{
												pos:  position{line: 190, col: 78, offset: 5350},
												name: "RawStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 190, col: 93, offset: 5365},
											val:        "`",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 192, col: 5, offset: 5438},
						run: (*parser).callonStringLiteral18,
						expr: &choiceExpr{
							pos: position{line: 192, col: 7, offset: 5440},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 192, col: 9, offset: 5442},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 192, col: 9, offset: 5442},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 192, col: 13, offset: 5446},
											expr: &ruleRefExpr{
												pos:  position{line: 192, col: 13, offset: 5446},
												name: "DoubleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 192, col: 33, offset: 5466},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 192, col: 33, offset: 5466},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 192, col: 39, offset: 5472},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 192, col: 51, offset: 5484},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 192, col: 51, offset: 5484},
											val:        "'",
											ignoreCase: false,
										},
										&zeroOrOneExpr{
											pos: position{line: 192, col: 55, offset: 5488},
											expr: &ruleRefExpr{
												pos:  position{line: 192, col: 55, offset: 5488},
												name: "SingleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 192, col: 75, offset: 5508},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 192, col: 75, offset: 5508},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 192, col: 81, offset: 5514},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 192, col: 91, offset: 5524},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 192, col: 91, offset: 5524},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 192, col: 95, offset: 5528},
											expr: &ruleRefExpr{
												pos:  position{line: 192, col: 95, offset: 5528},
												name: "RawStringChar",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 192, col: 110, offset: 5543},
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
			pos:  position{line: 196, col: 1, offset: 5645},
			expr: &choiceExpr{
				pos: position{line: 196, col: 20, offset: 5666},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 196, col: 20, offset: 5666},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 196, col: 20, offset: 5666},
								expr: &choiceExpr{
									pos: position{line: 196, col: 23, offset: 5669},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 196, col: 23, offset: 5669},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 196, col: 29, offset: 5675},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 196, col: 36, offset: 5682},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 196, col: 42, offset: 5688},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 196, col: 55, offset: 5701},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 196, col: 55, offset: 5701},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 196, col: 60, offset: 5706},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 197, col: 1, offset: 5725},
			expr: &choiceExpr{
				pos: position{line: 197, col: 20, offset: 5746},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 197, col: 20, offset: 5746},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 197, col: 20, offset: 5746},
								expr: &choiceExpr{
									pos: position{line: 197, col: 23, offset: 5749},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 197, col: 23, offset: 5749},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 197, col: 29, offset: 5755},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 197, col: 36, offset: 5762},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 197, col: 42, offset: 5768},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 197, col: 55, offset: 5781},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 197, col: 55, offset: 5781},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 197, col: 60, offset: 5786},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 198, col: 1, offset: 5805},
			expr: &seqExpr{
				pos: position{line: 198, col: 17, offset: 5823},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 198, col: 17, offset: 5823},
						expr: &litMatcher{
							pos:        position{line: 198, col: 18, offset: 5824},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 22, offset: 5828},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 200, col: 1, offset: 5840},
			expr: &choiceExpr{
				pos: position{line: 200, col: 22, offset: 5863},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 200, col: 24, offset: 5865},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 200, col: 24, offset: 5865},
								val:        "\"",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 200, col: 30, offset: 5871},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 201, col: 7, offset: 5900},
						run: (*parser).callonDoubleStringEscape5,
						expr: &ruleRefExpr{
							pos:  position{line: 201, col: 7, offset: 5900},
							name: "SourceChar",
						},
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 204, col: 1, offset: 5970},
			expr: &choiceExpr{
				pos: position{line: 204, col: 22, offset: 5993},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 204, col: 24, offset: 5995},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 204, col: 24, offset: 5995},
								val:        "'",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 30, offset: 6001},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 205, col: 7, offset: 6030},
						run: (*parser).callonSingleStringEscape5,
						expr: &ruleRefExpr{
							pos:  position{line: 205, col: 7, offset: 6030},
							name: "SourceChar",
						},
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 209, col: 1, offset: 6101},
			expr: &choiceExpr{
				pos: position{line: 209, col: 24, offset: 6126},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 209, col: 24, offset: 6126},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 43, offset: 6145},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 57, offset: 6159},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 69, offset: 6171},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 89, offset: 6191},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 210, col: 1, offset: 6210},
			expr: &choiceExpr{
				pos: position{line: 210, col: 20, offset: 6231},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 210, col: 20, offset: 6231},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 26, offset: 6237},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 32, offset: 6243},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 38, offset: 6249},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 44, offset: 6255},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 50, offset: 6261},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 56, offset: 6267},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 62, offset: 6273},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 211, col: 1, offset: 6278},
			expr: &choiceExpr{
				pos: position{line: 211, col: 15, offset: 6294},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 211, col: 15, offset: 6294},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 211, col: 15, offset: 6294},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 26, offset: 6305},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 37, offset: 6316},
								name: "OctalDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 212, col: 7, offset: 6333},
						run: (*parser).callonOctalEscape6,
						expr: &seqExpr{
							pos: position{line: 212, col: 7, offset: 6333},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 212, col: 7, offset: 6333},
									name: "OctalDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 212, col: 18, offset: 6344},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 215, col: 1, offset: 6410},
			expr: &choiceExpr{
				pos: position{line: 215, col: 13, offset: 6424},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 215, col: 13, offset: 6424},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 215, col: 13, offset: 6424},
								val:        "x",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 17, offset: 6428},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 26, offset: 6437},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 216, col: 7, offset: 6452},
						run: (*parser).callonHexEscape6,
						expr: &seqExpr{
							pos: position{line: 216, col: 7, offset: 6452},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 216, col: 7, offset: 6452},
									val:        "x",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 216, col: 11, offset: 6456},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 219, col: 1, offset: 6528},
			expr: &choiceExpr{
				pos: position{line: 219, col: 21, offset: 6550},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 219, col: 21, offset: 6550},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 219, col: 21, offset: 6550},
								val:        "U",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 25, offset: 6554},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 34, offset: 6563},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 43, offset: 6572},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 52, offset: 6581},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 61, offset: 6590},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 70, offset: 6599},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 79, offset: 6608},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 88, offset: 6617},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 220, col: 7, offset: 6632},
						run: (*parser).callonLongUnicodeEscape12,
						expr: &seqExpr{
							pos: position{line: 220, col: 7, offset: 6632},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 220, col: 7, offset: 6632},
									val:        "U",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 220, col: 11, offset: 6636},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 223, col: 1, offset: 6704},
			expr: &choiceExpr{
				pos: position{line: 223, col: 22, offset: 6727},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 223, col: 22, offset: 6727},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 223, col: 22, offset: 6727},
								val:        "u",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 26, offset: 6731},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 35, offset: 6740},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 44, offset: 6749},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 53, offset: 6758},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 224, col: 7, offset: 6773},
						run: (*parser).callonShortUnicodeEscape8,
						expr: &seqExpr{
							pos: position{line: 224, col: 7, offset: 6773},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 224, col: 7, offset: 6773},
									val:        "u",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 224, col: 11, offset: 6777},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 228, col: 1, offset: 6846},
			expr: &charClassMatcher{
				pos:        position{line: 228, col: 14, offset: 6861},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 229, col: 1, offset: 6867},
			expr: &charClassMatcher{
				pos:        position{line: 229, col: 16, offset: 6884},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 230, col: 1, offset: 6890},
			expr: &charClassMatcher{
				pos:        position{line: 230, col: 12, offset: 6903},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 232, col: 1, offset: 6914},
			expr: &choiceExpr{
				pos: position{line: 232, col: 20, offset: 6935},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 232, col: 20, offset: 6935},
						run: (*parser).callonCharClassMatcher2,
						expr: &seqExpr{
							pos: position{line: 232, col: 20, offset: 6935},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 232, col: 20, offset: 6935},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 232, col: 24, offset: 6939},
									expr: &choiceExpr{
										pos: position{line: 232, col: 26, offset: 6941},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 232, col: 26, offset: 6941},
												name: "ClassCharRange",
											},
											&ruleRefExpr{
												pos:  position{line: 232, col: 43, offset: 6958},
												name: "ClassChar",
											},
											&seqExpr{
												pos: position{line: 232, col: 55, offset: 6970},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 232, col: 55, offset: 6970},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 232, col: 60, offset: 6975},
														name: "UnicodeClassEscape",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 232, col: 82, offset: 6997},
									val:        "]",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 232, col: 86, offset: 7001},
									expr: &litMatcher{
										pos:        position{line: 232, col: 86, offset: 7001},
										val:        "i",
										ignoreCase: false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 236, col: 5, offset: 7108},
						run: (*parser).callonCharClassMatcher15,
						expr: &seqExpr{
							pos: position{line: 236, col: 5, offset: 7108},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 236, col: 5, offset: 7108},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 236, col: 9, offset: 7112},
									expr: &seqExpr{
										pos: position{line: 236, col: 11, offset: 7114},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 236, col: 11, offset: 7114},
												expr: &ruleRefExpr{
													pos:  position{line: 236, col: 14, offset: 7117},
													name: "EOL",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 236, col: 20, offset: 7123},
												name: "SourceChar",
											},
										},
									},
								},
								&choiceExpr{
									pos: position{line: 236, col: 36, offset: 7139},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 236, col: 36, offset: 7139},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 236, col: 42, offset: 7145},
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
			name: "ClassCharRange",
			pos:  position{line: 240, col: 1, offset: 7255},
			expr: &seqExpr{
				pos: position{line: 240, col: 18, offset: 7274},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 240, col: 18, offset: 7274},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 240, col: 28, offset: 7284},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 240, col: 32, offset: 7288},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 241, col: 1, offset: 7298},
			expr: &choiceExpr{
				pos: position{line: 241, col: 13, offset: 7312},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 241, col: 13, offset: 7312},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 241, col: 13, offset: 7312},
								expr: &choiceExpr{
									pos: position{line: 241, col: 16, offset: 7315},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 241, col: 16, offset: 7315},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 241, col: 22, offset: 7321},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 241, col: 29, offset: 7328},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 35, offset: 7334},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 241, col: 48, offset: 7347},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 241, col: 48, offset: 7347},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 53, offset: 7352},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 242, col: 1, offset: 7368},
			expr: &choiceExpr{
				pos: position{line: 242, col: 19, offset: 7388},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 242, col: 21, offset: 7390},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 242, col: 21, offset: 7390},
								val:        "]",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 242, col: 27, offset: 7396},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 243, col: 7, offset: 7425},
						run: (*parser).callonCharClassEscape5,
						expr: &seqExpr{
							pos: position{line: 243, col: 7, offset: 7425},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 243, col: 7, offset: 7425},
									expr: &litMatcher{
										pos:        position{line: 243, col: 8, offset: 7426},
										val:        "p",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 243, col: 12, offset: 7430},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 247, col: 1, offset: 7501},
			expr: &seqExpr{
				pos: position{line: 247, col: 22, offset: 7524},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 247, col: 22, offset: 7524},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 248, col: 7, offset: 7537},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 248, col: 7, offset: 7537},
								name: "SingleCharUnicodeClass",
							},
							&actionExpr{
								pos: position{line: 249, col: 7, offset: 7566},
								run: (*parser).callonUnicodeClassEscape5,
								expr: &seqExpr{
									pos: position{line: 249, col: 7, offset: 7566},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 249, col: 7, offset: 7566},
											expr: &choiceExpr{
												pos: position{line: 249, col: 10, offset: 7569},
												alternatives: []interface{}{
													&litMatcher{
														pos:        position{line: 249, col: 10, offset: 7569},
														val:        "{",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 249, col: 16, offset: 7575},
														name: "EOL",
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 249, col: 22, offset: 7581},
											name: "SourceChar",
										},
									},
								},
							},
							&seqExpr{
								pos: position{line: 250, col: 7, offset: 7657},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 250, col: 7, offset: 7657},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 250, col: 11, offset: 7661},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 250, col: 24, offset: 7674},
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
			pos:  position{line: 252, col: 1, offset: 7684},
			expr: &charClassMatcher{
				pos:        position{line: 252, col: 26, offset: 7711},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 255, col: 1, offset: 7784},
			expr: &choiceExpr{
				pos: position{line: 255, col: 16, offset: 7801},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 255, col: 16, offset: 7801},
						val:        "Other_Default_Ignorable_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 7, offset: 7844},
						val:        "Noncharacter_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 257, col: 7, offset: 7876},
						val:        "Logical_Order_Exception",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 258, col: 7, offset: 7908},
						val:        "Inscriptional_Parthian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 259, col: 7, offset: 7939},
						val:        "Other_Grapheme_Extend",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 260, col: 7, offset: 7969},
						val:        "Inscriptional_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 261, col: 7, offset: 7999},
						val:        "Terminal_Punctuation",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 7, offset: 8028},
						val:        "Meroitic_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 263, col: 7, offset: 8057},
						val:        "IDS_Trinary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 264, col: 7, offset: 8086},
						val:        "Egyptian_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 265, col: 7, offset: 8115},
						val:        "Pattern_White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 266, col: 7, offset: 8143},
						val:        "IDS_Binary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 267, col: 7, offset: 8171},
						val:        "Canadian_Aboriginal",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 268, col: 7, offset: 8199},
						val:        "Variation_Selector",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 269, col: 7, offset: 8226},
						val:        "Caucasian_Albanian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 270, col: 7, offset: 8253},
						val:        "Unified_Ideograph",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 271, col: 7, offset: 8279},
						val:        "Other_ID_Continue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 272, col: 7, offset: 8305},
						val:        "Old_South_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 273, col: 7, offset: 8331},
						val:        "Old_North_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 274, col: 7, offset: 8357},
						val:        "Other_Alphabetic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 275, col: 7, offset: 8382},
						val:        "Meroitic_Cursive",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 276, col: 7, offset: 8407},
						val:        "Imperial_Aramaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 277, col: 7, offset: 8432},
						val:        "Psalter_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 278, col: 7, offset: 8456},
						val:        "Other_Uppercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 7, offset: 8480},
						val:        "Other_Lowercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 280, col: 7, offset: 8504},
						val:        "ASCII_Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 281, col: 7, offset: 8528},
						val:        "Quotation_Mark",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 282, col: 7, offset: 8551},
						val:        "Pattern_Syntax",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 283, col: 7, offset: 8574},
						val:        "Other_ID_Start",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 284, col: 7, offset: 8597},
						val:        "Mende_Kikakui",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 7, offset: 8619},
						val:        "Syloti_Nagri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 286, col: 7, offset: 8640},
						val:        "Sora_Sompeng",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 287, col: 7, offset: 8661},
						val:        "Pahawh_Hmong",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 288, col: 7, offset: 8682},
						val:        "Meetei_Mayek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 289, col: 7, offset: 8703},
						val:        "Join_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 290, col: 7, offset: 8724},
						val:        "Bidi_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 291, col: 7, offset: 8745},
						val:        "White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 292, col: 7, offset: 8765},
						val:        "Warang_Citi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 293, col: 7, offset: 8785},
						val:        "Soft_Dotted",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 294, col: 7, offset: 8805},
						val:        "Pau_Cin_Hau",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 295, col: 7, offset: 8825},
						val:        "Old_Persian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 296, col: 7, offset: 8845},
						val:        "New_Tai_Lue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 297, col: 7, offset: 8865},
						val:        "Ideographic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 298, col: 7, offset: 8885},
						val:        "Saurashtra",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 299, col: 7, offset: 8904},
						val:        "Phoenician",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 300, col: 7, offset: 8923},
						val:        "Other_Math",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 301, col: 7, offset: 8942},
						val:        "Old_Turkic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 302, col: 7, offset: 8961},
						val:        "Old_Permic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 303, col: 7, offset: 8980},
						val:        "Old_Italic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 304, col: 7, offset: 8999},
						val:        "Manichaean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 305, col: 7, offset: 9018},
						val:        "Kharoshthi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 306, col: 7, offset: 9037},
						val:        "Glagolitic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 307, col: 7, offset: 9056},
						val:        "Devanagari",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 308, col: 7, offset: 9075},
						val:        "Deprecated",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 309, col: 7, offset: 9094},
						val:        "Sundanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 310, col: 7, offset: 9112},
						val:        "Samaritan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 311, col: 7, offset: 9130},
						val:        "Palmyrene",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 312, col: 7, offset: 9148},
						val:        "Nabataean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 313, col: 7, offset: 9166},
						val:        "Mongolian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 314, col: 7, offset: 9184},
						val:        "Malayalam",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 315, col: 7, offset: 9202},
						val:        "Khudawadi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 316, col: 7, offset: 9220},
						val:        "Inherited",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 317, col: 7, offset: 9238},
						val:        "Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 318, col: 7, offset: 9256},
						val:        "Diacritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 319, col: 7, offset: 9274},
						val:        "Cuneiform",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 320, col: 7, offset: 9292},
						val:        "Bassa_Vah",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 321, col: 7, offset: 9310},
						val:        "Ugaritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 322, col: 7, offset: 9327},
						val:        "Tifinagh",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 323, col: 7, offset: 9344},
						val:        "Tai_Viet",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 324, col: 7, offset: 9361},
						val:        "Tai_Tham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 325, col: 7, offset: 9378},
						val:        "Tagbanwa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 326, col: 7, offset: 9395},
						val:        "Phags_Pa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 327, col: 7, offset: 9412},
						val:        "Ol_Chiki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 328, col: 7, offset: 9429},
						val:        "Mahajani",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 329, col: 7, offset: 9446},
						val:        "Linear_B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 330, col: 7, offset: 9463},
						val:        "Linear_A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 331, col: 7, offset: 9480},
						val:        "Kayah_Li",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 332, col: 7, offset: 9497},
						val:        "Katakana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 333, col: 7, offset: 9514},
						val:        "Javanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 334, col: 7, offset: 9531},
						val:        "Hiragana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 335, col: 7, offset: 9548},
						val:        "Gurmukhi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 336, col: 7, offset: 9565},
						val:        "Gujarati",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 337, col: 7, offset: 9582},
						val:        "Georgian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 338, col: 7, offset: 9599},
						val:        "Extender",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 339, col: 7, offset: 9616},
						val:        "Ethiopic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 340, col: 7, offset: 9633},
						val:        "Duployan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 341, col: 7, offset: 9650},
						val:        "Cyrillic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 342, col: 7, offset: 9667},
						val:        "Cherokee",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 343, col: 7, offset: 9684},
						val:        "Buginese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 344, col: 7, offset: 9701},
						val:        "Bopomofo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 345, col: 7, offset: 9718},
						val:        "Balinese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 346, col: 7, offset: 9735},
						val:        "Armenian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 347, col: 7, offset: 9752},
						val:        "Tirhuta",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 348, col: 7, offset: 9768},
						val:        "Tibetan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 349, col: 7, offset: 9784},
						val:        "Tagalog",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 350, col: 7, offset: 9800},
						val:        "Sinhala",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 351, col: 7, offset: 9816},
						val:        "Siddham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 352, col: 7, offset: 9832},
						val:        "Shavian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 353, col: 7, offset: 9848},
						val:        "Sharada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 354, col: 7, offset: 9864},
						val:        "Radical",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 355, col: 7, offset: 9880},
						val:        "Osmanya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 356, col: 7, offset: 9896},
						val:        "Myanmar",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 357, col: 7, offset: 9912},
						val:        "Mandaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 358, col: 7, offset: 9928},
						val:        "Kannada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 359, col: 7, offset: 9944},
						val:        "Hanunoo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 360, col: 7, offset: 9960},
						val:        "Grantha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 361, col: 7, offset: 9976},
						val:        "Elbasan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 362, col: 7, offset: 9992},
						val:        "Deseret",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 363, col: 7, offset: 10008},
						val:        "Cypriot",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 364, col: 7, offset: 10024},
						val:        "Braille",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 365, col: 7, offset: 10040},
						val:        "Bengali",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 366, col: 7, offset: 10056},
						val:        "Avestan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 367, col: 7, offset: 10072},
						val:        "Thaana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 368, col: 7, offset: 10087},
						val:        "Telugu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 369, col: 7, offset: 10102},
						val:        "Tai_Le",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 370, col: 7, offset: 10117},
						val:        "Syriac",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 371, col: 7, offset: 10132},
						val:        "Rejang",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 372, col: 7, offset: 10147},
						val:        "Lydian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 373, col: 7, offset: 10162},
						val:        "Lycian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 374, col: 7, offset: 10177},
						val:        "Lepcha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 375, col: 7, offset: 10192},
						val:        "Khojki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 376, col: 7, offset: 10207},
						val:        "Kaithi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 377, col: 7, offset: 10222},
						val:        "Hyphen",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 378, col: 7, offset: 10237},
						val:        "Hebrew",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 379, col: 7, offset: 10252},
						val:        "Hangul",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 380, col: 7, offset: 10267},
						val:        "Gothic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 381, col: 7, offset: 10282},
						val:        "Coptic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 382, col: 7, offset: 10297},
						val:        "Common",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 383, col: 7, offset: 10312},
						val:        "Chakma",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 384, col: 7, offset: 10327},
						val:        "Carian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 385, col: 7, offset: 10342},
						val:        "Brahmi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 386, col: 7, offset: 10357},
						val:        "Arabic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 387, col: 7, offset: 10372},
						val:        "Tamil",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 388, col: 7, offset: 10386},
						val:        "Takri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 389, col: 7, offset: 10400},
						val:        "STerm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 390, col: 7, offset: 10414},
						val:        "Runic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 7, offset: 10428},
						val:        "Oriya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 392, col: 7, offset: 10442},
						val:        "Ogham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 393, col: 7, offset: 10456},
						val:        "Limbu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 394, col: 7, offset: 10470},
						val:        "Latin",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 395, col: 7, offset: 10484},
						val:        "Khmer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 396, col: 7, offset: 10498},
						val:        "Greek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 7, offset: 10512},
						val:        "Buhid",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 398, col: 7, offset: 10526},
						val:        "Batak",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 399, col: 7, offset: 10540},
						val:        "Bamum",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 400, col: 7, offset: 10554},
						val:        "Thai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 401, col: 7, offset: 10567},
						val:        "Modi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 402, col: 7, offset: 10580},
						val:        "Miao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 403, col: 7, offset: 10593},
						val:        "Lisu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 404, col: 7, offset: 10606},
						val:        "Dash",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 405, col: 7, offset: 10619},
						val:        "Cham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 406, col: 7, offset: 10632},
						val:        "Vai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 407, col: 7, offset: 10644},
						val:        "Nko",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 408, col: 7, offset: 10656},
						val:        "Mro",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 409, col: 7, offset: 10668},
						val:        "Lao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 410, col: 7, offset: 10680},
						val:        "Han",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 411, col: 7, offset: 10692},
						val:        "Zs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 412, col: 7, offset: 10703},
						val:        "Zp",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 413, col: 7, offset: 10714},
						val:        "Zl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 414, col: 7, offset: 10725},
						val:        "Yi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 415, col: 7, offset: 10736},
						val:        "So",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 416, col: 7, offset: 10747},
						val:        "Sm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 417, col: 7, offset: 10758},
						val:        "Sk",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 418, col: 7, offset: 10769},
						val:        "Sc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 419, col: 7, offset: 10780},
						val:        "Ps",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 420, col: 7, offset: 10791},
						val:        "Po",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 421, col: 7, offset: 10802},
						val:        "Pi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 422, col: 7, offset: 10813},
						val:        "Pf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 423, col: 7, offset: 10824},
						val:        "Pe",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 424, col: 7, offset: 10835},
						val:        "Pd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 425, col: 7, offset: 10846},
						val:        "Pc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 426, col: 7, offset: 10857},
						val:        "No",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 427, col: 7, offset: 10868},
						val:        "Nl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 428, col: 7, offset: 10879},
						val:        "Nd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 429, col: 7, offset: 10890},
						val:        "Mn",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 430, col: 7, offset: 10901},
						val:        "Me",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 431, col: 7, offset: 10912},
						val:        "Mc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 432, col: 7, offset: 10923},
						val:        "Lu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 433, col: 7, offset: 10934},
						val:        "Lt",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 434, col: 7, offset: 10945},
						val:        "Lo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 435, col: 7, offset: 10956},
						val:        "Lm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 436, col: 7, offset: 10967},
						val:        "Ll",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 437, col: 7, offset: 10978},
						val:        "Cs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 438, col: 7, offset: 10989},
						val:        "Co",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 439, col: 7, offset: 11000},
						val:        "Cf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 440, col: 7, offset: 11011},
						val:        "Cc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 441, col: 7, offset: 11022},
						val:        "Z",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 442, col: 7, offset: 11032},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 443, col: 7, offset: 11042},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 444, col: 7, offset: 11052},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 445, col: 7, offset: 11062},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 446, col: 7, offset: 11072},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 447, col: 7, offset: 11082},
						val:        "C",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 448, col: 7, offset: 11092},
						run: (*parser).callonUnicodeClass195,
						expr: &ruleRefExpr{
							pos:  position{line: 448, col: 7, offset: 11092},
							name: "SourceChar",
						},
					},
				},
			},
		},
		{
			name: "ReservedWord",
			pos:  position{line: 452, col: 1, offset: 11167},
			expr: &choiceExpr{
				pos: position{line: 455, col: 2, offset: 11238},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 455, col: 2, offset: 11238},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 455, col: 2, offset: 11238},
								val:        "break",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 455, col: 10, offset: 11246},
								expr: &ruleRefExpr{
									pos:  position{line: 455, col: 11, offset: 11247},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 456, col: 4, offset: 11265},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 456, col: 4, offset: 11265},
								val:        "case",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 456, col: 11, offset: 11272},
								expr: &ruleRefExpr{
									pos:  position{line: 456, col: 12, offset: 11273},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 457, col: 4, offset: 11291},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 457, col: 4, offset: 11291},
								val:        "chan",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 457, col: 11, offset: 11298},
								expr: &ruleRefExpr{
									pos:  position{line: 457, col: 12, offset: 11299},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 458, col: 4, offset: 11317},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 458, col: 4, offset: 11317},
								val:        "const",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 458, col: 12, offset: 11325},
								expr: &ruleRefExpr{
									pos:  position{line: 458, col: 13, offset: 11326},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 459, col: 4, offset: 11344},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 459, col: 4, offset: 11344},
								val:        "continue",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 459, col: 15, offset: 11355},
								expr: &ruleRefExpr{
									pos:  position{line: 459, col: 16, offset: 11356},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 460, col: 4, offset: 11374},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 460, col: 4, offset: 11374},
								val:        "default",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 460, col: 14, offset: 11384},
								expr: &ruleRefExpr{
									pos:  position{line: 460, col: 15, offset: 11385},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 461, col: 4, offset: 11403},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 461, col: 4, offset: 11403},
								val:        "defer",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 461, col: 12, offset: 11411},
								expr: &ruleRefExpr{
									pos:  position{line: 461, col: 13, offset: 11412},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 462, col: 4, offset: 11430},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 462, col: 4, offset: 11430},
								val:        "else",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 462, col: 11, offset: 11437},
								expr: &ruleRefExpr{
									pos:  position{line: 462, col: 12, offset: 11438},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 463, col: 4, offset: 11456},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 463, col: 4, offset: 11456},
								val:        "fallthrough",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 463, col: 18, offset: 11470},
								expr: &ruleRefExpr{
									pos:  position{line: 463, col: 19, offset: 11471},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 464, col: 4, offset: 11489},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 464, col: 4, offset: 11489},
								val:        "for",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 464, col: 10, offset: 11495},
								expr: &ruleRefExpr{
									pos:  position{line: 464, col: 11, offset: 11496},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 465, col: 4, offset: 11514},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 465, col: 4, offset: 11514},
								val:        "func",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 465, col: 11, offset: 11521},
								expr: &ruleRefExpr{
									pos:  position{line: 465, col: 12, offset: 11522},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 466, col: 4, offset: 11540},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 466, col: 4, offset: 11540},
								val:        "goto",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 466, col: 11, offset: 11547},
								expr: &ruleRefExpr{
									pos:  position{line: 466, col: 12, offset: 11548},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 467, col: 4, offset: 11566},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 467, col: 4, offset: 11566},
								val:        "go",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 467, col: 9, offset: 11571},
								expr: &ruleRefExpr{
									pos:  position{line: 467, col: 10, offset: 11572},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 468, col: 4, offset: 11590},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 468, col: 4, offset: 11590},
								val:        "if",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 468, col: 9, offset: 11595},
								expr: &ruleRefExpr{
									pos:  position{line: 468, col: 10, offset: 11596},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 469, col: 4, offset: 11614},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 469, col: 4, offset: 11614},
								val:        "import",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 469, col: 13, offset: 11623},
								expr: &ruleRefExpr{
									pos:  position{line: 469, col: 14, offset: 11624},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 470, col: 4, offset: 11642},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 470, col: 4, offset: 11642},
								val:        "interface",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 470, col: 16, offset: 11654},
								expr: &ruleRefExpr{
									pos:  position{line: 470, col: 17, offset: 11655},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 471, col: 4, offset: 11673},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 471, col: 4, offset: 11673},
								val:        "map",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 471, col: 10, offset: 11679},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 11, offset: 11680},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 472, col: 4, offset: 11698},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 472, col: 4, offset: 11698},
								val:        "package",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 472, col: 14, offset: 11708},
								expr: &ruleRefExpr{
									pos:  position{line: 472, col: 15, offset: 11709},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 473, col: 4, offset: 11727},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 473, col: 4, offset: 11727},
								val:        "range",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 473, col: 12, offset: 11735},
								expr: &ruleRefExpr{
									pos:  position{line: 473, col: 13, offset: 11736},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 474, col: 4, offset: 11754},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 474, col: 4, offset: 11754},
								val:        "return",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 474, col: 13, offset: 11763},
								expr: &ruleRefExpr{
									pos:  position{line: 474, col: 14, offset: 11764},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 475, col: 4, offset: 11782},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 475, col: 4, offset: 11782},
								val:        "select",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 475, col: 13, offset: 11791},
								expr: &ruleRefExpr{
									pos:  position{line: 475, col: 14, offset: 11792},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 476, col: 4, offset: 11810},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 476, col: 4, offset: 11810},
								val:        "struct",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 476, col: 13, offset: 11819},
								expr: &ruleRefExpr{
									pos:  position{line: 476, col: 14, offset: 11820},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 477, col: 4, offset: 11838},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 477, col: 4, offset: 11838},
								val:        "switch",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 477, col: 13, offset: 11847},
								expr: &ruleRefExpr{
									pos:  position{line: 477, col: 14, offset: 11848},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 478, col: 4, offset: 11866},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 478, col: 4, offset: 11866},
								val:        "type",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 478, col: 11, offset: 11873},
								expr: &ruleRefExpr{
									pos:  position{line: 478, col: 12, offset: 11874},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 479, col: 4, offset: 11892},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 479, col: 4, offset: 11892},
								val:        "var",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 479, col: 10, offset: 11898},
								expr: &ruleRefExpr{
									pos:  position{line: 479, col: 11, offset: 11899},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 483, col: 4, offset: 11998},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 483, col: 4, offset: 11998},
								val:        "bool",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 483, col: 11, offset: 12005},
								expr: &ruleRefExpr{
									pos:  position{line: 483, col: 12, offset: 12006},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 484, col: 4, offset: 12024},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 484, col: 4, offset: 12024},
								val:        "byte",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 484, col: 11, offset: 12031},
								expr: &ruleRefExpr{
									pos:  position{line: 484, col: 12, offset: 12032},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 485, col: 4, offset: 12050},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 485, col: 4, offset: 12050},
								val:        "complex64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 485, col: 16, offset: 12062},
								expr: &ruleRefExpr{
									pos:  position{line: 485, col: 17, offset: 12063},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 486, col: 4, offset: 12081},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 486, col: 4, offset: 12081},
								val:        "complex128",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 486, col: 17, offset: 12094},
								expr: &ruleRefExpr{
									pos:  position{line: 486, col: 18, offset: 12095},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 487, col: 4, offset: 12113},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 487, col: 4, offset: 12113},
								val:        "error",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 487, col: 12, offset: 12121},
								expr: &ruleRefExpr{
									pos:  position{line: 487, col: 13, offset: 12122},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 488, col: 4, offset: 12140},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 488, col: 4, offset: 12140},
								val:        "float32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 488, col: 14, offset: 12150},
								expr: &ruleRefExpr{
									pos:  position{line: 488, col: 15, offset: 12151},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 489, col: 4, offset: 12169},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 489, col: 4, offset: 12169},
								val:        "float64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 489, col: 14, offset: 12179},
								expr: &ruleRefExpr{
									pos:  position{line: 489, col: 15, offset: 12180},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 490, col: 4, offset: 12198},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 490, col: 4, offset: 12198},
								val:        "int8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 490, col: 11, offset: 12205},
								expr: &ruleRefExpr{
									pos:  position{line: 490, col: 12, offset: 12206},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 491, col: 4, offset: 12224},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 491, col: 4, offset: 12224},
								val:        "int16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 491, col: 12, offset: 12232},
								expr: &ruleRefExpr{
									pos:  position{line: 491, col: 13, offset: 12233},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 492, col: 4, offset: 12251},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 492, col: 4, offset: 12251},
								val:        "int32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 492, col: 12, offset: 12259},
								expr: &ruleRefExpr{
									pos:  position{line: 492, col: 13, offset: 12260},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 493, col: 4, offset: 12278},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 493, col: 4, offset: 12278},
								val:        "int64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 493, col: 12, offset: 12286},
								expr: &ruleRefExpr{
									pos:  position{line: 493, col: 13, offset: 12287},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 494, col: 4, offset: 12305},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 494, col: 4, offset: 12305},
								val:        "int",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 494, col: 10, offset: 12311},
								expr: &ruleRefExpr{
									pos:  position{line: 494, col: 11, offset: 12312},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 495, col: 4, offset: 12330},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 495, col: 4, offset: 12330},
								val:        "rune",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 495, col: 11, offset: 12337},
								expr: &ruleRefExpr{
									pos:  position{line: 495, col: 12, offset: 12338},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 496, col: 4, offset: 12356},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 496, col: 4, offset: 12356},
								val:        "string",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 496, col: 13, offset: 12365},
								expr: &ruleRefExpr{
									pos:  position{line: 496, col: 14, offset: 12366},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 497, col: 4, offset: 12384},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 497, col: 4, offset: 12384},
								val:        "uint8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 497, col: 12, offset: 12392},
								expr: &ruleRefExpr{
									pos:  position{line: 497, col: 13, offset: 12393},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 498, col: 4, offset: 12411},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 498, col: 4, offset: 12411},
								val:        "uint16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 498, col: 13, offset: 12420},
								expr: &ruleRefExpr{
									pos:  position{line: 498, col: 14, offset: 12421},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 499, col: 4, offset: 12439},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 499, col: 4, offset: 12439},
								val:        "uint32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 499, col: 13, offset: 12448},
								expr: &ruleRefExpr{
									pos:  position{line: 499, col: 14, offset: 12449},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 500, col: 4, offset: 12467},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 500, col: 4, offset: 12467},
								val:        "uint64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 500, col: 13, offset: 12476},
								expr: &ruleRefExpr{
									pos:  position{line: 500, col: 14, offset: 12477},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 501, col: 4, offset: 12495},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 501, col: 4, offset: 12495},
								val:        "uintptr",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 501, col: 14, offset: 12505},
								expr: &ruleRefExpr{
									pos:  position{line: 501, col: 15, offset: 12506},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 502, col: 4, offset: 12524},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 502, col: 4, offset: 12524},
								val:        "uint",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 502, col: 11, offset: 12531},
								expr: &ruleRefExpr{
									pos:  position{line: 502, col: 12, offset: 12532},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 503, col: 4, offset: 12550},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 503, col: 4, offset: 12550},
								val:        "true",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 503, col: 11, offset: 12557},
								expr: &ruleRefExpr{
									pos:  position{line: 503, col: 12, offset: 12558},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 504, col: 4, offset: 12576},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 504, col: 4, offset: 12576},
								val:        "false",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 504, col: 12, offset: 12584},
								expr: &ruleRefExpr{
									pos:  position{line: 504, col: 13, offset: 12585},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 505, col: 4, offset: 12603},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 505, col: 4, offset: 12603},
								val:        "iota",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 505, col: 11, offset: 12610},
								expr: &ruleRefExpr{
									pos:  position{line: 505, col: 12, offset: 12611},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 506, col: 4, offset: 12629},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 506, col: 4, offset: 12629},
								val:        "nil",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 506, col: 10, offset: 12635},
								expr: &ruleRefExpr{
									pos:  position{line: 506, col: 11, offset: 12636},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 507, col: 4, offset: 12654},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 507, col: 4, offset: 12654},
								val:        "append",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 507, col: 13, offset: 12663},
								expr: &ruleRefExpr{
									pos:  position{line: 507, col: 14, offset: 12664},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 508, col: 4, offset: 12682},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 508, col: 4, offset: 12682},
								val:        "cap",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 508, col: 10, offset: 12688},
								expr: &ruleRefExpr{
									pos:  position{line: 508, col: 11, offset: 12689},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 509, col: 4, offset: 12707},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 509, col: 4, offset: 12707},
								val:        "close",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 509, col: 12, offset: 12715},
								expr: &ruleRefExpr{
									pos:  position{line: 509, col: 13, offset: 12716},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 510, col: 4, offset: 12734},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 510, col: 4, offset: 12734},
								val:        "complex",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 510, col: 14, offset: 12744},
								expr: &ruleRefExpr{
									pos:  position{line: 510, col: 15, offset: 12745},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 511, col: 4, offset: 12763},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 511, col: 4, offset: 12763},
								val:        "copy",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 511, col: 11, offset: 12770},
								expr: &ruleRefExpr{
									pos:  position{line: 511, col: 12, offset: 12771},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 512, col: 4, offset: 12789},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 512, col: 4, offset: 12789},
								val:        "delete",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 512, col: 13, offset: 12798},
								expr: &ruleRefExpr{
									pos:  position{line: 512, col: 14, offset: 12799},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 513, col: 4, offset: 12817},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 513, col: 4, offset: 12817},
								val:        "imag",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 513, col: 11, offset: 12824},
								expr: &ruleRefExpr{
									pos:  position{line: 513, col: 12, offset: 12825},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 514, col: 4, offset: 12843},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 514, col: 4, offset: 12843},
								val:        "len",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 514, col: 10, offset: 12849},
								expr: &ruleRefExpr{
									pos:  position{line: 514, col: 11, offset: 12850},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 515, col: 4, offset: 12868},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 515, col: 4, offset: 12868},
								val:        "make",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 515, col: 11, offset: 12875},
								expr: &ruleRefExpr{
									pos:  position{line: 515, col: 12, offset: 12876},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 516, col: 4, offset: 12894},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 516, col: 4, offset: 12894},
								val:        "new",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 516, col: 10, offset: 12900},
								expr: &ruleRefExpr{
									pos:  position{line: 516, col: 11, offset: 12901},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 517, col: 4, offset: 12919},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 517, col: 4, offset: 12919},
								val:        "panic",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 517, col: 12, offset: 12927},
								expr: &ruleRefExpr{
									pos:  position{line: 517, col: 13, offset: 12928},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 518, col: 4, offset: 12946},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 518, col: 4, offset: 12946},
								val:        "println",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 518, col: 14, offset: 12956},
								expr: &ruleRefExpr{
									pos:  position{line: 518, col: 15, offset: 12957},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 519, col: 4, offset: 12975},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 519, col: 4, offset: 12975},
								val:        "print",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 519, col: 12, offset: 12983},
								expr: &ruleRefExpr{
									pos:  position{line: 519, col: 13, offset: 12984},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 520, col: 4, offset: 13002},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 520, col: 4, offset: 13002},
								val:        "real",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 520, col: 11, offset: 13009},
								expr: &ruleRefExpr{
									pos:  position{line: 520, col: 12, offset: 13010},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 521, col: 4, offset: 13028},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 521, col: 4, offset: 13028},
								val:        "recover",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 521, col: 14, offset: 13038},
								expr: &ruleRefExpr{
									pos:  position{line: 521, col: 15, offset: 13039},
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
			pos:  position{line: 523, col: 1, offset: 13055},
			expr: &actionExpr{
				pos: position{line: 523, col: 14, offset: 13070},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 523, col: 14, offset: 13070},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 528, col: 1, offset: 13145},
			expr: &choiceExpr{
				pos: position{line: 528, col: 13, offset: 13159},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 528, col: 13, offset: 13159},
						run: (*parser).callonCodeBlock2,
						expr: &seqExpr{
							pos: position{line: 528, col: 13, offset: 13159},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 528, col: 13, offset: 13159},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 528, col: 17, offset: 13163},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 528, col: 22, offset: 13168},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 532, col: 5, offset: 13267},
						run: (*parser).callonCodeBlock7,
						expr: &seqExpr{
							pos: position{line: 532, col: 5, offset: 13267},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 532, col: 5, offset: 13267},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 532, col: 9, offset: 13271},
									name: "Code",
								},
								&ruleRefExpr{
									pos:  position{line: 532, col: 14, offset: 13276},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 536, col: 1, offset: 13341},
			expr: &zeroOrMoreExpr{
				pos: position{line: 536, col: 8, offset: 13350},
				expr: &choiceExpr{
					pos: position{line: 536, col: 10, offset: 13352},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 536, col: 10, offset: 13352},
							expr: &seqExpr{
								pos: position{line: 536, col: 12, offset: 13354},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 536, col: 12, offset: 13354},
										expr: &charClassMatcher{
											pos:        position{line: 536, col: 13, offset: 13355},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 536, col: 18, offset: 13360},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 536, col: 34, offset: 13376},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 536, col: 34, offset: 13376},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 536, col: 38, offset: 13380},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 536, col: 43, offset: 13385},
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
			pos:  position{line: 538, col: 1, offset: 13393},
			expr: &zeroOrMoreExpr{
				pos: position{line: 538, col: 6, offset: 13400},
				expr: &choiceExpr{
					pos: position{line: 538, col: 8, offset: 13402},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 538, col: 8, offset: 13402},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 538, col: 21, offset: 13415},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 538, col: 27, offset: 13421},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 539, col: 1, offset: 13432},
			expr: &zeroOrMoreExpr{
				pos: position{line: 539, col: 5, offset: 13438},
				expr: &choiceExpr{
					pos: position{line: 539, col: 7, offset: 13440},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 539, col: 7, offset: 13440},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 539, col: 20, offset: 13453},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 541, col: 1, offset: 13490},
			expr: &charClassMatcher{
				pos:        position{line: 541, col: 14, offset: 13505},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 542, col: 1, offset: 13513},
			expr: &litMatcher{
				pos:        position{line: 542, col: 7, offset: 13521},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 543, col: 1, offset: 13526},
			expr: &choiceExpr{
				pos: position{line: 543, col: 7, offset: 13534},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 543, col: 7, offset: 13534},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 543, col: 7, offset: 13534},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 543, col: 10, offset: 13537},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 543, col: 16, offset: 13543},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 543, col: 16, offset: 13543},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 543, col: 18, offset: 13545},
								expr: &ruleRefExpr{
									pos:  position{line: 543, col: 18, offset: 13545},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 543, col: 37, offset: 13564},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 543, col: 43, offset: 13570},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 543, col: 43, offset: 13570},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 543, col: 46, offset: 13573},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 545, col: 1, offset: 13578},
			expr: &notExpr{
				pos: position{line: 545, col: 7, offset: 13586},
				expr: &anyMatcher{
					line: 545, col: 8, offset: 13587,
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
		// an invalid string literal raises an error in the escape rules,
		// so simply replace the literal with an empty string here to
		// avoid a cascade of errors.
		s = ""
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

func (c *current) onDoubleStringEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonDoubleStringEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleStringEscape5()
}

func (c *current) onSingleStringEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonSingleStringEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleStringEscape5()
}

func (c *current) onOctalEscape6() (interface{}, error) {
	return nil, errors.New("invalid octal escape")
}

func (p *parser) callonOctalEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOctalEscape6()
}

func (c *current) onHexEscape6() (interface{}, error) {
	return nil, errors.New("invalid hexadecimal escape")
}

func (p *parser) callonHexEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexEscape6()
}

func (c *current) onLongUnicodeEscape12() (interface{}, error) {
	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonLongUnicodeEscape12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLongUnicodeEscape12()
}

func (c *current) onShortUnicodeEscape8() (interface{}, error) {
	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonShortUnicodeEscape8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onShortUnicodeEscape8()
}

func (c *current) onCharClassMatcher2() (interface{}, error) {
	pos := c.astPos()
	cc := ast.NewCharClassMatcher(pos, string(c.text))
	return cc, nil
}

func (p *parser) callonCharClassMatcher2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher2()
}

func (c *current) onCharClassMatcher15() (interface{}, error) {
	return ast.NewCharClassMatcher(c.astPos(), "[]"), errors.New("character class not terminated")
}

func (p *parser) callonCharClassMatcher15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher15()
}

func (c *current) onCharClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonCharClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassEscape5()
}

func (c *current) onUnicodeClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid Unicode class escape")
}

func (p *parser) callonUnicodeClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape5()
}

func (c *current) onUnicodeClass195() (interface{}, error) {
	return nil, errors.New("invalid Unicode class escape")
}

func (p *parser) callonUnicodeClass195() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClass195()
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

func (c *current) onCodeBlock2() (interface{}, error) {
	pos := c.astPos()
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonCodeBlock2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock2()
}

func (c *current) onCodeBlock7() (interface{}, error) {
	return nil, errors.New("code block not terminated")
}

func (p *parser) callonCodeBlock7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock7()
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
