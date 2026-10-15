package foundation

import "github.com/php-any/origami/data"

// BindRequest 把当前 HTTP Request 实例绑到 Application 容器。
//
// 只绑 'request'：Illuminate\Http\Request 已由 Application::registerCoreContainerAliases
// 注册为 'request' 的别名，解析它会经 getAlias 落到同一个实例上。多绑一次
// instance("Illuminate\\Http\\Request", ...) 反而会 unset 掉该别名，并让
// rebound('request') 回调与 removeAbstractAlias 扫描多跑一遍。
func BindRequest(app *data.ClassValue, request data.Value) data.Control {
	if app == nil || request == nil {
		return nil
	}
	if m, ok := app.GetMethod("instance"); ok && m != nil {
		nctx := app.CreateContext(m.GetVariables())
		if ctl := data.BindDeclaredArgs(nctx, m, []data.Value{data.NewStringValue("request"), request}); ctl != nil {
			return ctl
		}
		if _, ctl := m.Call(nctx); ctl != nil {
			return ctl
		}
	}
	return nil
}

// Load 注册 illuminate/foundation 加速件（Application 本体仍走 vendor）。
func Load(vm data.VM) {
	_ = vm
}

const ComposerName = "illuminate/foundation"
const TargetVersion = "v13.23.0"
