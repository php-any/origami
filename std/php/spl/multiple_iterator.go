package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	miIteratorsKey = "__mi_iterators__"
	miFlagsKey     = "__mi_flags__"
	miValidKey     = "__mi_valid__"
	miCurValKey    = "__mi_curval__"
	miCurKeyKey    = "__mi_curkey__"
)

// multipleIterEntry 存储 attachIterator 的迭代器�?info
type multipleIterEntry struct {
	iter data.Value
	info data.Value
}

// multipleIteratorsValue 存储 MultipleIterator 附加的迭代器
type multipleIteratorsValue struct {
	entries []multipleIterEntry
}

func (s *multipleIteratorsValue) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return s, nil
}
func (s *multipleIteratorsValue) AsString() string                                     { return "multipleIterators" }
func (s *multipleIteratorsValue) Marshal(serializer data.Serializer) ([]byte, error)   { return nil, nil }
func (s *multipleIteratorsValue) Unmarshal(b []byte, serializer data.Serializer) error { return nil }
func (s *multipleIteratorsValue) ToGoValue(serializer data.Serializer) (any, error)    { return nil, nil }

// MultipleIteratorClass 实现 PHP MultipleIterator
type MultipleIteratorClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewMultipleIteratorClass() *MultipleIteratorClass {
	return &MultipleIteratorClass{
		StaticProperty: map[string]data.Value{
			"MIT_NEED_ANY":     data.NewIntValue(0),
			"MIT_NEED_ALL":     data.NewIntValue(1),
			"MIT_KEYS_NUMERIC": data.NewIntValue(0),
			"MIT_KEYS_ASSOC":   data.NewIntValue(2),
		},
	}
}

func (c *MultipleIteratorClass) GetName() string    { return "MultipleIterator" }
func (c *MultipleIteratorClass) GetExtend() *string { return nil }
func (c *MultipleIteratorClass) GetImplements() []string {
	return []string{"Iterator"}
}
func (c *MultipleIteratorClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MultipleIteratorClass) GetPropertyList() []data.Property              { return nil }
func (c *MultipleIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *MultipleIteratorClass) GetConstruct() data.Method { return &MIConstructMethod{} }

func (c *MultipleIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(miIteratorsKey, &multipleIteratorsValue{})
	cv.SetProperty(miFlagsKey, data.NewIntValue(0))
	cv.SetProperty(miValidKey, data.NewBoolValue(false))
	cv.SetProperty(miCurValKey, data.NewNullValue())
	cv.SetProperty(miCurKeyKey, data.NewNullValue())
	return cv, nil
}

func (c *MultipleIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &MIConstructMethod{}, true
	case "attachIterator":
		return &MIAttachMethod{}, true
	case "detachIterator":
		return &MIDetachMethod{}, true
	case "countIterators":
		return &MICountMethod{}, true
	case "rewind":
		return &MIRewindMethod{}, true
	case "next":
		return &MINextMethod{}, true
	case "valid":
		return &MIValidMethod{}, true
	case "current":
		return &MICurrentMethod{}, true
	case "key":
		return &MIKeyMethod{}, true
	case "getFlags":
		return &MIGetFlagsMethod{}, true
	case "setFlags":
		return &MISetFlagsMethod{}, true
	}
	return nil, false
}

func (c *MultipleIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&MIConstructMethod{}, &MIAttachMethod{}, &MIDetachMethod{}, &MICountMethod{},
		&MIRewindMethod{}, &MINextMethod{}, &MIValidMethod{}, &MICurrentMethod{},
		&MIKeyMethod{}, &MIGetFlagsMethod{}, &MISetFlagsMethod{},
	}
}

func miGetList(cv *data.ClassValue) *multipleIteratorsValue {
	v, _ := cv.ObjectValue.GetProperty(miIteratorsKey)
	if av, ok := v.(*multipleIteratorsValue); ok {
		return av
	}
	av := &multipleIteratorsValue{}
	cv.ObjectValue.SetProperty(miIteratorsKey, av)
	return av
}

func miGetFlags(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(miFlagsKey)
	return splAsInt(v)
}

func miAllValid(list *multipleIteratorsValue) bool {
	if len(list.entries) == 0 {
		return false
	}
	for _, e := range list.entries {
		if !filterInnerValid(e.iter) {
			return false
		}
	}
	return true
}

func miAnyValid(list *multipleIteratorsValue) bool {
	for _, e := range list.entries {
		if filterInnerValid(e.iter) {
			return true
		}
	}
	return false
}

func miCheckValid(cv *data.ClassValue) bool {
	list := miGetList(cv)
	flags := miGetFlags(cv)
	if flags&1 != 0 { // MIT_NEED_ALL
		return miAllValid(list)
	}
	return miAnyValid(list)
}

func miBuildCurrent(cv *data.ClassValue) {
	list := miGetList(cv)
	curItems := make([]data.Value, 0, len(list.entries))
	keyItems := make([]data.Value, 0, len(list.entries))
	for _, e := range list.entries {
		if filterInnerValid(e.iter) {
			curItems = append(curItems, filterInnerCurrent(e.iter))
			keyItems = append(keyItems, filterInnerKeyVal(e.iter))
		}
	}
	cv.ObjectValue.SetProperty(miCurValKey, data.NewArrayValue(curItems))
	cv.ObjectValue.SetProperty(miCurKeyKey, data.NewArrayValue(keyItems))
	cv.ObjectValue.SetProperty(miValidKey, data.NewBoolValue(miCheckValid(cv)))
}

type MIConstructMethod struct{}

func (m *MIConstructMethod) GetName() string            { return "__construct" }
func (m *MIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *MIConstructMethod) GetIsStatic() bool          { return false }
func (m *MIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *MIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *MIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, data.NewBaseType("int")),
	}
}
func (m *MIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	flags, _ := ctx.GetIndexValue(0)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	cv.ObjectValue.SetProperty(miIteratorsKey, &multipleIteratorsValue{})
	cv.ObjectValue.SetProperty(miFlagsKey, data.NewIntValue(splAsInt(flags)))
	return nil, nil
}

type MIAttachMethod struct{}

func (m *MIAttachMethod) GetName() string            { return "attachIterator" }
func (m *MIAttachMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *MIAttachMethod) GetIsStatic() bool          { return false }
func (m *MIAttachMethod) GetReturnType() data.Types  { return nil }
func (m *MIAttachMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
		node.NewParameter(nil, "info", 1, data.NewNullValue(), data.Mixed{}),
	}
}
func (m *MIAttachMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
		node.NewVariable(nil, "info", 1, data.Mixed{}),
	}
}
func (m *MIAttachMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	info, _ := ctx.GetIndexValue(1)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	list := miGetList(cv)
	iterVal, ok := it.(data.Value)
	if !ok {
		return nil, nil
	}
	infoVal := data.NewNullValue()
	if v, ok := info.(data.Value); ok {
		infoVal = v
	}
	list.entries = append(list.entries, multipleIterEntry{iter: iterVal, info: infoVal})
	return nil, nil
}

type MIDetachMethod struct{}

func (m *MIDetachMethod) GetName() string            { return "detachIterator" }
func (m *MIDetachMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *MIDetachMethod) GetIsStatic() bool          { return false }
func (m *MIDetachMethod) GetReturnType() data.Types  { return nil }
func (m *MIDetachMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *MIDetachMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *MIDetachMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	list := miGetList(cv)
	for i, e := range list.entries {
		if e.iter == it {
			list.entries = append(list.entries[:i], list.entries[i+1:]...)
			break
		}
	}
	return nil, nil
}

type MICountMethod struct{}

func (m *MICountMethod) GetName() string               { return "countIterators" }
func (m *MICountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MICountMethod) GetIsStatic() bool             { return false }
func (m *MICountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *MICountMethod) GetParams() []data.GetValue    { return nil }
func (m *MICountMethod) GetVariables() []data.Variable { return nil }
func (m *MICountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(miGetList(cv).entries)), nil
}

type MIRewindMethod struct{}

func (m *MIRewindMethod) GetName() string               { return "rewind" }
func (m *MIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MIRewindMethod) GetIsStatic() bool             { return false }
func (m *MIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *MIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *MIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *MIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	for _, e := range miGetList(cv).entries {
		filterCallInnerMethod(e.iter, "rewind")
	}
	miBuildCurrent(cv)
	return nil, nil
}

type MINextMethod struct{}

func (m *MINextMethod) GetName() string               { return "next" }
func (m *MINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MINextMethod) GetIsStatic() bool             { return false }
func (m *MINextMethod) GetReturnType() data.Types     { return nil }
func (m *MINextMethod) GetParams() []data.GetValue    { return nil }
func (m *MINextMethod) GetVariables() []data.Variable { return nil }
func (m *MINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	for _, e := range miGetList(cv).entries {
		filterCallInnerMethod(e.iter, "next")
	}
	miBuildCurrent(cv)
	return nil, nil
}

type MIValidMethod struct{}

func (m *MIValidMethod) GetName() string               { return "valid" }
func (m *MIValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MIValidMethod) GetIsStatic() bool             { return false }
func (m *MIValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *MIValidMethod) GetParams() []data.GetValue    { return nil }
func (m *MIValidMethod) GetVariables() []data.Variable { return nil }
func (m *MIValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	v, _ := cv.ObjectValue.GetProperty(miValidKey)
	if bv, ok := v.(*data.BoolValue); ok {
		return bv, nil
	}
	return data.NewBoolValue(false), nil
}

type MICurrentMethod struct{}

func (m *MICurrentMethod) GetName() string               { return "current" }
func (m *MICurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MICurrentMethod) GetIsStatic() bool             { return false }
func (m *MICurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *MICurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *MICurrentMethod) GetVariables() []data.Variable { return nil }
func (m *MICurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	v, _ := cv.ObjectValue.GetProperty(miCurValKey)
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

type MIKeyMethod struct{}

func (m *MIKeyMethod) GetName() string               { return "key" }
func (m *MIKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MIKeyMethod) GetIsStatic() bool             { return false }
func (m *MIKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *MIKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *MIKeyMethod) GetVariables() []data.Variable { return nil }
func (m *MIKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	v, _ := cv.ObjectValue.GetProperty(miCurKeyKey)
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

type MIGetFlagsMethod struct{}

func (m *MIGetFlagsMethod) GetName() string               { return "getFlags" }
func (m *MIGetFlagsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *MIGetFlagsMethod) GetIsStatic() bool             { return false }
func (m *MIGetFlagsMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *MIGetFlagsMethod) GetParams() []data.GetValue    { return nil }
func (m *MIGetFlagsMethod) GetVariables() []data.Variable { return nil }
func (m *MIGetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(miGetFlags(cv)), nil
}

type MISetFlagsMethod struct{}

func (m *MISetFlagsMethod) GetName() string            { return "setFlags" }
func (m *MISetFlagsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *MISetFlagsMethod) GetIsStatic() bool          { return false }
func (m *MISetFlagsMethod) GetReturnType() data.Types  { return nil }
func (m *MISetFlagsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, nil, data.NewBaseType("int")),
	}
}
func (m *MISetFlagsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, data.NewBaseType("int")),
	}
}
func (m *MISetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	flags, _ := ctx.GetIndexValue(0)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	cv.ObjectValue.SetProperty(miFlagsKey, data.NewIntValue(splAsInt(flags)))
	return nil, nil
}
