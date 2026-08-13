package annotation

import "github.com/php-any/origami/data"

// CompiledConstraintValue 构建预编译的校验约束注解实例。
func CompiledConstraintValue(fullName string, state map[string]data.GetValue) *data.ClassValue {
	var spec ConstraintSpec
	for _, s := range constraintSpecs {
		if s.FullName == fullName {
			spec = s
			break
		}
	}
	c := &ConstraintClass{
		spec:  spec,
		state: state,
	}
	c.construct = &constraintConstructMethod{class: c}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: c}
}
