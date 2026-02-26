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
		// 本轮打印内的对象序号，每次 var_dump 调用从 1 自增，不影响并发
		objectID := uint64(0)
		if arr, ok := argv.(*data.ArrayValue); ok {
			for _, zval := range arr.List {
				if zval != nil && zval.Value != nil {
					fmt.Println(loc)
					varDumpValue(zval.Value, "", &objectID)
				}
			}
			continue
		}
		if v, ok := argument.(data.Variable); ok {
			val, acl := ctx.GetVariableValue(v)
			if acl != nil {
				return nil, acl
			}
			if val != nil {
				fmt.Println(loc)
				varDumpValue(val, "", &objectID)
			}
			continue
		}
		if val, ok := argv.(data.Value); ok {
			fmt.Println(loc)
			varDumpValue(val, "", &objectID)
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

// varDumpValue 输出单个值的 PHP var_dump 格式；objectID 仅在本轮打印内自增，用于 class Name#id
func varDumpValue(v data.Value, indent string, objectID *uint64) {
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
			varDumpValue(zval.Value, inner, objectID)
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
				varDumpValue(val, inner, objectID)
			} else {
				fmt.Printf("%sNULL\n", inner)
			}
			return true
		})
		fmt.Printf("%s}\n", indent)
	case *data.ClassValue:
		n := 0
		arg.RangeProperties(func(string, data.Value) bool { n++; return true })
		*objectID++
		fmt.Printf("%sclass %s#%d (%d) {\n", indent, arg.Class.GetName(), *objectID, n)
		inner := indent + "  "
		arg.RangeProperties(func(k string, val data.Value) bool {
			fmt.Printf("%spublic $%s =>\n", inner, k)
			if val != nil {
				varDumpValue(val, inner, objectID)
			} else {
				fmt.Printf("%sNULL\n", inner)
			}
			return true
		})
		fmt.Printf("%s}\n", indent)
	default:
		if s, ok := v.(data.AsString); ok {
			fmt.Printf("%s%s\n", indent, strings.TrimSpace(s.AsString()))
		} else {
			fmt.Printf("%s%s\n", indent, fmt.Sprint(v))
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
