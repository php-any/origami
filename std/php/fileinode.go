package php

import (
	"os"
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FileInodeFunction 实现 fileinode(string $filename): int|false
type FileInodeFunction struct{}

func NewFileInodeFunction() data.FuncStmt {
	return &FileInodeFunction{}
}

func (f *FileInodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathVal, _ := ctx.GetIndexValue(0)
	if pathVal == nil {
		return data.NewBoolValue(false), nil
	}
	path := pathVal.AsString()
	fi, err := os.Stat(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		return data.NewIntValue(int(st.Ino)), nil
	}
	// 非 Unix：用 ModTime 纳秒作弱替代，仅满足「有 inode 号」调用方
	return data.NewIntValue(int(fi.ModTime().UnixNano())), nil
}

func (f *FileInodeFunction) GetName() string { return "fileinode" }

func (f *FileInodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, data.String{}),
	}
}

func (f *FileInodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.String{}),
	}
}
