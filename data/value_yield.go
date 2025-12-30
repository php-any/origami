package data

import (
	"fmt"
)

// NewYieldControlWithContext 创建一个带上下文的 yield 控制流
func NewYieldControlWithContext(key Value, value Value, ctx Context) YieldValueControl {
	return &YieldValue{
		Key:     key,
		Value:   value,
		Context: ctx,
	}
}

// YieldValue 表示 yield 控制流值
type YieldValue struct {
	Key     Value
	Value   Value
	Context Context
}

// GetValue 实现 Value 接口
func (y *YieldValue) GetValue(ctx Context) (GetValue, Control) {
	return y.Value, nil
}

// AsString 实现 Value 接口
func (y *YieldValue) AsString() string {
	return fmt.Sprintf("yield %s => %s", y.Key.AsString(), y.Value.AsString())
}

// GetYieldKey 实现 YieldValueControl 接口
func (y *YieldValue) GetYieldKey() Value {
	return y.Key
}

// GetYieldValue 实现 YieldValueControl 接口
func (y *YieldValue) GetYieldValue() Value {
	return y.Value
}
