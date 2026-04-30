package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArraySpliceFunction() data.FuncStmt {
	return &ArraySpliceFunction{}
}

type ArraySpliceFunction struct{}

func (f *ArraySpliceFunction) GetName() string {
	return "array_splice"
}

func (f *ArraySpliceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil),
		node.NewParameter(nil, "offset", 1, nil, nil),
		node.NewParameter(nil, "length", 2, nil, nil),
		node.NewParameter(nil, "replacement", 3, data.NewArrayValue(nil), nil),
	}
}

func (f *ArraySpliceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, nil),
		node.NewVariable(nil, "offset", 1, nil),
		node.NewVariable(nil, "length", 2, nil),
		node.NewVariable(nil, "replacement", 3, nil),
	}
}

func (f *ArraySpliceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	offsetValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil || offsetValue == nil {
		return data.NewArrayValue(nil), nil
	}

	// Get offset
	offset := 0
	if asInt, ok := offsetValue.(data.AsInt); ok {
		offset, _ = asInt.AsInt()
	}

	// Extract items from the array
	var items []data.Value
	var keys []string
	isAssoc := false

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		items = v.ToValueList()
	case *data.ObjectValue:
		isAssoc = true
		v.RangeProperties(func(key string, val data.Value) bool {
			keys = append(keys, key)
			items = append(items, val)
			return true
		})
	default:
		return data.NewArrayValue(nil), nil
	}

	arrLen := len(items)

	// Normalize offset
	if offset < 0 {
		offset = arrLen + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset > arrLen {
		offset = arrLen
	}

	// Get length
	length := arrLen - offset // default: remove all from offset
	lengthValue, exists := ctx.GetIndexValue(2)
	if exists && lengthValue != nil {
		if _, isNull := lengthValue.(*data.NullValue); !isNull {
			if asInt, ok := lengthValue.(data.AsInt); ok {
				length, _ = asInt.AsInt()
				if length < 0 {
					length = arrLen - offset + length
					if length < 0 {
						length = 0
					}
				}
			}
		}
	}

	if offset+length > arrLen {
		length = arrLen - offset
	}

	// Get replacement
	var replacement []data.Value
	replacementValue, exists := ctx.GetIndexValue(3)
	if exists && replacementValue != nil {
		switch rv := replacementValue.(type) {
		case *data.ArrayValue:
			replacement = rv.ToValueList()
		case *data.ObjectValue:
			rv.RangeProperties(func(key string, val data.Value) bool {
				replacement = append(replacement, val)
				return true
			})
		default:
			replacement = []data.Value{replacementValue}
		}
	}

	// Extract removed elements (this is what we return)
	removed := make([]data.Value, length)
	copy(removed, items[offset:offset+length])

	// Build new array: items[:offset] + replacement + items[offset+length:]
	newItems := make([]data.Value, 0, len(items)-length+len(replacement))
	newItems = append(newItems, items[:offset]...)
	newItems = append(newItems, replacement...)
	newItems = append(newItems, items[offset+length:]...)

	// Modify the original array in place
	// We need to update the ZVal that was passed by reference
	zval := ctx.GetIndexZVal(0)
	if zval != nil {
		if isAssoc {
			newObj := data.NewObjectValue()
			for i, val := range newItems {
				if i < len(keys) && i < offset {
					newObj.SetProperty(keys[i], val)
				} else {
					newObj.SetProperty(fmt.Sprintf("%d", i), val)
				}
			}
			zval.Value = newObj
		} else {
			zval.Value = data.NewArrayValue(newItems)
		}
	}

	return data.NewArrayValue(removed), nil
}
