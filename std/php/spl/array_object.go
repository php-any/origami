package spl

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	aoStorageKey   = "__ao_storage__"
	aoFlagsKey     = "__ao_flags__"
	aoIterClassKey = "__ao_iter_class__"
)

// ArrayObjectClass 实现 PHP SPL �?ArrayObject
type ArrayObjectClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewArrayObjectClass() *ArrayObjectClass {
	return &ArrayObjectClass{
		StaticProperty: map[string]data.Value{
			"STD_PROP_LIST":  data.NewIntValue(1),
			"ARRAY_AS_PROPS": data.NewIntValue(2),
		},
	}
}

func (c *ArrayObjectClass) GetName() string    { return "ArrayObject" }
func (c *ArrayObjectClass) GetExtend() *string { return nil }
func (c *ArrayObjectClass) GetImplements() []string {
	return []string{"IteratorAggregate", "ArrayAccess", "Countable"}
}
func (c *ArrayObjectClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ArrayObjectClass) GetPropertyList() []data.Property              { return nil }
func (c *ArrayObjectClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *ArrayObjectClass) GetConstruct() data.Method {
	return &ArrayObjectConstructMethod{}
}
func (c *ArrayObjectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.ObjectValue.SetProperty(aoStorageKey, &data.ArrayValue{List: []*data.ZVal{}})
	cv.ObjectValue.SetProperty(aoFlagsKey, data.NewIntValue(0))
	cv.ObjectValue.SetProperty(aoIterClassKey, data.NewStringValue("ArrayIterator"))
	return cv, nil
}

func (c *ArrayObjectClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ArrayObjectConstructMethod{}, true
	case "append":
		return &ArrayObjectAppendMethod{}, true
	case "count":
		return &ArrayObjectCountMethod{}, true
	case "exchangeArray":
		return &ArrayObjectExchangeArrayMethod{}, true
	case "getArrayCopy":
		return &ArrayObjectGetArrayCopyMethod{}, true
	case "getIterator":
		return &ArrayObjectGetIteratorMethod{}, true
	case "offsetExists":
		return &ArrayObjectOffsetExistsMethod{}, true
	case "offsetGet":
		return &ArrayObjectOffsetGetMethod{}, true
	case "offsetSet":
		return &ArrayObjectOffsetSetMethod{}, true
	case "offsetUnset":
		return &ArrayObjectOffsetUnsetMethod{}, true
	case "asort", "ksort", "natsort", "natcasesort":
		return &aoSortMethod{name: name}, true
	case "uasort":
		return &ArrayObjectUasortMethod{}, true
	case "uksort":
		return &ArrayObjectUksortMethod{}, true
	case "getFlags":
		return &ArrayObjectGetFlagsMethod{}, true
	case "setFlags":
		return &ArrayObjectSetFlagsMethod{}, true
	case "getIteratorClass":
		return &ArrayObjectGetIteratorClassMethod{}, true
	case "setIteratorClass":
		return &ArrayObjectSetIteratorClassMethod{}, true
	case "__serialize":
		return &ArrayObjectSerializeMethod{}, true
	case "__unserialize":
		return &ArrayObjectUnserializeMethod{}, true
	}
	return nil, false
}

func (c *ArrayObjectClass) GetMethods() []data.Method {
	return []data.Method{
		&ArrayObjectConstructMethod{},
		&ArrayObjectAppendMethod{},
		&ArrayObjectCountMethod{},
		&ArrayObjectExchangeArrayMethod{},
		&ArrayObjectGetArrayCopyMethod{},
		&ArrayObjectGetIteratorMethod{},
		&ArrayObjectOffsetExistsMethod{},
		&ArrayObjectOffsetGetMethod{},
		&ArrayObjectOffsetSetMethod{},
		&ArrayObjectOffsetUnsetMethod{},
		&aoSortMethod{name: "asort"},
		&aoSortMethod{name: "ksort"},
		&aoSortMethod{name: "natsort"},
		&aoSortMethod{name: "natcasesort"},
		&ArrayObjectUasortMethod{},
		&ArrayObjectUksortMethod{},
		&ArrayObjectGetFlagsMethod{},
		&ArrayObjectSetFlagsMethod{},
		&ArrayObjectGetIteratorClassMethod{},
		&ArrayObjectSetIteratorClassMethod{},
		&ArrayObjectSerializeMethod{},
		&ArrayObjectUnserializeMethod{},
	}
}

func aoGetClassValue(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

func aoGetStorage(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.ObjectValue.GetProperty(aoStorageKey)
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr
	}
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	cv.ObjectValue.SetProperty(aoStorageKey, arr)
	return arr
}

func aoObjectToArrayValue(obj *data.ObjectValue) *data.ArrayValue {
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	obj.RangeProperties(func(key string, v data.Value) bool {
		arr.List = append(arr.List, data.NewNamedZVal(key, v))
		return true
	})
	return arr
}

func aoStorageFromInput(input data.Value) *data.ArrayValue {
	switch v := input.(type) {
	case *data.ArrayValue:
		if len(v.List) > 0 {
			return data.CloneArrayValue(v)
		}
		// �?ArrayValue 可能是关联数组经参数传递时的占位，忽略
		return &data.ArrayValue{List: []*data.ZVal{}}
	case *data.ObjectValue:
		return aoObjectToArrayValue(v)
	}
	return &data.ArrayValue{List: []*data.ZVal{}}
}

func aoOffsetExists(arr *data.ArrayValue, offset data.Value) bool {
	if iv, ok := offset.(data.AsInt); ok {
		i, err := iv.AsInt()
		if err == nil {
			if z, _ := arr.FindSlotByIntKey(i); z != nil {
				return true
			}
		}
	}
	if sv, ok := offset.(data.AsString); ok {
		key := sv.AsString()
		for _, z := range arr.List {
			if z != nil && z.Name == key {
				return true
			}
		}
	}
	return false
}

func aoOffsetGet(arr *data.ArrayValue, offset data.Value) data.Value {
	if iv, ok := offset.(data.AsInt); ok {
		i, err := iv.AsInt()
		if err == nil {
			if z, _ := arr.FindSlotByIntKey(i); z != nil {
				return z.Value
			}
		}
	}
	if sv, ok := offset.(data.AsString); ok {
		key := sv.AsString()
		for _, z := range arr.List {
			if z != nil && z.Name == key {
				return z.Value
			}
		}
	}
	return nil
}

func aoOffsetSet(arr *data.ArrayValue, offset, value data.Value) {
	if offset == nil {
		arr.List = append(arr.List, data.NewZVal(value))
		return
	}
	if iv, ok := offset.(data.AsInt); ok {
		i, err := iv.AsInt()
		if err == nil {
			arr.SetIntKey(i, value)
			return
		}
	}
	if sv, ok := offset.(data.AsString); ok {
		key := sv.AsString()
		for _, z := range arr.List {
			if z != nil && z.Name == key {
				z.Value = value
				return
			}
		}
		arr.List = append(arr.List, data.NewNamedZVal(key, value))
	}
}

func aoOffsetUnset(arr *data.ArrayValue, offset data.Value) {
	arr.UnsetKey(offset)
}

func aoNewArrayIterator(ctx data.Context, storage data.Value) (data.GetValue, data.Control) {
	stmt, ok := ctx.GetVM().GetClass("ArrayIterator")
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ArrayIterator class not found"))
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ArrayIterator invalid"))
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		return cv, nil
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		fnCtx.SetVariableValue(method.GetVariables()[0], storage)
	}
	_, acl = method.Call(fnCtx)
	if acl != nil {
		return nil, acl
	}
	return cv, nil
}

type ArrayObjectConstructMethod struct{}

func (m *ArrayObjectConstructMethod) GetName() string            { return "__construct" }
func (m *ArrayObjectConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectConstructMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, data.NewArrayValue(nil), data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
		node.NewParameter(nil, "iteratorClass", 2, data.NewStringValue("ArrayIterator"), data.String{}),
	}
}
func (m *ArrayObjectConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
		node.NewVariable(nil, "iteratorClass", 2, data.String{}),
	}
}
func (m *ArrayObjectConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	input, _ := ctx.GetIndexValue(0)
	cv.ObjectValue.SetProperty(aoStorageKey, aoStorageFromInput(input))
	if flags, ok := ctx.GetIndexValue(1); ok && flags != nil {
		cv.ObjectValue.SetProperty(aoFlagsKey, flags)
	}
	if iterClass, ok := ctx.GetIndexValue(2); ok && iterClass != nil {
		cv.ObjectValue.SetProperty(aoIterClassKey, data.NewStringValue(iterClass.AsString()))
	}
	return nil, nil
}

type ArrayObjectAppendMethod struct{}

func (m *ArrayObjectAppendMethod) GetName() string            { return "append" }
func (m *ArrayObjectAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectAppendMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectAppendMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *ArrayObjectAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *ArrayObjectAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	val, _ := ctx.GetIndexValue(0)
	arr := aoGetStorage(cv)
	arr.List = append(arr.List, data.NewZVal(val))
	return nil, nil
}

type ArrayObjectCountMethod struct{}

func (m *ArrayObjectCountMethod) GetName() string               { return "count" }
func (m *ArrayObjectCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayObjectCountMethod) GetIsStatic() bool             { return false }
func (m *ArrayObjectCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *ArrayObjectCountMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayObjectCountMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayObjectCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	arr := aoGetStorage(cv)
	return data.NewIntValue(len(arr.List)), nil
}

type ArrayObjectExchangeArrayMethod struct{}

func (m *ArrayObjectExchangeArrayMethod) GetName() string            { return "exchangeArray" }
func (m *ArrayObjectExchangeArrayMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectExchangeArrayMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectExchangeArrayMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *ArrayObjectExchangeArrayMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, data.NewArrayValue(nil), data.Mixed{})}
}
func (m *ArrayObjectExchangeArrayMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.Mixed{})}
}
func (m *ArrayObjectExchangeArrayMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewArrayValue(nil), nil
	}
	old := aoGetStorage(cv)
	oldCopy := data.CloneArrayValue(old)
	input, _ := ctx.GetIndexValue(0)
	cv.ObjectValue.SetProperty(aoStorageKey, aoStorageFromInput(input))
	return oldCopy, nil
}

type ArrayObjectGetArrayCopyMethod struct{}

func (m *ArrayObjectGetArrayCopyMethod) GetName() string               { return "getArrayCopy" }
func (m *ArrayObjectGetArrayCopyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayObjectGetArrayCopyMethod) GetIsStatic() bool             { return false }
func (m *ArrayObjectGetArrayCopyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayObjectGetArrayCopyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayObjectGetArrayCopyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayObjectGetArrayCopyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewArrayValue(nil), nil
	}
	return data.CloneArrayValue(aoGetStorage(cv)), nil
}

type ArrayObjectGetIteratorMethod struct{}

func (m *ArrayObjectGetIteratorMethod) GetName() string               { return "getIterator" }
func (m *ArrayObjectGetIteratorMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayObjectGetIteratorMethod) GetIsStatic() bool             { return false }
func (m *ArrayObjectGetIteratorMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayObjectGetIteratorMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayObjectGetIteratorMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayObjectGetIteratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	return aoNewIteratorByClass(ctx, aoGetStorage(cv), aoGetFlags(cv), aoGetIteratorClass(cv))
}

type ArrayObjectOffsetExistsMethod struct{}

func (m *ArrayObjectOffsetExistsMethod) GetName() string            { return "offsetExists" }
func (m *ArrayObjectOffsetExistsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectOffsetExistsMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectOffsetExistsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ArrayObjectOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayObjectOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayObjectOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(aoOffsetExists(aoGetStorage(cv), offset)), nil
}

type ArrayObjectOffsetGetMethod struct{}

func (m *ArrayObjectOffsetGetMethod) GetName() string            { return "offsetGet" }
func (m *ArrayObjectOffsetGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectOffsetGetMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectOffsetGetMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *ArrayObjectOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayObjectOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayObjectOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	val := aoOffsetGet(aoGetStorage(cv), offset)
	if val == nil {
		return data.NewNullValue(), nil
	}
	return val, nil
}

type ArrayObjectOffsetSetMethod struct{}

func (m *ArrayObjectOffsetSetMethod) GetName() string            { return "offsetSet" }
func (m *ArrayObjectOffsetSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectOffsetSetMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectOffsetSetMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "offset", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "value", 1, nil, data.Mixed{}),
	}
}
func (m *ArrayObjectOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "offset", 0, data.Mixed{}),
		node.NewVariable(nil, "value", 1, data.Mixed{}),
	}
}
func (m *ArrayObjectOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	value, _ := ctx.GetIndexValue(1)
	aoOffsetSet(aoGetStorage(cv), offset, value)
	return nil, nil
}

type ArrayObjectOffsetUnsetMethod struct{}

func (m *ArrayObjectOffsetUnsetMethod) GetName() string            { return "offsetUnset" }
func (m *ArrayObjectOffsetUnsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayObjectOffsetUnsetMethod) GetIsStatic() bool          { return false }
func (m *ArrayObjectOffsetUnsetMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayObjectOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayObjectOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayObjectOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := aoGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	aoOffsetUnset(aoGetStorage(cv), offset)
	return nil, nil
}
