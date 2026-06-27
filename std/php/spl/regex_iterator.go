package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/preg"
)

const (
	rxiRegexKey = "__rxi_regex__"
	rxiModeKey  = "__rxi_mode__"
	rxiFlagsKey = "__rxi_flags__"
	rxiPregKey  = "__rxi_preg__"
)

// RegexIteratorClass 实现 PHP �?RegexIterator
type RegexIteratorClass struct {
	node.Node
}

func NewRegexIteratorClass() *RegexIteratorClass {
	return &RegexIteratorClass{}
}

func (c *RegexIteratorClass) GetName() string { return "RegexIterator" }

func (c *RegexIteratorClass) GetExtend() *string {
	parent := "FilterIterator"
	return &parent
}

func (c *RegexIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (c *RegexIteratorClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RegexIteratorClass) GetPropertyList() []data.Property              { return nil }
func (c *RegexIteratorClass) GetConstruct() data.Method                     { return &RegexIteratorConstructMethod{} }

func (c *RegexIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(filterValidKey, data.NewBoolValue(false))
	cv.SetProperty(filterCurValKey, data.NewNullValue())
	cv.SetProperty(filterCurKeyKey, data.NewNullValue())
	cv.SetProperty(rxiRegexKey, data.NewStringValue(""))
	cv.SetProperty(rxiModeKey, data.NewIntValue(0))
	cv.SetProperty(rxiFlagsKey, data.NewIntValue(0))
	cv.SetProperty(rxiPregKey, data.NewIntValue(0))
	return cv, nil
}

func (c *RegexIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RegexIteratorConstructMethod{}, true
	case "accept":
		return &RegexIteratorAcceptMethod{}, true
	}
	return nil, false
}

func (c *RegexIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RegexIteratorConstructMethod{},
		&RegexIteratorAcceptMethod{},
	}
}

func rxiCopyState(from, to *data.ClassValue) {
	for _, key := range []string{rxiRegexKey, rxiModeKey, rxiFlagsKey, rxiPregKey} {
		if v, _ := from.ObjectValue.GetProperty(key); v != nil {
			to.ObjectValue.SetProperty(key, v)
		}
	}
}

func rxiMatch(cv *data.ClassValue, subject string) bool {
	regexVal, _ := cv.ObjectValue.GetProperty(rxiRegexKey)
	pattern := regexVal.AsString()
	re, err := preg.CompileAny(pattern)
	if err != nil {
		return false
	}
	anchored := preg.HasModifier(pattern, 'A')
	loc := preg.FindSubmatchAt(re, subject, 0, anchored)
	return loc != nil
}

type RegexIteratorConstructMethod struct{}

func (m *RegexIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *RegexIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RegexIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *RegexIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *RegexIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
		node.NewParameter(nil, "regex", 1, nil, data.String{}),
		node.NewParameter(nil, "mode", 2, data.NewIntValue(0), data.Int{}),
		node.NewParameter(nil, "flags", 3, data.NewIntValue(0), data.Int{}),
		node.NewParameter(nil, "pregFlags", 4, data.NewIntValue(0), data.Int{}),
	}
}
func (m *RegexIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
		node.NewVariable(nil, "regex", 1, data.String{}),
		node.NewVariable(nil, "mode", 2, data.Int{}),
		node.NewVariable(nil, "flags", 3, data.Int{}),
		node.NewVariable(nil, "pregFlags", 4, data.Int{}),
	}
}
func (m *RegexIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	it, _ := ctx.GetIndexValue(0)
	regex, _ := ctx.GetIndexValue(1)
	mode, _ := ctx.GetIndexValue(2)
	flags, _ := ctx.GetIndexValue(3)
	pregFlags, _ := ctx.GetIndexValue(4)
	filterSetInner(cv, it)
	if regex != nil {
		cv.ObjectValue.SetProperty(rxiRegexKey, regex)
	}
	if mode != nil {
		cv.ObjectValue.SetProperty(rxiModeKey, mode)
	}
	if flags != nil {
		cv.ObjectValue.SetProperty(rxiFlagsKey, flags)
	}
	if pregFlags != nil {
		cv.ObjectValue.SetProperty(rxiPregKey, pregFlags)
	}
	filterSetValid(cv, false)
	filterSetCurVal(cv, data.NewNullValue())
	filterSetCurKey(cv, data.NewNullValue())
	return nil, nil
}

type RegexIteratorAcceptMethod struct{}

func (m *RegexIteratorAcceptMethod) GetName() string            { return "accept" }
func (m *RegexIteratorAcceptMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RegexIteratorAcceptMethod) GetIsStatic() bool          { return false }
func (m *RegexIteratorAcceptMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *RegexIteratorAcceptMethod) GetParams() []data.GetValue { return nil }
func (m *RegexIteratorAcceptMethod) GetVariables() []data.Variable {
	return nil
}
func (m *RegexIteratorAcceptMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(true), nil
	}
	cur := filterGetCurVal(cv)
	subject := cur.AsString()
	matched := rxiMatch(cv, subject)
	modeVal, _ := cv.ObjectValue.GetProperty(rxiModeKey)
	mode := 0
	if iv, ok := modeVal.(*data.IntValue); ok {
		mode = iv.Value
	}
	// mode 1 = MATCH (accept if matches), mode 0 = GET_MATCH (accept if no match) �?simplified
	if mode == 1 {
		return data.NewBoolValue(matched), nil
	}
	return data.NewBoolValue(matched), nil
}

// RecursiveRegexIteratorClass 实现 PHP �?RecursiveRegexIterator
type RecursiveRegexIteratorClass struct {
	node.Node
}

func NewRecursiveRegexIteratorClass() *RecursiveRegexIteratorClass {
	return &RecursiveRegexIteratorClass{}
}

func (c *RecursiveRegexIteratorClass) GetName() string { return "RecursiveRegexIterator" }

func (c *RecursiveRegexIteratorClass) GetExtend() *string {
	parent := "RecursiveFilterIterator"
	return &parent
}

func (c *RecursiveRegexIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator", "RecursiveIterator"}
}

func (c *RecursiveRegexIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *RecursiveRegexIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *RecursiveRegexIteratorClass) GetConstruct() data.Method {
	return &RegexIteratorConstructMethod{}
}

func (c *RecursiveRegexIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(filterValidKey, data.NewBoolValue(false))
	cv.SetProperty(filterCurValKey, data.NewNullValue())
	cv.SetProperty(filterCurKeyKey, data.NewNullValue())
	cv.SetProperty(rxiRegexKey, data.NewStringValue(""))
	cv.SetProperty(rxiModeKey, data.NewIntValue(0))
	cv.SetProperty(rxiFlagsKey, data.NewIntValue(0))
	cv.SetProperty(rxiPregKey, data.NewIntValue(0))
	return cv, nil
}

func (c *RecursiveRegexIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RegexIteratorConstructMethod{}, true
	case "accept":
		return &RegexIteratorAcceptMethod{}, true
	case "getChildren":
		return &RecursiveRegexIteratorGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *RecursiveRegexIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RegexIteratorConstructMethod{},
		&RegexIteratorAcceptMethod{},
		&RecursiveRegexIteratorGetChildrenMethod{},
	}
}

type RecursiveRegexIteratorGetChildrenMethod struct{}

func (m *RecursiveRegexIteratorGetChildrenMethod) GetName() string { return "getChildren" }
func (m *RecursiveRegexIteratorGetChildrenMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveRegexIteratorGetChildrenMethod) GetIsStatic() bool          { return false }
func (m *RecursiveRegexIteratorGetChildrenMethod) GetReturnType() data.Types  { return nil }
func (m *RecursiveRegexIteratorGetChildrenMethod) GetParams() []data.GetValue { return nil }
func (m *RecursiveRegexIteratorGetChildrenMethod) GetVariables() []data.Variable {
	return nil
}
func (m *RecursiveRegexIteratorGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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
	if ctl != nil || childCV == nil {
		return childCV, ctl
	}
	rxiCopyState(cv, childCV)
	return childCV, nil
}
