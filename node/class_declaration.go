package node

import (
	"github.com/php-any/origami/data"
	"maps"
	"slices"
)

// A cached declaration is a template. Deferred traits, inherited constructors,
// annotations and static initializers belong to the executing VM, not its AST.
func CloneClassDeclaration(source data.ClassStmt) data.ClassStmt {
	core := coreClass(source)
	if core == nil {
		return source
	}
	clone := &ClassStatement{
		Node: core.Node, Name: core.Name, Extends: core.Extends,
		Implements:      slices.Clone(core.Implements),
		PropertiesIndex: slices.Clone(core.PropertiesIndex), Properties: maps.Clone(core.Properties),
		Methods: maps.Clone(core.Methods), StaticMethods: maps.Clone(core.StaticMethods),
		Construct: core.Construct, IsAbstract: core.IsAbstract, Flags: core.Flags,
		EnumCases:   slices.Clone(core.EnumCases),
		Annotations: slices.Clone(core.Annotations),
		Traits:      slices.Clone(core.Traits), DeferredTraits: slices.Clone(core.DeferredTraits),
		DeferredTraitAliases: slices.Clone(core.DeferredTraitAliases),
		DeferredTraitsMerged: core.DeferredTraitsMerged,
		StaticProperties:     maps.Clone(core.StaticProperties), StaticPropertiesIndex: slices.Clone(core.StaticPropertiesIndex),
		lookupCache: data.NewMethodLookupCache(),
	}
	clone.descriptorTemplate.Store(core.descriptorTemplate.Load())
	for name, property := range clone.Properties {
		clone.Properties[name] = PropertyInClass(property, clone.Name)
	}
	for name, property := range clone.StaticProperties {
		if property, ok := property.(*ClassProperty); ok {
			copy := *property
			clone.StaticProperties[name] = &copy
		}
	}
	// Preserve native static slots without PHP declarations. PHP defaults,
	// including enum cases, are evaluated in the executing VM.
	core.StaticProperty.Range(func(key, value any) bool {
		if _, declared := core.StaticProperties[key.(string)]; !declared {
			clone.StaticProperty.Store(key, value)
		}
		return true
	})
	switch source := source.(type) {
	case *ClassGeneric:
		return &ClassGeneric{ClassStatement: clone, Generic: slices.Clone(source.Generic), GenericMap: maps.Clone(source.GenericMap)}
	case *AbstractClassStatement:
		return &AbstractClassStatement{ClassStatement: clone}
	default:
		return clone
	}
}
