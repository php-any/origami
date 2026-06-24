package container

import "github.com/php-any/origami/data"

// Binding 描述 abstract → concrete 的注册关系。
type Binding struct {
	Abstract string
	Concrete string
	Factory  data.GetValue
	Lifetime Lifetime
}

func (b *Binding) isFactory() bool {
	return b != nil && b.Factory != nil
}
