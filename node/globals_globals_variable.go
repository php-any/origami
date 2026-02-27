package node

import "github.com/php-any/origami/data"

// $GLOBALS

type GlobalsArrayVariable struct {
	*Node `pp:"-"`
}

var globalsValue *data.ObjectValue

func NewGlobalsArrayVariable(from data.From) data.Variable {
	return &GlobalsArrayVariable{Node: NewNode(from)}
}

func (v *GlobalsArrayVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if globalsValue == nil {
		globalsValue = data.NewObjectValue()
	}
	return globalsValue, nil
}

func (v *GlobalsArrayVariable) GetIndex() int       { return 0 }
func (v *GlobalsArrayVariable) GetName() string     { return "$GLOBALS" }
func (v *GlobalsArrayVariable) GetType() data.Types { return nil }
func (v *GlobalsArrayVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	// 目前 $GLOBALS 只读；如需支持写入，可在此扩展
	return data.NewErrorThrow(v.from, nil)
}
