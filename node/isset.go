package node

import (
	"github.com/php-any/origami/data"
)

// IssetStatement 表示 isset 语句
type IssetStatement struct {
	*Node `pp:"-"`
	Args  []data.GetValue // 参数表达式列表
}

// NewIssetStatement 创建一个新的 isset 语句
func NewIssetStatement(token *TokenFrom, args []data.GetValue) *IssetStatement {
	return &IssetStatement{
		Node: NewNode(token),
		Args: args,
	}
}

// GetValue 获取 isset 语句的值
// isset 检查变量是否存在且不为 null，返回 bool 值
// 如果所有参数都已设置且不为 null，返回 true；否则返回 false
func (i *IssetStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 如果没有参数，返回 false
	if len(i.Args) == 0 {
		return data.NewBoolValue(false), nil
	}

	// 遍历所有参数表达式，检查每个变量是否存在且不为 null
	for _, argExpr := range i.Args {
		// 获取参数值
		varValue, ctl := argExpr.GetValue(ctx)

		// 如果获取值出错，返回 false
		if ctl != nil {
			if acl, ok := ctl.(data.GetName); ok && "UndefinedIndexExpression" == acl.GetName() {
				return data.NewBoolValue(false), nil
			}
			return nil, ctl
		}

		// 检查值是否为 null
		if varValue == nil {
			return data.NewBoolValue(false), nil
		}

		// 检查是否为 null 值
		if _, isNull := varValue.(*data.NullValue); isNull {
			return data.NewBoolValue(false), nil
		}
	}

	// 所有参数都已设置且不为 null，返回 true
	return data.NewBoolValue(true), nil
}
