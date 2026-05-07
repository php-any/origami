package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/exception"
	"github.com/php-any/origami/std/php/array"
	"github.com/php-any/origami/std/php/attribute"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/directory"
	"github.com/php-any/origami/std/php/file"
	"github.com/php-any/origami/std/php/iconv"
	"github.com/php-any/origami/std/php/intl"
	"github.com/php-any/origami/std/php/math"
	"github.com/php-any/origami/std/php/pdo"
	"github.com/php-any/origami/std/php/preg"
	"github.com/php-any/origami/std/php/proc"
	"github.com/php-any/origami/std/php/reflection"
	"github.com/php-any/origami/std/php/stream"
)

func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewErrorReportingFunction(),
		NewSetErrorHandlerFunction(),
		NewRegisterShutdownFunctionFunction(),
		NewTimeFunction(),
		NewStrftimeFunction(),
		NewDateDefaultTimezoneGetFunction(),
		NewDateDefaultTimezoneSetFunction(),
		NewTimezoneNameFromAbbrFunction(),
		NewTimezoneNameGetFunction(),
		NewTimezoneOpenFunction(),
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
		NewInterfaceExistsFunction(),
		NewPropertyExistsFunction(),
		NewMethodExistsFunction(),
		NewClassAliasFunction(),
		NewIsAFunction(),
		NewIsSubclassOfFunction(),
		NewGetClassFunction(),
		NewGettypeFunction(),
		NewGetDebugTypeFunction(),
		NewJsonEncodeFunction(),
		NewJsonDecodeFunction(),
		NewSerializeFunction(),
		NewUnserializeFunction(),
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
		NewStrSplitFunction(),
		NewExplodeFunction(),
		NewImplodeFunction(),
		NewPackFunction(),
		NewUnpackFunction(),
		NewCountFunction(),
		NewInArrayFunction(),
		array.NewArrayKeyExistsFunction(),
		array.NewArrayRandFunction(),
		array.NewArrayKeysFunction(),
		array.NewArrayKeyFirstFunction(),
		array.NewArraySearchFunction(),
		array.NewArrayFillKeysFunction(),
		array.NewArrayIsListFunction(),
		NewMd5Function(),
		NewMd5FileFunction(),
		NewBase64EncodeFunction(),
		NewBase64DecodeFunction(),
		NewUrlencodeFunction(),
		NewUrldecodeFunction(),
		NewRawurlencodeFunction(),
		NewRawurldecodeFunction(),
		array.NewArrayMergeFunction(),
		array.NewArrayCombineFunction(),
		array.NewArrayReplaceRecursiveFunction(),
		array.NewArrayReplaceFunction(),
		array.NewArrayMergeRecursiveFunction(),
		array.NewArrayPushFunction(),
		array.NewArrayPopFunction(),
		array.NewArrayValuesFunction(),
		array.NewArrayUniqueFunction(),
		array.NewArrayIntersectFunction(),
		array.NewArrayReverseFunction(),
		array.NewSortFunction(),
		array.NewRsortFunction(),
		array.NewUsortFunction(),
		array.NewKsortFunction(),
		array.NewKrsortFunction(),
		array.NewArrayDiffUkeyFunction(),
		array.NewArrayIntersectKeyFunction(),
		array.NewArrayFlipFunction(),
		math.NewMinFunction(),
		array.NewArrayMapFunction(),
		array.NewArrayReduceFunction(),
		NewStrReplaceFunction(),
		NewStrIreplaceFunction(),
		NewStrtolowerFunction(),
		NewStrcasecmpFunction(),
		NewStrtoupperFunction(),
		NewOrdFunction(),
		NewChrFunction(),
		NewStrRepeatFunction(),
		core.NewPhpVersionFunction(),
		NewStrcspnFunction(),
		NewStrspnFunction(),
		NewStrpbrkFunction(),
		NewStrtokFunction(),
		NewStripslashesFunction(),
		NewStripsCslashesFunction(),
		NewHttpResponseCodeFunction(),
		NewHeaderFunction(),
		NewHeadersSentFunction(),
		NewMbConvertCaseFunction(),
		NewMbConvertEncodingFunction(),
		NewMbListEncodingsFunction(),
		NewMbStrtoupperFunction(),
		NewMbStrtolowerFunction(),
		NewMbStrlenFunction(),
		NewMbStrposFunction(),
		NewMbSubstrFunction(),
		NewCtypeSpaceFunction(),
		NewCtypeDigitFunction(),
		NewCtypeAlphaFunction(),
		NewCtypeAlnumFunction(),
		NewCeilFunction(),
		NewFloorFunction(),
		NewRoundFunction(),
		NewPowFunction(),
		NewRandomBytesFunction(),
		NewRandomIntFunction(),
		NewStrtotimeFunction(),
		NewGmdateFunction(),
		NewLevenshteinFunction(),
		NewMaxFunction(),
		NewNormalizerIsNormalizedFunction(),
		NewNormalizerNormalizeFunction(),
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
		core.NewHtmlspecialcharsFunction(),
		core.NewStripTagsFunction(),
		core.NewSetlocaleFunction(),
		core.NewSplObjectHashFunction(),

		core.NewSetExceptionHandlerFunction(),
		core.NewRestoreExceptionHandlerFunction(),

		NewStrrposFunction(),
		NewStrriposFunction(),
		NewStriposFunction(),
		NewPregMatchFunction(),
		core.NewIsCallableFunction(),
		NewIsStringFunction(),
		NewIsIntFunction(),
		NewIsScalarFunction(),
		NewIsArrayFunction(),
		NewIsBoolFunction(),
		NewIsFloatFunction(),
		NewIsNullFunction(),
		NewIsNumericFunction(),
		NewIsObjectFunction(),
		NewIsIterableFunction(),
		array.NewArrayShiftFunction(),
		array.NewArrayUnshiftFunction(),
		array.NewArraySliceFunction(),
		array.NewArrayDiffFunction(),
		array.NewArraySpliceFunction(),
		array.NewArrayPadFunction(),
		array.NewArrayWalkFunction(),
		NewIteratorToArrayFunction(),
		array.NewEndFunction(),
		array.NewResetFunction(),
		array.NewNextFunction(),
		array.NewPrevFunction(),
		array.NewCurrentFunction(),
		array.NewKeyFunction(),
		NewSprintfFunction(),
		NewVsprintfFunction(),
		NewVarDumpFunction(),
		NewChmodFunction(),
		NewClassImplementsFunction(),
		NewClassParentsFunction(),
		NewClassUsesFunction(),
		NewClearstatcacheFunction(),
		NewFilterVarFunction(),
		NewFuncNumArgsFunction(),
		NewGetCfgVarFunction(),
		NewParseUrlFunction(),
		NewTempnamFunction(),
		NewUmaskFunction(),
		NewVarExportFunction(),
		core.NewStreamResolveIncludePathFunction(),
		core.NewDefinedFunction(),
		core.NewDefineFunction(),
		core.NewTriggerErrorFunction(),
		core.NewHeadersSentFunction(),
		core.NewExtensionLoadedFunction(),
		core.NewExitFunction(),
		core.NewUnlinkFunction(),
		core.NewRmdirFunction(),
		core.NewMkdirFunction(),
		core.NewCopyFunction(),
		core.NewRenameFunction(),
		core.NewPutenvFunction(),
		core.NewGetenvFunction(),
		core.NewIniSetFunction(),
		core.NewIniGetFunction(),
		core.NewSapiWindowsVt100SupportFunction(),
		core.NewObStartFunction(),
		core.NewObGetCleanFunction(),
		core.NewObGetContentsFunction(),
		core.NewObEndCleanFunction(),
		core.NewObGetLevelFunction(),
		core.NewCliSetProcessTitleFunction(),
		core.NewChdirFunction(),
		file.NewFileExistsFunction(),
		file.NewIsReadableFunction(),
		file.NewIsWritableFunction(),
		file.NewFilesizeFunction(),
		file.NewFilemtimeFunction(),
		NewIsResourceFunction(),
		NewGetResourceTypeFunction(),
		proc.NewProcOpenFunction(),
		proc.NewProcCloseFunction(),
		proc.NewProcGetStatusFunction(),
		proc.NewProcTerminateFunction(),
		proc.NewShellExecFunction(),
		stream.NewFopenFunction(),
		stream.NewFcloseFunction(),
		stream.NewFwriteFunction(),
		stream.NewFflushFunction(),
		stream.NewStreamGetContentsFunction(),
		stream.NewStreamIsattyFunction(),
		NewJoinPathsFunction(),
		NewPathinfoFunction(),
		NewExtractFunction(),
	} {
		vm.AddFunc(fun)
	}

	// 初始化 pathinfo 常量
	InitPathinfoConstants(vm)

	// 注册核心类
	vm.AddClass(&core.ClosureClass{})
	vm.AddClass(&core.BackedEnumClass{})
	vm.AddClass(&core.StdClass{})
	vm.AddClass(&core.NormalizerClass{})
	vm.AddClass(&core.WeakMapClass{})

	// 注册 DOM 类
	vm.AddClass(core.NewDOMNodeClass())
	vm.AddClass(core.NewDOMDocumentClass())
	vm.AddClass(core.NewDOMElementClass())
	vm.AddClass(core.NewDOMTextClass())
	vm.AddClass(core.NewDOMCommentClass())
	vm.AddClass(core.NewDOMNodeListClass())
	vm.AddInterface(core.NewTraversableInterface())
	vm.AddInterface(core.NewIteratorAggregateInterface())
	vm.AddInterface(core.NewIteratorInterface())
	vm.AddInterface(core.NewRecursiveIteratorInterface())
	vm.AddInterface(core.NewOuterIteratorInterface())
	vm.AddClass(core.NewRecursiveDirectoryIteratorClass())
	vm.AddClass(core.NewRecursiveIteratorIteratorClass())
	vm.AddClass(&reflection.ReflectionClassClass{})
	vm.AddClass(&reflection.ReflectionMethodClass{})
	vm.AddClass(&reflection.ReflectionParameterClass{})
	vm.AddClass(&reflection.ReflectionPropertyClass{})
	vm.AddClass(&reflection.ReflectionAttributeClass{})
	vm.AddClass(&reflection.ReflectionTypeClass{})
	vm.AddClass(&reflection.ReflectionNamedTypeClass{})
	vm.AddClass(&reflection.ReflectionFunctionClass{})
	vm.AddClass(directory.NewSplFileInfoClass())
	vm.AddClass(&directory.DirectoryIteratorClass{})
	vm.AddClass(directory.NewFilesystemIteratorClass())
	vm.AddClass(&core.ArrayIteratorClass{})
	vm.AddClass(core.NewFilterIteratorClass())

	// 注册 DateTime 类
	vm.AddClass(NewDateTimeClass())

	// 注册 PHP 内置接口
	vm.AddInterface(NewArrayAccessInterface())
	vm.AddInterface(NewCountableInterface())
	vm.AddInterface(directory.NewSeekableIteratorInterface())
	vm.AddInterface(NewSerializableInterface())
	vm.AddInterface(NewSessionHandlerInterface())
	vm.AddInterface(exception.NewThrowableInterface())

	// 注册异常类
	vm.AddClass(exception.NewLogicExceptionClass())
	vm.AddClass(exception.NewInvalidArgumentExceptionClass())
	vm.AddClass(exception.NewRuntimeExceptionClass())
	vm.AddClass(exception.NewBadMethodCallExceptionClass())

	initPhpDefaultDefines(vm)

	// 加载 PDO 扩展
	pdo.Load(vm)

	// 加载 preg 包（注册函数和常量）
	preg.Load(vm)

	// 加载 iconv 系列函数
	iconv.Load(vm)

	// 加载 intl 扩展（grapheme_strlen、grapheme_substr 等）
	intl.Load(vm)

	// mb_convert_case 常量
	vm.SetConstant("MB_CASE_UPPER", data.NewIntValue(MB_CASE_UPPER))
	vm.SetConstant("MB_CASE_LOWER", data.NewIntValue(MB_CASE_LOWER))
	vm.SetConstant("MB_CASE_TITLE", data.NewIntValue(MB_CASE_TITLE))
	vm.SetConstant("MB_CASE_FOLD", data.NewIntValue(MB_CASE_FOLD))
	vm.SetConstant("MB_CASE_UPPER_SIMPLE", data.NewIntValue(MB_CASE_UPPER_SIMPLE))
	vm.SetConstant("MB_CASE_LOWER_SIMPLE", data.NewIntValue(MB_CASE_LOWER_SIMPLE))
	vm.SetConstant("MB_CASE_TITLE_SIMPLE", data.NewIntValue(MB_CASE_TITLE_SIMPLE))
	vm.SetConstant("MB_CASE_FOLD_SIMPLE", data.NewIntValue(MB_CASE_FOLD_SIMPLE))

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
	vm.SetConstant("SORT_REGULAR", data.NewIntValue(0))
	vm.SetConstant("SORT_NUMERIC", data.NewIntValue(1))
	vm.SetConstant("SORT_STRING", data.NewIntValue(2))
	vm.SetConstant("SORT_LOCALE_STRING", data.NewIntValue(3))
	vm.SetConstant("SORT_NATURAL", data.NewIntValue(5))
	vm.SetConstant("SORT_FLAG_CASE", data.NewIntValue(6-5)) // 组合时常用 SORT_NATURAL | SORT_FLAG_CASE

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

	// 布尔値常量
	vm.SetConstant("TRUE", data.NewBoolValue(true))
	vm.SetConstant("FALSE", data.NewBoolValue(false))
	vm.SetConstant("NULL", data.NewNullValue())

	// extract 相关常量（PHP extract flags）
	vm.SetConstant("EXTR_OVERWRITE", data.NewIntValue(0))        // 默认：覆盖已有变量
	vm.SetConstant("EXTR_SKIP", data.NewIntValue(1))             // 跳过已有变量（不覆盖）
	vm.SetConstant("EXTR_PREFIX_SAME", data.NewIntValue(2))      // 同名时加前缀
	vm.SetConstant("EXTR_PREFIX_ALL", data.NewIntValue(3))       // 所有变量都加前缀
	vm.SetConstant("EXTR_PREFIX_INVALID", data.NewIntValue(4))   // 非法标识符时加前缀
	vm.SetConstant("EXTR_IF_EXISTS", data.NewIntValue(6))        // 仅导入已存在的变量
	vm.SetConstant("EXTR_PREFIX_IF_EXISTS", data.NewIntValue(7)) // 已存在时加前缀导入
	vm.SetConstant("EXTR_REFS", data.NewIntValue(256))           // 以引用方式导入

	// FilesystemIterator flags 常量
	vm.SetConstant("FilesystemIterator::CURRENT_AS_PATHNAME", data.NewIntValue(directory.FSI_CURRENT_AS_PATHNAME))
	vm.SetConstant("FilesystemIterator::CURRENT_AS_FILEINFO", data.NewIntValue(directory.FSI_CURRENT_AS_FILEINFO))
	vm.SetConstant("FilesystemIterator::CURRENT_AS_SELF", data.NewIntValue(directory.FSI_CURRENT_AS_SELF))
	vm.SetConstant("FilesystemIterator::KEY_AS_PATHNAME", data.NewIntValue(directory.FSI_KEY_AS_PATHNAME))
	vm.SetConstant("FilesystemIterator::KEY_AS_FILENAME", data.NewIntValue(directory.FSI_KEY_AS_FILENAME))
	vm.SetConstant("FilesystemIterator::FOLLOW_SYMLINKS", data.NewIntValue(directory.FSI_FOLLOW_SYMLINKS))
	vm.SetConstant("FilesystemIterator::SKIP_DOTS", data.NewIntValue(directory.FSI_SKIP_DOTS))
	vm.SetConstant("FilesystemIterator::UNIX_PATHS", data.NewIntValue(directory.FSI_UNIX_PATHS))
	vm.SetConstant("FilesystemIterator::NEW_CURRENT_AND_KEY", data.NewIntValue(directory.FSI_NEW_CURRENT_AND_KEY))

	// parse_url constants
	vm.SetConstant("PHP_URL_FRAGMENT", data.NewIntValue(7))
	vm.SetConstant("PHP_URL_HOST", data.NewIntValue(1))
	vm.SetConstant("PHP_URL_PASS", data.NewIntValue(4))
	vm.SetConstant("PHP_URL_PATH", data.NewIntValue(5))
	vm.SetConstant("PHP_URL_PORT", data.NewIntValue(2))
	vm.SetConstant("PHP_URL_QUERY", data.NewIntValue(6))
	vm.SetConstant("PHP_URL_SCHEME", data.NewIntValue(0))
	vm.SetConstant("PHP_URL_USER", data.NewIntValue(3))

	// setlocale constants
	vm.SetConstant("LC_ALL", data.NewIntValue(0))
	vm.SetConstant("LC_COLLATE", data.NewIntValue(1))
	vm.SetConstant("LC_CTYPE", data.NewIntValue(2))
	vm.SetConstant("LC_MONETARY", data.NewIntValue(3))
	vm.SetConstant("LC_NUMERIC", data.NewIntValue(4))
	vm.SetConstant("LC_TIME", data.NewIntValue(5))
	vm.SetConstant("LC_MESSAGES", data.NewIntValue(6))

	// File constants
	// Filter constants
	vm.SetConstant("FILE_APPEND", data.NewIntValue(8))
	vm.SetConstant("FILE_USE_INCLUDE_PATH", data.NewIntValue(1))
	vm.SetConstant("FILTER_CALLBACK", data.NewIntValue(1024))
	vm.SetConstant("FILTER_NULL_ON_FAILURE", data.NewIntValue(134217728))
	vm.SetConstant("FILTER_REQUIRE_ARRAY", data.NewIntValue(8))
	vm.SetConstant("FILTER_VALIDATE_BOOLEAN", data.NewIntValue(258))
	vm.SetConstant("FILTER_VALIDATE_INT", data.NewIntValue(257))
}
