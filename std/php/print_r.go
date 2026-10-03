package php

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PrintRFunction 实现 print_r 函数
// print_r(mixed $value, bool $return = false): string|bool
type PrintRFunction struct{}

func NewPrintRFunction() data.FuncStmt { return &PrintRFunction{} }

func (f *PrintRFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	valueV, _ := ctx.GetIndexValue(0)
	returnV, _ := ctx.GetIndexValue(1)

	if valueV == nil {
		return data.NewBoolValue(true), nil
	}

	doReturn := false
	if returnV != nil {
		if bv, ok := returnV.(*data.BoolValue); ok {
			doReturn = bv.Value
		}
	}

	str := printRValue(valueV, 0)

	if doReturn {
		return data.NewStringValue(str), nil
	}

	// 直接输出
	if c := data.EmitOutput(ctx, str); c != nil {
		return nil, c
	}
	return data.NewBoolValue(true), nil
}

// printRValue 递归生成 print_r 输出
func printRValue(v data.GetValue, indent int) string {
	pad := strings.Repeat(" ", indent)

	// 获取值
	var val data.Value
	switch vv := v.(type) {
	case data.Value:
		val = vv
	default:
		s, acl := v.(data.GetValue).GetValue(nil)
		if acl != nil {
			return ""
		}
		if sv, ok := s.(data.Value); ok {
			val = sv
		} else if str, ok := s.(data.AsString); ok {
			return str.AsString()
		} else {
			return ""
		}
	}

	switch av := val.(type) {
	case *data.ArrayValue:
		if av.Len() == 0 {
			return "Array\n(\n" + pad + ")\n"
		}
		var sb strings.Builder
		sb.WriteString("Array\n(\n")
		for arraySlots127, arrayPosition127 := av.View(), 0; arrayPosition127 < arraySlots127.Len(); arrayPosition127++ {
			item := arraySlots127.At(arrayPosition127)
			sb.WriteString(pad + "    [" + item.Name + "] => ")
			if subArr, ok := item.ReadValue().(*data.ArrayValue); ok {
				sb.WriteString("Array\n")
				sb.WriteString(pad + "    (\n")
				sb.WriteString(printRValue(subArr, indent+8))
				sb.WriteString(pad + "    )\n")
			} else {
				sb.WriteString(item.ReadValue().AsString() + "\n")
			}
		}
		sb.WriteString(pad + ")\n")
		return sb.String()

	default:
		return av.AsString() + "\n"
	}
}

func (f *PrintRFunction) GetName() string { return "print_r" }

var printRFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
	node.NewParameter(nil, "return", 1, node.NewBooleanLiteral(nil, false), data.NewBaseType("bool")),
}

func (f *PrintRFunction) GetParams() []data.GetValue {
	return printRFunctionGetParams
}

var printRFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, nil),
	node.NewVariable(nil, "return", 1, data.NewBaseType("bool")),
}

func (f *PrintRFunction) GetVariables() []data.Variable {
	return printRFunctionGetVariables
}

// PrintfFunction 实现 printf 函数
// printf(string $format, mixed ...$values): int
type PrintfFunction struct{}

func NewPrintfFunction() data.FuncStmt { return &PrintfFunction{} }

func (f *PrintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatV, _ := ctx.GetIndexValue(0)
	if formatV == nil {
		return data.NewIntValue(0), nil
	}
	format := formatV.AsString()

	// 收集参数（通过 Parameters 打包在 index 1）
	var args []data.GetValue
	valuesV, _ := ctx.GetIndexValue(1)
	if valuesV != nil {
		if paramsArray, ok := valuesV.(*data.ArrayValue); ok {
			valueList := paramsArray.ToValueList()
			for _, val := range valueList {
				args = append(args, val)
			}
		}
	}

	result := doPhpSprintf(ctx, format, args)
	if result == "" {
		return data.NewIntValue(0), nil
	}
	if c := data.EmitOutput(ctx, result); c != nil {
		return nil, c
	}
	return data.NewIntValue(len([]rune(result))), nil
}

func (f *PrintfFunction) GetName() string { return "printf" }

var printfFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "format", 0, nil, nil),
	node.NewParameters(nil, "values", 1, nil, nil),
}

func (f *PrintfFunction) GetParams() []data.GetValue {
	return printfFunctionGetParams
}

var printfFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "format", 0, nil),
	node.NewVariable(nil, "values", 1, nil),
}

func (f *PrintfFunction) GetVariables() []data.Variable {
	return printfFunctionGetVariables
}

// VprintfFunction 实现 vprintf 函数
// vprintf(string $format, array $values): int
type VprintfFunction struct{}

func NewVprintfFunction() data.FuncStmt { return &VprintfFunction{} }

func (f *VprintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatV, _ := ctx.GetIndexValue(0)
	valuesV, _ := ctx.GetIndexValue(1)

	if formatV == nil {
		return data.NewIntValue(0), nil
	}
	format := formatV.AsString()

	var args []data.GetValue
	if arr, ok := valuesV.(*data.ArrayValue); ok {
		for arraySlots128, arrayPosition128 := arr.View(), 0; arrayPosition128 < arraySlots128.Len(); arrayPosition128++ {
			item := arraySlots128.At(arrayPosition128)
			args = append(args, item.ReadValue())
		}
	}

	result := doPhpSprintf(ctx, format, args)
	if result == "" {
		return data.NewIntValue(0), nil
	}
	if c := data.EmitOutput(ctx, result); c != nil {
		return nil, c
	}
	return data.NewIntValue(len([]rune(result))), nil
}

func (f *VprintfFunction) GetName() string { return "vprintf" }

var vprintfFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "format", 0, nil, nil),
	node.NewParameter(nil, "values", 1, nil, nil),
}

func (f *VprintfFunction) GetParams() []data.GetValue {
	return vprintfFunctionGetParams
}

var vprintfFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "format", 0, nil),
	node.NewVariable(nil, "values", 1, nil),
}

func (f *VprintfFunction) GetVariables() []data.Variable {
	return vprintfFunctionGetVariables
}

// doPhpSprintf 复用 sprintf 的格式化逻辑
func doPhpSprintf(ctx data.Context, format string, args []data.GetValue) string {
	if format == "" {
		return ""
	}

	// 转换为 Go interface
	goArgs := make([]interface{}, 0, len(args))
	for _, val := range args {
		if v, ok := val.(*data.IntValue); ok {
			goArgs = append(goArgs, v.Value)
		} else if v, ok := val.(*data.FloatValue); ok {
			goArgs = append(goArgs, v.Value)
		} else if v, ok := val.(*data.StringValue); ok {
			goArgs = append(goArgs, v.Value)
		} else if v, ok := val.(*data.BoolValue); ok {
			if v.Value {
				goArgs = append(goArgs, 1)
			} else {
				goArgs = append(goArgs, 0)
			}
		} else {
			s, acl := node.ValueToDisplayString(ctx, val)
			if acl != nil {
				return ""
			}
			goArgs = append(goArgs, s)
		}
	}

	// 转换 PHP 格式串为 Go 格式
	goFormat := phpToGoFormat(format)
	goArgs = coerceArgsForFormat(goFormat, goArgs)

	if hasIncompletePrintfVerb(goFormat) {
		return ""
	}

	result := fmt.Sprintf(goFormat, goArgs...)
	if strings.Contains(result, "%!(NOVERB)") || strings.Contains(result, "%!(EXTRA") {
		return ""
	}

	return result
}

// phpPositionalRe 已在 sprintf.go 中定义
var _ = regexp.MustCompile
var _ = strconv.Itoa
