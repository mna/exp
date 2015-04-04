package main

import (
	"fmt"
	"sort"
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
	fmt.Println("package main")
	fmt.Println("\nvar unicodeClasses = map[string]bool{")
	for _, s := range classes {
		fmt.Printf("\t%q: true,\n", s)
	}
	fmt.Println("}")
}

// lenSorter was used to generate Unicode classes directly in the PEG
// grammar (where longer classes had to come first).
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
