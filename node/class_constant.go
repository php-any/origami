package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ClassConstant 表示 ::class 语法节点
type ClassConstant struct {
	From *TokenFrom
	Expr data.GetValue
}

// NewClassConstant 创建新的 ClassConstant 节点
func NewClassConstant(from *TokenFrom, expr data.GetValue) *ClassConstant {
	return &ClassConstant{
		From: from,
		Expr: expr,
	}
}

// GetValue 获取类常量值
func (cc *ClassConstant) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取表达式值
	exprValue, control := cc.Expr.GetValue(ctx)
	if control != nil {
		return nil, control
	}

	// 如果表达式是字符串字面量，直接返回类名
	if strValue, ok := exprValue.(*StringLiteral); ok {
		className := strValue.Value
		// 尝试获取完整的类地址（包括命名空间）
		class, acl := ctx.GetVM().GetOrLoadClass(className)
		if acl == nil {
			// 返回完整的类名（包括命名空间）
			return data.NewStringValue(class.GetName()), nil
		} else {
			return nil, acl
		}
	}

	// 如果表达式是变量，获取变量的值
	if varExpr, ok := exprValue.(*VariableExpression); ok {
		// 获取变量的值
		if varValue, control := varExpr.GetValue(ctx); control == nil {
			switch v := varValue.(type) {
			case *data.ClassValue:
				// 返回完整的类名（包括命名空间）
				return data.NewStringValue(v.Class.GetName()), nil
			case *data.StringValue:
				// 尝试获取完整的类地址
				if vm := ctx.GetVM(); vm != nil {
					if class, acl := vm.GetOrLoadClass(v.AsString()); acl == nil {
						return data.NewStringValue(class.GetName()), nil
					} else {
						return nil, acl
					}
				}
			}
		}
		return nil, data.NewErrorThrow(cc.From, fmt.Errorf("无法获取变量类型"))
	}

	// 将表达式转换为字符串
	if value, ok := exprValue.(data.Value); ok {
		className := value.AsString()
		// 尝试获取完整的类地址
		if vm := ctx.GetVM(); vm != nil {
			if class, acl := vm.GetOrLoadClass(className); acl == nil {
				return data.NewStringValue(class.GetName()), nil
			} else {
				return nil, acl
			}
		}
		return data.NewStringValue(className), nil
	}

	// 如果无法转换，返回空字符串
	return data.NewStringValue(""), nil
}
