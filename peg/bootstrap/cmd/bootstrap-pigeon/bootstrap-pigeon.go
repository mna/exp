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

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 51, col: 1, offset: 889},
			expr: &actionExpr{
				pos: position{line: 51, col: 11, offset: 901},
				run: (*parser).callonGrammar_1,
				expr: &seqExpr{
					pos: position{line: 51, col: 11, offset: 901},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 51, col: 11, offset: 901},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 51, col: 14, offset: 904},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 51, col: 28, offset: 918},
								expr: &seqExpr{
									pos: position{line: 51, col: 28, offset: 918},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 51, col: 28, offset: 918},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 51, col: 40, offset: 930},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 51, col: 46, offset: 936},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 51, col: 54, offset: 944},
								expr: &seqExpr{
									pos: position{line: 51, col: 54, offset: 944},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 51, col: 54, offset: 944},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 51, col: 59, offset: 949},
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
			pos:  position{line: 70, col: 1, offset: 1392},
			expr: &actionExpr{
				pos: position{line: 70, col: 15, offset: 1408},
				run: (*parser).callonInitializer_1,
				expr: &seqExpr{
					pos: position{line: 70, col: 15, offset: 1408},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 70, col: 15, offset: 1408},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 70, col: 20, offset: 1413},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 70, col: 30, offset: 1423},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 74, col: 1, offset: 1453},
			expr: &actionExpr{
				pos: position{line: 74, col: 8, offset: 1462},
				run: (*parser).callonRule_1,
				expr: &seqExpr{
					pos: position{line: 74, col: 8, offset: 1462},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 74, col: 8, offset: 1462},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 74, col: 13, offset: 1467},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 74, col: 28, offset: 1482},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 74, col: 31, offset: 1485},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 74, col: 41, offset: 1495},
								expr: &seqExpr{
									pos: position{line: 74, col: 41, offset: 1495},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 74, col: 41, offset: 1495},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 74, col: 55, offset: 1509},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 74, col: 61, offset: 1515},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 74, col: 71, offset: 1525},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 74, col: 74, offset: 1528},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 74, col: 79, offset: 1533},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 74, col: 90, offset: 1544},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 87, col: 1, offset: 1828},
			expr: &ruleRefExpr{
				pos:  position{line: 87, col: 14, offset: 1843},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 89, col: 1, offset: 1855},
			expr: &actionExpr{
				pos: position{line: 89, col: 14, offset: 1870},
				run: (*parser).callonChoiceExpr_1,
				expr: &seqExpr{
					pos: position{line: 89, col: 14, offset: 1870},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 89, col: 14, offset: 1870},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 89, col: 20, offset: 1876},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 89, col: 31, offset: 1887},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 89, col: 38, offset: 1894},
								expr: &seqExpr{
									pos: position{line: 89, col: 38, offset: 1894},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 89, col: 38, offset: 1894},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 89, col: 41, offset: 1897},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 89, col: 45, offset: 1901},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 89, col: 48, offset: 1904},
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
			pos:  position{line: 104, col: 1, offset: 2309},
			expr: &actionExpr{
				pos: position{line: 104, col: 14, offset: 2324},
				run: (*parser).callonActionExpr_1,
				expr: &seqExpr{
					pos: position{line: 104, col: 14, offset: 2324},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 104, col: 14, offset: 2324},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 19, offset: 2329},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 104, col: 27, offset: 2337},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 104, col: 34, offset: 2344},
								expr: &seqExpr{
									pos: position{line: 104, col: 34, offset: 2344},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 104, col: 34, offset: 2344},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 104, col: 37, offset: 2347},
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
			pos:  position{line: 118, col: 1, offset: 2613},
			expr: &actionExpr{
				pos: position{line: 118, col: 11, offset: 2625},
				run: (*parser).callonSeqExpr_1,
				expr: &seqExpr{
					pos: position{line: 118, col: 11, offset: 2625},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 118, col: 11, offset: 2625},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 17, offset: 2631},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 118, col: 29, offset: 2643},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 118, col: 36, offset: 2650},
								expr: &seqExpr{
									pos: position{line: 118, col: 36, offset: 2650},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 118, col: 36, offset: 2650},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 118, col: 39, offset: 2653},
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
			pos:  position{line: 131, col: 1, offset: 3004},
			expr: &choiceExpr{
				pos: position{line: 131, col: 15, offset: 3020},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 131, col: 15, offset: 3020},
						run: (*parser).callonLabeledExpr_2,
						expr: &seqExpr{
							pos: position{line: 131, col: 15, offset: 3020},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 131, col: 15, offset: 3020},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 131, col: 21, offset: 3026},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 131, col: 32, offset: 3037},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 131, col: 35, offset: 3040},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 131, col: 39, offset: 3044},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 131, col: 42, offset: 3047},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 131, col: 47, offset: 3052},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 137, col: 5, offset: 3225},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 139, col: 1, offset: 3239},
			expr: &choiceExpr{
				pos: position{line: 139, col: 16, offset: 3256},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 139, col: 16, offset: 3256},
						run: (*parser).callonPrefixedExpr_2,
						expr: &seqExpr{
							pos: position{line: 139, col: 16, offset: 3256},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 139, col: 16, offset: 3256},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 139, col: 19, offset: 3259},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 139, col: 30, offset: 3270},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 139, col: 33, offset: 3273},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 139, col: 38, offset: 3278},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 150, col: 5, offset: 3560},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 152, col: 1, offset: 3574},
			expr: &actionExpr{
				pos: position{line: 152, col: 14, offset: 3589},
				run: (*parser).callonPrefixedOp_1,
				expr: &choiceExpr{
					pos: position{line: 152, col: 16, offset: 3591},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 152, col: 16, offset: 3591},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 152, col: 22, offset: 3597},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 156, col: 1, offset: 3639},
			expr: &choiceExpr{
				pos: position{line: 156, col: 16, offset: 3656},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 156, col: 16, offset: 3656},
						run: (*parser).callonSuffixedExpr_2,
						expr: &seqExpr{
							pos: position{line: 156, col: 16, offset: 3656},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 156, col: 16, offset: 3656},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 156, col: 21, offset: 3661},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 156, col: 33, offset: 3673},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 156, col: 36, offset: 3676},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 156, col: 39, offset: 3679},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 175, col: 5, offset: 4209},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 177, col: 1, offset: 4223},
			expr: &actionExpr{
				pos: position{line: 177, col: 14, offset: 4238},
				run: (*parser).callonSuffixedOp_1,
				expr: &choiceExpr{
					pos: position{line: 177, col: 16, offset: 4240},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 177, col: 16, offset: 4240},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 177, col: 22, offset: 4246},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 177, col: 28, offset: 4252},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 181, col: 1, offset: 4294},
			expr: &choiceExpr{
				pos: position{line: 181, col: 15, offset: 4310},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 181, col: 15, offset: 4310},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 181, col: 28, offset: 4323},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 181, col: 47, offset: 4342},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 181, col: 60, offset: 4355},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 181, col: 74, offset: 4369},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 181, col: 93, offset: 4388},
						run: (*parser).callonPrimaryExpr_7,
						expr: &seqExpr{
							pos: position{line: 181, col: 93, offset: 4388},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 181, col: 93, offset: 4388},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 181, col: 97, offset: 4392},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 181, col: 100, offset: 4395},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 181, col: 105, offset: 4400},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 181, col: 116, offset: 4411},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 181, col: 119, offset: 4414},
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
			pos:  position{line: 184, col: 1, offset: 4443},
			expr: &actionExpr{
				pos: position{line: 184, col: 15, offset: 4459},
				run: (*parser).callonRuleRefExpr_1,
				expr: &seqExpr{
					pos: position{line: 184, col: 15, offset: 4459},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 184, col: 15, offset: 4459},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 184, col: 20, offset: 4464},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 184, col: 35, offset: 4479},
							expr: &seqExpr{
								pos: position{line: 184, col: 38, offset: 4482},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 184, col: 38, offset: 4482},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 184, col: 43, offset: 4487},
										expr: &seqExpr{
											pos: position{line: 184, col: 43, offset: 4487},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 184, col: 43, offset: 4487},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 184, col: 57, offset: 4501},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 184, col: 63, offset: 4507},
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
			pos:  position{line: 189, col: 1, offset: 4623},
			expr: &actionExpr{
				pos: position{line: 189, col: 20, offset: 4644},
				run: (*parser).callonSemanticPredExpr_1,
				expr: &seqExpr{
					pos: position{line: 189, col: 20, offset: 4644},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 189, col: 20, offset: 4644},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 189, col: 23, offset: 4647},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 189, col: 38, offset: 4662},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 189, col: 41, offset: 4665},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 189, col: 46, offset: 4670},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 200, col: 1, offset: 4947},
			expr: &actionExpr{
				pos: position{line: 200, col: 18, offset: 4966},
				run: (*parser).callonSemanticPredOp_1,
				expr: &choiceExpr{
					pos: position{line: 200, col: 20, offset: 4968},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 200, col: 20, offset: 4968},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 200, col: 26, offset: 4974},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 204, col: 1, offset: 5016},
			expr: &choiceExpr{
				pos: position{line: 204, col: 13, offset: 5030},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 204, col: 13, offset: 5030},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 204, col: 19, offset: 5036},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 204, col: 26, offset: 5043},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 204, col: 37, offset: 5054},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 206, col: 1, offset: 5064},
			expr: &anyMatcher{
				line: 206, col: 14, offset: 5079,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 207, col: 1, offset: 5081},
			expr: &choiceExpr{
				pos: position{line: 207, col: 11, offset: 5093},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 207, col: 11, offset: 5093},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 207, col: 30, offset: 5112},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 208, col: 1, offset: 5130},
			expr: &seqExpr{
				pos: position{line: 208, col: 20, offset: 5151},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 208, col: 20, offset: 5151},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 208, col: 27, offset: 5158},
						expr: &seqExpr{
							pos: position{line: 208, col: 27, offset: 5158},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 208, col: 27, offset: 5158},
									expr: &litMatcher{
										pos:        position{line: 208, col: 28, offset: 5159},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 208, col: 33, offset: 5164},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 208, col: 47, offset: 5178},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 209, col: 1, offset: 5183},
			expr: &seqExpr{
				pos: position{line: 209, col: 36, offset: 5220},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 209, col: 36, offset: 5220},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 209, col: 43, offset: 5227},
						expr: &seqExpr{
							pos: position{line: 209, col: 43, offset: 5227},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 209, col: 43, offset: 5227},
									expr: &choiceExpr{
										pos: position{line: 209, col: 46, offset: 5230},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 209, col: 46, offset: 5230},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 209, col: 53, offset: 5237},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 209, col: 59, offset: 5243},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 209, col: 73, offset: 5257},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 210, col: 1, offset: 5262},
			expr: &seqExpr{
				pos: position{line: 210, col: 21, offset: 5284},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 210, col: 21, offset: 5284},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 210, col: 28, offset: 5291},
						expr: &seqExpr{
							pos: position{line: 210, col: 28, offset: 5291},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 210, col: 28, offset: 5291},
									expr: &ruleRefExpr{
										pos:  position{line: 210, col: 29, offset: 5292},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 210, col: 33, offset: 5296},
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
			pos:  position{line: 212, col: 1, offset: 5311},
			expr: &ruleRefExpr{
				pos:  position{line: 212, col: 14, offset: 5326},
				name: "IdentifierName",
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 213, col: 1, offset: 5341},
			expr: &actionExpr{
				pos: position{line: 213, col: 18, offset: 5360},
				run: (*parser).callonIdentifierName_1,
				expr: &seqExpr{
					pos: position{line: 213, col: 18, offset: 5360},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 213, col: 18, offset: 5360},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 213, col: 34, offset: 5376},
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 34, offset: 5376},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 216, col: 1, offset: 5458},
			expr: &charClassMatcher{
				pos:        position{line: 216, col: 19, offset: 5478},
				val:        "[a-z_]i",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 217, col: 1, offset: 5486},
			expr: &choiceExpr{
				pos: position{line: 217, col: 18, offset: 5505},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 217, col: 18, offset: 5505},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 217, col: 36, offset: 5523},
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
			pos:  position{line: 219, col: 1, offset: 5530},
			expr: &actionExpr{
				pos: position{line: 219, col: 14, offset: 5545},
				run: (*parser).callonLitMatcher_1,
				expr: &seqExpr{
					pos: position{line: 219, col: 14, offset: 5545},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 219, col: 14, offset: 5545},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 219, col: 18, offset: 5549},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 219, col: 32, offset: 5563},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 219, col: 39, offset: 5570},
								expr: &litMatcher{
									pos:        position{line: 219, col: 39, offset: 5570},
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
			pos:  position{line: 224, col: 1, offset: 5695},
			expr: &actionExpr{
				pos: position{line: 224, col: 17, offset: 5713},
				run: (*parser).callonStringLiteral_1,
				expr: &choiceExpr{
					pos: position{line: 224, col: 19, offset: 5715},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 224, col: 19, offset: 5715},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 224, col: 19, offset: 5715},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 224, col: 23, offset: 5719},
									expr: &ruleRefExpr{
										pos:  position{line: 224, col: 23, offset: 5719},
										name: "DoubleStringChar",
									},
								},
								&litMatcher{
									pos:        position{line: 224, col: 41, offset: 5737},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 224, col: 47, offset: 5743},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 224, col: 47, offset: 5743},
									val:        "'",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 224, col: 51, offset: 5747},
									name: "SingleStringChar",
								},
								&litMatcher{
									pos:        position{line: 224, col: 68, offset: 5764},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 224, col: 74, offset: 5770},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 224, col: 74, offset: 5770},
									val:        "`",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 224, col: 78, offset: 5774},
									name: "RawStringChar",
								},
								&litMatcher{
									pos:        position{line: 224, col: 92, offset: 5788},
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
			pos:  position{line: 227, col: 1, offset: 5859},
			expr: &choiceExpr{
				pos: position{line: 227, col: 20, offset: 5880},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 227, col: 20, offset: 5880},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 227, col: 20, offset: 5880},
								expr: &choiceExpr{
									pos: position{line: 227, col: 23, offset: 5883},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 227, col: 23, offset: 5883},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 227, col: 29, offset: 5889},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 227, col: 36, offset: 5896},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 227, col: 42, offset: 5902},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 227, col: 55, offset: 5915},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 227, col: 55, offset: 5915},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 227, col: 60, offset: 5920},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 228, col: 1, offset: 5939},
			expr: &choiceExpr{
				pos: position{line: 228, col: 20, offset: 5960},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 228, col: 20, offset: 5960},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 228, col: 20, offset: 5960},
								expr: &choiceExpr{
									pos: position{line: 228, col: 23, offset: 5963},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 228, col: 23, offset: 5963},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 228, col: 29, offset: 5969},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 228, col: 36, offset: 5976},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 228, col: 42, offset: 5982},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 228, col: 55, offset: 5995},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 228, col: 55, offset: 5995},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 228, col: 60, offset: 6000},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 229, col: 1, offset: 6019},
			expr: &seqExpr{
				pos: position{line: 229, col: 17, offset: 6037},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 229, col: 17, offset: 6037},
						expr: &litMatcher{
							pos:        position{line: 229, col: 18, offset: 6038},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 22, offset: 6042},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 231, col: 1, offset: 6054},
			expr: &choiceExpr{
				pos: position{line: 231, col: 22, offset: 6077},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 231, col: 22, offset: 6077},
						val:        "'",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 28, offset: 6083},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 232, col: 1, offset: 6104},
			expr: &choiceExpr{
				pos: position{line: 232, col: 22, offset: 6127},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 232, col: 22, offset: 6127},
						val:        "\"",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 232, col: 28, offset: 6133},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 234, col: 1, offset: 6155},
			expr: &choiceExpr{
				pos: position{line: 234, col: 24, offset: 6180},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 234, col: 24, offset: 6180},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 234, col: 43, offset: 6199},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 234, col: 57, offset: 6213},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 234, col: 69, offset: 6225},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 234, col: 89, offset: 6245},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 235, col: 1, offset: 6264},
			expr: &choiceExpr{
				pos: position{line: 235, col: 20, offset: 6285},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 235, col: 20, offset: 6285},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 26, offset: 6291},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 32, offset: 6297},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 38, offset: 6303},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 44, offset: 6309},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 50, offset: 6315},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 56, offset: 6321},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 235, col: 62, offset: 6327},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 236, col: 1, offset: 6332},
			expr: &seqExpr{
				pos: position{line: 236, col: 15, offset: 6348},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 236, col: 15, offset: 6348},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 236, col: 26, offset: 6359},
						name: "OctalDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 236, col: 37, offset: 6370},
						name: "OctalDigit",
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 237, col: 1, offset: 6381},
			expr: &seqExpr{
				pos: position{line: 237, col: 13, offset: 6395},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 237, col: 13, offset: 6395},
						val:        "x",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 237, col: 17, offset: 6399},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 237, col: 26, offset: 6408},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 238, col: 1, offset: 6417},
			expr: &seqExpr{
				pos: position{line: 238, col: 21, offset: 6439},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 238, col: 21, offset: 6439},
						val:        "U",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 25, offset: 6443},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 34, offset: 6452},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 43, offset: 6461},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 52, offset: 6470},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 61, offset: 6479},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 70, offset: 6488},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 79, offset: 6497},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 238, col: 88, offset: 6506},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 239, col: 1, offset: 6515},
			expr: &seqExpr{
				pos: position{line: 239, col: 22, offset: 6538},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 239, col: 22, offset: 6538},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 239, col: 26, offset: 6542},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 239, col: 35, offset: 6551},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 239, col: 44, offset: 6560},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 239, col: 53, offset: 6569},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 241, col: 1, offset: 6579},
			expr: &charClassMatcher{
				pos:        position{line: 241, col: 14, offset: 6594},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 242, col: 1, offset: 6600},
			expr: &charClassMatcher{
				pos:        position{line: 242, col: 16, offset: 6617},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 243, col: 1, offset: 6623},
			expr: &charClassMatcher{
				pos:        position{line: 243, col: 12, offset: 6636},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 245, col: 1, offset: 6647},
			expr: &actionExpr{
				pos: position{line: 245, col: 20, offset: 6668},
				run: (*parser).callonCharClassMatcher_1,
				expr: &seqExpr{
					pos: position{line: 245, col: 20, offset: 6668},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 245, col: 20, offset: 6668},
							val:        "[",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 245, col: 26, offset: 6674},
							expr: &choiceExpr{
								pos: position{line: 245, col: 26, offset: 6674},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 245, col: 26, offset: 6674},
										name: "ClassCharRange",
									},
									&ruleRefExpr{
										pos:  position{line: 245, col: 43, offset: 6691},
										name: "ClassChar",
									},
									&seqExpr{
										pos: position{line: 245, col: 55, offset: 6703},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 245, col: 55, offset: 6703},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 245, col: 60, offset: 6708},
												name: "UnicodeClassEscape",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 245, col: 82, offset: 6730},
							val:        "]",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 245, col: 86, offset: 6734},
							expr: &litMatcher{
								pos:        position{line: 245, col: 86, offset: 6734},
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
			pos:  position{line: 250, col: 1, offset: 6839},
			expr: &seqExpr{
				pos: position{line: 250, col: 18, offset: 6858},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 250, col: 18, offset: 6858},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 250, col: 28, offset: 6868},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 250, col: 32, offset: 6872},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 251, col: 1, offset: 6882},
			expr: &choiceExpr{
				pos: position{line: 251, col: 13, offset: 6896},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 251, col: 13, offset: 6896},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 251, col: 13, offset: 6896},
								expr: &choiceExpr{
									pos: position{line: 251, col: 16, offset: 6899},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 251, col: 16, offset: 6899},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 251, col: 22, offset: 6905},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 251, col: 29, offset: 6912},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 251, col: 35, offset: 6918},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 251, col: 48, offset: 6931},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 251, col: 48, offset: 6931},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 251, col: 53, offset: 6936},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 252, col: 1, offset: 6952},
			expr: &choiceExpr{
				pos: position{line: 252, col: 19, offset: 6972},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 252, col: 19, offset: 6972},
						val:        "]",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 252, col: 25, offset: 6978},
						name: "CommonEscapeSequence",
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 254, col: 1, offset: 7000},
			expr: &seqExpr{
				pos: position{line: 254, col: 22, offset: 7023},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 254, col: 22, offset: 7023},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 254, col: 28, offset: 7029},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 254, col: 28, offset: 7029},
								name: "SingleCharUnicodeClass",
							},
							&seqExpr{
								pos: position{line: 254, col: 53, offset: 7054},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 254, col: 53, offset: 7054},
										val:        "{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 254, col: 57, offset: 7058},
										name: "UnicodeClass",
									},
									&litMatcher{
										pos:        position{line: 254, col: 70, offset: 7071},
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
			pos:  position{line: 255, col: 1, offset: 7077},
			expr: &charClassMatcher{
				pos:        position{line: 255, col: 26, offset: 7104},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeClass",
			pos:  position{line: 256, col: 1, offset: 7114},
			expr: &oneOrMoreExpr{
				pos: position{line: 256, col: 16, offset: 7131},
				expr: &charClassMatcher{
					pos:        position{line: 256, col: 16, offset: 7131},
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
			pos:  position{line: 258, col: 1, offset: 7141},
			expr: &actionExpr{
				pos: position{line: 258, col: 14, offset: 7156},
				run: (*parser).callonAnyMatcher_1,
				expr: &litMatcher{
					pos:        position{line: 258, col: 14, offset: 7156},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 263, col: 1, offset: 7231},
			expr: &actionExpr{
				pos: position{line: 263, col: 13, offset: 7245},
				run: (*parser).callonCodeBlock_1,
				expr: &seqExpr{
					pos: position{line: 263, col: 13, offset: 7245},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 263, col: 13, offset: 7245},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 263, col: 17, offset: 7249},
							name: "Code",
						},
						&litMatcher{
							pos:        position{line: 263, col: 22, offset: 7254},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 269, col: 1, offset: 7352},
			expr: &zeroOrMoreExpr{
				pos: position{line: 269, col: 10, offset: 7363},
				expr: &choiceExpr{
					pos: position{line: 269, col: 10, offset: 7363},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 269, col: 12, offset: 7365},
							expr: &seqExpr{
								pos: position{line: 269, col: 12, offset: 7365},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 269, col: 12, offset: 7365},
										val:        "[^{}]",
										chars:      []rune{'{', '}'},
										ignoreCase: false,
										inverted:   true,
									},
									&ruleRefExpr{
										pos:  position{line: 269, col: 18, offset: 7371},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 269, col: 34, offset: 7387},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 269, col: 34, offset: 7387},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 269, col: 38, offset: 7391},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 269, col: 43, offset: 7396},
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
			pos:  position{line: 271, col: 1, offset: 7404},
			expr: &zeroOrMoreExpr{
				pos: position{line: 271, col: 8, offset: 7413},
				expr: &choiceExpr{
					pos: position{line: 271, col: 8, offset: 7413},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 271, col: 8, offset: 7413},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 271, col: 21, offset: 7426},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 271, col: 27, offset: 7432},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 272, col: 1, offset: 7443},
			expr: &zeroOrMoreExpr{
				pos: position{line: 272, col: 7, offset: 7451},
				expr: &choiceExpr{
					pos: position{line: 272, col: 7, offset: 7451},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 272, col: 7, offset: 7451},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 272, col: 20, offset: 7464},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 274, col: 1, offset: 7501},
			expr: &charClassMatcher{
				pos:        position{line: 274, col: 14, offset: 7516},
				val:        "[ \\n\\t\\r]",
				chars:      []rune{' ', '\n', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 275, col: 1, offset: 7526},
			expr: &litMatcher{
				pos:        position{line: 275, col: 7, offset: 7534},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 276, col: 1, offset: 7539},
			expr: &choiceExpr{
				pos: position{line: 276, col: 7, offset: 7547},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 276, col: 7, offset: 7547},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 276, col: 7, offset: 7547},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 276, col: 10, offset: 7550},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 276, col: 16, offset: 7556},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 276, col: 16, offset: 7556},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 276, col: 18, offset: 7558},
								expr: &ruleRefExpr{
									pos:  position{line: 276, col: 18, offset: 7558},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 276, col: 37, offset: 7577},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 276, col: 43, offset: 7583},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 276, col: 43, offset: 7583},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 276, col: 46, offset: 7586},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 277, col: 1, offset: 7590},
			expr: &notExpr{
				pos: position{line: 277, col: 7, offset: 7598},
				expr: &anyMatcher{
					line: 277, col: 8, offset: 7599,
				},
			},
		},
	},
}

func (c *current) onGrammar_1(initializer, rules interface{}) (interface{}, error) {
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
	displaySlice := toIfaceSlice(display)
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
	codeSlice := toIfaceSlice(code)
	act.Code = codeSlice[1].(*ast.CodeBlock)

	return act, nil
}

func (p *parser) callonActionExpr_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onActionExpr_1(stack["expr"], stack["code"])
}

func (c *current) onSeqExpr_1(first, rest interface{}) (interface{}, error) {
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
