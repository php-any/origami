package php

import (
	"encoding/csv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrGetcsvFunction() data.FuncStmt {
	return &StrGetcsvFunction{}
}

type StrGetcsvFunction struct{}

func (f *StrGetcsvFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	separatorValue, _ := ctx.GetIndexValue(1)
	enclosureValue, _ := ctx.GetIndexValue(2)
	escapeValue, _ := ctx.GetIndexValue(3)

	if stringValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	separator := ","
	if separatorValue != nil {
		if _, ok := separatorValue.(*data.NullValue); !ok {
			separator = separatorValue.AsString()
		}
	}
	enclosure := "\""
	if enclosureValue != nil {
		if _, ok := enclosureValue.(*data.NullValue); !ok {
			enclosure = enclosureValue.AsString()
		}
	}
	escape := "\\"
	if escapeValue != nil {
		if _, ok := escapeValue.(*data.NullValue); !ok {
			escape = escapeValue.AsString()
		}
	}

	parts, err := parseCSVLine(stringValue.AsString(), separator, enclosure, escape)
	if err != nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	values := make([]data.Value, len(parts))
	for i, part := range parts {
		values[i] = data.NewStringValue(part)
	}
	return data.NewArrayValue(values), nil
}

func parseCSVLine(line, separator, enclosure, escape string) ([]string, error) {
	if separator == "" {
		separator = ","
	}
	if enclosure == "" {
		enclosure = "\""
	}
	if escape == "" {
		escape = "\\"
	}

	r := csv.NewReader(strings.NewReader(line))
	r.Comma = []rune(separator)[0]
	r.LazyQuotes = true
	_ = enclosure
	_ = escape

	record, err := r.Read()
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (f *StrGetcsvFunction) GetName() string {
	return "str_getcsv"
}

func (f *StrGetcsvFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "separator", 1, node.NewStringLiteral(nil, ","), nil),
		node.NewParameter(nil, "enclosure", 2, node.NewStringLiteral(nil, "\""), nil),
		node.NewParameter(nil, "escape", 3, node.NewStringLiteral(nil, "\\"), nil),
	}
}

func (f *StrGetcsvFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "separator", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "enclosure", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "escape", 3, data.NewBaseType("string")),
	}
}
