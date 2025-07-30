package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ForeachStatement 表示foreach语句
type ForeachStatement struct {
	*Node `pp:"-"`
	Array data.GetValue // 要遍历的数组
	Key   data.Variable // 键变量名（可选）
	Value data.Variable // 值变量名
	Body  []data.GetValue
}

func (u *ForeachStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数组值
	arrayValue, ctl := u.Array.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 检查数组值是否为数组类型
	switch array := arrayValue.(type) {
	case *data.ArrayValue:
		var v data.GetValue
		var c data.Control

		// 遍历数组
		for i, element := range array.Value {
			// 设置值变量
			ctx.SetVariableValue(u.Value, element)

			// 如果有键变量，设置键变量
			if u.Key != nil {
				keyValue := data.NewIntValue(i)
				ctx.SetVariableValue(u.Key, keyValue)
			}

			// 执行循环体
			for _, statement := range u.Body {
				v, c = statement.GetValue(ctx)
				if c != nil {
					// break 跳出循环
					if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
						return nil, nil
					}
					// continue 跳到下一次迭代
					if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
						continue
					}
					// return/throw 直接返回
					return nil, c
				}
			}
		}
		return v, nil
	case *data.ObjectValue:
		var v data.GetValue
		var c data.Control

		// 遍历数组
		for i, element := range array.GetProperties() {
			// 设置值变量
			ctx.SetVariableValue(u.Value, element)

			// 如果有键变量，设置键变量
			if u.Key != nil {
				keyValue := data.NewStringValue(i)
				ctx.SetVariableValue(u.Key, keyValue)
			}

			// 执行循环体
			for _, statement := range u.Body {
				v, c = statement.GetValue(ctx)
				if c != nil {
					// break 跳出循环
					if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
						return nil, nil
					}
					// continue 跳到下一次迭代
					if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
						continue
					}
					// return/throw 直接返回
					return nil, c
				}
			}
		}
		return v, nil
	case *data.NullValue:
		return nil, nil
	}

	return nil, data.NewErrorThrow(u.from, fmt.Errorf("foreach 只能遍历数组"))
}

// NewForeachStatement 创建一个新的foreach语句
func NewForeachStatement(token *TokenFrom, array data.GetValue, key data.Variable, value data.Variable, body []data.GetValue) *ForeachStatement {
	return &ForeachStatement{
		Node:  NewNode(token),
		Array: array,
		Key:   key,
		Value: value,
		Body:  body,
	}
}
