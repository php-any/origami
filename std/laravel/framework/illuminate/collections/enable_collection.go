package collections

import (
	"fmt"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

// maybeRegisterCollection：默认关闭（方法面未齐时会挡住 vendor）；ORIGAMI_STD_COLLECTION=1 启用。
func maybeRegisterCollection(vm data.VM) {
	switch os.Getenv("ORIGAMI_STD_COLLECTION") {
	case "1", "true", "yes":
		vm.AddInterface(NewEnumerableInterface())
		vm.AddClass(NewCollectionClass())
	}
}

func collectionMissing(ctx data.Context) (data.GetValue, data.Control) {
	name := ""
	if v := kit.Arg(ctx, 0); v != nil {
		name = v.AsString()
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Method Illuminate\\Support\\Collection::%s does not exist.", name))
}
