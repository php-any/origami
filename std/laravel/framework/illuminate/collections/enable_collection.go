package collections

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

func collectionMissing(ctx data.Context) (data.GetValue, data.Control) {
	name := ""
	if v := kit.Arg(ctx, 0); v != nil {
		name = v.AsString()
	}
	// 对齐 Macroable::__call：抛 BadMethodCallException（LogicException 的子类）。
	return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Method Illuminate\\Support\\Collection::%s does not exist.", name), "BadMethodCallException")
}
