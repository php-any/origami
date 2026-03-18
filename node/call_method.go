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

	switch fv := call.(type) {
	case *data.FuncValue:
		return pe.handleFuncValue(ctx, call)
	case *StaticMethodFuncValue:
		// 静态方法包装器，调用 GetValue 获取 FuncValue 然后继续处理
		funcValue, acl := fv.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		// 递归处理，现在应该是 FuncValue 了
		return pe.handleFuncValue(ctx, funcValue)
	default:
		// 魔法方法 __invoke：对象作为可调用时调用 $object->__invoke(...$args)
		if obj, ok := call.(data.GetMethod); ok {
			if invoke, has := obj.GetMethod("__invoke"); has {
				if objCtx, ok := call.(data.Context); ok {
					return pe.invokeMagicInvoke(ctx, objCtx, invoke)
				}
			}
		}
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("不存在对应函数:"+TryGetCallClassName(pe.Method)))
}

// handleFuncValue 处理 FuncValue 类型的调用
func (pe *CallMethod) handleFuncValue(ctx data.Context, call data.GetValue) (data.GetValue, data.Control) {
	fv, ok := call.(*data.FuncValue)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("期望 FuncValue 类型"))
	}
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
			args, acl := fnCtx.GetVariableValue(argObj)
			if acl != nil {
				return nil, acl
			}
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
				ares.List = append(ares.List, data.NewZVal(tempV.(data.Value)))
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
					switch val := paramTV.Value.(type) {
					case *CallObjectProperty:
						// $obj->prop 作为引用参数：共享 ZVal 指针
						zv, acl := val.GetZVal(ctx)
						if acl != nil {
							return nil, acl
						}
						fnCtx.SetIndexZVal(vari.(*ParameterReference).Index, zv)
					case data.Variable:
						acl := vari.SetValue(fnCtx, data.NewReferenceValue(val, ctx))
						if acl != nil {
							return nil, acl
						}
					default:
						return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能传入变量, fn: %s", pe.Method))
					}
				case *CallObjectProperty:
					// $obj->prop 作为引用参数：通过 GetZVal 共享 ZVal 指针，而非走 ReferenceValue 路径
					zv, acl := paramTV.GetZVal(ctx)
					if acl != nil {
						return nil, acl
					}
					fnCtx.SetIndexZVal(argObj.Index, zv)
				default:
					if val, ok := paramTV.(data.Variable); ok {
						zv := ctx.GetIndexZVal(val.GetIndex())
						fnCtx.SetIndexZVal(argObj.Index, zv)
					} else {
						return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能传入变量, fn: %s", pe.Method))
					}
				}
			} else {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能是必传参数, fn: %s", pe.Method))
			}
		}
	}

	// 将本次调用的参数表达式列表记录到方法上下文中
	fnCtx.SetCallArgs(pe.Args)

	return fn.Call(fnCtx)
}

// invokeMagicInvoke 调用对象的 __invoke(...$args)，用于对象作为可调用时的魔法分发
func (pe *CallMethod) invokeMagicInvoke(ctx data.Context, object data.Context, invoke data.Method) (data.GetValue, data.Control) {
	varies := invoke.GetVariables()
	fnCtx := object.CreateContext(varies)
	for i, arg := range pe.Args {
		if i >= len(varies) {
			break
		}
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			fnCtx.SetVariableValue(varies[i], val)
		}
	}
	return invoke.Call(fnCtx)
}

func (pe *CallMethod) newFunParamsError(from data.From, name string, paramName string) data.Control {
	if name == "" {
		return data.NewErrorThrow(from, errors.New("无法调用匿名函数, 缺少参数:"+paramName))
	}
	return data.NewErrorThrow(from, errors.New("无法调用("+name+")函数, 缺少参数: "+paramName))
}
