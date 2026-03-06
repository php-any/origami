package php

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewVarDumpFunction 创建 var_dump 函数
func NewVarDumpFunction() data.FuncStmt {
	return &VarDumpFunction{}
}

// VarDumpFunction 实现 PHP 风格的 var_dump
type VarDumpFunction struct{}

func (f *VarDumpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	file, line, _ := getCallLocation(ctx)
	loc := file + ":" + strconv.Itoa(line) + ":"

	for _, argument := range f.GetParams() {
		argv, _ := argument.GetValue(ctx)
		if argv == nil {
			continue
		}
		if arr, ok := argv.(*data.ArrayValue); ok {
			for _, zval := range arr.List {
				if zval != nil && zval.Value != nil {
					fmt.Println(loc)
					varDumpValue(zval.Value, "", 0)
				}
			}
			continue
		}
		if v, ok := argument.(data.Variable); ok {
			val, acl := ctx.GetVariableValue(v)
			if acl != nil {
				return nil, acl
			}
			fmt.Println(loc)
			if val != nil {
				varDumpValue(val, "", 0)
			} else {
				varDumpValue(data.NewNullValue(), "", 0)
			}
			continue
		}
		if val, ok := argv.(data.Value); ok {
			fmt.Println(loc)
			varDumpValue(val, "", 0)
		}
	}
	return nil, nil
}

func getCallLocation(ctx data.Context) (file string, line int, pos int) {
	for _, arg := range ctx.GetCallArgs() {
		if g, ok := arg.(node.GetFrom); ok {
			from := g.GetFrom()
			if from != nil {
				file = from.GetSource()
				line, pos, _, _ = from.GetRange()
				line++
				return file, line, pos
			}
			break
		}
	}
	return "", 0, 0
}

// escapeSingleQuoted 将键中的 \ 和 ' 转义，用于 PHP 风格 'key' 输出
func escapeSingleQuoted(s string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s)
}

// maxVarDumpDepth 防止循环引用导致栈溢出或长时间无输出
const maxVarDumpDepth = 5

// dumpClassValue 输出 ClassValue 的 PHP var_dump 格式，供 ClassValue/ClassMethodContext/ThisValue 复用
func dumpClassValue(arg *data.ClassValue, indent string, depth int) {
	if arg == nil {
		fmt.Printf("%sNULL\n", indent)
		return
	}
	if depth > maxVarDumpDepth {
		fmt.Printf("%s... (max depth)\n", indent)
		return
	}
	n := 0
	arg.RangeProperties(func(string, data.Value) bool { n++; return true })
	ptrAddr := fmt.Sprintf("%p", arg)
	fmt.Printf("%sclass %s#%s (%d) {\n", indent, arg.Class.GetName(), ptrAddr, n)
	inner := indent + "  "
	arg.RangeProperties(func(k string, val data.Value) bool {
		fmt.Printf("%spublic $%s =>\n", inner, k)
		if val != nil {
			varDumpValue(val, inner, depth+1)
		} else {
			fmt.Printf("%sNULL\n", inner)
		}
		return true
	})
	fmt.Printf("%s}\n", indent)
}

// varDumpValue 输出单个值的 PHP var_dump 格式；对象使用 Go 指针地址作为 ID
func varDumpValue(v data.Value, indent string, depth int) {
	switch arg := v.(type) {
	case *data.IntValue:
		fmt.Printf("%sint(%d)\n", indent, arg.Value)
	case *data.FloatValue:
		fmt.Printf("%sfloat(%v)\n", indent, arg.Value)
	case *data.BoolValue:
		if arg.Value {
			fmt.Printf("%sbool(true)\n", indent)
		} else {
			fmt.Printf("%sbool(false)\n", indent)
		}
	case *data.StringValue:
		fmt.Printf("%sstring(%d) %q\n", indent, len(arg.Value), arg.Value)
	case *data.NullValue:
		fmt.Printf("%sNULL\n", indent)
	case *data.ArrayValue:
		fmt.Printf("%sarray(%d) {\n", indent, len(arg.List))
		inner := indent + "  "
		for i, zval := range arg.List {
			if zval == nil || zval.Value == nil {
				fmt.Printf("%s[%d] =>\n%sNULL\n", inner, i, inner)
				continue
			}
			fmt.Printf("%s[%d] =>\n", inner, i)
			varDumpValue(zval.Value, inner, depth+1)
		}
		fmt.Printf("%s}\n", indent)
	case *data.ObjectValue:
		n := 0
		arg.RangeProperties(func(string, data.Value) bool { n++; return true })
		fmt.Printf("%sarray(%d) {\n", indent, n)
		inner := indent + "  "
		arg.RangeProperties(func(k string, val data.Value) bool {
			fmt.Printf("%s'%s' =>\n", inner, escapeSingleQuoted(k))
			if val != nil {
				varDumpValue(val, inner, depth+1)
			} else {
				fmt.Printf("%sNULL\n", inner)
			}
			return true
		})
		fmt.Printf("%s}\n", indent)
	case *data.ClassValue:
		dumpClassValue(arg, indent, depth)
	case *data.ClassMethodContext:
		// 类方法上下文包装了 ClassValue，按对象递归输出以保持格式一致
		if arg.ClassValue == nil {
			fmt.Printf("%sNULL\n", indent)
		} else {
			dumpClassValue(arg.ClassValue, indent, depth)
		}
	case *data.ThisValue:
		// $this 包装了 ClassValue，按对象递归输出以保持格式一致
		if arg.ClassValue == nil {
			fmt.Printf("%sNULL\n", indent)
		} else {
			dumpClassValue(arg.ClassValue, indent, depth)
		}
	default:
		if s, ok := v.(data.AsString); ok {
			out := strings.TrimSpace(s.AsString())
			if out == "" {
				out = "(empty)"
			}
			fmt.Printf("%s%s\n", indent, out)
		} else {
			out := fmt.Sprint(v)
			if out == "" {
				out = "(empty)"
			}
			fmt.Printf("%s%s\n", indent, out)
		}
	}
}

func (f *VarDumpFunction) GetName() string {
	return "var_dump"
}

func (f *VarDumpFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

func (f *VarDumpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}
