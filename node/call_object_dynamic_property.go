package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CallObjectDynamicProperty 表示 $obj->$name 或 $obj->{expr} 动态属性访问
// 运行时优先查找声明的属性，找不到时回退到 ArrayAccess/offsetGet 或 __get
type CallObjectDynamicProperty struct {
	*Node    `pp:"-"`
	Object   data.GetValue // 对象表达式
	NameExpr data.GetValue // 属性名表达式（求值后为字符串）
}

func NewCallObjectDynamicProperty(from data.From, object data.GetValue, nameExpr data.GetValue) *CallObjectDynamicProperty {
	return &CallObjectDynamicProperty{
		Node:     NewNode(from),
		Object:   object,
		NameExpr: nameExpr,
	}
}

// GetZVal 返回动态属性的共享 ZVal，供 by-ref 参数（如 data_set($obj->{$key}, ...)）写回。
func (pe *CallObjectDynamicProperty) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	temp, acl := pe.Object.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	raw, acl := pe.NameExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	name := raw.(data.Value).AsString()

	if object, ok := temp.(*data.ThisValue); ok {
		temp = object.ClassValue
	}
	if object, ok := temp.(*data.ClassValue); ok {
		return propertyReferenceSlot(ctx, object, name, pe.Node)
	}
	return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("Cannot reference property on %T", temp))
}

func (pe *CallObjectDynamicProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 求值对象表达式
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 求值属性名表达式
	raw, acl := pe.NameExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	name := raw.(data.Value).AsString()

	proxy := &CallObjectProperty{Node: pe.Node, Object: &literalGetValue{v: o}, Property: name}
	return proxy.GetValue(ctx)
}

func (pe *CallObjectDynamicProperty) SetValue(ctx data.Context, value data.Value) data.Control {
	_, ctl := pe.AssignValue(ctx, value)
	return ctl
}

func (pe *CallObjectDynamicProperty) AssignValue(ctx data.Context, value data.Value) (data.Value, data.Control) {
	temp, acl := pe.Object.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	raw, acl := pe.NameExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	name := raw.(data.Value).AsString()

	proxy := &CallObjectProperty{Node: pe.Node, Object: &literalGetValue{v: temp}, Property: name}
	return proxy.AssignValue(ctx, value)
}
