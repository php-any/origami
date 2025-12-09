package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/core"
)

func Load(vm data.VM) {
	core.Load(vm)
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
		NewEmptyFunction(),
		NewStrlenFunction(),
		NewStrposFunction(),
		NewSubstrFunction(),
		NewTrimFunction(),
		NewExplodeFunction(),
		NewImplodeFunction(),
		NewCountFunction(),
		NewInArrayFunction(),
		NewArrayKeyExistsFunction(),
		NewMd5Function(),
		NewBase64EncodeFunction(),
		NewBase64DecodeFunction(),
		NewUrlencodeFunction(),
		NewUrldecodeFunction(),
		NewArrayMergeFunction(),
		NewArrayPushFunction(),
		NewArrayPopFunction(),
		NewStrReplaceFunction(),
		NewStrtolowerFunction(),
		NewStrtoupperFunction(),
	} {
		vm.AddFunc(fun)
	}
}
