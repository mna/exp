package main

import (
	"strings"
	"testing"
)

// With Unicode classes in the grammar:
// BenchmarkParseUnicodeClass          2000            548233 ns/op           96615 B/op        978 allocs/op
func BenchmarkParseUnicodeClass(b *testing.B) {
	const input = "a = [\\p{Latin}]"
	sr := strings.NewReader(input)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", sr); err != nil {
			b.Fatal(err)
		}
		if _, err := sr.Seek(0, 0); err != nil {
			b.Fatal(err)
		}
	}
}

// With keywords in the grammar:
// BenchmarkParseKeyword       5000            315189 ns/op           50175 B/op        530 allocs/op
func BenchmarkParseKeyword(b *testing.B) {
	const input = "a = uint32:'a'"
	sr := strings.NewReader(input)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", sr); err == nil {
			// error IS expected, fatal if none
			b.Fatal(err)
		}
		if _, err := sr.Seek(0, 0); err != nil {
			b.Fatal(err)
		}
	}
}
