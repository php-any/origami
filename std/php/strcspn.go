package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrcspnFunction 实现 strcspn: 返回字符串中不包含指定字符的初始段长度
func NewStrcspnFunction() data.FuncStmt {
	return &StrcspnFunction{}
}

type StrcspnFunction struct{}

func (f *StrcspnFunction) GetName() string {
	return "strcspn"
}

func (f *StrcspnFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "characters", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, data.NewIntValue(0), nil),
		node.NewParameter(nil, "length", 3, nil, nil),
	}
}

func (f *StrcspnFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "characters", 1, nil),
		node.NewVariable(nil, "offset", 2, nil),
		node.NewVariable(nil, "length", 3, nil),
	}
}

func (f *StrcspnFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strValue, _ := ctx.GetIndexValue(0)
	charsValue, _ := ctx.GetIndexValue(1)

	if strValue == nil || charsValue == nil {
		return data.NewIntValue(0), nil
	}

	str := strValue.AsString()
	chars := charsValue.AsString()

	// Handle offset
	offset := 0
	offsetValue, exists := ctx.GetIndexValue(2)
	if exists && offsetValue != nil {
		if _, isNull := offsetValue.(*data.NullValue); !isNull {
			if asInt, ok := offsetValue.(data.AsInt); ok {
				offset, _ = asInt.AsInt()
			}
		}
	}

	if offset < 0 {
		offset = len(str) + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset >= len(str) {
		return data.NewIntValue(0), nil
	}

	str = str[offset:]

	// Handle length - NullValue means "no length specified"
	lengthValue, exists := ctx.GetIndexValue(3)
	if exists && lengthValue != nil {
		if _, isNull := lengthValue.(*data.NullValue); !isNull {
			if asInt, ok := lengthValue.(data.AsInt); ok {
				l, _ := asInt.AsInt()
				if l < 0 {
					l = len(str) + l
				}
				if l <= 0 {
					return data.NewIntValue(0), nil
				}
				if l < len(str) {
					str = str[:l]
				}
			}
		}
	}

	// Find length of initial segment NOT containing any char from chars
	for i, c := range str {
		if strings.ContainsRune(chars, c) {
			return data.NewIntValue(i), nil
		}
	}
	return data.NewIntValue(len(str)), nil
}

// StrspnFunction 实现 strspn: 返回字符串中只包含指定字符的初始段长度
func NewStrspnFunction() data.FuncStmt {
	return &StrspnFunction{}
}

type StrspnFunction struct{}

func (f *StrspnFunction) GetName() string {
	return "strspn"
}

func (f *StrspnFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "characters", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, data.NewIntValue(0), nil),
		node.NewParameter(nil, "length", 3, nil, nil),
	}
}

func (f *StrspnFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "characters", 1, nil),
		node.NewVariable(nil, "offset", 2, nil),
		node.NewVariable(nil, "length", 3, nil),
	}
}

func (f *StrspnFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strValue, _ := ctx.GetIndexValue(0)
	charsValue, _ := ctx.GetIndexValue(1)

	if strValue == nil || charsValue == nil {
		return data.NewIntValue(0), nil
	}

	str := strValue.AsString()
	chars := charsValue.AsString()

	// Handle offset
	offset := 0
	offsetValue, exists := ctx.GetIndexValue(2)
	if exists && offsetValue != nil {
		if _, isNull := offsetValue.(*data.NullValue); !isNull {
			if asInt, ok := offsetValue.(data.AsInt); ok {
				offset, _ = asInt.AsInt()
			}
		}
	}

	if offset < 0 {
		offset = len(str) + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset >= len(str) {
		return data.NewIntValue(0), nil
	}

	str = str[offset:]

	// Handle length - NullValue means "no length specified"
	lengthValue, exists := ctx.GetIndexValue(3)
	if exists && lengthValue != nil {
		if _, isNull := lengthValue.(*data.NullValue); !isNull {
			if asInt, ok := lengthValue.(data.AsInt); ok {
				l, _ := asInt.AsInt()
				if l < 0 {
					l = len(str) + l
				}
				if l <= 0 {
					return data.NewIntValue(0), nil
				}
				if l < len(str) {
					str = str[:l]
				}
			}
		}
	}

	// Find length of initial segment containing ONLY chars from chars
	for i, c := range str {
		if !strings.ContainsRune(chars, c) {
			return data.NewIntValue(i), nil
		}
	}
	return data.NewIntValue(len(str)), nil
}
