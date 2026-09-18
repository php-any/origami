package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassIsAbstractMethod 实现 ReflectionClass::isAbstract
type ReflectionClassIsAbstractMethod struct{}

func (m *ReflectionClassIsAbstractMethod) GetName() string               { return "isAbstract" }
func (m *ReflectionClassIsAbstractMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ReflectionClassIsAbstractMethod) GetIsStatic() bool             { return false }
func (m *ReflectionClassIsAbstractMethod) GetParams() []data.GetValue    { return nil }
func (m *ReflectionClassIsAbstractMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassIsAbstractMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *ReflectionClassIsAbstractMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}
	if _, ok := classStmt.(*node.AbstractClassStatement); ok {
		return data.NewBoolValue(true), nil
	}
	// 抽象类注册为 *ClassStatement，用 IsAbstract 标记（见 class_register.go）
	if c, ok := classStmt.(*node.ClassStatement); ok {
		return data.NewBoolValue(c.IsAbstract), nil
	}
	if c, ok := classStmt.(*node.ClassGeneric); ok {
		return data.NewBoolValue(c.IsAbstract), nil
	}
	return data.NewBoolValue(false), nil
}

// ReflectionClassIsInterfaceMethod 实现 ReflectionClass::isInterface
type ReflectionClassIsInterfaceMethod struct{}

func (m *ReflectionClassIsInterfaceMethod) GetName() string               { return "isInterface" }
func (m *ReflectionClassIsInterfaceMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ReflectionClassIsInterfaceMethod) GetIsStatic() bool             { return false }
func (m *ReflectionClassIsInterfaceMethod) GetParams() []data.GetValue    { return nil }
func (m *ReflectionClassIsInterfaceMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassIsInterfaceMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *ReflectionClassIsInterfaceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, _ := getReflectionClassInfo(ctx)
	if className == "" {
		return data.NewBoolValue(false), nil
	}
	vm := ctx.GetVM()
	if _, ok := vm.GetInterface(className); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

// ReflectionClassIsTraitMethod 实现 ReflectionClass::isTrait
type ReflectionClassIsTraitMethod struct{}

func (m *ReflectionClassIsTraitMethod) GetName() string               { return "isTrait" }
func (m *ReflectionClassIsTraitMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ReflectionClassIsTraitMethod) GetIsStatic() bool             { return false }
func (m *ReflectionClassIsTraitMethod) GetParams() []data.GetValue    { return nil }
func (m *ReflectionClassIsTraitMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassIsTraitMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *ReflectionClassIsTraitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}
	// Trait 在 Origami 中通常以 ClassStatement 存储；若有 Trait 标记则识别
	if marker, ok := classStmt.(interface{ IsTrait() bool }); ok && marker.IsTrait() {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}
