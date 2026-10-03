package node

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

func init() { data.CallableResolver = ResolveCallback }

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
		name, ok := second.ReadValue().(*data.StringValue)
		if !ok {
			return invalid()
		}
		switch receiver := first.ReadValue().(type) {
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
	var lexical data.ClassStmt
	if receiver != nil {
		method, lexical, found = lookupObjectMethod(ctx, receiver, methodName)
	} else {
		method, found = class.GetMethod(methodName)
		if !found {
			if lookup, ok := class.(data.GetStaticMethod); ok {
				method, found = lookup.GetStaticMethod(methodName)
			}
		}
	}
	self := findDeclaringClassForMethod(ctx.GetVM(), class, methodName)
	if lexical != nil {
		self = lexical
	}
	if !found || method == nil || !callbackMethodAccessible(ctx, method, self) {
		magicName := "__callStatic"
		if receiver != nil {
			magicName = "__call"
		}
		var magic data.Method
		if receiver != nil {
			magic, _ = receiver.GetMethod(magicName)
		} else {
			magic, _ = class.GetMethod(magicName)
			if magic == nil {
				if lookup, ok := class.(data.GetStaticMethod); ok {
					magic, _ = lookup.GetStaticMethod(magicName)
				}
			}
		}
		if magic != nil && magic.GetModifier() == data.ModifierPublic && (receiver != nil || magic.GetIsStatic()) {
			base := &resolvedCallbackMethod{Method: magic, receiver: receiver, self: findDeclaringClassForMethod(ctx.GetVM(), class, magicName), static: class}
			return data.NewFuncValue(&resolvedMagicCallback{resolvedCallbackMethod: base, name: methodName}), nil
		}
		return nil, data.NewTypeError(nil, fmt.Errorf("invalid callback %s::%s", class.GetName(), methodName))
	}
	if receiver == nil && !method.GetIsStatic() {
		return nil, data.NewTypeError(nil, fmt.Errorf("non-static callback %s::%s", class.GetName(), methodName))
	}
	if _, abstract := method.(*AbstractMethod); abstract {
		return nil, data.NewTypeError(nil, fmt.Errorf("abstract callback %s::%s", class.GetName(), methodName))
	}
	return data.NewFuncValue(&resolvedCallbackMethod{Method: method, receiver: receiver, self: self, static: class}), nil
}

func callbackMethodAccessible(ctx data.Context, method data.Method, self data.ClassStmt) bool {
	if method == nil {
		return false
	}
	return memberAccessible(ctx, method.GetModifier(), self)
}

func memberAccessible(ctx data.Context, modifier data.Modifier, self data.ClassStmt) bool {
	if modifier != data.ModifierPublic {
		var caller data.ClassStmt
		if owner, ok := ctx.(*data.ClassMethodContext); ok {
			caller = owner.SelfClass
			if caller == nil {
				caller = owner.Class
			}
		}
		if owner, ok := ctx.(*data.ClassValue); ok {
			caller = owner.Class
		}
		if bound, ok := ctx.(*data.BoundContext); caller == nil && ok && bound.ScopeClass != "" {
			caller, _ = ctx.GetVM().GetClass(bound.ScopeClass)
		}
		if caller == nil {
			// Native functions have their own argument frame but do not create
			// a PHP lexical scope. The immediate PHP frame supplies that scope;
			// stop at ordinary functions instead of inheriting an older method.
			if recorder, ok := ctx.(data.CallStackTracker); ok {
				frames := recorder.SnapshotCallStack()
				if len(frames) != 0 {
					frame := frames[len(frames)-1]
					if frame.Class != "" {
						caller, _ = ctx.GetVM().GetClass(frame.Class)
						if caller != nil {
							caller = findDeclaringClassForMethod(ctx.GetVM(), caller, frame.Function)
						}
					}
				}
			}
		}
		allowed := caller != nil && self != nil && data.SameNominalClass(caller, self.GetName(), ctx.GetVM())
		if modifier == data.ModifierProtected && caller != nil && self != nil {
			allowed = allowed || data.NominalIsA(caller, self.GetName(), ctx.GetVM()) || data.NominalIsA(self, caller.GetName(), ctx.GetVM())
		}
		if !allowed {
			return false
		}
	}
	return true
}

type resolvedMagicCallback struct {
	*resolvedCallbackMethod
	name string
}

var magicCallbackParams = []data.GetValue{NewParameters(nil, "arguments", 0, nil, data.TypeMixed)}
var magicCallbackVariables = []data.Variable{NewVariable(nil, "arguments", 0, data.TypeMixed)}

func (f *resolvedMagicCallback) GetParams() []data.GetValue    { return magicCallbackParams }
func (f *resolvedMagicCallback) GetVariables() []data.Variable { return magicCallbackVariables }
func (f *resolvedMagicCallback) GetName() string               { return f.name }
func (f *resolvedMagicCallback) Call(ctx data.Context) (data.GetValue, data.Control) {
	arguments, _ := ctx.GetIndexValue(0)
	if arguments == nil {
		arguments = data.NewArrayValue(nil)
	}
	frame := f.ParameterTypeContext(ctx.CreateContext(f.Method.GetVariables()))
	defer tryReleaseCallContext(f.Method, frame)
	if ctl := data.BindDeclaredArgs(frame, f.Method, []data.Value{data.NewStringValue(f.name), arguments}); ctl != nil {
		return nil, ctl
	}
	return f.Method.Call(frame)
}
func (f *resolvedMagicCallback) BindRequestCapture(scope *data.RequestCaptureScope) data.FuncStmt {
	clone := *f
	clone.resolvedCallbackMethod = f.resolvedCallbackMethod.BindRequestCapture(scope).(*resolvedCallbackMethod)
	return &clone
}
func (f *resolvedMagicCallback) BindRequestScope(ctx data.Context, objects map[*data.ObjectValue]*data.ClassValue) data.FuncStmt {
	clone := *f
	clone.resolvedCallbackMethod = f.resolvedCallbackMethod.BindRequestScope(ctx, objects).(*resolvedCallbackMethod)
	return &clone
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
	frame := f.ParameterTypeContext(ctx).(*data.ClassMethodContext)
	defer frame.ReleaseBorrowedFrame()
	return f.Method.Call(frame)
}

func (f *resolvedCallbackMethod) RequestScopeObjects() []*data.ClassValue {
	if f.receiver == nil {
		return nil
	}
	return []*data.ClassValue{f.receiver}
}
func (f *resolvedCallbackMethod) BindRequestScope(ctx data.Context, objects map[*data.ObjectValue]*data.ClassValue) data.FuncStmt {
	clone := *f
	if f.receiver != nil {
		if target := objects[f.receiver.ObjectValue]; target != nil {
			clone.receiver = target
		}
	}
	return &clone
}
func (f *resolvedCallbackMethod) BindRequestCapture(scope *data.RequestCaptureScope) data.FuncStmt {
	clone := *f
	if f.receiver != nil {
		if scope.ScopeObject != nil {
			clone.receiver = scope.ScopeObject(f.receiver)
		} else if target := scope.Objects[f.receiver.ObjectValue]; target != nil {
			clone.receiver = target
		}
	}
	return &clone
}
