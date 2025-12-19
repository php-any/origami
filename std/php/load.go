package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/preg"
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
		NewClassExistsFunction(),
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
		NewRawurlencodeFunction(),
		NewRawurldecodeFunction(),
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
		core.NewStrStartsWithFunction(),
		core.NewStrEndsWithFunction(),
		core.NewStrContainsFunction(),
		core.NewArrayFilterFunction(),

		NewStrrposFunction(),
		NewStrriposFunction(),
		NewPregMatchFunction(),
		preg.NewPregMatchAllFunction(),
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
		core.NewUnlinkFunction(),
		core.NewRmdirFunction(),
		core.NewCopyFunction(),
		core.NewRenameFunction(),
	} {
		vm.AddFunc(fun)
	}

	// 注册核心类
	vm.AddClass(&core.ClosureClass{})
	vm.AddClass(&core.BackedEnumClass{})
}
