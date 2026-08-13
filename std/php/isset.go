package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIssetFunction() data.FuncStmt {
	return &IssetFunction{}
}

type IssetFunction struct{}

func (f *IssetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	params := f.GetParams()
	if len(params) == 0 {
		return data.NewBoolValue(false), nil
	}

	// 获取第一个参数
	varValue, _ := ctx.GetIndexValue(0)

	// 如果参数是 ASTValue，手动计算并处理错误
	if astValue, ok := varValue.(*data.ASTValue); ok {
		// 计算 AST 值
		val, acl := astValue.Node.GetValue(astValue.Ctx)
		if acl != nil {
			// 如果计算过程中报错（如未定义变量），isset 应该返回 false 而不是抛出异常
			return data.NewBoolValue(false), nil
		}

		if val == nil {
			return data.NewBoolValue(false), nil
		}

		// 检查是否为 NullValue
		if _, ok := val.(*data.NullValue); ok {
			return data.NewBoolValue(false), nil
		}

		return data.NewBoolValue(true), nil
	}

	// 兼容逻辑
	if varValue == nil {
		// 尝试从 ctx 获取变量（如果通过 parameter 传递并在 ctx 中设置了）
		// 但由于 parameter 是 RawAST，如果不走上面的分支，说明有问题
		return data.NewBoolValue(false), nil
	}

	switch varValue.(type) {
	case *data.NullValue:
		return data.NewBoolValue(false), nil
	default:
		// 如果值存在且不是 null
		return data.NewBoolValue(true), nil
	}
}

func (f *IssetFunction) GetName() string {
	return "isset"
}

func (f *IssetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterRawAST(nil, "var", 0, data.Mixed{}),
	}
}

func (f *IssetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "var", 0, data.Mixed{}),
	}
}
