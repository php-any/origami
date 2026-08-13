package data

import "fmt"

// Generic 泛型
type Generic struct {
	Name  string
	Types []Types
}

func (i Generic) Is(value Value) bool {
	return true // TODO
}

func (i Generic) String() string {
	return fmt.Sprintf("%v", i.Name)
}
