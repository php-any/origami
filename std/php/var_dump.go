package php

import (
	"fmt"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewVarDumpFunction 创建 var_dump 函数
func NewVarDumpFunction() data.FuncStmt {
	return &VarDumpFunction{}
}

// VarDumpFunction 实现 PHP 风格的 var_dump
// 为了简单，当前实现主要是把参数按顺序打印出来，
// 字符串按内容打印，其它类型用 fmt.Println 输出其值。
type VarDumpFunction struct{}

func (f *VarDumpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 从调用处取输出位置（文件:行:列）
	file, line, pos := "", 0, 0
	for _, arg := range ctx.GetCallArgs() {
		if g, ok := arg.(node.GetFrom); ok {
			from := g.GetFrom()
			if from != nil {
				file = from.GetSource()
				line, pos, _, _ = from.GetRange()
				line++
				pos++
			}
			break
		}
	}
	loc := file + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos)

	for _, argument := range f.GetParams() {
		argv, _ := argument.GetValue(ctx)
		if argv == nil {
			continue
		}
		// 可变参数得到的是 ArrayValue，逐项按 PHP 风格输出，并带位置
		if arr, ok := argv.(*data.ArrayValue); ok {
			for _, zval := range arr.List {
				if zval != nil && zval.Value != nil {
					varDumpOne(loc, zval.Value)
				}
			}
			continue
		}
		if v, ok := argv.(data.Variable); ok {
			val, acl := ctx.GetVariableValue(v)
			if acl != nil {
				return nil, acl
			}
			if val != nil {
				varDumpOne(loc, val)
			}
			continue
		}
		if val, ok := argv.(data.Value); ok {
			varDumpOne(loc, val)
		}
	}
	return nil, nil
}

// varDumpOne 先输出位置，再输出单个值（PHP 风格，如 int(17)）
func varDumpOne(loc string, v data.Value) {
	fmt.Println(loc)
	switch arg := v.(type) {
	case *data.IntValue:
		fmt.Printf("int(%d)\n", arg.Value)
	case *data.FloatValue:
		fmt.Printf("float(%v)\n", arg.Value)
	case *data.BoolValue:
		if arg.Value {
			fmt.Println("bool(true)")
		} else {
			fmt.Println("bool(false)")
		}
	case *data.StringValue:
		fmt.Printf("string(%d) %q\n", len(arg.Value), arg.Value)
	case *data.NullValue:
		fmt.Println("NULL")
	default:
		if s, ok := v.(data.AsString); ok {
			fmt.Println(s.AsString())
		} else {
			fmt.Println(v)
		}
	}
}

func (f *VarDumpFunction) GetName() string {
	return "var_dump"
}

func (f *VarDumpFunction) GetParams() []data.GetValue {
	// 使用 node.Parameters 支持任意数量的参数
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

func (f *VarDumpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}
