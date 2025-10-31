package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// 内联引导：扫描控制器并以 serveHTTP 方式分发当前请求
type springInlineFunc struct{}

func newSpringInlineFunc() data.FuncStmt { return &springInlineFunc{} }

func (f *springInlineFunc) GetName() string            { return "Annotation\\__internal_spring_inline" }
func (f *springInlineFunc) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *springInlineFunc) GetIsStatic() bool          { return true }
func (f *springInlineFunc) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
	}
}
func (f *springInlineFunc) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
	}
}
func (f *springInlineFunc) GetReturnType() data.Types { return data.NewBaseType("void") }

func (f *springInlineFunc) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 直接复用脚本层引导：通过 include("examples/spring/main.zy") 中的 __spring_bootstrap 实现
	// 这里为简化，仅调用 Net\Http\app 以当前工作目录的 main.zy 作为分发入口
	// 实际项目可将扫描逻辑迁移到 Go 层以提升性能
	app, ok := ctx.GetVM().GetFunc("Net\\Http\\app")
	if !ok {
		return nil, utils.NewThrow(data.NewError(nil, "缺少 Net\\Http\\app", nil))
	}
	// 传递 request/response 即可，让 app 加载当前目录 main.zy 并执行
	vars := app.GetVariables()
	fnCtx := ctx.GetVM().CreateContext(vars)
	v0, _ := ctx.GetIndexValue(0)
	v1, _ := ctx.GetIndexValue(1)
	fnCtx.SetVariableValue(vars[0], v0)
	fnCtx.SetVariableValue(vars[1], v1)
	return app.Call(fnCtx)
}
