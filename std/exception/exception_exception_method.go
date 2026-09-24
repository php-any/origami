package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ExceptionExceptionMethod struct {
	source *Exception
}

func (h *ExceptionExceptionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	msg, acl := ctx.GetVariableValue(node.NewVariable(nil, "msg", 0, nil))
	if acl != nil {
		return nil, acl
	}

	message := ""
	if msg != nil {
		if sv, ok := msg.(*data.StringValue); ok {
			message = sv.AsString()
		} else if v, ok := msg.(data.Value); ok {
			message = v.AsString()
		}
	}

	code := 0
	if codeVal, acl := ctx.GetVariableValue(node.NewVariable(nil, "code", 1, nil)); acl == nil && codeVal != nil {
		if iv, ok := codeVal.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				code = n
			}
		}
	}

	var previous data.Value = data.NewNullValue()
	if prevVal, acl := ctx.GetVariableValue(node.NewVariable(nil, "previous", 2, nil)); acl == nil && prevVal != nil {
		if v, ok := prevVal.(data.Value); ok {
			previous = v
		}
	}

	// skip=2：跳过 captureGoTrace 与本 Call，对齐原先 captureTrace 的 Callers(3)。
	file, line, frames := captureGoTrace(2)
	if h.source != nil {
		h.source.msg = message
		h.source.file = file
		h.source.line = line
		h.source.trace = frames
	}

	// 同步到 PHP 可见的受保护属性（子类可 $this->message = ...）
	setInstanceProperty(ctx, "message", data.NewStringValue(message))
	setInstanceProperty(ctx, "code", data.NewIntValue(code))
	setInstanceProperty(ctx, "previous", previous)
	setInstanceProperty(ctx, "file", data.NewStringValue(file))
	setInstanceProperty(ctx, "line", data.NewIntValue(line))
	setInstanceProperty(ctx, "trace", data.NewArrayValue(traceFramesToValues(frames)))

	return nil, nil
}

func (h *ExceptionExceptionMethod) GetName() string {
	return "exception"
}

func (h *ExceptionExceptionMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ExceptionExceptionMethod) GetIsStatic() bool {
	return false
}

var exceptionExceptionMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "msg", 0, data.NewStringValue(""), nil),
	node.NewParameter(nil, "code", 1, data.NewIntValue(0), nil),
	node.NewParameter(nil, "previous", 2, data.NewNullValue(), nil),
}

func (h *ExceptionExceptionMethod) GetParams() []data.GetValue {
	return exceptionExceptionMethodGetParams
}

var exceptionExceptionMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "msg", 0, nil),
	node.NewVariable(nil, "code", 1, nil),
	node.NewVariable(nil, "previous", 2, nil),
}

func (h *ExceptionExceptionMethod) GetVariables() []data.Variable {
	return exceptionExceptionMethodGetVariables
}

// GetReturnType 返回方法返回类型
func (h *ExceptionExceptionMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
