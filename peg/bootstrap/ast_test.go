package bootstrap

import (
	"strings"
	"testing"
	"unicode"
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
	`[\p{Greek}\tz\\\pN]`,
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

var expUnicodeClasses = [][]*unicode.RangeTable{
	9:  {unicode.Categories["L"]},
	10: {unicode.Scripts["Greek"], unicode.Categories["N"]},
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
		} else if n > 0 {
			want := expUnicodeClasses[i]
			got := m.UnicodeClasses
			for j, wantClass := range want {
				if wantClass != got[j] {
					t.Errorf("%d: range table %d: want %v, got %v", i, j, wantClass, got[j])
				}
			}
		}
	}
}
