package spl

import (
	"sort"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func aoGetFlags(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(aoFlagsKey)
	if v != nil {
		return v
	}
	return data.NewIntValue(0)
}

func aoGetIteratorClass(cv *data.ClassValue) string {
	v, _ := cv.ObjectValue.GetProperty(aoIterClassKey)
	if v != nil {
		s := v.AsString()
		if s != "" {
			return s
		}
	}
	return "ArrayIterator"
}

func aoCompareZValKeys(a, b *data.ZVal) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	ka := a.Name
	kb := b.Name
	if ka == "" && kb == "" {
		return 0
	}
	if ka == "" {
		return -1
	}
	if kb == "" {
		return 1
	}
	if na, ok := data.ParseIntArrayKeyName(ka); ok {
		if nb, ok2 := data.ParseIntArrayKeyName(kb); ok2 {
			if na < nb {
				return -1
			}
			if na > nb {
				return 1
			}
			return 0
		}
	}
	if ka < kb {
		return -1
	}
	if ka > kb {
		return 1
	}
	return 0
}

func aoSortByValue(arr *data.ArrayValue) {
	sort.SliceStable(arr.List, func(i, j int) bool {
		return data.Compare(arr.List[i].Value, arr.List[j].Value) < 0
	})
}

func aoSortByKey(arr *data.ArrayValue) {
	sort.SliceStable(arr.List, func(i, j int) bool {
		return aoCompareZValKeys(arr.List[i], arr.List[j]) < 0
	})
}

func aoSortByValueNatural(arr *data.ArrayValue) {
	sort.SliceStable(arr.List, func(i, j int) bool {
		return strings.Compare(arr.List[i].Value.AsString(), arr.List[j].Value.AsString()) < 0
	})
}

func aoSortByValueNaturalCase(arr *data.ArrayValue) {
	sort.SliceStable(arr.List, func(i, j int) bool {
		return strings.Compare(strings.ToLower(arr.List[i].Value.AsString()), strings.ToLower(arr.List[j].Value.AsString())) < 0
	})
}

func aoUserSort(arr *data.ArrayValue, ctx data.Context, callback data.GetValue, byKey bool) bool {
	if callback == nil {
		return false
	}
	var callbackVars []data.Variable
	switch cb := callback.(type) {
	case *data.FuncValue:
		callbackVars = cb.Value.GetVariables()
	}
	sort.SliceStable(arr.List, func(i, j int) bool {
		fnCtx := ctx.CreateContext(callbackVars)
		if byKey {
			if len(callbackVars) > 0 {
				fnCtx.SetIndexZVal(0, data.NewZVal(aiKeyAt(arr, i)))
			}
			if len(callbackVars) > 1 {
				fnCtx.SetIndexZVal(1, data.NewZVal(aiKeyAt(arr, j)))
			}
		} else {
			if len(callbackVars) > 0 {
				fnCtx.SetIndexZVal(0, data.NewZVal(arr.List[i].Value))
			}
			if len(callbackVars) > 1 {
				fnCtx.SetIndexZVal(1, data.NewZVal(arr.List[j].Value))
			}
		}
		switch cb := callback.(type) {
		case *data.FuncValue:
			ret, ctl := cb.Call(fnCtx)
			if ctl != nil {
				if rv, ok := ctl.(data.ReturnControl); ok {
					if v, ok := rv.ReturnValue().(data.Value); ok {
						if iv, ok := v.(*data.IntValue); ok {
							return iv.Value < 0
						}
					}
				}
				return false
			}
			if ret != nil {
				if v, ok := ret.(data.Value); ok {
					if iv, ok := v.(*data.IntValue); ok {
						return iv.Value < 0
					}
				}
			}
		}
		return false
	})
	return true
}

func aoNewIteratorByClass(ctx data.Context, storage data.Value, flags data.Value, className string) (data.GetValue, data.Control) {
	if className == "" {
		className = "ArrayIterator"
	}
	stmt, ok := ctx.GetVM().GetClass(className)
	if !ok {
		return aoNewArrayIterator(ctx, storage)
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return aoNewArrayIterator(ctx, storage)
	}
	method := cv.Class.GetConstruct()
	if method == nil {
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

type aoSortMethod struct{ name string }

func (m *aoSortMethod) GetName() string            { return m.name }
func (m *aoSortMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *aoSortMethod) GetIsStatic() bool          { return false }
func (m *aoSortMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *aoSortMethod) GetParams() []data.GetValue { return nil }
func (m *aoSortMethod) GetVariables() []data.Variable {
	return nil
}
func (m *aoSortMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	arr := aoGetStorage(cv)
	switch m.name {
	case "asort":
		aoSortByValue(arr)
	case "ksort":
		aoSortByKey(arr)
	case "natsort":
		aoSortByValueNatural(arr)
	case "natcasesort":
		aoSortByValueNaturalCase(arr)
	}
	return data.NewBoolValue(true), nil
}

type ArrayObjectUasortMethod struct{}

func (m *ArrayObjectUasortMethod) GetName() string            { return "uasort" }
func (m *ArrayObjectUasortMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectUasortMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectUasortMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ArrayObjectUasortMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callback", 0, nil, data.Mixed{})}
}
func (m *ArrayObjectUasortMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callback", 0, data.Mixed{})}
}
func (m *ArrayObjectUasortMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	cb, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(aoUserSort(aoGetStorage(cv), ctx, cb, false)), nil
}

type ArrayObjectUksortMethod struct{}

func (m *ArrayObjectUksortMethod) GetName() string            { return "uksort" }
func (m *ArrayObjectUksortMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectUksortMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectUksortMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ArrayObjectUksortMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callback", 0, nil, data.Mixed{})}
}
func (m *ArrayObjectUksortMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callback", 0, data.Mixed{})}
}
func (m *ArrayObjectUksortMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	cb, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(aoUserSort(aoGetStorage(cv), ctx, cb, true)), nil
}

type ArrayObjectGetFlagsMethod struct{}

func (m *ArrayObjectGetFlagsMethod) GetName() string               { return "getFlags" }
func (m *ArrayObjectGetFlagsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayObjectGetFlagsMethod) GetIsStatic() bool             { return false }
func (m *ArrayObjectGetFlagsMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *ArrayObjectGetFlagsMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayObjectGetFlagsMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayObjectGetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return aoGetFlags(cv), nil
}

type ArrayObjectSetFlagsMethod struct{}

func (m *ArrayObjectSetFlagsMethod) GetName() string            { return "setFlags" }
func (m *ArrayObjectSetFlagsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectSetFlagsMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectSetFlagsMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectSetFlagsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "flags", 0, data.NewIntValue(0), data.Int{})}
}
func (m *ArrayObjectSetFlagsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "flags", 0, data.Int{})}
}
func (m *ArrayObjectSetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	flags, _ := ctx.GetIndexValue(0)
	if flags != nil {
		cv.ObjectValue.SetProperty(aoFlagsKey, flags)
	}
	return nil, nil
}

type ArrayObjectGetIteratorClassMethod struct{}

func (m *ArrayObjectGetIteratorClassMethod) GetName() string               { return "getIteratorClass" }
func (m *ArrayObjectGetIteratorClassMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayObjectGetIteratorClassMethod) GetIsStatic() bool             { return false }
func (m *ArrayObjectGetIteratorClassMethod) GetReturnType() data.Types     { return data.String{} }
func (m *ArrayObjectGetIteratorClassMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayObjectGetIteratorClassMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayObjectGetIteratorClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewStringValue("ArrayIterator"), nil
	}
	return data.NewStringValue(aoGetIteratorClass(cv)), nil
}

type ArrayObjectSetIteratorClassMethod struct{}

func (m *ArrayObjectSetIteratorClassMethod) GetName() string            { return "setIteratorClass" }
func (m *ArrayObjectSetIteratorClassMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectSetIteratorClassMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectSetIteratorClassMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectSetIteratorClassMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "iteratorClass", 0, data.NewStringValue("ArrayIterator"), data.String{})}
}
func (m *ArrayObjectSetIteratorClassMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "iteratorClass", 0, data.String{})}
}
func (m *ArrayObjectSetIteratorClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	className, _ := ctx.GetIndexValue(0)
	if className != nil {
		cv.ObjectValue.SetProperty(aoIterClassKey, data.NewStringValue(className.AsString()))
	}
	return nil, nil
}

type ArrayObjectSerializeMethod struct{}

func (m *ArrayObjectSerializeMethod) GetName() string               { return "__serialize" }
func (m *ArrayObjectSerializeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayObjectSerializeMethod) GetIsStatic() bool             { return false }
func (m *ArrayObjectSerializeMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayObjectSerializeMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayObjectSerializeMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayObjectSerializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewArrayValue(nil), nil
	}
	return &data.ArrayValue{List: []*data.ZVal{
		data.NewNamedZVal("storage", data.CloneArrayValue(aoGetStorage(cv))),
		data.NewNamedZVal("flags", aoGetFlags(cv)),
		data.NewNamedZVal("iteratorClass", data.NewStringValue(aoGetIteratorClass(cv))),
	}}, nil
}

type ArrayObjectUnserializeMethod struct{}

func (m *ArrayObjectUnserializeMethod) GetName() string            { return "__unserialize" }
func (m *ArrayObjectUnserializeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectUnserializeMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectUnserializeMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectUnserializeMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "data", 0, data.NewArrayValue(nil), data.Mixed{})}
}
func (m *ArrayObjectUnserializeMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "data", 0, data.Mixed{})}
}
func (m *ArrayObjectUnserializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	input, _ := ctx.GetIndexValue(0)
	arr, ok := input.(*data.ArrayValue)
	if !ok {
		return nil, nil
	}
	for _, z := range arr.List {
		if z == nil {
			continue
		}
		switch z.Name {
		case "storage":
			cv.ObjectValue.SetProperty(aoStorageKey, aoStorageFromInput(z.Value))
		case "flags":
			cv.ObjectValue.SetProperty(aoFlagsKey, z.Value)
		case "iteratorClass":
			cv.ObjectValue.SetProperty(aoIterClassKey, data.NewStringValue(z.Value.AsString()))
		}
	}
	return nil, nil
}
