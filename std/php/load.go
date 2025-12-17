package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/core"
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
		core.NewSplAutoloadRegisterFunction(),
		core.NewSplAutoloadUnregisterFunction(),
		core.NewArrayFunction(),
		core.NewDirnameFunction(),
		core.NewCallUserFuncFunction(),
		core.NewStrtrFunction(),

		NewStrrposFunction(),
		NewStrriposFunction(),
		NewPregMatchFunction(),
		core.NewIsCallableFunction(),
		NewIsStringFunction(),
		NewIsIntFunction(),
		NewIsArrayFunction(),
		NewIsBoolFunction(),
		NewIsFloatFunction(),
		NewIsNullFunction(),
		NewIsNumericFunction(),
		NewIsObjectFunction(),
		NewArrayShiftFunction(),
		NewArrayUnshiftFunction(),
		NewSprintfFunction(),
		core.NewStreamResolveIncludePathFunction(),
		core.NewDefinedFunction(),
		core.NewDefineFunction(),
		core.NewTriggerErrorFunction(),
		core.NewExtensionLoadedFunction(),
	} {
		vm.AddFunc(fun)
	}

	// 注册核心类
	vm.AddClass(&core.ClosureClass{})
}
