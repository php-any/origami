package node

import "github.com/php-any/origami/data"

// ClassGeneric 泛型类
type ClassGeneric struct {
	*ClassStatement

	Generic    []data.Types
	GenericMap map[string]data.Types
}

func (c *ClassGeneric) Clone(mT map[string]data.Types) data.ClassGeneric {
	return &ClassGeneric{
		ClassStatement: c.ClassStatement,
		Generic:        c.Generic,
		GenericMap:     mT,
	}
}
func (c *ClassGeneric) GenericList() []data.Types {
	return c.Generic
}

func (c *ClassGeneric) GetProperty(name string) (data.Property, bool) {
	if f, ok := c.Properties[name]; ok {
		// 替换泛型类型
		if f.GetType() != nil {
			t := f.GetType()
			if gt, ok := t.(data.Generic); ok {
				f.SetType(c.GenericMap[gt.Name])
			}
		}
		return f, true
	}
	return nil, false
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
		ext, acl := vm.GetOrLoadClass(*c.Extends)
		if acl != nil {
			return nil, acl
		}

		_, ctl := ext.GetValue(object)
		if ctl != nil {
			return nil, ctl
		}
	}

	return object, nil
}
