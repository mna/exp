package main

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

func main() {
	set := make(map[string]bool)
	for k := range unicode.Categories {
		set[k] = true
	}
	for k := range unicode.Properties {
		set[k] = true
	}
	for k := range unicode.Scripts {
		set[k] = true
	}
	classes := make([]string, 0, len(set))
	for k := range set {
		classes = append(classes, k)
	}
	sort.Strings(classes)
	for i, s := range classes {
		fmt.Print("\t")
		if i > 0 {
			fmt.Print("/ ")
		}
		fmt.Println(strconv.Quote(s))
	}
}
