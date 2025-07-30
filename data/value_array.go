package data

import (
	"fmt"
)

func NewArrayValue(v []Value) Value {
	return &ArrayValue{
		Value: v,
	}
}

type ArrayValue struct {
	Value []Value
}

func (a *ArrayValue) GetValue(ctx Context) (GetValue, Control) {
	return a, nil
}

func (a *ArrayValue) AsString() string {
	str := "["
	for _, value := range a.Value {
		str = str + value.AsString() + ", "
	}
	if len(str) > 2 {
		str = str[:len(str)-2]
	}

	str = str + "]"
	return fmt.Sprintf("%s", str)
}

func (a *ArrayValue) AsBool() (bool, error) {
	return len(a.Value) > 0, nil
}

func (a *ArrayValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "push":
		return &ArrayValuePush{&a.Value}, true
	case "pop":
		return &ArrayValuePop{&a.Value}, true
	case "shift":
		return &ArrayValueShift{&a.Value}, true
	case "unshift":
		return &ArrayValueUnshift{&a.Value}, true
	case "slice":
		return &ArrayValueSlice{a.Value}, true
	case "splice":
		return &ArrayValueSplice{&a.Value}, true
	case "join":
		return &ArrayValueJoin{a.Value}, true
	case "reverse":
		return &ArrayValueReverse{a.Value}, true
	case "sort":
		return &ArrayValueSort{&a.Value}, true
	case "indexOf":
		return &ArrayValueIndexOf{a.Value}, true
	case "includes":
		return &ArrayValueIncludes{a.Value}, true
	case "forEach":
		return &ArrayValueForEach{a.Value}, true
	case "map":
		return &ArrayValueMap{a.Value}, true
	case "filter":
		return &ArrayValueFilter{a.Value}, true
	case "reduce":
		return &ArrayValueReduce{a.Value}, true
	case "concat":
		return &ArrayValueConcat{a.Value}, true
	case "every":
		return &ArrayValueEvery{a.Value}, true
	case "some":
		return &ArrayValueSome{a.Value}, true
	case "find":
		return &ArrayValueFind{a.Value}, true
	case "findIndex":
		return &ArrayValueFindIndex{a.Value}, true
	case "flat":
		return &ArrayValueFlat{a.Value}, true
	case "flatMap":
		return &ArrayValueFlatMap{a.Value}, true
	}

	return nil, false
}

func (a *ArrayValue) GetProperty(name string) (Value, bool) {
	switch name {
	case "length":
		return NewIntValue(len(a.Value)), true
	}
	return nil, false
}
