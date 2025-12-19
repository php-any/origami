package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// NewExpression 表示 new 表达式
type NewExpression struct {
	*Node     `pp:"-"`
	ClassName string
	Arguments []data.GetValue
	// 是否执行构造函数
}

// NewNewExpression 创建一个新的 new 表达式节点
func NewNewExpression(from *TokenFrom, className string, arguments []data.GetValue) *NewExpression {
	return &NewExpression{
		Node:      NewNode(from),
		ClassName: className,
		Arguments: arguments,
	}
}

// GetValue 实现 Value 接口
func (n *NewExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	stmt, acl := vm.GetOrLoadClass(n.ClassName)
	if acl != nil {
		if throwValue, ok := acl.(*data.ThrowValue); ok {
			throwValue.AddStackWithInfo(n.from, n.ClassName, "__construct")
		}
		return nil, acl
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

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

// NewGenerated new T()
type NewGenerated struct {
	*NewExpression
	T string
}

func (n *NewGenerated) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if object, ok := ctx.(*data.ClassMethodContext); ok {
		if classGeneric, ok := object.Class.(*ClassGeneric); ok {
			types, ok := classGeneric.GenericMap[n.T]
			if ok {
				switch t := types.(type) {
				case data.Class:
					className := t.Name
					vm := ctx.GetVM()
					stmt, acl := vm.GetOrLoadClass(className)
					if acl != nil {
						return nil, acl
					}
					return data.NewClassValue(stmt, object.CreateBaseContext()), nil
				}
			}
		}
	}
	return nil, data.NewErrorThrow(n.from, fmt.Errorf("泛型(%v)无法实例化", n.T))
}

// NewClassGenerated DB<Users>
type NewClassGenerated struct {
	*NewExpression
	T []string
}

func (n *NewClassGenerated) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	stmt, acl := vm.GetOrLoadClass(n.ClassName)
	if acl != nil {
		return nil, acl
	}
	mT := make(map[string]data.Types)
	if classGeneric, ok := stmt.(data.ClassGeneric); ok {
		for i, types := range classGeneric.GenericList() {
			newType := n.T[i]
			switch t := types.(type) {
			case data.Generic:
				mT[t.Name] = data.NewBaseType(newType)
			default:
				panic("TODO 未支持的泛型类型")
			}
		}

		stmt = classGeneric.Clone(mT)
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

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
						return nil, data.NewErrorThrow(n.from, fmt.Errorf("对象(%v)构造函数参数数量超出限制: %d", object.Class.GetName(), index))
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

// NewVariableExpression 表示使用变量作为类名的 new 表达式
type NewVariableExpression struct {
	*Node
	ClassNameExpr data.GetValue // 类名表达式（变量）
	Arguments     []data.GetValue
}

// NewNewVariableExpression 创建一个新的使用变量类名的 new 表达式节点
func NewNewVariableExpression(from *TokenFrom, classNameExpr data.GetValue, arguments []data.GetValue) *NewVariableExpression {
	return &NewVariableExpression{
		Node:          NewNode(from),
		ClassNameExpr: classNameExpr,
		Arguments:     arguments,
	}
}

// GetValue 实现 Value 接口
func (n *NewVariableExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类名表达式的值
	classNameValue, acl := n.ClassNameExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	// 将类名值转换为字符串
	var className string
	switch v := classNameValue.(type) {
	case *data.StringValue:
		className = v.Value
	case data.Value:
		// 尝试转换为字符串
		if strValue, ok := v.(*data.StringValue); ok {
			className = strValue.Value
		} else {
			// 尝试调用 AsString 方法
			className = v.AsString()
		}
	default:
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new表达式中的类名变量必须是字符串类型，当前类型: %T", v))
	}

	if className == "" {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new表达式中的类名变量不能为空"))
	}

	vm := ctx.GetVM()
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

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
						return nil, data.NewErrorThrow(n.from, fmt.Errorf("对象(%v)构造函数参数数量超出限制: %d", object.Class.GetName(), index))
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
