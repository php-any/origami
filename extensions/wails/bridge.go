package wails

import (
	"fmt"
	"strconv"

	"github.com/php-any/origami/data"
)

// ============================================================================
// PHP <-> Go 桥接层
//
// Wails 在自己的协程上触发生命周期 / 事件回调，需要一个根上下文来为每次
// 回调创建子上下文并执行 PHP 闭包。wailsRootCtx 在 Application::run 调用时设置。
// ============================================================================

// wailsRootCtx 是 Application::run 调用时的上下文，用于驱动 PHP 回调。
var wailsRootCtx data.Context

// phpCallable 是 *data.FuncValue / *data.BoundFuncValue 共同实现的调用接口。
type phpCallable interface {
	Call(ctx data.Context) (data.GetValue, data.Control)
}

// invokeCallback 在 Wails 回调线程中执行一个 PHP 闭包。
// args 会按位置注入到闭包的参数中。返回闭包的返回值（如果有）。
func invokeCallback(cb data.Value, args ...data.Value) data.Value {
	if cb == nil || wailsRootCtx == nil {
		return nil
	}
	caller, ok := cb.(phpCallable)
	if !ok {
		return nil
	}
	callCtx := wailsRootCtx.CreateContext(make([]data.Variable, len(args)))
	for i, a := range args {
		if a == nil {
			a = data.NewNullValue()
		}
		callCtx.SetIndexZVal(i, data.NewZVal(a))
	}
	ret, ctrl := caller.Call(callCtx)
	if ctrl != nil {
		// 事件回调中不要让 PHP 异常冒泡到 Wails 线程（否则会弹原生错误框）
		if wailsApp != nil {
			wailsApp.Logger.Error("PHP callback error: " + fmt.Sprint(ctrl))
		}
		return nil
	}
	if v, ok := ret.(data.Value); ok {
		return v
	}
	return nil
}

// isCallable 判断一个值是否为可调用的 PHP 闭包。
func isCallable(v data.Value) bool {
	if v == nil {
		return false
	}
	_, ok := v.(phpCallable)
	return ok
}

// ============================================================================
// 值转换：PHP Value <-> Go any（用于事件数据的传递与 JSON 化）
// ============================================================================

// valueToGo 将 PHP 值转换为 Go 原生值。
func valueToGo(v data.Value) any {
	switch tv := v.(type) {
	case nil:
		return nil
	case *data.NullValue:
		return nil
	case *data.StringValue:
		return tv.Value
	case *data.BoolValue:
		b, _ := tv.AsBool()
		return b
	case *data.IntValue:
		n, _ := tv.AsInt()
		return n
	case *data.FloatValue:
		f, _ := tv.AsFloat()
		return f
	case *data.ArrayValue:
		return arrayValueToGo(tv)
	case *data.ObjectValue:
		// Origami 中 PHP 关联数组常以 ObjectValue 表示，必须在此转换，
		// 否则会落到下方 AsString 回退，得到 "Object {...}" 字符串。
		return objectValueToGo(tv)
	}
	// 回退：尽量按字符串处理
	if sv, ok := v.(data.AsString); ok {
		return sv.AsString()
	}
	return nil
}

// objectValueToGo 将 ObjectValue（关联数组 / 列表）转换为 Go 值。
// 键为连续整数 0..n-1 时视为列表（[]any），否则视为关联数组（map）。
func objectValueToGo(ov *data.ObjectValue) any {
	if ov == nil {
		return nil
	}
	keys := make([]string, 0)
	vals := make(map[string]any)
	ov.RangeProperties(func(k string, val data.Value) bool {
		keys = append(keys, k)
		vals[k] = valueToGo(val)
		return true
	})

	isList := true
	for i, k := range keys {
		if k != strconv.Itoa(i) {
			isList = false
			break
		}
	}
	if isList {
		list := make([]any, len(keys))
		for i, k := range keys {
			list[i] = vals[k]
		}
		return list
	}
	return vals
}

func arrayValueToGo(av *data.ArrayValue) any {
	if av == nil {
		return nil
	}
	isAssoc := false
	for _, z := range av.List {
		if z != nil && z.Name != "" {
			isAssoc = true
			break
		}
	}
	if isAssoc {
		m := make(map[string]any, len(av.List))
		for i, z := range av.List {
			if z == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			m[key] = valueToGo(z.Value)
		}
		return m
	}
	list := make([]any, 0, len(av.List))
	for _, z := range av.List {
		if z == nil {
			list = append(list, nil)
			continue
		}
		list = append(list, valueToGo(z.Value))
	}
	return list
}

// goToValue 将 Go 原生值（含 JSON 反序列化结果）转换为 PHP 值。
func goToValue(v any) data.Value {
	switch tv := v.(type) {
	case nil:
		return data.NewNullValue()
	case string:
		return data.NewStringValue(tv)
	case bool:
		return data.NewBoolValue(tv)
	case int:
		return data.NewIntValue(tv)
	case int64:
		return data.NewIntValue(int(tv))
	case float64:
		// JSON 数字统一为 float64；若为整数则转为 int
		if tv == float64(int(tv)) {
			return data.NewIntValue(int(tv))
		}
		return data.NewFloatValue(tv)
	case float32:
		return data.NewFloatValue(float64(tv))
	case []any:
		vals := make([]data.Value, 0, len(tv))
		for _, item := range tv {
			vals = append(vals, goToValue(item))
		}
		return data.NewArrayValue(vals)
	case map[string]any:
		av := &data.ArrayValue{}
		for k, item := range tv {
			av.List = append(av.List, data.NewNamedZVal(k, goToValue(item)))
		}
		return av
	}
	return data.NewNullValue()
}

// ============================================================================
// 延迟事件注册队列
//
// 用户脚本在 Application::run 之前调用 Events::on 注册监听器，此时 wailsApp
// 尚未创建。这里把注册请求排队，待 RunApp 创建 app 后再统一注册。
// ============================================================================

type pendingEventListener struct {
	name string
	cb   data.Value
	once bool
}

var pendingEventListeners []pendingEventListener

// registerEventListener 注册一个事件监听器。若 app 已创建则立即注册，
// 否则加入待处理队列。
func registerEventListener(name string, cb data.Value, once bool) {
	if wailsApp != nil {
		bindWailsEvent(name, cb, once)
		return
	}
	pendingEventListeners = append(pendingEventListeners, pendingEventListener{name: name, cb: cb, once: once})
}

// flushPendingEventListeners 在 app 创建后注册所有待处理的监听器。
func flushPendingEventListeners() {
	for _, p := range pendingEventListeners {
		bindWailsEvent(p.name, p.cb, p.once)
	}
	pendingEventListeners = nil
}
