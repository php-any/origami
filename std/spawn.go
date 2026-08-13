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
	switch c := cb.(type) {
	case *data.BoundFuncValue:
		if !isClosure(c.Value) {
			return nil, utils.NewThrow(errors.New("spawn 只接受闭包"))
		}
		call = c.Call
	case *data.FuncValue:
		if !isClosure(c.Value) {
			return nil, utils.NewThrow(errors.New("spawn 只接受闭包"))
		}
		call = c.Call
	default:
		return nil, utils.NewThrow(errors.New("spawn 只接受闭包"))
	}

	vm := ctx.GetVM()
	go func() {
		callCtx := ctx.CreateContext(nil)
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

func (f *SpawnFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "closure", 0, nil, nil),
	}
}

func (f *SpawnFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "closure", 0, data.Mixed{}),
	}
}
