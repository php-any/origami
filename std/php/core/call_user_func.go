package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
	if fn == nil {
		return data.NewBoolValue(false), nil
	}
	// 创建调用上下文，传入实参（index 0 为 callback）
	callCtx := ctx.CreateContext(make([]data.Variable, len(argValues)))
	for i := range argValues {
		callCtx.SetIndexZVal(i, ctx.GetIndexZVal(i+1))
	}
	// BoundFuncValue 需要保留以确保 BoundContext 被创建
	if bfv, ok := cb.(*data.BoundFuncValue); ok {
		return bfv.Call(callCtx)
	}
	return fn.Call(callCtx)
}

func (f *CallUserFuncFunction) resolveCallback(ctx data.Context, cb data.GetValue) (*data.FuncValue, data.Control) {
	switch c := cb.(type) {
	case *data.FuncValue:
		return c, nil
	case *data.BoundFuncValue:
		return &c.FuncValue, nil
	case *data.ArrayValue:
		valueList := c.ToValueList()
		if len(valueList) < 2 {
			return nil, utils.NewThrow(errors.New("call_user_func 回调数组长度不足"))
		}
		methodName := valueList[1].AsString()
		if cv, ok := valueList[0].(*data.ClassValue); ok {
			return f.resolveObjectCallback(ctx, cv, methodName)
		}
		if tv, ok := valueList[0].(*data.ThisValue); ok && tv.ClassValue != nil {
			return f.resolveObjectCallback(ctx, tv.ClassValue, methodName)
		}
		className := valueList[0].AsString()

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

func (f *CallUserFuncFunction) resolveObjectCallback(ctx data.Context, cv *data.ClassValue, methodName string) (*data.FuncValue, data.Control) {
	if strings.Contains(methodName, "::") {
		parts := strings.SplitN(methodName, "::", 2)
		if len(parts) == 2 {
			className, mName := parts[0], parts[1]
			file, line := callUserFuncCallSite(ctx)
			objName := ""
			if cv.Class != nil {
				objName = cv.Class.GetName()
			}
			prefix := ""
			if data.HasUserOutput() {
				prefix = "\n"
			}
			_, _ = fmt.Fprintf(os.Stderr,
				"%sDeprecated: Callables of the form [\"%s\", \"%s\"] are deprecated in %s on line %d\n",
				prefix, objName, methodName, file, line,
			)
			stmt, acl := ctx.GetVM().GetOrLoadClass(className)
			if acl != nil {
				return nil, acl
			}
			method, ok := stmt.GetMethod(mName)
			if !ok {
				if sm, ok2 := stmt.(interface {
					GetStaticMethod(string) (data.Method, bool)
				}); ok2 {
					method, ok = sm.GetStaticMethod(mName)
				}
			}
			if !ok {
				return nil, utils.NewThrow(errors.New("call_user_func 未找到方法: " + methodName))
			}
			if _, isAbstract := method.(*node.AbstractMethod); isAbstract {
				_, _ = fmt.Fprintf(os.Stderr,
					"call_user_func(): Argument #1 ($callback) must be a valid callback, cannot call abstract method %s::%s()\n",
					className, mName,
				)
				return nil, nil
			}
			return data.NewFuncValue(node.NewObjectMethodCallable(cv, mName)), nil
		}
	}
	return data.NewFuncValue(node.NewObjectMethodCallable(cv, methodName)), nil
}

func callUserFuncCallSite(ctx data.Context) (file string, line int) {
	for _, arg := range ctx.GetCallArgs() {
		if g, ok := arg.(node.GetFrom); ok && g.GetFrom() != nil {
			from := g.GetFrom()
			file = from.GetSource()
			line, _, _, _ = from.GetRange()
			line++
			return file, line
		}
	}
	return "", 0
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
