package annotation

import "github.com/php-any/origami/data"

var constraintSpecs = []ConstraintSpec{
	{
		FullName:   "Validation\\Annotation\\Name",
		Repeatable: false,
		Params: []ConstraintParam{
			{Name: "value", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
	{
		FullName:   "Validation\\Annotation\\Min",
		Repeatable: true,
		Params: []ConstraintParam{
			{Name: "value", Type: data.NewBaseType("int"), DefaultVal: data.NewIntValue(0)},
			{Name: "message", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
	{
		FullName:   "Validation\\Annotation\\Max",
		Repeatable: true,
		Params: []ConstraintParam{
			{Name: "value", Type: data.NewBaseType("int"), DefaultVal: data.NewIntValue(0)},
			{Name: "message", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
	{
		FullName:   "Validation\\Annotation\\Email",
		Repeatable: true,
		Params: []ConstraintParam{
			{Name: "message", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
	{
		FullName:   "Validation\\Annotation\\NotBlank",
		Repeatable: true,
		Params: []ConstraintParam{
			{Name: "message", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
	{
		FullName:   "Validation\\Annotation\\Size",
		Repeatable: true,
		Params: []ConstraintParam{
			{Name: "min", Type: data.NewBaseType("int"), DefaultVal: data.NewIntValue(0)},
			{Name: "max", Type: data.NewBaseType("int"), DefaultVal: data.NewIntValue(0)},
			{Name: "message", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
	{
		FullName:   "Validation\\Annotation\\Pattern",
		Repeatable: true,
		Params: []ConstraintParam{
			{Name: "regexp", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
			{Name: "message", Type: data.NewBaseType("string"), DefaultVal: data.NewStringValue("")},
		},
	},
}
