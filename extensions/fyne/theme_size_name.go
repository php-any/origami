package fyne

import "github.com/php-any/origami/data"

// ThemeSizeNameClass 是 Fyne\Theme\SizeName 枚举类
type ThemeSizeNameClass struct{}

func NewThemeSizeNameClass() data.ClassStmt { return &ThemeSizeNameClass{} }

func (c *ThemeSizeNameClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ThemeSizeNameClass) GetFrom() data.From                              { return nil }
func (c *ThemeSizeNameClass) GetName() string                                 { return "Fyne\\Theme\\SizeName" }
func (c *ThemeSizeNameClass) GetExtend() *string                              { return nil }
func (c *ThemeSizeNameClass) GetImplements() []string                         { return nil }
func (c *ThemeSizeNameClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ThemeSizeNameClass) GetPropertyList() []data.Property                { return nil }
func (c *ThemeSizeNameClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ThemeSizeNameClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ThemeSizeNameClass) GetMethods() []data.Method                       { return nil }
func (c *ThemeSizeNameClass) GetConstruct() data.Method                       { return nil }

var sizeNames = map[string]string{
	"TEXT":                "text",
	"HEADING_TEXT":        "headingText",
	"SUB_HEADING_TEXT":    "subHeadingText",
	"CAPTION_TEXT":        "captionText",
	"INLINE_ICON":         "inlineIcon",
	"PADDING":             "padding",
	"INNER_PADDING":       "innerPadding",
	"SCROLL_BAR":          "scrollBar",
	"SCROLL_BAR_SMALL":    "scrollBarSmall",
	"SEPARATOR_THICKNESS": "separatorThickness",
	"INPUT_RADIUS":        "inputRadius",
	"SELECTION_RADIUS":    "selectionRadius",
	"WINDOW_BUTTON_SIZE":  "windowButtonSize",
}

func (c *ThemeSizeNameClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := sizeNames[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}
