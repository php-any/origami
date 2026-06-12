package core

import (
	"os"

	"github.com/php-any/origami/data"
)

type SysGetTempDirFunction struct{}

func NewSysGetTempDirFunction() data.FuncStmt {
	return &SysGetTempDirFunction{}
}

func (f *SysGetTempDirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if val, ok := IniGet("sys_temp_dir"); ok && val != "" {
		return data.NewStringValue(val), nil
	}
	dir := os.TempDir()
	return data.NewStringValue(dir), nil
}

func (f *SysGetTempDirFunction) GetName() string {
	return "sys_get_temp_dir"
}

func (f *SysGetTempDirFunction) GetParams() []data.GetValue {
	return nil
}

func (f *SysGetTempDirFunction) GetVariables() []data.Variable {
	return nil
}
