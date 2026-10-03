package runtime_test

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
	"testing"
)

func TestRequestDeclarationDescriptorsIncludeStaticDefaults(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	request := runtime.NewRequestVM(base).(*runtime.RequestVM)
	request.SetOutputWriter(func(string) {})
	if _, ctl := request.LoadAndRun("../tests/php/property_reference_identity_test.php"); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	classes := request.AddedClasses()
	staticProperties := 0
	for _, class := range classes {
		descriptor, found := request.ClassRegistry().Descriptor(data.Symbols.Intern(class.GetName()))
		if !found {
			t.Fatalf("descriptor missing: %s", class.GetName())
		}
		if metadata, ok := class.(interface{ StaticPropertyDeclarations() []data.Property }); ok {
			for _, declaration := range metadata.StaticPropertyDeclarations() {
				found := false
				for _, property := range descriptor.Properties().Range() {
					if property.Name() == declaration.GetName() && property.IsStatic() {
						found = true
						staticProperties++
						if property.DefaultValue() != declaration.GetDefaultValue() || property.Type() != data.DeclaredTypeRef(declaration.GetType()) {
							t.Error("static property descriptor lost default or type")
						}
					}
				}
				if !found {
					t.Errorf("static declaration missing: %s::$%s", class.GetName(), declaration.GetName())
				}
			}
		}
	}
	if staticProperties == 0 {
		t.Fatal("fixture did not exercise static declarations")
	}
}
