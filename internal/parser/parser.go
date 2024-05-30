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

func (p *Parser) ParseExpression() (expression Expression, err error) {
	token, err := p.tr.GetToken()

	if err != nil {
		err = fmt.Errorf("failed to get token: %w", err)
		return
	}

	switch token.Type {
	case tokreader.TokenType_Undefined:
		err = fmt.Errorf("encountered undefined token: %q", token.Value)
		return
	case tokreader.TokenType_Number:
		return parseNumber(token)
	default:
		err = fmt.Errorf("unknown token type: %v", token.Type)
		return
	}
}

func parseNumber(token tokreader.Token) (number NumberNode, err error) {
	number.Value, err = decimal.NewFromString(token.Value)

	if err != nil {
		err = fmt.Errorf("failed to convert token value %q to number: %w", token.Value, err)
		return
	}

	return number, nil
}
