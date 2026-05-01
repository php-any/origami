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
			for index, param := range params {
				if len(arguments) > index {
					arg := arguments[index]
					switch argTV := arg.(type) {
					case *NamedArgument:
						param, err := findParams(params, argTV.Name)
						if err != nil {
							return nil, data.NewErrorThrow(from, err)
						}
						acl = paramSetValue(fnCtx, ctx, object, param, argTV, varies, index, arguments)
					default:
						acl = paramSetValue(fnCtx, ctx, object, param, argTV, varies, index, arguments)
					}
				} else {
					switch param := param.(type) {
					case *PromotedParameter:
						// 触发初始化默认值
						_, acl = param.GetValue(object)
					case *CallerContextParameter:
						fnCtx = ctx
					default:
						_, acl = param.GetValue(fnCtx)
					}
				}
				if acl != nil {
					return nil, acl
				}
			}
			if acl != nil {
				return nil, acl
			}
			// 将本次调用的参数表达式列表记录到函数上下文中
			fnCtx.SetCallArgs(arguments)
			_, acl = method.Call(fnCtx)
			if acl != nil {
				return nil, acl
			}
		}
	}

	return object, acl
}

func paramSetValue(fnCtx, ctx, object data.Context, param, argTV data.GetValue, varies []data.Variable, index int, arguments []data.GetValue) data.Control {
	switch param := param.(type) {
	case *ParameterReference:
		switch val := arguments[index].(type) {
		case *CallObjectProperty:
			zv, acl := val.GetZVal(ctx)
			if acl != nil {
				return acl
			}
			fnCtx.SetIndexZVal(param.Index, zv)
		case *IndexExpression:
			zv, acl := val.GetZVal(ctx)
			if acl != nil {
				return acl
			}
			fnCtx.SetIndexZVal(param.Index, zv)
		case data.Variable:
			// 引用参数传入变量时，需要保证存在一个共享的 ZVal：
			zv := ctx.GetIndexZVal(val.GetIndex())
			fnCtx.SetIndexZVal(param.Index, zv)
		default:
			return data.NewErrorThrow(param.GetFrom(), fmt.Errorf("引用参数只能传入变量"))
		}
		return nil
	case *Parameters: // 可变参数
		args, acl := fnCtx.GetVariableValue(param)
		var ares *data.ArrayValue
		var ok bool
		if ares, ok = args.(*data.ArrayValue); !ok {
			ares = data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
		}

		for i := index; i < len(arguments); i++ {
			arg := arguments[i]
			// 支持展开实参 ...expr
			if spread, ok := arg.(*SpreadArgument); ok {
				tempV, acl := spread.Expr.GetValue(ctx)
				if acl != nil {
					return acl
				}
				if tempV == nil {
					continue
				}
				switch v := tempV.(type) {
				case *data.ArrayValue:
					for _, z := range v.List {
						ares.List = append(ares.List, data.NewZVal(z.Value))
					}
					fnCtx.SetVariableValue(param, ares)
				case *data.ObjectValue:
					// 关联数组展开：按属性遍历值（键在具体函数内部再决策如何使用）
					v.RangeProperties(func(_ string, val data.Value) bool {
						ares.List = append(ares.List, data.NewZVal(val))
						return true
					})
					fnCtx.SetVariableValue(param, ares)
				default:
					// 其他类型退化为普通单值参数
					if value, ok := tempV.(data.Value); ok {
						ares.List = append(ares.List, data.NewZVal(value))
						fnCtx.SetVariableValue(param, ares)
					}
				}
				continue
			}

			tempV, acl := arg.GetValue(ctx)
			if acl != nil {
				return acl
			}
			if tempV == nil {
				ares.List = append(ares.List, data.NewZVal(data.NewNullValue()))
				fnCtx.SetVariableValue(param, ares)
			} else {
				ares.List = append(ares.List, data.NewZVal(tempV.(data.Value)))
				fnCtx.SetVariableValue(param, ares)
			}
		}
		return acl
	case *PromotedParameter: // 属性提升
		tempV, acl := argTV.GetValue(ctx)
		if acl != nil {
			return acl
		}
		if index >= len(varies) {
			return data.NewErrorThrow(nil, fmt.Errorf("对象构造函数参数数量超出限制"))
		}
		acl = fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
		if acl == nil {
			acl = param.SetValue(object, tempV.(data.Value))
		}
		return acl
	case *ParameterRawAST:
		return fnCtx.SetVariableValue(param.Parameter, data.NewASTValue(argTV, ctx))
	case *Parameter:
		// 普通参数：求值实参，然后设置到上下文
		tempV, acl := argTV.GetValue(ctx)
		if acl != nil {
			return acl
		}
		if tempV == nil {
			return nil
		}
		return param.SetValue(fnCtx, tempV.(data.Value))
	case data.Variable:
		tempV, acl := argTV.GetValue(ctx)
		if acl != nil {
			return acl
		}
		return fnCtx.SetVariableValue(param, tempV.(data.Value))
	}

	return data.NewErrorThrow(nil, fmt.Errorf("无法识别的参数类型 %T", param))
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
	// 检查是否在类上下文中（类方法或类级初始化器）
	var currentClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		currentClass = classCtx.Class
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		currentClass = classVal.Class
	} else {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new self 只能在类方法中使用"))
	}

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
	// 检查是否在类上下文中（类方法或类级初始化器）
	var currentClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		currentClass = classCtx.Class
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		currentClass = classVal.Class
	} else {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new static 只能在类方法中使用"))
	}

	// 获取当前类的类名
	// TODO: 实现真正的 late static binding，返回实际调用时的类名（子类）
	className := currentClass.GetName()

	return createInstanceAndCallConstructor(n.from, className, n.Arguments, ctx)
}

// NewExpressionDynamic 表示 new $expr(...) 动态类名实例化
type NewExpressionDynamic struct {
	*Node     `pp:"-"`
	ClassExpr data.GetValue // 类名表达式（运行时求值为字符串）
	Arguments []data.GetValue
}

func NewNewExpressionDynamic(from *TokenFrom, classExpr data.GetValue, arguments []data.GetValue) *NewExpressionDynamic {
	return &NewExpressionDynamic{
		Node:      NewNode(from),
		ClassExpr: classExpr,
		Arguments: arguments,
	}
}

func (n *NewExpressionDynamic) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 求值类名表达式
	classVal, acl := n.ClassExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	className := ""
	if s, ok := classVal.(data.AsString); ok {
		className = s.AsString()
	}
	if className == "" {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("new 表达式类名求值结果为空"))
	}
	return createInstanceAndCallConstructor(n.from, className, n.Arguments, ctx)
}
