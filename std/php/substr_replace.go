package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SubstrReplaceFunction 实现 PHP substr_replace()
type SubstrReplaceFunction struct{}

func NewSubstrReplaceFunction() data.FuncStmt {
	return &SubstrReplaceFunction{}
}

func (f *SubstrReplaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	str, _ := ctx.GetIndexValue(0)
	repl, _ := ctx.GetIndexValue(1)
	offsetVal, _ := ctx.GetIndexValue(2)
	lengthVal, _ := ctx.GetIndexValue(3)

	if str == nil || repl == nil || offsetVal == nil {
		return data.NewStringValue(""), nil
	}

	s := str.AsString()
	r := repl.AsString()

	asInt, ok := offsetVal.(data.AsInt)
	if !ok {
		return data.NewStringValue(s), nil
	}
	offset, err := asInt.AsInt()
	if err != nil {
		return data.NewStringValue(s), nil
	}

	hasLength := false
	length := 0
	if lengthVal != nil {
		if _, isNull := lengthVal.(*data.NullValue); !isNull {
			if li, ok := lengthVal.(data.AsInt); ok {
				if l, err := li.AsInt(); err == nil {
					length = l
					hasLength = true
				}
			}
		}
	}

	n := len(s)

	// Normalize offset: negative counts from end
	if offset < 0 {
		offset = n + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset > n {
		offset = n
	}

	var actualLen int
	if !hasLength {
		// Default: replace to end of string
		actualLen = n - offset
	} else if length < 0 {
		// Explicit negative length: leave that many chars from end
		stop := n + length
		if stop < offset {
			actualLen = 0
		} else {
			actualLen = stop - offset
		}
	} else {
		actualLen = length
	}
	if offset+actualLen > n {
		actualLen = n - offset
	}
	if actualLen < 0 {
		actualLen = 0
	}

	result := s[:offset] + r + s[offset+actualLen:]
	return data.NewStringValue(result), nil
}

func (f *SubstrReplaceFunction) GetName() string            { return "substr_replace" }
func (f *SubstrReplaceFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *SubstrReplaceFunction) GetIsStatic() bool          { return false }
func (f *SubstrReplaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
		node.NewParameter(nil, "replacement", 1, nil, data.String{}),
		node.NewParameter(nil, "offset", 2, nil, data.Int{}),
		node.NewParameter(nil, "length", 3, node.NewNullLiteral(nil), data.Int{}),
	}
}
func (f *SubstrReplaceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "replacement", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 3, data.NewBaseType("int")),
	}
}
func (f *SubstrReplaceFunction) GetReturnType() data.Types { return data.NewBaseType("string") }
