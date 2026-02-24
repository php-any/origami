package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// safeAsString 将 Value 转为字符串，避免 AsString() 内部 panic 导致进程退出
func safeAsString(v data.Value) string {
	if v == nil {
		return ""
	}
	defer func() { _ = recover() }()
	return v.AsString()
}

// StrStartsWithFunction 实现 str_starts_with 函数
//
// 使用 str_starts_with 检查字符串是否以指定子串开头
//
// 语法: str_starts_with(string $haystack, string $needle): bool
//
// 参数:
//   - haystack: 要搜索的字符串
//   - needle: 要查找的子串
//
// 返回值: 如果 haystack 以 needle 开头返回 true，否则返回 false
//
//	如果 needle 为空字符串，返回 true（PHP 8.0+ 行为）
//
// 使用示例:
//
//	str_starts_with('Hello World', 'Hello');  // true
//	str_starts_with('Hello World', 'World'); // false
//	str_starts_with('Hello World', '');      // true
type StrStartsWithFunction struct{}

func NewStrStartsWithFunction() data.FuncStmt {
	return &StrStartsWithFunction{}
}

func (f *StrStartsWithFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 转换为字符串（使用安全转换，避免某些 Value 实现 AsString 时 panic）
	haystack := safeAsString(haystackValue)
	needle := safeAsString(needleValue)

	// 如果 needle 为空字符串，PHP 8.0+ 返回 true
	if needle == "" {
		return data.NewBoolValue(true), nil
	}

	// 检查 haystack 是否以 needle 开头
	result := strings.HasPrefix(haystack, needle)
	return data.NewBoolValue(result), nil
}

func (f *StrStartsWithFunction) GetName() string {
	return "str_starts_with"
}

func (f *StrStartsWithFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, data.String{}),
		node.NewParameter(nil, "needle", 1, nil, data.String{}),
	}
}

func (f *StrStartsWithFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.String{}),
		node.NewVariable(nil, "needle", 1, data.String{}),
	}
}
