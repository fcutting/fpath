package fpath

import "testing"

func Test_bufferedRuneReader_Get(t *testing.T) {
	input := "hello world"
	expectedRunes := []rune{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
	brr := newBufferedRuneReader(input)

	for i, expected := range expectedRunes {
		r := brr.Get()
		if r != expected {
			t.Fatalf("Unexpected rune at position %d\nExpected: %q\n Actual: %q", i, expected, r)
		}
	}
}
