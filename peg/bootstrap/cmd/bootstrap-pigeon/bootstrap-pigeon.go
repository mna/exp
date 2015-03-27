// Command blah does blah.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
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
	dbgFlag := flag.Bool("debug", false, "set debug mode")
	//noBuildFlag := flag.Bool("x", false, "do not build, only parse")
	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "USAGE: <cmd> FILE")
		os.Exit(1)
	}

	var in io.Reader

	nm := "stdin"
	if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer f.Close()
		in = f
		nm = flag.Arg(0)
	} else {
		in = bufio.NewReader(os.Stdin)
	}

	debug = *dbgFlag
	res, err := Parse(nm, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func (c *current) astPos() ast.Pos {
	return ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 44, col: 1, offset: 765},
			expr: &actionExpr{
				pos: position{line: 44, col: 11, offset: 777},
				run: (*parser).callonGrammar_1,
				expr: &seqExpr{
					pos: position{line: 44, col: 11, offset: 777},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 44, col: 11, offset: 777},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 44, col: 14, offset: 780},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 44, col: 28, offset: 794},
								expr: &seqExpr{
									pos: position{line: 44, col: 28, offset: 794},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 44, col: 28, offset: 794},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 44, col: 40, offset: 806},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 44, col: 46, offset: 812},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 44, col: 54, offset: 820},
								expr: &seqExpr{
									pos: position{line: 44, col: 54, offset: 820},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 44, col: 54, offset: 820},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 44, col: 59, offset: 825},
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
			pos:  position{line: 63, col: 1, offset: 1272},
			expr: &actionExpr{
				pos: position{line: 63, col: 15, offset: 1288},
				run: (*parser).callonInitializer_1,
				expr: &seqExpr{
					pos: position{line: 63, col: 15, offset: 1288},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 63, col: 15, offset: 1288},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 63, col: 20, offset: 1293},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 63, col: 30, offset: 1303},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 67, col: 1, offset: 1333},
			expr: &actionExpr{
				pos: position{line: 67, col: 8, offset: 1342},
				run: (*parser).callonRule_1,
				expr: &seqExpr{
					pos: position{line: 67, col: 8, offset: 1342},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 67, col: 8, offset: 1342},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 67, col: 13, offset: 1347},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 67, col: 28, offset: 1362},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 67, col: 31, offset: 1365},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 67, col: 41, offset: 1375},
								expr: &seqExpr{
									pos: position{line: 67, col: 41, offset: 1375},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 67, col: 41, offset: 1375},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 67, col: 55, offset: 1389},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 67, col: 61, offset: 1395},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 67, col: 71, offset: 1405},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 67, col: 74, offset: 1408},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 67, col: 79, offset: 1413},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 67, col: 90, offset: 1424},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 80, col: 1, offset: 1710},
			expr: &ruleRefExpr{
				pos:  position{line: 80, col: 14, offset: 1725},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 82, col: 1, offset: 1737},
			expr: &actionExpr{
				pos: position{line: 82, col: 14, offset: 1752},
				run: (*parser).callonChoiceExpr_1,
				expr: &seqExpr{
					pos: position{line: 82, col: 14, offset: 1752},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 82, col: 14, offset: 1752},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 82, col: 20, offset: 1758},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 82, col: 31, offset: 1769},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 82, col: 38, offset: 1776},
								expr: &seqExpr{
									pos: position{line: 82, col: 38, offset: 1776},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 82, col: 38, offset: 1776},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 82, col: 41, offset: 1779},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 82, col: 45, offset: 1783},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 82, col: 48, offset: 1786},
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
			pos:  position{line: 97, col: 1, offset: 2193},
			expr: &actionExpr{
				pos: position{line: 97, col: 14, offset: 2208},
				run: (*parser).callonActionExpr_1,
				expr: &seqExpr{
					pos: position{line: 97, col: 14, offset: 2208},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 97, col: 14, offset: 2208},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 97, col: 19, offset: 2213},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 97, col: 27, offset: 2221},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 97, col: 34, offset: 2228},
								expr: &seqExpr{
									pos: position{line: 97, col: 34, offset: 2228},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 97, col: 34, offset: 2228},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 97, col: 37, offset: 2231},
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
			pos:  position{line: 111, col: 1, offset: 2499},
			expr: &actionExpr{
				pos: position{line: 111, col: 11, offset: 2511},
				run: (*parser).callonSeqExpr_1,
				expr: &seqExpr{
					pos: position{line: 111, col: 11, offset: 2511},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 111, col: 11, offset: 2511},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 111, col: 17, offset: 2517},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 111, col: 29, offset: 2529},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 111, col: 36, offset: 2536},
								expr: &seqExpr{
									pos: position{line: 111, col: 36, offset: 2536},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 111, col: 36, offset: 2536},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 111, col: 39, offset: 2539},
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
			pos:  position{line: 124, col: 1, offset: 2892},
			expr: &choiceExpr{
				pos: position{line: 124, col: 15, offset: 2908},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 124, col: 15, offset: 2908},
						run: (*parser).callonLabeledExpr_2,
						expr: &seqExpr{
							pos: position{line: 124, col: 15, offset: 2908},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 124, col: 15, offset: 2908},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 124, col: 21, offset: 2914},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 124, col: 32, offset: 2925},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 124, col: 35, offset: 2928},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 124, col: 39, offset: 2932},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 124, col: 42, offset: 2935},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 124, col: 47, offset: 2940},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 130, col: 5, offset: 3113},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 132, col: 1, offset: 3127},
			expr: &choiceExpr{
				pos: position{line: 132, col: 16, offset: 3144},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 132, col: 16, offset: 3144},
						run: (*parser).callonPrefixedExpr_2,
						expr: &seqExpr{
							pos: position{line: 132, col: 16, offset: 3144},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 132, col: 16, offset: 3144},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 132, col: 19, offset: 3147},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 132, col: 30, offset: 3158},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 132, col: 33, offset: 3161},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 132, col: 38, offset: 3166},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 143, col: 5, offset: 3448},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 145, col: 1, offset: 3462},
			expr: &actionExpr{
				pos: position{line: 145, col: 14, offset: 3477},
				run: (*parser).callonPrefixedOp_1,
				expr: &choiceExpr{
					pos: position{line: 145, col: 16, offset: 3479},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 145, col: 16, offset: 3479},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 145, col: 22, offset: 3485},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 149, col: 1, offset: 3527},
			expr: &choiceExpr{
				pos: position{line: 149, col: 16, offset: 3544},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 149, col: 16, offset: 3544},
						run: (*parser).callonSuffixedExpr_2,
						expr: &seqExpr{
							pos: position{line: 149, col: 16, offset: 3544},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 149, col: 16, offset: 3544},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 149, col: 21, offset: 3549},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 149, col: 33, offset: 3561},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 149, col: 36, offset: 3564},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 149, col: 39, offset: 3567},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 168, col: 5, offset: 4097},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 170, col: 1, offset: 4111},
			expr: &actionExpr{
				pos: position{line: 170, col: 14, offset: 4126},
				run: (*parser).callonSuffixedOp_1,
				expr: &choiceExpr{
					pos: position{line: 170, col: 16, offset: 4128},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 170, col: 16, offset: 4128},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 170, col: 22, offset: 4134},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 170, col: 28, offset: 4140},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 174, col: 1, offset: 4182},
			expr: &choiceExpr{
				pos: position{line: 174, col: 15, offset: 4198},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 174, col: 15, offset: 4198},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 174, col: 28, offset: 4211},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 174, col: 47, offset: 4230},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 174, col: 60, offset: 4243},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 174, col: 74, offset: 4257},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 174, col: 93, offset: 4276},
						run: (*parser).callonPrimaryExpr_7,
						expr: &seqExpr{
							pos: position{line: 174, col: 93, offset: 4276},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 174, col: 93, offset: 4276},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 174, col: 97, offset: 4280},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 174, col: 100, offset: 4283},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 174, col: 105, offset: 4288},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 174, col: 116, offset: 4299},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 174, col: 119, offset: 4302},
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
			pos:  position{line: 177, col: 1, offset: 4331},
			expr: &actionExpr{
				pos: position{line: 177, col: 15, offset: 4347},
				run: (*parser).callonRuleRefExpr_1,
				expr: &seqExpr{
					pos: position{line: 177, col: 15, offset: 4347},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 177, col: 15, offset: 4347},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 177, col: 20, offset: 4352},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 177, col: 35, offset: 4367},
							expr: &seqExpr{
								pos: position{line: 177, col: 38, offset: 4370},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 177, col: 38, offset: 4370},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 177, col: 43, offset: 4375},
										expr: &seqExpr{
											pos: position{line: 177, col: 43, offset: 4375},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 177, col: 43, offset: 4375},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 177, col: 57, offset: 4389},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 177, col: 63, offset: 4395},
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
			pos:  position{line: 182, col: 1, offset: 4511},
			expr: &actionExpr{
				pos: position{line: 182, col: 20, offset: 4532},
				run: (*parser).callonSemanticPredExpr_1,
				expr: &seqExpr{
					pos: position{line: 182, col: 20, offset: 4532},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 182, col: 20, offset: 4532},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 182, col: 23, offset: 4535},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 182, col: 38, offset: 4550},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 182, col: 41, offset: 4553},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 182, col: 46, offset: 4558},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 193, col: 1, offset: 4835},
			expr: &actionExpr{
				pos: position{line: 193, col: 18, offset: 4854},
				run: (*parser).callonSemanticPredOp_1,
				expr: &choiceExpr{
					pos: position{line: 193, col: 20, offset: 4856},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 193, col: 20, offset: 4856},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 193, col: 26, offset: 4862},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 197, col: 1, offset: 4904},
			expr: &choiceExpr{
				pos: position{line: 197, col: 13, offset: 4918},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 197, col: 13, offset: 4918},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 197, col: 19, offset: 4924},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 197, col: 26, offset: 4931},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 197, col: 37, offset: 4942},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 199, col: 1, offset: 4952},
			expr: &anyMatcher{
				line: 199, col: 14, offset: 4967,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 200, col: 1, offset: 4969},
			expr: &choiceExpr{
				pos: position{line: 200, col: 11, offset: 4981},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 200, col: 11, offset: 4981},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 200, col: 30, offset: 5000},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 201, col: 1, offset: 5018},
			expr: &seqExpr{
				pos: position{line: 201, col: 20, offset: 5039},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 201, col: 20, offset: 5039},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 201, col: 27, offset: 5046},
						expr: &seqExpr{
							pos: position{line: 201, col: 27, offset: 5046},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 201, col: 27, offset: 5046},
									expr: &litMatcher{
										pos:        position{line: 201, col: 28, offset: 5047},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 201, col: 33, offset: 5052},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 201, col: 47, offset: 5066},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 202, col: 1, offset: 5071},
			expr: &seqExpr{
				pos: position{line: 202, col: 36, offset: 5108},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 202, col: 36, offset: 5108},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 202, col: 43, offset: 5115},
						expr: &seqExpr{
							pos: position{line: 202, col: 43, offset: 5115},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 202, col: 43, offset: 5115},
									expr: &choiceExpr{
										pos: position{line: 202, col: 46, offset: 5118},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 202, col: 46, offset: 5118},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 202, col: 53, offset: 5125},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 202, col: 59, offset: 5131},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 202, col: 73, offset: 5145},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 203, col: 1, offset: 5150},
			expr: &seqExpr{
				pos: position{line: 203, col: 21, offset: 5172},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 203, col: 21, offset: 5172},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 203, col: 28, offset: 5179},
						expr: &seqExpr{
							pos: position{line: 203, col: 28, offset: 5179},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 203, col: 28, offset: 5179},
									expr: &ruleRefExpr{
										pos:  position{line: 203, col: 29, offset: 5180},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 203, col: 33, offset: 5184},
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
			pos:  position{line: 205, col: 1, offset: 5199},
			expr: &ruleRefExpr{
				pos:  position{line: 205, col: 14, offset: 5214},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 206, col: 1, offset: 5229},
			expr: &actionExpr{
				pos: position{line: 206, col: 18, offset: 5248},
				run: (*parser).callonIdentifierName_1,
				expr: &seqExpr{
					pos: position{line: 206, col: 18, offset: 5248},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 206, col: 18, offset: 5248},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 206, col: 34, offset: 5264},
							expr: &ruleRefExpr{
								pos:  position{line: 206, col: 34, offset: 5264},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 209, col: 1, offset: 5346},
			expr: &charClassMatcher{
				pos:        position{line: 209, col: 19, offset: 5366},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 210, col: 1, offset: 5374},
			expr: &choiceExpr{
				pos: position{line: 210, col: 18, offset: 5393},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 210, col: 18, offset: 5393},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 210, col: 36, offset: 5411},
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
			pos:  position{line: 212, col: 1, offset: 5418},
			expr: &actionExpr{
				pos: position{line: 212, col: 14, offset: 5433},
				run: (*parser).callonLitMatcher_1,
				expr: &seqExpr{
					pos: position{line: 212, col: 14, offset: 5433},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 212, col: 14, offset: 5433},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 212, col: 18, offset: 5437},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 212, col: 32, offset: 5451},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 212, col: 39, offset: 5458},
								expr: &litMatcher{
									pos:        position{line: 212, col: 39, offset: 5458},
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
			pos:  position{line: 217, col: 1, offset: 5583},
			expr: &actionExpr{
				pos: position{line: 217, col: 17, offset: 5601},
				run: (*parser).callonStringLiteral_1,
				expr: &choiceExpr{
					pos: position{line: 217, col: 19, offset: 5603},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 217, col: 19, offset: 5603},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 217, col: 19, offset: 5603},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 217, col: 23, offset: 5607},
									expr: &ruleRefExpr{
										pos:  position{line: 217, col: 23, offset: 5607},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 217, col: 41, offset: 5625},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 217, col: 47, offset: 5631},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 217, col: 47, offset: 5631},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 217, col: 51, offset: 5635},
									name: "SingleStringChar",
								},
								&litMatcher{
									pos:        position{line: 217, col: 68, offset: 5652},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 217, col: 74, offset: 5658},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 217, col: 74, offset: 5658},
									val:        "`",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 217, col: 78, offset: 5662},
									name: "RawStringChar",
								},
								&litMatcher{
									pos:        position{line: 217, col: 92, offset: 5676},
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
			pos:  position{line: 220, col: 1, offset: 5747},
			expr: &choiceExpr{
				pos: position{line: 220, col: 20, offset: 5768},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 220, col: 20, offset: 5768},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 220, col: 20, offset: 5768},
								expr: &choiceExpr{
									pos: position{line: 220, col: 23, offset: 5771},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 220, col: 23, offset: 5771},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 220, col: 29, offset: 5777},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 220, col: 36, offset: 5784},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 220, col: 42, offset: 5790},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 220, col: 55, offset: 5803},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 220, col: 55, offset: 5803},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 220, col: 60, offset: 5808},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 221, col: 1, offset: 5827},
			expr: &choiceExpr{
				pos: position{line: 221, col: 20, offset: 5848},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 221, col: 20, offset: 5848},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 221, col: 20, offset: 5848},
								expr: &choiceExpr{
									pos: position{line: 221, col: 23, offset: 5851},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 221, col: 23, offset: 5851},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 221, col: 29, offset: 5857},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 221, col: 36, offset: 5864},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 221, col: 42, offset: 5870},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 221, col: 55, offset: 5883},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 221, col: 55, offset: 5883},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 221, col: 60, offset: 5888},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 222, col: 1, offset: 5907},
			expr: &seqExpr{
				pos: position{line: 222, col: 17, offset: 5925},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 222, col: 17, offset: 5925},
						expr: &litMatcher{
							pos:        position{line: 222, col: 18, offset: 5926},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 222, col: 22, offset: 5930},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 224, col: 1, offset: 5942},
			expr: &choiceExpr{
				pos: position{line: 224, col: 22, offset: 5965},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 224, col: 22, offset: 5965},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 224, col: 28, offset: 5971},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 225, col: 1, offset: 5992},
			expr: &choiceExpr{
				pos: position{line: 225, col: 22, offset: 6015},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 225, col: 22, offset: 6015},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 225, col: 28, offset: 6021},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 227, col: 1, offset: 6043},
			expr: &choiceExpr{
				pos: position{line: 227, col: 24, offset: 6068},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 227, col: 24, offset: 6068},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 227, col: 43, offset: 6087},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 227, col: 57, offset: 6101},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 227, col: 69, offset: 6113},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 227, col: 89, offset: 6133},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 228, col: 1, offset: 6152},
			expr: &choiceExpr{
				pos: position{line: 228, col: 20, offset: 6173},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 228, col: 20, offset: 6173},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 26, offset: 6179},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 32, offset: 6185},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 38, offset: 6191},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 44, offset: 6197},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 50, offset: 6203},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 56, offset: 6209},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 62, offset: 6215},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 229, col: 1, offset: 6220},
			expr: &seqExpr{
				pos: position{line: 229, col: 15, offset: 6236},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 229, col: 15, offset: 6236},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 26, offset: 6247},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 37, offset: 6258},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 230, col: 1, offset: 6269},
			expr: &seqExpr{
				pos: position{line: 230, col: 13, offset: 6283},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 230, col: 13, offset: 6283},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 230, col: 17, offset: 6287},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 230, col: 26, offset: 6296},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 231, col: 1, offset: 6305},
			expr: &seqExpr{
				pos: position{line: 231, col: 21, offset: 6327},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 231, col: 21, offset: 6327},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 25, offset: 6331},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 34, offset: 6340},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 43, offset: 6349},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 52, offset: 6358},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 61, offset: 6367},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 70, offset: 6376},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 79, offset: 6385},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 88, offset: 6394},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 232, col: 1, offset: 6403},
			expr: &seqExpr{
				pos: position{line: 232, col: 22, offset: 6426},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 232, col: 22, offset: 6426},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 232, col: 26, offset: 6430},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 232, col: 35, offset: 6439},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 232, col: 44, offset: 6448},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 232, col: 53, offset: 6457},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 234, col: 1, offset: 6467},
			expr: &charClassMatcher{
				pos:        position{line: 234, col: 14, offset: 6482},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 235, col: 1, offset: 6488},
			expr: &charClassMatcher{
				pos:        position{line: 235, col: 16, offset: 6505},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 236, col: 1, offset: 6511},
			expr: &charClassMatcher{
				pos:        position{line: 236, col: 12, offset: 6524},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 238, col: 1, offset: 6535},
			expr: &actionExpr{
				pos: position{line: 238, col: 20, offset: 6556},
				run: (*parser).callonCharClassMatcher_1,
				expr: &seqExpr{
					pos: position{line: 238, col: 20, offset: 6556},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 238, col: 20, offset: 6556},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 238, col: 26, offset: 6562},
							expr: &choiceExpr{
								pos: position{line: 238, col: 26, offset: 6562},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 238, col: 26, offset: 6562},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 238, col: 43, offset: 6579},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 238, col: 55, offset: 6591},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 238, col: 55, offset: 6591},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 238, col: 60, offset: 6596},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 238, col: 82, offset: 6618},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 238, col: 86, offset: 6622},
							expr: &litMatcher{
								pos:        position{line: 238, col: 86, offset: 6622},
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
			pos:  position{line: 243, col: 1, offset: 6727},
			expr: &seqExpr{
				pos: position{line: 243, col: 18, offset: 6746},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 243, col: 18, offset: 6746},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 243, col: 28, offset: 6756},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 243, col: 32, offset: 6760},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 244, col: 1, offset: 6770},
			expr: &choiceExpr{
				pos: position{line: 244, col: 13, offset: 6784},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 244, col: 13, offset: 6784},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 244, col: 13, offset: 6784},
								expr: &choiceExpr{
									pos: position{line: 244, col: 16, offset: 6787},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 244, col: 16, offset: 6787},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 244, col: 22, offset: 6793},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 244, col: 29, offset: 6800},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 244, col: 35, offset: 6806},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 244, col: 48, offset: 6819},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 244, col: 48, offset: 6819},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 244, col: 53, offset: 6824},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 245, col: 1, offset: 6840},
			expr: &choiceExpr{
				pos: position{line: 245, col: 19, offset: 6860},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 245, col: 19, offset: 6860},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 245, col: 25, offset: 6866},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 247, col: 1, offset: 6888},
			expr: &seqExpr{
				pos: position{line: 247, col: 22, offset: 6911},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 247, col: 22, offset: 6911},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 247, col: 28, offset: 6917},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 247, col: 28, offset: 6917},
								name: "SingleCharUnicodeClass",
							},
							&seqExpr{
								pos: position{line: 247, col: 53, offset: 6942},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 247, col: 53, offset: 6942},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 247, col: 57, offset: 6946},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 247, col: 70, offset: 6959},
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
			pos:  position{line: 248, col: 1, offset: 6965},
			expr: &charClassMatcher{
				pos:        position{line: 248, col: 26, offset: 6992},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 249, col: 1, offset: 7002},
			expr: &oneOrMoreExpr{
				pos: position{line: 249, col: 16, offset: 7019},
				expr: &charClassMatcher{
					pos:        position{line: 249, col: 16, offset: 7019},
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
			pos:  position{line: 251, col: 1, offset: 7029},
			expr: &actionExpr{
				pos: position{line: 251, col: 14, offset: 7044},
				run: (*parser).callonAnyMatcher_1,
				expr: &litMatcher{
					pos:        position{line: 251, col: 14, offset: 7044},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 256, col: 1, offset: 7119},
			expr: &actionExpr{
				pos: position{line: 256, col: 13, offset: 7133},
				run: (*parser).callonCodeBlock_1,
				expr: &seqExpr{
					pos: position{line: 256, col: 13, offset: 7133},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 256, col: 13, offset: 7133},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 256, col: 17, offset: 7137},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 256, col: 22, offset: 7142},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 262, col: 1, offset: 7240},
			expr: &zeroOrMoreExpr{
				pos: position{line: 262, col: 10, offset: 7251},
				expr: &choiceExpr{
					pos: position{line: 262, col: 10, offset: 7251},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 262, col: 12, offset: 7253},
							expr: &seqExpr{
								pos: position{line: 262, col: 12, offset: 7253},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 262, col: 12, offset: 7253},
										val:        "[^{}]",
										chars:      []rune{'{', '}'},
										ignoreCase: false,
										inverted:   true,
									},
									&ruleRefExpr{
										pos:  position{line: 262, col: 18, offset: 7259},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 262, col: 34, offset: 7275},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 262, col: 34, offset: 7275},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 262, col: 38, offset: 7279},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 262, col: 43, offset: 7284},
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
			pos:  position{line: 264, col: 1, offset: 7292},
			expr: &zeroOrMoreExpr{
				pos: position{line: 264, col: 8, offset: 7301},
				expr: &choiceExpr{
					pos: position{line: 264, col: 8, offset: 7301},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 264, col: 8, offset: 7301},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 264, col: 21, offset: 7314},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 264, col: 27, offset: 7320},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 265, col: 1, offset: 7331},
			expr: &zeroOrMoreExpr{
				pos: position{line: 265, col: 7, offset: 7339},
				expr: &choiceExpr{
					pos: position{line: 265, col: 7, offset: 7339},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 265, col: 7, offset: 7339},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 265, col: 20, offset: 7352},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 267, col: 1, offset: 7389},
			expr: &charClassMatcher{
				pos:        position{line: 267, col: 14, offset: 7404},
				val:        "[ \\n\\t\\r]",
				chars:      []rune{' ', '\n', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 268, col: 1, offset: 7414},
			expr: &litMatcher{
				pos:        position{line: 268, col: 7, offset: 7422},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 269, col: 1, offset: 7427},
			expr: &choiceExpr{
				pos: position{line: 269, col: 7, offset: 7435},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 269, col: 7, offset: 7435},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 269, col: 7, offset: 7435},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 269, col: 10, offset: 7438},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 269, col: 16, offset: 7444},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 269, col: 16, offset: 7444},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 269, col: 18, offset: 7446},
								expr: &ruleRefExpr{
									pos:  position{line: 269, col: 18, offset: 7446},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 269, col: 37, offset: 7465},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 269, col: 43, offset: 7471},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 269, col: 43, offset: 7471},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 269, col: 46, offset: 7474},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 270, col: 1, offset: 7478},
			expr: &notExpr{
				pos: position{line: 270, col: 7, offset: 7486},
				expr: &anyMatcher{
					line: 270, col: 8, offset: 7487,
				},
			},
		},
	},
}

func (c *current) onGrammar_1(initializer, rules interface{}) (interface{}, error) {
	pos := c.astPos()

	// create the grammar, assign its initializer
	g := ast.NewGrammar(pos)
	initSlice := initializer.([]interface{})
	if len(initSlice) > 0 {
		g.Init = initSlice[0].(*ast.CodeBlock)
	}

	rulesSlice := rules.([]interface{})
	g.Rules = make([]*ast.Rule, len(rulesSlice))
	for i, duo := range rulesSlice {
		g.Rules[i] = duo.([]interface{})[0].(*ast.Rule)
	}

	return g, nil
}

func (p *parser) callonGrammar_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar_1(stack["initializer"], stack["rules"])
}

func (c *current) onInitializer_1(code interface{}) (interface{}, error) {
	return code, nil
}

func (p *parser) callonInitializer_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInitializer_1(stack["code"])
}

func (c *current) onRule_1(name, display, expr interface{}) (interface{}, error) {
	pos := c.astPos()

	rule := ast.NewRule(pos, name.(*ast.Identifier))
	displaySlice := display.([]interface{})
	if len(displaySlice) > 0 {
		rule.DisplayName = displaySlice[0].(*ast.StringLit)
	}
	rule.Expr = expr.(ast.Expression)

	return rule, nil
}

func (p *parser) callonRule_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRule_1(stack["name"], stack["display"], stack["expr"])
}

func (c *current) onChoiceExpr_1(first, rest interface{}) (interface{}, error) {
	restSlice := rest.([]interface{})
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

func (p *parser) callonChoiceExpr_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onChoiceExpr_1(stack["first"], stack["rest"])
}

func (c *current) onActionExpr_1(expr, code interface{}) (interface{}, error) {
	if code == nil {
		return expr, nil
	}

	pos := c.astPos()
	act := ast.NewActionExpr(pos)
	act.Expr = expr.(ast.Expression)
	codeSlice := code.([]interface{})
	act.Code = codeSlice[1].(*ast.CodeBlock)

	return act, nil
}

func (p *parser) callonActionExpr_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onActionExpr_1(stack["expr"], stack["code"])
}

func (c *current) onSeqExpr_1(first, rest interface{}) (interface{}, error) {
	restSlice := rest.([]interface{})
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

func (p *parser) callonSeqExpr_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSeqExpr_1(stack["first"], stack["rest"])
}

func (c *current) onLabeledExpr_2(label, expr interface{}) (interface{}, error) {
	pos := c.astPos()
	lab := ast.NewLabeledExpr(pos)
	lab.Label = label.(*ast.Identifier)
	lab.Expr = expr.(ast.Expression)
	return lab, nil
}

func (p *parser) callonLabeledExpr_2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabeledExpr_2(stack["label"], stack["expr"])
}

func (c *current) onPrefixedExpr_2(op, expr interface{}) (interface{}, error) {
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

func (p *parser) callonPrefixedExpr_2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedExpr_2(stack["op"], stack["expr"])
}

func (c *current) onPrefixedOp_1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonPrefixedOp_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedOp_1()
}

func (c *current) onSuffixedExpr_2(expr, op interface{}) (interface{}, error) {
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

func (p *parser) callonSuffixedExpr_2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedExpr_2(stack["expr"], stack["op"])
}

func (c *current) onSuffixedOp_1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonSuffixedOp_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedOp_1()
}

func (c *current) onPrimaryExpr_7(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonPrimaryExpr_7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimaryExpr_7(stack["expr"])
}

func (c *current) onRuleRefExpr_1(name interface{}) (interface{}, error) {
	ref := ast.NewRuleRefExpr(c.astPos())
	ref.Name = name.(*ast.Identifier)
	return ref, nil
}

func (p *parser) callonRuleRefExpr_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRuleRefExpr_1(stack["name"])
}

func (c *current) onSemanticPredExpr_1(op, code interface{}) (interface{}, error) {
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

func (p *parser) callonSemanticPredExpr_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredExpr_1(stack["op"], stack["code"])
}

func (c *current) onSemanticPredOp_1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonSemanticPredOp_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredOp_1()
}

func (c *current) onIdentifierName_1() (interface{}, error) {
	return ast.NewIdentifier(c.astPos(), string(c.text)), nil
}

func (p *parser) callonIdentifierName_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierName_1()
}

func (c *current) onLitMatcher_1(lit, ignore interface{}) (interface{}, error) {
	m := ast.NewLitMatcher(c.astPos(), lit.(*ast.StringLit).Val)
	m.IgnoreCase = ignore != nil
	return m, nil
}

func (p *parser) callonLitMatcher_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLitMatcher_1(stack["lit"], stack["ignore"])
}

func (c *current) onStringLiteral_1() (interface{}, error) {
	return ast.NewStringLit(c.astPos(), string(c.text)), nil
}

func (p *parser) callonStringLiteral_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral_1()
}

func (c *current) onCharClassMatcher_1() (interface{}, error) {
	pos := c.astPos()
	cc := ast.NewCharClassMatcher(pos, string(c.text))
	return cc, nil
}

func (p *parser) callonCharClassMatcher_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher_1()
}

func (c *current) onAnyMatcher_1() (interface{}, error) {
	any := ast.NewAnyMatcher(c.astPos(), ".")
	return any, nil
}

func (p *parser) callonAnyMatcher_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnyMatcher_1()
}

func (c *current) onCodeBlock_1() (interface{}, error) {
	pos := c.astPos()
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonCodeBlock_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock_1()
}

var (
	ErrNoRule          = errors.New("grammar has no rule")
	ErrInvalidEncoding = errors.New("invalid encoding")
	ErrNoMatch         = errors.New("no match found")
)

var debug = false

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
	defer p.out(p.in("save"))
	return p.pt
}

func (p *parser) restore(pt savepoint) {
	defer p.out(p.in("restore"))
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
			defer p.out(p.in("panic handler"))
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
	defer p.out(p.in("parseRule " + rule.name))

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
	defer p.out(p.in("parseActionExpr"))

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
	defer p.out(p.in("parseAndCodeExpr"))

	ok, err := and.run(p)
	if err != nil {
		p.errs.add(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	defer p.out(p.in("parseAndExpr"))

	pt := p.save()
	_, ok := p.parseExpr(and.expr)
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	defer p.out(p.in("parseAnyMatcher"))

	if p.pt.rn != utf8.RuneError {
		p.read()
		return string(p.pt.rn), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	defer p.out(p.in("parseCharClassMatcher"))

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
	defer p.out(p.in("parseChoiceExpr"))

	for _, alt := range ch.alternatives {
		val, ok := p.parseExpr(alt)
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	defer p.out(p.in("parseLabeledExpr"))

	val, ok := p.parseExpr(lab.expr)
	if ok && lab.label != "" && len(p.vstack) > 0 {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	defer p.out(p.in("parseLitMatcher"))

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
	defer p.out(p.in("parseNotCodeExpr"))

	ok, err := not.run(p)
	if err != nil {
		p.errs.add(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	defer p.out(p.in("parseNotExpr"))

	pt := p.save()
	_, ok := p.parseExpr(not.expr)
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	defer p.out(p.in("parseOneOrMoreExpr"))

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
	defer p.out(p.in("parseRuleRefExpr " + ref.name))

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
	defer p.out(p.in("parseSeqExpr"))

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
	defer p.out(p.in("parseZeroOrMoreExpr"))

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
	defer p.out(p.in("parseZeroOrOneExpr"))

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
