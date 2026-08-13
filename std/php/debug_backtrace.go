package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	debugBacktraceProvideObject = 1
	debugBacktraceIgnoreArgs    = 2
)

// DebugBacktraceFunction 实现 debug_backtrace([int $options = DEBUG_BACKTRACE_PROVIDE_OBJECT [, int $limit = 0]]): array
type DebugBacktraceFunction struct{}

func NewDebugBacktraceFunction() data.FuncStmt { return &DebugBacktraceFunction{} }

func (f *DebugBacktraceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	limit := 0
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if asInt, ok := v.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				limit = n
			}
		}
	}

	var frames []data.CallFrame
	if vm := ctx.GetVM(); vm != nil {
		if tracker, ok := vm.(data.CallStackTracker); ok {
			frames = tracker.SnapshotCallStack()
		}
	}

	// Snapshot 为从底到顶；PHP debug_backtrace 为从当前到外层（最近调用在前）
	out := make([]data.Value, 0, len(frames))
	for i := len(frames) - 1; i >= 0; i-- {
		fr := frames[i]
		obj := data.NewObjectValue()
		if fr.File != "" {
			obj.SetProperty("file", data.NewStringValue(fr.File))
		}
		if fr.Line > 0 {
			obj.SetProperty("line", data.NewIntValue(fr.Line))
		}
		if fr.Function != "" {
			obj.SetProperty("function", data.NewStringValue(fr.Function))
		}
		if fr.Class != "" {
			obj.SetProperty("class", data.NewStringValue(fr.Class))
		}
		if fr.Type != "" {
			obj.SetProperty("type", data.NewStringValue(fr.Type))
		}
		out = append(out, obj)
		if limit > 0 && len(out) >= limit {
			break
		}
	}

	return data.NewArrayValue(out), nil
}

func (f *DebugBacktraceFunction) GetName() string { return "debug_backtrace" }
func (f *DebugBacktraceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewIntValue(debugBacktraceProvideObject), nil),
		node.NewParameter(nil, "limit", 1, data.NewIntValue(0), nil),
	}
}
func (f *DebugBacktraceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, nil),
		node.NewVariable(nil, "limit", 1, nil),
	}
}
