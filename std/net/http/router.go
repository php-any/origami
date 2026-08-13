package http

import (
	"errors"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	netannotation "github.com/php-any/origami/std/net/annotation"
	netdata "github.com/php-any/origami/std/net/data"
	"github.com/php-any/origami/utils"
)

// RouteCatalogEntry 已声明的路由条目（供 Router::getRoutes 查询）。
type RouteCatalogEntry struct {
	Method      string
	Path        string
	Controller  string
	Action      string
	Middlewares []string
}

type routerGroupState struct {
	prefix      string
	middlewares []string
}

var (
	routerStateMu    sync.RWMutex
	routerGroupStack []routerGroupState
	routeCatalog     []RouteCatalogEntry
)

func NewRouterClass() data.ClassStmt {
	return &RouterClass{}
}

type RouterClass struct {
	node.Node
}

func (r *RouterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(r, ctx.CreateBaseContext()), nil
}

func (r *RouterClass) GetName() string           { return "Net\\Http\\Router" }
func (r *RouterClass) GetExtend() *string        { return nil }
func (r *RouterClass) GetImplements() []string   { return nil }
func (r *RouterClass) GetConstruct() data.Method { return nil }
func (r *RouterClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}
func (r *RouterClass) GetPropertyList() []data.Property { return []data.Property{} }
func (r *RouterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "get", "post", "put", "delete", "patch":
		return &RouterMapMethod{name: strings.ToUpper(name)}, true
	case "group":
		return &RouterGroupMethod{}, true
	case "getRoutes":
		return &RouterGetRoutesMethod{}, true
	}
	return nil, false
}

func (r *RouterClass) GetMethods() []data.Method {
	return []data.Method{
		&RouterMapMethod{name: "GET"},
		&RouterMapMethod{name: "POST"},
		&RouterMapMethod{name: "PUT"},
		&RouterMapMethod{name: "DELETE"},
		&RouterMapMethod{name: "PATCH"},
		&RouterGroupMethod{},
		&RouterGetRoutesMethod{},
	}
}

func currentRouterGroup() routerGroupState {
	routerStateMu.RLock()
	defer routerStateMu.RUnlock()
	if len(routerGroupStack) == 0 {
		return routerGroupState{}
	}
	current := routerGroupStack[len(routerGroupStack)-1]
	current.middlewares = append([]string(nil), current.middlewares...)
	return current
}

func pushRouterGroup(state routerGroupState) {
	routerStateMu.Lock()
	defer routerStateMu.Unlock()
	routerGroupStack = append(routerGroupStack, state)
}

func popRouterGroup() {
	routerStateMu.Lock()
	defer routerStateMu.Unlock()
	if len(routerGroupStack) > 0 {
		routerGroupStack = routerGroupStack[:len(routerGroupStack)-1]
	}
}

func normalizeRoutePrefix(prefix string) string {
	if prefix == "" {
		return ""
	}
	if prefix == "/" {
		return "/"
	}
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	for len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	return prefix
}

func joinRoutePath(prefix, path string) string {
	P := normalizeRoutePrefix(prefix)
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	if P == "" || P == "/" {
		return path
	}
	if path == "/" {
		return P
	}
	return P + path
}

func parseRouteAction(action data.GetValue) (controllerClass, methodName string, acl data.Control) {
	arr, ok := action.(*data.ArrayValue)
	if !ok {
		return "", "", utils.NewThrow(errors.New("路由 action 必须是 [Controller::class, 'method'] 数组"))
	}
	items := arr.ToValueList()
	if len(items) < 2 {
		return "", "", utils.NewThrow(errors.New("路由 action 数组至少需要 2 个元素"))
	}
	classVal, ok := items[0].(data.AsString)
	if !ok {
		return "", "", utils.NewThrow(errors.New("路由 action[0] 必须是控制器类名"))
	}
	methodVal, ok := items[1].(data.AsString)
	if !ok {
		return "", "", utils.NewThrow(errors.New("路由 action[1] 必须是方法名"))
	}
	controllerClass = classVal.AsString()
	methodName = methodVal.AsString()
	if controllerClass == "" || methodName == "" {
		return "", "", utils.NewThrow(errors.New("控制器类名与方法名不能为空"))
	}
	return controllerClass, methodName, nil
}

func parseMiddlewareList(v data.GetValue) ([]string, data.Control) {
	if v == nil {
		return nil, nil
	}
	switch val := v.(type) {
	case *data.ArrayValue:
		out := make([]string, 0, len(val.List))
		for _, z := range val.List {
			if z == nil {
				continue
			}
			s, ok := z.Value.(data.AsString)
			if !ok {
				return nil, utils.NewThrow(errors.New("middleware 数组元素必须是类名字符串"))
			}
			name := s.AsString()
			if name != "" {
				out = append(out, name)
			}
		}
		return out, nil
	default:
		s, ok := v.(data.AsString)
		if !ok {
			return nil, utils.NewThrow(errors.New("middleware 必须是类名或类名数组"))
		}
		name := s.AsString()
		if name == "" {
			return nil, nil
		}
		return []string{name}, nil
	}
}

func registerProgrammaticRoute(ctx data.Context, method, path string, action data.GetValue, extraMiddlewares []string) data.Control {
	vm := ctx.GetVM()
	if !netdata.SupportsHTTPRoutes(vm) {
		return utils.NewThrow(errors.New("Router 路由注册需在 HTTP 应用引导上下文中调用"))
	}

	controllerClass, methodName, acl := parseRouteAction(action)
	if acl != nil {
		return acl
	}

	group := currentRouterGroup()
	fullPath := joinRoutePath(group.prefix, path)
	middlewares := append(append([]string{}, group.middlewares...), extraMiddlewares...)

	classStmt, acl := vm.GetOrLoadClass(controllerClass)
	if acl != nil {
		return acl
	}
	if classStmt == nil {
		return utils.NewThrowf("无法加载控制器: %s", controllerClass)
	}

	targetMethod, ok := lookupClassMethod(vm, classStmt, methodName)
	if !ok {
		return utils.NewThrowf("控制器 %s 不存在方法 %s", controllerClass, methodName)
	}

	netdata.RegisterDeferredController(controllerClass, classStmt, ctx)
	classValue := data.NewClassValue(classStmt, ctx)

	var staticReceiver data.GetValue
	if targetMethod.GetIsStatic() {
		staticReceiver = classValue
	}

	effective, spec, acl := netannotation.ExpandHTTPHandlerMethod(targetMethod, ctx, fullPath)
	if acl != nil {
		return acl
	}

	netdata.AddPendingRoute(netdata.PendingRoute{
		Method:         method,
		Path:           fullPath,
		Target:         effective,
		ControllerName: controllerClass,
		StaticReceiver: staticReceiver,
		HandlerSpec:    spec,
		Middlewares:    middlewares,
	})

	routerStateMu.Lock()
	routeCatalog = append(routeCatalog, RouteCatalogEntry{
		Method:      method,
		Path:        fullPath,
		Controller:  controllerClass,
		Action:      methodName,
		Middlewares: append([]string{}, middlewares...),
	})
	routerStateMu.Unlock()
	return nil
}

func lookupClassMethod(vm data.VM, classStmt data.ClassStmt, methodName string) (data.Method, bool) {
	if classStmt == nil {
		return nil, false
	}
	if m, ok := classStmt.GetMethod(methodName); ok && m != nil {
		return m, true
	}
	last := classStmt
	for last.GetExtend() != nil {
		parentName := last.GetExtend()
		if parentName == nil || *parentName == "" {
			break
		}
		parent, acl := vm.GetOrLoadClass(*parentName)
		if acl != nil || parent == nil {
			return nil, false
		}
		if m, ok := parent.GetMethod(methodName); ok && m != nil {
			return m, true
		}
		last = parent
	}
	return nil, false
}

// RouterMapMethod GET/POST/PUT/DELETE/PATCH 路由注册。
type RouterMapMethod struct {
	name string
}

func (m *RouterMapMethod) GetName() string            { return strings.ToLower(m.name) }
func (m *RouterMapMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RouterMapMethod) GetIsStatic() bool          { return true }
func (m *RouterMapMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.String{}),
		node.NewParameter(nil, "action", 1, nil, data.Arrays{}),
	}
}
func (m *RouterMapMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.String{}),
		node.NewVariable(nil, "action", 1, data.Arrays{}),
	}
}
func (m *RouterMapMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (m *RouterMapMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	path, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	action, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 action 参数"))
	}
	if acl := registerProgrammaticRoute(ctx, m.name, path, action, nil); acl != nil {
		return nil, acl
	}
	return nil, nil
}

// RouterGroupMethod 路由分组：Router::group(['prefix' => 'api', 'middleware' => [...]], fn)
type RouterGroupMethod struct{}

func (m *RouterGroupMethod) GetName() string            { return "group" }
func (m *RouterGroupMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RouterGroupMethod) GetIsStatic() bool          { return true }
func (m *RouterGroupMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "attributes", 0, nil, data.Arrays{}),
		node.NewParameter(nil, "callback", 1, nil, nil),
	}
}
func (m *RouterGroupMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "attributes", 0, data.Arrays{}),
		node.NewVariable(nil, "callback", 1, nil),
	}
}
func (m *RouterGroupMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (m *RouterGroupMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	attrsVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 attributes 参数"))
	}
	callback, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 callback 参数"))
	}

	parent := currentRouterGroup()
	next := routerGroupState{
		prefix:      parent.prefix,
		middlewares: append([]string{}, parent.middlewares...),
	}

	switch val := attrsVal.(type) {
	case *data.ArrayValue:
		for _, z := range val.List {
			if z == nil {
				continue
			}
			switch z.Name {
			case "prefix":
				if ps, ok := z.Value.(data.AsString); ok {
					next.prefix = joinRoutePath(next.prefix, ps.AsString())
				}
			case "middleware":
				mws, acl := parseMiddlewareList(z.Value)
				if acl != nil {
					return nil, acl
				}
				next.middlewares = append(next.middlewares, mws...)
			}
		}
	case *data.ObjectValue:
		if prefixVal, ctl := val.GetProperty("prefix"); ctl == nil && prefixVal != nil {
			if ps, ok := prefixVal.(data.AsString); ok {
				next.prefix = joinRoutePath(next.prefix, ps.AsString())
			}
		}
		if mwVal, ctl := val.GetProperty("middleware"); ctl == nil && mwVal != nil {
			mws, acl := parseMiddlewareList(mwVal)
			if acl != nil {
				return nil, acl
			}
			next.middlewares = append(next.middlewares, mws...)
		}
	}

	pushRouterGroup(next)
	defer popRouterGroup()

	fn, acl := callback.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if fv, ok := fn.(*data.FuncValue); ok {
		_, acl = fv.Call(ctx)
		return nil, acl
	}
	return nil, utils.NewThrow(errors.New("group callback 必须是可调用闭包"))
}

// RouterGetRoutesMethod 返回已声明的路由列表。
type RouterGetRoutesMethod struct{}

func (m *RouterGetRoutesMethod) GetName() string            { return "getRoutes" }
func (m *RouterGetRoutesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RouterGetRoutesMethod) GetIsStatic() bool          { return true }
func (m *RouterGetRoutesMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *RouterGetRoutesMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *RouterGetRoutesMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *RouterGetRoutesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	routerStateMu.RLock()
	catalog := append([]RouteCatalogEntry(nil), routeCatalog...)
	routerStateMu.RUnlock()
	catalogIndex := make(map[string]RouteCatalogEntry, len(catalog))
	for _, rt := range catalog {
		catalogIndex[rt.Method+" "+rt.Path] = rt
	}

	routes := netdata.HTTPRoutes()
	if len(routes) == 0 {
		// 尚未 RegisterPendingRoutes 时，仅返回声明式目录
		items := make([]data.Value, 0, len(catalog))
		for _, rt := range catalog {
			items = append(items, routeEntryToObject(rt))
		}
		return data.NewArrayValue(items), nil
	}

	items := make([]data.Value, 0, len(routes))
	for _, rt := range routes {
		key := rt.Method + " " + rt.Path
		if cat, ok := catalogIndex[key]; ok {
			items = append(items, routeEntryToObject(cat))
			continue
		}
		obj := data.NewObjectValue()
		obj.SetProperty("method", data.NewStringValue(rt.Method))
		obj.SetProperty("path", data.NewStringValue(rt.Path))
		obj.SetProperty("controller", data.NewStringValue("(annotation)"))
		action := ""
		if rt.Target != nil {
			action = rt.Target.GetName()
		}
		obj.SetProperty("action", data.NewStringValue(action))
		mwItems := make([]data.Value, 0, len(rt.Middlewares))
		for _, mw := range rt.Middlewares {
			mwItems = append(mwItems, data.NewStringValue(mw.ClassName))
		}
		obj.SetProperty("middleware", data.NewArrayValue(mwItems))
		items = append(items, obj)
	}
	return data.NewArrayValue(items), nil
}

func routeEntryToObject(rt RouteCatalogEntry) data.Value {
	obj := data.NewObjectValue()
	obj.SetProperty("method", data.NewStringValue(rt.Method))
	obj.SetProperty("path", data.NewStringValue(rt.Path))
	obj.SetProperty("controller", data.NewStringValue(rt.Controller))
	obj.SetProperty("action", data.NewStringValue(rt.Action))
	mwItems := make([]data.Value, 0, len(rt.Middlewares))
	for _, mw := range rt.Middlewares {
		mwItems = append(mwItems, data.NewStringValue(mw))
	}
	obj.SetProperty("middleware", data.NewArrayValue(mwItems))
	return obj
}

// ResetRouterState 清空路由注册状态，供热重载使用。
func ResetRouterState() {
	routerStateMu.Lock()
	defer routerStateMu.Unlock()
	routerGroupStack = nil
	routeCatalog = nil
}
