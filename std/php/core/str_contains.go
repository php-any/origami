package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrContainsFunction 实现 str_contains 函数
//
// 使用 str_contains 检查字符串是否包含指定子串
//
// 语法: str_contains(string $haystack, string $needle): bool
//
// 参数:
//   - haystack: 要搜索的字符串
//   - needle: 要查找的子串
//
// 返回值: 如果 haystack 包含 needle 返回 true，否则返回 false
//
//	如果 needle 为空字符串，返回 true（PHP 8.0+ 行为）
//
// 使用示例:
//
//	str_contains('Hello World', 'World');  // true
//	str_contains('Hello World', 'PHP');   // false
//	str_contains('Hello World', '');       // true
type StrContainsFunction struct{}

func NewStrContainsFunction() data.FuncStmt {
	return &StrContainsFunction{}
}

func (f *StrContainsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：haystack（要搜索的字符串）
	haystackValue, _ := ctx.GetIndexValue(0)
	if haystackValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取第二个参数：needle（要查找的子串）
	needleValue, _ := ctx.GetIndexValue(1)
	if needleValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 转换为字符串
	haystack := haystackValue.AsString()
	needle := needleValue.AsString()

	// 如果 needle 为空字符串，PHP 8.0+ 返回 true
	if needle == "" {
		return data.NewBoolValue(true), nil
	}

	// 检查 haystack 是否包含 needle
	result := strings.Contains(haystack, needle)
	return data.NewBoolValue(result), nil
}

func (f *StrContainsFunction) GetName() string {
	return "str_contains"
}

func (f *StrContainsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, data.String{}),
		node.NewParameter(nil, "needle", 1, nil, data.String{}),
	}
}

func (f *StrContainsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.String{}),
		node.NewVariable(nil, "needle", 1, data.String{}),
	}
}
