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
		if srcValue == nil {
			// 空源，直接构造一个已关闭的生成器
			return NewFuncYieldStackState(ctx, fn, originalBody, bodyIndex+1, nil, nil)
		}

		switch v := srcValue.(type) {
		case data.Generator:
			y.iter = v
		default:
			// 如果不是生成器，尝试将其视为数组或可迭代对象，并转换为内置生成器
			if arr, ok := srcValue.(*data.ArrayValue); ok {
				y.iter = newArrayGenerator(ctx, arr)
			}
		}
	}

	// 如果仍然无法获得有效的 iter，则退化为一个空生成器
	if y.iter == nil {
		return NewFuncYieldStackState(ctx, fn, originalBody, bodyIndex+1, nil, nil)
	}

	// 构建新的函数体，将 bodyIndex 位置的 yield from 语句替换为 y 自身。
	newBody := originalBody[:bodyIndex]
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

// advance 将内部迭代器向前推进一次，并更新当前的 key/value。
// 返回值 ok 表示是否成功获取到一个新的元素。
func (y *YieldFromControl) advance(ctx data.Context) (ok bool, ctl data.Control) {
	if y.iter == nil {
		return false, nil
	}

	valid, ctl := y.iter.Valid(ctx)
	if ctl != nil {
		return false, ctl
	}
	if b, okBool := valid.(*data.BoolValue); !okBool || !b.Value {
		// 已经结束
		y.key = nil
		y.value = nil
		return false, nil
	}

	// 读取当前元素
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

	// 将内部迭代器推进到下一位，为下一次 advance/Next 做准备
	ctl = y.iter.Next(ctx)
	if ctl != nil {
		return false, ctl
	}

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
