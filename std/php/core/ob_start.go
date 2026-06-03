package core

import (
	"strings"

	"github.com/php-any/origami/data"
)

// outputBufferStack 输出缓冲栈
type outputBufferStack struct {
	buffers []*strings.Builder
}

var obStack = &outputBufferStack{buffers: []*strings.Builder{{}}}

func (s *outputBufferStack) push() {
	s.buffers = append(s.buffers, &strings.Builder{})
	s.syncWriter()
}

func (s *outputBufferStack) pop() string {
	if len(s.buffers) <= 1 {
		return ""
	}
	last := s.buffers[len(s.buffers)-1].String()
	s.buffers = s.buffers[:len(s.buffers)-1]
	s.syncWriter()
	return last
}

func (s *outputBufferStack) contents() string {
	if len(s.buffers) == 0 {
		return ""
	}
	return s.buffers[len(s.buffers)-1].String()
}

func (s *outputBufferStack) syncWriter() {
	if len(s.buffers) <= 1 {
		data.WriteOutput = data.DefaultOutputWriter
		return
	}
	buf := s.buffers[len(s.buffers)-1]
	data.WriteOutput = func(str string) {
		buf.WriteString(str)
	}
}

// FlushAllBuffers 刷新所有输出缓冲区（脚本结束时调用）
func FlushAllBuffers() {
	for len(obStack.buffers) > 1 {
		content := obStack.buffers[len(obStack.buffers)-1].String()
		obStack.buffers = obStack.buffers[:len(obStack.buffers)-1]
		if content != "" {
			// 写入到父缓冲区（当前栈顶）
			if len(obStack.buffers) > 1 {
				obStack.buffers[len(obStack.buffers)-1].WriteString(content)
			} else {
				// 如果没有父缓冲区，写入到 stdout
				data.DefaultOutputWriter(content)
			}
		}
	}
	obStack.syncWriter()
}

func init() {
	data.FlushAllBuffersFn = FlushAllBuffers
}

// ObStartFunction 实现 ob_start
type ObStartFunction struct{}

func NewObStartFunction() data.FuncStmt { return &ObStartFunction{} }
func (f *ObStartFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	obStack.push()
	return data.NewBoolValue(true), nil
}
func (f *ObStartFunction) GetName() string                    { return "ob_start" }
func (f *ObStartFunction) GetModifier() data.Modifier         { return data.ModifierPublic }
func (f *ObStartFunction) GetIsStatic() bool                  { return false }
func (f *ObStartFunction) GetParams() []data.GetValue         { return nil }
func (f *ObStartFunction) GetVariables() []data.Variable      { return nil }
func (f *ObStartFunction) GetReturnType() data.Types          { return nil }

// ObGetCleanFunction 实现 ob_get_clean
type ObGetCleanFunction struct{}

func NewObGetCleanFunction() data.FuncStmt { return &ObGetCleanFunction{} }
func (f *ObGetCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	content := obStack.pop()
	return data.NewStringValue(content), nil
}
func (f *ObGetCleanFunction) GetName() string                    { return "ob_get_clean" }
func (f *ObGetCleanFunction) GetModifier() data.Modifier         { return data.ModifierPublic }
func (f *ObGetCleanFunction) GetIsStatic() bool                  { return false }
func (f *ObGetCleanFunction) GetParams() []data.GetValue         { return nil }
func (f *ObGetCleanFunction) GetVariables() []data.Variable      { return nil }
func (f *ObGetCleanFunction) GetReturnType() data.Types          { return nil }

// ObGetContentsFunction 实现 ob_get_contents
type ObGetContentsFunction struct{}

func NewObGetContentsFunction() data.FuncStmt { return &ObGetContentsFunction{} }
func (f *ObGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	content := obStack.contents()
	return data.NewStringValue(content), nil
}
func (f *ObGetContentsFunction) GetName() string                    { return "ob_get_contents" }
func (f *ObGetContentsFunction) GetModifier() data.Modifier         { return data.ModifierPublic }
func (f *ObGetContentsFunction) GetIsStatic() bool                  { return false }
func (f *ObGetContentsFunction) GetParams() []data.GetValue         { return nil }
func (f *ObGetContentsFunction) GetVariables() []data.Variable      { return nil }
func (f *ObGetContentsFunction) GetReturnType() data.Types          { return nil }

// ObEndCleanFunction 实现 ob_end_clean
type ObEndCleanFunction struct{}

func NewObEndCleanFunction() data.FuncStmt { return &ObEndCleanFunction{} }
func (f *ObEndCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	obStack.pop()
	return data.NewBoolValue(true), nil
}
func (f *ObEndCleanFunction) GetName() string                    { return "ob_end_clean" }
func (f *ObEndCleanFunction) GetModifier() data.Modifier         { return data.ModifierPublic }
func (f *ObEndCleanFunction) GetIsStatic() bool                  { return false }
func (f *ObEndCleanFunction) GetParams() []data.GetValue         { return nil }
func (f *ObEndCleanFunction) GetVariables() []data.Variable      { return nil }
func (f *ObEndCleanFunction) GetReturnType() data.Types          { return nil }

// ObGetLevelFunction 实现 ob_get_level
type ObGetLevelFunction struct{}

func NewObGetLevelFunction() data.FuncStmt { return &ObGetLevelFunction{} }
func (f *ObGetLevelFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(len(obStack.buffers) - 1), nil
}
func (f *ObGetLevelFunction) GetName() string                    { return "ob_get_level" }
func (f *ObGetLevelFunction) GetModifier() data.Modifier         { return data.ModifierPublic }
func (f *ObGetLevelFunction) GetIsStatic() bool                  { return false }
func (f *ObGetLevelFunction) GetParams() []data.GetValue         { return nil }
func (f *ObGetLevelFunction) GetVariables() []data.Variable      { return nil }
func (f *ObGetLevelFunction) GetReturnType() data.Types          { return nil }
