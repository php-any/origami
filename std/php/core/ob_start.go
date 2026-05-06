package core

import (
	"github.com/php-any/origami/data"
)

// outputBufferStack is a simple per-VM output buffer stack.
// We store it as a VM-level constant map to keep state across calls.
var obStack = &outputBufferStack{buffers: []string{""}}

type outputBufferStack struct {
	buffers []string
}

func (s *outputBufferStack) push() {
	s.buffers = append(s.buffers, "")
}

func (s *outputBufferStack) pop() string {
	if len(s.buffers) <= 1 {
		return ""
	}
	last := s.buffers[len(s.buffers)-1]
	s.buffers = s.buffers[:len(s.buffers)-1]
	return last
}

func (s *outputBufferStack) current() *string {
	if len(s.buffers) == 0 {
		s.buffers = []string{""}
	}
	return &s.buffers[len(s.buffers)-1]
}

func (s *outputBufferStack) write(data string) {
	cur := s.current()
	*cur += data
}

// ObStartFunction 实现 ob_start 函数
type ObStartFunction struct{}

func NewObStartFunction() data.FuncStmt {
	return &ObStartFunction{}
}

func (f *ObStartFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	obStack.push()
	return data.NewBoolValue(true), nil
}

func (f *ObStartFunction) GetName() string {
	return "ob_start"
}

func (f *ObStartFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *ObStartFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// ObGetCleanFunction 实现 ob_get_clean 函数
type ObGetCleanFunction struct{}

func NewObGetCleanFunction() data.FuncStmt {
	return &ObGetCleanFunction{}
}

func (f *ObGetCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	content := obStack.pop()
	return data.NewStringValue(content), nil
}

func (f *ObGetCleanFunction) GetName() string {
	return "ob_get_clean"
}

func (f *ObGetCleanFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *ObGetCleanFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// ObGetContentsFunction 实现 ob_get_contents 函数
type ObGetContentsFunction struct{}

func NewObGetContentsFunction() data.FuncStmt {
	return &ObGetContentsFunction{}
}

func (f *ObGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	content := *obStack.current()
	return data.NewStringValue(content), nil
}

func (f *ObGetContentsFunction) GetName() string {
	return "ob_get_contents"
}

func (f *ObGetContentsFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *ObGetContentsFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// ObEndCleanFunction 实现 ob_end_clean 函数
type ObEndCleanFunction struct{}

func NewObEndCleanFunction() data.FuncStmt {
	return &ObEndCleanFunction{}
}

func (f *ObEndCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	obStack.pop()
	return data.NewBoolValue(true), nil
}

func (f *ObEndCleanFunction) GetName() string {
	return "ob_end_clean"
}

func (f *ObEndCleanFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *ObEndCleanFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// ObGetLevelFunction 实现 ob_get_level 函数
type ObGetLevelFunction struct{}

func NewObGetLevelFunction() data.FuncStmt {
	return &ObGetLevelFunction{}
}

func (f *ObGetLevelFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(len(obStack.buffers) - 1), nil
}

func (f *ObGetLevelFunction) GetName() string {
	return "ob_get_level"
}

func (f *ObGetLevelFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *ObGetLevelFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
