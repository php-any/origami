package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// NewAnonymousClassExpression 表示匿名类的 new 表达式
type NewAnonymousClassExpression struct {
	*Node
	ClassStmt    data.ClassStmt  // 类定义
	Arguments    []data.GetValue // 构造函数参数
	GenericTypes []data.Types    // 泛型类型（如果有）
}

// NewNewAnonymousClassExpression 创建一个新的匿名类 new 表达式节点
func NewNewAnonymousClassExpression(from data.From, classStmt data.ClassStmt, arguments []data.GetValue, genericTypes []data.Types) *NewAnonymousClassExpression {
	return &NewAnonymousClassExpression{
		Node:         NewNode(from),
		ClassStmt:    classStmt,
		Arguments:    arguments,
		GenericTypes: genericTypes,
	}
}

// GetValue 实现 Value 接口
// 在执行阶段才注册类、实例化对象并调用构造函数
func (n *NewAnonymousClassExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 处理泛型
	var classStmt data.ClassStmt = n.ClassStmt
	if n.GenericTypes != nil {
		if cg, ok := n.ClassStmt.(*ClassGeneric); ok {
			classStmt = cg
		} else if c, ok := n.ClassStmt.(*ClassStatement); ok {
			// 创建泛型类
			cg := &ClassGeneric{
				ClassStatement: c,
				Generic:        n.GenericTypes,
			}
			classStmt = cg
		}
	}

	// 匿名类不需要注册到 VM，直接实例化
	object, acl := classStmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	// 如果有构造函数，调用构造函数
	if object, ok := object.(*data.ClassValue); ok {
		if method := object.Class.GetConstruct(); method != nil {
			varies := method.GetVariables()
			fnCtx := object.CreateContext(varies)
			// 入参的值设置到上下文中
			for index, arg := range n.Arguments {
				switch argTV := arg.(type) {
				case *NamedArgument:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, argTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(n.from, err)
					}
					fnCtx.SetVariableValue(vari, tempV.(data.Value))
				default:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}

					if index >= len(varies) {
						return nil, data.NewErrorThrow(n.from, fmt.Errorf("匿名类构造函数参数数量超出限制: %d", index))
					}

					fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
				}
			}

			_, acl = method.Call(fnCtx)
			if acl != nil {
				return nil, acl
			}
		}
	}

	return object, acl
}
