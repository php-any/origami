package php

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewTimeFunction(),
		NewSleepFunction(),
		NewIsDirFunction(),
		NewIsFileFunction(),
		NewScandirFunction(),
		NewFileGetContentsFunction(),
		NewMicrotimeFunction(),
		NewNumberFormatFunction(),
		NewFunctionExistsFunction(),
		NewGettypeFunction(),
		NewJsonEncodeFunction(),
		NewJsonDecodeFunction(),
		NewIssetFunction(),
	} {
		vm.AddFunc(fun)
	}
}
