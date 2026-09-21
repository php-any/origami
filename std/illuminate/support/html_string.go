package support

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const htmlStringClassName = "Illuminate\\Support\\HtmlString"

// HtmlStringClass 对齐 Laravel HtmlString（体积小、接口固定，可整类替换）。
type HtmlStringClass struct {
	node.Node
	methods map[string]data.Method
}

func NewHtmlStringClass() data.ClassStmt {
	c := &HtmlStringClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = newInstanceMethod("__construct", []string{"html"}, htmlStringConstruct, false)
	c.methods["tohtml"] = newInstanceMethod("toHtml", nil, htmlStringToHtml, false)
	c.methods["isempty"] = newInstanceMethod("isEmpty", nil, htmlStringIsEmpty, false)
	c.methods["isnotempty"] = newInstanceMethod("isNotEmpty", nil, htmlStringIsNotEmpty, false)
	c.methods["__tostring"] = newInstanceMethod("__toString", nil, htmlStringToHtml, false)
	return c
}

func (c *HtmlStringClass) GetName() string { return htmlStringClassName }
func (c *HtmlStringClass) GetExtend() *string {
	return nil
}
func (c *HtmlStringClass) GetImplements() []string {
	return []string{
		"Illuminate\\Contracts\\Support\\Htmlable",
		"Stringable",
	}
}
func (c *HtmlStringClass) GetProperty(name string) (data.Property, bool) {
	if name == "html" {
		return node.NewProperty(nil, "html", "protected", false, data.NewStringValue("")), true
	}
	return nil, false
}
func (c *HtmlStringClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "html", "protected", false, data.NewStringValue("")),
	}
}
func (c *HtmlStringClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *HtmlStringClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *HtmlStringClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *HtmlStringClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *HtmlStringClass) GetStaticMethod(string) (data.Method, bool) {
	return nil, false
}

func htmlStringConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := overlayReceiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	html := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil && !isNull(v) {
		html = v.AsString()
	}
	_ = cv.SetProperty("html", data.NewStringValue(html))
	return cv, nil
}

func htmlStringToHtml(ctx data.Context) (data.GetValue, data.Control) {
	cv := overlayReceiver(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	v, ctl := cv.GetProperty("html")
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(v.AsString()), nil
}

func htmlStringIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := htmlStringToHtml(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(s.(data.Value).AsString() == ""), nil
}

func htmlStringIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	empty, ctl := htmlStringIsEmpty(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := empty.(*data.BoolValue); ok {
		return data.NewBoolValue(!bv.Value), nil
	}
	return data.NewBoolValue(true), nil
}
