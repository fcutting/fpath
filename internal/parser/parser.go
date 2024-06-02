package parser

import (
	"fmt"
	"io"

	"github.com/fcutting/fpath/internal/lexer"
	"github.com/shopspring/decimal"
)

func NewParser(lexer *lexer.Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

// Parser parses a tokenized string into an executable AST.
type Parser struct {
	lexer *lexer.Lexer
}

// ParseBlock returns the next block in the query.
func (p *Parser) ParseBlock() (block BlockNode, err error) {
	block.Expression, err = p.ParseExpression()

	if err != nil {
		err = fmt.Errorf("failed to parse expression: %w", err)
	}

	for {
		var operation Operation
		operation, err = p.ParseOperation()

		if err == io.EOF {
			break
		}

		if err != nil {
			err = fmt.Errorf("failed to parser operation: %w", err)
			return
		}

		block.Operations = append(block.Operations, operation)
	}

	return block, nil
}

// ParseOperation returns the next operation in the query.
// If the next token is not an operation, this step will return an error.
func (p *Parser) ParseOperation() (operation Operation, err error) {
	token, err := p.lexer.GetToken()

	if err == io.EOF {
		return nil, err
	}

	if err != nil {
		err = fmt.Errorf("failed to get token: %w", err)
		return
	}

	switch token.Type {
	case lexer.TokenType_Undefined:
		err = fmt.Errorf("encountered undefined token: %q", token.Value)
		return
	case lexer.TokenType_Equals:
		return p.ParseEquals()
	default:
		err = fmt.Errorf("unsupported token type: %s", lexer.TokenTypeString[token.Type])
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
	token, err := p.lexer.GetToken()

	if err != nil {
		err = fmt.Errorf("failed to get token: %w", err)
		return
	}

	switch token.Type {
	case lexer.TokenType_Undefined:
		err = fmt.Errorf("encountered undefined token: %q", token.Value)
		return
	case lexer.TokenType_Number:
		return parseNumber(token)
	default:
		err = fmt.Errorf("unsupported token type: %s", lexer.TokenTypeString[token.Type])
		return
	}
}

// parseNumber accepts a number token and converts it to a NumberNode.
func parseNumber(token lexer.Token) (number NumberNode, err error) {
	if token.Type != lexer.TokenType_Number {
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
