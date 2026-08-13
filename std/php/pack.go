package php

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PackFunction 实现 PHP 内置函数 pack（支持的格式子集：C / C* / n / n*）
//
//	pack(string $format, mixed ...$values): string
//
// 目前实现:
//   - "C"  / "C*" : 无符号字节，取参数的 AsInt 值并截断到 0-255
//   - "n"  / "n*" : 无符号 16 位大端整数
type PackFunction struct{}

func NewPackFunction() data.FuncStmt { return &PackFunction{} }

func packArgs(ctx data.Context) []data.Value {
	// NewParameters 把可变实参打成 index=1 的 ArrayValue（同 sprintf）
	valuesValue, _ := ctx.GetIndexValue(1)
	if valuesValue == nil {
		return nil
	}
	if arr, ok := valuesValue.(*data.ArrayValue); ok {
		return arr.ToValueList()
	}
	return []data.Value{valuesValue}
}

func (f *PackFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatVal, _ := ctx.GetIndexValue(0)
	if formatVal == nil {
		return data.NewBoolValue(false), nil
	}
	format := formatVal.AsString()
	args := packArgs(ctx)

	var buf bytes.Buffer

	switch format {
	case "C", "C*":
		for _, v := range args {
			b := byte(0)
			if asInt, ok := v.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil {
					b = byte(n)
				}
			}
			buf.WriteByte(b)
		}
	case "n", "n*":
		for _, v := range args {
			n := uint16(0)
			if asInt, ok := v.(data.AsInt); ok {
				if iv, err := asInt.AsInt(); err == nil {
					n = uint16(iv)
				}
			}
			var b [2]byte
			binary.BigEndian.PutUint16(b[:], n)
			buf.Write(b[:])
		}
	default:
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

// UnpackFunction 实现 PHP 内置函数 unpack（支持的格式子集：C / C* / n / n*）
//
//	unpack(string $format, string $string, int $offset = 0): array|false
//
// 目前实现:
//   - "C"  / "C*" : 无符号字节
//   - "n"  / "n*" : 无符号 16 位大端整数
//
// 返回的数组索引从 1 开始（使用 ObjectValue，键为 "1".."n"），
// 以匹配 PHP 对 unpack("C*") 的行为，便于在 PHP 代码中使用 $a[1], $a[2] 访问。
type UnpackFunction struct{}

func NewUnpackFunction() data.FuncStmt { return &UnpackFunction{} }

func (f *UnpackFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatVal, _ := ctx.GetIndexValue(0)
	dataVal, _ := ctx.GetIndexValue(1)

	if formatVal == nil || dataVal == nil {
		return data.NewBoolValue(false), nil
	}
	format := formatVal.AsString()
	raw := []byte(dataVal.AsString())

	switch format {
	case "C":
		if len(raw) == 0 {
			return data.NewBoolValue(false), nil
		}
		obj := data.NewObjectValue()
		obj.SetProperty("1", data.NewIntValue(int(raw[0])))
		return obj, nil
	case "C*":
		obj := data.NewObjectValue()
		for i := 0; i < len(raw); i++ {
			key := fmt.Sprintf("%d", i+1)
			obj.SetProperty(key, data.NewIntValue(int(raw[i])))
		}
		return obj, nil
	case "n":
		if len(raw) < 2 {
			return data.NewBoolValue(false), nil
		}
		obj := data.NewObjectValue()
		obj.SetProperty("1", data.NewIntValue(int(binary.BigEndian.Uint16(raw[:2]))))
		return obj, nil
	case "n*":
		obj := data.NewObjectValue()
		idx := 1
		for i := 0; i+1 < len(raw); i += 2 {
			key := fmt.Sprintf("%d", idx)
			obj.SetProperty(key, data.NewIntValue(int(binary.BigEndian.Uint16(raw[i:i+2]))))
			idx++
		}
		return obj, nil
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
