package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	var (
		outputFlag  = flag.String("output", "", "output file, defaults to stdout")
		packageFlag = flag.String("package", "main", "package holding the stringified file, defaults to main")
		varFlag     = flag.String("var", "", "name of the variable holding the stringified file, defaults to the name of the input file without extension")
	)

	flag.Usage = usage
	flag.Parse()

	in := os.Getenv("GOFILE")
	switch flag.NArg() {
	case 0:
		if in == "" {
			fmt.Fprintf(os.Stderr, "expected an input file")
			fmt.Println()
			flag.Usage()
			os.Exit(1)
		}
	case 1:
		in = flag.Arg(0)
	default:
		fmt.Fprintf(os.Stderr, "expected one input file")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	// open input
	f, err := os.Open(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer f.Close()

	varName := *varFlag
	if varName == "" {
		ext := filepath.Ext(f.Name())
		varName = filepath.Base(f.Name())
		varName = varName[:len(varName)-len(ext)]
	}

	pkg := *packageFlag
	if pkg == "" {
		pkg = "main"
	}

	// open output
	out := os.Stdout
	if *outputFlag != "" {
		outf, err := os.Create(*outputFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		out = outf
	}

	fmt.Fprintf(out, `package %s

var %s = `, pkg, varName)
	fmt.Fprint(out, "`")
	io.Copy(out, f)
	fmt.Fprint(out, "`\n")
}

func usage() {
	fmt.Println("usage: stringify -output=FILE -package=PKG -var=IDENTIFIER INPUTFILE")
}
