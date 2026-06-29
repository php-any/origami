package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ContainerClass 是 Fyne\Container 类，提供静态工厂方法
type ContainerClass struct{}

func NewContainerClass() data.ClassStmt { return &ContainerClass{} }

func (c *ContainerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ContainerClass) GetFrom() data.From                            { return nil }
func (c *ContainerClass) GetName() string                               { return "Fyne\\Container" }
func (c *ContainerClass) GetExtend() *string                            { return nil }
func (c *ContainerClass) GetImplements() []string                       { return nil }
func (c *ContainerClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ContainerClass) GetPropertyList() []data.Property              { return nil }
func (c *ContainerClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *ContainerClass) GetMethods() []data.Method                     { return nil }
func (c *ContainerClass) GetConstruct() data.Method                     { return nil }

func (c *ContainerClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "newVBox":
		return &containerNewVBoxMethod{}, true
	case "newHBox":
		return &containerNewHBoxMethod{}, true
	case "newGridWithColumns":
		return &containerNewGridWithColumnsMethod{}, true
	case "newGridWithRows":
		return &containerNewGridWithRowsMethod{}, true
	case "newMax":
		return &containerNewMaxMethod{}, true
	case "newStack":
		return &containerNewStackMethod{}, true
	case "newScroll":
		return &containerNewScrollMethod{}, true
	case "newCenter":
		return &containerNewCenterMethod{}, true
	case "newPadded":
		return &containerNewPaddedMethod{}, true
	case "newBorder":
		return &containerNewBorderMethod{}, true
	case "newAdaptiveGrid":
		return &containerNewAdaptiveGridMethod{}, true
	default:
		return nil, false
	}
}

// extractFyneObjects 从参数中提取所有 fyneLib.CanvasObject
func extractFyneObjects(args []data.Value) []fyneLib.CanvasObject {
	var objs []fyneLib.CanvasObject
	for _, arg := range args {
		if cv, ok := arg.(*data.ClassValue); ok {
			if obj := getFyneObject(cv); obj != nil {
				objs = append(objs, obj)
			}
		}
	}
	return objs
}

// wrapContainer 将 fyneLib.CanvasObject 包装为 Fyne\Container 的 ClassValue
func wrapContainer(obj fyneLib.CanvasObject, ctx data.Context) *data.ClassValue {
	cl := NewContainerClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		setFyneObject(classVal, obj)
		return classVal
	}
	return nil
}

// ====== VBox ======

type containerNewVBoxMethod struct{}

func (m *containerNewVBoxMethod) GetName() string            { return "newVBox" }
func (m *containerNewVBoxMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewVBoxMethod) GetIsStatic() bool          { return true }
func (m *containerNewVBoxMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewVBoxMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "objects", 0, nil, nil),
	}
}
func (m *containerNewVBoxMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "objects", 0, nil),
	}
}
func (m *containerNewVBoxMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var args []data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewVBox(objs...)
	return wrapContainer(c, ctx), nil
}

// ====== HBox ======

type containerNewHBoxMethod struct{}

func (m *containerNewHBoxMethod) GetName() string            { return "newHBox" }
func (m *containerNewHBoxMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewHBoxMethod) GetIsStatic() bool          { return true }
func (m *containerNewHBoxMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewHBoxMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "objects", 0, nil, nil),
	}
}
func (m *containerNewHBoxMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "objects", 0, nil),
	}
}
func (m *containerNewHBoxMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var args []data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewHBox(objs...)
	return wrapContainer(c, ctx), nil
}

// ====== GridWithColumns ======

type containerNewGridWithColumnsMethod struct{}

func (m *containerNewGridWithColumnsMethod) GetName() string            { return "newGridWithColumns" }
func (m *containerNewGridWithColumnsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewGridWithColumnsMethod) GetIsStatic() bool          { return true }
func (m *containerNewGridWithColumnsMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewGridWithColumnsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "cols", 0, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "objects", 1, nil, nil),
	}
}
func (m *containerNewGridWithColumnsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "cols", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "objects", 1, nil),
	}
}
func (m *containerNewGridWithColumnsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cols := 2
	if v, ok := ctx.GetIndexValue(0); ok {
		if i, ok := v.(data.AsInt); ok {
			cols, _ = i.AsInt()
		}
	}
	var args []data.Value
	if v, ok := ctx.GetIndexValue(1); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewGridWithColumns(cols, objs...)
	return wrapContainer(c, ctx), nil
}

// ====== GridWithRows ======

type containerNewGridWithRowsMethod struct{}

func (m *containerNewGridWithRowsMethod) GetName() string            { return "newGridWithRows" }
func (m *containerNewGridWithRowsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewGridWithRowsMethod) GetIsStatic() bool          { return true }
func (m *containerNewGridWithRowsMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewGridWithRowsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "rows", 0, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "objects", 1, nil, nil),
	}
}
func (m *containerNewGridWithRowsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "rows", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "objects", 1, nil),
	}
}
func (m *containerNewGridWithRowsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	rows := 2
	if v, ok := ctx.GetIndexValue(0); ok {
		if i, ok := v.(data.AsInt); ok {
			rows, _ = i.AsInt()
		}
	}
	var args []data.Value
	if v, ok := ctx.GetIndexValue(1); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewGridWithRows(rows, objs...)
	return wrapContainer(c, ctx), nil
}

// ====== Max ======

type containerNewMaxMethod struct{}

func (m *containerNewMaxMethod) GetName() string            { return "newMax" }
func (m *containerNewMaxMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewMaxMethod) GetIsStatic() bool          { return true }
func (m *containerNewMaxMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewMaxMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "objects", 0, nil, nil),
	}
}
func (m *containerNewMaxMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "objects", 0, nil),
	}
}
func (m *containerNewMaxMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var args []data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewMax(objs...)
	return wrapContainer(c, ctx), nil
}

// ====== Stack ======

type containerNewStackMethod struct{}

func (m *containerNewStackMethod) GetName() string            { return "newStack" }
func (m *containerNewStackMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewStackMethod) GetIsStatic() bool          { return true }
func (m *containerNewStackMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewStackMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "objects", 0, nil, nil),
	}
}
func (m *containerNewStackMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "objects", 0, nil),
	}
}
func (m *containerNewStackMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var args []data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewStack(objs...)
	return wrapContainer(c, ctx), nil
}

// ====== Scroll ======

type containerNewScrollMethod struct{}

func (m *containerNewScrollMethod) GetName() string            { return "newScroll" }
func (m *containerNewScrollMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewScrollMethod) GetIsStatic() bool          { return true }
func (m *containerNewScrollMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewScrollMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "content", 0, nil, nil),
	}
}
func (m *containerNewScrollMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "content", 0, nil),
	}
}
func (m *containerNewScrollMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if v, ok := ctx.GetIndexValue(0); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			if obj := getFyneObject(cv); obj != nil {
				c := container.NewScroll(obj)
				return wrapContainer(c, ctx), nil
			}
		}
	}
	return nil, nil
}

// ====== Center ======

type containerNewCenterMethod struct{}

func (m *containerNewCenterMethod) GetName() string            { return "newCenter" }
func (m *containerNewCenterMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewCenterMethod) GetIsStatic() bool          { return true }
func (m *containerNewCenterMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewCenterMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "objects", 0, nil, nil),
	}
}
func (m *containerNewCenterMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "objects", 0, nil),
	}
}
func (m *containerNewCenterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var args []data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewCenter(objs...)
	return wrapContainer(c, ctx), nil
}

// ====== Padded ======

type containerNewPaddedMethod struct{}

func (m *containerNewPaddedMethod) GetName() string            { return "newPadded" }
func (m *containerNewPaddedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewPaddedMethod) GetIsStatic() bool          { return true }
func (m *containerNewPaddedMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewPaddedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "objects", 0, nil, nil),
	}
}
func (m *containerNewPaddedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "objects", 0, nil),
	}
}
func (m *containerNewPaddedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var args []data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewPadded(objs...)
	return wrapContainer(c, ctx), nil
}

// ====== Border ======

type containerNewBorderMethod struct{}

func (m *containerNewBorderMethod) GetName() string            { return "newBorder" }
func (m *containerNewBorderMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewBorderMethod) GetIsStatic() bool          { return true }
func (m *containerNewBorderMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewBorderMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "top", 0, data.NewNullValue(), nil),
		node.NewParameter(nil, "bottom", 1, data.NewNullValue(), nil),
		node.NewParameter(nil, "left", 2, data.NewNullValue(), nil),
		node.NewParameter(nil, "right", 3, data.NewNullValue(), nil),
		node.NewParameter(nil, "center", 4, data.NewNullValue(), nil),
	}
}
func (m *containerNewBorderMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "top", 0, nil),
		node.NewVariable(nil, "bottom", 1, nil),
		node.NewVariable(nil, "left", 2, nil),
		node.NewVariable(nil, "right", 3, nil),
		node.NewVariable(nil, "center", 4, nil),
	}
}
func (m *containerNewBorderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var top, bottom, left, right, center fyneLib.CanvasObject
	if v, ok := ctx.GetIndexValue(0); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			top = getFyneObject(cv)
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			bottom = getFyneObject(cv)
		}
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			left = getFyneObject(cv)
		}
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			right = getFyneObject(cv)
		}
	}
	if v, ok := ctx.GetIndexValue(4); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			center = getFyneObject(cv)
		}
	}
	c := container.NewBorder(top, bottom, left, right, center)
	return wrapContainer(c, ctx), nil
}

// ====== AdaptiveGrid ======

type containerNewAdaptiveGridMethod struct{}

func (m *containerNewAdaptiveGridMethod) GetName() string            { return "newAdaptiveGrid" }
func (m *containerNewAdaptiveGridMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *containerNewAdaptiveGridMethod) GetIsStatic() bool          { return true }
func (m *containerNewAdaptiveGridMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Container")
}
func (m *containerNewAdaptiveGridMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "rowCols", 0, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "objects", 1, nil, nil),
	}
}
func (m *containerNewAdaptiveGridMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "rowCols", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "objects", 1, nil),
	}
}
func (m *containerNewAdaptiveGridMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	rowCols := 2
	if v, ok := ctx.GetIndexValue(0); ok {
		if i, ok := v.(data.AsInt); ok {
			rowCols, _ = i.AsInt()
		}
	}
	var args []data.Value
	if v, ok := ctx.GetIndexValue(1); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			args = arr.ToValueList()
		} else {
			args = []data.Value{v}
		}
	}
	objs := extractFyneObjects(args)
	c := container.NewAdaptiveGrid(rowCols, objs...)
	return wrapContainer(c, ctx), nil
}
