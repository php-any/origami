package reflection

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetTraitsMethod 实现 ReflectionClass::getTraits
// 返回类直接使用的 trait 列表：key 为 trait 短名，value 为 trait 的 ReflectionClass 实例
type ReflectionClassGetTraitsMethod struct{}

// GetName 返回方法名 "getTraits"
func (m *ReflectionClassGetTraitsMethod) GetName() string { return "getTraits" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetTraitsMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetTraitsMethod) GetIsStatic() bool { return false }

var reflectionClassGetTraitsMethodGetParams = []data.GetValue{}

// GetParams 返回参数列表
func (m *ReflectionClassGetTraitsMethod) GetParams() []data.GetValue {
	return reflectionClassGetTraitsMethodGetParams
}

var reflectionClassGetTraitsMethodGetVariables = []data.Variable{}

// GetVariables 返回变量列表
func (m *ReflectionClassGetTraitsMethod) GetVariables() []data.Variable {
	return reflectionClassGetTraitsMethodGetVariables
}

// GetReturnType 返回返回类型，返回数组类型
func (m *ReflectionClassGetTraitsMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// Call 执行 getTraits 方法
func (m *ReflectionClassGetTraitsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	list := make([]*data.ZVal, 0)
	if classStmt == nil {
		return &data.ArrayValue{List: list}, nil
	}
	cs, ok := classStmt.(*node.ClassStatement)
	if !ok {
		return &data.ArrayValue{List: list}, nil
	}
	for _, traitName := range cs.Traits {
		shortName := traitName
		if idx := strings.LastIndex(traitName, "\\"); idx >= 0 {
			shortName = traitName[idx+1:]
		}
		traitValue := newReflectionClassValue(ctx, traitName)
		list = append(list, data.NewNamedZVal(shortName, traitValue))
	}
	return &data.ArrayValue{List: list}, nil
}
