package container

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// instantiateClass 在容器内完成类实例化与构造器调用；arguments 为已解析的依赖值。
func instantiateClass(stmt data.ClassStmt, arguments []data.GetValue, ctx data.Context) (data.GetValue, data.Control) {
	if node.IsAbstractClassStmt(stmt) {
		msg := fmt.Sprintf("Uncaught Error: Cannot instantiate abstract class %s", stmt.GetName())
		return nil, data.NewPHPUncaughtError(nil, msg)
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	cv, ok := object.(*data.ClassValue)
	if !ok {
		return object, nil
	}

	method := cv.Class.GetConstruct()
	if method == nil {
		return cv, nil
	}

	varies := method.GetVariables()
	params := method.GetParams()
	fnCtx := cv.CreateContext(varies)

	for index, param := range params {
		if index < len(arguments) && arguments[index] != nil {
			acl = setConstructorArg(fnCtx, cv, param, arguments[index], varies, index)
		} else {
			acl = applyConstructorDefault(fnCtx, cv, param)
		}
		if acl != nil {
			return nil, acl
		}
	}

	fnCtx.SetCallArgs(arguments)
	_, acl = method.Call(fnCtx)
	if acl != nil {
		return nil, acl
	}
	return cv, nil
}

func setConstructorArg(fnCtx, object data.Context, param data.GetValue, arg data.GetValue, varies []data.Variable, index int) data.Control {
	val, acl := valueFromArg(arg, fnCtx)
	if acl != nil {
		return acl
	}

	switch p := param.(type) {
	case *node.PromotedParameter:
		if index >= len(varies) {
			return data.NewErrorThrow(nil, fmt.Errorf("对象构造函数参数数量超出限制"))
		}
		if acl := fnCtx.SetVariableValue(varies[index], val); acl != nil {
			return acl
		}
		return p.SetValue(object, val)
	case *node.Parameter:
		return p.SetValue(fnCtx, val)
	case *node.CallerContextParameter:
		return nil
	default:
		if v, ok := param.(data.Variable); ok {
			return fnCtx.SetVariableValue(v, val)
		}
		return nil
	}
}

func applyConstructorDefault(fnCtx, object data.Context, param data.GetValue) data.Control {
	switch p := param.(type) {
	case *node.PromotedParameter:
		_, acl := p.GetValue(object)
		return acl
	case *node.CallerContextParameter:
		return nil
	default:
		if gp, ok := param.(data.GetValue); ok {
			_, acl := gp.GetValue(fnCtx)
			return acl
		}
		return nil
	}
}

func valueFromArg(arg data.GetValue, ctx data.Context) (data.Value, data.Control) {
	if v, ok := arg.(data.Value); ok {
		return v, nil
	}
	v, acl := arg.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if v == nil {
		return data.NewNullValue(), nil
	}
	val, ok := v.(data.Value)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("构造器实参类型错误"))
	}
	return val, nil
}
