package parser

import "github.com/shopspring/decimal"

const (
	NodeType_Block = iota
	NodeType_Number
	NodeType_Add
)

// Node is the most atomic piece of fpath syntax, describing both expressions
// and operations.
type Node interface {
	Type() int
}

func (BlockNode) Type() int  { return NodeType_Block }
func (NumberNode) Type() int { return NodeType_Number }
func (AddNode) Type() int    { return NodeType_Add }

// Expression nodes are evaluable in isolation of other nodes and don't depend
// on external data.
type Expression interface {
	expression()
}

func (BlockNode) expression()  {}
func (NumberNode) expression() {}

// Operation nodes require an additional input to evaluate to a value.
type Operation interface {
	operation()
}

func (AddNode) operation() {}

// BlockNode represents an executable fpath block that contains a base
// expression and a collection of operations to perform on the expression.
type BlockNode struct {
	Expression Expression
	Operations []Operation
}

// NumberNode represents a number literal.
type NumberNode struct {
	Value decimal.Decimal
}

// AddNode represents an operation that adds the current value of the block
// with an expression.
type AddNode struct {
	Value Expression
}
