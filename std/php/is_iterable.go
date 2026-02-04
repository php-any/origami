package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsIterableFunction 实现 is_iterable 函数
// 规则与 PHP 保持一致：
// - array 一定是 iterable
// - 实现 Traversable（此处近似为 data.Iterator 或 data.Generator）的对象是 iterable
// - 其他类型一律返回 false
type IsIterableFunction struct{}

func NewIsIterableFunction() data.FuncStmt {
	return &IsIterableFunction{}
}

func (f *IsIterableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	// 1. 数组一定是可迭代的
	if _, ok := value.(*data.ArrayValue); ok {
		return data.NewBoolValue(true), nil
	}

	// 2. 生成器（Generator）是 Traversable
	if _, ok := value.(data.Generator); ok {
		return data.NewBoolValue(true), nil
	}

	// 3. 脚本类：只有实现 Traversable（Iterator / IteratorAggregate）才算 iterable
	if obj, ok := value.(*data.ClassValue); ok && obj != nil {
		if cls := obj.Class; cls != nil {
			for _, impl := range cls.GetImplements() {
				if isTraversableInterface(impl) {
					return data.NewBoolValue(true), nil
				}
			}
		}
		// 普通对象（含 stdClass）不视为 iterable
		return data.NewBoolValue(false), nil
	}

	// 4. 其他类型一律不是 iterable
	return data.NewBoolValue(false), nil
}

func (f *IsIterableFunction) GetName() string {
	return "is_iterable"
}

func (f *IsIterableFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsIterableFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

// isTraversableInterface 判断接口名是否为 Traversable 相关
// 支持：
// - Iterator / \Iterator / 完整命名空间结尾为 \Iterator
// - IteratorAggregate / \IteratorAggregate / ...\IteratorAggregate
// - Traversable / \Traversable / ...\Traversable
func isTraversableInterface(name string) bool {
	if name == "" {
		return false
	}
	trimmed := strings.TrimPrefix(name, "\\")

	switch trimmed {
	case "Iterator", "IteratorAggregate", "Traversable":
		return true
	}

	if strings.HasSuffix(trimmed, "\\Iterator") ||
		strings.HasSuffix(trimmed, "\\IteratorAggregate") ||
		strings.HasSuffix(trimmed, "\\Traversable") {
		return true
	}

	return false
}


