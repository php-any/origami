package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ErrorExceptionConstructMethod 对齐 PHP ErrorException::__construct：
//
//	__construct(
//	    string $message = "",
//	    int $code = 0,
//	    int $severity = E_ERROR,
//	    ?string $filename = null,
//	    ?int $line = null,
//	    ?Throwable $previous = null
//	)
//
// Laravel CompilerEngine 以 6 参包装视图异常：
// new ViewException($msg, 0, 1, $e->getFile(), $e->getLine(), $e)
// 若误用 Exception 的 3 参签名，第 3 实参 severity(int) 会写进 previous，
// ViewException::report() → Reflector::isCallable → ReflectionMethod(int, ...) 类型错误。
type ErrorExceptionConstructMethod struct {
	source *Exception
}

func (h *ErrorExceptionConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	message := ""
	if msg := indexArg(ctx, 0); msg != nil {
		if _, isNull := msg.(*data.NullValue); !isNull {
			message = msg.AsString()
		}
	}

	code := 0
	if codeVal := indexArg(ctx, 1); codeVal != nil {
		if iv, ok := codeVal.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				code = n
			}
		}
	}

	severity := 1 // E_ERROR
	if sevVal := indexArg(ctx, 2); sevVal != nil {
		if iv, ok := sevVal.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				severity = n
			}
		}
	}

	var previous data.Value = data.NewNullValue()
	if prevVal := indexArg(ctx, 5); prevVal != nil {
		previous = prevVal
	}

	file, line, frames := captureGoTrace(2)
	if filenameVal := indexArg(ctx, 3); filenameVal != nil {
		if _, isNull := filenameVal.(*data.NullValue); !isNull {
			file = filenameVal.AsString()
		}
	}
	if lineVal := indexArg(ctx, 4); lineVal != nil {
		if _, isNull := lineVal.(*data.NullValue); !isNull {
			if iv, ok := lineVal.(data.AsInt); ok {
				if n, err := iv.AsInt(); err == nil {
					line = n
				}
			}
		}
	}

	if h.source != nil {
		h.source.msg = message
		h.source.file = file
		h.source.line = line
		h.source.trace = frames
	}

	setInstanceProperty(ctx, "message", data.NewStringValue(message))
	setInstanceProperty(ctx, "code", data.NewIntValue(code))
	setInstanceProperty(ctx, "severity", data.NewIntValue(severity))
	setInstanceProperty(ctx, "previous", previous)
	setInstanceProperty(ctx, "file", data.NewStringValue(file))
	setInstanceProperty(ctx, "line", data.NewIntValue(line))
	setInstanceProperty(ctx, "trace", data.NewArrayValue(traceFramesToValues(frames)))

	return nil, nil
}

func indexArg(ctx data.Context, i int) data.Value {
	if ctx == nil {
		return nil
	}
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil {
		return nil
	}
	return v
}

func (h *ErrorExceptionConstructMethod) GetName() string { return "__construct" }

func (h *ErrorExceptionConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (h *ErrorExceptionConstructMethod) GetIsStatic() bool { return false }

func (h *ErrorExceptionConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "message", 0, data.NewStringValue(""), nil),
		node.NewParameter(nil, "code", 1, data.NewIntValue(0), nil),
		node.NewParameter(nil, "severity", 2, data.NewIntValue(1), nil),
		node.NewParameter(nil, "filename", 3, data.NewNullValue(), nil),
		node.NewParameter(nil, "line", 4, data.NewNullValue(), nil),
		node.NewParameter(nil, "previous", 5, data.NewNullValue(), nil),
	}
}

func (h *ErrorExceptionConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "message", 0, nil),
		node.NewVariable(nil, "code", 1, nil),
		node.NewVariable(nil, "severity", 2, nil),
		node.NewVariable(nil, "filename", 3, nil),
		node.NewVariable(nil, "line", 4, nil),
		node.NewVariable(nil, "previous", 5, nil),
	}
}

func (h *ErrorExceptionConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
