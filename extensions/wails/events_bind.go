package wails

import (
	"github.com/php-any/origami/data"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// bindWailsEvent 把一个 PHP 回调绑定到 Wails 的自定义事件总线上。
// 当事件被 Emit 时，事件数据会被转换为 PHP 值并传入回调。
func bindWailsEvent(name string, cb data.Value, once bool) {
	if wailsApp == nil || !isCallable(cb) {
		return
	}
	handler := func(e *application.CustomEvent) {
		var arg data.Value = data.NewNullValue()
		if e != nil && e.Data != nil {
			arg = goToValue(e.Data)
		}
		invokeCallback(cb, arg)
	}
	if once {
		wailsApp.Event.OnMultiple(name, handler, 1)
	} else {
		wailsApp.Event.On(name, handler)
	}
}

// emitWailsEvent 从 PHP 端发送一个自定义事件到事件总线。
func emitWailsEvent(name string, payload data.Value) {
	if wailsApp == nil {
		return
	}
	var goData any
	if payload != nil {
		goData = valueToGo(payload)
	}
	wailsApp.Event.EmitEvent(&application.CustomEvent{Name: name, Data: goData})
}
