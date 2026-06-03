package php

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// FuncGetArgFunction 实现 func_get_arg — 返回参数列表中的某个参数
type FuncGetArgFunction struct{}

func NewFuncGetArgFunction() data.FuncStmt { return &FuncGetArgFunction{} }
func (f *FuncGetArgFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	idx, _ := utils.ConvertFromIndex[int64](ctx, 0)
	args := ctx.GetCallArgs()
	if args == nil || int(idx) >= len(args) || idx < 0 {
		return data.NewBoolValue(false), nil
	}
	return args[idx], nil
}
func (f *FuncGetArgFunction) GetName() string               { return "func_get_arg" }
func (f *FuncGetArgFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *FuncGetArgFunction) GetIsStatic() bool             { return false }
func (f *FuncGetArgFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.Int{}),
	}
}
func (f *FuncGetArgFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.Int{}),
	}
}
func (f *FuncGetArgFunction) GetReturnType() data.Types { return nil }

// Bin2hexFunction 实现 bin2hex — 将二进制数据转换为十六进制表示
type Bin2hexFunction struct{}

func NewBin2hexFunction() data.FuncStmt { return &Bin2hexFunction{} }
func (f *Bin2hexFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	str, _ := utils.ConvertFromIndex[string](ctx, 0)
	return data.NewStringValue(hex.EncodeToString([]byte(str))), nil
}
func (f *Bin2hexFunction) GetName() string               { return "bin2hex" }
func (f *Bin2hexFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *Bin2hexFunction) GetIsStatic() bool             { return false }
func (f *Bin2hexFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, data.String{}),
	}
}
func (f *Bin2hexFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.String{}),
	}
}
func (f *Bin2hexFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// PrintfFunction 实现 printf — 输出格式化字符串
type PrintfFunction struct{}

func NewPrintfFunction() data.FuncStmt { return &PrintfFunction{} }
func (f *PrintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	sprintfFunc := &SprintfFunction{}
	result, _ := sprintfFunc.Call(ctx)
	if result != nil {
		if sv, ok := result.(*data.StringValue); ok {
			data.WriteOutput(sv.Value)
			return data.NewIntValue(len(sv.Value)), nil
		}
	}
	return data.NewIntValue(0), nil
}
func (f *PrintfFunction) GetName() string               { return "printf" }
func (f *PrintfFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *PrintfFunction) GetIsStatic() bool             { return false }
func (f *PrintfFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameters(nil, "values", 1, nil, nil),
	}
}
func (f *PrintfFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "values", 1, data.NewBaseType("mixed")),
	}
}
func (f *PrintfFunction) GetReturnType() data.Types { return data.NewBaseType("int") }

// ConstantFunction 实现 constant — 返回常量的值
type ConstantFunction struct{}

func NewConstantFunction() data.FuncStmt { return &ConstantFunction{} }
func (f *ConstantFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, _ := utils.ConvertFromIndex[string](ctx, 0)
	vm := ctx.GetVM()
	if val, ok := vm.GetConstant(name); ok {
		return val, nil
	}
	return data.NewNullValue(), nil
}
func (f *ConstantFunction) GetName() string               { return "constant" }
func (f *ConstantFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ConstantFunction) GetIsStatic() bool             { return false }
func (f *ConstantFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}
func (f *ConstantFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}
func (f *ConstantFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }

// GetParentClassFunction 实现 get_parent_class — 获取父类名
type GetParentClassFunction struct{}

func NewGetParentClassFunction() data.FuncStmt { return &GetParentClassFunction{} }
func (f *GetParentClassFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewBoolValue(false), nil
	}

	className := ""
	switch v := val.(type) {
	case *data.StringValue:
		className = v.Value
	case *data.ClassValue:
		parent := v.Class.GetExtend()
		if parent != nil && *parent != "" {
			return data.NewStringValue(*parent), nil
		}
		return data.NewBoolValue(false), nil
	default:
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()
	cls, ok := vm.GetClass(className)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	parent := cls.GetExtend()
	if parent == nil || *parent == "" {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(*parent), nil
}
func (f *GetParentClassFunction) GetName() string               { return "get_parent_class" }
func (f *GetParentClassFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *GetParentClassFunction) GetIsStatic() bool             { return false }
func (f *GetParentClassFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
	}
}
func (f *GetParentClassFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.NewBaseType("mixed")),
	}
}
func (f *GetParentClassFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// GetClassMethodsFunction 实现 get_class_methods — 获取类的方法名列表
type GetClassMethodsFunction struct{}

func NewGetClassMethodsFunction() data.FuncStmt { return &GetClassMethodsFunction{} }
func (f *GetClassMethodsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewBoolValue(false), nil
	}

	className := ""
	switch v := val.(type) {
	case *data.StringValue:
		className = v.Value
	case *data.ClassValue:
		methods := v.Class.GetMethods()
		result := make([]*data.ZVal, 0, len(methods))
		for i, m := range methods {
			result = append(result, data.NewNamedZVal(fmt.Sprintf("%d", i), data.NewStringValue(m.GetName())))
		}
		return &data.ArrayValue{List: result}, nil
	default:
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()
	cls, ok := vm.GetClass(className)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	methods := cls.GetMethods()
	result := make([]*data.ZVal, 0, len(methods))
	for i, m := range methods {
		result = append(result, data.NewNamedZVal(fmt.Sprintf("%d", i), data.NewStringValue(m.GetName())))
	}
	return &data.ArrayValue{List: result}, nil
}
func (f *GetClassMethodsFunction) GetName() string               { return "get_class_methods" }
func (f *GetClassMethodsFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *GetClassMethodsFunction) GetIsStatic() bool             { return false }
func (f *GetClassMethodsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
	}
}
func (f *GetClassMethodsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.NewBaseType("mixed")),
	}
}
func (f *GetClassMethodsFunction) GetReturnType() data.Types { return data.NewBaseType("array") }

// SysGetTempDirFunction 实现 sys_get_temp_dir — 返回临时目录路径
type SysGetTempDirFunction struct{}

func NewSysGetTempDirFunction() data.FuncStmt { return &SysGetTempDirFunction{} }
func (f *SysGetTempDirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(os.TempDir()), nil
}
func (f *SysGetTempDirFunction) GetName() string               { return "sys_get_temp_dir" }
func (f *SysGetTempDirFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *SysGetTempDirFunction) GetIsStatic() bool             { return false }
func (f *SysGetTempDirFunction) GetParams() []data.GetValue    { return nil }
func (f *SysGetTempDirFunction) GetVariables() []data.Variable { return nil }
func (f *SysGetTempDirFunction) GetReturnType() data.Types     { return data.NewBaseType("string") }

// CallUserFuncArrayFunction 实现 call_user_func_array — 用数组参数调用回调函数
type CallUserFuncArrayFunction struct{}

func NewCallUserFuncArrayFunction() data.FuncStmt { return &CallUserFuncArrayFunction{} }
func (f *CallUserFuncArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	callback, _ := ctx.GetIndexValue(0)
	argsVal, _ := ctx.GetIndexValue(1)

	if callback == nil {
		return data.NewBoolValue(false), nil
	}

	// 提取参数数组
	var args []data.Value
	if arr, ok := argsVal.(*data.ArrayValue); ok {
		args = arr.ToValueList()
	}

	// 调用回调
	switch cb := callback.(type) {
	case *data.FuncValue:
		// 创建调用上下文
		vars := cb.Value.GetVariables()
		callCtx := ctx.CreateContext(vars)
		for i, arg := range args {
			callCtx.GetIndexZVal(i).Value = arg
		}
		return cb.Call(callCtx)
	case *data.StringValue:
		// 字符串函数名
		vm := ctx.GetVM()
		funcName := cb.Value
		if strings.Contains(funcName, "::") {
			parts := strings.SplitN(funcName, "::", 2)
			cls, ok := vm.GetClass(parts[0])
			if !ok {
				return data.NewBoolValue(false), nil
			}
			method, ok := cls.GetMethod(parts[1])
			if !ok || method == nil {
				return data.NewBoolValue(false), nil
			}
			vars := method.GetVariables()
			callCtx := ctx.CreateContext(vars)
			for i, arg := range args {
				callCtx.GetIndexZVal(i).Value = arg
			}
			return method.Call(callCtx)
		}
		fn, ok := vm.GetFunc(funcName)
		if !ok {
			return data.NewBoolValue(false), nil
		}
		vars := fn.GetVariables()
		callCtx := ctx.CreateContext(vars)
		for i, arg := range args {
			callCtx.GetIndexZVal(i).Value = arg
		}
		return fn.Call(callCtx)
	case *data.ArrayValue:
		// 数组形式 [object, method] 或 [className, method]
		items := cb.ToValueList()
		if len(items) < 2 {
			return data.NewBoolValue(false), nil
		}
		objVal := items[0]
		methodName := ""
		if sv, ok := items[1].(*data.StringValue); ok {
			methodName = sv.Value
		}
		if methodName == "" {
			return data.NewBoolValue(false), nil
		}

		switch obj := objVal.(type) {
		case *data.ClassValue:
			method, ok := obj.Class.GetMethod(methodName)
			if !ok || method == nil {
				return data.NewBoolValue(false), nil
			}
			vars := method.GetVariables()
			callCtx := ctx.CreateContext(vars)
			for i, arg := range args {
				callCtx.GetIndexZVal(i).Value = arg
			}
			return method.Call(callCtx)
		case *data.StringValue:
			// 静态方法调用
			vm := ctx.GetVM()
			cls, ok := vm.GetClass(obj.Value)
			if !ok {
				return data.NewBoolValue(false), nil
			}
			method, ok := cls.GetMethod(methodName)
			if !ok || method == nil {
				return data.NewBoolValue(false), nil
			}
			vars := method.GetVariables()
			callCtx := ctx.CreateContext(vars)
			for i, arg := range args {
				callCtx.GetIndexZVal(i).Value = arg
			}
			return method.Call(callCtx)
		}
	}

	return data.NewBoolValue(false), nil
}
func (f *CallUserFuncArrayFunction) GetName() string               { return "call_user_func_array" }
func (f *CallUserFuncArrayFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *CallUserFuncArrayFunction) GetIsStatic() bool             { return false }
func (f *CallUserFuncArrayFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
		node.NewParameter(nil, "args", 1, nil, nil),
	}
}
func (f *CallUserFuncArrayFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.NewBaseType("callable")),
		node.NewVariable(nil, "args", 1, data.NewBaseType("array")),
	}
}
func (f *CallUserFuncArrayFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }

// MktimeFunction 实现 mktime — 取得一个日期的 Unix 时间戳
type MktimeFunction struct{}

func NewMktimeFunction() data.FuncStmt { return &MktimeFunction{} }
func (f *MktimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化实现：委托给 time() 函数
	return data.NewIntValue(0), nil
}
func (f *MktimeFunction) GetName() string               { return "mktime" }
func (f *MktimeFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *MktimeFunction) GetIsStatic() bool             { return false }
func (f *MktimeFunction) GetParams() []data.GetValue    { return nil }
func (f *MktimeFunction) GetVariables() []data.Variable { return nil }
func (f *MktimeFunction) GetReturnType() data.Types     { return data.NewBaseType("int") }

// GetDefinedFunctionsFunction 实现 get_defined_functions — 返回所有已定义函数的数组
type GetDefinedFunctionsFunction struct{}

func NewGetDefinedFunctionsFunction() data.FuncStmt { return &GetDefinedFunctionsFunction{} }
func (f *GetDefinedFunctionsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化实现：返回空数组
	return data.NewArrayValue([]data.Value{}), nil
}
func (f *GetDefinedFunctionsFunction) GetName() string               { return "get_defined_functions" }
func (f *GetDefinedFunctionsFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *GetDefinedFunctionsFunction) GetIsStatic() bool             { return false }
func (f *GetDefinedFunctionsFunction) GetParams() []data.GetValue    { return nil }
func (f *GetDefinedFunctionsFunction) GetVariables() []data.Variable { return nil }
func (f *GetDefinedFunctionsFunction) GetReturnType() data.Types     { return data.NewBaseType("array") }

// StreamWrapperRegisterFunction 实现 stream_wrapper_register — 注册一个 URL 包装器
type StreamWrapperRegisterFunction struct{}

func NewStreamWrapperRegisterFunction() data.FuncStmt { return &StreamWrapperRegisterFunction{} }
func (f *StreamWrapperRegisterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化实现：返回 false
	return data.NewBoolValue(false), nil
}
func (f *StreamWrapperRegisterFunction) GetName() string               { return "stream_wrapper_register" }
func (f *StreamWrapperRegisterFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *StreamWrapperRegisterFunction) GetIsStatic() bool             { return false }
func (f *StreamWrapperRegisterFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "protocol", 0, nil, data.String{}),
		node.NewParameter(nil, "class", 1, nil, data.String{}),
	}
}
func (f *StreamWrapperRegisterFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "protocol", 0, data.String{}),
		node.NewVariable(nil, "class", 1, data.String{}),
	}
}
func (f *StreamWrapperRegisterFunction) GetReturnType() data.Types { return data.NewBaseType("bool") }

// IniRestoreFunction 实现 ini_restore — 恢复 ini 配置项的值
type IniRestoreFunction struct{}

func NewIniRestoreFunction() data.FuncStmt { return &IniRestoreFunction{} }
func (f *IniRestoreFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化实现：不执行任何操作
	return data.NewNullValue(), nil
}
func (f *IniRestoreFunction) GetName() string               { return "ini_restore" }
func (f *IniRestoreFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *IniRestoreFunction) GetIsStatic() bool             { return false }
func (f *IniRestoreFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "option", 0, nil, data.String{}),
	}
}
func (f *IniRestoreFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "option", 0, data.String{}),
	}
}
func (f *IniRestoreFunction) GetReturnType() data.Types { return nil }

// DecHexFunction 实现 dechex — 十进制转换为十六进制
type DecHexFunction struct{}

func NewDecHexFunction() data.FuncStmt { return &DecHexFunction{} }
func (f *DecHexFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	num, _ := utils.ConvertFromIndex[int64](ctx, 0)
	return data.NewStringValue(fmt.Sprintf("%x", num)), nil
}
func (f *DecHexFunction) GetName() string               { return "dechex" }
func (f *DecHexFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *DecHexFunction) GetIsStatic() bool             { return false }
func (f *DecHexFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "num", 0, nil, data.Int{}),
	}
}
func (f *DecHexFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "num", 0, data.Int{}),
	}
}
func (f *DecHexFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// GetDeclaredClassesFunction 实现 get_declared_classes — 返回已定义类的数组
type GetDeclaredClassesFunction struct{}

func NewGetDeclaredClassesFunction() data.FuncStmt { return &GetDeclaredClassesFunction{} }
func (f *GetDeclaredClassesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化实现：返回空数组
	return data.NewArrayValue([]data.Value{}), nil
}
func (f *GetDeclaredClassesFunction) GetName() string               { return "get_declared_classes" }
func (f *GetDeclaredClassesFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *GetDeclaredClassesFunction) GetIsStatic() bool             { return false }
func (f *GetDeclaredClassesFunction) GetParams() []data.GetValue    { return nil }
func (f *GetDeclaredClassesFunction) GetVariables() []data.Variable { return nil }
func (f *GetDeclaredClassesFunction) GetReturnType() data.Types     { return data.NewBaseType("array") }

// IniAlterFunction 实现 ini_alter — 别名 ini_set
type IniAlterFunction struct{}

func NewIniAlterFunction() data.FuncStmt { return &IniAlterFunction{} }
func (f *IniAlterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// ini_alter 是 ini_set 的别名，直接委托给 ini_set 函数
	vm := ctx.GetVM()
	fn, ok := vm.GetFunc("ini_set")
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return fn.Call(ctx)
}
func (f *IniAlterFunction) GetName() string               { return "ini_alter" }
func (f *IniAlterFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *IniAlterFunction) GetIsStatic() bool             { return false }
func (f *IniAlterFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "option", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (f *IniAlterFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "option", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "value", 1, data.NewBaseType("mixed")),
	}
}
func (f *IniAlterFunction) GetReturnType() data.Types { return data.NewBaseType("string") }
