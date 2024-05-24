package value

import (
	"fmt"

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

func (v *ValueHolder) PutNumberValue(token string) (value Value, err error) {
	number, err := decimal.NewFromString(token)

	if err != nil {
		err = fmt.Errorf("failed to convert %q to decimal.Decimal: %w", token, err)
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
