package php

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DechexFunction 实现 PHP 内置函数 dechex
//
//	dechex(int $num): string
//
// 将整数转为小写十六进制字符串（无 0x 前缀）。
func NewDechexFunction() data.FuncStmt {
	return &DechexFunction{}
}

type DechexFunction struct{}

func (f *DechexFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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
	// PHP dechex 对负数按无符号 32/64 位语义输出；此处按 Go uint64 位模式对齐常见用法
	return data.NewStringValue(fmt.Sprintf("%x", uint64(n))), nil
}

func (f *DechexFunction) GetName() string { return "dechex" }
func (f *DechexFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *DechexFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, data.NewBaseType("int"))}
}

// HexdecFunction 实现 PHP 内置函数 hexdec
//
//	hexdec(string $hex_string): int|float
//
// 将十六进制字符串转为十进制整数。
func NewHexdecFunction() data.FuncStmt {
	return &HexdecFunction{}
}

type HexdecFunction struct{}

func (f *HexdecFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	s := strings.TrimSpace(v.AsString())
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	// PHP hexdec 会忽略任意位置的非十六进制字符（hexdec("a-b") === 171）
	cleaned := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			cleaned = append(cleaned, c)
		}
	}
	if len(cleaned) == 0 {
		return data.NewIntValue(0), nil
	}
	n, err := strconv.ParseUint(string(cleaned), 16, 64)
	if err != nil {
		return data.NewIntValue(0), nil
	}
	// 超出 int 范围时 PHP 返回 float；Origami IntValue 用平台 int，过大时截断为 int 位宽可接受于 uuid 场景
	return data.NewIntValue(int(n)), nil
}

func (f *HexdecFunction) GetName() string { return "hexdec" }
func (f *HexdecFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "hex_string", 0, nil, nil)}
}
func (f *HexdecFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "hex_string", 0, data.NewBaseType("string"))}
}

// OctdecFunction 实现 PHP 内置函数 octdec
//
//	octdec(string $octal_string): int|float
func NewOctdecFunction() data.FuncStmt {
	return &OctdecFunction{}
}

type OctdecFunction struct{}

func (f *OctdecFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	s := strings.TrimSpace(v.AsString())
	cleaned := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '7' {
			cleaned = append(cleaned, c)
		}
	}
	if len(cleaned) == 0 {
		return data.NewIntValue(0), nil
	}
	n, err := strconv.ParseUint(string(cleaned), 8, 64)
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(n)), nil
}

func (f *OctdecFunction) GetName() string { return "octdec" }
func (f *OctdecFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "octal_string", 0, nil, nil)}
}
func (f *OctdecFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "octal_string", 0, data.NewBaseType("string"))}
}

// DecoctFunction 实现 PHP 内置函数 decoct
//
//	decoct(int $num): string
func NewDecoctFunction() data.FuncStmt {
	return &DecoctFunction{}
}

type DecoctFunction struct{}

func (f *DecoctFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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
	return data.NewStringValue(fmt.Sprintf("%o", uint64(n))), nil
}

func (f *DecoctFunction) GetName() string { return "decoct" }
func (f *DecoctFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *DecoctFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, data.NewBaseType("int"))}
}
