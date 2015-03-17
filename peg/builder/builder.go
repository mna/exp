package builder

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/exp/peg/ast"
)

func BuildParser(w io.Writer, g *ast.Grammar) error {
	b := &builder{w: w}
	return b.buildParser(g)
}

type builder struct {
	w   io.Writer
	err error
}

func (b *builder) buildParser(g *ast.Grammar) error {
	b.writePackage(g.Package)
	b.writeInit(g.Init)
	return b.err
}

func (b *builder) writePackage(pkg *ast.Package) {
	if pkg == nil {
		return
	}
	b.writelnf("package %s\n", pkg.Name.Val)
}

func (b *builder) writeInit(init *ast.CodeBlock) {
	if init == nil {
		return
	}

	// remove opening and closing braces
	val := init.Val[1 : len(init.Val)-1]
	b.writelnf("%s\n", val)
}

func (b *builder) writef(f string, args ...interface{}) {
	if b.err == nil {
		_, b.err = fmt.Fprintf(b.w, f, args...)
	}
}

func (b *builder) writelnf(f string, args ...interface{}) {
	b.writef(f+"\n", args...)
}
