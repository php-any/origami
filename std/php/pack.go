package php

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

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

var packFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "format", 0, nil, nil),
	node.NewParameters(nil, "values", 1, nil, nil),
}

func (f *PackFunction) GetParams() []data.GetValue {
	return packFunctionGetParams
}

var packFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "values", 1, data.NewBaseType("mixed")),
}

func (f *PackFunction) GetVariables() []data.Variable {
	return packFunctionGetVariables
}

// UnpackFunction 实现 PHP 内置函数 unpack（支持的格式子集：C / C* / n / n*）
//
//	unpack(string $format, string $string, int $offset = 0): array|false
//
// 目前实现:
//   - "C"  / "C*" : 无符号字节
//   - "n"  / "n*" : 无符号 16 位大端整数
//
// 返回的 PHP 数组使用从 1 开始的整数键，
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
	offset := 0
	if v, _ := ctx.GetIndexValue(2); v != nil {
		if n, ok := v.(data.AsInt); ok {
			offset, _ = n.AsInt()
		}
	}
	if offset < 0 || offset > len(raw) {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("unpack(): Argument #3 ($offset) must be contained in argument #2 ($data)"), "ValueError")
	}
	result, err := unpackIntegers(format, raw[offset:])
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return result, nil
}

func unpackIntegers(format string, raw []byte) (*data.ArrayValue, error) {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, field := range strings.Split(format, "/") {
		if field == "" {
			continue
		}
		code := field[0]
		size := 0
		var order binary.ByteOrder = binary.NativeEndian
		switch code {
		case 'c', 'C':
			size = 1
		case 's', 'S', 'n', 'v':
			size = 2
		case 'i', 'I', 'l', 'L', 'N', 'V':
			size = 4
		case 'q', 'Q', 'J', 'P':
			size = 8
		default:
			return nil, fmt.Errorf("unsupported unpack format %c", code)
		}
		if code == 'n' || code == 'N' || code == 'J' {
			order = binary.BigEndian
		}
		if code == 'v' || code == 'V' || code == 'P' {
			order = binary.LittleEndian
		}
		end, count, star := 1, 1, false
		if end < len(field) && field[end] == '*' {
			star, end, count = true, end+1, len(raw)/size
		} else {
			for end < len(field) && field[end] >= '0' && field[end] <= '9' {
				end++
			}
			if end > 1 {
				var err error
				count, err = strconv.Atoi(field[1:end])
				if err != nil {
					return nil, err
				}
			}
		}
		if count > len(raw)/size {
			return nil, fmt.Errorf("not enough data")
		}
		name := field[end:]
		for i := 1; i <= count; i++ {
			var value uint64
			switch size {
			case 1:
				value = uint64(raw[0])
			case 2:
				value = uint64(order.Uint16(raw))
			case 4:
				value = uint64(order.Uint32(raw))
			case 8:
				value = order.Uint64(raw)
			}
			n := int(value)
			switch code {
			case 'c':
				n = int(int8(value))
			case 's':
				n = int(int16(value))
			case 'i', 'l':
				n = int(int32(value))
			}
			key := name
			if name == "" || count != 1 || star {
				key += strconv.Itoa(i)
			}
			out.SetStringKey(key, data.NewIntValue(n))
			raw = raw[size:]
		}
	}
	return out, nil
}

func (f *UnpackFunction) GetName() string { return "unpack" }

var unpackFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "format", 0, nil, nil),
	node.NewParameter(nil, "string", 1, nil, nil),
	node.NewParameter(nil, "offset", 2, node.NewIntLiteral(nil, "0"), nil),
}

func (f *UnpackFunction) GetParams() []data.GetValue {
	return unpackFunctionGetParams
}

var unpackFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "string", 1, data.NewBaseType("string")),
	node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
}

func (f *UnpackFunction) GetVariables() []data.Variable {
	return unpackFunctionGetVariables
}
