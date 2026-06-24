package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewScopeClass() data.ClassStmt {
	return &ScopeClass{
		makeMethod:    &ScopeMakeMethod{},
		disposeMethod: &ScopeDisposeMethod{},
	}
}

type ScopeClass struct {
	node.Node
	engine        *Engine
	makeMethod    data.Method
	disposeMethod data.Method
	disposed      bool
}

func (s *ScopeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, utils.NewThrow(errors.New("Container\\Scope 请通过 Container::createScope() 创建"))
}

func (s *ScopeClass) GetName() string                            { return "Container\\Scope" }
func (s *ScopeClass) GetExtend() *string                         { return nil }
func (s *ScopeClass) GetImplements() []string                    { return nil }
func (s *ScopeClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *ScopeClass) GetPropertyList() []data.Property           { return nil }
func (s *ScopeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "make":
		return s.makeMethod, true
	case "dispose":
		return s.disposeMethod, true
	}
	return nil, false
}
func (s *ScopeClass) GetMethods() []data.Method {
	return []data.Method{s.makeMethod, s.disposeMethod}
}
func (s *ScopeClass) GetConstruct() data.Method { return nil }

type ScopeMakeMethod struct{}

func (m *ScopeMakeMethod) GetName() string            { return "make" }
func (m *ScopeMakeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ScopeMakeMethod) GetIsStatic() bool          { return false }
func (m *ScopeMakeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("abstract", 0),
		data.NewParameterDefault("parameters", 1, data.NewNullValue(), nil),
	}
}
func (m *ScopeMakeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("abstract", 0, data.NewBaseType("string")),
		data.NewVariable("parameters", 1, nil),
	}
}
func (m *ScopeMakeMethod) GetReturnType() data.Types { return data.NewBaseType("object") }
func (m *ScopeMakeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	scope, acl := scopeFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	if scope.disposed {
		return nil, utils.NewThrow(errors.New("Container\\Scope 已 dispose"))
	}
	abstract, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	var params []data.GetValue
	if v, ok := ctx.GetIndexValue(1); ok {
		if arr, isArr := v.(*data.ArrayValue); isArr {
			for _, item := range arr.List {
				params = append(params, item.Value)
			}
		}
	}
	return scope.engine.Make(ctx, abstract, params)
}

type ScopeDisposeMethod struct{}

func (m *ScopeDisposeMethod) GetName() string               { return "dispose" }
func (m *ScopeDisposeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ScopeDisposeMethod) GetIsStatic() bool             { return false }
func (m *ScopeDisposeMethod) GetParams() []data.GetValue    { return nil }
func (m *ScopeDisposeMethod) GetVariables() []data.Variable { return nil }
func (m *ScopeDisposeMethod) GetReturnType() data.Types     { return nil }
func (m *ScopeDisposeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	scope, acl := scopeFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	scope.disposed = true
	scope.engine.mu.Lock()
	scope.engine.instances = make(map[string]data.GetValue)
	scope.engine.mu.Unlock()
	return data.NewNullValue(), nil
}

func scopeFromCtx(ctx data.Context) (*ScopeClass, data.Control) {
	cv, ok := classValueFromCtx(ctx)
	if !ok {
		return nil, utils.NewThrow(errors.New("Scope 方法必须在 Container\\Scope 实例上调用"))
	}
	s, ok := cv.Class.(*ScopeClass)
	if !ok {
		return nil, utils.NewThrow(errors.New("Scope 方法必须在 Container\\Scope 实例上调用"))
	}
	return s, nil
}

func newScopeValue(parent *Engine, ctx data.Context) (data.GetValue, data.Control) {
	scopeEngine := parent.CreateScope()
	scope := &ScopeClass{
		engine:        scopeEngine,
		makeMethod:    &ScopeMakeMethod{},
		disposeMethod: &ScopeDisposeMethod{},
	}
	return data.NewClassValue(scope, ctx), nil
}
