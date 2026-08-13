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

func (m *ReflectionMethodInvokeMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ReflectionMethodInvokeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ReflectionMethodInvokeMethod) GetReturnType() data.Types { return nil }

func (m *ReflectionMethodInvokeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, methodName := "", ""
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		props := cmc.ObjectValue.GetProperties()
		if cv, ok := props["_className"]; ok {
			if sv, ok := cv.(*data.StringValue); ok {
				className = sv.AsString()
			}
		}
		if mv, ok := props["_methodName"]; ok {
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

	// Set up call context with proper class binding for self resolution
	targetValue, hasTarget := ctx.GetIndexValue(0)
	varies := method.GetVariables()

	var fnCtx data.Context
	if hasTarget && targetValue != nil {
		if objCtx, ok := targetValue.(data.Context); ok {
			fnCtx = objCtx.CreateContext(varies)
		} else {
			fnCtx = ctx.CreateContext(varies)
		}
	} else {
		// Static method call: create class context so self:: resolves correctly
		classVal := data.NewClassValue(stmt, ctx.CreateBaseContext())
		fnCtx = classVal.CreateContext(varies)
	}

	callArgs := ctx.GetCallArgs()
	for i := 1; i < len(callArgs); i++ {
		argIdx := i - 1
		if argIdx >= len(varies) {
			break
		}
		v, acl := callArgs[i].GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			fnCtx.SetVariableValue(varies[argIdx], val)
		}
	}

	return method.Call(fnCtx)
}
