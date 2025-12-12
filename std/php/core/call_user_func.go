package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// CallUserFuncFunction 实现 call_user_func
type CallUserFuncFunction struct{}

func NewCallUserFuncFunction() data.FuncStmt { return &CallUserFuncFunction{} }

func (f *CallUserFuncFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cb, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("call_user_func 缺少回调参数"))
	}

	// 收集实参（从 index 1 开始）
	var argValues []data.Value
	for i := 1; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok {
			break
		}
		argValues = append(argValues, v)
	}

	// 解析回调为 FuncValue
	fn, acl := f.resolveCallback(ctx, cb)
	if acl != nil {
		return nil, acl
	}
	// 创建调用上下文，传入实参
	callCtx := ctx.CreateContext(make([]data.Variable, len(argValues)))
	for i, v := range argValues {
		callCtx.GetIndexZVal(i).Value = v
	}
	return fn.Call(callCtx)
}

func (f *CallUserFuncFunction) resolveCallback(ctx data.Context, cb data.GetValue) (*data.FuncValue, data.Control) {
	switch c := cb.(type) {
	case *data.FuncValue:
		return c, nil
	case *data.ArrayValue:
		if len(c.Value) < 2 {
			return nil, utils.NewThrow(errors.New("call_user_func 回调数组长度不足"))
		}
		className := c.Value[0].AsString()
		methodName := c.Value[1].AsString()

		stmt, acl := ctx.GetVM().GetOrLoadClass(className)
		if acl != nil {
			return nil, acl
		}
		var method data.Method
		var ok bool
		method, ok = stmt.GetMethod(methodName)
		if !ok {
			if sm, ok2 := stmt.(data.GetStaticMethod); ok2 {
				method, ok = sm.GetStaticMethod(methodName)
			}
		}
		if !ok {
			return nil, utils.NewThrow(errors.New("call_user_func 未找到方法: " + className + "::" + methodName))
		}
		fn, acl := node.NewStaticMethodFuncValue(stmt, method).GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if fv, ok := fn.(*data.FuncValue); ok {
			return fv, nil
		}
		return nil, utils.NewThrow(errors.New("call_user_func 回调不是函数值"))
	default:
		return nil, utils.NewThrow(errors.New("call_user_func 回调不可调用"))
	}
}

func (f *CallUserFuncFunction) GetName() string { return "call_user_func" }

func (f *CallUserFuncFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (f *CallUserFuncFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
		node.NewVariable(nil, "args", 1, data.Mixed{}),
	}
}
