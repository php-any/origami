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

func (pe *CallObjectProperty) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	temp, acl := pe.Object.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	switch object := temp.(type) {
	case data.GetPropertyStmt: // 需要检查属性类型
		property, ok := object.GetPropertyStmt(pe.Property)
		if ok {
			return property.GetZVal(object)
		}
	default:
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("object is not get property"))
	}
	return nil, nil
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
	case *data.ThisValue:
		property, ok := object.GetPropertyStmt(pe.Property)
		if ok {
			if property.GetType() != nil && !property.GetType().Is(value) {
				test := property.GetType()
				test.Is(value)
				return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("%s 属性 %s 因为类型不一致无法赋值", TryGetCallClassName(object), pe.Property))
			}
			return object.SetProperty(pe.Property, value)
		}
		// 无声明属性时尝试 __set(string $name, mixed $value)
		if magic, hasSet := object.GetMethod("__set"); hasSet {
			return pe.invokeMagicSet(object, magic, pe.Property, value)
		}
		return object.SetProperty(pe.Property, value)
	case *data.ClassValue:
		property, ok := object.GetPropertyStmt(pe.Property)
		if ok {
			if property.GetType() != nil && !property.GetType().Is(value) {
				return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("%s 属性 %s 因为类型不一致无法赋值", TryGetCallClassName(object), pe.Property))
			}
			if property.GetModifier() == data.ModifierPrivate {
				if !isCallerInClassHierarchy(ctx, object.Class) {
					return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("对象(%s)属性(%s)是私有的", object.Class.GetName(), pe.Property))
				}
			} else if property.GetModifier() == data.ModifierProtected {
				if !isCallerInClassHierarchy(ctx, object.Class) {
					return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("对象(%s)属性(%s)不是公开的", object.Class.GetName(), pe.Property))
				}
			}
			return object.SetProperty(pe.Property, value)
		}
		// 无声明属性时尝试 __set(string $name, mixed $value)
		if magic, hasSet := object.GetMethod("__set"); hasSet {
			return pe.invokeMagicSet(object, magic, pe.Property, value)
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

// invokeMagicGet 调用 __get(string $name)，用于读取不存在或不可见属性时的魔法分发
func (pe *CallObjectProperty) invokeMagicGet(object data.Context, magic data.Method, name string) (data.GetValue, data.Control) {
	varies := magic.GetVariables()
	if len(varies) < 1 {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__get 需要至少 1 个参数 (name)"))
	}
	fnCtx := object.CreateContext(varies)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(name))
	return magic.Call(fnCtx)
}

// invokeMagicSet 调用 __set(string $name, mixed $value)，用于写入不存在或不可见属性时的魔法分发
func (pe *CallObjectProperty) invokeMagicSet(object data.Context, magic data.Method, name string, value data.Value) data.Control {
	varies := magic.GetVariables()
	if len(varies) < 2 {
		return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__set 需要至少 2 个参数 (name, value)"))
	}
	fnCtx := object.CreateContext(varies)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(name))
	fnCtx.SetVariableValue(varies[1], value)
	_, acl := magic.Call(fnCtx)
	return acl
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallObjectProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		if ctl, ok := ctl.(data.AddStack); ok {
			ctl.AddStackWithInfo(pe.from, TryGetCallClassName(pe.Object), pe.Property)
		}
		return nil, ctl
	}
	switch v := o.(type) {
	case *data.NullValue:
		return data.NewNullValue(), nil
	case *data.ThisValue:
		property, ok := v.GetPropertyStmt(pe.Property)
		if ok {
			return property.GetValue(v)
		}
		// 无声明属性时尝试 __get(string $name)
		if magic, hasGet := v.GetMethod("__get"); hasGet {
			return pe.invokeMagicGet(v, magic, pe.Property)
		}
		return nil, data.NewErrorThrow(pe.from, fmt.Errorf("对象(%s)不存在属性(%s)", v.Class.GetName(), pe.Property))
	case *data.ClassValue:
		property, ok := v.GetPropertyStmt(pe.Property)
		if ok {
			if property.GetModifier() == data.ModifierPrivate {
				if !isCallerInClassHierarchy(ctx, v.Class) {
					return nil, data.NewErrorThrow(pe.from, fmt.Errorf("对象(%s)属性(%s)是私有的", v.Class.GetName(), pe.Property))
				}
			} else if property.GetModifier() == data.ModifierProtected {
				if !isCallerInClassHierarchy(ctx, v.Class) {
					return nil, data.NewErrorThrow(pe.from, fmt.Errorf("对象(%s)属性(%s)不是公开的", v.Class.GetName(), pe.Property))
				}
			}
			return property.GetValue(v)
		}
		// 无声明属性时尝试 __get(string $name)
		if magic, hasGet := v.GetMethod("__get"); hasGet {
			return pe.invokeMagicGet(v, magic, pe.Property)
		}
		dynVal, acl := v.ObjectValue.GetProperty(pe.Property)
		if acl != nil {
			return nil, acl
		}
		return dynVal, nil
	case data.GetProperty:
		ov, acl := v.GetProperty(pe.Property)
		if acl != nil {
			return nil, acl
		}
		return ov.GetValue(ctx)
	default:
		return nil, data.NewErrorThrow(pe.from, fmt.Errorf("值(%s)不是对象, 不能操作属性(%s)", TryGetCallClassName(pe.Object), pe.Property))
	}
}
