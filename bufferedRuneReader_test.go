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
