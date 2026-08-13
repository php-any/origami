package node

import (
	"errors"
	"fmt"
	"strings"

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
		ret, acl := pe.handleFuncValue(ctx, call)
		if acl != nil {
			if _, ok := acl.(ToClosure); ok {
				return data.NewFuncValue(fv.Value), nil
			}
		}
		return ret, acl
	case *data.BoundFuncValue:
		return pe.handleBoundFuncValue(ctx, fv)
	case *StaticMethodFuncValue:
		// 静态方法包装器，调用 GetValue 获取 FuncValue 然后继续处理
		funcValue, acl := fv.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		// 递归处理，现在应该是 FuncValue 了
		return pe.handleFuncValue(ctx, funcValue)
	case *staticMethodFuncWithLateBinding:
		// 后期静态绑定静态方法包装器
		return pe.handleStaticMethodWithLateBinding(ctx, fv)
	default:
		// 检查是否是 first-class callable: 单个 SpreadArgument(nil)
		allSpread := len(pe.Args) == 1
		if allSpread {
			spread, ok := pe.Args[0].(*SpreadArgument)
			if !ok || spread.Expr != nil {
				allSpread = false
			}
		}
		if allSpread {
			// first-class callable: 返回包装的调用
			return call, nil
		}
		// PHP 可变函数：$fn = 'login_page'; $fn();
		if name := callableFunctionName(call); name != "" {
			fn, ok := ctx.GetVM().GetFunc(name)
			if !ok && len(name) > 0 && name[0] == '\\' {
				fn, ok = ctx.GetVM().GetFunc(name[1:])
			}
			if ok {
				if isFirstClassCallableArgs(pe.Args) {
					return data.NewFuncValue(fn), nil
				}
				return pe.invokeFuncStmt(ctx, fn, fn.Call)
			}
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("Call to undefined function %s()", name))
		}
		// PHP 数组可调用: [$obj, 'method'](...$args) 或 ['ClassName', 'method'](...)
		if arr, ok2 := call.(*data.ArrayValue); ok2 && len(arr.List) >= 2 {
			objVal := arr.List[0].Value
			methodVal := arr.List[1].Value
			methodName := ""
			if sv, ok4 := methodVal.(*data.StringValue); ok4 {
				methodName = sv.Value
			} else if methodVal != nil {
				methodName = methodVal.AsString()
			}
			if methodName != "" {
				if obj, ok3 := objVal.(data.GetMethod); ok3 {
					if method, has := obj.GetMethod(methodName); has {
						return pe.doCallWithArgs(ctx, obj, method)
					}
				}
				if className := callableFunctionName(objVal); className != "" {
					stmt, acl := ctx.GetVM().GetOrLoadClass(className)
					if acl != nil {
						return nil, acl
					}
					if stmt != nil {
						var method data.Method
						var ok bool
						method, ok = stmt.GetMethod(methodName)
						if !ok {
							if sm, ok2 := stmt.(data.GetStaticMethod); ok2 {
								method, ok = sm.GetStaticMethod(methodName)
							}
						}
						if ok {
							fv, acl := NewStaticMethodFuncValue(stmt, method).GetValue(ctx)
							if acl != nil {
								return nil, acl
							}
							return pe.handleFuncValue(ctx, fv)
						}
					}
				}
			}
		}
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

// callableFunctionName 从字符串/可转为字符串的值取出函数名；空串表示不是函数名回调。
func callableFunctionName(v data.GetValue) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case *data.StringValue:
		return strings.TrimSpace(s.Value)
	case *data.NullValue:
		return ""
	case data.AsString:
		// 避免把数组/对象误当成函数名
		if _, ok := v.(*data.ArrayValue); ok {
			return ""
		}
		if _, ok := v.(*data.ObjectValue); ok {
			return ""
		}
		if _, ok := v.(*data.ClassValue); ok {
			return ""
		}
		name := strings.TrimSpace(s.AsString())
		if name == "" || name == "Array" || strings.HasPrefix(name, "Object") {
			return ""
		}
		return name
	default:
		return ""
	}
}

// handleStaticMethodWithLateBinding 处理后期静态绑定的静态方法调用
func (pe *CallMethod) handleStaticMethodWithLateBinding(ctx data.Context, sm *staticMethodFuncWithLateBinding) (data.GetValue, data.Control) {
	fn := sm.method
	varies := fn.GetVariables()

	// 创建类方法上下文，绑定调用时的类（用于 static:: 后期静态绑定）
	classValue := data.NewClassValue(sm.callClass, ctx)
	fnCtx := classValue.CreateContext(varies)
	// 设置后期静态绑定类
	if cmc, ok := fnCtx.(*data.ClassMethodContext); ok {
		cmc.StaticClass = sm.callClass
	}

	params := fn.GetParams()
	flatArgs, namedArgs, acl := flattenCallArgsForBinding(ctx, pe.Args)
	if acl != nil {
		return nil, acl
	}

	bound := make([]bool, len(params))
	for _, named := range namedArgs {
		variable, err := findVariable(varies, named.Name)
		if err != nil {
			return nil, data.NewErrorThrow(pe.from, err)
		}
		for i, v := range varies {
			if v.GetName() == named.Name && i < len(bound) {
				bound[i] = true
				break
			}
		}
		if acl := variable.SetValue(fnCtx, named.Value); acl != nil {
			return nil, acl
		}
	}

	pos := 0
	for index, param := range params {
		if index < len(bound) && bound[index] {
			continue
		}
		switch p := param.(type) {
		case *ParameterReference:
			rawArg := nthPositionalExpr(pe.Args, pos)
			if rawArg == nil {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能是必传参数, fn: %s", fn.GetName()))
			}
			pos++
			if acl := bindByRefParam(fnCtx, ctx, p, rawArg); acl != nil {
				return nil, acl
			}
		case *Parameter:
			if pos < len(flatArgs) {
				if acl := p.SetValue(fnCtx, flatArgs[pos]); acl != nil {
					return nil, acl
				}
				pos++
			} else if p.DefaultValue == nil {
				return nil, pe.newFunParamsError(pe.GetFrom(), fn.GetName(), p.Name)
			} else if _, acl := p.GetValue(fnCtx); acl != nil {
				return nil, acl
			}
		case *Parameters:
			fnCtx.SetVariableValue(p, data.NewArrayValue(flatArgs[pos:]))
			pos = len(flatArgs)
		case data.Parameter:
			if pos < len(flatArgs) {
				if acl := p.SetValue(fnCtx, flatArgs[pos]); acl != nil {
					return nil, acl
				}
				pos++
			} else if p.GetDefaultValue() == nil {
				return nil, pe.newFunParamsError(pe.GetFrom(), fn.GetName(), p.GetName())
			} else if _, acl := p.GetValue(fnCtx); acl != nil {
				return nil, acl
			}
		}
	}

	fnCtx.SetCallArgs(pe.Args)

	return fn.Call(fnCtx)
}

// handleFuncValue 处理 FuncValue 类型的调用
func (pe *CallMethod) handleFuncValue(ctx data.Context, call data.GetValue) (data.GetValue, data.Control) {
	fv, ok := call.(*data.FuncValue)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("期望 FuncValue 类型"))
	}
	// PHP 8.1: $fn(...) / $listener(...) 一等可调用，返回闭包本身
	if isFirstClassCallableArgs(pe.Args) {
		return fv, nil
	}
	return pe.invokeFuncStmt(ctx, fv.Value, fv.Call)
}

// handleBoundFuncValue 处理 bindTo/bind 后的 BoundFuncValue 变量调用 $fn(...)
func (pe *CallMethod) handleBoundFuncValue(ctx data.Context, bfv *data.BoundFuncValue) (data.GetValue, data.Control) {
	if isFirstClassCallableArgs(pe.Args) {
		return bfv, nil
	}
	return pe.invokeFuncStmt(ctx, bfv.Value, bfv.Call)
}

func (pe *CallMethod) invokeFuncStmt(ctx data.Context, fn data.FuncStmt, invoke func(data.Context) (data.GetValue, data.Control)) (data.GetValue, data.Control) {
	varies := fn.GetVariables()
	fnCtx := ctx.CreateContext(varies)
	params := fn.GetParams()

	// 先展开 ...$arr，再按位置/命名绑定（Laravel Event：$listener(...array_values($payload))）
	flatArgs, namedArgs, acl := flattenCallArgsForBinding(ctx, pe.Args)
	if acl != nil {
		return nil, acl
	}

	bound := make([]bool, len(params))
	for _, na := range namedArgs {
		vari, err := findVariable(varies, na.Name)
		if err != nil {
			return nil, data.NewErrorThrow(pe.from, err)
		}
		idx := -1
		for i, v := range varies {
			if v.GetName() == na.Name {
				idx = i
				break
			}
		}
		if idx >= 0 && idx < len(bound) {
			bound[idx] = true
		}
		if acl := vari.SetValue(fnCtx, na.Value); acl != nil {
			return nil, acl
		}
	}

	pos := 0
	for index, arg := range params {
		if index < len(bound) && bound[index] {
			continue
		}
		switch argObj := arg.(type) {
		case *Parameter:
			if pos < len(flatArgs) {
				acl := argObj.SetValue(fnCtx, flatArgs[pos])
				pos++
				if acl != nil {
					return nil, acl
				}
			} else if argObj.DefaultValue == nil {
				return nil, pe.newFunParamsError(pe.GetFrom(), fn.GetName(), argObj.Name)
			} else {
				argObj.GetValue(fnCtx)
			}
		case *Parameters:
			remaining := flatArgs[pos:]
			fnCtx.SetVariableValue(argObj, data.NewArrayValue(remaining))
			pos = len(flatArgs)
		case *ParameterReference:
			rawArg := nthPositionalExpr(pe.Args, pos)
			if rawArg == nil {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能是必传参数, fn: %s", pe.Method))
			}
			pos++
			if acl := bindByRefParam(fnCtx, ctx, argObj, rawArg); acl != nil {
				return nil, acl
			}
		case data.Parameter:
			if pos < len(flatArgs) {
				acl := argObj.SetValue(fnCtx, flatArgs[pos])
				pos++
				if acl != nil {
					return nil, acl
				}
			} else if argObj.GetDefaultValue() == nil {
				return nil, pe.newFunParamsError(pe.GetFrom(), fn.GetName(), argObj.GetName())
			} else {
				argObj.GetValue(fnCtx)
			}
		}
	}

	// 将本次调用的参数表达式列表记录到方法上下文中
	fnCtx.SetCallArgs(pe.Args)

	return invoke(fnCtx)
}

type namedArgValue struct {
	Name  string
	Value data.Value
}

// flattenCallArgsForBinding 将调用实参展开为位置实参列表 + 命名实参
func flattenCallArgsForBinding(ctx data.Context, args []data.GetValue) ([]data.Value, []namedArgValue, data.Control) {
	var flat []data.Value
	var named []namedArgValue
	for _, arg := range args {
		switch a := arg.(type) {
		case *NamedArgument:
			v, acl := a.GetValue(ctx)
			if acl != nil {
				return nil, nil, acl
			}
			val, _ := v.(data.Value)
			if val == nil {
				val = data.NewNullValue()
			}
			named = append(named, namedArgValue{Name: a.Name, Value: val})
		case *SpreadArgument:
			if a.Expr == nil {
				continue
			}
			spreadVal, acl := a.GetValue(ctx)
			if acl != nil {
				return nil, nil, acl
			}
			if arr, ok := spreadVal.(*data.ArrayValue); ok {
				for _, z := range arr.List {
					flat = append(flat, z.Value)
				}
			} else if objVal, ok := spreadVal.(*data.ObjectValue); ok {
				objVal.RangeProperties(func(_ string, value data.Value) bool {
					flat = append(flat, value)
					return true
				})
			}
		default:
			v, acl := arg.GetValue(ctx)
			if acl != nil {
				return nil, nil, acl
			}
			if val, ok := v.(data.Value); ok && val != nil {
				flat = append(flat, val)
			} else {
				flat = append(flat, data.NewNullValue())
			}
		}
	}
	return flat, named, nil
}

// nthPositionalExpr 返回第 n 个非命名位置实参表达式（不求值，供引用参数共享 ZVal）
func nthPositionalExpr(args []data.GetValue, n int) data.GetValue {
	i := 0
	for _, a := range args {
		if _, ok := a.(*NamedArgument); ok {
			continue
		}
		if _, ok := a.(*SpreadArgument); ok {
			// 展开实参无法作为引用目标；跳过计数由 flatArgs 路径覆盖
			continue
		}
		if i == n {
			return a
		}
		i++
	}
	return nil
}

// bindByRefParam 将实参以引用方式绑定到形参（共享 ZVal，写回调用方）
func bindByRefParam(fnCtx, callCtx data.Context, param *ParameterReference, rawArg data.GetValue) data.Control {
	switch v := rawArg.(type) {
	case *NamedArgument:
		return bindByRefParam(fnCtx, callCtx, param, v.Value)
	case *CallObjectProperty:
		zv, acl := v.GetZVal(callCtx)
		if acl != nil {
			return acl
		}
		fnCtx.SetIndexZVal(param.Index, zv)
		return nil
	case *CallStaticProperty:
		val, acl := v.GetValue(callCtx)
		if acl != nil {
			return acl
		}
		fnCtx.SetIndexZVal(param.Index, data.NewZVal(val.(data.Value)))
		return nil
	case *CallStaticPropertyLater:
		val, acl := v.GetValue(callCtx)
		if acl != nil {
			return acl
		}
		fnCtx.SetIndexZVal(param.Index, data.NewZVal(val.(data.Value)))
		return nil
	case *CallStaticKeywordProperty:
		val, acl := v.GetValue(callCtx)
		if acl != nil {
			return acl
		}
		fnCtx.SetIndexZVal(param.Index, data.NewZVal(val.(data.Value)))
		return nil
	case data.Variable:
		zv := callCtx.GetIndexZVal(v.GetIndex())
		if zv == nil {
			zv = data.NewZVal(data.NewNullValue())
			callCtx.SetIndexZVal(v.GetIndex(), zv)
		}
		fnCtx.SetIndexZVal(param.Index, zv)
		return nil
	default:
		val, acl := rawArg.GetValue(callCtx)
		if acl != nil {
			return acl
		}
		if val == nil {
			fnCtx.SetIndexZVal(param.Index, data.NewZVal(data.NewNullValue()))
			return nil
		}
		return fnCtx.SetVariableValue(param, val.(data.Value))
	}
}

// doCallWithArgs PHP 数组可调用 [$obj, 'method'](...$args) 的支持
func (pe *CallMethod) doCallWithArgs(ctx data.Context, object data.GetMethod, method data.Method) (data.GetValue, data.Control) {
	varies := method.GetVariables()
	var fnCtx data.Context
	if objCtx, ok := object.(data.Context); ok {
		fnCtx = objCtx.CreateContext(varies)
	} else if cv, ok := object.(*data.ClassValue); ok {
		fnCtx = cv.CreateContext(varies)
	} else {
		fnCtx = ctx.CreateContext(varies)
	}
	fnCtx.SetVM(ctx.GetVM())

	// 先展开所有参数中的 ...$arr (SpreadArgument)，构建展平后的实参列表
	var flatArgs []data.Value
	for _, arg := range pe.Args {
		if spread, ok := arg.(*SpreadArgument); ok {
			spreadVal, acl := spread.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if arr, ok := spreadVal.(*data.ArrayValue); ok {
				for _, z := range arr.List {
					flatArgs = append(flatArgs, z.Value)
				}
			} else if objVal, ok := spreadVal.(*data.ObjectValue); ok {
				objVal.RangeProperties(func(key string, value data.Value) bool {
					flatArgs = append(flatArgs, value)
					return true
				})
			}
			continue
		}
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			flatArgs = append(flatArgs, val)
		} else {
			flatArgs = append(flatArgs, data.NewNullValue())
		}
	}

	// 将展平的实参绑定到方法参数
	for i := 0; i < len(flatArgs) && i < len(varies); i++ {
		fnCtx.SetVariableValue(varies[i], flatArgs[i])
	}

	fnCtx.SetCallArgs(pe.Args)
	return method.Call(fnCtx)
}

// invokeMagicInvoke 调用对象的 __invoke(...$args)，用于对象作为可调用时的魔法分发
func (pe *CallMethod) invokeMagicInvoke(ctx data.Context, object data.Context, invoke data.Method) (data.GetValue, data.Control) {
	varies := invoke.GetVariables()
	fnCtx := object.CreateContext(varies)
	fnCtx.SetVM(ctx.GetVM())
	params := invoke.GetParams()

	var flatArgs []data.Value
	for _, arg := range pe.Args {
		if spread, ok := arg.(*SpreadArgument); ok {
			spreadVal, acl := spread.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if arr, ok := spreadVal.(*data.ArrayValue); ok {
				for _, z := range arr.List {
					flatArgs = append(flatArgs, z.Value)
				}
			} else if objVal, ok := spreadVal.(*data.ObjectValue); ok {
				objVal.RangeProperties(func(key string, value data.Value) bool {
					flatArgs = append(flatArgs, value)
					return true
				})
			}
			continue
		}
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			flatArgs = append(flatArgs, val)
		} else {
			flatArgs = append(flatArgs, data.NewNullValue())
		}
	}

	for index, param := range params {
		if index < len(flatArgs) {
			switch param.(type) {
			case *Parameters:
				remaining := flatArgs[index:]
				fnCtx.SetVariableValue(varies[index], data.NewArrayValue(remaining))
			default:
				fnCtx.SetVariableValue(varies[index], flatArgs[index])
			}
		} else if _, ok := param.(*Parameters); ok {
			fnCtx.SetVariableValue(varies[index], data.NewArrayValue([]data.Value{}))
		} else if argObj, ok := param.(*Parameter); ok {
			if argObj.DefaultValue == nil {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("调用 __invoke 时参数 %s 缺少值和默认值", argObj.Name))
			}
			if _, acl := argObj.GetValue(fnCtx); acl != nil {
				return nil, acl
			}
		}
	}

	fnCtx.SetCallArgs(pe.Args)
	return invoke.Call(fnCtx)
}

func (pe *CallMethod) newFunParamsError(from data.From, name string, paramName string) data.Control {
	if name == "" {
		return data.NewErrorThrow(from, errors.New("无法调用匿名函数, 缺少参数:"+paramName))
	}
	return data.NewErrorThrow(from, errors.New("无法调用("+name+")函数, 缺少参数: "+paramName))
}
