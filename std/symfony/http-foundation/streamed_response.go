package httpfoundation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnStreamedResponse = "Symfony\\Component\\HttpFoundation\\StreamedResponse"

// StreamedResponseClass 实现 Symfony\Component\HttpFoundation\StreamedResponse。
type StreamedResponseClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewStreamedResponseClass() data.ClassStmt {
	c := &StreamedResponseClass{}
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("callback", 0, data.NewNullValue(), nil),
				param("status", 1, data.NewIntValue(200), nil),
				param("headers", 2, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				variable("callback", 0, nil),
				variable("status", 1, nil),
				variable("headers", 2, nil),
			},
			nil, streamedConstruct),
		pubMethod("setCallback",
			[]data.GetValue{param("callback", 0, nil, nil)},
			[]data.Variable{variable("callback", 0, nil)},
			nil, streamedSetCallback),
		pubMethod("getCallback", nil, nil, nil, streamedGetCallback),
		pubMethod("setChunks",
			[]data.GetValue{param("chunks", 0, nil, nil)},
			[]data.Variable{variable("chunks", 0, nil)},
			nil, streamedSetChunks),
		pubMethod("sendContent", nil, nil, nil, streamedSendContent),
		pubMethod("send",
			[]data.GetValue{param("flush", 0, data.NewBoolValue(true), nil)},
			[]data.Variable{variable("flush", 0, nil)},
			nil, streamedSend),
		pubMethod("setContent",
			[]data.GetValue{param("content", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("content", 0, nil)},
			nil, streamedSetContent),
		pubMethod("getContent", nil, nil, nil, streamedGetContent),
	}
	c.methods = indexMethods(list)
	c.methodList = list
	return c
}

func (c *StreamedResponseClass) GetName() string { return fqnStreamedResponse }
func (c *StreamedResponseClass) GetExtend() *string {
	parent := fqnSymfonyResponse
	return &parent
}
func (c *StreamedResponseClass) GetImplements() []string { return nil }
func (c *StreamedResponseClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.GetPropertyList() {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *StreamedResponseClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "callback", "protected", false, data.NewNullValue()),
		node.NewProperty(nil, "streamed", "protected", false, data.NewBoolValue(false)),
		node.NewProperty(nil, "headersSent", "private", false, data.NewBoolValue(false)),
		node.NewProperty(nil, "chunks", "private", false, data.NewNullValue()),
	}
}
func (c *StreamedResponseClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *StreamedResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *StreamedResponseClass) GetMethod(name string) (data.Method, bool) {
	return getIndexedMethod(c.methods, name)
}
func (c *StreamedResponseClass) GetMethods() []data.Method { return c.methodList }

func streamedConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	callback, _ := ctx.GetIndexValue(0)
	status := intParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)
	headers := createResponseHeaders(ctx, headersVal)
	_ = cv.SetProperty("headers", headers)
	_ = cv.SetProperty("version", data.NewStringValue("1.0"))
	_ = cv.SetProperty("content", data.NewStringValue(""))
	if ctl := applyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("streamed", data.NewBoolValue(false))
	_ = cv.SetProperty("headersSent", data.NewBoolValue(false))
	_ = cv.SetProperty("chunks", data.NewNullValue())
	_ = cv.SetProperty("callback", data.NewNullValue())
	if callback != nil && !isNull(callback) {
		if isPHPCallable(ctx, callback) {
			_ = cv.SetProperty("callback", callback)
		} else {
			_ = cv.SetProperty("chunks", callback)
		}
	}
	return data.NewNullValue(), nil
}

func streamedSetCallback(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	cb, _ := ctx.GetIndexValue(0)
	if cb == nil {
		cb = data.NewNullValue()
	}
	_ = cv.SetProperty("callback", cb)
	_ = cv.SetProperty("chunks", data.NewNullValue())
	return responseSelf(ctx), nil
}

func streamedGetCallback(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := cv.GetProperty("callback")
	if v == nil || isNull(v) {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func streamedSetChunks(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	chunks, _ := ctx.GetIndexValue(0)
	if chunks == nil {
		chunks = data.NewNullValue()
	}
	_ = cv.SetProperty("chunks", chunks)
	_ = cv.SetProperty("callback", data.NewNullValue())
	return responseSelf(ctx), nil
}

func streamedSend(ctx data.Context) (data.GetValue, data.Control) {
	if _, ctl := symfonyResponseSendHeaders(ctx); ctl != nil {
		return nil, ctl
	}
	return streamedSendContent(ctx)
}

func streamedSendContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	if cookiePropBool(cv, "streamed", false) {
		return responseSelf(ctx), nil
	}
	_ = cv.SetProperty("streamed", data.NewBoolValue(true))

	chunks, _ := cv.GetProperty("chunks")
	if chunks != nil && !isNull(chunks) {
		if ctl := emitChunks(ctx, chunks); ctl != nil {
			return nil, ctl
		}
		return responseSelf(ctx), nil
	}

	cb, _ := cv.GetProperty("callback")
	if cb == nil || isNull(cb) {
		return nil, throwNamed("LogicException", "The Response callback must be set.")
	}
	if !isPHPCallable(ctx, cb) {
		return nil, throwNamed("LogicException", "The Response callback must be set.")
	}
	if ctl := invokeCallable(ctx, cb); ctl != nil {
		return nil, ctl
	}
	return responseSelf(ctx), nil
}

func streamedSetContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	content, _ := ctx.GetIndexValue(0)
	if content != nil && !isNull(content) {
		return nil, throwNamed("LogicException", "The content cannot be set on a StreamedResponse instance.")
	}
	_ = cv.SetProperty("streamed", data.NewBoolValue(true))
	return responseSelf(ctx), nil
}

func streamedGetContent(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

func emitChunks(ctx data.Context, chunks data.Value) data.Control {
	switch t := chunks.(type) {
	case *data.ArrayValue:
		for _, z := range t.List {
			if z == nil || z.Value == nil {
				continue
			}
			if ctl := data.EmitOutput(ctx, z.Value.AsString()); ctl != nil {
				return ctl
			}
		}
	default:
		m, err := valueToAssocMap(chunks)
		if err != nil {
			return data.EmitOutput(ctx, chunks.AsString())
		}
		for _, v := range m {
			if v == nil {
				continue
			}
			if ctl := data.EmitOutput(ctx, v.AsString()); ctl != nil {
				return ctl
			}
		}
	}
	return nil
}

func isPHPCallable(ctx data.Context, v data.Value) bool {
	if v == nil || isNull(v) {
		return false
	}
	switch t := v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return true
	case *data.ClassValue:
		_, ok := t.GetMethod("__invoke")
		return ok
	case *data.StringValue:
		if ctx == nil || ctx.GetVM() == nil {
			return false
		}
		_, ok := ctx.GetVM().GetFunc(t.AsString())
		return ok
	case *data.ArrayValue:
		vals := t.ToValueList()
		if len(vals) != 2 {
			return false
		}
		if _, ok := vals[1].(data.AsString); !ok {
			return false
		}
		if _, ok := vals[0].(*data.ClassValue); ok {
			return true
		}
		if sv, ok := vals[0].(*data.StringValue); ok && ctx != nil && ctx.GetVM() != nil {
			_, found := ctx.GetVM().GetClass(sv.AsString())
			return found
		}
	}
	return false
}

func invokeCallable(ctx data.Context, cb data.Value) data.Control {
	if cb == nil || isNull(cb) {
		return nil
	}
	switch v := cb.(type) {
	case *data.BoundFuncValue:
		_, ctl := v.Call(ctx.CreateBaseContext())
		return ctl
	case *data.FuncValue:
		_, ctl := v.Call(ctx.CreateBaseContext())
		return ctl
	case *data.ClassValue:
		if m, ok := v.GetMethod("__invoke"); ok && m != nil {
			_, ctl := m.Call(v.CreateContext(m.GetVariables()))
			return ctl
		}
	case *data.ArrayValue:
		vals := v.ToValueList()
		if len(vals) >= 2 {
			name := vals[1].AsString()
			if obj, ok := vals[0].(*data.ClassValue); ok {
				_, ctl := callObjMethod(obj, name)
				return ctl
			}
			if ctx.GetVM() != nil {
				if stmt, found := ctx.GetVM().GetClass(vals[0].AsString()); found && stmt != nil {
					if m, ok := stmt.GetMethod(name); ok && m != nil {
						_, ctl := m.Call(ctx.CreateBaseContext())
						return ctl
					}
				}
			}
		}
	case *data.StringValue:
		if ctx.GetVM() != nil {
			if fn, ok := ctx.GetVM().GetFunc(v.AsString()); ok && fn != nil {
				_, ctl := fn.Call(ctx.CreateBaseContext())
				return ctl
			}
		}
	}
	return nil
}
