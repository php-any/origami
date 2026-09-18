package routing

import (
	"fmt"
	"sort"
	"strings"

	"github.com/php-any/origami/data"
)

const (
	routeCollectionName = "Symfony\\Component\\Routing\\RouteCollection"
	aliasName           = "Symfony\\Component\\Routing\\Alias"
)

func newAliasClass() *rtClass {
	c := newRtClass(aliasName, nil, nil, []data.Property{
		privProp("id", data.NewStringValue("")),
		privProp("deprecation", emptyArray()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{param("id", 0, nil, nil)},
		[]data.Variable{variable("id", 0, nil)},
		func(ctx data.Context) (data.GetValue, data.Control) {
			setProp(rtSelf(ctx), "id", data.NewStringValue(argString(ctx, 0, "")))
			setProp(rtSelf(ctx), "deprecation", phpList())
			return data.NewNullValue(), nil
		}))
	c.add(meth("withId", []string{"id"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		neu, ctl := instantiate(ctx, aliasName, data.NewStringValue(argString(ctx, 0, "")))
		if ctl != nil {
			return nil, ctl
		}
		setProp(neu, "deprecation", cloneArray(propArray(cv, "deprecation")))
		return neu, nil
	}))
	c.add(meth("getId", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "id", "")), nil
	}))
	c.add(meth("setDeprecated", []string{"package", "version", "message"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		msg := argString(ctx, 2, "")
		if msg != "" {
			if strings.ContainsAny(msg, "\r\n") || strings.Contains(msg, "*/") {
				return nil, throwNamed(ctx, exInvalidArgument, "Invalid characters found in deprecation template.")
			}
			if !strings.Contains(msg, "%alias_id%") {
				return nil, throwNamed(ctx, exInvalidArgument, "The deprecation template must contain the \"%alias_id%\" placeholder.")
			}
		}
		if msg == "" {
			msg = "The \"%alias_id%\" route alias is deprecated. You should stop using it, as it will be removed in the future."
		}
		setProp(cv, "deprecation", phpAssoc(
			"package", data.NewStringValue(argString(ctx, 0, "")),
			"version", data.NewStringValue(argString(ctx, 1, "")),
			"message", data.NewStringValue(msg),
		))
		return cv, nil
	}))
	c.add(meth("isDeprecated", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewBoolValue(len(propArray(rtSelf(ctx), "deprecation").List) > 0), nil
	}))
	c.add(meth("getDeprecation", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		dep := propArray(rtSelf(ctx), "deprecation")
		pkg, _ := assocGet(dep, "package")
		ver, _ := assocGet(dep, "version")
		msg, _ := assocGet(dep, "message")
		message := ""
		if msg != nil {
			message = strings.ReplaceAll(msg.AsString(), "%alias_id%", argString(ctx, 0, ""))
		}
		return phpAssoc(
			"package", pkg,
			"version", ver,
			"message", data.NewStringValue(message),
		), nil
	}))
	return c
}

func newRouteCollectionClass() *rtClass {
	c := newRtClass(routeCollectionName, nil, []string{"IteratorAggregate", "Countable"}, []data.Property{
		privProp("routes", emptyArray()),
		privProp("aliases", emptyArray()),
		privProp("resources", emptyArray()),
		privProp("priorities", emptyArray()),
	})
	c.add(meth("__construct", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "routes", phpList())
		setProp(cv, "aliases", phpList())
		setProp(cv, "resources", phpList())
		setProp(cv, "priorities", phpList())
		return data.NewNullValue(), nil
	}))
	c.add(meth("add", []string{"name", "route", "priority"}, collectionAdd))
	c.add(meth("all", nil, collectionAll))
	c.add(meth("get", []string{"name"}, collectionGet))
	c.add(meth("remove", []string{"name"}, collectionRemove))
	c.add(meth("addCollection", []string{"collection"}, collectionAddCollection))
	c.add(meth("addPrefix", []string{"prefix", "defaults", "requirements"}, collectionAddPrefix))
	c.add(meth("addNamePrefix", []string{"prefix"}, collectionAddNamePrefix))
	c.add(meth("setHost", []string{"pattern", "defaults", "requirements"}, collectionSetHost))
	c.add(meth("setCondition", []string{"condition"}, collectionSetCondition))
	c.add(meth("addDefaults", []string{"defaults"}, collectionAddDefaults))
	c.add(meth("addRequirements", []string{"requirements"}, collectionAddRequirements))
	c.add(meth("addOptions", []string{"options"}, collectionAddOptions))
	c.add(meth("setSchemes", []string{"schemes"}, collectionSetSchemes))
	c.add(meth("setMethods", []string{"methods"}, collectionSetMethods))
	c.add(meth("getResources", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		src := propArray(rtSelf(ctx), "resources")
		vals := make([]data.Value, 0, len(src.List))
		for _, z := range src.List {
			if z != nil && z.Value != nil {
				vals = append(vals, z.Value)
			}
		}
		return phpList(vals...), nil
	}))
	c.add(meth("addResource", []string{"resource"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		res := arg(ctx, 0)
		if res == nil {
			return data.NewNullValue(), nil
		}
		key := res.AsString()
		resources := cloneArray(propArray(cv, "resources"))
		if !assocHas(resources, key) {
			assocSet(resources, key, res)
			setProp(cv, "resources", resources)
		}
		return data.NewNullValue(), nil
	}))
	c.add(meth("addAlias", []string{"name", "alias"}, collectionAddAlias))
	c.add(meth("getAliases", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "aliases"), nil
	}))
	c.add(meth("getAlias", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		v, ok := assocGet(propArray(rtSelf(ctx), "aliases"), argString(ctx, 0, ""))
		if !ok {
			return data.NewNullValue(), nil
		}
		return v, nil
	}))
	c.add(meth("getPriority", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		v, ok := assocGet(propArray(rtSelf(ctx), "priorities"), argString(ctx, 0, ""))
		if !ok {
			return data.NewNullValue(), nil
		}
		return v, nil
	}))
	c.add(meth("count", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewIntValue(len(propArray(rtSelf(ctx), "routes").List)), nil
	}))
	c.add(meth("getIterator", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		all, ctl := collectionAll(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return newArrayIterator(ctx, all.(data.Value))
	}))
	return c
}

func collectionAdd(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	name := argString(ctx, 0, "")
	route := arg(ctx, 1)
	priority := argInt(ctx, 2, 0)
	routes := cloneArray(propArray(cv, "routes"))
	aliases := cloneArray(propArray(cv, "aliases"))
	priorities := cloneArray(propArray(cv, "priorities"))
	assocUnset(routes, name)
	assocUnset(aliases, name)
	assocUnset(priorities, name)
	assocSet(routes, name, route)
	if priority != 0 {
		assocSet(priorities, name, data.NewIntValue(priority))
	}
	setProp(cv, "routes", routes)
	setProp(cv, "aliases", aliases)
	setProp(cv, "priorities", priorities)
	return data.NewNullValue(), nil
}

func collectionAll(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	routes := cloneArray(propArray(cv, "routes"))
	pri := propArray(cv, "priorities")
	if len(pri.List) == 0 {
		return routes, nil
	}
	type item struct {
		name string
		val  data.Value
		pri  int
		ord  int
	}
	items := make([]item, 0, len(routes.List))
	for i, z := range routes.List {
		if z == nil {
			continue
		}
		p := 0
		if pv, ok := assocGet(pri, z.Name); ok {
			if iv, ok := pv.(data.AsInt); ok {
				p, _ = iv.AsInt()
			}
		}
		items = append(items, item{name: z.Name, val: z.Value, pri: p, ord: i})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].pri != items[j].pri {
			return items[i].pri > items[j].pri
		}
		return items[i].ord < items[j].ord
	})
	out := &data.ArrayValue{List: make([]*data.ZVal, 0, len(items))}
	for _, it := range items {
		out.List = append(out.List, data.NewNamedZVal(it.name, it.val))
	}
	return out, nil
}

func collectionGet(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	name := argString(ctx, 0, "")
	aliases := propArray(cv, "aliases")
	visited := []string{}
	for {
		aliasVal, ok := assocGet(aliases, name)
		if !ok {
			break
		}
		for i, v := range visited {
			if v == name {
				path := append(visited[i:], name)
				vals := make([]data.Value, len(path))
				for j, s := range path {
					vals[j] = data.NewStringValue(s)
				}
				return nil, throwNamedArgs(ctx, exRouteCircularReference, data.NewStringValue(name), phpList(vals...))
			}
		}
		aliasCV := asClassValue(aliasVal)
		if aliasCV == nil {
			break
		}
		visited = append(visited, name)
		id, ctl := callNamed(aliasCV, "getId")
		if ctl != nil {
			return nil, ctl
		}
		name = id.AsString()
	}
	v, ok := assocGet(propArray(cv, "routes"), name)
	if !ok {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func collectionRemove(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	names := toStringSlice(arg(ctx, 0))
	routes := cloneArray(propArray(cv, "routes"))
	aliases := cloneArray(propArray(cv, "aliases"))
	priorities := cloneArray(propArray(cv, "priorities"))
	removed := map[string]bool{}
	for _, n := range names {
		if assocHas(routes, n) {
			removed[n] = true
		}
		assocUnset(routes, n)
		assocUnset(priorities, n)
		assocUnset(aliases, n)
	}
	if len(removed) > 0 {
		keep := &data.ArrayValue{}
		for _, z := range aliases.List {
			if z == nil {
				continue
			}
			if acv := asClassValue(z.Value); acv != nil {
				id, _ := callNamed(acv, "getId")
				if id != nil && removed[id.AsString()] {
					continue
				}
			}
			keep.List = append(keep.List, z)
		}
		aliases = keep
	}
	setProp(cv, "routes", routes)
	setProp(cv, "aliases", aliases)
	setProp(cv, "priorities", priorities)
	return data.NewNullValue(), nil
}

func collectionAddCollection(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	other := asClassValue(arg(ctx, 0))
	if other == nil {
		return data.NewNullValue(), nil
	}
	all, ctl := callNamed(other, "all")
	if ctl != nil {
		return nil, ctl
	}
	allArr, _ := all.(*data.ArrayValue)
	routes := cloneArray(propArray(cv, "routes"))
	aliases := cloneArray(propArray(cv, "aliases"))
	priorities := cloneArray(propArray(cv, "priorities"))
	otherPri := propArray(other, "priorities")
	if allArr != nil {
		for i, z := range allArr.List {
			if z == nil {
				continue
			}
			name := arrayKey(z, i)
			assocUnset(routes, name)
			assocUnset(priorities, name)
			assocUnset(aliases, name)
			assocSet(routes, name, z.Value)
			if pv, ok := assocGet(otherPri, name); ok {
				assocSet(priorities, name, pv)
			}
		}
	}
	otherAliases := propArray(other, "aliases")
	for i, z := range otherAliases.List {
		if z == nil {
			continue
		}
		name := arrayKey(z, i)
		assocUnset(routes, name)
		assocUnset(priorities, name)
		assocUnset(aliases, name)
		assocSet(aliases, name, z.Value)
	}
	setProp(cv, "routes", routes)
	setProp(cv, "aliases", aliases)
	setProp(cv, "priorities", priorities)
	return data.NewNullValue(), nil
}

func eachRoute(cv *data.ClassValue, fn func(*data.ClassValue) data.Control) data.Control {
	routes := propArray(cv, "routes")
	for _, z := range routes.List {
		if z == nil {
			continue
		}
		r := asClassValue(z.Value)
		if r == nil {
			continue
		}
		if ctl := fn(r); ctl != nil {
			return ctl
		}
	}
	return nil
}

func collectionAddPrefix(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	prefix := strings.Trim(strings.TrimSpace(argString(ctx, 0, "")), "/")
	if prefix == "" {
		return data.NewNullValue(), nil
	}
	defs := argArray(ctx, 1)
	reqs := argArray(ctx, 2)
	return data.NewNullValue(), eachRoute(cv, func(r *data.ClassValue) data.Control {
		path := "/" + prefix + propString(r, "path", "/")
		routeSetPath(r, path)
		routeAddDefaults(r, defs)
		return routeAddRequirements(r, reqs)
	})
}

func collectionAddNamePrefix(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	prefix := argString(ctx, 0, "")
	routes := propArray(cv, "routes")
	priorities := propArray(cv, "priorities")
	aliases := propArray(cv, "aliases")
	prefixed := &data.ArrayValue{}
	prefixedPri := &data.ArrayValue{}
	for _, z := range routes.List {
		if z == nil {
			continue
		}
		newName := prefix + z.Name
		if r := asClassValue(z.Value); r != nil {
			if canon := routeGetDefault(r, "_canonical_route"); !isNull(canon) {
				routeSetDefault(r, "_canonical_route", data.NewStringValue(prefix+canon.AsString()))
			}
		}
		prefixed.List = append(prefixed.List, data.NewNamedZVal(newName, z.Value))
		if pv, ok := assocGet(priorities, z.Name); ok {
			assocSet(prefixedPri, newName, pv)
		}
	}
	prefixedAliases := &data.ArrayValue{}
	for _, z := range aliases.List {
		if z == nil {
			continue
		}
		acv := asClassValue(z.Value)
		if acv == nil {
			continue
		}
		id, _ := callNamed(acv, "getId")
		target := ""
		if id != nil {
			target = id.AsString()
		}
		if assocHas(routes, target) || assocHas(aliases, target) {
			target = prefix + target
		}
		neu, ctl := callNamed(acv, "withId", data.NewStringValue(target))
		if ctl != nil {
			return nil, ctl
		}
		prefixedAliases.List = append(prefixedAliases.List, data.NewNamedZVal(prefix+z.Name, neu))
	}
	setProp(cv, "routes", prefixed)
	setProp(cv, "priorities", prefixedPri)
	setProp(cv, "aliases", prefixedAliases)
	return data.NewNullValue(), nil
}

func collectionSetHost(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	pattern := argString(ctx, 0, "")
	defs := argArray(ctx, 1)
	reqs := argArray(ctx, 2)
	return data.NewNullValue(), eachRoute(cv, func(r *data.ClassValue) data.Control {
		routeSetHost(r, pattern)
		routeAddDefaults(r, defs)
		return routeAddRequirements(r, reqs)
	})
}

func collectionSetCondition(ctx data.Context) (data.GetValue, data.Control) {
	cond := argString(ctx, 0, "")
	return data.NewNullValue(), eachRoute(rtSelf(ctx), func(r *data.ClassValue) data.Control {
		setProp(r, "condition", data.NewStringValue(cond))
		setProp(r, "compiled", data.NewNullValue())
		return nil
	})
}

func collectionAddDefaults(ctx data.Context) (data.GetValue, data.Control) {
	defs := argArray(ctx, 0)
	if len(defs.List) == 0 {
		return data.NewNullValue(), nil
	}
	return data.NewNullValue(), eachRoute(rtSelf(ctx), func(r *data.ClassValue) data.Control {
		routeAddDefaults(r, defs)
		return nil
	})
}

func collectionAddRequirements(ctx data.Context) (data.GetValue, data.Control) {
	reqs := argArray(ctx, 0)
	if len(reqs.List) == 0 {
		return data.NewNullValue(), nil
	}
	return data.NewNullValue(), eachRoute(rtSelf(ctx), func(r *data.ClassValue) data.Control {
		return routeAddRequirements(r, reqs)
	})
}

func collectionAddOptions(ctx data.Context) (data.GetValue, data.Control) {
	opts := argArray(ctx, 0)
	if len(opts.List) == 0 {
		return data.NewNullValue(), nil
	}
	return data.NewNullValue(), eachRoute(rtSelf(ctx), func(r *data.ClassValue) data.Control {
		routeAddOptions(r, opts)
		return nil
	})
}

func collectionSetSchemes(ctx data.Context) (data.GetValue, data.Control) {
	v := arg(ctx, 0)
	return data.NewNullValue(), eachRoute(rtSelf(ctx), func(r *data.ClassValue) data.Control {
		routeSetSchemes(r, v)
		return nil
	})
}

func collectionSetMethods(ctx data.Context) (data.GetValue, data.Control) {
	v := arg(ctx, 0)
	return data.NewNullValue(), eachRoute(rtSelf(ctx), func(r *data.ClassValue) data.Control {
		routeSetMethods(r, v)
		return nil
	})
}

func collectionAddAlias(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	name := argString(ctx, 0, "")
	alias := argString(ctx, 1, "")
	if name == alias {
		return nil, throwNamed(ctx, exInvalidArgument, fmt.Sprintf("Route alias \"%s\" can not reference itself.", name))
	}
	routes := cloneArray(propArray(cv, "routes"))
	priorities := cloneArray(propArray(cv, "priorities"))
	assocUnset(routes, name)
	assocUnset(priorities, name)
	obj, ctl := instantiate(ctx, aliasName, data.NewStringValue(alias))
	if ctl != nil {
		return nil, ctl
	}
	aliases := cloneArray(propArray(cv, "aliases"))
	assocSet(aliases, name, obj)
	setProp(cv, "routes", routes)
	setProp(cv, "priorities", priorities)
	setProp(cv, "aliases", aliases)
	return obj, nil
}
