package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/fcutting/fpath/internal/tokreader"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(m *testing.M) {
	r := m.Run()
	snaps.Clean(m, snaps.CleanOpts{Sort: true})
	os.Exit(r)
}

func Test_Parse_ParseBlock(t *testing.T) {
	testCases := map[string]struct {
		input string
	}{
		"Arithmetic": {
			input: "2 + 4",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tr := tokreader.NewTokenReader(tc.input)
			parser := NewParser(tr)
			block, err := parser.ParseBlock()

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, fmt.Sprintf("Expression: %s: %v", NodeTypeString[block.Expression.Type()], block.Expression))
		})
	}
}

func Test_Parser_ParseExpression(t *testing.T) {
	testCases := map[string]struct {
		input string
	}{
		"Integer": {
			input: "123",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tr := tokreader.NewTokenReader(tc.input)
			parser := NewParser(tr)
			expression, err := parser.ParseExpression()

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, fmt.Sprintf("%T: %v", expression, expression))
		})
	}
}

func Test_Parser_ParseExpression_Error(t *testing.T) {
	testCases := map[string]struct {
		input string
	}{
		"Unknown": {
			input: "(",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tr := tokreader.NewTokenReader(tc.input)
			parser := NewParser(tr)
			_, err := parser.ParseExpression()

			if err == nil {
				t.Fatalf("Expected error but none returned")
			}

			snaps.MatchSnapshot(t, err.Error())
		})
	}
}

func Test_parseNumber(t *testing.T) {
	testCases := map[string]struct {
		value string
	}{
		"Integer": {
			value: "123",
		},
		"Float": {
			value: "123.456",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			numberNode, err := parseNumber(tokreader.Token{Type: tokreader.TokenType_Number, Value: tc.value})

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, numberNode.Value.String())
		})
	}
}

func Test_parseNumber_Error(t *testing.T) {
	testCases := map[string]struct {
		typ   int
		value string
	}{
		"Bad float": {
			typ:   tokreader.TokenType_Number,
			value: "123,456",
		},
		"Word": {
			typ:   tokreader.TokenType_Number,
			value: "kachow",
		},
		"Wrong type": {
			typ: tokreader.TokenType_OpenParan,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := parseNumber(tokreader.Token{Type: tc.typ, Value: tc.value})

			if err == nil {
				t.Fatalf("Expected error but none returned")
			}

			snaps.MatchSnapshot(t, err.Error())
		})
	}
}
