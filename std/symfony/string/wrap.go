package sfstring

import (
	"github.com/php-any/origami/data"
)

func methodWrap(ctx data.Context) (data.GetValue, data.Control) {
	av := asArrayValue(argValue(ctx, 0))
	if av == nil {
		return data.NewArrayValue(nil), nil
	}
	return wrapArray(ctx, lateClass(ctx), av)
}

func methodUnwrap(ctx data.Context) (data.GetValue, data.Control) {
	av := asArrayValue(argValue(ctx, 0))
	if av == nil {
		return data.NewArrayValue(nil), nil
	}
	return unwrapArray(av), nil
}

func methodFromCodePoints(ctx data.Context) (data.GetValue, data.Control) {
	s := encodePHPCodePoints(variadicInts(ctx, 0))
	cls := lateClass(ctx)
	if sc, ok := cls.(*strClass); ok && sc.kind != kindByte {
		if ctl := mustUTF8(s); ctl != nil {
			return nil, ctl
		}
	}
	return newTyped(ctx, cls, s), nil
}

func newStaticString(ctx data.Context, cls data.ClassStmt, s string) (*data.ClassValue, data.Control) {
	if sc, ok := cls.(*strClass); ok && sc.kind != kindByte {
		if ctl := mustUTF8(s); ctl != nil {
			return nil, ctl
		}
	}
	return newTyped(ctx, cls, s), nil
}

func wrapArray(ctx data.Context, cls data.ClassStmt, av *data.ArrayValue) (*data.ArrayValue, data.Control) {
	list := make([]*data.ZVal, len(av.List))
	for i, z := range av.List {
		name := ""
		var val data.Value
		if z != nil {
			name = z.Name
			val = z.Value
		}
		if name != "" {
			if _, isInt := data.ParseIntArrayKeyName(name); !isInt {
				obj, ctl := newStaticString(ctx, cls, name)
				if ctl != nil {
					return nil, ctl
				}
				if j := getString(obj); j != name {
					name = j
				}
			}
		}
		if sv, ok := val.(*data.StringValue); ok {
			obj, ctl := newStaticString(ctx, cls, sv.Value)
			if ctl != nil {
				return nil, ctl
			}
			val = obj
		} else if nested := phpArray(val); nested != nil {
			w, ctl := wrapArray(ctx, cls, nested)
			if ctl != nil {
				return nil, ctl
			}
			val = w
		}
		if name != "" {
			list[i] = data.NewNamedZVal(name, val)
		} else {
			list[i] = data.NewZVal(val)
		}
	}
	return &data.ArrayValue{List: list}, nil
}

func unwrapArray(av *data.ArrayValue) *data.ArrayValue {
	list := make([]*data.ZVal, len(av.List))
	for i, z := range av.List {
		name := ""
		var val data.Value
		if z != nil {
			name = z.Name
			val = z.Value
		}
		if cv, ok := val.(*data.ClassValue); ok && isStringObj(cv) {
			val = data.NewStringValue(getString(cv))
		} else if nested := phpArray(val); nested != nil {
			val = unwrapArray(nested)
		}
		if name != "" {
			list[i] = data.NewNamedZVal(name, val)
		} else {
			list[i] = data.NewZVal(val)
		}
	}
	return &data.ArrayValue{List: list}
}

func phpArray(v data.Value) *data.ArrayValue {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return asArrayValue(v)
	default:
		return nil
	}
}

func asArrayValue(v data.Value) *data.ArrayValue {
	if v == nil {
		return nil
	}
	if zv, ok := v.(*data.ZValValue); ok && zv != nil && zv.ZVal != nil {
		return asArrayValue(zv.ZVal.Value)
	}
	switch t := v.(type) {
	case *data.ArrayValue:
		return t
	case *data.ObjectValue:
		out := &data.ArrayValue{}
		t.RangeProperties(func(key string, value data.Value) bool {
			out.List = append(out.List, data.NewNamedZVal(key, value))
			return true
		})
		return out
	}
	return nil
}

func encodePHPCodePoints(codes []int) string {
	var b []byte
	for _, code := range codes {
		code %= 0x200000
		if code < 0 {
			code += 0x200000
		}
		switch {
		case code < 0x80:
			b = append(b, byte(code))
		case code < 0x800:
			b = append(b, byte(0xC0|code>>6), byte(0x80|code&0x3F))
		case code < 0x10000:
			b = append(b, byte(0xE0|code>>12), byte(0x80|code>>6&0x3F), byte(0x80|code&0x3F))
		default:
			b = append(b, byte(0xF0|code>>18), byte(0x80|code>>12&0x3F), byte(0x80|code>>6&0x3F), byte(0x80|code&0x3F))
		}
	}
	return string(b)
}

func variadicInts(ctx data.Context, i int) []int {
	v, ok := ctx.GetIndexValue(i)
	if ok && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			out := make([]int, 0, len(av.List))
			for _, z := range av.List {
				if z == nil || z.Value == nil {
					continue
				}
				if iv, ok := z.Value.(data.AsInt); ok {
					n, err := iv.AsInt()
					if err == nil {
						out = append(out, n)
					}
				}
			}
			return out
		}
		if iv, ok := v.(data.AsInt); ok {
			n, err := iv.AsInt()
			if err == nil {
				return []int{n}
			}
		}
	}
	if g, ok := ctx.(interface{ GetFlatCallArgs() []data.Value }); ok {
		flat := g.GetFlatCallArgs()
		if len(flat) > i {
			out := make([]int, 0, len(flat)-i)
			for _, a := range flat[i:] {
				if a == nil {
					continue
				}
				if iv, ok := a.(data.AsInt); ok {
					n, err := iv.AsInt()
					if err == nil {
						out = append(out, n)
					}
				}
			}
			return out
		}
	}
	return nil
}
