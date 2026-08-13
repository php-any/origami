package container

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ResolveConstructor 根据构造器类型提示自动装配参数；显式传入的 params 按位置覆盖。
func ResolveConstructor(ctx data.Context, e *Engine, stmt data.ClassStmt, provided []data.GetValue) ([]data.GetValue, data.Control) {
	ctor := stmt.GetConstruct()
	if ctor == nil {
		return provided, nil
	}
	defs := ctor.GetParams()
	if len(defs) == 0 {
		return provided, nil
	}

	args := make([]data.GetValue, len(defs))
	for i, def := range defs {
		if i < len(provided) && provided[i] != nil {
			args[i] = provided[i]
			continue
		}

		abstract, named, ok := paramResolveName(def, stmt.GetName())
		if !ok {
			if name := paramName(def); name != "" && e.Has(name) {
				abstract = name
				ok = true
			}
		}
		if !ok {
			if val, acl := defaultParamValue(ctx, def); acl != nil {
				return nil, acl
			} else if val != nil {
				args[i] = val
			}
			continue
		}
		if named != "" {
			abstract = named
		}

		resolved, acl := e.Make(ctx, abstract, nil)
		if acl != nil {
			if val, defACL := defaultParamValue(ctx, def); defACL == nil && val != nil {
				args[i] = val
				continue
			}
			return nil, acl
		}
		args[i] = resolved
	}
	return args, nil
}

func defaultParamValue(ctx data.Context, def data.GetValue) (data.GetValue, data.Control) {
	switch p := def.(type) {
	case *node.Parameter:
		if p.DefaultValue == nil {
			return nil, nil
		}
		v, acl := p.DefaultValue.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		return v, nil
	case *node.PromotedParameter:
		if p.DefaultValue == nil {
			return nil, nil
		}
		v, acl := p.DefaultValue.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		return v, nil
	default:
		return nil, nil
	}
}

func paramResolveName(def data.GetValue, className string) (abstract, named string, ok bool) {
	if name := metadataNamedService(className, paramIndex(def)); name != "" {
		return name, name, true
	}
	ty := paramType(def)
	if ty == nil {
		return "", "", false
	}
	if className != "" {
		ty = resolveSelfType(ty, className)
	}
	abstract, ok = typeHintAbstract(ty)
	return abstract, "", ok
}

func paramName(def data.GetValue) string {
	switch p := def.(type) {
	case *node.Parameter:
		return p.Name
	case *node.PromotedParameter:
		return p.Name
	default:
		return ""
	}
}

func paramIndex(def data.GetValue) int {
	switch p := def.(type) {
	case *node.Parameter:
		return p.Index
	case *node.PromotedParameter:
		return p.Index
	default:
		return -1
	}
}

func paramType(def data.GetValue) data.Types {
	switch p := def.(type) {
	case *node.Parameter:
		return p.Type
	case *node.PromotedParameter:
		return p.Type
	default:
		return nil
	}
}

func resolveSelfType(ty data.Types, className string) data.Types {
	switch ty.(type) {
	case data.StaticType:
		return data.Class{Name: className}
	default:
		return ty
	}
}

func typeHintAbstract(ty data.Types) (string, bool) {
	if ty == nil {
		return "", false
	}
	switch t := ty.(type) {
	case data.Class:
		if data.ISBaseType(t.Name) || t.Name == "self" || t.Name == "static" {
			return "", false
		}
		return t.Name, true
	case data.NullableType:
		return typeHintAbstract(t.BaseType)
	case data.UnionType:
		for _, ut := range t.Types {
			if name, ok := typeHintAbstract(ut); ok {
				return name, true
			}
		}
		return "", false
	case data.Generic:
		if data.ISBaseType(t.Name) {
			return "", false
		}
		return t.Name, true
	default:
		name := ty.String()
		if name == "" || data.ISBaseType(name) {
			return "", false
		}
		return name, true
	}
}

func scopedResolveError(abstract string) data.Control {
	return data.NewErrorThrow(nil, fmt.Errorf("Cannot resolve scoped service [%s] from the root container. Create a scope first.", abstract))
}
