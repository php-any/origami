package php

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewTimeFunction(),
		NewSleepFunction(),
		NewIsDirFunction(),
		NewScandirFunction(),
		NewMicrotimeFunction(),
		NewNumberFormatFunction(),
		NewFunctionExistsFunction(),
		NewGettypeFunction(),
	} {
		vm.AddFunc(fun)
	}
}
