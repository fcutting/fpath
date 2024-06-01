package tokreader

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(m *testing.M) {
	r := m.Run()
	snaps.Clean(m, snaps.CleanOpts{Sort: true})
	os.Exit(r)
}

func _tokensMatch(expected, actual Token) (err error) {
	if expected.Type != actual.Type {
		err = fmt.Errorf("Unexpected type\nExpected: %s\nActual: %s", TokenTypeString[expected.Type], TokenTypeString[actual.Type])
		return
	}

	if expected.Value != actual.Value {
		err = fmt.Errorf("Unexpected value\nExpected: %s\nActual: %s", expected.Value, actual.Value)
	}

	return nil
}

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
	tr := NewTokenReader(input)

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
	tr := NewTokenReader(input)

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
	tr := NewTokenReader(input)

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
	tr := NewTokenReader(input)

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
		expectedTokens []Token
	}{
		"Whitespace": {
			input: "  123  ",
			expectedTokens: []Token{
				{Type: TokenType_Number, Value: "123"},
			},
		},
		"Number": {
			input: "123",
			expectedTokens: []Token{
				{Type: TokenType_Number, Value: "123"},
			},
		},
		"Label": {
			input: "fletcher",
			expectedTokens: []Token{
				{Type: TokenType_Label, Value: "fletcher"},
			},
		},
		"StringLiteral": {
			input: `"hello world"`,
			expectedTokens: []Token{
				{Type: TokenType_StringLiteral, Value: "hello world"},
			},
		},
		"Keyword Not": {
			input: "not",
			expectedTokens: []Token{
				{Type: TokenType_Not},
			},
		},
		"Keyword Equals": {
			input: "equals",
			expectedTokens: []Token{
				{Type: TokenType_Equals},
			},
		},
		"Keyword Contains": {
			input: "contains",
			expectedTokens: []Token{
				{Type: TokenType_Contains},
			},
		},
		"Keyword Greater": {
			input: "greater",
			expectedTokens: []Token{
				{Type: TokenType_Greater},
			},
		},
		"Keyword Lesser": {
			input: "lesser",
			expectedTokens: []Token{
				{Type: TokenType_Lesser},
			},
		},
		"OpenParan": {
			input: "(",
			expectedTokens: []Token{
				{Type: TokenType_OpenParan},
			},
		},
		"CloseParan": {
			input: ")",
			expectedTokens: []Token{
				{Type: TokenType_CloseParan},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tr := NewTokenReader(tc.input)

			for _, expected := range tc.expectedTokens {
				tok, err := tr.GetToken()

				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}

				if err := _tokensMatch(expected, tok); err != nil {
					t.Fatalf("Unexpected result: %s", err)
				}
			}
		})
	}
}

func Test_tokenReader_getToken_EOF(t *testing.T) {
	input := "  123  "
	expected := Token{
		Type:  TokenType_Number,
		Value: "123",
	}
	tr := NewTokenReader(input)

	tok, err := tr.GetToken()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := _tokensMatch(expected, tok); err != nil {
		t.Fatalf("Unexpected result: %s", err)
	}

	for i := 0; i < 100; i++ {
		if _, err := tr.GetToken(); err != io.EOF {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}

func Test_tokenReader_getToken_InvalidRune(t *testing.T) {
	input := "  123  `"
	expected := Token{
		Type:  TokenType_Number,
		Value: "123",
	}
	tr := NewTokenReader(input)

	tok, err := tr.GetToken()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := _tokensMatch(expected, tok); err != nil {
		t.Fatalf("Unexpected result: %s", err)
	}

	_, err = tr.GetToken()

	if err == nil {
		t.Fatalf("Error expected but not returned")
	}

	snaps.MatchSnapshot(t, err.Error())
}

func Test_tokenReader_peekToken(t *testing.T) {
	input := "123 equals"
	firstExpected := Token{
		Type:  TokenType_Number,
		Value: "123",
	}
	secondExpected := Token{
		Type: TokenType_Equals,
	}
	tr := NewTokenReader(input)
	shouldBreak := false

	for i := 0; i < 100; i++ {
		t.Run("First peek", func(t *testing.T) {
			tok, err := tr.PeekToken()

			if err != nil {
				shouldBreak = true
				t.Fatalf("Unexpected error: %s", err)
			}

			if err := _tokensMatch(firstExpected, tok); err != nil {
				shouldBreak = true
				t.Fatalf("Unexpected result: %s", err)
			}
		})

		if shouldBreak {
			break
		}
	}

	t.Run("First get", func(t *testing.T) {
		tok, err := tr.GetToken()

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if err := _tokensMatch(firstExpected, tok); err != nil {
			t.Fatalf("Unexpected result: %s", err)
		}
	})

	for i := 0; i < 100; i++ {
		t.Run("Second peek", func(t *testing.T) {
			tok, err := tr.PeekToken()

			if err != nil {
				shouldBreak = true
				t.Fatalf("Unexpected error: %s", err)
			}

			if err := _tokensMatch(secondExpected, tok); err != nil {
				shouldBreak = true
				t.Fatalf("Unexpected result: %s", err)
			}
		})

		if shouldBreak {
			break
		}
	}

	t.Run("Second get", func(t *testing.T) {
		tok, err := tr.GetToken()

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if err := _tokensMatch(secondExpected, tok); err != nil {
			t.Fatalf("Unexpected result: %s", err)
		}
	})
}

func Test_tokenReader_getTokenStringLiteral_UnexpectedEOF(t *testing.T) {
	input := `"hello `
	tr := NewTokenReader(input)

	_, err := tr.GetToken()

	if err == nil {
		t.Fatalf("Error expected but not returned")
	}

	snaps.MatchSnapshot(t, err.Error())
}
