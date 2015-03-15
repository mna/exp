package ast

import (
	"strings"
	"testing"
	"unicode/utf8"
)

var charClasses = []string{
	"[]",
	"[]i",
	"[^]",
	"[^]i",
	"[a]",
	"[ab]i",
	"[^abc]i",
	`[\a]`,
	`[\b\nt]`,
	`[\b\nt\pL]`,
	`[\p{Greek}\tz\\\pD]`,
}

var expChars = []string{
	"",
	"",
	"",
	"",
	"a",
	"ab",
	"abc",
	"\a",
	"\b\nt",
	"\b\nt",
	"\tz\\",
}

var expUnicodeClasses = [][]string{
	9:  {"L"},
	10: {"Greek", "D"},
}

func TestCharClassParse(t *testing.T) {
	for i, c := range charClasses {
		m := NewCharClassMatcher(Pos{}, c)

		ic := strings.HasSuffix(c, "i")
		if m.IgnoreCase != ic {
			t.Errorf("%d: want ignore case: %t, got %t", i, ic, m.IgnoreCase)
		}
		iv := c[1] == '^'
		if m.Inverted != iv {
			t.Errorf("%d: want inverted: %t, got %t", i, iv, m.Inverted)
		}

		if n := utf8.RuneCountInString(expChars[i]); len(m.Chars) != n {
			t.Errorf("%d: want %d chars, got %d", i, n, len(m.Chars))
		} else if string(m.Chars) != expChars[i] {
			t.Errorf("%d: want %q, got %q", i, expChars[i], string(m.Chars))
		}

		if n := len(expUnicodeClasses[i]); len(m.UnicodeClasses) != n {
			t.Errorf("%d: want %d Unicode classes, got %d", i, n, len(m.UnicodeClasses))
		} else {
			want := strings.Join(expUnicodeClasses[i], "\n")
			got := strings.Join(m.UnicodeClasses, "\n")
			if want != got {
				t.Errorf("%d: want %v, got %v", i, expUnicodeClasses[i], m.UnicodeClasses)
			}
		}
	}
}
