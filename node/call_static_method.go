package node

import (
	"fmt"
	"sync"

	"github.com/php-any/origami/data"
)

// CallStaticMethod 表示对象属性访问表达式
type CallStaticMethod struct {
	*Node  `pp:"-"`
	stmt   data.GetValue // 类名称 Class::fn() or Class::test::one
	Method string        // 函数名
}

func NewCallStaticMethod(from *TokenFrom, path data.GetValue, method string) *CallStaticMethod {
	return &CallStaticMethod{
		Node:   NewNode(from),
		stmt:   path,
		Method: method,
	}
}

// GetStmt 获取类引用表达式
func (pe *CallStaticMethod) GetStmt() data.GetValue {
	return pe.stmt
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var method data.Method
	var classStmt data.ClassStmt // 方法定义所在的类
	var callClass data.ClassStmt // 实际调用的类（用于 late static binding）
	var has bool

	switch expr := pe.stmt.(type) {
	case data.GetStaticMethod:
		// 先在当前类上查找静态方法
		method, has = expr.GetStaticMethod(pe.Method)
		if has {
			if cls, ok := expr.(data.ClassStmt); ok {
				classStmt = cls
				callClass = cls
			}
		} else if cls, ok := expr.(data.ClassStmt); ok {
			callClass = cls // 记录原始调用类
			// 若当前类未找到，再沿继承链向上查找。
			// 父类即使未实现 GetStaticMethod（原生基类）也要继续往上，否则
			// BackedEnum::tryFrom 这类挂在祖先上的静态方法会找不到。
			extend := cls.GetExtend()
			for extend != nil {
				vm := ctx.GetVM()
				if vm == nil {
					break
				}
				ext, acl := vm.GetOrLoadClass(*extend)
				if acl != nil {
					return nil, acl
				}
				if ext == nil {
					break
				}
				extend = ext.GetExtend()
				if getter, ok := ext.(data.GetStaticMethod); ok {
					if m, ok := getter.GetStaticMethod(pe.Method); ok {
						method = m
						classStmt = ext
						has = true
						break
					}
				}
			}
			if !has {
				// PHP: 实例方法上下文中可用 self::foo() / ClassName::foo() 调用非静态方法
				if objCtx, ok := ctx.(*data.ClassMethodContext); ok && objCtx.ClassValue != nil && objCtx.ObjectValue != nil {
					checkClass := cls
					for checkClass != nil {
						if m, ok := checkClass.GetMethod(pe.Method); ok {
							return data.NewFuncValue(&instanceViaSelfFunc{
								this:   objCtx.ClassValue,
								method: m,
							}), nil
						}
						if checkClass.GetExtend() == nil {
							break
						}
						vm := ctx.GetVM()
						parent, acl := vm.GetOrLoadClass(*checkClass.GetExtend())
						if acl != nil || parent == nil {
							break
						}
						checkClass = parent
					}
				}

				// 检查 __callStatic 魔术方法（包括父类）
				checkClass := cls
				for checkClass != nil {
					if getter, ok := checkClass.(data.GetStaticMethod); ok {
						if magic, hasMagic := getter.GetStaticMethod("__callStatic"); hasMagic {
							method = magic
							classStmt = cls // use original class for static
							has = true
							break
						}
					}
					if checkClass.GetExtend() == nil {
						break
					}
					vm := ctx.GetVM()
					parent, acl := vm.GetOrLoadClass(*checkClass.GetExtend())
					if acl != nil || parent == nil {
						break
					}
					checkClass = parent
				}
				if !has {
					if fn, ok := tryNewInstanceMagicCallViaStaticFunc(ctx, pe.Method, cls); ok {
						return data.NewFuncValue(fn), nil
					}
					return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("(%s)无法调用函数(%s)。", cls.GetName(), pe.Method))
				}
			}
		} else {
			if fn, ok := tryNewInstanceMagicCallViaStaticFunc(ctx, pe.Method, nil); ok {
				return data.NewFuncValue(fn), nil
			}
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法调用函数(%s)。", pe.Method))
		}
	default:
		c, acl := pe.stmt.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		switch expr := c.(type) {
		case *data.StringValue:
			// 动态字符串类名：$className::method()
			className := expr.Value
			vm := ctx.GetVM()
			stmt, acl := vm.GetOrLoadClass(className)
			if acl != nil {
				return nil, acl
			}
			if stmt == nil {
				return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法调用静态方法(%s::%s), 未找到类", className, pe.Method))
			}
			if tokenFrom, ok := pe.GetFrom().(*TokenFrom); ok {
				callStaticMethod := NewCallStaticMethod(tokenFrom, stmt, pe.Method)
				return callStaticMethod.GetValue(ctx)
			}
			// fallback: 直接查找静态方法
			if getter, ok := stmt.(data.GetStaticMethod); ok {
				if m, ok := getter.GetStaticMethod(pe.Method); ok {
					classStmt = stmt
					callClass = stmt
					method = m
					has = true
				}
			}
		case *data.ClassValue:
			// $instance::staticMethod()：后期静态绑定必须用实例的类，不能沿用调用方 ClassMethodContext。
			// Livewire：$synth::getKey() 在 HandleSynths 内调用时 static::class 应为 Synth 子类。
			callClass = expr.Class
			if gsm, ok := expr.Class.(data.GetStaticMethod); ok {
				method, has = gsm.GetStaticMethod(pe.Method)
			}
			if !has {
				method, has = expr.GetMethod(pe.Method)
			}
			if has {
				classStmt = findDeclaringClassForMethod(ctx.GetVM(), expr.Class, pe.Method)
				if classStmt == nil {
					classStmt = expr.Class
				}
			}
		case data.GetStaticMethod:
			// 先在当前类上查找静态方法
			method, has = expr.GetStaticMethod(pe.Method)
			if has {
				if cls, ok := expr.(data.ClassStmt); ok {
					classStmt = cls
					callClass = cls
				}
			} else if cls, ok := expr.(data.ClassStmt); ok {
				callClass = cls
				// 若当前类未找到，再沿继承链向上查找（原生父类可没有 GetStaticMethod）
				extend := cls.GetExtend()
				for extend != nil {
					vm := ctx.GetVM()
					if vm == nil {
						break
					}
					ext, acl := vm.GetOrLoadClass(*extend)
					if acl != nil {
						return nil, acl
					}
					if ext == nil {
						break
					}
					extend = ext.GetExtend()
					if getter, ok := ext.(data.GetStaticMethod); ok {
						if m, ok := getter.GetStaticMethod(pe.Method); ok {
							method = m
							classStmt = ext
							has = true
							break
						}
					}
				}
			}
		case data.GetMethod:
			method, has = expr.GetMethod(pe.Method)
			if has {
				// 实例方法，直接返回 FuncValue
				return data.NewFuncValue(method), nil
			}
		}
	}

	if !has {
		var staticCls data.ClassStmt
		if cls, ok := pe.stmt.(data.ClassStmt); ok {
			staticCls = cls
		}
		if fn, ok := tryNewInstanceMagicCallViaStaticFunc(ctx, pe.Method, staticCls); ok {
			return data.NewFuncValue(fn), nil
		}
		name := ""
		if getName, ok := pe.stmt.(data.ClassStmt); ok {
			name = getName.GetName()
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("(%v)无法调用函数(%s)。", name, pe.Method))
	}

	// 静态方法需要 ClassMethodContext，返回包装器让 CallMethod 正确处理
	if classStmt != nil {
		// __callStatic 需要特殊处理：调用方传入的实参需要重打包为 [methodName, args]
		if method.GetName() == "__callStatic" {
			return data.NewFuncValue(&callStaticFunc{
				class:          classStmt,
				method:         method,
				originalMethod: pe.Method,
			}), nil
		}
		return data.NewFuncValue(&staticMethodFunc{
			class:          classStmt,
			callClass:      callClass,
			method:         method,
			originalMethod: pe.Method,
		}), nil
	}

	// 如果没有类信息，直接返回 FuncValue（向后兼容）
	return data.NewFuncValue(method), nil
}

// CallStaticMethodLater 延迟的静态方法调用（类未加载时）
type CallStaticMethodLater struct {
	*Node
	className string            // 类名（字符串形式）
	method    string            // 方法名
	namespace string            // 命名空间
	call      *CallStaticMethod `pp:"-"` // 解析后缓存
	resolveMu sync.Mutex
}

// NewCallStaticMethodLater 创建延迟的静态方法调用
func NewCallStaticMethodLater(from *TokenFrom, className, method, namespace string) *CallStaticMethodLater {
	return &CallStaticMethodLater{
		Node:      NewNode(from),
		className: className,
		method:    method,
		namespace: namespace,
	}
}

func (pe *CallStaticMethodLater) resolveCall(ctx data.Context) (*CallStaticMethod, data.Control) {
	pe.resolveMu.Lock()
	defer pe.resolveMu.Unlock()
	if pe.call != nil {
		return pe.call, nil
	}
	stmt, acl := ctx.GetVM().GetOrLoadClass(pe.className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		fullClassName := pe.className
		if pe.namespace != "" {
			fullClassName = pe.namespace + "\\" + pe.className
		}
		stmt, acl = ctx.GetVM().GetOrLoadClass(fullClassName)
		if acl != nil {
			return nil, acl
		}
		if stmt == nil {
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法调用静态方法(%s::%s), 未找到类", pe.className, pe.method))
		}
	}
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	pe.call = NewCallStaticMethod(tokenFrom, stmt, pe.method)
	return pe.call, nil
}

// GetValue 获取延迟静态方法调用的值
func (pe *CallStaticMethodLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	call, acl := pe.resolveCall(ctx)
	if acl != nil {
		return nil, acl
	}
	return call.GetValue(ctx)
}

func NewStaticMethodFuncValue(class data.ClassStmt, method data.Method) *StaticMethodFuncValue {
	return &StaticMethodFuncValue{
		class:  class,
		method: method,
	}
}

// StaticMethodFuncValue 静态方法函数值包装器，确保调用时使用 ClassMethodContext
type StaticMethodFuncValue struct {
	class  data.ClassStmt
	method data.Method
}

func (s *StaticMethodFuncValue) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 返回 FuncValue，但内部使用 staticMethodFunc 包装，确保调用时使用 ClassMethodContext
	return data.NewFuncValue(&staticMethodFunc{class: s.class, method: s.method}), nil
}

// staticMethodFunc 适配器：将 data.Method 包装为 data.FuncStmt，并在调用时切换到 ClassMethodContext
type staticMethodFunc struct {
	class          data.ClassStmt // 方法定义所在的类
	callClass      data.ClassStmt // 实际调用的类（用于 late static binding）
	method         data.Method
	originalMethod string // 用于 __callStatic 时保存原始方法名
}

func (s *staticMethodFunc) GetName() string               { return s.method.GetName() }
func (s *staticMethodFunc) GetParams() []data.GetValue    { return s.method.GetParams() }
func (s *staticMethodFunc) GetVariables() []data.Variable { return s.method.GetVariables() }
func (s *staticMethodFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	if s.method.GetName() == "__callStatic" {
		return s.callStatic(callCtx)
	}
	return s.method.Call(data.NewStaticMethodContext(callCtx, s.class, s.callClass))
}

// callStatic 走新建帧：__callStatic 需要把实参重排成 [$name, [$args]]。
func (s *staticMethodFunc) callStatic(callCtx data.Context) (data.GetValue, data.Control) {
	classValue := data.NewClassValue(s.class, callCtx)
	fnCtx := classValue.CreateContext(s.method.GetVariables())
	if cmc, ok := fnCtx.(*data.ClassMethodContext); ok {
		if s.callClass != nil {
			cmc.StaticClass = s.callClass
		}
		cmc.SelfClass = s.class
	}
	vars := s.method.GetVariables()
	if len(vars) >= 2 {
		fnCtx.SetVariableValue(vars[0], data.NewStringValue(s.originalMethod))
		argList := make([]data.Value, 0)
		for i := 0; ; i++ {
			zv, ok := callCtx.GetIndexValue(i)
			if !ok || zv == nil {
				break
			}
			argList = append(argList, zv)
		}
		fnCtx.SetVariableValue(vars[1], data.NewArrayValue(argList))
	}
	return s.method.Call(fnCtx)
}

// instanceViaSelfFunc 支持实例方法上下文中的 self::nonStaticMethod() 调用
type instanceViaSelfFunc struct {
	this   *data.ClassValue
	method data.Method
}

func (s *instanceViaSelfFunc) GetName() string               { return s.method.GetName() }
func (s *instanceViaSelfFunc) GetParams() []data.GetValue    { return s.method.GetParams() }
func (s *instanceViaSelfFunc) GetVariables() []data.Variable { return s.method.GetVariables() }
func (s *instanceViaSelfFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	cv := s.this.CloneWithContext(callCtx)
	cmc := &data.ClassMethodContext{ClassValue: cv}
	return s.method.Call(cmc)
}

// callStaticFunc 专门用于 __callStatic，将调用方实参重打包为 [methodName, args]
type callStaticFunc struct {
	class          data.ClassStmt
	method         data.Method
	originalMethod string
}

func (s *callStaticFunc) GetName() string { return s.method.GetName() }
func (s *callStaticFunc) GetParams() []data.GetValue {
	// 接受任意数量的任意参数
	return []data.GetValue{NewParametersNoName(0)}
}
func (s *callStaticFunc) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("args", 0, nil)}
}
func (s *callStaticFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	classValue := data.NewClassValue(s.class, callCtx)
	fnCtx := classValue.CreateContext(s.method.GetVariables())

	// __callStatic($method, $args)：保留命名实参键，供 Facade 转发 $instance->$method(...$args)
	vars := s.method.GetVariables()
	if len(vars) >= 2 {
		fnCtx.SetVariableValue(vars[0], data.NewStringValue(s.originalMethod))
		if v, ok := callCtx.GetIndexValue(0); ok && v != nil {
			if arr, isArr := v.(*data.ArrayValue); isArr {
				fnCtx.SetVariableValue(vars[1], arr)
			} else {
				fnCtx.SetVariableValue(vars[1], data.NewArrayValue([]data.Value{v}))
			}
		} else {
			fnCtx.SetVariableValue(vars[1], data.NewArrayValue(nil))
		}
	}
	return s.method.Call(fnCtx)
}
