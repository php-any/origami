package node

import "github.com/php-any/origami/data"

// $_COOKIE

type CookieVariable struct {
	*Node `pp:"-"`
}

var cookieValue *data.ArrayValue

func NewCookieVariable(from data.From) data.Variable {
	return &CookieVariable{Node: NewNode(from)}
}

func (v *CookieVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_COOKIE", func() *data.ArrayValue {
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			obj := data.NewArrayValueFromSlots(nil)
			for _, cookie := range httpReq.Cookies() {
				obj.SetStringKey(cookie.Name, data.NewStringValue(cookie.Value))
			}
			return obj
		}
		return data.NewArrayValueFromSlots(nil)

	}), nil
}

func (v *CookieVariable) GetIndex() int       { return -1 }
func (v *CookieVariable) GetName() string     { return "$_COOKIE" }
func (v *CookieVariable) GetType() data.Types { return nil }
func (v *CookieVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_COOKIE", value)
}

func (v *CookieVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_COOKIE"), nil
}
func (v *CookieVariable) SuperglobalName() string { return "_COOKIE" }
