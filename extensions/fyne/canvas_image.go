package fyne

import (
	"fyne.io/fyne/v2/canvas"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CanvasImageClass 是 Fyne\Canvas\Image 类
type CanvasImageClass struct {
	*CanvasObjectClass
}

func NewCanvasImageClass() data.ClassStmt {
	return &CanvasImageClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Canvas\\Image", nil),
	}
}

func (c *CanvasImageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CanvasImageClass) GetConstruct() data.Method { return &canvasImageConstruct{} }

func (c *CanvasImageClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setFile":
		return &canvasImageSetFileMethod{}, true
	case "setFillMode":
		return &canvasImageSetFillModeMethod{}, true
	case "setTranslucency":
		return &canvasImageSetTranslucencyMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CanvasImageClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&canvasImageSetFileMethod{},
		&canvasImageSetFillModeMethod{},
		&canvasImageSetTranslucencyMethod{},
	)
}

type canvasImageConstruct struct{}

func (m *canvasImageConstruct) GetName() string            { return token.ConstructName }
func (m *canvasImageConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasImageConstruct) GetIsStatic() bool          { return false }
func (m *canvasImageConstruct) GetReturnType() data.Types  { return nil }
func (m *canvasImageConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "file", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}
func (m *canvasImageConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "file", 0, data.NewBaseType("string")),
	}
}

func (m *canvasImageConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	file := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			file = s.AsString()
		}
	}
	var img *canvas.Image
	if file != "" {
		img = canvas.NewImageFromFile(file)
	} else {
		img = &canvas.Image{}
	}
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, img)
			classVal.SetProperty("_image", data.NewAnyValue(img))
		}
	}
	return nil, nil
}

func getCanvasImage(cv *data.ClassValue) *canvas.Image {
	if v, _ := cv.GetProperty("_image"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if i, ok := av.Value.(*canvas.Image); ok {
				return i
			}
		}
	}
	return nil
}

type canvasImageSetFileMethod struct{}

func (m *canvasImageSetFileMethod) GetName() string            { return "setFile" }
func (m *canvasImageSetFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasImageSetFileMethod) GetIsStatic() bool          { return false }
func (m *canvasImageSetFileMethod) GetReturnType() data.Types  { return nil }
func (m *canvasImageSetFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
	}
}
func (m *canvasImageSetFileMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
	}
}
func (m *canvasImageSetFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if img := getCanvasImage(classVal); img != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						newImg := canvas.NewImageFromFile(s.AsString())
						img.File = newImg.File
						img.Resource = newImg.Resource
					}
				}
			}
		}
	}
	return nil, nil
}

type canvasImageSetFillModeMethod struct{}

func (m *canvasImageSetFillModeMethod) GetName() string            { return "setFillMode" }
func (m *canvasImageSetFillModeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasImageSetFillModeMethod) GetIsStatic() bool          { return false }
func (m *canvasImageSetFillModeMethod) GetReturnType() data.Types  { return nil }
func (m *canvasImageSetFillModeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "mode", 0, nil, data.NewBaseType("string")),
	}
}
func (m *canvasImageSetFillModeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "mode", 0, data.NewBaseType("string")),
	}
}
func (m *canvasImageSetFillModeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if img := getCanvasImage(classVal); img != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						switch s.AsString() {
						case "contain":
							img.FillMode = canvas.ImageFillContain
						case "original":
							img.FillMode = canvas.ImageFillOriginal
						default:
							img.FillMode = canvas.ImageFillStretch
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type canvasImageSetTranslucencyMethod struct{}

func (m *canvasImageSetTranslucencyMethod) GetName() string            { return "setTranslucency" }
func (m *canvasImageSetTranslucencyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasImageSetTranslucencyMethod) GetIsStatic() bool          { return false }
func (m *canvasImageSetTranslucencyMethod) GetReturnType() data.Types  { return nil }
func (m *canvasImageSetTranslucencyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "amount", 0, nil, data.NewBaseType("float")),
	}
}
func (m *canvasImageSetTranslucencyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "amount", 0, data.NewBaseType("float")),
	}
}
func (m *canvasImageSetTranslucencyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if img := getCanvasImage(classVal); img != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						a, _ := f.AsFloat()
						img.Translucency = float32(a)
					}
				}
			}
		}
	}
	return nil, nil
}
