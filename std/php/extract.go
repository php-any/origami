package php

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// extract flags 常量
const (
	extrOverwrite      = 0 // 默认：覆盖已有变量
	extrSkip           = 1 // 跳过已有变量（不覆盖）
	extrPrefixSame     = 2 // 同名时加前缀
	extrPrefixAll      = 3 // 所有变量都加前缀
	extrPrefixInvalid  = 4 // 非法标识符时加前缀
	extrIfExists       = 6 // 仅导入已存在的变量
	extrPrefixIfExists = 7 // 已存在时加前缀导入
)

func NewExtractFunction() data.FuncStmt {
	return &ExtractFunction{}
}

// ExtractFunction 实现 PHP 内置函数 extract
//
//	extract(array &$array, int $flags = EXTR_OVERWRITE, string $prefix = ""): int
//
// 将数组中的键值对以变量名=>值的形式导入到当前符号表（调用者作用域）中。
// 使用 CallerContextParameter，函数直接在调用者上下文中执行。
type ExtractFunction struct{}

func (f *ExtractFunction) GetName() string { return "extract" }

func (f *ExtractFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 通过 GetCallArgs() 获取调用时传入的实际参数表达式（因为使用了 CallerContextParameter）
	callArgs := ctx.GetCallArgs()
	if len(callArgs) == 0 {
		return nil, utils.NewThrowf("extract() expects at least 1 argument, 0 given")
	}

	// 求值第一个参数（数组）
	arrayVal, acl := callArgs[0].GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if arrayVal == nil {
		return nil, utils.NewThrowf("extract(): argument #1 must not be nil")
	}

	// 获取 flags 参数（第二个，可选，默认 EXTR_OVERWRITE=0）
	flags := extrOverwrite
	if len(callArgs) >= 2 {
		flagsVal, acl2 := callArgs[1].GetValue(ctx)
		if acl2 != nil {
			return nil, acl2
		}
		iv, ok := flagsVal.(*data.IntValue)
		if !ok {
			return nil, utils.NewThrowf("extract(): argument #2 (flags) must be int, got %T", flagsVal)
		}
		n, err := iv.AsInt()
		if err != nil {
			return nil, utils.NewThrowf("extract(): argument #2 (flags) invalid: %v", err)
		}
		flags = n
	}

	// 获取 prefix 参数（第三个，可选，默认空字符串）
	prefix := ""
	if len(callArgs) >= 3 {
		prefixVal, acl3 := callArgs[2].GetValue(ctx)
		if acl3 != nil {
			return nil, acl3
		}
		pv, ok := prefixVal.(data.Value)
		if !ok {
			return nil, utils.NewThrowf("extract(): argument #3 (prefix) must be a value, got %T", prefixVal)
		}
		prefix = pv.AsString()
		if prefix != "" {
			prefix += "_"
		}
	}

	// 需要 prefix 的 flags 要求 prefix 非空
	if (flags == extrPrefixAll || flags == extrPrefixSame || flags == extrPrefixIfExists) && prefix == "" {
		return nil, utils.NewThrowf("extract(): flags EXTR_PREFIX_* requires a non-empty prefix (argument #3)")
	}

	count := 0

	doExtract := func(key string, val data.Value) data.Control {
		var varName string
		switch flags {
		case extrSkip:
			if ctx.HasVariableByName(key) {
				return nil
			}
			varName = key
		case extrPrefixAll:
			varName = prefix + key
		case extrPrefixSame:
			if ctx.HasVariableByName(key) {
				varName = prefix + key
			} else {
				varName = key
			}
		case extrIfExists:
			if !ctx.HasVariableByName(key) {
				return nil
			}
			varName = key
		case extrPrefixIfExists:
			if !ctx.HasVariableByName(key) {
				return nil
			}
			varName = prefix + key
		case extrOverwrite:
			varName = key
		default:
			return utils.NewThrowf("extract(): unsupported flags value %d", flags)
		}
		if !isValidIdentifier(varName) {
			// 非法标识符静默跳过（PHP 原生行为）
			return nil
		}
		ctx.SetVariableByName(varName, val)
		count++
		return nil
	}

	switch v := arrayVal.(type) {
	case *data.ObjectValue:
		// 关联数组由 ObjectValue 表示（PHP 中 ['key' => 'val'] 即此类型）
		v.RangeProperties(func(key string, val data.Value) bool {
			acl = doExtract(key, val)
			return acl == nil
		})
		if acl != nil {
			return nil, acl
		}
	case *data.ArrayValue:
		// 顺序数组：仅处理带有字符串键名（ZVal.Name != ""）的条目
		for _, zv := range v.List {
			if zv == nil || zv.Name == "" {
				continue
			}
			if acl = doExtract(zv.Name, zv.Value); acl != nil {
				return nil, acl
			}
		}
	default:
		return nil, utils.NewThrowf("extract(): argument #1 must be an array, got %T", arrayVal)
	}

	return data.NewIntValue(count), nil
}

// isValidIdentifier 检查字符串是否是合法的 PHP 变量名标识符
// PHP 变量名必须以字母或下划线开头
func isValidIdentifier(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_') {
		return false
	}
	for i := 1; i < len(name); i++ {
		c = name[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

func (f *ExtractFunction) GetParams() []data.GetValue {
	// 使用 CallerContextParameter：调用时 fnCtx = ctx（调用者上下文），
	// 使得 Call() 中可以直接通过 SetVariableByName 修改调用者的符号表。
	return []data.GetValue{
		node.NewCallerContextParameter(nil),
	}
}

func (f *ExtractFunction) GetVariables() []data.Variable {
	// CallerContextParameter 模式下无需额外变量槽
	return []data.Variable{}
}

// 确保 fmt 包被引用（用于错误信息格式化）
var _ = fmt.Sprintf
