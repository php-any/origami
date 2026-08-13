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
func (f *GetcwdFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (f *GetcwdFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
