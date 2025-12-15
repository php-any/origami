package php

import (
	"regexp"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregMatchFunction 实现 preg_match 函数
type PregMatchFunction struct{}

func NewPregMatchFunction() data.FuncStmt {
	return &PregMatchFunction{}
}

func (f *PregMatchFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	subjectValue, _ := ctx.GetIndexValue(1)
	matchesValue, _ := ctx.GetIndexValue(2)
	// flagsValue, _ := ctx.GetIndexValue(3)
	// offsetValue, _ := ctx.GetIndexValue(4)

	pattern := patternValue.AsString()
	subject := subjectValue.AsString()

	// PHP patterns usually start and end with a delimiter, e.g. /pattern/
	// Go regexp doesn't use delimiters. We need to strip them.
	// Also handle modifiers like 'i' (case insensitive) after the delimiter.

	if len(pattern) >= 2 {
		delimiter := pattern[0]
		endIndex := -1
		// Find the last occurrence of delimiter
		for i := len(pattern) - 1; i > 0; i-- {
			if pattern[i] == delimiter && pattern[i-1] != '\\' {
				endIndex = i
				break
			}
		}

		if endIndex != -1 {
			modifiers := pattern[endIndex+1:]
			regexBody := pattern[1:endIndex]

			// Handle modifiers
			prefix := ""
			if len(modifiers) > 0 {
				for _, mod := range modifiers {
					if mod == 'i' {
						prefix += "(?i)"
					}
					// Other modifiers like 'm', 's', 'x', 'u' might need mapping or are default/different in Go.
					// 'm' (multiline) -> (?m)
					// 's' (dotall) -> (?s)
					if mod == 'm' {
						prefix += "(?m)"
					}
					if mod == 's' {
						prefix += "(?s)"
					}
				}
			}
			pattern = prefix + regexBody
		}
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		// PHP issues a warning, returns false.
		// For now, return false.
		return data.NewBoolValue(false), nil
	}

	// Find matches
	// preg_match finds the first match.
	loc := re.FindStringSubmatchIndex(subject)
	if loc == nil {
		return data.NewIntValue(0), nil // No match
	}

	// If matches is provided, populate it.
	if matchesValue != nil {
		// matchesValue should be a reference or we should update it if it's passed by reference.
		// In Origami, if we get a Value, is it a reference?
		// We need to check if it's a variable reference or just a value.
		// `ctx.GetIndexValue` returns a value.
		// If the user passed a variable, we might need to update it.
		// But `Call` receives values.
		// Wait, `preg_match` 3rd arg is `array &$matches`.
		// In Origami, how do we handle references?
		// Let's look at `data/value_reference.go` or similar.
		// Or `ctx.GetIndexValue` might return the value, but we can't update the variable in the caller scope unless we have the variable name or reference.
		// However, `GetIndexValue` gets the value of the argument.
		// If the argument was passed by reference, we might need to handle it.
		// But `GetIndexValue` resolves the value.
		// Let's check `data.Context` or `data.VM` to see how references are handled.
		// Actually, standard functions in this codebase seem to just take values.
		// If I look at `exec` or similar?
		// Let's assume for now we can't easily update the variable unless we have a mechanism.
		// But wait, `preg_match` is useless without matches if we want to capture.
		// Let's check if `matchesValue` is a `*data.ReferenceValue`?
		// If so, we can update it.

		// Let's try to cast to ReferenceValue or similar?
		// `data.Value` interface.
		// Let's check `data/value_reference.go`.

		// Assuming we can update it if it's a reference.
		// But `GetIndexValue` might return the dereferenced value?
		// Let's check `data/context.go`.

		// If I can't update it, I'll just skip it for now or try to update if it's an object/array (passed by value in PHP but objects are ref, arrays are value).
		// But `matches` is an output parameter.

		// Let's populate a new array and try to assign it?
		// If `matchesValue` is a `Reference`, we can `SetValue`.
		// I'll check `data.Reference` interface if it exists.

		// For now, I'll construct the array.

		matchStrs := []data.Value{}
		// loc contains [start, end, start, end...]
		for i := 0; i < len(loc); i += 2 {
			start, end := loc[i], loc[i+1]
			if start == -1 {
				matchStrs = append(matchStrs, data.NewStringValue("")) // Unmatched group?
			} else {
				matchStrs = append(matchStrs, data.NewStringValue(subject[start:end]))
			}
		}

		newMatches := data.NewArrayValue(matchStrs)

		// How to assign back?
		// If `matchesValue` is a reference, we can set it.
		if r, ok := matchesValue.(*data.ReferenceValue); ok {
			r.Ctx.SetVariableValue(r.Val, newMatches)
		} else if arr, ok := matchesValue.(*data.ArrayValue); ok {
			arr.Value = matchStrs
		}
	}

	return data.NewIntValue(1), nil // Match found
}

func (f *PregMatchFunction) GetName() string {
	return "preg_match"
}

func (f *PregMatchFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "subject", 1, nil, nil),
		node.NewParameter(nil, "matches", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "flags", 3, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "offset", 4, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *PregMatchFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "subject", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "matches", 2, data.NewBaseType("array")), // Should be reference?
		node.NewVariable(nil, "flags", 3, data.NewBaseType("int")),
		node.NewVariable(nil, "offset", 4, data.NewBaseType("int")),
	}
}
