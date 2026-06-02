package http

import (
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// AppFlashFunction 性能优先：无实例状态，仅在全局 VM 尚无注解路由时加载 $dir/main.php，随后直接分发。
type AppFlashFunction struct{}

func NewAppFlashFunction() data.FuncStmt { return &AppFlashFunction{} }

func (h *AppFlashFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if _, ok := ctx.GetIndexValue(0); !ok {
		return nil, utils.NewThrowf("缺少参数: request")
	}
	if _, ok := ctx.GetIndexValue(1); !ok {
		return nil, utils.NewThrowf("缺少参数: response")
	}

	vm := ctx.GetVM()
	if len(runtime.HTTPRoutes(vm)) == 0 {
		dir, acl := resolveAppDirPath(ctx, 2, "./src")
		if acl != nil {
			return nil, acl
		}
		mainFile := filepath.Join(dir, "main.php")
		if _, acl = vm.LoadAndRun(mainFile); acl != nil {
			return nil, acl
		}
	}

	vars := []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
	}
	dispatchCtx := vm.CreateContext(vars)
	v0, _ := ctx.GetIndexValue(0)
	v1, _ := ctx.GetIndexValue(1)
	dispatchCtx.SetVariableValue(vars[0], v0)
	dispatchCtx.SetVariableValue(vars[1], v1)

	return DispatchHTTPRoutes(vm, dispatchCtx)
}

func (h *AppFlashFunction) GetName() string            { return "Net\\Http\\app_flash" }
func (h *AppFlashFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *AppFlashFunction) GetIsStatic() bool          { return true }
func (h *AppFlashFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
		node.NewParameter(nil, "dir", 2, data.NewStringValue("./src"), data.String{}),
	}
}
func (h *AppFlashFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
		node.NewVariable(nil, "dir", 2, nil),
	}
}
func (h *AppFlashFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }
