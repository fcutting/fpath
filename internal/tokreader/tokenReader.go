package tokreader

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

// NewTokenReader returns a new TokenReader configured to read from a slice
// []rune value of the input string.
func NewTokenReader(input string) *TokenReader {
	return &TokenReader{
		input: []rune(input),
	}
}

// tokenReader adds the functionality to get and peek tokens from a
// string using a buffer.
type TokenReader struct {
	input []rune
	index int
	buf   *Token
}

// getRune returns the rune at the current index of the input and increments the
// index.
// If the index is larger than the length of the input, getRune returns an
// io.EOF error.
func (tr *TokenReader) getRune() (r rune, err error) {
	if tr.index == len(tr.input) {
		return 0, io.EOF
	}

	r = tr.input[tr.index]
	tr.index++
	return r, nil
}

// peekRune returns the rune at the current index of the input but doesn't
// increment the index.
// If the index is larger than the length of the input, peekRune returns an
// io.EOF error.
func (tr *TokenReader) peekRune() (r rune, err error) {
	if tr.index == len(tr.input) {
		return 0, io.EOF
	}

	return tr.input[tr.index], nil
}

// getToken returns the next token in the input string.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (tr *TokenReader) GetToken() (tok Token, err error) {
	if tr.buf != nil {
		tok = *tr.buf
		tr.buf = nil
		return tok, nil
	}

	var r rune

	for {
		r, err = tr.peekRune()

		if err != nil {
			return tok, err
		}

		if unicode.IsSpace(r) {
			tr.index++
			continue
		}

		if unicode.IsNumber(r) {
			return tr.getTokenNumber()
		}

		if isLabelRune(r) {
			return tr.getTokenLabel()
		}

		switch r {
		case '"':
			tr.index++
			return tr.getTokenStringLiteral()
		case '(':
			tr.index++
			return Token{
				Type: TokenType_OpenParan,
			}, nil
		case ')':
			tr.index++
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
func (tr *TokenReader) PeekToken() (tok Token, err error) {
	if tr.buf != nil {
		tok = *tr.buf
		return tok, nil
	}

	tok, err = tr.GetToken()
	tr.buf = &tok
	return tok, err
}

// getTokenNumber returns the current number token in the input string.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (tr *TokenReader) getTokenNumber() (tok Token, err error) {
	tok.Type = TokenType_Number
	var r rune

	for {
		r, err = tr.peekRune()

		if err == io.EOF {
			return tok, nil
		}

		if err != nil {
			return tok, err
		}

		if unicode.IsNumber(r) {
			tr.index++
			tok.Value += string(r)
			continue
		}

		return tok, nil
	}
}

// getTokenLabel returns the current label token in the input string.
// If there are no more tokens to process in the string, getToken returns an
// io.EOF error.
func (tr *TokenReader) getTokenLabel() (tok Token, err error) {
	tok.Type = TokenType_Label
	var r rune

	for {
		r, err = tr.peekRune()

		if err != nil {
			break
		}

		if isLabelRune(r) {
			tr.index++
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
func (tr *TokenReader) getTokenStringLiteral() (tok Token, err error) {
	tok.Type = TokenType_StringLiteral
	var r rune

	for {
		r, err = tr.getRune()

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
