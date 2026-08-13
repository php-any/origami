package data

// IsScalarAssignFast reports values assignable without reference binding or structural clone.
func IsScalarAssignFast(v Value) bool {
	switch v.(type) {
	case *IntValue, *FloatValue, *BoolValue, *StringValue, *NullValue:
		return true
	default:
		return false
	}
}

// AssignIntToZVal 将整数按值写入变量槽（始终使用独立 *IntValue，避免与数组快照等共享指针）。
func AssignIntToZVal(zv *ZVal, n int) {
	if zv == nil {
		return
	}
	if iv, ok := zv.Value.(*IntValue); ok {
		if iv.Value == n {
			return
		}
	}
	zv.Value = &IntValue{Value: n}
}

// AssignFloatToZVal 将浮点按值写入变量槽。
func AssignFloatToZVal(zv *ZVal, f float64) {
	if zv == nil {
		return
	}
	if fv, ok := zv.Value.(*FloatValue); ok {
		if fv.Value == f {
			return
		}
	}
	zv.Value = &FloatValue{Value: f}
}

// AssignScalarToZVal 按 PHP 标量语义赋值（整数/浮点按值复制）。
func AssignScalarToZVal(zv *ZVal, value Value) {
	if zv == nil {
		return
	}
	switch v := value.(type) {
	case *IntValue:
		AssignIntToZVal(zv, v.Value)
	case *FloatValue:
		AssignFloatToZVal(zv, v.Value)
	case *BoolValue:
		zv.Value = NewBoolValue(v.Value)
	case *NullValue:
		zv.Value = NewNullValue()
	default:
		zv.Value = value
	}
}
