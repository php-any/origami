package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

// EnumCaseInitializer constructs the canonical case in the executing VM.
// A cached declaration stores expressions, never parsing-VM object handles.
type EnumCaseInitializer struct {
	*Node        `pp:"-"`
	ClassName    string
	CaseName     string
	BackingType  data.TypeRef
	BackingValue data.GetValue
}

func (e *EnumCaseInitializer) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	class, ctl := ctx.GetVM().GetOrLoadClass(e.ClassName)
	if ctl != nil {
		return nil, ctl
	}
	raw, ctl := class.GetValue(ctx.CreateBaseContext())
	if ctl != nil {
		return nil, ctl
	}
	object := raw.(*data.ClassValue)
	object.SetProperty("name", data.NewStringValue(e.CaseName))
	if e.BackingValue != nil {
		raw, ctl := e.BackingValue.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		strict := ctx.CreateContext(nil)
		strict.SetStrictTypes(true)
		value, accepted, ctl := data.PrepareDeclaredValueInContext(e.BackingType, raw.(data.Value), strict)
		if ctl != nil {
			return nil, ctl
		}
		if !accepted {
			return nil, data.NewTypeError(e.GetFrom(), fmt.Errorf("Enum case backing type mismatch"))
		}
		object.SetProperty("value", value)
	}
	return object, nil
}
