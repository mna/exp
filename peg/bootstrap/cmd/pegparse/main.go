package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/exp/peg/bootstrap"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: pegparse FILE")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer f.Close()

	p := bootstrap.New()
	if _, err := p.Parse(os.Args[1], f); err != nil {
		log.Fatal(err)
	}
}
