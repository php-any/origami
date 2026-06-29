package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CanvasObjectClass 是 Fyne\CanvasObject 基类
// 所有 widget 和 canvas 对象的公共基类，提供 show/hide/move/resize/refresh 方法
type CanvasObjectClass struct {
	name    string
	extend  *string
	methods map[string]data.Method
}

func newCanvasObjectClass(name string, extend *string) *CanvasObjectClass {
	c := &CanvasObjectClass{
		name:    name,
		extend:  extend,
		methods: make(map[string]data.Method),
	}
	c.methods["show"] = &canvasObjectShowMethod{}
	c.methods["hide"] = &canvasObjectHideMethod{}
	c.methods["move"] = &canvasObjectMoveMethod{}
	c.methods["resize"] = &canvasObjectResizeMethod{}
	c.methods["refresh"] = &canvasObjectRefreshMethod{}
	c.methods["visible"] = &canvasObjectVisibleMethod{}
	c.methods["getPosition"] = &canvasObjectGetPositionMethod{}
	c.methods["getSize"] = &canvasObjectGetSizeMethod{}
	return c
}

func (c *CanvasObjectClass) GetFrom() data.From                              { return nil }
func (c *CanvasObjectClass) GetName() string                                 { return c.name }
func (c *CanvasObjectClass) GetExtend() *string                              { return c.extend }
func (c *CanvasObjectClass) GetImplements() []string                         { return nil }
func (c *CanvasObjectClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *CanvasObjectClass) GetPropertyList() []data.Property                { return nil }
func (c *CanvasObjectClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *CanvasObjectClass) GetConstruct() data.Method                       { return nil }

func (c *CanvasObjectClass) GetMethod(name string) (data.Method, bool) {
	if m, ok := c.methods[name]; ok {
		return m, true
	}
	return nil, false
}

func (c *CanvasObjectClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		methods = append(methods, m)
	}
	return methods
}

// getFyneObject 从 ClassValue 中获取底层 fyneLib.CanvasObject
func getFyneObject(cv *data.ClassValue) fyneLib.CanvasObject {
	if obj, _ := cv.GetProperty("fyneObj"); obj != nil {
		if av, ok := obj.(*data.AnyValue); ok {
			if fyneObj, ok := av.Value.(fyneLib.CanvasObject); ok {
				return fyneObj
			}
		}
	}
	return nil
}

// setFyneObject 将 fyneLib.CanvasObject 存入 ClassValue
func setFyneObject(cv *data.ClassValue, obj fyneLib.CanvasObject) {
	cv.SetProperty("fyneObj", data.NewAnyValue(obj))
}

// ====== 公共方法实现 ======

type canvasObjectShowMethod struct{}

func (m *canvasObjectShowMethod) GetName() string               { return "show" }
func (m *canvasObjectShowMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *canvasObjectShowMethod) GetIsStatic() bool             { return false }
func (m *canvasObjectShowMethod) GetReturnType() data.Types     { return nil }
func (m *canvasObjectShowMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasObjectShowMethod) GetVariables() []data.Variable { return nil }
func (m *canvasObjectShowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				obj.Show()
			}
		}
	}
	return nil, nil
}

type canvasObjectHideMethod struct{}

func (m *canvasObjectHideMethod) GetName() string               { return "hide" }
func (m *canvasObjectHideMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *canvasObjectHideMethod) GetIsStatic() bool             { return false }
func (m *canvasObjectHideMethod) GetReturnType() data.Types     { return nil }
func (m *canvasObjectHideMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasObjectHideMethod) GetVariables() []data.Variable { return nil }
func (m *canvasObjectHideMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				obj.Hide()
			}
		}
	}
	return nil, nil
}

type canvasObjectMoveMethod struct{}

func (m *canvasObjectMoveMethod) GetName() string            { return "move" }
func (m *canvasObjectMoveMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasObjectMoveMethod) GetIsStatic() bool          { return false }
func (m *canvasObjectMoveMethod) GetReturnType() data.Types  { return nil }
func (m *canvasObjectMoveMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "position", 0, nil, data.NewBaseType("Fyne\\Position")),
	}
}
func (m *canvasObjectMoveMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "position", 0, data.NewBaseType("Fyne\\Position")),
	}
}
func (m *canvasObjectMoveMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if posCV, ok := v.(*data.ClassValue); ok {
						if p, _ := posCV.GetProperty("_pos"); p != nil {
							if av, ok := p.(*data.AnyValue); ok {
								if pos, ok := av.Value.(fyneLib.Position); ok {
									obj.Move(pos)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type canvasObjectResizeMethod struct{}

func (m *canvasObjectResizeMethod) GetName() string            { return "resize" }
func (m *canvasObjectResizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasObjectResizeMethod) GetIsStatic() bool          { return false }
func (m *canvasObjectResizeMethod) GetReturnType() data.Types  { return nil }
func (m *canvasObjectResizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "size", 0, nil, data.NewBaseType("Fyne\\Size")),
	}
}
func (m *canvasObjectResizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "size", 0, data.NewBaseType("Fyne\\Size")),
	}
}
func (m *canvasObjectResizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if sizeCV, ok := v.(*data.ClassValue); ok {
						if s, _ := sizeCV.GetProperty("_size"); s != nil {
							if av, ok := s.(*data.AnyValue); ok {
								if size, ok := av.Value.(fyneLib.Size); ok {
									obj.Resize(size)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type canvasObjectRefreshMethod struct{}

func (m *canvasObjectRefreshMethod) GetName() string               { return "refresh" }
func (m *canvasObjectRefreshMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *canvasObjectRefreshMethod) GetIsStatic() bool             { return false }
func (m *canvasObjectRefreshMethod) GetReturnType() data.Types     { return nil }
func (m *canvasObjectRefreshMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasObjectRefreshMethod) GetVariables() []data.Variable { return nil }
func (m *canvasObjectRefreshMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				obj.Refresh()
			}
		}
	}
	return nil, nil
}

type canvasObjectVisibleMethod struct{}

func (m *canvasObjectVisibleMethod) GetName() string               { return "visible" }
func (m *canvasObjectVisibleMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *canvasObjectVisibleMethod) GetIsStatic() bool             { return false }
func (m *canvasObjectVisibleMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *canvasObjectVisibleMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasObjectVisibleMethod) GetVariables() []data.Variable { return nil }
func (m *canvasObjectVisibleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				return data.NewBoolValue(obj.Visible()), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}

type canvasObjectGetPositionMethod struct{}

func (m *canvasObjectGetPositionMethod) GetName() string            { return "getPosition" }
func (m *canvasObjectGetPositionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasObjectGetPositionMethod) GetIsStatic() bool          { return false }
func (m *canvasObjectGetPositionMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Position")
}
func (m *canvasObjectGetPositionMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasObjectGetPositionMethod) GetVariables() []data.Variable { return nil }
func (m *canvasObjectGetPositionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				return NewPositionValue(obj.Position(), ctx), nil
			}
		}
	}
	return nil, nil
}

type canvasObjectGetSizeMethod struct{}

func (m *canvasObjectGetSizeMethod) GetName() string               { return "getSize" }
func (m *canvasObjectGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *canvasObjectGetSizeMethod) GetIsStatic() bool             { return false }
func (m *canvasObjectGetSizeMethod) GetReturnType() data.Types     { return data.NewBaseType("Fyne\\Size") }
func (m *canvasObjectGetSizeMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasObjectGetSizeMethod) GetVariables() []data.Variable { return nil }
func (m *canvasObjectGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if obj := getFyneObject(classVal); obj != nil {
				return NewSizeValue(obj.Size(), ctx), nil
			}
		}
	}
	return nil, nil
}
