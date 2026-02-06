package php

import (
	"bytes"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PackFunction 实现 PHP 内置函数 pack（支持的格式子集：C / C*）
//
//	pack(string $format, mixed ...$values): string
//
// 目前仅实现:
//   - "C"  / "C*" : 无符号字节，取参数的 AsInt 值并截断到 0-255
type PackFunction struct{}

func NewPackFunction() data.FuncStmt { return &PackFunction{} }

func (f *PackFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatVal, _ := ctx.GetIndexValue(0)
	if formatVal == nil {
		return data.NewBoolValue(false), nil
	}
	format := formatVal.AsString()

	var buf bytes.Buffer

	switch format {
	case "C", "C*":
		// 后续所有参数按字节写入
		for i := 1; ; i++ {
			v, ok := ctx.GetIndexValue(i)
			if !ok {
				break
			}
			b := byte(0)
			if asInt, ok := v.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil {
					b = byte(n)
				}
			}
			buf.WriteByte(b)
		}
	default:
		// 未实现的格式，返回 false
		return data.NewBoolValue(false), nil
	}

	return data.NewStringValue(buf.String()), nil
}

func (f *PackFunction) GetName() string { return "pack" }

func (f *PackFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameters(nil, "values", 1, nil, nil),
	}
}

func (f *PackFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "values", 1, data.NewBaseType("mixed")),
	}
}

// UnpackFunction 实现 PHP 内置函数 unpack（支持的格式子集：C / C*）
//
//	unpack(string $format, string $string, int $offset = 0): array|false
//
// 目前仅实现:
//   - "C"  : 返回数组 [第一个字节]
//   - "C*" : 返回数组 [所有字节] （数组索引从 0 开始，近似 PHP 行为）
type UnpackFunction struct{}

func NewUnpackFunction() data.FuncStmt { return &UnpackFunction{} }

func (f *UnpackFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatVal, _ := ctx.GetIndexValue(0)
	dataVal, _ := ctx.GetIndexValue(1)
	// offsetVal, _ := ctx.GetIndexValue(2) // 当前暂不处理 offset

	if formatVal == nil || dataVal == nil {
		return data.NewBoolValue(false), nil
	}
	format := formatVal.AsString()
	s := dataVal.AsString()

	switch format {
	case "C":
		if len(s) == 0 {
			return data.NewBoolValue(false), nil
		}
		return data.NewArrayValue([]data.Value{
			data.NewIntValue(int(s[0])),
		}), nil
	case "C*":
		if len(s) == 0 {
			return data.NewArrayValue([]data.Value{}), nil
		}
		values := make([]data.Value, len(s))
		for i := 0; i < len(s); i++ {
			values[i] = data.NewIntValue(int(s[i]))
		}
		return data.NewArrayValue(values), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *UnpackFunction) GetName() string { return "unpack" }

func (f *UnpackFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameter(nil, "string", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *UnpackFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "string", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
	}
}
