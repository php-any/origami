package node

import "github.com/php-any/origami/data"

// $_REQUEST

type RequestVariable struct {
	*Node `pp:"-"`
}

var requestValue *data.ArrayValue

func NewRequestVariable(from data.From) data.Variable {
	return &RequestVariable{Node: NewNode(from)}
}

func (v *RequestVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_REQUEST", func() *data.ArrayValue {
		result := data.NewArrayValueFromSlots(nil)
		get, _ := (&GetVariable{Node: v.Node}).GetValue(ctx)
		post, _ := (&PostVariable{Node: v.Node}).GetValue(ctx)
		for _, source := range []data.GetValue{get, post} {
			if array, ok := source.(*data.ArrayValue); ok {
				for position, slot := range array.Range() {
					result.SetKey(slot.PHPArrayKey(position), data.CowAddRef(slot.ReadValue()))
				}
			}
		}
		return result
	}), nil
}

func (v *RequestVariable) GetIndex() int       { return -1 }
func (v *RequestVariable) GetName() string     { return "$_REQUEST" }
func (v *RequestVariable) GetType() data.Types { return nil }
func (v *RequestVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_REQUEST", value)
}

func (v *RequestVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_REQUEST"), nil
}
func (v *RequestVariable) SuperglobalName() string { return "_REQUEST" }
