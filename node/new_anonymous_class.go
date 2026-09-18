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

	// PHP：匿名类在首次 new class 后可通过 get_class/class_exists/new $name 使用（Livewire __mountParamsContainer）
	if vm := ctx.GetVM(); vm != nil {
		if acl := vm.AddClass(classStmt); acl != nil {
			// 同一定义重复 new：AddClass 可能报已存在；忽略并继续实例化
			if _, ok := vm.GetClass(classStmt.GetName()); !ok {
				return nil, acl
			}
		}
	}

	object, acl := classStmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	// 如果有构造函数，调用构造函数
	if object, ok := object.(*data.ClassValue); ok {
		if method := object.Class.GetConstruct(); method != nil {
			params := method.GetParams()
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

			// 处理未传递的参数，设置默认值
			for index := len(n.Arguments); index < len(params); index++ {
				if index >= len(varies) {
					break
				}
				if argObj, ok := params[index].(*Parameter); ok {
					if argObj.DefaultValue == nil {
						continue
					}
					// 调用 GetValue 来触发默认值的设置
					_, acl := argObj.GetValue(fnCtx)
					if acl != nil {
						return nil, acl
					}
				}
			}

			// 将构造函数参数属性的值赋值给对象属性（PHP 8 构造函数参数属性提升）
			for index, param := range params {
				// 检查是否是属性提升的参数
				if promotedParam, ok := param.(*PromotedParameter); ok {
					// 从函数上下文获取参数值
					if index < len(varies) {
						paramValue, acl := fnCtx.GetVariableValue(varies[index])
						if acl != nil {
							// 如果获取失败，尝试使用默认值
							if promotedParam.DefaultValue != nil {
								paramValueGet, acl := promotedParam.DefaultValue.GetValue(fnCtx)
								if acl != nil {
									return nil, acl
								}
								if paramValueGet != nil {
									paramValue = paramValueGet.(data.Value)
								}
							} else {
								// 没有默认值，跳过
								continue
							}
						}
						// 将参数值赋值给对象属性
						if paramValue != nil {
							object.SetProperty(promotedParam.PropertyName, paramValue.(data.Value))
						}
					}
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
