package container

import (
	"sync"

	"github.com/php-any/origami/node"
)

type InjectPoint struct {
	ParamIndex  int
	ParamName   string
	ServiceName string
	HasInject   bool
}

type ClassMetadata struct {
	ClassName string
	Lifetime  Lifetime
	Alias     string
	Params    map[int]InjectPoint
}

var (
	metaMu   sync.RWMutex
	metaByVM = make(map[uintptr]map[string]*ClassMetadata)
)

func metadataStore(className string) *ClassMetadata {
	metaMu.Lock()
	defer metaMu.Unlock()
	if metaByVM[0] == nil {
		metaByVM[0] = make(map[string]*ClassMetadata)
	}
	if m, ok := metaByVM[0][className]; ok {
		return m
	}
	m := &ClassMetadata{
		ClassName: className,
		Lifetime:  LifetimeTransient,
		Params:    make(map[int]InjectPoint),
	}
	metaByVM[0][className] = m
	return m
}

func metadataGet(className string) *ClassMetadata {
	metaMu.RLock()
	defer metaMu.RUnlock()
	if metaByVM[0] == nil {
		return nil
	}
	return metaByVM[0][className]
}

func metadataSetLifetime(className string, lifetime Lifetime, alias string) {
	m := metadataStore(className)
	m.Lifetime = lifetime
	if alias != "" {
		m.Alias = alias
	}
}

func metadataMarkConstructorInject(className string, paramIndex int, paramName, serviceName string, hasInject bool) {
	m := metadataStore(className)
	m.Params[paramIndex] = InjectPoint{
		ParamIndex:  paramIndex,
		ParamName:   paramName,
		ServiceName: serviceName,
		HasInject:   hasInject,
	}
}

func metadataNamedService(className string, paramIndex int) string {
	m := metadataGet(className)
	if m == nil {
		return ""
	}
	if p, ok := m.Params[paramIndex]; ok && p.ServiceName != "" {
		return p.ServiceName
	}
	return ""
}

func metadataLifetime(className string) Lifetime {
	m := metadataGet(className)
	if m == nil {
		return LifetimeTransient
	}
	return m.Lifetime
}

func metadataAlias(className string) string {
	m := metadataGet(className)
	if m == nil {
		return ""
	}
	return m.Alias
}

func metadataFromClassTarget(target any) (className string, ok bool) {
	switch t := target.(type) {
	case *node.ClassStatement:
		return t.Name, true
	default:
		return "", false
	}
}

func metadataFromParamTarget(target any) (className string, paramIndex int, paramName string, ok bool) {
	switch t := target.(type) {
	case *node.Parameter:
		return "", t.Index, t.Name, true
	case *node.PromotedParameter:
		return "", t.Index, t.Name, true
	default:
		return "", -1, "", false
	}
}
