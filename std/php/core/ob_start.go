package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// 对齐 PHP 8.4 ob_*：ob_start($callback, $chunk_size, $flags)、handler 改写、
// CLEANABLE/FLUSHABLE/REMOVABLE、ob_flush 原位冒泡。回调在独立 Context 中调用，
// 不捕获 ob_start 的帧（避免 pooled Context 悬挂）。chunk_size 热路径不加追踪。

func outputBufferHost(ctx data.Context) (data.OutputBufferHost, bool) {
	if ctx == nil {
		return nil, false
	}
	if host, ok := ctx.(data.OutputBufferHost); ok {
		return host, true
	}
	host, ok := ctx.GetVM().(data.OutputBufferHost)
	return host, ok
}

func finishOb(host data.OutputBufferHost, ret data.GetValue) (data.GetValue, data.Control) {
	if host != nil {
		if ctl := host.TakeOutputControl(); ctl != nil {
			return nil, ctl
		}
	}
	return ret, nil
}

func isNullish(v data.GetValue) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func asIntArg(v data.GetValue, def int) int {
	if v == nil {
		return def
	}
	if asInt, ok := v.(data.AsInt); ok {
		if n, err := asInt.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func outputHandlerName(cb data.Value) string {
	switch c := cb.(type) {
	case *data.StringValue:
		if s := c.AsString(); s != "" {
			return s
		}
	case *data.FuncValue:
		if c.Value != nil {
			if n := c.Value.GetName(); n != "" {
				return n
			}
		}
		return "{closure}"
	case *data.BoundFuncValue:
		if c.Value != nil {
			if n := c.Value.GetName(); n != "" {
				return n
			}
		}
		return "{closure}"
	}
	return "user output handler"
}

func outputHandlerArgc(cb data.Value) int {
	var fn data.FuncStmt
	switch c := cb.(type) {
	case *data.FuncValue:
		fn = c.Value
	case *data.BoundFuncValue:
		fn = c.Value
	default:
		return 2
	}
	if fn == nil {
		return 2
	}
	n := len(fn.GetParams())
	if n <= 0 {
		return 1
	}
	return n
}

func wrapOutputHandler(vm data.VM, cb data.Value) data.OutputBufferHandler {
	return func(buffer string, phase int) (string, bool, data.Control) {
		return invokeOutputHandler(vm, cb, buffer, phase)
	}
}

func invokeOutputHandler(vm data.VM, cb data.Value, buffer string, phase int) (string, bool, data.Control) {
	if vm == nil || cb == nil {
		return buffer, true, nil
	}
	cu := NewCallUserFuncFunction()
	ctx := vm.CreateContext(cu.GetVariables())
	args := []data.Value{data.NewStringValue(buffer)}
	if outputHandlerArgc(cb) != 1 {
		args = append(args, data.NewIntValue(phase))
	}
	ctx.SetIndexZVal(0, data.NewZVal(cb))
	ctx.SetIndexZVal(1, data.NewZVal(data.NewArrayValue(args)))
	ret, ctl := cu.Call(ctx)
	if ctl != nil {
		return buffer, false, ctl
	}
	if ret == nil {
		return "", true, nil
	}
	if bv, ok := ret.(*data.BoolValue); ok {
		if b, err := bv.AsBool(); err == nil && !b {
			return buffer, true, nil
		}
	}
	if val, ok := ret.(data.Value); ok {
		return val.AsString(), true, nil
	}
	return buffer, true, nil
}

// ob_start(?callable $callback = null, int $chunk_size = 0, int $flags = PHP_OUTPUT_HANDLER_STDFLAGS): bool
type ObStartFunction struct{}

func NewObStartFunction() data.FuncStmt { return &ObStartFunction{} }
func (f *ObStartFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	spec := data.OutputBufferStartSpec{
		Flags: data.PHPOutputHandlerStdFlags,
		Name:  "default output handler",
		Type:  data.PHPOutputHandlerInternal,
	}
	if v, has := ctx.GetIndexValue(0); has && !isNullish(v) {
		if cb, ok := v.(data.Value); ok {
			spec.Handler = wrapOutputHandler(ctx.GetVM(), cb)
			spec.Name = outputHandlerName(cb)
			spec.Type = data.PHPOutputHandlerUser
		}
	}
	if v, has := ctx.GetIndexValue(1); has {
		spec.ChunkSize = asIntArg(v, 0)
	}
	if v, has := ctx.GetIndexValue(2); has && !isNullish(v) {
		spec.Flags = asIntArg(v, data.PHPOutputHandlerStdFlags)
	}
	return data.NewBoolValue(host.StartOutputBufferSpec(spec)), nil
}
func (f *ObStartFunction) GetName() string            { return "ob_start" }
func (f *ObStartFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *ObStartFunction) GetIsStatic() bool          { return false }
func (f *ObStartFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, data.NewNullValue(), data.Mixed{}),
		node.NewParameter(nil, "chunk_size", 1, data.NewIntValue(0), data.Int{}),
		node.NewParameter(nil, "flags", 2, data.NewIntValue(data.PHPOutputHandlerStdFlags), data.Int{}),
	}
}
func (f *ObStartFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
		node.NewVariable(nil, "chunk_size", 1, data.Int{}),
		node.NewVariable(nil, "flags", 2, data.Int{}),
	}
}
func (f *ObStartFunction) GetReturnType() data.Types { return data.Bool{} }

type ObGetCleanFunction struct{}

func NewObGetCleanFunction() data.FuncStmt { return &ObGetCleanFunction{} }
func (f *ObGetCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	content, active := host.CleanOutputBuffer()
	if !active {
		return finishOb(host, data.NewBoolValue(false))
	}
	return finishOb(host, data.NewStringValue(content))
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
	return finishOb(host, data.NewBoolValue(active))
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

type ObCleanFunction struct{}

func NewObCleanFunction() data.FuncStmt { return &ObCleanFunction{} }
func (f *ObCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return finishOb(host, data.NewBoolValue(host.CleanCurrentBuffer()))
}
func (f *ObCleanFunction) GetName() string               { return "ob_clean" }
func (f *ObCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObCleanFunction) GetIsStatic() bool             { return false }
func (f *ObCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObCleanFunction) GetReturnType() data.Types     { return data.Bool{} }

type ObEndFlushFunction struct{}

func NewObEndFlushFunction() data.FuncStmt { return &ObEndFlushFunction{} }
func (f *ObEndFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	_, active := host.FlushOutputBuffer()
	return finishOb(host, data.NewBoolValue(active))
}
func (f *ObEndFlushFunction) GetName() string               { return "ob_end_flush" }
func (f *ObEndFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObEndFlushFunction) GetIsStatic() bool             { return false }
func (f *ObEndFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObEndFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObEndFlushFunction) GetReturnType() data.Types     { return data.Bool{} }

type ObFlushFunction struct{}

func NewObFlushFunction() data.FuncStmt { return &ObFlushFunction{} }
func (f *ObFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	_, active := host.FlushCurrentBuffer()
	return finishOb(host, data.NewBoolValue(active))
}
func (f *ObFlushFunction) GetName() string               { return "ob_flush" }
func (f *ObFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObFlushFunction) GetIsStatic() bool             { return false }
func (f *ObFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObFlushFunction) GetReturnType() data.Types     { return data.Bool{} }

type ObGetFlushFunction struct{}

func NewObGetFlushFunction() data.FuncStmt { return &ObGetFlushFunction{} }
func (f *ObGetFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, ok := outputBufferHost(ctx)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	content, active := host.FlushOutputBuffer()
	if !active {
		return finishOb(host, data.NewBoolValue(false))
	}
	return finishOb(host, data.NewStringValue(content))
}
func (f *ObGetFlushFunction) GetName() string               { return "ob_get_flush" }
func (f *ObGetFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetFlushFunction) GetIsStatic() bool             { return false }
func (f *ObGetFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetFlushFunction) GetReturnType() data.Types     { return data.Mixed{} }

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
		return statusToArray(status[0]), nil
	}
	list := make([]data.Value, len(status))
	for i, st := range status {
		list[i] = statusToArray(st)
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

func statusToArray(st data.OutputBufferStatusInfo) data.Value {
	return &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("level", data.NewIntValue(st.Level)),
			data.NewNamedZVal("type", data.NewIntValue(st.Type)),
			data.NewNamedZVal("flags", data.NewIntValue(st.Flags)),
			data.NewNamedZVal("chunk_size", data.NewIntValue(st.ChunkSize)),
			data.NewNamedZVal("buffer_size", data.NewIntValue(st.BufferSize)),
			data.NewNamedZVal("buffer_used", data.NewIntValue(st.BufferUsed)),
			data.NewNamedZVal("name", data.NewStringValue(st.Name)),
		},
	}
}

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

// FlushFunction 对齐 PHP flush(): 刷新 SAPI，不弹出 ob 栈。
type FlushFunction struct{}

func NewFlushFunction() data.FuncStmt { return &FlushFunction{} }
func (f *FlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if host, ok := outputBufferHost(ctx); ok {
		host.FlushSAPI()
	}
	return data.NewBoolValue(true), nil
}
func (f *FlushFunction) GetName() string               { return "flush" }
func (f *FlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *FlushFunction) GetIsStatic() bool             { return false }
func (f *FlushFunction) GetParams() []data.GetValue    { return nil }
func (f *FlushFunction) GetVariables() []data.Variable { return nil }
func (f *FlushFunction) GetReturnType() data.Types     { return data.Bool{} }
