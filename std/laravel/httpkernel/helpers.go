package httpkernel

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// kernelMethod 是共享的方法适配器。
type kernelMethod struct {
	name     string
	modifier data.Modifier
	static   bool
	params   []data.GetValue
	vars     []data.Variable
	ret      data.Types
	fn       func(ctx data.Context) (data.GetValue, data.Control)
}

func (m *kernelMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.fn(ctx)
}
func (m *kernelMethod) GetName() string               { return m.name }
func (m *kernelMethod) GetModifier() data.Modifier    { return m.modifier }
func (m *kernelMethod) GetIsStatic() bool             { return m.static }
func (m *kernelMethod) GetParams() []data.GetValue    { return m.params }
func (m *kernelMethod) GetVariables() []data.Variable { return m.vars }
func (m *kernelMethod) GetReturnType() data.Types     { return m.ret }

func pubMethod(name string, params []data.GetValue, vars []data.Variable, ret data.Types, fn func(ctx data.Context) (data.GetValue, data.Control)) data.Method {
	return &kernelMethod{
		name:     name,
		modifier: data.ModifierPublic,
		params:   params,
		vars:     vars,
		ret:      ret,
		fn:       fn,
	}
}

func param(name string, index int, def data.GetValue, ty data.Types) data.GetValue {
	return node.NewParameter(nil, name, index, def, ty)
}

func variable(name string, index int, ty data.Types) data.Variable {
	return node.NewVariable(nil, name, index, ty)
}

func thisClassValue(ctx data.Context) *data.ClassValue {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ClassValue
	case *data.ClassValue:
		return c
	default:
		return nil
	}
}

func getState(ctx data.Context) (*kernelState, data.Control) {
	cv := thisClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 无法获取 Kernel 实例"))
	}
	if s, ok := cv.GetSource().(*kernelState); ok && s != nil {
		return s, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: Kernel 状态未初始化"))
}

func indexValue(ctx data.Context, index int) data.Value {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return nil
	}
	return v
}

func valueAsString(v data.Value) string {
	if v == nil {
		return ""
	}
	return v.AsString()
}

func isTruthy(v data.Value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		bv, err := b.AsBool()
		return err == nil && bv
	}
	s := v.AsString()
	return s != "" && s != "0" && s != "false"
}

func stringsToArrayValue(ss []string) *data.ArrayValue {
	vals := make([]data.Value, len(ss))
	for i, s := range ss {
		vals[i] = data.NewStringValue(s)
	}
	return data.NewArrayValue(vals).(*data.ArrayValue)
}

func stringListFromValue(v data.Value) []string {
	if v == nil {
		return nil
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		out := make([]string, 0, len(arr.List))
		for _, z := range arr.List {
			if z == nil || z.Value == nil {
				continue
			}
			out = append(out, z.Value.AsString())
		}
		return out
	case *data.ObjectValue:
		out := make([]string, 0)
		arr.RangeProperties(func(_ string, value data.Value) bool {
			if value != nil {
				out = append(out, value.AsString())
			}
			return true
		})
		return out
	default:
		s := v.AsString()
		if s == "" {
			return nil
		}
		return []string{s}
	}
}

func stringMapFromValue(v data.Value) map[string]string {
	out := make(map[string]string)
	if v == nil {
		return out
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, z := range arr.List {
			if z == nil || z.Value == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			out[key] = z.Value.AsString()
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, value data.Value) bool {
			if value != nil {
				out[key] = value.AsString()
			}
			return true
		})
	}
	return out
}

func groupsFromValue(v data.Value) map[string][]string {
	out := make(map[string][]string)
	if v == nil {
		return out
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, z := range arr.List {
			if z == nil || z.Value == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			out[key] = stringListFromValue(z.Value)
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, value data.Value) bool {
			out[key] = stringListFromValue(value)
			return true
		})
	}
	return out
}

func groupsToArrayValue(groups map[string][]string) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(groups))
	for k, items := range groups {
		list = append(list, &data.ZVal{Name: k, Value: stringsToArrayValue(items)})
	}
	return &data.ArrayValue{List: list}
}

func aliasesToArrayValue(aliases map[string]string) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(aliases))
	for k, v := range aliases {
		list = append(list, &data.ZVal{Name: k, Value: data.NewStringValue(v)})
	}
	return &data.ArrayValue{List: list}
}

func syncProperties(cv *data.ClassValue, s *kernelState) {
	if cv == nil || s == nil {
		return
	}
	if s.app != nil {
		_ = cv.SetProperty("app", s.app)
	}
	if s.router != nil {
		_ = cv.SetProperty("router", s.router)
	}
	_ = cv.SetProperty("bootstrappers", stringsToArrayValue(s.bootstrappers))
	_ = cv.SetProperty("middleware", stringsToArrayValue(s.middleware))
	_ = cv.SetProperty("middlewareGroups", groupsToArrayValue(s.middlewareGroups))
	_ = cv.SetProperty("middlewareAliases", aliasesToArrayValue(s.middlewareAliases))
	_ = cv.SetProperty("middlewarePriority", stringsToArrayValue(s.middlewarePriority))
}

const (
	fqnConfigRepository = "Illuminate\\Contracts\\Config\\Repository"
	fqnTelescope        = "Laravel\\Telescope\\Telescope"
)

// telescopeCache 缓存「请求无关」的 Telescope 判定结果，由常驻 kernelState 持有、
// 请求沙箱共享同一个指针（见 Sandbox）。字段在 once.Do 完成后只读，并发请求无锁可读。
//
// config('telescope.enabled') 自 bootstrap 起不再变化，匹配用的 patterns 也完全由 config
// 决定，Telescope 是否真被 boot 同样只在启动期确定一次。原实现每请求都重算：
// 2 次 $app->make() + 3 次 $config->get() + 1 次多 pattern 的 $request->is()，
// 实测占单请求 CPU 的 24%（examples/laravel13 /hello，p50=7ms）。
//
// disabled 为真时 Telescope 根本没有 watcher 注册（见下），整个同步退化成一次 bool 判断。
type telescopeCache struct {
	once     sync.Once
	disabled bool
	obj      data.Value
	patterns []string // 默认模式：命中即「不该记录」
	only     []string // telescope.only_paths 非空时改用它判定「该记录」
}

func (t *telescopeCache) resolve(ctx data.Context, app data.Value) {
	if app == nil {
		t.disabled = true
		return
	}
	config, ctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue(fqnConfigRepository))
	if ctl != nil || config == nil {
		t.disabled = true
		return
	}
	configObj := asValue(config)
	if configObj == nil {
		t.disabled = true
		return
	}
	// 与 TelescopeServiceProvider::boot 的 `if (! config('telescope.enabled'))` 同义：
	// 取不到时按假处理（provider 读到 null 同样是假）。
	enabled, ctl := callObjectMethodInContext(ctx, configObj, "get",
		data.NewStringValue("telescope.enabled"), data.NewBoolValue(false))
	if ctl != nil || !isTruthy(asValue(enabled)) {
		// TelescopeServiceProvider::boot 在此直接 return：Telescope::start() 从未执行、
		// listenForStorageOpportunities() 从未调用，没有任何 watcher 注册，
		// 对 startRecording/stopRecording 的调用没有任何可观察效果。
		// 反过来跳过它还更接近 php-fpm 语义——那里 $shouldRecord 会一直保持初值 false。
		t.disabled = true
		return
	}
	cfgString := func(key, def string) string {
		ret, ctl := callObjectMethodInContext(ctx, configObj, "get",
			data.NewStringValue(key), data.NewStringValue(def))
		if ctl != nil || ret == nil {
			return def
		}
		return valueAsString(asValue(ret))
	}
	cfgList := func(key string) []string {
		ret, ctl := callObjectMethodInContext(ctx, configObj, "get",
			data.NewStringValue(key), stringsToArrayValue(nil))
		if ctl != nil || ret == nil {
			return nil
		}
		return stringListFromValue(asValue(ret))
	}

	// 对齐 Telescope::requestIsToApprovedUri 的判断逻辑
	if onlyList := cfgList("telescope.only_paths"); len(onlyList) > 0 {
		t.only = onlyList
	} else {
		patterns := []string{"telescope-api*", "vendor/telescope*", "horizon*", "vendor/horizon*"}
		if path := cfgString("telescope.path", "telescope"); path != "" {
			patterns = append(patterns, path+"*")
		}
		t.patterns = append(patterns, cfgList("telescope.ignore_paths")...)
	}

	telescope, ctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue(fqnTelescope))
	if ctl != nil || telescope == nil {
		t.disabled = true
		return
	}
	obj := asValue(telescope)
	if obj == nil {
		t.disabled = true
		return
	}
	t.obj = obj
}

// syncTelescopeRecording 在每次请求开始时同步 Telescope 的记录状态。
// 常驻模式下 Telescope::start() 只在首次 bootstrap 执行一次，$shouldRecord 不会按请求重置，
// 导致后续 /telescope/* 请求也被记录。这里模拟 Octane 的 RequestReceived 语义：
// 当前请求命中 telescope/ignore 路径则停止记录，否则开始记录。
//
// 除「本次请求命中哪些 pattern」外，其余判定都是启动期常量，一律走 telescopeCache。
func syncTelescopeRecording(ctx data.Context, s *kernelState, request data.Value) {
	if s == nil || s.app == nil || request == nil || s.tel == nil {
		return
	}
	t := s.tel
	t.once.Do(func() { t.resolve(ctx, s.app) })
	if t.disabled {
		return
	}

	approved := false
	if len(t.only) > 0 {
		approved = requestMatches(ctx, request, t.only)
	} else {
		approved = !requestMatches(ctx, request, t.patterns)
	}
	if approved {
		_, _ = callObjectMethodInContext(ctx, t.obj, "startRecording", data.NewBoolValue(false))
	} else {
		_, _ = callObjectMethodInContext(ctx, t.obj, "stopRecording")
	}
}

// requestMatches 即 $request->is($patterns)。
func requestMatches(ctx data.Context, request data.Value, patterns []string) bool {
	ret, ctl := callObjectMethodInContext(ctx, request, "is", stringsToArrayValue(patterns))
	if ctl != nil || ret == nil {
		return false
	}
	return isTruthy(asValue(ret))
}

// callObjectMethod 在 ClassValue 上调用实例方法；失败时原样返回 control。
func callObjectMethod(obj data.Value, name string, args ...data.Value) (data.GetValue, data.Control) {
	return callObjectMethodInContext(nil, obj, name, args...)
}

// callObjectMethodInContext 让常驻 Application/Kernel 的方法调用使用当前请求 VM，
// 但继续共享对象属性和类定义。ctx 为 nil 时保持对象原有上下文。
func callObjectMethodInContext(ctx data.Context, obj data.Value, name string, args ...data.Value) (data.GetValue, data.Control) {
	var cv *data.ClassValue
	switch value := obj.(type) {
	case *data.ClassValue:
		cv = value
	case *data.ThisValue:
		cv = value.ClassValue
	}
	ok := cv != nil
	if !ok || cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 期望对象以调用 %s", name))
	}
	if ctx != nil {
		cv = cv.CloneWithContext(ctx.CreateBaseContext())
	}
	method, exists := cv.GetMethod(name)
	if !exists || method == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 方法 %s 不存在", name))
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	data.PreferVM(fnCtx, ctx, cv)
	vars := method.GetVariables()
	for i, arg := range args {
		if i >= len(vars) {
			break
		}
		if arg == nil {
			arg = data.NewNullValue()
		}
		if acl := fnCtx.SetVariableValue(vars[i], arg); acl != nil {
			return nil, acl
		}
	}
	return method.Call(fnCtx)
}

func cloneStringMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func cloneGroupsMap(m map[string][]string) map[string][]string {
	if m == nil {
		return nil
	}
	out := make(map[string][]string, len(m))
	for k, v := range m {
		out[k] = append([]string(nil), v...)
	}
	return out
}

func cloneIfClass(v data.Value, ctx data.Context) data.Value {
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil {
		return v
	}
	return cv.CloneSandbox(ctx)
}

func containsString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func indexOfString(list []string, item string) int {
	for i, v := range list {
		if v == item {
			return i
		}
	}
	return -1
}

func existingNames(v data.Value) []string {
	if v == nil {
		return nil
	}
	if arr, ok := v.(*data.ArrayValue); ok {
		out := make([]string, 0, len(arr.List))
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				out = append(out, z.Value.AsString())
			}
		}
		return out
	}
	s := v.AsString()
	if s == "" {
		return nil
	}
	return []string{s}
}
