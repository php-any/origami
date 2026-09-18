package php

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	fileIgnoreNewLines = 2
	fileSkipEmptyLines = 4
)

// FileFunction 实现 PHP file()：把文件读成行数组。
func NewFileFunction() data.FuncStmt {
	return &FileFunction{}
}

type FileFunction struct{}

func (f *FileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathVal, _ := ctx.GetIndexValue(0)
	if pathVal == nil {
		return data.NewBoolValue(false), nil
	}
	path := pathVal.AsString()
	if path == "" {
		return data.NewBoolValue(false), nil
	}

	flags := 0
	if fv, ok := ctx.GetIndexValue(1); ok && fv != nil {
		if asInt, ok := fv.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				flags = n
			}
		}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewArrayValue(phpFileLines(string(raw), flags)), nil
}

func phpFileLines(s string, flags int) []data.Value {
	ignoreNL := flags&fileIgnoreNewLines != 0
	skipEmpty := flags&fileSkipEmptyLines != 0
	if s == "" {
		return []data.Value{data.NewStringValue("")}
	}
	out := make([]data.Value, 0, 64)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			continue
		}
		line := s[start : i+1]
		start = i + 1
		if ignoreNL {
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")
		}
		if skipEmpty && phpFileLineEmpty(line, ignoreNL) {
			continue
		}
		out = append(out, data.NewStringValue(line))
	}
	if start < len(s) {
		line := s[start:]
		if ignoreNL {
			line = strings.TrimSuffix(line, "\r")
		}
		if !(skipEmpty && phpFileLineEmpty(line, ignoreNL)) {
			out = append(out, data.NewStringValue(line))
		}
	}
	return out
}

func phpFileLineEmpty(line string, ignoreNL bool) bool {
	if ignoreNL {
		return line == ""
	}
	return line == "" || line == "\n" || line == "\r\n"
}

func (f *FileFunction) GetName() string { return "file" }

func (f *FileFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), nil),
		node.NewParameter(nil, "context", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *FileFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, nil),
		node.NewVariable(nil, "flags", 1, nil),
		node.NewVariable(nil, "context", 2, nil),
	}
}
