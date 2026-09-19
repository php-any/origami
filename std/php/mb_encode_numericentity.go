package php

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MbEncodeNumericentityFunction 实现 mb_encode_numericentity。
// Symfony HtmlDumper 每行都调用；PHP polyfill 按字节解释执行会把异常页拖过 30s，
// 无效 UTF-8 上 ulenMask 还可能 i+=0 死循环。必须走原生实现。
type MbEncodeNumericentityFunction struct{}

func NewMbEncodeNumericentityFunction() data.FuncStmt {
	return &MbEncodeNumericentityFunction{}
}

func (f *MbEncodeNumericentityFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	mapVal, _ := ctx.GetIndexValue(1)
	if strVal == nil {
		return data.NewStringValue(""), nil
	}
	s := strVal.AsString()
	if s == "" {
		return data.NewStringValue(""), nil
	}

	convmap := convmapFromValue(mapVal)
	if len(convmap) < 4 {
		return data.NewBoolValue(false), nil
	}

	hex := false
	if hexVal, ok := ctx.GetIndexValue(3); ok && hexVal != nil {
		if as, ok := hexVal.(data.AsBool); ok {
			if b, err := as.AsBool(); err == nil {
				hex = b
			}
		}
	}

	return data.NewStringValue(mbEncodeNumericentity(s, convmap, hex)), nil
}

func convmapFromValue(v data.Value) []int {
	if v == nil {
		return nil
	}
	arr, ok := v.(*data.ArrayValue)
	if !ok || arr == nil {
		return nil
	}
	out := make([]int, 0, len(arr.List))
	for _, z := range arr.List {
		if z == nil || z.Value == nil {
			out = append(out, 0)
			continue
		}
		if as, ok := z.Value.(data.AsInt); ok {
			n, err := as.AsInt()
			if err == nil {
				out = append(out, n)
				continue
			}
		}
		out = append(out, 0)
	}
	n := (len(out) / 4) * 4
	return out[:n]
}

func mbEncodeNumericentity(s string, convmap []int, hex bool) string {
	var b strings.Builder
	b.Grow(len(s) + len(s)/4)
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		var c int
		var raw string
		if r == utf8.RuneError && size == 1 {
			c = int(s[i])
			raw = s[i : i+1]
			i++
		} else {
			c = int(r)
			raw = s[i : i+size]
			i += size
		}
		if enc, ok := numericEntityFor(c, convmap, hex); ok {
			b.WriteString(enc)
		} else {
			b.WriteString(raw)
		}
	}
	return b.String()
}

func numericEntityFor(c int, convmap []int, hex bool) (string, bool) {
	for j := 0; j+3 < len(convmap); j += 4 {
		if c >= convmap[j] && c <= convmap[j+1] {
			off := (c + convmap[j+2]) & convmap[j+3]
			if hex {
				return "&#x" + strings.ToUpper(strconv.FormatInt(int64(uint32(off)), 16)) + ";", true
			}
			return "&#" + strconv.Itoa(off) + ";", true
		}
	}
	return "", false
}

func (f *MbEncodeNumericentityFunction) GetName() string { return "mb_encode_numericentity" }

func (f *MbEncodeNumericentityFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "map", 1, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "encoding", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "hex", 3, data.NewBoolValue(false), nil),
	}
}

func (f *MbEncodeNumericentityFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "map", 1, data.NewBaseType("array")),
		node.NewVariable(nil, "encoding", 2, data.NewNullableType(data.NewBaseType("string"))),
		node.NewVariable(nil, "hex", 3, data.NewBaseType("bool")),
	}
}
