package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticKeywordProperty 表示 static::$prop （late static binding 风格）的静态属性访问表达式
// 实现了 PHP 的后期静态绑定语义：static::$property 会查找定义该属性的类（沿继承链向上），
// 并从定义该属性的类中读取属性值
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

// findPropertyDefiningClass 沿继承链查找定义了指定静态属性的类
// 这是实现后期静态绑定的关键：static::$property 应该访问定义该属性的类中的属性
func (pe *CallStaticKeywordProperty) findPropertyDefiningClass(vm data.VM, startClass data.ClassStmt) (data.ClassStmt, data.Control) {
	// 首先检查当前类是否定义了该静态属性
	if cs, ok := startClass.(*ClassStatement); ok {
		if _, has := cs.StaticProperty.Load(pe.Property); has {
			return startClass, nil
		}
	}

	// 沿继承链向上查找
	extend := startClass.GetExtend()
	for extend != nil {
		parentClass, acl := vm.GetOrLoadClass(*extend)
		if acl != nil {
			return nil, acl
		}

		// 检查父类是否定义了该静态属性
		if cs, ok := parentClass.(*ClassStatement); ok {
			if _, has := cs.StaticProperty.Load(pe.Property); has {
				return parentClass, nil
			}
		}

		extend = parentClass.GetExtend()
	}

	// 没有找到定义该属性的类，返回起始类（让后续代码报错）
	return startClass, nil
}

// GetValue 获取 static::$prop 访问的值
func (pe *CallStaticKeywordProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 与 self:: 一样，必须在类上下文中使用
	var lateStaticClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		lateStaticClass = classCtx.Class
		if classCtx.StaticClass != nil {
			lateStaticClass = classCtx.StaticClass
		}
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		lateStaticClass = classVal.Class
	} else {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	vm := ctx.GetVM()

	// 查找定义该属性的类（沿继承链向上）
	definingClass, acl := pe.findPropertyDefiningClass(vm, lateStaticClass)
	if acl != nil {
		return nil, acl
	}

	// 从定义该属性的类中获取静态属性
	getter, ok := definingClass.(data.GetStaticProperty)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类 %s 不支持静态属性访问", definingClass.GetName()))
	}

	property, has := getter.GetStaticProperty(pe.Property)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类 %s 没有静态属性 %s", definingClass.GetName(), pe.Property))
	}

	return property, nil
}

// SetProperty 设置 static::$prop 的值
func (pe *CallStaticKeywordProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	// 与 self:: 一样，必须在类上下文中使用
	var lateStaticClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		lateStaticClass = classCtx.Class
		if classCtx.StaticClass != nil {
			lateStaticClass = classCtx.StaticClass
		}
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		lateStaticClass = classVal.Class
	} else {
		return data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	vm := ctx.GetVM()

	// 查找定义该属性的类（沿继承链向上）
	definingClass, acl := pe.findPropertyDefiningClass(vm, lateStaticClass)
	if acl != nil {
		return acl
	}

	// 在定义该属性的类中设置静态属性
	switch c := definingClass.(type) {
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
	if getName, ok := definingClass.(data.ClassStmt); ok {
		cname = getName.GetName()
	}
	return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)没有静态属性(%s)。", cname, pe.Property))
}
