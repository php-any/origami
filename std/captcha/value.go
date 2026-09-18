package captcha

import (
	"net/http"

	"github.com/php-any/origami/data"
)

func valueAsString(v data.GetValue) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(data.AsString); ok {
		return s.AsString()
	}
	if val, ok := v.(data.Value); ok {
		return val.AsString()
	}
	return ""
}

func isNullish(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func argInt(ctx data.Context, index, def int) int {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return def
	}
	return valueInt(v, def)
}

func argString(ctx data.Context, index int, def string) string {
	v, ok := ctx.GetIndexValue(index)
	if !ok || isNullish(v) {
		return def
	}
	s := v.AsString()
	if s == "" && def != "" {
		if _, isStr := v.(*data.StringValue); !isStr {
			return def
		}
	}
	return s
}

func argBool(ctx data.Context, index int, def bool) bool {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return def
	}
	return valueBool(v, def)
}

func valueInt(v data.Value, def int) int {
	if v == nil || isNullish(v) {
		return def
	}
	if ai, ok := v.(data.AsInt); ok {
		if n, err := ai.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func valueBool(v data.Value, def bool) bool {
	if v == nil || isNullish(v) {
		return def
	}
	if ab, ok := v.(data.AsBool); ok {
		if b, err := ab.AsBool(); err == nil {
			return b
		}
	}
	return def
}

func propString(cv *data.ClassValue, name, def string) string {
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	if v == nil || isNullish(v) {
		return def
	}
	s := v.AsString()
	if s == "" {
		return def
	}
	return s
}

func propInt(cv *data.ClassValue, name string, def int) int {
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	return valueInt(v, def)
}

func propBool(cv *data.ClassValue, name string, def bool) bool {
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	if v == nil || isNullish(v) {
		return def
	}
	return valueBool(v, def)
}

func optionMap(v data.Value) map[string]data.Value {
	out := map[string]data.Value{}
	switch t := v.(type) {
	case *data.ArrayValue:
		for _, z := range t.List {
			if z == nil || z.Value == nil || z.Name == "" {
				continue
			}
			out[z.Name] = z.Value
		}
	case *data.ObjectValue:
		t.RangeProperties(func(k string, val data.Value) bool {
			out[k] = val
			return true
		})
	}
	return out
}

func optionValue(v data.Value, key string) (data.Value, bool) {
	m := optionMap(v)
	got, ok := m[key]
	return got, ok
}

func setImageHeader(ctx data.Context, contentType string) {
	if ctx == nil {
		return
	}
	type host interface {
		HTTPResponseWriter() http.ResponseWriter
	}
	if h, ok := ctx.GetVM().(host); ok {
		if w := h.HTTPResponseWriter(); w != nil {
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
		}
	}
}
