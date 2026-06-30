package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// NextFunc 实现 data.FuncStmt，作为洋葱模型中 $next 回调
type NextFunc struct {
	name     string
	fn       func(request data.Value, response data.Value) (data.GetValue, data.Control)
	variable []data.Variable
}

func (f NextFunc) Call(ctx data.Context) (_ data.GetValue, acl data.Control) {
	request, err := utils.ConvertFromIndex[data.Value](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	response, err := utils.ConvertFromIndex[data.Value](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	defer func() {
		if r := recover(); r != nil {
			if acl2, ok2 := r.(data.Control); ok2 {
				acl = acl2
				return
			}
			panic(r)
		}
	}()

	return f.fn(request, response)
}

func (f NextFunc) GetName() string { return f.name }
func (f NextFunc) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
	}
}
func (f NextFunc) GetVariables() []data.Variable {
	return f.variable
}

// DispatchHTTPRoutes 根据 VM 中已注册的注解路由分发当前请求。
// 使用洋葱模型中间件：middleware.handle($request, $response, $next)
func DispatchHTTPRoutes(vm data.VM, ctx data.Context) (data.GetValue, data.Control) {
	routes := runtime.HTTPRoutes(vm)
	if len(routes) == 0 {
		return nil, utils.NewThrowf("未注册 HTTP 路由，请确认应用入口已加载且控制器带有 @Controller/@*Mapping 注解")
	}

	mux := http.NewServeMux()
	var lastACL data.Control
	for _, rt := range routes {
		rt := rt
		// 预提取路由路径中的参数名列表
		pathParamKeys := extractPathParamKeys(rt.Path)
		mux.HandleFunc(rt.Method+" "+rt.Path, func(w http.ResponseWriter, r *http.Request) {
			rw, response := beginResponse(w, r)
			defer rw.commitPending()
			r, request := beginRequest(r)
			defer detachRequestAttrs(r)

			// 将路径参数键名关联到请求，以便 all/input/only/except/has/post/route 等方法获取
			setPathValueKeys(r, pathParamKeys)

			reqProxy := data.NewProxyValue(request, ctx)
			resProxy := data.NewProxyValue(response, ctx)

			// 使用洋葱模型执行中间件链 + 控制器
			_, acl := executeMiddlewareChain(vm, ctx, rt, reqProxy, resProxy)
			if acl != nil {
				lastACL = acl
			}
		})
	}

	req, err := utils.ConvertFromIndex[*http.Request](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	res, err := utils.ConvertFromIndex[http.ResponseWriter](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	mux.ServeHTTP(res, req)
	return nil, lastACL
}

// executeMiddlewareChain 使用洋葱模型执行中间件链和控制器
// 链结构: middleware[0].handle → middleware[1].handle → ... → controller
func executeMiddlewareChain(vm data.VM, ctx data.Context, rt runtime.Route, request data.Value, response data.Value) (data.GetValue, data.Control) {
	middlewares := rt.Middlewares

	// 如果没有中间件，直接执行控制器
	if len(middlewares) == 0 {
		return executeControllerMethod(rt, request, response, ctx)
	}

	// 构建洋葱链，从内到外
	chainIdx := 0

	var buildNext func() (data.GetValue, data.Control)
	buildNext = func() (data.GetValue, data.Control) {
		if chainIdx >= len(middlewares) {
			// 链末端：执行控制器方法
			return executeControllerMethod(rt, request, response, ctx)
		}

		mw := middlewares[chainIdx]
		chainIdx++

		cls, ok := vm.GetClass(mw.ClassName)
		if !ok {
			// 中间件类不存在，跳过，继续下一个
			return buildNext()
		}

		// 实例化中间件类
		inst, acl := cls.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}

		cv, ok := inst.(*data.ClassValue)
		if !ok {
			return buildNext()
		}

		// 获取 handle 方法
		method, has := cv.GetMethod("handle")
		if !has {
			// 没有 handle 方法，跳过
			return buildNext()
		}

		vars := method.GetVariables()
		if len(vars) < 3 {
			return buildNext()
		}

		// 创建 $next 回调，指向链中的下一个节点
		nextFunc := data.NewFuncValue(NextFunc{
			name: "next",
			fn:   func(request data.Value, response data.Value) (data.GetValue, data.Control) { return buildNext() },
			variable: []data.Variable{
				node.NewVariable(nil, "request", 0, nil),
				node.NewVariable(nil, "response", 1, nil),
			},
		})

		// 绑定参数并调用 handle($request, $response, $next)
		fnCtx := cv.CreateContext(vars)
		fnCtx.SetVariableValue(vars[0], request)
		fnCtx.SetVariableValue(vars[1], response)
		fnCtx.SetVariableValue(vars[2], nextFunc)

		return method.Call(fnCtx)
	}

	return buildNext()
}

// paramSource 参数绑定来源
type paramSource int

const (
	srcAuto        paramSource = iota // 按名称自动绑定 (string/int/float/bool/mixed)
	srcRequestObj                     // Request 对象
	srcResponseObj                    // Response 对象
	srcAllData                        // array 类型 → 全部请求数据
	srcBodyClass                      // 自定义类 → JSON body 绑定
)

// classifyParam 根据参数类型决定绑定来源
func classifyParam(v data.Variable) (paramSource, string) {
	ty := v.GetType()
	if ty == nil {
		return srcAuto, ""
	}
	typeName := ty.String()
	switch typeName {
	case "int", "string", "float", "bool", "mixed", "":
		return srcAuto, ""
	case "array":
		return srcAllData, ""
	}

	// 检查是否为 Request / Response 类型
	lower := strings.ToLower(typeName)
	if strings.Contains(lower, "request") && strings.Contains(lower, "http") {
		return srcRequestObj, ""
	}
	if strings.Contains(lower, "response") && strings.Contains(lower, "http") {
		return srcResponseObj, ""
	}

	// 检查可空类型 ?TypeName
	if len(typeName) > 1 && typeName[0] == '?' {
		inner := typeName[1:]
		switch inner {
		case "int", "string", "float", "bool", "mixed", "array":
			return srcAuto, ""
		}
		if data.ISBaseType(inner) {
			return srcAuto, ""
		}
		// 可空类类型
		return srcBodyClass, inner
	}

	// 检查联合类型 type1|type2|...
	if strings.Contains(typeName, "|") {
		// 联合类型中的 null 表示可空，按第一个非 null 类型处理
		parts := strings.Split(typeName, "|")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "null" && p != "mixed" {
				if data.ISBaseType(p) {
					return srcAuto, ""
				}
				return srcBodyClass, p
			}
		}
		return srcAuto, ""
	}

	// 其他类型视为自定义类 — 从 JSON body 绑定
	if data.ISBaseType(typeName) {
		return srcAuto, ""
	}
	return srcBodyClass, typeName
}

// ParamBind 参数绑定来源（纯 std 层定义，导出供 annotation 包使用）
type ParamBind int

const (
	BindAuto  ParamBind = iota // 自动（按名称从 POST > Query > Route 查找）
	BindPath                   // 来自路由路径参数
	BindQuery                  // 来自 URL 查询参数
	BindBody                   // 来自请求体
)

// methodParamBindings 存储方法的参数绑定映射（key: 方法的指针地址字符串）
var methodParamBindings = make(map[string][]ParamBind)

// RegisterParamBindings 由 annotation 包在扫描时调用，注册方法的参数绑定
func RegisterParamBindings(methodID string, bindings []ParamBind) {
	if len(bindings) == 0 {
		return
	}
	methodParamBindings[methodID] = bindings
}

// GetParamBindings 获取方法的参数绑定
func GetParamBindings(methodID string) []ParamBind {
	return methodParamBindings[methodID]
}

// executeControllerMethod 执行控制器方法
// 参数根据类型自动绑定（位置无关），绑定来源通过类型推断决定
func executeControllerMethod(rt runtime.Route, reqProxy data.Value, resProxy data.Value, ctx data.Context) (data.GetValue, data.Control) {
	vars := rt.Target.GetVariables()

	// 提取底层 *http.Request 用于参数绑定
	var httpReq *http.Request
	if cv, ok := reqProxy.(*data.ClassValue); ok {
		if src, ok2 := cv.Class.(data.GetSource); ok2 {
			if r, ok3 := src.GetSource().(*http.Request); ok3 {
				httpReq = r
			}
		}
	}

	// 获取方法级注解声明的参数绑定
	methodBindings := GetParamBindings(methodID(rt.Target))

	args := make([]data.Value, len(vars))
	for i, v := range vars {
		// 1. 注解声明的绑定来源优先
		if methodBindings != nil && i < len(methodBindings) {
			switch methodBindings[i] {
			case BindPath:
				args[i] = bindFromPath(httpReq, v)
				continue
			case BindQuery:
				args[i] = bindFromQuery(httpReq, v)
				continue
			case BindBody:
				args[i] = bindFromBody(httpReq, v, ctx)
				continue
			}
		}

		// 2. 类型推断
		src, className := classifyParam(v)
		switch src {
		case srcRequestObj:
			args[i] = reqProxy
		case srcResponseObj:
			args[i] = resProxy
		case srcAllData:
			args[i] = buildAllRequestData(httpReq)
		case srcBodyClass:
			args[i] = bindBodyToClass(httpReq, className, ctx)
		default:
			args[i] = bindRequestParam(httpReq, v)
		}
	}

	if rt.Receiver != nil {
		return node.CallHTTPControllerMethod(rt.Receiver, rt.Target, args)
	}

	mute := ctx.CreateContext(vars)
	for i, arg := range args {
		if i < len(vars) {
			mute.SetVariableValue(vars[i], arg)
		}
	}
	return rt.Target.Call(mute)
}

// buildAllRequestData 合并所有请求数据（路由参数 + 查询参数 + POST 表单）
func buildAllRequestData(r *http.Request) data.Value {
	if r == nil {
		return data.NewObjectValue()
	}
	result := data.NewObjectValue()

	// 路由参数
	for key, val := range collectPathValues(r) {
		result.SetProperty(key, data.NewStringValue(val))
	}
	// 查询参数
	for key, vals := range r.URL.Query() {
		if len(vals) == 1 {
			result.SetProperty(key, data.NewStringValue(vals[0]))
		} else {
			result.SetProperty(key, data.NewStringValue(strings.Join(vals, ",")))
		}
	}
	// POST 表单（最高优先级）
	if r.PostForm != nil {
		for key, vals := range r.PostForm {
			if len(vals) == 1 {
				result.SetProperty(key, data.NewStringValue(vals[0]))
			} else {
				result.SetProperty(key, data.NewStringValue(strings.Join(vals, ",")))
			}
		}
	}
	return result
}

// bindBodyToClass 从 JSON body 绑定到指定类，也支持表单数据
func bindBodyToClass(r *http.Request, className string, ctx data.Context) data.Value {
	if r == nil || className == "" {
		return data.NewNullValue()
	}

	vm := ctx.GetVM()
	if vm == nil {
		return data.NewNullValue()
	}

	classStmt, acl := vm.GetOrLoadClass(className)
	if acl != nil || classStmt == nil {
		return data.NewNullValue()
	}

	classInst, _ := classStmt.GetValue(ctx)
	cv, ok := classInst.(*data.ClassValue)
	if !ok {
		return data.NewNullValue()
	}

	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		// JSON body
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil || len(body) == 0 {
				return cv
			}
			// 将 JSON 数据设置到类属性上
			var rawMap map[string]interface{}
			if err := json.Unmarshal(body, &rawMap); err == nil {
				for key, val := range rawMap {
					cv.SetProperty(key, convertGoToDataValue(val))
				}
			}
		}
	} else {
		// 表单数据 → 类属性
		if r.PostForm != nil {
			for key, vals := range r.PostForm {
				if len(vals) > 0 {
					cv.SetProperty(key, data.NewStringValue(vals[0]))
				}
			}
		}
		// 也加入查询参数
		for key, vals := range r.URL.Query() {
			if len(vals) > 0 {
				cv.SetProperty(key, data.NewStringValue(vals[0]))
			}
		}
		// 也加入路由参数
		for key, val := range collectPathValues(r) {
			cv.SetProperty(key, data.NewStringValue(val))
		}
	}

	return cv
}

// convertGoToDataValue 将 Go 值转换为 data.Value
func convertGoToDataValue(v interface{}) data.Value {
	switch val := v.(type) {
	case map[string]interface{}:
		obj := data.NewObjectValue()
		for k, v := range val {
			obj.SetProperty(k, convertGoToDataValue(v))
		}
		return obj
	case []interface{}:
		values := make([]data.Value, len(val))
		for i, v := range val {
			values[i] = convertGoToDataValue(v)
		}
		return data.NewArrayValue(values)
	case string:
		return data.NewStringValue(val)
	case float64:
		return data.NewFloatValue(val)
	case bool:
		return data.NewBoolValue(val)
	case nil:
		return data.NewNullValue()
	default:
		return data.NewStringValue("")
	}
}

// bindRequestParam 从请求数据中查找值并转换为方法参数所需的类型
// 查找优先级：POST 表单 > URL 查询参数 > 路由参数
func bindRequestParam(r *http.Request, v data.Variable) data.Value {
	name := v.GetName()
	if name == "" || r == nil {
		return data.NewNullValue()
	}

	var raw string

	// 1. 从 POST 表单获取（最高优先级）
	if r.PostForm != nil {
		if vals, exists := r.PostForm[name]; exists && len(vals) > 0 {
			raw = vals[0]
		}
	}

	// 2. 从查询参数获取
	if raw == "" {
		if vals, exists := r.URL.Query()[name]; exists && len(vals) > 0 {
			raw = vals[0]
		}
	}

	// 3. 从路由参数获取（最低优先级）
	if raw == "" {
		raw = r.PathValue(name)
	}

	// 4. 未找到，返回 null
	if raw == "" {
		return data.NewNullValue()
	}

	// 5. 根据参数声明的类型进行转换
	return convertParamValue(raw, v.GetType())
}

// convertParamValue 将字符串值转换为目标类型，支持可空类型
func convertParamValue(raw string, ty data.Types) data.Value {
	if ty == nil {
		return data.NewStringValue(raw)
	}

	typeName := ty.String()

	// 处理可空类型 ?type
	if len(typeName) > 1 && typeName[0] == '?' {
		typeName = typeName[1:]
	}

	// 处理联合类型 type|null 或 null|type
	if strings.Contains(typeName, "|") {
		parts := strings.Split(typeName, "|")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "null" && p != "mixed" {
				typeName = p
				break
			}
		}
	}

	switch typeName {
	case "int":
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return data.NewIntValue(int(n))
		}
		return data.NewIntValue(0)
	case "float":
		if f, err := strconv.ParseFloat(raw, 64); err == nil {
			return data.NewFloatValue(f)
		}
		return data.NewFloatValue(0)
	case "bool", "false":
		switch raw {
		case "1", "true", "True", "TRUE", "yes", "Yes", "YES", "on", "On", "ON":
			return data.NewBoolValue(true)
		default:
			return data.NewBoolValue(false)
		}
	case "string", "mixed", "":
		return data.NewStringValue(raw)
	default:
		return data.NewStringValue(raw)
	}
}

// MethodID 为方法生成唯一标识符（导出供 annotation 包使用）
func MethodID(m data.Method) string {
	return fmt.Sprintf("%p", m)
}

// methodID 内部使用的别名
func methodID(m data.Method) string {
	return MethodID(m)
}

// bindFromPath 仅从路由路径参数绑定
func bindFromPath(r *http.Request, v data.Variable) data.Value {
	name := v.GetName()
	if name == "" || r == nil {
		return data.NewNullValue()
	}
	raw := r.PathValue(name)
	if raw == "" {
		return data.NewNullValue()
	}
	return convertParamValue(raw, v.GetType())
}

// bindFromQuery 仅从 URL 查询参数绑定
func bindFromQuery(r *http.Request, v data.Variable) data.Value {
	name := v.GetName()
	if name == "" || r == nil {
		return data.NewNullValue()
	}
	vals, exists := r.URL.Query()[name]
	if !exists || len(vals) == 0 || vals[0] == "" {
		return data.NewNullValue()
	}
	return convertParamValue(vals[0], v.GetType())
}

// bindFromBody 从请求体绑定（POST 表单或 JSON）
func bindFromBody(r *http.Request, v data.Variable, ctx data.Context) data.Value {
	name := v.GetName()
	if name == "" || r == nil {
		return data.NewNullValue()
	}

	// 先尝试 POST 表单
	if r.PostForm != nil {
		if vals, exists := r.PostForm[name]; exists && len(vals) > 0 {
			return convertParamValue(vals[0], v.GetType())
		}
	}

	ty := v.GetType()
	if ty == nil {
		return data.NewNullValue()
	}

	// array 或自定义类 → 从请求体绑定整个对象
	typeName := ty.String()
	if typeName == "array" {
		return buildAllRequestData(r)
	}
	if !data.ISBaseType(typeName) {
		return bindBodyToClass(r, typeName, ctx)
	}

	return data.NewNullValue()
}
