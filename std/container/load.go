package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	containerannotation "github.com/php-any/origami/std/container/annotation"
	"github.com/php-any/origami/std/net/annotation"
	netdata "github.com/php-any/origami/std/net/data"
)

func Load(vm data.VM) {
	annotation.OnApplicationScanStart = onApplicationScanStart
	containerannotation.RegisterClassLifetime = func(ctx data.Context, lifetime int) data.Control {
		return RegisterClassAnnotation(ctx, Lifetime(lifetime))
	}
	containerannotation.BindClassAnnotation = bindClassAnnotation
	containerannotation.InjectParameterAnnotation = injectParameterAnnotation
	containerannotation.NamedParameterAnnotation = namedParameterAnnotation
	netdata.ControllerInstantiator = instantiateController

	vm.AddClass(NewContainerClass())
	vm.AddClass(NewServiceProviderClass())
	vm.AddClass(NewScopeClass())
	containerannotation.Load(vm)
	vm.AddClass(NewCircularDependencyExceptionClass())
}

func onApplicationScanStart(ctx data.Context) (func(), data.Control) {
	e := NewEngine()
	restoreEngine := setRegisteringEngine(e)

	inst := newContainerClass(e)
	cv := data.NewClassValue(inst, ctx.CreateBaseContext())
	e.setHost(cv)
	restoreContainer := setApplicationContainer(cv)

	return func() {
		restoreContainer()
		restoreEngine()
	}, nil
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
