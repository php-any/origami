package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewResponseWriterClass() data.ClassStmt {
	return &ResponseWriterClass{w: newBufferedWriter(nil)}
}

func NewResponseWriterClassFrom(w httpsrc.ResponseWriter) data.ClassStmt {
	return &ResponseWriterClass{w: newBufferedWriter(w)}
}

type ResponseWriterClass struct {
	node.Node
	w *bufferedWriter
}

func (s *ResponseWriterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewResponseWriterClass(), ctx.CreateBaseContext()), nil
}

func (s *ResponseWriterClass) GetName() string         { return "Net\\Http\\Response" }
func (s *ResponseWriterClass) GetExtend() *string      { return nil }
func (s *ResponseWriterClass) GetImplements() []string { return nil }
func (s *ResponseWriterClass) AsString() string        { return "Response{}" }
func (s *ResponseWriterClass) GetSource() any          { return s.w }

func (s *ResponseWriterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "json":
		return &ResponseWriterJsonMethod{w: s.w}, true
	case "header":
		return &ResponseWriterHeaderMethod{w: s.w}, true
	case "write":
		return &ResponseWriterWriteMethod{w: s.w}, true
	case "writeHeader":
		return &ResponseWriterWriteHeaderMethod{w: s.w}, true
	case "view":
		return &ResponseWriterViewMethod{w: s.w}, true
	case "status":
		return &ResponseWriterStatusMethod{w: s.w}, true
	case "redirect":
		return &ResponseWriterRedirectMethod{w: s.w}, true
	case "noContent":
		return &ResponseWriterNoContentMethod{w: s.w}, true
	case "cookie":
		return &ResponseWriterCookieMethod{w: s.w}, true
	case "html":
		return &ResponseWriterHtmlMethod{w: s.w}, true
	case "file":
		return &ResponseWriterFileMethod{w: s.w}, true
	}
	return nil, false
}

func (s *ResponseWriterClass) GetMethods() []data.Method {
	return []data.Method{
		&ResponseWriterJsonMethod{w: s.w},
		&ResponseWriterHeaderMethod{w: s.w},
		&ResponseWriterWriteMethod{w: s.w},
		&ResponseWriterWriteHeaderMethod{w: s.w},
		&ResponseWriterViewMethod{w: s.w},
		&ResponseWriterStatusMethod{w: s.w},
		&ResponseWriterRedirectMethod{w: s.w},
		&ResponseWriterNoContentMethod{w: s.w},
		&ResponseWriterCookieMethod{w: s.w},
		&ResponseWriterHtmlMethod{w: s.w},
		&ResponseWriterFileMethod{w: s.w},
	}
}

func (s *ResponseWriterClass) GetConstruct() data.Method { return nil }

func (s *ResponseWriterClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (s *ResponseWriterClass) GetPropertyList() []data.Property {
	return []data.Property{}
}
