package value

import (
	"fmt"

	"github.com/fcutting/fpath/internal/tokreader"
	"github.com/shopspring/decimal"
)

const (
	ValueType_Undefined = iota
	ValueType_Number
)

type Value struct {
	Type int
	Key  int
}

func NewValueHolder() *ValueHolder {
	return &ValueHolder{
		numberValues: map[int]decimal.Decimal{},
	}
}

type ValueHolder struct {
	numberValues map[int]decimal.Decimal
	numberIndex  int
}

func (v *ValueHolder) PutValue(tokenType int, tokenValue string) (value Value, err error) {
	switch tokenType {
	case tokreader.TokenType_Undefined:
		err = fmt.Errorf("Undefined token type provided")
		return
	case tokreader.TokenType_Number:
		return v.PutNumberValue(tokenValue)
	default:
		err = fmt.Errorf("Unknown token type %q", tokenType)
		return
	}
}

func (v *ValueHolder) PutNumberValue(tokenValue string) (value Value, err error) {
	number, err := decimal.NewFromString(tokenValue)

	if err != nil {
		err = fmt.Errorf("failed to convert %q to decimal.Decimal: %w", tokenValue, err)
		return
	}

	key := v.numberIndex
	v.numberIndex++
	v.numberValues[key] = number

	return Value{
		Type: ValueType_Number,
		Key:  key,
	}, nil
}

func (v ValueHolder) GetNumberValue(key int) (number decimal.Decimal, ok bool) {
	number, ok = v.numberValues[key]
	return number, ok
}
