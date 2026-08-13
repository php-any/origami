package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayMergeFunction() data.FuncStmt {
	return &ArrayMergeFunction{}
}

type ArrayMergeFunction struct{}

// isStringArrayKey 判断是否为 PHP array_merge 中的“字符串键”。
// 空 Name、以及纯数字字符串 Name（如 "0"）都按整数键处理，会被重新编号。
func isStringArrayKey(name string) bool {
	if name == "" {
		return false
	}
	_, isInt := data.ParseIntArrayKeyName(name)
	return !isInt
}

func (f *ArrayMergeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsValue, _ := ctx.GetIndexValue(0)
	if paramsValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsArray, ok := paramsValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// PHP 语义：
	// - string 键：后面的覆盖前面的；
	// - int 键（含 "0"/"1" 这类数字字符串键）：按出现顺序从 0 重新编号。
	result := make([]*data.ZVal, 0)
	stringKeyIndex := map[string]int{} // string key -> slot in result
	nextInt := 0

	appendInt := func(v data.Value) {
		result = append(result, data.NewNamedZVal(data.IntArrayKeyName(nextInt), v))
		nextInt++
	}

	setString := func(key string, v data.Value) {
		if idx, ok := stringKeyIndex[key]; ok {
			result[idx].Value = v
			return
		}
		stringKeyIndex[key] = len(result)
		result = append(result, data.NewNamedZVal(key, v))
	}

	for _, paramValue := range paramsArray.ToValueList() {
		switch v := paramValue.(type) {
		case *data.ArrayValue:
			for _, zval := range v.List {
				if zval == nil {
					continue
				}
				if isStringArrayKey(zval.Name) {
					setString(zval.Name, zval.Value)
				} else {
					appendInt(zval.Value)
				}
			}

		case *data.ObjectValue:
			for key, val := range v.GetProperties() {
				if isStringArrayKey(key) {
					setString(key, val)
				} else {
					appendInt(val)
				}
			}

		default:
			// 非数组参数：PHP 会 warning；这里按值附加为下一个 int 键
			appendInt(v)
		}
	}

	return &data.ArrayValue{List: result}, nil
}

func (f *ArrayMergeFunction) GetName() string {
	return "array_merge"
}

func (f *ArrayMergeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayMergeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
