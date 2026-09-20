package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type ReflectionPropertyClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func (c *ReflectionPropertyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ReflectionPropertyClass) GetName() string                               { return "ReflectionProperty" }
func (c *ReflectionPropertyClass) GetExtend() *string                            { return nil }
func (c *ReflectionPropertyClass) GetImplements() []string                       { return nil }
func (c *ReflectionPropertyClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ReflectionPropertyClass) GetPropertyList() []data.Property              { return nil }

// GetStaticProperty 对齐 PHP ReflectionProperty::IS_* 类常量（:: 访问走静态属性查找）
func (c *ReflectionPropertyClass) GetStaticProperty(name string) (data.Value, bool) {
	if c.StaticProperty == nil {
		c.StaticProperty = reflectionPropertyConstants()
	}
	v, ok := c.StaticProperty[name]
	return v, ok
}

func reflectionPropertyConstants() map[string]data.Value {
	return map[string]data.Value{
		"IS_PUBLIC":        data.NewIntValue(1),
		"IS_PROTECTED":     data.NewIntValue(2),
		"IS_PRIVATE":       data.NewIntValue(4),
		"IS_FINAL":         data.NewIntValue(32),
		"IS_ABSTRACT":      data.NewIntValue(64),
		"IS_STATIC":        data.NewIntValue(16),
		"IS_READONLY":      data.NewIntValue(128),
		"IS_PROTECTED_SET": data.NewIntValue(2048),
		"IS_PRIVATE_SET":   data.NewIntValue(4096),
		"IS_VIRTUAL":       data.NewIntValue(16384),
	}
}
func (c *ReflectionPropertyClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return &ReflectionPropertyConstructMethod{}, true
	case "setAccessible":
		return &ReflectionPropertySetAccessibleMethod{}, true
	case "getName":
		return &ReflectionPropertyGetNameMethod{}, true
	case "getDeclaringClass":
		return &ReflectionPropertyGetDeclaringClassMethod{}, true
	case "getValue":
		return &ReflectionPropertyGetValueMethod{}, true
	case "setValue":
		return &ReflectionPropertySetValueMethod{}, true
	case "setRawValue":
		// PHP 8.4：绕过 set hook；当前无 hook 运行时，行为同 setValue
		return &ReflectionPropertySetValueMethod{raw: true}, true
	case "getRawValue":
		return &ReflectionPropertyGetValueMethod{raw: true}, true
	case "isPublic":
		return &ReflectionPropertyIsPublicMethod{}, true
	case "isProtected":
		return &ReflectionPropertyIsProtectedMethod{}, true
	case "isPrivate":
		return &ReflectionPropertyIsPrivateMethod{}, true
	case "isStatic":
		return &ReflectionPropertyIsStaticMethod{}, true
	case "isDefault":
		return &ReflectionPropertyIsDefaultMethod{}, true
	case "isInitialized":
		return &ReflectionPropertyIsInitializedMethod{}, true
	case "hasType":
		return &ReflectionPropertyHasTypeMethod{}, true
	case "getType":
		return &ReflectionPropertyGetTypeMethod{}, true
	case "getAttributes":
		return &ReflectionPropertyGetAttributesMethod{}, true
	case "isVirtual":
		return &ReflectionPropertyIsVirtualMethod{}, true
	}
	return nil, false
}
func (c *ReflectionPropertyClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionPropertyConstructMethod{},
		&ReflectionPropertySetAccessibleMethod{},
		&ReflectionPropertyGetNameMethod{},
		&ReflectionPropertyGetDeclaringClassMethod{},
		&ReflectionPropertyGetValueMethod{},
		&ReflectionPropertySetValueMethod{},
		&ReflectionPropertySetValueMethod{raw: true},
		&ReflectionPropertyIsPublicMethod{},
		&ReflectionPropertyIsProtectedMethod{},
		&ReflectionPropertyIsPrivateMethod{},
		&ReflectionPropertyIsStaticMethod{},
		&ReflectionPropertyIsDefaultMethod{},
		&ReflectionPropertyIsInitializedMethod{},
		&ReflectionPropertyHasTypeMethod{},
		&ReflectionPropertyGetTypeMethod{},
		&ReflectionPropertyGetAttributesMethod{},
		&ReflectionPropertyIsVirtualMethod{},
	}
}

func newReflectionProperty(ctx data.Context, className, propertyName string) *data.ClassValue {
	propertyClass := &ReflectionPropertyClass{}
	propertyValue := data.NewClassValue(propertyClass, ctx.CreateBaseContext())
	propertyValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	propertyValue.ObjectValue.SetProperty("_propertyName", data.NewStringValue(propertyName))
	propertyValue.ObjectValue.SetProperty("name", data.NewStringValue(propertyName))
	propertyValue.ObjectValue.SetProperty("class", data.NewStringValue(className))
	return propertyValue
}
func (c *ReflectionPropertyClass) GetConstruct() data.Method {
	return &ReflectionPropertyConstructMethod{}
}

type ReflectionPropertyConstructMethod struct{}

func (m *ReflectionPropertyConstructMethod) GetName() string            { return "__construct" }
func (m *ReflectionPropertyConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyConstructMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertyConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, nil),
		node.NewParameter(nil, "property", 1, nil, nil),
	}
}
func (m *ReflectionPropertyConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.Mixed{}),
		node.NewVariable(nil, "property", 1, data.Mixed{}),
	}
}
func (m *ReflectionPropertyConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	classVal, _ := ctx.GetIndexValue(0)
	propVal, _ := ctx.GetIndexValue(1)
	if classVal == nil || propVal == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionProperty::__construct() expects class and property"))
	}

	className := classVal.AsString()
	if cv, ok := classVal.(*data.ClassValue); ok && cv.Class != nil {
		className = cv.Class.GetName()
	}
	propName := propVal.AsString()

	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		cmc.ObjectValue.SetProperty("_className", data.NewStringValue(className))
		cmc.ObjectValue.SetProperty("_propertyName", data.NewStringValue(propName))
		cmc.ObjectValue.SetProperty("name", data.NewStringValue(propName))
		cmc.ObjectValue.SetProperty("class", data.NewStringValue(className))
	}
	return nil, nil
}

type ReflectionPropertySetAccessibleMethod struct{}

func (m *ReflectionPropertySetAccessibleMethod) GetName() string { return "setAccessible" }
func (m *ReflectionPropertySetAccessibleMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionPropertySetAccessibleMethod) GetIsStatic() bool         { return false }
func (m *ReflectionPropertySetAccessibleMethod) GetReturnType() data.Types { return nil }
func (m *ReflectionPropertySetAccessibleMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "accessible", 0, nil, nil),
	}
}
func (m *ReflectionPropertySetAccessibleMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "accessible", 0, data.Mixed{}),
	}
}
func (m *ReflectionPropertySetAccessibleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

type ReflectionPropertyGetNameMethod struct{}

func (m *ReflectionPropertyGetNameMethod) GetName() string            { return "getName" }
func (m *ReflectionPropertyGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyGetNameMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyGetNameMethod) GetReturnType() data.Types  { return data.String{} }
func (m *ReflectionPropertyGetNameMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyGetNameMethod) GetVariables() []data.Variable {
	return nil
}

type ReflectionPropertyGetDeclaringClassMethod struct{}

func (m *ReflectionPropertyGetDeclaringClassMethod) GetName() string { return "getDeclaringClass" }
func (m *ReflectionPropertyGetDeclaringClassMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionPropertyGetDeclaringClassMethod) GetIsStatic() bool { return false }
func (m *ReflectionPropertyGetDeclaringClassMethod) GetReturnType() data.Types {
	return data.Mixed{}
}
func (m *ReflectionPropertyGetDeclaringClassMethod) GetParams() []data.GetValue {
	return nil
}
func (m *ReflectionPropertyGetDeclaringClassMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyGetDeclaringClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if className, ctl := cmc.ObjectValue.GetProperty("_className"); ctl == nil && className != nil {
			return newReflectionClassValue(ctx, className.AsString()), nil
		}
	}
	return data.NewNullValue(), nil
}
func (m *ReflectionPropertyGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if v, ctl := cmc.ObjectValue.GetProperty("_propertyName"); ctl == nil && v != nil {
			return v, nil
		}
	}
	return data.NewStringValue(""), nil
}

type ReflectionPropertyGetValueMethod struct {
	raw bool
}

func (m *ReflectionPropertyGetValueMethod) GetName() string {
	if m.raw {
		return "getRawValue"
	}
	return "getValue"
}
func (m *ReflectionPropertyGetValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyGetValueMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyGetValueMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertyGetValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, data.NewNullValue(), nil),
	}
}
func (m *ReflectionPropertyGetValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.Mixed{}),
	}
}
func (m *ReflectionPropertyGetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propName := reflectionPropertyName(ctx)
	objVal, _ := ctx.GetIndexValue(0)
	if propName == "" || objVal == nil {
		return data.NewNullValue(), nil
	}
	switch o := objVal.(type) {
	case *data.ClassValue:
		if o.ObjectValue != nil {
			if v, _ := o.ObjectValue.GetProperty(propName); v != nil {
				return v, nil
			}
		}
		if z, ctl := o.GetPropertyZVal(propName); ctl == nil && z != nil && z.Value != nil {
			return z.Value, nil
		}
	case data.GetProperty:
		if v, ctl := o.GetProperty(propName); ctl == nil && v != nil {
			return v, nil
		}
	}
	return data.NewNullValue(), nil
}

type ReflectionPropertySetValueMethod struct {
	raw bool
}

func (m *ReflectionPropertySetValueMethod) GetName() string {
	if m.raw {
		return "setRawValue"
	}
	return "setValue"
}
func (m *ReflectionPropertySetValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertySetValueMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertySetValueMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertySetValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (m *ReflectionPropertySetValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.Mixed{}),
		node.NewVariable(nil, "value", 1, data.Mixed{}),
	}
}
func (m *ReflectionPropertySetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propName := reflectionPropertyName(ctx)
	objVal, _ := ctx.GetIndexValue(0)
	value, _ := ctx.GetIndexValue(1)
	if propName == "" || objVal == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionProperty::setValue() missing object or property"))
	}
	if value == nil {
		value = data.NewNullValue()
	}
	switch o := objVal.(type) {
	case *data.ClassValue:
		return nil, o.SetProperty(propName, value)
	case data.SetProperty:
		return nil, o.SetProperty(propName, value)
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionProperty::setValue() expects object"))
}

func reflectionPropertyName(ctx data.Context) string {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if v, ctl := cmc.ObjectValue.GetProperty("_propertyName"); ctl == nil && v != nil {
			return v.AsString()
		}
	}
	return ""
}

func reflectionPropertyClassName(ctx data.Context) string {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if v, ctl := cmc.ObjectValue.GetProperty("_className"); ctl == nil && v != nil {
			return v.AsString()
		}
	}
	return ""
}

// reflectionPropertyInfo 沿继承链查找属性定义
func reflectionPropertyInfo(ctx data.Context) data.Property {
	className := reflectionPropertyClassName(ctx)
	propName := reflectionPropertyName(ctx)
	if className == "" || propName == "" {
		return nil
	}
	vm := ctx.GetVM()
	if vm == nil {
		return nil
	}
	v, acl := vm.LoadPkg(className)
	if acl != nil || v == nil {
		return nil
	}
	stmt, ok := v.(data.ClassStmt)
	if !ok {
		return nil
	}
	current := stmt
	for current != nil {
		if p, found := current.GetProperty(propName); found {
			return p
		}
		ext := current.GetExtend()
		if ext == nil || *ext == "" {
			break
		}
		ev, eacl := vm.LoadPkg(*ext)
		if eacl != nil || ev == nil {
			break
		}
		es, ok := ev.(data.ClassStmt)
		if !ok {
			break
		}
		current = es
	}
	return nil
}

type ReflectionPropertyIsPublicMethod struct{}

func (m *ReflectionPropertyIsPublicMethod) GetName() string            { return "isPublic" }
func (m *ReflectionPropertyIsPublicMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyIsPublicMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyIsPublicMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyIsPublicMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyIsPublicMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyIsPublicMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prop := reflectionPropertyInfo(ctx)
	if prop == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(prop.GetModifier() == data.ModifierPublic), nil
}

type ReflectionPropertyIsProtectedMethod struct{}

func (m *ReflectionPropertyIsProtectedMethod) GetName() string            { return "isProtected" }
func (m *ReflectionPropertyIsProtectedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyIsProtectedMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyIsProtectedMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyIsProtectedMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyIsProtectedMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyIsProtectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prop := reflectionPropertyInfo(ctx)
	if prop == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(prop.GetModifier() == data.ModifierProtected), nil
}

type ReflectionPropertyIsPrivateMethod struct{}

func (m *ReflectionPropertyIsPrivateMethod) GetName() string            { return "isPrivate" }
func (m *ReflectionPropertyIsPrivateMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyIsPrivateMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyIsPrivateMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyIsPrivateMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyIsPrivateMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyIsPrivateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prop := reflectionPropertyInfo(ctx)
	if prop == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(prop.GetModifier() == data.ModifierPrivate), nil
}

type ReflectionPropertyIsStaticMethod struct{}

func (m *ReflectionPropertyIsStaticMethod) GetName() string            { return "isStatic" }
func (m *ReflectionPropertyIsStaticMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyIsStaticMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyIsStaticMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyIsStaticMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyIsStaticMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyIsStaticMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prop := reflectionPropertyInfo(ctx)
	if prop == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(prop.GetIsStatic()), nil
}

type ReflectionPropertyIsDefaultMethod struct{}

func (m *ReflectionPropertyIsDefaultMethod) GetName() string            { return "isDefault" }
func (m *ReflectionPropertyIsDefaultMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyIsDefaultMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyIsDefaultMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyIsDefaultMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyIsDefaultMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyIsDefaultMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// PHP：isDefault 表示属性在类里声明（非运行期动态属性），与是否有默认值无关。
	// Livewire getPublicPropertiesDefinedOnSubclass 会过滤 !isDefault()。
	if reflectionPropertyName(ctx) != "" {
		return data.NewBoolValue(true), nil
	}
	prop := reflectionPropertyInfo(ctx)
	return data.NewBoolValue(prop != nil), nil
}

type ReflectionPropertyIsInitializedMethod struct{}

func (m *ReflectionPropertyIsInitializedMethod) GetName() string { return "isInitialized" }
func (m *ReflectionPropertyIsInitializedMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionPropertyIsInitializedMethod) GetIsStatic() bool         { return false }
func (m *ReflectionPropertyIsInitializedMethod) GetReturnType() data.Types { return data.Bool{} }
func (m *ReflectionPropertyIsInitializedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, node.NewNullLiteral(nil), nil),
	}
}
func (m *ReflectionPropertyIsInitializedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.Mixed{}),
	}
}
func unwrapReflectionValue(v data.Value) data.Value {
	for i := 0; i < 4 && v != nil; i++ {
		zv, ok := v.(*data.ZValValue)
		if !ok || zv.ZVal == nil {
			return v
		}
		v = zv.ZVal.Value
	}
	return v
}

func classValueFromReflectionObject(objVal data.Value) *data.ClassValue {
	switch o := unwrapReflectionValue(objVal).(type) {
	case *data.ThisValue:
		return o.ClassValue
	case *data.ClassValue:
		return o
	default:
		return nil
	}
}

func (m *ReflectionPropertyIsInitializedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propName := reflectionPropertyName(ctx)
	objVal, _ := ctx.GetIndexValue(0)
	cv := classValueFromReflectionObject(objVal)
	if propName == "" || cv == nil {
		return data.NewBoolValue(false), nil
	}
	if cv.ObjectValue != nil && cv.ObjectValue.HasProperty(propName) {
		return data.NewBoolValue(true), nil
	}
	stmt, ok := cv.GetPropertyStmt(propName)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	// PHP：无类型属性视为已初始化（默认 null）；有类型无默认值则在赋值前为未初始化。
	if stmt.GetType() == nil || stmt.GetDefaultValue() != nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

type ReflectionPropertyHasTypeMethod struct{}

func (m *ReflectionPropertyHasTypeMethod) GetName() string            { return "hasType" }
func (m *ReflectionPropertyHasTypeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyHasTypeMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyHasTypeMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyHasTypeMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyHasTypeMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyHasTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prop := reflectionPropertyInfo(ctx)
	if prop == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(prop.GetType() != nil), nil
}

type ReflectionPropertyIsVirtualMethod struct{}

func (m *ReflectionPropertyIsVirtualMethod) GetName() string            { return "isVirtual" }
func (m *ReflectionPropertyIsVirtualMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyIsVirtualMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyIsVirtualMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ReflectionPropertyIsVirtualMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyIsVirtualMethod) GetVariables() []data.Variable {
	return nil
}

// Call 实现 isVirtual。PHP 8.4 的虚拟属性由 get/set hook 定义，origami 尚未支持，恒返回 false。
func (m *ReflectionPropertyIsVirtualMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

type ReflectionPropertyGetTypeMethod struct{}

func (m *ReflectionPropertyGetTypeMethod) GetName() string            { return "getType" }
func (m *ReflectionPropertyGetTypeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyGetTypeMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyGetTypeMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *ReflectionPropertyGetTypeMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyGetTypeMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyGetTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prop := reflectionPropertyInfo(ctx)
	if prop == nil || prop.GetType() == nil {
		return data.NewNullValue(), nil
	}
	return newPhpReflectionType(ctx, prop.GetType()), nil
}
