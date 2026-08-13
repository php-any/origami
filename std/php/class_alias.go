package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewClassAliasFunction 创建 class_alias 函数
// PHP 语义：
//
//	class_alias(string $original, string $alias, bool $autoload = true): bool
//
// 在 VM 中为别名类名注册一个代理 ClassStmt，并按 $autoload 参数加载原类。
// 若原类不存在或别名已被其它类占用，则返回 false。
func NewClassAliasFunction() data.FuncStmt {
	return &ClassAliasFunction{}
}

type ClassAliasFunction struct{}

func (f *ClassAliasFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	origVal, ctl := ctx.GetVariableValue(node.NewVariable(nil, "original", 0, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}
	aliasVal, ctl := ctx.GetVariableValue(node.NewVariable(nil, "alias", 1, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}

	original := origVal.(data.AsString).AsString()
	alias := aliasVal.(data.AsString).AsString()

	vm := ctx.GetVM()

	autoload := true
	if autoloadVal, ok := ctx.GetIndexValue(2); ok && autoloadVal != nil {
		if value, ok := autoloadVal.(data.AsBool); ok {
			if parsed, err := value.AsBool(); err == nil {
				autoload = parsed
			}
		}
	}

	originalClass, ok := vm.GetClass(original)
	if !ok && autoload {
		loaded, control := vm.GetOrLoadClass(original)
		if control == nil && loaded != nil {
			originalClass, ok = loaded, true
		}
	}
	if !ok || originalClass == nil {
		return data.NewBoolValue(false), nil
	}

	if _, exists := vm.GetClass(alias); exists {
		return data.NewBoolValue(false), nil
	}
	if _, exists := vm.GetInterface(alias); exists {
		return data.NewBoolValue(false), nil
	}

	if control := vm.AddClass(&classAliasStmt{alias: alias, original: originalClass}); control != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *ClassAliasFunction) GetName() string {
	return "class_alias"
}

func (f *ClassAliasFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "original", 0, nil, data.String{}),
		node.NewParameter(nil, "alias", 1, nil, data.String{}),
		node.NewParameter(nil, "autoload", 2, data.NewBoolValue(true), data.NewBaseType("bool")),
	}
}

func (f *ClassAliasFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "original", 0, data.String{}),
		node.NewVariable(nil, "alias", 1, data.String{}),
		node.NewVariable(nil, "autoload", 2, data.NewBaseType("bool")),
	}
}

// classAliasStmt 只改变 VM 中的查找名称，其余行为委托给原类。
type classAliasStmt struct {
	alias    string
	original data.ClassStmt
}

func (c *classAliasStmt) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return c.original.GetValue(ctx)
}
func (c *classAliasStmt) GetFrom() data.From      { return c.original.GetFrom() }
func (c *classAliasStmt) GetName() string         { return c.alias }
func (c *classAliasStmt) GetExtend() *string      { return c.original.GetExtend() }
func (c *classAliasStmt) GetImplements() []string { return c.original.GetImplements() }
func (c *classAliasStmt) GetProperty(name string) (data.Property, bool) {
	return c.original.GetProperty(name)
}
func (c *classAliasStmt) GetPropertyList() []data.Property { return c.original.GetPropertyList() }
func (c *classAliasStmt) GetMethod(name string) (data.Method, bool) {
	return c.original.GetMethod(name)
}
func (c *classAliasStmt) GetMethods() []data.Method { return c.original.GetMethods() }
func (c *classAliasStmt) GetConstruct() data.Method { return c.original.GetConstruct() }

func (c *classAliasStmt) GetStaticMethod(name string) (data.Method, bool) {
	if class, ok := c.original.(data.GetStaticMethod); ok {
		return class.GetStaticMethod(name)
	}
	return nil, false
}

func (c *classAliasStmt) GetStaticProperty(name string) (data.Value, bool) {
	if class, ok := c.original.(data.GetStaticProperty); ok {
		return class.GetStaticProperty(name)
	}
	return nil, false
}
