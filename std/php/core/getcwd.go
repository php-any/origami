package core

import (
	"os"

	"github.com/php-any/origami/data"
)

// GetcwdFunction 实现 PHP getcwd(): string|false
type GetcwdFunction struct{}

func NewGetcwdFunction() data.FuncStmt { return &GetcwdFunction{} }

func (f *GetcwdFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dir, err := os.Getwd()
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(dir), nil
}

func (f *GetcwdFunction) GetName() string { return "getcwd" }
var getcwdFunctionGetParams = []data.GetValue{}

func (f *GetcwdFunction) GetParams() []data.GetValue {
	return getcwdFunctionGetParams
}
var getcwdFunctionGetVariables = []data.Variable{}

func (f *GetcwdFunction) GetVariables() []data.Variable {
	return getcwdFunctionGetVariables
}
