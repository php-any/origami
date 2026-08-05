package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// HashUpdateFunction 实现 PHP hash_update()
type HashUpdateFunction struct{}

func NewHashUpdateFunction() data.FuncStmt {
	return &HashUpdateFunction{}
}

func (f *HashUpdateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hc, ctrl := getHashContext(ctx, 0)
	if ctrl != nil {
		return nil, ctrl
	}

	dataStr, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	_, err = hc.Hash.Write([]byte(dataStr))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *HashUpdateFunction) GetName() string            { return "hash_update" }
func (f *HashUpdateFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashUpdateFunction) GetIsStatic() bool          { return false }
func (f *HashUpdateFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "context", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, nil),
	}
}
func (f *HashUpdateFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "context", 0, nil),
		node.NewVariable(nil, "data", 1, nil),
	}
}
func (f *HashUpdateFunction) GetReturnType() data.Types { return data.NewBaseType("bool") }
