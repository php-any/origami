package routing

import (
	"github.com/php-any/origami/data"
)

const attributeRouteName = "Symfony\\Component\\Routing\\Attribute\\Route"

func newAttributeRouteClass() *rtClass {
	c := newRtClass(attributeRouteName, nil, nil, []data.Property{
		pubProp("path", data.NewNullValue()),
		pubProp("name", data.NewNullValue()),
		pubProp("requirements", emptyArray()),
		pubProp("options", emptyArray()),
		pubProp("defaults", emptyArray()),
		pubProp("host", data.NewNullValue()),
		pubProp("methods", emptyArray()),
		pubProp("schemes", emptyArray()),
		pubProp("condition", data.NewNullValue()),
		pubProp("priority", data.NewNullValue()),
		pubProp("envs", emptyArray()),
		pubProp("aliases", emptyArray()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{
			param("path", 0, data.NewNullValue(), nil),
			param("name", 1, data.NewNullValue(), nil),
			param("requirements", 2, emptyArray(), nil),
			param("options", 3, emptyArray(), nil),
			param("defaults", 4, emptyArray(), nil),
			param("host", 5, data.NewNullValue(), nil),
			param("methods", 6, emptyArray(), nil),
			param("schemes", 7, emptyArray(), nil),
			param("condition", 8, data.NewNullValue(), nil),
			param("priority", 9, data.NewNullValue(), nil),
			param("locale", 10, data.NewNullValue(), nil),
			param("format", 11, data.NewNullValue(), nil),
			param("utf8", 12, data.NewNullValue(), nil),
			param("stateless", 13, data.NewNullValue(), nil),
			param("env", 14, data.NewNullValue(), nil),
			param("alias", 15, emptyArray(), nil),
		},
		[]data.Variable{
			variable("path", 0, nil), variable("name", 1, nil), variable("requirements", 2, nil),
			variable("options", 3, nil), variable("defaults", 4, nil), variable("host", 5, nil),
			variable("methods", 6, nil), variable("schemes", 7, nil), variable("condition", 8, nil),
			variable("priority", 9, nil), variable("locale", 10, nil), variable("format", 11, nil),
			variable("utf8", 12, nil), variable("stateless", 13, nil), variable("env", 14, nil),
			variable("alias", 15, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			setProp(cv, "path", nullOrValue(arg(ctx, 0)))
			setProp(cv, "name", nullOrValue(arg(ctx, 1)))
			setProp(cv, "requirements", argArray(ctx, 2))
			setProp(cv, "options", argArray(ctx, 3))
			defs := argArray(ctx, 4)
			setProp(cv, "host", nullOrValue(arg(ctx, 5)))
			setProp(cv, "methods", stringsToArray(toStringSlice(arg(ctx, 6))))
			setProp(cv, "schemes", stringsToArray(toStringSlice(arg(ctx, 7))))
			setProp(cv, "condition", nullOrValue(arg(ctx, 8)))
			setProp(cv, "priority", nullOrValue(arg(ctx, 9)))
			if loc := arg(ctx, 10); loc != nil && !isNull(loc) {
				assocSet(defs, "_locale", loc)
			}
			if format := arg(ctx, 11); format != nil && !isNull(format) {
				assocSet(defs, "_format", format)
			}
			opts := propArray(cv, "options")
			if utf8 := arg(ctx, 12); utf8 != nil && !isNull(utf8) {
				assocSet(opts, "utf8", utf8)
				setProp(cv, "options", opts)
			}
			if stateless := arg(ctx, 13); stateless != nil && !isNull(stateless) {
				assocSet(defs, "_stateless", stateless)
			}
			setProp(cv, "defaults", defs)
			env := arg(ctx, 14)
			if env == nil || isNull(env) {
				setProp(cv, "envs", phpList())
			} else {
				setProp(cv, "envs", stringsToArray(toStringSlice(env)))
			}
			alias := arg(ctx, 15)
			if av, ok := alias.(*data.ArrayValue); ok {
				setProp(cv, "aliases", av)
			} else if alias != nil && !isNull(alias) {
				setProp(cv, "aliases", phpList(alias))
			} else {
				setProp(cv, "aliases", phpList())
			}
			return data.NewNullValue(), nil
		}))
	return c
}

func nullOrValue(v data.Value) data.Value {
	if v == nil {
		return data.NewNullValue()
	}
	return v
}
