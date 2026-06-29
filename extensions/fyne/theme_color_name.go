package fyne

import "github.com/php-any/origami/data"

// ThemeColorNameClass 是 Fyne\Theme\ColorName 枚举类
type ThemeColorNameClass struct{}

func NewThemeColorNameClass() data.ClassStmt { return &ThemeColorNameClass{} }

func (c *ThemeColorNameClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ThemeColorNameClass) GetFrom() data.From                              { return nil }
func (c *ThemeColorNameClass) GetName() string                                 { return "Fyne\\Theme\\ColorName" }
func (c *ThemeColorNameClass) GetExtend() *string                              { return nil }
func (c *ThemeColorNameClass) GetImplements() []string                         { return nil }
func (c *ThemeColorNameClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ThemeColorNameClass) GetPropertyList() []data.Property                { return nil }
func (c *ThemeColorNameClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ThemeColorNameClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ThemeColorNameClass) GetMethods() []data.Method                       { return nil }
func (c *ThemeColorNameClass) GetConstruct() data.Method                       { return nil }

var colorNames = map[string]string{
	"BACKGROUND":         "background",
	"FOREGROUND":         "foreground",
	"PRIMARY":            "primary",
	"BUTTON":             "button",
	"DISABLED":           "disabled",
	"ERROR":              "error",
	"FOCUS":              "focus",
	"HOVER":              "hover",
	"INPUT_BACKGROUND":   "inputBackground",
	"INPUT_BORDER":       "inputBorder",
	"MENU_BACKGROUND":    "menuBackground",
	"OVERLAY_BACKGROUND": "overlayBackground",
	"PLACEHOLDER":        "placeholder",
	"PRESSED":            "pressed",
	"SCROLL_BAR":         "scrollBar",
	"SHADOW":             "shadow",
	"SEPARATOR":          "separator",
}

func (c *ThemeColorNameClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := colorNames[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}
