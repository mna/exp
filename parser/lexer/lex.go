package lexer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	// ErrInvalidBack is returned when a call to back is made and no
	// previous rune is available.
	ErrInvalidBack = errors.New("lexer: invalid call to back: no previous rune")
)

const (
	EOF     = -(iota + 1) // reserved token ID for EOF
	Invalid               // reserved token ID for Invalid
)

// Pos represents a position in an arbitrary source of text.
type Pos struct {
	// Name is to identify the source of text.
	Name string

	// Line number.
	Line int

	// Column number.
	Col int

	// Offset.
	Off int
}

// String returns a textual representation of the position.
func (p Pos) String() string {
	args := []interface{}{p.Name, p.Line, p.Col, p.Off}
	f := "%s:%d:%d [%d]"
	if p.Name == "" {
		f = f[3:]
		args = args[1:]
	}
	return fmt.Sprintf(f, args...)
}

// Token is a token recognized by the lexer.
type Token struct {
	ID  int
	Lit string
	Pos Pos
}

// String is the textual representation of a token.
func (t Token) String() string {
	return fmt.Sprintf("%s: (%d) %q", t.Pos, t.ID, t.Lit)
}

// Lexer holds the state to tokenize text input. Init must be called to
// start emitting tokens.
type Lexer struct {
	r io.RuneReader

	eof bool // has reached EOF

	ch    chan Token
	stack []StateFn
	lit   bytes.Buffer
	intok bool

	tpos Pos  // pos of start of token
	cpos Pos  // pos of cur
	ppos Pos  // pos of pv
	cur  rune // current rune
	pv   rune // previous rune
}

// Init initializes the Lexer to read from r. The name is informational, it is
// used to report errors. The returned channel must be drained otherwise it may
// leak a goroutine. Init should not be called on a running lexer.
func (l *Lexer) Init(name string, r io.Reader, startFn StateFn) <-chan Token {
	l.r = runeReader(r)
	l.eof = false
	l.cpos = Pos{Name: name, Line: 1, Off: -1}
	l.tpos, l.ppos = Pos{}, Pos{}
	l.cur = -1
	l.pv = -1
	l.lit.Reset()
	l.intok = false

	l.ch = make(chan Token)
	go l.run(startFn)
	return l.ch
}

func runeReader(r io.Reader) io.RuneReader {
	if rr, ok := r.(io.RuneReader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

func (l *Lexer) run(startFn StateFn) {
	l.Push(startFn)

	var fn StateFn
	for {
		if fn == nil {
			fn = l.pop()
		}
		if fn == nil {
			break
		}
		fn = fn(l)
	}
	// always emit the EOF
	l.Emit(EOF, "")
	close(l.ch)
}

// Emit emits the token with the provided id and literal value.
func (l *Lexer) Emit(id int, lit string) {
	l.ch <- Token{id, lit, l.tpos}
	l.lit.Reset()
	l.intok = false
}

// Push adds a new StateFn on the stack. To start running this function, it
// should also be returned from the calling StateFn.
func (l *Lexer) Push(fn StateFn) {
	l.stack = append(l.stack, fn)
}

func (l *Lexer) pop() StateFn {
	if n := len(l.stack); n > 0 {
		fn := l.stack[n-1]
		l.stack = l.stack[:n-1]
		return fn
	}
	return nil
}

// Next advances the lexer to the next rune.
func (l *Lexer) Next() (rune, bool) {
	if l.eof {
		return l.cur, false
	}

	rn, w, err := l.r.ReadRune()
	if err != nil {
		l.fatalError(err)
		return l.cur, false
	}

	l.pv = l.cur
	l.ppos = l.cpos
	l.cur = rn
	l.cpos.Off += w

	switch rn {
	case '\n':
		l.cpos.Line++
		l.cpos.Col = 0
	case '\r':
		// ignore carriage returns, only increment line on \n so \r\n does
		// the right thing
	default:
		l.cpos.Col++
	}

	if l.intok {
		l.lit.WriteRune(l.cur)
	}
	return l.cur, true
}

// Back unreads the current rune, making it again the next rune to read.
// It can only be called once, and only after a successful call to Next was
// made.
func (l *Lexer) Back() {
	if l.pv == -1 {
		l.fatalError(ErrInvalidBack)
	}
	l.cur = l.pv
	l.cpos = l.ppos
	l.pv = -1
	if l.intok {
		l.lit.Truncate(l.lit.Len() - 1)
	}
}

func (l *Lexer) Lit() string {
	return l.lit.String()
}

// SkipWhile advances the lexer to the first rune that doesn't satisfy the
// provided predicate. It returns that rune and true if there was no read
// error, like a call to Next.
func (l *Lexer) SkipWhile(p RunePredicate) (rune, bool) {
	for p(l.cur) {
		if _, ok := l.Next(); !ok {
			return l.cur, false
		}
	}
	return l.cur, true
}

// SkipUntil advances the lexer to the first rune that satisfies the
// provided predicate. It returns that rune and true if there was no read
// error, like a call to Next.
func (l *Lexer) SkipUntil(p RunePredicate) (rune, bool) {
	for !p(l.cur) {
		if _, ok := l.Next(); !ok {
			return l.cur, false
		}
	}
	return l.cur, true
}

// StartToken saves the token position at the current position.
func (l *Lexer) StartToken() {
	l.tpos = l.cpos
	l.lit.WriteRune(l.cur)
	l.intok = true
}

func (l *Lexer) Expect(s string) bool {
	for _, want := range s {
		r, ok := l.Next()
		if !ok || r != want {
			l.Errorf("expected %c", want)
			return false
		}
	}
	return true
}

// acceptOneIn consumes the next rune if it is in the valid set.
func (l *Lexer) acceptOneIn(valid string) bool {
	l.Next()
	if strings.IndexRune(valid, l.cur) >= 0 {
		return true
	}
	l.Back()
	return false
}

// acceptManyIn consumes runes as long as it is in the valid set.
func (l *Lexer) acceptManyIn(valid string) int {
	n := 0
	for {
		l.Next()
		if strings.IndexRune(valid, l.cur) < 0 {
			l.Back()
			return n
		}
		n++
	}
}

func (l *Lexer) error(err error) {
	// TODO : report errors, default to stderr
	fmt.Fprintf(os.Stderr, "%s: %v\n", l.cpos, err)
}

func (l *Lexer) fatalError(err error) {
	l.cur = -1
	l.eof = true
	if err != io.EOF {
		l.error(err)
	}
}

// Errorf emits an error with the provided format string and arguments to
// use as error message.
func (l *Lexer) Errorf(f string, args ...interface{}) {
	l.error(fmt.Errorf(f, args...))
}
