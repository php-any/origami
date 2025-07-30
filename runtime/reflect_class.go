package runtime

import (
	"fmt"
	"reflect"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectClass 表示通过反射生成的类
type ReflectClass struct {
	name         string
	instanceType reflect.Type
	methods      map[string]data.Method
	properties   map[string]data.Property
	instance     interface{} // 新增：用于存储被代理的实例
}

// NewReflectClass 创建一个新的反射类
func NewReflectClass(name string, instance interface{}) *ReflectClass {
	instanceType := reflect.TypeOf(instance)
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	rc := &ReflectClass{
		name:         name,
		instanceType: instanceType,
		methods:      make(map[string]data.Method),
		properties:   make(map[string]data.Property),
	}

	// 只分析类型，不分析具体实例
	// rc.analyzeMethods() // 不再在类级别分析方法

	return rc
}

// analyzeMethods 分析结构体的公开方法
func (rc *ReflectClass) analyzeMethods() {
	instanceType := reflect.TypeOf(rc.instance)
	for i := 0; i < instanceType.NumMethod(); i++ {
		method := instanceType.Method(i)
		if !isPublicMethod(method.Name) {
			continue
		}
		methodWrapper := &ReflectMethod{
			name:         method.Name,
			method:       method,
			instance:     rc.instance,
			instanceType: rc.instanceType,
		}
		rc.methods[method.Name] = methodWrapper
	}
}

// GetName 返回类名
func (rc *ReflectClass) GetName() string {
	return rc.name
}

// GetExtend 返回父类名（无继承）
func (rc *ReflectClass) GetExtend() *string {
	return nil
}

// GetImplements 返回实现的接口列表
func (rc *ReflectClass) GetImplements() []string {
	return nil
}

// GetProperty 获取属性
func (rc *ReflectClass) GetProperty(name string) (data.Property, bool) {
	prop, exists := rc.properties[name]
	return prop, exists
}

// GetProperties 获取所有属性
func (rc *ReflectClass) GetProperties() map[string]data.Property {
	return rc.properties
}

// GetMethod 获取方法
func (rc *ReflectClass) GetMethod(name string) (data.Method, bool) {
	method, exists := rc.methods[name]
	return method, exists
}

// GetMethods 获取所有方法
func (rc *ReflectClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(rc.methods))
	for _, method := range rc.methods {
		methods = append(methods, method)
	}
	return methods
}

// GetConstruct 获取构造函数
func (rc *ReflectClass) GetConstruct() data.Method {
	return &ReflectConstructor{
		className:    rc.name,
		instanceType: rc.instanceType,
		instance:     rc.instance, // 传递实例
	}
}

// GetValue 获取类的值 (implements data.GetValue)
func (rc *ReflectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 每次都创建新的被代理实例
	newInstance := reflect.New(rc.instanceType).Interface()

	// 创建新的代理类，所有方法共享同一个被代理实例
	newReflectClass := &ReflectClass{
		name:         rc.name,
		instanceType: rc.instanceType,
		instance:     newInstance, // 共享的被代理实例
		methods:      make(map[string]data.Method),
		properties:   rc.properties,
	}

	// 为新的代理类分析方法
	newReflectClass.analyzeMethods()

	return data.NewClassValue(newReflectClass, ctx), nil
}

// AsString 返回类的字符串表示
func (rc *ReflectClass) AsString() string {
	return fmt.Sprintf("ReflectClass{%s}", rc.name)
}

// GetFrom 获取来源信息
func (rc *ReflectClass) GetFrom() data.From {
	return nil
}

// ReflectMethod 表示反射方法的包装器
type ReflectMethod struct {
	name         string
	method       reflect.Method
	instance     interface{}
	instanceType reflect.Type
}

// GetName 返回方法名
func (rm *ReflectMethod) GetName() string {
	return rm.name
}

// GetModifier 返回访问修饰符
func (rm *ReflectMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否是静态方法
func (rm *ReflectMethod) GetIsStatic() bool {
	return false
}

// GetParams 返回参数列表
func (rm *ReflectMethod) GetParams() []data.GetValue {
	params := make([]data.GetValue, 0)

	// 获取结构体类型信息
	instanceType := rm.instanceType
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	// 跳过第一个参数（接收者）
	for i := 1; i < rm.method.Type.NumIn(); i++ {
		paramName := fmt.Sprintf("param%d", i-1)

		// 尝试使用结构体字段名称作为参数名称
		if i-1 < instanceType.NumField() {
			field := instanceType.Field(i - 1)
			if field.PkgPath == "" { // 公开字段
				paramName = field.Name
			}
		}

		param := node.NewParameter(nil, paramName, i-1, nil, nil)
		params = append(params, param)
	}

	return params
}

// GetVariables 返回变量列表
func (rm *ReflectMethod) GetVariables() []data.Variable {
	variables := make([]data.Variable, 0)

	// 获取结构体类型信息
	instanceType := rm.instanceType
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	// 跳过第一个参数（接收者）
	for i := 1; i < rm.method.Type.NumIn(); i++ {
		paramName := fmt.Sprintf("param%d", i-1)

		// 尝试使用结构体字段名称作为参数名称
		if i-1 < instanceType.NumField() {
			field := instanceType.Field(i - 1)
			if field.PkgPath == "" { // 公开字段
				paramName = field.Name
			}
		}

		variable := node.NewVariable(nil, paramName, i-1, nil)
		variables = append(variables, variable)
	}

	return variables
}

// Call 调用方法
func (rm *ReflectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 准备参数
	args := make([]reflect.Value, 0)

	// 添加接收者 - 使用被代理实例
	args = append(args, reflect.ValueOf(rm.instance))

	// 获取方法参数
	params := rm.GetParams()
	for i, param := range params {
		// 从上下文获取参数值
		paramValue, ctl := param.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}

		// 转换为Go类型
		goValue, err := rm.convertToGoValue(paramValue, rm.method.Type.In(i+1))
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("参数类型转换失败: %v", err))
		}

		args = append(args, goValue)
	}

	// 调用Go方法
	results := rm.method.Func.Call(args)

	// 处理返回值
	if len(results) == 0 {
		return nil, nil
	}

	// 转换第一个返回值为脚本值
	if len(results) > 0 {
		scriptValue, err := rm.convertToScriptValue(results[0])
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("返回值类型转换失败: %v", err))
		}
		return scriptValue, nil
	}

	return nil, nil
}

// convertToGoValue 将脚本值转换为Go值
func (rm *ReflectMethod) convertToGoValue(scriptValue data.GetValue, goType reflect.Type) (reflect.Value, error) {
	// 获取实际值
	value, ctl := scriptValue.GetValue(nil)
	if ctl != nil {
		return reflect.Value{}, fmt.Errorf("获取值失败: %v", ctl)
	}

	switch goType.Kind() {
	case reflect.String:
		if strValue, ok := value.(data.AsString); ok {
			return reflect.ValueOf(strValue.AsString()), nil
		}
		if val, ok := value.(data.Value); ok {
			return reflect.ValueOf(val.AsString()), nil
		}
		return reflect.Value{}, fmt.Errorf("无法转换为string类型")

	case reflect.Int, reflect.Int64:
		if intValue, ok := value.(data.AsInt); ok {
			intVal, err := intValue.AsInt()
			if err != nil {
				return reflect.Value{}, fmt.Errorf("无法转换为int类型: %v", err)
			}
			return reflect.ValueOf(intVal), nil
		}
		return reflect.Value{}, fmt.Errorf("无法转换为int类型")

	case reflect.Float64:
		if floatValue, ok := value.(data.AsFloat); ok {
			floatVal, err := floatValue.AsFloat()
			if err != nil {
				return reflect.Value{}, fmt.Errorf("无法转换为float64类型: %v", err)
			}
			return reflect.ValueOf(floatVal), nil
		}
		return reflect.Value{}, fmt.Errorf("无法转换为float64类型")

	case reflect.Bool:
		if boolValue, ok := value.(data.AsBool); ok {
			boolVal, err := boolValue.AsBool()
			if err != nil {
				return reflect.Value{}, fmt.Errorf("无法转换为bool类型: %v", err)
			}
			return reflect.ValueOf(boolVal), nil
		}
		return reflect.Value{}, fmt.Errorf("无法转换为bool类型")

	default:
		return reflect.Value{}, fmt.Errorf("不支持的Go类型: %v", goType)
	}
}

// convertToScriptValue 将Go值转换为脚本值
func (rm *ReflectMethod) convertToScriptValue(goValue reflect.Value) (data.GetValue, error) {
	switch goValue.Kind() {
	case reflect.String:
		return data.NewStringValue(goValue.String()), nil

	case reflect.Int, reflect.Int64:
		return data.NewIntValue(int(goValue.Int())), nil

	case reflect.Float64:
		return data.NewFloatValue(goValue.Float()), nil

	case reflect.Bool:
		return data.NewBoolValue(goValue.Bool()), nil

	default:
		return data.NewStringValue(fmt.Sprintf("%v", goValue.Interface())), nil
	}
}

// ReflectConstructor 表示反射构造函数
type ReflectConstructor struct {
	className    string
	instanceType reflect.Type
	instance     interface{} // 新增：用于存储被代理的实例
}

// GetName 返回构造函数名
func (rc *ReflectConstructor) GetName() string {
	return "__construct"
}

// GetModifier 返回访问修饰符
func (rc *ReflectConstructor) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否是静态方法
func (rc *ReflectConstructor) GetIsStatic() bool {
	return false
}

// GetParams 返回参数列表
func (rc *ReflectConstructor) GetParams() []data.GetValue {
	// 分析结构体字段作为参数
	params := make([]data.GetValue, 0)

	// 获取结构体的字段
	instanceType := rc.instanceType
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		// 只处理公开字段
		if field.PkgPath == "" {
			paramName := field.Name
			param := node.NewParameter(nil, paramName, i, nil, nil)
			params = append(params, param)
		}
	}

	return params
}

// GetVariables 返回变量列表
func (rc *ReflectConstructor) GetVariables() []data.Variable {
	variables := make([]data.Variable, 0)

	// 获取结构体的字段
	instanceType := rc.instanceType
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		// 只处理公开字段
		if field.PkgPath == "" {
			paramName := field.Name
			variable := node.NewVariable(nil, paramName, i, nil)
			variables = append(variables, variable)
		}
	}

	return variables
}

// Call 调用构造函数
func (rc *ReflectConstructor) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数并设置字段值
	variables := rc.GetVariables()

	for _, variable := range variables {
		// 从上下文获取参数值
		paramValue, ctl := ctx.GetVariableValue(variable)
		if ctl != nil {
			// 如果获取失败，可能是参数没有传递，跳过
			continue
		}

		// 如果参数值不为空，设置到被代理实例中
		if paramValue != nil {
			// 设置字段值
			if err := setFieldValue(rc.instance, variable.GetName(), paramValue); err != nil {
				return nil, data.NewErrorThrow(nil, fmt.Errorf("设置字段值失败: %v", err))
			}
		}
	}

	// 返回 nil，表示构造函数执行完成
	return nil, nil
}

// setFieldValue 设置结构体字段的值
func setFieldValue(instance interface{}, fieldName string, value data.Value) error {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	field := instanceValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("字段 %s 不存在", fieldName)
	}

	// 根据字段类型转换值
	switch field.Kind() {
	case reflect.String:
		if strValue, ok := value.(data.AsString); ok {
			field.SetString(strValue.AsString())
		} else {
			field.SetString(value.AsString())
		}

	case reflect.Int, reflect.Int64:
		if intValue, ok := value.(data.AsInt); ok {
			intVal, err := intValue.AsInt()
			if err != nil {
				return fmt.Errorf("无法转换为int类型: %v", err)
			}
			field.SetInt(int64(intVal))
		} else {
			return fmt.Errorf("无法转换为int类型")
		}

	case reflect.Float64:
		if floatValue, ok := value.(data.AsFloat); ok {
			floatVal, err := floatValue.AsFloat()
			if err != nil {
				return fmt.Errorf("无法转换为float64类型: %v", err)
			}
			field.SetFloat(floatVal)
		} else {
			return fmt.Errorf("无法转换为float64类型")
		}

	case reflect.Bool:
		if boolValue, ok := value.(data.AsBool); ok {
			boolVal, err := boolValue.AsBool()
			if err != nil {
				return fmt.Errorf("无法转换为bool类型: %v", err)
			}
			field.SetBool(boolVal)
		} else {
			return fmt.Errorf("无法转换为bool类型")
		}

	default:
		return fmt.Errorf("不支持的字段类型: %v", field.Kind())
	}

	return nil
}

// isPublicMethod 判断方法是否为公开方法
func isPublicMethod(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}

// RegisterReflectClass 注册反射类到VM
func (vm *VM) RegisterReflectClass(name string, instance interface{}) data.Control {
	reflectClass := NewReflectClass(name, instance)
	return vm.AddClass(reflectClass)
}
