package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IteratorIterator 内部状态属性名（OuterIterator 委托模式�?
const (
	iiInnerKey  = "__ii_inner__"
	iiValidKey  = "__ii_valid__"
	iiCurValKey = "__ii_curval__"
	iiCurKeyKey = "__ii_curkey__"
)

// IteratorIteratorClass 实现 PHP �?IteratorIterator
type IteratorIteratorClass struct {
	node.Node
}

func NewIteratorIteratorClass() *IteratorIteratorClass {
	return &IteratorIteratorClass{}
}

func (c *IteratorIteratorClass) GetName() string { return "IteratorIterator" }

func (c *IteratorIteratorClass) GetExtend() *string { return nil }

func (c *IteratorIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (c *IteratorIteratorClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *IteratorIteratorClass) GetPropertyList() []data.Property              { return nil }
func (c *IteratorIteratorClass) GetConstruct() data.Method                     { return &IteratorIteratorConstructMethod{} }

func (c *IteratorIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	iiSetInner(cv, nil)
	return cv, nil
}

func (c *IteratorIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &IteratorIteratorConstructMethod{}, true
	case "rewind":
		return &IteratorIteratorRewindMethod{}, true
	case "valid":
		return &IteratorIteratorValidMethod{}, true
	case "current":
		return &IteratorIteratorCurrentMethod{}, true
	case "key":
		return &IteratorIteratorKeyMethod{}, true
	case "next":
		return &IteratorIteratorNextMethod{}, true
	case "getInnerIterator":
		return &IteratorIteratorGetInnerIteratorMethod{}, true
	}
	return nil, false
}

func (c *IteratorIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&IteratorIteratorConstructMethod{},
		&IteratorIteratorRewindMethod{},
		&IteratorIteratorValidMethod{},
		&IteratorIteratorCurrentMethod{},
		&IteratorIteratorKeyMethod{},
		&IteratorIteratorNextMethod{},
		&IteratorIteratorGetInnerIteratorMethod{},
	}
}

func iiGetInner(cv *data.ClassValue) data.GetValue {
	v, _ := cv.ObjectValue.GetProperty(iiInnerKey)
	if v == nil {
		return nil
	}
	if _, ok := v.(*data.NullValue); ok {
		return nil
	}
	return v
}

func iiSetInner(cv *data.ClassValue, inner data.GetValue) {
	if inner == nil {
		cv.ObjectValue.SetProperty(iiInnerKey, data.NewNullValue())
		return
	}
	if v, ok := inner.(data.Value); ok {
		cv.ObjectValue.SetProperty(iiInnerKey, v)
	}
}

func iiCallInner(inner data.GetValue, name string) (data.GetValue, data.Control) {
	return filterCallInnerMethod(inner, name)
}

func iiInnerValid(inner data.GetValue) bool {
	return filterInnerValid(inner)
}

func iiInnerCurrent(inner data.GetValue) data.Value {
	return filterInnerCurrent(inner)
}

func iiInnerKeyVal(inner data.GetValue) data.Value {
	return filterInnerKeyVal(inner)
}

func iiCallInnerMethod(inner data.GetValue, name string) (data.GetValue, data.Control) {
	return filterCallInnerMethod(inner, name)
}

func iiIsValid(cv *data.ClassValue) bool {
	v, _ := cv.ObjectValue.GetProperty(iiValidKey)
	if bv, ok := v.(*data.BoolValue); ok {
		return bv.Value
	}
	return iiInnerValid(iiGetInner(cv))
}

func iiSetValid(cv *data.ClassValue, valid bool) {
	cv.ObjectValue.SetProperty(iiValidKey, data.NewBoolValue(valid))
}

func iiGetCurVal(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(iiCurValKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func iiSetCurVal(cv *data.ClassValue, val data.Value) {
	if val == nil {
		val = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(iiCurValKey, val)
}

func iiGetCurKey(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(iiCurKeyKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func iiSetCurKey(cv *data.ClassValue, key data.Value) {
	if key == nil {
		key = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(iiCurKeyKey, key)
}

func iiSyncFromInner(cv *data.ClassValue) {
	inner := iiGetInner(cv)
	valid := iiInnerValid(inner)
	iiSetValid(cv, valid)
	if valid {
		iiSetCurVal(cv, iiInnerCurrent(inner))
		iiSetCurKey(cv, iiInnerKeyVal(inner))
	} else {
		iiSetCurVal(cv, data.NewNullValue())
		iiSetCurKey(cv, data.NewNullValue())
	}
}

// splInstantiateWithArgs 创建类实例并调用 __construct
func splInstantiateWithArgs(ctx data.Context, classStmt data.ClassStmt, args []data.Value) (*data.ClassValue, data.Control) {
	obj, acl := classStmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return nil, nil
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		return cv, nil
	}
	varies := method.GetVariables()
	fnCtx := cv.CreateContext(varies)
	for i, arg := range args {
		if i >= len(varies) {
			break
		}
		if ctl := fnCtx.SetVariableValue(varies[i], arg); ctl != nil {
			return nil, ctl
		}
	}
	if _, acl := method.Call(fnCtx); acl != nil {
		return nil, acl
	}
	return cv, nil
}

type IteratorIteratorConstructMethod struct{}

func (m *IteratorIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *IteratorIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *IteratorIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *IteratorIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *IteratorIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *IteratorIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *IteratorIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	iiSetInner(cv, it)
	return nil, nil
}

type IteratorIteratorRewindMethod struct{}

func (m *IteratorIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *IteratorIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *IteratorIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *IteratorIteratorRewindMethod) GetReturnType() data.Types     { return nil }
func (m *IteratorIteratorRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *IteratorIteratorRewindMethod) GetVariables() []data.Variable { return nil }
func (m *IteratorIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	_, ctl := iiCallInner(iiGetInner(cv), "rewind")
	return nil, ctl
}

type IteratorIteratorValidMethod struct{}

func (m *IteratorIteratorValidMethod) GetName() string               { return "valid" }
func (m *IteratorIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *IteratorIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *IteratorIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *IteratorIteratorValidMethod) GetParams() []data.GetValue    { return nil }
func (m *IteratorIteratorValidMethod) GetVariables() []data.Variable { return nil }
func (m *IteratorIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iiInnerValid(iiGetInner(cv))), nil
}

type IteratorIteratorCurrentMethod struct{}

func (m *IteratorIteratorCurrentMethod) GetName() string               { return "current" }
func (m *IteratorIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *IteratorIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *IteratorIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *IteratorIteratorCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *IteratorIteratorCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *IteratorIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return iiInnerCurrent(iiGetInner(cv)), nil
}

type IteratorIteratorKeyMethod struct{}

func (m *IteratorIteratorKeyMethod) GetName() string               { return "key" }
func (m *IteratorIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *IteratorIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *IteratorIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *IteratorIteratorKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *IteratorIteratorKeyMethod) GetVariables() []data.Variable { return nil }
func (m *IteratorIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return iiInnerKeyVal(iiGetInner(cv)), nil
}

type IteratorIteratorNextMethod struct{}

func (m *IteratorIteratorNextMethod) GetName() string               { return "next" }
func (m *IteratorIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *IteratorIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *IteratorIteratorNextMethod) GetReturnType() data.Types     { return nil }
func (m *IteratorIteratorNextMethod) GetParams() []data.GetValue    { return nil }
func (m *IteratorIteratorNextMethod) GetVariables() []data.Variable { return nil }
func (m *IteratorIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	_, ctl := iiCallInner(iiGetInner(cv), "next")
	return nil, ctl
}

type IteratorIteratorGetInnerIteratorMethod struct{}

func (m *IteratorIteratorGetInnerIteratorMethod) GetName() string { return "getInnerIterator" }
func (m *IteratorIteratorGetInnerIteratorMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *IteratorIteratorGetInnerIteratorMethod) GetIsStatic() bool             { return false }
func (m *IteratorIteratorGetInnerIteratorMethod) GetReturnType() data.Types     { return nil }
func (m *IteratorIteratorGetInnerIteratorMethod) GetParams() []data.GetValue    { return nil }
func (m *IteratorIteratorGetInnerIteratorMethod) GetVariables() []data.Variable { return nil }
func (m *IteratorIteratorGetInnerIteratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := iiGetInner(cv)
	if inner == nil {
		return data.NewNullValue(), nil
	}
	return inner, nil
}
