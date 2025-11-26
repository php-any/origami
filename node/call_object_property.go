package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallObjectProperty 表示对象属性访问表达式
type CallObjectProperty struct {
	*Node    `pp:"-"`
	Object   data.GetValue // 对象表达式
	Property string        // 属性名
}

func (pe *CallObjectProperty) GetIndex() int {
	panic("不支持获取调用类属性过程获取属性索引")
}

func (pe *CallObjectProperty) GetName() string {
	return pe.Property
}

func (pe *CallObjectProperty) GetType() data.Types {
	return data.NewBaseType("")
}

func (pe *CallObjectProperty) SetValue(ctx data.Context, value data.Value) data.Control {
	temp, acl := pe.Object.GetValue(ctx)
	if acl != nil {
		return acl
	}
	switch object := temp.(type) {
	case *data.ClassValue: // 需要检查属性类型
		property, ok := object.GetProperty(pe.Property)
		if ok {
			if property.GetType() != nil && !property.GetType().Is(value) {
				return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("%s 属性 %s 因为类型不一致无法赋值", TryGetCallClassName(object), pe.Property))
			}
		}
		return object.SetProperty(pe.Property, value)
	case data.SetProperty:
		return object.SetProperty(pe.Property, value)
	default:
		return data.NewErrorThrow(pe.GetFrom(), errors.New("object is not set property"))
	}
}

// NewObjectProperty 创建一个新的对象属性访问表达式
func NewObjectProperty(token *TokenFrom, object data.GetValue, property string) *CallObjectProperty {
	return &CallObjectProperty{
		Node:     NewNode(token),
		Object:   object,
		Property: property,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallObjectProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	switch v := o.(type) {
	case *data.ThisValue:
		property, ok := v.Class.GetProperty(pe.Property)
		if ok {
			return property.GetValue(ctx)
		}
	case *data.ClassValue:
		property, ok := v.Class.GetProperty(pe.Property)
		if ok {
			if property.GetModifier() != data.ModifierPublic {
				return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("对象(%s)属性(%s)不是公开的", v.Class.GetName(), pe.Property)))
			}
			return property.GetValue(v)
		} else {
			property, ok := v.ObjectValue.GetProperty(pe.Property)
			if ok {
				return property.GetValue(v)
			}
		}
		return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("对象(%s)不存在属性(%s)", v.Class.GetName(), pe.Property)))
	case *data.ObjectValue:
		ov, has := v.GetProperty(pe.Property)
		if has {
			return ov.GetValue(v)
		}
	default:
		if obj, ok := v.(data.GetProperty); ok {
			ov, has := obj.GetProperty(pe.Property)
			if has {
				return ov.GetValue(ctx)
			}
		} else {
			return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("值(%s)不是对象, 不能操作属性(%s)", TryGetCallClassName(pe.Object), pe.Property)))
		}
	}
	return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("对象(%s)不存在属性(%s)", TryGetCallClassName(pe.Object), pe.Property)))
}
