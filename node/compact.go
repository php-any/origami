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
// PHP compact() 将变量名和它们的值打包成一个关联数组
func (c *CompactStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建结果关联数组（ObjectValue 表示关联数组/字典）
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
			varName = val.AsString()
		} else {
			continue
		}

		// 移除可能的 $ 前缀
		if len(varName) > 0 && varName[0] == '$' {
			varName = varName[1:]
		}

		// 方法1: 如果参数本身就是变量表达式（如 compact($var)），直接获取值
		if variable, ok := varNameExpr.(data.Variable); ok {
			varValue, ctl := variable.GetValue(ctx)
			if ctl == nil && varValue != nil {
				if val, ok := varValue.(data.Value); ok {
					varName := variable.GetName()
					result.SetProperty(varName, val)
				}
			}
			continue
		}

		// 方法2: 参数是字符串，通过变量名查找变量
		varExpr := NewVariable(nil, varName, 0, nil)
		varValue, ctl := varExpr.GetValue(ctx)

		if ctl == nil && varValue != nil {
			if _, isNull := varValue.(*data.NullValue); !isNull {
				if val, ok := varValue.(data.Value); ok {
					result.SetProperty(varName, val)
				}
			}
		}
	}

	// 返回关联数组（ObjectValue）
	return result, nil
}
