package lexer

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	TokenType_Undefined = iota
	TokenType_Number
	TokenType_Label
	TokenType_StringLiteral
	TokenType_Not
	TokenType_Equals
	TokenType_Contains
	TokenType_Greater
	TokenType_Lesser
	TokenType_OpenParan
	TokenType_CloseParan
)

var TokenTypeString map[int]string = map[int]string{
	TokenType_Undefined:     "Undefined",
	TokenType_Number:        "Number",
	TokenType_Label:         "Label",
	TokenType_StringLiteral: "StringLiteral",
	TokenType_Not:           "Not",
	TokenType_Equals:        "Equals",
	TokenType_Contains:      "Contains",
	TokenType_Greater:       "Greater",
	TokenType_Lesser:        "Lesser",
	TokenType_OpenParan:     "OpenParan",
	TokenType_CloseParan:    "CloseParan",
}

var UnexpectedEOF = errors.New("Unexpected EOF")

var keywords = map[string]int{
	"not":      TokenType_Not,
	"equals":   TokenType_Equals,
	"contains": TokenType_Contains,
	"greater":  TokenType_Greater,
	"lesser":   TokenType_Lesser,
}

// isLabelRune returns whether the provided rune is a valid label rune.
// Valid label runes are letters, numbers, and underscores.
func isLabelRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_'
}

type Token struct {
	Type  int
	Value string
}

// NewLexer returns a new Lexer configured to read from a slice
// []rune value of the input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
	}
}

// Lexer adds the functionality to get and peek tokens from a
// string using a buffer.
type Lexer struct {
	input []rune
	index int
	buf   *Token
}

// getRune returns the rune at the current index of the input and increments the
// index.
// If the index is larger than the length of the input, getRune returns an
// io.EOF error.
func (l *Lexer) getRune() (r rune, err error) {
	if l.index == len(l.input) {
		return 0, io.EOF
	}

	r = l.input[l.index]
	l.index++
	return r, nil
}

// peekRune returns the rune at the current index of the input but doesn't
// increment the index.
// If the index is larger than the length of the input, peekRune returns an
// io.EOF error.
func (l *Lexer) peekRune() (r rune, err error) {
	if l.index == len(l.input) {
		return 0, io.EOF
	}

	return l.input[l.index], nil
}

// getToken returns the next token in the input string.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (l *Lexer) GetToken() (tok Token, err error) {
	if l.buf != nil {
		tok = *l.buf
		l.buf = nil
		return tok, nil
	}

	var r rune

	for {
		r, err = l.peekRune()

		if err != nil {
			return tok, err
		}

		if unicode.IsSpace(r) {
			l.index++
			continue
		}

		if unicode.IsNumber(r) {
			return l.getTokenNumber()
		}

		if isLabelRune(r) {
			return l.getTokenLabel()
		}

		switch r {
		case '"':
			l.index++
			return l.getTokenStringLiteral()
		case '(':
			l.index++
			return Token{
				Type: TokenType_OpenParan,
			}, nil
		case ')':
			l.index++
			return Token{
				Type: TokenType_CloseParan,
			}, nil
		default:
			err = fmt.Errorf("Invalid rune %q", r)
			return
		}
	}
}

// peekToken returns the current token in the input string but does not
// increment the index.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (l *Lexer) PeekToken() (tok Token, err error) {
	if l.buf != nil {
		tok = *l.buf
		return tok, nil
	}

	tok, err = l.GetToken()
	l.buf = &tok
	return tok, err
}

// getTokenNumber returns the current number token in the input string.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (l *Lexer) getTokenNumber() (tok Token, err error) {
	tok.Type = TokenType_Number
	var r rune

	for {
		r, err = l.peekRune()

		if err == io.EOF {
			return tok, nil
		}

		if err != nil {
			return tok, err
		}

		if unicode.IsNumber(r) {
			l.index++
			tok.Value += string(r)
			continue
		}

		return tok, nil
	}
}

// getTokenLabel returns the current label token in the input string.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (l *Lexer) getTokenLabel() (tok Token, err error) {
	tok.Type = TokenType_Label
	var r rune

	for {
		r, err = l.peekRune()

		if err != nil {
			break
		}

		if isLabelRune(r) {
			l.index++
			tok.Value += string(r)
			continue
		}

		break
	}

	if key, ok := keywords[strings.ToLower(tok.Value)]; ok {
		return Token{
			Type: key,
		}, nil
	}

	if err == io.EOF {
		return tok, nil
	}

	return tok, err
}

// getTokenStringLiteral returns the current string literal token in the input
// string.
// If the token reaches the end of the string, getTokenStringLiteral returns an
// UnexpectedEOF error.
func (l *Lexer) getTokenStringLiteral() (tok Token, err error) {
	tok.Type = TokenType_StringLiteral
	var r rune

	for {
		r, err = l.getRune()

		if err == io.EOF {
			err = UnexpectedEOF
			return
		}

		if err != nil {
			err = fmt.Errorf("unexpected error: %s", err)
			return
		}

		if r == '"' {
			break
		}

		tok.Value += string(r)
	}

	return tok, nil
}
