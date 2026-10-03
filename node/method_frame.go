package node

import "github.com/php-any/origami/data"

// implicitMethodFrame preserves the caller's VM and deadline for language
// hooks (ArrayAccess, magic properties and iteration). The receiver's saved
// context belongs to its construction, which may have happened during boot.
func implicitMethodFrame(caller data.Context, receiver *data.ClassValue, method data.Method) data.Context {
	self := findDeclaringClassForMethod(caller.GetVM(), receiver.Class, method.GetName())
	return data.WrapMethodFrame(caller.CreateContext(method.GetVariables()), receiver, self, receiver.Class)
}

func implicitReceiverFrame(caller, receiver data.Context, method data.Method) data.Context {
	switch object := receiver.(type) {
	case *data.ClassValue:
		return implicitMethodFrame(caller, object, method)
	case *data.ThisValue:
		return implicitMethodFrame(caller, object.ClassValue, method)
	default:
		return receiver.CreateContext(method.GetVariables())
	}
}
