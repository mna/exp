// Package builder generates the parser code for a given grammar. It makes
// no attempt to verify the correctness of the grammar.
package builder

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/PuerkitoBio/exp/peg/ast"
)

var funcTemplate = `func %s(%s) (interface{}, error) {
%s
}
`

var requiredImports = []string{
	"bytes",
	"errors",
	"fmt",
	"io",
	"io/ioutil",
	"os",
	"strings",
	"unicode",
	"unicode/utf8",
}

func BuildParser(w io.Writer, g *ast.Grammar, imports ...string) error {
	b := &builder{w: w}
	return b.buildParser(g, append(requiredImports, imports...))
}

type builder struct {
	w   io.Writer
	err error

	ruleName  string
	exprIndex int
	argsStack [][]string
}

func (b *builder) buildParser(g *ast.Grammar, imports []string) error {
	b.writePackageAndImports(g.Package, imports)
	b.writeInit(g.Init)
	b.writeGrammar(g)

	for _, rule := range g.Rules {
		b.writeRuleCode(rule)
	}

	return b.err
}

func (b *builder) writePackageAndImports(pkg *ast.Package, imports []string) {
	if pkg == nil {
		return
	}
	b.writelnf("package %s", pkg.Name.Val)
	b.writelnf("import (")
	for _, imp := range imports {
		b.writelnf("\t%q", imp)
	}
	b.writelnf(")")
}

func (b *builder) writeInit(init *ast.CodeBlock) {
	if init == nil {
		return
	}

	// remove opening and closing braces
	val := init.Val[1 : len(init.Val)-1]
	b.writelnf("%s", val)
}

func (b *builder) writeGrammar(g *ast.Grammar) {
	// transform the ast grammar to the self-contained, no dependency version
	// of the parser-generator grammar.
	b.writelnf("var g = &grammar {")
	b.writelnf("\trules = []*rule{")
	for _, r := range g.Rules {
		b.writeRule(r)
	}
	b.writelnf("\t},")
	b.writelnf("}")
}

func (b *builder) writeRule(r *ast.Rule) {
	if r == nil || r.Name == nil {
		return
	}

	b.exprIndex = 0
	b.ruleName = r.Name.Val

	b.writelnf("{")
	b.writelnf("\tname: %q,", r.Name.Val)
	if r.DisplayName != nil && r.DisplayName.Val != "" {
		b.writelnf("\tdisplayName: %q,", r.DisplayName.Val)
	}
	pos := r.Pos()
	b.writelnf("\tpos: position{line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(r.Expr)
	b.writelnf("},")
}

func (b *builder) writeExpr(expr ast.Expression) {
	b.exprIndex++
	switch expr := expr.(type) {
	case *ast.ActionExpr:
		b.writeActionExpr(expr)
	case *ast.AndCodeExpr:
		b.writeAndCodeExpr(expr)
	case *ast.AndExpr:
		b.writeAndExpr(expr)
	case *ast.AnyMatcher:
	case *ast.CharClassMatcher:
	case *ast.ChoiceExpr:
	case *ast.LabeledExpr:
	case *ast.LitMatcher:
	case *ast.NotCodeExpr:
	case *ast.NotExpr:
	case *ast.OneOrMoreExpr:
	case *ast.RuleRefExpr:
	case *ast.SeqExpr:
	case *ast.ZeroOrMoreExpr:
	case *ast.ZeroOrOneExpr:
	default:
		b.err = fmt.Errorf("builder: unknown expression type %T", expr)
	}
}

func (b *builder) writeActionExpr(act *ast.ActionExpr) {
	if act == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("&actionExpr{")
	pos := act.Pos()
	b.writelnf("\tpos: position{line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: (*parser).callon%s_%d,", b.ruleName, b.exprIndex)
	b.writef("\texpr: ")
	b.writeExpr(act.Expr)
	b.writelnf("},")
}

func (b *builder) writeAndCodeExpr(and *ast.AndCodeExpr) {
	if and == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("&andCodeExpr{")
	pos := and.Pos()
	b.writelnf("\tpos: position{line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: (*parser).callon%s_%d,", b.ruleName, b.exprIndex)
	b.writelnf("},")
}

func (b *builder) writeRuleCode(rule *ast.Rule) {
	if rule == nil || rule.Name == nil {
		return
	}

	// keep trace of the current rule, as the code blocks are created
	// in functions named "on<RuleName><#ExprIndex>".
	b.ruleName = rule.Name.Val
	b.exprIndex = 0
	b.writeExprCode(rule.Expr)
}

func (b *builder) pushArgsSet() {
	b.argsStack = append(b.argsStack, nil)
}

func (b *builder) popArgsSet() {
	b.argsStack = b.argsStack[:len(b.argsStack)-1]
}

func (b *builder) addArg(arg *ast.Identifier) {
	if arg == nil {
		return
	}
	ix := len(b.argsStack) - 1
	b.argsStack[ix] = append(b.argsStack[ix], arg.Val)
}

func (b *builder) writeExprCode(expr ast.Expression) {
	b.exprIndex++
	switch expr := expr.(type) {
	case *ast.ActionExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.writeActionExprCode(expr)
		b.popArgsSet()

	case *ast.AndCodeExpr:
		// TODO : should be able to access labeled vars too, but when to
		// start a new args set?
		b.writeAndCodeExprCode(expr)

	case *ast.LabeledExpr:
		b.addArg(expr.Label)
		b.writeExprCode(expr.Expr)

	case *ast.NotCodeExpr:
		// TODO : should be able to access labeled vars too, but when to
		// start a new args set?
		b.writeNotCodeExprCode(expr)

	case *ast.AndExpr:
		b.writeExprCode(expr.Expr)
	case *ast.ChoiceExpr:
		for _, alt := range expr.Alternatives {
			b.writeExprCode(alt)
		}
	case *ast.NotExpr:
		b.writeExprCode(expr.Expr)
	case *ast.OneOrMoreExpr:
		b.writeExprCode(expr.Expr)
	case *ast.SeqExpr:
		for _, sub := range expr.Exprs {
			b.writeExprCode(sub)
		}
	case *ast.ZeroOrMoreExpr:
		b.writeExprCode(expr.Expr)
	case *ast.ZeroOrOneExpr:
		b.writeExprCode(expr.Expr)
	}
}

func (b *builder) writeActionExprCode(act *ast.ActionExpr) {
	if act == nil {
		return
	}
	b.writeFunc(act.Code)
}

func (b *builder) writeAndCodeExprCode(and *ast.AndCodeExpr) {
	if and == nil {
		return
	}
	b.writeFunc(and.Code)
}

func (b *builder) writeNotCodeExprCode(not *ast.NotCodeExpr) {
	if not == nil {
		return
	}
	b.writeFunc(not.Code)
}

func (b *builder) writeFunc(code *ast.CodeBlock) {
	if code == nil {
		return
	}
	val := code.Val[1 : len(code.Val)-1]
	if val[0] == '\n' {
		val = val[1:]
	}
	if val[len(val)-1] == '\n' {
		val = val[:len(val)-1]
	}
	var args bytes.Buffer
	ix := len(b.argsStack) - 1
	for i, arg := range b.argsStack[ix] {
		if i > 0 {
			args.WriteString(", ")
		}
		args.WriteString(arg)
	}
	if args.Len() > 0 {
		args.WriteString(" interface{}")
	}

	b.writelnf(funcTemplate, b.funcName(), args.String(), val)
}

func (b *builder) funcName() string {
	return "on" + b.ruleName + "_" + strconv.Itoa(b.exprIndex)
}

func (b *builder) writef(f string, args ...interface{}) {
	if b.err == nil {
		_, b.err = fmt.Fprintf(b.w, f, args...)
	}
}

func (b *builder) writelnf(f string, args ...interface{}) {
	b.writef(f+"\n", args...)
}
