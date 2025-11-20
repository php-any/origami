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
	// 获取参数值（引用参数）
	varValue, _ := ctx.GetIndexValue(0)

	// 检查参数值是否是 IndexReferenceValue（数组元素引用）
	if indexRefValue, ok := varValue.(*data.IndexReferenceValue); ok {
		// 处理数组元素访问
		parentCtx := indexRefValue.Ctx
		indexExpr := indexRefValue.Expr

		// 类型断言为 IndexExpression
		// indexExpr 本身就是 IndexExpression，不需要调用 GetValue
		indexExpression, ok := indexExpr.(*node.IndexExpression)
		if !ok {
			return data.NewBoolValue(true), nil // 无法解析，视为空
		}

		// 获取数组/对象
		arrayValue, acl := indexExpression.Array.GetValue(parentCtx)
		if acl != nil {
			return data.NewBoolValue(true), nil // 数组不存在，视为空
		}

		if arrayValue == nil {
			return data.NewBoolValue(true), nil
		}

		// 获取索引值
		indexValue, acl := indexExpression.Index.GetValue(parentCtx)
		if acl != nil {
			return data.NewBoolValue(true), nil
		}

		if indexValue == nil {
			return data.NewBoolValue(true), nil
		}

		// 检查数组或对象中是否存在该键
		switch arr := arrayValue.(type) {
		case *data.ArrayValue:
			// 数组索引访问
			if iv, ok := indexValue.(data.AsInt); ok {
				// 整数索引
				i, err := iv.AsInt()
				if err != nil {
					return data.NewBoolValue(true), nil
				}
				// 检查索引是否在范围内
				if i < 0 || i >= len(arr.Value) {
					return data.NewBoolValue(true), nil // 索引越界，视为空
				}
				// 获取元素值
				elementValue := arr.Value[i]
				return f.isEmptyValue(elementValue), nil
			} else if sv, ok := indexValue.(data.AsString); ok {
				// 字符串索引（关联数组），使用 GetProperty
				indexStr := sv.AsString()
				val, exists := arr.GetProperty(indexStr)
				if !exists {
					return data.NewBoolValue(true), nil // 键不存在，视为空
				}
				return f.isEmptyValue(val), nil
			} else {
				return data.NewBoolValue(true), nil
			}
		case *data.ObjectValue:
			// 对象属性访问
			if sv, ok := indexValue.(data.AsString); ok {
				propValue, has := arr.GetProperty(sv.AsString())
				if !has {
					return data.NewBoolValue(true), nil // 属性不存在，视为空
				}
				// 即使 has 为 true，propValue 也可能是 NullValue（表示属性存在但值为 null）
				return f.isEmptyValue(propValue), nil
			}
			return data.NewBoolValue(true), nil
		}

		return data.NewBoolValue(true), nil
	}

	// 检查参数值是否是 ReferenceValue（变量引用）
	if refValue, ok := varValue.(*data.ReferenceValue); ok {
		parentCtx := refValue.Ctx
		varRef := refValue.Val

		v, acl := varRef.GetValue(parentCtx)
		if acl != nil {
			// 变量不存在，视为空
			return data.NewBoolValue(true), nil
		}

		if v == nil {
			return data.NewBoolValue(true), nil
		}

		internalV, intervalCtl := v.GetValue(parentCtx)
		if intervalCtl != nil {
			return data.NewBoolValue(true), nil
		}

		return f.isEmptyValue(internalV), nil
	}

	// 如果不是引用类型，直接检查值
	if varValue == nil {
		return data.NewBoolValue(true), nil
	}

	internalV, intervalCtl := varValue.GetValue(ctx)
	if intervalCtl != nil {
		return data.NewBoolValue(true), nil
	}

	return f.isEmptyValue(internalV), nil
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

	// 先检查具体类型，避免接口类型匹配的顺序问题
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
		if len(arrayVal.Value) == 0 {
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
		node.NewParameterReference(nil, "var", 0, data.Mixed{}),
	}
}

func (f *EmptyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "var", 0, data.Mixed{}),
	}
}
