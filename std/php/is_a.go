package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	"strings"
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
	source, _ := ctx.GetIndexValue(0)
	target, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil || target == "" {
		return data.NewBoolValue(false), nil
	}
	allowString, _ := utils.ConvertFromIndex[bool](ctx, 2)
	result, ctl := checkNominalRelation(ctx, source, target, allowString, false)
	return data.NewBoolValue(result), ctl
}

// Only the source string may autoload. Target lookup and ancestry comparisons
// are read-only; exceptions raised by an autoloader propagate to PHP.
func checkNominalRelation(ctx data.Context, source data.Value, target string, allowString, strict bool) (bool, data.Control) {
	vm := ctx.GetVM()
	var class data.ClassStmt
	switch value := source.(type) {
	case *data.ClassValue:
		class = value.Class
	case *data.ThisValue:
		class = value.Class
	case *data.ThrowValue:
		if value.Object != nil {
			class = value.Object.Class
		} else {
			name := value.Name
			if name == "" {
				name = "Exception"
			}
			return data.InternalTypeIsA(name, target) && (!strict || !data.TypeNameEqual(name, target)), nil
		}
	case *data.StringValue:
		if !allowString || vm == nil || value.Value == "" {
			return false, nil
		}
		name := strings.TrimPrefix(value.Value, "\\")
		var iface data.InterfaceStmt
		class, _ = vm.GetClass(name)
		iface, _ = vm.GetInterface(name)
		if class == nil && iface == nil {
			_, ctl := vm.LoadPkg(name)
			if ctl != nil {
				return false, ctl
			}
			class, _ = vm.GetClass(name)
			iface, _ = vm.GetInterface(name)
		}
		if iface != nil {
			return data.InterfaceIsA(iface.GetName(), target, vm) && (!strict || !data.TypeNameEqual(iface.GetName(), target)), nil
		}
	case *data.FuncValue, *data.BoundFuncValue:
		return !strict && data.TypeNameEqual(target, "Closure"), nil
	case data.Generator:
		return data.InternalTypeIsA("Generator", target) && (!strict || !data.TypeNameEqual(target, "Generator")), nil
	}
	if class == nil {
		return false, nil
	}
	if strict && data.SameNominalClass(class, target, vm) {
		return false, nil
	}
	return data.NominalIsA(class, target, vm), nil
}

func (f *IsAFunction) GetName() string {
	return "is_a"
}

var isAFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "object_or_class", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "class_name", 1, nil, data.String{}),
	node.NewParameter(nil, "allow_string", 2, data.NewBoolValue(false), data.Bool{}),
}

func (f *IsAFunction) GetParams() []data.GetValue {
	return isAFunctionGetParams
}

var isAFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "object_or_class", 0, data.Mixed{}),
	node.NewVariable(nil, "class_name", 1, data.String{}),
	node.NewVariable(nil, "allow_string", 2, data.Bool{}),
}

func (f *IsAFunction) GetVariables() []data.Variable {
	return isAFunctionGetVariables
}
