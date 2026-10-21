package core

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"runtime"
	"sort"
	"strings"
)

type PHPInfoFunction struct{}

func (*PHPInfoFunction) GetName() string { return "phpinfo" }
func (*PHPInfoFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "flags", 0, data.NewIntValue(127), data.Int{})}
}
func (*PHPInfoFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "flags", 0, data.Int{})}
}
func (*PHPInfoFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	flags := 127
	if n, ok := v.(data.AsInt); ok {
		flags, _ = n.AsInt()
	}
	var out strings.Builder
	out.WriteString("phpinfo()\n")
	if flags&1 != 0 {
		version, _ := ctx.GetVM().GetConstant("PHP_VERSION")
		fmt.Fprintf(&out, "PHP Version => %s\nRuntime => Origami (%s)\nSystem => %s %s\n", version.AsString(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
		if sapi, ok := ctx.GetVM().GetConstant("PHP_SAPI"); ok {
			fmt.Fprintf(&out, "Server API => %s\n", sapi.AsString())
		}
		out.WriteString("Loaded Configuration File => (none)\nAdditional .ini files parsed => (none)\n")
	}
	if flags&4 != 0 {
		keys := make([]string, 0, len(iniDefaults))
		for key := range iniDefaults {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value, _ := IniGetInContext(ctx, key)
			fmt.Fprintf(&out, "%s => %s\n", key, value)
		}
	}
	if flags&8 != 0 {
		list := loadedExtensions(ctx)
		{
			for _, value := range list.ToValueList() {
				fmt.Fprintf(&out, "%s support => enabled\n", value.AsString())
			}
		}
	}
	return data.NewBoolValue(true), data.EmitOutput(ctx, out.String())
}
