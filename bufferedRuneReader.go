package fpath

// bufferedRuneReader adds the functionality to get and peek runes from a
// string using a buffer.
type bufferedRuneReader struct {
	input  []rune
	buffer rune
	index  int
}

// newBufferedRuneReader returns a new bufferedRuneReader configured to read
// from a []rune value of the input string.
func newBufferedRuneReader(input string) (brr *bufferedRuneReader) {
	return &bufferedRuneReader{
		input: []rune(input),
	}
}
