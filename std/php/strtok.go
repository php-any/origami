package php

import (
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type strtokState struct {
	mu     sync.Mutex
	remain string
}

var globalStrtokState strtokState

// StrtokFunction 实现 strtok 函数
// 首次调用: strtok(string, token) — 初始化并返回第一个 token
// 后续调用: strtok(token) — 继续分割前一个 string（仅传一个参数时 token 在首参）
// token 中每个字符作为独立分隔符
type StrtokFunction struct{}

func NewStrtokFunction() data.FuncStmt {
	return &StrtokFunction{}
}

func (f *StrtokFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	tokenVal, _ := ctx.GetIndexValue(1)

	if strVal == nil {
		return data.NewBoolValue(false), nil
	}

	globalStrtokState.mu.Lock()
	defer globalStrtokState.mu.Unlock()

	// 2 参数: strtok(string, token) — 重置存储的字符串
	if _, isNull := tokenVal.(*data.NullValue); !isNull {
		globalStrtokState.remain = strVal.AsString()
		result, newRemain := nextToken(globalStrtokState.remain, tokenVal.AsString())
		globalStrtokState.remain = newRemain
		if result == "" {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(result), nil
	}

	// 1 参数: strtok(token) — strVal 是 token，继续分割
	token := strVal.AsString()
	if token == "" || globalStrtokState.remain == "" {
		return data.NewBoolValue(false), nil
	}
	result, newRemain := nextToken(globalStrtokState.remain, token)
	globalStrtokState.remain = newRemain
	if result == "" {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(result), nil
}

// nextToken 跳过开头的分隔符字符，返回第一个 token 及其剩余部分。
// token 中每个字符都是独立分隔符。
func nextToken(s, token string) (string, string) {
	start := 0
	for start < len(s) {
		if strings.ContainsRune(token, rune(s[start])) {
			start++
		} else {
			break
		}
	}
	if start >= len(s) {
		return "", ""
	}
	for i := start; i < len(s); i++ {
		if strings.ContainsRune(token, rune(s[i])) {
			return s[start:i], s[i+1:]
		}
	}
	return s[start:], ""
}

func (f *StrtokFunction) GetName() string {
	return "strtok"
}

func (f *StrtokFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "token", 1, nil, nil),
	}
}

func (f *StrtokFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "token", 1, data.NewBaseType("string")),
	}
}
