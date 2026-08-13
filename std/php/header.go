package php

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HeaderFunction 实现 header 函数（与 Net/Http 服务器集成）
type HeaderFunction struct{}

func NewHeaderFunction() data.FuncStmt { return &HeaderFunction{} }

func (f *HeaderFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	writer := responseWriterFromContext(ctx)
	if writer == nil {
		return data.NewNullValue(), nil
	}

	headerValue, _ := ctx.GetIndexValue(0)
	if headerValue == nil {
		return data.NewNullValue(), nil
	}
	line := headerValue.AsString()

	replace := true
	if value, ok := ctx.GetIndexValue(1); ok {
		if boolean, ok := value.(data.AsBool); ok {
			replace, _ = boolean.AsBool()
		}
	}
	responseCode := 0
	if value, ok := ctx.GetIndexValue(2); ok {
		if integer, ok := value.(data.AsInt); ok {
			responseCode, _ = integer.AsInt()
		}
	}

	if strings.HasPrefix(strings.ToUpper(line), "HTTP/") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			if code, err := strconv.Atoi(parts[1]); err == nil {
				writer.WriteHeader(code)
			}
		}
	} else if name, value, ok := strings.Cut(line, ":"); ok {
		name = http.CanonicalHeaderKey(strings.TrimSpace(name))
		value = strings.TrimSpace(value)
		if replace {
			writer.Header().Set(name, value)
		} else {
			writer.Header().Add(name, value)
		}
	}
	if responseCode > 0 {
		writer.WriteHeader(responseCode)
	}
	return data.NewNullValue(), nil
}

type httpResponseHost interface {
	HTTPResponseWriter() http.ResponseWriter
}

func responseWriterFromContext(ctx data.Context) http.ResponseWriter {
	if ctx != nil {
		if host, ok := ctx.GetVM().(httpResponseHost); ok {
			if w := host.HTTPResponseWriter(); w != nil {
				return w
			}
		}
	}
	return node.HTTPResponseWriter()
}

func (f *HeaderFunction) GetName() string { return "header" }
func (f *HeaderFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "header", 0, nil, nil),
		node.NewParameter(nil, "replace", 1, node.NewBooleanLiteral(nil, true), nil),
		node.NewParameter(nil, "response_code", 2, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (f *HeaderFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "header", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "replace", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "response_code", 2, data.NewBaseType("int")),
	}
}
