package node

import "github.com/php-any/origami/data"

// $_REQUEST

type RequestVariable struct {
	*Node `pp:"-"`
}

var requestValue *data.ObjectValue

func NewRequestVariable(from data.From) data.Variable {
	return &RequestVariable{Node: NewNode(from)}
}

func (v *RequestVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if httpReq := getHTTPRequest(ctx); httpReq != nil {
		obj := data.NewObjectValue()
		merge := func(src data.GetValue) {
			if src == nil {
				return
			}
			if o, ok := src.(*data.ObjectValue); ok {
				o.RangeProperties(func(key string, value data.Value) bool {
					obj.SetProperty(key, value)
					return true
				})
			}
		}
		getVal, _ := (&GetVariable{Node: v.Node}).GetValue(ctx)
		postVal, _ := (&PostVariable{Node: v.Node}).GetValue(ctx)
		cookieVal, _ := (&CookieVariable{Node: v.Node}).GetValue(ctx)
		merge(getVal)
		merge(postVal)
		merge(cookieVal)
		return obj, nil
	}
	if requestValue == nil {
		requestValue = data.NewObjectValue()
		getVal, _ := (&GetVariable{Node: v.Node}).GetValue(ctx)
		postVal, _ := (&PostVariable{Node: v.Node}).GetValue(ctx)
		cookieVal, _ := (&CookieVariable{Node: v.Node}).GetValue(ctx)
		for _, src := range []data.GetValue{getVal, postVal, cookieVal} {
			if o, ok := src.(*data.ObjectValue); ok {
				o.RangeProperties(func(key string, value data.Value) bool {
					requestValue.SetProperty(key, value)
					return true
				})
			}
		}
	}
	return requestValue, nil
}

func (v *RequestVariable) GetIndex() int       { return 0 }
func (v *RequestVariable) GetName() string     { return "$_REQUEST" }
func (v *RequestVariable) GetType() data.Types { return nil }
func (v *RequestVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}
