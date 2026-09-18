package vardumper

import (
	"fmt"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)


const varDumperName = "Symfony\\Component\\VarDumper\\VarDumper"

type VarDumperClass struct{ node.Node }

func NewVarDumperClass() data.ClassStmt { return &VarDumperClass{} }
func (c *VarDumperClass) GetName() string {
	return varDumperName
}
func (c *VarDumperClass) GetExtend() *string                       { return nil }
func (c *VarDumperClass) GetImplements() []string                  { return nil }
func (c *VarDumperClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *VarDumperClass) GetPropertyList() []data.Property         { return nil }
func (c *VarDumperClass) GetConstruct() data.Method                { return nil }
func (c *VarDumperClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *VarDumperClass) GetMethod(name string) (data.Method, bool) {
	if name == "dump" {
		return &vdMethod{name: "dump", static: true, params: []string{"var"}, fn: vdDump}, true
	}
	return nil, false
}
func (c *VarDumperClass) GetMethods() []data.Method { return nil }
func (c *VarDumperClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

type vdMethod struct {
	name   string
	params []string
	static bool
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *vdMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *vdMethod) GetName() string                                     { return m.name }
func (m *vdMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *vdMethod) GetIsStatic() bool                                   { return m.static }
func (m *vdMethod) GetReturnType() data.Types                           { return nil }
func (m *vdMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p, i, nil, nil)
	}
	return out
}
func (m *vdMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}

func vdDump(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		fmt.Fprintln(os.Stdout, "null")
		return data.NewNullValue(), nil
	}
	fmt.Fprintln(os.Stdout, v.AsString())
	return v, nil
}
