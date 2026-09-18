package httpfoundation

import (
	"net/http"

	"github.com/php-any/origami/data"
)

// HTTPRequestHolder 允许跨包（Illuminate Request state）提供底层 *http.Request。
type HTTPRequestHolder interface {
	HTTPRequest() *http.Request
}

// HTTPRequestFromValue 从 Symfony 或 Illuminate Request 对象取出 *http.Request。
func HTTPRequestFromValue(cv *data.ClassValue) *http.Request {
	if cv == nil {
		return nil
	}
	if class, ok := cv.Class.(*SymfonyRequestClass); ok && class.source != nil {
		return class.source.request
	}
	value, _ := cv.GetProperty("__illuminate_request_state")
	if anyValue, ok := value.(*data.AnyValue); ok {
		if holder, ok := anyValue.Value.(HTTPRequestHolder); ok {
			return holder.HTTPRequest()
		}
	}
	return nil
}

func tryRender(ctx data.Context, content data.Value) (string, bool, data.Control) {
	cv, ok := content.(*data.ClassValue)
	if !ok {
		return "", false, nil
	}

	callStringMethod := func(name string) (string, bool, data.Control) {
		if _, has := cv.GetMethod(name); !has {
			return "", false, nil
		}
		ret, ctl := callObjMethod(cv, name)
		if ctl != nil {
			return "", false, ctl
		}
		s, ok := stringifyRenderedValue(ret)
		return s, ok, nil
	}

	if s, ok, ctl := callStringMethod("render"); ctl != nil || ok {
		return s, ok, ctl
	}
	if s, ok, ctl := callStringMethod("toHtml"); ctl != nil || ok {
		return s, ok, ctl
	}
	if s, ok, ctl := callStringMethod("__toString"); ctl != nil || ok {
		return s, ok, ctl
	}
	return "", false, nil
}

func stringifyRenderedValue(ret data.GetValue) (string, bool) {
	if ret == nil {
		return "", false
	}
	switch t := ret.(type) {
	case *data.StringValue:
		return t.Value, true
	case *data.NullValue:
		return "", true
	case *data.ClassValue:
		return "", false
	case data.Value:
		return t.AsString(), true
	default:
		return "", false
	}
}
