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
	name := ""
	if getName, ok := pe.Stmt.(data.ClassStmt); ok {
		name = getName.GetName()
	}
	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("(%v)没有静态属性(%s)。", name, pe.Property)))
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
	cname := ""
	if getName, ok := pe.Stmt.(data.ClassStmt); ok {
		cname = getName.GetName()
	}
	return data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)没有静态属性(%s)。", cname, pe.Property)))
}

// CallStaticPropertyLater 延迟的静态属性访问（类未加载时）
type CallStaticPropertyLater struct {
	*Node
	className string // 类名（字符串形式）
	property  string // 属性名
	namespace string // 命名空间
}

// NewCallStaticPropertyLater 创建延迟的静态属性访问
func NewCallStaticPropertyLater(from *TokenFrom, className, property, namespace string) *CallStaticPropertyLater {
	return &CallStaticPropertyLater{
		Node:      NewNode(from),
		className: className,
		property:  property,
		namespace: namespace,
	}
}

// GetValue 获取延迟静态属性访问的值
func (pe *CallStaticPropertyLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 尝试加载类
	stmt, acl := ctx.GetVM().GetOrLoadClass(pe.className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		// 如果还是找不到，尝试使用命名空间
		fullClassName := pe.className
		if pe.namespace != "" {
			fullClassName = pe.namespace + "\\" + pe.className
		}
		stmt, acl = ctx.GetVM().GetOrLoadClass(fullClassName)
		if acl != nil {
			return nil, acl
		}
		if stmt == nil {
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法访问静态属性(%s::%s), 未找到类", pe.className, pe.property))
		}
	}

	// 创建实际的静态属性访问
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	callStaticProperty := NewCallStaticProperty(tokenFrom, stmt, pe.property)
	return callStaticProperty.GetValue(ctx)
}

// SetProperty 设置延迟静态属性的值
func (pe *CallStaticPropertyLater) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	// 尝试加载类
	stmt, acl := ctx.GetVM().GetOrLoadClass(pe.className)
	if acl != nil {
		return acl
	}
	if stmt == nil {
		// 如果还是找不到，尝试使用命名空间
		fullClassName := pe.className
		if pe.namespace != "" {
			fullClassName = pe.namespace + "\\" + pe.className
		}
		stmt, acl = ctx.GetVM().GetOrLoadClass(fullClassName)
		if acl != nil {
			return acl
		}
		if stmt == nil {
			return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法设置静态属性(%s::%s), 未找到类", pe.className, pe.property))
		}
	}

	// 创建实际的静态属性访问并设置值
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	callStaticProperty := NewCallStaticProperty(tokenFrom, stmt, pe.property)
	return callStaticProperty.SetProperty(ctx, name, value)
}
