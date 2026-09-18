package console

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)


const argvInputName = "Symfony\\Component\\Console\\Input\\ArgvInput"
const arrayInputName = "Symfony\\Component\\Console\\Input\\ArrayInput"
const consoleOutputName = "Symfony\\Component\\Console\\Output\\ConsoleOutput"

type ArgvInputClass struct {
	node.Node
	methods map[string]data.Method
}

func NewArgvInputClass() data.ClassStmt {
	c := &ArgvInputClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = &consoleMethod{name: "__construct", params: []string{"argv", "definition"}, fn: argvConstruct}
	c.methods["getfirstargument"] = &consoleMethod{name: "getFirstArgument", fn: argvGetFirst}
	c.methods["hasparameteroption"] = &consoleMethod{name: "hasParameterOption", params: []string{"values", "onlyParams"}, fn: argvHasOption}
	c.methods["getparameteroption"] = &consoleMethod{name: "getParameterOption", params: []string{"values", "default", "onlyParams"}, fn: argvGetOption}
	c.methods["__tostring"] = &consoleMethod{name: "__toString", fn: argvToString}
	return c
}

func (c *ArgvInputClass) GetName() string { return argvInputName }
func (c *ArgvInputClass) GetExtend() *string {
	parent := "Symfony\\Component\\Console\\Input\\Input"
	return &parent
}
func (c *ArgvInputClass) GetImplements() []string {
	return []string{"Symfony\\Component\\Console\\Input\\InputInterface"}
}
func (c *ArgvInputClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *ArgvInputClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "tokens", "private", false, data.NewArrayValue(nil)),
	}
}
func (c *ArgvInputClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *ArgvInputClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ArgvInputClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *ArgvInputClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}

type ArrayInputClass struct{ node.Node }

func NewArrayInputClass() data.ClassStmt { return &ArrayInputClass{} }
func (c *ArrayInputClass) GetName() string {
	return arrayInputName
}
func (c *ArrayInputClass) GetExtend() *string {
	parent := "Symfony\\Component\\Console\\Input\\Input"
	return &parent
}
func (c *ArrayInputClass) GetImplements() []string {
	return []string{"Symfony\\Component\\Console\\Input\\InputInterface"}
}
func (c *ArrayInputClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *ArrayInputClass) GetPropertyList() []data.Property         { return nil }
func (c *ArrayInputClass) GetConstruct() data.Method {
	return &consoleMethod{name: "__construct", params: []string{"parameters", "definition"}, fn: func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewNullValue(), nil
	}}
}
func (c *ArrayInputClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ArrayInputClass) GetMethod(name string) (data.Method, bool) {
	if strings.EqualFold(name, "__construct") {
		return c.GetConstruct(), true
	}
	return nil, false
}
func (c *ArrayInputClass) GetMethods() []data.Method { return []data.Method{c.GetConstruct()} }

type ConsoleOutputClass struct{ node.Node }

func NewConsoleOutputClass() data.ClassStmt { return &ConsoleOutputClass{} }
func (c *ConsoleOutputClass) GetName() string {
	return consoleOutputName
}
func (c *ConsoleOutputClass) GetExtend() *string {
	parent := "Symfony\\Component\\Console\\Output\\Output"
	return &parent
}
func (c *ConsoleOutputClass) GetImplements() []string {
	return []string{"Symfony\\Component\\Console\\Output\\OutputInterface"}
}
func (c *ConsoleOutputClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *ConsoleOutputClass) GetPropertyList() []data.Property         { return nil }
func (c *ConsoleOutputClass) GetConstruct() data.Method {
	return &consoleMethod{name: "__construct", fn: func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewNullValue(), nil
	}}
}
func (c *ConsoleOutputClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ConsoleOutputClass) GetMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "__construct":
		return c.GetConstruct(), true
	case "writeln":
		return &consoleMethod{name: "writeln", params: []string{"messages", "options"}, fn: consoleWriteln}, true
	case "write":
		return &consoleMethod{name: "write", params: []string{"messages", "newline", "options"}, fn: consoleWrite}, true
	}
	return nil, false
}
func (c *ConsoleOutputClass) GetMethods() []data.Method {
	return []data.Method{c.GetConstruct()}
}

type consoleMethod struct {
	name   string
	params []string
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *consoleMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *consoleMethod) GetName() string                                     { return m.name }
func (m *consoleMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *consoleMethod) GetIsStatic() bool                                   { return false }
func (m *consoleMethod) GetReturnType() data.Types                           { return nil }
func (m *consoleMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p, i, nil, nil)
	}
	return out
}
func (m *consoleMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}

func argvSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func argvConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := argvSelf(ctx)
	argv, _ := ctx.GetIndexValue(0)
	tokens := data.NewArrayValue(nil).(*data.ArrayValue)
	if argv == nil || isNull(argv) {
		for _, a := range os.Args {
			tokens.List = append(tokens.List, data.NewZVal(data.NewStringValue(a)))
		}
	} else if av, ok := argv.(*data.ArrayValue); ok {
		tokens = data.CloneArrayValue(av)
	}
	_ = cv.SetProperty("tokens", tokens)
	return data.NewNullValue(), nil
}

func isNull(v data.Value) bool {
	_, ok := v.(*data.NullValue)
	return ok
}

func argvTokens(cv *data.ClassValue) []string {
	v, _ := cv.GetProperty("tokens")
	out := make([]string, 0)
	if av, ok := v.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil {
				out = append(out, z.Value.AsString())
			}
		}
	}
	return out
}

func argvGetFirst(ctx data.Context) (data.GetValue, data.Control) {
	tokens := argvTokens(argvSelf(ctx))
	for _, t := range tokens {
		if t == "" || strings.HasPrefix(t, "-") {
			continue
		}
		return data.NewStringValue(t), nil
	}
	return data.NewNullValue(), nil
}

func argvHasOption(ctx data.Context) (data.GetValue, data.Control) {
	values, _ := ctx.GetIndexValue(0)
	tokens := argvTokens(argvSelf(ctx))
	needles := []string{}
	if av, ok := values.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil {
				needles = append(needles, z.Value.AsString())
			}
		}
	} else if values != nil {
		needles = append(needles, values.AsString())
	}
	for _, t := range tokens {
		for _, n := range needles {
			if t == n || strings.HasPrefix(t, n+"=") {
				return data.NewBoolValue(true), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}

func argvGetOption(ctx data.Context) (data.GetValue, data.Control) {
	values, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	tokens := argvTokens(argvSelf(ctx))
	needles := []string{}
	if av, ok := values.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil {
				needles = append(needles, z.Value.AsString())
			}
		}
	} else if values != nil {
		needles = append(needles, values.AsString())
	}
	for i, t := range tokens {
		for _, n := range needles {
			if t == n {
				if i+1 < len(tokens) {
					return data.NewStringValue(tokens[i+1]), nil
				}
				return data.NewBoolValue(true), nil
			}
			if strings.HasPrefix(t, n+"=") {
				return data.NewStringValue(strings.TrimPrefix(t, n+"=")), nil
			}
		}
	}
	if def == nil {
		return data.NewBoolValue(false), nil
	}
	return def, nil
}

func argvToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(strings.Join(argvTokens(argvSelf(ctx)), " ")), nil
}

func consoleWriteln(ctx data.Context) (data.GetValue, data.Control) {
	msg, _ := ctx.GetIndexValue(0)
	if msg != nil {
		os.Stdout.WriteString(msg.AsString() + "\n")
	}
	return data.NewNullValue(), nil
}

func consoleWrite(ctx data.Context) (data.GetValue, data.Control) {
	msg, _ := ctx.GetIndexValue(0)
	nl, _ := ctx.GetIndexValue(1)
	s := ""
	if msg != nil {
		s = msg.AsString()
	}
	if nl != nil {
		if b, ok := nl.(data.AsBool); ok {
			if okv, _ := b.AsBool(); okv {
				s += "\n"
			}
		}
	}
	os.Stdout.WriteString(s)
	return data.NewNullValue(), nil
}
