package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/exception"
	"github.com/php-any/origami/std/php/array"
	"github.com/php-any/origami/std/php/attribute"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/directory"
	"github.com/php-any/origami/std/php/file"
	"github.com/php-any/origami/std/php/preg"
	"github.com/php-any/origami/std/php/proc"
	"github.com/php-any/origami/std/php/reflection"
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
		NewPropertyExistsFunction(),
		NewGetClassFunction(),
		NewGettypeFunction(),
		NewGetDebugTypeFunction(),
		NewJsonEncodeFunction(),
		NewJsonDecodeFunction(),
		NewEmptyFunction(),
		NewStrlenFunction(),
		NewStrposFunction(),
		NewSubstrFunction(),
		NewSubstrCountFunction(),
		NewTrimFunction(),
		NewLtrimFunction(),
		NewRtrimFunction(),
		NewUcfirstFunction(),
		NewLcfirstFunction(),
		NewUcwordsFunction(),
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
		array.NewArrayValuesFunction(),
		array.NewSortFunction(),
		NewStrReplaceFunction(),
		NewStrtolowerFunction(),
		NewStrtoupperFunction(),
		core.NewSplAutoloadRegisterFunction(),
		core.NewSplAutoloadUnregisterFunction(),
		core.NewArrayFunction(),
		core.NewDirnameFunction(),
		core.NewBasenameFunction(),
		core.NewRealpathFunction(),
		core.NewCallUserFuncFunction(),
		core.NewStrtrFunction(),
		core.NewStrStartsWithFunction(),
		core.NewStrEndsWithFunction(),
		core.NewStrContainsFunction(),
		core.NewArrayFilterFunction(),

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
		NewVsprintfFunction(),
		NewVarDumpFunction(),
		core.NewStreamResolveIncludePathFunction(),
		core.NewDefinedFunction(),
		core.NewDefineFunction(),
		core.NewTriggerErrorFunction(),
		core.NewHeadersSentFunction(),
		core.NewExtensionLoadedFunction(),
		core.NewUnlinkFunction(),
		core.NewRmdirFunction(),
		core.NewCopyFunction(),
		core.NewRenameFunction(),
		core.NewFuncGetArgsFunction(),
		core.NewPutenvFunction(),
		core.NewCliSetProcessTitleFunction(),
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
		stream.NewStreamIsattyFunction(),
	} {
		vm.AddFunc(fun)
	}

	// 注册核心类
	vm.AddClass(&core.ClosureClass{})
	vm.AddClass(&core.BackedEnumClass{})
	vm.AddClass(&reflection.ReflectionClassClass{})
	vm.AddClass(&reflection.ReflectionMethodClass{})
	vm.AddClass(&reflection.ReflectionParameterClass{})
	vm.AddClass(&reflection.ReflectionAttributeClass{})
	vm.AddClass(&reflection.ReflectionTypeClass{})
	vm.AddClass(&reflection.ReflectionNamedTypeClass{})
	vm.AddClass(&directory.DirectoryIteratorClass{})

	// 注册异常类
	vm.AddClass(exception.NewLogicExceptionClass())
	vm.AddClass(exception.NewInvalidArgumentExceptionClass())

	initPhpDefaultDefines(vm)

	// 加载 preg 包（注册函数和常量）
	preg.Load(vm)

	// 加载 PHP 原生注解类
	attribute.Load(vm)
}

func initPhpDefaultDefines(vm data.VM) {
	// 目录和路径相关常量
	vm.SetConstant("DIRECTORY_SEPARATOR", data.NewStringValue("/"))
	vm.SetConstant("PATH_SEPARATOR", data.NewStringValue(":"))

	// 数组相关常量
	vm.SetConstant("ARRAY_FILTER_USE_KEY", data.NewIntValue(1))
	vm.SetConstant("ARRAY_FILTER_USE_BOTH", data.NewIntValue(2))

	// 错误级别常量
	vm.SetConstant("E_ERROR", data.NewIntValue(1))
	vm.SetConstant("E_WARNING", data.NewIntValue(2))
	vm.SetConstant("E_PARSE", data.NewIntValue(4))
	vm.SetConstant("E_NOTICE", data.NewIntValue(8))
	vm.SetConstant("E_CORE_ERROR", data.NewIntValue(16))
	vm.SetConstant("E_CORE_WARNING", data.NewIntValue(32))
	vm.SetConstant("E_COMPILE_ERROR", data.NewIntValue(64))
	vm.SetConstant("E_COMPILE_WARNING", data.NewIntValue(128))
	vm.SetConstant("E_USER_ERROR", data.NewIntValue(256))
	vm.SetConstant("E_USER_WARNING", data.NewIntValue(512))
	vm.SetConstant("E_USER_NOTICE", data.NewIntValue(1024))
	vm.SetConstant("E_STRICT", data.NewIntValue(2048))
	vm.SetConstant("E_RECOVERABLE_ERROR", data.NewIntValue(4096))
	vm.SetConstant("E_DEPRECATED", data.NewIntValue(8192))
	vm.SetConstant("E_USER_DEPRECATED", data.NewIntValue(16384))
	vm.SetConstant("E_ALL", data.NewIntValue(32767))

	// PHP 版本和系统信息常量
	vm.SetConstant("PHP_VERSION", data.NewStringValue("8.0.0"))
	vm.SetConstant("PHP_MAJOR_VERSION", data.NewIntValue(8))
	vm.SetConstant("PHP_MINOR_VERSION", data.NewIntValue(0))
	vm.SetConstant("PHP_RELEASE_VERSION", data.NewIntValue(0))
	vm.SetConstant("PHP_VERSION_ID", data.NewIntValue(80225))
	vm.SetConstant("PHP_EXTRA_VERSION", data.NewStringValue(""))
	vm.SetConstant("PHP_OS", data.NewStringValue("Linux"))
	vm.SetConstant("PHP_OS_FAMILY", data.NewStringValue("Linux"))
	vm.SetConstant("PHP_SAPI", data.NewStringValue("cli"))
	vm.SetConstant("PHP_EOL", data.NewStringValue("\n"))

	// 整数相关常量
	vm.SetConstant("PHP_INT_MAX", data.NewIntValue(9223372036854775807))
	vm.SetConstant("PHP_INT_MIN", data.NewIntValue(-9223372036854775808))
	vm.SetConstant("PHP_INT_SIZE", data.NewIntValue(8))

	// 浮点数相关常量
	vm.SetConstant("PHP_FLOAT_MAX", data.NewFloatValue(1.7976931348623157e+308))
	vm.SetConstant("PHP_FLOAT_MIN", data.NewFloatValue(2.2250738585072014e-308))
	vm.SetConstant("PHP_FLOAT_DIG", data.NewIntValue(15))
	vm.SetConstant("PHP_FLOAT_EPSILON", data.NewFloatValue(2.220446049250313e-16))

	// 数学常量
	vm.SetConstant("M_PI", data.NewFloatValue(3.14159265358979323846))
	vm.SetConstant("M_E", data.NewFloatValue(2.7182818284590452354))
	vm.SetConstant("M_LOG2E", data.NewFloatValue(1.4426950408889634074))
	vm.SetConstant("M_LOG10E", data.NewFloatValue(0.43429448190325182765))
	vm.SetConstant("M_LN2", data.NewFloatValue(0.69314718055994530942))
	vm.SetConstant("M_LN10", data.NewFloatValue(2.30258509299404568402))
	vm.SetConstant("M_PI_2", data.NewFloatValue(1.57079632679489661923))
	vm.SetConstant("M_PI_4", data.NewFloatValue(0.78539816339744830962))
	vm.SetConstant("M_1_PI", data.NewFloatValue(0.31830988618379067154))
	vm.SetConstant("M_2_PI", data.NewFloatValue(0.63661977236758134308))
	vm.SetConstant("M_SQRTPI", data.NewFloatValue(1.77245385090551602729))
	vm.SetConstant("M_2_SQRTPI", data.NewFloatValue(1.12837916709551257390))
	vm.SetConstant("M_SQRT2", data.NewFloatValue(1.41421356237309504880))
	vm.SetConstant("M_SQRT3", data.NewFloatValue(1.73205080756887729353))
	vm.SetConstant("M_SQRT1_2", data.NewFloatValue(0.70710678118654752440))
	vm.SetConstant("M_LNPI", data.NewFloatValue(1.14472988584940017414))
	vm.SetConstant("M_EULER", data.NewFloatValue(0.57721566490153286061))

	// 布尔值常量
	vm.SetConstant("TRUE", data.NewBoolValue(true))
	vm.SetConstant("FALSE", data.NewBoolValue(false))
	vm.SetConstant("NULL", data.NewNullValue())
}
