package node

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
)

// CallObjectMethod 表示对象属性访问表达式
type CallObjectMethod struct {
	*Node  `pp:"-"`
	Object data.GetValue // 对象表达式
	Method string        // 函数名
	Args   []data.GetValue
}

// NewObjectMethod 创建一个新的对象属性访问表达式
func NewObjectMethod(from *TokenFrom, object data.GetValue, method string, args []data.GetValue) *CallObjectMethod {
	return &CallObjectMethod{
		Node:   NewNode(from),
		Object: object,
		Method: method,
		Args:   args,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallObjectMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	switch class := o.(type) {
	case *data.ThisValue:
		method, has := class.GetMethod(pe.Method)
		if has {
			fnCtx, acl := pe.callMethodParams(class, ctx, method)
			if acl != nil {
				return nil, acl
			}

			return method.Call(fnCtx)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("this 对象不存在对应函数: "+pe.Method))
	case *data.ClassValue:
		method, has := class.GetMethod(pe.Method)
		if has {
			if method.GetModifier() != data.ModifierPublic {
				return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("对象属性访问表达式对象属性访问函数非公开"))
			}
			fnCtx, acl := pe.callMethodParams(class, ctx, method)
			if acl != nil {
				return nil, acl
			}

			return method.Call(fnCtx)
		}

		errStr := fmt.Sprintf("类(%s)不存在对应函数(%s)", class.Class.GetName(), pe.Method)
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(errStr))
	default:
		if class, ok := o.(data.GetMethod); ok {
			method, has := class.GetMethod(pe.Method)
			if has {
				if method.GetModifier() != data.ModifierPublic {
					errStr := fmt.Sprintf("对象属性访问表达式对象属性访问函数(%s)非公开", pe.Method)
					return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(errStr))
				}
				fnCtx, acl := pe.callMethodParams(ctx, ctx, method)
				if acl != nil {
					return nil, acl
				}

				return method.Call(fnCtx)
			}
			return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("当前值不存在函数, 你调用的函数(%s)", pe.Method)))
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("当前值不支持调用函数, 你调用的函数(%s)", pe.Method)))
	}
}

func (pe *CallObjectMethod) callMethodParams(class, ctx data.Context, method data.Method) (data.Context, data.Control) {
	varies := method.GetVariables()
	fnCtx := class.CreateContext(varies)
	// 入参的值设置到上下文中
	for index, arg := range method.GetParams() {
		argClone := arg
		switch argObj := argClone.(type) {
		case *Parameter:
			if index < len(pe.Args) {
				param := pe.Args[index]
				switch paramTV := param.(type) {
				case *NamedArgument:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, paramTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(pe.from, err)
					}
					acl = vari.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				default:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					acl = argObj.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				}
			} else if argObj.DefaultValue == nil {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("调用 %s 函数时参数 %s 缺少值", pe.Method, argObj.Name))
			} else {
				argObj.GetValue(fnCtx)
			}
		case *Parameters:
			args, _ := fnCtx.GetVariableValue(argObj)
			var ares *data.ArrayValue
			var ok bool
			if ares, ok = args.(*data.ArrayValue); !ok {
				ares = data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
				fnCtx.SetVariableValue(argObj, ares)
			}

			for i := index; i < len(pe.Args); i++ {
				param := pe.Args[i]
				tempV, acl := param.GetValue(ctx)
				if acl != nil {
					return nil, acl
				}
				ares.Value = append(ares.Value, tempV.(data.Value))
				fnCtx.SetVariableValue(argObj, ares)
			}
		case *data.ParameterTODO:
			if index < len(pe.Args) {
				param := pe.Args[index]
				switch paramTV := param.(type) {
				case *NamedArgument:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, paramTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(pe.from, err)
					}
					acl = vari.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				default:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					acl = argObj.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				}
			} else if argObj.DefaultValue == nil {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("调用 %s 函数时参数 %s 缺少值", pe.Method, argObj.Name))
			} else {
				argObj.GetValue(fnCtx)
			}
		case *data.ParametersTODO:
			args, _ := fnCtx.GetVariableValue(argObj)
			var ares *data.ArrayValue
			var ok bool
			if ares, ok = args.(*data.ArrayValue); !ok {
				ares = data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
				fnCtx.SetVariableValue(argObj, ares)
			}

			for i := index; i < len(pe.Args); i++ {
				param := pe.Args[i]
				tempV, acl := param.GetValue(ctx)
				if acl != nil {
					return nil, acl
				}
				ares.Value = append(ares.Value, tempV.(data.Value))
				fnCtx.SetVariableValue(argObj, ares)
			}
		case *ParameterReference:
			if index < len(pe.Args) {
				param := pe.Args[index]
				switch paramTV := param.(type) {
				case *NamedArgument:
					vari, err := findVariable(varies, paramTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(pe.from, err)
					}
					if val, ok := paramTV.Value.(data.Variable); ok {
						acl := vari.SetValue(fnCtx, data.NewReferenceValue(val, ctx))
						if acl != nil {
							return nil, acl
						}
					} else {
						return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能传入变量, fn: %s", pe.Method))
					}
				default:
					if val, ok := paramTV.(data.Variable); ok {
						acl := argObj.SetValue(fnCtx, data.NewReferenceValue(val, ctx))
						if acl != nil {
							return nil, acl
						}
					} else {
						return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能传入变量, fn: %s", pe.Method))
					}
				}
			} else {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能是必传参数, fn: %s", pe.Method))
			}
		case data.Variable:
			if index < len(pe.Args) {
				param := pe.Args[index]
				switch paramTV := param.(type) {
				case *NamedArgument:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, paramTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(pe.from, err)
					}
					acl = vari.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				default:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					acl = argObj.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				}
			} else {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("无法调用函数(%s), 缺少参数", pe.Method))
			}
		}
	}

	return fnCtx, nil
}

func findVariable(varies []data.Variable, name string) (data.Variable, error) {
	for _, vary := range varies {
		check := vary.GetName()
		if check == name {
			return vary, nil
		}
	}
	return nil, errors.New("无法找到变量: " + name)
}
