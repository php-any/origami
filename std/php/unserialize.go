package php

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"math"
	"strconv"
	"strings"
)

type serializedIncompleteClass struct{ data.StdClass }

func (*serializedIncompleteClass) GetName() string { return "__PHP_Incomplete_Class" }

type phpDecoder struct {
	input    string
	pos      int
	ctx      data.Context
	values   []*data.ZVal
	control  data.Control
	hooks    []func() data.Control
	allowed  data.Value
	maxDepth int
}

func (d *phpDecoder) take(s string) bool {
	if !strings.HasPrefix(d.input[d.pos:], s) {
		return false
	}
	d.pos += len(s)
	return true
}
func (d *phpDecoder) number(delimiter byte) (int, bool) {
	start := d.pos
	for d.pos < len(d.input) && d.input[d.pos] >= '0' && d.input[d.pos] <= '9' {
		d.pos++
	}
	if start == d.pos || d.pos >= len(d.input) || d.input[d.pos] != delimiter {
		return 0, false
	}
	n, err := strconv.Atoi(d.input[start:d.pos])
	d.pos++
	return n, err == nil && n >= 0
}
func (d *phpDecoder) quoted(end string) (string, bool) {
	length, ok := d.number(':')
	if !ok || !d.take("\"") || length > len(d.input)-d.pos-len(end) {
		return "", false
	}
	value := d.input[d.pos : d.pos+length]
	d.pos += length
	return value, d.take(end)
}
func (d *phpDecoder) key() (data.Value, bool) {
	if d.take("s:") {
		s, ok := d.quoted("\";")
		return data.NewStringValue(s), ok
	}
	if !d.take("i:") {
		return nil, false
	}
	start := d.pos
	end := strings.IndexByte(d.input[start:], ';')
	if end < 0 {
		return nil, false
	}
	n, err := strconv.Atoi(d.input[start : start+end])
	d.pos = start + end + 1
	return data.NewIntValue(n), err == nil
}
func (d *phpDecoder) value(depth int) (*data.ZVal, int, bool, bool) {
	if d.pos >= len(d.input) {
		return nil, 0, false, false
	}
	kind := d.input[d.pos]
	d.pos++
	if kind == 'R' || kind == 'r' {
		if !d.take(":") {
			return nil, 0, false, false
		}
		n, ok := d.number(';')
		if !ok || n == 0 || n >= len(d.values) {
			return nil, 0, false, false
		}
		source := d.values[n]
		if source == nil {
			return nil, 0, false, false
		}
		if kind == 'R' {
			return source, n, true, true
		}
		if _, ok := source.ReadValue().(*data.ClassValue); !ok {
			return nil, 0, false, false
		}
		slot := data.NewZVal(source.ReadValue())
		d.values = append(d.values, slot)
		return slot, len(d.values) - 1, false, true
	}
	slot := data.NewZVal(data.NewNullValue())
	d.values = append(d.values, slot)
	id := len(d.values) - 1
	switch kind {
	case 'N':
		return slot, id, false, d.take(";")
	case 'b':
		if d.take(":0;") {
			slot.StoreRaw(data.NewBoolValue(false))
			return slot, id, false, true
		}
		if d.take(":1;") {
			slot.StoreRaw(data.NewBoolValue(true))
			return slot, id, false, true
		}
		d.pos-- // PHP reports malformed booleans at their type marker.
	case 'i', 'd':
		if !d.take(":") {
			break
		}
		start := d.pos
		end := strings.IndexByte(d.input[start:], ';')
		if end < 0 {
			break
		}
		raw := d.input[start : start+end]
		d.pos = start + end + 1
		if kind == 'i' {
			n, err := strconv.Atoi(raw)
			if errors.Is(err, strconv.ErrRange) {
				if ctl := data.EmitPHPError(d.ctx, 2, "unserialize(): Numerical result out of range", nil); ctl != nil {
					d.control = ctl
					return nil, 0, false, false
				}
				err = nil // Atoi returns the saturated machine integer on overflow.
			}
			if err != nil {
				d.pos = start - 2
				break
			}
			slot.StoreRaw(data.NewIntValue(n))
		} else {
			var n float64
			var err error
			switch raw {
			case "NAN":
				n = math.NaN()
			case "INF":
				n = math.Inf(1)
			case "-INF":
				n = math.Inf(-1)
			default:
				if !serializedDecimal(raw) {
					d.pos = start - 2
					return nil, 0, false, false
				}
				n, err = strconv.ParseFloat(raw, 64)
			}
			if err != nil && !math.IsInf(n, 0) {
				break
			}
			slot.StoreRaw(data.NewFloatValue(n))
		}
		return slot, id, false, true
	case 's':
		if !d.take(":") {
			break
		}
		s, ok := d.quoted("\";")
		if !ok {
			break
		}
		slot.StoreRaw(data.NewStringValue(s))
		return slot, id, false, true
	case 'E':
		if !d.take(":") {
			break
		}
		name, ok := d.quoted("\";")
		if !ok || d.ctx == nil {
			break
		}
		className, caseName, ok := strings.Cut(name, ":")
		if !ok {
			break
		}
		class, ctl := d.ctx.GetVM().GetOrLoadClass(className)
		if ctl != nil {
			d.control = ctl
			break
		}
		metadata, isEnum := class.(interface{ DeclarationFlags() data.ClassFlags })
		if !isEnum || metadata.DeclarationFlags()&data.ClassEnum == 0 {
			break
		}
		getter, ok := class.(data.GetStaticProperty)
		if !ok {
			break
		}
		value, ok := getter.GetStaticProperty(caseName)
		if !ok {
			break
		}
		object, ok := value.(*data.ClassValue)
		if !ok || !data.SameNominalClass(object.Class, className, d.ctx.GetVM()) {
			break
		}
		slot.StoreRaw(object)
		return slot, id, false, true
	case 'C':
		if !d.take(":") {
			break
		}
		name, ok := d.quoted("\":")
		if !ok {
			break
		}
		length, ok := d.number(':')
		if !ok || !d.take("{") || length > len(d.input)-d.pos-1 {
			break
		}
		payload := d.input[d.pos : d.pos+length]
		d.pos += length
		if !d.take("}") {
			break
		}
		object := d.object(name)
		if object == nil {
			break
		}
		slot.StoreRaw(object)
		if hook, found := object.GetMethod("unserialize"); found && d.ctx != nil && data.NominalIsA(object.Class, "Serializable", d.ctx.GetVM()) {
			_, ctl := callSerializationHook(d.ctx, object, hook, []data.Value{data.NewStringValue(payload)})
			if ctl != nil {
				d.control = ctl
				return nil, 0, false, false
			}
		} else {
			if ctl := data.EmitPHPError(d.ctx, 2, fmt.Sprintf("Class %s has no unserializer", name), nil); ctl != nil {
				d.control = ctl
				return nil, 0, false, false
			}
		}
		return slot, id, false, true
	case 'a', 'O':
		if !d.take(":") {
			break
		}
		var object *data.ClassValue
		var array *data.ArrayValue
		name := ""
		if kind == 'O' {
			var ok bool
			name, ok = d.quoted("\":")
			if !ok || name == "" {
				break
			}
			object = d.object(name)
			if object == nil {
				break
			}
			slot.StoreRaw(object)
		} else {
			array = data.NewArrayValue(nil).(*data.ArrayValue)
			slot.StoreRaw(array)
		}
		count, ok := d.number(':')
		if !ok || count > len(d.input)-d.pos || !d.take("{") {
			break
		}
		if (count > 0 || kind == 'O') && d.maxDepth > 0 && depth >= d.maxDepth {
			d.control = data.EmitPHPError(d.ctx, 2, fmt.Sprintf("unserialize(): Maximum depth of %d exceeded. The depth limit can be changed using the max_depth unserialize() option or the unserialize_max_depth ini setting", d.maxDepth), nil)
			break
		}
		properties := data.NewArrayValue(nil).(*data.ArrayValue)
		custom := false
		if object != nil {
			_, custom = object.GetMethod("__unserialize")
		}
		for i := 0; i < count; i++ {
			key, ok := d.key()
			if !ok {
				return nil, 0, false, false
			}
			child, childID, ref, ok := d.value(depth + 1)
			if !ok {
				return nil, 0, false, false
			}
			if array != nil || custom {
				target := array
				if target == nil {
					target = properties
				}
				if ref {
					target.BindReference(key, child)
				} else {
					target.SetKey(key, child.ReadValue())
					installed, _ := target.LookupZValByStringKey(key.AsString())
					d.values[childID] = installed
				}
			} else {
				property := key.AsString()
				if object.Class.GetName() != "__PHP_Incomplete_Class" && strings.HasPrefix(property, "\x00*\x00") {
					property = strings.TrimPrefix(property, "\x00*\x00")
				}
				_, declared := object.GetPropertyStmt(property)
				created := !declared && !object.ObjectValue.HasProperty(property)
				if created {
					if ctl := node.CheckDynamicPropertyCreation(object, property, nil); ctl != nil {
						d.control = ctl
						return nil, 0, false, false
					}
				}
				if ref {
					strict := d.ctx
					if strict != nil {
						strict = strict.CreateContext(nil)
						strict.SetStrictTypes(true)
					}
					if ctl := object.BindPropertyReference(strict, property, child); ctl != nil {
						d.control = ctl
						return nil, 0, false, false
					}
				} else {
					if d.ctx != nil {
						prepared := child.ReadValue()
						if declaration, found := object.GetPropertyStmt(property); found {
							strict := d.ctx.CreateContext(nil)
							strict.SetStrictTypes(true)
							ref := data.PropertyType(object, property, data.DeclaredTypeRef(declaration.GetType()))
							var accepted bool
							var ctl data.Control
							prepared, accepted, ctl = data.PrepareDeclaredValueInContext(ref, prepared, strict)
							if ctl != nil {
								d.control = ctl
								return nil, 0, false, false
							}
							if !accepted {
								d.control = data.NewTypeError(nil, fmt.Errorf("Cannot assign unserialized value to property %s::$%s", object.Class.GetName(), property))
								return nil, 0, false, false
							}
						}
						object.SetProperty(property, prepared)
					} else {
						object.SetProperty(property, child.ReadValue())
					}
					installed, _ := object.ObjectValue.GetZVal(property)
					d.values[childID] = installed
				}
				if created {
					if ctl := node.ReportDynamicPropertyCreation(d.ctx, object, property, nil); ctl != nil {
						d.control = ctl
						return nil, 0, false, false
					}
				}
			}
		}
		if !d.take("}") {
			break
		}
		if object != nil && d.ctx != nil {
			if hook, found := object.GetMethod("__unserialize"); found {
				d.hooks = append(d.hooks, func() data.Control {
					_, ctl := callSerializationHook(d.ctx, object, hook, []data.Value{properties})
					return ctl
				})
			} else if hook, found := object.GetMethod("__wakeup"); found {
				d.hooks = append(d.hooks, func() data.Control { _, ctl := callSerializationHook(d.ctx, object, hook, nil); return ctl })
			}
		}
		return slot, id, false, true
	}
	return nil, 0, false, false
}
func (d *phpDecoder) object(name string) *data.ClassValue {
	allowed := true
	switch option := d.allowed.(type) {
	case *data.BoolValue:
		allowed = option.Value
	case *data.ArrayValue:
		allowed = false
		for view, i := option.View(), 0; i < view.Len(); i++ {
			if strings.EqualFold(view.At(i).ReadValue().AsString(), name) {
				allowed = true
				break
			}
		}
	}
	if allowed && d.ctx != nil {
		class, found := d.ctx.GetVM().GetClass(name)
		if !found {
			loaded, ctl := d.ctx.GetVM().LoadPkg(name)
			if ctl != nil {
				d.control = ctl
				return nil
			}
			class, _ = loaded.(data.ClassStmt)
		}
		if class != nil {
			raw, ctl := class.GetValue(d.ctx)
			if ctl != nil {
				d.control = ctl
				return nil
			}
			object, _ := raw.(*data.ClassValue)
			return object
		}
	}
	object := data.NewClassValue(&serializedIncompleteClass{}, d.ctx)
	object.SetProperty("__PHP_Incomplete_Class_Name", data.NewStringValue(name))
	return object
}
func newPHPDecoder(input string, ctx data.Context) *phpDecoder {
	depth := 4096
	if raw, ok := core.IniGetInContext(ctx, "unserialize_max_depth"); ok {
		if configured, err := strconv.Atoi(raw); err == nil && configured >= 0 {
			depth = configured
		}
	}
	return &phpDecoder{input: input, ctx: ctx, values: []*data.ZVal{nil}, maxDepth: depth}
}

// PHP accepts decimal doubles only; strconv.ParseFloat also accepts Go's hex
// floats and case-insensitive infinities. Those must not enter the wire format.
func serializedDecimal(raw string) bool {
	i := 0
	if i < len(raw) && (raw[i] == '+' || raw[i] == '-') {
		i++
	}
	digits := 0
	for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
		i++
		digits++
	}
	if i < len(raw) && raw[i] == '.' {
		i++
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
			digits++
		}
	}
	if digits == 0 {
		return false
	}
	if i < len(raw) && (raw[i] == 'e' || raw[i] == 'E') {
		i++
		if i < len(raw) && (raw[i] == '+' || raw[i] == '-') {
			i++
		}
		start := i
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
		if i == start {
			return false
		}
	}
	return i == len(raw)
}
func parsePhpSerializedValue(input string) (data.Value, bool) {
	decoder := newPHPDecoder(input, nil)
	value, _, _, ok := decoder.value(0)
	if !ok || decoder.pos != len(input) {
		return nil, false
	}
	return value.ReadValue(), true
}

type UnserializeFunction struct{}

func NewUnserializeFunction() data.FuncStmt { return &UnserializeFunction{} }
func (f *UnserializeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	input, _ := ctx.GetIndexValue(0)
	decoder := newPHPDecoder(input.AsString(), ctx)
	if options, _ := ctx.GetIndexValue(1); options != nil {
		if options, ok := options.(*data.ArrayValue); ok {
			if slot, ok := options.LookupZValByStringKey("allowed_classes"); ok {
				decoder.allowed = slot.ReadValue()
			}
			if slot, ok := options.LookupZValByStringKey("max_depth"); ok {
				if limit, ok := slot.ReadValue().(*data.IntValue); ok && limit.Value >= 0 {
					decoder.maxDepth = limit.Value
				} else {
					return nil, data.NewErrorThrowByName(nil, fmt.Errorf("unserialize(): max_depth must be a non-negative integer"), "ValueError")
				}
			}
		}
	}
	switch decoder.allowed.(type) {
	case nil, *data.BoolValue, *data.ArrayValue:
	default:
		return nil, data.NewTypeError(nil, fmt.Errorf("unserialize(): allowed_classes must be of type array|bool"))
	}
	value, _, _, ok := decoder.value(0)
	if decoder.control != nil {
		return nil, decoder.control
	}
	if !ok {
		return data.NewBoolValue(false), data.EmitPHPError(ctx, 2, fmt.Sprintf("unserialize(): Error at offset %d of %d bytes", decoder.pos, len(decoder.input)), nil)
	}
	if decoder.pos < len(decoder.input) {
		if ctl := data.EmitPHPError(ctx, 2, "unserialize(): Extra data starting at offset "+strconv.Itoa(decoder.pos)+" of "+strconv.Itoa(len(decoder.input))+" bytes", nil); ctl != nil {
			return nil, ctl
		}
	}
	for _, hook := range decoder.hooks {
		if ctl := hook(); ctl != nil {
			return nil, ctl
		}
	}
	result := value.ReadValue()
	// The decoder's root is a temporary zval. When only one array bucket
	// references it, PHP unwraps that sole reference while releasing the root;
	// the recursive bucket sees null. Multiple aliases keep the shared cell.
	if _, array := result.(*data.ArrayValue); array && value.ReferenceIdentity() != nil && value.RefCount() == 1 {
		value.StoreRaw(data.NewNullValue())
	}
	return result, nil
}
func (f *UnserializeFunction) GetName() string { return "unserialize" }

var unserializeFunctionGetParams = []data.GetValue{node.NewParameter(nil, "data", 0, nil, data.String{}), node.NewParameter(nil, "options", 1, data.NewArrayValue(nil), data.NewBaseType("array"))}

func (f *UnserializeFunction) GetParams() []data.GetValue { return unserializeFunctionGetParams }

var unserializeFunctionGetVariables = []data.Variable{node.NewVariable(nil, "data", 0, data.String{}), node.NewVariable(nil, "options", 1, data.NewBaseType("array"))}

func (f *UnserializeFunction) GetVariables() []data.Variable { return unserializeFunctionGetVariables }
