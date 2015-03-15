package json

import "github.com/PuerkitoBio/exp/parser/lexer"

const (
	StartObject = iota
	EndObject
	StartArray
	EndArray
	ObjectKey
	ObjectValue
	String
	Number
	True
	False
	Null
)

func startValue(l *lexer.Lexer) lexer.StateFn {
	// can be an object, an array, true, false, null
	// a string or a number.
	r, ok := l.SkipWhile(lexer.IsGoWhitespace)
	if !ok {
		return nil
	}
	switch r {
	case '{':
		l.Push(startObject)
		return startObject
		/*
			case '[':
				l.Push(startArray)
				return startArray
		*/
	case 't':
		l.Emit(scanTrue(l))
		return startValue
	case 'f':
		l.Emit(scanFalse(l))
		return startValue
	case 'n':
		l.Emit(scanNull(l))
		return startValue
		/*
			case '"':
				l.Emit(scanString)
				return startValue
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
				l.Emit(scanNumber)
				return startValue
		*/
	default:
		l.Errorf("invalid value: %c", r)
		l.Emit(lexer.Invalid, string(r))
		return startValue
	}
}

func scanTrue(l *lexer.Lexer) (int, string) {
	l.StartToken()
	if !l.Expect("rue") {
		return lexer.Invalid, l.Lit()
	}
	return True, l.Lit()
}

func scanFalse(l *lexer.Lexer) (int, string) {
	l.StartToken()
	if !l.Expect("alse") {
		return lexer.Invalid, l.Lit()
	}
	return False, l.Lit()
}

func scanNull(l *lexer.Lexer) (int, string) {
	l.StartToken()
	if !l.Expect("ull") {
		return lexer.Invalid, l.Lit()
	}
	return Null, l.Lit()
}

func startObject(l *lexer.Lexer) lexer.StateFn {
	l.Emit(StartObject, "{")
	return keyOrEndObject
}

func keyOrEndObject(l *lexer.Lexer) lexer.StateFn {
	r, ok := l.SkipWhile(lexer.IsGoWhitespace)
	if !ok {
		return nil
	}
	switch r {
	case '}':
		l.Emit(EndObject, "}")
		return nil
	case '"':
		return scanString
	default:
		l.Errorf("invalid value: %c", r)
		l.Emit(lexer.Invalid, string(r))
		return keyOrEndObject
	}
}
