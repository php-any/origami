package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
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

// ob_clean — 清空当前（栈顶）输出缓冲内容，但不结束缓冲。
type ObCleanFunction struct{}

func NewObCleanFunction() data.FuncStmt { return &ObCleanFunction{} }
func (f *ObCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	// PHP 8：成功清理（存在活动缓冲）返回 true，否则 false。
	return data.NewBoolValue(host.CleanCurrentBuffer()), nil
}
func (f *ObCleanFunction) GetName() string               { return "ob_clean" }
func (f *ObCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObCleanFunction) GetIsStatic() bool             { return false }
func (f *ObCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObCleanFunction) GetReturnType() data.Types     { return data.Bool{} }

// ob_end_flush — 输出并结束当前输出缓冲。
type ObEndFlushFunction struct{}

func NewObEndFlushFunction() data.FuncStmt { return &ObEndFlushFunction{} }
func (f *ObEndFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	_, active := host.FlushOutputBuffer()
	return data.NewBoolValue(active), nil
}
func (f *ObEndFlushFunction) GetName() string               { return "ob_end_flush" }
func (f *ObEndFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObEndFlushFunction) GetIsStatic() bool             { return false }
func (f *ObEndFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObEndFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObEndFlushFunction) GetReturnType() data.Types     { return data.Bool{} }

// ob_flush — 输出当前缓冲内容到上一层（或最终输出），但不结束缓冲。
type ObFlushFunction struct{}

func NewObFlushFunction() data.FuncStmt { return &ObFlushFunction{} }
func (f *ObFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	// ob_flush 语义：把当前缓冲内容输出到上一层，但保留当前缓冲层（清空后继续可写）。
	// 通过“弹出并写出到上层，再压入空缓冲”保持层级不变。
	if _, active := host.FlushOutputBuffer(); active {
		host.StartOutputBuffer()
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}
func (f *ObFlushFunction) GetName() string               { return "ob_flush" }
func (f *ObFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObFlushFunction) GetIsStatic() bool             { return false }
func (f *ObFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObFlushFunction) GetReturnType() data.Types     { return data.Bool{} }

// ob_get_flush — 获取当前缓冲内容并结束缓冲，同时把内容输出到上一层。
type ObGetFlushFunction struct{}

func NewObGetFlushFunction() data.FuncStmt { return &ObGetFlushFunction{} }
func (f *ObGetFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	content, active := host.FlushOutputBuffer()
	if !active {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(content), nil
}
func (f *ObGetFlushFunction) GetName() string               { return "ob_get_flush" }
func (f *ObGetFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetFlushFunction) GetIsStatic() bool             { return false }
func (f *ObGetFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetFlushFunction) GetReturnType() data.Types     { return data.Mixed{} }

// ob_get_length — 获取当前缓冲内容的字节长度。
type ObGetLengthFunction struct{}

func NewObGetLengthFunction() data.FuncStmt { return &ObGetLengthFunction{} }
func (f *ObGetLengthFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	length, active := host.OutputBufferLength()
	if !active {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(length), nil
}
func (f *ObGetLengthFunction) GetName() string               { return "ob_get_length" }
func (f *ObGetLengthFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetLengthFunction) GetIsStatic() bool             { return false }
func (f *ObGetLengthFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetLengthFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetLengthFunction) GetReturnType() data.Types     { return data.Mixed{} }

// ob_get_status — 返回输出缓冲状态数组。
type ObGetStatusFunction struct{}

func NewObGetStatusFunction() data.FuncStmt { return &ObGetStatusFunction{} }
func (f *ObGetStatusFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}
	full := false
	if v, has := ctx.GetIndexValue(0); has && v != nil {
		if asBool, ok := v.(data.AsBool); ok {
			if b, err := asBool.AsBool(); err == nil {
				full = b
			}
		}
	}
	status := host.OutputBufferStatus(full)
	if len(status) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}
	if !full {
		// 仅最顶层：返回单个关联数组。
		return statusToArray(status[0]), nil
	}
	// 全部层：返回按层级为键的数组。
	list := make([]data.Value, len(status))
	for i, st := range status {
		list[i] = &data.ArrayValue{
			List: []*data.ZVal{
				data.NewNamedZVal("level", data.NewIntValue(st.Level)),
				data.NewNamedZVal("type", data.NewIntValue(st.Type)),
				data.NewNamedZVal("flags", data.NewIntValue(st.Flags)),
				data.NewNamedZVal("chunk_size", data.NewIntValue(st.ChunkSize)),
				data.NewNamedZVal("buffer_size", data.NewIntValue(st.BufferSize)),
				data.NewNamedZVal("name", data.NewStringValue(st.Name)),
			},
		}
	}
	return data.NewArrayValue(list), nil
}
func (f *ObGetStatusFunction) GetName() string            { return "ob_get_status" }
func (f *ObGetStatusFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *ObGetStatusFunction) GetIsStatic() bool          { return false }
func (f *ObGetStatusFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "full", 0, nil, nil),
	}
}
func (f *ObGetStatusFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "full", 0, data.NewBaseType("bool")),
	}
}
func (f *ObGetStatusFunction) GetReturnType() data.Types { return data.Arrays{} }

// statusToArray 将单个缓冲层状态转换为 PHP 关联数组。
func statusToArray(st data.OutputBufferStatusInfo) data.Value {
	return &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("level", data.NewIntValue(st.Level)),
			data.NewNamedZVal("type", data.NewIntValue(st.Type)),
			data.NewNamedZVal("flags", data.NewIntValue(st.Flags)),
			data.NewNamedZVal("chunk_size", data.NewIntValue(st.ChunkSize)),
			data.NewNamedZVal("buffer_size", data.NewIntValue(st.BufferSize)),
			data.NewNamedZVal("name", data.NewStringValue(st.Name)),
		},
	}
}

// ob_list_handlers — 返回所有激活输出处理器的名称列表。
type ObListHandlersFunction struct{}

func NewObListHandlersFunction() data.FuncStmt { return &ObListHandlersFunction{} }
func (f *ObListHandlersFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}
	handlers := host.ListOutputHandlers()
	list := make([]data.Value, 0, len(handlers))
	for _, h := range handlers {
		list = append(list, data.NewStringValue(h))
	}
	return data.NewArrayValue(list), nil
}
func (f *ObListHandlersFunction) GetName() string               { return "ob_list_handlers" }
func (f *ObListHandlersFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObListHandlersFunction) GetIsStatic() bool             { return false }
func (f *ObListHandlersFunction) GetParams() []data.GetValue    { return nil }
func (f *ObListHandlersFunction) GetVariables() []data.Variable { return nil }
func (f *ObListHandlersFunction) GetReturnType() data.Types     { return data.Arrays{} }

// ob_implicit_flush — 打开/关闭隐式刷新。
type ObImplicitFlushFunction struct{}

func NewObImplicitFlushFunction() data.FuncStmt { return &ObImplicitFlushFunction{} }
func (f *ObImplicitFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewNullValue(), nil
	}
	on := true
	if v, has := ctx.GetIndexValue(0); has && v != nil {
		if asBool, ok := v.(data.AsBool); ok {
			if b, err := asBool.AsBool(); err == nil {
				on = b
			}
		}
	}
	host.SetImplicitFlush(on)
	return data.NewNullValue(), nil
}
func (f *ObImplicitFlushFunction) GetName() string            { return "ob_implicit_flush" }
func (f *ObImplicitFlushFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *ObImplicitFlushFunction) GetIsStatic() bool          { return false }
func (f *ObImplicitFlushFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "enable", 0, nil, nil),
	}
}
func (f *ObImplicitFlushFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "enable", 0, data.NewBaseType("bool")),
	}
}
func (f *ObImplicitFlushFunction) GetReturnType() data.Types { return data.NewBaseType("null") }
