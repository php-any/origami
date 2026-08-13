package core

import (
	"github.com/php-any/origami/data"
)

func outputBufferHost(ctx data.Context) (data.OutputBufferHost, bool) {
	if ctx == nil {
		return nil, false
	}
	host, ok := ctx.GetVM().(data.OutputBufferHost)
	return host, ok
}

type ObStartFunction struct{}

func NewObStartFunction() data.FuncStmt { return &ObStartFunction{} }
func (f *ObStartFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	host.StartOutputBuffer()
	return data.NewBoolValue(true), nil
}
func (f *ObStartFunction) GetName() string               { return "ob_start" }
func (f *ObStartFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObStartFunction) GetIsStatic() bool             { return false }
func (f *ObStartFunction) GetParams() []data.GetValue    { return nil }
func (f *ObStartFunction) GetVariables() []data.Variable { return nil }
func (f *ObStartFunction) GetReturnType() data.Types     { return data.Bool{} }

type ObGetCleanFunction struct{}

func NewObGetCleanFunction() data.FuncStmt { return &ObGetCleanFunction{} }
func (f *ObGetCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	content, active := host.CleanOutputBuffer()
	if !active {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(content), nil
}
func (f *ObGetCleanFunction) GetName() string               { return "ob_get_clean" }
func (f *ObGetCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetCleanFunction) GetIsStatic() bool             { return false }
func (f *ObGetCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetCleanFunction) GetReturnType() data.Types     { return data.Mixed{} }

type ObGetContentsFunction struct{}

func NewObGetContentsFunction() data.FuncStmt { return &ObGetContentsFunction{} }
func (f *ObGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	content, active := host.OutputBufferContents()
	if !active {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(content), nil
}
func (f *ObGetContentsFunction) GetName() string               { return "ob_get_contents" }
func (f *ObGetContentsFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetContentsFunction) GetIsStatic() bool             { return false }
func (f *ObGetContentsFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetContentsFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetContentsFunction) GetReturnType() data.Types     { return data.Mixed{} }

type ObEndCleanFunction struct{}

func NewObEndCleanFunction() data.FuncStmt { return &ObEndCleanFunction{} }
func (f *ObEndCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	_, active := host.CleanOutputBuffer()
	return data.NewBoolValue(active), nil
}
func (f *ObEndCleanFunction) GetName() string               { return "ob_end_clean" }
func (f *ObEndCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObEndCleanFunction) GetIsStatic() bool             { return false }
func (f *ObEndCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObEndCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObEndCleanFunction) GetReturnType() data.Types     { return data.Bool{} }

type ObGetLevelFunction struct{}

func NewObGetLevelFunction() data.FuncStmt { return &ObGetLevelFunction{} }
func (f *ObGetLevelFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(host.OutputBufferLevel()), nil
}
func (f *ObGetLevelFunction) GetName() string               { return "ob_get_level" }
func (f *ObGetLevelFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetLevelFunction) GetIsStatic() bool             { return false }
func (f *ObGetLevelFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetLevelFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetLevelFunction) GetReturnType() data.Types     { return data.Int{} }
