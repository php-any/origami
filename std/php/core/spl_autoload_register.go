package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// SplAutoloadRegisterFunction 实现 spl_autoload_register
// 当前实现为占位，记录回调的功能可后续扩展
type SplAutoloadRegisterFunction struct{}

func NewSplAutoloadRegisterFunction() data.FuncStmt { return &SplAutoloadRegisterFunction{} }

func (f *SplAutoloadRegisterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	switch f := a1.(type) {
	case *data.ArrayValue:
		valueList := f.ToValueList()
		class := valueList[0]
		methodName := valueList[1].AsString()
		var method data.Method
		var ok bool

		switch class := class.(type) {
		case *data.ThisValue:
			fn, acl := node.NewClassClosure(class.ClassValue, methodName)
			if acl != nil {
				return nil, acl
			}
			runtime.AddAutoLoad(fn)
		case *data.ClassValue:
			fn, acl := node.NewClassClosure(class, methodName)
			if acl != nil {
				return nil, acl
			}
			runtime.AddAutoLoad(fn)

		case *data.StringValue:
			stmt, acl := ctx.GetVM().GetOrLoadClass(class.AsString())
			if acl != nil {
				return nil, acl
			}

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
			runtime.AddAutoLoad(fn.(*data.FuncValue))
		default:
			_ = ok
			panic("spl_autoload_register 无法处理的类型，请添加支持")
		}
	case *data.FuncValue:
		fn := a1.(*data.FuncValue)

		runtime.AddAutoLoad(fn)
	}

	return data.NewBoolValue(true), nil
}

func (f *SplAutoloadRegisterFunction) GetName() string { return "spl_autoload_register" }

func (f *SplAutoloadRegisterFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}

func (f *SplAutoloadRegisterFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
	}
}
