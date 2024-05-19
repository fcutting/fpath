package fpath

// bufferedRuneReader adds the functionality to get and peek runes from a
// string using a buffer.
type bufferedRuneReader struct {
	input []rune
	index int
}

// newBufferedRuneReader returns a new bufferedRuneReader configured to read
// from a []rune value of the input string.
func newBufferedRuneReader(input string) (brr *bufferedRuneReader) {
	return &bufferedRuneReader{
		input: []rune(input),
	}
}

// Get returns the rune at the current index of the input and increments the
// index.
func (brr *bufferedRuneReader) get() (r rune) {
	r = brr.input[brr.index]
	brr.index++
	return r
}

// Peek returns the rune at the current index of the input but doesn't
// increment the index.
func (brr *bufferedRuneReader) peek() (r rune) {
	return brr.input[brr.index]
}
