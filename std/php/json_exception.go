package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/exception"
)

const jsonThrowOnError = 4194304
const jsonPartialOutputOnError = 512

type jsonExceptionClass struct{ *exception.ExceptionClass }

func (*jsonExceptionClass) GetName() string    { return "JsonException" }
func (*jsonExceptionClass) GetExtend() *string { name := "Exception"; return &name }
func (c *jsonExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func jsonException(ctx data.Context, code int) data.Control {
	class, _ := ctx.GetVM().GetClass("JsonException")
	value := data.NewClassValue(class, ctx)
	value.SetProperty("message", data.NewStringValue(jsonErrorMessage(code)))
	value.SetProperty("code", data.NewIntValue(code))
	return data.NewErrorThrowFromClassValue(nil, value)
}

func jsonDecodeFailure(ctx data.Context, flags, code int) (data.GetValue, data.Control) {
	if flags&jsonThrowOnError != 0 {
		return nil, jsonException(ctx, code)
	}
	setJsonLastError(code, "")
	return data.NewNullValue(), nil
}

func jsonEncodeFailure(ctx data.Context, flags, code int) (string, bool, data.Control) {
	if flags&jsonThrowOnError != 0 && flags&jsonPartialOutputOnError == 0 {
		return "", false, jsonException(ctx, code)
	}
	setJsonLastError(code, "")
	return "", false, nil
}
