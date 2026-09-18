package core

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	backedEnumTryFromMethod = &BackedEnumTryFromMethod{}
	backedEnumFromMethod    = &BackedEnumFromMethod{}
	backedEnumCasesMethod   = &BackedEnumCasesMethod{}
)

func backedEnumCalledClass(ctx data.Context) data.ClassStmt {
	cmc, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil
	}
	if cmc.StaticClass != nil {
		return cmc.StaticClass
	}
	if cmc.ClassValue != nil {
		return cmc.ClassValue.Class
	}
	return nil
}

func collectBackedEnumCases(cls data.ClassStmt) []data.Value {
	cs, ok := cls.(*node.ClassStatement)
	if !ok || cs == nil {
		return nil
	}
	enumName := cs.GetName()
	var cases []data.Value
	cs.StaticProperty.Range(func(key, value any) bool {
		name, _ := key.(string)
		if name == "" {
			return true
		}
		if _, isConst := cs.StaticProperties[name]; isConst {
			return true
		}
		cv, ok := value.(*data.ClassValue)
		if !ok || cv == nil || cv.Class == nil {
			return true
		}
		if cv.Class.GetName() != enumName {
			return true
		}
		cases = append(cases, cv)
		return true
	})
	return cases
}

func enumBackingValue(cv *data.ClassValue) data.Value {
	if cv == nil {
		return nil
	}
	v, ctl := cv.GetProperty("value")
	if ctl != nil || v == nil {
		return nil
	}
	return unwrapEnumValue(v)
}

func enumBackingIdentical(needle, backing data.Value) bool {
	if needle == nil || backing == nil {
		return false
	}
	switch n := needle.(type) {
	case *data.StringValue:
		b, ok := backing.(*data.StringValue)
		return ok && n.Value == b.Value
	case *data.IntValue:
		b, ok := backing.(*data.IntValue)
		return ok && n.Value == b.Value
	default:
		return false
	}
}

func unwrapEnumValue(v data.Value) data.Value {
	for i := 0; i < 4 && v != nil; i++ {
		switch x := v.(type) {
		case *data.ZValValue:
			if x.ZVal == nil {
				return nil
			}
			v = x.ZVal.Value
		default:
			return v
		}
	}
	return v
}

func backedEnumLookup(ctx data.Context, method string) (data.Value, data.Control) {
	cls := backedEnumCalledClass(ctx)
	if cls == nil {
		return data.NewNullValue(), nil
	}
	val, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("%s::%s() expects exactly 1 argument, 0 given", cls.GetName(), method), "ArgumentCountError")
	}
	val = unwrapEnumValue(val)
	if val == nil {
		return data.NewNullValue(), nil
	}
	switch val.(type) {
	case *data.StringValue, *data.IntValue:
	case *data.NullValue:
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("%s::%s(): Argument #1 ($value) must be of type string|int, null given", cls.GetName(), method), "TypeError")
	default:
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("%s::%s(): Argument #1 ($value) must be of type string|int, %T given", cls.GetName(), method, val), "TypeError")
	}
	for _, c := range collectBackedEnumCases(cls) {
		cv, ok := c.(*data.ClassValue)
		if !ok {
			continue
		}
		if enumBackingIdentical(val, enumBackingValue(cv)) {
			return cv, nil
		}
	}
	return data.NewNullValue(), nil
}

// BackedEnumTryFromMethod 实现 BackedEnum::tryFrom(int|string $value): ?static
type BackedEnumTryFromMethod struct{}

func (m *BackedEnumTryFromMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return backedEnumLookup(ctx, "tryFrom")
}
func (m *BackedEnumTryFromMethod) GetName() string            { return "tryFrom" }
func (m *BackedEnumTryFromMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *BackedEnumTryFromMethod) GetIsStatic() bool          { return true }
func (m *BackedEnumTryFromMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *BackedEnumTryFromMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *BackedEnumTryFromMethod) GetReturnType() data.Types { return nil }

// BackedEnumFromMethod 实现 BackedEnum::from(int|string $value): static
type BackedEnumFromMethod struct{}

func (m *BackedEnumFromMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	found, ctl := backedEnumLookup(ctx, "from")
	if ctl != nil {
		return nil, ctl
	}
	if found == nil {
		return nil, backedEnumFromValueError(ctx, nil)
	}
	if _, ok := found.(*data.NullValue); ok {
		return nil, backedEnumFromValueError(ctx, ctxIndexValue(ctx, 0))
	}
	return found, nil
}
func (m *BackedEnumFromMethod) GetName() string            { return "from" }
func (m *BackedEnumFromMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *BackedEnumFromMethod) GetIsStatic() bool          { return true }
func (m *BackedEnumFromMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *BackedEnumFromMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *BackedEnumFromMethod) GetReturnType() data.Types { return nil }

func ctxIndexValue(ctx data.Context, i int) data.Value {
	v, ok := ctx.GetIndexValue(i)
	if !ok {
		return nil
	}
	return v
}

func backedEnumFromValueError(ctx data.Context, val data.Value) data.Control {
	cls := backedEnumCalledClass(ctx)
	name := "BackedEnum"
	if cls != nil {
		name = cls.GetName()
	}
	shown := "null"
	if val != nil {
		shown = val.AsString()
	}
	msg := fmt.Sprintf("%q is not a valid backing value for enum %q", shown, name)
	return data.NewErrorThrowByName(nil, fmt.Errorf("%s", msg), "ValueError")
}

// BackedEnumCasesMethod 实现 UnitEnum::cases(): array
type BackedEnumCasesMethod struct{}

func (m *BackedEnumCasesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cls := backedEnumCalledClass(ctx)
	if cls == nil {
		return data.NewArrayValue(nil), nil
	}
	return data.NewArrayValue(collectBackedEnumCases(cls)), nil
}
func (m *BackedEnumCasesMethod) GetName() string            { return "cases" }
func (m *BackedEnumCasesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *BackedEnumCasesMethod) GetIsStatic() bool          { return true }
func (m *BackedEnumCasesMethod) GetParams() []data.GetValue { return nil }
func (m *BackedEnumCasesMethod) GetVariables() []data.Variable {
	return nil
}
func (m *BackedEnumCasesMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
