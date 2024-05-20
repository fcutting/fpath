package fpath

import (
	"io"
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
