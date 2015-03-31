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
										&ruleRefExpr{
											pos:  position{line: 192, col: 55, offset: 5488},
											name: "SingleStringChar",
										},
										&choiceExpr{
											pos: position{line: 192, col: 74, offset: 5507},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 192, col: 74, offset: 5507},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 192, col: 80, offset: 5513},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 192, col: 90, offset: 5523},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 192, col: 90, offset: 5523},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 192, col: 94, offset: 5527},
											expr: &ruleRefExpr{
												pos:  position{line: 192, col: 94, offset: 5527},
												name: "RawStringChar",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 192, col: 109, offset: 5542},
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
			pos:  position{line: 196, col: 1, offset: 5644},
			expr: &choiceExpr{
				pos: position{line: 196, col: 20, offset: 5665},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 196, col: 20, offset: 5665},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 196, col: 20, offset: 5665},
								expr: &choiceExpr{
									pos: position{line: 196, col: 23, offset: 5668},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 196, col: 23, offset: 5668},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 196, col: 29, offset: 5674},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 196, col: 36, offset: 5681},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 196, col: 42, offset: 5687},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 196, col: 55, offset: 5700},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 196, col: 55, offset: 5700},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 196, col: 60, offset: 5705},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 197, col: 1, offset: 5724},
			expr: &choiceExpr{
				pos: position{line: 197, col: 20, offset: 5745},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 197, col: 20, offset: 5745},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 197, col: 20, offset: 5745},
								expr: &choiceExpr{
									pos: position{line: 197, col: 23, offset: 5748},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 197, col: 23, offset: 5748},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 197, col: 29, offset: 5754},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 197, col: 36, offset: 5761},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 197, col: 42, offset: 5767},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 197, col: 55, offset: 5780},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 197, col: 55, offset: 5780},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 197, col: 60, offset: 5785},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 198, col: 1, offset: 5804},
			expr: &seqExpr{
				pos: position{line: 198, col: 17, offset: 5822},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 198, col: 17, offset: 5822},
						expr: &litMatcher{
							pos:        position{line: 198, col: 18, offset: 5823},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 198, col: 22, offset: 5827},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 200, col: 1, offset: 5839},
			expr: &choiceExpr{
				pos: position{line: 200, col: 22, offset: 5862},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 200, col: 24, offset: 5864},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 200, col: 24, offset: 5864},
								val:        "\"",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 200, col: 30, offset: 5870},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 201, col: 7, offset: 5899},
						run: (*parser).callonDoubleStringEscape5,
						expr: &ruleRefExpr{
							pos:  position{line: 201, col: 7, offset: 5899},
							name: "SourceChar",
						},
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 204, col: 1, offset: 5969},
			expr: &choiceExpr{
				pos: position{line: 204, col: 22, offset: 5992},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 204, col: 24, offset: 5994},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 204, col: 24, offset: 5994},
								val:        "'",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 30, offset: 6000},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 205, col: 7, offset: 6029},
						run: (*parser).callonSingleStringEscape5,
						expr: &ruleRefExpr{
							pos:  position{line: 205, col: 7, offset: 6029},
							name: "SourceChar",
						},
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 209, col: 1, offset: 6100},
			expr: &choiceExpr{
				pos: position{line: 209, col: 24, offset: 6125},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 209, col: 24, offset: 6125},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 43, offset: 6144},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 57, offset: 6158},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 69, offset: 6170},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 209, col: 89, offset: 6190},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 210, col: 1, offset: 6209},
			expr: &choiceExpr{
				pos: position{line: 210, col: 20, offset: 6230},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 210, col: 20, offset: 6230},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 26, offset: 6236},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 32, offset: 6242},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 38, offset: 6248},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 44, offset: 6254},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 50, offset: 6260},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 56, offset: 6266},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 62, offset: 6272},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 211, col: 1, offset: 6277},
			expr: &choiceExpr{
				pos: position{line: 211, col: 15, offset: 6293},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 211, col: 15, offset: 6293},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 211, col: 15, offset: 6293},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 26, offset: 6304},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 211, col: 37, offset: 6315},
								name: "OctalDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 212, col: 7, offset: 6332},
						run: (*parser).callonOctalEscape6,
						expr: &seqExpr{
							pos: position{line: 212, col: 7, offset: 6332},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 212, col: 7, offset: 6332},
									name: "OctalDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 212, col: 18, offset: 6343},
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
			pos:  position{line: 215, col: 1, offset: 6409},
			expr: &choiceExpr{
				pos: position{line: 215, col: 13, offset: 6423},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 215, col: 13, offset: 6423},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 215, col: 13, offset: 6423},
								val:        "x",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 17, offset: 6427},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 26, offset: 6436},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 216, col: 7, offset: 6451},
						run: (*parser).callonHexEscape6,
						expr: &seqExpr{
							pos: position{line: 216, col: 7, offset: 6451},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 216, col: 7, offset: 6451},
									val:        "x",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 216, col: 11, offset: 6455},
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
			pos:  position{line: 219, col: 1, offset: 6527},
			expr: &choiceExpr{
				pos: position{line: 219, col: 21, offset: 6549},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 219, col: 21, offset: 6549},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 219, col: 21, offset: 6549},
								val:        "U",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 25, offset: 6553},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 34, offset: 6562},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 43, offset: 6571},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 52, offset: 6580},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 61, offset: 6589},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 70, offset: 6598},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 79, offset: 6607},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 88, offset: 6616},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 220, col: 7, offset: 6631},
						run: (*parser).callonLongUnicodeEscape12,
						expr: &seqExpr{
							pos: position{line: 220, col: 7, offset: 6631},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 220, col: 7, offset: 6631},
									val:        "U",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 220, col: 11, offset: 6635},
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
			pos:  position{line: 223, col: 1, offset: 6703},
			expr: &choiceExpr{
				pos: position{line: 223, col: 22, offset: 6726},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 223, col: 22, offset: 6726},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 223, col: 22, offset: 6726},
								val:        "u",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 26, offset: 6730},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 35, offset: 6739},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 44, offset: 6748},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 53, offset: 6757},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 224, col: 7, offset: 6772},
						run: (*parser).callonShortUnicodeEscape8,
						expr: &seqExpr{
							pos: position{line: 224, col: 7, offset: 6772},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 224, col: 7, offset: 6772},
									val:        "u",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 224, col: 11, offset: 6776},
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
			pos:  position{line: 228, col: 1, offset: 6845},
			expr: &charClassMatcher{
				pos:        position{line: 228, col: 14, offset: 6860},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 229, col: 1, offset: 6866},
			expr: &charClassMatcher{
				pos:        position{line: 229, col: 16, offset: 6883},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 230, col: 1, offset: 6889},
			expr: &charClassMatcher{
				pos:        position{line: 230, col: 12, offset: 6902},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 232, col: 1, offset: 6913},
			expr: &choiceExpr{
				pos: position{line: 232, col: 20, offset: 6934},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 232, col: 20, offset: 6934},
						run: (*parser).callonCharClassMatcher2,
						expr: &seqExpr{
							pos: position{line: 232, col: 20, offset: 6934},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 232, col: 20, offset: 6934},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 232, col: 24, offset: 6938},
									expr: &choiceExpr{
										pos: position{line: 232, col: 26, offset: 6940},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 232, col: 26, offset: 6940},
												name: "ClassCharRange",
											},
											&ruleRefExpr{
												pos:  position{line: 232, col: 43, offset: 6957},
												name: "ClassChar",
											},
											&seqExpr{
												pos: position{line: 232, col: 55, offset: 6969},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 232, col: 55, offset: 6969},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 232, col: 60, offset: 6974},
														name: "UnicodeClassEscape",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 232, col: 82, offset: 6996},
									val:        "]",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 232, col: 86, offset: 7000},
									expr: &litMatcher{
										pos:        position{line: 232, col: 86, offset: 7000},
										val:        "i",
										ignoreCase: false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 236, col: 5, offset: 7107},
						run: (*parser).callonCharClassMatcher15,
						expr: &seqExpr{
							pos: position{line: 236, col: 5, offset: 7107},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 236, col: 5, offset: 7107},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 236, col: 9, offset: 7111},
									expr: &seqExpr{
										pos: position{line: 236, col: 11, offset: 7113},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 236, col: 11, offset: 7113},
												expr: &choiceExpr{
													pos: position{line: 236, col: 14, offset: 7116},
													alternatives: []interface{}{
														&ruleRefExpr{
															pos:  position{line: 236, col: 14, offset: 7116},
															name: "EOL",
														},
														&litMatcher{
															pos:        position{line: 236, col: 20, offset: 7122},
															val:        "]",
															ignoreCase: false,
														},
													},
												},
											},
											&ruleRefExpr{
												pos:  position{line: 236, col: 26, offset: 7128},
												name: "SourceChar",
											},
										},
									},
								},
								&choiceExpr{
									pos: position{line: 236, col: 42, offset: 7144},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 236, col: 42, offset: 7144},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 236, col: 48, offset: 7150},
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
			pos:  position{line: 240, col: 1, offset: 7260},
			expr: &seqExpr{
				pos: position{line: 240, col: 18, offset: 7279},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 240, col: 18, offset: 7279},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 240, col: 28, offset: 7289},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 240, col: 32, offset: 7293},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 241, col: 1, offset: 7303},
			expr: &choiceExpr{
				pos: position{line: 241, col: 13, offset: 7317},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 241, col: 13, offset: 7317},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 241, col: 13, offset: 7317},
								expr: &choiceExpr{
									pos: position{line: 241, col: 16, offset: 7320},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 241, col: 16, offset: 7320},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 241, col: 22, offset: 7326},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 241, col: 29, offset: 7333},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 35, offset: 7339},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 241, col: 48, offset: 7352},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 241, col: 48, offset: 7352},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 53, offset: 7357},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 242, col: 1, offset: 7373},
			expr: &choiceExpr{
				pos: position{line: 242, col: 19, offset: 7393},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 242, col: 21, offset: 7395},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 242, col: 21, offset: 7395},
								val:        "]",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 242, col: 27, offset: 7401},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 243, col: 7, offset: 7430},
						run: (*parser).callonCharClassEscape5,
						expr: &seqExpr{
							pos: position{line: 243, col: 7, offset: 7430},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 243, col: 7, offset: 7430},
									expr: &litMatcher{
										pos:        position{line: 243, col: 8, offset: 7431},
										val:        "p",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 243, col: 12, offset: 7435},
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
			pos:  position{line: 247, col: 1, offset: 7506},
			expr: &seqExpr{
				pos: position{line: 247, col: 22, offset: 7529},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 247, col: 22, offset: 7529},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 248, col: 7, offset: 7542},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 248, col: 7, offset: 7542},
								name: "SingleCharUnicodeClass",
							},
							&actionExpr{
								pos: position{line: 249, col: 7, offset: 7571},
								run: (*parser).callonUnicodeClassEscape5,
								expr: &seqExpr{
									pos: position{line: 249, col: 7, offset: 7571},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 249, col: 7, offset: 7571},
											expr: &choiceExpr{
												pos: position{line: 249, col: 10, offset: 7574},
												alternatives: []interface{}{
													&litMatcher{
														pos:        position{line: 249, col: 10, offset: 7574},
														val:        "{",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 249, col: 16, offset: 7580},
														name: "EOL",
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 249, col: 22, offset: 7586},
											name: "SourceChar",
										},
									},
								},
							},
							&seqExpr{
								pos: position{line: 250, col: 7, offset: 7662},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 250, col: 7, offset: 7662},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 250, col: 11, offset: 7666},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 250, col: 24, offset: 7679},
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
			pos:  position{line: 252, col: 1, offset: 7689},
			expr: &charClassMatcher{
				pos:        position{line: 252, col: 26, offset: 7716},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 255, col: 1, offset: 7789},
			expr: &choiceExpr{
				pos: position{line: 255, col: 16, offset: 7806},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 255, col: 16, offset: 7806},
						val:        "Other_Default_Ignorable_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 256, col: 7, offset: 7849},
						val:        "Noncharacter_Code_Point",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 257, col: 7, offset: 7881},
						val:        "Logical_Order_Exception",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 258, col: 7, offset: 7913},
						val:        "Inscriptional_Parthian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 259, col: 7, offset: 7944},
						val:        "Other_Grapheme_Extend",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 260, col: 7, offset: 7974},
						val:        "Inscriptional_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 261, col: 7, offset: 8004},
						val:        "Terminal_Punctuation",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 262, col: 7, offset: 8033},
						val:        "Meroitic_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 263, col: 7, offset: 8062},
						val:        "IDS_Trinary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 264, col: 7, offset: 8091},
						val:        "Egyptian_Hieroglyphs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 265, col: 7, offset: 8120},
						val:        "Pattern_White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 266, col: 7, offset: 8148},
						val:        "IDS_Binary_Operator",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 267, col: 7, offset: 8176},
						val:        "Canadian_Aboriginal",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 268, col: 7, offset: 8204},
						val:        "Variation_Selector",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 269, col: 7, offset: 8231},
						val:        "Caucasian_Albanian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 270, col: 7, offset: 8258},
						val:        "Unified_Ideograph",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 271, col: 7, offset: 8284},
						val:        "Other_ID_Continue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 272, col: 7, offset: 8310},
						val:        "Old_South_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 273, col: 7, offset: 8336},
						val:        "Old_North_Arabian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 274, col: 7, offset: 8362},
						val:        "Other_Alphabetic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 275, col: 7, offset: 8387},
						val:        "Meroitic_Cursive",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 276, col: 7, offset: 8412},
						val:        "Imperial_Aramaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 277, col: 7, offset: 8437},
						val:        "Psalter_Pahlavi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 278, col: 7, offset: 8461},
						val:        "Other_Uppercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 279, col: 7, offset: 8485},
						val:        "Other_Lowercase",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 280, col: 7, offset: 8509},
						val:        "ASCII_Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 281, col: 7, offset: 8533},
						val:        "Quotation_Mark",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 282, col: 7, offset: 8556},
						val:        "Pattern_Syntax",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 283, col: 7, offset: 8579},
						val:        "Other_ID_Start",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 284, col: 7, offset: 8602},
						val:        "Mende_Kikakui",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 285, col: 7, offset: 8624},
						val:        "Syloti_Nagri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 286, col: 7, offset: 8645},
						val:        "Sora_Sompeng",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 287, col: 7, offset: 8666},
						val:        "Pahawh_Hmong",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 288, col: 7, offset: 8687},
						val:        "Meetei_Mayek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 289, col: 7, offset: 8708},
						val:        "Join_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 290, col: 7, offset: 8729},
						val:        "Bidi_Control",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 291, col: 7, offset: 8750},
						val:        "White_Space",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 292, col: 7, offset: 8770},
						val:        "Warang_Citi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 293, col: 7, offset: 8790},
						val:        "Soft_Dotted",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 294, col: 7, offset: 8810},
						val:        "Pau_Cin_Hau",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 295, col: 7, offset: 8830},
						val:        "Old_Persian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 296, col: 7, offset: 8850},
						val:        "New_Tai_Lue",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 297, col: 7, offset: 8870},
						val:        "Ideographic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 298, col: 7, offset: 8890},
						val:        "Saurashtra",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 299, col: 7, offset: 8909},
						val:        "Phoenician",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 300, col: 7, offset: 8928},
						val:        "Other_Math",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 301, col: 7, offset: 8947},
						val:        "Old_Turkic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 302, col: 7, offset: 8966},
						val:        "Old_Permic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 303, col: 7, offset: 8985},
						val:        "Old_Italic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 304, col: 7, offset: 9004},
						val:        "Manichaean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 305, col: 7, offset: 9023},
						val:        "Kharoshthi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 306, col: 7, offset: 9042},
						val:        "Glagolitic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 307, col: 7, offset: 9061},
						val:        "Devanagari",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 308, col: 7, offset: 9080},
						val:        "Deprecated",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 309, col: 7, offset: 9099},
						val:        "Sundanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 310, col: 7, offset: 9117},
						val:        "Samaritan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 311, col: 7, offset: 9135},
						val:        "Palmyrene",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 312, col: 7, offset: 9153},
						val:        "Nabataean",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 313, col: 7, offset: 9171},
						val:        "Mongolian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 314, col: 7, offset: 9189},
						val:        "Malayalam",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 315, col: 7, offset: 9207},
						val:        "Khudawadi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 316, col: 7, offset: 9225},
						val:        "Inherited",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 317, col: 7, offset: 9243},
						val:        "Hex_Digit",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 318, col: 7, offset: 9261},
						val:        "Diacritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 319, col: 7, offset: 9279},
						val:        "Cuneiform",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 320, col: 7, offset: 9297},
						val:        "Bassa_Vah",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 321, col: 7, offset: 9315},
						val:        "Ugaritic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 322, col: 7, offset: 9332},
						val:        "Tifinagh",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 323, col: 7, offset: 9349},
						val:        "Tai_Viet",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 324, col: 7, offset: 9366},
						val:        "Tai_Tham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 325, col: 7, offset: 9383},
						val:        "Tagbanwa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 326, col: 7, offset: 9400},
						val:        "Phags_Pa",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 327, col: 7, offset: 9417},
						val:        "Ol_Chiki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 328, col: 7, offset: 9434},
						val:        "Mahajani",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 329, col: 7, offset: 9451},
						val:        "Linear_B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 330, col: 7, offset: 9468},
						val:        "Linear_A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 331, col: 7, offset: 9485},
						val:        "Kayah_Li",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 332, col: 7, offset: 9502},
						val:        "Katakana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 333, col: 7, offset: 9519},
						val:        "Javanese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 334, col: 7, offset: 9536},
						val:        "Hiragana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 335, col: 7, offset: 9553},
						val:        "Gurmukhi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 336, col: 7, offset: 9570},
						val:        "Gujarati",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 337, col: 7, offset: 9587},
						val:        "Georgian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 338, col: 7, offset: 9604},
						val:        "Extender",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 339, col: 7, offset: 9621},
						val:        "Ethiopic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 340, col: 7, offset: 9638},
						val:        "Duployan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 341, col: 7, offset: 9655},
						val:        "Cyrillic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 342, col: 7, offset: 9672},
						val:        "Cherokee",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 343, col: 7, offset: 9689},
						val:        "Buginese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 344, col: 7, offset: 9706},
						val:        "Bopomofo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 345, col: 7, offset: 9723},
						val:        "Balinese",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 346, col: 7, offset: 9740},
						val:        "Armenian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 347, col: 7, offset: 9757},
						val:        "Tirhuta",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 348, col: 7, offset: 9773},
						val:        "Tibetan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 349, col: 7, offset: 9789},
						val:        "Tagalog",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 350, col: 7, offset: 9805},
						val:        "Sinhala",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 351, col: 7, offset: 9821},
						val:        "Siddham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 352, col: 7, offset: 9837},
						val:        "Shavian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 353, col: 7, offset: 9853},
						val:        "Sharada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 354, col: 7, offset: 9869},
						val:        "Radical",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 355, col: 7, offset: 9885},
						val:        "Osmanya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 356, col: 7, offset: 9901},
						val:        "Myanmar",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 357, col: 7, offset: 9917},
						val:        "Mandaic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 358, col: 7, offset: 9933},
						val:        "Kannada",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 359, col: 7, offset: 9949},
						val:        "Hanunoo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 360, col: 7, offset: 9965},
						val:        "Grantha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 361, col: 7, offset: 9981},
						val:        "Elbasan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 362, col: 7, offset: 9997},
						val:        "Deseret",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 363, col: 7, offset: 10013},
						val:        "Cypriot",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 364, col: 7, offset: 10029},
						val:        "Braille",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 365, col: 7, offset: 10045},
						val:        "Bengali",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 366, col: 7, offset: 10061},
						val:        "Avestan",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 367, col: 7, offset: 10077},
						val:        "Thaana",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 368, col: 7, offset: 10092},
						val:        "Telugu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 369, col: 7, offset: 10107},
						val:        "Tai_Le",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 370, col: 7, offset: 10122},
						val:        "Syriac",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 371, col: 7, offset: 10137},
						val:        "Rejang",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 372, col: 7, offset: 10152},
						val:        "Lydian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 373, col: 7, offset: 10167},
						val:        "Lycian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 374, col: 7, offset: 10182},
						val:        "Lepcha",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 375, col: 7, offset: 10197},
						val:        "Khojki",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 376, col: 7, offset: 10212},
						val:        "Kaithi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 377, col: 7, offset: 10227},
						val:        "Hyphen",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 378, col: 7, offset: 10242},
						val:        "Hebrew",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 379, col: 7, offset: 10257},
						val:        "Hangul",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 380, col: 7, offset: 10272},
						val:        "Gothic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 381, col: 7, offset: 10287},
						val:        "Coptic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 382, col: 7, offset: 10302},
						val:        "Common",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 383, col: 7, offset: 10317},
						val:        "Chakma",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 384, col: 7, offset: 10332},
						val:        "Carian",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 385, col: 7, offset: 10347},
						val:        "Brahmi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 386, col: 7, offset: 10362},
						val:        "Arabic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 387, col: 7, offset: 10377},
						val:        "Tamil",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 388, col: 7, offset: 10391},
						val:        "Takri",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 389, col: 7, offset: 10405},
						val:        "STerm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 390, col: 7, offset: 10419},
						val:        "Runic",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 391, col: 7, offset: 10433},
						val:        "Oriya",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 392, col: 7, offset: 10447},
						val:        "Ogham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 393, col: 7, offset: 10461},
						val:        "Limbu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 394, col: 7, offset: 10475},
						val:        "Latin",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 395, col: 7, offset: 10489},
						val:        "Khmer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 396, col: 7, offset: 10503},
						val:        "Greek",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 397, col: 7, offset: 10517},
						val:        "Buhid",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 398, col: 7, offset: 10531},
						val:        "Batak",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 399, col: 7, offset: 10545},
						val:        "Bamum",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 400, col: 7, offset: 10559},
						val:        "Thai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 401, col: 7, offset: 10572},
						val:        "Modi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 402, col: 7, offset: 10585},
						val:        "Miao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 403, col: 7, offset: 10598},
						val:        "Lisu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 404, col: 7, offset: 10611},
						val:        "Dash",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 405, col: 7, offset: 10624},
						val:        "Cham",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 406, col: 7, offset: 10637},
						val:        "Vai",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 407, col: 7, offset: 10649},
						val:        "Nko",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 408, col: 7, offset: 10661},
						val:        "Mro",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 409, col: 7, offset: 10673},
						val:        "Lao",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 410, col: 7, offset: 10685},
						val:        "Han",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 411, col: 7, offset: 10697},
						val:        "Zs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 412, col: 7, offset: 10708},
						val:        "Zp",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 413, col: 7, offset: 10719},
						val:        "Zl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 414, col: 7, offset: 10730},
						val:        "Yi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 415, col: 7, offset: 10741},
						val:        "So",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 416, col: 7, offset: 10752},
						val:        "Sm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 417, col: 7, offset: 10763},
						val:        "Sk",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 418, col: 7, offset: 10774},
						val:        "Sc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 419, col: 7, offset: 10785},
						val:        "Ps",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 420, col: 7, offset: 10796},
						val:        "Po",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 421, col: 7, offset: 10807},
						val:        "Pi",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 422, col: 7, offset: 10818},
						val:        "Pf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 423, col: 7, offset: 10829},
						val:        "Pe",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 424, col: 7, offset: 10840},
						val:        "Pd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 425, col: 7, offset: 10851},
						val:        "Pc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 426, col: 7, offset: 10862},
						val:        "No",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 427, col: 7, offset: 10873},
						val:        "Nl",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 428, col: 7, offset: 10884},
						val:        "Nd",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 429, col: 7, offset: 10895},
						val:        "Mn",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 430, col: 7, offset: 10906},
						val:        "Me",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 431, col: 7, offset: 10917},
						val:        "Mc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 432, col: 7, offset: 10928},
						val:        "Lu",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 433, col: 7, offset: 10939},
						val:        "Lt",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 434, col: 7, offset: 10950},
						val:        "Lo",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 435, col: 7, offset: 10961},
						val:        "Lm",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 436, col: 7, offset: 10972},
						val:        "Ll",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 437, col: 7, offset: 10983},
						val:        "Cs",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 438, col: 7, offset: 10994},
						val:        "Co",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 439, col: 7, offset: 11005},
						val:        "Cf",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 440, col: 7, offset: 11016},
						val:        "Cc",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 441, col: 7, offset: 11027},
						val:        "Z",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 442, col: 7, offset: 11037},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 443, col: 7, offset: 11047},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 444, col: 7, offset: 11057},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 445, col: 7, offset: 11067},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 446, col: 7, offset: 11077},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 447, col: 7, offset: 11087},
						val:        "C",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 448, col: 7, offset: 11097},
						run: (*parser).callonUnicodeClass195,
						expr: &ruleRefExpr{
							pos:  position{line: 448, col: 7, offset: 11097},
							name: "SourceChar",
						},
					},
				},
			},
		},
		{
			name: "ReservedWord",
			pos:  position{line: 452, col: 1, offset: 11172},
			expr: &choiceExpr{
				pos: position{line: 455, col: 2, offset: 11243},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 455, col: 2, offset: 11243},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 455, col: 2, offset: 11243},
								val:        "break",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 455, col: 10, offset: 11251},
								expr: &ruleRefExpr{
									pos:  position{line: 455, col: 11, offset: 11252},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 456, col: 4, offset: 11270},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 456, col: 4, offset: 11270},
								val:        "case",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 456, col: 11, offset: 11277},
								expr: &ruleRefExpr{
									pos:  position{line: 456, col: 12, offset: 11278},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 457, col: 4, offset: 11296},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 457, col: 4, offset: 11296},
								val:        "chan",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 457, col: 11, offset: 11303},
								expr: &ruleRefExpr{
									pos:  position{line: 457, col: 12, offset: 11304},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 458, col: 4, offset: 11322},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 458, col: 4, offset: 11322},
								val:        "const",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 458, col: 12, offset: 11330},
								expr: &ruleRefExpr{
									pos:  position{line: 458, col: 13, offset: 11331},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 459, col: 4, offset: 11349},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 459, col: 4, offset: 11349},
								val:        "continue",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 459, col: 15, offset: 11360},
								expr: &ruleRefExpr{
									pos:  position{line: 459, col: 16, offset: 11361},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 460, col: 4, offset: 11379},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 460, col: 4, offset: 11379},
								val:        "default",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 460, col: 14, offset: 11389},
								expr: &ruleRefExpr{
									pos:  position{line: 460, col: 15, offset: 11390},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 461, col: 4, offset: 11408},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 461, col: 4, offset: 11408},
								val:        "defer",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 461, col: 12, offset: 11416},
								expr: &ruleRefExpr{
									pos:  position{line: 461, col: 13, offset: 11417},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 462, col: 4, offset: 11435},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 462, col: 4, offset: 11435},
								val:        "else",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 462, col: 11, offset: 11442},
								expr: &ruleRefExpr{
									pos:  position{line: 462, col: 12, offset: 11443},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 463, col: 4, offset: 11461},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 463, col: 4, offset: 11461},
								val:        "fallthrough",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 463, col: 18, offset: 11475},
								expr: &ruleRefExpr{
									pos:  position{line: 463, col: 19, offset: 11476},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 464, col: 4, offset: 11494},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 464, col: 4, offset: 11494},
								val:        "for",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 464, col: 10, offset: 11500},
								expr: &ruleRefExpr{
									pos:  position{line: 464, col: 11, offset: 11501},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 465, col: 4, offset: 11519},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 465, col: 4, offset: 11519},
								val:        "func",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 465, col: 11, offset: 11526},
								expr: &ruleRefExpr{
									pos:  position{line: 465, col: 12, offset: 11527},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 466, col: 4, offset: 11545},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 466, col: 4, offset: 11545},
								val:        "goto",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 466, col: 11, offset: 11552},
								expr: &ruleRefExpr{
									pos:  position{line: 466, col: 12, offset: 11553},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 467, col: 4, offset: 11571},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 467, col: 4, offset: 11571},
								val:        "go",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 467, col: 9, offset: 11576},
								expr: &ruleRefExpr{
									pos:  position{line: 467, col: 10, offset: 11577},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 468, col: 4, offset: 11595},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 468, col: 4, offset: 11595},
								val:        "if",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 468, col: 9, offset: 11600},
								expr: &ruleRefExpr{
									pos:  position{line: 468, col: 10, offset: 11601},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 469, col: 4, offset: 11619},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 469, col: 4, offset: 11619},
								val:        "import",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 469, col: 13, offset: 11628},
								expr: &ruleRefExpr{
									pos:  position{line: 469, col: 14, offset: 11629},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 470, col: 4, offset: 11647},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 470, col: 4, offset: 11647},
								val:        "interface",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 470, col: 16, offset: 11659},
								expr: &ruleRefExpr{
									pos:  position{line: 470, col: 17, offset: 11660},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 471, col: 4, offset: 11678},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 471, col: 4, offset: 11678},
								val:        "map",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 471, col: 10, offset: 11684},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 11, offset: 11685},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 472, col: 4, offset: 11703},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 472, col: 4, offset: 11703},
								val:        "package",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 472, col: 14, offset: 11713},
								expr: &ruleRefExpr{
									pos:  position{line: 472, col: 15, offset: 11714},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 473, col: 4, offset: 11732},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 473, col: 4, offset: 11732},
								val:        "range",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 473, col: 12, offset: 11740},
								expr: &ruleRefExpr{
									pos:  position{line: 473, col: 13, offset: 11741},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 474, col: 4, offset: 11759},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 474, col: 4, offset: 11759},
								val:        "return",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 474, col: 13, offset: 11768},
								expr: &ruleRefExpr{
									pos:  position{line: 474, col: 14, offset: 11769},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 475, col: 4, offset: 11787},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 475, col: 4, offset: 11787},
								val:        "select",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 475, col: 13, offset: 11796},
								expr: &ruleRefExpr{
									pos:  position{line: 475, col: 14, offset: 11797},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 476, col: 4, offset: 11815},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 476, col: 4, offset: 11815},
								val:        "struct",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 476, col: 13, offset: 11824},
								expr: &ruleRefExpr{
									pos:  position{line: 476, col: 14, offset: 11825},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 477, col: 4, offset: 11843},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 477, col: 4, offset: 11843},
								val:        "switch",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 477, col: 13, offset: 11852},
								expr: &ruleRefExpr{
									pos:  position{line: 477, col: 14, offset: 11853},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 478, col: 4, offset: 11871},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 478, col: 4, offset: 11871},
								val:        "type",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 478, col: 11, offset: 11878},
								expr: &ruleRefExpr{
									pos:  position{line: 478, col: 12, offset: 11879},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 479, col: 4, offset: 11897},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 479, col: 4, offset: 11897},
								val:        "var",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 479, col: 10, offset: 11903},
								expr: &ruleRefExpr{
									pos:  position{line: 479, col: 11, offset: 11904},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 483, col: 4, offset: 12003},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 483, col: 4, offset: 12003},
								val:        "bool",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 483, col: 11, offset: 12010},
								expr: &ruleRefExpr{
									pos:  position{line: 483, col: 12, offset: 12011},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 484, col: 4, offset: 12029},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 484, col: 4, offset: 12029},
								val:        "byte",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 484, col: 11, offset: 12036},
								expr: &ruleRefExpr{
									pos:  position{line: 484, col: 12, offset: 12037},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 485, col: 4, offset: 12055},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 485, col: 4, offset: 12055},
								val:        "complex64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 485, col: 16, offset: 12067},
								expr: &ruleRefExpr{
									pos:  position{line: 485, col: 17, offset: 12068},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 486, col: 4, offset: 12086},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 486, col: 4, offset: 12086},
								val:        "complex128",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 486, col: 17, offset: 12099},
								expr: &ruleRefExpr{
									pos:  position{line: 486, col: 18, offset: 12100},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 487, col: 4, offset: 12118},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 487, col: 4, offset: 12118},
								val:        "error",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 487, col: 12, offset: 12126},
								expr: &ruleRefExpr{
									pos:  position{line: 487, col: 13, offset: 12127},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 488, col: 4, offset: 12145},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 488, col: 4, offset: 12145},
								val:        "float32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 488, col: 14, offset: 12155},
								expr: &ruleRefExpr{
									pos:  position{line: 488, col: 15, offset: 12156},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 489, col: 4, offset: 12174},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 489, col: 4, offset: 12174},
								val:        "float64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 489, col: 14, offset: 12184},
								expr: &ruleRefExpr{
									pos:  position{line: 489, col: 15, offset: 12185},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 490, col: 4, offset: 12203},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 490, col: 4, offset: 12203},
								val:        "int8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 490, col: 11, offset: 12210},
								expr: &ruleRefExpr{
									pos:  position{line: 490, col: 12, offset: 12211},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 491, col: 4, offset: 12229},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 491, col: 4, offset: 12229},
								val:        "int16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 491, col: 12, offset: 12237},
								expr: &ruleRefExpr{
									pos:  position{line: 491, col: 13, offset: 12238},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 492, col: 4, offset: 12256},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 492, col: 4, offset: 12256},
								val:        "int32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 492, col: 12, offset: 12264},
								expr: &ruleRefExpr{
									pos:  position{line: 492, col: 13, offset: 12265},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 493, col: 4, offset: 12283},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 493, col: 4, offset: 12283},
								val:        "int64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 493, col: 12, offset: 12291},
								expr: &ruleRefExpr{
									pos:  position{line: 493, col: 13, offset: 12292},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 494, col: 4, offset: 12310},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 494, col: 4, offset: 12310},
								val:        "int",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 494, col: 10, offset: 12316},
								expr: &ruleRefExpr{
									pos:  position{line: 494, col: 11, offset: 12317},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 495, col: 4, offset: 12335},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 495, col: 4, offset: 12335},
								val:        "rune",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 495, col: 11, offset: 12342},
								expr: &ruleRefExpr{
									pos:  position{line: 495, col: 12, offset: 12343},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 496, col: 4, offset: 12361},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 496, col: 4, offset: 12361},
								val:        "string",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 496, col: 13, offset: 12370},
								expr: &ruleRefExpr{
									pos:  position{line: 496, col: 14, offset: 12371},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 497, col: 4, offset: 12389},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 497, col: 4, offset: 12389},
								val:        "uint8",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 497, col: 12, offset: 12397},
								expr: &ruleRefExpr{
									pos:  position{line: 497, col: 13, offset: 12398},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 498, col: 4, offset: 12416},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 498, col: 4, offset: 12416},
								val:        "uint16",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 498, col: 13, offset: 12425},
								expr: &ruleRefExpr{
									pos:  position{line: 498, col: 14, offset: 12426},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 499, col: 4, offset: 12444},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 499, col: 4, offset: 12444},
								val:        "uint32",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 499, col: 13, offset: 12453},
								expr: &ruleRefExpr{
									pos:  position{line: 499, col: 14, offset: 12454},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 500, col: 4, offset: 12472},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 500, col: 4, offset: 12472},
								val:        "uint64",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 500, col: 13, offset: 12481},
								expr: &ruleRefExpr{
									pos:  position{line: 500, col: 14, offset: 12482},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 501, col: 4, offset: 12500},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 501, col: 4, offset: 12500},
								val:        "uintptr",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 501, col: 14, offset: 12510},
								expr: &ruleRefExpr{
									pos:  position{line: 501, col: 15, offset: 12511},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 502, col: 4, offset: 12529},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 502, col: 4, offset: 12529},
								val:        "uint",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 502, col: 11, offset: 12536},
								expr: &ruleRefExpr{
									pos:  position{line: 502, col: 12, offset: 12537},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 503, col: 4, offset: 12555},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 503, col: 4, offset: 12555},
								val:        "true",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 503, col: 11, offset: 12562},
								expr: &ruleRefExpr{
									pos:  position{line: 503, col: 12, offset: 12563},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 504, col: 4, offset: 12581},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 504, col: 4, offset: 12581},
								val:        "false",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 504, col: 12, offset: 12589},
								expr: &ruleRefExpr{
									pos:  position{line: 504, col: 13, offset: 12590},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 505, col: 4, offset: 12608},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 505, col: 4, offset: 12608},
								val:        "iota",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 505, col: 11, offset: 12615},
								expr: &ruleRefExpr{
									pos:  position{line: 505, col: 12, offset: 12616},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 506, col: 4, offset: 12634},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 506, col: 4, offset: 12634},
								val:        "nil",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 506, col: 10, offset: 12640},
								expr: &ruleRefExpr{
									pos:  position{line: 506, col: 11, offset: 12641},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 507, col: 4, offset: 12659},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 507, col: 4, offset: 12659},
								val:        "append",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 507, col: 13, offset: 12668},
								expr: &ruleRefExpr{
									pos:  position{line: 507, col: 14, offset: 12669},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 508, col: 4, offset: 12687},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 508, col: 4, offset: 12687},
								val:        "cap",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 508, col: 10, offset: 12693},
								expr: &ruleRefExpr{
									pos:  position{line: 508, col: 11, offset: 12694},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 509, col: 4, offset: 12712},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 509, col: 4, offset: 12712},
								val:        "close",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 509, col: 12, offset: 12720},
								expr: &ruleRefExpr{
									pos:  position{line: 509, col: 13, offset: 12721},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 510, col: 4, offset: 12739},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 510, col: 4, offset: 12739},
								val:        "complex",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 510, col: 14, offset: 12749},
								expr: &ruleRefExpr{
									pos:  position{line: 510, col: 15, offset: 12750},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 511, col: 4, offset: 12768},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 511, col: 4, offset: 12768},
								val:        "copy",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 511, col: 11, offset: 12775},
								expr: &ruleRefExpr{
									pos:  position{line: 511, col: 12, offset: 12776},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 512, col: 4, offset: 12794},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 512, col: 4, offset: 12794},
								val:        "delete",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 512, col: 13, offset: 12803},
								expr: &ruleRefExpr{
									pos:  position{line: 512, col: 14, offset: 12804},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 513, col: 4, offset: 12822},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 513, col: 4, offset: 12822},
								val:        "imag",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 513, col: 11, offset: 12829},
								expr: &ruleRefExpr{
									pos:  position{line: 513, col: 12, offset: 12830},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 514, col: 4, offset: 12848},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 514, col: 4, offset: 12848},
								val:        "len",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 514, col: 10, offset: 12854},
								expr: &ruleRefExpr{
									pos:  position{line: 514, col: 11, offset: 12855},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 515, col: 4, offset: 12873},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 515, col: 4, offset: 12873},
								val:        "make",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 515, col: 11, offset: 12880},
								expr: &ruleRefExpr{
									pos:  position{line: 515, col: 12, offset: 12881},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 516, col: 4, offset: 12899},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 516, col: 4, offset: 12899},
								val:        "new",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 516, col: 10, offset: 12905},
								expr: &ruleRefExpr{
									pos:  position{line: 516, col: 11, offset: 12906},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 517, col: 4, offset: 12924},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 517, col: 4, offset: 12924},
								val:        "panic",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 517, col: 12, offset: 12932},
								expr: &ruleRefExpr{
									pos:  position{line: 517, col: 13, offset: 12933},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 518, col: 4, offset: 12951},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 518, col: 4, offset: 12951},
								val:        "println",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 518, col: 14, offset: 12961},
								expr: &ruleRefExpr{
									pos:  position{line: 518, col: 15, offset: 12962},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 519, col: 4, offset: 12980},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 519, col: 4, offset: 12980},
								val:        "print",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 519, col: 12, offset: 12988},
								expr: &ruleRefExpr{
									pos:  position{line: 519, col: 13, offset: 12989},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 520, col: 4, offset: 13007},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 520, col: 4, offset: 13007},
								val:        "real",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 520, col: 11, offset: 13014},
								expr: &ruleRefExpr{
									pos:  position{line: 520, col: 12, offset: 13015},
									name: "IdentifierPart",
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 521, col: 4, offset: 13033},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 521, col: 4, offset: 13033},
								val:        "recover",
								ignoreCase: false,
							},
							&notExpr{
								pos: position{line: 521, col: 14, offset: 13043},
								expr: &ruleRefExpr{
									pos:  position{line: 521, col: 15, offset: 13044},
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
			pos:  position{line: 523, col: 1, offset: 13060},
			expr: &actionExpr{
				pos: position{line: 523, col: 14, offset: 13075},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 523, col: 14, offset: 13075},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 528, col: 1, offset: 13150},
			expr: &choiceExpr{
				pos: position{line: 528, col: 13, offset: 13164},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 528, col: 13, offset: 13164},
						run: (*parser).callonCodeBlock2,
						expr: &seqExpr{
							pos: position{line: 528, col: 13, offset: 13164},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 528, col: 13, offset: 13164},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 528, col: 17, offset: 13168},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 528, col: 22, offset: 13173},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 532, col: 5, offset: 13272},
						run: (*parser).callonCodeBlock7,
						expr: &seqExpr{
							pos: position{line: 532, col: 5, offset: 13272},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 532, col: 5, offset: 13272},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 532, col: 9, offset: 13276},
									name: "Code",
								},
								&ruleRefExpr{
									pos:  position{line: 532, col: 14, offset: 13281},
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
			pos:  position{line: 536, col: 1, offset: 13346},
			expr: &zeroOrMoreExpr{
				pos: position{line: 536, col: 8, offset: 13355},
				expr: &choiceExpr{
					pos: position{line: 536, col: 10, offset: 13357},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 536, col: 10, offset: 13357},
							expr: &seqExpr{
								pos: position{line: 536, col: 12, offset: 13359},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 536, col: 12, offset: 13359},
										expr: &charClassMatcher{
											pos:        position{line: 536, col: 13, offset: 13360},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 536, col: 18, offset: 13365},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 536, col: 34, offset: 13381},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 536, col: 34, offset: 13381},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 536, col: 38, offset: 13385},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 536, col: 43, offset: 13390},
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
			pos:  position{line: 538, col: 1, offset: 13398},
			expr: &zeroOrMoreExpr{
				pos: position{line: 538, col: 6, offset: 13405},
				expr: &choiceExpr{
					pos: position{line: 538, col: 8, offset: 13407},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 538, col: 8, offset: 13407},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 538, col: 21, offset: 13420},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 538, col: 27, offset: 13426},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 539, col: 1, offset: 13437},
			expr: &zeroOrMoreExpr{
				pos: position{line: 539, col: 5, offset: 13443},
				expr: &choiceExpr{
					pos: position{line: 539, col: 7, offset: 13445},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 539, col: 7, offset: 13445},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 539, col: 20, offset: 13458},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 541, col: 1, offset: 13495},
			expr: &charClassMatcher{
				pos:        position{line: 541, col: 14, offset: 13510},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 542, col: 1, offset: 13518},
			expr: &litMatcher{
				pos:        position{line: 542, col: 7, offset: 13526},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 543, col: 1, offset: 13531},
			expr: &choiceExpr{
				pos: position{line: 543, col: 7, offset: 13539},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 543, col: 7, offset: 13539},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 543, col: 7, offset: 13539},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 543, col: 10, offset: 13542},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 543, col: 16, offset: 13548},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 543, col: 16, offset: 13548},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 543, col: 18, offset: 13550},
								expr: &ruleRefExpr{
									pos:  position{line: 543, col: 18, offset: 13550},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 543, col: 37, offset: 13569},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 543, col: 43, offset: 13575},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 543, col: 43, offset: 13575},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 543, col: 46, offset: 13578},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 545, col: 1, offset: 13583},
			expr: &notExpr{
				pos: position{line: 545, col: 7, offset: 13591},
				expr: &anyMatcher{
					line: 545, col: 8, offset: 13592,
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
