package node

import (
	"github.com/php-any/origami/data"
)

// UnsetStatement 表示 unset 语句
type UnsetStatement struct {
	*Node `pp:"-"`
	Args  []data.GetValue // 参数表达式列表
}

// NewUnsetStatement 创建一个新的 unset 语句
func NewUnsetStatement(token *TokenFrom, args []data.GetValue) *UnsetStatement {
	return &UnsetStatement{
		Node: NewNode(token),
		Args: args,
	}
}

// GetValue 获取 unset 语句的值
func (u *UnsetStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 如果没有参数，返回 null
	if len(u.Args) == 0 {
		return data.NewNullValue(), nil
	}

	// 遍历所有参数表达式，将每个变量设置为 null
	for _, argExpr := range u.Args {
		// 检查参数表达式是否是 Variable 类型
		if variable, ok := argExpr.(data.Variable); ok {
			// 直接将变量设置为 null
			ctl := variable.SetValue(ctx, data.NewNullValue())
			if ctl != nil {
				return nil, ctl
			}
		} else if indexExpr, ok := argExpr.(*IndexExpression); ok {
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
					if err == nil && i >= 0 && i < len(arr.Value) {
						// 删除元素（设置为 null）
						arr.Value[i] = data.NewNullValue()
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
		} else if callProp, ok := argExpr.(*CallObjectProperty); ok {
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
