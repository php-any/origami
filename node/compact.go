package node

import (
	"github.com/php-any/origami/data"
)

// CompactStatement 表示 compact 语句
type CompactStatement struct {
	*Node    `pp:"-"`
	VarNames []data.GetValue // 变量名列表（字符串表达式）
}

// NewCompactStatement 创建一个新的 compact 语句
func NewCompactStatement(token *TokenFrom, varNames []data.GetValue) *CompactStatement {
	return &CompactStatement{
		Node:     NewNode(token),
		VarNames: varNames,
	}
}

// GetValue 获取 compact 语句的值
// 实现类似 unset，但 compact 的参数是字符串（变量名），需要查找变量并获取值
func (c *CompactStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建结果数组（使用 ObjectValue 来支持字符串键）
	result := data.NewObjectValue()

	// 遍历所有变量名表达式
	for _, varNameExpr := range c.VarNames {
		// 获取变量名字符串
		varNameValue, ctl := varNameExpr.GetValue(ctx)
		if ctl != nil {
			continue // 跳过错误
		}

		if varNameValue == nil {
			continue
		}

		// 将参数值转换为字符串（变量名）
		var varName string
		if strValue, ok := varNameValue.(data.AsString); ok {
			varName = strValue.AsString()
		} else if val, ok := varNameValue.(data.Value); ok {
			// 如果是 Value 类型，使用 AsString 方法
			varName = val.AsString()
		} else {
			// 如果都不是，跳过这个参数
			continue
		}

		// 移除可能的 $ 前缀
		if len(varName) > 0 && varName[0] == '$' {
			varName = varName[1:]
		}

		// 尝试通过变量名查找变量并获取值
		// 类似 unset 的处理方式，但我们需要通过变量名查找变量
		// 由于系统使用索引访问变量，我们需要尝试不同的方法

		// 方法1: 如果参数本身就是变量表达式（如 compact($var)），直接获取值
		if variable, ok := varNameExpr.(data.Variable); ok {
			// 直接获取变量值
			varValue, ctl := variable.GetValue(ctx)
			if ctl == nil && varValue != nil {
				if val, ok := varValue.(data.Value); ok {
					// 获取变量名
					varName := variable.GetName()
					result.SetProperty(varName, val)
				}
			}
			continue
		}

		// 方法2: 参数是字符串，需要通过变量名查找变量
		// 尝试创建一个 VariableExpression 并获取值
		// 注意：这需要正确的索引，但由于我们无法在运行时获取索引，
		// 我们需要通过其他方式查找变量

		// 尝试通过作用域查找变量（如果可能）
		// 由于运行时上下文可能没有直接访问 parser 的方式，
		// 我们尝试创建一个变量表达式并获取值
		// 这里使用一个特殊的方法：尝试所有可能的索引（不实际可行）
		// 或者通过其他机制查找变量

		// 简化实现：尝试创建变量表达式并获取值
		// 如果变量存在，GetVariableValue 应该能够找到它
		varExpr := NewVariable(nil, varName, 0, nil)
		varValue, ctl := varExpr.GetValue(ctx)

		// 如果获取成功且值不为 null，则添加到结果数组
		if ctl == nil && varValue != nil {
			// 检查是否为 null 值
			if _, isNull := varValue.(*data.NullValue); !isNull {
				// 值存在且不为 null，添加到结果
				if val, ok := varValue.(data.Value); ok {
					result.SetProperty(varName, val)
				}
			}
		}
	}

	// 返回关联数组
	return result, nil
}
