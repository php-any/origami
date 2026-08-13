package php

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

type stubClass struct {
	name string
}

func (c *stubClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
func (c *stubClass) GetFrom() data.From                       { return nil }
func (c *stubClass) GetName() string                          { return c.name }
func (c *stubClass) GetExtend() *string                       { return nil }
func (c *stubClass) GetImplements() []string                  { return nil }
func (c *stubClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *stubClass) GetPropertyList() []data.Property         { return nil }
func (c *stubClass) GetMethod(string) (data.Method, bool)     { return nil, false }
func (c *stubClass) GetMethods() []data.Method                { return nil }
func (c *stubClass) GetConstruct() data.Method                { return nil }

func TestClassAliasRegistersProxyClass(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	ctx := vm.CreateContext(nil)

	if control := vm.AddClass(&stubClass{name: "Illuminate\\Support\\Facades\\Route"}); control != nil {
		t.Fatalf("AddClass failed: %v", control)
	}

	fn := NewClassAliasFunction()
	callCtx := ctx.CreateContext(fn.GetVariables())
	callCtx.SetVariableValue(node.NewVariable(nil, "original", 0, data.String{}), data.NewStringValue("Illuminate\\Support\\Facades\\Route"))
	callCtx.SetVariableValue(node.NewVariable(nil, "alias", 1, data.String{}), data.NewStringValue("Route"))
	callCtx.SetVariableValue(node.NewVariable(nil, "autoload", 2, data.NewBaseType("bool")), data.NewBoolValue(true))

	result, control := fn.Call(callCtx)
	if control != nil {
		t.Fatalf("class_alias returned control: %v", control)
	}
	ok, err := result.(data.AsBool).AsBool()
	if err != nil || !ok {
		t.Fatalf("class_alias = false, want true")
	}

	aliasClass, exists := vm.GetClass("Route")
	if !exists {
		t.Fatal("Route alias was not registered")
	}
	if aliasClass.GetName() != "Route" {
		t.Fatalf("alias class name = %q, want Route", aliasClass.GetName())
	}
}
