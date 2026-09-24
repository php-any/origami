package pagination

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const paginatorName = "Illuminate\\Pagination\\Paginator"

type PaginatorClass struct {
	node.Node
	methods map[string]data.Method
}

func NewPaginatorClass() data.ClassStmt {
	c := &PaginatorClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *PaginatorClass) GetName() string { return paginatorName }
func (c *PaginatorClass) GetExtend() *string {
	s := abstractPaginatorName
	return &s
}
func (c *PaginatorClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Pagination\\Paginator"}
}
func (c *PaginatorClass) GetProperty(name string) (data.Property, bool) {
	if name == "hasMore" {
		return node.NewProperty(nil, name, "protected", false, data.NewBoolValue(false)), true
	}
	return nil, false
}
func (c *PaginatorClass) GetPropertyList() []data.Property {
	p, _ := c.GetProperty("hasMore")
	return []data.Property{p}
}
func (c *PaginatorClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *PaginatorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	_ = cv.SetProperty("path", data.NewStringValue("/"))
	_ = cv.SetProperty("pageName", data.NewStringValue("page"))
	_ = cv.SetProperty("query", data.NewArrayValue(nil))
	return cv, nil
}
func (c *PaginatorClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *PaginatorClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *PaginatorClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *PaginatorClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"items", "perPage", "currentPage", "options"}, 2, simpleConstruct)
	c.methods["hasmorepages"] = kit.InstanceMethod("hasMorePages", nil, simpleHasMorePages)
	c.methods["links"] = kit.InstanceMethodOpt("links", []string{"view", "data"}, 0, pagLinksStub)
	c.methods["toarray"] = kit.InstanceMethod("toArray", nil, simpleToArray)
}

func simpleConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := kit.Arg(ctx, 0)
	perPage := pagIntArg(ctx, 1, 15)
	current := pagIntArg(ctx, 2, 1)
	if !pagValidPage(current) {
		current = 1
	}
	col, ctl := newCollection(ctx, items)
	if ctl != nil {
		return nil, ctl
	}
	cnt := collectionCount(ctx, col)
	hasMore := cnt > perPage
	sliced, ctl := collectionSlice(ctx, col, 0, perPage)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("items", sliced)
	_ = cv.SetProperty("perPage", data.NewIntValue(perPage))
	_ = cv.SetProperty("currentPage", data.NewIntValue(current))
	_ = cv.SetProperty("hasMore", data.NewBoolValue(hasMore))
	pagApplyOptions(cv, kit.Arg(ctx, 3))
	path := pagStrProp(cv, "path", "/")
	if path != "/" {
		_ = cv.SetProperty("path", data.NewStringValue(strings.TrimRight(path, "/")))
	}
	return cv, nil
}

func simpleHasMorePages(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := cv.GetProperty("hasMore")
	if b, ok := kit.Unwrap(v).(*data.BoolValue); ok {
		return b, nil
	}
	return data.NewBoolValue(false), nil
}

func simpleToArray(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cur := pagIntProp(cv, "currentPage", 1)
	per := pagIntProp(cv, "perPage", 15)
	col := pagCollectionProp(cv)
	dataItems := collectionToArray(ctx, col)
	path := pagStrProp(cv, "path", "/")
	pageName := pagStrProp(cv, "pageName", "page")
	frag := pagStrProp(cv, "fragment", "")
	q := pagQueryArray(cv)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	out.SetStringKey("current_page", data.NewIntValue(cur))
	out.SetStringKey("current_page_url", data.NewStringValue(pagBuildURL(path, pageName, cur, q, frag)))
	out.SetStringKey("data", dataItems)
	out.SetStringKey("first_page_url", data.NewStringValue(pagBuildURL(path, pageName, 1, q, frag)))
	firstItem, _ := pagFirstItem(ctx)
	out.SetStringKey("from", pagAsValue(firstItem))
	next := ""
	if v, _ := cv.GetProperty("hasMore"); kit.Truthy(kit.Unwrap(v)) {
		next = pagBuildURL(path, pageName, cur+1, q, frag)
	}
	out.SetStringKey("next_page_url", data.NewStringValue(next))
	out.SetStringKey("path", data.NewStringValue(path))
	out.SetStringKey("per_page", data.NewIntValue(per))
	prev := ""
	if cur > 1 {
		prev = pagBuildURL(path, pageName, cur-1, q, frag)
	}
	out.SetStringKey("prev_page_url", data.NewStringValue(prev))
	lastItem, _ := pagLastItem(ctx)
	out.SetStringKey("to", pagAsValue(lastItem))
	return out, nil
}
