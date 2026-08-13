package data

import (
	"fmt"
)

// NewYieldControlWithContext 创建一个带上下文的 yield 控制流（仅用于普通 yield）
func NewYieldControlWithContext(key Value, value Value, ctx Context) YieldValueControl {
	return &YieldValue{
		Key:     key,
		Value:   value,
		Context: ctx,
	}
}

// YieldValue 表示普通 yield 的控制流值
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
	keyStr := "nil"
	if y.Key != nil {
		keyStr = y.Key.AsString()
	}
	valueStr := "nil"
	if y.Value != nil {
		valueStr = y.Value.AsString()
	}
	return fmt.Sprintf("yield %s => %s", keyStr, valueStr)
}

// GetYieldKey 实现 YieldValueControl 接口
func (y *YieldValue) GetYieldKey() Value {
	return y.Key
}

// GetYieldValue 实现 YieldValueControl 接口
func (y *YieldValue) GetYieldValue() Value {
	return y.Value
}
