package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ReflectionMethodInvokeMethod implements ReflectionMethod::invoke.
type ReflectionMethodInvokeMethod struct{}

func (m *ReflectionMethodInvokeMethod) GetName() string { return "invoke" }

func (m *ReflectionMethodInvokeMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ReflectionMethodInvokeMethod) GetIsStatic() bool { return false }

var reflectionMethodInvokeMethodGetParams = []data.GetValue{}

func (m *ReflectionMethodInvokeMethod) GetParams() []data.GetValue {
	return reflectionMethodInvokeMethodGetParams
}

var reflectionMethodInvokeMethodGetVariables = []data.Variable{}

func (m *ReflectionMethodInvokeMethod) GetVariables() []data.Variable {
	return reflectionMethodInvokeMethodGetVariables
}

func (m *ReflectionMethodInvokeMethod) GetReturnType() data.Types { return nil }

func (m *ReflectionMethodInvokeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, methodName := "", ""
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		props := cmc.ObjectValue
		if cv, ok := props.LookupProperty("_className"); ok {
			if sv, ok := cv.(*data.StringValue); ok {
				className = sv.AsString()
			}
		}
		if mv, ok := props.LookupProperty("_methodName"); ok {
			if sv, ok := mv.(*data.StringValue); ok {
				methodName = sv.AsString()
			}
		}
	}
	if className == "" || methodName == "" {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionMethod::invoke(): no reflection info"))
	}

	vm := ctx.GetVM()

	pkgVal, acl := vm.LoadPkg(className)
	if acl != nil {
		return nil, acl
	}
	stmt, ok := pkgVal.(data.ClassStmt)
	if !ok || stmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionMethod::invoke(): class %s not found", className))
	}

	method, exists := stmt.GetMethod(methodName)
	if !exists {
		if gsm, ok := stmt.(data.GetStaticMethod); ok {
			method, exists = gsm.GetStaticMethod(methodName)
		}
	}
	if !exists {
		last := stmt
		for last.GetExtend() != nil {
			ext := last.GetExtend()
			pv, acl := vm.LoadPkg(*ext)
			if acl != nil {
				return nil, acl
			}
			ps, ok := pv.(data.ClassStmt)
			if !ok || ps == nil {
				break
			}
			method, exists = ps.GetMethod(methodName)
			if exists && method.GetModifier() != data.ModifierPrivate {
				break
			}
			if gsm, ok := ps.(data.GetStaticMethod); ok {
				method, exists = gsm.GetStaticMethod(methodName)
				if exists {
					break
				}
			}
			exists = false
			last = ps
		}
	}

	if !exists || method == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionMethod::invoke(): method %s::%s not found", className, methodName))
	}

	// ReflectionMethod 没有 PHP 形参槽。invoke($object, ...$args) 的实参
	// 在扁平参数列表里；按被调帧重放 AST 会越界，漏掉对象则会新建空实例，
	// $this 上的属性读成 null（Livewire invade()->getListeners()）。
	flat := ctx.GetFlatCallArgs()
	targetValue, hasTarget := ctx.GetIndexValue(0)
	if (!hasTarget || targetValue == nil) && len(flat) > 0 {
		val := flat[0]
		if tv, ok := val.(*data.ThisValue); ok && tv.ClassValue != nil {
			val = tv.ClassValue
		}
		if _, isNull := val.(*data.NullValue); !isNull && val != nil {
			targetValue = val
			hasTarget = true
		}
	}

	// Set up call context with proper class binding for self resolution
	varies := method.GetVariables()

	var fnCtx data.Context
	if hasTarget && targetValue != nil {
		if objCtx, ok := targetValue.(data.Context); ok {
			fnCtx = objCtx.CreateContext(varies)
			data.PreferVM(fnCtx, ctx, objCtx)
		} else {
			fnCtx = ctx.CreateContext(varies)
			data.PreferVM(fnCtx, ctx)
		}
	} else {
		// Static method call: create class context so self:: resolves correctly
		classVal := data.NewClassValue(stmt, ctx.CreateBaseContext())
		fnCtx = classVal.CreateContext(varies)
		data.PreferVM(fnCtx, ctx, classVal)
	}

	for i := 1; i < len(flat); i++ {
		argIdx := i - 1
		if argIdx >= len(varies) {
			break
		}
		if flat[i] != nil {
			fnCtx.SetVariableValue(varies[argIdx], flat[i])
		}
	}

	return method.Call(fnCtx)
}
