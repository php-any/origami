package php

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// json_validate 可接受的旗标（与 PHP 的 JSON_INVALID_UTF8_* 取值一致）。
const (
	JSON_INVALID_UTF8_IGNORE     = 1048576
	JSON_INVALID_UTF8_SUBSTITUTE = 2097152
)

func NewJsonValidateFunction() data.FuncStmt {
	return &JsonValidateFunction{}
}

type JsonValidateFunction struct{}

// json_validate 允许的 flags：只接受改写非法 UTF-8 行为的两面旗标，其余按 PHP 抛 ValueError。
const jsonValidateAllowedFlags = JSON_INVALID_UTF8_IGNORE | JSON_INVALID_UTF8_SUBSTITUTE

// Call 实现 json_validate(string $json, int $depth = 512, int $flags = 0): bool
//
// 只校验不建值：比 json_decode() === null 快，也比它更准（合法的 "null" / "false" 也算通过）。
func (f *JsonValidateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	jsonValue, _ := ctx.GetIndexValue(0)
	depthValue, _ := ctx.GetIndexValue(1)
	flagsValue, _ := ctx.GetIndexValue(2)

	var src string
	if jsonValue != nil {
		src = jsonValue.AsString()
	}

	depth := 512
	if depthValue != nil {
		if _, isNull := depthValue.(*data.NullValue); !isNull {
			if iv, ok := depthValue.(data.AsInt); ok {
				d, err := iv.AsInt()
				if err != nil {
					return nil, data.NewErrorThrow(nil, fmt.Errorf("json_validate(): Argument #2 ($depth) must be of type int"))
				}
				depth = d
			}
		}
	}
	if depth <= 0 {
		return nil, data.NewErrorThrowByName(nil,
			fmt.Errorf("json_validate(): Argument #2 ($depth) must be greater than 0"), "ValueError")
	}

	flags := 0
	if flagsValue != nil {
		if _, isNull := flagsValue.(*data.NullValue); !isNull {
			if iv, ok := flagsValue.(data.AsInt); ok {
				fl, err := iv.AsInt()
				if err != nil {
					return nil, data.NewErrorThrow(nil, fmt.Errorf("json_validate(): Argument #3 ($flags) must be of type int"))
				}
				flags = fl
			}
		}
	}
	if flags&^jsonValidateAllowedFlags != 0 {
		return nil, data.NewErrorThrowByName(nil,
			fmt.Errorf("json_validate(): Argument #3 ($flags) must be a valid flag (allowed flags: JSON_INVALID_UTF8_IGNORE, JSON_INVALID_UTF8_SUBSTITUTE)"), "ValueError")
	}

	// 非法 UTF-8 的两种旗标先按 PHP 语义改写，再做语法校验。
	if !utf8.ValidString(src) {
		switch {
		case flags&JSON_INVALID_UTF8_IGNORE != 0:
			src = strings.ToValidUTF8(src, "")
		case flags&JSON_INVALID_UTF8_SUBSTITUTE != 0:
			src = strings.ToValidUTF8(src, "�")
		}
	}

	if !json.Valid([]byte(src)) {
		setJsonLastError(JSON_ERROR_SYNTAX, "")
		return data.NewBoolValue(false), nil
	}
	// PHP 的 depth 含最外层容器（[] 需要 depth >= 2），标量只需要 depth >= 1。
	if jsonNestingDepth(src)+1 > depth {
		setJsonLastError(JSON_ERROR_DEPTH, "")
		return data.NewBoolValue(false), nil
	}
	clearJsonLastError()
	return data.NewBoolValue(true), nil
}

// jsonNestingDepth 返回最深的容器嵌套层数（[] 为 1，[[]] 为 2，标量为 0）。
// 纯字节扫描，跳过字符串内的括号与转义字符。
func jsonNestingDepth(src string) int {
	depth, maxDepth := 0, 0
	inString, escaped := false, false
	for i := 0; i < len(src); i++ {
		c := src[i]
		if inString {
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}
		switch c {
		case '"':
			inString = true
		case '{', '[':
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		case '}', ']':
			if depth > 0 {
				depth--
			}
		}
	}
	return maxDepth
}

func (f *JsonValidateFunction) GetName() string { return "json_validate" }

var jsonValidateFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "json", 0, nil, data.String{}),
	node.NewParameter(nil, "depth", 1, data.NewIntValue(512), nil),
	node.NewParameter(nil, "flags", 2, data.NewIntValue(0), nil),
}

func (f *JsonValidateFunction) GetParams() []data.GetValue {
	return jsonValidateFunctionGetParams
}

var jsonValidateFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "json", 0, nil),
	node.NewVariable(nil, "depth", 1, nil),
	node.NewVariable(nil, "flags", 2, nil),
}

func (f *JsonValidateFunction) GetVariables() []data.Variable {
	return jsonValidateFunctionGetVariables
}
