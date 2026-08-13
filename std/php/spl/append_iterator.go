package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	aiIteratorsKey = "__ai_iterators__"
	aiIndexKey     = "__ai_idx__"
)

// appendIteratorsValue 存储 AppendIterator 的迭代器列表
type appendIteratorsValue struct {
	items []data.Value
}

func (s *appendIteratorsValue) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return s, nil
}
func (s *appendIteratorsValue) AsString() string                                     { return "appendIterators" }
func (s *appendIteratorsValue) Marshal(serializer data.Serializer) ([]byte, error)   { return nil, nil }
func (s *appendIteratorsValue) Unmarshal(b []byte, serializer data.Serializer) error { return nil }
func (s *appendIteratorsValue) ToGoValue(serializer data.Serializer) (any, error)    { return nil, nil }

// AppendIteratorClass 实现 PHP AppendIterator
type AppendIteratorClass struct {
	node.Node
}

func NewAppendIteratorClass() *AppendIteratorClass {
	return &AppendIteratorClass{}
}

func (c *AppendIteratorClass) GetName() string { return "AppendIterator" }
func (c *AppendIteratorClass) GetExtend() *string {
	parent := "IteratorIterator"
	return &parent
}
func (c *AppendIteratorClass) GetImplements() []string { return nil }
func (c *AppendIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *AppendIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *AppendIteratorClass) GetConstruct() data.Method        { return &AppIConstructMethod{} }

func (c *AppendIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(iiValidKey, data.NewBoolValue(false))
	cv.SetProperty(iiCurValKey, data.NewNullValue())
	cv.SetProperty(iiCurKeyKey, data.NewNullValue())
	cv.SetProperty(aiIteratorsKey, &appendIteratorsValue{})
	cv.SetProperty(aiIndexKey, data.NewIntValue(0))
	return cv, nil
}

func (c *AppendIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &AppIConstructMethod{}, true
	case "append":
		return &AppIAppendMethod{}, true
	case "rewind":
		return &AppIRewindMethod{}, true
	case "next":
		return &AppINextMethod{}, true
	case "valid":
		return &AppIValidMethod{}, true
	case "current":
		return &AppICurrentMethod{}, true
	case "key":
		return &AppIKeyMethod{}, true
	case "getInnerIterator":
		return &AppIGetInnerMethod{}, true
	case "getArrayIterator":
		return &AppIGetArrayIteratorMethod{}, true
	}
	return nil, false
}

func (c *AppendIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&AppIConstructMethod{}, &AppIAppendMethod{}, &AppIRewindMethod{},
		&AppINextMethod{}, &AppIValidMethod{}, &AppICurrentMethod{},
		&AppIKeyMethod{}, &AppIGetInnerMethod{}, &AppIGetArrayIteratorMethod{},
	}
}

func appGetList(cv *data.ClassValue) *appendIteratorsValue {
	v, _ := cv.ObjectValue.GetProperty(aiIteratorsKey)
	if av, ok := v.(*appendIteratorsValue); ok {
		return av
	}
	av := &appendIteratorsValue{}
	cv.ObjectValue.SetProperty(aiIteratorsKey, av)
	return av
}

func appGetIdx(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(aiIndexKey)
	return splAsInt(v)
}

func appSetIdx(cv *data.ClassValue, idx int) {
	cv.ObjectValue.SetProperty(aiIndexKey, data.NewIntValue(idx))
}

func appCurrentInner(cv *data.ClassValue) data.GetValue {
	list := appGetList(cv)
	idx := appGetIdx(cv)
	if idx < 0 || idx >= len(list.items) {
		return nil
	}
	return list.items[idx]
}

func appAdvanceToValid(cv *data.ClassValue) {
	list := appGetList(cv)
	for idx := appGetIdx(cv); idx < len(list.items); {
		inner := list.items[idx]
		iiCallInnerMethod(inner, "rewind")
		if iiInnerValid(inner) {
			iiSetInner(cv, inner)
			iiSyncFromInner(cv)
			appSetIdx(cv, idx)
			return
		}
		idx++
		appSetIdx(cv, idx)
	}
	iiSetValid(cv, false)
}

type AppIConstructMethod struct{}

func (m *AppIConstructMethod) GetName() string            { return "__construct" }
func (m *AppIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *AppIConstructMethod) GetIsStatic() bool          { return false }
func (m *AppIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *AppIConstructMethod) GetParams() []data.GetValue { return nil }
func (m *AppIConstructMethod) GetVariables() []data.Variable {
	return nil
}
func (m *AppIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	cv.ObjectValue.SetProperty(aiIteratorsKey, &appendIteratorsValue{})
	appSetIdx(cv, 0)
	iiSetValid(cv, false)
	return nil, nil
}

type AppIAppendMethod struct{}

func (m *AppIAppendMethod) GetName() string            { return "append" }
func (m *AppIAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *AppIAppendMethod) GetIsStatic() bool          { return false }
func (m *AppIAppendMethod) GetReturnType() data.Types  { return nil }
func (m *AppIAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *AppIAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *AppIAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	list := appGetList(cv)
	if v, ok := it.(data.Value); ok {
		list.items = append(list.items, v)
	}
	return nil, nil
}

type AppIRewindMethod struct{}

func (m *AppIRewindMethod) GetName() string               { return "rewind" }
func (m *AppIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppIRewindMethod) GetIsStatic() bool             { return false }
func (m *AppIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *AppIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *AppIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *AppIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	appSetIdx(cv, 0)
	appAdvanceToValid(cv)
	return nil, nil
}

type AppINextMethod struct{}

func (m *AppINextMethod) GetName() string               { return "next" }
func (m *AppINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppINextMethod) GetIsStatic() bool             { return false }
func (m *AppINextMethod) GetReturnType() data.Types     { return nil }
func (m *AppINextMethod) GetParams() []data.GetValue    { return nil }
func (m *AppINextMethod) GetVariables() []data.Variable { return nil }
func (m *AppINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := appCurrentInner(cv)
	if inner != nil {
		iiCallInnerMethod(inner, "next")
		if iiInnerValid(inner) {
			iiSyncFromInner(cv)
			return nil, nil
		}
	}
	appSetIdx(cv, appGetIdx(cv)+1)
	appAdvanceToValid(cv)
	return nil, nil
}

type AppIValidMethod struct{}

func (m *AppIValidMethod) GetName() string               { return "valid" }
func (m *AppIValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppIValidMethod) GetIsStatic() bool             { return false }
func (m *AppIValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *AppIValidMethod) GetParams() []data.GetValue    { return nil }
func (m *AppIValidMethod) GetVariables() []data.Variable { return nil }
func (m *AppIValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iiIsValid(cv)), nil
}

type AppICurrentMethod struct{}

func (m *AppICurrentMethod) GetName() string               { return "current" }
func (m *AppICurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppICurrentMethod) GetIsStatic() bool             { return false }
func (m *AppICurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *AppICurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *AppICurrentMethod) GetVariables() []data.Variable { return nil }
func (m *AppICurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return iiGetCurVal(cv), nil
}

type AppIKeyMethod struct{}

func (m *AppIKeyMethod) GetName() string               { return "key" }
func (m *AppIKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppIKeyMethod) GetIsStatic() bool             { return false }
func (m *AppIKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *AppIKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *AppIKeyMethod) GetVariables() []data.Variable { return nil }
func (m *AppIKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return iiGetCurKey(cv), nil
}

type AppIGetInnerMethod struct{}

func (m *AppIGetInnerMethod) GetName() string               { return "getInnerIterator" }
func (m *AppIGetInnerMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppIGetInnerMethod) GetIsStatic() bool             { return false }
func (m *AppIGetInnerMethod) GetReturnType() data.Types     { return nil }
func (m *AppIGetInnerMethod) GetParams() []data.GetValue    { return nil }
func (m *AppIGetInnerMethod) GetVariables() []data.Variable { return nil }
func (m *AppIGetInnerMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := appCurrentInner(cv)
	if inner == nil {
		return data.NewNullValue(), nil
	}
	return inner, nil
}

type AppIGetArrayIteratorMethod struct{}

func (m *AppIGetArrayIteratorMethod) GetName() string               { return "getArrayIterator" }
func (m *AppIGetArrayIteratorMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *AppIGetArrayIteratorMethod) GetIsStatic() bool             { return false }
func (m *AppIGetArrayIteratorMethod) GetReturnType() data.Types     { return nil }
func (m *AppIGetArrayIteratorMethod) GetParams() []data.GetValue    { return nil }
func (m *AppIGetArrayIteratorMethod) GetVariables() []data.Variable { return nil }
func (m *AppIGetArrayIteratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	list := appGetList(cv)
	zvals := make([]*data.ZVal, len(list.items))
	for i, item := range list.items {
		zvals[i] = data.NewZVal(item)
	}
	return &data.ArrayValue{List: zvals}, nil
}
