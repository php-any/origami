package node

import "github.com/php-any/origami/data"

// YieldFromControl 实现 yield from 的控制流。
// 它本身是一个 YieldControl，用来在生成器中委托另一个可迭代对象/生成器。
type YieldFromControl struct {
	// 保存 yield from 时的上下文信息
	ctx    data.Context
	source data.GetValue

	// 解析后的迭代器/生成器
	iter data.Generator

	// 当前 key/value（实现 YieldValueControl）
	key   data.Value
	value data.Value

	// primed 表示已取出当前元素但尚未 next 过去。
	// 对 SplFileInfo 等可变迭代器，必须先 yield 再 next，否则 current 会被原地改掉。
	primed bool
}

// NewYieldFromControl 创建一个新的 YieldFromControl。
func NewYieldFromControl(ctx data.Context, source data.GetValue) data.YieldControl {
	return &YieldFromControl{
		ctx:    ctx,
		source: source,
	}
}

// GetValue 实现 data.Value 接口。
// 作为函数体中的一条语句被执行时，负责驱动内部迭代器向前推进一次：
// - 如果还有元素：通过返回自身（YieldControl）来让外层生成器产生一次 yield；
// - 如果已经结束：什么也不做，直接返回 nil,nil。
func (y *YieldFromControl) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	ctl := y.Next(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return nil, nil
}

// AsString 实现 data.Value 接口
func (y *YieldFromControl) AsString() string {
	return "yield from"
}

// GetYieldKey 实现 YieldValueControl 接口
func (y *YieldFromControl) GetYieldKey() data.Value {
	if y.key == nil {
		return data.NewNullValue()
	}
	return y.key
}

// GetYieldValue 实现 YieldValueControl 接口
func (y *YieldFromControl) GetYieldValue() data.Value {
	if y.value == nil {
		return data.NewNullValue()
	}
	return y.value
}

// CreateStackState 实现 YieldControl 接口，用于把 yield from 嵌入函数生成器执行流中。
func (y *YieldFromControl) CreateStackState(ctx data.Context, fn data.FuncStmt, originalBody []data.GetValue, bodyIndex int) data.Generator {
	// 解析 source 为 Generator
	if y.iter == nil {
		srcValue, ctl := y.source.GetValue(ctx)
		if ctl != nil {
			// 出现控制流时，由外层处理；此处仅简单忽略，生成器在 Next 时会重新触发
			return NewFuncYieldStackState(ctx, fn, originalBody, bodyIndex+1, nil, nil)
		}
		y.iter = resolveYieldFromSource(ctx, srcValue)
	}

	// 如果仍然无法获得有效的 iter，则退化为一个空生成器
	if y.iter == nil {
		return NewFuncYieldStackState(ctx, fn, originalBody, bodyIndex+1, nil, nil)
	}

	// 构建新的函数体，将 bodyIndex 位置的 yield from 语句替换为 y 自身。
	newBody := make([]data.GetValue, 0, len(originalBody))
	newBody = append(newBody, originalBody[:bodyIndex]...)
	newBody = append(newBody, y)
	newBody = append(newBody, originalBody[bodyIndex+1:]...)

	// 预先推进一次，获得第一个元素（如果有），这样外部生成器在第一次 valid()/current() 时就能看到值。
	ok, ctl := y.advance(ctx)
	if ctl != nil {
		// 出错时，直接交给外层生成器在执行过程中处理，这里先不设置当前值。
		return NewFuncYieldStackState(ctx, fn, newBody, bodyIndex+1, nil, nil)
	}
	if !ok {
		// 没有任何元素，直接从下一条语句开始。
		return NewFuncYieldStackState(ctx, fn, originalBody, bodyIndex+1, nil, nil)
	}

	return NewFuncYieldStackState(ctx, fn, newBody, bodyIndex, y.key, y.value)
}

// resolveYieldFromSource 将 yield from 右侧解析为可委托的 Generator。
// 生成器函数/方法返回的是 *ClassValue{Class:*GeneratorClass}，需解包内部 data.Generator。
// 也支持实现 Iterator / IteratorAggregate 的对象（如 RecursiveIteratorIterator）。
func resolveYieldFromSource(ctx data.Context, srcValue data.GetValue) data.Generator {
	if srcValue == nil {
		return nil
	}
	switch v := srcValue.(type) {
	case data.Generator:
		return v
	case *data.ClassValue:
		if gc, ok := v.Class.(*GeneratorClass); ok && gc.generator != nil {
			return gc.generator
		}
		return newObjectIteratorGenerator(ctx, v)
	case *data.ThisValue:
		if v.ClassValue != nil {
			if gc, ok := v.ClassValue.Class.(*GeneratorClass); ok && gc.generator != nil {
				return gc.generator
			}
			return newObjectIteratorGenerator(ctx, v.ClassValue)
		}
	case *data.ArrayValue:
		return newArrayGenerator(ctx, v)
	}
	return nil
}

// objectIteratorGenerator 把实现 Iterator（或 IteratorAggregate）的 ClassValue 适配为 Generator。
type objectIteratorGenerator struct {
	obj     *data.ClassValue
	started bool
}

func newObjectIteratorGenerator(ctx data.Context, obj *data.ClassValue) data.Generator {
	if obj == nil || obj.Class == nil {
		return nil
	}

	// IteratorAggregate → getIterator()
	isAggregate, ctl := checkClassIs(ctx, obj.Class, "IteratorAggregate")
	if ctl == nil && isAggregate {
		inner, ictl := callValueMethod(obj, "getIterator")
		if ictl != nil {
			return nil
		}
		switch it := inner.(type) {
		case *data.ClassValue:
			return newObjectIteratorGenerator(ctx, it)
		case *data.ThisValue:
			if it.ClassValue != nil {
				return newObjectIteratorGenerator(ctx, it.ClassValue)
			}
		case data.Generator:
			return it
		}
	}

	isIterator, ctl := checkClassIs(ctx, obj.Class, "Iterator")
	if ctl != nil || !isIterator {
		return nil
	}

	g := &objectIteratorGenerator{obj: obj}
	// PHP Iterator：委托前先 rewind
	_ = callVoidMethod(obj, "rewind")
	g.started = true
	return g
}

func (g *objectIteratorGenerator) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return g.Current(ctx)
}
func (g *objectIteratorGenerator) AsString() string { return "ObjectIteratorGenerator" }

func (g *objectIteratorGenerator) Current(ctx data.Context) (data.Value, data.Control) {
	return callValueMethod(g.obj, "current")
}
func (g *objectIteratorGenerator) Key(ctx data.Context) (data.Value, data.Control) {
	return callValueMethod(g.obj, "key")
}
func (g *objectIteratorGenerator) Next(ctx data.Context) data.Control {
	return callVoidMethod(g.obj, "next")
}
func (g *objectIteratorGenerator) Rewind(ctx data.Context) (data.Value, data.Control) {
	if ctl := callVoidMethod(g.obj, "rewind"); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}
func (g *objectIteratorGenerator) Valid(ctx data.Context) (data.Value, data.Control) {
	ok, ctl := callBoolMethod(g.obj, "valid")
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(ok), nil
}
func (g *objectIteratorGenerator) Send(ctx data.Context, value data.Value) data.Control {
	return nil
}
func (g *objectIteratorGenerator) Throw(ctx data.Context) data.Control {
	return nil
}
func (g *objectIteratorGenerator) GetReturn(ctx data.Context) (data.Value, data.Control) {
	return data.NewNullValue(), nil
}

// advance 读取内部迭代器的当前元素（不提前 next）。
// 若此前已取出过元素（primed），则先 next 再读取。
func (y *YieldFromControl) advance(ctx data.Context) (ok bool, ctl data.Control) {
	if y.iter == nil {
		return false, nil
	}

	if y.primed {
		ctl = y.iter.Next(ctx)
		if ctl != nil {
			return false, ctl
		}
	}

	valid, ctl := y.iter.Valid(ctx)
	if ctl != nil {
		return false, ctl
	}
	if b, okBool := valid.(*data.BoolValue); !okBool || !b.Value {
		y.key = nil
		y.value = nil
		y.primed = false
		return false, nil
	}

	current, ctl := y.iter.Current(ctx)
	if ctl != nil {
		return false, ctl
	}
	key, ctl := y.iter.Key(ctx)
	if ctl != nil {
		return false, ctl
	}

	y.key = key
	y.value = current
	y.primed = true
	return true, nil
}

// Next 将内部迭代器推进到下一个元素，并在成功时返回自身作为控制流（YieldControl）。
func (y *YieldFromControl) Next(ctx data.Context) data.Control {
	ok, ctl := y.advance(ctx)
	if ctl != nil {
		return ctl
	}
	if !ok {
		// 没有更多元素了
		return nil
	}
	// 有新的元素，返回自身，交给外层生成器设置当前 key/value。
	return y
}

// newArrayGenerator 将数组包装成一个简单的 Generator。
// 这里实现一个最小可用的数组生成器，供 yield from 使用。
type arrayGenerator struct {
	ctx   data.Context
	array *data.ArrayValue
	index int
}

func newArrayGenerator(ctx data.Context, array *data.ArrayValue) data.Generator {
	return &arrayGenerator{
		ctx:   ctx,
		array: array,
		index: 0,
	}
}

func (a *arrayGenerator) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if a.index >= len(a.array.List) {
		return data.NewNullValue(), nil
	}
	return a.array.List[a.index].Value, nil
}

func (a *arrayGenerator) AsString() string {
	return "ArrayGenerator"
}

func (a *arrayGenerator) Current(ctx data.Context) (data.Value, data.Control) {
	if a.index >= len(a.array.List) {
		return data.NewNullValue(), nil
	}
	return a.array.List[a.index].Value, nil
}

func (a *arrayGenerator) Key(ctx data.Context) (data.Value, data.Control) {
	if a.index >= len(a.array.List) {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(a.index), nil
}

func (a *arrayGenerator) Next(ctx data.Context) data.Control {
	a.index++
	return nil
}

func (a *arrayGenerator) Rewind(ctx data.Context) (data.Value, data.Control) {
	a.index = 0
	return data.NewNullValue(), nil
}

func (a *arrayGenerator) Valid(ctx data.Context) (data.Value, data.Control) {
	return data.NewBoolValue(a.index < len(a.array.List)), nil
}

func (a *arrayGenerator) Send(ctx data.Context, value data.Value) data.Control {
	// 对数组生成器，不支持 send，直接忽略
	return nil
}

func (a *arrayGenerator) Throw(ctx data.Context) data.Control {
	// 对数组生成器，不支持 throw，直接结束
	a.index = len(a.array.List)
	return nil
}

func (a *arrayGenerator) GetReturn(ctx data.Context) (data.Value, data.Control) {
	// 数组生成器没有返回值
	return data.NewNullValue(), nil
}
