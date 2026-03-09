package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallSelfProperty 表示当前类静态属性访问表达式
type CallSelfProperty struct {
	*Node    `pp:"-"`
	Property string // 属性名
}

func NewCallSelfProperty(from data.From, property string) *CallSelfProperty {
	return &CallSelfProperty{
		Node:     NewNode(from),
		Property: property,
	}
}

// GetValue 获取当前类静态属性访问表达式的值
func (pe *CallSelfProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("self:: 只能在类方法中使用"))
	}

	// 获取当前类
	currentClass := classCtx.Class

	// 先检查当前类是否实现了 GetStaticProperty 接口
	getter, ok := currentClass.(data.GetStaticProperty)
	if ok {
		// 获取当前类的静态属性
		property, has := getter.GetStaticProperty(pe.Property)
		if has {
			return property, nil
		}
	}

	vm := ctx.GetVM()

	// 首先在当前类及其所有父类中查找静态属性/常量
	extend := currentClass.GetExtend()
	for extend != nil {
		parentClass, acl := vm.GetOrLoadClass(*extend)
		if acl != nil {
			return nil, acl
		}

		// 检查父类是否实现了 GetStaticProperty 接口
		if parentGetter, ok := parentClass.(data.GetStaticProperty); ok {
			property, has := parentGetter.GetStaticProperty(pe.Property)
			if has {
				return property, nil
			}
		}

		extend = parentClass.GetExtend()
	}

	// 然后在当前类及其所有父类实现的接口中查找常量
	classToCheck := currentClass
	for classToCheck != nil {
		implements := classToCheck.GetImplements()
		for _, interfaceName := range implements {
			// 递归查找接口及其所有父接口的常量
			property, acl := pe.findInInterfaceAndParents(vm, interfaceName)
			if acl != nil {
				return nil, acl
			}
			if property != nil {
				return property, nil
			}
		}
		// 向上遍历父类继承链
		if classToCheck.GetExtend() == nil {
			break
		}
		parentClass, acl := vm.GetOrLoadClass(*classToCheck.GetExtend())
		if acl != nil {
			return nil, acl
		}
		classToCheck = parentClass
	}

	// 所有地方都没有找到，返回错误
	return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 及其父类和接口都没有静态属性或常量 %s", currentClass.GetName(), pe.Property))
}

// findInInterfaceAndParents 递归查找接口及其所有父接口的常量
func (pe *CallSelfProperty) findInInterfaceAndParents(vm data.VM, interfaceName string) (data.Value, data.Control) {
	interfaceStmt, acl := vm.GetOrLoadInterface(interfaceName)
	if acl != nil {
		return nil, acl
	}

	// 检查当前接口是否实现了 GetStaticProperty 接口
	if interfaceGetter, ok := interfaceStmt.(data.GetStaticProperty); ok {
		property, has := interfaceGetter.GetStaticProperty(pe.Property)
		if has {
			return property, nil
		}
	}

	// 递归检查接口的父接口
	for _, parentName := range interfaceStmt.GetExtends() {
		property, acl := pe.findInInterfaceAndParents(vm, parentName)
		if acl != nil {
			return nil, acl
		}
		if property != nil {
			return property, nil
		}
	}

	return nil, nil
}

// SetProperty 设置当前类静态属性的值
func (pe *CallSelfProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewErrorThrow(pe.GetFrom(), errors.New("self:: 只能在类方法中使用"))
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
	cname := currentClass.GetName()
	return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)没有静态属性(%s)", cname, pe.Property))
}
