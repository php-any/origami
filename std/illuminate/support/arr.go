package support

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const arrClassName = "Illuminate\\Support\\Arr"

// ArrClass 实现 Illuminate\Support\Arr 静态工具（vendor 加速）。
type ArrClass struct {
	node.Node
	methods map[string]data.Method
}

func NewArrClass() data.ClassStmt {
	c := &ArrClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *ArrClass) GetName() string                               { return arrClassName }
func (c *ArrClass) GetExtend() *string                            { return nil }
func (c *ArrClass) GetImplements() []string                       { return nil }
func (c *ArrClass) GetProperty(string) (data.Property, bool)      { return nil, false }
func (c *ArrClass) GetPropertyList() []data.Property              { return nil }
func (c *ArrClass) GetConstruct() data.Method                     { return nil }
func (c *ArrClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ArrClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *ArrClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *ArrClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *ArrClass) register() {
	add := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = newStaticMethod(name, params, optionalFrom, fn, false)
	}
	addRef := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = newStaticMethod(name, params, optionalFrom, fn, true)
	}
	add("accessible", []string{"value"}, -1, arrAccessible)
	add("exists", []string{"array", "key"}, -1, arrExists)
	add("has", []string{"array", "keys"}, -1, arrHas)
	add("get", []string{"array", "key", "default"}, 2, arrGet)
	addRef("set", []string{"array", "key", "value"}, -1, arrSet)
	addRef("add", []string{"array", "key", "value"}, -1, arrAdd)
	addRef("forget", []string{"array", "keys"}, -1, arrForget)
	add("only", []string{"array", "keys"}, -1, arrOnly)
	add("except", []string{"array", "keys"}, -1, arrExcept)
	add("first", []string{"array", "callback", "default"}, 1, arrFirst)
	add("last", []string{"array", "callback", "default"}, 1, arrLast)
	add("wrap", []string{"value"}, -1, arrWrap)
	add("from", []string{"items"}, -1, arrFrom)
	add("partition", []string{"array", "callback"}, -1, arrPartition)
	add("collapse", []string{"array"}, -1, arrCollapse)
	add("flatten", []string{"array", "depth"}, 1, arrFlatten)
	add("dot", []string{"array", "prepend"}, 1, arrDot)
	add("undot", []string{"array"}, -1, arrUndot)
	add("pluck", []string{"array", "value", "key"}, 2, arrPluck)
	add("map", []string{"array", "callback"}, -1, arrMap)
	add("mapWithKeys", []string{"array", "callback"}, -1, arrMapWithKeys)
	add("isList", []string{"array"}, -1, arrIsList)
	add("isAssoc", []string{"array"}, -1, arrIsAssoc)
	add("where", []string{"array", "callback"}, -1, arrWhere)
	add("prepend", []string{"array", "value", "key"}, 2, arrPrepend)
	addRef("pull", []string{"array", "key", "default"}, 2, arrPull)
	add("query", []string{"array"}, -1, arrQuery)
	add("random", []string{"array", "number"}, 1, arrRandom)
	add("shuffle", []string{"array", "seed"}, 1, arrShuffle)
	add("sort", []string{"array", "callback"}, 1, arrSort)
	add("sortDesc", []string{"array", "options"}, 1, arrSortDesc)
	add("sortRecursive", []string{"array", "options"}, 1, arrSortRecursive)
	add("toCssClasses", []string{"array"}, -1, arrToCssClasses)
	add("toCssStyles", []string{"array"}, -1, arrToCssStyles)
	add("join", []string{"array", "glue", "finalGlue"}, 2, arrJoin)
	add("keyBy", []string{"array", "keyBy"}, -1, arrKeyBy)
	add("divide", []string{"array"}, -1, arrDivide)
}

func newStaticMethod(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control), firstByRef bool) data.Method {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		var def data.GetValue
		if optionalFrom >= 0 && i >= optionalFrom {
			def = data.NewNullValue()
		}
		if i == 0 && firstByRef {
			ps[i] = node.NewParameterReference(nil, p, i, def, nil)
		} else {
			ps[i] = node.NewParameter(nil, p, i, def, nil)
		}
		vs[i] = node.NewVariable(nil, p, i, nil)
	}
	return &arrMethod{name: name, params: ps, vars: vs, fn: fn}
}

type arrMethod struct {
	name   string
	params []data.GetValue
	vars   []data.Variable
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *arrMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *arrMethod) GetName() string                                     { return m.name }
func (m *arrMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *arrMethod) GetIsStatic() bool                                   { return true }
func (m *arrMethod) GetParams() []data.GetValue                          { return m.params }
func (m *arrMethod) GetVariables() []data.Variable                       { return m.vars }
func (m *arrMethod) GetReturnType() data.Types                           { return nil }
