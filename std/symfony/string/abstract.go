package sfstring

import "github.com/php-any/origami/data"

func NewAbstractStringClass() data.ClassStmt {
	c := &strClass{
		name:       abstractStringName,
		parent:     "",
		kind:       kindGrapheme,
		methods:    map[string]data.Method{},
		implements: []string{"Stringable", "JsonSerializable"},
		abstract:   true,
	}
	initStringProps(c)
	c.add("wrap", true, []data.GetValue{p("values", 0, nil)}, []data.Variable{v("values", 0)}, methodWrap)
	c.add("unwrap", true, []data.GetValue{p("values", 0, nil)}, []data.Variable{v("values", 0)}, methodUnwrap)
	return c
}

func NewAbstractUnicodeStringClass() data.ClassStmt {
	c := &strClass{
		name:       abstractUnicodeName,
		parent:     abstractStringName,
		kind:       kindGrapheme,
		methods:    map[string]data.Method{},
		implements: nil,
		abstract:   true,
	}
	initStringProps(c)
	c.add("fromCodePoints", true, []data.GetValue{variadic("codes", 0)}, []data.Variable{v("codes", 0)}, methodFromCodePoints)
	c.add("wrap", true, []data.GetValue{p("values", 0, nil)}, []data.Variable{v("values", 0)}, methodWrap)
	c.add("unwrap", true, []data.GetValue{p("values", 0, nil)}, []data.Variable{v("values", 0)}, methodUnwrap)
	return c
}
