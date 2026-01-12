package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// createInstanceAndCallConstructor 创建类实例并调用构造函数
// 这是一个公共辅助函数，用于减少重复代码
func createInstanceAndCallConstructor(
	from data.From,
	className string,
	arguments []data.GetValue,
	ctx data.Context,
) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		if throwValue, ok := acl.(*data.ThrowValue); ok {
			throwValue.AddStackWithInfo(from, className, "__construct")
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
			params := method.GetParams()
			fnCtx := object.CreateContext(varies)
			// 入参的值设置到上下文中
			for index, arg := range arguments {
				switch argTV := arg.(type) {
				case *NamedArgument:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, argTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(from, err)
					}
					fnCtx.SetVariableValue(vari, tempV.(data.Value))
				default:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}

					if index >= len(varies) {
						return nil, data.NewErrorThrow(from, fmt.Errorf("对象(%v)构造函数参数数量超出限制: %d", object.Class.GetName(), index))
					}

					fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
				}
			}

			// 处理未传递的参数，设置默认值
			for index := len(arguments); index < len(params); index++ {
				if index >= len(varies) {
					break
				}
				if argObj, ok := params[index].(*Parameter); ok {
					if argObj.DefaultValue == nil {
						return nil, data.NewErrorThrow(from, fmt.Errorf("调用 %s 构造函数时参数 %s 缺少值和默认值", object.Class.GetName(), argObj.Name))
					}
					// 调用 GetValue 来触发默认值的设置
					_, acl := argObj.GetValue(fnCtx)
					if acl != nil {
						return nil, acl
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

// createInstanceAndCallConstructorWithStmt 使用已加载的类语句创建实例并调用构造函数
// 这是 createInstanceAndCallConstructor 的变体，用于已经加载并处理过的类（如泛型类）
func createInstanceAndCallConstructorWithStmt(
	from data.From,
	stmt data.ClassStmt,
	arguments []data.GetValue,
	ctx data.Context,
) (data.GetValue, data.Control) {
	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	if object, ok := object.(*data.ClassValue); ok {
		if method := object.Class.GetConstruct(); method != nil {
			varies := method.GetVariables()
			params := method.GetParams()
			fnCtx := object.CreateContext(varies)
			// 入参的值设置到上下文中
			for index, arg := range arguments {
				switch argTV := arg.(type) {
				case *NamedArgument:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, argTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(from, err)
					}
					fnCtx.SetVariableValue(vari, tempV.(data.Value))
				default:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}

					if index >= len(varies) {
						return nil, data.NewErrorThrow(from, fmt.Errorf("对象(%v)构造函数参数数量超出限制: %d", object.Class.GetName(), index))
					}

					fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
				}
			}

			// 处理未传递的参数，设置默认值
			for index := len(arguments); index < len(params); index++ {
				if index >= len(varies) {
					break
				}
				if argObj, ok := params[index].(*Parameter); ok {
					if argObj.DefaultValue == nil {
						return nil, data.NewErrorThrow(from, fmt.Errorf("调用 %s 构造函数时参数 %s 缺少值和默认值", object.Class.GetName(), argObj.Name))
					}
					// 调用 GetValue 来触发默认值的设置
					_, acl := argObj.GetValue(fnCtx)
					if acl != nil {
						return nil, acl
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
	return createInstanceAndCallConstructor(n.from, n.ClassName, n.Arguments, ctx)
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

	return createInstanceAndCallConstructorWithStmt(n.from, stmt, n.Arguments, ctx)
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

	return createInstanceAndCallConstructor(n.from, className, n.Arguments, ctx)
}

// NewSelfExpression 表示 new self 表达式
type NewSelfExpression struct {
	*Node     `pp:"-"`
	Arguments []data.GetValue
}

// NewNewSelfExpression 创建一个新的 new self 表达式节点
func NewNewSelfExpression(from *TokenFrom, arguments []data.GetValue) *NewSelfExpression {
	return &NewSelfExpression{
		Node:      NewNode(from),
		Arguments: arguments,
	}
}

// GetValue 实现 Value 接口
func (n *NewSelfExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new self 只能在类方法中使用"))
	}

	// 获取当前类名
	currentClass := classCtx.Class
	className := currentClass.GetName()

	return createInstanceAndCallConstructor(n.from, className, n.Arguments, ctx)
}

// NewStaticExpression 表示 new static 表达式
// 在 PHP 中，new static 使用 late static binding，创建实际调用时的类实例
// 当前实现暂时使用当前类的类名，后续可以增强为真正的 late static binding
type NewStaticExpression struct {
	*Node     `pp:"-"`
	Arguments []data.GetValue
}

// NewNewStaticExpression 创建一个新的 new static 表达式节点
func NewNewStaticExpression(from *TokenFrom, arguments []data.GetValue) *NewStaticExpression {
	return &NewStaticExpression{
		Node:      NewNode(from),
		Arguments: arguments,
	}
}

// GetValue 实现 Value 接口
func (n *NewStaticExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new static 只能在类方法中使用"))
	}

	// 获取当前类的类名
	// TODO: 实现真正的 late static binding，返回实际调用时的类名（子类）
	currentClass := classCtx.Class
	className := currentClass.GetName()

	return createInstanceAndCallConstructor(n.from, className, n.Arguments, ctx)
}
