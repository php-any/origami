package php

import (
	"os"

	"github.com/php-any/origami/data"
)

// GethostnameFunction 实现 PHP gethostname()。
type GethostnameFunction struct{}

func NewGethostnameFunction() data.FuncStmt { return &GethostnameFunction{} }

func (f *GethostnameFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hostname, err := os.Hostname()
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(hostname), nil
}

func (f *GethostnameFunction) GetName() string               { return "gethostname" }
func (f *GethostnameFunction) GetParams() []data.GetValue    { return nil }
func (f *GethostnameFunction) GetVariables() []data.Variable { return nil }
