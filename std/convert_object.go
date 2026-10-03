package std

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ObjectFunction 实现 PHP 的 (object) 类型转换。
// PHP：(object)$array 得到 stdClass 实例（属性为数组键）；Blade $loop = (object)$last 依赖此语义。
type ObjectFunction struct{}

func NewObjectFunction() data.FuncStmt { return &ObjectFunction{} }

func (f *ObjectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewStdClassValue(ctx), nil
	}
	switch v := value.(type) {
	case *data.NullValue:
		return data.NewStdClassValue(ctx), nil
	case *data.ThisValue:
		return v.ClassValue, nil
	case *data.ArrayValue:
		return arrayToStdClass(ctx, v), nil
	default:
		if data.ValueKindOf(value) == data.ValueObject {
			return value, nil
		}
		object := data.NewStdClassValue(ctx)
		object.SetProperty("scalar", value)
		return object, nil
	}
}

func newStdClassInstance(ctx data.Context) data.GetValue { return data.NewStdClassValue(ctx) }

func arrayToStdClass(ctx data.Context, arr *data.ArrayValue) data.GetValue {
	object := data.NewStdClassValue(ctx)
	for slots, i := arr.View(), 0; i < slots.Len(); i++ {
		slot := slots.At(i)
		object.SetProperty(slot.PHPArrayKey(i).AsString(), slot.ReadValue())
	}
	return object
}

func (f *ObjectFunction) GetName() string { return "object" }

var objectFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *ObjectFunction) GetParams() []data.GetValue {
	return objectFunctionGetParams
}

var objectFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
}

func (f *ObjectFunction) GetVariables() []data.Variable {
	return objectFunctionGetVariables
}
