package pagination

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const lengthAwareName = "Illuminate\\Pagination\\LengthAwarePaginator"

type LengthAwarePaginatorClass struct {
	node.Node
	methods map[string]data.Method
}

func NewLengthAwarePaginatorClass() data.ClassStmt {
	c := &LengthAwarePaginatorClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *LengthAwarePaginatorClass) GetName() string { return lengthAwareName }
func (c *LengthAwarePaginatorClass) GetExtend() *string {
	s := abstractPaginatorName
	return &s
}
func (c *LengthAwarePaginatorClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Pagination\\LengthAwarePaginator"}
}
func (c *LengthAwarePaginatorClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "total", "lastPage":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *LengthAwarePaginatorClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 2)
	for _, n := range []string{"total", "lastPage"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *LengthAwarePaginatorClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *LengthAwarePaginatorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	_ = cv.SetProperty("path", data.NewStringValue("/"))
	_ = cv.SetProperty("pageName", data.NewStringValue("page"))
	_ = cv.SetProperty("query", data.NewArrayValue(nil))
	return cv, nil
}
func (c *LengthAwarePaginatorClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *LengthAwarePaginatorClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *LengthAwarePaginatorClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *LengthAwarePaginatorClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"items", "total", "perPage", "currentPage", "options"}, 3, laConstruct)
	c.methods["total"] = kit.InstanceMethod("total", nil, laTotal)
	c.methods["lastpage"] = kit.InstanceMethod("lastPage", nil, laLastPage)
	c.methods["hasmorepages"] = kit.InstanceMethod("hasMorePages", nil, laHasMorePages)
	c.methods["links"] = kit.InstanceMethodOpt("links", []string{"view", "data"}, 0, pagLinksStub)
	c.methods["toarray"] = kit.InstanceMethod("toArray", nil, laToArray)
}

func laConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := kit.Arg(ctx, 0)
	total := pagIntArg(ctx, 1, 0)
	perPage := pagIntArg(ctx, 2, 15)
	current := pagIntArg(ctx, 3, 1)
	if !pagValidPage(current) {
		current = 1
	}
	col, ctl := newCollection(ctx, items)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("items", col)
	_ = cv.SetProperty("total", data.NewIntValue(total))
	_ = cv.SetProperty("perPage", data.NewIntValue(perPage))
	_ = cv.SetProperty("currentPage", data.NewIntValue(current))
	last := pagLastPage(total, perPage)
	_ = cv.SetProperty("lastPage", data.NewIntValue(last))
	pagApplyOptions(cv, kit.Arg(ctx, 4))
	path := pagStrProp(cv, "path", "/")
	if path != "/" {
		_ = cv.SetProperty("path", data.NewStringValue(strings.TrimRight(path, "/")))
	}
	return cv, nil
}

func laTotal(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(pagIntProp(cv, "total", 0)), nil
}

func laLastPage(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(pagIntProp(cv, "lastPage", 1)), nil
}

func laHasMorePages(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cur := pagIntProp(cv, "currentPage", 1)
	last := pagIntProp(cv, "lastPage", 1)
	return data.NewBoolValue(cur < last), nil
}

func pagLinksStub(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(""), nil
}

func laToArray(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cur := pagIntProp(cv, "currentPage", 1)
	per := pagIntProp(cv, "perPage", 15)
	total := pagIntProp(cv, "total", 0)
	last := pagIntProp(cv, "lastPage", 1)
	col := pagCollectionProp(cv)
	dataItems := collectionToArray(ctx, col)
	path := pagStrProp(cv, "path", "/")
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	out.SetStringKey("current_page", data.NewIntValue(cur))
	out.SetStringKey("data", dataItems)
	out.SetStringKey("first_page_url", data.NewStringValue(pagBuildURL(path, pagStrProp(cv, "pageName", "page"), 1, pagQueryArray(cv), pagStrProp(cv, "fragment", ""))))
	firstItem, _ := pagFirstItem(ctx)
	out.SetStringKey("from", pagAsValue(firstItem))
	out.SetStringKey("last_page", data.NewIntValue(last))
	out.SetStringKey("last_page_url", data.NewStringValue(pagBuildURL(path, pagStrProp(cv, "pageName", "page"), last, pagQueryArray(cv), pagStrProp(cv, "fragment", ""))))
	out.SetStringKey("links", data.NewArrayValue(nil))
	next := ""
	if cur < last {
		next = pagBuildURL(path, pagStrProp(cv, "pageName", "page"), cur+1, pagQueryArray(cv), pagStrProp(cv, "fragment", ""))
	}
	out.SetStringKey("next_page_url", data.NewStringValue(next))
	out.SetStringKey("path", data.NewStringValue(path))
	out.SetStringKey("per_page", data.NewIntValue(per))
	prev := ""
	if cur > 1 {
		prev = pagBuildURL(path, pagStrProp(cv, "pageName", "page"), cur-1, pagQueryArray(cv), pagStrProp(cv, "fragment", ""))
	}
	out.SetStringKey("prev_page_url", data.NewStringValue(prev))
	lastItem, _ := pagLastItem(ctx)
	out.SetStringKey("to", pagAsValue(lastItem))
	out.SetStringKey("total", data.NewIntValue(total))
	return out, nil
}
