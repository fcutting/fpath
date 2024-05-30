package parser

import (
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
			numberNode, err := parseNumber(tokreader.Token{Value: tc.value})

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			snaps.MatchSnapshot(t, numberNode.Value.String())
		})
	}
}

func Test_parseNumber_Error(t *testing.T) {
	testCases := map[string]struct {
		value string
	}{
		"Bad float": {
			value: "123,456",
		},
		"Word": {
			value: "kachow",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := parseNumber(tokreader.Token{Value: tc.value})

			if err == nil {
				t.Fatalf("Expected error but none returned")
			}

			snaps.MatchSnapshot(t, err.Error())
		})
	}
}
