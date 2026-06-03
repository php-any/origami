package core

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ObEndFlushFunction 实现 ob_end_flush — 刷新输出缓冲并删除顶层缓冲区
type ObEndFlushFunction struct{}

func NewObEndFlushFunction() data.FuncStmt { return &ObEndFlushFunction{} }
func (f *ObEndFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	content := obStack.buffers[len(obStack.buffers)-1].String()
	obStack.buffers = obStack.buffers[:len(obStack.buffers)-1]
	obStack.syncWriter()
	if content != "" {
		// 写入到父缓冲区（当前栈顶）
		if len(obStack.buffers) > 1 {
			obStack.buffers[len(obStack.buffers)-1].WriteString(content)
		} else {
			// 如果没有父缓冲区，写入到 stdout
			data.DefaultOutputWriter(content)
		}
	}
	return data.NewBoolValue(true), nil
}
func (f *ObEndFlushFunction) GetName() string               { return "ob_end_flush" }
func (f *ObEndFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObEndFlushFunction) GetIsStatic() bool             { return false }
func (f *ObEndFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObEndFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObEndFlushFunction) GetReturnType() data.Types     { return nil }

// ObFlushFunction 实现 ob_flush — 刷新（发送）输出缓冲区到父缓冲区
type ObFlushFunction struct{}

func NewObFlushFunction() data.FuncStmt { return &ObFlushFunction{} }
func (f *ObFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	content := obStack.buffers[len(obStack.buffers)-1].String()
	obStack.buffers[len(obStack.buffers)-1].Reset()
	if content != "" {
		// 写入到父缓冲区（倒数第二个）
		if len(obStack.buffers) > 2 {
			obStack.buffers[len(obStack.buffers)-2].WriteString(content)
		} else {
			// 如果没有父缓冲区，写入到 stdout
			data.DefaultOutputWriter(content)
		}
	}
	return data.NewBoolValue(true), nil
}
func (f *ObFlushFunction) GetName() string               { return "ob_flush" }
func (f *ObFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObFlushFunction) GetIsStatic() bool             { return false }
func (f *ObFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObFlushFunction) GetReturnType() data.Types     { return nil }

// ObCleanFunction 实现 ob_clean — 清空输出缓冲区
type ObCleanFunction struct{}

func NewObCleanFunction() data.FuncStmt { return &ObCleanFunction{} }
func (f *ObCleanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	obStack.buffers[len(obStack.buffers)-1].Reset()
	return data.NewBoolValue(true), nil
}
func (f *ObCleanFunction) GetName() string               { return "ob_clean" }
func (f *ObCleanFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObCleanFunction) GetIsStatic() bool             { return false }
func (f *ObCleanFunction) GetParams() []data.GetValue    { return nil }
func (f *ObCleanFunction) GetVariables() []data.Variable { return nil }
func (f *ObCleanFunction) GetReturnType() data.Types     { return nil }

// ObGetStatusFunction 实现 ob_get_status — 获取输出缓冲区状态
type ObGetStatusFunction struct{}

func NewObGetStatusFunction() data.FuncStmt { return &ObGetStatusFunction{} }
func (f *ObGetStatusFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	level := len(obStack.buffers) - 1
	// 返回一个关联数组（用 ObjectValue 表示）
	result := data.NewObjectValue()
	result.SetProperty("level", data.NewIntValue(level))
	result.SetProperty("type", data.NewIntValue(0))
	result.SetProperty("status", data.NewIntValue(0))
	result.SetProperty("name", data.NewStringValue("default output handler"))
	result.SetProperty("flags", data.NewIntValue(0))
	result.SetProperty("chunk_size", data.NewIntValue(0))
	if level > 0 {
		result.SetProperty("buffer_size", data.NewIntValue(obStack.buffers[level].Len()))
	} else {
		result.SetProperty("buffer_size", data.NewIntValue(0))
	}
	result.SetProperty("bytes_deleted", data.NewIntValue(0))
	return result, nil
}
func (f *ObGetStatusFunction) GetName() string               { return "ob_get_status" }
func (f *ObGetStatusFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetStatusFunction) GetIsStatic() bool             { return false }
func (f *ObGetStatusFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetStatusFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetStatusFunction) GetReturnType() data.Types     { return nil }

// ObGetLengthFunction 实现 ob_get_length — 返回输出缓冲区长度
type ObGetLengthFunction struct{}

func NewObGetLengthFunction() data.FuncStmt { return &ObGetLengthFunction{} }
func (f *ObGetLengthFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(obStack.buffers[len(obStack.buffers)-1].Len()), nil
}
func (f *ObGetLengthFunction) GetName() string               { return "ob_get_length" }
func (f *ObGetLengthFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetLengthFunction) GetIsStatic() bool             { return false }
func (f *ObGetLengthFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetLengthFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetLengthFunction) GetReturnType() data.Types     { return nil }

// ObGetFlushFunction 实现 ob_get_flush — 刷新缓冲区内容并关闭缓冲区
type ObGetFlushFunction struct{}

func NewObGetFlushFunction() data.FuncStmt { return &ObGetFlushFunction{} }
func (f *ObGetFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if len(obStack.buffers) <= 1 {
		return data.NewBoolValue(false), nil
	}
	content := obStack.buffers[len(obStack.buffers)-1].String()
	obStack.buffers = obStack.buffers[:len(obStack.buffers)-1]
	obStack.syncWriter()
	if content != "" {
		// 写入到父缓冲区（当前栈顶）
		if len(obStack.buffers) > 1 {
			obStack.buffers[len(obStack.buffers)-1].WriteString(content)
		} else {
			// 如果没有父缓冲区，写入到 stdout
			data.DefaultOutputWriter(content)
		}
	}
	return data.NewStringValue(content), nil
}
func (f *ObGetFlushFunction) GetName() string               { return "ob_get_flush" }
func (f *ObGetFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObGetFlushFunction) GetIsStatic() bool             { return false }
func (f *ObGetFlushFunction) GetParams() []data.GetValue    { return nil }
func (f *ObGetFlushFunction) GetVariables() []data.Variable { return nil }
func (f *ObGetFlushFunction) GetReturnType() data.Types     { return nil }

// ObImplicitFlushFunction 实现 ob_implicit_flush — 打开/关闭隐式刷新
type ObImplicitFlushFunction struct{}

func NewObImplicitFlushFunction() data.FuncStmt { return &ObImplicitFlushFunction{} }
func (f *ObImplicitFlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// origami 不需要实际处理隐式刷新，因为输出是直接写入的
	return data.NewNullValue(), nil
}
func (f *ObImplicitFlushFunction) GetName() string               { return "ob_implicit_flush" }
func (f *ObImplicitFlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObImplicitFlushFunction) GetIsStatic() bool             { return false }
func (f *ObImplicitFlushFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flag", 0, data.NewBoolValue(true), data.Bool{}),
	}
}
func (f *ObImplicitFlushFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flag", 0, data.Bool{}),
	}
}
func (f *ObImplicitFlushFunction) GetReturnType() data.Types { return nil }

// OutputAddRewriteVarFunction 实现 output_add_rewrite_var — 添加 URL 重写变量
type OutputAddRewriteVarFunction struct{}

func NewOutputAddRewriteVarFunction() data.FuncStmt { return &OutputAddRewriteVarFunction{} }
func (f *OutputAddRewriteVarFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// origami 不支持 URL 重写变量
	return data.NewBoolValue(false), nil
}
func (f *OutputAddRewriteVarFunction) GetName() string               { return "output_add_rewrite_var" }
func (f *OutputAddRewriteVarFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *OutputAddRewriteVarFunction) GetIsStatic() bool             { return false }
func (f *OutputAddRewriteVarFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (f *OutputAddRewriteVarFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "value", 1, data.NewBaseType("string")),
	}
}
func (f *OutputAddRewriteVarFunction) GetReturnType() data.Types { return nil }

// ObListHandlersFunction 实现 ob_list_handlers — 返回正在使用中的输出处理程序列表
type ObListHandlersFunction struct{}

func NewObListHandlersFunction() data.FuncStmt { return &ObListHandlersFunction{} }
func (f *ObListHandlersFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	level := len(obStack.buffers) - 1
	if level <= 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}
	result := make([]data.Value, level)
	for i := 0; i < level; i++ {
		result[i] = data.NewStringValue("default output handler")
	}
	return data.NewArrayValue(result), nil
}
func (f *ObListHandlersFunction) GetName() string               { return "ob_list_handlers" }
func (f *ObListHandlersFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *ObListHandlersFunction) GetIsStatic() bool             { return false }
func (f *ObListHandlersFunction) GetParams() []data.GetValue    { return nil }
func (f *ObListHandlersFunction) GetVariables() []data.Variable { return nil }
func (f *ObListHandlersFunction) GetReturnType() data.Types     { return nil }

// FlushFunction 实现 flush — 刷新输出缓冲区
type FlushFunction struct{}

func NewFlushFunction() data.FuncStmt { return &FlushFunction{} }
func (f *FlushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 在 origami 中，flush 是一个空操作
	return data.NewNullValue(), nil
}
func (f *FlushFunction) GetName() string               { return "flush" }
func (f *FlushFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *FlushFunction) GetIsStatic() bool             { return false }
func (f *FlushFunction) GetParams() []data.GetValue    { return nil }
func (f *FlushFunction) GetVariables() []data.Variable { return nil }
func (f *FlushFunction) GetReturnType() data.Types     { return nil }

// PrintRFunction 实现 print_r — 打印变量的可读信息
type PrintRFunction struct{}

func NewPrintRFunction() data.FuncStmt { return &PrintRFunction{} }
func (f *PrintRFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	ret, _ := ctx.GetIndexValue(1)

	output := printRFormat(val, "")

	if ret != nil {
		if b, ok := ret.(*data.BoolValue); ok && b.Value {
			return data.NewStringValue(output), nil
		}
	}

	data.WriteOutput(output)
	return data.NewBoolValue(true), nil
}

func printRFormat(val data.Value, indent string) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case *data.ArrayValue:
		var sb strings.Builder
		sb.WriteString(indent + "Array\n" + indent + "(\n")
		for i, zval := range v.List {
			if zval == nil || zval.Value == nil {
				continue
			}
			keyStr := fmt.Sprintf("%d", i)
			if zval.Name != "" {
				keyStr = zval.Name
			}
			sb.WriteString(indent + "    [" + keyStr + "] => " + printRFormat(zval.Value, indent+"    "))
		}
		sb.WriteString(indent + ")\n")
		return sb.String()
	case *data.ObjectValue:
		var sb strings.Builder
		sb.WriteString(indent + "Object\n" + indent + "(\n")
		v.RangeProperties(func(k string, val data.Value) bool {
			if val != nil {
				sb.WriteString(indent + "    [" + k + "] => " + printRFormat(val, indent+"    "))
			}
			return true
		})
		sb.WriteString(indent + ")\n")
		return sb.String()
	case *data.ClassValue:
		var sb strings.Builder
		className := v.Class.GetName()
		sb.WriteString(indent + className + " Object\n" + indent + "(\n")
		props := v.GetProperties()
		for _, p := range v.Class.GetPropertyList() {
			if p.GetIsStatic() {
				continue
			}
			name := p.GetName()
			mod := p.GetModifier()
			var key string
			switch mod {
			case data.ModifierPrivate:
				key = fmt.Sprintf(`"%s":"%s":private`, name, className)
			case data.ModifierProtected:
				key = fmt.Sprintf(`"%s":protected`, name)
			default:
				key = name
			}
			if val, ok := props[name]; ok && val != nil {
				sb.WriteString(indent + "    [" + key + "] => " + printRFormat(val, indent+"    "))
			}
		}
		sb.WriteString(indent + ")\n")
		return sb.String()
	case *data.StringValue:
		return v.Value + "\n"
	case *data.IntValue:
		return fmt.Sprintf("%d\n", v.Value)
	case *data.FloatValue:
		return fmt.Sprintf("%g\n", v.Value)
	case *data.BoolValue:
		if v.Value {
			return "1\n"
		}
		return ""
	case *data.NullValue:
		return "\n"
	default:
		return val.AsString() + "\n"
	}
}

func (f *PrintRFunction) GetName() string               { return "print_r" }
func (f *PrintRFunction) GetModifier() data.Modifier    { return data.ModifierPublic }
func (f *PrintRFunction) GetIsStatic() bool             { return false }
func (f *PrintRFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "return", 1, data.NewBoolValue(false), data.Bool{}),
	}
}
func (f *PrintRFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "return", 1, data.Bool{}),
	}
}
func (f *PrintRFunction) GetReturnType() data.Types { return nil }
