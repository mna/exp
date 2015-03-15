package parser

import (
	"io"
	"os"
)

func ParseFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(filename, f)
}

func Parse(filename string, r io.Reader) (interface{}, error) {
	return nil, nil
}

type parser struct {
	filename string
	data     []byte
	cur      int
}
