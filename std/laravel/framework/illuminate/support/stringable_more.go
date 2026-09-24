package support

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/illuminate/conditionable"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

var stringableStrSingleton data.ClassStmt

func strClassForStringable() *StrClass {
	if stringableStrSingleton == nil {
		stringableStrSingleton = NewStrClass()
	}
	return stringableStrSingleton.(*StrClass)
}

func (c *StringableClass) registerStringableMore() {
	inst := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethod(name, params, fn)
	}
	inst("append", nil, stringableAppend)
	inst("prepend", nil, stringablePrepend)
	inst("newline", []string{"count"}, stringableNewLine)
	inst("finish", []string{"cap"}, stringableFinish)
	inst("start", []string{"prefix"}, stringableStart)
	inst("pipe", []string{"callback"}, stringablePipe)
	inst("tap", []string{"callback"}, stringableTap)
	inst("value", nil, stringableToString)
	inst("isempty", nil, stringableIsEmpty)
	inst("isnotempty", nil, stringableIsNotEmpty)
	inst("basename", []string{"suffix"}, stringableBasename)
	inst("offsetexists", []string{"offset"}, stringableOffsetExists)
	inst("offsetget", []string{"offset"}, stringableOffsetGet)
	inst("offsetset", []string{"offset", "value"}, stringableOffsetSet)
	inst("offsetunset", []string{"offset"}, stringableOffsetUnset)
	inst("__get", []string{"key"}, stringableGet)
	inst("tohtmlstring", nil, stringableToHtmlString)
	c.methods["explode"] = kit.InstanceMethodOpt("explode", []string{"delimiter", "limit"}, 1, stringableExplode)

	kit.RegisterConditionable(c.methods, conditionable.NewWhenProxy)

	valueFirst := []struct {
		name   string
		params []string
	}{
		{"afterlast", []string{"search"}},
		{"ascii", []string{"language"}},
		{"between", []string{"from", "to"}},
		{"betweenfirst", []string{"from", "to"}},
		{"charat", []string{"index"}},
		{"chopstart", []string{"needle"}},
		{"chopend", []string{"needle"}},
		{"containsall", []string{"needles", "ignoreCase"}},
		{"doesntcontain", []string{"needles", "ignoreCase"}},
		{"convertcase", []string{"mode", "encoding"}},
		{"counted", []string{"count"}},
		{"deduplicate", []string{"characters"}},
		{"doesntendwith", []string{"needles"}},
		{"doesntstartwith", []string{"needles"}},
		{"endswith", []string{"needles"}},
		{"excerpt", []string{"phrase", "options"}},
		{"headline", nil},
		{"initials", []string{"capitalize"}},
		{"apa", nil},
		{"singular", nil},
		{"isascii", nil},
		{"isjson", nil},
		{"isurl", []string{"protocols"}},
		{"isuuid", []string{"version"}},
		{"isulid", nil},
		{"markdown", []string{"options", "extensions"}},
		{"inlinemarkdown", []string{"options", "extensions"}},
		{"mask", []string{"character", "index", "length", "encoding"}},
		{"numbers", nil},
		{"padboth", []string{"length", "pad"}},
		{"padleft", []string{"length", "pad"}},
		{"padright", []string{"length", "pad"}},
		{"plural", []string{"count", "prependCount"}},
		{"pluralstudly", []string{"count"}},
		{"pluralpascal", []string{"count"}},
		{"repeat", []string{"times"}},
		{"reverse", nil},
		{"squish", nil},
		{"startswith", []string{"needles"}},
		{"pascal", []string{"normalize"}},
		{"title", nil},
		{"take", []string{"limit"}},
		{"tobase64", nil},
		{"frombase64", []string{"strict"}},
		{"transliterate", []string{"unknown", "strict"}},
		{"lcfirst", nil},
		{"ucfirst", nil},
		{"ucwords", []string{"separators"}},
		{"ucsplit", nil},
		{"words", []string{"words", "end"}},
		{"wordcount", []string{"characters"}},
		{"wordwrap", []string{"characters", "break", "cutLongWords"}},
		{"wrap", []string{"before", "after"}},
		{"unwrap", []string{"before", "after"}},
		{"ltrim", []string{"characters"}},
		{"rtrim", []string{"characters"}},
		{"replacefirst", []string{"search", "replace"}},
		{"replacelast", []string{"search", "replace"}},
		{"replacestart", []string{"search", "replace"}},
		{"replaceend", []string{"search", "replace"}},
		{"replacearray", []string{"search", "replace"}},
		{"replacematches", []string{"pattern", "replace", "limit"}},
		{"remove", []string{"search", "caseSensitive"}},
		{"replace", []string{"search", "replace", "caseSensitive"}},
		{"match", []string{"pattern"}},
		{"matchall", []string{"pattern"}},
		{"ismatch", []string{"pattern"}},
		{"test", []string{"pattern"}},
		{"position", []string{"needle", "offset", "encoding"}},
		{"substrcount", []string{"needle", "offset", "length"}},
		{"substrreplace", []string{"replace", "offset", "length"}},
		{"swap", []string{"map"}},
		{"parsecallback", []string{"default"}},
	}
	for _, m := range valueFirst {
		name := m.name
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethod(name, m.params, func(ctx data.Context) (data.GetValue, data.Control) {
			return stringableInvokeStr(ctx, name, stringableCallArgs(ctx))
		})
	}
}

func stringableCallArgs(ctx data.Context) []data.Value {
	if args := ctx.GetFlatCallArgs(); args != nil {
		return args
	}
	var out []data.Value
	for i := 0; i < 12; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok || v == nil || kit.IsNull(v) {
			continue
		}
		out = append(out, v)
	}
	return out
}

func stringableBuildStrArgs(method, value string, extra []data.Value) []data.Value {
	m := strings.ToLower(method)
	val := data.NewStringValue(value)
	switch m {
	case "is", "ismatch":
		if len(extra) == 0 {
			return []data.Value{data.NewStringValue(""), val}
		}
		out := []data.Value{extra[0], val}
		if len(extra) > 1 {
			out = append(out, extra[1:]...)
		}
		return out
	case "replace", "replacefirst", "replacelast", "replacestart", "replaceend", "replacearray", "replacematches":
		out := append([]data.Value{}, extra...)
		out = append(out, val)
		return out
	case "remove":
		if len(extra) == 0 {
			return []data.Value{data.NewStringValue(""), val}
		}
		out := []data.Value{extra[0], val}
		if len(extra) > 1 {
			out = append(out, extra[1:]...)
		}
		return out
	case "match", "matchall", "test":
		if len(extra) == 0 {
			return []data.Value{data.NewStringValue(""), val}
		}
		return []data.Value{extra[0], val}
	case "swap":
		if len(extra) == 0 {
			return []data.Value{data.NewArrayValue(nil), val}
		}
		return []data.Value{extra[0], val}
	default:
		out := []data.Value{val}
		return append(out, extra...)
	}
}

func stringableInvokeStr(ctx data.Context, method string, extra []data.Value) (data.GetValue, data.Control) {
	value := stringableValue(ctx)
	args := stringableBuildStrArgs(method, value, extra)
	strStmt := strClassForStringable()
	m, ok := strStmt.GetStaticMethod(method)
	if !ok || m == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s::%s does not exist.", stringableClassName, method))
	}
	nctx := ctx.CreateContext(m.GetVariables())
	data.BindDeclaredArgs(nctx, m, args)
	ret, ctl := m.Call(nctx)
	if ctl != nil {
		return nil, ctl
	}
	return stringableCoerceReturn(ctx, ret)
}

func stringableCoerceReturn(ctx data.Context, ret data.GetValue) (data.GetValue, data.Control) {
	if ret == nil {
		return data.NewNullValue(), nil
	}
	v, ok := ret.(data.Value)
	if !ok {
		return ret, nil
	}
	switch v.(type) {
	case *data.StringValue:
		return stringableWrap(ctx, v.AsString())
	default:
		if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok && cv.Class != nil && cv.Class.GetName() == stringableClassName {
			return cv, nil
		}
		return v, nil
	}
}

func stringableExplode(ctx data.Context) (data.GetValue, data.Control) {
	delimiter := ""
	if v := kit.Arg(ctx, 0); v != nil {
		delimiter = v.AsString()
	}
	limit := int(^uint(0) >> 1)
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		if iv, ok := v.(*data.IntValue); ok {
			limit = iv.Value
		} else if b, ok := v.(data.AsInt); ok {
			if n, err := b.AsInt(); err == nil {
				limit = n
			}
		}
	}
	parts := explodeLimited(stringableValue(ctx), delimiter, limit)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for i, p := range parts {
		out.SetIntKey(i, data.NewStringValue(p))
	}
	// 对齐 PHP：return new Collection(...)
	vm := ctx.GetVM()
	stmt, ctl := vm.GetOrLoadClass("Illuminate\\Support\\Collection")
	if ctl != nil {
		return out, nil
	}
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	_ = cv.SetProperty("items", out)
	return cv, nil
}

func stringableAppend(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	for _, a := range stringableCallArgs(ctx) {
		if a != nil {
			v += a.AsString()
		}
	}
	return stringableWrap(ctx, v)
}

func stringablePrepend(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	prefix := ""
	for _, a := range stringableCallArgs(ctx) {
		if a != nil {
			prefix += a.AsString()
		}
	}
	return stringableWrap(ctx, prefix+v)
}

func stringableNewLine(ctx data.Context) (data.GetValue, data.Control) {
	count := 1
	if a := kit.Arg(ctx, 0); a != nil && !kit.IsNull(a) {
		count = strIntArgFromValue(a, 1)
	}
	if count < 1 {
		count = 1
	}
	v := stringableValue(ctx) + strings.Repeat("\n", count)
	return stringableWrap(ctx, v)
}

func stringableFinish(ctx data.Context) (data.GetValue, data.Control) {
	return stringableInvokeStr(ctx, "finish", stringableCallArgs(ctx))
}

func stringableStart(ctx data.Context) (data.GetValue, data.Control) {
	return stringableInvokeStr(ctx, "start", stringableCallArgs(ctx))
}

func stringablePipe(ctx data.Context) (data.GetValue, data.Control) {
	recv := kit.Receiver(ctx)
	cb := kit.Arg(ctx, 0)
	if cb == nil {
		return recv, nil
	}
	ret, ctl := kit.Call(ctx, cb, recv)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil {
		return stringableWrap(ctx, "")
	}
	if v, ok := ret.(data.Value); ok {
		if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok && cv.Class != nil && cv.Class.GetName() == stringableClassName {
			return cv, nil
		}
		return stringableWrap(ctx, v.AsString())
	}
	return stringableWrap(ctx, "")
}

func stringableTap(ctx data.Context) (data.GetValue, data.Control) {
	recv := kit.Receiver(ctx)
	cb := kit.Arg(ctx, 0)
	if cb != nil {
		_, ctl := kit.Call(ctx, cb, recv)
		if ctl != nil {
			return nil, ctl
		}
	}
	return recv, nil
}

func stringableIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(stringableValue(ctx) == ""), nil
}

func stringableIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(stringableValue(ctx) != ""), nil
}

func stringableBasename(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	suffix := ""
	if a := kit.Arg(ctx, 0); a != nil && !kit.IsNull(a) {
		suffix = a.AsString()
	}
	if suffix != "" && strings.HasSuffix(v, suffix) {
		v = v[:len(v)-len(suffix)]
	}
	return stringableWrap(ctx, filepath.Base(v))
}

func stringableOffsetExists(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	off := kit.Arg(ctx, 0)
	if off == nil {
		return data.NewBoolValue(false), nil
	}
	idx := strIntArgFromValue(off, -1)
	return data.NewBoolValue(idx >= 0 && idx < len(v)), nil
}

func stringableOffsetGet(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	off := kit.Arg(ctx, 0)
	if off == nil {
		return data.NewStringValue(""), nil
	}
	idx := strIntArgFromValue(off, -1)
	if idx < 0 || idx >= len(v) {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(v[idx])), nil
}

func stringableOffsetSet(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	v := stringableValue(ctx)
	off := kit.Arg(ctx, 0)
	val := kit.Arg(ctx, 1)
	if off == nil || val == nil {
		return data.NewNullValue(), nil
	}
	idx := strIntArgFromValue(off, 0)
	ch := val.AsString()
	if len(ch) == 0 {
		return data.NewNullValue(), nil
	}
	r := []rune(v)
	for len(r) <= idx {
		r = append(r, ' ')
	}
	r[idx] = []rune(ch)[0]
	_ = cv.SetProperty("value", data.NewStringValue(string(r)))
	return data.NewNullValue(), nil
}

func stringableOffsetUnset(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func stringableGet(ctx data.Context) (data.GetValue, data.Control) {
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	recv := kit.Receiver(ctx)
	if recv == nil || key == "" {
		return data.NewNullValue(), nil
	}
	return kit.CallInstanceMethod(ctx, recv, key)
}

func stringableToHtmlString(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(NewHtmlStringClass(), ctx.CreateBaseContext())
	_ = cv.SetProperty("html", data.NewStringValue(stringableValue(ctx)))
	return cv, nil
}
