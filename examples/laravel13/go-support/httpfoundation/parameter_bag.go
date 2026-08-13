package httpfoundation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ParameterBagClass 实现 Symfony\Component\HttpFoundation\ParameterBag。
type ParameterBagClass struct {
	node.Node
	source     *ParamBagData
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewParameterBagClass() data.ClassStmt {
	return NewParameterBagClassFrom(nil)
}

func NewParameterBagClassFrom(source *ParamBagData) data.ClassStmt {
	c := &ParameterBagClass{
		source: source,
		properties: []data.Property{
			protectedArrayProp("parameters"),
		},
	}
	c.methods, c.methodList = parameterBagMethods()
	return c
}

func (c *ParameterBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := c.source
	if src == nil {
		src = newParamBagData()
	} else {
		src = src.clone()
	}
	return data.NewProxyValue(NewParameterBagClassFrom(src), ctx.CreateBaseContext()), nil
}

func (c *ParameterBagClass) GetName() string    { return fqnParameterBag }
func (c *ParameterBagClass) GetExtend() *string { return nil }
func (c *ParameterBagClass) GetImplements() []string {
	return []string{"IteratorAggregate", "Countable"}
}
func (c *ParameterBagClass) GetSource() any                   { return c.source }
func (c *ParameterBagClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *ParameterBagClass) GetPropertyList() []data.Property { return c.properties }
func (c *ParameterBagClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *ParameterBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *ParameterBagClass) GetMethods() []data.Method { return c.methodList }

func parameterBagMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("parameters", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("parameters", 0, nil)},
			nil, parameterBagConstruct),
		pubMethod("all",
			[]data.GetValue{param("key", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("array"), parameterBagAll),
		pubMethod("keys", nil, nil, data.NewBaseType("array"), parameterBagKeys),
		pubMethod("replace",
			[]data.GetValue{param("parameters", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("parameters", 0, nil)},
			nil, parameterBagReplace),
		pubMethod("add",
			[]data.GetValue{param("parameters", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("parameters", 0, nil)},
			nil, parameterBagAdd),
		pubMethod("get",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			nil, parameterBagGet),
		pubMethod("set",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("value", 1, nil, nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("value", 1, nil),
			},
			nil, parameterBagSet),
		pubMethod("has",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("bool"), parameterBagHas),
		pubMethod("remove",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, parameterBagRemove),
		pubMethod("getAlpha",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("string"), parameterBagGetAlpha),
		pubMethod("getAlnum",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("string"), parameterBagGetAlnum),
		pubMethod("getDigits",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("string"), parameterBagGetDigits),
		pubMethod("getString",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("string"), parameterBagGetString),
		pubMethod("getInt",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewIntValue(0), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("int"), parameterBagGetInt),
		pubMethod("getBoolean",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("bool"), parameterBagGetBoolean),
		pubMethod("getEnum",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("class", 1, nil, nil),
				param("default", 2, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("class", 1, nil),
				variable("default", 2, nil),
			},
			nil, parameterBagGetEnum),
		pubMethod("filter",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewNullValue(), nil),
				param("filter", 2, data.NewIntValue(filterDefault), nil),
				param("options", 3, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
				variable("filter", 2, nil),
				variable("options", 3, nil),
			},
			nil, parameterBagFilter),
		pubMethod("getIterator", nil, nil, nil, parameterBagGetIterator),
		pubMethod("count", nil, nil, data.NewBaseType("int"), parameterBagCount),
	}
	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func parameterBagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	keys, m, err := valueToOrderedAssoc(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	store.replaceOrdered(keys, m)
	syncParametersProperty(ctx, store)
	return nil, nil
}

func syncParametersProperty(ctx data.Context, store *ParamBagData) {
	if cv := bagClassValue(ctx); cv != nil && store != nil {
		cv.ObjectValue.SetProperty("parameters", store.toArrayValue())
	}
}

func parameterBagAll(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return data.NewArrayValue(nil), nil
	}
	key, present, isStr := optionalStringParam(ctx, 0)
	if !present || !isStr {
		return store.toArrayValue(), nil
	}
	v, ok := store.get(key)
	if !ok {
		return data.NewArrayValue(nil), nil
	}
	if !isArrayValue(v) {
		return nil, throwNamed(
			"Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException",
			`Unexpected value for parameter "%s": expecting "array", got "%T".`, key, v,
		)
	}
	return v, nil
}

func parameterBagKeys(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return data.NewArrayValue(nil), nil
	}
	keys := store.keyList()
	vals := make([]data.Value, len(keys))
	for i, k := range keys {
		vals[i] = data.NewStringValue(k)
	}
	return data.NewArrayValue(vals), nil
}

func parameterBagReplace(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	keys, m, err := valueToOrderedAssoc(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	store.replaceOrdered(keys, m)
	syncParametersProperty(ctx, store)
	return nil, nil
}

func parameterBagAdd(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	store.add(m)
	syncParametersProperty(ctx, store)
	return nil, nil
}

func parameterBagGet(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	def := defaultValueParam(ctx, 1, data.NewNullValue())
	if store == nil {
		return def, nil
	}
	if v, exists := store.get(key); exists {
		return v, nil
	}
	return def, nil
}

func parameterBagSet(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	val := defaultValueParam(ctx, 1, data.NewNullValue())
	store.set(key, val)
	syncParametersProperty(ctx, store)
	return nil, nil
}

func parameterBagHas(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok || store == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(store.has(key)), nil
}

func parameterBagRemove(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	store.remove(key)
	syncParametersProperty(ctx, store)
	return nil, nil
}

func parameterBagGetAlpha(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewStringValue(""), nil
	}
	def := defaultValueParam(ctx, 1, data.NewStringValue("")).AsString()
	s, ctl := paramGetString(store, key, def)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(reAlpha.ReplaceAllString(s, "")), nil
}

func parameterBagGetAlnum(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewStringValue(""), nil
	}
	def := defaultValueParam(ctx, 1, data.NewStringValue("")).AsString()
	s, ctl := paramGetString(store, key, def)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(reAlnum.ReplaceAllString(s, "")), nil
}

func parameterBagGetDigits(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewStringValue(""), nil
	}
	def := defaultValueParam(ctx, 1, data.NewStringValue("")).AsString()
	s, ctl := paramGetString(store, key, def)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(reDigits.ReplaceAllString(s, "")), nil
}

func parameterBagGetString(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewStringValue(""), nil
	}
	def := defaultValueParam(ctx, 1, data.NewStringValue("")).AsString()
	s, ctl := paramGetString(store, key, def)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(s), nil
}

func parameterBagGetInt(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewIntValue(0), nil
	}
	def := defaultValueParam(ctx, 1, data.NewIntValue(0))
	opts := assocMapToArrayValue(map[string]data.Value{
		"flags": data.NewIntValue(filterRequireScalar | filterNullOnFailure),
	})
	ret, ctl := paramBagFilter(ctx, store, key, def, filterValidateInt, opts, false)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil || isNull(ret) {
		return nil, throwNamed("UnexpectedValueException", `Parameter value "%s" cannot be converted to "int".`, key)
	}
	return ret, nil
}

func parameterBagGetBoolean(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	def := defaultValueParam(ctx, 1, data.NewBoolValue(false))
	opts := assocMapToArrayValue(map[string]data.Value{
		"flags": data.NewIntValue(filterRequireScalar | filterNullOnFailure),
	})
	ret, ctl := paramBagFilter(ctx, store, key, def, filterValidateBool, opts, false)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil || isNull(ret) {
		return nil, throwNamed("UnexpectedValueException", `Parameter value "%s" cannot be converted to "bool".`, key)
	}
	return ret, nil
}

func isNull(v data.GetValue) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func parameterBagGetEnum(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	className, _, ok := optionalStringParam(ctx, 1)
	if !ok {
		return data.NewNullValue(), nil
	}
	def := defaultValueParam(ctx, 2, data.NewNullValue())
	if store == nil || !store.has(key) {
		return def, nil
	}
	value, _ := store.get(key)
	if value == nil || isNull(value) {
		return def, nil
	}

	cls, ctl := ctx.GetVM().GetOrLoadClass(className)
	if ctl != nil {
		return nil, throwNamed("UnexpectedValueException", `Parameter "%s" cannot be converted to enum: %v.`, key, ctl)
	}
	if gsm, ok := cls.(data.GetStaticMethod); ok {
		if from, ok := gsm.GetStaticMethod("from"); ok && from != nil {
			callCtx := ctx.CreateContext(from.GetVariables())
			if len(from.GetVariables()) > 0 {
				_ = callCtx.SetVariableValue(from.GetVariables()[0], value)
			}
			ret, c2 := from.Call(callCtx)
			if c2 != nil {
				return nil, throwNamed("UnexpectedValueException", `Parameter "%s" cannot be converted to enum.`, key)
			}
			return ret, nil
		}
	}
	return nil, throwNamed("UnexpectedValueException", `Parameter "%s" cannot be converted to enum: from() not available.`, key)
}

func parameterBagFilter(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	def := defaultValueParam(ctx, 1, data.NewNullValue())
	filter := intParam(ctx, 2, filterDefault)
	opts := defaultValueParam(ctx, 3, data.NewArrayValue(nil))
	return paramBagFilter(ctx, store, key, def, filter, opts, false)
}

func parameterBagGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	var arr data.Value = data.NewArrayValue(nil)
	if store != nil {
		arr = store.toArrayValue()
	}
	return newArrayIterator(ctx, arr)
}

func parameterBagCount(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(store.count()), nil
}
