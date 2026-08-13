package php

import (
	"github.com/php-any/origami/data"
)

// MbListEncodingsFunction 实现 mb_list_encodings 函数
type MbListEncodingsFunction struct{}

func NewMbListEncodingsFunction() data.FuncStmt { return &MbListEncodingsFunction{} }

func (f *MbListEncodingsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	encodings := []data.Value{
		data.NewStringValue("UTF-8"),
		data.NewStringValue("UTF-7"),
		data.NewStringValue("ASCII"),
		data.NewStringValue("ISO-8859-1"),
		data.NewStringValue("Windows-1252"),
	}
	return data.NewArrayValue(encodings), nil
}

func (f *MbListEncodingsFunction) GetName() string               { return "mb_list_encodings" }
func (f *MbListEncodingsFunction) GetParams() []data.GetValue    { return []data.GetValue{} }
func (f *MbListEncodingsFunction) GetVariables() []data.Variable { return []data.Variable{} }
