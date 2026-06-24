package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/net/annotation"
)

func Load(vm data.VM) {
	annotation.OnApplicationScanStart = onApplicationScanStart
	annotation.ControllerInstantiator = instantiateController

	vm.AddClass(NewContainerClass())
	vm.AddClass(NewServiceProviderClass())
	vm.AddClass(NewScopeClass())
	vm.AddClass(NewComponentClass())
	vm.AddClass(NewSingletonAnnotationClass())
	vm.AddClass(NewScopedAnnotationClass())
	vm.AddClass(NewBindClass())
	vm.AddClass(NewCircularDependencyExceptionClass())
}

func onApplicationScanStart(ctx data.Context) (func(), data.Control) {
	e := NewEngine()
	restore := setRegisteringEngine(e)
	return func() { restore() }, nil
}

func instantiateController(stmt data.ClassStmt, ctx data.Context) (data.GetValue, data.Control) {
	e := activeEngine(ctx)
	if e == nil {
		return node.InstantiateController(stmt, ctx)
	}
	obj, acl := e.Make(ctx, stmt.GetName(), nil)
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("容器 make 未返回对象实例"))
	}
	return cv, nil
}
