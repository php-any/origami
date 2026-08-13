package stream

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// StreamContextCreateFunction 实现 stream_context_create 函数
//
//	stream_context_create(?array $options = null, ?array $params = null): resource
type StreamContextCreateFunction struct{}

func NewStreamContextCreateFunction() data.FuncStmt {
	return &StreamContextCreateFunction{}
}

func (f *StreamContextCreateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	var options map[string]map[string]string
	if v, _ := ctx.GetIndexValue(0); v != nil {
		if val, ok := v.(data.Value); ok {
			options = ParseStreamContextOptions(val)
		}
	} else {
		options = make(map[string]map[string]string)
	}

	var params map[string]string
	if v, _ := ctx.GetIndexValue(1); v != nil {
		if val, ok := v.(data.Value); ok {
			params = ParseStreamContextParams(val)
		}
	} else {
		params = make(map[string]string)
	}

	streamCtx := NewStreamContext(options, params)
	resourceClass := core.NewResourceClass("stream-context", streamCtx, allocStreamContextID())
	return core.NewResourceValue(resourceClass, ctx), nil
}

func (f *StreamContextCreateFunction) GetName() string {
	return "stream_context_create"
}

func (f *StreamContextCreateFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "params", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *StreamContextCreateFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewNullableType(data.NewBaseType("array"))),
		node.NewVariable(nil, "params", 1, data.NewNullableType(data.NewBaseType("array"))),
	}
}
