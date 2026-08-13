package php

import (
	"time"

	"github.com/ncruces/go-strftime"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewStrftimeFunction 创建 strftime 函数。
// PHP 语义：
//
//	strftime(string $format, ?int $timestamp = null): string|false
//
// 按本地时区格式化时间。format 为空时返回 false；timestamp 为 null 或省略时使用当前时间。
// 使用 time.Local 作为默认时区。
func NewStrftimeFunction() data.FuncStmt {
	return &StrftimeFunction{}
}

type StrftimeFunction struct {
	data.Function
}

func (f *StrftimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatVar := node.NewVariable(nil, "format", 0, nil)
	formatVal, ctl := ctx.GetVariableValue(formatVar)
	if ctl != nil {
		return nil, ctl
	}
	if formatVal == nil {
		return data.NewBoolValue(false), nil
	}
	asStr, ok := formatVal.(data.AsString)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	format := asStr.AsString()
	if format == "" {
		return data.NewBoolValue(false), nil
	}

	// 可选 timestamp：null 或省略时使用当前时间
	loc := time.Local
	if loc == nil {
		loc = time.UTC
	}
	var t time.Time
	timestampVar := node.NewVariable(nil, "timestamp", 1, nil)
	timestampVal, _ := ctx.GetVariableValue(timestampVar)
	if timestampVal == nil {
		t = time.Now().In(loc)
	} else if _, isNull := timestampVal.(*data.NullValue); isNull {
		t = time.Now().In(loc)
	} else if asInt, ok := timestampVal.(data.AsInt); ok {
		unix, _ := asInt.AsInt()
		t = time.Unix(int64(unix), 0).In(loc)
	} else {
		t = time.Now().In(loc)
	}

	result := strftime.Format(format, t)
	if len(result) > 4095 {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(result), nil
}

func (f *StrftimeFunction) GetName() string {
	return "strftime"
}

func (f *StrftimeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameter(nil, "timestamp", 1, data.NewNullValue(), nil),
	}
}

func (f *StrftimeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, nil),
		node.NewVariable(nil, "timestamp", 1, nil),
	}
}
