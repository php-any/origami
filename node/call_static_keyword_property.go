package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticKeywordProperty 表示 static::$prop （late static binding 风格）的静态属性访问表达式
// 注意：当前实现语义上仍等同于 self::$prop，但通过单独节点类型与 self:: 区分，便于后续增强
type CallStaticKeywordProperty struct {
	*Node    `pp:"-"`
	Property string // 属性名
}

func NewCallStaticKeywordProperty(from data.From, property string) *CallStaticKeywordProperty {
	return &CallStaticKeywordProperty{
		Node:     NewNode(from),
		Property: property,
	}
}

// GetValue 获取 static::$prop 访问的值
func (pe *CallStaticKeywordProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 与 self:: 一样，必须在类方法上下文中使用
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	// 获取当前类
	currentClass := classCtx.Class

	// 检查类是否实现了 GetStaticProperty 接口
	getter, ok := currentClass.(data.GetStaticProperty)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 不支持静态属性访问", currentClass.GetName()))
	}

	// 获取当前类的静态属性
	property, has := getter.GetStaticProperty(pe.Property)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态属性 %s", currentClass.GetName(), pe.Property))
	}

	return property, nil
}

// SetProperty 设置 static::$prop 的值
func (pe *CallStaticKeywordProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	// 与 self:: 一样，必须在类方法上下文中使用
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	// 获取当前类
	currentClass := classCtx.Class

	switch c := currentClass.(type) {
	case *ClassStatement:
		c.StaticProperty.Store(name, value)
		return nil
	case *ClassGeneric:
		c.StaticProperty.Store(name, value)
		return nil
	case data.SetProperty:
		return c.SetProperty(name, value)
	}

	cname := ""
	if getName, ok := currentClass.(data.ClassStmt); ok {
		cname = getName.GetName()
	}
	return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)没有静态属性(%s)。", cname, pe.Property))
}
