package data

import (
	"net/url"
	"strings"
)

// PHP form/query fields preserve input order and use last-value-wins for
// repeated scalar fields. Brackets create nested PHP arrays, including [].
func ParseFormFields(query string) *ArrayValue {
	array := NewArrayValueFromSlots(nil)
	for _, field := range strings.Split(query, "&") {
		name, value, _ := strings.Cut(field, "=")
		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			continue
		}
		decodedValue, err := url.QueryUnescape(value)
		if err != nil {
			continue
		}
		SetFormField(array, decodedName, NewStringValue(decodedValue))
	}
	return array
}

func FormFieldPath(name string) []string {
	name = strings.TrimLeft(name, " ")
	if nul := strings.IndexByte(name, 0); nul >= 0 {
		name = name[:nul]
	}
	base, brackets, found := strings.Cut(name, "[")
	base = strings.NewReplacer(".", "_", " ", "_").Replace(base)
	if base == "" {
		return nil
	}
	path := []string{base}
	if !found {
		return path
	}
	for depth := 0; depth < 64; depth++ {
		key, rest, closed := strings.Cut(brackets, "]")
		if !closed {
			if len(path) == 1 {
				path[0] += "_" + strings.NewReplacer(".", "_", " ", "_").Replace(brackets)
			}
			return path
		}
		path = append(path, key)
		if !strings.HasPrefix(rest, "[") {
			return path
		}
		brackets = rest[1:]
	}
	return nil
}

func SetFormField(array *ArrayValue, name string, value Value) {
	SetFormFieldPath(array, FormFieldPath(name), value)
}
func SetFormFieldPath(array *ArrayValue, path []string, value Value) {
	if len(path) == 0 {
		return
	}
	for position, key := range path {
		if position == len(path)-1 {
			if key == "" && position != 0 {
				array.AppendValue(value)
			} else {
				array.SetStringKey(key, value)
			}
			return
		}
		var slot *ZVal
		if key == "" && position != 0 {
			child := NewArrayValueFromSlots(nil)
			slot = array.AppendSlot(child)
		} else {
			slot, _ = array.LookupZValByStringKey(key)
			if slot == nil {
				array.SetStringKey(key, NewArrayValueFromSlots(nil))
				slot, _ = array.LookupZValByStringKey(key)
			}
		}
		if slot == nil {
			return
		}
		child, ok := slot.ReadValue().(*ArrayValue)
		if !ok {
			child = NewArrayValueFromSlots(nil)
			CowAssign(slot, child)
		}
		array = child
	}
}
