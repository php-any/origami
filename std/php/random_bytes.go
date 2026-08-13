package php

import (
	"crypto/rand"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RandomBytesFunction 实现 random_bytes 函数
type RandomBytesFunction struct{}

func NewRandomBytesFunction() data.FuncStmt { return &RandomBytesFunction{} }

func (f *RandomBytesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	lengthValue, _ := ctx.GetIndexValue(0)
	if lengthValue == nil {
		return data.NewStringValue(""), nil
	}

	var length int
	if iv, ok := lengthValue.(*data.IntValue); ok {
		length = iv.Value
	} else {
		return data.NewStringValue(""), nil
	}

	if length <= 0 {
		return data.NewStringValue(""), nil
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return data.NewStringValue(""), nil
	}

	return data.NewStringValue(string(bytes)), nil
}

func (f *RandomBytesFunction) GetName() string { return "random_bytes" }
func (f *RandomBytesFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "length", 0, nil, nil)}
}
func (f *RandomBytesFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "length", 0, data.NewBaseType("int"))}
}
