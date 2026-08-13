package php

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	varDumpObjMu   sync.Mutex
	varDumpObjIDs  = map[uintptr]int{}
	varDumpObjNext int
)

func varDumpObjectHandle(v *data.ClassValue) int {
	if v == nil {
		return 0
	}
	// 同一 PHP 对象可能在每次方法调用时包装为新的 ClassValue，用 ObjectValue 标识身份
	key := uintptr(unsafe.Pointer(v.ObjectValue))
	varDumpObjMu.Lock()
	defer varDumpObjMu.Unlock()
	if id, ok := varDumpObjIDs[key]; ok {
		return id
	}
	varDumpObjNext++
	varDumpObjIDs[key] = varDumpObjNext
	return varDumpObjNext
}

// NewVarDumpFunction 创建 var_dump 函数
func NewVarDumpFunction() data.FuncStmt {
	return &VarDumpFunction{}
}

// VarDumpFunction 实现 PHP 风格的 var_dump
type VarDumpFunction struct{}

func (f *VarDumpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	data.MarkUserOutput()
	file, line, _ := getCallLocation(ctx)
	loc := file + ":" + strconv.Itoa(line) + ":"

	for _, v := range varDumpCollectArgs(ctx) {
		fmt.Println(loc)
		varDumpValue(v, "", 0)
	}
	return nil, nil
}

// varDumpCollectArgs 收集 var_dump 实参（兼容 Parameters 打包为单个 array）
func varDumpCollectArgs(ctx data.Context) []data.Value {
	var out []data.Value
	for i := 0; ; i++ {
		val, ok := ctx.GetIndexValue(i)
		if !ok {
			break
		}
		if val == nil {
			continue
		}
		if arr, ok := val.(*data.ArrayValue); ok {
			for _, z := range arr.List {
				if z != nil && z.Value != nil {
					out = append(out, z.Value)
				}
			}
			continue
		}
		if v, ok := val.(data.Value); ok {
			out = append(out, v)
		}
	}
	return out
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

// varDumpArrayKeyLine 输出 array 元素键（稀疏整数键 Name 为 "6" 时显示 [6]=>）
func varDumpArrayKeyLine(inner string, listIndex int, zval *data.ZVal) {
	if zval != nil && zval.Name != "" {
		if n, ok := data.ParseIntArrayKeyName(zval.Name); ok {
			fmt.Printf("%s[%d]=>\n", inner, n)
			return
		}
		fmt.Printf("%s[\"%s\"]=>\n", inner, escapeSingleQuoted(zval.Name))
		return
	}
	fmt.Printf("%s[%d]=>\n", inner, listIndex)
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
	propList := arg.Class.GetPropertyList()
	props := arg.GetProperties()
	className := arg.Class.GetName()
	n := 0
	for _, p := range propList {
		if p.GetIsStatic() {
			continue
		}
		n++
	}
	fmt.Printf("%sobject(%s)#%d (%d) {\n", indent, className, varDumpObjectHandle(arg), n)
	inner := indent + "  "
	for _, p := range propList {
		if p.GetIsStatic() {
			continue
		}
		val := props[p.GetName()]
		fmt.Printf("%s%s=>\n", inner, varDumpPropertyKey(p, className))
		if val != nil {
			varDumpValue(val, inner, depth+1)
		} else {
			fmt.Printf("%sNULL\n", inner)
		}
	}
	fmt.Printf("%s}\n", indent)
}

func varDumpPropertyKey(p data.Property, className string) string {
	switch p.GetModifier() {
	case data.ModifierProtected:
		return fmt.Sprintf(`["%s":protected]`, p.GetName())
	case data.ModifierPrivate:
		return fmt.Sprintf(`["%s":"%s":private]`, p.GetName(), className)
	default:
		return fmt.Sprintf(`["%s"]`, p.GetName())
	}
}

// varDumpZVal 输出数组槽位（含 PHP 引用标记 &）
func varDumpZVal(zval *data.ZVal, indent string, depth int) {
	if zval == nil || zval.Value == nil {
		fmt.Printf("%sNULL\n", indent)
		return
	}
	if zval.RefSlotCount > 0 {
		if sv, ok := zval.Value.(*data.StringValue); ok {
			fmt.Printf("%s&string(%d) \"%s\"\n", indent, len(sv.Value), sv.Value)
			return
		}
	}
	varDumpValue(zval.Value, indent, depth)
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
		fmt.Printf("%sstring(%d) \"%s\"\n", indent, len(arg.Value), arg.Value)
	case *data.NullValue:
		fmt.Printf("%sNULL\n", indent)
	case *data.ArrayValue:
		fmt.Printf("%sarray(%d) {\n", indent, len(arg.List))
		inner := indent + "  "
		for i, zval := range arg.List {
			if zval == nil || zval.Value == nil {
				varDumpArrayKeyLine(inner, i, zval)
				fmt.Printf("%sNULL\n", inner)
				continue
			}
			varDumpArrayKeyLine(inner, i, zval)
			varDumpZVal(zval, inner, depth+1)
		}
		fmt.Printf("%s}\n", indent)
	case *data.ObjectValue:
		n := 0
		arg.RangeProperties(func(string, data.Value) bool { n++; return true })
		fmt.Printf("%sarray(%d) {\n", indent, n)
		inner := indent + "  "
		arg.RangeProperties(func(k string, val data.Value) bool {
			if idx, err := strconv.Atoi(k); err == nil {
				fmt.Printf("%s[%d]=>\n", inner, idx)
			} else {
				fmt.Printf("%s[\"%s\"]=>\n", inner, k)
			}
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
