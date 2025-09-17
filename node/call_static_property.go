package node

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
)

type CallStaticProperty struct {
	*Node    `pp:"-"`
	Class    data.ClassStmt
	Property string // 属性名
}

func NewCallStaticProperty(token *TokenFrom, stmt data.ClassStmt, property string) *CallStaticProperty {
	return &CallStaticProperty{
		Node:     NewNode(token),
		Class:    stmt,
		Property: property,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if o, ok := pe.Class.(data.GetStaticProperty); ok {
		property, ok := o.GetStaticProperty(pe.Property)
		if ok {
			return property, nil
		}
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)没有静态属性(%s)。", pe.Class, pe.Property)))
}

func (pe *CallStaticProperty) SetProperty(name string, value data.Value) data.Control {
	switch c := pe.Class.(type) {
	case *ClassStatement:
		c.StaticProperty.Store(name, value)
		return nil
	case *ClassGeneric:
		c.StaticProperty.Store(name, value)
		return nil
	case data.SetProperty:
		return c.SetProperty(name, value)
	}

	return data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)没有静态属性(%s)。", pe.Class, pe.Property)))
}
