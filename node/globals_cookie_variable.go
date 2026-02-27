package node

import "github.com/php-any/origami/data"

// $_COOKIE

type CookieVariable struct {
	*Node `pp:"-"`
}

var cookieValue *data.ObjectValue

func NewCookieVariable(from data.From) data.Variable {
	return &CookieVariable{Node: NewNode(from)}
}

func (v *CookieVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if cookieValue == nil {
		cookieValue = data.NewObjectValue()
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			for _, cookie := range httpReq.Cookies() {
				cookieValue.SetProperty(cookie.Name, data.NewStringValue(cookie.Value))
			}
		}
	}
	return cookieValue, nil
}

func (v *CookieVariable) GetIndex() int       { return 0 }
func (v *CookieVariable) GetName() string     { return "$_COOKIE" }
func (v *CookieVariable) GetType() data.Types { return nil }
func (v *CookieVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}
