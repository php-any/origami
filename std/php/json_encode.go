package php

import (
	jsonpkg "encoding/json"
	"unsafe"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/serializer/json"
)

func NewJsonEncodeFunction() data.FuncStmt {
	return &JsonEncodeFunction{}
}

type JsonEncodeFunction struct {
	data.Function
}

func (f *JsonEncodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	params := f.GetParams()
	if len(params) == 0 {
		clearJsonLastError()
		return data.NewStringValue("null"), nil
	}

	// 获取第一个参数（要编码的值）
	valueParam := params[0]
	raw, _ := valueParam.GetValue(ctx)
	if raw == nil {
		clearJsonLastError()
		return data.NewStringValue("null"), nil
	}

	// 递归解析值：对实现了 JsonSerializable 的对象（含数组/对象内嵌的）优先调用 jsonSerialize()
	value := raw
	if v, ok := raw.(data.Value); ok {
		resolved, jerr, ctl := resolveJSONValue(ctx, v, map[uintptr]bool{}, 0)
		if ctl != nil {
			return nil, ctl
		}
		if jerr != JSON_ERROR_NONE {
			setJsonLastError(jerr, "")
			return data.NewBoolValue(false), nil
		}
		value = resolved
	}

	// 创建 JSON 序列化器
	serializer := json.NewJsonSerializer()

	// 根据值的类型进行序列化
	var result []byte
	var err error

	switch v := value.(type) {
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
		// 对于其他类型，转换为字符串
		result, err = jsonpkg.Marshal(value.(data.Value).AsString())
	}

	if err != nil {
		setJsonLastError(JSON_ERROR_UNSUPPORTED_TYPE, "")
		return data.NewBoolValue(false), nil
	}

	clearJsonLastError()
	return data.NewStringValue(string(result)), nil
}

func (f *JsonEncodeFunction) GetName() string {
	return "json_encode"
}

func (f *JsonEncodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), nil),
		node.NewParameter(nil, "depth", 2, data.NewIntValue(512), nil),
	}
}

func (f *JsonEncodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
		node.NewVariable(nil, "flags", 1, nil),
		node.NewVariable(nil, "depth", 2, nil),
	}
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
		for i, z := range v.List {
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
			v.List[i] = data.NewNamedZVal(z.Name, nv)
		}
	case *data.ObjectValue:
		ptr := uintptr(unsafe.Pointer(v))
		if visiting[ptr] {
			return nil, JSON_ERROR_RECURSION, nil
		}
		visiting[ptr] = true
		defer delete(visiting, ptr)
		for k, pv := range v.GetProperties() {
			nv, jerr, ctl := resolveJSONValue(ctx, pv, visiting, depth+1)
			if ctl != nil {
				return nil, JSON_ERROR_NONE, ctl
			}
			if jerr != JSON_ERROR_NONE {
				return nil, jerr, nil
			}
			v.SetProperty(k, nv)
		}
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
