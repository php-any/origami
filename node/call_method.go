package node

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
)

type CallMethod struct {
	*Node  `pp:"-"`
	Method data.GetValue // 指向一个函数
	Args   []data.GetValue
}

// NewCallMethod 创建一个新的对象属性访问表达式
func NewCallMethod(token *TokenFrom, method data.GetValue, args []data.GetValue) *CallMethod {
	return &CallMethod{
		Node:   NewNode(token),
		Method: method,
		Args:   args,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	call, acl := pe.Method.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	if fv, ok := call.(*data.FuncValue); ok {
		fn := fv.Value
		varies := fn.GetVariables()
		fnCtx := ctx.CreateContext(varies)
		// 入参的值设置到上下文中
		for index, arg := range fn.GetParams() {
			switch argObj := arg.(type) {
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
					return nil, pe.newFunParamsError(pe.GetFrom(), fn.GetName(), argObj.Name)
				} else {
					argObj.GetValue(fnCtx)
				}
			case *Parameters:
				args, _ := fnCtx.GetVariableValue(argObj)
				var ares *data.ArrayValue
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
			}
		}

		return fn.Call(fnCtx)
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("不存在对应函数"))
}

func (pe *CallMethod) newFunParamsError(from data.From, name string, paramName string) data.Control {
	if name == "" {
		return data.NewErrorThrow(from, errors.New("无法调用匿名函数, 缺少参数:"+paramName))
	}
	return data.NewErrorThrow(from, errors.New("无法调用"+name+"函数, 缺少参数:"+paramName))
}
