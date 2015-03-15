package lexer

import (
	"crypto/rand"
	"io"
	"strings"
	"testing"
)

func TestPosString(t *testing.T) {
	cases := []struct {
		in  Pos
		out string
	}{
		{Pos{}, "0:0 [0]"},
		{Pos{Name: "a"}, "a:0:0 [0]"},
		{Pos{Name: "a", Line: 1}, "a:1:0 [0]"},
		{Pos{Name: "a", Line: 1, Col: 2, Off: 3}, "a:1:2 [3]"},
		{Pos{Line: 1, Col: 2, Off: 3}, "1:2 [3]"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.out {
			t.Errorf("%q: got %q", c.out, got)
		}
	}
}

func TestLexerPassthrough(t *testing.T) {
	cases := []struct {
		in   string
		toks []Token
	}{
		{
			"",
			[]Token{{-1, "", Pos{}}},
		},
		{
			"a",
			[]Token{
				{int('a'), "a", Pos{Line: 1, Col: 1, Off: 0}},
				{-1, "", Pos{Line: 1, Col: 1, Off: 0}},
			},
		},
		{
			"abc",
			[]Token{
				{int('a'), "a", Pos{Line: 1, Col: 1, Off: 0}},
				{int('b'), "b", Pos{Line: 1, Col: 2, Off: 1}},
				{int('c'), "c", Pos{Line: 1, Col: 3, Off: 2}},
				{-1, "", Pos{Line: 1, Col: 3, Off: 2}},
			},
		},
		{
			"a\nb\r\nc\n\rd",
			[]Token{
				{int('a'), "a", Pos{Line: 1, Col: 1, Off: 0}},
				{int('\n'), "\n", Pos{Line: 2, Col: 0, Off: 1}},
				{int('b'), "b", Pos{Line: 2, Col: 1, Off: 2}},
				{int('\r'), "\r", Pos{Line: 2, Col: 1, Off: 3}},
				{int('\n'), "\n", Pos{Line: 3, Col: 0, Off: 4}},
				{int('c'), "c", Pos{Line: 3, Col: 1, Off: 5}},
				{int('\n'), "\n", Pos{Line: 4, Col: 0, Off: 6}},
				{int('\r'), "\r", Pos{Line: 4, Col: 0, Off: 7}},
				{int('d'), "d", Pos{Line: 4, Col: 1, Off: 8}},
				{-1, "", Pos{Line: 4, Col: 1, Off: 8}},
			},
		},
	}

	var l Lexer
	for _, c := range cases {
		ch := l.Init("", strings.NewReader(c.in), stateRune)
		i := 0
		for got := range ch {
			if i >= len(c.toks) {
				t.Errorf("want %d tokens, got #%d (%s)", len(c.toks), i+1, got)
				i++
				continue
			}

			want := c.toks[i]
			if want != got {
				t.Errorf("%q: token %d: want %s, got %s", c.in, i, want, got)
			}
			i++
		}
	}
}

func benchmarkLexerPassthrough(b *testing.B, n int64) {
	var l Lexer

	r := io.LimitReader(rand.Reader, n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := l.Init("", r, stateRune)
		for range ch {
		}
	}
}

func BenchmarkLexerPassthrough1KB(b *testing.B) {
	benchmarkLexerPassthrough(b, 1<<10)
}

func BenchmarkLexerPassthrough100KB(b *testing.B) {
	benchmarkLexerPassthrough(b, 100<<10)
}

func BenchmarkLexerPassthrough1MB(b *testing.B) {
	benchmarkLexerPassthrough(b, 1<<20)
}
