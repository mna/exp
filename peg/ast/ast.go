package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Pos represents a position in a source file.
type Pos struct {
	Filename string
	Line     int
	Col      int
	Off      int
}

// String returns the textual representation of a position.
func (p Pos) String() string {
	if p.Filename != "" {
		return fmt.Sprintf("%s:%d:%d (%d)", p.Filename, p.Line, p.Col, p.Off)
	}
	return fmt.Sprintf("%d:%d (%d)", p.Line, p.Col, p.Off)
}

type Grammar struct {
	p       Pos        // identical to Grammar.Package.Pos()
	Package *Package   // package declaration
	Init    *CodeBlock // initializer code block
	Rules   []*Rule    // all rules
}

func NewGrammar(p Pos, pkg *Package) *Grammar {
	return &Grammar{p: p, Package: pkg}
}

func (g *Grammar) Pos() Pos { return g.p }

func (g *Grammar) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s: %T{Package: %v, Init: %v, Rules: [\n",
		g.p, g, g.Package, g.Init))
	for _, r := range g.Rules {
		buf.WriteString(fmt.Sprintf("%s,\n", r))
	}
	buf.WriteString("]}")
	return buf.String()
}

type Rule struct {
	p           Pos // identical to Rule.Name.Pos()
	Name        *Identifier
	DisplayName *StringLit
	Expr        Expression
}

func NewRule(p Pos, name *Identifier) *Rule {
	return &Rule{p: p, Name: name}
}

func (r *Rule) Pos() Pos { return r.p }

func (r *Rule) String() string {
	return fmt.Sprintf("%s: %T{Name: %v, DisplayName: %v, Expr: %v}",
		r.p, r, r.Name, r.DisplayName, r.Expr)
}

type Expression interface {
	Pos() Pos
}

type ChoiceExpr struct {
	p            Pos
	Alternatives []Expression
}

func NewChoiceExpr(p Pos) *ChoiceExpr {
	return &ChoiceExpr{p: p}
}

func (c *ChoiceExpr) Pos() Pos { return c.p }

func (c *ChoiceExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s: %T{Alternatives: [\n", c.p, c))
	for _, e := range c.Alternatives {
		buf.WriteString(fmt.Sprintf("%s,\n", e))
	}
	buf.WriteString("]}")
	return buf.String()
}

type ActionExpr struct {
	p    Pos
	Expr Expression
	Code *CodeBlock
}

func NewActionExpr(p Pos) *ActionExpr {
	return &ActionExpr{p: p}
}

func (a *ActionExpr) Pos() Pos { return a.p }

func (a *ActionExpr) String() string {
	return fmt.Sprintf("%s: %T{Expr: %v, Code: %v}", a.p, a, a.Expr, a.Code)
}

type SeqExpr struct {
	p     Pos
	Exprs []Expression
}

func NewSeqExpr(p Pos) *SeqExpr {
	return &SeqExpr{p: p}
}

func (s *SeqExpr) Pos() Pos { return s.p }

func (s *SeqExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s: %T{Exprs: [\n", s.p, s))
	for _, e := range s.Exprs {
		buf.WriteString(fmt.Sprintf("%s,\n", e))
	}
	buf.WriteString("]}")
	return buf.String()
}

type LabeledExpr struct {
	p     Pos
	Label *Identifier
	Expr  Expression
}

func NewLabeledExpr(p Pos) *LabeledExpr {
	return &LabeledExpr{p: p}
}

func (l *LabeledExpr) Pos() Pos { return l.p }

func (l *LabeledExpr) String() string {
	return fmt.Sprintf("%s: %T{Label: %v, Expr: %v}", l.p, l, l.Label, l.Expr)
}

type AndExpr struct {
	p    Pos
	Expr Expression
}

func NewAndExpr(p Pos) *AndExpr {
	return &AndExpr{p: p}
}

func (a *AndExpr) Pos() Pos { return a.p }

func (a *AndExpr) String() string {
	return fmt.Sprintf("%s: %T{Expr: %v}", a.p, a, a.Expr)
}

type NotExpr struct {
	p    Pos
	Expr Expression
}

func NewNotExpr(p Pos) *NotExpr {
	return &NotExpr{p: p}
}

func (n *NotExpr) Pos() Pos { return n.p }

func (n *NotExpr) String() string {
	return fmt.Sprintf("%s: %T{Expr: %v}", n.p, n, n.Expr)
}

type ZeroOrOneExpr struct {
	p    Pos
	Expr Expression
}

func NewZeroOrOneExpr(p Pos) *ZeroOrOneExpr {
	return &ZeroOrOneExpr{p: p}
}

func (z *ZeroOrOneExpr) Pos() Pos { return z.p }

func (z *ZeroOrOneExpr) String() string {
	return fmt.Sprintf("%s: %T{Expr: %v}", z.p, z, z.Expr)
}

type ZeroOrMoreExpr struct {
	p    Pos
	Expr Expression
}

func NewZeroOrMoreExpr(p Pos) *ZeroOrMoreExpr {
	return &ZeroOrMoreExpr{p: p}
}

func (z *ZeroOrMoreExpr) Pos() Pos { return z.p }

func (z *ZeroOrMoreExpr) String() string {
	return fmt.Sprintf("%s: %T{Expr: %v}", z.p, z, z.Expr)
}

type OneOrMoreExpr struct {
	p    Pos
	Expr Expression
}

func NewOneOrMoreExpr(p Pos) *OneOrMoreExpr {
	return &OneOrMoreExpr{p: p}
}

func (o *OneOrMoreExpr) Pos() Pos { return o.p }

func (o *OneOrMoreExpr) String() string {
	return fmt.Sprintf("%s: %T{Expr: %v}", o.p, o, o.Expr)
}

type RuleRefExpr struct {
	p    Pos
	Name *Identifier
}

func NewRuleRefExpr(p Pos) *RuleRefExpr {
	return &RuleRefExpr{p: p}
}

func (r *RuleRefExpr) Pos() Pos { return r.p }

func (r *RuleRefExpr) String() string {
	return fmt.Sprintf("%s: %T{Name: %v}", r.p, r, r.Name)
}

type AndCodeExpr struct {
	p    Pos
	Code *CodeBlock
}

func NewAndCodeExpr(p Pos) *AndCodeExpr {
	return &AndCodeExpr{p: p}
}

func (a *AndCodeExpr) Pos() Pos { return a.p }

func (a *AndCodeExpr) String() string {
	return fmt.Sprintf("%s: %T{Code: %v}", a.p, a, a.Code)
}

type NotCodeExpr struct {
	p    Pos
	Code *CodeBlock
}

func NewNotCodeExpr(p Pos) *NotCodeExpr {
	return &NotCodeExpr{p: p}
}

func (n *NotCodeExpr) Pos() Pos { return n.p }

func (n *NotCodeExpr) String() string {
	return fmt.Sprintf("%s: %T{Code: %v}", n.p, n, n.Code)
}

type LitMatcher struct {
	posValue   // can be str, rstr or char
	IgnoreCase bool
}

func NewLitMatcher(p Pos, v string) *LitMatcher {
	return &LitMatcher{posValue: posValue{p: p, Val: v}}
}

func (l *LitMatcher) Pos() Pos { return l.p }

func (l *LitMatcher) String() string {
	return fmt.Sprintf("%s: %T{Val: %q, IgnoreCase: %t}", l.p, l, l.Val, l.IgnoreCase)
}

type CharClassMatcher struct {
	posValue
	IgnoreCase     bool
	Inverted       bool
	Chars          []rune
	Ranges         []rune // pairs of low/high range
	UnicodeClasses []string
}

func NewCharClassMatcher(p Pos, raw string) *CharClassMatcher {
	c := &CharClassMatcher{posValue: posValue{p: p, Val: raw}}
	c.parse()
	return c
}

func (c *CharClassMatcher) parse() {
	raw := c.Val
	c.IgnoreCase = strings.HasSuffix(raw, "i")
	if c.IgnoreCase {
		raw = raw[:len(raw)-1]
	}

	// "unquote" the character classes
	raw = raw[1 : len(raw)-1]
	if len(raw) == 0 {
		return
	}

	c.Inverted = raw[0] == '^'
	if c.Inverted {
		raw = raw[1:]
		if len(raw) == 0 {
			return
		}
	}

	// content of char class is necessarily valid, so escapes are correct
	r := strings.NewReader(raw)
outer:
	for {
		rn, _, err := r.ReadRune()
		if err != nil {
			break outer
		}

		switch rn {
		case '\\':
			rn, _, _ := r.ReadRune()
			switch rn {
			case ']':
				c.Chars = append(c.Chars, rn)
			case 'p':
				rn, _, _ := r.ReadRune()
				if rn == '{' {
					var class bytes.Buffer
					for {
						rn, _, _ := r.ReadRune()
						if rn == '}' {
							break
						}
						class.WriteRune(rn)
					}
					c.UnicodeClasses = append(c.UnicodeClasses, class.String())
				} else {
					c.UnicodeClasses = append(c.UnicodeClasses, string(rn))
				}
			default:
				rn, _, _, _ := strconv.UnquoteChar("\\"+string(rn), 0)
				c.Chars = append(c.Chars, rn)
			}
			// TODO : Implement range
		default:
			if c.IgnoreCase {
				rn = unicode.ToLower(rn)
			}
			c.Chars = append(c.Chars, rn)
		}
	}
}

func (c *CharClassMatcher) Pos() Pos { return c.p }

func (c *CharClassMatcher) String() string {
	return fmt.Sprintf("%s: %T{Val: %q, IgnoreCase: %t, Inverted: %t}",
		c.p, c, c.Val, c.IgnoreCase, c.Inverted)
}

type AnyMatcher struct {
	posValue
}

func NewAnyMatcher(p Pos, v string) *AnyMatcher {
	return &AnyMatcher{posValue{p, v}}
}

func (a *AnyMatcher) Pos() Pos { return a.p }

func (a *AnyMatcher) String() string {
	return fmt.Sprintf("%s: %T{Val: %q}", a.p, a, a.Val)
}

type Package struct {
	p    Pos // starting pos of the package keyword
	Name *Identifier
}

func NewPackage(p Pos) *Package {
	return &Package{p: p}
}

func (p *Package) Pos() Pos { return p.p }

func (p *Package) String() string {
	return fmt.Sprintf("%s: %T{Name: %v}", p.p, p, p.Name)
}

type CodeBlock struct {
	posValue
}

func NewCodeBlock(p Pos, code string) *CodeBlock {
	return &CodeBlock{posValue{p, code}}
}

func (c *CodeBlock) Pos() Pos { return c.p }

func (c *CodeBlock) String() string {
	return fmt.Sprintf("%s: %T{Val: %q}", c.p, c, c.Val)
}

type Identifier struct {
	posValue
}

func NewIdentifier(p Pos, name string) *Identifier {
	return &Identifier{posValue{p: p, Val: name}}
}

func (i *Identifier) Pos() Pos { return i.p }

func (i *Identifier) String() string {
	return fmt.Sprintf("%s: %T{Val: %q}", i.p, i, i.Val)
}

type StringLit struct {
	posValue
}

func NewStringLit(p Pos, val string) *StringLit {
	return &StringLit{posValue{p: p, Val: val}}
}

func (s *StringLit) Pos() Pos { return s.p }

func (s *StringLit) String() string {
	return fmt.Sprintf("%s: %T{Val: %q}", s.p, s, s.Val)
}

type posValue struct {
	p   Pos
	Val string
}
