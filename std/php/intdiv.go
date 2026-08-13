package php

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IntdivFunction 实现 PHP intdiv() 的整数除法语义。
type IntdivFunction struct{}

func NewIntdivFunction() data.FuncStmt { return &IntdivFunction{} }

func (f *IntdivFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	numeratorValue, _ := ctx.GetIndexValue(0)
	denominatorValue, _ := ctx.GetIndexValue(1)

	numeratorInt, numeratorOK := numeratorValue.(data.AsInt)
	denominatorInt, denominatorOK := denominatorValue.(data.AsInt)
	if !numeratorOK || !denominatorOK {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("intdiv(): arguments must be integers"), "TypeError")
	}
	numerator, numeratorErr := numeratorInt.AsInt()
	denominator, denominatorErr := denominatorInt.AsInt()
	if numeratorErr != nil || denominatorErr != nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("intdiv(): arguments must be integers"), "TypeError")
	}
	if denominator == 0 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Division by zero"), "DivisionByZeroError")
	}

	return data.NewIntValue(numerator / denominator), nil
}

func (f *IntdivFunction) GetName() string { return "intdiv" }

func (f *IntdivFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "num1", 0, nil, data.Int{}),
		node.NewParameter(nil, "num2", 1, nil, data.Int{}),
	}
}

func (f *IntdivFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "num1", 0, data.Int{}),
		node.NewVariable(nil, "num2", 1, data.Int{}),
	}
}
