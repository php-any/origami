package php

import (
	jsonpkg "encoding/json"

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
		return data.NewStringValue("null"), nil
	}

	// 获取第一个参数（要编码的值）
	valueParam := params[0]
	value, _ := valueParam.GetValue(ctx)
	if value == nil {
		return data.NewStringValue("null"), nil
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
		return data.NewStringValue("null"), nil
	}

	return data.NewStringValue(string(result)), nil
}

func (f *JsonEncodeFunction) GetName() string {
	return "json_encode"
}

func (f *JsonEncodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *JsonEncodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}
