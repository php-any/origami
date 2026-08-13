package core

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SetlocaleFunction 实现 setlocale 函数
type SetlocaleFunction struct{}

func NewSetlocaleFunction() data.FuncStmt {
	return &SetlocaleFunction{}
}

func (f *SetlocaleFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	category, _ := ctx.GetIndexValue(0)
	locale, _ := ctx.GetIndexValue(1)

	if category == nil {
		return data.NewBoolValue(false), nil
	}

	// Simple stub: just set environment variable and return the locale
	if locale != nil {
		loc := locale.AsString()
		if loc == "0" || loc == "" {
			// Return current locale
			current := os.Getenv("LC_ALL")
			if current == "" {
				current = os.Getenv("LANG")
			}
			if current == "" {
				current = "C"
			}
			return data.NewStringValue(current), nil
		}
		// Accept comma-separated locales
		parts := strings.Split(loc, ",")
		if len(parts) > 0 {
			chosen := strings.TrimSpace(parts[0])
			return data.NewStringValue(chosen), nil
		}
		return data.NewStringValue(loc), nil
	}

	return data.NewBoolValue(false), nil
}

func (f *SetlocaleFunction) GetName() string {
	return "setlocale"
}

func (f *SetlocaleFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "category", 0, nil, data.Int{}),
		node.NewParameter(nil, "locale", 1, nil, data.Mixed{}),
	}
}

func (f *SetlocaleFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "category", 0, data.Int{}),
		node.NewVariable(nil, "locale", 1, data.Mixed{}),
	}
}
