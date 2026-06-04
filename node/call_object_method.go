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

			fnCtx.SetCallArgs(pe.Args)
			ret, acl := method.Call(fnCtx)
			return pe.wrapMethodCallResult(class.ClassValue, ret, acl)
		}
		// 方法未找到时尝试魔法方法 __call(string $name, array $arguments)
		if magic, hasCall := class.GetMethod("__call"); hasCall {
			return pe.invokeMagicCall(class, ctx, magic, pe.Method, pe.Args)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("this 对象不存在对应函数: "+pe.Method))
	case *data.ClassValue:
		method, has := class.GetMethod(pe.Method)
		if has {
			if method.GetModifier() == data.ModifierPrivate {
				if !isCallerInClassHierarchy(ctx, class.Class) {
					return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("不能调用 private 方法: "+pe.Method))
				}
			} else if method.GetModifier() == data.ModifierProtected {
				if !isCallerInClassHierarchy(ctx, class.Class) {
					return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("对象属性访问表达式对象属性访问函数非公开"))
				}
			}

			fnCtx, acl := pe.callMethodParams(class, ctx, method)
			if acl != nil {
				if _, ok := acl.(ToClosure); ok {
					return data.NewFuncValue(method), nil
				}
				return nil, acl
			}

			fnCtx.SetCallArgs(pe.Args)
			ret, acl := method.Call(fnCtx)
			return pe.wrapMethodCallResult(class, ret, acl)
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

				fnCtx.SetCallArgs(pe.Args)
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

func (pe *CallObjectMethod) wrapMethodCallResult(class *data.ClassValue, ret data.GetValue, acl data.Control) (data.GetValue, data.Control) {
	if acl != nil {
		if tv, ok := acl.(*data.ThrowValue); ok && tv.PHPUncaughtError {
			tv.AddStackWithInfo(pe.GetFrom(), class.Class.GetName(), pe.Method)
		}
	}
	return ret, acl
}

// magicCallArgValue 为 __call 的 $arguments 数组准备实参（数组需 COW 克隆，保留元素级引用）
func magicCallArgValue(val data.Value) data.Value {
	if arr, ok := val.(*data.ArrayValue); ok {
		return data.CloneArrayValueForCallArgs(arr)
	}
	return val
}

// invokeMagicCall 调用 __call(string $name, array $arguments)，用于未定义方法时的魔法分发
func (pe *CallObjectMethod) invokeMagicCall(object data.Context, ctx data.Context, magic data.Method, methodName string, args []data.GetValue) (data.GetValue, data.Control) {
	var argsList []data.Value
	for _, arg := range args {
		// 展开 ...$arr (SpreadArgument)
		if spread, ok := arg.(*SpreadArgument); ok {
			spreadVal, acl := spread.GetValue(ctx)
			if acl != nil {
				if _, isToClosure := acl.(ToClosure); isToClosure {
					continue
				}
				return nil, acl
			}
			if arr, ok := spreadVal.(*data.ArrayValue); ok {
				for _, z := range arr.List {
					argsList = append(argsList, magicCallArgValue(z.Value))
				}
			} else if objVal, ok := spreadVal.(*data.ObjectValue); ok {
				objVal.RangeProperties(func(key string, value data.Value) bool {
					argsList = append(argsList, magicCallArgValue(value))
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
			argsList = append(argsList, magicCallArgValue(val))
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
	params := method.GetParams()

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
	for index, param := range params {
		if index < len(flatArgs) {
			var acl data.Control
			switch p := param.(type) {
			case *Parameter:
				fnCtx.SetVariableValue(varies[index], flatArgs[index])
			case *ParameterReference:
				// 引用参数：直接设置值
				fnCtx.SetVariableValue(varies[index], flatArgs[index])
			case *Parameters:
				// 可变参数：收集剩余的所有实参
				remaining := flatArgs[index:]
				arr := data.NewArrayValue(remaining)
				fnCtx.SetVariableValue(varies[index], arr)
				index = len(params) // 跳过后续参数
			case *PromotedParameter:
				fnCtx.SetVariableValue(varies[index], flatArgs[index])
				acl = p.SetValue(object, flatArgs[index])
			default:
				fnCtx.SetVariableValue(varies[index], flatArgs[index])
			}
			if acl != nil {
				return nil, acl
			}
		} else {
			// 实参不足
			if pVar, ok := param.(*Parameters); ok {
				// Variadic 带 0 实参 → 空数组（PHP 语义）
				arr := data.NewArrayValue([]data.Value{})
				fnCtx.SetVariableValue(pVar, arr)
			} else if promotedParam, ok := param.(*PromotedParameter); ok {
				_, acl := promotedParam.GetValue(object)
				if acl != nil {
					return nil, acl
				}
			} else if argObj, ok := param.(*Parameter); ok {
				if argObj.DefaultValue == nil {
					return nil, data.NewErrorThrow(pe.from, fmt.Errorf("调用 %s 构造函数时参数 %s 缺少值和默认值", object, argObj.Name))
				}
				_, acl := argObj.GetValue(fnCtx)
				if acl != nil {
					return nil, acl
				}
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

// isCallerInClassHierarchy 检查调用者是否在目标类的类层次结构中
// 用于确定是否允许调用 protected 方法
func isCallerInClassHierarchy(ctx data.Context, targetClass data.ClassStmt) bool {
	// 检查是否通过 Closure::bind() 绑定了作用域（允许访问私有成员）
	if bc, ok := ctx.(*data.BoundContext); ok {
		if bc.ScopeClass == targetClass.GetName() {
			return true
		}
	}

	// 从上下文链中查找 ClassMethodContext 或 ClassValue
	var callerClass data.ClassStmt
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		callerClass = cmc.Class
	} else if cv, ok := ctx.(*data.ClassValue); ok {
		callerClass = cv.Class
	} else {
		return false
	}

	// 检查调用者类是否与目标类相同
	if callerClass.GetName() == targetClass.GetName() {
		return true
	}

	vm := ctx.GetVM()

	// 检查调用者类是否是目标类的子类（调用者继承自目标）
	extend := callerClass.GetExtend()
	for extend != nil {
		if *extend == targetClass.GetName() {
			return true
		}
		cls, acl := vm.GetOrLoadClass(*extend)
		if acl != nil {
			return false
		}
		extend = cls.GetExtend()
	}

	// 检查目标类是否是调用者类的子类（目标继承自调用者）
	// 父类可以访问子类实例上的 protected 属性
	targetExtend := targetClass.GetExtend()
	for targetExtend != nil {
		if *targetExtend == callerClass.GetName() {
			return true
		}
		cls, acl := vm.GetOrLoadClass(*targetExtend)
		if acl != nil {
			return false
		}
		targetExtend = cls.GetExtend()
	}

	return false
}

// findClassMethodContext 从上下文链中查找含 $this 的类方法上下文
func findClassMethodContext(ctx data.Context) *data.ClassMethodContext {
	for c := ctx; c != nil; {
		switch v := c.(type) {
		case *data.ClassMethodContext:
			return v
		case *data.ClassValue:
			if v.Context != nil {
				c = v.Context
				continue
			}
			return nil
		default:
			return nil
		}
	}
	return nil
}

// findMagicCallOnClass 沿继承链查找实例 __call
func findMagicCallOnClass(class data.ClassStmt, vm data.VM) (data.Method, bool) {
	for class != nil {
		if m, ok := class.GetMethod("__call"); ok {
			return m, true
		}
		if class.GetExtend() == nil {
			break
		}
		parent, acl := vm.GetOrLoadClass(*class.GetExtend())
		if acl != nil || parent == nil {
			break
		}
		class = parent
	}
	return nil, false
}

// tryNewInstanceMagicCallViaStaticFunc 在实例方法内 Class::missing() 时转调当前作用域的 __call($this)
func tryNewInstanceMagicCallViaStaticFunc(ctx data.Context, methodName string) (data.FuncStmt, bool) {
	cmc := findClassMethodContext(ctx)
	if cmc == nil {
		return nil, false
	}
	magic, ok := findMagicCallOnClass(cmc.Class, ctx.GetVM())
	if !ok {
		return nil, false
	}
	return &instanceMagicCallViaStaticFunc{
		objectCtx:      cmc,
		magic:          magic,
		originalMethod: methodName,
	}, true
}

// instanceMagicCallViaStaticFunc 将 A::foo() 形式的未定义静态调用转为 $this->__call('foo', $args)
type instanceMagicCallViaStaticFunc struct {
	objectCtx      *data.ClassMethodContext
	magic          data.Method
	originalMethod string
}

func (s *instanceMagicCallViaStaticFunc) GetName() string { return s.magic.GetName() }
func (s *instanceMagicCallViaStaticFunc) GetParams() []data.GetValue {
	return []data.GetValue{NewParametersNoName(0)}
}
func (s *instanceMagicCallViaStaticFunc) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("args", 0, nil)}
}

// objectMethodCallable 将 [$obj, 'method'] 转为可调用（支持 __call）
type objectMethodCallable struct {
	obj    *data.ClassValue
	method string
}

// NewObjectMethodCallable 创建实例方法回调（供 call_user_func 等使用）
func NewObjectMethodCallable(obj *data.ClassValue, method string) data.FuncStmt {
	return &objectMethodCallable{obj: obj, method: method}
}

func (o *objectMethodCallable) GetName() string { return o.method }
func (o *objectMethodCallable) GetParams() []data.GetValue {
	return []data.GetValue{NewParametersNoName(0)}
}
func (o *objectMethodCallable) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("args", 0, nil)}
}

func (o *objectMethodCallable) Call(callCtx data.Context) (data.GetValue, data.Control) {
	proxy := &CallObjectMethod{Object: o.obj, Method: o.method}
	if method, has := o.obj.GetMethod(o.method); has {
		fnCtx, acl := proxy.callMethodParams(o.obj, callCtx, method)
		if acl != nil {
			return nil, acl
		}
		return method.Call(fnCtx)
	}
	if magic, has := o.obj.GetMethod("__call"); has {
		return proxy.invokeMagicCallFromCallCtx(o.obj, callCtx, magic, o.method)
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("未找到方法 %s", o.method))
}

// invokeMagicCallFromCallCtx 从调用上下文收集实参并调用 __call
func (pe *CallObjectMethod) invokeMagicCallFromCallCtx(object data.Context, callCtx data.Context, magic data.Method, methodName string) (data.GetValue, data.Control) {
	var argsList []data.Value
	for i := 0; ; i++ {
		v, ok := callCtx.GetIndexValue(i)
		if !ok || v == nil {
			break
		}
		argsList = append(argsList, magicCallArgValue(v))
	}
	varies := magic.GetVariables()
	if len(varies) < 2 {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__call 需要至少 2 个参数"))
	}
	fnCtx := object.CreateContext(varies)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(methodName))
	fnCtx.SetVariableValue(varies[1], data.NewArrayValue(argsList))
	return magic.Call(fnCtx)
}

func (s *instanceMagicCallViaStaticFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	callerArgs := make([]data.Value, 0)
	for i := 0; ; i++ {
		v, ok := callCtx.GetIndexValue(i)
		if !ok || v == nil {
			break
		}
		if arr, isArr := v.(*data.ArrayValue); isArr {
			for _, z := range arr.List {
				callerArgs = append(callerArgs, magicCallArgValue(z.Value))
			}
		} else {
			callerArgs = append(callerArgs, v)
		}
	}
	varies := s.magic.GetVariables()
	if len(varies) < 2 {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("__call 需要至少 2 个参数"))
	}
	fnCtx := s.objectCtx.CreateContext(varies)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(s.originalMethod))
	fnCtx.SetVariableValue(varies[1], data.NewArrayValue(callerArgs))
	return s.magic.Call(fnCtx)
}
