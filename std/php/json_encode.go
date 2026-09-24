package php

import (
	"bytes"
	jsonpkg "encoding/json"
	"unsafe"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/serializer/json"
)

// json_encode 对齐 PHP：只读遍历，不改写入参。JsonSerializable 的返回值写到拷贝上。
func NewJsonEncodeFunction() data.FuncStmt {
	return &JsonEncodeFunction{}
}

type JsonEncodeFunction struct {
	data.Function
}

func (f *JsonEncodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	params := f.GetParams()
	if len(params) == 0 {
		clearJsonLastError()
		return data.NewStringValue("null"), nil
	}

	raw, _ := params[0].GetValue(ctx)
	if raw == nil {
		clearJsonLastError()
		return data.NewStringValue("null"), nil
	}
	v, ok := raw.(data.Value)
	if !ok {
		clearJsonLastError()
		return data.NewStringValue("null"), nil
	}

	flags := jsonEncodeFlags(ctx)
	encoded, ok, ctl := jsonEncode(ctx, v, flags)
	if ctl != nil {
		return nil, ctl
	}
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(encoded), nil
}

func jsonEncodeFlags(ctx data.Context) int {
	if ctx == nil {
		return 0
	}
	fv, ok := ctx.GetIndexValue(1)
	if !ok || fv == nil {
		return 0
	}
	iv, ok := fv.(data.AsInt)
	if !ok {
		return 0
	}
	n, err := iv.AsInt()
	if err != nil {
		return 0
	}
	return n
}

// JsonEncode 对齐 PHP json_encode：嵌套对象走 JsonSerializable / 公共属性，不把 ClassValue 编成 Object(FQCN) 字符串。
func JsonEncode(ctx data.Context, value data.Value) (string, bool, data.Control) {
	return jsonEncode(ctx, value, 0)
}

func jsonEncode(ctx data.Context, value data.Value, flags int) (string, bool, data.Control) {
	if value == nil {
		clearJsonLastError()
		return "null", true, nil
	}

	resolved, jerr, ctl := resolveJSONValue(ctx, value, map[uintptr]bool{}, 0)
	if ctl != nil {
		return "", false, ctl
	}
	if jerr != JSON_ERROR_NONE {
		setJsonLastError(jerr, "")
		return "", false, nil
	}

	serializer := json.NewJsonSerializer()
	var result []byte
	var err error

	switch v := resolved.(type) {
	case *data.IntValue:
		result, err = serializer.MarshalInt(v)
	case *data.StringValue:
		result, err = serializer.MarshalString(v)
	case *data.BoolValue:
		result, err = serializer.MarshalBool(v)
	case *data.FloatValue:
		result, err = serializer.MarshalFloat(v)
	case *data.NullValue:
		result, err = serializer.MarshalNull(v)
	case *data.ArrayValue:
		result, err = serializer.MarshalArray(v)
	case *data.ObjectValue:
		result, err = serializer.MarshalObject(v)
	case *data.ClassValue:
		result, err = serializer.MarshalClass(v)
	case data.ValueSerializer:
		result, err = v.Marshal(serializer)
	default:
		result, err = jsonpkg.Marshal(resolved.AsString())
	}

	if err != nil {
		setJsonLastError(JSON_ERROR_UNSUPPORTED_TYPE, "")
		return "", false, nil
	}

	result = applyJsonEncodeFlags(result, flags)
	clearJsonLastError()
	return string(result), true, nil
}

// applyJsonEncodeFlags 对齐 PHP JSON_HEX_*：只改 JSON 字符串内容，结构分隔用的 " 保持原样。
// Laravel Js::from 二次 json_encode 依赖 HEX_QUOT 把内容里的 " 变成 \u0022，否则
// JSON.parse('{"table":true}') 写进 HTML 双引号属性会截断，Livewire 报 mountAction 缺右括号。
func applyJsonEncodeFlags(in []byte, flags int) []byte {
	const (
		hexTag  = 1 // JSON_HEX_TAG
		hexAmp  = 2 // JSON_HEX_AMP
		hexApos = 4 // JSON_HEX_APOS
		hexQuot = 8 // JSON_HEX_QUOT
	)
	if flags&(hexTag|hexAmp|hexApos|hexQuot) == 0 {
		return in
	}
	var out bytes.Buffer
	out.Grow(len(in) + 24)
	inString := false
	for i := 0; i < len(in); i++ {
		c := in[i]
		if !inString {
			if c == '"' {
				inString = true
			}
			out.WriteByte(c)
			continue
		}
		if c == '\\' {
			if i+1 < len(in) {
				n := in[i+1]
				if n == '"' && flags&hexQuot != 0 {
					out.WriteString(`\u0022`)
					i++
					continue
				}
				out.WriteByte('\\')
				out.WriteByte(n)
				i++
				continue
			}
			out.WriteByte(c)
			continue
		}
		if c == '"' {
			inString = false
			out.WriteByte(c)
			continue
		}
		if flags&hexApos != 0 && c == '\'' {
			out.WriteString(`\u0027`)
			continue
		}
		if flags&hexAmp != 0 && c == '&' {
			out.WriteString(`\u0026`)
			continue
		}
		if flags&hexTag != 0 && c == '<' {
			out.WriteString(`\u003C`)
			continue
		}
		if flags&hexTag != 0 && c == '>' {
			out.WriteString(`\u003E`)
			continue
		}
		out.WriteByte(c)
	}
	return out.Bytes()
}

func (f *JsonEncodeFunction) GetName() string {
	return "json_encode"
}

var jsonEncodeFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
	node.NewParameter(nil, "flags", 1, data.NewIntValue(0), nil),
	node.NewParameter(nil, "depth", 2, data.NewIntValue(512), nil),
}

func (f *JsonEncodeFunction) GetParams() []data.GetValue {
	return jsonEncodeFunctionGetParams
}

var jsonEncodeFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, nil),
	node.NewVariable(nil, "flags", 1, nil),
	node.NewVariable(nil, "depth", 2, nil),
}

func (f *JsonEncodeFunction) GetVariables() []data.Variable {
	return jsonEncodeFunctionGetVariables
}

// resolveJSONValue 递归解析值：对实现了 JsonSerializable 的对象（含数组内嵌、顶层对象）调用其
// jsonSerialize() 方法，使序列化行为与 PHP 原生 json_encode 一致。
func resolveJSONValue(ctx data.Context, value data.Value, visiting map[uintptr]bool, depth int) (data.Value, int, data.Control) {
	if depth > 512 {
		return nil, JSON_ERROR_DEPTH, nil
	}
	if tv, ok := value.(*data.ThisValue); ok && tv != nil && tv.ClassValue != nil {
		value = tv.ClassValue
	}
	switch v := value.(type) {
	case *data.ArrayValue:
		ptr := uintptr(unsafe.Pointer(v))
		if visiting[ptr] {
			return nil, JSON_ERROR_RECURSION, nil
		}
		visiting[ptr] = true
		defer delete(visiting, ptr)
		// 必须拷贝后再写解析结果。就地替换会把 Livewire/视图共享数组改成 JSON 树，
		// 下次再 encode 字符串套字符串，直到 OOM / 十几秒卡死。
		cloned := data.CloneArrayValue(v)
		for i, z := range cloned.List {
			if z == nil {
				continue
			}
			nv, jerr, ctl := resolveJSONValue(ctx, z.Value, visiting, depth+1)
			if ctl != nil {
				return nil, JSON_ERROR_NONE, ctl
			}
			if jerr != JSON_ERROR_NONE {
				return nil, jerr, nil
			}
			cloned.List[i] = data.CopyZValKeepName(z, nv)
		}
		return cloned, JSON_ERROR_NONE, nil
	case *data.ObjectValue:
		ptr := uintptr(unsafe.Pointer(v))
		if visiting[ptr] {
			return nil, JSON_ERROR_RECURSION, nil
		}
		visiting[ptr] = true
		defer delete(visiting, ptr)
		cloned := data.NewObjectValue()
		var jerr int
		var ctl data.Control
		v.RangeProperties(func(k string, pv data.Value) bool {
			var nv data.Value
			nv, jerr, ctl = resolveJSONValue(ctx, pv, visiting, depth+1)
			if ctl != nil || jerr != JSON_ERROR_NONE {
				return false
			}
			cloned.SetProperty(k, nv)
			return true
		})
		if ctl != nil {
			return nil, JSON_ERROR_NONE, ctl
		}
		if jerr != JSON_ERROR_NONE {
			return nil, jerr, nil
		}
		return cloned, JSON_ERROR_NONE, nil
	case *data.ClassValue:
		ptr := uintptr(unsafe.Pointer(v))
		if visiting[ptr] {
			return nil, JSON_ERROR_RECURSION, nil
		}
		visiting[ptr] = true
		defer delete(visiting, ptr)
		jsonSerializable := data.Class{Name: "JsonSerializable"}
		if jsonSerializable.Is(v) {
			if method, has := v.GetMethod("jsonSerialize"); has {
				res, acl := method.Call(v.CreateContext(method.GetVariables()))
				if acl != nil {
					return nil, JSON_ERROR_NONE, acl
				}
				if res != nil {
					if rv, ok := res.(data.Value); ok {
						return resolveJSONValue(ctx, rv, visiting, depth+1)
					}
				}
			}
		}
	}
	return value, JSON_ERROR_NONE, nil
}
