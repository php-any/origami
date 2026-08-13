package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CountCharsFunction 实现 PHP count_chars() 的字节统计语义。
type CountCharsFunction struct{}

func NewCountCharsFunction() data.FuncStmt { return &CountCharsFunction{} }

func (f *CountCharsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	modeValue, _ := ctx.GetIndexValue(1)

	input := ""
	if value != nil {
		input = value.AsString()
	}
	mode := 0
	if iv, ok := modeValue.(data.AsInt); ok {
		if parsed, err := iv.AsInt(); err == nil {
			mode = parsed
		}
	}

	var counts [256]int
	for _, b := range []byte(input) {
		counts[int(b)]++
	}

	switch mode {
	case 0, 1, 2:
		list := make([]*data.ZVal, 0, 256)
		for i, count := range counts {
			include := mode == 0 || (mode == 1 && count > 0) || (mode == 2 && count == 0)
			if include {
				list = append(list, data.NewNamedZVal(data.IntArrayKeyName(i), data.NewIntValue(count)))
			}
		}
		return &data.ArrayValue{List: list}, nil
	case 3, 4:
		output := make([]byte, 0, 256)
		for i, count := range counts {
			if (mode == 3 && count > 0) || (mode == 4 && count == 0) {
				output = append(output, byte(i))
			}
		}
		return data.NewStringValue(string(output)), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *CountCharsFunction) GetName() string { return "count_chars" }

func (f *CountCharsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
		node.NewParameter(nil, "mode", 1, data.NewIntValue(0), data.Int{}),
	}
}

func (f *CountCharsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.String{}),
		node.NewVariable(nil, "mode", 1, data.Int{}),
	}
}
