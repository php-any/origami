package php

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RandomIntFunction 实现 random_int 函数
type RandomIntFunction struct{}

func NewRandomIntFunction() data.FuncStmt { return &RandomIntFunction{} }

func (f *RandomIntFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	minValue, _ := ctx.GetIndexValue(0)
	maxValue, _ := ctx.GetIndexValue(1)

	min := 0
	max := 0
	if iv, ok := minValue.(*data.IntValue); ok {
		min = iv.Value
	}
	if iv, ok := maxValue.(*data.IntValue); ok {
		max = iv.Value
	}

	if min > max {
		return data.NewIntValue(0), nil
	}

	diff := max - min
	if diff < 0 {
		return data.NewIntValue(min), nil
	}

	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return data.NewIntValue(min), nil
	}

	r := int(binary.LittleEndian.Uint64(b[:]))
	if r < 0 {
		r = -r
	}

	result := min + (r % (diff + 1))
	return data.NewIntValue(result), nil
}

func (f *RandomIntFunction) GetName() string { return "random_int" }
func (f *RandomIntFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "min", 0, nil, nil),
		node.NewParameter(nil, "max", 1, nil, nil),
	}
}
func (f *RandomIntFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "min", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "max", 1, data.NewBaseType("int")),
	}
}
