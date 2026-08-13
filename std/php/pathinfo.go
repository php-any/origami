package php

import (
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PathinfoFunction 实现 pathinfo 全局函数
type PathinfoFunction struct{}

// NewPathinfoFunction 创建一个新的 pathinfo 函数实例
func NewPathinfoFunction() data.FuncStmt {
	return &PathinfoFunction{}
}

// GetName 返回函数名
func (f *PathinfoFunction) GetName() string {
	return "pathinfo"
}

// GetParams 返回参数列表
func (f *PathinfoFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "options", 1, nil, data.Mixed{}),
	}
}

// GetVariables 返回变量列表
func (f *PathinfoFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.Mixed{}),
		node.NewVariable(nil, "options", 1, data.Mixed{}),
	}
}

// Call 执行 pathinfo 函数
// pathinfo 用于返回文件路径的信息
func (f *PathinfoFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 path 参数
	pathValue, _ := ctx.GetIndexValue(0)
	if pathValue == nil {
		return data.NewStringValue(""), nil
	}

	path := ""
	if s, ok := pathValue.(data.AsString); ok {
		path = s.AsString()
	} else {
		path = pathValue.AsString()
	}

	// 获取 options 参数（可选）
	optionsValue, _ := ctx.GetIndexValue(1)

	// 解析路径信息
	dirname := filepath.Dir(path)
	basename := filepath.Base(path)
	ext := filepath.Ext(basename)
	filename := basename
	if ext != "" {
		filename = basename[:len(basename)-len(ext)]
	}

	// 如果指定了选项，只返回对应的部分
	if optionsValue != nil {
		if optVal, ok := optionsValue.(data.AsInt); ok {
			options, _ := optVal.AsInt()
			if options == 1 { // PATHINFO_DIRNAME
				return data.NewStringValue(dirname), nil
			} else if options == 2 { // PATHINFO_BASENAME
				return data.NewStringValue(basename), nil
			} else if options == 4 { // PATHINFO_EXTENSION
				if len(ext) > 1 {
					return data.NewStringValue(ext[1:]), nil // 去掉前面的点
				}
				return data.NewStringValue(""), nil
			} else if options == 8 { // PATHINFO_FILENAME
				return data.NewStringValue(filename), nil
			}
		}
	}

	// 默认返回 filename
	return data.NewStringValue(filename), nil
}
