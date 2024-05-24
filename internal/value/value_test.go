package value

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

func Test_ValueHolder_PutValue(t *testing.T) {
	v := NewValueHolder()
	tokenType := tokreader.TokenType_Number
	tokenValue := "123"
	value, err := v.PutValue(tokenType, tokenValue)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("Value", func(t *testing.T) {
		snaps.MatchJSON(t, value)
	})

	t.Run("NumberValues", func(t *testing.T) {
		snaps.MatchJSON(t, v.numberValues)
	})
}

func Test_ValueHolder_PutNumberValue(t *testing.T) {
	v := NewValueHolder()
	token := "123"
	value, err := v.PutNumberValue(token)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("NumberValues", func(t *testing.T) {
		snaps.MatchJSON(t, v.numberValues)
	})

	t.Run("Value", func(t *testing.T) {
		snaps.MatchJSON(t, value)
	})
}

func Test_ValueHolder_GetNumberValue(t *testing.T) {
	v := NewValueHolder()
	token := "123"
	value, err := v.PutNumberValue(token)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	number, ok := v.GetNumberValue(value.Key)

	if !ok {
		t.Fatalf("Value not found")
	}

	snaps.MatchJSON(t, number)
}
