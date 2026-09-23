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
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Method Illuminate\\Support\\Collection::%s does not exist.", name))
}
