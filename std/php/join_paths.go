package php

import (
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// JoinPathsFunction 实现 join_paths 全局函数
type JoinPathsFunction struct{}

// NewJoinPathsFunction 创建一个新的 join_paths 函数实例
func NewJoinPathsFunction() data.FuncStmt {
	return &JoinPathsFunction{}
}

// GetName 返回函数名
func (f *JoinPathsFunction) GetName() string {
	return "join_paths"
}

// GetParams 返回参数列表
func (f *JoinPathsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "basePath", 0, nil, data.Mixed{}),
		node.NewParameters(nil, "paths", 1, nil, data.Mixed{}),
	}
}

// GetVariables 返回变量列表
func (f *JoinPathsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "basePath", 0, data.Mixed{}),
		node.NewVariable(nil, "paths", 1, data.Mixed{}),
	}
}

// Call 执行 join_paths 函数
// join_paths 用于连接路径，类似于 PHP 的 DIRECTORY_SEPARATOR 处理
// 签名：join_paths(string $basePath, string ...$paths): string
func (f *JoinPathsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 basePath 参数
	basePathValue, _ := ctx.GetIndexValue(0)
	if basePathValue == nil {
		return data.NewStringValue(""), nil
	}

	basePath := ""
	if s, ok := basePathValue.(data.AsString); ok {
		basePath = s.AsString()
	} else {
		basePath = basePathValue.AsString()
	}

	// 获取 paths 参数（可变参数数组）
	pathsValue, _ := ctx.GetIndexValue(1)

	// 将所有路径部分收集到一个切片中
	parts := []string{basePath}

	// 如果提供了 paths 数组，遍历它添加所有路径部分
	if pathsValue != nil {
		if arrayVal, ok := pathsValue.(*data.ArrayValue); ok {
			for _, zval := range arrayVal.List {
				if zval != nil && zval.Value != nil {
					if str, ok := zval.Value.(data.AsString); ok {
						parts = append(parts, str.AsString())
					} else if zval.Value != nil {
						parts = append(parts, zval.Value.AsString())
					}
				}
			}
		}
	}

	// 使用 filepath.Join 连接所有路径
	result := filepath.Join(parts...)

	// 将 Windows 风格的路径转换为 Unix 风格（如果需要）
	result = strings.ReplaceAll(result, "\\", "/")

	return data.NewStringValue(result), nil
}
