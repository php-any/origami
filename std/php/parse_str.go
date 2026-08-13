package php

import (
	"net/url"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ParseStrFunction 实现 parse_str。
type ParseStrFunction struct{}

func NewParseStrFunction() data.FuncStmt { return &ParseStrFunction{} }

func (f *ParseStrFunction) GetName() string { return "parse_str" }

func (f *ParseStrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.NewBaseType("string")),
		node.NewParameterReference(nil, "result", 1, nil, data.NewBaseType("array")),
	}
}

func (f *ParseStrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "result", 1, data.NewBaseType("array")),
	}
}

func (f *ParseStrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	s := ""
	if strVal != nil {
		s = strVal.AsString()
	}
	values, err := url.ParseQuery(s)
	if err != nil {
		values = url.Values{}
	}
	list := make([]*data.ZVal, 0, len(values))
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		key = strings.TrimSuffix(key, "[]")
		list = append(list, data.NewNamedZVal(key, data.NewStringValue(vals[len(vals)-1])))
	}
	parsed := &data.ArrayValue{List: list}

	resultVal, _ := ctx.GetIndexValue(1)
	if dest, ok := resultVal.(*data.ArrayValue); ok {
		dest.List = parsed.List
	} else if dest, ok := resultVal.(*data.ObjectValue); ok {
		keys := make([]string, 0)
		dest.RangeProperties(func(k string, _ data.Value) bool {
			keys = append(keys, k)
			return true
		})
		for _, k := range keys {
			dest.UnsetProperty(k)
		}
		for _, zv := range parsed.List {
			if zv != nil && zv.Name != "" {
				_ = dest.SetProperty(zv.Name, zv.Value)
			}
		}
	} else if zv := ctx.GetIndexZVal(1); zv != nil {
		zv.Value = parsed
	}
	return data.NewNullValue(), nil
}
