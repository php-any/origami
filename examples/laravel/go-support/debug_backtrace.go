package gosupport

import (
	"runtime"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DebugBacktraceFunction implements debug_backtrace() for Telescope QueryWatcher.
type DebugBacktraceFunction struct{}

func (f *DebugBacktraceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	limit := 0
	if _, has := ctx.GetIndexValue(1); has {
		limit = valueAsInt(mustIndex(ctx, 1), 0)
	}

	max := 32
	if limit > 0 && limit < max {
		max = limit
	}

	pcs := make([]uintptr, max)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	list := make([]*data.ZVal, 0, n)
	for {
		frame, more := frames.Next()
		if frame.File == "" && !more {
			break
		}

		obj := data.NewObjectValue()
		obj.SetProperty("file", data.NewStringValue(frame.File))
		obj.SetProperty("line", data.NewIntValue(frame.Line))
		if frame.Function != "" {
			obj.SetProperty("function", data.NewStringValue(frame.Function))
		}
		list = append(list, data.NewZVal(obj))

		if !more {
			break
		}
	}

	return &data.ArrayValue{List: list}, nil
}

func (f *DebugBacktraceFunction) GetName() string            { return "debug_backtrace" }
func (f *DebugBacktraceFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *DebugBacktraceFunction) GetIsStatic() bool          { return false }
func (f *DebugBacktraceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewIntValue(0), nil),
		node.NewParameter(nil, "limit", 1, data.NewIntValue(0), nil),
	}
}
func (f *DebugBacktraceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, nil),
		node.NewVariable(nil, "limit", 1, nil),
	}
}
func (f *DebugBacktraceFunction) GetReturnType() data.Types { return data.NewBaseType("array") }
