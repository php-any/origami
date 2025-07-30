package runtime

import (
	"fmt"
	"reflect"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectFunction 表示通过反射注册的函数
type ReflectFunction struct {
	name      string
	fn        reflect.Value
	fnType    reflect.Type
	params    []data.GetValue
	variables []data.Variable
}

// NewReflectFunction 创建一个新的反射函数
func NewReflectFunction(name string, fn interface{}) *ReflectFunction {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// 分析函数参数
	params := make([]data.GetValue, 0)
	variables := make([]data.Variable, 0)

	for i := 0; i < fnType.NumIn(); i++ {
		paramName := fmt.Sprintf("param%d", i)

		// 创建参数节点
		param := node.NewParameter(nil, paramName, i, nil, nil)
		params = append(params, param)

		// 创建变量节点
		variable := node.NewVariable(nil, paramName, i, nil)
		variables = append(variables, variable)
	}

	return &ReflectFunction{
		name:      name,
		fn:        fnValue,
		fnType:    fnType,
		params:    params,
		variables: variables,
	}
}

// GetName 返回函数名
func (rf *ReflectFunction) GetName() string {
	return rf.name
}

// GetParams 返回参数列表
func (rf *ReflectFunction) GetParams() []data.GetValue {
	return rf.params
}

// GetVariables 返回变量列表
func (rf *ReflectFunction) GetVariables() []data.Variable {
	return rf.variables
}

// Call 调用反射函数
func (rf *ReflectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 准备参数
	args := make([]reflect.Value, 0)

	for i, param := range rf.params {
		// 从上下文获取参数值
		paramValue, ctl := param.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}

		// 转换为Go类型
		goValue, err := rf.convertToGoValue(paramValue, rf.fnType.In(i))
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("参数类型转换失败: %v", err))
		}

		args = append(args, goValue)
	}

	// 调用Go函数
	results := rf.fn.Call(args)

	// 处理返回值
	if len(results) == 0 {
		return nil, nil
	}

	// 转换第一个返回值为脚本值
	if len(results) > 0 {
		scriptValue, err := rf.convertToScriptValue(results[0])
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("返回值类型转换失败: %v", err))
		}
		return scriptValue, nil
	}

	return nil, nil
}

// convertToGoValue 将脚本值转换为Go值
func (rf *ReflectFunction) convertToGoValue(scriptValue data.GetValue, goType reflect.Type) (reflect.Value, error) {
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
func (rf *ReflectFunction) convertToScriptValue(goValue reflect.Value) (data.GetValue, error) {
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

// RegisterReflectFunction 注册反射函数到VM
func (vm *VM) RegisterReflectFunction(name string, fn interface{}) data.Control {
	reflectFn := NewReflectFunction(name, fn)
	return vm.AddFunc(reflectFn)
}

// RegisterFunction 注册函数的简化接口
func (vm *VM) RegisterFunction(name string, fn interface{}) data.Control {
	return vm.RegisterReflectFunction(name, fn)
}

// RegisterReflectFunctions 批量注册反射函数
func (vm *VM) RegisterReflectFunctions(functions map[string]interface{}) {
	for name, fn := range functions {
		if ctl := vm.RegisterReflectFunction(name, fn); ctl != nil {
			// 记录错误但不中断
			fmt.Printf("注册函数 %s 失败: %v\n", name, ctl)
		}
	}
}
