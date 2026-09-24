package php

import (
	"crypto/rand"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RandomBytesFunction 实现 random_bytes(int $length): string
// 长度须走 AsInt：PHP 里 `$length - 6` 在 Origami 可能是 float，只认 *IntValue 会得到空串，
// Ramsey CombGenerator 就会在同一秒内生成相同 UUID（Telescope UNIQUE 失败）。
type RandomBytesFunction struct{}

func NewRandomBytesFunction() data.FuncStmt { return &RandomBytesFunction{} }

func (f *RandomBytesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	lengthValue, _ := ctx.GetIndexValue(0)
	length, ok := asRandomBytesLength(lengthValue)
	if !ok || length < 1 {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("random_bytes(): Argument #1 ($length) must be greater than 0"))
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	return data.NewStringValue(string(bytes)), nil
}

func asRandomBytesLength(v data.Value) (int, bool) {
	if v == nil {
		return 0, false
	}
	if ai, ok := v.(data.AsInt); ok {
		n, err := ai.AsInt()
		return n, err == nil
	}
	return 0, false
}

func (f *RandomBytesFunction) GetName() string { return "random_bytes" }
var randomBytesFunctionGetParams = []data.GetValue{node.NewParameter(nil, "length", 0, nil, nil)}

func (f *RandomBytesFunction) GetParams() []data.GetValue {
	return randomBytesFunctionGetParams
}
var randomBytesFunctionGetVariables = []data.Variable{node.NewVariable(nil, "length", 0, data.NewBaseType("int"))}

func (f *RandomBytesFunction) GetVariables() []data.Variable {
	return randomBytesFunctionGetVariables
}
