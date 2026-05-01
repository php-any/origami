package php

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewVarExportFunction() data.FuncStmt {
	return &VarExportFunction{}
}

type VarExportFunction struct{}

func (f *VarExportFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	returnVal := false
	if retArg, ok := ctx.GetIndexValue(1); ok && retArg != nil {
		if b, ok := retArg.(*data.BoolValue); ok {
			returnVal = b.Value
		}
	}

	result := varExportValue(val)

	if returnVal {
		return data.NewStringValue(result), nil
	}
	// Without second arg, output directly
	fmt.Print(result)
	return data.NewNullValue(), nil
}

func varExportValue(v data.Value) string {
	if v == nil {
		return "NULL"
	}

	switch val := v.(type) {
	case *data.NullValue:
		return "NULL"
	case *data.BoolValue:
		if val.Value {
			return "true"
		}
		return "false"
	case *data.IntValue:
		if i, err := val.AsInt(); err == nil {
			return fmt.Sprintf("%d", i)
		}
		return "0"
	case *data.FloatValue:
		if f, err := val.AsFloat(); err == nil {
			return fmt.Sprintf("%.6f", f)
		}
		return "0.0"
	case *data.StringValue:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(val.Value, "'", "\\'"))
	case *data.ArrayValue:
		items := make([]string, 0, len(val.List))
		for i, z := range val.List {
			items = append(items, fmt.Sprintf("  %d => %s,", i, varExportValue(z.Value)))
		}
		return "array (\n" + strings.Join(items, "\n") + "\n)"
	case *data.ObjectValue:
		items := make([]string, 0)
		val.RangeProperties(func(key string, value data.Value) bool {
			items = append(items, fmt.Sprintf("  '%s' => %s,", key, varExportValue(value)))
			return true
		})
		return "array (\n" + strings.Join(items, "\n") + "\n)"
	case *data.ClassValue:
		return fmt.Sprintf("'%s'", val.Class.GetName())
	default:
		return fmt.Sprintf("'%s'", val.AsString())
	}
}

func (f *VarExportFunction) GetName() string {
	return "var_export"
}

func (f *VarExportFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "return", 1, data.NewBoolValue(false), nil),
	}
}

func (f *VarExportFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
		node.NewVariable(nil, "return", 1, data.Mixed{}),
	}
}
