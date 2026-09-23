package routing

import "github.com/php-any/origami/data"

func Load(vm data.VM) {
	vm.AddClass(NewUrlGeneratorClass())
}

const ComposerName = "illuminate/routing"
const TargetVersion = "v13.23.0"
