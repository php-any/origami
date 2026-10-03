package node

import "github.com/php-any/origami/data"

// $_GET

type GetVariable struct {
	*Node `pp:"-"`
}

var getValue *data.ArrayValue

func NewGetVariable(from data.From) data.Variable {
	return &GetVariable{Node: NewNode(from)}
}

func (v *GetVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_GET", func() *data.ArrayValue {
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			return data.ParseFormFields(httpReq.URL.RawQuery)
		}
		return data.NewArrayValueFromSlots(nil)

	}), nil
}

func (v *GetVariable) GetIndex() int       { return -1 }
func (v *GetVariable) GetName() string     { return "$_GET" }
func (v *GetVariable) GetType() data.Types { return nil }
func (v *GetVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_GET", value)
}

func (v *GetVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_GET"), nil
}
func (v *GetVariable) SuperglobalName() string { return "_GET" }
