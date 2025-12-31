package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/array"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/file"
	"github.com/php-any/origami/std/php/preg"
	"github.com/php-any/origami/std/php/proc"
	"github.com/php-any/origami/std/php/stream"
)

func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewTimeFunction(),
		NewSleepFunction(),
		NewIsDirFunction(),
		NewIsFileFunction(),
		NewScandirFunction(),
		NewFileGetContentsFunction(),
		NewFilePutContentsFunction(),
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
		NewSubstrCountFunction(),
		NewTrimFunction(),
		NewRtrimFunction(),
		NewExplodeFunction(),
		NewImplodeFunction(),
		NewCountFunction(),
		NewInArrayFunction(),
		array.NewArrayKeyExistsFunction(),
		NewMd5Function(),
		NewBase64EncodeFunction(),
		NewBase64DecodeFunction(),
		NewUrlencodeFunction(),
		NewUrldecodeFunction(),
		NewRawurlencodeFunction(),
		NewRawurldecodeFunction(),
		array.NewArrayMergeFunction(),
		array.NewArrayPushFunction(),
		array.NewArrayPopFunction(),
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
		array.NewArrayShiftFunction(),
		array.NewArrayUnshiftFunction(),
		array.NewArraySliceFunction(),
		array.NewEndFunction(),
		array.NewResetFunction(),
		array.NewNextFunction(),
		array.NewPrevFunction(),
		array.NewCurrentFunction(),
		array.NewKeyFunction(),
		NewSprintfFunction(),
		NewVarDumpFunction(),
		core.NewStreamResolveIncludePathFunction(),
		core.NewDefinedFunction(),
		core.NewDefineFunction(),
		core.NewTriggerErrorFunction(),
		core.NewExtensionLoadedFunction(),
		core.NewUnlinkFunction(),
		core.NewRmdirFunction(),
		core.NewCopyFunction(),
		core.NewRenameFunction(),
		core.NewFuncGetArgsFunction(),
		file.NewFileExistsFunction(),
		file.NewIsReadableFunction(),
		file.NewIsWritableFunction(),
		file.NewFilesizeFunction(),
		file.NewFilemtimeFunction(),
		NewIsResourceFunction(),
		proc.NewProcOpenFunction(),
		proc.NewProcCloseFunction(),
		proc.NewProcGetStatusFunction(),
		proc.NewProcTerminateFunction(),
		stream.NewFopenFunction(),
		stream.NewFcloseFunction(),
		stream.NewFwriteFunction(),
		stream.NewStreamGetContentsFunction(),
	} {
		vm.AddFunc(fun)
	}

	// 注册核心类
	vm.AddClass(&core.ClosureClass{})
	vm.AddClass(&core.BackedEnumClass{})

	initPhpDefaultDefines(vm)
}

func initPhpDefaultDefines(vm data.VM) {
	vm.SetConstant("DIRECTORY_SEPARATOR", data.NewStringValue("/"))
}
