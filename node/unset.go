package node

import (
	"fmt"

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
			// PHP 8.1: null is treated as empty string (no deprecation for unset)
			if _, isNull := indexValue.(*data.NullValue); isNull {
				indexValue = data.NewStringValue("")
			}
			switch arr := arrayValue.(type) {
			case *data.ArrayValue:
				if iv, ok := indexValue.(data.Value); ok {
					arr.UnsetKey(iv)
				}
				writeBackArrayProperty(ctx, indexExpr.Array, arr)
			case *data.ObjectValue:
				if sv, ok := indexValue.(data.AsString); ok {
					arr.UnsetProperty(sv.AsString())
				} else if iv, ok := indexValue.(data.AsInt); ok {
					if i, err := iv.AsInt(); err == nil {
						arr.UnsetProperty(fmt.Sprintf("%d", i))
					}
				}
			case *data.ClassValue:
				if iv, ok := indexValue.(data.Value); ok && CheckArrayAccess(ctx, arr.Class) {
					if ctl := CallArrayAccessOffsetUnset(ctx, arr, iv); ctl != nil {
						return nil, ctl
					}
					continue
				}
				if sv, ok := indexValue.(data.AsString); ok {
					arr.SetProperty(sv.AsString(), data.NewNullValue())
				}
			case *data.ThisValue:
				if arr.ClassValue != nil && CheckArrayAccess(ctx, arr.Class) {
					if iv, ok := indexValue.(data.Value); ok {
						if ctl := CallArrayAccessOffsetUnset(ctx, arr.ClassValue, iv); ctl != nil {
							return nil, ctl
						}
					}
				}
			}
		}
	}

	return data.NewNullValue(), nil
}

func writeBackArrayProperty(ctx data.Context, arrayExpr data.GetValue, arr *data.ArrayValue) {
	switch a := arrayExpr.(type) {
	case *CallObjectProperty:
		obj, acl := a.Object.GetValue(ctx)
		if acl != nil {
			return
		}
		if cv, ok := obj.(*data.ClassValue); ok {
			cv.SetProperty(a.Property, arr)
		} else if tv, ok := obj.(*data.ThisValue); ok && tv.ClassValue != nil {
			tv.ClassValue.SetProperty(a.Property, arr)
		}
	case *IndexExpression:
		writeBackArrayProperty(ctx, a.Array, arr)
	}
}
