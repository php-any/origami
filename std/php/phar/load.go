package phar

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func Load(vm data.VM) { vm.AddClass(&Class{}) }

type Class struct{ node.Node }

func (*Class) GetName() string                          { return "Phar" }
func (*Class) GetExtend() *string                       { return nil }
func (*Class) GetImplements() []string                  { return nil }
func (*Class) GetProperty(string) (data.Property, bool) { return nil, false }
func (*Class) GetPropertyList() []data.Property         { return nil }
func (*Class) GetConstruct() data.Method                { return nil }
func (c *Class) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *Class) GetMethods() []data.Method {
	return []data.Method{&method{name: "mapPhar"}, &method{name: "running"}}
}
func (c *Class) GetMethod(name string) (data.Method, bool) { return c.GetStaticMethod(name) }
func (*Class) GetStaticMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "mapphar":
		return &method{name: "mapPhar"}, true
	case "running":
		return &method{name: "running"}, true
	}
	return nil, false
}

type method struct{ name string }

func (m *method) GetName() string          { return m.name }
func (*method) GetModifier() data.Modifier { return data.ModifierPublic }
func (*method) GetIsStatic() bool          { return true }
func (*method) GetReturnType() data.Types  { return nil }
func (m *method) GetParams() []data.GetValue {
	if m.name == "running" {
		return []data.GetValue{node.NewParameter(nil, "retphar", 0, data.NewBoolValue(true), nil)}
	}
	return []data.GetValue{node.NewParameter(nil, "alias", 0, data.NewNullValue(), nil), node.NewParameter(nil, "offset", 1, data.NewIntValue(0), nil)}
}
func (m *method) GetVariables() []data.Variable {
	if m.name == "running" {
		return []data.Variable{node.NewVariable(nil, "retphar", 0, nil)}
	}
	return []data.Variable{node.NewVariable(nil, "alias", 0, nil), node.NewVariable(nil, "offset", 1, nil)}
}
func (m *method) Call(ctx data.Context) (data.GetValue, data.Control) {
	var source string
	if tracker, ok := ctx.(data.CallStackTracker); ok {
		frames := tracker.SnapshotCallStack()
		if len(frames) > 0 {
			source = frames[len(frames)-1].File
		}
	}
	if source == "" {
		if tracker, ok := ctx.GetVM().(data.CallStackTracker); ok {
			frames := tracker.SnapshotCallStack()
			if len(frames) > 0 {
				source = frames[len(frames)-1].File
			}
		}
	}
	if source == "" {
		if host, ok := ctx.GetVM().(interface{ CurrentPHPFile() string }); ok {
			source = host.CurrentPHPFile()
		}
	}
	if m.name == "running" {
		if !strings.HasPrefix(source, "phar://") {
			return data.NewStringValue(""), nil
		}
		// Resolve the archive portion from the canonical source URI.
		for i := len(source) - 1; i >= 7; i-- {
			if source[i] == '/' {
				if err := ctx.GetVM().(data.PharFileHost).PharFS().Map(source[7:i], ""); err == nil {
					if v, ok := ctx.GetIndexValue(0); ok {
						if b, ok := v.(data.AsBool); ok {
							ret, _ := b.AsBool()
							if !ret {
								return data.NewStringValue(source[7:i]), nil
							}
						}
					}
					return data.NewStringValue(source[:i]), nil
				}
			}
		}
		return data.NewStringValue(""), nil
	}
	host, ok := ctx.GetVM().(data.PharFileHost)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Phar: VM does not support archives"))
	}
	alias := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		alias = v.AsString()
	}
	if err := host.PharFS().Map(source, alias); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewBoolValue(true), nil
}
