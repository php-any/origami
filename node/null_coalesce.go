package node

import (
	"github.com/php-any/origami/data"
)

// NullCoalesceExpression 表示空合并运算符表达式
type NullCoalesceExpression struct {
	*Node `pp:"-"`
	Left  data.GetValue // 左操作数
	Right data.GetValue // 右操作数
}

// NewNullCoalesceExpression 创建一个新的空合并运算符表达式
func NewNullCoalesceExpression(from *TokenFrom, left, right data.GetValue) *NullCoalesceExpression {
	return &NullCoalesceExpression{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

// GetValue 获取空合并运算符表达式的值
func (n *NullCoalesceExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if index, ok := n.Left.(*IndexExpression); ok {
		value, exists, ctl := ReadDimensionQuiet(ctx, index, true, true)
		if ctl != nil {
			return nil, ctl
		}
		if exists && value != nil {
			if _, null := value.(*data.NullValue); !null {
				return value, nil
			}
		}
		return n.Right.GetValue(ctx)
	}

	// PHP：$obj->prop ?? $b 使用 isset 语义（__isset 或属性存在性检查），
	// 属性不存在时不调用 __get，因此不会触发 Undefined array key Warning
	left := n.Left
	if cop, ok := left.(*CallObjectProperty); ok {
		frozen, ctl := freezeLvalue(ctx, cop)
		if ctl != nil {
			return nil, ctl
		}
		cop = frozen.(*CallObjectProperty)
		left = cop
		if exists, handled := coalesceObjectPropertyExists(ctx, cop); handled {
			if !exists {
				return n.Right.GetValue(ctx)
			}
			// 属性存在但值为 null：仍需走 GetValue（会正确返回 null 或 __get 值）
		}
	}

	// 计算左操作数的值
	leftValue, ctl := left.GetValue(ctx)
	if ctl != nil {
		if acl, ok := ctl.(data.GetName); ok && "UndefinedIndexExpression" == acl.GetName() {
			return n.Right.GetValue(ctx)
		}
		return nil, ctl
	}

	// 检查左操作数是否为 null
	if leftValue == nil {
		// 如果左操作数为 nil，返回右操作数
		return n.Right.GetValue(ctx)
	}

	// 检查是否为 NullValue 类型
	switch leftValue.(type) {
	case *data.NullValue:
		// 如果左操作数为 null，返回右操作数
		return n.Right.GetValue(ctx)
	default:
		// 如果左操作数不为 null，返回左操作数
		return leftValue, nil
	}
}

// AsString 返回空合并运算符表达式的字符串表示
func (n *NullCoalesceExpression) AsString() string {
	return "null_coalesce_expression"
}

// coalesceObjectPropertyExists 判断 $obj->prop ?? $default 中左操作数是否“已设置”
// （isset 语义：不调用 __get，因此不触发 Undefined array key Warning）
func coalesceObjectPropertyExists(ctx data.Context, pe *CallObjectProperty) (exists bool, handled bool) {
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		return false, true
	}
	var object *data.ClassValue
	switch owner := o.(type) {
	case *data.ClassValue:
		object = owner
	case *data.ThisValue:
		object = owner.ClassValue
	default:
		return false, false
	}
	if object == nil || object.ObjectValue == nil {
		return false, true
	}
	property, declared := lookupObjectProperty(ctx, object, pe.Property)
	if declared && propertyAccessible(ctx, object, property, pe.Property) {
		value, ctl := object.ObjectValue.GetProperty(data.PropertyStorageName(property))
		if ctl != nil || value == nil {
			return false, true
		}
		_, isNull := value.(*data.NullValue)
		return !isNull, true
	}
	if !declared && object.ObjectValue.HasProperty(pe.Property) {
		value, _ := object.ObjectValue.GetProperty(pe.Property)
		_, isNull := value.(*data.NullValue)
		return value != nil && !isNull, true
	}
	if magic, found := object.GetMethod("__isset"); found {
		isSet, ctl := pe.invokeMagicIsset(ctx, object, magic, pe.Property)
		return ctl == nil && isSet, true
	}
	return false, true
}
