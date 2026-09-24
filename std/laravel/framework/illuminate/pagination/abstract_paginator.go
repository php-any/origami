package pagination

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const abstractPaginatorName = "Illuminate\\Pagination\\AbstractPaginator"

type AbstractPaginatorClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAbstractPaginatorClass() data.ClassStmt {
	c := &AbstractPaginatorClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *AbstractPaginatorClass) GetName() string    { return abstractPaginatorName }
func (c *AbstractPaginatorClass) GetExtend() *string { return nil }
func (c *AbstractPaginatorClass) GetImplements() []string {
	return nil
}
func (c *AbstractPaginatorClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "items", "perPage", "currentPage", "path", "query", "fragment", "pageName", "options", "onEachSide":
		def := data.NewNullValue()
		if name == "path" {
			def = data.NewStringValue("/")
		}
		if name == "pageName" {
			def = data.NewStringValue("page")
		}
		if name == "onEachSide" {
			def = data.NewIntValue(3)
		}
		return node.NewProperty(nil, name, "protected", false, def), true
	}
	return nil, false
}
func (c *AbstractPaginatorClass) GetPropertyList() []data.Property {
	names := []string{"items", "perPage", "currentPage", "path", "query", "fragment", "pageName", "options", "onEachSide"}
	out := make([]data.Property, 0, len(names))
	for _, n := range names {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *AbstractPaginatorClass) GetConstruct() data.Method { return nil }
func (c *AbstractPaginatorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *AbstractPaginatorClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *AbstractPaginatorClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *AbstractPaginatorClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *AbstractPaginatorClass) register() {
	c.methods["currentpage"] = kit.InstanceMethod("currentPage", nil, pagCurrentPage)
	c.methods["perpage"] = kit.InstanceMethod("perPage", nil, pagPerPage)
	c.methods["path"] = kit.InstanceMethod("path", nil, pagPath)
	c.methods["url"] = kit.InstanceMethod("url", []string{"page"}, pagURL)
	c.methods["items"] = kit.InstanceMethod("items", nil, pagItems)
	c.methods["count"] = kit.InstanceMethod("count", nil, pagCount)
	c.methods["firstitem"] = kit.InstanceMethod("firstItem", nil, pagFirstItem)
	c.methods["lastitem"] = kit.InstanceMethod("lastItem", nil, pagLastItem)
	c.methods["getcollection"] = kit.InstanceMethod("getCollection", nil, pagGetCollection)
}

func pagCurrentPage(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(pagIntProp(cv, "currentPage", 1)), nil
}

func pagPerPage(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(pagIntProp(cv, "perPage", 15)), nil
}

func pagPath(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(pagStrProp(cv, "path", "/")), nil
}

func pagURL(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	page := pagIntArg(ctx, 0, 1)
	path := pagStrProp(cv, "path", "/")
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	pageName := pagStrProp(cv, "pageName", "page")
	frag := pagStrProp(cv, "fragment", "")
	return data.NewStringValue(pagBuildURL(path, pageName, page, pagQueryArray(cv), frag)), nil
}

func pagItems(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	col := pagCollectionProp(cv)
	return collectionToArray(ctx, col), nil
}

func pagCount(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(collectionCount(ctx, pagCollectionProp(cv))), nil
}

func pagFirstItem(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if collectionCount(ctx, pagCollectionProp(cv)) == 0 {
		return data.NewNullValue(), nil
	}
	cur := pagIntProp(cv, "currentPage", 1)
	per := pagIntProp(cv, "perPage", 15)
	return data.NewIntValue((cur-1)*per + 1), nil
}

func pagLastItem(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cnt := collectionCount(ctx, pagCollectionProp(cv))
	if cnt == 0 {
		return data.NewNullValue(), nil
	}
	first, _ := pagFirstItem(ctx)
	firstN := pagIntProp(cv, "currentPage", 1)
	if iv, ok := first.(*data.IntValue); ok {
		firstN = iv.Value
	}
	return data.NewIntValue(firstN + cnt - 1), nil
}

func pagGetCollection(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := pagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	col := pagCollectionProp(cv)
	if col != nil {
		return col, nil
	}
	return newCollection(ctx, data.NewArrayValue(nil))
}
