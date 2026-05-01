package node

import (
	"errors"
	"fmt"

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
	return nil, data.NewErrorThrow(v.from, errors.New("只能引用变量"))
}

// resolveIndexRef 解析 &$array[$key] 或 &$array[] 引用
func (v *ValueReference) resolveIndexRef(ctx data.Context, ie *IndexExpression) (data.GetValue, data.Control) {
	arrVal, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	idxVal, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	switch arr := arrVal.(type) {
	case *data.ArrayValue:
		idx, ok := idxVal.(data.Value)
		if !ok {
			idx = data.NewNullValue()
		}
		i, err := toArrayIndex(idx)
		if err != nil {
			return nil, data.NewErrorThrow(v.from, err)
		}
		// &$array[] 推空槽位
		if i == len(arr.List) {
			arr.List = append(arr.List, data.NewZVal(data.NewNullValue()))
		}
		if i < 0 || i >= len(arr.List) {
			return nil, data.NewErrorThrow(v.from, errors.New("数组索引超出范围"))
		}
		return &data.ArraySlotRef{Arr: arr, Idx: i}, nil
	case *data.ObjectValue:
		key := ""
		if sv, ok := idxVal.(data.AsString); ok {
			key = sv.AsString()
		} else if iv, ok := idxVal.(data.AsInt); ok {
			i, _ := iv.AsInt()
			key = fmt.Sprintf("%d", i)
		} else {
			key = idxVal.(data.Value).AsString()
		}
		// 确保槽位存在
		if _, acl := arr.GetProperty(key); acl != nil {
			arr.SetProperty(key, data.NewNullValue())
		}
		// ObjectValue 用 IndexReferenceValue 保证 SetValue 行为
		return data.NewIndexReferenceValue(ie, ctx), nil
	default:
		return nil, data.NewErrorThrow(v.from, errors.New("只能引用数组或对象的索引"))
	}
}

// toArrayIndex 将索引值转为数组整数下标
func toArrayIndex(idx data.Value) (int, error) {
	if iv, ok := idx.(data.AsInt); ok {
		return iv.AsInt()
	}
	return 0, errors.New("数组索引必须是整数")
}
