package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/utils"
)

// HashInitFunction 实现 PHP hash_init()
type HashInitFunction struct{}

func NewHashInitFunction() data.FuncStmt {
	return &HashInitFunction{}
}

func (f *HashInitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	algo, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	h, err := newHashByAlgo(algo)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	hc := &HashContext{Hash: h, Algo: algo}
	rc := core.NewResourceClass("Hash Context", hc, allocHashContextID())
	return core.NewResourceValue(rc, ctx), nil
}

func (f *HashInitFunction) GetName() string            { return "hash_init" }
func (f *HashInitFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashInitFunction) GetIsStatic() bool          { return false }
func (f *HashInitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "algo", 0, nil, nil),
	}
}
func (f *HashInitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "algo", 0, nil),
	}
}
func (f *HashInitFunction) GetReturnType() data.Types { return data.NewBaseType("resource") }
