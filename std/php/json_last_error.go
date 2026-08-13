package php

import (
	"sync"

	"github.com/php-any/origami/data"
)

// PHP json_last_error / json_last_error_msg 错误码（与 PHP 一致）
const (
	JSON_ERROR_NONE                  = 0
	JSON_ERROR_DEPTH                 = 1
	JSON_ERROR_STATE_MISMATCH        = 2
	JSON_ERROR_CTRL_CHAR             = 3
	JSON_ERROR_SYNTAX                = 4
	JSON_ERROR_UTF8                  = 5
	JSON_ERROR_RECURSION             = 6
	JSON_ERROR_INF_OR_NAN            = 7
	JSON_ERROR_UNSUPPORTED_TYPE      = 8
	JSON_ERROR_INVALID_PROPERTY_NAME = 9
	JSON_ERROR_UTF16                 = 10
)

var (
	jsonLastErrorMu  sync.Mutex
	jsonLastError    = JSON_ERROR_NONE
	jsonLastErrorMsg = "No error"
)

func setJsonLastError(code int, msg string) {
	jsonLastErrorMu.Lock()
	defer jsonLastErrorMu.Unlock()
	jsonLastError = code
	if msg == "" {
		jsonLastErrorMsg = jsonErrorMessage(code)
	} else {
		jsonLastErrorMsg = msg
	}
}

func clearJsonLastError() {
	setJsonLastError(JSON_ERROR_NONE, "No error")
}

func getJsonLastError() (int, string) {
	jsonLastErrorMu.Lock()
	defer jsonLastErrorMu.Unlock()
	return jsonLastError, jsonLastErrorMsg
}

func jsonErrorMessage(code int) string {
	switch code {
	case JSON_ERROR_NONE:
		return "No error"
	case JSON_ERROR_DEPTH:
		return "Maximum stack depth exceeded"
	case JSON_ERROR_STATE_MISMATCH:
		return "State mismatch (invalid or malformed JSON)"
	case JSON_ERROR_CTRL_CHAR:
		return "Control character error, possibly incorrectly encoded"
	case JSON_ERROR_SYNTAX:
		return "Syntax error"
	case JSON_ERROR_UTF8:
		return "Malformed UTF-8 characters, possibly incorrectly encoded"
	case JSON_ERROR_RECURSION:
		return "Recursion detected"
	case JSON_ERROR_INF_OR_NAN:
		return "Inf and NaN cannot be JSON encoded"
	case JSON_ERROR_UNSUPPORTED_TYPE:
		return "Type is not supported"
	case JSON_ERROR_INVALID_PROPERTY_NAME:
		return "The decoded property name is invalid"
	case JSON_ERROR_UTF16:
		return "Single unpaired UTF-16 surrogate in unicode escape"
	default:
		return "Unknown error"
	}
}

// JsonLastErrorFunction 实现 json_last_error(): int
type JsonLastErrorFunction struct{}

func NewJsonLastErrorFunction() data.FuncStmt {
	return &JsonLastErrorFunction{}
}

func (f *JsonLastErrorFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	code, _ := getJsonLastError()
	return data.NewIntValue(code), nil
}

func (f *JsonLastErrorFunction) GetName() string               { return "json_last_error" }
func (f *JsonLastErrorFunction) GetParams() []data.GetValue    { return nil }
func (f *JsonLastErrorFunction) GetVariables() []data.Variable { return nil }

// JsonLastErrorMsgFunction 实现 json_last_error_msg(): string
type JsonLastErrorMsgFunction struct{}

func NewJsonLastErrorMsgFunction() data.FuncStmt {
	return &JsonLastErrorMsgFunction{}
}

func (f *JsonLastErrorMsgFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, msg := getJsonLastError()
	return data.NewStringValue(msg), nil
}

func (f *JsonLastErrorMsgFunction) GetName() string               { return "json_last_error_msg" }
func (f *JsonLastErrorMsgFunction) GetParams() []data.GetValue    { return nil }
func (f *JsonLastErrorMsgFunction) GetVariables() []data.Variable { return nil }
