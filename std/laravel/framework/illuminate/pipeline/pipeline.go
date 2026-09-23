package pipeline

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/illuminate/conditionable"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const pipelineClassName = "Illuminate\\Pipeline\\Pipeline"

type pipeState struct {
	passable          data.Value
	pipes             []data.Value
	method            string
	container         *data.ClassValue
	finallyCb         data.Value
	withinTransaction data.Value
}

var (
	pipeStates sync.Map
	pipeNextID atomic.Int64
)

func stateOf(cv *data.ClassValue) *pipeState {
	if cv == nil {
		return &pipeState{method: "handle", withinTransaction: data.NewBoolValue(false)}
	}
	if v, ctl := cv.GetProperty("__origami_pipeline_id"); ctl == nil && v != nil {
		if iv, ok := kit.Unwrap(v).(*data.IntValue); ok && iv.Value > 0 {
			if s, ok := pipeStates.Load(int64(iv.Value)); ok {
				return s.(*pipeState)
			}
		}
	}
	id := pipeNextID.Add(1)
	s := &pipeState{method: "handle", withinTransaction: data.NewBoolValue(false)}
	pipeStates.Store(id, s)
	_ = cv.SetProperty("__origami_pipeline_id", data.NewIntValue(int(id)))
	return s
}

type PipelineClass struct {
	node.Node
	methods map[string]data.Method
}

func NewPipelineClass() data.ClassStmt {
	c := &PipelineClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *PipelineClass) GetName() string    { return pipelineClassName }
func (c *PipelineClass) GetExtend() *string { return nil }
func (c *PipelineClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Pipeline\\Pipeline"}
}
func (c *PipelineClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "container", "passable", "pipes", "method", "finally", "withinTransaction", "__origami_pipeline_id":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *PipelineClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 7)
	for _, n := range []string{"container", "passable", "pipes", "method", "finally", "withinTransaction", "__origami_pipeline_id"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *PipelineClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *PipelineClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *PipelineClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *PipelineClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *PipelineClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *PipelineClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"container"}, 0, pipeConstruct)
	c.methods["send"] = kit.InstanceMethod("send", []string{"passable"}, pipeSend)
	c.methods["through"] = kit.InstanceMethod("through", []string{"pipes"}, pipeThrough)
	c.methods["pipe"] = kit.InstanceMethod("pipe", []string{"pipes"}, pipePipe)
	c.methods["via"] = kit.InstanceMethod("via", []string{"method"}, pipeVia)
	c.methods["then"] = kit.InstanceMethod("then", []string{"destination"}, pipeThen)
	c.methods["thenreturn"] = kit.InstanceMethod("thenReturn", nil, pipeThenReturn)
	c.methods["finally"] = kit.InstanceMethod("finally", []string{"callback"}, pipeFinally)
	c.methods["setcontainer"] = kit.InstanceMethod("setContainer", []string{"container"}, pipeSetContainer)
	c.methods["withintransaction"] = kit.InstanceMethodOpt("withinTransaction", []string{"withinTransaction"}, 0, pipeWithinTransaction)
	kit.RegisterMacroable(c.methods, pipelineClassName)
	kit.RegisterConditionable(c.methods, conditionable.NewWhenProxy)
}

func pipeRecv(ctx data.Context) (*data.ClassValue, *pipeState, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return nil, nil, data.NewErrorThrow(nil, fmt.Errorf("Pipeline missing $this"))
	}
	return cv, stateOf(cv), nil
}

func pipeConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if c := kit.Arg(ctx, 0); c != nil {
		if ccv, ok := kit.Unwrap(c).(*data.ClassValue); ok {
			st.container = ccv
			_ = cv.SetProperty("container", c)
		}
	}
	return cv, nil
}

func pipeSend(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st.passable = kit.Arg(ctx, 0)
	_ = cv.SetProperty("passable", st.passable)
	return cv, nil
}

func pipeThrough(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st.pipes = collectPipes(ctx, 0)
	syncPipesProperty(cv, st)
	return cv, nil
}

func pipePipe(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st.pipes = append(st.pipes, collectPipes(ctx, 0)...)
	syncPipesProperty(cv, st)
	return cv, nil
}

func collectPipes(ctx data.Context, start int) []data.Value {
	args := ctx.GetFlatCallArgs()
	if args != nil && len(args) > start {
		first := args[start]
		if av, ok := kit.Unwrap(first).(*data.ArrayValue); ok && av != nil {
			return av.ToValueList()
		}
		out := make([]data.Value, 0, len(args)-start)
		for i := start; i < len(args); i++ {
			out = append(out, args[i])
		}
		return out
	}
	first := kit.Arg(ctx, start)
	if av, ok := kit.Unwrap(first).(*data.ArrayValue); ok && av != nil {
		return av.ToValueList()
	}
	if first != nil {
		return []data.Value{first}
	}
	return nil
}

func syncPipesProperty(cv *data.ClassValue, st *pipeState) {
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	for i, p := range st.pipes {
		arr.SetIntKey(i, p)
	}
	_ = cv.SetProperty("pipes", arr)
}

func pipeVia(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	m := kit.Arg(ctx, 0)
	if m != nil {
		st.method = m.AsString()
		_ = cv.SetProperty("method", data.NewStringValue(st.method))
	}
	return cv, nil
}

func pipeThenReturn(ctx data.Context) (data.GetValue, data.Control) {
	cv, _, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return pipeThenWithDest(ctx, cv, data.NewFuncValue(&pipeIdentityFunc{}))
}

type pipeIdentityFunc struct{}

func (pipeIdentityFunc) GetName() string                 { return "pipeline_then_return" }
func (pipeIdentityFunc) GetParams() []data.GetValue      { return []data.GetValue{node.NewParameter(nil, "passable", 0, nil, nil)} }
func (pipeIdentityFunc) GetVariables() []data.Variable { return []data.Variable{node.NewVariable(nil, "passable", 0, nil)} }
func (pipeIdentityFunc) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return v, nil
}

func pipeThen(ctx data.Context) (data.GetValue, data.Control) {
	cv, _, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	dest := kit.Arg(ctx, 0)
	return pipeThenWithDest(ctx, cv, dest)
}

func pipeThenWithDest(ctx data.Context, cv *data.ClassValue, dest data.Value) (data.GetValue, data.Control) {
	st := stateOf(cv)
	stack := func(passable data.Value) (data.GetValue, data.Control) {
		return kit.Call(ctx, dest, passable)
	}
	for i := len(st.pipes) - 1; i >= 0; i-- {
		pipe := st.pipes[i]
		next := stack
		stack = func(passable data.Value) (data.GetValue, data.Control) {
			return runPipe(ctx, cv, st, pipe, passable, next)
		}
	}
	var result data.GetValue
	var runCtl data.Control
	if shouldWrapTransaction(st) {
		result, runCtl = runWithinTransaction(ctx, st, func() (data.GetValue, data.Control) {
			return stack(st.passable)
		})
	} else {
		result, runCtl = stack(st.passable)
	}
	if st.finallyCb != nil {
		_, _ = kit.Call(ctx, st.finallyCb, st.passable)
	}
	return result, runCtl
}

func shouldWrapTransaction(st *pipeState) bool {
	if st.withinTransaction == nil {
		return false
	}
	if bv, ok := st.withinTransaction.(*data.BoolValue); ok && !bv.Value {
		return false
	}
	return !kit.IsNull(st.withinTransaction)
}

func runWithinTransaction(ctx data.Context, st *pipeState, fn func() (data.GetValue, data.Control)) (data.GetValue, data.Control) {
	if st.container == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("A container instance has not been passed to the Pipeline."))
	}
	dbRet, ctl := kit.CallInstanceMethod(ctx, st.container, "make", data.NewStringValue("db"))
	if ctl != nil {
		return nil, ctl
	}
	db := kit.Unwrap(dbRet.(data.Value))
	connName := st.withinTransaction
	if bv, ok := st.withinTransaction.(*data.BoolValue); ok && bv.Value {
		connName = data.NewNullValue()
	}
	dbCV, ok := db.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Pipeline db binding is not a container object"))
	}
	connRet, ctl := kit.CallInstanceMethod(ctx, dbCV, "connection", connName)
	if ctl != nil {
		return nil, ctl
	}
	connCV, ok := kit.Unwrap(connRet.(data.Value)).(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Pipeline db connection is not an object"))
	}
	txFn := data.NewFuncValue(&pipeTxWrapper{inner: fn})
	return kit.CallInstanceMethod(ctx, connCV, "transaction", txFn)
}

type pipeTxWrapper struct {
	inner func() (data.GetValue, data.Control)
}

func (w *pipeTxWrapper) GetName() string                 { return "pipeline_tx" }
func (w *pipeTxWrapper) GetParams() []data.GetValue      { return nil }
func (w *pipeTxWrapper) GetVariables() []data.Variable { return nil }
func (w *pipeTxWrapper) Call(data.Context) (data.GetValue, data.Control) {
	return w.inner()
}

func runPipe(ctx data.Context, cv *data.ClassValue, st *pipeState, pipe, passable data.Value, stack func(data.Value) (data.GetValue, data.Control)) (data.GetValue, data.Control) {
	pipe = kit.Unwrap(pipe)
	stackFn := data.NewFuncValue(&pipeStackFunc{fn: stack})
	if isCallable(pipe) {
		return kit.Call(ctx, pipe, passable, stackFn)
	}
	if _, ok := pipe.(*data.ClassValue); ok {
		return callPipeHandle(ctx, st, pipe, passable, stackFn)
	}
	if sv, ok := pipe.(*data.StringValue); ok {
		name, extra := parsePipeString(sv.AsString())
		if st.container == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("A container instance has not been passed to the Pipeline."))
		}
		made, ctl := kit.CallInstanceMethod(ctx, st.container, "make", data.NewStringValue(name))
		if ctl != nil {
			return nil, ctl
		}
		pipe = kit.Unwrap(made.(data.Value))
		args := append([]data.Value{passable, stackFn}, extra...)
		return callPipeHandleArgs(ctx, st, pipe, args...)
	}
	args := []data.Value{passable, stackFn}
	return callPipeHandleArgs(ctx, st, pipe, args...)
}

type pipeStackFunc struct {
	fn func(data.Value) (data.GetValue, data.Control)
}

func (p *pipeStackFunc) GetName() string                 { return "pipeline_stack" }
func (p *pipeStackFunc) GetParams() []data.GetValue      { return []data.GetValue{node.NewParameter(nil, "passable", 0, nil, nil)} }
func (p *pipeStackFunc) GetVariables() []data.Variable { return []data.Variable{node.NewVariable(nil, "passable", 0, nil)} }
func (p *pipeStackFunc) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return p.fn(v)
}

func callPipeHandle(ctx data.Context, st *pipeState, pipe, passable, stack data.Value) (data.GetValue, data.Control) {
	return callPipeHandleArgs(ctx, st, pipe, passable, stack)
}

func callPipeHandleArgs(ctx data.Context, st *pipeState, pipe data.Value, args ...data.Value) (data.GetValue, data.Control) {
	pipe = kit.Unwrap(pipe)
	cv, ok := pipe.(*data.ClassValue)
	if !ok {
		return kit.Call(ctx, pipe, args...)
	}
	if _, ok := cv.GetMethod(strings.ToLower(st.method)); ok {
		return kit.CallInstanceMethod(ctx, cv, st.method, args...)
	}
	return kit.Call(ctx, pipe, args...)
}

func isCallable(v data.Value) bool {
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return true
	default:
		return false
	}
}

func parsePipeString(pipe string) (string, []data.Value) {
	parts := strings.SplitN(pipe, ":", 2)
	name := parts[0]
	var params []data.Value
	if len(parts) == 2 && parts[1] != "" {
		for _, p := range strings.Split(parts[1], ",") {
			params = append(params, data.NewStringValue(strings.TrimSpace(p)))
		}
	}
	return name, params
}

func pipeFinally(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st.finallyCb = kit.Arg(ctx, 0)
	_ = cv.SetProperty("finally", st.finallyCb)
	return cv, nil
}

func pipeSetContainer(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	c := kit.Arg(ctx, 0)
	if ccv, ok := kit.Unwrap(c).(*data.ClassValue); ok {
		st.container = ccv
	}
	_ = cv.SetProperty("container", c)
	return cv, nil
}

func pipeWithinTransaction(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := pipeRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	n := kitArgCount(ctx)
	if n == 0 {
		st.withinTransaction = data.NewNullValue()
	} else {
		st.withinTransaction = kit.Arg(ctx, 0)
	}
	_ = cv.SetProperty("withinTransaction", st.withinTransaction)
	return cv, nil
}

func kitArgCount(ctx data.Context) int {
	if args := ctx.GetFlatCallArgs(); args != nil {
		return len(args)
	}
	return 0
}
