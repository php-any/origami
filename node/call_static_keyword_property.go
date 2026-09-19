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
	if hasStaticPropertySlot(startClass, pe.Property) {
		return startClass, nil
	}

	extend := startClass.GetExtend()
	for extend != nil {
		parentClass, acl := vm.GetOrLoadClass(*extend)
		if acl != nil {
			return nil, acl
		}

		if hasStaticPropertySlot(parentClass, pe.Property) {
			return parentClass, nil
		}

		extend = parentClass.GetExtend()
	}

	return startClass, nil
}

func hasStaticPropertySlot(class data.ClassStmt, name string) bool {
	switch c := class.(type) {
	case *ClassStatement:
		// 只看声明，不看运行时 StaticProperty.Load：子类若曾被误写入父类静态值，
		// Load 会让 Application 被当成 $instance 的定义类，HTTP overlay 写成
		// Application.instance，而 app()/csrf_token 走 Container::getInstance() 读不到。
		_, has := c.StaticProperties[name]
		return has
	case *AbstractClassStatement:
		_, has := c.StaticProperties[name]
		return has
	case *ClassGeneric:
		_, has := c.StaticProperties[name]
		return has
	default:
		if gsp, ok := class.(data.GetStaticProperty); ok {
			_, has := gsp.GetStaticProperty(name)
			return has
		}
	}
	return false
}

// GetValue 获取 static::$prop 访问的值
func (pe *CallStaticKeywordProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lateStaticClass, acl, ok := resolveLateStaticClass(ctx)
	if acl != nil {
		return nil, acl
	}
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	vm := ctx.GetVM()

	// 查找定义该属性的类（沿继承链向上）
	definingClass, acl := pe.findPropertyDefiningClass(vm, lateStaticClass)
	if acl != nil {
		return nil, acl
	}

	// 从定义该属性的类中获取静态属性（含接口类常量）
	if property, has := LookupStaticProperty(vm, definingClass, pe.Property); has {
		return property, nil
	}

	// PHP 兼容：未定义的静态属性访问返回 null（warning 级别，非致命错误）
	return data.NewNullValue(), nil
}

// SetProperty 设置 static::$prop 的值
func (pe *CallStaticKeywordProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	lateStaticClass, acl, ok := resolveLateStaticClass(ctx)
	if acl != nil {
		return acl
	}
	if !ok {
		return data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	vm := ctx.GetVM()

	// 查找定义该属性的类（沿继承链向上）
	definingClass, acl := pe.findPropertyDefiningClass(vm, lateStaticClass)
	if acl != nil {
		return acl
	}

	// 在定义该属性的类中设置静态属性
	if storeClassStatic(definingClass, name, value) {
		return nil
	}
	if sp, ok := definingClass.(data.SetProperty); ok {
		return sp.SetProperty(name, value)
	}

	cname := ""
	if getName, ok := definingClass.(data.ClassStmt); ok {
		cname = getName.GetName()
	}
	return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)没有静态属性(%s)。", cname, pe.Property))
}
