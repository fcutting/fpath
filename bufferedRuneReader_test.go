package fpath

import "testing"

func Test_bufferedRuneReader_get(t *testing.T) {
	input := "hello world"
	expectedRunes := []rune{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
	brr := newBufferedRuneReader(input)

	for i, expected := range expectedRunes {
		r := brr.get()
		if r != expected {
			t.Fatalf("Unexpected rune at position %d\nExpected: %q\n Actual: %q", i, expected, r)
		}
	}
}

func Test_bufferedRuneReader_peek(t *testing.T) {
	input := "hello world"
	expected := 'h'
	brr := newBufferedRuneReader(input)
	for i := 0; i < 10; i++ {
		r := brr.peek()
		if r != expected {
			t.Fatalf("Unexpected rune\nExpected: %q\nActual: %q", expected, r)
		}
	}
}
