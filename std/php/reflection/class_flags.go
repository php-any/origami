package reflection

import "github.com/php-any/origami/data"

func reflectionClassFlags(ctx data.Context) data.ClassFlags {
	name, class := getReflectionClassInfo(ctx)
	if provider, ok := ctx.GetVM().(data.DescriptorProvider); ok {
		if symbol, found := data.Symbols.Lookup(name); found {
			if descriptor, found := provider.ClassRegistry().Descriptor(symbol); found {
				return descriptor.Flags()
			}
		}
	}
	if metadata, ok := class.(interface{ DeclarationFlags() data.ClassFlags }); ok {
		return metadata.DeclarationFlags()
	}
	return 0
}

type ReflectionClassFlagsMethod struct {
	Name string
	Flag data.ClassFlags
}

func (m *ReflectionClassFlagsMethod) GetName() string             { return m.Name }
func (*ReflectionClassFlagsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (*ReflectionClassFlagsMethod) GetIsStatic() bool             { return false }
func (*ReflectionClassFlagsMethod) GetParams() []data.GetValue    { return nil }
func (*ReflectionClassFlagsMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassFlagsMethod) GetReturnType() data.Types {
	if m.Name == "getModifiers" {
		return data.Int{}
	}
	return data.Bool{}
}
func (m *ReflectionClassFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	flags := reflectionClassFlags(ctx)
	if m.Name != "getModifiers" {
		return data.NewBoolValue(flags&m.Flag != 0), nil
	}
	modifiers := 0
	if flags&data.ClassAbstract != 0 {
		modifiers |= 64
	}
	if flags&data.ClassFinal != 0 {
		modifiers |= 32
	}
	if flags&data.ClassReadonly != 0 {
		modifiers |= 65536
	}
	return data.NewIntValue(modifiers), nil
}
