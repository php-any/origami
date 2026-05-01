package node

import (
	"github.com/php-any/origami/data"
)

// IssetStatement 表示 isset 语句
type IssetStatement struct {
	*Node `pp:"-"`
	Args  []data.GetValue // 参数表达式列表
}

// NewIssetStatement 创建一个新的 isset 语句
func NewIssetStatement(token *TokenFrom, args []data.GetValue) *IssetStatement {
	return &IssetStatement{
		Node: NewNode(token),
		Args: args,
	}
}

// GetValue 获取 isset 语句的值
func (i *IssetStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if len(i.Args) == 0 {
		return data.NewBoolValue(false), nil
	}

	for _, argExpr := range i.Args {
		// 对 IndexExpression，尝试使用 ArrayAccess::offsetExists
		if ie, ok := argExpr.(*IndexExpression); ok {
			if isSet, handled := isSetViaOffsetExists(ctx, ie); handled {
				return data.NewBoolValue(isSet), nil
			}
		}

		varValue, ctl := argExpr.GetValue(ctx)
		if ctl != nil {
			if acl, ok := ctl.(data.GetName); ok && "UndefinedIndexExpression" == acl.GetName() {
				return data.NewBoolValue(false), nil
			}
			return nil, ctl
		}
		if varValue == nil {
			return data.NewBoolValue(false), nil
		}
		if _, isNull := varValue.(*data.NullValue); isNull {
			return data.NewBoolValue(false), nil
		}
	}

	return data.NewBoolValue(true), nil
}

// isSetViaOffsetExists 对实现了 ArrayAccess 的对象使用 offsetExists 检查
func isSetViaOffsetExists(ctx data.Context, ie *IndexExpression) (bool, bool) {
	array, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return false, false
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return false, false
	}

	var obj *data.ClassValue
	switch v := array.(type) {
	case *data.ClassValue:
		obj = v
	case *data.ThisValue:
		obj = v.ClassValue
	default:
		return false, false
	}

	method, exists := obj.GetMethod("offsetExists")
	if !exists {
		return false, false
	}

	fnCtx := obj.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		if iv, ok := index.(data.Value); ok {
			fnCtx.SetVariableValue(method.GetVariables()[0], iv)
		}
	}
	ret, ctl := method.Call(fnCtx)
	if ctl != nil {
		return false, false
	}
	if bv, ok := ret.(*data.BoolValue); ok {
		return bv.Value, true
	}
	return false, true
}
