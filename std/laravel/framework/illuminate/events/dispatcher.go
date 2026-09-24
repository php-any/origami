package events

import (
	"fmt"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const dispatcherClassName = "Illuminate\\Events\\Dispatcher"

type pushedEvent struct {
	event   string
	payload data.Value
}

type dispatcherState struct {
	mu             sync.Mutex
	listeners      map[string][]data.Value
	wildcards      map[string][]data.Value
	wildcardsCache map[string][]data.Value
	pushed         map[string][]pushedEvent
	queueResolver  data.Value
	txResolver     data.Value
	deferring      bool
	deferred       []deferredEvent
	eventsToDefer  map[string]struct{} // nil = all
}

type deferredEvent struct {
	event   string
	payload data.Value
	halt    bool
}

// dispatcherStateProp：Go 侧状态挂在对象隐蔽属性上，由 GC 随对象回收（见 kit.CachedState）。
const dispatcherStateProp = "__origami_dispatcher_state"

func newDispatcherState() *dispatcherState {
	return &dispatcherState{
		listeners:      map[string][]data.Value{},
		wildcards:      map[string][]data.Value{},
		wildcardsCache: map[string][]data.Value{},
		pushed:         map[string][]pushedEvent{},
	}
}

func stateOf(cv *data.ClassValue) *dispatcherState {
	if cv == nil {
		return newDispatcherState()
	}
	return kit.CachedState(cv, dispatcherStateProp, newDispatcherState)
}

// DispatcherClass 对齐 Illuminate\Events\Dispatcher 核心 API。
type DispatcherClass struct {
	node.Node
	methods map[string]data.Method
}

func NewDispatcherClass() data.ClassStmt {
	c := &DispatcherClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *DispatcherClass) GetName() string    { return dispatcherClassName }
func (c *DispatcherClass) GetExtend() *string { return nil }
func (c *DispatcherClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Events\\Dispatcher"}
}
func (c *DispatcherClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "container", "queueResolver", "transactionManagerResolver":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	case "listeners", "wildcards", "wildcardsCache", "deferredEvents":
		return node.NewProperty(nil, name, "protected", false, data.NewArrayValue(nil)), true
	case "deferringEvents":
		return node.NewProperty(nil, name, "protected", false, data.NewBoolValue(false)), true
	case "__origami_dispatcher_id":
		return node.NewProperty(nil, name, "protected", false, data.NewIntValue(0)), true
	}
	return nil, false
}
func (c *DispatcherClass) GetPropertyList() []data.Property {
	names := []string{"container", "listeners", "wildcards", "wildcardsCache", "queueResolver", "transactionManagerResolver", "deferredEvents", "deferringEvents"}
	out := make([]data.Property, 0, len(names))
	for _, n := range names {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *DispatcherClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *DispatcherClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DispatcherClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *DispatcherClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *DispatcherClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *DispatcherClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"container"}, 0, dispConstruct)
	c.methods["listen"] = kit.InstanceMethodOpt("listen", []string{"events", "listener"}, 1, dispListen)
	c.methods["haslisteners"] = kit.InstanceMethod("hasListeners", []string{"eventName"}, dispHasListeners)
	c.methods["haswildcardlisteners"] = kit.InstanceMethod("hasWildcardListeners", []string{"eventName"}, dispHasWildcardListeners)
	c.methods["push"] = kit.InstanceMethodOpt("push", []string{"event", "payload"}, 1, dispPush)
	c.methods["flush"] = kit.InstanceMethod("flush", []string{"event"}, dispFlush)
	c.methods["subscribe"] = kit.InstanceMethod("subscribe", []string{"subscriber"}, dispSubscribe)
	c.methods["until"] = kit.InstanceMethodOpt("until", []string{"event", "payload"}, 1, dispUntil)
	c.methods["dispatch"] = kit.InstanceMethodOpt("dispatch", []string{"event", "payload", "halt"}, 1, dispDispatch)
	c.methods["getlisteners"] = kit.InstanceMethod("getListeners", []string{"eventName"}, dispGetListeners)
	c.methods["forget"] = kit.InstanceMethod("forget", []string{"event"}, dispForget)
	c.methods["forgetpushed"] = kit.InstanceMethod("forgetPushed", nil, dispForgetPushed)
	c.methods["setqueueresolver"] = kit.InstanceMethod("setQueueResolver", []string{"resolver"}, dispSetQueueResolver)
	c.methods["settransactionmanagerresolver"] = kit.InstanceMethod("setTransactionManagerResolver", []string{"resolver"}, dispSetTxResolver)
	c.methods["defer"] = kit.InstanceMethodOpt("defer", []string{"callback", "events"}, 1, dispDefer)
	c.methods["getrawlisteners"] = kit.InstanceMethod("getRawListeners", nil, dispGetRawListeners)
	kit.RegisterMacroable(c.methods, dispatcherClassName)
}

func dispRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Events\\Dispatcher missing $this"))
}

func dispConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = stateOf(cv)
	container := kit.Arg(ctx, 0)
	if container == nil || kit.IsNull(container) {
		// 尝试 new Illuminate\Container\Container
		if vm := ctx.GetVM(); vm != nil {
			if cls, ok := vm.GetClass("Illuminate\\Container\\Container"); ok && cls != nil {
				container = data.NewClassValue(cls, ctx.CreateBaseContext())
				if ctor := cls.GetConstruct(); ctor != nil {
					nctx := container.(*data.ClassValue).CreateContext(ctor.GetVariables())
					if _, ctl := ctor.Call(nctx); ctl != nil {
						return nil, ctl
					}
				}
			}
		}
	}
	if container != nil {
		_ = cv.SetProperty("container", container)
	}
	return cv, nil
}

func eventNameOf(v data.Value) string {
	v = kit.Unwrap(v)
	if v == nil {
		return ""
	}
	if s, ok := v.(*data.StringValue); ok {
		return s.AsString()
	}
	if cv, ok := v.(*data.ClassValue); ok && cv.Class != nil {
		return cv.Class.GetName()
	}
	if ov, ok := v.(*data.ObjectValue); ok {
		_ = ov
	}
	return v.AsString()
}

func strIs(pattern, value string) bool {
	if pattern == value {
		return true
	}
	if !strings.Contains(pattern, "*") {
		return false
	}
	parts := strings.Split(pattern, "*")
	if !strings.HasPrefix(value, parts[0]) {
		return false
	}
	value = value[len(parts[0]):]
	for i := 1; i < len(parts)-1; i++ {
		idx := strings.Index(value, parts[i])
		if idx < 0 {
			return false
		}
		value = value[idx+len(parts[i]):]
	}
	return strings.HasSuffix(value, parts[len(parts)-1])
}

func dispListen(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := stateOf(cv)
	events := kit.Arg(ctx, 0)
	listener := kit.Arg(ctx, 1)
	names := eventNames(events)
	st.mu.Lock()
	defer st.mu.Unlock()
	for _, name := range names {
		if name == "" {
			continue
		}
		wildcard := strings.Contains(name, "*")
		wrapped := makeListenerWrapper(cv, listener, wildcard)
		if wildcard {
			st.wildcards[name] = append(st.wildcards[name], wrapped)
			st.wildcardsCache = map[string][]data.Value{}
		} else {
			st.listeners[name] = append(st.listeners[name], wrapped)
		}
	}
	return data.NewNullValue(), nil
}

func eventNames(v data.Value) []string {
	v = kit.Unwrap(v)
	if v == nil {
		return nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		out := make([]string, 0)
		for _, e := range kit.Entries(av) {
			out = append(out, eventNameOf(e.Value))
		}
		return out
	}
	return []string{eventNameOf(v)}
}

func dispHasListeners(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name := eventNameOf(kit.Arg(ctx, 0))
	st := stateOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	if len(st.listeners[name]) > 0 || len(st.wildcards[name]) > 0 {
		return data.NewBoolValue(true), nil
	}
	for pat := range st.wildcards {
		if strIs(pat, name) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func dispHasWildcardListeners(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name := eventNameOf(kit.Arg(ctx, 0))
	st := stateOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	for pat := range st.wildcards {
		if strIs(pat, name) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func dispPush(ctx data.Context) (data.GetValue, data.Control) {
	event := eventNameOf(kit.Arg(ctx, 0))
	payload := kit.Arg(ctx, 1)
	if payload == nil {
		payload = data.NewArrayValue(nil)
	}
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := stateOf(cv)
	st.mu.Lock()
	if st.pushed == nil {
		st.pushed = map[string][]pushedEvent{}
	}
	st.pushed[event] = append(st.pushed[event], pushedEvent{event: event, payload: payload})
	st.mu.Unlock()
	return data.NewNullValue(), nil
}

func dispFlush(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	event := eventNameOf(kit.Arg(ctx, 0))
	st := stateOf(cv)
	st.mu.Lock()
	list := append([]pushedEvent(nil), st.pushed[event]...)
	delete(st.pushed, event)
	st.mu.Unlock()
	for _, item := range list {
		if _, ctl := dispatchNamed(ctx, cv, item.event, item.payload, false); ctl != nil {
			return nil, ctl
		}
	}
	return data.NewNullValue(), nil
}

func dispSubscribe(ctx data.Context) (data.GetValue, data.Control) {
	// 最小实现：若 subscriber 有 subscribe 方法则调用
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	sub := kit.Unwrap(kit.Arg(ctx, 0))
	if scv, ok := sub.(*data.ClassValue); ok {
		if m, ok := scv.GetMethod("subscribe"); ok && m != nil {
			nctx := scv.CreateContext(m.GetVariables())
			data.BindDeclaredArgs(nctx, m, []data.Value{cv})
			ret, ctl := m.Call(nctx)
			if ctl != nil {
				return nil, ctl
			}
			if av, ok := asValue(ret).(*data.ArrayValue); ok {
				for _, e := range kit.Entries(av) {
					event := e.KeyStr
					for _, lis := range eventNames(e.Value) {
						_ = lis
					}
					// listen each
					st := stateOf(cv)
					st.mu.Lock()
					for _, lisE := range kit.Entries(e.Value) {
						st.listeners[event] = append(st.listeners[event], makeListenerWrapper(cv, lisE.Value, false))
					}
					if _, isArr := e.Value.(*data.ArrayValue); !isArr {
						st.listeners[event] = append(st.listeners[event], makeListenerWrapper(cv, e.Value, false))
					}
					st.mu.Unlock()
				}
			}
		}
	}
	return data.NewNullValue(), nil
}

func asValue(v data.GetValue) data.Value {
	if v == nil {
		return data.NewNullValue()
	}
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}

func dispUntil(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name, payload := parseEventAndPayload(kit.Arg(ctx, 0), kit.Arg(ctx, 1))
	return dispatchNamed(ctx, cv, name, payload, true)
}

func dispDispatch(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	event := kit.Arg(ctx, 0)
	payload := kit.Arg(ctx, 1)
	halt := false
	if h := kit.Arg(ctx, 2); h != nil {
		if b, ok := h.(data.AsBool); ok {
			halt, _ = b.AsBool()
		}
	}
	name, payload := parseEventAndPayload(event, payload)
	return dispatchNamed(ctx, cv, name, payload, halt)
}

// parseEventAndPayload 对齐 Dispatcher::parseEventAndPayload。
func parseEventAndPayload(event, payload data.Value) (string, data.Value) {
	event = kit.Unwrap(event)
	if cv, ok := event.(*data.ClassValue); ok && cv != nil && cv.Class != nil {
		return cv.Class.GetName(), data.NewArrayValue([]data.Value{event})
	}
	return eventNameOf(event), ensureArrayPayload(payload)
}

func ensureArrayPayload(payload data.Value) data.Value {
	payload = kit.Unwrap(payload)
	if payload == nil || kit.IsNull(payload) {
		return data.NewArrayValue(nil)
	}
	if _, ok := payload.(*data.ArrayValue); ok {
		return payload
	}
	return data.NewArrayValue([]data.Value{payload})
}

func dispatchNamed(ctx data.Context, cv *data.ClassValue, name string, payload data.Value, halt bool) (data.GetValue, data.Control) {
	st := stateOf(cv)
	st.mu.Lock()
	if st.deferring {
		if st.eventsToDefer == nil {
			st.deferred = append(st.deferred, deferredEvent{event: name, payload: payload, halt: halt})
			st.mu.Unlock()
			return data.NewArrayValue(nil), nil
		}
		if _, ok := st.eventsToDefer[name]; ok {
			st.deferred = append(st.deferred, deferredEvent{event: name, payload: payload, halt: halt})
			st.mu.Unlock()
			return data.NewArrayValue(nil), nil
		}
	}
	exact, wild := collectListenersSplit(st, name)
	st.mu.Unlock()

	responses := data.NewArrayValue(nil).(*data.ArrayValue)
	run := func(lis data.Value) (data.GetValue, data.Control, bool) {
		// 对齐 invokeListeners：始终 $listener($event, $payload)
		ret, ctl := invokeWrappedListener(ctx, lis, name, payload)
		if ctl != nil {
			return nil, ctl, false
		}
		val := asValue(ret)
		if halt && val != nil && !kit.IsNull(val) {
			return val, nil, true
		}
		if !halt {
			responses.List = append(responses.List, data.NewZVal(val))
		}
		return nil, nil, false
	}
	for _, lis := range exact {
		if ret, ctl, done := run(lis); ctl != nil {
			return nil, ctl
		} else if done {
			return ret, nil
		}
	}
	for _, lis := range wild {
		if ret, ctl, done := run(lis); ctl != nil {
			return nil, ctl
		} else if done {
			return ret, nil
		}
	}
	if halt {
		return data.NewNullValue(), nil
	}
	return responses, nil
}

func collectListeners(st *dispatcherState, name string) []data.Value {
	exact, wild := collectListenersSplit(st, name)
	return append(exact, wild...)
}

func collectListenersSplit(st *dispatcherState, name string) (exact, wild []data.Value) {
	exact = append([]data.Value(nil), st.listeners[name]...)
	if cached, ok := st.wildcardsCache[name]; ok {
		return exact, append([]data.Value(nil), cached...)
	}
	for pat, list := range st.wildcards {
		if strIs(pat, name) {
			wild = append(wild, list...)
		}
	}
	st.wildcardsCache[name] = wild
	return exact, wild
}

func invokeWrappedListener(ctx data.Context, listener data.Value, event string, payload data.Value) (data.GetValue, data.Control) {
	listener = kit.Unwrap(listener)
	payload = ensureArrayPayload(payload)
	args := []data.Value{data.NewStringValue(event), payload}
	switch t := listener.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return kit.Call(ctx, t, args...)
	default:
		// 兼容未包装的旧 listener（subscribe 等路径遗漏时）
		return invokeRawListener(ctx, nil, listener, listenerArgs(event, payload, false))
	}
}

// listenerArgs 对齐 makeListener：通配符传 ($event, $payload)，否则 ...$payload。
func listenerArgs(event string, payload data.Value, wildcard bool) []data.Value {
	payload = ensureArrayPayload(payload)
	if wildcard {
		return []data.Value{data.NewStringValue(event), payload}
	}
	if av, ok := payload.(*data.ArrayValue); ok {
		return av.ToValueList()
	}
	return []data.Value{payload}
}

func invokeStringListener(ctx data.Context, cv *data.ClassValue, listener string, args []data.Value) (data.GetValue, data.Control) {
	className, method := listener, "handle"
	if i := strings.Index(listener, "@"); i >= 0 {
		className, method = listener[:i], listener[i+1:]
	}
	return invokeClassMethod(ctx, cv, className, method, args)
}

func invokeClassMethod(ctx data.Context, cv *data.ClassValue, className, method string, args []data.Value) (data.GetValue, data.Control) {
	container, _ := cv.GetProperty("container")
	var instance data.Value
	if ccv, ok := kit.Unwrap(container).(*data.ClassValue); ok {
		if m, ok := ccv.GetMethod("make"); ok && m != nil {
			nctx := ccv.CreateContext(m.GetVariables())
			data.BindDeclaredArgs(nctx, m, []data.Value{data.NewStringValue(className)})
			ret, ctl := m.Call(nctx)
			if ctl != nil {
				return nil, ctl
			}
			instance = asValue(ret)
		}
	}
	if instance == nil {
		vm := ctx.GetVM()
		cls, ctl := vm.GetOrLoadClass(className)
		if ctl != nil {
			return nil, ctl
		}
		instance = data.NewClassValue(cls, ctx.CreateBaseContext())
	}
	scv, ok := kit.Unwrap(instance).(*data.ClassValue)
	if !ok {
		return data.NewNullValue(), nil
	}
	m, ok := scv.GetMethod(method)
	if !ok || m == nil {
		return data.NewNullValue(), nil
	}
	nctx := scv.CreateContext(m.GetVariables())
	data.BindDeclaredArgs(nctx, m, args)
	return m.Call(nctx)
}

func dispGetListeners(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name := eventNameOf(kit.Arg(ctx, 0))
	st := stateOf(cv)
	st.mu.Lock()
	list := collectListeners(st, name)
	st.mu.Unlock()
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, lis := range list {
		out.List = append(out.List, data.NewZVal(lis))
	}
	return out, nil
}

func dispForget(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	event := eventNameOf(kit.Arg(ctx, 0))
	st := stateOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	if strings.Contains(event, "*") {
		delete(st.wildcards, event)
	} else {
		delete(st.listeners, event)
	}
	st.wildcardsCache = map[string][]data.Value{}
	return data.NewNullValue(), nil
}

func dispForgetPushed(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := stateOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	for k := range st.listeners {
		if strings.HasSuffix(k, "_pushed") {
			delete(st.listeners, k)
		}
	}
	return data.NewNullValue(), nil
}

func dispSetQueueResolver(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	stateOf(cv).queueResolver = kit.Arg(ctx, 0)
	_ = cv.SetProperty("queueResolver", kit.Arg(ctx, 0))
	return cv, nil
}

func dispSetTxResolver(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	stateOf(cv).txResolver = kit.Arg(ctx, 0)
	_ = cv.SetProperty("transactionManagerResolver", kit.Arg(ctx, 0))
	return cv, nil
}

func dispDefer(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := stateOf(cv)
	cb := kit.Arg(ctx, 0)
	eventsArg := kit.Arg(ctx, 1)
	st.mu.Lock()
	st.deferring = true
	st.deferred = nil
	if eventsArg != nil && !kit.IsNull(eventsArg) {
		st.eventsToDefer = map[string]struct{}{}
		for _, n := range eventNames(eventsArg) {
			st.eventsToDefer[n] = struct{}{}
		}
	} else {
		st.eventsToDefer = nil
	}
	st.mu.Unlock()

	ret, ctl := kit.Call(ctx, cb)
	st.mu.Lock()
	st.deferring = false
	pending := append([]deferredEvent(nil), st.deferred...)
	st.deferred = nil
	st.eventsToDefer = nil
	st.mu.Unlock()
	if ctl != nil {
		return nil, ctl
	}
	for _, ev := range pending {
		if _, ctl := dispatchNamed(ctx, cv, ev.event, ev.payload, ev.halt); ctl != nil {
			return nil, ctl
		}
	}
	return ret, nil
}

func dispGetRawListeners(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := dispRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := stateOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for name, list := range st.listeners {
		arr := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, lis := range list {
			arr.List = append(arr.List, data.NewZVal(lis))
		}
		zv := data.NewZVal(arr)
		zv.Name = name
		out.List = append(out.List, zv)
	}
	return out, nil
}
