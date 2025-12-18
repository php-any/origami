package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// IndexExpression 表示数组访问表达式
type IndexExpression struct {
	*Node `pp:"-"`
	Array data.GetValue // 数组表达式
	Index data.GetValue // 索引表达式
}

// NewIndexExpression 创建一个新的数组访问表达式
func NewIndexExpression(token *TokenFrom, array data.GetValue, index data.GetValue) *IndexExpression {
	return &IndexExpression{
		Node:  NewNode(token),
		Array: array,
		Index: index,
	}
}

// GetValue 获取数组访问表达式的值
func (ie *IndexExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	temp, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	switch v := temp.(type) {
	case *data.ArrayValue:
		i := 0
		if iv, ok := index.(data.AsInt); ok {
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
		}
		if i >= len(v.Value) {
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("数组索引超出范围"))
		}

		return v.Value[i], nil
	case *data.ObjectValue:
		if iv, ok := index.(data.AsString); ok {
			ov, has := v.GetProperty(iv.AsString())
			if has {
				return ov, nil
			}
			return ov, nil
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ObjectValue无法处理索引的类型值"))
	case *data.ClassValue:
		// 支持对类实例通过字符串索引访问公开属性：
		// $obj[$name]，在动态属性语法 $obj->$name 降级为索引访问后会走到这里
		if iv, ok := index.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := v.GetProperty(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				return prop.GetValue(v)
			}
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ClassValue无法处理索引的类型值"))
	case *data.ThisValue:
		// $this[$name] 动态访问当前对象属性
		if iv, ok := index.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := v.Class.GetProperty(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				return prop.GetValue(ctx)
			}
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ThisValue无法处理索引的类型值"))
	case *data.StringValue:
		// 获取字符串指定位置符号
		if iv, ok := index.(data.AsInt); ok {
			var err error
			i, err := iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.Value) {
				return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("字符串索引超出范围"))
			}
			return data.NewStringValue(string(v.Value[i])), nil
		} else {
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("字符串无法处理非int值"))
		}
	}
	return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
}
