package fpath

import (
	"fmt"
	"io"
	"unicode"
)

const (
	TokenType_Undefined = iota
	TokenType_Number
)

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

		err = fmt.Errorf("Invalid rune %q", r)
		return
	}
}

// getTokenNumber returns the current number token in the input string.
// If the token reaches the end of the string, getTokenNumber also returns an
// io.EOF error.
func (tr *tokenReader) getTokenNumber() (tok token, err error) {
	tok.typ = TokenType_Number
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
