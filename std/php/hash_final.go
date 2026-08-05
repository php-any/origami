package php

import (
	"encoding/hex"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HashFinalFunction 实现 PHP hash_final()
type HashFinalFunction struct{}

func NewHashFinalFunction() data.FuncStmt {
	return &HashFinalFunction{}
}

func (f *HashFinalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hc, ctrl := getHashContext(ctx, 0)
	if ctrl != nil {
		return nil, ctrl
	}

	rawOutput := false
	if rawVal, _ := ctx.GetIndexValue(1); rawVal != nil {
		if _, isNull := rawVal.(*data.NullValue); !isNull {
			if asBool, ok := rawVal.(data.AsBool); ok {
				if b, err := asBool.AsBool(); err == nil {
					rawOutput = b
				}
			}
		}
	}

	sum := hc.Hash.Sum(nil)
	if rawOutput {
		return data.NewStringValue(string(sum)), nil
	}
	return data.NewStringValue(hex.EncodeToString(sum)), nil
}

func (f *HashFinalFunction) GetName() string            { return "hash_final" }
func (f *HashFinalFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashFinalFunction) GetIsStatic() bool          { return false }
func (f *HashFinalFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "context", 0, nil, nil),
		node.NewParameter(nil, "binary", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *HashFinalFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "context", 0, nil),
		node.NewVariable(nil, "binary", 1, nil),
	}
}
func (f *HashFinalFunction) GetReturnType() data.Types { return data.NewBaseType("string") }
