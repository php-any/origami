package node

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
)

type CallStaticProperty struct {
	*Node    `pp:"-"`
	Stmt     data.GetValue
	Property string // 属性名
}

func NewCallStaticProperty(token *TokenFrom, stmt data.GetValue, property string) *CallStaticProperty {
	return &CallStaticProperty{
		Node:     NewNode(token),
		Stmt:     stmt,
		Property: property,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	switch expr := pe.Stmt.(type) {
	case data.GetStaticProperty:
		property, ok := expr.GetStaticProperty(pe.Property)
		if ok {
			return property, nil
		}

		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("无法调用属性(%s)。", pe.Property)))
	default:
		next, acl := pe.Stmt.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		switch expr := next.(type) {
		case *data.ClassValue:
			if c, ok := expr.Class.(data.GetStaticProperty); ok {
				property, ok := c.GetStaticProperty(pe.Property)
				if ok {
					return property, nil
				}
			}

		case data.GetStaticProperty:
			property, ok := expr.GetStaticProperty(pe.Property)
			if ok {
				return property, nil
			}
		}
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("没有静态属性(%s)。", pe.Property)))
}

func (pe *CallStaticProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	switch c := pe.Stmt.(type) {
	case *ClassStatement:
		c.StaticProperty.Store(name, value)
		return nil
	case *ClassGeneric:
		c.StaticProperty.Store(name, value)
		return nil
	case data.SetProperty:
		return c.SetProperty(name, value)
	default:
		c, acl := pe.Stmt.GetValue(ctx)
		if acl != nil {
			return acl
		}
		switch c := c.(type) {
		case data.SetProperty:
			return c.SetProperty(name, value)
		}
	}

	return data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)没有静态属性(%s)。", pe.Stmt, pe.Property)))
}
