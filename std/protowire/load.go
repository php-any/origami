package protowire

import (
	"github.com/php-any/origami/data"
)

// Load registers all protowire functions and constants into the VM.
func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewParseRawFieldsFunction(),
		NewEncodeVarintFunction(),
		NewEncodeTagFunction(),
		NewEncodeBytesFunction(),
		NewEncodeFixed32Function(),
		NewEncodeFixed64Function(),
	} {
		vm.AddFunc(fun)
	}

	// Wire type constants
	constants := map[string]int{
		"PROTOWIRE_WIRE_VARINT":           WireVarint,
		"PROTOWIRE_WIRE_FIXED64":          WireFixed64,
		"PROTOWIRE_WIRE_LENGTH_DELIMITED": WireLengthDelimited,
		"PROTOWIRE_WIRE_START_GROUP":      WireStartGroup,
		"PROTOWIRE_WIRE_END_GROUP":        WireEndGroup,
		"PROTOWIRE_WIRE_FIXED32":          WireFixed32,
	}
	for name, val := range constants {
		vm.SetConstant(name, data.NewIntValue(val))
	}
}
