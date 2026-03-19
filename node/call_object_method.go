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
				if _, ok := acl.(ToClosure); ok {
					return data.NewFuncValue(method), nil
				}
				return nil, acl
			}

			return method.Call(fnCtx)
		}
		// 方法未找到时尝试魔法方法 __call(string $name, array $arguments)
		if magic, hasCall := class.GetMethod("__call"); hasCall {
			return pe.invokeMagicCall(class, ctx, magic, pe.Method, pe.Args)
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
				if _, ok := acl.(ToClosure); ok {
					return data.NewFuncValue(method), nil
				}
				return nil, acl
			}

			return method.Call(fnCtx)
		}
		// 方法未找到时尝试魔法方法 __call(string $name, array $arguments)
		if magic, hasCall := class.GetMethod("__call"); hasCall {
			return pe.invokeMagicCall(class, ctx, magic, pe.Method, pe.Args)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)不存在对应函数(%s)", class.Class.GetName(), pe.Method))
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
					if _, ok := acl.(ToClosure); ok {
						return data.NewFuncValue(method), nil
					}
					return nil, acl
				}

				return method.Call(fnCtx)
			}
			// 方法未找到时尝试魔法方法 __call，$this 为当前对象
			if magic, hasCall := class.GetMethod("__call"); hasCall {
				if objCtx, ok := o.(data.Context); ok {
					return pe.invokeMagicCall(objCtx, ctx, magic, pe.Method, pe.Args)
				}
			}
		}
	}
	return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前值(%#v)不支持调用函数, 你调用的函数(%s)", TryGetCallClassName(o), pe.Method))
}

// invokeMagicCall 调用 __call(string $name, array $arguments)，用于未定义方法时的魔法分发
func (pe *CallObjectMethod) invokeMagicCall(object data.Context, ctx data.Context, magic data.Method, methodName string, args []data.GetValue) (data.GetValue, data.Control) {
	var argsList []data.Value
	for _, arg := range args {
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			argsList = append(argsList, val)
		} else {
			argsList = append(argsList, data.NewNullValue())
		}
	}
	varies := magic.GetVariables()
	if len(varies) < 2 {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__call 需要至少 2 个参数 (name, arguments)"))
	}
	fnCtx := object.CreateContext(varies)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(methodName))
	fnCtx.SetVariableValue(varies[1], data.NewArrayValue(argsList))
	return magic.Call(fnCtx)
}

func (pe *CallObjectMethod) callMethodParams(object, ctx data.Context, method data.Method) (data.Context, data.Control) {
	varies := method.GetVariables()
	fnCtx := object.CreateContext(varies)
	// 入参的值设置到上下文中
	for index, param := range method.GetParams() {
		if len(pe.Args) > index {
			var acl data.Control
			arg := pe.Args[index]
			var tempV data.GetValue
			switch argTV := arg.(type) {
			case *NamedArgument:
				tempV, acl = argTV.GetValue(ctx)
				if acl != nil {
					return nil, acl
				}
				vari, err := findVariable(varies, argTV.Name)
				if err != nil {
					return nil, data.NewErrorThrow(pe.from, err)
				}
				fnCtx.SetVariableValue(vari, tempV.(data.Value))
				if promotedParam, ok := param.(*PromotedParameter); ok {
					acl = promotedParam.SetValue(object, tempV.(data.Value))
				}
			default:
				tempV, acl = argTV.GetValue(ctx)
				if acl != nil {
					return nil, acl
				}
				if index >= len(varies) {
					return nil, data.NewErrorThrow(pe.from, fmt.Errorf("对象 (%v) 构造函数参数数量超出限制：%d", object, index))
				}
				fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
				if promotedParam, ok := param.(*PromotedParameter); ok {
					acl = promotedParam.SetValue(object, tempV.(data.Value))
				}
			}
			if acl != nil {
				return nil, acl
			}
		} else if promotedParam, ok := param.(*PromotedParameter); ok {
			// 触发初始化默认值
			_, acl := promotedParam.GetValue(object)
			if acl != nil {
				return nil, acl
			}
		} else if argObj, ok := param.(*Parameter); ok {
			if argObj.DefaultValue == nil {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("调用 %s 构造函数时参数 %s 缺少值和默认值", object, argObj.Name))
			}
			// 调用 GetValue 来触发默认值的设置
			_, acl := argObj.GetValue(fnCtx)
			if acl != nil {
				return nil, acl
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

func findParams(varies []data.GetValue, name string) (data.GetValue, error) {
	for _, vary := range varies {
		if check, ok := vary.(data.GetName); ok {
			if check.GetName() == name {
				return vary, nil
			}
		}
	}
	return nil, errors.New("无法找到变量: " + name)
}
