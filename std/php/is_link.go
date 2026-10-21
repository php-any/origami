package php

import (
	"github.com/php-any/origami/data"
	"os"
)

type IsLinkFunction struct{ IsFileFunction }

func NewIsLinkFunction() data.FuncStmt  { return &IsLinkFunction{} }
func (*IsLinkFunction) GetName() string { return "is_link" }
func (*IsLinkFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Lstat(value.AsString())
	return data.NewBoolValue(err == nil && info.Mode()&os.ModeSymlink != 0), nil
}
