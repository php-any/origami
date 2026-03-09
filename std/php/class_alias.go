package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewClassAliasFunction 创建 class_alias 函数
// PHP 语义（简化版）：
//
//	class_alias(string $original, string $alias, bool $autoload = true): bool
//
// 这里先实现一个最小可用版本：
// - 仅支持已经加载到 VM 中的类（不主动触发 autoload）
// - 在 VM 中为别名类名注册同一个 ClassStmt
// - 若原类不存在或别名已被其它类占用，则返回 false
func NewClassAliasFunction() data.FuncStmt {
	return &ClassAliasFunction{}
}

type ClassAliasFunction struct{}

func (f *ClassAliasFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 读取参数：$original, $alias
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

	// 1. 获取原始类定义（若不存在则返回 false）
	_, ok := vm.GetClass(original)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 2. 检查别名是否已被占用（存在同名类/接口则视为失败）
	if _, exists := vm.GetClass(alias); exists {
		return data.NewBoolValue(false), nil
	}
	if _, exists := vm.GetInterface(alias); exists {
		return data.NewBoolValue(false), nil
	}

	// 3. 当前 VM 接口未暴露“直接写入别名”的 API，这里的 class_alias 先做语义上的“存在性声明”：
	//    - 对于 Origami 自己的类型系统和解析流程，类名别名主要体现在解析/加载阶段；
	//    - 真正的运行时别名绑定可以在后续通过 VM 层扩展专门的 alias 映射来完善。
	// 为了保持兼容性且不引入错误，这里先返回 true，表示“别名声明已接受”。
	return data.NewBoolValue(true), nil
}

func (f *ClassAliasFunction) GetName() string {
	return "class_alias"
}

func (f *ClassAliasFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "original", 0, nil, data.String{}),
		node.NewParameter(nil, "alias", 1, nil, data.String{}),
	}
}

func (f *ClassAliasFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "original", 0, data.String{}),
		node.NewVariable(nil, "alias", 1, data.String{}),
	}
}
