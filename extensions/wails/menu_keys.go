package wails

import (
	"github.com/php-any/origami/data"
)

// ============================================================================
// Wails\Menu\Keys — 键盘修饰键和快捷键
//
// 静态方法用于构造 Accelerator：
//
//	$accel = Wails\Menu\Keys::cmdOrCtrl("o");       // CmdOrCtrl+O
//	$accel = Wails\Menu\Keys::combo("a", ["shift"]); // Shift+A
//	$accel = Wails\Menu\Keys::parse("Ctrl+Shift+N");
//
// 修饰键常量：
//
//	Wails\Menu\Keys::CMD_OR_CTRL
//	Wails\Menu\Keys::OPTION_OR_ALT
//	Wails\Menu\Keys::SHIFT
//	Wails\Menu\Keys::CONTROL
// ============================================================================

type MenuKeysClass struct{}

func NewMenuKeysClass() data.ClassStmt { return &MenuKeysClass{} }

func (c *MenuKeysClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MenuKeysClass) GetFrom() data.From                            { return nil }
func (c *MenuKeysClass) GetName() string                               { return "Wails\\Menu\\Keys" }
func (c *MenuKeysClass) GetExtend() *string                            { return nil }
func (c *MenuKeysClass) GetImplements() []string                       { return nil }
func (c *MenuKeysClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MenuKeysClass) GetPropertyList() []data.Property              { return nil }
func (c *MenuKeysClass) GetConstruct() data.Method                     { return nil }

var menuKeyModifiers = map[string]string{
	"CMD_OR_CTRL":   "cmdorctrl",
	"OPTION_OR_ALT": "optionoralt",
	"SHIFT":         "shift",
	"CONTROL":       "ctrl",
}

func (c *MenuKeysClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := menuKeyModifiers[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}

func (c *MenuKeysClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "cmdOrCtrl":
		return &keysCmdOrCtrlMethod{}, true
	case "optionOrAlt":
		return &keysOptionOrAltMethod{}, true
	case "shift":
		return &keysShiftMethod{}, true
	case "control":
		return &keysControlMethod{}, true
	case "combo":
		return &keysComboMethod{}, true
	case "parse":
		return &keysParseMethod{}, true
	default:
		return nil, false
	}
}

func (c *MenuKeysClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "cmdOrCtrl":
		return &keysCmdOrCtrlMethod{}, true
	case "optionOrAlt":
		return &keysOptionOrAltMethod{}, true
	case "shift":
		return &keysShiftMethod{}, true
	case "control":
		return &keysControlMethod{}, true
	case "combo":
		return &keysComboMethod{}, true
	case "parse":
		return &keysParseMethod{}, true
	default:
		return nil, false
	}
}

func (c *MenuKeysClass) GetMethods() []data.Method {
	return []data.Method{
		&keysCmdOrCtrlMethod{},
		&keysOptionOrAltMethod{},
		&keysShiftMethod{},
		&keysControlMethod{},
		&keysComboMethod{},
		&keysParseMethod{},
	}
}

// ====== cmdOrCtrl(key) ======

type keysCmdOrCtrlMethod struct{}

func (m *keysCmdOrCtrlMethod) GetName() string            { return "cmdOrCtrl" }
func (m *keysCmdOrCtrlMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *keysCmdOrCtrlMethod) GetIsStatic() bool          { return true }
func (m *keysCmdOrCtrlMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }

func (m *keysCmdOrCtrlMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("key", 0),
	}
}
func (m *keysCmdOrCtrlMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("key", 0, data.NewBaseType("string")),
	}
}

func (m *keysCmdOrCtrlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if key := firstArgString(ctx); key != "" {
		return data.NewStringValue("cmdorctrl+" + key), nil
	}
	return data.NewStringValue(""), nil
}

// ====== optionOrAlt(key) ======

type keysOptionOrAltMethod struct{}

func (m *keysOptionOrAltMethod) GetName() string            { return "optionOrAlt" }
func (m *keysOptionOrAltMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *keysOptionOrAltMethod) GetIsStatic() bool          { return true }
func (m *keysOptionOrAltMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }

func (m *keysOptionOrAltMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("key", 0),
	}
}
func (m *keysOptionOrAltMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("key", 0, data.NewBaseType("string")),
	}
}

func (m *keysOptionOrAltMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if key := firstArgString(ctx); key != "" {
		return data.NewStringValue("optionoralt+" + key), nil
	}
	return data.NewStringValue(""), nil
}

// ====== shift(key) ======

type keysShiftMethod struct{}

func (m *keysShiftMethod) GetName() string            { return "shift" }
func (m *keysShiftMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *keysShiftMethod) GetIsStatic() bool          { return true }
func (m *keysShiftMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }

func (m *keysShiftMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("key", 0),
	}
}
func (m *keysShiftMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("key", 0, data.NewBaseType("string")),
	}
}

func (m *keysShiftMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if key := firstArgString(ctx); key != "" {
		return data.NewStringValue("shift+" + key), nil
	}
	return data.NewStringValue(""), nil
}

// ====== control(key) ======

type keysControlMethod struct{}

func (m *keysControlMethod) GetName() string            { return "control" }
func (m *keysControlMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *keysControlMethod) GetIsStatic() bool          { return true }
func (m *keysControlMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }

func (m *keysControlMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("key", 0),
	}
}
func (m *keysControlMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("key", 0, data.NewBaseType("string")),
	}
}

func (m *keysControlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if key := firstArgString(ctx); key != "" {
		return data.NewStringValue("ctrl+" + key), nil
	}
	return data.NewStringValue(""), nil
}

// ====== combo(key, modifiers...) ======

type keysComboMethod struct{}

func (m *keysComboMethod) GetName() string            { return "combo" }
func (m *keysComboMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *keysComboMethod) GetIsStatic() bool          { return true }
func (m *keysComboMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }

func (m *keysComboMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("key", 0),
		data.NewParameter("modifiers", 1),
	}
}
func (m *keysComboMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("key", 0, data.NewBaseType("string")),
		data.NewVariable("modifiers", 1, data.NewBaseType("array")),
	}
}

func (m *keysComboMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	key := firstArgString(ctx)
	prefix := modifierPrefixFromValue(argValueAt(ctx, 1))
	if prefix == "" {
		return data.NewStringValue(key), nil
	}
	return data.NewStringValue(prefix + key), nil
}

// ====== parse(acceleratorString) ======

type keysParseMethod struct{}

func (m *keysParseMethod) GetName() string            { return "parse" }
func (m *keysParseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *keysParseMethod) GetIsStatic() bool          { return true }
func (m *keysParseMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }

func (m *keysParseMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("accelerator", 0),
	}
}
func (m *keysParseMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("accelerator", 0, data.NewBaseType("string")),
	}
}

func (m *keysParseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if accel := firstArgString(ctx); accel != "" {
		return data.NewStringValue(accel), nil
	}
	return data.NewStringValue(""), nil
}
