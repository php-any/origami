package container

import (
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	defaultEngine   *Engine
	defaultInstance data.GetValue
	defaultOnce     sync.Once
)

func defaultContainerClass() *ContainerClass {
	return newContainerClass(defaultEngine)
}

func newContainerClass(engine *Engine) *ContainerClass {
	return &ContainerClass{
		engine:            engine,
		bindMethod:        &ContainerBindMethod{},
		singletonMethod:   &ContainerSingletonMethod{},
		scopedMethod:      &ContainerScopedMethod{},
		instanceMethod:    &ContainerInstanceMethod{},
		makeMethod:        &ContainerMakeMethod{},
		hasMethod:         &ContainerHasMethod{},
		isSharedMethod:    &ContainerIsSharedMethod{},
		aliasMethod:       &ContainerAliasMethod{},
		getInstanceMethod: &ContainerGetInstanceMethod{},
		registerProviders: &ContainerRegisterProvidersMethod{},
		createScopeMethod: &ContainerCreateScopeMethod{},
		scanMethod:        &ContainerScanMethod{},
	}
}

// ContainerClass PHP 类 Container\Container
type ContainerClass struct {
	node.Node
	engine            *Engine
	bindMethod        data.Method
	singletonMethod   data.Method
	scopedMethod      data.Method
	instanceMethod    data.Method
	makeMethod        data.Method
	hasMethod         data.Method
	isSharedMethod    data.Method
	aliasMethod       data.Method
	getInstanceMethod data.Method
	registerProviders data.Method
	createScopeMethod data.Method
	scanMethod        data.Method
}

func NewContainerClass() data.ClassStmt {
	return newContainerClass(nil)
}

func (c *ContainerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	clone := *c
	clone.engine = NewEngine()
	cv := data.NewClassValue(&clone, ctx)
	clone.engine.setHost(cv)
	return cv, nil
}

func (c *ContainerClass) GetName() string { return "Container\\Container" }

func (c *ContainerClass) GetExtend() *string { return nil }

func (c *ContainerClass) GetImplements() []string { return nil }

func (c *ContainerClass) GetProperty(_ string) (data.Property, bool) { return nil, false }

func (c *ContainerClass) GetPropertyList() []data.Property { return nil }

func (c *ContainerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "bind":
		return c.bindMethod, true
	case "singleton":
		return c.singletonMethod, true
	case "scoped":
		return c.scopedMethod, true
	case "instance":
		return c.instanceMethod, true
	case "make":
		return c.makeMethod, true
	case "has":
		return c.hasMethod, true
	case "isShared":
		return c.isSharedMethod, true
	case "alias":
		return c.aliasMethod, true
	case "registerProviders":
		return c.registerProviders, true
	case "createScope":
		return c.createScopeMethod, true
	case "scan":
		return c.scanMethod, true
	}
	return nil, false
}

func (c *ContainerClass) GetStaticMethod(name string) (data.Method, bool) {
	if name == "getInstance" {
		return c.getInstanceMethod, true
	}
	return nil, false
}

func (c *ContainerClass) GetMethods() []data.Method {
	return []data.Method{
		c.bindMethod,
		c.singletonMethod,
		c.scopedMethod,
		c.instanceMethod,
		c.makeMethod,
		c.hasMethod,
		c.isSharedMethod,
		c.aliasMethod,
		c.registerProviders,
		c.createScopeMethod,
		c.scanMethod,
	}
}

func (c *ContainerClass) GetConstruct() data.Method { return nil }

func ensureDefaultInstance(ctx data.Context) (data.GetValue, data.Control) {
	var acl data.Control
	defaultOnce.Do(func() {
		if defaultEngine == nil {
			defaultEngine = NewEngine()
		}
		inst := defaultContainerClass()
		defaultInstance = data.NewClassValue(inst, ctx.CreateBaseContext())
		if inst.engine != nil {
			inst.engine.setHost(defaultInstance)
		}
	})
	return defaultInstance, acl
}
