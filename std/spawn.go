package std

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewSpawnFunction() data.FuncStmt {
	return &SpawnFunction{}
}

// SpawnFunction 封装 Go 的 go 关键字，异步执行闭包。
type SpawnFunction struct{}

func (f *SpawnFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cb, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("spawn 缺少闭包参数"))
	}
	if _, has := ctx.GetIndexValue(1); has {
		return nil, utils.NewThrow(errors.New("spawn 只接受一个闭包参数"))
	}

	var call func(data.Context) (data.GetValue, data.Control)
	var variables []data.Variable
	switch c := cb.(type) {
	case *data.BoundFuncValue:
		if !isClosure(c.Value) {
			return nil, utils.NewThrow(errors.New("spawn 只接受闭包"))
		}
		call = c.Call
		variables = c.Value.GetVariables()
	case *data.FuncValue:
		if !isClosure(c.Value) {
			return nil, utils.NewThrow(errors.New("spawn 只接受闭包"))
		}
		call = c.Call
		variables = c.Value.GetVariables()
	default:
		return nil, utils.NewThrow(errors.New("spawn 只接受闭包"))
	}

	vm := ctx.GetVM()
	// Allocate the closure frame before the native call frame can be released.
	callCtx := ctx.CreateContext(variables)
	go func() {
		_, acl := call(callCtx)
		if acl != nil {
			vm.ThrowControl(acl)
		}
	}()

	return nil, nil
}

func isClosure(v data.FuncStmt) bool {
	_, ok := v.(*node.LambdaExpression)
	return ok
}

func (f *SpawnFunction) GetName() string { return "spawn" }

var spawnFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "closure", 0, nil, nil),
}

func (f *SpawnFunction) GetParams() []data.GetValue {
	return spawnFunctionGetParams
}

var spawnFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "closure", 0, data.Mixed{}),
}

func (f *SpawnFunction) GetVariables() []data.Variable {
	return spawnFunctionGetVariables
}
