package fpath

import (
	"io"
	"testing"
)

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
				{typ: TokenType_Number, value: "123"},
			},
		},
		"Number": {
			input: "123",
			expectedTokens: []token{
				{typ: TokenType_Number, value: "123"},
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
		typ:   TokenType_Number,
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
