package node

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

// ResolveCallback validates a callback at registration and retains the resolved
// method scope. It does not execute the callback or modify normal call paths.
func ResolveCallback(ctx data.Context, value data.Value) (data.Value, data.Control) {
	invalid := func() (data.Value, data.Control) {
		return nil, data.NewTypeError(nil, fmt.Errorf("callback must be a valid callable"))
	}
	switch callback := value.(type) {
	case *data.FuncValue:
		if callback != nil && callback.Value != nil {
			return callback, nil
		}
	case *data.BoundFuncValue:
		if callback != nil && callback.Value != nil {
			return callback, nil
		}
	case *data.StringValue:
		if class, method, found := strings.Cut(callback.Value, "::"); found {
			return resolveCallbackMethod(ctx, nil, class, method)
		}
		if function, found := ctx.GetVM().GetFunc(strings.TrimPrefix(callback.Value, "\\")); found {
			return data.NewFuncValue(function), nil
		}
	case *data.ThisValue:
		return resolveCallbackMethod(ctx, callback.ClassValue, "", "__invoke")
	case *data.ClassValue:
		return resolveCallbackMethod(ctx, callback, "", "__invoke")
	case *data.ArrayValue:
		if callback.Len() != 2 {
			return invalid()
		}
		first, _ := callback.FindSlotByIntKey(0)
		second, _ := callback.FindSlotByIntKey(1)
		if first == nil || second == nil {
			return invalid()
		}
		name, ok := second.Value.(*data.StringValue)
		if !ok {
			return invalid()
		}
		switch receiver := first.Value.(type) {
		case *data.ClassValue:
			return resolveCallbackMethod(ctx, receiver, "", name.Value)
		case *data.ThisValue:
			return resolveCallbackMethod(ctx, receiver.ClassValue, "", name.Value)
		case *data.StringValue:
			return resolveCallbackMethod(ctx, nil, receiver.Value, name.Value)
		}
	}
	return invalid()
}

func resolveCallbackMethod(ctx data.Context, receiver *data.ClassValue, name, methodName string) (data.Value, data.Control) {
	var class data.ClassStmt
	if receiver != nil {
		class = receiver.Class
	} else {
		name = strings.TrimPrefix(name, "\\")
		class, _ = ctx.GetVM().GetClass(name)
		if class == nil {
			// LoadPkg is a probe: an autoloader may decline a symbol without
			// manufacturing an Error before callback validation reports TypeError.
			loaded, ctl := ctx.GetVM().LoadPkg(name)
			if ctl != nil {
				return nil, ctl
			}
			class, _ = loaded.(data.ClassStmt)
		}
	}
	if class == nil {
		return nil, data.NewTypeError(nil, fmt.Errorf("invalid callback class"))
	}
	var method data.Method
	var found bool
	if receiver != nil {
		method, found = receiver.GetMethod(methodName)
	} else {
		method, found = class.GetMethod(methodName)
		if !found {
			if lookup, ok := class.(data.GetStaticMethod); ok {
				method, found = lookup.GetStaticMethod(methodName)
			}
		}
	}
	if !found || method == nil || (receiver == nil && !method.GetIsStatic()) {
		return nil, data.NewTypeError(nil, fmt.Errorf("invalid callback %s::%s", class.GetName(), methodName))
	}
	if _, abstract := method.(*AbstractMethod); abstract {
		return nil, data.NewTypeError(nil, fmt.Errorf("abstract callback %s::%s", class.GetName(), methodName))
	}
	self := findDeclaringClassForMethod(ctx.GetVM(), class, methodName)
	if method.GetModifier() != data.ModifierPublic {
		var caller data.ClassStmt
		if owner, ok := ctx.(*data.ClassMethodContext); ok {
			caller = owner.SelfClass
		}
		if bound := data.FindBoundContext(ctx); bound != nil && bound.ScopeClass != "" {
			caller, _ = ctx.GetVM().GetClass(bound.ScopeClass)
		}
		allowed := caller != nil && self != nil && data.SameNominalClass(caller, self.GetName(), ctx.GetVM())
		if method.GetModifier() == data.ModifierProtected && caller != nil && self != nil {
			allowed = allowed || data.NominalIsA(caller, self.GetName(), ctx.GetVM()) || data.NominalIsA(self, caller.GetName(), ctx.GetVM())
		}
		if !allowed {
			return nil, data.NewTypeError(nil, fmt.Errorf("inaccessible callback %s::%s", class.GetName(), methodName))
		}
	}
	return data.NewFuncValue(&resolvedCallbackMethod{Method: method, receiver: receiver, self: self, static: class}), nil
}

type resolvedCallbackMethod struct {
	data.Method
	receiver     *data.ClassValue
	self, static data.ClassStmt
}

func (f *resolvedCallbackMethod) ParameterTypeContext(ctx data.Context) data.Context {
	if f.receiver == nil || f.Method.GetIsStatic() {
		return data.NewStaticMethodContext(ctx, f.self, f.static)
	}
	return data.WrapMethodFrame(ctx, f.receiver, f.self, f.static)
}

func (f *resolvedCallbackMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return f.Method.Call(f.ParameterTypeContext(ctx))
}
