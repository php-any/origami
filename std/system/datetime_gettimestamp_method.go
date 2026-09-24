package system

import (
	"github.com/php-any/origami/data"
)

type DateTimeGetTimestampMethod struct {
	source *DateTime
}

func (h *DateTimeGetTimestampMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(int(h.source.GetTimestamp())), nil
}

func (h *DateTimeGetTimestampMethod) GetName() string {
	return "getTimestamp"
}

func (h *DateTimeGetTimestampMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *DateTimeGetTimestampMethod) GetIsStatic() bool {
	return false
}

var dateTimeGetTimestampMethodGetParams = []data.GetValue{}

func (h *DateTimeGetTimestampMethod) GetParams() []data.GetValue {
	return dateTimeGetTimestampMethodGetParams
}

var dateTimeGetTimestampMethodGetVariables = []data.Variable{}

func (h *DateTimeGetTimestampMethod) GetVariables() []data.Variable {
	return dateTimeGetTimestampMethodGetVariables
}

func (h *DateTimeGetTimestampMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
