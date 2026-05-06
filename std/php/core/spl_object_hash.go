package core

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplObjectHashFunction 实现 spl_object_hash 函数
type SplObjectHashFunction struct{}

func NewSplObjectHashFunction() data.FuncStmt {
	return &SplObjectHashFunction{}
}

func (f *SplObjectHashFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	objVal, _ := ctx.GetIndexValue(0)
	if objVal == nil {
		return data.NewStringValue(""), nil
	}

	// Generate a hash based on the object's pointer address
	hash := fmt.Sprintf("%p", objVal)
	return data.NewStringValue(hash), nil
}

func (f *SplObjectHashFunction) GetName() string {
	return "spl_object_hash"
}

func (f *SplObjectHashFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}

func (f *SplObjectHashFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
	}
}
