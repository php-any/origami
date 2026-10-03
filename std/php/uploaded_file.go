package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"os"
)

type UploadedFileFunction struct{ name string }

func (f *UploadedFileFunction) GetName() string { return f.name }
func (f *UploadedFileFunction) GetParams() []data.GetValue {
	params := []data.GetValue{node.NewParameter(nil, "filename", 0, nil, data.String{})}
	if f.name == "move_uploaded_file" {
		params = append(params, node.NewParameter(nil, "to", 1, nil, data.String{}))
	}
	return params
}
func (f *UploadedFileFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "filename", 0, nil), node.NewVariable(nil, "to", 1, nil)}
}
func (f *UploadedFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	path := value.AsString()
	host, ok := ctx.GetVM().(interface {
		IsUploadedFile(string) bool
		ForgetUploadedFile(string)
	})
	if !ok || !host.IsUploadedFile(path) {
		return data.NewBoolValue(false), nil
	}
	if _, err := os.Stat(path); err != nil {
		return data.NewBoolValue(false), nil
	}
	if f.name == "is_uploaded_file" {
		return data.NewBoolValue(true), nil
	}
	target, _ := ctx.GetIndexValue(1)
	if err := os.Rename(path, target.AsString()); err != nil {
		return data.NewBoolValue(false), nil
	}
	host.ForgetUploadedFile(path)
	return data.NewBoolValue(true), nil
}
