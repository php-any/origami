package reflection

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterAllowsNullMethod 实现 ReflectionParameter::allowsNull。
// PHP：无类型、?T、T|null、mixed，或默认值为 null 时返回 true。
type ReflectionParameterAllowsNullMethod struct{}

func (m *ReflectionParameterAllowsNullMethod) GetName() string { return "allowsNull" }

func (m *ReflectionParameterAllowsNullMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ReflectionParameterAllowsNullMethod) GetIsStatic() bool { return false }

func (m *ReflectionParameterAllowsNullMethod) GetParams() []data.GetValue { return []data.GetValue{} }

func (m *ReflectionParameterAllowsNullMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ReflectionParameterAllowsNullMethod) GetReturnType() data.Types { return data.Bool{} }

func (m *ReflectionParameterAllowsNullMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(paramAllowsNull(param)), nil
}

func paramAllowsNull(param data.GetValue) bool {
	if param == nil {
		return true
	}
	if typesAllowsNull(paramTypeOf(param)) {
		return true
	}
	return paramDefaultIsNull(param)
}

func paramTypeOf(param data.GetValue) data.Types {
	if vp, ok := param.(*virtualParam); ok {
		return vp.GetType()
	}
	if paramVar, ok := param.(data.Variable); ok {
		return paramVar.GetType()
	}
	if paramInterface, ok := param.(data.Parameter); ok {
		return paramInterface.GetType()
	}
	if paramNode, ok := param.(*node.Parameter); ok {
		return paramNode.GetType()
	}
	if paramNodes, ok := param.(*node.Parameters); ok {
		return paramNodes.GetType()
	}
	return nil
}

func typesAllowsNull(t data.Types) bool {
	if t == nil {
		return true
	}
	switch x := t.(type) {
	case data.NullableType:
		return true
	case *data.NullableType:
		return true
	case data.UnionType:
		for _, m := range x.Types {
			if typesAllowsNull(m) {
				return true
			}
		}
		return false
	case *data.UnionType:
		return typesAllowsNull(data.UnionType{Types: x.Types})
	default:
		s := strings.ToLower(strings.TrimSpace(t.String()))
		if s == "" || s == "mixed" || s == "null" {
			return true
		}
		if strings.HasPrefix(s, "?") {
			return true
		}
		for _, part := range strings.Split(s, "|") {
			if strings.TrimSpace(part) == "null" {
				return true
			}
		}
		return false
	}
}

func paramDefaultIsNull(param data.GetValue) bool {
	type defaultValGetter interface{ GetDefaultValue() data.GetValue }
	dg, ok := param.(defaultValGetter)
	if !ok {
		return false
	}
	def := dg.GetDefaultValue()
	if def == nil {
		return false
	}
	if _, ok := def.(*data.NullValue); ok {
		return true
	}
	return false
}
