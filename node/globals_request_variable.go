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
	if requestValue == nil {
		requestValue = data.NewObjectValue()

		if getVal, _ := (&GetVariable{Node: v.Node}).GetValue(ctx); getVal != nil {
			if getObj, ok := getVal.(*data.ObjectValue); ok {
				getObj.RangeProperties(func(key string, value data.Value) bool {
					requestValue.SetProperty(key, value)
					return true
				})
			}
		}

		if postVal, _ := (&PostVariable{Node: v.Node}).GetValue(ctx); postVal != nil {
			if postObj, ok := postVal.(*data.ObjectValue); ok {
				postObj.RangeProperties(func(key string, value data.Value) bool {
					requestValue.SetProperty(key, value)
					return true
				})
			}
		}

		if cookieVal, _ := (&CookieVariable{Node: v.Node}).GetValue(ctx); cookieVal != nil {
			if cookieObj, ok := cookieVal.(*data.ObjectValue); ok {
				cookieObj.RangeProperties(func(key string, value data.Value) bool {
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
