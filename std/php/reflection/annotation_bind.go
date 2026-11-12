package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func isAnnotationTargetParam(param data.GetValue) bool {
	switch p := param.(type) {
	case *node.AnnotationTargetParameter:
		return true
	case *node.Parameter:
		return p.Name == node.TargetName
	default:
		return false
}
}

// bindUserAnnotationProperties 将注解声明中的用户参数写入实例属性，
// 便于 newInstance() 后通过 $instance->prop 读取（兼容无 GetPropertyList 的原生注解类）。
func bindUserAnnotationProperties(ctx data.Context, inst *data.ClassValue, decl *node.Annotation) {
	if inst == nil || decl == nil || inst.Class == nil {
		return
	}
	construct := inst.Class.GetConstruct()
	if construct == nil {
		return
	}
	params := construct.GetParams()
	if len(params) == 0 {
		return
	}

	pos := 0
	for _, arg := range decl.Arguments {
		switch a := arg.(type) {
		case *node.NamedArgument:
			if a.Name == node.TargetName {
				continue
			}
			val, ctl := a.GetValue(ctx)
			if ctl != nil || val == nil {
				continue
			}
			if v, ok := val.(data.Value); ok {
				inst.SetProperty(a.Name, v)
			}
		default:
			val, ctl := arg.GetValue(ctx)
			if ctl != nil || val == nil {
				continue
			}
			v, ok := val.(data.Value)
			if !ok {
				continue
			}
			for pos < len(params) {
				if isAnnotationTargetParam(params[pos]) {
					pos++
					continue
				}
				name := constructorParamPropertyName(params[pos])
				if name != "" {
					inst.SetProperty(name, v)
				}
				pos++
				break
			}
		}
	}
}

func constructorParamPropertyName(param data.GetValue) string {
	switch p := param.(type) {
	case *node.PromotedParameter:
		return p.PropertyName
	case *node.Parameter:
		return p.Name
	default:
		if v, ok := param.(data.Variable); ok {
			return v.GetName()
		}
	}
	return ""
}

// userAnnotationArguments 把声明中的实参求值为 getArguments() 应返回的数组（排除编译器注入的 target）。
func userAnnotationArguments(ctx data.Context, decl *node.Annotation) (data.Value, data.Control) {
	array := data.NewArrayValueFromSlots(nil)
	if decl == nil {
		return array, nil
	}
	for _, argument := range decl.Arguments {
		if named, ok := argument.(*node.NamedArgument); ok {
			if named.Name == node.TargetName {
				continue
			}
		}
		raw, ctl := argument.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		value, ok := raw.(data.Value)
		if !ok || value == nil {
			value = data.NewNullValue()
		}
		if named, ok := argument.(*node.NamedArgument); ok {
			array.SetStringKey(named.Name, value)
		} else {
			array.SetKey(data.NewIntValue(array.NextAppendIntKey()), value)
		}
	}
	return array, nil
}
