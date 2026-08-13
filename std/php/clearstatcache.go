package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewClearstatcacheFunction() data.FuncStmt {
	return &ClearstatcacheFunction{}
}

type ClearstatcacheFunction struct{}

func (f *ClearstatcacheFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// Go 没有文件 stat 缓存，此函数为 no-op
	return nil, nil
}

func (f *ClearstatcacheFunction) GetName() string {
	return "clearstatcache"
}

func (f *ClearstatcacheFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "clear_realpath_cache", 0, data.NewBoolValue(false), nil),
		node.NewParameter(nil, "filename", 1, data.NewStringValue(""), nil),
	}
}

func (f *ClearstatcacheFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "clear_realpath_cache", 0, data.NewBaseType("bool")),
		node.NewVariable(nil, "filename", 1, data.NewBaseType("string")),
	}
}
