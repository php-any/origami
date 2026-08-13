package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ChdirFunction 实现 PHP 的 chdir 函数
// 签名近似：
//
//	chdir(string $directory): bool
//
// 当前实现：
//   - 仅接受可转换为字符串的第一个参数作为目录路径
//   - 使用 os.Chdir 修改当前工作目录
//   - 成功返回 true，失败返回 false（暂不触发 PHP 级别的 warning）
type ChdirFunction struct{}

func NewChdirFunction() data.FuncStmt {
	return &ChdirFunction{}
}

func (f *ChdirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dirValue, _ := ctx.GetIndexValue(0)
	if dirValue == nil {
		return data.NewBoolValue(false), nil
	}

	var dir string
	if str, ok := dirValue.(data.AsString); ok {
		dir = str.AsString()
	} else {
		dir = dirValue.AsString()
	}

	if dir == "" {
		return data.NewBoolValue(false), nil
	}

	if err := os.Chdir(dir); err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *ChdirFunction) GetName() string {
	return "chdir"
}

func (f *ChdirFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, data.String{}),
	}
}

func (f *ChdirFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.String{}),
	}
}
