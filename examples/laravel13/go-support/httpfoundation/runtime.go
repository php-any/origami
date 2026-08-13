package httpfoundation

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// SendResponse 调用 Symfony Response::send，将 Header 与 Body 写入当前 Go writer。
func SendResponse(response data.GetValue) (*data.ClassValue, data.Control) {
	value, ok := response.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: Kernel 返回值不是 Response"))
	}
	method, ok := value.GetMethod("send")
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: Response 缺少 send 方法"))
	}
	callCtx := value.CreateContext(method.GetVariables())
	result, control := method.Call(callCtx)
	if control != nil {
		return nil, control
	}
	if sent, ok := result.(*data.ClassValue); ok {
		return sent, nil
	}
	return value, nil
}
