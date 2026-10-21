package reflection

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ReflectionClassConstantsMethod struct{ Values bool }

func (m *ReflectionClassConstantsMethod) GetName() string {
	if m.Values {
		return "getConstants"
	}
	return "getReflectionConstants"
}
func (*ReflectionClassConstantsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (*ReflectionClassConstantsMethod) GetIsStatic() bool          { return false }
func (*ReflectionClassConstantsMethod) GetReturnType() data.Types  { return nil }
func (*ReflectionClassConstantsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "filter", 0, data.NewNullValue(), nil)}
}
func (*ReflectionClassConstantsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "filter", 0, nil)}
}
func constantProperties(class data.ClassStmt) []data.Property {
	if provider, ok := class.(interface{ StaticPropertyDeclarations() []data.Property }); ok {
		return provider.StaticPropertyDeclarations()
	}
	return nil
}
func (m *ReflectionClassConstantsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, class := getReflectionClassInfo(ctx)
	out := data.NewArrayValueFromSlots(nil)
	seen := map[string]bool{}
	filter := 0
	if v, ok := ctx.GetIndexValue(0); ok {
		if n, ok := v.(data.AsInt); ok {
			filter, _ = n.AsInt()
		}
	}
	inherited := false
	for class != nil {
		for _, p := range constantProperties(class) {
			decl, ok := p.(interface{ IsClassConstant() bool })
			if !ok || !decl.IsClassConstant() || seen[p.GetName()] || (inherited && p.GetModifier() == data.ModifierPrivate) {
				continue
			}
			seen[p.GetName()] = true
			flags := constantVisibility(p)
			if filter != 0 && filter&flags == 0 {
				continue
			}
			if m.Values {
				if v, ok := node.LookupStaticProperty(ctx.GetVM(), class, p.GetName()); ok {
					out.SetStringKey(p.GetName(), v)
				}
			} else {
				v := data.NewClassValue(&ReflectionClassConstantClass{}, ctx.CreateBaseContext())
				v.SetProperty("name", data.NewStringValue(p.GetName()))
				v.SetProperty("class", data.NewStringValue(class.GetName()))
				out.AppendValue(v)
			}
		}
		if class.GetExtend() == nil {
			break
		}
		parent, ctl := ctx.GetVM().GetOrLoadClass(*class.GetExtend())
		if ctl != nil {
			return nil, ctl
		}
		class = parent
		inherited = true
	}
	return out, nil
}
func constantVisibility(p data.Property) int {
	switch p.GetModifier() {
	case data.ModifierPrivate:
		return 4
	case data.ModifierProtected:
		return 2
	default:
		return 1
	}
}

type ReflectionClassConstantClass struct{ node.Node }

func (*ReflectionClassConstantClass) GetName() string                          { return "ReflectionClassConstant" }
func (*ReflectionClassConstantClass) GetExtend() *string                       { return nil }
func (*ReflectionClassConstantClass) GetImplements() []string                  { return nil }
func (*ReflectionClassConstantClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (*ReflectionClassConstantClass) GetPropertyList() []data.Property         { return nil }
func (*ReflectionClassConstantClass) GetConstruct() data.Method {
	return &classConstantMethod{name: "__construct"}
}
func (c *ReflectionClassConstantClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ReflectionClassConstantClass) GetMethods() []data.Method {
	var out []data.Method
	for _, name := range []string{"__construct", "getName", "getValue", "getDeclaringClass", "isPublic", "isProtected", "isPrivate", "getModifiers"} {
		out = append(out, &classConstantMethod{name: name})
	}
	return out
}
func (c *ReflectionClassConstantClass) GetMethod(name string) (data.Method, bool) {
	for _, m := range c.GetMethods() {
		if strings.EqualFold(name, m.GetName()) {
			return m, true
		}
	}
	return nil, false
}

type classConstantMethod struct{ name string }

func (m *classConstantMethod) GetName() string          { return m.name }
func (*classConstantMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (*classConstantMethod) GetIsStatic() bool          { return false }
func (*classConstantMethod) GetReturnType() data.Types  { return nil }
func (m *classConstantMethod) GetParams() []data.GetValue {
	if m.name != "__construct" {
		return nil
	}
	return []data.GetValue{node.NewParameter(nil, "class", 0, nil, nil), node.NewParameter(nil, "constant", 1, nil, nil)}
}
func (m *classConstantMethod) GetVariables() []data.Variable {
	if m.name != "__construct" {
		return nil
	}
	return []data.Variable{node.NewVariable(nil, "class", 0, nil), node.NewVariable(nil, "constant", 1, nil)}
}
func (m *classConstantMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	receiver, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionClassConstant: missing receiver"))
	}
	className := phpValueAsString(reflectionProperty(receiver.ObjectValue, "class"))
	name := phpValueAsString(reflectionProperty(receiver.ObjectValue, "name"))
	if m.name == "__construct" {
		className, _ = classNameFromObjectOrString(reflectionConstructArg(ctx, 0))
		name = phpValueAsString(reflectionConstructArg(ctx, 1))
	}
	class, ctl := ctx.GetVM().GetOrLoadClass(className)
	if ctl != nil {
		return nil, ctl
	}
	for class != nil {
		for _, p := range constantProperties(class) {
			decl, ok := p.(interface{ IsClassConstant() bool })
			if p.GetName() != name || !ok || !decl.IsClassConstant() {
				continue
			}
			switch m.name {
			case "__construct":
				receiver.ObjectValue.SetProperty("class", data.NewStringValue(class.GetName()))
				receiver.ObjectValue.SetProperty("name", data.NewStringValue(name))
				return nil, nil
			case "getName":
				return data.NewStringValue(name), nil
			case "getValue":
				if v, ok := node.LookupStaticProperty(ctx.GetVM(), class, name); ok {
					return v, nil
				}
			case "getDeclaringClass":
				return newReflectionClassValue(ctx, class.GetName()), nil
			case "getModifiers":
				return data.NewIntValue(constantVisibility(p)), nil
			case "isPublic":
				return data.NewBoolValue(constantVisibility(p) == 1), nil
			case "isProtected":
				return data.NewBoolValue(constantVisibility(p) == 2), nil
			case "isPrivate":
				return data.NewBoolValue(constantVisibility(p) == 4), nil
			}
		}
		if class.GetExtend() == nil {
			break
		}
		class, ctl = ctx.GetVM().GetOrLoadClass(*class.GetExtend())
		if ctl != nil {
			return nil, ctl
		}
	}
	return nil, createReflectionException("Constant "+className+"::"+name+" does not exist", ctx, nil)
}
