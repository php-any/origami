package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	valannotation "github.com/php-any/origami/std/validation/annotation"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ValidateObject 校验 DTO 实例上属性声明的约束注解，返回违规列表（空表示通过）。
func ValidateObject(dto *data.ClassValue) []Violation {
	if dto == nil || dto.Class == nil {
		return nil
	}
	var violations []Violation
	for _, prop := range dto.Class.GetPropertyList() {
		cp, ok := prop.(*node.ClassProperty)
		if !ok || len(cp.Annotations) == 0 {
			continue
		}
		fieldLabel := fieldDisplayName(cp)
		val, acl := cp.GetValue(dto)
		if acl != nil {
			violations = append(violations, Violation{
				Field:   cp.Name,
				Message: acl.AsString(),
			})
			continue
		}
		var valData data.Value
		if val != nil {
			valData, _ = val.(data.Value)
		}
		for _, ann := range cp.Annotations {
			cc, ok := ann.Class.(*valannotation.ConstraintClass)
			if !ok {
				continue
			}
			if msg, fail := checkConstraint(cc, fieldLabel, cp.Name, valData); fail {
				violations = append(violations, Violation{Field: cp.Name, Message: msg})
			}
		}
	}
	return violations
}

// ConstraintAnnotations 从注解列表中筛出 Validation 约束（不含 Name）。
func ConstraintAnnotations(anns []*data.ClassValue) []*data.ClassValue {
	if len(anns) == 0 {
		return nil
	}
	out := make([]*data.ClassValue, 0, len(anns))
	for _, ann := range anns {
		if ann == nil || ann.Class == nil {
			continue
		}
		cc, ok := ann.Class.(*valannotation.ConstraintClass)
		if !ok || cc.Spec().FullName == "Validation\\Annotation\\Name" {
			continue
		}
		out = append(out, ann)
	}
	return out
}

// AnnotationDisplayName 从注解列表读取 #[Name] 显示名，否则返回 fallback。
func AnnotationDisplayName(anns []*data.ClassValue, fallback string) string {
	for _, ann := range anns {
		cc, ok := ann.Class.(*valannotation.ConstraintClass)
		if !ok || cc.Spec().FullName != "Validation\\Annotation\\Name" {
			continue
		}
		if v := stateString(cc.State(), "value", ""); v != "" {
			return v
		}
	}
	return fallback
}

// ValidateConstraints 对单个值应用约束注解列表。
func ValidateConstraints(label, field string, constraints []*data.ClassValue, value data.Value) []Violation {
	if len(constraints) == 0 {
		return nil
	}
	var violations []Violation
	for _, ann := range constraints {
		cc, ok := ann.Class.(*valannotation.ConstraintClass)
		if !ok {
			continue
		}
		if msg, fail := checkConstraint(cc, label, field, value); fail {
			violations = append(violations, Violation{Field: field, Message: msg})
		}
	}
	return violations
}

func fieldDisplayName(cp *node.ClassProperty) string {
	for _, ann := range cp.Annotations {
		cc, ok := ann.Class.(*valannotation.ConstraintClass)
		if !ok || cc.Spec().FullName != "Validation\\Annotation\\Name" {
			continue
		}
		if v := stateString(cc.State(), "value", ""); v != "" {
			return v
		}
	}
	return cp.Name
}

func checkConstraint(cc *valannotation.ConstraintClass, label, field string, value data.Value) (string, bool) {
	state := cc.State()
	switch cc.Spec().FullName {
	case "Validation\\Annotation\\Name":
		return "", false
	case "Validation\\Annotation\\NotBlank":
		if isBlank(value) {
			return defaultMessage(state, "message", fmt.Sprintf("%s 不能为空", label)), true
		}
	case "Validation\\Annotation\\Email":
		if isBlank(value) {
			return "", false
		}
		if !emailPattern.MatchString(valueAsString(value)) {
			return defaultMessage(state, "message", fmt.Sprintf("%s 必须是合法邮箱", label)), true
		}
	case "Validation\\Annotation\\Min":
		if isBlank(value) {
			return "", false
		}
		min := stateInt(state, "value", 0)
		if n, ok := valueAsInt(value); !ok || n < min {
			return defaultMessage(state, "message", fmt.Sprintf("%s 不能小于 %d", label, min)), true
		}
	case "Validation\\Annotation\\Max":
		if isBlank(value) {
			return "", false
		}
		max := stateInt(state, "value", 0)
		if n, ok := valueAsInt(value); !ok || n > max {
			return defaultMessage(state, "message", fmt.Sprintf("%s 不能大于 %d", label, max)), true
		}
	case "Validation\\Annotation\\Size":
		if isBlank(value) {
			return "", false
		}
		min := stateInt(state, "min", 0)
		max := stateInt(state, "max", 0)
		size := valueSize(value)
		if min > 0 && size < min {
			return defaultMessage(state, "message", fmt.Sprintf("%s 长度不能小于 %d", label, min)), true
		}
		if max > 0 && size > max {
			return defaultMessage(state, "message", fmt.Sprintf("%s 长度不能大于 %d", label, max)), true
		}
	case "Validation\\Annotation\\Pattern":
		if isBlank(value) {
			return "", false
		}
		expr := stateString(state, "regexp", "")
		re, err := compilePattern(expr)
		if err != nil {
			return fmt.Sprintf("%s 的正则表达式无效", label), true
		}
		if !re.MatchString(valueAsString(value)) {
			return defaultMessage(state, "message", fmt.Sprintf("%s 格式不正确", label)), true
		}
	}
	return "", false
}

func defaultMessage(state map[string]data.GetValue, key, fallback string) string {
	if msg := stateString(state, key, ""); msg != "" {
		return msg
	}
	return fallback
}

func stateString(state map[string]data.GetValue, key, def string) string {
	if v, ok := state[key]; ok && v != nil {
		if sv, ok := v.(data.AsString); ok {
			return sv.AsString()
		}
	}
	return def
}

func stateInt(state map[string]data.GetValue, key string, def int) int {
	if v, ok := state[key]; ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				return n
			}
		}
	}
	return def
}

func isBlank(value data.Value) bool {
	if value == nil {
		return true
	}
	if _, ok := value.(*data.NullValue); ok {
		return true
	}
	if sv, ok := value.(*data.StringValue); ok {
		return sv.AsString() == ""
	}
	if av, ok := value.(*data.ArrayValue); ok {
		return len(av.List) == 0
	}
	if pv, ok := value.(*data.ProxyValue); ok && pv.Class != nil {
		if pv.Class.GetName() == "Net\\Http\\UploadedFile" {
			if getter, ok := pv.Class.(interface{ GetSource() any }); ok {
				return getter.GetSource() == nil
			}
		}
	}
	return false
}

func valueAsString(value data.Value) string {
	if value == nil {
		return ""
	}
	if sv, ok := value.(data.AsString); ok {
		return sv.AsString()
	}
	return ""
}

func valueAsInt(value data.Value) (int, bool) {
	if value == nil {
		return 0, false
	}
	if iv, ok := value.(data.AsInt); ok {
		n, err := iv.AsInt()
		return n, err == nil
	}
	if fv, ok := value.(data.AsFloat); ok {
		f, err := fv.AsFloat()
		return int(f), err == nil
	}
	return 0, false
}

func valueSize(value data.Value) int {
	if value == nil {
		return 0
	}
	if av, ok := value.(*data.ArrayValue); ok {
		return len(av.List)
	}
	return utf8.RuneCountInString(valueAsString(value))
}

func compilePattern(expr string) (*regexp.Regexp, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty pattern")
	}
	if expr[0] == '/' {
		if end := strings.LastIndex(expr, "/"); end > 0 {
			return regexp.Compile(expr[1:end])
		}
	}
	return regexp.Compile(expr)
}
