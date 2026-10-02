package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// ValueReference 表示引用取值表达式 &$var
type ValueReference struct {
	*Node `pp:"-"`
	Value data.GetValue
}

// NewValueReference 创建一个新的引用取值表达式
func NewValueReference(token *TokenFrom, value data.GetValue) *ValueReference {
	return &ValueReference{
		Node:  NewNode(token),
		Value: value,
	}
}

// GetValue 获取引用的值
func (v *ValueReference) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 变量引用：&$var
	if variable, ok := v.Value.(data.Variable); ok {
		return data.NewReferenceValue(variable, ctx), nil
	}
	// 数组/对象属性索引引用：&$array[$key] 或 &$array[]
	if ie, ok := v.Value.(*IndexExpression); ok {
		return v.resolveIndexRef(ctx, ie)
	}
	// 按引用返回的函数调用：&getRef($x)
	if call, ok := v.Value.(*CallExpression); ok {
		ret, ctl := call.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if ref, ok := ret.(data.Value); ok {
			switch ref.(type) {
			case *data.ReferenceValue, *data.ArraySlotRef, *data.IndexReferenceValue:
				return ref, nil
			}
		}
		return nil, data.NewErrorThrow(v.from, errors.New("只能引用变量"))
	}
	if call, ok := v.Value.(*CallLater); ok {
		ret, ctl := call.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if ref, ok := ret.(data.Value); ok {
			switch ref.(type) {
			case *data.ReferenceValue, *data.ArraySlotRef, *data.IndexReferenceValue:
				return ref, nil
			}
		}
		return nil, data.NewErrorThrow(v.from, errors.New("只能引用变量"))
	}
	return nil, data.NewErrorThrow(v.from, errors.New("只能引用变量"))
}

// resolveIndexRef binds the same stable slot used by by-reference arguments.
func (v *ValueReference) resolveIndexRef(ctx data.Context, ie *IndexExpression) (data.GetValue, data.Control) {
	slot, ctl := ie.GetOrCreateZVal(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if slot == nil {
		return nil, data.NewErrorThrow(v.from, errors.New("只能引用数组或对象的索引"))
	}
	return &data.ArraySlotRef{Slot: slot}, nil
}
