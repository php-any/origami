package core

import (
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DirnameFunction 实现 dirname 函数
type DirnameFunction struct{}

func NewDirnameFunction() data.FuncStmt {
	return &DirnameFunction{}
}

func (f *DirnameFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathValue, _ := ctx.GetIndexValue(0)

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		path = pathValue.AsString()
	}

	dir := filepath.Dir(path)

	return data.NewStringValue(dir), nil
}

func (f *DirnameFunction) GetName() string {
	return "dirname"
}

func (f *DirnameFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, nil),
	}
}

func (f *DirnameFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
	}
}
