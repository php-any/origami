package node

import "github.com/php-any/origami/data"

type initializedProperty struct {
	property data.Property
	name     string
}

type classInitialization struct {
	descriptor   *data.ClassDescriptor
	properties   []initializedProperty
	nativeParent data.ClassStmt
}

func (c *ClassStatement) linkedDescriptor(vm data.VM) *data.ClassDescriptor {
	if host, ok := vm.(data.DescriptorProvider); ok {
		declaration := c.DeclarationDescriptor()
		descriptor, _ := host.ClassRegistry().Descriptor(data.SymbolID(declaration.ID()))
		return descriptor
	}
	return nil
}

// Only immutable property declarations are cached. Defaults are evaluated
// anew on each instance, including arrays and request-dependent constants.
func (c *ClassStatement) instanceInitialization(vm data.VM) (*classInitialization, data.Control) {
	descriptor := c.linkedDescriptor(vm)
	if cached := c.initialization.Load(); descriptor != nil && cached != nil && cached.descriptor == descriptor {
		return cached, nil
	}
	plan := &classInitialization{}
	seen := make(map[string]bool)
	add := func(property data.Property) {
		if property == nil || property.GetIsStatic() {
			return
		}
		name := data.PropertyStorageName(property)
		if seen[name] {
			return
		}
		// A redeclared uninitialized typed property also hides its ancestor's
		// default. Private properties have distinct declaring-class keys.
		seen[name] = true
		plan.properties = append(plan.properties, initializedProperty{property, name})
	}
	for _, name := range c.PropertiesIndex {
		add(c.Properties[name])
	}
	var owner data.ClassStmt = c
	for owner.GetExtend() != nil {
		parent, ctl := vm.GetOrLoadClass(*owner.GetExtend())
		if ctl != nil {
			return nil, ctl
		}
		if owner == c {
			original := parent
			for {
				alias, ok := original.(data.OriginalClass)
				if !ok {
					break
				}
				original = alias.OriginalClass()
			}
			if coreClass(original) == nil {
				plan.nativeParent = parent
			}
		}
		for _, property := range parent.GetPropertyList() {
			add(property)
		}
		owner = parent
	}
	// Loading a parent may have linked the child's descriptor in the meantime.
	plan.descriptor = c.linkedDescriptor(vm)
	if plan.descriptor != nil {
		c.initialization.Store(plan)
	}
	return plan, nil
}
