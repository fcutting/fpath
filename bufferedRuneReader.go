package fpath

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	tokenType_Undefined = iota
	tokenType_Number
	tokenType_Label
	tokenType_StringLiteral
	tokenType_Not
	tokenType_Equals
	tokenType_Contains
	tokenType_Greater
	tokenType_Lesser
	tokenType_OpenParan
	tokenType_CloseParan
)

var UnexpectedEOF = errors.New("Unexpected EOF")

var keywords = map[string]int{
	"not":      tokenType_Not,
	"equals":   tokenType_Equals,
	"contains": tokenType_Contains,
	"greater":  tokenType_Greater,
	"lesser":   tokenType_Lesser,
}

func isLabelRune(r rune) bool {
	return unicode.IsNumber(r) || unicode.IsLetter(r) || r == '_'
}

type token struct {
	typ   int
	value string
}

// tokenReader adds the functionality to get and peek tokens from a
// string using a buffer.
type tokenReader struct {
	input []rune
	index int
}

// newTokenReader returns a new bufferedRuneReader configured to read from a
// []rune value of the input string.
func newTokenReader(input string) (tr *tokenReader) {
	return &tokenReader{
		input: []rune(input),
	}
}

// Get returns the rune at the current index of the input and increments the
// index.
func (tr *tokenReader) getRune() (r rune, err error) {
	if tr.index == len(tr.input) {
		return 0, io.EOF
	}

	r = tr.input[tr.index]
	tr.index++
	return r, nil
}

// Peek returns the rune at the current index of the input but doesn't
// increment the index.
func (tr *tokenReader) peekRune() (r rune, err error) {
	if tr.index == len(tr.input) {
		return 0, io.EOF
	}

	return tr.input[tr.index], nil
}

// getToken returns the next token in the input string.
// At the end of the input string, getToken returns an io.EOF error.
func (tr *tokenReader) getToken() (tok token, err error) {
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
			return token{
				typ: tokenType_OpenParan,
			}, nil
		case ')':
			tr.index++
			return token{
				typ: tokenType_CloseParan,
			}, nil
		default:
			err = fmt.Errorf("Invalid rune %q", r)
			return
		}
	}
}

// getTokenNumber returns the current number token in the input string.
// If the token reaches the end of the string, getTokenNumber also returns an
// io.EOF error.
func (tr *tokenReader) getTokenNumber() (tok token, err error) {
	tok.typ = tokenType_Number
	var r rune

	for {
		r, err = tr.peekRune()

		if err != nil {
			return tok, err
		}

		if unicode.IsNumber(r) {
			tr.index++
			tok.value += string(r)
			continue
		}

		return tok, nil
	}
}

// getTokenLabel returns the current label token in the input string.
// If the token reaches the end of the string, getTokenLabel also returns an
// io.EOF error.
func (tr *tokenReader) getTokenLabel() (tok token, err error) {
	tok.typ = tokenType_Label
	var r rune

	for {
		r, err = tr.peekRune()

		if err != nil {
			break
		}

		if isLabelRune(r) {
			tr.index++
			tok.value += string(r)
			continue
		}

		break
	}

	if key, ok := keywords[strings.ToLower(tok.value)]; ok {
		return token{
			typ: key,
		}, nil
	}

	return tok, err
}

// getTokenStringLiteral returns the current string literal token in the input
// string.
// If the token reaches the end of the string, getTokenStringLiteral returns an
// UnexpectedEOF error.
func (tr *tokenReader) getTokenStringLiteral() (tok token, err error) {
	tok.typ = tokenType_StringLiteral
	var r rune

	for {
		r, err = tr.getRune()

		if err != nil {
			err = UnexpectedEOF
			return
		}

		if r == '"' {
			break
		}

		tok.value += string(r)
	}

	return tok, nil
}
