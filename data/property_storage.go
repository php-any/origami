package data

import "strings"

func PropertyStorageName(property Property) string {
	if stored, ok := property.(interface{ PropertyStorageName() string }); ok {
		return stored.PropertyStorageName()
	}
	return property.GetName()
}

func privatePropertyParts(name string) (string, string, bool) {
	if len(name) == 0 || name[0] != 0 {
		return "", name, false
	}
	owner, property, ok := strings.Cut(name[1:], "\x00")
	return owner, property, ok
}
