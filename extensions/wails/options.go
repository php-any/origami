package wails

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ============================================================================
// Wails\Options\App — Wails 主配置类
//
// 对应 Wails v3 的 application.Options + application.WebviewWindowOptions
//
// 用法示例 (PHP):
//
//	$app = new Wails\Options\App([
//	    'Title'  => 'My App',
//	    'Width'  => 1024,
//	    'Height' => 768,
//	    'Bind'   => [$myService],
//	]);
//	$app->onDomReady(function() {
//	    echo "App is ready!\n";
//	});
// ============================================================================

type OptionsAppClass struct{}

func NewOptionsAppClass() data.ClassStmt { return &OptionsAppClass{} }

func (c *OptionsAppClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *OptionsAppClass) GetFrom() data.From                            { return nil }
func (c *OptionsAppClass) GetName() string                               { return "Wails\\Options\\App" }
func (c *OptionsAppClass) GetExtend() *string                            { return nil }
func (c *OptionsAppClass) GetImplements() []string                       { return nil }
func (c *OptionsAppClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *OptionsAppClass) GetPropertyList() []data.Property              { return nil }
func (c *OptionsAppClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *OptionsAppClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "onStartup":
		return &optionsAppOnStartupMethod{}, true
	case "onDomReady":
		return &optionsAppOnDomReadyMethod{}, true
	case "onShutdown":
		return &optionsAppOnShutdownMethod{}, true
	case "onBeforeClose":
		return &optionsAppOnBeforeCloseMethod{}, true
	case "bind":
		return &optionsAppBindMethod{}, true
	case "setMenu":
		return &optionsAppSetMenuMethod{}, true
	case "bindingsAllowedOrigins":
		return &optionsAppBindingsAllowedOriginsMethod{}, true
	case "setErrorFormatter":
		return &optionsAppSetErrorFormatterMethod{}, true
	default:
		return nil, false
	}
}

func (c *OptionsAppClass) GetMethods() []data.Method {
	return []data.Method{
		&optionsAppOnStartupMethod{},
		&optionsAppOnDomReadyMethod{},
		&optionsAppOnShutdownMethod{},
		&optionsAppOnBeforeCloseMethod{},
		&optionsAppBindMethod{},
		&optionsAppSetMenuMethod{},
		&optionsAppBindingsAllowedOriginsMethod{},
		&optionsAppSetErrorFormatterMethod{},
	}
}
func (c *OptionsAppClass) GetConstruct() data.Method { return &optionsAppConstruct{} }

// ====== 构造函数 ======

type optionsAppConstruct struct{}

func (m *optionsAppConstruct) GetName() string            { return token.ConstructName }
func (m *optionsAppConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppConstruct) GetIsStatic() bool          { return false }
func (m *optionsAppConstruct) GetReturnType() data.Types  { return nil }

func (m *optionsAppConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *optionsAppConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *optionsAppConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}

	// 设置所有默认值 (映射到 WebviewWindowOptions)
	setDefaultStringProperty(cv, "Title", "")
	setDefaultIntProperty(cv, "Width", 1024)
	setDefaultIntProperty(cv, "Height", 768)
	setDefaultBoolProperty(cv, "DisableResize", false)
	setDefaultBoolProperty(cv, "Frameless", false)
	setDefaultIntProperty(cv, "MinWidth", 0)
	setDefaultIntProperty(cv, "MinHeight", 0)
	setDefaultIntProperty(cv, "MaxWidth", 0)
	setDefaultIntProperty(cv, "MaxHeight", 0)
	setDefaultBoolProperty(cv, "StartHidden", false)
	setDefaultBoolProperty(cv, "HideWindowOnClose", false)
	setDefaultBoolProperty(cv, "AlwaysOnTop", false)
	setDefaultStringProperty(cv, "CSSDragProperty", "--wails-draggable")
	setDefaultStringProperty(cv, "CSSDragValue", "drag")
	setDefaultBoolProperty(cv, "EnableDefaultContextMenu", false)
	setDefaultIntProperty(cv, "WindowStartState", 0) // Normal

	// 存储回调占位
	cv.SetProperty("_onStartup", data.NewArrayValue(nil))
	cv.SetProperty("_onDomReady", data.NewArrayValue(nil))
	cv.SetProperty("_onShutdown", data.NewArrayValue(nil))
	cv.SetProperty("_onBeforeClose", data.NewArrayValue(nil))
	cv.SetProperty("_bind", data.NewArrayValue(nil))
	cv.SetProperty("_menu", data.NewNullValue())
	cv.SetProperty("_bindingsAllowedOrigins", data.NewStringValue(""))
	cv.SetProperty("_errorFormatter", data.NewNullValue())

	// 从传入的数组选项中读取并覆盖
	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"Title", "Width", "Height", "DisableResize", "Frameless",
				"MinWidth", "MinHeight", "MaxWidth", "MaxHeight",
				"StartHidden", "HideWindowOnClose", "AlwaysOnTop",
				"CSSDragProperty", "CSSDragValue",
				"EnableDefaultContextMenu",
				"WindowStartState",
			})
			// 特殊处理嵌套对象
			for _, key := range []string{
				"BackgroundColour", "SingleInstanceLock", "DragAndDrop",
				"Windows", "Mac", "Linux", "AssetServer", "Debug",
			} {
				if v, ok := arrayGet(av, key); ok {
					cv.SetProperty(key, v)
				}
			}
			// _bind 和 _menu 使用不同属性名
			if v, ok := arrayGet(av, "Bind"); ok {
				cv.SetProperty("_bind", v)
			}
			if v, ok := arrayGet(av, "Menu"); ok {
				cv.SetProperty("_menu", v)
			}
			if v, ok := arrayGet(av, "BindingsAllowedOrigins"); ok {
				cv.SetProperty("_bindingsAllowedOrigins", v)
			}
		}
	}
	return nil, nil
}

// ====== onStartup 方法 ======

type optionsAppOnStartupMethod struct{}

func (m *optionsAppOnStartupMethod) GetName() string            { return "onStartup" }
func (m *optionsAppOnStartupMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppOnStartupMethod) GetIsStatic() bool          { return false }
func (m *optionsAppOnStartupMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppOnStartupMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callback", 0, nil, data.NewBaseType("callable"))}
}
func (m *optionsAppOnStartupMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callback", 0, data.NewBaseType("callable"))}
}
func (m *optionsAppOnStartupMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_onStartup", v)
	}
	return nil, nil
}

// ====== onDomReady 方法 ======

type optionsAppOnDomReadyMethod struct{}

func (m *optionsAppOnDomReadyMethod) GetName() string            { return "onDomReady" }
func (m *optionsAppOnDomReadyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppOnDomReadyMethod) GetIsStatic() bool          { return false }
func (m *optionsAppOnDomReadyMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppOnDomReadyMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callback", 0, nil, data.NewBaseType("callable"))}
}
func (m *optionsAppOnDomReadyMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callback", 0, data.NewBaseType("callable"))}
}
func (m *optionsAppOnDomReadyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_onDomReady", v)
	}
	return nil, nil
}

// ====== onShutdown 方法 ======

type optionsAppOnShutdownMethod struct{}

func (m *optionsAppOnShutdownMethod) GetName() string            { return "onShutdown" }
func (m *optionsAppOnShutdownMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppOnShutdownMethod) GetIsStatic() bool          { return false }
func (m *optionsAppOnShutdownMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppOnShutdownMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callback", 0, nil, data.NewBaseType("callable"))}
}
func (m *optionsAppOnShutdownMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callback", 0, data.NewBaseType("callable"))}
}
func (m *optionsAppOnShutdownMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_onShutdown", v)
	}
	return nil, nil
}

// ====== onBeforeClose 方法 ======

type optionsAppOnBeforeCloseMethod struct{}

func (m *optionsAppOnBeforeCloseMethod) GetName() string            { return "onBeforeClose" }
func (m *optionsAppOnBeforeCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppOnBeforeCloseMethod) GetIsStatic() bool          { return false }
func (m *optionsAppOnBeforeCloseMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppOnBeforeCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callback", 0, nil, data.NewBaseType("callable"))}
}
func (m *optionsAppOnBeforeCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callback", 0, data.NewBaseType("callable"))}
}
func (m *optionsAppOnBeforeCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_onBeforeClose", v)
	}
	return nil, nil
}

// ====== bind 方法 ======

type optionsAppBindMethod struct{}

func (m *optionsAppBindMethod) GetName() string            { return "bind" }
func (m *optionsAppBindMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppBindMethod) GetIsStatic() bool          { return false }
func (m *optionsAppBindMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppBindMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "service", 0, nil, data.NewBaseType("object"))}
}
func (m *optionsAppBindMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "service", 0, data.NewBaseType("object"))}
}
func (m *optionsAppBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		existing, _ := cv.GetProperty("_bind")
		if av, ok := existing.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(v))
		}
	}
	return nil, nil
}

// ====== setMenu 方法 ======

type optionsAppSetMenuMethod struct{}

func (m *optionsAppSetMenuMethod) GetName() string            { return "setMenu" }
func (m *optionsAppSetMenuMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppSetMenuMethod) GetIsStatic() bool          { return false }
func (m *optionsAppSetMenuMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppSetMenuMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "menu", 0, nil, data.NewBaseType("Wails\\Menu\\Menu"))}
}
func (m *optionsAppSetMenuMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "menu", 0, data.NewBaseType("Wails\\Menu\\Menu"))}
}
func (m *optionsAppSetMenuMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_menu", v)
	}
	return nil, nil
}

// ====== bindingsAllowedOrigins 方法 ======

type optionsAppBindingsAllowedOriginsMethod struct{}

func (m *optionsAppBindingsAllowedOriginsMethod) GetName() string { return "bindingsAllowedOrigins" }
func (m *optionsAppBindingsAllowedOriginsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *optionsAppBindingsAllowedOriginsMethod) GetIsStatic() bool         { return false }
func (m *optionsAppBindingsAllowedOriginsMethod) GetReturnType() data.Types { return nil }
func (m *optionsAppBindingsAllowedOriginsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "origins", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *optionsAppBindingsAllowedOriginsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "origins", 0, data.NewBaseType("string"))}
}
func (m *optionsAppBindingsAllowedOriginsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_bindingsAllowedOrigins", v)
	}
	return nil, nil
}

// ====== setErrorFormatter 方法 ======

type optionsAppSetErrorFormatterMethod struct{}

func (m *optionsAppSetErrorFormatterMethod) GetName() string            { return "setErrorFormatter" }
func (m *optionsAppSetErrorFormatterMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *optionsAppSetErrorFormatterMethod) GetIsStatic() bool          { return false }
func (m *optionsAppSetErrorFormatterMethod) GetReturnType() data.Types  { return nil }
func (m *optionsAppSetErrorFormatterMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "formatter", 0, nil, data.NewBaseType("callable"))}
}
func (m *optionsAppSetErrorFormatterMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "formatter", 0, data.NewBaseType("callable"))}
}
func (m *optionsAppSetErrorFormatterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_errorFormatter", v)
	}
	return nil, nil
}
