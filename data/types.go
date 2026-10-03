package data

import (
	"strings"

	"github.com/php-any/origami/tooling/typeinfo"
)

// Types is a cold compatibility view; runtime declarations store TypeRef.
type Types = typeinfo.Type

func NewLspTypes(t Types) *LspTypes {
	return &LspTypes{
		Types: []Types{t},
	}
}

// Deprecated compatibility names for tooling, with no runtime predicates.
type LspTypes = typeinfo.Inferred

// NullableType 表示可空类型
type NullableType struct {
	BaseType Types
}

func (n NullableType) Is(value Value) bool {
	// 可空类型可以接受 null 值或基础类型的值
	if _, ok := value.(*NullValue); ok {
		return true
	}
	return DeclaredTypeRef(n.BaseType).Matches(value, nil)
}

func (n NullableType) String() string {
	return "?" + n.BaseType.String()
}

type MultipleReturnType = typeinfo.Tuple

// UnionType 表示联合类型（type1|type2|...）
type UnionType struct {
	Types []Types
}

func (u UnionType) Is(value Value) bool {
	for _, t := range u.Types {
		if DeclaredTypeRef(t).Matches(value, nil) {
			return true
		}
	}
	return false
}

func (u UnionType) String() string {
	result := ""
	for i, t := range u.Types {
		if i > 0 {
			result += "|"
		}
		result += t.String()
	}
	return result
}

func NewUnionType(types []Types) Types {
	return UnionType{Types: types}
}

// IntersectionType 表示交集类型（type1&type2&...），值须同时满足所有成员类型
type IntersectionType struct {
	Types []Types
}

func (i IntersectionType) Is(value Value) bool {
	for _, t := range i.Types {
		if !DeclaredTypeRef(t).Matches(value, nil) {
			return false
		}
	}
	return len(i.Types) > 0
}

func (i IntersectionType) String() string {
	result := ""
	for idx, t := range i.Types {
		if idx > 0 {
			result += "&"
		}
		result += t.String()
	}
	return result
}

func NewIntersectionType(types []Types) Types {
	return IntersectionType{Types: types}
}

func ISBaseType(ty string) bool {
	if _, ok := declaredBuiltinRef(ty); ok {
		return true
	}
	switch ty {
	case "":
		return true
	case "void":
		return true
	case "mixed":
		return true
	case "int":
		return true
	case "string":
		return true
	case "bool":
		return true
	case "false":
		return true
	case "array":
		return true
	case "object":
		return true
	case "float":
		return true
	case "callable":
		return true
	default:
		return false
	}
}

func NewBaseType(ty string) Types {
	switch ty {
	case "":
		return nil
	case "void", "mixed", "never", "false", "true", "self", "parent", "iterable":
		return NewDeclaredType(ty)
	case "int":
		return Int{}
	case "float":
		return Float{}
	case "string":
		return String{}
	case "bool":
		return Bool{}
	case "array":
		return Arrays{}
	case "object":
		return Object{}
	case "callable":
		return Callable{}
	case "static":
		return StaticType{}
	case "null":
		return NullType{}
	case "closure", "Closure":
		return ClosureType{}
	default:
		if strings.ContainsAny(ty, "|&?()") {
			return NewDeclaredType(ty)
		}
		return Class{Name: ty}
	}
}

// NewNullableType 创建可空类型
func NewNullableType(baseType Types) Types {
	return NullableType{BaseType: baseType}
}

// NewMultipleReturnType 创建多返回值类型
func NewMultipleReturnType(types []Types) Types {
	return MultipleReturnType{Types: types}
}

func NewGenericType(name string, types []Types) Types {
	switch name {
	case "":
		return nil
	case "void":
		return TypeVoid
	case "mixed":
		return TypeMixed
	case "int":
		return Int{}
	case "string":
		return String{}
	case "bool":
		return Bool{}
	case "array":
		return Arrays{}
	case "object":
		return Object{}
	default:
		return Generic{Name: name, Types: types}
	}
}

// StaticType 表示 static 返回类型（PHP 8.0+）
// static 类型表示返回调用该方法的类的实例
type StaticType struct{}

func (s StaticType) Is(value Value) bool {
	// static 类型直接返回 true，实际的类型检查在方法调用时进行（在 ClassMethod.Call 中）
	return true
}

func (s StaticType) String() string {
	return "static"
}

type ClosureType struct{}

func (s ClosureType) Is(value Value) bool {
	switch v := value.(type) {
	case *FuncValue, *BoundFuncValue:
		return true
	case *ClassValue:
		return NominalIsA(v.Class, "Closure", v.GetVM())
	}
	return false
}

func (s ClosureType) String() string {
	return "closure"
}

type NullType struct{}

func (s NullType) Is(value Value) bool {
	if _, ok := value.(*NullValue); ok {
		return true
	}
	return false
}

func (s NullType) String() string {
	return "null"
}
