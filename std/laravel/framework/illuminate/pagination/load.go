package pagination

import "github.com/php-any/origami/data"

// Load：Paginator 视图/URL 工厂未齐前不占名。
func Load(vm data.VM) {
	_ = vm
	_ = NewAbstractPaginatorClass
	_ = NewLengthAwarePaginatorClass
	_ = NewPaginatorClass
}
