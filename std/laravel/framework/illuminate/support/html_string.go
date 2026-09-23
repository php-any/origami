package support

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const htmlStringClassName = "Illuminate\\Support\\HtmlString"

type HtmlStringClass struct {
	node.Node
	methods map[string]data.Method
}

func NewHtmlStringClass() data.ClassStmt {
	c := &HtmlStringClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *HtmlStringClass) GetName() string                          { return htmlStringClassName }
func (c *HtmlStringClass) GetExtend() *string                       { return nil }
func (c *HtmlStringClass) GetImplements() []string                  {
	return []string{"Illuminate\\Contracts\\Support\\Htmlable", "Stringable"}
}
func (c *HtmlStringClass) GetProperty(name string) (data.Property, bool) {
	if name == "html" {
		return node.NewProperty(nil, "html", "protected", false, data.NewStringValue("")), true
	}
	return nil, false
}
func (c *HtmlStringClass) GetPropertyList() []data.Property {
	p, _ := c.GetProperty("html")
	return []data.Property{p}
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
func (c *HtmlStringClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *HtmlStringClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"html"}, 0, htmlStringConstruct)
	c.methods["tohtml"] = kit.InstanceMethod("toHtml", nil, htmlStringToHtml)
	c.methods["isempty"] = kit.InstanceMethod("isEmpty", nil, htmlStringIsEmpty)
	c.methods["isnotempty"] = kit.InstanceMethod("isNotEmpty", nil, htmlStringIsNotEmpty)
	c.methods["__tostring"] = kit.InstanceMethod("__toString", nil, htmlStringToHtml)
}

func htmlStringConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	html := ""
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		html = v.AsString()
	}
	_ = cv.SetProperty("html", data.NewStringValue(html))
	return data.NewNullValue(), nil
}

func htmlStringHtml(cv *data.ClassValue) string {
	v, _ := cv.GetProperty("html")
	if v == nil {
		return ""
	}
	return v.AsString()
}

func htmlStringToHtml(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(htmlStringHtml(cv)), nil
}

func htmlStringIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(htmlStringHtml(cv) == ""), nil
}

func htmlStringIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	b, ctl := htmlStringIsEmpty(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := b.(data.AsBool).AsBool()
	return data.NewBoolValue(!v), nil
}
