package wails

import (
	"github.com/php-any/origami/data"
)

// ============================================================================
// Wails\Application — Wails 应用程序入口
//
// 静态方法类，用于启动 Wails 桌面应用程序。
//
// 用法示例 (PHP):
//
//	$app = new Wails\Options\App([
//	    'Title' => 'My Desktop App',
//	    'Width' => 1024,
//	    'Height' => 768,
//	]);
//	$app->onDomReady(function() {
//	    Wails\Runtime\Window::center();
//	});
//	Wails\Application::run($app);
// ============================================================================

type ApplicationClass struct{}

func NewApplicationClass() data.ClassStmt { return &ApplicationClass{} }

func (c *ApplicationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ApplicationClass) GetFrom() data.From                            { return nil }
func (c *ApplicationClass) GetName() string                               { return "Wails\\Application" }
func (c *ApplicationClass) GetExtend() *string                            { return nil }
func (c *ApplicationClass) GetImplements() []string                       { return nil }
func (c *ApplicationClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ApplicationClass) GetPropertyList() []data.Property              { return nil }
func (c *ApplicationClass) GetConstruct() data.Method                     { return nil }

var applicationStaticMethods = map[string]data.Method{
	"run":  &appRunMethod{},
	"quit": &appQuitMethod{},
	"hide": &appHideMethod{},
	"show": &appShowMethod{},
}

func (c *ApplicationClass) GetMethod(name string) (data.Method, bool) {
	m, ok := applicationStaticMethods[name]
	return m, ok
}
func (c *ApplicationClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := applicationStaticMethods[name]
	return m, ok
}
func (c *ApplicationClass) GetMethods() []data.Method {
	return []data.Method{
		&appRunMethod{},
		&appQuitMethod{},
		&appHideMethod{},
		&appShowMethod{},
	}
}

// ====== run ======

type appRunMethod struct{}

func (m *appRunMethod) GetName() string            { return "run" }
func (m *appRunMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *appRunMethod) GetIsStatic() bool          { return true }
func (m *appRunMethod) GetReturnType() data.Types  { return nil }

func (m *appRunMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("options", 0),
	}
}
func (m *appRunMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("options", 0, data.NewBaseType("Wails\\Options\\App")),
	}
}

func (m *appRunMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok || v == nil {
		return nil, nil
	}
	// 调用 Wails v3 的 Run — 这会阻塞直到应用退出
	if err := RunApp(ctx, v); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

// ====== quit ======

type appQuitMethod struct{}

func (m *appQuitMethod) GetName() string               { return "quit" }
func (m *appQuitMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *appQuitMethod) GetIsStatic() bool             { return true }
func (m *appQuitMethod) GetReturnType() data.Types     { return nil }
func (m *appQuitMethod) GetParams() []data.GetValue    { return nil }
func (m *appQuitMethod) GetVariables() []data.Variable { return nil }

func (m *appQuitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsApp != nil {
		wailsApp.Quit()
	}
	return nil, nil
}

// ====== hide ======

type appHideMethod struct{}

func (m *appHideMethod) GetName() string               { return "hide" }
func (m *appHideMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *appHideMethod) GetIsStatic() bool             { return true }
func (m *appHideMethod) GetReturnType() data.Types     { return nil }
func (m *appHideMethod) GetParams() []data.GetValue    { return nil }
func (m *appHideMethod) GetVariables() []data.Variable { return nil }

func (m *appHideMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Hide()
	}
	return nil, nil
}

// ====== show ======

type appShowMethod struct{}

func (m *appShowMethod) GetName() string               { return "show" }
func (m *appShowMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *appShowMethod) GetIsStatic() bool             { return true }
func (m *appShowMethod) GetReturnType() data.Types     { return nil }
func (m *appShowMethod) GetParams() []data.GetValue    { return nil }
func (m *appShowMethod) GetVariables() []data.Variable { return nil }

func (m *appShowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Show()
	}
	return nil, nil
}
