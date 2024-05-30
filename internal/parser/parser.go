package parser

import (
	"fmt"

	"github.com/fcutting/fpath/internal/tokreader"
	"github.com/shopspring/decimal"
)

func NewParser(tr *tokreader.TokenReader) *Parser {
	return &Parser{
		tr: tr,
	}
}

type Parser struct {
	tr *tokreader.TokenReader
}

func parseNumber(token tokreader.Token) (number NumberNode, err error) {
	number.Value, err = decimal.NewFromString(token.Value)

	if err != nil {
		err = fmt.Errorf("failed to convert token value %q to number: %w", token.Value, err)
		return
	}

	return number, nil
}
