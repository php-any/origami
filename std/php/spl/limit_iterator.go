package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	liOffsetKey = "__li_offset__"
	liCountKey  = "__li_count__"
	liPosKey    = "__li_pos__"
)

// LimitIteratorClass 实现 PHP LimitIterator
type LimitIteratorClass struct {
	node.Node
}

func NewLimitIteratorClass() *LimitIteratorClass {
	return &LimitIteratorClass{}
}

func (c *LimitIteratorClass) GetName() string { return "LimitIterator" }
func (c *LimitIteratorClass) GetExtend() *string {
	parent := "IteratorIterator"
	return &parent
}
func (c *LimitIteratorClass) GetImplements() []string { return nil }
func (c *LimitIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *LimitIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *LimitIteratorClass) GetConstruct() data.Method        { return &LIConstructMethod{} }

func (c *LimitIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(iiValidKey, data.NewBoolValue(false))
	cv.SetProperty(iiCurValKey, data.NewNullValue())
	cv.SetProperty(iiCurKeyKey, data.NewNullValue())
	cv.SetProperty(liPosKey, data.NewIntValue(0))
	return cv, nil
}

func (c *LimitIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &LIConstructMethod{}, true
	case "rewind":
		return &LIRewindMethod{}, true
	case "valid":
		return &LIValidMethod{}, true
	case "next":
		return &LINextMethod{}, true
	case "getPosition":
		return &LIGetPositionMethod{}, true
	case "seek":
		return &LISeekMethod{}, true
	}
	return nil, false
}

func (c *LimitIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&LIConstructMethod{},
		&LIRewindMethod{},
		&LIValidMethod{},
		&LINextMethod{},
		&LIGetPositionMethod{},
		&LISeekMethod{},
	}
}

func liGetOffset(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(liOffsetKey)
	return splAsInt(v)
}

func liGetCount(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(liCountKey)
	return splAsInt(v)
}

func liGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(liPosKey)
	return splAsInt(v)
}

func liSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(liPosKey, data.NewIntValue(pos))
}

func liSkipOffset(cv *data.ClassValue) {
	offset := liGetOffset(cv)
	inner := iiGetInner(cv)
	iiCallInnerMethod(inner, "rewind")
	for i := 0; i < offset; i++ {
		if !iiInnerValid(inner) {
			break
		}
		iiCallInnerMethod(inner, "next")
	}
}

func liCheckValid(cv *data.ClassValue) bool {
	count := liGetCount(cv)
	pos := liGetPos(cv)
	inner := iiGetInner(cv)
	if !iiInnerValid(inner) {
		return false
	}
	if count >= 0 && pos >= count {
		return false
	}
	return true
}

type LIConstructMethod struct{}

func (m *LIConstructMethod) GetName() string            { return "__construct" }
func (m *LIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *LIConstructMethod) GetIsStatic() bool          { return false }
func (m *LIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *LIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
		node.NewParameter(nil, "offset", 1, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "count", 2, data.NewIntValue(-1), data.NewBaseType("int")),
	}
}
func (m *LIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
		node.NewVariable(nil, "offset", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "count", 2, data.NewBaseType("int")),
	}
}
func (m *LIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	offset, _ := ctx.GetIndexValue(1)
	count, _ := ctx.GetIndexValue(2)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	iiSetInner(cv, it)
	cv.ObjectValue.SetProperty(liOffsetKey, data.NewIntValue(splAsInt(offset)))
	cv.ObjectValue.SetProperty(liCountKey, data.NewIntValue(splAsInt(count)))
	liSetPos(cv, 0)
	iiSetValid(cv, false)
	return nil, nil
}

type LIRewindMethod struct{}

func (m *LIRewindMethod) GetName() string               { return "rewind" }
func (m *LIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *LIRewindMethod) GetIsStatic() bool             { return false }
func (m *LIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *LIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *LIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *LIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	liSetPos(cv, 0)
	liSkipOffset(cv)
	if liCheckValid(cv) {
		iiSyncFromInner(cv)
	} else {
		iiSetValid(cv, false)
	}
	return nil, nil
}

type LIValidMethod struct{}

func (m *LIValidMethod) GetName() string               { return "valid" }
func (m *LIValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *LIValidMethod) GetIsStatic() bool             { return false }
func (m *LIValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *LIValidMethod) GetParams() []data.GetValue    { return nil }
func (m *LIValidMethod) GetVariables() []data.Variable { return nil }
func (m *LIValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iiIsValid(cv) && liCheckValid(cv)), nil
}

type LINextMethod struct{}

func (m *LINextMethod) GetName() string               { return "next" }
func (m *LINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *LINextMethod) GetIsStatic() bool             { return false }
func (m *LINextMethod) GetReturnType() data.Types     { return nil }
func (m *LINextMethod) GetParams() []data.GetValue    { return nil }
func (m *LINextMethod) GetVariables() []data.Variable { return nil }
func (m *LINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := iiGetInner(cv)
	iiCallInnerMethod(inner, "next")
	liSetPos(cv, liGetPos(cv)+1)
	if liCheckValid(cv) {
		iiSyncFromInner(cv)
	} else {
		iiSetValid(cv, false)
	}
	return nil, nil
}

type LIGetPositionMethod struct{}

func (m *LIGetPositionMethod) GetName() string               { return "getPosition" }
func (m *LIGetPositionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *LIGetPositionMethod) GetIsStatic() bool             { return false }
func (m *LIGetPositionMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *LIGetPositionMethod) GetParams() []data.GetValue    { return nil }
func (m *LIGetPositionMethod) GetVariables() []data.Variable { return nil }
func (m *LIGetPositionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(liGetPos(cv)), nil
}

type LISeekMethod struct{}

func (m *LISeekMethod) GetName() string            { return "seek" }
func (m *LISeekMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *LISeekMethod) GetIsStatic() bool          { return false }
func (m *LISeekMethod) GetReturnType() data.Types  { return nil }
func (m *LISeekMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "position", 0, nil, data.NewBaseType("int")),
	}
}
func (m *LISeekMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "position", 0, data.NewBaseType("int")),
	}
}
func (m *LISeekMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	posVal, _ := ctx.GetIndexValue(0)
	target := splAsInt(posVal)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	// seek 到相�?position（在 limit 窗口内）
	current := liGetPos(cv)
	if target > current {
		for i := current; i < target; i++ {
			inner := iiGetInner(cv)
			iiCallInnerMethod(inner, "next")
			liSetPos(cv, liGetPos(cv)+1)
			if !liCheckValid(cv) {
				iiSetValid(cv, false)
				return nil, nil
			}
		}
	}
	iiSyncFromInner(cv)
	return nil, nil
}
