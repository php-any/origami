package fpm

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

// RequestVM is the language RequestVM; this package only adapts its SAPI output.
type RequestVM = runtime.RequestVM

func New(base *runtime.VM, output data.OutputWriter) *runtime.RequestVM {
	request := runtime.NewRequestVM(base).(*runtime.RequestVM)
	request.SetOutputWriter(output)
	return request
}
