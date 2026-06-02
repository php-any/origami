package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// InstantiateController 在路由注册阶段实例化控制器（仅调用一次）。
func InstantiateController(stmt data.ClassStmt, ctx data.Context) (data.GetValue, data.Control) {
	obj, acl := createInstanceFromClassStmt(nil, stmt, nil, ctx)
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("控制器实例化失败"))
	}
	return cv, nil
}

// CallHTTPControllerMethod 在已准备好的接收者上调用路由方法（分发阶段不再 new）。
func CallHTTPControllerMethod(receiver data.GetValue, method data.Method, args []data.Value) (data.GetValue, data.Control) {
	if method.GetIsStatic() {
		cv, ok := receiver.(*data.ClassValue)
		if !ok {
			return nil, data.NewErrorThrow(nil, errors.New("静态路由缺少 ClassValue"))
		}
		fnCtx := cv.CreateContext(method.GetVariables())
		for i, arg := range args {
			if i < len(method.GetVariables()) {
				fnCtx.SetVariableValue(method.GetVariables()[i], arg)
			}
		}
		return method.Call(fnCtx)
	}

	cv, ok := receiver.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("控制器实例类型错误"))
	}

	m, has := cv.GetMethod(method.GetName())
	if !has {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("控制器不存在方法 %s", method.GetName()))
	}

	fnCtx := cv.CreateContext(m.GetVariables())
	for i, arg := range args {
		if i < len(m.GetVariables()) {
			fnCtx.SetVariableValue(m.GetVariables()[i], arg)
		}
	}
	return m.Call(fnCtx)
}
