package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIteratorClass 实现 PHP 的 ArrayIterator 类
// 将数组包装为 Iterator 对象
type ArrayIteratorClass struct {
	node.Node
}

func (c *ArrayIteratorClass) GetName() string    { return "ArrayIterator" }
func (c *ArrayIteratorClass) GetExtend() *string { return nil }
func (c *ArrayIteratorClass) GetImplements() []string {
	return []string{"Iterator", "ArrayAccess", "Countable"}
}
func (c *ArrayIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *ArrayIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *ArrayIteratorClass) GetConstruct() data.Method {
	return &ArrayIteratorConstructMethod{}
}
func (c *ArrayIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ArrayIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ArrayIteratorConstructMethod{}, true
	case "rewind":
		return &ArrayIteratorRewindMethod{}, true
	case "valid":
		return &ArrayIteratorValidMethod{}, true
	case "current":
		return &ArrayIteratorCurrentMethod{}, true
	case "key":
		return &ArrayIteratorKeyMethod{}, true
	case "next":
		return &ArrayIteratorNextMethod{}, true
	case "offsetExists":
		return &ArrayIteratorOffsetExistsMethod{}, true
	case "offsetGet":
		return &ArrayIteratorOffsetGetMethod{}, true
	case "offsetSet":
		return &ArrayIteratorOffsetSetMethod{}, true
	case "offsetUnset":
		return &ArrayIteratorOffsetUnsetMethod{}, true
	case "count":
		return &ArrayIteratorCountMethod{}, true
	case "append":
		return &ArrayIteratorAppendMethod{}, true
	case "getArrayCopy":
		return &ArrayIteratorGetArrayCopyMethod{}, true
	}
	return nil, false
}

func (c *ArrayIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&ArrayIteratorConstructMethod{},
		&ArrayIteratorRewindMethod{},
		&ArrayIteratorValidMethod{},
		&ArrayIteratorCurrentMethod{},
		&ArrayIteratorKeyMethod{},
		&ArrayIteratorNextMethod{},
		&ArrayIteratorOffsetExistsMethod{},
		&ArrayIteratorOffsetGetMethod{},
		&ArrayIteratorOffsetSetMethod{},
		&ArrayIteratorOffsetUnsetMethod{},
		&ArrayIteratorCountMethod{},
		&ArrayIteratorAppendMethod{},
		&ArrayIteratorGetArrayCopyMethod{},
	}
}

// 辅助：从 ArrayIterator 实例上下文获取数组和游标
func getArrayIteratorData(ctx data.Context) (*data.ArrayValue, int) {
	if cv, ok := ctx.(*data.ClassValue); ok {
		arr, _ := cv.GetProperty("__array")
		pos, _ := cv.GetProperty("__pos")
		arrVal, ok1 := arr.(*data.ArrayValue)
		posVal, ok2 := pos.(*data.IntValue)
		if ok1 && ok2 {
			return arrVal, posVal.Value
		}
		if ok1 {
			return arrVal, 0
		}
	}
	return data.NewArrayValue(nil).(*data.ArrayValue), 0
}

func setArrayIteratorPos(ctx data.Context, pos int) {
	if cv, ok := ctx.(*data.ClassValue); ok {
		cv.SetProperty("__pos", data.NewIntValue(pos))
	}
}

// __construct(array $array = [], int $flags = 0)
type ArrayIteratorConstructMethod struct{}

func (m *ArrayIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *ArrayIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, data.NewArrayValue(nil), data.Mixed{}),
	}
}
func (m *ArrayIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}
func (m *ArrayIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, _ := ctx.GetIndexValue(0)
	var arr *data.ArrayValue
	switch v := arrVal.(type) {
	case *data.ArrayValue:
		arr = v
	default:
		arr = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		cv.SetProperty("__array", arr)
		cv.SetProperty("__pos", data.NewIntValue(0))
	}
	return nil, nil
}

// rewind()
type ArrayIteratorRewindMethod struct{}

func (m *ArrayIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *ArrayIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorRewindMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayIteratorRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorRewindMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	setArrayIteratorPos(ctx, 0)
	return nil, nil
}

// valid()
type ArrayIteratorValidMethod struct{}

func (m *ArrayIteratorValidMethod) GetName() string               { return "valid" }
func (m *ArrayIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *ArrayIteratorValidMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorValidMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arr, pos := getArrayIteratorData(ctx)
	return data.NewBoolValue(pos >= 0 && pos < len(arr.List)), nil
}

// current()
type ArrayIteratorCurrentMethod struct{}

func (m *ArrayIteratorCurrentMethod) GetName() string               { return "current" }
func (m *ArrayIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arr, pos := getArrayIteratorData(ctx)
	if pos < 0 || pos >= len(arr.List) {
		return data.NewBoolValue(false), nil
	}
	return arr.List[pos].Value, nil
}

// key()
type ArrayIteratorKeyMethod struct{}

func (m *ArrayIteratorKeyMethod) GetName() string               { return "key" }
func (m *ArrayIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorKeyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, pos := getArrayIteratorData(ctx)
	return data.NewIntValue(pos), nil
}

// next()
type ArrayIteratorNextMethod struct{}

func (m *ArrayIteratorNextMethod) GetName() string               { return "next" }
func (m *ArrayIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorNextMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayIteratorNextMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorNextMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, pos := getArrayIteratorData(ctx)
	setArrayIteratorPos(ctx, pos+1)
	return nil, nil
}

// offsetExists($offset)
type ArrayIteratorOffsetExistsMethod struct{}

func (m *ArrayIteratorOffsetExistsMethod) GetName() string            { return "offsetExists" }
func (m *ArrayIteratorOffsetExistsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetExistsMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetExistsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ArrayIteratorOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayIteratorOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	offsetVal, _ := ctx.GetIndexValue(0)
	arr, _ := getArrayIteratorData(ctx)
	if intVal, ok := offsetVal.(*data.IntValue); ok {
		idx := int(intVal.Value)
		return data.NewBoolValue(idx >= 0 && idx < len(arr.List)), nil
	}
	return data.NewBoolValue(false), nil
}

// offsetGet($offset)
type ArrayIteratorOffsetGetMethod struct{}

func (m *ArrayIteratorOffsetGetMethod) GetName() string            { return "offsetGet" }
func (m *ArrayIteratorOffsetGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetGetMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetGetMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *ArrayIteratorOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayIteratorOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	offsetVal, _ := ctx.GetIndexValue(0)
	arr, _ := getArrayIteratorData(ctx)
	if intVal, ok := offsetVal.(*data.IntValue); ok {
		idx := int(intVal.Value)
		if idx >= 0 && idx < len(arr.List) {
			return arr.List[idx].Value, nil
		}
	}
	return data.NewNullValue(), nil
}

// offsetSet($offset, $value)
type ArrayIteratorOffsetSetMethod struct{}

func (m *ArrayIteratorOffsetSetMethod) GetName() string            { return "offsetSet" }
func (m *ArrayIteratorOffsetSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetSetMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetSetMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "offset", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "value", 1, nil, data.Mixed{}),
	}
}
func (m *ArrayIteratorOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "offset", 0, data.Mixed{}),
		node.NewVariable(nil, "value", 1, data.Mixed{}),
	}
}
func (m *ArrayIteratorOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// offsetUnset($offset)
type ArrayIteratorOffsetUnsetMethod struct{}

func (m *ArrayIteratorOffsetUnsetMethod) GetName() string            { return "offsetUnset" }
func (m *ArrayIteratorOffsetUnsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetUnsetMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetUnsetMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayIteratorOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// count()
type ArrayIteratorCountMethod struct{}

func (m *ArrayIteratorCountMethod) GetName() string               { return "count" }
func (m *ArrayIteratorCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorCountMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *ArrayIteratorCountMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorCountMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := getArrayIteratorData(ctx)
	return data.NewIntValue(len(arr.List)), nil
}

// append($value)
type ArrayIteratorAppendMethod struct{}

func (m *ArrayIteratorAppendMethod) GetName() string            { return "append" }
func (m *ArrayIteratorAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorAppendMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorAppendMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *ArrayIteratorAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if cv, ok := ctx.(*data.ClassValue); ok {
		arr, _ := cv.GetProperty("__array")
		if arrVal, ok := arr.(*data.ArrayValue); ok {
			if val != nil {
				arrVal.List = append(arrVal.List, data.NewZVal(val))
			}
		}
	}
	return nil, nil
}

// getArrayCopy()
type ArrayIteratorGetArrayCopyMethod struct{}

func (m *ArrayIteratorGetArrayCopyMethod) GetName() string               { return "getArrayCopy" }
func (m *ArrayIteratorGetArrayCopyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorGetArrayCopyMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorGetArrayCopyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorGetArrayCopyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorGetArrayCopyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorGetArrayCopyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := getArrayIteratorData(ctx)
	return data.CloneArrayValue(arr), nil
}
