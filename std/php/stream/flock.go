package stream

import (
	"os"
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// PHP flock 操作常量
const (
	LockSH = 1 // LOCK_SH
	LockEX = 2 // LOCK_EX
	LockUN = 3 // LOCK_UN
	LockNB = 4 // LOCK_NB（可与 SH/EX 按位或）
)

// FlockFunction 实现 flock($handle, $operation, &$would_block = null): bool
type FlockFunction struct{}

func NewFlockFunction() data.FuncStmt {
	return &FlockFunction{}
}

func (f *FlockFunction) GetName() string { return "flock" }

func (f *FlockFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "operation", 1, nil, data.NewBaseType("int")),
		node.NewParameterReference(nil, "would_block", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *FlockFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, nil),
		node.NewVariable(nil, "operation", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "would_block", 2, nil),
	}
}

func (f *FlockFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	streamValue, _ := ctx.GetIndexValue(0)
	file := fileFromStreamValue(streamValue)
	if file == nil {
		return data.NewBoolValue(false), nil
	}

	op := LockEX
	if opVal, ok := ctx.GetIndexValue(1); ok && opVal != nil {
		if ai, ok := opVal.(interface{ AsInt() (int, error) }); ok {
			if v, err := ai.AsInt(); err == nil {
				op = v
			}
		}
	}

	how, ok := phpFlockHow(op)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	err := syscall.Flock(int(file.Fd()), how)
	if err != nil {
		if op&LockNB != 0 && (err == syscall.EWOULDBLOCK || err == syscall.EAGAIN) {
			_ = ctx.SetVariableValue(node.NewVariable(nil, "would_block", 2, nil), data.NewIntValue(1))
			return data.NewBoolValue(false), nil
		}
		return data.NewBoolValue(false), nil
	}
	if op&LockNB != 0 {
		_ = ctx.SetVariableValue(node.NewVariable(nil, "would_block", 2, nil), data.NewIntValue(0))
	}
	return data.NewBoolValue(true), nil
}

func phpFlockHow(op int) (int, bool) {
	nb := op&LockNB != 0
	base := op &^ LockNB
	var how int
	switch base {
	case LockSH:
		how = syscall.LOCK_SH
	case LockEX:
		how = syscall.LOCK_EX
	case LockUN:
		how = syscall.LOCK_UN
	default:
		return 0, false
	}
	if nb {
		how |= syscall.LOCK_NB
	}
	return how, true
}

func fileFromStreamValue(streamValue data.Value) *os.File {
	res, ok := streamValue.(*core.ResourceValue)
	if !ok {
		return nil
	}
	resource := res.GetResource()
	if info, ok := resource.(*StreamInfo); ok && info != nil && !info.IsClosed() {
		return info.File
	}
	return nil
}
