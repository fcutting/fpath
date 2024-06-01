package parser

import "github.com/shopspring/decimal"

const (
	NodeType_Undefined = iota
	NodeType_Block
	NodeType_Number
	NodeType_Equals
)

var NodeTypeString map[int]string = map[int]string{
	NodeType_Undefined: "Undefined",
	NodeType_Block:     "Block",
	NodeType_Number:    "Number",
	NodeType_Equals:    "Equals",
}

// Node is the most atomic piece of fpath syntax, describing both expressions
// and operations.
type Node interface {
	Type() int
}

func (BlockNode) Type() int  { return NodeType_Block }
func (NumberNode) Type() int { return NodeType_Number }
func (EqualsNode) Type() int { return NodeType_Equals }

// Expression nodes are evaluable in isolation of other nodes and don't depend
// on external data.
type Expression interface {
	Node
	expression()
}

func (BlockNode) expression()  {}
func (NumberNode) expression() {}

// Operation nodes require an additional input to evaluate to a value.
type Operation interface {
	Node
	operation()
}

func (EqualsNode) operation() {}

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

// EqualsNode represents an operation that compares the current value with an
// expression and updates the current value with the result.
type EqualsNode struct {
	Expresion Expression
}
