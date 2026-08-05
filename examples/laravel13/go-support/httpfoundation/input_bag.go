package httpfoundation

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// InputBagClass 实现 Symfony\Component\HttpFoundation\InputBag。
type InputBagClass struct {
	node.Node
	source     *ParamBagData
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewInputBagClass() data.ClassStmt {
	return NewInputBagClassFrom(nil)
}

func NewInputBagClassFrom(source *ParamBagData) data.ClassStmt {
	c := &InputBagClass{
		source:     source,
		properties: []data.Property{protectedArrayProp("parameters")},
	}
	c.methods, c.methodList = inputBagMethods()
	return c
}

func (c *InputBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := c.source
	if src == nil {
		src = newParamBagData()
	} else {
		src = src.clone()
	}
	return data.NewProxyValue(NewInputBagClassFrom(src), ctx.CreateBaseContext()), nil
}

func (c *InputBagClass) GetName() string {
	return fqnInputBag
}
func (c *InputBagClass) GetExtend() *string {
	parent := fqnParameterBag
	return &parent
}
func (c *InputBagClass) GetImplements() []string          { return nil }
func (c *InputBagClass) GetSource() any                   { return c.source }
func (c *InputBagClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *InputBagClass) GetPropertyList() []data.Property { return c.properties }
func (c *InputBagClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *InputBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *InputBagClass) GetMethods() []data.Method { return c.methodList }

func inputBagMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("parameters", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("parameters", 0, nil)},
			nil, parameterBagConstruct),
		pubMethod("get",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			nil, inputBagGet),
		pubMethod("replace",
			[]data.GetValue{param("inputs", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("inputs", 0, nil)},
			nil, inputBagReplace),
		pubMethod("add",
			[]data.GetValue{param("inputs", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("inputs", 0, nil)},
			nil, inputBagAdd),
		pubMethod("set",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("value", 1, nil, nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("value", 1, nil),
			},
			nil, inputBagSet),
		pubMethod("getInt",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewIntValue(0), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("int"), inputBagGetInt),
		pubMethod("getBoolean",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("bool"), inputBagGetBoolean),
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
			nil, inputBagGetEnum),
		pubMethod("getString",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			data.NewBaseType("string"), inputBagGetString),
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
			nil, inputBagFilter),
	}
	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func inputBagGet(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	def := defaultValueParam(ctx, 1, data.NewNullValue())
	if def != nil && !isNull(def) && !isScalarOrStringable(def) {
		return nil, data.NewErrorThrow(nil, fmt.Errorf(
			`Expected a scalar value as a 2nd argument to "InputBag::get()", "%T" given.`, def,
		))
	}
	if store == nil || !store.has(key) {
		return def, nil
	}
	value, _ := store.get(key)
	if value != nil && !isNull(value) && !isScalarOrStringable(value) {
		return nil, throwNamed(
			"Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException",
			`Input value "%s" contains a non-scalar value.`, key,
		)
	}
	return value, nil
}

func inputBagReplace(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	store.replace(map[string]data.Value{})
	return inputBagAdd(ctx)
}

func inputBagAdd(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	for k, v := range m {
		if ctl := inputBagSetValue(store, k, v); ctl != nil {
			return nil, ctl
		}
	}
	syncParametersProperty(ctx, store)
	return nil, nil
}

func inputBagSet(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	val := defaultValueParam(ctx, 1, data.NewNullValue())
	if ctl := inputBagSetValue(store, key, val); ctl != nil {
		return nil, ctl
	}
	syncParametersProperty(ctx, store)
	return nil, nil
}

func inputBagSetValue(store *ParamBagData, key string, value data.Value) data.Control {
	if value != nil && !isNull(value) && !isScalarOrStringable(value) && !isArrayValue(value) {
		return data.NewErrorThrow(nil, fmt.Errorf(
			`Expected a scalar, or an array as a 2nd argument to "InputBag::set()", "%T" given.`, value,
		))
	}
	store.set(key, value)
	return nil
}

func inputBagGetInt(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewIntValue(0), nil
	}
	def := defaultValueParam(ctx, 1, data.NewIntValue(0))
	opts := assocMapToArrayValue(map[string]data.Value{
		"flags": data.NewIntValue(filterRequireScalar | filterNullOnFailure),
	})
	ret, ctl := paramBagFilter(ctx, store, key, def, filterValidateInt, opts, true)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil || isNull(ret) {
		return nil, throwNamed(
			"Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException",
			`Input value "%s" cannot be converted to "int".`, key,
		)
	}
	return ret, nil
}

func inputBagGetBoolean(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	def := defaultValueParam(ctx, 1, data.NewBoolValue(false))
	opts := assocMapToArrayValue(map[string]data.Value{
		"flags": data.NewIntValue(filterRequireScalar | filterNullOnFailure),
	})
	ret, ctl := paramBagFilter(ctx, store, key, def, filterValidateBool, opts, true)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil || isNull(ret) {
		return nil, throwNamed(
			"Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException",
			`Input value "%s" cannot be converted to "bool".`, key,
		)
	}
	return ret, nil
}

func inputBagGetEnum(ctx data.Context) (data.GetValue, data.Control) {
	ret, ctl := parameterBagGetEnum(ctx)
	if ctl != nil {
		// UnexpectedValueException -> BadRequestException
		return nil, throwNamed(
			"Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException",
			"%v", ctl,
		)
	}
	return ret, nil
}

func inputBagGetString(ctx data.Context) (data.GetValue, data.Control) {
	ret, ctl := inputBagGet(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil || isNull(ret) {
		def := defaultValueParam(ctx, 1, data.NewStringValue(""))
		return data.NewStringValue(def.AsString()), nil
	}
	if val, ok := ret.(data.Value); ok {
		return data.NewStringValue(val.AsString()), nil
	}
	return data.NewStringValue(""), nil
}

func inputBagFilter(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	def := defaultValueParam(ctx, 1, data.NewNullValue())
	filter := intParam(ctx, 2, filterDefault)
	opts := defaultValueParam(ctx, 3, data.NewArrayValue(nil))

	var value data.Value = def
	if store != nil && store.has(key) {
		value, _ = store.get(key)
	}

	flags := 0
	if opts != nil {
		if arr, ok := opts.(*data.ArrayValue); ok {
			m, _ := valueToAssocMap(arr)
			if f, ok := m["flags"]; ok {
				if iv, ok := f.(data.AsInt); ok {
					if n, err := iv.AsInt(); err == nil {
						flags = n
					}
				}
			}
		} else if iv, ok := opts.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				flags = n
			}
		}
	}

	if isArrayValue(value) && flags&(filterRequireArray|filterForceArray) == 0 {
		return nil, throwNamed(
			"Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException",
			`Input value "%s" contains an array, but "FILTER_REQUIRE_ARRAY" or "FILTER_FORCE_ARRAY" flags were not set.`, key,
		)
	}

	return paramBagFilter(ctx, store, key, def, filter, opts, true)
}
