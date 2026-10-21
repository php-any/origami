package php

import "github.com/php-any/origami/data"

type HTTPResponseHeadersFunction struct{ Clear bool }

func (f *HTTPResponseHeadersFunction) GetName() string {
	if f.Clear {
		return "http_clear_last_response_headers"
	}
	return "http_get_last_response_headers"
}
func (*HTTPResponseHeadersFunction) GetParams() []data.GetValue    { return nil }
func (*HTTPResponseHeadersFunction) GetVariables() []data.Variable { return nil }
func (f *HTTPResponseHeadersFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := ctx.GetVM().(interface {
		SetLastHTTPResponseHeaders([]string)
		LastHTTPResponseHeaders() []string
	})
	if !ok {
		return data.NewNullValue(), nil
	}
	if f.Clear {
		host.SetLastHTTPResponseHeaders(nil)
		return data.NewNullValue(), nil
	}
	headers := host.LastHTTPResponseHeaders()
	if headers == nil {
		return data.NewNullValue(), nil
	}
	values := make([]data.Value, len(headers))
	for i, header := range headers {
		values[i] = data.NewStringValue(header)
	}
	return data.NewArrayValue(values), nil
}
