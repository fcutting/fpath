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

// ParseBlock returns the next block in the query.
func (p *Parser) ParseBlock() (block BlockNode, err error) {
	block.Expression, err = p.ParseExpression()

	if err != nil {
		err = fmt.Errorf("failed to parse expression: %w", err)
	}

	return block, nil
}

// ParseOperation returns the next operation in the query.
// If the next token is not an operation, this step will return an error.
func (p *Parser) ParseOperation() (operation Operation, err error) {
	token, err := p.tr.GetToken()

	if err != nil {
		err = fmt.Errorf("failed to get token: %w", err)
		return
	}

	switch token.Type {
	case tokreader.TokenType_Undefined:
		err = fmt.Errorf("encountered undefined token: %q", token.Value)
		return
	case tokreader.TokenType_Equals:
		return p.ParseEquals()
	default:
		err = fmt.Errorf("unsupported token type: %s", tokreader.TokenTypeString[token.Type])
		return
	}
}

// ParseEquals returns a parsed EqualsNode assuming the current operation is an
// equals operation.
func (p *Parser) ParseEquals() (equals EqualsNode, err error) {
	equals.Expression, err = p.ParseExpression()

	if err != nil {
		err = fmt.Errorf("failed to parse expression: %w", err)
		return
	}

	return equals, nil
}

// ParseExpression returns the next expression in the query.
// If the next token is not an expression, this step will return an error.
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
		err = fmt.Errorf("unsupported token type: %s", tokreader.TokenTypeString[token.Type])
		return
	}
}

// parseNumber accepts a number token and converts it to a NumberNode.
func parseNumber(token tokreader.Token) (number NumberNode, err error) {
	if token.Type != tokreader.TokenType_Number {
		err = fmt.Errorf("Token type is not a number: %v", token.Type)
		return
	}

	number.Value, err = decimal.NewFromString(token.Value)

	if err != nil {
		err = fmt.Errorf("failed to convert token value %q to number: %w", token.Value, err)
		return
	}

	return number, nil
}
