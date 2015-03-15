package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/exp/peg/bootstrap"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: pegscan FILE")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer f.Close()

	var s bootstrap.Scanner
	s.Init(os.Args[1], f, nil)
	for {
		tok, ok := s.Scan()
		fmt.Println(tok)
		if !ok {
			break
		}
	}
}
