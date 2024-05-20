package fpath

import "testing"

func Test_tokenReader_get(t *testing.T) {
	input := "hello world"
	expectedRunes := []rune{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
	tr := newTokenReader(input)

	for i, expected := range expectedRunes {
		r := tr.getRune()
		if r != expected {
			t.Fatalf("Unexpected rune at position %d\nExpected: %q\n Actual: %q", i, expected, r)
		}
	}
}

func Test_tokenReader_peek(t *testing.T) {
	input := "hello world"
	expected := 'h'
	tr := newTokenReader(input)
	for i := 0; i < 10; i++ {
		r := tr.peekRune()
		if r != expected {
			t.Fatalf("Unexpected rune\nExpected: %q\nActual: %q", expected, r)
		}
	}
}
