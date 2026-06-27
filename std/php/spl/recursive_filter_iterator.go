package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveFilterIteratorClass 实现 PHP �?RecursiveFilterIterator 抽象�?
type RecursiveFilterIteratorClass struct {
	node.Node
}

func NewRecursiveFilterIteratorClass() *RecursiveFilterIteratorClass {
	return &RecursiveFilterIteratorClass{}
}

func (c *RecursiveFilterIteratorClass) IsBuiltinAbstractClass() bool { return true }

func (c *RecursiveFilterIteratorClass) GetName() string { return "RecursiveFilterIterator" }

func (c *RecursiveFilterIteratorClass) GetExtend() *string {
	parent := "FilterIterator"
	return &parent
}

func (c *RecursiveFilterIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator", "RecursiveIterator"}
}

func (c *RecursiveFilterIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *RecursiveFilterIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *RecursiveFilterIteratorClass) GetConstruct() data.Method {
	return &FilterIteratorConstructMethod{}
}

func (c *RecursiveFilterIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(filterValidKey, data.NewBoolValue(false))
	cv.SetProperty(filterCurValKey, data.NewNullValue())
	cv.SetProperty(filterCurKeyKey, data.NewNullValue())
	return cv, nil
}

func (c *RecursiveFilterIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "hasChildren":
		return &RecursiveFilterIteratorHasChildrenMethod{}, true
	case "getChildren":
		return &RecursiveFilterIteratorGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *RecursiveFilterIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RecursiveFilterIteratorHasChildrenMethod{},
		&RecursiveFilterIteratorGetChildrenMethod{},
	}
}

func rfiInnerCV(inner data.GetValue) *data.ClassValue {
	if inner == nil {
		return nil
	}
	if tv, ok := inner.(*data.ThisValue); ok {
		inner = tv.ClassValue
	}
	if cv, ok := inner.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

type RecursiveFilterIteratorHasChildrenMethod struct{}

func (m *RecursiveFilterIteratorHasChildrenMethod) GetName() string { return "hasChildren" }
func (m *RecursiveFilterIteratorHasChildrenMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveFilterIteratorHasChildrenMethod) GetIsStatic() bool          { return false }
func (m *RecursiveFilterIteratorHasChildrenMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *RecursiveFilterIteratorHasChildrenMethod) GetParams() []data.GetValue { return nil }
func (m *RecursiveFilterIteratorHasChildrenMethod) GetVariables() []data.Variable {
	return nil
}
func (m *RecursiveFilterIteratorHasChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	inner := filterGetInner(cv)
	result, ctl := filterCallInnerMethod(inner, "hasChildren")
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := result.(*data.BoolValue); ok {
		return bv, nil
	}
	return data.NewBoolValue(false), nil
}

type RecursiveFilterIteratorGetChildrenMethod struct{}

func (m *RecursiveFilterIteratorGetChildrenMethod) GetName() string { return "getChildren" }
func (m *RecursiveFilterIteratorGetChildrenMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveFilterIteratorGetChildrenMethod) GetIsStatic() bool          { return false }
func (m *RecursiveFilterIteratorGetChildrenMethod) GetReturnType() data.Types  { return nil }
func (m *RecursiveFilterIteratorGetChildrenMethod) GetParams() []data.GetValue { return nil }
func (m *RecursiveFilterIteratorGetChildrenMethod) GetVariables() []data.Variable {
	return nil
}
func (m *RecursiveFilterIteratorGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := filterGetInner(cv)
	childInner, ctl := filterCallInnerMethod(inner, "getChildren")
	if ctl != nil {
		return nil, ctl
	}
	if childInner == nil {
		return data.NewNullValue(), nil
	}
	childCV, ctl := splInstantiateWithArgs(ctx, cv.Class, []data.Value{splAsValue(childInner)})
	if ctl != nil {
		return nil, ctl
	}
	return childCV, nil
}
