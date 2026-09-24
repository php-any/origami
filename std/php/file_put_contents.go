package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// PHP file_put_contents flags (subset)
const (
	fileAppend = 8
)

func NewFilePutContentsFunction() data.FuncStmt {
	return &FilePutContentsFunction{}
}

type FilePutContentsFunction struct{}

func (f *FilePutContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathValue, _ := ctx.GetIndexValue(0)
	dataValue, _ := ctx.GetIndexValue(1)

	if pathValue == nil {
		return nil, utils.NewThrowf("FilePutContentsFunction called with no file path")
	}

	var filePath string
	if s, ok := pathValue.(data.AsString); ok {
		filePath = s.AsString()
	} else {
		filePath = pathValue.AsString()
	}

	if filePath == "" {
		return nil, utils.NewThrowf("FilePutContentsFunction called with no file path")
	}

	var content string
	if dataValue == nil {
		content = ""
	} else if s, ok := dataValue.(data.AsString); ok {
		content = s.AsString()
	} else {
		content = dataValue.AsString()
	}

	flags := 0
	if flagsVal, ok := ctx.GetIndexValue(2); ok && flagsVal != nil {
		if iv, ok := flagsVal.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				flags = n
			}
		}
	}

	var n int
	var err error
	if flags&fileAppend != 0 {
		n, err = appendFile(filePath, []byte(content))
	} else {
		err = os.WriteFile(filePath, []byte(content), 0644)
		n = len(content)
	}
	if err != nil {
		return nil, utils.NewThrowf("FilePutContentsFunction called with file path '%s': %v", filePath, err)
	}
	return data.NewIntValue(n), nil
}

func appendFile(path string, content []byte) (int, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(content)
}

func (f *FilePutContentsFunction) GetName() string { return "file_put_contents" }

var filePutContentsFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "filename", 0, nil, nil),
	node.NewParameter(nil, "data", 1, nil, nil),
	node.NewParameter(nil, "flags", 2, data.NewIntValue(0), nil),
	node.NewParameter(nil, "context", 3, nil, nil),
}

func (f *FilePutContentsFunction) GetParams() []data.GetValue {
	return filePutContentsFunctionGetParams
}

var filePutContentsFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "data", 1, data.NewBaseType("string")),
	node.NewVariable(nil, "flags", 2, data.NewBaseType("int")),
	node.NewVariable(nil, "context", 3, nil),
}

func (f *FilePutContentsFunction) GetVariables() []data.Variable {
	return filePutContentsFunctionGetVariables
}
