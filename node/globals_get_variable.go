package node

import "github.com/php-any/origami/data"

// $_GET

type GetVariable struct {
	*Node `pp:"-"`
}

var getValue *data.ObjectValue

func NewGetVariable(from data.From) data.Variable {
	return &GetVariable{Node: NewNode(from)}
}

func (v *GetVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if getValue == nil {
		getValue = data.NewObjectValue()
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			for key, values := range httpReq.URL.Query() {
				if len(values) > 0 {
					getValue.SetProperty(key, data.NewStringValue(values[0]))
				}
			}
		}
	}
	return getValue, nil
}

func (v *GetVariable) GetIndex() int       { return 0 }
func (v *GetVariable) GetName() string     { return "$_GET" }
func (v *GetVariable) GetType() data.Types { return nil }
func (v *GetVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}
