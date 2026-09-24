package php

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
						setJsonLastError(JSON_ERROR_SYNTAX, "")
						return data.NewNullValue(), nil
					}
					clearJsonLastError()
					return cv, nil
				}
			} else {
				return nil, acl
			}
		}
	}

	// PHP: assoc=true → 关联数组；assoc=false 时数组/标量仍为 PHP array/标量，仅 {} 为对象
	trimmed := strings.TrimSpace(jsonString)
	useObject := !asArray && len(trimmed) > 0 && trimmed[0] == '{'

	if !useObject {
		v, err := goJsonDecode(jsonString)
		if err != nil {
			setJsonLastError(JSON_ERROR_SYNTAX, "")
			return data.NewNullValue(), nil
		}
		clearJsonLastError()
		return v, nil
	}

	serializer := origamiJson.NewJsonSerializer()
	value := data.NewObjectValue()
	if err := value.Unmarshal([]byte(jsonString), serializer); err != nil {
		setJsonLastError(JSON_ERROR_SYNTAX, "")
		return data.NewNullValue(), nil
	}
	clearJsonLastError()
	return value, nil
}

func goJsonDecode(js string) (data.Value, error) {
	dec := json.NewDecoder(strings.NewReader(js))
	dec.UseNumber()
	return decodeJSONToken(dec)
}

// decodeJSONToken 按 JSON 出现顺序解析，保留对象键序（Livewire checksum 依赖 json_encode 键序稳定）。
func decodeJSONToken(dec *json.Decoder) (data.Value, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			list := make([]*data.ZVal, 0)
			for dec.More() {
				keyTok, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyTok.(string)
				if !ok {
					return nil, fmt.Errorf("json: expected object key string")
				}
				val, err := decodeJSONToken(dec)
				if err != nil {
					return nil, err
				}
				list = append(list, data.NewNamedZVal(key, val))
			}
			if _, err := dec.Token(); err != nil { // consume '}'
				return nil, err
			}
			return &data.ArrayValue{List: list}, nil
		case '[':
			items := make([]data.Value, 0)
			for dec.More() {
				val, err := decodeJSONToken(dec)
				if err != nil {
					return nil, err
				}
				items = append(items, val)
			}
			if _, err := dec.Token(); err != nil { // consume ']'
				return nil, err
			}
			return data.NewArrayValue(items), nil
		default:
			return nil, fmt.Errorf("json: unexpected delimiter %v", t)
		}
	case nil:
		return data.NewNullValue(), nil
	case bool:
		return data.NewBoolValue(t), nil
	case string:
		return data.NewStringValue(t), nil
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return data.NewIntValue(int(i)), nil
		}
		f, err := t.Float64()
		if err != nil {
			return data.NewStringValue(t.String()), nil
		}
		return data.NewFloatValue(f), nil
	case float64:
		if t == float64(int64(t)) {
			return data.NewIntValue(int(t)), nil
		}
		return data.NewFloatValue(t), nil
	default:
		return data.NewStringValue(fmt.Sprint(t)), nil
	}
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

var jsonDecodeFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "json", 0, nil, data.String{}),
	// PHP 8 正式参数名为 $associative（Livewire 等用 named arg associative: true）
	node.NewParameter(nil, "associative", 1, data.NewNullValue(), nil),
}

func (f *JsonDecodeFunction) GetParams() []data.GetValue {
	return jsonDecodeFunctionGetParams
}

var jsonDecodeFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "json", 0, nil),
	node.NewVariable(nil, "associative", 1, nil),
}

func (f *JsonDecodeFunction) GetVariables() []data.Variable {
	return jsonDecodeFunctionGetVariables
}
