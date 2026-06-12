package core

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// outputBuffer 缓冲条目
type outputBuffer struct {
	builder   strings.Builder
	flushable bool
	callback  string
}

// outputBufferStack 输出缓冲栈
type outputBufferStack struct {
	buffers []outputBuffer
}

var obStack = &outputBufferStack{buffers: []outputBuffer{{flushable: true}}}

func (s *outputBufferStack) push() {
	s.buffers = append(s.buffers, outputBuffer{flushable: true})
	s.syncWriter()
}

func (s *outputBufferStack) pushWith(callback string, flushable bool) {
	s.buffers = append(s.buffers, outputBuffer{flushable: flushable, callback: callback})
	s.syncWriter()
}

func (s *outputBufferStack) pop() string {
	if len(s.buffers) <= 1 {
		return ""
	}
	last := s.buffers[len(s.buffers)-1].builder.String()
	s.buffers = s.buffers[:len(s.buffers)-1]
	s.syncWriter()
	return last
}

func (s *outputBufferStack) contents() string {
	if len(s.buffers) == 0 {
		return ""
	}
	return s.buffers[len(s.buffers)-1].builder.String()
}

func (s *outputBufferStack) flush() (string, bool) {
	if len(s.buffers) <= 1 {
		return "", false
	}
	top := &s.buffers[len(s.buffers)-1]
	if !top.flushable {
		return "", false
	}
	content := top.builder.String()
	top.builder.Reset()
	return content, true
}

func (s *outputBufferStack) syncWriter() {
	if len(s.buffers) <= 1 {
		data.WriteOutput = data.DefaultOutputWriter
		return
	}
	buf := &s.buffers[len(s.buffers)-1].builder
	data.WriteOutput = func(str string) {
		buf.WriteString(str)
	}
}

// ObStartFunction 实现 ob_start
type ObStartFunction struct{}

func NewObStartFunction() data.FuncStmt { return &ObStartFunction{} }
func (f *ObStartFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	callback, _ := ctx.GetIndexValue(0)
	flagsVal, _ := ctx.GetIndexValue(2)

	// PHP's ob_start 第 3 参数 flags 是输出处理器标志位。
	// PHP_OUTPUT_HANDLER_FLUSHABLE = 2（可刷新）
	flags := 7 // default: all flags enabled
	if flagsVal != nil {
		switch v := flagsVal.(type) {
		case *data.BoolValue:
			if v.Value {
				flags = 7
			} else {
				flags = 0
			}
		default:
			raw := strings.TrimSpace(flagsVal.AsString())
			if raw != "" {
				f, err := strconv.Atoi(raw)
				if err == nil {
					flags = f
				}
			}
		}
	}

	callbackName := ""
	if callback != nil {
		callbackName = callback.AsString()
	}

	obStack.pushWith(callbackName, flags&2 != 0)
	return data.NewBoolValue(true), nil
}
func (f *ObStartFunction) GetName() string            { return "ob_start" }
func (f *ObStartFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *ObStartFunction) GetIsStatic() bool          { return false }
func (f *ObStartFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
		node.NewParameter(nil, "chunk_size", 1, data.NewIntValue(0), nil),
		node.NewParameter(nil, "flags", 2, data.NewIntValue(0), nil),
	}
}
func (f *ObStartFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
		node.NewVariable(nil, "chunk_size", 1, nil),
		node.NewVariable(nil, "flags", 2, nil),
	}
}
func (f *ObStartFunction) GetReturnType() data.Types { return nil }

// ObFlushFunction 实现 ob_flush
type ObFlushFunction struct{}

func NewObFlushFunction() data.FuncStmt { return &ObFlushFunction{} }
func (f *ObFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, ok := obStack.flush()
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}
func (f *ObFlushFunction) GetName() string               { return "ob_flush" }
func (f *ObFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObFlushFunction) GetIsStatic() bool             { return false }
func (f *ObFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObFlushFunction) GetReturnType() data.Types     { return nil }

// ObGetCleanFunction 实现 ob_get_clean
type ObGetCleanFunction struct{}

func NewObGetCleanFunction() data.FuncStmt { return &ObGetCleanFunction{} }
func (f *ObGetCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	content := obStack.pop()
	return data.NewStringValue(content), nil
}
func (f *ObGetCleanFunction) GetName() string               { return "ob_get_clean" }
func (f *ObGetCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetCleanFunction) GetIsStatic() bool             { return false }
func (f *ObGetCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetCleanFunction) GetReturnType() data.Types     { return nil }

// ObGetContentsFunction 实现 ob_get_contents
type ObGetContentsFunction struct{}

func NewObGetContentsFunction() data.FuncStmt { return &ObGetContentsFunction{} }
func (f *ObGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	content := obStack.contents()
	return data.NewStringValue(content), nil
}
func (f *ObGetContentsFunction) GetName() string               { return "ob_get_contents" }
func (f *ObGetContentsFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetContentsFunction) GetIsStatic() bool             { return false }
func (f *ObGetContentsFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetContentsFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetContentsFunction) GetReturnType() data.Types     { return nil }

// ObGetStatusFunction 实现 ob_get_status
type ObGetStatusFunction struct{}

func NewObGetStatusFunction() data.FuncStmt { return &ObGetStatusFunction{} }
func (f *ObGetStatusFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arr := data.NewArrayValue(nil)
	return arr, nil
}
func (f *ObGetStatusFunction) GetName() string               { return "ob_get_status" }
func (f *ObGetStatusFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetStatusFunction) GetIsStatic() bool             { return false }
func (f *ObGetStatusFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetStatusFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetStatusFunction) GetReturnType() data.Types     { return nil }

// ObEndCleanFunction 实现 ob_end_clean
type ObEndCleanFunction struct{}

func NewObEndCleanFunction() data.FuncStmt { return &ObEndCleanFunction{} }
func (f *ObEndCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	obStack.pop()
	return data.NewBoolValue(true), nil
}
func (f *ObEndCleanFunction) GetName() string               { return "ob_end_clean" }
func (f *ObEndCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObEndCleanFunction) GetIsStatic() bool             { return false }
func (f *ObEndCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObEndCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObEndCleanFunction) GetReturnType() data.Types     { return nil }

// ObGetLevelFunction 实现 ob_get_level
type ObGetLevelFunction struct{}

func NewObGetLevelFunction() data.FuncStmt { return &ObGetLevelFunction{} }
func (f *ObGetLevelFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(len(obStack.buffers) - 1), nil
}
func (f *ObGetLevelFunction) GetName() string               { return "ob_get_level" }
func (f *ObGetLevelFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetLevelFunction) GetIsStatic() bool             { return false }
func (f *ObGetLevelFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetLevelFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetLevelFunction) GetReturnType() data.Types     { return nil }
