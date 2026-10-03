package data

// ArrayValue 的 copy-on-write；PropertyBag 不参与 PHP 数组值复制。
// rc 为指向该容器的 zval 数。字面量/新建为 0，写入变量时 CowAddRef。
// 写入结构前必须 CowSeparateZVal，避免共享容器被原地改掉。

func CowAddRef(v Value) Value {
	switch a := v.(type) {
	case *ArrayValue:
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
	}
}

// CowAssign 把值写入 zval：释放旧数组引用并 AddRef 新数组（标量/对象按指针赋值）。
func CowAssign(zv *ZVal, v Value) {
	if zv == nil {
		return
	}
	zv.Defined = true
	if zv.ReadValue() == v {
		return
	}
	CowRelease(zv.ReadValue())
	zv.StoreRaw(CowAddRef(v))
}

func CowSeparateZVal(zv *ZVal) {
	if zv == nil {
		return
	}
	switch a := zv.ReadValue().(type) {
	case *ArrayValue:
		if a.rc > 1 {
			a.rc--
			cloned := CloneArrayValue(a)
			cloned.rc = 1
			zv.StoreRaw(cloned)
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
	return zv.ReadValue()
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
