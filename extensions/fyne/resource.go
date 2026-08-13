package fyne

import (
	"os"

	fyneLib "fyne.io/fyne/v2"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ResourceClass 是 Fyne\Resource 类
type ResourceClass struct{}

func NewResourceClass() data.ClassStmt { return &ResourceClass{} }

func (c *ResourceClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ResourceClass) GetFrom() data.From                            { return nil }
func (c *ResourceClass) GetName() string                               { return "Fyne\\Resource" }
func (c *ResourceClass) GetExtend() *string                            { return nil }
func (c *ResourceClass) GetImplements() []string                       { return nil }
func (c *ResourceClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ResourceClass) GetPropertyList() []data.Property              { return nil }
func (c *ResourceClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *ResourceClass) GetMethods() []data.Method                     { return nil }
func (c *ResourceClass) GetConstruct() data.Method                     { return nil }

func (c *ResourceClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "newStatic":
		return &resourceNewStaticMethod{}, true
	case "fromFile":
		return &resourceFromFileMethod{}, true
	default:
		return nil, false
	}
}

type resourceNewStaticMethod struct{}

func (m *resourceNewStaticMethod) GetName() string            { return "newStatic" }
func (m *resourceNewStaticMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *resourceNewStaticMethod) GetIsStatic() bool          { return true }
func (m *resourceNewStaticMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Resource")
}
func (m *resourceNewStaticMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "content", 1, nil, data.NewBaseType("string")),
	}
}
func (m *resourceNewStaticMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "content", 1, data.NewBaseType("string")),
	}
}

func (m *resourceNewStaticMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name := ""
	content := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			name = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if s, ok := v.(data.AsString); ok {
			content = s.AsString()
		}
	}
	res := fyneLib.NewStaticResource(name, []byte(content))
	cl := NewResourceClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("_resource", data.NewAnyValue(res))
		return classVal, nil
	}
	return nil, nil
}

type resourceFromFileMethod struct{}

func (m *resourceFromFileMethod) GetName() string            { return "fromFile" }
func (m *resourceFromFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *resourceFromFileMethod) GetIsStatic() bool          { return true }
func (m *resourceFromFileMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Resource")
}
func (m *resourceFromFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
	}
}
func (m *resourceFromFileMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
	}
}

func (m *resourceFromFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	path := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			path = s.AsString()
		}
	}
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	res := fyneLib.NewStaticResource(path, fileData)
	cl := NewResourceClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("_resource", data.NewAnyValue(res))
		return classVal, nil
	}
	return nil, nil
}
