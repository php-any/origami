package fpm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

// RequestVM 为经典 PHP 入口提供请求级隔离：函数/类/常量/文件缓存本地化，
// 输出与 HTTP 绑定到当前请求，避免污染共享 base VM。
// 每个 HTTP 请求独占一个实例、单 goroutine，无需额外锁。
type RequestVM struct {
	base   *runtime.VM
	output data.OutputWriter
	ob     obStack

	request   *http.Request
	response  http.ResponseWriter
	lastThrow data.Control

	classes            map[string]data.ClassStmt
	interfaces         map[string]data.InterfaceStmt
	funcs              map[string]data.FuncStmt
	constants          map[string]data.Value
	phpFileCache       map[string]struct{}
	includeOnceResults map[string]data.GetValue
	globalVars         map[string]*data.ZVal
	globalsArray       *data.ObjectValue
	sessionArray       *data.ObjectValue
	callDepth          int
	callStack          []data.CallFrame
}

// New 创建请求级 VM。output 为本请求 HTTP body 写入目标；nil 时使用 data.DefaultOutputWriter。
func New(base *runtime.VM, output data.OutputWriter) *RequestVM {
	if output == nil {
		output = data.DefaultOutputWriter
	}
	return &RequestVM{
		base:               base,
		output:             output,
		classes:            make(map[string]data.ClassStmt),
		interfaces:         make(map[string]data.InterfaceStmt),
		funcs:              make(map[string]data.FuncStmt),
		constants:          make(map[string]data.Value),
		phpFileCache:       make(map[string]struct{}),
		includeOnceResults: make(map[string]data.GetValue),
		globalVars:         make(map[string]*data.ZVal),
	}
}

// BindHTTP 绑定本请求的 Go Request/ResponseWriter。
func (v *RequestVM) BindHTTP(req *http.Request, resp http.ResponseWriter) {
	v.request = req
	v.response = resp
}

// HTTPRequest 返回当前绑定的 HTTP 请求。
func (v *RequestVM) HTTPRequest() *http.Request {
	return v.request
}

// HTTPResponseWriter 返回当前绑定的 HTTP 响应写入器。
func (v *RequestVM) HTTPResponseWriter() http.ResponseWriter {
	return v.response
}

func (v *RequestVM) WriteOutput(s string) {
	if v.ob.write(s) {
		return
	}
	v.output(s)
}

func (v *RequestVM) StartOutputBuffer() {
	v.ob.push()
}

func (v *RequestVM) CleanOutputBuffer() (string, bool) {
	if v.ob.level() == 0 {
		return "", false
	}
	return v.ob.pop(), true
}

func (v *RequestVM) OutputBufferContents() (string, bool) {
	if v.ob.level() == 0 {
		return "", false
	}
	return v.ob.contents(), true
}

func (v *RequestVM) OutputBufferLevel() int {
	return v.ob.level()
}

func (v *RequestVM) AddClass(c data.ClassStmt) data.Control {
	if _, ok := v.classes[c.GetName()]; ok {
		return data.NewErrorThrow(nil, fmt.Errorf("已存在同名的 class: %s", c.GetName()))
	}
	v.classes[c.GetName()] = c
	return nil
}

func (v *RequestVM) GetClass(pkg string) (data.ClassStmt, bool) {
	if c, ok := v.classes[pkg]; ok {
		return c, true
	}
	return v.base.GetClass(pkg)
}

func (v *RequestVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	if c, ok := v.GetClass(pkg); ok {
		return c, nil
	}
	if name := strings.TrimPrefix(pkg, "\\"); name != "" {
		ok, acl := parser.CallAutoLoad(name, v.CreateContext(nil))
		if acl != nil {
			return nil, acl
		}
		if ok {
			if c, found := v.GetClass(name); found {
				return c, nil
			}
		}
	}
	return v.base.GetOrLoadClass(pkg)
}

func (v *RequestVM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	return v.base.LoadPkg(pkg)
}

func (v *RequestVM) AddInterface(i data.InterfaceStmt) data.Control {
	if _, ok := v.interfaces[i.GetName()]; ok {
		return data.NewErrorThrow(nil, fmt.Errorf("已存在同名的 interface: %s", i.GetName()))
	}
	v.interfaces[i.GetName()] = i
	return nil
}

func (v *RequestVM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	if i, ok := v.interfaces[pkg]; ok {
		return i, true
	}
	return v.base.GetInterface(pkg)
}

func (v *RequestVM) GetOrLoadInterface(pkg string) (data.InterfaceStmt, data.Control) {
	if i, ok := v.GetInterface(pkg); ok {
		return i, nil
	}
	return v.base.GetOrLoadInterface(pkg)
}

func (v *RequestVM) AddFunc(f data.FuncStmt) data.Control {
	if _, ok := v.funcs[f.GetName()]; ok {
		return data.NewErrorThrow(nil, fmt.Errorf("已存在同名的 function: %s", f.GetName()))
	}
	v.funcs[f.GetName()] = f
	return nil
}

func (v *RequestVM) GetFunc(pkg string) (data.FuncStmt, bool) {
	if f, ok := v.funcs[pkg]; ok {
		return f, true
	}
	return v.base.GetFunc(pkg)
}

func (v *RequestVM) RegisterFunction(name string, fn interface{}) data.Control {
	return v.base.RegisterFunction(name, fn)
}

func (v *RequestVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	return v.base.RegisterReflectClass(name, instance)
}

func (v *RequestVM) CreateContext(vars []data.Variable) data.Context {
	ctx := v.base.CreateContext(vars)
	if rctx, ok := ctx.(interface{ SetVM(data.VM) }); ok {
		rctx.SetVM(v)
	}
	return ctx
}

func (v *RequestVM) SetThrowControl(fn func(acl data.Control)) {}

func (v *RequestVM) ThrowControl(acl data.Control) {
	v.lastThrow = acl
}

func (v *RequestVM) takeThrow() data.Control {
	acl := v.lastThrow
	v.lastThrow = nil
	return acl
}

// TakeThrow 取出并清空本请求暂存的 throw 控制流。
func (v *RequestVM) TakeThrow() data.Control {
	return v.takeThrow()
}

func (v *RequestVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	return v.runFile(file, false)
}

// LoadAndRunFresh 强制重新读盘并解析（热重载/调试）；正常 HTTP 请求应使用 LoadAndRun。
func (v *RequestVM) LoadAndRunFresh(file string) (data.GetValue, data.Control) {
	return v.runFile(file, true)
}

func (v *RequestVM) runFile(file string, fresh bool) (data.GetValue, data.Control) {
	_ = v.takeThrow()

	var program data.GetValue
	var vars []data.Variable
	var acl data.Control

	if fresh {
		p := v.base.GetParser().Clone()
		p.SetVM(v)
		program, acl = p.ParseFile(file)
		if acl != nil {
			return nil, acl
		}
		vars = p.GetVariables()
	} else {
		program, vars, acl = v.base.ParseFileCached(file)
		if acl != nil {
			return nil, acl
		}
	}

	ctx := v.CreateContext(vars)
	v.registerGlobalContext(vars, ctx)
	result, ctrl := program.GetValue(ctx)
	if ctrl == nil {
		if thrown := v.takeThrow(); thrown != nil {
			return result, thrown
		}
	}
	return result, ctrl
}

func (v *RequestVM) LoadInCallerContext(parent data.Context, file string) (data.GetValue, data.Control) {
	program, vars, acl := v.base.ParseFileCached(file)
	if acl != nil {
		return nil, acl
	}
	ctx := v.CreateContext(vars)
	for _, variable := range vars {
		name := variable.GetName()
		if name == "" {
			continue
		}
		if val, ok := parent.GetVariableByName(name); ok && val != nil {
			if ctl := variable.SetValue(ctx, val); ctl != nil {
				return nil, ctl
			}
			continue
		}
		ctx.SetIndexZVal(variable.GetIndex(), v.EnsureGlobalZVal(name))
	}
	return program.GetValue(ctx)
}

func (v *RequestVM) CompileLoad(file string) data.Control {
	return v.base.CompileLoad(file)
}

func (v *RequestVM) RegisterCompiledFile(file string, fn func() (data.GetValue, []data.Variable)) {
	v.base.RegisterCompiledFile(file, fn)
}

func (v *RequestVM) RunCompiledFile(file string) (data.GetValue, data.Control) {
	return v.base.RunCompiledFile(file)
}

func (v *RequestVM) ParseFile(file string, value data.Value) (data.Value, data.Control) {
	return v.base.ParseFile(file, value)
}

func (v *RequestVM) SetPhpFileCache(file string) {
	v.phpFileCache[file] = struct{}{}
}

func (v *RequestVM) GetPhpFileCache(file string) bool {
	_, ok := v.phpFileCache[file]
	return ok
}

func (v *RequestVM) GetIncludeOnceResult(file string) (data.GetValue, bool) {
	result, ok := v.includeOnceResults[file]
	return result, ok
}

func (v *RequestVM) SetIncludeOnceResult(file string, result data.GetValue) {
	v.includeOnceResults[file] = result
}

func (v *RequestVM) registerGlobalContext(vars []data.Variable, ctx data.Context) {
	for _, variable := range vars {
		if variable == nil {
			continue
		}
		name := variable.GetName()
		if name == "" {
			continue
		}
		if _, exists := v.globalVars[name]; exists {
			continue
		}
		if zv := ctx.GetIndexZVal(variable.GetIndex()); zv != nil {
			v.globalVars[name] = zv
		}
	}
}

func (v *RequestVM) AddNamespace(namespace string, path string) {
	v.base.AddNamespace(namespace, path)
}

func (v *RequestVM) EnterCall() int {
	v.callDepth++
	return v.callDepth
}

func (v *RequestVM) LeaveCall() {
	if v.callDepth > 0 {
		v.callDepth--
	}
}

func (v *RequestVM) SetConstant(name string, value data.Value) data.Control {
	if _, ok := v.constants[name]; ok {
		return data.NewErrorThrow(nil, fmt.Errorf("常量 %s 已经定义，不能重新定义", name))
	}
	v.constants[name] = value
	return nil
}

func (v *RequestVM) GetConstant(name string) (data.Value, bool) {
	if c, ok := v.constants[name]; ok {
		return c, true
	}
	return v.base.GetConstant(name)
}

// EnsureGlobalZVal 返回请求级全局变量槽位；不存在时创建为 null。
func (v *RequestVM) EnsureGlobalZVal(name string) *data.ZVal {
	if zv, ok := v.globalVars[name]; ok {
		return zv
	}
	zv := data.NewZVal(data.NewNullValue())
	v.globalVars[name] = zv
	return zv
}

// EnsureGlobalsArray 返回本请求的 $GLOBALS 数组（不跨请求共享）。
func (v *RequestVM) EnsureGlobalsArray() *data.ObjectValue {
	if v.globalsArray == nil {
		v.globalsArray = data.NewObjectValue()
	}
	return v.globalsArray
}

// EnsureSessionArray 返回本请求的 $_SESSION 数组（不跨请求共享）。
func (v *RequestVM) EnsureSessionArray() *data.ObjectValue {
	if v.sessionArray == nil {
		v.sessionArray = data.NewObjectValue()
	}
	return v.sessionArray
}

func (v *RequestVM) SetExceptionHandler(handler data.Value) data.Value {
	return v.base.SetExceptionHandler(handler)
}

func (v *RequestVM) GetExceptionHandler() data.Value {
	return v.base.GetExceptionHandler()
}

func (v *RequestVM) SetErrorHandler(handler data.Value) data.Value {
	return v.base.SetErrorHandler(handler)
}

func (v *RequestVM) RestoreErrorHandler() bool {
	return v.base.RestoreErrorHandler()
}

func (v *RequestVM) GetErrorHandler() data.Value {
	return v.base.GetErrorHandler()
}

func (v *RequestVM) AddShutdownCallback(cb data.Value) {
	v.base.AddShutdownCallback(cb)
}

func (v *RequestVM) RunShutdownCallbacks() {
	v.base.RunShutdownCallbacks()
}

func (v *RequestVM) PushCallFrame(frame data.CallFrame) {
	v.callStack = append(v.callStack, frame)
}

func (v *RequestVM) PopCallFrame() {
	if n := len(v.callStack); n > 0 {
		v.callStack = v.callStack[:n-1]
	}
}

func (v *RequestVM) SnapshotCallStack() []data.CallFrame {
	if len(v.callStack) == 0 {
		return nil
	}
	out := make([]data.CallFrame, len(v.callStack))
	copy(out, v.callStack)
	return out
}
