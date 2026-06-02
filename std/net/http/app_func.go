package http

import (
	"errors"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

type AppFunction struct {
	// runMain 按解析后的绝对 filePath 缓存（hotReload=false 时每个路径仅 LoadAndRun 一次）
	runMain sync.Map
}

type mainLoadEntry struct {
	once sync.Once
	acl  data.Control
}

func NewAppFunction() data.FuncStmt { return &AppFunction{} }

func (h *AppFunction) loadMainFile(vm data.VM, filePath string, hotReload bool) data.Control {
	if hotReload {
		_, acl := vm.LoadAndRun(filePath)
		return acl
	}
	v, _ := h.runMain.LoadOrStore(filePath, &mainLoadEntry{})
	entry := v.(*mainLoadEntry)
	entry.once.Do(func() {
		_, entry.acl = vm.LoadAndRun(filePath)
	})
	return entry.acl
}

func (h *AppFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 request 和 response 参数
	requestValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: request"))
	}

	responseValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: response"))
	}

	filePath, acl := resolveAppFilePath(ctx, 2, "./src/main.php")
	if acl != nil {
		return nil, acl
	}

	hotReload := true
	if hotReloadArg, ok := ctx.GetIndexValue(4); ok && hotReloadArg != nil {
		if bv, ok := hotReloadArg.(data.AsBool); ok {
			hotReload, _ = bv.AsBool()
		}
	}

	baseVM := ctx.GetVM()
	vm := baseVM
	if hotReload {
		vm = runtime.NewTempVM(baseVM)
	}
	if acl := h.loadMainFile(vm, filePath, hotReload); acl != nil {
		return nil, acl
	}

	// 查找 app 函数
	fullName, err := utils.ConvertFromIndex[string](ctx, 3)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	fn, ok := vm.GetFunc(fullName)
	if !ok {
		return nil, utils.NewThrowf("未找到%s($request, $response) 函数", fullName)
	}

	// 调用 app($request, $response) - 参考 NextHandler 的简洁方式
	vars := fn.GetVariables()
	if len(vars) < 2 {
		return nil, utils.NewThrow(errors.New("app 函数需要至少 2 个参数: $request, $response"))
	}

	fnCtx := vm.CreateContext(vars)
	fnCtx.SetVariableValue(vars[0], requestValue)
	fnCtx.SetVariableValue(vars[1], responseValue)

	return fn.Call(fnCtx)
}

func (h *AppFunction) GetName() string            { return "Net\\Http\\app" }
func (h *AppFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *AppFunction) GetIsStatic() bool          { return true }
func (h *AppFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
		node.NewParameter(nil, "filePath", 2, data.NewStringValue("./src/main.php"), data.String{}),
		node.NewParameter(nil, "fun", 3, data.NewStringValue("App\\main"), data.String{}),
		node.NewParameter(nil, "hotReload", 4, data.NewBoolValue(true), data.NewBaseType("bool")),
	}
}
func (h *AppFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
		node.NewVariable(nil, "filePath", 2, nil),
		node.NewVariable(nil, "fun", 3, nil),
		node.NewVariable(nil, "hotReload", 4, nil),
	}
}
func (h *AppFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }
