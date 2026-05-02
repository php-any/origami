package php

import (
	"encoding/json"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	origamiJson "github.com/php-any/origami/std/serializer/json"
)

func NewJsonDecodeFunction() data.FuncStmt {
	return &JsonDecodeFunction{}
}

type JsonDecodeFunction struct {
	data.Function
}

func (f *JsonDecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	jsonValue, _ := ctx.GetIndexValue(0)
	classValue, _ := ctx.GetIndexValue(1)

	var jsonString string
	if jsonValue != nil {
		switch v := jsonValue.(type) {
		case *data.StringValue:
			jsonString = v.Value
		case *data.IntValue:
			if i, err := v.AsInt(); err == nil {
				jsonString = string(rune(i))
			}
		default:
			jsonString = jsonValue.AsString()
		}
	}

	// 第二个参数：true=关联数组, false/null=对象, string=类名
	asArray := false
	if classValue != nil {
		if bv, ok := classValue.(*data.BoolValue); ok {
			asArray = bv.Value
			classValue = nil
		} else if _, ok := classValue.(*data.NullValue); ok {
			classValue = nil
		}
	}

	if classValue != nil {
		var className string
		switch v := classValue.(type) {
		case *data.StringValue:
			className = v.Value
		case *data.ClassValue:
			className = v.Class.GetName()
		default:
			className = v.(data.Value).AsString()
		}

		vm := ctx.GetVM()
		if vm != nil {
			if classStmt, acl := vm.GetOrLoadClass(className); acl == nil {
				classInstance, _ := classStmt.GetValue(ctx)
				if cv, ok := classInstance.(*data.ClassValue); ok {
					serializer := origamiJson.NewJsonSerializer()
					if err := serializer.UnmarshalClass([]byte(jsonString), cv); err != nil {
						return data.NewNullValue(), nil
					}
					return cv, nil
				}
			} else {
				return nil, acl
			}
		}
	}

	if asArray {
		v, err := goJsonDecode(jsonString)
		if err != nil {
			return data.NewNullValue(), nil
		}
		return v, nil
	}

	serializer := origamiJson.NewJsonSerializer()
	value := data.NewObjectValue()
	if err := value.Unmarshal([]byte(jsonString), serializer); err != nil {
		return data.NewNullValue(), nil
	}
	return value, nil
}

func goJsonDecode(js string) (data.Value, error) {
	var result interface{}
	if err := json.Unmarshal([]byte(js), &result); err != nil {
		return nil, err
	}
	return convertGoValue(result), nil
}

func convertGoValue(v interface{}) data.Value {
	switch val := v.(type) {
	case nil:
		return data.NewNullValue()
	case bool:
		return data.NewBoolValue(val)
	case float64:
		if val == float64(int64(val)) {
			return data.NewIntValue(int(val))
		}
		return data.NewFloatValue(val)
	case string:
		return data.NewStringValue(val)
	case []interface{}:
		list := make([]data.Value, len(val))
		for i, item := range val {
			list[i] = convertGoValue(item)
		}
		return data.NewArrayValue(list)
	case map[string]interface{}:
		arrList := make([]*data.ZVal, 0, len(val))
		for k, item := range val {
			arrList = append(arrList, &data.ZVal{Name: k, Value: convertGoValue(item)})
		}
		return &data.ArrayValue{List: arrList}
	default:
		s := strconv.FormatFloat(v.(float64), 'f', -1, 64)
		return data.NewStringValue(s)
	}
}

func (f *JsonDecodeFunction) GetName() string {
	return "json_decode"
}

func (f *JsonDecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "json", 0, nil, data.String{}),
		node.NewParameter(nil, "assoc", 1, data.NewNullValue(), nil),
	}
}

func (f *JsonDecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "json", 0, nil),
		node.NewVariable(nil, "assoc", 1, nil),
	}
}
