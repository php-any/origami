package http

import (
	"net/http"

	"github.com/php-any/origami/data"
	jsonSerializer "github.com/php-any/origami/std/serializer/json"
	"github.com/php-any/origami/utils"
)

// writeHandlerReturnValue 当控制器返回 array/object 且尚未写出响应时，自动 JSON 序列化。
func writeHandlerReturnValue(resProxy data.Value, ret data.GetValue) data.Control {
	if ret == nil {
		return nil
	}
	if _, isNull := ret.(*data.NullValue); isNull {
		return nil
	}

	bw := bufferedWriterFromProxy(resProxy)
	if bw == nil || bw.headerSent {
		return nil
	}

	val, ok := ret.(data.Value)
	if !ok {
		return nil
	}

	serializer, ok := val.(data.ValueSerializer)
	if !ok {
		return nil
	}

	bytes, err := serializer.Marshal(jsonSerializer.NewJsonSerializer())
	if err != nil {
		return utils.NewThrow(err)
	}
	if !bw.statusSet {
		bw.SetStatus(http.StatusOK)
	}
	if err := bw.WriteJSON(bytes); err != nil {
		return utils.NewThrow(err)
	}
	return nil
}

func bufferedWriterFromProxy(resProxy data.Value) *bufferedWriter {
	resCV, ok := resProxy.(*data.ProxyValue)
	if !ok {
		return nil
	}
	if bw, ok := resCV.GetSource().(*bufferedWriter); ok {
		return bw
	}
	if rw, ok := resCV.Class.(interface{ GetSource() any }); ok {
		if bw, ok := rw.GetSource().(*bufferedWriter); ok {
			return bw
		}
	}
	return nil
}
