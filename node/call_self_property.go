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

	// 如果当前类没有，检查父类
	extend := currentClass.GetExtend()
	vm := ctx.GetVM()
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

		// 继续向上查找父类
		extend = parentClass.GetExtend()
	}

	// 如果父类也没有，检查实现的接口
	implements := currentClass.GetImplements()
	for _, interfaceName := range implements {
		// 递归查找接口及其所有父接口的常量
		property := pe.findInInterfaceAndParents(vm, interfaceName)
		if property != nil {
			return property, nil
		}
	}

	// 所有地方都没有找到，返回错误
	return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 及其父类和接口都没有静态属性或常量 %s", currentClass.GetName(), pe.Property))
}

// findInInterfaceAndParents 递归查找接口及其所有父接口的常量
func (pe *CallSelfProperty) findInInterfaceAndParents(vm data.VM, interfaceName string) data.Value {
	interfaceStmt, acl := vm.GetOrLoadInterface(interfaceName)
	if acl != nil || interfaceStmt == nil {
		return nil
	}

	// 检查当前接口是否实现了 GetStaticProperty 接口
	if interfaceGetter, ok := interfaceStmt.(data.GetStaticProperty); ok {
		property, has := interfaceGetter.GetStaticProperty(pe.Property)
		if has {
			return property
		}
	}

	// 递归检查接口的父接口
	interfaceExtend := interfaceStmt.GetExtend()
	if interfaceExtend != nil {
		return pe.findInInterfaceAndParents(vm, *interfaceExtend)
	}

	return nil
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
