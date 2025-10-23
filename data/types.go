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
		if len(arr.Value) != len(m.Types) {
			return false
		}
		for i, typ := range m.Types {
			if !typ.Is(arr.Value[i]) {
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

func ISBaseType(ty string) bool {
	switch ty {
	case "":
		return true
	case "void":
		return true
	case "int":
		return true
	case "string":
		return true
	case "bool":
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
	case "void":
		return nil
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
