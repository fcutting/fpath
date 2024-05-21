package fpath

import (
	"io"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func Test_isLabelRune(t *testing.T) {
	testCases := map[string]struct {
		r        rune
		expected bool
	}{
		"1": {
			r:        '1',
			expected: true,
		},
		"f": {
			r:        'f',
			expected: true,
		},
		"_": {
			r:        '_',
			expected: true,
		},
		"$": {
			r:        '$',
			expected: false,
		},
		"-": {
			r:        '-',
			expected: false,
		},
		"Tab": {
			r:        '\t',
			expected: false,
		},
		"Space": {
			r:        ' ',
			expected: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := isLabelRune(tc.r)
			if result != tc.expected {
				t.Fatalf("Unexpected result\nExpected: %v\nActual: %v", tc.expected, result)
			}
		})
	}
}

func Test_tokenReader_getRune(t *testing.T) {
	input := "hello world"
	expectedRunes := []rune{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
	tr := newTokenReader(input)

	for i, expected := range expectedRunes {
		r, err := tr.getRune()

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if r != expected {
			t.Fatalf("Unexpected rune at position %d\nExpected: %q\n Actual: %q", i, expected, r)
		}
	}
}

func Test_tokenReader_getRune_EOF(t *testing.T) {
	input := "h"
	expected := 'h'
	tr := newTokenReader(input)

	r, err := tr.getRune()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if r != expected {
		t.Fatalf("Unexpected rune\nExpected: %q\nActual: %q", expected, r)
	}

	for i := 0; i < 100; i++ {
		if _, err = tr.getRune(); err != io.EOF {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}

func Test_tokenReader_peekRune(t *testing.T) {
	input := "hello world"
	expected := 'h'
	tr := newTokenReader(input)

	for i := 0; i < 10; i++ {
		r, err := tr.peekRune()

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if r != expected {
			t.Fatalf("Unexpected rune\nExpected: %q\nActual: %q", expected, r)
		}
	}
}

func Test_tokenReader_peekRune_EOF(t *testing.T) {
	input := "h"
	expected := 'h'
	tr := newTokenReader(input)

	r, err := tr.getRune()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if r != expected {
		t.Fatalf("Unexpected rune\nExpected: %q\nActual: %q", expected, r)
	}

	for i := 0; i < 100; i++ {
		if _, err = tr.peekRune(); err != io.EOF {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}

func Test_tokenReader_getToken(t *testing.T) {
	testCases := map[string]struct {
		input          string
		expectedTokens []token
	}{
		"Whitespace": {
			input: "  123  ",
			expectedTokens: []token{
				{typ: tokenType_Number, value: "123"},
			},
		},
		"Number": {
			input: "123",
			expectedTokens: []token{
				{typ: tokenType_Number, value: "123"},
			},
		},
		"Label": {
			input: "fletcher",
			expectedTokens: []token{
				{typ: tokenType_Label, value: "fletcher"},
			},
		},
		"StringLiteral": {
			input: `"hello world"`,
			expectedTokens: []token{
				{typ: tokenType_StringLiteral, value: "hello world"},
			},
		},
		"Keyword Not": {
			input: "not",
			expectedTokens: []token{
				{typ: tokenType_Not},
			},
		},
		"Keyword Equals": {
			input: "equals",
			expectedTokens: []token{
				{typ: tokenType_Equals},
			},
		},
		"Keyword Contains": {
			input: "contains",
			expectedTokens: []token{
				{typ: tokenType_Contains},
			},
		},
		"Keyword Greater": {
			input: "greater",
			expectedTokens: []token{
				{typ: tokenType_Greater},
			},
		},
		"Keyword Lesser": {
			input: "lesser",
			expectedTokens: []token{
				{typ: tokenType_Lesser},
			},
		},
		"OpenParan": {
			input: "(",
			expectedTokens: []token{
				{typ: tokenType_OpenParan},
			},
		},
		"CloseParan": {
			input: ")",
			expectedTokens: []token{
				{typ: tokenType_CloseParan},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tr := newTokenReader(tc.input)

			for _, expected := range tc.expectedTokens {
				tok, err := tr.getToken()

				if err != nil && err != io.EOF {
					t.Fatalf("Unexpected error: %s", err)
				}

				if tok.typ != expected.typ {
					t.Fatalf("Unexpected token type\nExpected: %d\nActual: %d", expected.typ, tok.typ)
				}

				if tok.value != expected.value {
					t.Fatalf("Unexpected token value\nExpected: %s\nActual: %s", expected.value, tok.value)
				}
			}
		})
	}
}

func Test_tokenReader_getToken_EOF(t *testing.T) {
	input := "  123  "
	expected := token{
		typ:   tokenType_Number,
		value: "123",
	}
	tr := newTokenReader(input)

	tok, err := tr.getToken()

	if err != nil && err != io.EOF {
		t.Fatalf("Unexpected error: %s", err)
	}

	if tok.typ != expected.typ {
		t.Fatalf("Unexpected token type\nExpected: %d\nActual: %d", expected.typ, tok.typ)
	}

	if tok.value != expected.value {
		t.Fatalf("Unexpected token value\nExpected: %s\nActual: %s", expected.value, tok.value)
	}

	for i := 0; i < 100; i++ {
		if _, err := tr.getToken(); err != io.EOF {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}

func Test_tokenReader_getToken_InvalidRune(t *testing.T) {
	input := "  123  `"
	expected := token{
		typ:   tokenType_Number,
		value: "123",
	}
	tr := newTokenReader(input)

	tok, err := tr.getToken()

	if err != nil && err != io.EOF {
		t.Fatalf("Unexpected error: %s", err)
	}

	if tok.typ != expected.typ {
		t.Fatalf("Unexpected token type\nExpected: %d\nActual: %d", expected.typ, tok.typ)
	}

	if tok.value != expected.value {
		t.Fatalf("Unexpected token value\nExpected: %s\nActual: %s", expected.value, tok.value)
	}

	_, err = tr.getToken()

	if err == nil {
		t.Fatalf("Error expected but not returned")
	}

	snaps.MatchSnapshot(t, err.Error())
}

func Test_tokenReader_getTokenStringLiteral_UnexpectedEOF(t *testing.T) {
	input := `"hello `
	tr := newTokenReader(input)

	_, err := tr.getToken()

	if err == nil {
		t.Fatalf("Error expected but not returned")
	}

	snaps.MatchSnapshot(t, err.Error())
}
