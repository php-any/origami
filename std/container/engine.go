package container

import (
	"fmt"
	"sync"

	"github.com/php-any/origami/data"
)

// Engine IoC 容器核心。
type Engine struct {
	mu        sync.RWMutex
	parent    *Engine
	host      data.GetValue
	bindings  map[string]*Binding
	aliases   map[string]string
	instances map[string]data.GetValue
	resolving []string
}

func NewEngine() *Engine {
	return &Engine{
		bindings:  make(map[string]*Binding),
		aliases:   make(map[string]string),
		instances: make(map[string]data.GetValue),
	}
}

func (e *Engine) root() *Engine {
	if e.parent == nil {
		return e
	}
	return e.parent.root()
}

func (e *Engine) CreateScope() *Engine {
	root := e.root()
	return &Engine{
		parent:    root,
		host:      root.host,
		instances: make(map[string]data.GetValue),
	}
}

func (e *Engine) setHost(host data.GetValue) {
	e.host = host
}

func (e *Engine) containerHost(ctx data.Context) data.GetValue {
	if e.host != nil {
		return e.host
	}
	if e.parent != nil {
		return e.parent.containerHost(ctx)
	}
	if cv, ok := classValueFromCtx(ctx); ok {
		if cc, ok := cv.Class.(*ContainerClass); ok && cc.engine == e {
			return cv
		}
	}
	if defaultInstance != nil {
		return defaultInstance
	}
	return nil
}

func (e *Engine) resolveAbstract(abstract string) string {
	root := e.root()
	root.mu.RLock()
	defer root.mu.RUnlock()

	seen := make(map[string]bool, 4)
	for {
		if seen[abstract] {
			return abstract
		}
		seen[abstract] = true
		target, ok := root.aliases[abstract]
		if !ok {
			return abstract
		}
		abstract = target
	}
}

func (e *Engine) setBinding(abstract string, binding *Binding) {
	root := e.root()
	root.mu.Lock()
	defer root.mu.Unlock()
	root.bindings[abstract] = binding
}

func (e *Engine) Bind(abstract, concrete string) {
	e.setBinding(abstract, &Binding{
		Abstract: abstract,
		Concrete: concrete,
		Lifetime: LifetimeTransient,
	})
}

func (e *Engine) BindFactory(abstract string, factory data.GetValue) {
	e.setBinding(abstract, &Binding{
		Abstract: abstract,
		Factory:  factory,
		Lifetime: LifetimeTransient,
	})
}

func (e *Engine) Singleton(abstract, concrete string) {
	e.setBinding(abstract, &Binding{
		Abstract: abstract,
		Concrete: concrete,
		Lifetime: LifetimeSingleton,
	})
}

func (e *Engine) SingletonFactory(abstract string, factory data.GetValue) {
	e.setBinding(abstract, &Binding{
		Abstract: abstract,
		Factory:  factory,
		Lifetime: LifetimeSingleton,
	})
}

func (e *Engine) Scoped(abstract, concrete string) {
	e.setBinding(abstract, &Binding{
		Abstract: abstract,
		Concrete: concrete,
		Lifetime: LifetimeScoped,
	})
}

func (e *Engine) ScopedFactory(abstract string, factory data.GetValue) {
	e.setBinding(abstract, &Binding{
		Abstract: abstract,
		Factory:  factory,
		Lifetime: LifetimeScoped,
	})
}

func (e *Engine) Alias(alias, abstract string) {
	root := e.root()
	root.mu.Lock()
	defer root.mu.Unlock()
	root.aliases[alias] = abstract
}

func (e *Engine) Instance(abstract string, obj data.GetValue) {
	root := e.root()
	root.mu.Lock()
	defer root.mu.Unlock()
	root.instances[abstract] = obj
	root.bindings[abstract] = &Binding{
		Abstract: abstract,
		Concrete: abstract,
		Lifetime: LifetimeSingleton,
	}
}

func (e *Engine) Has(abstract string) bool {
	abstract = e.resolveAbstract(abstract)
	root := e.root()
	root.mu.RLock()
	defer root.mu.RUnlock()
	if _, ok := root.instances[abstract]; ok {
		return true
	}
	if _, ok := root.bindings[abstract]; ok {
		return true
	}
	return false
}

func (e *Engine) IsShared(abstract string) bool {
	abstract = e.resolveAbstract(abstract)
	root := e.root()
	root.mu.RLock()
	defer root.mu.RUnlock()
	if _, ok := root.instances[abstract]; ok {
		return true
	}
	if b, ok := root.bindings[abstract]; ok {
		return b.Lifetime.shared()
	}
	return false
}

func (e *Engine) lookupBinding(abstract string) (resolved string, binding *Binding) {
	resolved = e.resolveAbstract(abstract)
	root := e.root()
	root.mu.RLock()
	defer root.mu.RUnlock()

	if inst, exists := root.instances[resolved]; exists {
		concrete := resolved
		if cv, isClass := inst.(*data.ClassValue); isClass {
			concrete = cv.Class.GetName()
		}
		return resolved, &Binding{
			Abstract: resolved,
			Concrete: concrete,
			Lifetime: LifetimeSingleton,
		}
	}
	if b, exists := root.bindings[resolved]; exists {
		return resolved, b
	}
	lifetime := metadataLifetime(resolved)
	return resolved, &Binding{
		Abstract: resolved,
		Concrete: resolved,
		Lifetime: lifetime,
	}
}

func (e *Engine) cachedInstance(abstract string, lifetime Lifetime) (data.GetValue, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if inst, ok := e.instances[abstract]; ok {
		return inst, true
	}
	if lifetime == LifetimeSingleton && e.parent != nil {
		return e.parent.cachedInstance(abstract, lifetime)
	}
	return nil, false
}

func (e *Engine) storeInstance(abstract string, lifetime Lifetime, inst data.GetValue) {
	switch lifetime {
	case LifetimeSingleton:
		root := e.root()
		root.mu.Lock()
		root.instances[abstract] = inst
		root.mu.Unlock()
	case LifetimeScoped:
		if e.parent == nil {
			return
		}
		e.mu.Lock()
		e.instances[abstract] = inst
		e.mu.Unlock()
	}
}

// Make 解析 abstract 并返回实例。
func (e *Engine) Make(ctx data.Context, abstract string, params []data.GetValue) (data.GetValue, data.Control) {
	resolved, binding := e.lookupBinding(abstract)
	lifetime := binding.Lifetime

	e.mu.Lock()
	for _, a := range e.resolving {
		if a == resolved {
			e.mu.Unlock()
			return nil, circularDependencyError(abstract)
		}
	}
	e.mu.Unlock()

	if lifetime == LifetimeScoped && e.parent == nil {
		return nil, scopedResolveError(resolved)
	}

	if inst, ok := e.cachedInstance(resolved, lifetime); ok {
		return inst, nil
	}

	e.mu.Lock()
	e.resolving = append(e.resolving, resolved)
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		if n := len(e.resolving); n > 0 {
			e.resolving = e.resolving[:n-1]
		}
		e.mu.Unlock()
	}()

	var result data.GetValue
	var acl data.Control

	if binding.isFactory() {
		result, acl = invokeFactory(ctx, binding.Factory, e.containerHost(ctx))
	} else {
		concrete := binding.Concrete
		if concrete == "" {
			concrete = resolved
		}

		vm := ctx.GetVM()
		stmt, loadACL := vm.GetOrLoadClass(concrete)
		if loadACL != nil {
			return nil, loadACL
		}
		if stmt == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("Target class [%s] is not instantiable.", concrete))
		}

		args, resolveACL := ResolveConstructor(ctx, e, stmt, params)
		if resolveACL != nil {
			return nil, resolveACL
		}

		result, acl = instantiateClass(stmt, args, ctx)
	}

	if acl != nil {
		return nil, acl
	}

	if lifetime.shared() {
		if cached, ok := e.cachedInstance(resolved, lifetime); ok {
			return cached, nil
		}
		e.storeInstance(resolved, lifetime, result)
	}

	return result, nil
}

// RegisterClass 由 Component/Singleton/Scoped 注解调用。
func (e *Engine) RegisterClass(className string, lifetime Lifetime, alias string) {
	abstract := className
	if alias != "" {
		abstract = alias
	}
	switch lifetime {
	case LifetimeSingleton:
		e.Singleton(abstract, className)
	case LifetimeScoped:
		e.Scoped(abstract, className)
	default:
		e.Bind(abstract, className)
	}
}
