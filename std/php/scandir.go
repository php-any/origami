package php

import (
	"errors"
	"os"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewScandirFunction() data.FuncStmt {
	return &ScandirFunction{}
}

type ScandirFunction struct{}

func (f *ScandirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取路径参数
	pathValue, _ := ctx.GetIndexValue(0)

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 检查路径是否为空
	if path == "" {
		return nil, data.NewErrorThrow(nil, errors.New("scandir函数调用: path 是空的"))
	}

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	defer dir.Close()

	// 读取目录内容
	entries, err := dir.Readdirnames(0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	// 对文件名进行排序（PHP scandir 默认按字母顺序排序）
	sort.Strings(entries)

	// 将文件名转换为字符串值数组
	var values []data.Value
	for _, entry := range entries {
		values = append(values, data.NewStringValue(entry))
	}

	return data.NewArrayValue(values), nil
}

func (f *ScandirFunction) GetName() string {
	return "scandir"
}

func (f *ScandirFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, nil),
	}
}

func (f *ScandirFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.NewBaseType("string")),
	}
}
