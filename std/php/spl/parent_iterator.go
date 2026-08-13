package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ParentIteratorClass 实现 PHP ParentIterator（仅迭代有子节点的父节点�?
type ParentIteratorClass struct {
	node.Node
}

func NewParentIteratorClass() *ParentIteratorClass {
	return &ParentIteratorClass{}
}

func (c *ParentIteratorClass) GetName() string { return "ParentIterator" }
func (c *ParentIteratorClass) GetExtend() *string {
	parent := "RecursiveIteratorIterator"
	return &parent
}
func (c *ParentIteratorClass) GetImplements() []string { return nil }
func (c *ParentIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *ParentIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *ParentIteratorClass) GetConstruct() data.Method        { return &PIConstructMethod{} }

func (c *ParentIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(riiValidKey, data.NewBoolValue(false))
	cv.SetProperty(riiCurValKey, data.NewNullValue())
	cv.SetProperty(riiCurKeyKey, data.NewNullValue())
	cv.SetProperty(riiModeKey, data.NewIntValue(1))
	cv.SetProperty(riiStackKey, &riiStackValue{})
	return cv, nil
}

func (c *ParentIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &PIConstructMethod{}, true
	case "next":
		return &PINextMethod{}, true
	case "rewind":
		return &PIRewindMethod{}, true
	case "valid":
		return &PIValidMethod{}, true
	case "current":
		return &PICurrentMethod{}, true
	case "key":
		return &PIKeyMethod{}, true
	}
	return nil, false
}

func (c *ParentIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&PIConstructMethod{}, &PINextMethod{}, &PIRewindMethod{},
		&PIValidMethod{}, &PICurrentMethod{}, &PIKeyMethod{},
	}
}

func piCurrentIterHasChildren(cv *data.ClassValue) bool {
	iter := riiCurrentIter(cv)
	if iter == nil {
		return false
	}
	has, _ := riiCallBool(iter, "hasChildren")
	return has
}

func piAdvanceToParent(cv *data.ClassValue) data.Control {
	for {
		if ctl := riiAdvance(cv); ctl != nil {
			return ctl
		}
		if !riiIsValid(cv) {
			return nil
		}
		if piCurrentIterHasChildren(cv) {
			return nil
		}
		// skip non-parent nodes
		if ctl := piNextInternal(cv); ctl != nil {
			return ctl
		}
	}
}

func piNextInternal(cv *data.ClassValue) data.Control {
	inner := riiGetInner(cv)
	if inner == nil {
		riiSetValid(cv, false)
		return nil
	}
	mode := riiGetMode(cv)
	stack := riiGetStack(cv)
	iter := riiCurrentIter(cv)

	if mode == 1 || mode == 0 {
		hasChildren, ctl := riiCallBool(iter, "hasChildren")
		if ctl == nil && hasChildren {
			childVal, ctl2 := riiCallValue(iter, "getChildren")
			if ctl2 == nil {
				if childClass, ok := childVal.(*data.ClassValue); ok {
					if _, ctl3 := riiCallMethod(childClass, "rewind"); ctl3 == nil {
						childValid, ctl4 := riiCallBool(childClass, "valid")
						if ctl4 == nil && childValid {
							stack.frames = append(stack.frames, riiStackEntry{iter: childClass})
							return riiAdvance(cv)
						}
					}
				}
			}
		}
	}

	if _, ctl := riiCallMethod(iter, "next"); ctl != nil {
		return ctl
	}
	valid, ctl := riiCallBool(iter, "valid")
	if ctl != nil {
		return ctl
	}
	if valid {
		return riiAdvance(cv)
	}

	for len(stack.frames) > 0 {
		stack.frames = stack.frames[:len(stack.frames)-1]
		parentIter := riiCurrentIter(cv)
		if _, ctl := riiCallMethod(parentIter, "next"); ctl != nil {
			return ctl
		}
		parentValid, ctl := riiCallBool(parentIter, "valid")
		if ctl != nil {
			return ctl
		}
		if parentValid {
			return riiAdvance(cv)
		}
	}

	riiSetValid(cv, false)
	riiSetCurVal(cv, data.NewNullValue())
	return nil
}

type PIConstructMethod struct{}

func (m *PIConstructMethod) GetName() string            { return "__construct" }
func (m *PIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *PIConstructMethod) GetIsStatic() bool          { return false }
func (m *PIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *PIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("RecursiveIterator")),
	}
}
func (m *PIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("RecursiveIterator")),
	}
}
func (m *PIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	var classVal *data.ClassValue
	switch v := it.(type) {
	case *data.ClassValue:
		classVal = v
	case *data.ThisValue:
		classVal = v.ClassValue
	default:
		return nil, nil
	}
	riiSetInner(cv, classVal)
	riiSetMode(cv, 1)
	riiSetValid(cv, false)
	stack := riiGetStack(cv)
	stack.frames = stack.frames[:0]
	return nil, nil
}

type PIRewindMethod struct{}

func (m *PIRewindMethod) GetName() string               { return "rewind" }
func (m *PIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PIRewindMethod) GetIsStatic() bool             { return false }
func (m *PIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *PIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *PIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *PIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := riiGetInner(cv)
	if inner == nil {
		riiSetValid(cv, false)
		return nil, nil
	}
	stack := riiGetStack(cv)
	stack.frames = stack.frames[:0]
	if _, ctl := riiCallMethod(inner, "rewind"); ctl != nil {
		return nil, ctl
	}
	return nil, piAdvanceToParent(cv)
}

type PINextMethod struct{}

func (m *PINextMethod) GetName() string               { return "next" }
func (m *PINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PINextMethod) GetIsStatic() bool             { return false }
func (m *PINextMethod) GetReturnType() data.Types     { return nil }
func (m *PINextMethod) GetParams() []data.GetValue    { return nil }
func (m *PINextMethod) GetVariables() []data.Variable { return nil }
func (m *PINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	if ctl := piNextInternal(cv); ctl != nil {
		return nil, ctl
	}
	if !riiIsValid(cv) {
		return nil, nil
	}
	if piCurrentIterHasChildren(cv) {
		return nil, nil
	}
	return nil, piAdvanceToParent(cv)
}

type PIValidMethod struct{}

func (m *PIValidMethod) GetName() string               { return "valid" }
func (m *PIValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PIValidMethod) GetIsStatic() bool             { return false }
func (m *PIValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *PIValidMethod) GetParams() []data.GetValue    { return nil }
func (m *PIValidMethod) GetVariables() []data.Variable { return nil }
func (m *PIValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(riiIsValid(cv)), nil
}

type PICurrentMethod struct{}

func (m *PICurrentMethod) GetName() string               { return "current" }
func (m *PICurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PICurrentMethod) GetIsStatic() bool             { return false }
func (m *PICurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *PICurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *PICurrentMethod) GetVariables() []data.Variable { return nil }
func (m *PICurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return riiGetCurVal(cv), nil
}

type PIKeyMethod struct{}

func (m *PIKeyMethod) GetName() string               { return "key" }
func (m *PIKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PIKeyMethod) GetIsStatic() bool             { return false }
func (m *PIKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *PIKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *PIKeyMethod) GetVariables() []data.Variable { return nil }
func (m *PIKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return riiGetCurKey(cv), nil
}
