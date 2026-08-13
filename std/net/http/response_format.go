package http

import (
	"errors"
	"fmt"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/serializer/json"
)

type formatHandlerSlot struct {
	fn  data.FuncStmt
	ctx data.Context
}

func defaultFormattedPayload(code int, message string, payload data.Value) data.Value {
	obj := data.NewObjectValue()
	obj.SetProperty("code", data.NewIntValue(code))
	obj.SetProperty("message", data.NewStringValue(message))
	if payload != nil {
		obj.SetProperty("data", payload)
	} else {
		obj.SetProperty("data", data.NewNullValue())
	}
	obj.SetProperty("timestamp", data.NewIntValue(int(time.Now().Unix())))
	return obj
}

func invokeFormatHandler(slot *formatHandlerSlot, ctx data.Context, code int, message string, payload data.Value) (data.Value, error) {
	if slot == nil || slot.fn == nil {
		return defaultFormattedPayload(code, message, payload), nil
	}

	vars := slot.fn.GetVariables()
	if len(vars) < 3 {
		return nil, errors.New("onFormat 闭包需要 3 个参数: ($code, $message, $data)")
	}

	mctx := slot.ctx.CreateContext(vars)
	mctx.SetVariableValue(vars[0], data.NewIntValue(code))
	mctx.SetVariableValue(vars[1], data.NewStringValue(message))
	if payload != nil {
		mctx.SetVariableValue(vars[2], payload)
	} else {
		mctx.SetVariableValue(vars[2], data.NewNullValue())
	}

	ret, acl := slot.fn.Call(mctx)
	if acl != nil {
		return nil, fmt.Errorf("onFormat 闭包执行失败: %v", acl)
	}
	if ret == nil {
		return nil, errors.New("onFormat 闭包必须返回 array 或 object")
	}

	val, ok := ret.(data.Value)
	if !ok {
		return nil, errors.New("onFormat 闭包必须返回 array 或 object")
	}
	if _, ok := val.(data.ValueSerializer); !ok {
		return nil, errors.New("onFormat 返回值无法 JSON 序列化")
	}
	return val, nil
}

func buildFormattedPayload(bw *bufferedWriter, ctx data.Context, code int, message string, payload data.Value) (data.Value, error) {
	if bw != nil && bw.formatter != nil {
		return invokeFormatHandler(bw.formatter, ctx, code, message, payload)
	}
	return defaultFormattedPayload(code, message, payload), nil
}

func writeFormattedResponse(bw *bufferedWriter, ctx data.Context, code int, message string, payload data.Value) error {
	formatted, err := buildFormattedPayload(bw, ctx, code, message, payload)
	if err != nil {
		return err
	}

	serializer, ok := formatted.(data.ValueSerializer)
	if !ok {
		return errors.New("格式化结果无法 JSON 序列化")
	}

	bytes, err := serializer.Marshal(json.NewJsonSerializer())
	if err != nil {
		return err
	}

	bw.SetStatus(code)
	return bw.WriteJSON(bytes)
}
