package spl

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const raiChildArraysOnly = 4

// RecursiveArrayIteratorClass 实现 PHP RecursiveArrayIterator（extends ArrayIterator�?
type RecursiveArrayIteratorClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewRecursiveArrayIteratorClass() *RecursiveArrayIteratorClass {
	return &RecursiveArrayIteratorClass{
		StaticProperty: map[string]data.Value{
			"STD_PROP_LIST":     data.NewIntValue(1),
			"ARRAY_AS_PROPS":    data.NewIntValue(2),
			"CHILD_ARRAYS_ONLY": data.NewIntValue(raiChildArraysOnly),
		},
	}
}

func (c *RecursiveArrayIteratorClass) GetName() string { return "RecursiveArrayIterator" }
func (c *RecursiveArrayIteratorClass) GetExtend() *string {
	parent := "ArrayIterator"
	return &parent
}
func (c *RecursiveArrayIteratorClass) GetImplements() []string {
	return []string{"RecursiveIterator", "Iterator", "ArrayAccess", "Countable", "SeekableIterator"}
}
func (c *RecursiveArrayIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *RecursiveArrayIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *RecursiveArrayIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *RecursiveArrayIteratorClass) GetConstruct() data.Method {
	return &ArrayIteratorConstructMethod{}
}
func (c *RecursiveArrayIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.ObjectValue.SetProperty(aiStorageKey, &data.ArrayValue{List: []*data.ZVal{}})
	cv.ObjectValue.SetProperty(aiPosKey, data.NewIntValue(0))
	cv.ObjectValue.SetProperty(aiFlagsKey, data.NewIntValue(0))
	return cv, nil
}

func (c *RecursiveArrayIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "hasChildren":
		return &RAIHasChildrenMethod{}, true
	case "getChildren":
		return &RAIGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *RecursiveArrayIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RAIHasChildrenMethod{},
		&RAIGetChildrenMethod{},
	}
}

func raiGetFlags(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(aiFlagsKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func raiCurrentValue(cv *data.ClassValue) data.Value {
	arr := aiGetStorage(cv)
	pos := aiGetPos(cv)
	if pos < 0 || pos >= len(arr.List) {
		return data.NewNullValue()
	}
	if arr.List[pos] == nil {
		return data.NewNullValue()
	}
	return arr.List[pos].Value
}

func raiValueHasChildren(val data.Value, flags int) bool {
	if val == nil {
		return false
	}
	arraysOnly := flags&raiChildArraysOnly != 0
	switch v := val.(type) {
	case *data.ArrayValue:
		return len(v.List) > 0
	case *data.ObjectValue:
		if arraysOnly {
			return false
		}
		hasAny := false
		v.RangeProperties(func(_ string, _ data.Value) bool {
			hasAny = true
			return false
		})
		return hasAny
	case *data.ClassValue:
		if arraysOnly {
			return false
		}
		for _, iface := range v.Class.GetImplements() {
			if iface == "Iterator" || iface == "IteratorAggregate" || iface == "RecursiveIterator" {
				return true
			}
		}
	}
	return false
}

func raiValueToStorage(val data.Value) *data.ArrayValue {
	switch v := val.(type) {
	case *data.ArrayValue:
		return data.CloneArrayValue(v)
	case *data.ObjectValue:
		return aoObjectToArrayValue(v)
	case *data.ClassValue:
		iterMethod, ok := v.Class.GetMethod("getIterator")
		if !ok {
			return &data.ArrayValue{List: []*data.ZVal{}}
		}
		fnCtx := v.CreateContext(iterMethod.GetVariables())
		ret, ctl := iterMethod.Call(fnCtx)
		if ctl != nil {
			return &data.ArrayValue{List: []*data.ZVal{}}
		}
		if iterCV, ok := ret.(*data.ClassValue); ok {
			if curMethod, ok := iterCV.Class.GetMethod("current"); ok {
				out := &data.ArrayValue{List: []*data.ZVal{}}
				rewindM, _ := iterCV.Class.GetMethod("rewind")
				if rewindM != nil {
					rctx := iterCV.CreateContext(rewindM.GetVariables())
					rewindM.Call(rctx)
				}
				for {
					validM, _ := iterCV.Class.GetMethod("valid")
					vctx := iterCV.CreateContext(validM.GetVariables())
					validRet, _ := validM.Call(vctx)
					if bv, ok := validRet.(*data.BoolValue); !ok || !bv.Value {
						break
					}
					cctx := iterCV.CreateContext(curMethod.GetVariables())
					curRet, _ := curMethod.Call(cctx)
					if curVal, ok := curRet.(data.Value); ok {
						out.List = append(out.List, data.NewZVal(curVal))
					}
					nextM, _ := iterCV.Class.GetMethod("next")
					nctx := iterCV.CreateContext(nextM.GetVariables())
					nextM.Call(nctx)
				}
				return out
			}
		}
	}
	return &data.ArrayValue{List: []*data.ZVal{}}
}

func raiNewInstance(ctx data.Context, storage *data.ArrayValue, flags data.Value) (data.GetValue, data.Control) {
	stmt, ok := ctx.GetVM().GetClass("RecursiveArrayIterator")
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("RecursiveArrayIterator class not found"))
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("RecursiveArrayIterator invalid"))
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		cv.ObjectValue.SetProperty(aiStorageKey, storage)
		cv.ObjectValue.SetProperty(aiFlagsKey, flags)
		return cv, nil
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		fnCtx.SetVariableValue(method.GetVariables()[0], storage)
	}
	if len(method.GetVariables()) > 1 && flags != nil {
		fnCtx.SetVariableValue(method.GetVariables()[1], flags)
	}
	_, acl = method.Call(fnCtx)
	if acl != nil {
		return nil, acl
	}
	return cv, nil
}

type RAIHasChildrenMethod struct{}

func (m *RAIHasChildrenMethod) GetName() string               { return "hasChildren" }
func (m *RAIHasChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RAIHasChildrenMethod) GetIsStatic() bool             { return false }
func (m *RAIHasChildrenMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *RAIHasChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *RAIHasChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *RAIHasChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(raiValueHasChildren(raiCurrentValue(cv), raiGetFlags(cv))), nil
}

type RAIGetChildrenMethod struct{}

func (m *RAIGetChildrenMethod) GetName() string               { return "getChildren" }
func (m *RAIGetChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RAIGetChildrenMethod) GetIsStatic() bool             { return false }
func (m *RAIGetChildrenMethod) GetReturnType() data.Types     { return nil }
func (m *RAIGetChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *RAIGetChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *RAIGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	cur := raiCurrentValue(cv)
	if !raiValueHasChildren(cur, raiGetFlags(cv)) {
		return data.NewNullValue(), nil
	}
	flags, _ := cv.ObjectValue.GetProperty(aiFlagsKey)
	return raiNewInstance(ctx, raiValueToStorage(cur), flags)
}
