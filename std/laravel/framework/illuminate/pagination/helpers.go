package pagination

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

func pagIntArg(ctx data.Context, i int, def int) int {
	v := kit.Unwrap(kit.Arg(ctx, i))
	if v == nil || kit.IsNull(v) {
		return def
	}
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	s := strings.TrimSpace(v.AsString())
	if s == "" {
		return def
	}
	n, _ := strconv.Atoi(s)
	if n == 0 && s != "0" {
		return def
	}
	return n
}

func pagRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Paginator missing $this"))
}

func pagIntProp(cv *data.ClassValue, name string, def int) int {
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	u := kit.Unwrap(v)
	if iv, ok := u.(*data.IntValue); ok {
		return iv.Value
	}
	n, _ := strconv.Atoi(u.AsString())
	return n
}

func pagStrProp(cv *data.ClassValue, name string, def string) string {
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	s := kit.Unwrap(v).AsString()
	if s == "" {
		return def
	}
	return s
}

func pagCollectionProp(cv *data.ClassValue) *data.ClassValue {
	v, _ := cv.GetProperty("items")
	if cvItems, ok := kit.Unwrap(v).(*data.ClassValue); ok {
		return cvItems
	}
	return nil
}

func newCollection(ctx data.Context, items data.Value) (*data.ClassValue, data.Control) {
	vm := ctx.GetVM()
	cls, ok := vm.GetClass("Illuminate\\Support\\Collection")
	if !ok || cls == nil {
		var ctl data.Control
		cls, ctl = vm.GetOrLoadClass("Illuminate\\Support\\Collection")
		if ctl != nil {
			return nil, ctl
		}
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	if ctor := cls.GetConstruct(); ctor != nil {
		nctx := cv.CreateContext(ctor.GetVariables())
		data.BindDeclaredArgs(nctx, ctor, []data.Value{items})
		if _, ctl := ctor.Call(nctx); ctl != nil {
			return nil, ctl
		}
	} else {
		_ = cv.SetProperty("items", items)
	}
	return cv, nil
}

func collectionCount(ctx data.Context, col *data.ClassValue) int {
	if col == nil {
		return 0
	}
	ret, ctl := kit.CallInstanceMethod(ctx, col, "count")
	if ctl != nil {
		return 0
	}
	if iv, ok := kit.Unwrap(ret.(data.Value)).(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func collectionToArray(ctx data.Context, col *data.ClassValue) data.Value {
	if col == nil {
		return data.NewArrayValue(nil)
	}
	ret, ctl := kit.CallInstanceMethod(ctx, col, "all")
	if ctl != nil {
		return data.NewArrayValue(nil)
	}
	if v, ok := ret.(data.Value); ok {
		return v
	}
	return data.NewArrayValue(nil)
}

func collectionSlice(ctx data.Context, col *data.ClassValue, offset, length int) (*data.ClassValue, data.Control) {
	if col == nil {
		return newCollection(ctx, data.NewArrayValue(nil))
	}
	ret, ctl := kit.CallInstanceMethod(ctx, col, "slice", data.NewIntValue(offset), data.NewIntValue(length))
	if ctl != nil {
		return nil, ctl
	}
	if cv, ok := kit.Unwrap(ret.(data.Value)).(*data.ClassValue); ok {
		return cv, nil
	}
	return newCollection(ctx, data.NewArrayValue(nil))
}

func pagApplyOptions(cv *data.ClassValue, options data.Value) {
	av, ok := kit.Unwrap(options).(*data.ArrayValue)
	if !ok || av == nil {
		return
	}
	for _, e := range av.List {
		if e == nil || e.Name == "" {
			continue
		}
		_ = cv.SetProperty(e.Name, e.Value)
	}
}

func pagValidPage(n int) bool {
	return n >= 1
}

func pagBuildURL(path, pageName string, page int, query *data.ArrayValue, fragment string) string {
	if page <= 0 {
		page = 1
	}
	params := url.Values{}
	if query != nil {
		for _, e := range query.List {
			if e == nil || e.Name == "" {
				continue
			}
			params.Set(e.Name, kit.Unwrap(e.Value).AsString())
		}
	}
	params.Set(pageName, strconv.Itoa(page))
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	u := path + sep + params.Encode()
	if fragment != "" {
		u += "#" + fragment
	}
	return u
}

func pagQueryArray(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("query")
	if av, ok := kit.Unwrap(v).(*data.ArrayValue); ok && av != nil {
		return av
	}
	av := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("query", av)
	return av
}

func pagAsValue(v data.GetValue) data.Value {
	if v == nil {
		return data.NewNullValue()
	}
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}

func pagLastPage(total, perPage int) int {
	if perPage <= 0 {
		return 1
	}
	return int(math.Max(math.Ceil(float64(total)/float64(perPage)), 1))
}
