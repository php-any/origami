package php

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"math"
	"strconv"
	"strings"
)

// Keys do not occupy value numbers. R shares a cell, r shares an object.
type phpSerState struct {
	ctx       data.Context
	n         int
	objects   map[*data.PropertyBag]int
	arrays    map[*data.ArrayValue]int
	refs      map[*data.ReferenceCell]int
	control   data.Control
	precision int
}

func newPhpSerState() *phpSerState {
	return &phpSerState{objects: map[*data.PropertyBag]int{}, arrays: map[*data.ArrayValue]int{}, refs: map[*data.ReferenceCell]int{}, precision: -1}
}
func makeSerializedString(s string) string { return "s:" + strconv.Itoa(len(s)) + ":\"" + s + "\";" }
func serializeSlot(slot *data.ZVal, st *phpSerState) (string, bool) {
	if cell := slot.ReferenceIdentity(); cell != nil {
		if n, found := st.refs[cell]; found {
			return "R:" + strconv.Itoa(n) + ";", true
		}
		if array, ok := slot.ReadValue().(*data.ArrayValue); ok {
			if n, found := st.arrays[array]; found {
				st.refs[cell] = n
				return "R:" + strconv.Itoa(n) + ";", true
			}
		}
		st.refs[cell] = st.n + 1
	}
	return phpSerializeValue(slot.ReadValue(), st)
}
func serializedFloat(value float64, precision int) string {
	if math.IsNaN(value) {
		return "NAN"
	}
	if math.IsInf(value, 1) {
		return "INF"
	}
	if math.IsInf(value, -1) {
		return "-INF"
	}
	digits := precision - 1
	if precision < 0 {
		digits = -1
	}
	if precision == 0 {
		digits = 0
	}
	s := strconv.FormatFloat(value, 'E', digits, 64)
	if mantissa, exponent, found := strings.Cut(s, "E"); found {
		e, _ := strconv.Atoi(exponent)
		limit := precision
		if limit < 0 {
			limit = 17
		}
		mantissa = strings.TrimRight(strings.TrimRight(mantissa, "0"), ".")
		// Do not trim the integer zero in 0E+00.
		if mantissa == "" || mantissa == "-" {
			mantissa += "0"
		}
		if e >= -4 && e < limit {
			sign := ""
			if strings.HasPrefix(mantissa, "-") {
				sign, mantissa = "-", mantissa[1:]
			}
			digits := strings.ReplaceAll(mantissa, ".", "")
			point := e + 1
			if point <= 0 {
				return sign + "0." + strings.Repeat("0", -point) + digits
			}
			if point >= len(digits) {
				return sign + digits + strings.Repeat("0", point-len(digits))
			}
			return sign + digits[:point] + "." + digits[point:]
		}
		if !strings.Contains(mantissa, ".") {
			mantissa += ".0"
		}
		sign := "+"
		if e < 0 {
			sign = "-"
			e = -e
		}
		return mantissa + "E" + sign + strconv.Itoa(e)
	}
	return s
}
func phpSerializeValue(value data.Value, st *phpSerState) (string, bool) {
	if st == nil {
		st = newPhpSerState()
	}
	if this, ok := value.(*data.ThisValue); ok {
		value = this.ClassValue
	}
	st.n++
	switch v := value.(type) {
	case nil, *data.NullValue:
		return "N;", true
	case *data.BoolValue:
		if v.Value {
			return "b:1;", true
		}
		return "b:0;", true
	case *data.IntValue:
		return "i:" + strconv.Itoa(v.Value) + ";", true
	case *data.FloatValue:
		return "d:" + serializedFloat(v.Value, st.precision) + ";", true
	case *data.StringValue:
		return makeSerializedString(v.Value), true
	case *data.ArrayValue:
		if len(st.arrays) >= 512 {
			return "", false
		}
		st.arrays[v] = st.n
		defer delete(st.arrays, v)
		var b strings.Builder
		b.WriteString("a:" + strconv.Itoa(v.Len()) + ":{")
		for view, i := v.View(), 0; i < view.Len(); i++ {
			slot := view.At(i)
			key := slot.PHPArrayKey(i)
			if n, ok := key.(*data.IntValue); ok {
				b.WriteString("i:" + strconv.Itoa(n.Value) + ";")
			} else {
				b.WriteString(makeSerializedString(key.AsString()))
			}
			encoded, ok := serializeSlot(slot, st)
			if !ok {
				return "", false
			}
			b.WriteString(encoded)
		}
		b.WriteByte('}')
		return b.String(), true
	case *data.ClassValue:
		if n, exists := st.objects[v.ObjectValue]; exists {
			return "r:" + strconv.Itoa(n) + ";", true
		}
		st.objects[v.ObjectValue] = st.n
		if flags, ok := v.Class.(interface{ DeclarationFlags() data.ClassFlags }); ok && flags.DeclarationFlags()&data.ClassEnum != 0 {
			caseName, ctl := v.GetProperty("name")
			if ctl != nil {
				st.control = ctl
				return "", false
			}
			name := v.Class.GetName() + ":" + caseName.AsString()
			return "E:" + strconv.Itoa(len(name)) + ":\"" + name + "\";", true
		}
		if _, hasModern := v.GetMethod("__serialize"); !hasModern && st.ctx != nil && data.NominalIsA(v.Class, "Serializable", st.ctx.GetVM()) {
			if hook, found := v.GetMethod("serialize"); found {
				raw, ctl := callSerializationHook(st.ctx, v, hook, nil)
				if ctl != nil {
					st.control = ctl
					return "", false
				}
				if _, isNull := raw.(*data.NullValue); isNull {
					return "N;", true
				}
				payload, ok := raw.(*data.StringValue)
				if !ok {
					st.control = data.NewErrorThrowByName(nil, fmt.Errorf("%s::serialize() must return a string or NULL", v.Class.GetName()), "Exception")
					return "", false
				}
				name := v.Class.GetName()
				return "C:" + strconv.Itoa(len(name)) + ":\"" + name + "\":" + strconv.Itoa(len(payload.Value)) + ":{" + payload.Value + "}", true
			}
		}
		props, ctl := serializedObjectProperties(v, st.ctx)
		if ctl != nil {
			st.control = ctl
			return "", false
		}
		if props == nil {
			return "N;", true
		}
		name := v.Class.GetName()
		if name == "__PHP_Incomplete_Class" {
			original, _ := v.GetProperty("__PHP_Incomplete_Class_Name")
			name = original.AsString()
			filtered := props[:0]
			for _, p := range props {
				if p.name != "__PHP_Incomplete_Class_Name" {
					filtered = append(filtered, p)
				}
			}
			props = filtered
		}
		var b strings.Builder
		b.WriteString("O:" + strconv.Itoa(len(name)) + ":\"" + name + "\":" + strconv.Itoa(len(props)) + ":{")
		for _, p := range props {
			if p.key != nil {
				switch key := p.key.(type) {
				case *data.IntValue:
					b.WriteString("i:" + strconv.Itoa(key.Value) + ";")
				default:
					b.WriteString(makeSerializedString(key.AsString()))
				}
			} else {
				b.WriteString(makeSerializedString(p.name))
			}
			encoded, ok := serializeSlot(p.slot, st)
			if !ok {
				return "", false
			}
			b.WriteString(encoded)
		}
		b.WriteByte('}')
		return b.String(), true
	case *data.FuncValue, *data.BoundFuncValue:
		st.control = data.NewErrorThrowByName(nil, fmt.Errorf("Serialization of 'Closure' is not allowed"), "Exception")
		return "", false
	default:
		return "i:0;", true
	}
}

type serializedProperty struct {
	name string
	slot *data.ZVal
	key  data.Value
}

func callSerializationHook(ctx data.Context, object *data.ClassValue, method data.Method, args []data.Value) (data.GetValue, data.Control) {
	self := data.MethodDeclaringClass(ctx.GetVM(), object.Class, method.GetName())
	frame := data.WrapMethodFrame(ctx.CreateContext(method.GetVariables()), object, self, object.Class)
	if ctl := data.BindDeclaredArgs(frame, method, args); ctl != nil {
		return nil, ctl
	}
	return method.Call(frame)
}
func serializedObjectProperties(object *data.ClassValue, ctx data.Context) ([]serializedProperty, data.Control) {
	var sleepNames *data.ArrayValue
	if ctx != nil {
		if hook, found := object.GetMethod("__serialize"); found {
			raw, ctl := callSerializationHook(ctx, object, hook, nil)
			if ctl != nil {
				return nil, ctl
			}
			array, ok := raw.(*data.ArrayValue)
			if !ok {
				return nil, data.NewTypeError(nil, fmt.Errorf("%s::__serialize() must return an array", object.Class.GetName()))
			}
			props := make([]serializedProperty, 0, array.Len())
			for view, i := array.View(), 0; i < view.Len(); i++ {
				slot := view.At(i)
				props = append(props, serializedProperty{slot.PHPArrayKey(i).AsString(), slot, slot.PHPArrayKey(i)})
			}
			return props, nil
		}
		if hook, found := object.GetMethod("__sleep"); found {
			raw, ctl := callSerializationHook(ctx, object, hook, nil)
			if ctl != nil {
				return nil, ctl
			}
			var ok bool
			sleepNames, ok = raw.(*data.ArrayValue)
			if !ok {
				return nil, data.EmitPHPError(ctx, 2, fmt.Sprintf("serialize(): %s::__sleep() should return an array only containing the names of instance-variables to serialize", object.Class.GetName()), nil)
			}
		}
	}
	classes := []data.ClassStmt{object.Class}
	if ctx != nil {
		for class := object.Class; class.GetExtend() != nil; {
			parent, ctl := ctx.GetVM().GetOrLoadClass(*class.GetExtend())
			if ctl != nil {
				return nil, ctl
			}
			if parent == nil {
				break
			}
			classes = append(classes, parent)
			class = parent
		}
	}
	props := []serializedProperty{}
	seen := map[string]int{}
	for i := len(classes) - 1; i >= 0; i-- {
		for _, property := range classes[i].GetPropertyList() {
			if property == nil || property.GetIsStatic() {
				continue
			}
			stored := data.PropertyStorageName(property)
			if !object.HasProperty(stored) {
				continue
			}
			name := property.GetName()
			switch property.GetModifier() {
			case data.ModifierPrivate:
				name = "\x00" + classes[i].GetName() + "\x00" + name
			case data.ModifierProtected:
				name = "\x00*\x00" + name
			}
			slot, _ := object.ObjectValue.GetZVal(stored)
			p := serializedProperty{name, slot, nil}
			if n, ok := seen[stored]; ok {
				props[n] = p
			} else {
				seen[stored] = len(props)
				props = append(props, p)
			}
		}
	}
	object.RangeProperties(func(name string, value data.Value) bool {
		if _, found := seen[name]; found {
			return true
		}
		slot, _ := object.ObjectValue.GetZVal(name)
		props = append(props, serializedProperty{name, slot, nil})
		return true
	})
	if sleepNames != nil {
		byName := make(map[string]serializedProperty, len(props))
		for _, property := range props {
			byName[property.name] = property
		}
		selected := []serializedProperty{}
		used := map[string]bool{}
		for view, i := sleepNames.View(), 0; i < view.Len(); i++ {
			value := view.At(i).ReadValue()
			if _, ok := value.(*data.StringValue); !ok {
				if ctl := data.EmitPHPError(ctx, 2, fmt.Sprintf("serialize(): %s::__sleep() should return an array only containing the names of instance-variables to serialize", object.Class.GetName()), nil); ctl != nil {
					return nil, ctl
				}
			}
			name := value.AsString()
			property, found := byName[name]
			if !found {
				property, found = byName["\x00"+object.Class.GetName()+"\x00"+name]
			}
			if !found {
				property, found = byName["\x00*\x00"+name]
			}
			if !found {
				if ctl := data.EmitPHPError(ctx, 2, fmt.Sprintf("serialize(): %q returned as member variable from __sleep() but does not exist", name), nil); ctl != nil {
					return nil, ctl
				}
				continue
			}
			if used[property.name] {
				if ctl := data.EmitPHPError(ctx, 2, fmt.Sprintf("serialize(): %q is returned from __sleep() multiple times", name), nil); ctl != nil {
					return nil, ctl
				}
				continue
			}
			used[property.name] = true
			selected = append(selected, property)
		}
		props = selected
	}
	return props, nil
}

type SerializeFunction struct{ data.Function }

func NewSerializeFunction() data.FuncStmt { return &SerializeFunction{} }
func (f *SerializeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	// PHP passes the top-level array by value. A reference back to the input
	// points to the original array, not to this serialization argument snapshot.
	if array, ok := value.(*data.ArrayValue); ok {
		view := array.View()
		slots := make([]*data.ZVal, view.Len())
		for i := range slots {
			slots[i] = view.At(i)
		}
		value = data.NewArrayValueFromSlots(slots)
	}
	state := newPhpSerState()
	state.ctx = ctx
	if raw, ok := core.IniGetInContext(ctx, "serialize_precision"); ok {
		if precision, err := strconv.Atoi(raw); err == nil && precision >= -1 {
			state.precision = precision
		}
	}
	encoded, ok := phpSerializeValue(value, state)
	if state.control != nil {
		return nil, state.control
	}
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(encoded), nil
}
func (f *SerializeFunction) GetName() string { return "serialize" }

var serializeFunctionGetParams = []data.GetValue{node.NewParameter(nil, "value", 0, nil, nil)}

func (f *SerializeFunction) GetParams() []data.GetValue { return serializeFunctionGetParams }

var serializeFunctionGetVariables = []data.Variable{node.NewVariable(nil, "value", 0, nil)}

func (f *SerializeFunction) GetVariables() []data.Variable { return serializeFunctionGetVariables }
