package process

import (
	"os/exec"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)


const processName = "Symfony\\Component\\Process\\Process"

type ProcessClass struct{ node.Node }

func NewProcessClass() data.ClassStmt { return &ProcessClass{} }
func (c *ProcessClass) GetName() string {
	return processName
}
func (c *ProcessClass) GetExtend() *string                       { return nil }
func (c *ProcessClass) GetImplements() []string                  { return nil }
func (c *ProcessClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *ProcessClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "commandline", "private", false, data.NewArrayValue(nil)),
		node.NewProperty(nil, "output", "private", false, data.NewStringValue("")),
		node.NewProperty(nil, "exitcode", "private", false, data.NewIntValue(-1)),
	}
}
func (c *ProcessClass) GetConstruct() data.Method {
	return &procMethod{name: "__construct", params: []string{"command", "cwd", "env", "input", "timeout"}, fn: processConstruct}
}
func (c *ProcessClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ProcessClass) GetMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "__construct":
		return c.GetConstruct(), true
	case "run":
		return &procMethod{name: "run", params: []string{"callback", "env"}, fn: processRun}, true
	case "getoutput":
		return &procMethod{name: "getOutput", fn: processGetOutput}, true
	case "getexitcode":
		return &procMethod{name: "getExitCode", fn: processGetExitCode}, true
	case "fromshellcommandline":
		return &procMethod{name: "fromShellCommandline", static: true, params: []string{"command", "cwd", "env", "input", "timeout"}, fn: processFromShell}, true
	}
	return nil, false
}
func (c *ProcessClass) GetMethods() []data.Method            { return []data.Method{c.GetConstruct()} }
func (c *ProcessClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

type procMethod struct {
	name   string
	params []string
	static bool
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *procMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *procMethod) GetName() string                                     { return m.name }
func (m *procMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *procMethod) GetIsStatic() bool                                   { return m.static }
func (m *procMethod) GetReturnType() data.Types                           { return nil }
func (m *procMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p, i, nil, nil)
	}
	return out
}
func (m *procMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}

func procSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func processConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := procSelf(ctx)
	cmd, _ := ctx.GetIndexValue(0)
	_ = cv.SetProperty("commandline", cmd)
	return data.NewNullValue(), nil
}

func processFromShell(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(NewProcessClass(), ctx.CreateBaseContext())
	cmd, _ := ctx.GetIndexValue(0)
	_ = cv.SetProperty("commandline", data.NewArrayValue([]data.Value{data.NewStringValue(cmd.AsString())}))
	return cv, nil
}

func processRun(ctx data.Context) (data.GetValue, data.Control) {
	cv := procSelf(ctx)
	cmdV, _ := cv.GetProperty("commandline")
	var args []string
	if av, ok := cmdV.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil {
				args = append(args, z.Value.AsString())
			}
		}
	} else if cmdV != nil {
		args = []string{cmdV.AsString()}
	}
	if len(args) == 0 {
		_ = cv.SetProperty("exitcode", data.NewIntValue(1))
		return data.NewIntValue(1), nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	_ = cv.SetProperty("output", data.NewStringValue(string(out)))
	code := 0
	if err != nil {
		code = 1
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		}
	}
	_ = cv.SetProperty("exitcode", data.NewIntValue(code))
	return data.NewIntValue(code), nil
}

func processGetOutput(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := procSelf(ctx).GetProperty("output")
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(v.AsString()), nil
}

func processGetExitCode(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := procSelf(ctx).GetProperty("exitcode")
	if v == nil {
		return data.NewIntValue(-1), nil
	}
	return v, nil
}
