package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	aiStorageKey = "__ai_storage__"
	aiPosKey     = "__ai_pos__"
	aiFlagsKey   = "__ai_flags__"
)

// ArrayIteratorClass 实现 PHP 的 ArrayIterator 类（状态存于 ClassValue 属性）
type ArrayIteratorClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewArrayIteratorClass() *ArrayIteratorClass {
	return &ArrayIteratorClass{
		StaticProperty: map[string]data.Value{
			"STD_PROP_LIST":  data.NewIntValue(1),
			"ARRAY_AS_PROPS": data.NewIntValue(2),
		},
	}
}

func (c *ArrayIteratorClass) GetName() string    { return "ArrayIterator" }
func (c *ArrayIteratorClass) GetExtend() *string { return nil }
func (c *ArrayIteratorClass) GetImplements() []string {
	return []string{"Iterator", "ArrayAccess", "Countable", "SeekableIterator"}
}
func (c *ArrayIteratorClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ArrayIteratorClass) GetPropertyList() []data.Property              { return nil }
func (c *ArrayIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *ArrayIteratorClass) GetConstruct() data.Method { return &ArrayIteratorConstructMethod{} }
func (c *ArrayIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.ObjectValue.SetProperty(aiStorageKey, &data.ArrayValue{List: []*data.ZVal{}})
	cv.ObjectValue.SetProperty(aiPosKey, data.NewIntValue(0))
	cv.ObjectValue.SetProperty(aiFlagsKey, data.NewIntValue(0))
	return cv, nil
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
	case "seek":
		return &ArrayIteratorSeekMethod{}, true
	case "getFlags":
		return &ArrayIteratorGetFlagsMethod{}, true
	case "setFlags":
		return &ArrayIteratorSetFlagsMethod{}, true
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
		&ArrayIteratorSeekMethod{},
		&ArrayIteratorGetFlagsMethod{},
		&ArrayIteratorSetFlagsMethod{},
		&ArrayIteratorOffsetExistsMethod{},
		&ArrayIteratorOffsetGetMethod{},
		&ArrayIteratorOffsetSetMethod{},
		&ArrayIteratorOffsetUnsetMethod{},
		&ArrayIteratorCountMethod{},
		&ArrayIteratorAppendMethod{},
		&ArrayIteratorGetArrayCopyMethod{},
	}
}

func aiGetCV(ctx data.Context) *data.ClassValue { return aoGetClassValue(ctx) }

func aiGetStorage(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.ObjectValue.GetProperty(aiStorageKey)
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr
	}
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	cv.ObjectValue.SetProperty(aiStorageKey, arr)
	return arr
}

func aiGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(aiPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func aiSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(aiPosKey, data.NewIntValue(pos))
}

func aiKeyAt(arr *data.ArrayValue, pos int) data.Value {
	if pos < 0 || pos >= len(arr.List) {
		return data.NewNullValue()
	}
	z := arr.List[pos]
	if z != nil && z.Name != "" {
		if n, ok := data.ParseIntArrayKeyName(z.Name); ok {
			return data.NewIntValue(n)
		}
		return data.NewStringValue(z.Name)
	}
	return data.NewIntValue(pos)
}

type ArrayIteratorConstructMethod struct{}

func (m *ArrayIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *ArrayIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, data.NewArrayValue(nil), data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
	}
}
func (m *ArrayIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}
func (m *ArrayIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	input, _ := ctx.GetIndexValue(0)
	cv.ObjectValue.SetProperty(aiStorageKey, aoStorageFromInput(input))
	aiSetPos(cv, 0)
	if flags, ok := ctx.GetIndexValue(1); ok && flags != nil {
		cv.ObjectValue.SetProperty(aiFlagsKey, flags)
	}
	return nil, nil
}

type ArrayIteratorRewindMethod struct{}

func (m *ArrayIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *ArrayIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorRewindMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayIteratorRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorRewindMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := aiGetCV(ctx); cv != nil {
		aiSetPos(cv, 0)
	}
	return nil, nil
}

type ArrayIteratorValidMethod struct{}

func (m *ArrayIteratorValidMethod) GetName() string               { return "valid" }
func (m *ArrayIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *ArrayIteratorValidMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorValidMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	arr := aiGetStorage(cv)
	pos := aiGetPos(cv)
	return data.NewBoolValue(pos >= 0 && pos < len(arr.List)), nil
}

type ArrayIteratorCurrentMethod struct{}

func (m *ArrayIteratorCurrentMethod) GetName() string               { return "current" }
func (m *ArrayIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	arr := aiGetStorage(cv)
	pos := aiGetPos(cv)
	if pos < 0 || pos >= len(arr.List) || arr.List[pos] == nil {
		return data.NewNullValue(), nil
	}
	return arr.List[pos].Value, nil
}

type ArrayIteratorKeyMethod struct{}

func (m *ArrayIteratorKeyMethod) GetName() string               { return "key" }
func (m *ArrayIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorKeyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return aiKeyAt(aiGetStorage(cv), aiGetPos(cv)), nil
}

type ArrayIteratorNextMethod struct{}

func (m *ArrayIteratorNextMethod) GetName() string               { return "next" }
func (m *ArrayIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorNextMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayIteratorNextMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorNextMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := aiGetCV(ctx); cv != nil {
		aiSetPos(cv, aiGetPos(cv)+1)
	}
	return nil, nil
}

type ArrayIteratorSeekMethod struct{}

func (m *ArrayIteratorSeekMethod) GetName() string            { return "seek" }
func (m *ArrayIteratorSeekMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorSeekMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorSeekMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorSeekMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "position", 0, nil, data.Int{})}
}
func (m *ArrayIteratorSeekMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "position", 0, data.Int{})}
}
func (m *ArrayIteratorSeekMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	posVal, _ := ctx.GetIndexValue(0)
	if iv, ok := posVal.(data.AsInt); ok {
		if pos, err := iv.AsInt(); err == nil {
			aiSetPos(cv, pos)
		}
	}
	return nil, nil
}

type ArrayIteratorGetFlagsMethod struct{}

func (m *ArrayIteratorGetFlagsMethod) GetName() string               { return "getFlags" }
func (m *ArrayIteratorGetFlagsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorGetFlagsMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorGetFlagsMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *ArrayIteratorGetFlagsMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorGetFlagsMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorGetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	v, _ := cv.ObjectValue.GetProperty(aiFlagsKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv, nil
	}
	return data.NewIntValue(0), nil
}

type ArrayIteratorSetFlagsMethod struct{}

func (m *ArrayIteratorSetFlagsMethod) GetName() string            { return "setFlags" }
func (m *ArrayIteratorSetFlagsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorSetFlagsMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorSetFlagsMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorSetFlagsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "flags", 0, nil, data.Int{})}
}
func (m *ArrayIteratorSetFlagsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "flags", 0, data.Int{})}
}
func (m *ArrayIteratorSetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	flags, _ := ctx.GetIndexValue(0)
	if flags != nil {
		cv.ObjectValue.SetProperty(aiFlagsKey, flags)
	}
	return nil, nil
}

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
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(aoOffsetExists(aiGetStorage(cv), offset)), nil
}

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
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	val := aoOffsetGet(aiGetStorage(cv), offset)
	if val == nil {
		return data.NewNullValue(), nil
	}
	return val, nil
}

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
	cv := aiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	value, _ := ctx.GetIndexValue(1)
	aoOffsetSet(aiGetStorage(cv), offset, value)
	return nil, nil
}

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
	cv := aiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	aoOffsetUnset(aiGetStorage(cv), offset)
	return nil, nil
}

type ArrayIteratorCountMethod struct{}

func (m *ArrayIteratorCountMethod) GetName() string               { return "count" }
func (m *ArrayIteratorCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorCountMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *ArrayIteratorCountMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorCountMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(aiGetStorage(cv).List)), nil
}

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
	cv := aiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	val, _ := ctx.GetIndexValue(0)
	if val != nil {
		arr := aiGetStorage(cv)
		arr.List = append(arr.List, data.NewZVal(val))
	}
	return nil, nil
}

type ArrayIteratorGetArrayCopyMethod struct{}

func (m *ArrayIteratorGetArrayCopyMethod) GetName() string               { return "getArrayCopy" }
func (m *ArrayIteratorGetArrayCopyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorGetArrayCopyMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorGetArrayCopyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorGetArrayCopyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorGetArrayCopyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorGetArrayCopyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewArrayValue(nil), nil
	}
	return data.CloneArrayValue(aiGetStorage(cv)), nil
}
