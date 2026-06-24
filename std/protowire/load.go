package protowire

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/protowire/annotation"
)

// Load registers the Protowire class, the @Field annotation class,
// and wire type constants into the VM.
func Load(vm data.VM) {
	annotation.Load(vm)
	vm.AddClass(NewProtowireClass())

	constants := map[string]int{
		"PROTOWIRE_VARINT":           WireVarint,
		"PROTOWIRE_FIXED64":          WireFixed64,
		"PROTOWIRE_LENGTH_DELIMITED": WireLengthDelimited,
		"PROTOWIRE_START_GROUP":      WireStartGroup,
		"PROTOWIRE_END_GROUP":        WireEndGroup,
		"PROTOWIRE_FIXED32":          WireFixed32,
	}
	for name, val := range constants {
		vm.SetConstant(name, data.NewIntValue(val))
	}
}
