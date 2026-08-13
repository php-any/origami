package wails

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ============================================================================
// Wails\Menu\MenuItem — 菜单项
//
// 用法示例:
//
//	$item = new Wails\Menu\MenuItem([
//	    'Label' => 'Open File',
//	    'Accelerator' => Wails\Menu\Keys::cmdOrCtrl("o"),
//	    'Type' => Wails\MenuItemType::TEXT,
//	]);
//	$item->onClick(function($data) {
//	    echo "Clicked!";
//	});
//	$item->setChecked(true);
// ============================================================================

type MenuItemClass struct{}

func NewMenuItemClass() data.ClassStmt { return &MenuItemClass{} }

func (c *MenuItemClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MenuItemClass) GetFrom() data.From                            { return nil }
func (c *MenuItemClass) GetName() string                               { return "Wails\\Menu\\MenuItem" }
func (c *MenuItemClass) GetExtend() *string                            { return nil }
func (c *MenuItemClass) GetImplements() []string                       { return nil }
func (c *MenuItemClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MenuItemClass) GetPropertyList() []data.Property              { return nil }
func (c *MenuItemClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

var menuItemMethods = map[string]data.Method{
	"setLabel":       &menuItemSetLabelMethod{},
	"setAccelerator": &menuItemSetAcceleratorMethod{},
	"setChecked":     &menuItemSetCheckedMethod{},
	"setDisabled":    &menuItemSetDisabledMethod{},
	"setHidden":      &menuItemSetHiddenMethod{},
	"onClick":        &menuItemOnClickMethod{},
	"setSubMenu":     &menuItemSetSubMenuMethod{},
}

func (c *MenuItemClass) GetMethod(name string) (data.Method, bool) {
	m, ok := menuItemMethods[name]
	return m, ok
}

func (c *MenuItemClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(menuItemMethods))
	for _, m := range menuItemMethods {
		methods = append(methods, m)
	}
	return methods
}
func (c *MenuItemClass) GetConstruct() data.Method { return &menuItemConstruct{} }

type menuItemConstruct struct{}

func (m *menuItemConstruct) GetName() string            { return token.ConstructName }
func (m *menuItemConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemConstruct) GetIsStatic() bool          { return false }
func (m *menuItemConstruct) GetReturnType() data.Types  { return nil }

func (m *menuItemConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}
func (m *menuItemConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *menuItemConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultStringProperty(cv, "Label", "")
	setDefaultStringProperty(cv, "Accelerator", "")
	setDefaultStringProperty(cv, "Type", "Text")
	setDefaultBoolProperty(cv, "Disabled", false)
	setDefaultBoolProperty(cv, "Hidden", false)
	setDefaultBoolProperty(cv, "Checked", false)
	cv.SetProperty("_onClick", data.NewNullValue())
	cv.SetProperty("_subMenu", data.NewNullValue())

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"Label", "Accelerator", "Type",
				"Disabled", "Hidden", "Checked",
			})
			if v, ok := arrayGet(av, "SubMenu"); ok {
				cv.SetProperty("_subMenu", v)
			}
			if v, ok := arrayGet(av, "OnClick"); ok {
				cv.SetProperty("_onClick", v)
			}
		}
	}
	return nil, nil
}

// ====== setLabel ======

type menuItemSetLabelMethod struct{}

func (m *menuItemSetLabelMethod) GetName() string            { return "setLabel" }
func (m *menuItemSetLabelMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemSetLabelMethod) GetIsStatic() bool          { return false }
func (m *menuItemSetLabelMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemSetLabelMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}
func (m *menuItemSetLabelMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
	}
}

func (m *menuItemSetLabelMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("Label", data.NewStringValue(toString(v)))
	}
	return nil, nil
}

// ====== setAccelerator ======

type menuItemSetAcceleratorMethod struct{}

func (m *menuItemSetAcceleratorMethod) GetName() string            { return "setAccelerator" }
func (m *menuItemSetAcceleratorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemSetAcceleratorMethod) GetIsStatic() bool          { return false }
func (m *menuItemSetAcceleratorMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemSetAcceleratorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "accelerator", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}
func (m *menuItemSetAcceleratorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "accelerator", 0, data.NewBaseType("string")),
	}
}

func (m *menuItemSetAcceleratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("Accelerator", data.NewStringValue(toString(v)))
	}
	return nil, nil
}

// ====== setChecked ======

type menuItemSetCheckedMethod struct{}

func (m *menuItemSetCheckedMethod) GetName() string            { return "setChecked" }
func (m *menuItemSetCheckedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemSetCheckedMethod) GetIsStatic() bool          { return false }
func (m *menuItemSetCheckedMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemSetCheckedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "checked", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}
func (m *menuItemSetCheckedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "checked", 0, data.NewBaseType("bool")),
	}
}

func (m *menuItemSetCheckedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("Checked", data.NewBoolValue(toBool(v)))
	}
	return nil, nil
}

// ====== setDisabled ======

type menuItemSetDisabledMethod struct{}

func (m *menuItemSetDisabledMethod) GetName() string            { return "setDisabled" }
func (m *menuItemSetDisabledMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemSetDisabledMethod) GetIsStatic() bool          { return false }
func (m *menuItemSetDisabledMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemSetDisabledMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "disabled", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}
func (m *menuItemSetDisabledMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "disabled", 0, data.NewBaseType("bool")),
	}
}

func (m *menuItemSetDisabledMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("Disabled", data.NewBoolValue(toBool(v)))
	}
	return nil, nil
}

// ====== setHidden ======

type menuItemSetHiddenMethod struct{}

func (m *menuItemSetHiddenMethod) GetName() string            { return "setHidden" }
func (m *menuItemSetHiddenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemSetHiddenMethod) GetIsStatic() bool          { return false }
func (m *menuItemSetHiddenMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemSetHiddenMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "hidden", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}
func (m *menuItemSetHiddenMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "hidden", 0, data.NewBaseType("bool")),
	}
}

func (m *menuItemSetHiddenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("Hidden", data.NewBoolValue(toBool(v)))
	}
	return nil, nil
}

// ====== onClick ======

type menuItemOnClickMethod struct{}

func (m *menuItemOnClickMethod) GetName() string            { return "onClick" }
func (m *menuItemOnClickMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemOnClickMethod) GetIsStatic() bool          { return false }
func (m *menuItemOnClickMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemOnClickMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, data.NewBaseType("callable")),
	}
}
func (m *menuItemOnClickMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.NewBaseType("callable")),
	}
}

func (m *menuItemOnClickMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_onClick", v)
	}
	return nil, nil
}

// ====== setSubMenu ======

type menuItemSetSubMenuMethod struct{}

func (m *menuItemSetSubMenuMethod) GetName() string            { return "setSubMenu" }
func (m *menuItemSetSubMenuMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuItemSetSubMenuMethod) GetIsStatic() bool          { return false }
func (m *menuItemSetSubMenuMethod) GetReturnType() data.Types  { return nil }

func (m *menuItemSetSubMenuMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "menu", 0, nil, data.NewBaseType("Wails\\Menu\\Menu")),
	}
}
func (m *menuItemSetSubMenuMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "menu", 0, data.NewBaseType("Wails\\Menu\\Menu")),
	}
}

func (m *menuItemSetSubMenuMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("_subMenu", v)
	}
	return nil, nil
}

// ============================================================================
// Wails\Menu\Menu — 菜单容器
//
// 用法示例:
//
//	$menu = new Wails\Menu\Menu();
//	$menu->addText("&File", "", function($data) {});
//	$menu->addSeparator();
//	$menu->append($item);
// ============================================================================

type MenuClass struct{}

func NewMenuClass() data.ClassStmt { return &MenuClass{} }

func (c *MenuClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MenuClass) GetFrom() data.From                            { return nil }
func (c *MenuClass) GetName() string                               { return "Wails\\Menu\\Menu" }
func (c *MenuClass) GetExtend() *string                            { return nil }
func (c *MenuClass) GetImplements() []string                       { return nil }
func (c *MenuClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MenuClass) GetPropertyList() []data.Property              { return nil }
func (c *MenuClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

var menuMethods = map[string]data.Method{
	"append":        &menuAppendMethod{},
	"prepend":       &menuPrependMethod{},
	"addText":       &menuAddTextMethod{},
	"addSeparator":  &menuAddSeparatorMethod{},
	"addRadio":      &menuAddRadioMethod{},
	"addCheckbox":   &menuAddCheckboxMethod{},
	"addSubMenu":    &menuAddSubMenuMethod{},
	"insertAfter":   &menuInsertAfterMethod{},
	"insertBefore":  &menuInsertBeforeMethod{},
	"remove":        &menuRemoveMethod{},
	"setApplicationMenu": &menuSetApplicationMenuMethod{},
}

func (c *MenuClass) GetMethod(name string) (data.Method, bool) {
	m, ok := menuMethods[name]
	return m, ok
}

func (c *MenuClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(menuMethods))
	for _, m := range menuMethods {
		methods = append(methods, m)
	}
	return methods
}
func (c *MenuClass) GetConstruct() data.Method { return &menuConstruct{} }

type menuConstruct struct{}

func (m *menuConstruct) GetName() string            { return token.ConstructName }
func (m *menuConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuConstruct) GetIsStatic() bool          { return false }
func (m *menuConstruct) GetReturnType() data.Types  { return nil }

func (m *menuConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "items", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}
func (m *menuConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "items", 0, data.NewBaseType("array")),
	}
}

func (m *menuConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	cv.SetProperty("_items", data.NewArrayValue(nil))

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			cv.SetProperty("_items", av)
		}
	}
	return nil, nil
}

// ====== append ======

type menuAppendMethod struct{}

func (m *menuAppendMethod) GetName() string            { return "append" }
func (m *menuAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuAppendMethod) GetIsStatic() bool          { return false }
func (m *menuAppendMethod) GetReturnType() data.Types  { return nil }

func (m *menuAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}
func (m *menuAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}

func (m *menuAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(v))
		}
	}
	return nil, nil
}

// ====== prepend ======

type menuPrependMethod struct{}

func (m *menuPrependMethod) GetName() string            { return "prepend" }
func (m *menuPrependMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuPrependMethod) GetIsStatic() bool          { return false }
func (m *menuPrependMethod) GetReturnType() data.Types  { return nil }

func (m *menuPrependMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}
func (m *menuPrependMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}

func (m *menuPrependMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append([]*data.ZVal{data.NewZVal(v)}, av.List...)
		}
	}
	return nil, nil
}

// ====== addText ======

type menuAddTextMethod struct{}

func (m *menuAddTextMethod) GetName() string            { return "addText" }
func (m *menuAddTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuAddTextMethod) GetIsStatic() bool          { return false }
func (m *menuAddTextMethod) GetReturnType() data.Types  { return nil }

func (m *menuAddTextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "accelerator", 1, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "click", 2, nil, data.NewBaseType("callable")),
	}
}
func (m *menuAddTextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "accelerator", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "click", 2, data.NewBaseType("callable")),
	}
}

func (m *menuAddTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	label, accel := "", ""
	var onClick data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		label = toString(v)
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		accel = toString(v)
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		onClick = v
	}

	// 创建一个新的 MenuItem 并追加
	miClass := NewMenuItemClass()
	miCV, _ := miClass.GetValue(ctx)
	if classVal, ok := miCV.(*data.ClassValue); ok {
		classVal.SetProperty("Label", data.NewStringValue(label))
		classVal.SetProperty("Accelerator", data.NewStringValue(accel))
		classVal.SetProperty("Type", data.NewStringValue("Text"))
		classVal.SetProperty("_onClick", onClick)
		classVal.SetProperty("Disabled", data.NewBoolValue(false))
		classVal.SetProperty("Hidden", data.NewBoolValue(false))

		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(classVal))
		}
	}
	return nil, nil
}

// ====== addSeparator ======

type menuAddSeparatorMethod struct{}

func (m *menuAddSeparatorMethod) GetName() string               { return "addSeparator" }
func (m *menuAddSeparatorMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *menuAddSeparatorMethod) GetIsStatic() bool             { return false }
func (m *menuAddSeparatorMethod) GetReturnType() data.Types     { return nil }
func (m *menuAddSeparatorMethod) GetParams() []data.GetValue    { return nil }
func (m *menuAddSeparatorMethod) GetVariables() []data.Variable { return nil }

func (m *menuAddSeparatorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	miClass := NewMenuItemClass()
	miCV, _ := miClass.GetValue(ctx)
	if classVal, ok := miCV.(*data.ClassValue); ok {
		classVal.SetProperty("Type", data.NewStringValue("Separator"))

		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(classVal))
		}
	}
	return nil, nil
}

// ====== addRadio ======

type menuAddRadioMethod struct{}

func (m *menuAddRadioMethod) GetName() string            { return "addRadio" }
func (m *menuAddRadioMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuAddRadioMethod) GetIsStatic() bool          { return false }
func (m *menuAddRadioMethod) GetReturnType() data.Types  { return nil }

func (m *menuAddRadioMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "selected", 1, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "accelerator", 2, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "click", 3, nil, data.NewBaseType("callable")),
	}
}
func (m *menuAddRadioMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "selected", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "accelerator", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "click", 3, data.NewBaseType("callable")),
	}
}

func (m *menuAddRadioMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	label, accel := "", ""
	selected := false
	var onClick data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		label = toString(v)
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		selected = toBool(v)
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		accel = toString(v)
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		onClick = v
	}

	miClass := NewMenuItemClass()
	miCV, _ := miClass.GetValue(ctx)
	if classVal, ok := miCV.(*data.ClassValue); ok {
		classVal.SetProperty("Label", data.NewStringValue(label))
		classVal.SetProperty("Checked", data.NewBoolValue(selected))
		classVal.SetProperty("Accelerator", data.NewStringValue(accel))
		classVal.SetProperty("Type", data.NewStringValue("Radio"))
		classVal.SetProperty("_onClick", onClick)

		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(classVal))
		}
	}
	return nil, nil
}

// ====== addCheckbox ======

type menuAddCheckboxMethod struct{}

func (m *menuAddCheckboxMethod) GetName() string            { return "addCheckbox" }
func (m *menuAddCheckboxMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuAddCheckboxMethod) GetIsStatic() bool          { return false }
func (m *menuAddCheckboxMethod) GetReturnType() data.Types  { return nil }

func (m *menuAddCheckboxMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "checked", 1, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "accelerator", 2, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "click", 3, nil, data.NewBaseType("callable")),
	}
}
func (m *menuAddCheckboxMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "checked", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "accelerator", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "click", 3, data.NewBaseType("callable")),
	}
}

func (m *menuAddCheckboxMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	label, accel := "", ""
	checked := false
	var onClick data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		label = toString(v)
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		checked = toBool(v)
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		accel = toString(v)
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		onClick = v
	}

	miClass := NewMenuItemClass()
	miCV, _ := miClass.GetValue(ctx)
	if classVal, ok := miCV.(*data.ClassValue); ok {
		classVal.SetProperty("Label", data.NewStringValue(label))
		classVal.SetProperty("Checked", data.NewBoolValue(checked))
		classVal.SetProperty("Accelerator", data.NewStringValue(accel))
		classVal.SetProperty("Type", data.NewStringValue("Checkbox"))
		classVal.SetProperty("_onClick", onClick)

		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(classVal))
		}
	}
	return nil, nil
}

// ====== addSubMenu ======

type menuAddSubMenuMethod struct{}

func (m *menuAddSubMenuMethod) GetName() string            { return "addSubMenu" }
func (m *menuAddSubMenuMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuAddSubMenuMethod) GetIsStatic() bool          { return false }
func (m *menuAddSubMenuMethod) GetReturnType() data.Types  { return nil }

func (m *menuAddSubMenuMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "subMenu", 1, nil, data.NewBaseType("Wails\\Menu\\Menu")),
	}
}
func (m *menuAddSubMenuMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "subMenu", 1, data.NewBaseType("Wails\\Menu\\Menu")),
	}
}

func (m *menuAddSubMenuMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	label := ""
	var subMenuVal data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		label = toString(v)
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		subMenuVal = v
	}

	miClass := NewMenuItemClass()
	miCV, _ := miClass.GetValue(ctx)
	if classVal, ok := miCV.(*data.ClassValue); ok {
		classVal.SetProperty("Label", data.NewStringValue(label))
		classVal.SetProperty("Type", data.NewStringValue("Submenu"))
		classVal.SetProperty("_subMenu", subMenuVal)

		items, _ := cv.GetProperty("_items")
		if av, ok := items.(*data.ArrayValue); ok {
			av.List = append(av.List, data.NewZVal(classVal))
		}
	}
	return nil, nil
}

// ====== insertAfter ======

type menuInsertAfterMethod struct{}

func (m *menuInsertAfterMethod) GetName() string            { return "insertAfter" }
func (m *menuInsertAfterMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuInsertAfterMethod) GetIsStatic() bool          { return false }
func (m *menuInsertAfterMethod) GetReturnType() data.Types  { return nil }

func (m *menuInsertAfterMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "afterIndex", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "item", 1, nil, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}
func (m *menuInsertAfterMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "afterIndex", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "item", 1, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}

func (m *menuInsertAfterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	afterIdx := 0
	var item data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		afterIdx = toInt(v)
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		item = v
	}

	items, _ := cv.GetProperty("_items")
	if av, ok := items.(*data.ArrayValue); ok {
		if afterIdx+1 < len(av.List) {
			av.List = append(av.List[:afterIdx+1], append([]*data.ZVal{data.NewZVal(item)}, av.List[afterIdx+1:]...)...)
		} else {
			av.List = append(av.List, data.NewZVal(item))
		}
	}
	return nil, nil
}

// ====== insertBefore ======

type menuInsertBeforeMethod struct{}

func (m *menuInsertBeforeMethod) GetName() string            { return "insertBefore" }
func (m *menuInsertBeforeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuInsertBeforeMethod) GetIsStatic() bool          { return false }
func (m *menuInsertBeforeMethod) GetReturnType() data.Types  { return nil }

func (m *menuInsertBeforeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "beforeIndex", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "item", 1, nil, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}
func (m *menuInsertBeforeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "beforeIndex", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "item", 1, data.NewBaseType("Wails\\Menu\\MenuItem")),
	}
}

func (m *menuInsertBeforeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	beforeIdx := 0
	var item data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		beforeIdx = toInt(v)
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		item = v
	}

	items, _ := cv.GetProperty("_items")
	if av, ok := items.(*data.ArrayValue); ok {
		if beforeIdx <= len(av.List) {
			av.List = append(av.List[:beforeIdx], append([]*data.ZVal{data.NewZVal(item)}, av.List[beforeIdx:]...)...)
		} else {
			av.List = append(av.List, data.NewZVal(item))
		}
	}
	return nil, nil
}

// ====== remove ======

type menuRemoveMethod struct{}

func (m *menuRemoveMethod) GetName() string            { return "remove" }
func (m *menuRemoveMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuRemoveMethod) GetIsStatic() bool          { return false }
func (m *menuRemoveMethod) GetReturnType() data.Types  { return nil }

func (m *menuRemoveMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *menuRemoveMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}

func (m *menuRemoveMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	idx := 0
	if v, ok := ctx.GetIndexValue(0); ok {
		idx = toInt(v)
	}

	items, _ := cv.GetProperty("_items")
	if av, ok := items.(*data.ArrayValue); ok {
		if idx >= 0 && idx < len(av.List) {
			av.List = append(av.List[:idx], av.List[idx+1:]...)
		}
	}
	return nil, nil
}

// ====== setApplicationMenu ======

type menuSetApplicationMenuMethod struct{}

func (m *menuSetApplicationMenuMethod) GetName() string            { return "setApplicationMenu" }
func (m *menuSetApplicationMenuMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *menuSetApplicationMenuMethod) GetIsStatic() bool          { return false }
func (m *menuSetApplicationMenuMethod) GetReturnType() data.Types  { return nil }

func (m *menuSetApplicationMenuMethod) GetParams() []data.GetValue    { return nil }
func (m *menuSetApplicationMenuMethod) GetVariables() []data.Variable { return nil }

func (m *menuSetApplicationMenuMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	// 在 Application 层设置菜单; 这里存储菜单引用供后续使用
	cv.SetProperty("_isApplicationMenu", data.NewBoolValue(true))
	return nil, nil
}
