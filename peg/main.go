// Command pigeon generates a PEG parser from a PEG grammar.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/PuerkitoBio/exp/peg/ast"
	"github.com/PuerkitoBio/exp/peg/builder"
)

func main() {
	dbgFlag := flag.Bool("debug", false, "set debug mode")
	noBuildFlag := flag.Bool("x", false, "do not build, only parse")
	outputFlag := flag.String("o", "", "output file, defaults to stdout")
	curRecvrNmFlag := flag.String("current-receiver-name", "c", "receiver name for the `current` type's generated methods")
	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s [options] [FILE]\n", os.Args[0])
		os.Exit(1)
	}

	nm := "stdin"
	inf := os.Stdin
	if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer f.Close()
		inf = f
		nm = flag.Arg(0)
	}
	in := bufio.NewReader(inf)

	debug = *dbgFlag
	g, err := Parse(nm, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error: ", err)
		os.Exit(3)
	}

	if !*noBuildFlag {
		outw := os.Stdout
		if *outputFlag != "" {
			f, err := os.Create(*outputFlag)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(4)
			}
			defer f.Close()
			outw = f
		}

		curNmOpt := builder.CurrentReceiverName(*curRecvrNmFlag)
		if err := builder.BuildParser(outw, g.(*ast.Grammar), curNmOpt); err != nil {
			fmt.Fprintln(os.Stderr, "build error: ", err)
			os.Exit(5)
		}
	}
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
