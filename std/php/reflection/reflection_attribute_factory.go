package reflection

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	reflectionAttributeTargetClass     = 1
	reflectionAttributeTargetFunction  = 2
	reflectionAttributeTargetMethod    = 4
	reflectionAttributeTargetProperty  = 8
	reflectionAttributeTargetParameter = 16
)

// newReflectionAttribute 创建 ReflectionAttribute 实例；decl 为源码注解节点（可为 nil）。
func newReflectionAttribute(ctx data.Context, annotation *data.ClassValue, decl *node.Annotation, target int, repeated bool) *data.ClassValue {
	if decl != nil && annotation != nil {
		// 保持与运行时注解实例一致，供无声明时的回退路径使用
		_ = annotation
	}
	attrClass := &ReflectionAttributeClass{
		declaration: decl,
		target:      target,
		repeated:    repeated,
	}
	attrValue := data.NewClassValue(attrClass, ctx.CreateBaseContext())
	if annotation != nil {
		attrValue.ObjectValue.SetProperty("_annotation", annotation)
	}
	return attrValue
}

func attributeIsRepeatedAt(annotations []*data.ClassValue, index int) bool {
	if index < 0 || index >= len(annotations) || annotations[index] == nil || annotations[index].Class == nil {
		return false
	}
	name := strings.ToLower(annotations[index].Class.GetName())
	seen := 0
	for i, ann := range annotations {
		if ann == nil || ann.Class == nil {
			continue
		}
		if strings.ToLower(ann.Class.GetName()) != name {
			continue
		}
		if i == index {
			return seen > 0
		}
		seen++
	}
	return false
}
