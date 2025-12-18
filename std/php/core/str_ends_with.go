package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrEndsWithFunction 实现 str_ends_with 函数
//
// 使用 str_ends_with 检查字符串是否以指定子串结尾
//
// 语法: str_ends_with(string $haystack, string $needle): bool
//
// 参数:
//   - haystack: 要搜索的字符串
//   - needle: 要查找的子串
//
// 返回值: 如果 haystack 以 needle 结尾返回 true，否则返回 false
//
//	如果 needle 为空字符串，返回 true（PHP 8.0+ 行为）
//
// 使用示例:
//
//	str_ends_with('Hello World', 'World');  // true
//	str_ends_with('Hello World', 'Hello');  // false
//	str_ends_with('Hello World', '');       // true
type StrEndsWithFunction struct{}

func NewStrEndsWithFunction() data.FuncStmt {
	return &StrEndsWithFunction{}
}

func (f *StrEndsWithFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 检查 haystack 是否以 needle 结尾
	result := strings.HasSuffix(haystack, needle)
	return data.NewBoolValue(result), nil
}

func (f *StrEndsWithFunction) GetName() string {
	return "str_ends_with"
}

func (f *StrEndsWithFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, data.String{}),
		node.NewParameter(nil, "needle", 1, nil, data.String{}),
	}
}

func (f *StrEndsWithFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.String{}),
		node.NewVariable(nil, "needle", 1, data.String{}),
	}
}
