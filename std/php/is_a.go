package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// NewIsAFunction 创建 is_a 函数
// PHP 语义：
// is_a(mixed $object_or_class, string $class_name, bool $allow_string = false): bool
// 检查对象是否为指定类的实例或子类实例
// $allow_string = true 时，第一个参数可以是类名字符串
func NewIsAFunction() data.FuncStmt {
	return &IsAFunction{}
}

type IsAFunction struct{}

func (f *IsAFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	objectOrClass, ok := ctx.GetIndexValue(0)
	if !ok || objectOrClass == nil {
		return data.NewBoolValue(false), nil
	}

	className, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil || className == "" {
		return data.NewBoolValue(false), nil
	}

	allowString, _ := utils.ConvertFromIndex[bool](ctx, 2)

	vm := ctx.GetVM()

	// 如果是对象实例
	if classValue, ok := objectOrClass.(*data.ClassValue); ok {
		// 加载目标类/接口
		target, acl := vm.LoadPkg(className)
		if acl != nil {
			return data.NewBoolValue(false), nil
		}
		if target == nil {
			return data.NewBoolValue(false), nil
		}
		var targetName string
		switch t := target.(type) {
		case data.ClassStmt:
			targetName = t.GetName()
		case data.InterfaceStmt:
			targetName = t.GetName()
		default:
			return data.NewBoolValue(false), nil
		}
		result, acl := checkClassIsHelper(ctx, classValue.Class, targetName)
		if acl != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewBoolValue(result), nil
	}

	// 如果是 ThisValue
	if thisValue, ok := objectOrClass.(*data.ThisValue); ok {
		target, acl := vm.LoadPkg(className)
		if acl != nil {
			return data.NewBoolValue(false), nil
		}
		if target == nil {
			return data.NewBoolValue(false), nil
		}
		var targetName string
		switch t := target.(type) {
		case data.ClassStmt:
			targetName = t.GetName()
		case data.InterfaceStmt:
			targetName = t.GetName()
		default:
			return data.NewBoolValue(false), nil
		}
		result, acl := checkClassIsHelper(ctx, thisValue.Class, targetName)
		if acl != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewBoolValue(result), nil
	}

	// 如果 allow_string = true，第一个参数可以是类名字符串
	if allowString {
		if strVal, ok := objectOrClass.(*data.StringValue); ok {
			sourceClassName := strVal.Value
			if sourceClassName == "" {
				return data.NewBoolValue(false), nil
			}
			sourceStmt, acl := vm.GetOrLoadClass(sourceClassName)
			if acl != nil || sourceStmt == nil {
				return data.NewBoolValue(false), nil
			}
			result, acl := checkClassIsHelper(ctx, sourceStmt, className)
			if acl != nil {
				return data.NewBoolValue(false), nil
			}
			return data.NewBoolValue(result), nil
		}
	}

	return data.NewBoolValue(false), nil
}

// checkClassIsHelper 检查 source 类是否是 target 类/接口的实例或子类
// 复用 node.checkClassIs 的逻辑（通过接口遍历）
func checkClassIsHelper(ctx data.Context, source data.ClassStmt, target string) (bool, data.Control) {
	if source == nil {
		return false, nil
	}
	// 类名相同
	if source.GetName() == target {
		return true, nil
	}

	vm := ctx.GetVM()

	// 检查实现的接口
	for _, impl := range source.GetImplements() {
		if impl == target {
			return true, nil
		}
		if ifaceStmt, ok := vm.GetInterface(impl); ok {
			// 检查接口继承
			for _, parentIface := range ifaceStmt.GetExtends() {
				if parentIface == target {
					return true, nil
				}
			}
		}
	}

	// 检查父类继承链
	if source.GetExtend() != nil {
		extName := *source.GetExtend()
		if extName == target {
			return true, nil
		}
		parentStmt, acl := vm.GetOrLoadClass(extName)
		if acl != nil {
			return false, acl
		}
		if parentStmt != nil {
			return checkClassIsHelper(ctx, parentStmt, target)
		}
	}

	return false, nil
}

func (f *IsAFunction) GetName() string {
	return "is_a"
}

func (f *IsAFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "class_name", 1, nil, data.String{}),
		node.NewParameter(nil, "allow_string", 2, data.NewBoolValue(false), data.Bool{}),
	}
}

func (f *IsAFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.Mixed{}),
		node.NewVariable(nil, "class_name", 1, data.String{}),
		node.NewVariable(nil, "allow_string", 2, data.Bool{}),
	}
}
