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
	sort.Sort(lenSorter(classes))
	for i, s := range classes {
		fmt.Print("\t")
		if i > 0 {
			fmt.Print("/ ")
		}
		fmt.Println(strconv.Quote(s))
	}
}

type lenSorter []string

func (l lenSorter) Len() int      { return len(l) }
func (l lenSorter) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l lenSorter) Less(i, j int) bool {
	li, lj := len(l[i]), len(l[j])
	if lj < li {
		return true
	} else if li < lj {
		return false
	}
	return l[j] < l[i]
}
