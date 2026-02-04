package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewEmptyFunction() data.FuncStmt {
	return &EmptyFunction{}
}

type EmptyFunction struct{}

func (f *EmptyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数值（ASTWrapper）
	varValue, _ := ctx.GetIndexValue(0)

	// 如果参数是 ASTValue，我们需要自己计算它的值
	if astValue, ok := varValue.(*data.ASTValue); ok {
		// 使用 Call 时的 Context 来计算值，这样可以捕获未定义变量的错误
		// 但是我们需要禁用错误抛出，因为 empty 应该抑制未定义变量错误
		// 这里我们假设 GetValue 返回 error control 表示变量未定义或其他错误
		val, acl := astValue.Node.GetValue(astValue.Ctx)
		if acl != nil {
			// 发生了错误（例如未定义变量 ReferenceError），empty 返回 true
			return data.NewBoolValue(true), nil
		}

		// 递归检查计算出的值
		return f.isEmptyValue(val), nil
	}

	if varValue == nil {
		return data.NewBoolValue(true), nil
	}

	return f.isEmptyValue(varValue), nil
}

// isEmptyValue 检查值是否为空
func (f *EmptyFunction) isEmptyValue(v data.GetValue) data.GetValue {
	if v == nil {
		return data.NewBoolValue(true)
	}

	// 检查是否为 NullValue
	if _, ok := v.(*data.NullValue); ok {
		return data.NewBoolValue(true)
	}

	// ... (rest of isEmptyValue logic)

	// 检查浮点数（FloatValue 同时实现了 AsInt、AsFloat 和 AsString）
	if floatVal, ok := v.(*data.FloatValue); ok {
		if floatVal.Value == 0.0 {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 检查整数（IntValue 同时实现了 AsFloat 和 AsString）
	if intVal, ok := v.(*data.IntValue); ok {
		if intVal.Value == 0 {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 检查字符串（StringValue）
	if strVal, ok := v.(*data.StringValue); ok {
		str := strVal.Value
		if str == "" || str == "0" {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 检查整数（排除 FloatValue，因为它已经处理过了）
	if intVal, ok := v.(data.AsInt); ok {
		// 如果是 FloatValue，已经处理过了，跳过
		if _, isFloat := v.(*data.FloatValue); isFloat {
			return data.NewBoolValue(false)
		}
		if i, err := intVal.AsInt(); err == nil && i == 0 {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 检查其他实现了 AsFloat 的类型
	if floatVal, ok := v.(data.AsFloat); ok {
		if f64, err := floatVal.AsFloat(); err == nil && f64 == 0.0 {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 检查布尔值
	if boolVal, ok := v.(data.AsBool); ok {
		if b, err := boolVal.AsBool(); err == nil && !b {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 检查数组
	if arrayVal, ok := v.(*data.ArrayValue); ok {
		if len(arrayVal.List) == 0 {
			return data.NewBoolValue(true)
		}
		return data.NewBoolValue(false)
	}

	// 其他类型，非空
	return data.NewBoolValue(false)
}

func (f *EmptyFunction) GetName() string {
	return "empty"
}

func (f *EmptyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterRawAST(nil, "var", 0, data.Mixed{}),
	}
}

func (f *EmptyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "var", 0, data.Mixed{}),
	}
}
