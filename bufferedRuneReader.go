package fpath

// bufferedRuneReader adds the functionality to get and peek runes from a
// string using a buffer.
type bufferedRuneReader struct {
	input  []rune
	buffer rune
	index  int
}
