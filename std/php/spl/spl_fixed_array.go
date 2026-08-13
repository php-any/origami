package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	sfaStorageKey = "__sfa_storage__"
	sfaSizeKey    = "__sfa_size__"
	sfaPosKey     = "__sfa_pos__"
)

func sfaGetCV(ctx data.Context) *data.ClassValue {
	return aoGetClassValue(ctx)
}

func sfaGetStorage(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.ObjectValue.GetProperty(sfaStorageKey)
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr
	}
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	cv.ObjectValue.SetProperty(sfaStorageKey, arr)
	return arr
}

func sfaGetSize(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(sfaSizeKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func sfaGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(sfaPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func sfaSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(sfaPosKey, data.NewIntValue(pos))
}

// SplFixedArrayClass 实现 PHP SPL �?SplFixedArray
type SplFixedArrayClass struct {
	node.Node
}

func NewSplFixedArrayClass() *SplFixedArrayClass {
	return &SplFixedArrayClass{}
}

func (c *SplFixedArrayClass) GetName() string    { return "SplFixedArray" }
func (c *SplFixedArrayClass) GetExtend() *string { return nil }
func (c *SplFixedArrayClass) GetImplements() []string {
	return []string{"Iterator", "ArrayAccess", "Countable"}
}
func (c *SplFixedArrayClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplFixedArrayClass) GetPropertyList() []data.Property              { return nil }
func (c *SplFixedArrayClass) GetConstruct() data.Method {
	return &SplFixedArrayConstructMethod{}
}
func (c *SplFixedArrayClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(sfaStorageKey, &data.ArrayValue{List: []*data.ZVal{}})
	cv.SetProperty(sfaSizeKey, data.NewIntValue(0))
	cv.SetProperty(sfaPosKey, data.NewIntValue(0))
	return cv, nil
}

func (c *SplFixedArrayClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &SplFixedArrayConstructMethod{}, true
	case "count":
		return &SplFixedArrayCountMethod{}, true
	case "getSize":
		return &SplFixedArrayGetSizeMethod{}, true
	case "toArray":
		return &SplFixedArrayToArrayMethod{}, true
	case "rewind":
		return &SplFixedArrayRewindMethod{}, true
	case "valid":
		return &SplFixedArrayValidMethod{}, true
	case "current":
		return &SplFixedArrayCurrentMethod{}, true
	case "key":
		return &SplFixedArrayKeyMethod{}, true
	case "next":
		return &SplFixedArrayNextMethod{}, true
	case "offsetExists":
		return &SplFixedArrayOffsetExistsMethod{}, true
	case "offsetGet":
		return &SplFixedArrayOffsetGetMethod{}, true
	case "offsetSet":
		return &SplFixedArrayOffsetSetMethod{}, true
	case "offsetUnset":
		return &SplFixedArrayOffsetUnsetMethod{}, true
	}
	return nil, false
}

func (c *SplFixedArrayClass) GetMethods() []data.Method {
	return []data.Method{
		&SplFixedArrayConstructMethod{},
		&SplFixedArrayCountMethod{},
		&SplFixedArrayGetSizeMethod{},
		&SplFixedArrayToArrayMethod{},
		&SplFixedArrayRewindMethod{},
		&SplFixedArrayValidMethod{},
		&SplFixedArrayCurrentMethod{},
		&SplFixedArrayKeyMethod{},
		&SplFixedArrayNextMethod{},
		&SplFixedArrayOffsetExistsMethod{},
		&SplFixedArrayOffsetGetMethod{},
		&SplFixedArrayOffsetSetMethod{},
		&SplFixedArrayOffsetUnsetMethod{},
	}
}

func sfaInitSize(cv *data.ClassValue, size int) {
	if size < 0 {
		size = 0
	}
	list := make([]*data.ZVal, size)
	for i := range list {
		list[i] = data.NewZVal(data.NewNullValue())
	}
	cv.SetProperty(sfaStorageKey, &data.ArrayValue{List: list})
	cv.SetProperty(sfaSizeKey, data.NewIntValue(size))
	cv.SetProperty(sfaPosKey, data.NewIntValue(0))
}

type SplFixedArrayConstructMethod struct{}

func (m *SplFixedArrayConstructMethod) GetName() string            { return "__construct" }
func (m *SplFixedArrayConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFixedArrayConstructMethod) GetIsStatic() bool          { return false }
func (m *SplFixedArrayConstructMethod) GetReturnType() data.Types  { return nil }
func (m *SplFixedArrayConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "size", 0, data.NewIntValue(0), data.Int{}),
	}
}
func (m *SplFixedArrayConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "size", 0, data.Int{})}
}
func (m *SplFixedArrayConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	size := 0
	if sizeVal, ok := ctx.GetIndexValue(0); ok && sizeVal != nil {
		if iv, ok := sizeVal.(data.AsInt); ok {
			size, _ = iv.AsInt()
		}
	}
	sfaInitSize(cv, size)
	return nil, nil
}

type SplFixedArrayCountMethod struct{}

func (m *SplFixedArrayCountMethod) GetName() string               { return "count" }
func (m *SplFixedArrayCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayCountMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFixedArrayCountMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayCountMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(sfaGetSize(cv)), nil
}

type SplFixedArrayGetSizeMethod struct{}

func (m *SplFixedArrayGetSizeMethod) GetName() string               { return "getSize" }
func (m *SplFixedArrayGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayGetSizeMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayGetSizeMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFixedArrayGetSizeMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayGetSizeMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(sfaGetSize(cv)), nil
}

type SplFixedArrayToArrayMethod struct{}

func (m *SplFixedArrayToArrayMethod) GetName() string               { return "toArray" }
func (m *SplFixedArrayToArrayMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayToArrayMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayToArrayMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplFixedArrayToArrayMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayToArrayMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayToArrayMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewArrayValue(nil), nil
	}
	return data.CloneArrayValue(sfaGetStorage(cv)), nil
}

type SplFixedArrayRewindMethod struct{}

func (m *SplFixedArrayRewindMethod) GetName() string               { return "rewind" }
func (m *SplFixedArrayRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayRewindMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayRewindMethod) GetReturnType() data.Types     { return nil }
func (m *SplFixedArrayRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayRewindMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := sfaGetCV(ctx); cv != nil {
		sfaSetPos(cv, 0)
	}
	return nil, nil
}

type SplFixedArrayValidMethod struct{}

func (m *SplFixedArrayValidMethod) GetName() string               { return "valid" }
func (m *SplFixedArrayValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayValidMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFixedArrayValidMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayValidMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	pos := sfaGetPos(cv)
	return data.NewBoolValue(pos >= 0 && pos < sfaGetSize(cv)), nil
}

type SplFixedArrayCurrentMethod struct{}

func (m *SplFixedArrayCurrentMethod) GetName() string               { return "current" }
func (m *SplFixedArrayCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayCurrentMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplFixedArrayCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	pos := sfaGetPos(cv)
	arr := sfaGetStorage(cv)
	if pos < 0 || pos >= len(arr.List) {
		return data.NewNullValue(), nil
	}
	return arr.List[pos].Value, nil
}

type SplFixedArrayKeyMethod struct{}

func (m *SplFixedArrayKeyMethod) GetName() string               { return "key" }
func (m *SplFixedArrayKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayKeyMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayKeyMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFixedArrayKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayKeyMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(sfaGetPos(cv)), nil
}

type SplFixedArrayNextMethod struct{}

func (m *SplFixedArrayNextMethod) GetName() string               { return "next" }
func (m *SplFixedArrayNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFixedArrayNextMethod) GetIsStatic() bool             { return false }
func (m *SplFixedArrayNextMethod) GetReturnType() data.Types     { return nil }
func (m *SplFixedArrayNextMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFixedArrayNextMethod) GetVariables() []data.Variable { return nil }
func (m *SplFixedArrayNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := sfaGetCV(ctx); cv != nil {
		sfaSetPos(cv, sfaGetPos(cv)+1)
	}
	return nil, nil
}

type SplFixedArrayOffsetExistsMethod struct{}

func (m *SplFixedArrayOffsetExistsMethod) GetName() string            { return "offsetExists" }
func (m *SplFixedArrayOffsetExistsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFixedArrayOffsetExistsMethod) GetIsStatic() bool          { return false }
func (m *SplFixedArrayOffsetExistsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplFixedArrayOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, data.Mixed{})}
}
func (m *SplFixedArrayOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.Mixed{})}
}
func (m *SplFixedArrayOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	i, ok := splListOffsetIndex(offset)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(i >= 0 && i < sfaGetSize(cv)), nil
}

type SplFixedArrayOffsetGetMethod struct{}

func (m *SplFixedArrayOffsetGetMethod) GetName() string            { return "offsetGet" }
func (m *SplFixedArrayOffsetGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFixedArrayOffsetGetMethod) GetIsStatic() bool          { return false }
func (m *SplFixedArrayOffsetGetMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplFixedArrayOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, data.Mixed{})}
}
func (m *SplFixedArrayOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.Mixed{})}
}
func (m *SplFixedArrayOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	i, ok := splListOffsetIndex(offset)
	if !ok || i < 0 || i >= sfaGetSize(cv) {
		return data.NewNullValue(), nil
	}
	return sfaGetStorage(cv).List[i].Value, nil
}

type SplFixedArrayOffsetSetMethod struct{}

func (m *SplFixedArrayOffsetSetMethod) GetName() string            { return "offsetSet" }
func (m *SplFixedArrayOffsetSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFixedArrayOffsetSetMethod) GetIsStatic() bool          { return false }
func (m *SplFixedArrayOffsetSetMethod) GetReturnType() data.Types  { return nil }
func (m *SplFixedArrayOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "newval", 1, nil, data.Mixed{}),
	}
}
func (m *SplFixedArrayOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.Mixed{}),
		node.NewVariable(nil, "newval", 1, data.Mixed{}),
	}
}
func (m *SplFixedArrayOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	val, _ := ctx.GetIndexValue(1)
	i, ok := splListOffsetIndex(offset)
	if !ok || i < 0 || i >= sfaGetSize(cv) {
		return nil, nil
	}
	sfaGetStorage(cv).List[i].Value = val
	return nil, nil
}

type SplFixedArrayOffsetUnsetMethod struct{}

func (m *SplFixedArrayOffsetUnsetMethod) GetName() string            { return "offsetUnset" }
func (m *SplFixedArrayOffsetUnsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFixedArrayOffsetUnsetMethod) GetIsStatic() bool          { return false }
func (m *SplFixedArrayOffsetUnsetMethod) GetReturnType() data.Types  { return nil }
func (m *SplFixedArrayOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, data.Mixed{})}
}
func (m *SplFixedArrayOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.Mixed{})}
}
func (m *SplFixedArrayOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfaGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	i, ok := splListOffsetIndex(offset)
	if !ok || i < 0 || i >= sfaGetSize(cv) {
		return nil, nil
	}
	sfaGetStorage(cv).List[i].Value = data.NewNullValue()
	return nil, nil
}
