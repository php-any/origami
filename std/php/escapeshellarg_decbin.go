package php

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EscapeshellargFunction 实现 PHP 内置函数 escapeshellarg
//
//	escapeshellarg(string $arg): string
//
// 用单引号包裹参数，内部单引号转义为 '\”（POSIX 风格）。
func NewEscapeshellargFunction() data.FuncStmt {
	return &EscapeshellargFunction{}
}

type EscapeshellargFunction struct{}

func (f *EscapeshellargFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	s := ""
	if v != nil {
		s = v.AsString()
	}
	// PHP: "'" . str_replace("'", "'\\''", $arg) . "'"
	escaped := "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
	return data.NewStringValue(escaped), nil
}

func (f *EscapeshellargFunction) GetName() string { return "escapeshellarg" }
func (f *EscapeshellargFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "arg", 0, nil, nil)}
}
func (f *EscapeshellargFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "arg", 0, data.NewBaseType("string"))}
}

// DecbinFunction 实现 PHP 内置函数 decbin
//
//	decbin(int $num): string
func NewDecbinFunction() data.FuncStmt {
	return &DecbinFunction{}
}

type DecbinFunction struct{}

func (f *DecbinFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	n := int64(0)
	if v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if i, err := iv.AsInt(); err == nil {
				n = int64(i)
			}
		} else {
			s := strings.TrimSpace(v.AsString())
			if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
				n = parsed
			}
		}
	}
	return data.NewStringValue(strconv.FormatUint(uint64(n), 2)), nil
}

func (f *DecbinFunction) GetName() string { return "decbin" }
func (f *DecbinFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *DecbinFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, data.NewBaseType("int"))}
}
