package node

import "github.com/php-any/origami/data"

type privateMethodIndex struct {
	count   int
	methods map[string]data.Method
}

// Lexical private dispatch only needs private instance methods. Most scopes
// have none; avoid building and probing their complete method-name index.
func (c *ClassStatement) privateMethodByKey(key string) (data.Method, bool) {
	count := len(c.Methods) + len(c.StaticMethods)
	index := c.privateLookup.Load()
	if index == nil || index.count != count {
		index = &privateMethodIndex{count: count}
		add := func(name string, method data.Method) {
			if method == nil || method.GetModifier() != data.ModifierPrivate {
				return
			}
			if index.methods == nil {
				index.methods = make(map[string]data.Method)
			}
			index.methods[data.MethodLookupKey(name)] = method
		}
		for name, method := range c.Methods {
			add(name, method)
		}
		if _, declared := c.Methods["__construct"]; !declared {
			add("__construct", c.Construct)
		}
		c.privateLookup.Store(index)
	}
	method, found := index.methods[key]
	return method, found
}
