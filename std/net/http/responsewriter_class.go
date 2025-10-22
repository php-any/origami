package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewResponseWriterClass() data.ClassStmt {
	return &ResponseWriterClass{
		source: nil,
	}
}

func NewResponseWriterClassFrom(source httpsrc.ResponseWriter) data.ClassStmt {
	return &ResponseWriterClass{
		source: source,
	}
}

type ResponseWriterClass struct {
	node.Node
	source httpsrc.ResponseWriter
}

func (s *ResponseWriterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewResponseWriterClassFrom(nil), ctx.CreateBaseContext()), nil
}

func (s *ResponseWriterClass) GetName() string         { return "Net\\Http\\Response" }
func (s *ResponseWriterClass) GetExtend() *string      { return nil }
func (s *ResponseWriterClass) GetImplements() []string { return nil }
func (s *ResponseWriterClass) AsString() string        { return "Response{}" }
func (s *ResponseWriterClass) GetSource() any          { return s.source }
func (s *ResponseWriterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "json":
		return &ResponseWriterJsonMethod{source: s.source}, true
	case "header":
		return &ResponseWriterHeaderMethod{source: s.source}, true
	case "write":
		return &ResponseWriterWriteMethod{source: s.source}, true
	case "writeHeader":
		return &ResponseWriterWriteHeaderMethod{source: s.source}, true
	}
	return nil, false
}

func (s *ResponseWriterClass) GetMethods() []data.Method {
	return []data.Method{
		&ResponseWriterJsonMethod{source: s.source},
		&ResponseWriterHeaderMethod{source: s.source},
		&ResponseWriterWriteMethod{source: s.source},
		&ResponseWriterWriteHeaderMethod{source: s.source},
	}
}

func (s *ResponseWriterClass) GetConstruct() data.Method { return nil }

func (s *ResponseWriterClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (s *ResponseWriterClass) GetPropertyList() []data.Property {
	return []data.Property{}
}
