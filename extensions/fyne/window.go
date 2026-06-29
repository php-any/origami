package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// WindowClass 是 Fyne\Window 类
type WindowClass struct {
	inner fyneLib.Window
}

func NewWindowClass() data.ClassStmt {
	return &WindowClass{}
}

func (c *WindowClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *WindowClass) GetFrom() data.From                              { return nil }
func (c *WindowClass) GetName() string                                 { return "Fyne\\Window" }
func (c *WindowClass) GetExtend() *string                              { return nil }
func (c *WindowClass) GetImplements() []string                         { return nil }
func (c *WindowClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *WindowClass) GetPropertyList() []data.Property                { return nil }
func (c *WindowClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *WindowClass) GetConstruct() data.Method                       { return nil }

func (c *WindowClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setContent":
		return &windowSetContentMethod{}, true
	case "show":
		return &windowShowMethod{}, true
	case "hide":
		return &windowHideMethod{}, true
	case "close":
		return &windowCloseMethod{}, true
	case "showAndRun":
		return &windowShowAndRunMethod{}, true
	case "resize":
		return &windowResizeMethod{}, true
	case "centerOnScreen":
		return &windowCenterOnScreenMethod{}, true
	case "setTitle":
		return &windowSetTitleMethod{}, true
	case "getTitle":
		return &windowGetTitleMethod{}, true
	case "setFixedSize":
		return &windowSetFixedSizeMethod{}, true
	case "setPadded":
		return &windowSetPaddedMethod{}, true
	case "requestFocus":
		return &windowRequestFocusMethod{}, true
	case "fullScreen":
		return &windowFullScreenMethod{}, true
	case "setFullScreen":
		return &windowSetFullScreenMethod{}, true
	default:
		return nil, false
	}
}

func (c *WindowClass) GetMethods() []data.Method {
	return []data.Method{
		&windowSetContentMethod{},
		&windowShowMethod{},
		&windowHideMethod{},
		&windowCloseMethod{},
		&windowShowAndRunMethod{},
		&windowResizeMethod{},
		&windowCenterOnScreenMethod{},
		&windowSetTitleMethod{},
		&windowGetTitleMethod{},
		&windowSetFixedSizeMethod{},
		&windowSetPaddedMethod{},
		&windowRequestFocusMethod{},
		&windowFullScreenMethod{},
		&windowSetFullScreenMethod{},
	}
}

// getFyneWindow 从 ClassValue 中获取底层 fyne.Window
func getFyneWindow(cv *data.ClassValue) fyneLib.Window {
	if v, _ := cv.GetProperty("_window"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if w, ok := av.Value.(fyneLib.Window); ok {
				return w
			}
		}
	}
	return nil
}

// NewWindowValue 创建一个 Fyne\Window 的 ClassValue
func NewWindowValue(w fyneLib.Window, ctx data.Context) *data.ClassValue {
	cl := NewWindowClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("_window", data.NewAnyValue(w))
		return classVal
	}
	return nil
}

// ====== setContent ======

type windowSetContentMethod struct{}

func (m *windowSetContentMethod) GetName() string            { return "setContent" }
func (m *windowSetContentMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowSetContentMethod) GetIsStatic() bool          { return false }
func (m *windowSetContentMethod) GetReturnType() data.Types  { return nil }

func (m *windowSetContentMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "content", 0, nil, data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *windowSetContentMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "content", 0, data.NewBaseType("Fyne\\CanvasObject")),
	}
}

func (m *windowSetContentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			w := getFyneWindow(classVal)
			if w == nil {
				return nil, nil
			}
			if v, ok := ctx.GetIndexValue(0); ok {
				if contentCV, ok := v.(*data.ClassValue); ok {
					if obj := getFyneObject(contentCV); obj != nil {
						w.SetContent(obj)
					}
				}
			}
		}
	}
	return nil, nil
}

// ====== show ======

type windowShowMethod struct{}

func (m *windowShowMethod) GetName() string               { return "show" }
func (m *windowShowMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowShowMethod) GetIsStatic() bool             { return false }
func (m *windowShowMethod) GetReturnType() data.Types     { return nil }
func (m *windowShowMethod) GetParams() []data.GetValue    { return nil }
func (m *windowShowMethod) GetVariables() []data.Variable { return nil }

func (m *windowShowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				w.Show()
			}
		}
	}
	return nil, nil
}

// ====== hide ======

type windowHideMethod struct{}

func (m *windowHideMethod) GetName() string               { return "hide" }
func (m *windowHideMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowHideMethod) GetIsStatic() bool             { return false }
func (m *windowHideMethod) GetReturnType() data.Types     { return nil }
func (m *windowHideMethod) GetParams() []data.GetValue    { return nil }
func (m *windowHideMethod) GetVariables() []data.Variable { return nil }

func (m *windowHideMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				w.Hide()
			}
		}
	}
	return nil, nil
}

// ====== close ======

type windowCloseMethod struct{}

func (m *windowCloseMethod) GetName() string               { return "close" }
func (m *windowCloseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowCloseMethod) GetIsStatic() bool             { return false }
func (m *windowCloseMethod) GetReturnType() data.Types     { return nil }
func (m *windowCloseMethod) GetParams() []data.GetValue    { return nil }
func (m *windowCloseMethod) GetVariables() []data.Variable { return nil }

func (m *windowCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				w.Close()
			}
		}
	}
	return nil, nil
}

// ====== showAndRun ======

type windowShowAndRunMethod struct{}

func (m *windowShowAndRunMethod) GetName() string               { return "showAndRun" }
func (m *windowShowAndRunMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowShowAndRunMethod) GetIsStatic() bool             { return false }
func (m *windowShowAndRunMethod) GetReturnType() data.Types     { return nil }
func (m *windowShowAndRunMethod) GetParams() []data.GetValue    { return nil }
func (m *windowShowAndRunMethod) GetVariables() []data.Variable { return nil }

func (m *windowShowAndRunMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				w.ShowAndRun()
			}
		}
	}
	return nil, nil
}

// ====== resize ======

type windowResizeMethod struct{}

func (m *windowResizeMethod) GetName() string            { return "resize" }
func (m *windowResizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowResizeMethod) GetIsStatic() bool          { return false }
func (m *windowResizeMethod) GetReturnType() data.Types  { return nil }

func (m *windowResizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "size", 0, nil, data.NewBaseType("Fyne\\Size")),
	}
}
func (m *windowResizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "size", 0, data.NewBaseType("Fyne\\Size")),
	}
}

func (m *windowResizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if sizeCV, ok := v.(*data.ClassValue); ok {
						if s, _ := sizeCV.GetProperty("_size"); s != nil {
							if av, ok := s.(*data.AnyValue); ok {
								if size, ok := av.Value.(fyneLib.Size); ok {
									w.Resize(size)
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

// ====== centerOnScreen ======

type windowCenterOnScreenMethod struct{}

func (m *windowCenterOnScreenMethod) GetName() string               { return "centerOnScreen" }
func (m *windowCenterOnScreenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowCenterOnScreenMethod) GetIsStatic() bool             { return false }
func (m *windowCenterOnScreenMethod) GetReturnType() data.Types     { return nil }
func (m *windowCenterOnScreenMethod) GetParams() []data.GetValue    { return nil }
func (m *windowCenterOnScreenMethod) GetVariables() []data.Variable { return nil }

func (m *windowCenterOnScreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				w.CenterOnScreen()
			}
		}
	}
	return nil, nil
}

// ====== setTitle ======

type windowSetTitleMethod struct{}

func (m *windowSetTitleMethod) GetName() string            { return "setTitle" }
func (m *windowSetTitleMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowSetTitleMethod) GetIsStatic() bool          { return false }
func (m *windowSetTitleMethod) GetReturnType() data.Types  { return nil }

func (m *windowSetTitleMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
	}
}
func (m *windowSetTitleMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
	}
}

func (m *windowSetTitleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						w.SetTitle(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

// ====== getTitle ======

type windowGetTitleMethod struct{}

func (m *windowGetTitleMethod) GetName() string               { return "getTitle" }
func (m *windowGetTitleMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowGetTitleMethod) GetIsStatic() bool             { return false }
func (m *windowGetTitleMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *windowGetTitleMethod) GetParams() []data.GetValue    { return nil }
func (m *windowGetTitleMethod) GetVariables() []data.Variable { return nil }

func (m *windowGetTitleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				return data.NewStringValue(w.Title()), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}

// ====== setFixedSize ======

type windowSetFixedSizeMethod struct{}

func (m *windowSetFixedSizeMethod) GetName() string            { return "setFixedSize" }
func (m *windowSetFixedSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowSetFixedSizeMethod) GetIsStatic() bool          { return false }
func (m *windowSetFixedSizeMethod) GetReturnType() data.Types  { return nil }

func (m *windowSetFixedSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "fixed", 0, nil, data.NewBaseType("bool")),
	}
}
func (m *windowSetFixedSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "fixed", 0, data.NewBaseType("bool")),
	}
}

func (m *windowSetFixedSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if b, ok := v.(data.AsBool); ok {
						fixed, _ := b.AsBool()
						w.SetFixedSize(fixed)
					}
				}
			}
		}
	}
	return nil, nil
}

// ====== setPadded ======

type windowSetPaddedMethod struct{}

func (m *windowSetPaddedMethod) GetName() string            { return "setPadded" }
func (m *windowSetPaddedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowSetPaddedMethod) GetIsStatic() bool          { return false }
func (m *windowSetPaddedMethod) GetReturnType() data.Types  { return nil }

func (m *windowSetPaddedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "padded", 0, nil, data.NewBaseType("bool")),
	}
}
func (m *windowSetPaddedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "padded", 0, data.NewBaseType("bool")),
	}
}

func (m *windowSetPaddedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if b, ok := v.(data.AsBool); ok {
						padded, _ := b.AsBool()
						w.SetPadded(padded)
					}
				}
			}
		}
	}
	return nil, nil
}

// ====== requestFocus ======

type windowRequestFocusMethod struct{}

func (m *windowRequestFocusMethod) GetName() string               { return "requestFocus" }
func (m *windowRequestFocusMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowRequestFocusMethod) GetIsStatic() bool             { return false }
func (m *windowRequestFocusMethod) GetReturnType() data.Types     { return nil }
func (m *windowRequestFocusMethod) GetParams() []data.GetValue    { return nil }
func (m *windowRequestFocusMethod) GetVariables() []data.Variable { return nil }

func (m *windowRequestFocusMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				w.RequestFocus()
			}
		}
	}
	return nil, nil
}

// ====== fullScreen ======

type windowFullScreenMethod struct{}

func (m *windowFullScreenMethod) GetName() string               { return "fullScreen" }
func (m *windowFullScreenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *windowFullScreenMethod) GetIsStatic() bool             { return false }
func (m *windowFullScreenMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *windowFullScreenMethod) GetParams() []data.GetValue    { return nil }
func (m *windowFullScreenMethod) GetVariables() []data.Variable { return nil }

func (m *windowFullScreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				return data.NewBoolValue(w.FullScreen()), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}

// ====== setFullScreen ======

type windowSetFullScreenMethod struct{}

func (m *windowSetFullScreenMethod) GetName() string            { return "setFullScreen" }
func (m *windowSetFullScreenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowSetFullScreenMethod) GetIsStatic() bool          { return false }
func (m *windowSetFullScreenMethod) GetReturnType() data.Types  { return nil }

func (m *windowSetFullScreenMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "fullScreen", 0, nil, data.NewBaseType("bool")),
	}
}
func (m *windowSetFullScreenMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "fullScreen", 0, data.NewBaseType("bool")),
	}
}

func (m *windowSetFullScreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if w := getFyneWindow(classVal); w != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if b, ok := v.(data.AsBool); ok {
						fullScreen, _ := b.AsBool()
						w.SetFullScreen(fullScreen)
					}
				}
			}
		}
	}
	return nil, nil
}
