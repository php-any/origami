package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

type BinaryAdd struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryAdd(from data.From, left, right data.GetValue) *BinaryAdd {
	return &BinaryAdd{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

// hasRangeProperties 检查值是否有 RangeProperties 方法（ObjectValue 或 ClassValue）
type hasRangeProperties interface {
	RangeProperties(func(key string, value data.Value) bool)
}

// mergeObjects 合并两个对象/类实例，如果键相同，保留左边对象的值
func mergeObjects(left, right hasRangeProperties) *data.ObjectValue {
	result := data.NewObjectValue()

	// 先添加左边对象的所有属性
	left.RangeProperties(func(key string, value data.Value) bool {
		result.SetProperty(key, value)
		return true
	})

	// 获取左边对象的所有键，用于检查
	leftKeys := make(map[string]bool)
	left.RangeProperties(func(key string, value data.Value) bool {
		leftKeys[key] = true
		return true
	})

	// 然后添加右边对象的属性（只添加左边不存在的键）
	right.RangeProperties(func(key string, value data.Value) bool {
		if !leftKeys[key] {
			result.SetProperty(key, value)
		}
		return true
	})

	return result
}

// objectToArrayValues 将对象/类实例的属性值转换为数组元素
func objectToArrayValues(obj hasRangeProperties) []data.Value {
	var result []data.Value
	obj.RangeProperties(func(key string, value data.Value) bool {
		result = append(result, value)
		return true
	})
	return result
}

// mergeArrayWithObject 将数组与对象/类实例合并，返回新数组
func mergeArrayWithObject(arr *data.ArrayValue, obj hasRangeProperties) *data.ArrayValue {
	result := arr.ToValueList()
	objValues := objectToArrayValues(obj)
	result = append(result, objValues...)
	return data.NewArrayValue(result).(*data.ArrayValue)
}

// mergeObjectWithArray 将对象/类实例与数组合并，返回新数组
func mergeObjectWithArray(obj hasRangeProperties, arr *data.ArrayValue) *data.ArrayValue {
	result := objectToArrayValues(obj)
	result = append(result, arr.ToValueList()...)
	return data.NewArrayValue(result).(*data.ArrayValue)
}

func (b *BinaryAdd) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	switch l := lv.(type) {
	case *data.StringValue:
		lStr := l.AsString()
		rStr, rCtl := ValueToDisplayString(ctx, rv)
		if rCtl != nil {
			return nil, rCtl
		}
		return data.NewStringValue(lStr + rStr), nil
	case *data.IntValue:
		switch r := rv.(type) {
		case data.AsInt:
			li, err := l.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			ri, err := r.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}

			return data.NewIntValue(li + ri), nil
		case data.AsString:
			return data.NewStringValue(l.AsString() + r.AsString()), nil
		}
	case *data.FloatValue:
		switch r := rv.(type) {
		case data.AsInt:
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			ri, err := r.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}

			return data.NewFloatValue(lf + float64(ri)), nil
		case *data.StringValue:
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rf := r.AsString()
			return data.NewStringValue(fmt.Sprintf("%f", lf) + rf), nil
		case data.AsFloat:
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rf, err := r.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewFloatValue(lf + rf), nil
		}

	case *data.NullValue:
		if riv, ok := rv.(data.AsInt); ok {
			ri, err := riv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewIntValue(0 + ri), nil
		}

		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil

	case *data.ArrayValue:
		// 数组相加是合并操作
		switch r := rv.(type) {
		case *data.ArrayValue:
			// 合并两个数组：先复制左边数组，然后追加右边数组的元素
			result := l.ToValueList()
			result = append(result, r.ToValueList()...)
			return data.NewArrayValue(result), nil
		case *data.ObjectValue:
			// 数组与对象相加：将对象的属性值添加到数组中
			return mergeArrayWithObject(l, r), nil
		case *data.ClassValue:
			// 数组与类实例相加：将类的属性值添加到数组中
			return mergeArrayWithObject(l, r), nil
		default:
			// 如果右边不是数组、对象或类，将右边作为单个元素添加到数组中
			result := l.ToValueList()
			result = append(result, r.(data.Value))
			return data.NewArrayValue(result), nil
		}
	case *data.ObjectValue:
		// 对象相加是合并操作
		switch r := rv.(type) {
		case *data.ObjectValue:
			// 合并两个对象：如果键相同，保留左边的值
			return mergeObjects(l, r), nil
		case *data.ClassValue:
			// 对象与类实例相加：合并属性（如果键相同，保留左边对象的值）
			return mergeObjects(l, r), nil
		case *data.ArrayValue:
			// 对象与数组相加：先添加对象的属性值，然后添加数组元素
			return mergeObjectWithArray(l, r), nil
		default:
			// 如果右边不是对象、类或数组，返回错误
			return nil, data.NewErrorThrow(b.from, fmt.Errorf("对象不能与非对象/数组类型相加: %T", r))
		}
	case *data.ClassValue:
		// 类实例相加是合并操作（与对象类似）；否则尝试 __toString 后字符串拼接
		switch r := rv.(type) {
		case *data.ObjectValue:
			// 类实例与对象相加：合并属性（如果键相同，保留左边类实例的值）
			return mergeObjects(l, r), nil
		case *data.ClassValue:
			// 合并两个类实例：合并属性（如果键相同，保留左边的值）
			return mergeObjects(l, r), nil
		case *data.ArrayValue:
			// 类实例与数组相加：先添加类实例的属性值，然后添加数组元素
			return mergeObjectWithArray(l, r), nil
		default:
			// 右边非对象/类/数组：若类有 __toString，则转为字符串后拼接
			lStr, lCtl := ValueToDisplayString(ctx, l)
			if lCtl != nil {
				return nil, lCtl
			}
			rStr, rCtl := ValueToDisplayString(ctx, rv)
			if rCtl != nil {
				return nil, rCtl
			}
			return data.NewStringValue(lStr + rStr), nil
		}
	case *data.AnyValue:
		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil
	}

	return nil, data.NewErrorThrow(b.from, fmt.Errorf("TODO 有未支持的类型加法 %v", lv))
}
