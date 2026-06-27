package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveTreeIteratorClass 实现 PHP RecursiveTreeIterator
type RecursiveTreeIteratorClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewRecursiveTreeIteratorClass() *RecursiveTreeIteratorClass {
	return &RecursiveTreeIteratorClass{
		StaticProperty: map[string]data.Value{
			"PREORDER":  data.NewIntValue(0),
			"POSTORDER": data.NewIntValue(1),
		},
	}
}

func (c *RecursiveTreeIteratorClass) GetName() string { return "RecursiveTreeIterator" }
func (c *RecursiveTreeIteratorClass) GetExtend() *string {
	parent := "RecursiveIteratorIterator"
	return &parent
}
func (c *RecursiveTreeIteratorClass) GetImplements() []string { return nil }
func (c *RecursiveTreeIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *RecursiveTreeIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *RecursiveTreeIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *RecursiveTreeIteratorClass) GetConstruct() data.Method { return &RTIConstructMethod{} }

func (c *RecursiveTreeIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(riiValidKey, data.NewBoolValue(false))
	cv.SetProperty(riiCurValKey, data.NewNullValue())
	cv.SetProperty(riiCurKeyKey, data.NewNullValue())
	cv.SetProperty(riiModeKey, data.NewIntValue(0)) // PREORDER
	cv.SetProperty(riiStackKey, &riiStackValue{})
	return cv, nil
}

func (c *RecursiveTreeIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RTIConstructMethod{}, true
	case "next":
		return &RTINextMethod{}, true
	case "rewind":
		return &RTIRewindMethod{}, true
	case "valid":
		return &RTIValidMethod{}, true
	case "current":
		return &RTICurrentMethod{}, true
	case "key":
		return &RTIKeyMethod{}, true
	case "callHasChildren":
		return &RTICallHasChildrenMethod{}, true
	case "callGetChildren":
		return &RTICallGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *RecursiveTreeIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RTIConstructMethod{}, &RTINextMethod{}, &RTIRewindMethod{},
		&RTIValidMethod{}, &RTICurrentMethod{}, &RTIKeyMethod{},
		&RTICallHasChildrenMethod{}, &RTICallGetChildrenMethod{},
	}
}

func rtiIsPostOrder(cv *data.ClassValue) bool {
	return riiGetMode(cv) == 1
}

func rtiNextPostOrder(cv *data.ClassValue) data.Control {
	stack := riiGetStack(cv)
	iter := riiCurrentIter(cv)

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

type RTIConstructMethod struct{}

func (m *RTIConstructMethod) GetName() string            { return "__construct" }
func (m *RTIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RTIConstructMethod) GetIsStatic() bool          { return false }
func (m *RTIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *RTIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("RecursiveIterator")),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "mode", 2, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *RTIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("RecursiveIterator")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "mode", 2, data.NewBaseType("int")),
	}
}
func (m *RTIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	modeVal, hasMode := ctx.GetIndexValue(2)
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
	mode := 0
	if hasMode {
		mode = splAsInt(modeVal)
	}
	riiSetInner(cv, classVal)
	riiSetMode(cv, mode)
	riiSetValid(cv, false)
	stack := riiGetStack(cv)
	stack.frames = stack.frames[:0]
	if _, ctl := riiCallMethod(classVal, "rewind"); ctl != nil {
		return nil, ctl
	}
	return nil, riiAdvance(cv)
}

type RTIRewindMethod struct{}

func (m *RTIRewindMethod) GetName() string               { return "rewind" }
func (m *RTIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTIRewindMethod) GetIsStatic() bool             { return false }
func (m *RTIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *RTIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *RTIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *RTIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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
	return nil, riiAdvance(cv)
}

type RTINextMethod struct{}

func (m *RTINextMethod) GetName() string               { return "next" }
func (m *RTINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTINextMethod) GetIsStatic() bool             { return false }
func (m *RTINextMethod) GetReturnType() data.Types     { return nil }
func (m *RTINextMethod) GetParams() []data.GetValue    { return nil }
func (m *RTINextMethod) GetVariables() []data.Variable { return nil }
func (m *RTINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	if rtiIsPostOrder(cv) {
		return nil, rtiNextPostOrder(cv)
	}
	// PREORDER: 复用 RecursiveIteratorIterator �?SELF_FIRST 逻辑
	inner := riiGetInner(cv)
	if inner == nil {
		riiSetValid(cv, false)
		return nil, nil
	}
	stack := riiGetStack(cv)
	iter := riiCurrentIter(cv)
	hasChildren, ctl := riiCallBool(iter, "hasChildren")
	if ctl == nil && hasChildren {
		childVal, ctl2 := riiCallValue(iter, "getChildren")
		if ctl2 == nil {
			if childClass, ok := childVal.(*data.ClassValue); ok {
				if _, ctl3 := riiCallMethod(childClass, "rewind"); ctl3 == nil {
					childValid, ctl4 := riiCallBool(childClass, "valid")
					if ctl4 == nil && childValid {
						stack.frames = append(stack.frames, riiStackEntry{iter: childClass})
						return nil, riiAdvance(cv)
					}
				}
			}
		}
	}
	if _, ctl := riiCallMethod(iter, "next"); ctl != nil {
		return nil, ctl
	}
	valid, ctl := riiCallBool(iter, "valid")
	if ctl != nil {
		return nil, ctl
	}
	if valid {
		return nil, riiAdvance(cv)
	}
	for len(stack.frames) > 0 {
		stack.frames = stack.frames[:len(stack.frames)-1]
		parentIter := riiCurrentIter(cv)
		if _, ctl := riiCallMethod(parentIter, "next"); ctl != nil {
			return nil, ctl
		}
		parentValid, ctl := riiCallBool(parentIter, "valid")
		if ctl != nil {
			return nil, ctl
		}
		if parentValid {
			return nil, riiAdvance(cv)
		}
	}
	riiSetValid(cv, false)
	riiSetCurVal(cv, data.NewNullValue())
	return nil, nil
}

type RTIValidMethod struct{}

func (m *RTIValidMethod) GetName() string               { return "valid" }
func (m *RTIValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTIValidMethod) GetIsStatic() bool             { return false }
func (m *RTIValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *RTIValidMethod) GetParams() []data.GetValue    { return nil }
func (m *RTIValidMethod) GetVariables() []data.Variable { return nil }
func (m *RTIValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(riiIsValid(cv)), nil
}

type RTICurrentMethod struct{}

func (m *RTICurrentMethod) GetName() string               { return "current" }
func (m *RTICurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTICurrentMethod) GetIsStatic() bool             { return false }
func (m *RTICurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *RTICurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *RTICurrentMethod) GetVariables() []data.Variable { return nil }
func (m *RTICurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return riiGetCurVal(cv), nil
}

type RTIKeyMethod struct{}

func (m *RTIKeyMethod) GetName() string               { return "key" }
func (m *RTIKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTIKeyMethod) GetIsStatic() bool             { return false }
func (m *RTIKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *RTIKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *RTIKeyMethod) GetVariables() []data.Variable { return nil }
func (m *RTIKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return riiGetCurKey(cv), nil
}

type RTICallHasChildrenMethod struct{}

func (m *RTICallHasChildrenMethod) GetName() string               { return "callHasChildren" }
func (m *RTICallHasChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTICallHasChildrenMethod) GetIsStatic() bool             { return false }
func (m *RTICallHasChildrenMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *RTICallHasChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *RTICallHasChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *RTICallHasChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	iter := riiCurrentIter(cv)
	if iter == nil {
		return data.NewBoolValue(false), nil
	}
	has, _ := riiCallBool(iter, "hasChildren")
	return data.NewBoolValue(has), nil
}

type RTICallGetChildrenMethod struct{}

func (m *RTICallGetChildrenMethod) GetName() string               { return "callGetChildren" }
func (m *RTICallGetChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RTICallGetChildrenMethod) GetIsStatic() bool             { return false }
func (m *RTICallGetChildrenMethod) GetReturnType() data.Types     { return nil }
func (m *RTICallGetChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *RTICallGetChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *RTICallGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	iter := riiCurrentIter(cv)
	if iter == nil {
		return data.NewNullValue(), nil
	}
	val, _ := riiCallValue(iter, "getChildren")
	return val, nil
}
