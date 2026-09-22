package foundation

import "github.com/php-any/origami/data"

// BindRequest 把当前 HTTP Request 实例绑到 Application 容器。
func BindRequest(app *data.ClassValue, request data.Value) data.Control {
	if app == nil || request == nil {
		return nil
	}
	for _, abstract := range []string{"Illuminate\\Http\\Request", "request"} {
		if m, ok := app.GetMethod("instance"); ok && m != nil {
			nctx := app.CreateContext(m.GetVariables())
			data.BindDeclaredArgs(nctx, m, []data.Value{data.NewStringValue(abstract), request})
			if _, ctl := m.Call(nctx); ctl != nil {
				return ctl
			}
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
