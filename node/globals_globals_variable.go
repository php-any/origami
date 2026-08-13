package node

import "github.com/php-any/origami/data"

// $GLOBALS

type GlobalsArrayVariable struct {
	*Node `pp:"-"`
}

// globalsValue 仅作无 VM 作用域时的回退（CLI/测试）；HTTP 请求必须走 SuperglobalArrayProvider。
var globalsValue *data.ObjectValue

func NewGlobalsArrayVariable(from data.From) data.Variable {
	return &GlobalsArrayVariable{Node: NewNode(from)}
}

func (v *GlobalsArrayVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return globalsArrayFromContext(ctx), nil
}

func (v *GlobalsArrayVariable) GetIndex() int       { return 0 }
func (v *GlobalsArrayVariable) GetName() string     { return "$GLOBALS" }
func (v *GlobalsArrayVariable) GetType() data.Types { return nil }
func (v *GlobalsArrayVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	// 目前 $GLOBALS 只读；如需支持写入，可在此扩展
	return data.NewErrorThrow(v.from, nil)
}
