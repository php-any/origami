package node

import (
	"github.com/php-any/origami/data"
)

// NullsafeCall 表示空安全调用表达式 ?->
// 如果对象为 null，则返回 null，否则执行调用
type NullsafeCall struct {
	*Node    `pp:"-"`
	Object   data.GetValue // 对象表达式
	CallExpr data.GetValue // 调用表达式（方法调用或属性访问）
}

// NewNullsafeCall 创建一个新的空安全调用表达式
func NewNullsafeCall(from *TokenFrom, object data.GetValue, callExpr data.GetValue) *NullsafeCall {
	return &NullsafeCall{
		Node:     NewNode(from),
		Object:   object,
		CallExpr: callExpr,
	}
}

// GetValue 获取空安全调用表达式的值
func (n *NullsafeCall) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先获取对象的值
	objValue, ctl := n.Object.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 检查对象是否为 null
	if objValue == nil {
		return data.NewNullValue(), nil
	}

	// 检查是否为 NullValue 类型
	switch objValue.(type) {
	case *data.NullValue:
		// 如果对象为 null，返回 null
		return data.NewNullValue(), nil
	}

	// 对象不为 null，执行调用
	// 如果 CallExpr 本身也是 NullsafeCall，需要递归处理
	result, ctl := n.CallExpr.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 如果结果是 null，直接返回
	if result == nil {
		return data.NewNullValue(), nil
	}

	switch result.(type) {
	case *data.NullValue:
		return data.NewNullValue(), nil
	}

	return result, nil
}

// AsString 返回空安全调用表达式的字符串表示
func (n *NullsafeCall) AsString() string {
	return "nullsafe_call"
}
