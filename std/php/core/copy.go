package core

import (
	"io"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CopyFunction 实现 copy 函数
// 对标 PHP copy: 拷贝文件，成功返回 true，失败返回 false（不抛异常）
type CopyFunction struct{}

func NewCopyFunction() data.FuncStmt {
	return &CopyFunction{}
}

func (f *CopyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取源文件和目标文件
	srcValue, _ := ctx.GetIndexValue(0)
	dstValue, _ := ctx.GetIndexValue(1)
	if srcValue == nil || dstValue == nil {
		return data.NewBoolValue(false), nil
	}

	var src, dst string
	if s, ok := srcValue.(data.AsString); ok {
		src = s.AsString()
	} else {
		src = srcValue.AsString()
	}
	if s, ok := dstValue.(data.AsString); ok {
		dst = s.AsString()
	} else {
		dst = dstValue.AsString()
	}

	if src == "" || dst == "" {
		return data.NewBoolValue(false), nil
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *CopyFunction) GetName() string {
	return "copy"
}

func (f *CopyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "source", 0, nil, nil),
		node.NewParameter(nil, "dest", 1, nil, nil),
	}
}

func (f *CopyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "source", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "dest", 1, data.NewBaseType("string")),
	}
}
