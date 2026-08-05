package node

import (
	"fmt"
	"strconv"

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

// arraySlotKey 返回 PHP 数组槽位对应的键名（空 Name 的密集整数键用下标）。
func arraySlotKey(z *data.ZVal, index int) string {
	if z != nil && z.Name != "" {
		return z.Name
	}
	return data.IntArrayKeyName(index)
}

// objectToNamedArray 将对象/类属性转为带键名的 ArrayValue（用于 array + 语义）。
func objectToNamedArray(obj hasRangeProperties) *data.ArrayValue {
	list := make([]*data.ZVal, 0)
	obj.RangeProperties(func(key string, value data.Value) bool {
		list = append(list, data.NewNamedZVal(key, value))
		return true
	})
	return &data.ArrayValue{List: list}
}

// valueAsArrayForUnion 将 Array/Object/Class 统一为可按键并集的 ArrayValue。
func valueAsArrayForUnion(v data.Value) (*data.ArrayValue, bool) {
	switch x := v.(type) {
	case *data.ArrayValue:
		return x, true
	case *data.ObjectValue:
		return objectToNamedArray(x), true
	case *data.ClassValue:
		return objectToNamedArray(x), true
	default:
		return nil, false
	}
}

// mergeArrayUnion 实现 PHP 的 array + array：左侧键优先，仅追加右侧不存在的键。
func mergeArrayUnion(left, right *data.ArrayValue) *data.ArrayValue {
	result := make([]*data.ZVal, 0, len(left.List)+len(right.List))
	seen := make(map[string]struct{}, len(left.List)+len(right.List))

	appendSlot := func(key string, value data.Value) {
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		if _, isInt := data.ParseIntArrayKeyName(key); isInt {
			if n, err := strconv.Atoi(key); err == nil && n == len(result) {
				result = append(result, data.NewZVal(value))
				return
			}
		}
		result = append(result, data.NewNamedZVal(key, value))
	}

	for i, z := range left.List {
		if z == nil {
			continue
		}
		appendSlot(arraySlotKey(z, i), z.Value)
	}
	for i, z := range right.List {
		if z == nil {
			continue
		}
		appendSlot(arraySlotKey(z, i), z.Value)
	}
	return &data.ArrayValue{List: result}
}

func addOperandIsFloat(v data.GetValue) bool {
	val, ok := v.(data.Value)
	if !ok {
		return false
	}
	_, ok = val.(*data.FloatValue)
	return ok
}

func addOperandAsFloat64(v data.GetValue) (float64, bool) {
	val, ok := v.(data.Value)
	if !ok {
		return 0, false
	}
	if f, ok := val.(data.AsFloat); ok {
		fl, err := f.AsFloat()
		return fl, err == nil
	}
	return 0, false
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

	// PHP：任一侧为 float 时，+ 结果为 float
	if addOperandIsFloat(lv) || addOperandIsFloat(rv) {
		if lf, ok := addOperandAsFloat64(lv); ok {
			if rf, ok := addOperandAsFloat64(rv); ok {
				return data.NewFloatValue(lf + rf), nil
			}
		}
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
		case data.AsFloat:
			li, err := l.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rf, err := r.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewFloatValue(float64(li) + rf), nil
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
			return data.NewStringValue(data.NewFloatValue(lf).AsString() + rf), nil
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

	case *data.BoolValue:
		// 布尔值与任何类型相加：转换为整数（true->1, false->0）后相加
		var li int
		if l.Value {
			li = 1
		} else {
			li = 0
		}
		switch r := rv.(type) {
		case data.AsInt:
			ri, err := r.AsInt()
			if err != nil {
				// 如果无法转换为整数，尝试字符串拼接
				if str, ok := rv.(data.AsString); ok {
					return data.NewStringValue(fmt.Sprintf("%d", li) + str.AsString()), nil
				}
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewIntValue(li + ri), nil
		case data.AsFloat:
			rf, err := r.AsFloat()
			if err != nil {
				// 如果无法转换为浮点数，尝试字符串拼接
				if str, ok := rv.(data.AsString); ok {
					return data.NewStringValue(fmt.Sprintf("%d", li) + str.AsString()), nil
				}
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewFloatValue(float64(li) + rf), nil
		case data.AsString:
			// 布尔值与字符串相加：转换为字符串后拼接
			return data.NewStringValue(fmt.Sprintf("%d", li) + r.AsString()), nil
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
		// PHP 数组 + ：按键并集；关联字面量 ['a'=>…] 在 Origami 里可能是 ObjectValue。
		if ra, ok := valueAsArrayForUnion(rv.(data.Value)); ok {
			return mergeArrayUnion(l, ra), nil
		}
		// 右边非数组/对象：作为下一个元素追加
		result := l.ToValueList()
		result = append(result, rv.(data.Value))
		return data.NewArrayValue(result), nil
	case *data.ObjectValue:
		if ra, ok := valueAsArrayForUnion(rv.(data.Value)); ok {
			return mergeArrayUnion(objectToNamedArray(l), ra), nil
		}
		return nil, data.NewErrorThrow(b.from, fmt.Errorf("对象不能与非对象/数组类型相加: %T", rv))
	case *data.ClassValue:
		if ra, ok := valueAsArrayForUnion(rv.(data.Value)); ok {
			return mergeArrayUnion(objectToNamedArray(l), ra), nil
		}
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
	case *data.AnyValue:
		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil
	}

	return nil, data.NewErrorThrow(b.from, fmt.Errorf("TODO 有未支持的类型加法 %v", lv))
}
