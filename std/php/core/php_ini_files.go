package core

import "github.com/php-any/origami/data"

// Origami configures ini values through its host/environment, without reading
// a php.ini file or scanning PHP's additional configuration directories.
type PHPIniFilesFunction struct{ Name string }

func (f *PHPIniFilesFunction) GetName() string             { return f.Name }
func (*PHPIniFilesFunction) GetParams() []data.GetValue    { return nil }
func (*PHPIniFilesFunction) GetVariables() []data.Variable { return nil }
func (*PHPIniFilesFunction) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}
