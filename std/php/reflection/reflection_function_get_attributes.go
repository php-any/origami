package reflection

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strings"
)

type ReflectionFunctionGetAttributesMethod struct {
	ReflectionMethodGetAttributesMethod
}

func (*ReflectionFunctionGetAttributesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	object, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewArrayValueFromSlots(nil), nil
	}
	var function data.FuncStmt
	if value, _ := object.ObjectValue.GetProperty("_function"); value != nil {
		if value, ok := value.(*data.FuncValue); ok {
			function = value.Value
		}
	}
	declaration, ok := function.(*node.FunctionStatement)
	if !ok {
		return data.NewArrayValueFromSlots(nil), nil
	}
	filter, _ := ctx.GetIndexValue(0)
	flags, _ := ctx.GetIndexValue(1)
	name := ""
	if _, isNull := filter.(*data.NullValue); filter != nil && !isNull {
		name = filter.AsString()
	}
	instanceof := false
	if flags, ok := flags.(*data.IntValue); ok {
		if flags.Value != 0 && flags.Value != 2 {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("ReflectionFunction::getAttributes(): Argument #2 ($flags) must be a valid attribute filter flag"), "ValueError")
		}
		instanceof = flags.Value&2 != 0
	}
	counts := make(map[string]int)
	for _, attribute := range declaration.PHPAttributes {
		counts[strings.ToLower(attribute.Name)]++
	}
	result := []data.Value{}
	for _, attribute := range declaration.PHPAttributes {
		if name != "" && !strings.EqualFold(name, attribute.Name) {
			if !instanceof {
				continue
			}
			class, ctl := ctx.GetVM().GetOrLoadClass(attribute.Name)
			if ctl != nil {
				return nil, ctl
			}
			if !attributeMatchesName(ctx, class, name) {
				continue
			}
		}
		result = append(result, data.NewClassValue(&ReflectionAttributeClass{declaration: attribute, target: 2, repeated: counts[strings.ToLower(attribute.Name)] > 1}, ctx.CreateBaseContext()))
	}
	return data.NewArrayValue(result), nil
}
