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
		if ie, ok := argExpr.(*IndexExpression); ok {
			if isSet, handled := issetIndexExpression(ctx, ie); handled {
				if !isSet {
					return data.NewBoolValue(false), nil
				}
				continue
			}
		}

		varValue, ctl := argExpr.GetValue(ctx)
		if ctl != nil {
			if acl, ok := ctl.(data.GetName); ok && "UndefinedIndexExpression" == acl.GetName() {
				return data.NewBoolValue(false), nil
			}
			// PHP 中 isset 抑制 notice/warning 级错误（未定义属性、静态属性缺失等），
			// 但不抑制真正的 Exception 对象。origami 通过 Name=="Exception" 区分内部错误。
			if tv, isThrow := ctl.(*data.ThrowValue); isThrow && tv.Name == "Exception" {
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

// issetIndexExpression 检查数组/对象下标是否存在（不触发 PHP 8 未定义键 Warning）
func issetIndexExpression(ctx data.Context, ie *IndexExpression) (isSet bool, handled bool) {
	// isset($a[$k]['x'])：必须在 isSetViaOffsetExists 之前处理，因其会对 ie.Array 调用 GetValue
	if inner, ok := ie.Array.(*IndexExpression); ok {
		parentSet, handled := issetIndexExpression(ctx, inner)
		if !handled {
			return false, false
		}
		if !parentSet {
			return false, true
		}
		parentVal, ok := indexExpressionReadNoWarn(ctx, inner)
		if !ok {
			return false, true
		}
		index, acl := ie.Index.GetValue(ctx)
		if acl != nil {
			return false, true
		}
		return issetOnContainer(parentVal, index)
	}
	if isSet, handled := isSetViaOffsetExists(ctx, ie); handled {
		return isSet, true
	}
	container, acl := indexExpressionContainer(ctx, ie.Array)
	if acl != nil {
		return false, true
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return false, true
	}
	return issetOnContainer(container, index)
}

// indexExpressionContainer 求值下标表达式的“容器”部分（$a[$k] 中的 $a，或 $this->prop）
func indexExpressionContainer(ctx data.Context, arrayOperand data.GetValue) (data.GetValue, data.Control) {
	if inner, ok := arrayOperand.(*IndexExpression); ok {
		parentVal, ok := indexExpressionReadNoWarn(ctx, inner)
		if !ok {
			return data.NewNullValue(), nil
		}
		return parentVal, nil
	}
	return arrayOperand.GetValue(ctx)
}

// indexExpressionReadNoWarn 在已知 isset(IE) 为真时读取 IE 的值（不触发未定义键 Warning）
func indexExpressionReadNoWarn(ctx data.Context, ie *IndexExpression) (data.GetValue, bool) {
	container, acl := indexExpressionContainer(ctx, ie.Array)
	if acl != nil {
		return nil, false
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, false
	}
	val, ok := readIndexNoWarn(container, index)
	if !ok {
		return nil, false
	}
	return val, true
}

func issetOnContainer(container data.GetValue, index data.GetValue) (isSet bool, handled bool) {
	switch arr := container.(type) {
	case *data.StringValue:
		if iv, ok := index.(data.AsInt); ok {
			i, err := iv.AsInt()
			if err != nil {
				return false, true
			}
			if i < 0 {
				i = len(arr.Value) + i
			}
			return i >= 0 && i < len(arr.Value), true
		}
		return false, true
	case *data.ArrayValue:
		switch iv := index.(type) {
		case *data.StringValue:
			for _, z := range arr.List {
				if z != nil && z.Name == iv.Value {
					return z.Value != nil, true
				}
			}
			return false, true
		case data.AsInt:
			i, err := iv.AsInt()
			if err != nil {
				return false, true
			}
			z, _ := arr.FindSlotByIntKey(i)
			return z != nil && z.Value != nil, true
		}
	case *data.ObjectValue:
		key, ok := indexKeyString(index)
		if !ok {
			return false, true
		}
		val, acl := arr.GetProperty(key)
		if acl != nil {
			return false, true
		}
		_, isNull := val.(*data.NullValue)
		return val != nil && !isNull, true
	case data.GetProperty:
		if _, ok := arr.(*data.StringValue); ok {
			return false, false
		}
		key, ok := indexKeyString(index)
		if !ok {
			return false, true
		}
		val, acl := arr.GetProperty(key)
		if acl != nil {
			return false, true
		}
		_, isNull := val.(*data.NullValue)
		return val != nil && !isNull, true
	}
	return false, false
}

func readIndexNoWarn(container data.GetValue, index data.GetValue) (data.GetValue, bool) {
	switch arr := container.(type) {
	case *data.ArrayValue:
		switch iv := index.(type) {
		case *data.StringValue:
			for _, z := range arr.List {
				if z != nil && z.Name == iv.Value {
					if z.Value == nil {
						return data.NewNullValue(), true
					}
					return z.Value, true
				}
			}
			return nil, false
		case data.AsInt:
			i, err := iv.AsInt()
			if err != nil {
				return nil, false
			}
			z, _ := arr.FindSlotByIntKey(i)
			if z == nil {
				return nil, false
			}
			if z.Value == nil {
				return data.NewNullValue(), true
			}
			return z.Value, true
		}
	case *data.ObjectValue:
		key, ok := indexKeyString(index)
		if !ok {
			return nil, false
		}
		val, acl := arr.GetProperty(key)
		if acl != nil {
			return nil, false
		}
		return val, true
	case data.GetProperty:
		if _, ok := arr.(*data.StringValue); ok {
			return nil, false
		}
		key, ok := indexKeyString(index)
		if !ok {
			return nil, false
		}
		val, acl := arr.GetProperty(key)
		if acl != nil {
			return nil, false
		}
		return val, true
	}
	return nil, false
}

// indexExpressionKeyExists 判断下标槽位是否存在（与 isset 不同：值为 null 仍算存在）
func indexExpressionKeyExists(ctx data.Context, ie *IndexExpression) (exists bool, handled bool) {
	if inner, ok := ie.Array.(*IndexExpression); ok {
		parentExists, handled := indexExpressionKeyExists(ctx, inner)
		if !handled {
			return false, false
		}
		if !parentExists {
			return false, true
		}
		parentVal, ok := indexExpressionReadNoWarn(ctx, inner)
		if !ok {
			return false, true
		}
		index, acl := ie.Index.GetValue(ctx)
		if acl != nil {
			return false, true
		}
		return indexKeyExistsOnContainer(parentVal, index), true
	}
	if _, handled := isSetViaOffsetExists(ctx, ie); handled {
		return false, false
	}
	container, acl := indexExpressionContainer(ctx, ie.Array)
	if acl != nil {
		return false, true
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return false, true
	}
	return indexKeyExistsOnContainer(container, index), true
}

func indexKeyExistsOnContainer(container data.GetValue, index data.GetValue) bool {
	switch arr := container.(type) {
	case *data.StringValue:
		if iv, ok := index.(data.AsInt); ok {
			i, err := iv.AsInt()
			if err != nil {
				return false
			}
			if i < 0 {
				i = len(arr.Value) + i
			}
			return i >= 0 && i < len(arr.Value)
		}
		return false
	case *data.ArrayValue:
		switch iv := index.(type) {
		case *data.StringValue:
			for _, z := range arr.List {
				if z != nil && z.Name == iv.Value {
					return true
				}
			}
			return false
		case data.AsInt:
			i, err := iv.AsInt()
			if err != nil {
				return false
			}
			z, _ := arr.FindSlotByIntKey(i)
			return z != nil
		}
	case *data.ObjectValue:
		key, ok := indexKeyString(index)
		if !ok {
			return false
		}
		_, acl := arr.GetProperty(key)
		return acl == nil
	case data.GetProperty:
		if _, ok := arr.(*data.StringValue); ok {
			return false
		}
		key, ok := indexKeyString(index)
		if !ok {
			return false
		}
		_, acl := arr.GetProperty(key)
		return acl == nil
	}
	return false
}
