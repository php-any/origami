package httpfoundation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// bagMethod 是共享的方法适配器。
type bagMethod struct {
	name     string
	modifier data.Modifier
	static   bool
	params   []data.GetValue
	vars     []data.Variable
	ret      data.Types
	fn       func(ctx data.Context) (data.GetValue, data.Control)
}

func (m *bagMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.fn(ctx)
}
func (m *bagMethod) GetName() string               { return m.name }
func (m *bagMethod) GetModifier() data.Modifier    { return m.modifier }
func (m *bagMethod) GetIsStatic() bool             { return m.static }
func (m *bagMethod) GetParams() []data.GetValue    { return m.params }
func (m *bagMethod) GetVariables() []data.Variable { return m.vars }
func (m *bagMethod) GetReturnType() data.Types     { return m.ret }

func pubMethod(name string, params []data.GetValue, vars []data.Variable, ret data.Types, fn func(ctx data.Context) (data.GetValue, data.Control)) data.Method {
	return &bagMethod{
		name:     name,
		modifier: data.ModifierPublic,
		params:   params,
		vars:     vars,
		ret:      ret,
		fn:       fn,
	}
}

func param(name string, index int, def data.GetValue, ty data.Types) data.GetValue {
	return node.NewParameter(nil, name, index, def, ty)
}

func variable(name string, index int, ty data.Types) data.Variable {
	return node.NewVariable(nil, name, index, ty)
}

// PHP filter 常量（与 std/php/load.go 对齐）
const (
	filterDefault       = 516
	filterValidateInt   = 257
	filterValidateBool  = 258
	filterCallback      = 1024
	filterNullOnFailure = 134217728
	filterRequireArray  = 8
	filterForceArray    = 64 // PHP FILTER_FORCE_ARRAY
	filterRequireScalar = 33554432
)

var (
	reAlpha  = regexp.MustCompile(`[^[:alpha:]]`)
	reAlnum  = regexp.MustCompile(`[^[:alnum:]]`)
	reDigits = regexp.MustCompile(`[^[:digit:]]`)
)

func paramBagFilter(ctx data.Context, store *ParamBagData, key string, defaultVal data.Value, filter int, options data.Value, badRequest bool) (data.GetValue, data.Control) {
	value := defaultVal
	if store != nil {
		if v, ok := store.get(key); ok {
			value = v
		}
	}

	flags := 0
	var optMap map[string]data.Value
	if options != nil {
		if arr, ok := options.(*data.ArrayValue); ok {
			optMap, _ = valueToAssocMap(arr)
			if f, ok := optMap["flags"]; ok {
				if iv, ok := f.(data.AsInt); ok {
					if n, err := iv.AsInt(); err == nil {
						flags = n
					}
				}
			}
		} else if iv, ok := options.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil && n != 0 {
				flags = n
				optMap = map[string]data.Value{"flags": data.NewIntValue(n)}
			}
		}
	}
	if optMap == nil {
		optMap = map[string]data.Value{}
	}
	if isArrayValue(value) && flags == 0 {
		flags = filterRequireArray
		optMap["flags"] = data.NewIntValue(flags)
	}

	if value != nil && !isScalarOrStringable(value) && !isArrayValue(value) {
		name := "UnexpectedValueException"
		if badRequest {
			name = "Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException"
		}
		return nil, throwNamed(name, `Parameter value "%s" cannot be filtered.`, key)
	}

	if filter&filterCallback != 0 {
		cb, ok := optMap["options"]
		if !ok {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("A Closure must be passed when FILTER_CALLBACK is used"))
		}
		_ = cb
		return nil, data.NewErrorThrow(nil, fmt.Errorf("FILTER_CALLBACK is not fully supported in Go ParameterBag"))
	}

	nullOnFailure := flags&filterNullOnFailure != 0
	filtered, ok := applyFilter(value, filter, flags)
	if ok {
		return filtered, nil
	}
	if nullOnFailure {
		return data.NewNullValue(), nil
	}
	name := "UnexpectedValueException"
	msg := `Parameter value "%s" is invalid and flag "FILTER_NULL_ON_FAILURE" was not set.`
	if badRequest {
		name = "Symfony\\Component\\HttpFoundation\\Exception\\BadRequestException"
		msg = `Input value "%s" is invalid and flag "FILTER_NULL_ON_FAILURE" was not set.`
	}
	return nil, throwNamed(name, msg, key)
}

func applyFilter(value data.Value, filter, flags int) (data.Value, bool) {
	if value == nil {
		return data.NewNullValue(), true
	}
	switch filter {
	case filterValidateInt:
		if flags&filterRequireScalar != 0 && isArrayValue(value) {
			return nil, false
		}
		switch v := value.(type) {
		case *data.IntValue:
			return v, true
		case *data.StringValue:
			n, err := strconv.Atoi(strings.TrimSpace(v.Value))
			if err != nil {
				return nil, false
			}
			return data.NewIntValue(n), true
		case *data.FloatValue:
			f, _ := v.AsFloat()
			if float64(int(f)) != f {
				return nil, false
			}
			return data.NewIntValue(int(f)), true
		case *data.BoolValue:
			b, _ := v.AsBool()
			if b {
				return data.NewIntValue(1), true
			}
			return data.NewIntValue(0), true
		default:
			n, err := strconv.Atoi(strings.TrimSpace(value.AsString()))
			if err != nil {
				return nil, false
			}
			return data.NewIntValue(n), true
		}
	case filterValidateBool:
		if flags&filterRequireScalar != 0 && isArrayValue(value) {
			return nil, false
		}
		switch v := value.(type) {
		case *data.BoolValue:
			return v, true
		case *data.IntValue:
			n, _ := v.AsInt()
			return data.NewBoolValue(n != 0), true
		case *data.StringValue:
			s := strings.ToLower(strings.TrimSpace(v.Value))
			switch s {
			case "1", "true", "on", "yes":
				return data.NewBoolValue(true), true
			case "0", "false", "off", "no", "":
				return data.NewBoolValue(false), true
			default:
				return nil, false
			}
		default:
			s := strings.ToLower(strings.TrimSpace(value.AsString()))
			switch s {
			case "1", "true", "on", "yes":
				return data.NewBoolValue(true), true
			case "0", "false", "off", "no", "":
				return data.NewBoolValue(false), true
			default:
				return nil, false
			}
		}
	default:
		// FILTER_DEFAULT：原样返回
		return value, true
	}
}

func paramGetString(store *ParamBagData, key, def string) (string, data.Control) {
	v := data.NewStringValue(def)
	if store != nil {
		if got, ok := store.get(key); ok {
			v = got
		} else {
			return def, nil
		}
	}
	s, err := valueAsString(v)
	if err != nil {
		return "", throwNamed("UnexpectedValueException", `Parameter value "%s" cannot be converted to "string".`, key)
	}
	return s, nil
}
