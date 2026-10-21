package php

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"os"
	"sync/atomic"
)

type directoryHandle struct {
	file *os.File
	dots int
}

var directoryResourceID atomic.Int64

type DirectoryFunction struct{ name string }

func (*DirectoryFunction) GetReturnType() data.Types { return nil }
func (f *DirectoryFunction) GetName() string         { return f.name }
func (f *DirectoryFunction) GetParams() []data.GetValue {
	if f.name == "opendir" {
		return []data.GetValue{node.NewParameter(nil, "directory", 0, nil, data.String{}), node.NewParameter(nil, "context", 1, data.NewNullValue(), nil)}
	}
	return []data.GetValue{node.NewParameter(nil, "dir_handle", 0, nil, nil)}
}
func (f *DirectoryFunction) GetVariables() []data.Variable {
	if f.name == "opendir" {
		return []data.Variable{node.NewVariable(nil, "directory", 0, data.String{}), node.NewVariable(nil, "context", 1, nil)}
	}
	return []data.Variable{node.NewVariable(nil, "dir_handle", 0, nil)}
}
func (f *DirectoryFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if f.name == "opendir" {
		file, err := os.Open(v.AsString())
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		info, err := file.Stat()
		if err != nil || !info.IsDir() {
			file.Close()
			return data.NewBoolValue(false), nil
		}
		h := &directoryHandle{file: file}
		return core.NewResourceValue(core.NewResourceClass("stream", h, int(directoryResourceID.Add(1))+400000), ctx), nil
	}
	var h *directoryHandle
	if resource, ok := v.(*core.ResourceValue); ok {
		h, _ = resource.GetResource().(*directoryHandle)
	}
	if h == nil || h.file == nil {
		return nil, data.NewTypeError(nil, fmt.Errorf("%s(): argument must be a valid directory resource", f.name))
	}
	switch f.name {
	case "closedir":
		h.file.Close()
		h.file = nil
		return data.NewNullValue(), nil
	case "rewinddir":
		if _, err := h.file.Seek(0, 0); err != nil {
			return data.NewBoolValue(false), nil
		}
		h.dots = 0
		return data.NewNullValue(), nil
	case "readdir":
		if h.dots < 2 {
			h.dots++
			if h.dots == 1 {
				return data.NewStringValue("."), nil
			}
			return data.NewStringValue(".."), nil
		}
		names, _ := h.file.Readdirnames(1)
		if len(names) == 0 {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(names[0]), nil
	}
	return data.NewNullValue(), nil
}
