package reflection

import (
	"crypto/sha1"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionReferenceClass 提供 PHP ReflectionReference 类定义
// 用于检测数组元素是否为引用（Symfony VarDumper 等调试工具依赖）。
type ReflectionReferenceClass struct {
	node.Node
}

func (c *ReflectionReferenceClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ReflectionReferenceClass) GetName() string { return "ReflectionReference" }

func (c *ReflectionReferenceClass) GetExtend() *string { return nil }

func (c *ReflectionReferenceClass) GetImplements() []string { return nil }

func (c *ReflectionReferenceClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ReflectionReferenceClass) GetPropertyList() []data.Property {
	return nil
}

func (c *ReflectionReferenceClass) GetConstruct() data.Method {
	return nil
}

func (c *ReflectionReferenceClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionReferenceGetIdMethod{},
	}
}

func (c *ReflectionReferenceClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "fromArrayElement":
		return &ReflectionReferenceFromArrayElementMethod{}, true
	}
	return nil, false
}

func (c *ReflectionReferenceClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "getId":
		return &ReflectionReferenceGetIdMethod{}, true
	}
	return nil, false
}

// ReflectionReferenceFromArrayElementMethod 实现 ReflectionReference::fromArrayElement
// 若数组 $array 的 $key 元素为引用则返回 ReflectionReference 实例，否则返回 null
type ReflectionReferenceFromArrayElementMethod struct{}

func (m *ReflectionReferenceFromArrayElementMethod) GetName() string { return "fromArrayElement" }
func (m *ReflectionReferenceFromArrayElementMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionReferenceFromArrayElementMethod) GetIsStatic() bool { return true }

var reflectionReferenceFromArrayElementMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.TypeArray),
	node.NewParameter(nil, "key", 1, nil, data.NewDeclaredUnionType([]data.Types{data.TypeInt, data.TypeString})),
}

func (m *ReflectionReferenceFromArrayElementMethod) GetParams() []data.GetValue {
	return reflectionReferenceFromArrayElementMethodGetParams
}

var reflectionReferenceFromArrayElementMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.TypeArray),
	node.NewVariable(nil, "key", 1, data.TypeMixed),
}

func (m *ReflectionReferenceFromArrayElementMethod) GetVariables() []data.Variable {
	return reflectionReferenceFromArrayElementMethodGetVariables
}
func (m *ReflectionReferenceFromArrayElementMethod) GetReturnType() data.Types {
	return data.NewBaseType("?ReflectionReference")
}

func (m *ReflectionReferenceFromArrayElementMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrArg, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewNullValue(), nil
	}
	arr, ok := arrArg.(*data.ArrayValue)
	if !ok {
		return data.NewNullValue(), nil
	}
	keyArg, ok := ctx.GetIndexValue(1)
	if !ok {
		return data.NewNullValue(), nil
	}

	var slot *data.ZVal
	switch key := keyArg.(type) {
	case *data.IntValue:
		slot, _ = arr.FindSlotByIntKey(key.Value)
	case *data.StringValue:
		// ReflectionReference uses the literal key kind, unlike PHP's
		// ordinary array access which converts canonical numeric strings.
		if _, numeric := data.ParseIntArrayKeyName(key.Value); !numeric {
			slot, _ = arr.LookupZValByStringKey(key.Value)
		}
	default:
		return nil, data.NewTypeError(nil, fmt.Errorf("ReflectionReference::fromArrayElement(): key must be int|string"))
	}
	if slot == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Array key not found"), "ReflectionException")
	}
	if slot.RefCount() > 0 {
		return newReflectionReferenceValue(ctx, slot), nil
	}
	return data.NewNullValue(), nil
}
func newReflectionReferenceValue(ctx data.Context, slot *data.ZVal) data.GetValue {
	id := sha1.Sum([]byte(fmt.Sprintf("%p", slot.ReferenceIdentity())))
	value := data.NewClassValue(&ReflectionReferenceClass{}, ctx.CreateBaseContext())
	value.ObjectValue.SetProperty("_id", data.NewStringValue(string(id[:])))
	return value
}

// ReflectionReferenceGetIdMethod 实现 ReflectionReference::getId
type ReflectionReferenceGetIdMethod struct{}

func (m *ReflectionReferenceGetIdMethod) GetName() string { return "getId" }
func (m *ReflectionReferenceGetIdMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionReferenceGetIdMethod) GetIsStatic() bool { return false }

var reflectionReferenceGetIdMethodGetParams = []data.GetValue{}

func (m *ReflectionReferenceGetIdMethod) GetParams() []data.GetValue {
	return reflectionReferenceGetIdMethodGetParams
}

var reflectionReferenceGetIdMethodGetVariables = []data.Variable{}

func (m *ReflectionReferenceGetIdMethod) GetVariables() []data.Variable {
	return reflectionReferenceGetIdMethodGetVariables
}
func (m *ReflectionReferenceGetIdMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ReflectionReferenceGetIdMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if id, ctl := cmc.ObjectValue.GetProperty("_id"); ctl == nil && id != nil {
			return id, nil
		}
	}
	return data.NewStringValue(""), nil
}
