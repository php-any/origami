package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/serializer/json"
)

func NewJsonDecodeFunction() data.FuncStmt {
	return &JsonDecodeFunction{}
}

type JsonDecodeFunction struct {
	data.Function
}

func (f *JsonDecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	params := f.GetParams()
	if len(params) == 0 {
		return data.NewNullValue(), nil
	}

	// 获取第一个参数（JSON 字符串）
	jsonParam := params[0]
	jsonValue, _ := jsonParam.GetValue(ctx)
	if jsonValue == nil {
		return data.NewNullValue(), nil
	}

	// 获取 JSON 字符串值
	var jsonString string
	switch v := jsonValue.(type) {
	case *data.StringValue:
		jsonString = v.Value
	case *data.IntValue:
		jsonString = string(rune(v.Value))
	default:
		jsonString = v.(data.Value).AsString()
	}

	// 检查是否有第二个参数（目标类名）
	if len(params) > 1 && params[1] != nil {
		classParam := params[1]
		classValue, _ := classParam.GetValue(ctx)
		if classValue != nil {
			// 获取类名
			var className string
			switch v := classValue.(type) {
			case *data.StringValue:
				className = v.Value
			case *data.ClassValue:
				className = v.Class.GetName()
			case *data.NullValue:
				// 当第二个参数是 null 时，反序列化为 ObjectValue
				serializer := json.NewJsonSerializer()
				value := data.NewObjectValue()
				err := value.Unmarshal([]byte(jsonString), serializer)
				if err != nil {
					return data.NewNullValue(), nil
				}
				return value, nil
			default:
				className = v.(data.Value).AsString()
			}

			// 获取 VM 来查找类定义
			vm := ctx.GetVM()
			if vm != nil {
				// 根据类名获取类定义
				if classStmt, ok := vm.GetClass(className); ok {
					// 创建类实例
					classInstance, _ := classStmt.GetValue(ctx)
					if classValue, ok := classInstance.(*data.ClassValue); ok {
						// 创建 JSON 反序列化器
						serializer := json.NewJsonSerializer()
						// 反序列化到类实例
						err := serializer.UnmarshalClass([]byte(jsonString), classValue)
						if err != nil {
							return data.NewNullValue(), nil
						}
						return classValue, nil
					}
				}
			}
		}
	}

	// 创建 JSON 反序列化器
	serializer := json.NewJsonSerializer()

	// 默认情况下，反序列化为通用值
	value := data.NewObjectValue()
	err := value.Unmarshal([]byte(jsonString), serializer)
	if err != nil {
		return data.NewNullValue(), nil
	}

	return value, nil
}

func (f *JsonDecodeFunction) GetName() string {
	return "json_decode"
}

func (f *JsonDecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "json", 0, nil, data.String{}),
		node.NewParameter(nil, "class", 1, data.NewNullValue(), nil),
	}
}

func (f *JsonDecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "json", 0, nil),
		node.NewVariable(nil, "class", 1, nil),
	}
}
