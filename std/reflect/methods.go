package reflect

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetClassInfoMethod 获取类信息方法
type GetClassInfoMethod struct{}

func (g *GetClassInfoMethod) GetName() string {
	return "getClassInfo"
}

func (g *GetClassInfoMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetClassInfoMethod) GetIsStatic() bool {
	return false
}

func (g *GetClassInfoMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, data.NewBaseType("string")),
	}
}

func (g *GetClassInfoMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (g *GetClassInfoMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetClassInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类名参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 构建友好的类信息字符串
	info := fmt.Sprintf("类名: %s\n", className)

	if extend := class.GetExtend(); extend != nil {
		info += fmt.Sprintf("父类: %s\n", *extend)
	}

	if implements := class.GetImplements(); len(implements) > 0 {
		info += fmt.Sprintf("实现接口: %v\n", implements)
	}

	// 获取方法信息
	methods := class.GetMethods()
	info += fmt.Sprintf("方法数量: %d\n", len(methods))
	if len(methods) > 0 {
		info += "方法列表:\n"
		for i, method := range methods {
			info += fmt.Sprintf("  %d. %s (%s)\n", i+1, method.GetName(), getModifierString(method.GetModifier()))
		}
	}

	// 获取属性信息
	properties := class.GetPropertyList()
	info += fmt.Sprintf("属性数量: %d\n", len(properties))
	if len(properties) > 0 {
		info += "属性列表:\n"
		i := 1
		for _, property := range properties {
			info += fmt.Sprintf("  %d. %s\n", i, property.GetName())
			i++
		}
	}

	return data.NewStringValue(info), nil
}

// GetMethodInfoMethod 获取方法信息方法
type GetMethodInfoMethod struct{}

func (g *GetMethodInfoMethod) GetName() string {
	return "getMethodInfo"
}

func (g *GetMethodInfoMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetMethodInfoMethod) GetIsStatic() bool {
	return false
}

func (g *GetMethodInfoMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
		node.NewParameter(nil, "methodName", 1, nil, nil),
	}
}

func (g *GetMethodInfoMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
		node.NewVariable(nil, "methodName", 1, nil),
	}
}

// GetReturnType 返回方法返回类型
func (g *GetMethodInfoMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetMethodInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	methodNameValue, exists := ctx.GetIndexValue(1)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	methodName := methodNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 查找方法
	method, exists := class.GetMethod(methodName)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 构建方法信息对象
	methodInfo := make(map[string]interface{})
	methodInfo["name"] = method.GetName()
	methodInfo["modifier"] = getModifierString(method.GetModifier())
	methodInfo["isStatic"] = method.GetIsStatic()

	// 获取参数信息
	params := method.GetParams()
	methodInfo["paramCount"] = len(params)
	paramNames := make([]string, len(params))
	for i := range params {
		paramNames[i] = fmt.Sprintf("param%d", i+1)
	}
	methodInfo["paramNames"] = paramNames

	return data.NewStringValue(fmt.Sprintf("%v", methodInfo)), nil
}

// GetPropertyInfoMethod 获取属性信息方法
type GetPropertyInfoMethod struct{}

func (g *GetPropertyInfoMethod) GetName() string {
	return "getPropertyInfo"
}

func (g *GetPropertyInfoMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetPropertyInfoMethod) GetIsStatic() bool {
	return false
}

func (g *GetPropertyInfoMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
		node.NewParameter(nil, "propertyName", 1, nil, nil),
	}
}

func (g *GetPropertyInfoMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
		node.NewVariable(nil, "propertyName", 1, nil),
	}
}

// GetReturnType 返回方法返回类型
func (g *GetPropertyInfoMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetPropertyInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewObjectValue(), nil
	}

	propertyNameValue, exists := ctx.GetIndexValue(1)
	if !exists {
		return data.NewObjectValue(), nil
	}

	className := classNameValue.AsString()
	propertyName := propertyNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewObjectValue(), nil
	}

	// 查找属性
	property, exists := class.GetProperty(propertyName)
	if !exists {
		return data.NewObjectValue(), nil
	}

	obj := data.NewObjectValue()
	obj.SetProperty("name", data.NewStringValue(property.GetName()))
	obj.SetProperty("modifier", data.NewStringValue(getModifierString(property.GetModifier())))
	obj.SetProperty("isStatic", data.NewBoolValue(property.GetIsStatic()))

	// 从默认值推断类型
	defaultValue := property.GetDefaultValue()
	if defaultValue != nil {
		actualValue, _ := defaultValue.GetValue(ctx)
		if actualValue != nil {
			if strValue, ok := actualValue.(*data.StringValue); ok {
				obj.SetProperty("type", data.NewStringValue("string"))
				obj.SetProperty("defaultValue", strValue)
			} else if intValue, ok := actualValue.(*data.IntValue); ok {
				obj.SetProperty("type", data.NewStringValue("int"))
				obj.SetProperty("defaultValue", intValue)
			} else if floatValue, ok := actualValue.(*data.FloatValue); ok {
				obj.SetProperty("type", data.NewStringValue("float"))
				obj.SetProperty("defaultValue", floatValue)
			} else if boolValue, ok := actualValue.(*data.BoolValue); ok {
				obj.SetProperty("type", data.NewStringValue("bool"))
				obj.SetProperty("defaultValue", boolValue)
			} else if nullValue, ok := actualValue.(*data.NullValue); ok {
				obj.SetProperty("type", data.NewStringValue("null"))
				obj.SetProperty("defaultValue", nullValue)
			} else if v, ok := actualValue.(data.Value); ok {
				obj.SetProperty("type", data.NewStringValue("unknown"))
				obj.SetProperty("defaultValue", v)
			} else {
				obj.SetProperty("type", data.NewStringValue("unknown"))
				obj.SetProperty("defaultValue", data.NewStringValue("unknown"))
			}
		} else {
			obj.SetProperty("type", data.NewStringValue("mixed"))
			obj.SetProperty("defaultValue", data.NewNullValue())
		}
	} else {
		obj.SetProperty("type", data.NewStringValue("mixed"))
		obj.SetProperty("defaultValue", data.NewNullValue())
	}

	return obj, nil
}

// ListClassesMethod 列出所有类方法
type ListClassesMethod struct{}

func (l *ListClassesMethod) GetName() string {
	return "listClasses"
}

func (l *ListClassesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (l *ListClassesMethod) GetIsStatic() bool {
	return false
}

func (l *ListClassesMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (l *ListClassesMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (l *ListClassesMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (l *ListClassesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 由于VM接口没有提供获取所有类的方法，返回空数组
	return data.NewStringValue("[]"), nil
}

// ListMethodsMethod 列出类的所有方法
type ListMethodsMethod struct{}

func (l *ListMethodsMethod) GetName() string {
	return "listMethods"
}

func (l *ListMethodsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (l *ListMethodsMethod) GetIsStatic() bool {
	return false
}

func (l *ListMethodsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
	}
}

func (l *ListMethodsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (l *ListMethodsMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (l *ListMethodsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类名参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue("[]"), nil
	}

	className := classNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue("[]"), nil
	}

	// 获取所有方法
	methods := class.GetMethods()
	methodNames := make([]string, len(methods))
	for i, method := range methods {
		methodNames[i] = method.GetName()
	}

	return data.NewStringValue(fmt.Sprintf("%v", methodNames)), nil
}

// ListPropertiesMethod 列出类的所有属性
type ListPropertiesMethod struct{}

func (l *ListPropertiesMethod) GetName() string {
	return "listProperties"
}

func (l *ListPropertiesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (l *ListPropertiesMethod) GetIsStatic() bool {
	return false
}

func (l *ListPropertiesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
	}
}

func (l *ListPropertiesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (l *ListPropertiesMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (l *ListPropertiesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类名参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewArrayValue(nil), nil
	}

	className := classNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewArrayValue(nil), nil
	}

	// 获取所有属性
	properties := class.GetPropertyList()
	propertyValues := make([]data.Value, 0, len(properties))
	for _, property := range properties {
		propertyValues = append(propertyValues, data.NewStringValue(property.GetName()))
	}

	return data.NewArrayValue(propertyValues), nil
}

// getModifierString 获取修饰符字符串
func getModifierString(modifier data.Modifier) string {
	switch int(modifier) {
	case 0: // ModifierPublic
		return "public"
	case 1: // ModifierPrivate
		return "private"
	case 2: // ModifierProtected
		return "protected"
	default:
		return "unknown"
	}
}

// GetClassAnnotationsMethod 获取类注解信息方法
type GetClassAnnotationsMethod struct{}

func (g *GetClassAnnotationsMethod) GetName() string {
	return "getClassAnnotations"
}

func (g *GetClassAnnotationsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetClassAnnotationsMethod) GetIsStatic() bool {
	return false
}

func (g *GetClassAnnotationsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
	}
}

func (g *GetClassAnnotationsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
	}
}

func (g *GetClassAnnotationsMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetClassAnnotationsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类名参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 检查类是否有注解
	if classStmt, ok := class.(*node.ClassStatement); ok {
		if classStmt.Annotations != nil && len(classStmt.Annotations) > 0 {
			info := fmt.Sprintf("类 %s 的注解信息:\n", className)
			for i, annotation := range classStmt.Annotations {
				info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
			}
			return data.NewStringValue(info), nil
		}
	}

	return data.NewStringValue(fmt.Sprintf("类 %s 没有注解", className)), nil
}

// GetMethodAnnotationsMethod 获取方法注解信息方法
type GetMethodAnnotationsMethod struct{}

func (g *GetMethodAnnotationsMethod) GetName() string {
	return "getMethodAnnotations"
}

func (g *GetMethodAnnotationsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetMethodAnnotationsMethod) GetIsStatic() bool {
	return false
}

func (g *GetMethodAnnotationsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
		node.NewParameter(nil, "methodName", 1, nil, nil),
	}
}

func (g *GetMethodAnnotationsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
		node.NewVariable(nil, "methodName", 1, nil),
	}
}

func (g *GetMethodAnnotationsMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetMethodAnnotationsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	methodNameValue, exists := ctx.GetIndexValue(1)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	methodName := methodNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 查找方法
	method, exists := class.GetMethod(methodName)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 检查方法是否有注解
	if methodStmt, ok := method.(*node.ClassMethod); ok {
		if methodStmt.Annotations != nil && len(methodStmt.Annotations) > 0 {
			info := fmt.Sprintf("方法 %s::%s 的注解信息:\n", className, methodName)
			for i, annotation := range methodStmt.Annotations {
				info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
			}
			return data.NewStringValue(info), nil
		}
	}

	return data.NewStringValue(fmt.Sprintf("方法 %s::%s 没有注解", className, methodName)), nil
}

// GetPropertyAnnotationsMethod 获取属性注解信息方法
type GetPropertyAnnotationsMethod struct{}

func (g *GetPropertyAnnotationsMethod) GetName() string {
	return "getPropertyAnnotations"
}

func (g *GetPropertyAnnotationsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetPropertyAnnotationsMethod) GetIsStatic() bool {
	return false
}

func (g *GetPropertyAnnotationsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
		node.NewParameter(nil, "propertyName", 1, nil, nil),
	}
}

func (g *GetPropertyAnnotationsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
		node.NewVariable(nil, "propertyName", 1, nil),
	}
}

func (g *GetPropertyAnnotationsMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetPropertyAnnotationsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	propertyNameValue, exists := ctx.GetIndexValue(1)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	propertyName := propertyNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 查找属性
	property, exists := class.GetProperty(propertyName)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 检查属性是否有注解
	if propertyStmt, ok := property.(*node.ClassProperty); ok {
		if propertyStmt.Annotations != nil && len(propertyStmt.Annotations) > 0 {
			info := fmt.Sprintf("属性 %s::%s 的注解信息:\n", className, propertyName)
			for i, annotation := range propertyStmt.Annotations {
				info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
			}
			return data.NewStringValue(info), nil
		}
	}

	return data.NewStringValue(fmt.Sprintf("属性 %s::%s 没有注解", className, propertyName)), nil
}

// GetAllAnnotationsMethod 获取类的所有注解信息方法
type GetAllAnnotationsMethod struct{}

func (g *GetAllAnnotationsMethod) GetName() string {
	return "getAllAnnotations"
}

func (g *GetAllAnnotationsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetAllAnnotationsMethod) GetIsStatic() bool {
	return false
}

func (g *GetAllAnnotationsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
	}
}

func (g *GetAllAnnotationsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
	}
}

func (g *GetAllAnnotationsMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetAllAnnotationsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类名参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	// 检查类是否有注解
	if classStmt, ok := class.(*node.ClassStatement); ok {
		info := fmt.Sprintf("=== %s 类的完整注解信息 ===\n\n", className)

		// 类注解
		if classStmt.Annotations != nil && len(classStmt.Annotations) > 0 {
			info += "类注解:\n"
			for i, annotation := range classStmt.Annotations {
				info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
			}
			info += "\n"
		} else {
			info += "类注解: 无\n\n"
		}

		// 属性注解
		properties := class.GetPropertyList()
		if len(properties) > 0 {
			info += "属性注解:\n"
			for _, property := range properties {
				name := property.GetName()
				if propertyStmt, ok := property.(*node.ClassProperty); ok {
					if propertyStmt.Annotations != nil && len(propertyStmt.Annotations) > 0 {
						info += fmt.Sprintf("  %s:\n", name)
						for i, annotation := range propertyStmt.Annotations {
							info += fmt.Sprintf("    %d. %s\n", i+1, annotation.Class.GetName())
						}
					} else {
						info += fmt.Sprintf("  %s: 无注解\n", name)
					}
				}
			}
			info += "\n"
		}

		// 方法注解
		methods := class.GetMethods()
		if len(methods) > 0 {
			info += "方法注解:\n"
			for _, method := range methods {
				if methodStmt, ok := method.(*node.ClassMethod); ok {
					if methodStmt.Annotations != nil && len(methodStmt.Annotations) > 0 {
						info += fmt.Sprintf("  %s:\n", method.GetName())
						for i, annotation := range methodStmt.Annotations {
							info += fmt.Sprintf("    %d. %s\n", i+1, annotation.Class.GetName())
						}
					} else {
						info += fmt.Sprintf("  %s: 无注解\n", method.GetName())
					}
				}
			}
		}

		return data.NewStringValue(info), nil
	}

	return data.NewStringValue(fmt.Sprintf("类 %s 没有注解信息", className)), nil
}

// GetAnnotationDetailsMethod 获取注解详细信息方法
type GetAnnotationDetailsMethod struct{}

func (g *GetAnnotationDetailsMethod) GetName() string {
	return "getAnnotationDetails"
}

func (g *GetAnnotationDetailsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (g *GetAnnotationDetailsMethod) GetIsStatic() bool {
	return false
}

func (g *GetAnnotationDetailsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
		node.NewParameter(nil, "memberType", 1, nil, nil), // class, property, method
		node.NewParameter(nil, "memberName", 2, nil, nil),
	}
}

func (g *GetAnnotationDetailsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
		node.NewVariable(nil, "memberType", 1, nil),
		node.NewVariable(nil, "memberName", 2, nil),
	}
}

func (g *GetAnnotationDetailsMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (g *GetAnnotationDetailsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	classNameValue, exists := ctx.GetIndexValue(0)
	if !exists {
		return data.NewStringValue(""), nil
	}

	memberTypeValue, exists := ctx.GetIndexValue(1)
	if !exists {
		return data.NewStringValue(""), nil
	}

	memberNameValue, exists := ctx.GetIndexValue(2)
	if !exists {
		return data.NewStringValue(""), nil
	}

	className := classNameValue.AsString()
	memberType := memberTypeValue.AsString()
	memberName := memberNameValue.AsString()
	vm := ctx.GetVM()

	// 查找类
	class, exists := vm.GetClass(className)
	if !exists {
		return data.NewStringValue(""), nil
	}

	info := fmt.Sprintf("=== %s 的注解详细信息 ===\n\n", memberName)

	switch memberType {
	case "class":
		if classStmt, ok := class.(*node.ClassStatement); ok {
			if classStmt.Annotations != nil && len(classStmt.Annotations) > 0 {
				info += "类注解详细信息:\n"
				for i, annotation := range classStmt.Annotations {
					info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
					// 这里可以添加更多注解信息，如参数等
				}
			} else {
				info += "类没有注解\n"
			}
		}
	case "property":
		property, exists := class.GetProperty(memberName)
		if !exists {
			return data.NewStringValue(fmt.Sprintf("属性 %s 不存在", memberName)), nil
		}
		if propertyStmt, ok := property.(*node.ClassProperty); ok {
			if propertyStmt.Annotations != nil && len(propertyStmt.Annotations) > 0 {
				info += fmt.Sprintf("属性 %s 注解详细信息:\n", memberName)
				for i, annotation := range propertyStmt.Annotations {
					info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
					// 这里可以添加更多注解信息，如参数等
				}
			} else {
				info += fmt.Sprintf("属性 %s 没有注解\n", memberName)
			}
		}
	case "method":
		method, exists := class.GetMethod(memberName)
		if !exists {
			return data.NewStringValue(fmt.Sprintf("方法 %s 不存在", memberName)), nil
		}
		if methodStmt, ok := method.(*node.ClassMethod); ok {
			if methodStmt.Annotations != nil && len(methodStmt.Annotations) > 0 {
				info += fmt.Sprintf("方法 %s 注解详细信息:\n", memberName)
				for i, annotation := range methodStmt.Annotations {
					info += fmt.Sprintf("  %d. %s\n", i+1, annotation.Class.GetName())
					// 这里可以添加更多注解信息，如参数等
				}
			} else {
				info += fmt.Sprintf("方法 %s 没有注解\n", memberName)
			}
		}
	default:
		return data.NewStringValue("无效的成员类型，支持: class, property, method"), nil
	}

	return data.NewStringValue(info), nil
}
