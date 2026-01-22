package data

type Types interface {
	Is(value Value) bool
	// String 范围标识识别是什么类型, 泛型类返回不需要泛型信息
	String() string
}

func NewLspTypes(t Types) *LspTypes {
	return &LspTypes{
		Types: []Types{t},
	}
}

// LspTypes 多种可能的类型 - 只能 lsp 使用
type LspTypes struct {
	Types []Types
}

func (l *LspTypes) Is(_ Value) bool {
	return true
}

func (l *LspTypes) String() string {
	return "LspTypes"
}

func (l *LspTypes) Add(t Types) {
	l.Types = append(l.Types, t)
}

// NullableType 表示可空类型
type NullableType struct {
	BaseType Types
}

func (n NullableType) Is(value Value) bool {
	// 可空类型可以接受 null 值或基础类型的值
	if _, ok := value.(*NullValue); ok {
		return true
	}
	return n.BaseType.Is(value)
}

func (n NullableType) String() string {
	return "?" + n.BaseType.String()
}

// MultipleReturnType 表示多返回值类型
type MultipleReturnType struct {
	Types []Types
}

func (m MultipleReturnType) Is(value Value) bool {
	// 多返回值类型检查数组中的每个元素
	if arr, ok := value.(*ArrayValue); ok {
		if len(arr.List) != len(m.Types) {
			return false
		}
		for i, typ := range m.Types {
			if !typ.Is(arr.List[i].Value) {
				return false
			}
		}
		return true
	}
	return false
}

func (m MultipleReturnType) String() string {
	result := ""
	for i, typ := range m.Types {
		if i > 0 {
			result += ", "
		}
		result += typ.String()
	}
	return result
}

// UnionType 表示联合类型（type1|type2|...）
type UnionType struct {
	Types []Types
}

func (u UnionType) Is(value Value) bool {
	for _, t := range u.Types {
		if t.Is(value) {
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

func ISBaseType(ty string) bool {
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
		// 允许 false 作为类型声明，语义等同于 bool
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
	case "void", "mixed":
		return nil
	case "int":
		return Int{}
	case "float":
		return Float{}
	case "string":
		return String{}
	case "bool":
		return Bool{}
	case "false":
		// false 类型声明等同于 bool
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
	case "self":
		return StaticType{}
	case "closure", "\\Closure":
		return ClosureType{}
	default:
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
		return nil
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
	switch value.(type) {
	case *FuncValue, *ArrayValue:
		return true
	case *StringValue:
		return true
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
