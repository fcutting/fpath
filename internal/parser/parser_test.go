package parser

import (
	"os"
	"testing"

	"github.com/fcutting/fpath/internal/lexer"
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
		"Equals": {
			input: "2 equals 4",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			lexer := lexer.NewLexer(tc.input)
			parser := NewParser(lexer)
			block, err := parser.ParseBlock()

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, block.String())
		})
	}
}

func Test_Parse_ParseOperation(t *testing.T) {
	testCases := map[string]struct {
		input string
	}{
		"Equals": {
			input: "equals 2",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			lexer := lexer.NewLexer(tc.input)
			parser := NewParser(lexer)
			operation, err := parser.ParseOperation()

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, operation.String())
		})
	}
}

func Test_Parse_ParseEquals(t *testing.T) {
	input := "2"
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	equals, err := parser.ParseEquals()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	snaps.MatchSnapshot(t, equals.String())
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
			lexer := lexer.NewLexer(tc.input)
			parser := NewParser(lexer)
			expression, err := parser.ParseExpression()

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, expression.String())
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
			lexer := lexer.NewLexer(tc.input)
			parser := NewParser(lexer)
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
			numberNode, err := parseNumber(lexer.Token{Type: lexer.TokenType_Number, Value: tc.value})

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, numberNode.String())
		})
	}
}

func Test_parseNumber_Error(t *testing.T) {
	testCases := map[string]struct {
		typ   int
		value string
	}{
		"Bad float": {
			typ:   lexer.TokenType_Number,
			value: "123,456",
		},
		"Word": {
			typ:   lexer.TokenType_Number,
			value: "kachow",
		},
		"Wrong type": {
			typ: lexer.TokenType_OpenParan,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := parseNumber(lexer.Token{Type: tc.typ, Value: tc.value})

			if err == nil {
				t.Fatalf("Expected error but none returned")
			}

			snaps.MatchSnapshot(t, err.Error())
		})
	}
}
