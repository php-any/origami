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
	if len(u.Args) == 0 {
		return data.NewNullValue(), nil
	}

	for _, argExpr := range u.Args {
		// 优先处理对象属性：unset($obj->prop)，需要绕过类型检查
		if callProp, ok := argExpr.(*CallObjectProperty); ok {
			objValue, acl := callProp.Object.GetValue(ctx)
			if acl != nil || objValue == nil {
				continue
			}
			switch obj := objValue.(type) {
			case *data.ObjectValue:
				obj.SetProperty(callProp.Property, data.NewNullValue())
			case *data.ClassValue:
				// 直接通过 ObjectValue 设置，绕过可能的类型检查
				obj.ObjectValue.SetProperty(callProp.Property, data.NewNullValue())
			case *data.ThisValue:
				obj.ClassValue.SetProperty(callProp.Property, data.NewNullValue())
			}
			continue
		}

		// 变量类型（普通变量 $var）
		if variable, ok := argExpr.(data.Variable); ok {
			variable.SetValue(ctx, data.NewNullValue())
			continue
		}

		// 索引表达式：unset($arr['key'])
		if indexExpr, ok := argExpr.(*IndexExpression); ok {
			arrayValue, acl := indexExpr.Array.GetValue(ctx)
			if acl != nil || arrayValue == nil {
				continue
			}
			indexValue, acl := indexExpr.Index.GetValue(ctx)
			if acl != nil || indexValue == nil {
				continue
			}
			switch arr := arrayValue.(type) {
			case *data.ArrayValue:
				if iv, ok := indexValue.(data.AsInt); ok {
					i, err := iv.AsInt()
					if err == nil && i >= 0 && i < len(arr.List) {
						arr.List[i] = data.NewZVal(data.NewNullValue())
					}
				}
			case *data.ObjectValue:
				if sv, ok := indexValue.(data.AsString); ok {
					arr.SetProperty(sv.AsString(), data.NewNullValue())
				}
			case *data.ClassValue:
				if sv, ok := indexValue.(data.AsString); ok {
					arr.SetProperty(sv.AsString(), data.NewNullValue())
				}
			}
		}
	}

	return data.NewNullValue(), nil
}
