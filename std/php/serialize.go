package php

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	jsonser "github.com/php-any/origami/std/serializer/json"
)

// NewSerializeFunction 创建 serialize 函数
// 当前为最小实现，对标上面的 unserialize 子集，仅支持：
// - string  -> s:len:"...";   （len 为字节长度）
// - int     -> i:number;
// - bool    -> b:0; / b:1;
// - null    -> N;
// 其它类型暂不支持，返回 false。
func NewSerializeFunction() data.FuncStmt {
	return &SerializeFunction{}
}

type SerializeFunction struct {
	data.Function
}

// makeSerializedString 构造形如 s:len:"content"; 的序列化字符串
func makeSerializedString(content string) string {
	l := len([]byte(content))
	var sb strings.Builder
	sb.Grow(len(content) + 16)
	sb.WriteString("s:")
	sb.WriteString(fmt.Sprintf("%d", l))
	sb.WriteString(":\"")
	sb.WriteString(content)
	sb.WriteString("\";")
	return sb.String()
}

// phpSerializeValue 将 Origami 内部的 data.Value 按 PHP serialize 语义编码为字符串。
// 目前支持：
// - null / bool / int / string
// - 数组：按 PHP 的 a:len:{key;value;...} 语法编码（数值下标使用 i:n;，关联键使用 s:len:"key";）
// - 对象：ClassValue 按 PHP 的 O:...:...:{...} 格式序列化公共属性
// 其它复杂类型返回 false，交由上层处理。
func phpSerializeValue(v data.Value, serializer *jsonser.JsonSerializer) (string, bool) {
	switch val := v.(type) {
	case *data.NullValue:
		return "N;", true
	case *data.BoolValue:
		if val.Value {
			return "b:1;", true
		}
		return "b:0;", true
	case *data.IntValue:
		return fmt.Sprintf("i:%d;", val.Value), true
	case *data.StringValue:
		return makeSerializedString(val.Value), true
	case *data.ArrayValue:
		// 数值下标数组：a:<len>:{i:0;v0;i:1;v1;...}
		values := val.ToValueList()
		var sb strings.Builder
		sb.WriteString("a:")
		sb.WriteString(strconv.Itoa(len(values)))
		sb.WriteString(":{")
		for idx, elem := range values {
			// key
			sb.WriteString("i:")
			sb.WriteString(strconv.Itoa(idx))
			sb.WriteString(";")
			// value
			valStr, ok := phpSerializeValue(elem, serializer)
			if !ok {
				return "", false
			}
			sb.WriteString(valStr)
		}
		sb.WriteString("}")
		return sb.String(), true
	case *data.ObjectValue:
		// 关联数组语义：使用字符串键序列化为 PHP 数组
		type kv struct {
			key string
			val data.Value
		}
		props := make([]kv, 0)
		val.RangeProperties(func(k string, v data.Value) bool {
			if v == nil {
				return true
			}
			props = append(props, kv{key: k, val: v})
			return true
		})

		var sb strings.Builder
		sb.WriteString("a:")
		sb.WriteString(strconv.Itoa(len(props)))
		sb.WriteString(":{")
		for _, p := range props {
			// key 始终作为字符串键处理
			sb.WriteString(makeSerializedString(p.key))
			valStr, ok := phpSerializeValue(p.val, serializer)
			if !ok {
				return "", false
			}
			sb.WriteString(valStr)
		}
		sb.WriteString("}")
		return sb.String(), true
	case *data.ClassValue:
		// 对应 PHP 中的对象序列化：O:<len>:"ClassName":<propCount>:{...}
		className := val.Class.GetName()
		classNameLen := len([]byte(className))

		// 按插入顺序遍历当前实例上已存在的属性（不强求覆盖 PHP 对“未显式赋值默认属性”的行为）
		type kv struct {
			key string
			val data.Value
		}
		props := make([]kv, 0)
		val.RangeProperties(func(k string, v data.Value) bool {
			if v == nil {
				return true
			}
			props = append(props, kv{key: k, val: v})
			return true
		})

		var sb strings.Builder
		// O:<classNameLen>:"ClassName":<propCount>:{
		sb.WriteString("O:")
		sb.WriteString(strconv.Itoa(classNameLen))
		sb.WriteString(":\"")
		sb.WriteString(className)
		sb.WriteString("\":")
		sb.WriteString(strconv.Itoa(len(props)))
		sb.WriteString(":{")

		for _, p := range props {
			// 属性名总是按公共属性处理：s:len:"name";
			nameLen := len([]byte(p.key))
			sb.WriteString("s:")
			sb.WriteString(strconv.Itoa(nameLen))
			sb.WriteString(":\"")
			sb.WriteString(p.key)
			sb.WriteString("\";")

			valStr, ok := phpSerializeValue(p.val, serializer)
			if !ok {
				return "", false
			}
			sb.WriteString(valStr)
		}

		sb.WriteString("}")
		return sb.String(), true
	default:
		return "", false
	}
}

func (f *SerializeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	params := f.GetParams()
	if len(params) == 0 {
		return data.NewBoolValue(false), nil
	}
	raw, _ := params[0].GetValue(ctx)
	if raw == nil {
		return data.NewStringValue("N;"), nil
	}

	v, ok := raw.(data.Value)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	serializer := jsonser.NewJsonSerializer()

	s, ok := phpSerializeValue(v, serializer)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(s), nil
}

func (f *SerializeFunction) GetName() string {
	return "serialize"
}

func (f *SerializeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *SerializeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}
