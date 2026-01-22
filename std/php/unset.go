package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewUnsetFunction() data.FuncStmt {
	return &UnsetFunction{}
}

type UnsetFunction struct{}

func (f *UnsetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取实际传递的参数表达式列表
	callArgs := ctx.GetCallArgs()
	if len(callArgs) == 0 {
		// 没有参数，返回 null
		return data.NewNullValue(), nil
	}

	// 遍历所有参数表达式，将每个变量设置为 null
	for _, argExpr := range callArgs {
		// 检查参数表达式是否是 Variable 类型
		if variable, ok := argExpr.(data.Variable); ok {
			// 直接将变量设置为 null
			ctl := variable.SetValue(ctx, data.NewNullValue())
			if ctl != nil {
				return nil, ctl
			}
		} else if indexExpr, ok := argExpr.(*node.IndexExpression); ok {
			// 处理数组元素或对象属性：unset($arr['key']) 或 unset($obj->prop)
			// 获取数组/对象
			arrayValue, acl := indexExpr.Array.GetValue(ctx)
			if acl != nil {
				continue // 跳过错误
			}

			// 获取索引值
			indexValue, acl := indexExpr.Index.GetValue(ctx)
			if acl != nil {
				continue // 跳过错误
			}

			if indexValue == nil {
				continue
			}

			// 根据数组/对象类型处理
			switch arr := arrayValue.(type) {
			case *data.ArrayValue:
				// 数组元素删除
				if iv, ok := indexValue.(data.AsInt); ok {
					// 整数索引
					i, err := iv.AsInt()
					if err == nil && i >= 0 && i < len(arr.List) {
						// 删除元素（设置为 null）
						arr.List[i] = data.NewZVal(data.NewNullValue())
					}
				}
				// 注意：ArrayValue 不支持字符串索引，字符串索引应该使用 ObjectValue
			case *data.ObjectValue:
				// 对象属性删除
				if sv, ok := indexValue.(data.AsString); ok {
					propName := sv.AsString()
					// 设置为 null（实际上相当于删除）
					arr.SetProperty(propName, data.NewNullValue())
				}
			case *data.ClassValue:
				// 类对象属性删除
				if sv, ok := indexValue.(data.AsString); ok {
					propName := sv.AsString()
					// 设置为 null（实际上相当于删除）
					arr.SetProperty(propName, data.NewNullValue())
				}
			}
		} else if callProp, ok := argExpr.(*node.CallObjectProperty); ok {
			// 处理对象属性访问：unset($obj->prop)
			// 获取对象
			objValue, acl := callProp.Object.GetValue(ctx)
			if acl != nil {
				continue // 跳过错误
			}

			if objValue == nil {
				continue
			}

			// 根据对象类型处理
			switch obj := objValue.(type) {
			case *data.ObjectValue:
				// 对象属性设置为 null
				obj.SetProperty(callProp.Property, data.NewNullValue())
			case *data.ClassValue:
				// 类对象属性设置为 null
				obj.SetProperty(callProp.Property, data.NewNullValue())
			}
		}
		// 其他类型的参数表达式暂时不做处理
	}

	// unset 返回 null
	return data.NewNullValue(), nil
}

func (f *UnsetFunction) GetName() string {
	return "unset"
}

func (f *UnsetFunction) GetParams() []data.GetValue {
	// unset 可以接受可变数量的参数
	// 使用 CallerContextParameter 来在调用者上下文中执行，以便获取实际参数
	return []data.GetValue{
		node.NewCallerContextParameter(nil),
	}
}

func (f *UnsetFunction) GetVariables() []data.Variable {
	// unset 的参数是动态的，这里返回空数组
	return []data.Variable{}
}
