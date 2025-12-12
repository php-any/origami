package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// SplAutoloadUnregisterFunction 实现 spl_autoload_unregister
type SplAutoloadUnregisterFunction struct{}

func NewSplAutoloadUnregisterFunction() data.FuncStmt { return &SplAutoloadUnregisterFunction{} }

func (f *SplAutoloadUnregisterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	switch cb := a1.(type) {
	case *data.ArrayValue:
		className := cb.Value[0].AsString()
		methodName := cb.Value[1].AsString()

		stmt, acl := ctx.GetVM().GetOrLoadClass(className)
		if acl != nil {
			return nil, acl
		}

		var method data.Method
		var ok bool

		method, ok = stmt.GetMethod(methodName)
		if !ok {
			var c data.GetStaticMethod
			if c, ok = stmt.(data.GetStaticMethod); ok {
				method, ok = c.GetStaticMethod(methodName)
			}
		}
		fn, acl := node.NewStaticMethodFuncValue(stmt, method).GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if fun, ok := fn.(*data.FuncValue); ok {
			runtime.RemoveAutoLoad(fun)
		}
	case *data.FuncValue:
		runtime.RemoveAutoLoad(cb)
	default:
		return nil, utils.NewThrow(errors.New("spl_autoload_unregister 需要传入可调用类型"))
	}

	return data.NewBoolValue(true), nil
}

func (f *SplAutoloadUnregisterFunction) GetName() string { return "spl_autoload_unregister" }

func (f *SplAutoloadUnregisterFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}

func (f *SplAutoloadUnregisterFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
	}
}
