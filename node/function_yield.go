package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

func NewFuncYieldStackState(ctx data.Context, fn data.FuncStmt, body []data.GetValue, index int, currentKey, currentValue data.Value) data.Generator {
	generator := &FuncYieldStackState{
		ctx:          ctx,
		Fn:           fn,
		BodyIndex:    index,
		Body:         body,
		CurrentKey:   currentKey,
		CurrentValue: currentValue,
		autoKeyIndex: 0,
		initialized:  currentKey != nil || currentValue != nil, // 如果有初始值，表示已执行过第一次 yield
	}
	// 如果 currentKey 为 nil（第一次 yield 没有显式 key），分配 autoKey=0
	if currentKey == nil && currentValue != nil {
		generator.CurrentKey = data.NewIntValue(0)
		generator.autoKeyIndex = 1
	}
	return generator
}

type FuncYieldStackState struct {
	ctx          data.Context
	Fn           data.FuncStmt
	BodyIndex    int
	Body         []data.GetValue
	CurrentKey   data.Value
	CurrentValue data.Value
	autoKeyIndex int
	initialized  bool // true 表示已执行过第一次 yield（有 CurrentValue）
}

func (f *FuncYieldStackState) AsString() string {
	if f.CurrentValue == nil {
		return "Generator(closed)"
	}
	return fmt.Sprintf("Generator(%s)", f.Fn.GetName())
}

func (f *FuncYieldStackState) Current(_ data.Context) (data.Value, data.Control) {
	if f.CurrentValue == nil {
		return data.NewNullValue(), nil
	}
	return f.CurrentValue, nil
}

func (f *FuncYieldStackState) Key(_ data.Context) (data.Value, data.Control) {
	if f.CurrentKey == nil {
		return data.NewIntValue(0), nil
	}
	return f.CurrentKey, nil
}

func (f *FuncYieldStackState) Next(_ data.Context) data.Control {
	ctx := f.ctx
	// 从 BodyIndex 开始执行剩余的 body
	for bodyIndex := f.BodyIndex; bodyIndex < len(f.Body); bodyIndex++ {
		statement := f.Body[bodyIndex]

		var ctl data.Control
		_, ctl = statement.GetValue(ctx)

		if ctl != nil {
			switch rv := ctl.(type) {
			case data.YieldControl:
				// 判断是否是已经初始化过的 YieldFromControl（即 body 中当前语句就是这个控制流对象本身）
				if statement == rv {
					// 已初始化的 yield from：直接取已设置的 key/value
					key := rv.GetYieldKey()
					if key == nil {
						key = data.NewIntValue(f.autoKeyIndex)
						f.autoKeyIndex++
					}
					f.CurrentKey = key
					f.CurrentValue = rv.GetYieldValue()
					f.BodyIndex = bodyIndex // 保持在当前位置，下次继续执行这个 yield from
				} else {
					// 第一次遇到：调用 CreateStackState 初始化并推进到第一个元素
					newState := rv.CreateStackState(ctx, f.Fn, f.Body, bodyIndex)
					if ns, ok := newState.(*FuncYieldStackState); ok {
						f.Body = ns.Body
						f.BodyIndex = ns.BodyIndex
						f.CurrentKey = ns.CurrentKey
						f.CurrentValue = ns.CurrentValue
						f.initialized = true
					}
				}
				return nil
			case data.YieldValueControl:
				// 遇到 yield，更新键和值
				key := rv.GetYieldKey()
				if key == nil {
					key = data.NewIntValue(f.autoKeyIndex)
					f.autoKeyIndex++
				}
				f.CurrentKey = key
				f.CurrentValue = rv.GetYieldValue()
				f.BodyIndex = bodyIndex + 1
				return nil
			case data.AddStack:
				var from data.From
				if getFrom, ok := f.Fn.(GetFrom); ok {
					from = getFrom.GetFrom()
				}
				rv.AddStackWithInfo(from, "function", f.Fn.GetName())
				return ctl
			default:
				// 其他控制流，直接返回
				return ctl
			}
		}
		f.BodyIndex = bodyIndex + 1
	}
	f.CurrentKey = nil
	f.CurrentValue = nil
	return nil
}

func (f *FuncYieldStackState) Rewind(_ data.Context) (data.Value, data.Control) {
	// PHP 语义：第一次 rewind 时如果未初始化，先执行到第一个 yield
	if !f.initialized {
		f.initialized = true
		ctl := f.Next(f.ctx)
		return data.NewNullValue(), ctl
	}
	return data.NewNullValue(), nil
}

func (f *FuncYieldStackState) Valid(_ data.Context) (data.Value, data.Control) {
	// 如果未初始化，先执行到第一个 yield
	if !f.initialized {
		f.initialized = true
		ctl := f.Next(f.ctx)
		if ctl != nil {
			return data.NewBoolValue(false), ctl
		}
	}
	if f.CurrentValue == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *FuncYieldStackState) Send(_ data.Context, value data.Value) data.Control {
	// Send 将值保存到 CurrentValue，这样 Current() 可以获取到
	f.CurrentValue = value
	// 然后继续执行（类似 Next），这样 yield 语句可以获取到这个值
	return f.Next(f.ctx) // TODO 赋值语句需要匹配
}

func (f *FuncYieldStackState) Throw(_ data.Context) data.Control {
	// 抛出异常，清除键和值
	f.CurrentKey = nil
	f.CurrentValue = nil
	var from data.From
	if getFrom, ok := f.Fn.(GetFrom); ok {
		from = getFrom.GetFrom()
	}
	return data.NewErrorThrow(from, fmt.Errorf("生成器被异常终止"))
}

func (f *FuncYieldStackState) GetReturn(_ data.Context) (data.Value, data.Control) {
	// 执行完所有 body，查找返回值
	for f.BodyIndex < len(f.Body) {
		statement := f.Body[f.BodyIndex]
		f.BodyIndex++

		var ctl data.Control
		_, ctl = statement.GetValue(f.ctx)

		if ctl != nil {
			if rv, ok := ctl.(data.ReturnControl); ok {
				return rv.ReturnValue(), nil
			}
			// 其他控制流，继续执行
		}
	}
	// 没有找到 return，返回 null
	return data.NewNullValue(), nil
}

func (f *FuncYieldStackState) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 返回包装成类值的生成器，支持 $data->valid() 等调用
	generatorClass := NewGeneratorClass(f)
	return generatorClass.GetValue(ctx)
}
