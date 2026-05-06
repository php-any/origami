package node

import (
	"fmt"

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
			// 若当前类未找到，再沿继承链向上查找
			extend := cls.GetExtend()
			for extend != nil {
				vm := ctx.GetVM()
				ext, acl := vm.GetOrLoadClass(*extend)
				if acl != nil {
					return nil, acl
				}
				extend = nil
				if getter, ok := ext.(data.GetStaticMethod); ok {
					if m, ok := getter.GetStaticMethod(pe.Method); ok {
						method = m
						classStmt = ext
						has = true
						break
					}
					extend = ext.GetExtend()
				}
			}
			if !has {
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
					return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("(%s)无法调用函数(%s)。", cls.GetName(), pe.Method))
				}
			}
		} else {
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
				// 若当前类未找到，再沿继承链向上查找
				extend := cls.GetExtend()
				for extend != nil {
					vm := ctx.GetVM()
					ext, acl := vm.GetOrLoadClass(*extend)
					if acl != nil {
						return nil, acl
					}
					extend = nil
					if getter, ok := ext.(data.GetStaticMethod); ok {
						if m, ok := getter.GetStaticMethod(pe.Method); ok {
							method = m
							classStmt = ext
							has = true
							break
						}
						extend = ext.GetExtend()
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
	className string // 类名（字符串形式）
	method    string // 方法名
	namespace string // 命名空间
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

// GetValue 获取延迟静态方法调用的值
func (pe *CallStaticMethodLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 尝试加载类
	stmt, acl := ctx.GetVM().GetOrLoadClass(pe.className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		// 如果还是找不到，尝试使用命名空间
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

	// 创建实际的静态方法调用
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	callStaticMethod := NewCallStaticMethod(tokenFrom, stmt, pe.method)

	return callStaticMethod.GetValue(ctx)
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
	// 创建类方法上下文，使用传入的 callCtx（包含已设置的参数），绑定当前类，保证 self:: 可用
	classValue := data.NewClassValue(s.class, callCtx)
	fnCtx := classValue.CreateContext(s.method.GetVariables())
	// 设置后期静态绑定类
	if cmc, ok := fnCtx.(*data.ClassMethodContext); ok {
		if s.callClass != nil {
			cmc.StaticClass = s.callClass
		}
	}
	if s.method.GetName() == "__callStatic" {
		// __callStatic($method, $args): 将原始参数包装为 [$methodName, [$originalArgs...]]
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
	} else {
		for i := 0; i < len(s.method.GetVariables()); i++ {
			fnCtx.SetIndexZVal(i, callCtx.GetIndexZVal(i))
		}
	}
	return s.method.Call(fnCtx)
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

	// 获取调用方传入的所有实参（使用安全的 GetIndexValue）
	callerArgs := make([]data.Value, 0)
	for i := 0; ; i++ {
		v, ok := callCtx.GetIndexValue(i)
		if !ok || v == nil {
			break
		}
		// 如果实参是 Parameters 打包的数组，展开它
		if arr, isArr := v.(*data.ArrayValue); isArr {
			for _, z := range arr.List {
				callerArgs = append(callerArgs, z.Value)
			}
		} else {
			callerArgs = append(callerArgs, v)
		}
	}

	// __callStatic($method, $args)
	vars := s.method.GetVariables()
	if len(vars) >= 2 {
		fnCtx.SetVariableValue(vars[0], data.NewStringValue(s.originalMethod))
		fnCtx.SetVariableValue(vars[1], data.NewArrayValue(callerArgs))
	}
	return s.method.Call(fnCtx)
}
