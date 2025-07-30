package node

import "github.com/php-any/origami/data"

// ClassGeneric 泛型类
type ClassGeneric struct {
	*ClassStatement

	Generic    []data.Types
	GenericMap map[string]data.Types
}

func (c *ClassGeneric) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	object := data.NewClassValue(c, ctx)

	for _, property := range c.Properties {
		def := property.GetDefaultValue()
		if def == nil {
			continue
		}
		v, ctl := def.GetValue(object)
		if ctl != nil {
			return nil, ctl
		}
		object.SetProperty(property.GetName(), v.(data.Value))
	}
	if c.Extends != nil {
		vm := object.GetVM()
		ext, ok := vm.GetClass(*c.Extends)
		if ok {
			_, ctl := ext.GetValue(object)
			if ctl != nil {
				return nil, ctl
			}
		}
	}

	return object, nil
}
