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
	printRObjMu  sync.Mutex
	printRObjIDs = map[uintptr]bool{}
)

func printRObjectSeen(v *data.ClassValue) bool {
	if v == nil || v.ObjectValue == nil {
		return false
	}
	key := uintptr(unsafe.Pointer(v.ObjectValue))
	printRObjMu.Lock()
	defer printRObjMu.Unlock()
	if printRObjIDs[key] {
		return true
	}
	printRObjIDs[key] = true
	return false
}

func printRResetSeen() {
	printRObjMu.Lock()
	defer printRObjMu.Unlock()
	printRObjIDs = map[uintptr]bool{}
}

type PrintRFunction struct{}

func NewPrintRFunction() data.FuncStmt {
	return &PrintRFunction{}
}

func (f *PrintRFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	data.MarkUserOutput()
	val, _ := ctx.GetIndexValue(0)
	returnVal := false
	if retArg, ok := ctx.GetIndexValue(1); ok && retArg != nil {
		if b, ok := retArg.(*data.BoolValue); ok {
			returnVal = b.Value
		}
	}

	printRResetSeen()
	result := printRValue(val, 0)

	if returnVal {
		return data.NewStringValue(result), nil
	}
	fmt.Print(result)
	return nil, nil
}

func printRValue(v data.Value, depth int) string {
	if v == nil || depth > 10 {
		return "NULL"
	}
	indent := strings.Repeat("    ", depth)

	switch arg := v.(type) {
	case *data.NullValue:
		return ""
	case *data.BoolValue:
		if arg.Value {
			return "1"
		}
		return ""
	case *data.IntValue:
		return strconv.Itoa(arg.Value)
	case *data.FloatValue:
		return fmt.Sprintf("%v", arg.Value)
	case *data.StringValue:
		return arg.Value
	case *data.ArrayValue:
		if len(arg.List) == 0 {
			return "Array\n" + indent + "(\n" + indent + ")\n"
		}
		b := &strings.Builder{}
		b.WriteString("Array\n" + indent + "(\n")
		inner := indent + "    "
		for i, zval := range arg.List {
			if zval == nil || zval.Value == nil {
				continue
			}
			name := ""
			if zval.Name != "" {
				name = zval.Name
			} else {
				name = strconv.Itoa(i)
			}
			b.WriteString(fmt.Sprintf("%s[%s] => %s\n", inner, name, printRValue(zval.Value, depth+1)))
		}
		b.WriteString(indent + ")\n")
		return b.String()
	case *data.ObjectValue:
		items := []string{}
		arg.RangeProperties(func(key string, val data.Value) bool {
			items = append(items, fmt.Sprintf("%s[%s] => %s", strings.Repeat("    ", depth+1), key, printRValue(val, depth+1)))
			return true
		})
		return "Array\n" + strings.Repeat("    ", depth) + "(\n" + strings.Join(items, "\n") + "\n" + strings.Repeat("    ", depth) + ")\n"
	case *data.ClassValue:
		className := arg.Class.GetName()
		if printRObjectSeen(arg) {
			return className + " Object\n" + indent + "(\n" + indent + "    *RECURSION*\n" + indent + ")\n"
		}
		propList := arg.Class.GetPropertyList()
		props := arg.GetProperties()
		items := []string{}
		inner := indent + "    "
		for _, p := range propList {
			if p.GetIsStatic() {
				continue
			}
			val := props[p.GetName()]
			key := printRPropertyKey(p, className)
			items = append(items, fmt.Sprintf("%s[%s] => %s", inner, key, printRValue(val, depth+1)))
		}
		return className + " Object\n" + indent + "(\n" + strings.Join(items, "\n") + "\n" + indent + ")\n"
	case *data.ClassMethodContext:
		if arg.ClassValue == nil {
			return ""
		}
		return printRValue(arg.ClassValue, depth)
	case *data.ThisValue:
		if arg.ClassValue == nil {
			return ""
		}
		return printRValue(arg.ClassValue, depth)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func printRPropertyKey(p data.Property, className string) string {
	switch p.GetModifier() {
	case data.ModifierProtected:
		return fmt.Sprintf("%s:protected", p.GetName())
	case data.ModifierPrivate:
		return fmt.Sprintf("%s:%s:private", p.GetName(), className)
	default:
		return p.GetName()
	}
}

func (f *PrintRFunction) GetName() string {
	return "print_r"
}

func (f *PrintRFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "return", 1, data.NewBoolValue(false), nil),
	}
}

func (f *PrintRFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
		node.NewVariable(nil, "return", 1, data.Mixed{}),
	}
}
