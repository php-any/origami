package data

// ArrayValue / ObjectValue（关联数组）的 copy-on-write：
// rc 为指向该容器的 zval 数。字面量/新建为 0，写入变量时 CowAddRef。
// 写入结构前必须 CowSeparateZVal，避免共享容器被原地改掉。

func CowAddRef(v Value) Value {
	switch a := v.(type) {
	case *ArrayValue:
		a.rc++
	case *ObjectValue:
		a.rc++
	}
	return v
}

func CowRelease(v Value) {
	switch a := v.(type) {
	case *ArrayValue:
		if a.rc > 0 {
			a.rc--
		}
	case *ObjectValue:
		if a.rc > 0 {
			a.rc--
		}
	}
}

// CowAssign 把值写入 zval：释放旧数组引用并 AddRef 新数组（标量/对象按指针赋值）。
func CowAssign(zv *ZVal, v Value) {
	if zv == nil {
		return
	}
	zv.Defined = true
	if zv.Value == v {
		return
	}
	CowRelease(zv.Value)
	zv.Value = CowAddRef(v)
}

func CowSeparateZVal(zv *ZVal) {
	if zv == nil {
		return
	}
	switch a := zv.Value.(type) {
	case *ArrayValue:
		if a.rc > 1 {
			a.rc--
			cloned := CloneArrayValue(a)
			cloned.rc = 1
			zv.Value = cloned
		} else if a.rc == 0 {
			a.rc = 1
		}
	case *ObjectValue:
		if a.rc > 1 {
			a.rc--
			cloned := CloneObjectValue(a)
			cloned.rc = 1
			zv.Value = cloned
		} else if a.rc == 0 {
			a.rc = 1
		}
	}
}

// CowSeparateIndex 对调用帧第 index 个参数做写前分离（array_push 等引用参数）。
func CowSeparateIndex(ctx Context, index int) Value {
	zv := ctx.GetIndexZVal(index)
	if zv == nil {
		v, _ := ctx.GetIndexValue(index)
		return v
	}
	CowSeparateZVal(zv)
	return zv.Value
}

func SeparateArrayValue(a *ArrayValue) *ArrayValue {
	if a == nil {
		return nil
	}
	if a.rc > 1 {
		a.rc--
		cloned := CloneArrayValue(a)
		cloned.rc = 1
		return cloned
	}
	if a.rc == 0 {
		a.rc = 1
	}
	return a
}

func SeparateObjectValue(o *ObjectValue) *ObjectValue {
	if o == nil {
		return nil
	}
	if o.rc > 1 {
		o.rc--
		cloned := CloneObjectValue(o)
		cloned.rc = 1
		return cloned
	}
	if o.rc == 0 {
		o.rc = 1
	}
	return o
}
