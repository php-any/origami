package php

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	jsonser "github.com/php-any/origami/std/serializer/json"
)

// NewUnserializeFunction 创建 unserialize 函数
// 目前实现的子集：
// - 标量：N; / b:0; / b:1; / i:number; / s:len:"...";
// - 数组：a:len:{key;value;...}
//   - 数值下标 0..n-1 -> data.ArrayValue
//   - 其他键（字符串/非连续整数）-> data.ObjectValue（关联数组）
//
// 同时保留对旧版 "__origami_a:" / "__origami_o:" JSON 包装字符串的兼容。
func NewUnserializeFunction() data.FuncStmt {
	return &UnserializeFunction{}
}

type UnserializeFunction struct{}

func (f *UnserializeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 参数：string $data
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "data", 0, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}
	raw := v.(data.AsString).AsString()
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return data.NewBoolValue(false), nil
	}

	// 优先尝试解析 PHP 标准 serialize 语法（标量 + 数组）
	if strings.HasPrefix(raw, "N;") ||
		strings.HasPrefix(raw, "b:") ||
		strings.HasPrefix(raw, "i:") ||
		strings.HasPrefix(raw, "s:") ||
		strings.HasPrefix(raw, "a:") {

		if v, ok := parsePhpSerializedValue(raw); ok {
			return v, nil
		}
	}

	// 兼容旧实现中基于 __origami_* 前缀的 JSON 包装字符串
	if strings.HasPrefix(raw, "s:") {
		// 期望形如：s:<len>:"<content>";
		firstQuote := strings.IndexByte(raw, '"')
		if firstQuote < 0 {
			return data.NewBoolValue(false), nil
		}
		lastQuote := strings.LastIndexByte(raw, '"')
		if lastQuote <= firstQuote {
			return data.NewBoolValue(false), nil
		}
		content := raw[firstQuote+1 : lastQuote]

		// 尝试识别 Origami 内部为数组/对象生成的前缀，并用 JsonSerializer 还原。
		if strings.HasPrefix(content, "__origami_a:") {
			jsonPart := content[len("__origami_a:"):]
			serializer := jsonser.NewJsonSerializer()
			arr := &data.ArrayValue{}
			if err := serializer.UnmarshalArray([]byte(jsonPart), arr); err != nil {
				return data.NewBoolValue(false), nil
			}
			return arr, nil
		}
		if strings.HasPrefix(content, "__origami_o:") {
			jsonPart := content[len("__origami_o:"):]
			serializer := jsonser.NewJsonSerializer()
			obj := data.NewObjectValue()
			if err := serializer.UnmarshalObject([]byte(jsonPart), obj); err != nil {
				return data.NewBoolValue(false), nil
			}
			return obj, nil
		}

		// 普通字符串按原样返回
		return data.NewStringValue(content), nil
	}

	// 其它复杂类型（对象/引用等）暂不支持，直接返回 false
	return data.NewBoolValue(false), nil
}

// parsePhpSerializedValue 解析一整段 PHP serialize 字符串（支持标量 + 数组）。
// 这里只实现 Origami 目前需要的最小子集。
func parsePhpSerializedValue(s string) (data.Value, bool) {
	idx := 0
	v, ok := parsePhpValue(s, &idx)
	if !ok {
		return nil, false
	}
	// 必须完全消费完字符串
	if idx != len(s) {
		return nil, false
	}
	return v, true
}

func parsePhpValue(s string, idx *int) (data.Value, bool) {
	if *idx >= len(s) {
		return nil, false
	}
	switch s[*idx] {
	case 'N':
		// NULL: N;
		if *idx+2 <= len(s) && s[*idx:*idx+2] == "N;" {
			*idx += 2
			return data.NewNullValue(), true
		}
		return nil, false
	case 'b':
		// 布尔：b:0; / b:1;
		if *idx+4 <= len(s) && s[*idx:*idx+2] == "b:" && s[*idx+3] == ';' {
			switch s[*idx+2] {
			case '0':
				*idx += 4
				return data.NewBoolValue(false), true
			case '1':
				*idx += 4
				return data.NewBoolValue(true), true
			default:
				return nil, false
			}
		}
		return nil, false
	case 'i':
		// 整数：i:number;
		if *idx+3 > len(s) || s[*idx:*idx+2] != "i:" {
			return nil, false
		}
		j := *idx + 2
		sign := int64(1)
		if j < len(s) && (s[j] == '-' || s[j] == '+') {
			if s[j] == '-' {
				sign = -1
			}
			j++
		}
		start := j
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if start == j || j >= len(s) || s[j] != ';' {
			return nil, false
		}
		numStr := s[start:j]
		n, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return nil, false
		}
		*idx = j + 1
		return data.NewIntValue(int(n * sign)), true
	case 's':
		// 字符串：s:len:"...";  （这里不严格校验 len，只要能找到成对的引号就解析）
		if *idx+2 > len(s) || s[*idx:*idx+2] != "s:" {
			return nil, false
		}
		// 跳过 s:len:
		colon := strings.IndexByte(s[*idx+2:], ':')
		if colon < 0 {
			return nil, false
		}
		firstQuote := strings.IndexByte(s[*idx+2+colon+1:], '"')
		if firstQuote < 0 {
			return nil, false
		}
		firstQuote += *idx + 2 + colon + 1
		lastQuote := strings.LastIndexByte(s[firstQuote+1:], '"')
		if lastQuote < 0 {
			return nil, false
		}
		lastQuote += firstQuote + 1
		if lastQuote+2 > len(s) || s[lastQuote+1] != ';' {
			return nil, false
		}
		content := s[firstQuote+1 : lastQuote]
		*idx = lastQuote + 2
		return data.NewStringValue(content), true
	case 'a':
		// 数组：a:len:{key;value;...}
		return parsePhpArray(s, idx)
	default:
		return nil, false
	}
}

func parsePhpArray(s string, idx *int) (data.Value, bool) {
	if *idx+2 > len(s) || s[*idx:*idx+2] != "a:" {
		return nil, false
	}
	j := *idx + 2
	start := j
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	if start == j || j >= len(s) || s[j] != ':' {
		return nil, false
	}
	lenStr := s[start:j]
	n, err := strconv.Atoi(lenStr)
	if err != nil || n < 0 {
		return nil, false
	}
	j++ // 跳过 ':'
	if j >= len(s) || s[j] != '{' {
		return nil, false
	}
	j++
	*idx = j

	keys := make([]data.Value, 0, n)
	values := make([]data.Value, 0, n)
	for i := 0; i < n; i++ {
		key, ok := parsePhpValue(s, idx)
		if !ok {
			return nil, false
		}
		val, ok := parsePhpValue(s, idx)
		if !ok {
			return nil, false
		}
		keys = append(keys, key)
		values = append(values, val)
	}
	if *idx >= len(s) || s[*idx] != '}' {
		return nil, false
	}
	*idx++

	// 判断是否为 0..n-1 顺序的数值下标数组，若是则返回 ArrayValue，否则使用 ObjectValue 作为关联数组
	isSequential := true
	for i, k := range keys {
		ik, ok := k.(*data.IntValue)
		if !ok || ik.Value != i {
			isSequential = false
			break
		}
	}

	if isSequential {
		return data.NewArrayValue(values), true
	}

	obj := data.NewObjectValue()
	for i, k := range keys {
		var keyStr string
		switch kv := k.(type) {
		case *data.StringValue:
			keyStr = kv.Value
		case *data.IntValue:
			keyStr = strconv.Itoa(kv.Value)
		default:
			keyStr = k.AsString()
		}
		obj.SetProperty(keyStr, values[i])
	}
	return obj, true
}

func (f *UnserializeFunction) GetName() string {
	return "unserialize"
}

func (f *UnserializeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, data.String{}),
	}
}

func (f *UnserializeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.String{}),
	}
}
