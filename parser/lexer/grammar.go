package lexer

import (
	"strings"
	"unicode"
)

type RunePredicate func(rune) bool

type StateFn func(*Lexer) StateFn

func stateRune(l *Lexer) StateFn {
	r, ok := l.Next()
	if !ok {
		return nil
	}
	l.StartToken()
	l.Emit(int(r), string(r))
	return stateRune
}

func EqPredicate(want rune) RunePredicate {
	return func(r rune) bool {
		return r == want
	}
}

func IsGoWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func IsAsciiWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n' || r == '\v' || r == '\f'
}

func AndPredicate(ps ...RunePredicate) RunePredicate {
	return func(r rune) bool {
		for _, p := range ps {
			if !p(r) {
				return false
			}
		}
		return true
	}
}

func OrPredicate(ps ...RunePredicate) RunePredicate {
	return func(r rune) bool {
		for _, p := range ps {
			if p(r) {
				return true
			}
		}
		return false
	}
}

func RangeTablePredicate(t *unicode.RangeTable) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(t, r)
	}
}

func SetPredicate(set string) RunePredicate {
	return func(r rune) bool {
		return strings.IndexRune(set, r) >= 0
	}
}

func NotPredicate(p RunePredicate) RunePredicate {
	return func(r rune) bool {
		return !p(r)
	}
}
