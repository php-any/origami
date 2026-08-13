package validation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ValidatorClass PHP 校验器 Validation\Validator。
type ValidatorClass struct {
	node.Node
	validate data.Method
}

func (v *ValidatorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&ValidatorClass{validate: &validatorValidateMethod{}}, ctx.CreateBaseContext()), nil
}

func (v *ValidatorClass) GetName() string    { return "Validation\\Validator" }
func (v *ValidatorClass) GetExtend() *string { return nil }
func (v *ValidatorClass) GetImplements() []string {
	return nil
}
func (v *ValidatorClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (v *ValidatorClass) GetPropertyList() []data.Property           { return nil }
func (v *ValidatorClass) GetConstruct() data.Method                  { return nil }

func (v *ValidatorClass) GetMethod(name string) (data.Method, bool) {
	if name == "validate" {
		if v.validate == nil {
			v.validate = &validatorValidateMethod{}
		}
		return v.validate, true
	}
	return nil, false
}

func (v *ValidatorClass) GetMethods() []data.Method {
	if v.validate == nil {
		v.validate = &validatorValidateMethod{}
	}
	return []data.Method{v.validate}
}

type validatorValidateMethod struct{}

func (m *validatorValidateMethod) GetName() string            { return "validate" }
func (m *validatorValidateMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *validatorValidateMethod) GetIsStatic() bool          { return true }
func (m *validatorValidateMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}
func (m *validatorValidateMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, nil),
	}
}
func (m *validatorValidateMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (m *validatorValidateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj, ok := ctx.GetIndexValue(0)
	if !ok || obj == nil {
		return data.NewArrayValue(nil), nil
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return data.NewArrayValue(nil), nil
	}
	violations := ValidateObject(cv)
	items := make([]data.Value, 0, len(violations))
	for _, v := range violations {
		item := data.NewObjectValue()
		item.SetProperty("field", data.NewStringValue(v.Field))
		item.SetProperty("message", data.NewStringValue(v.Message))
		items = append(items, item)
	}
	return data.NewArrayValue(items), nil
}
