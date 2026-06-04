package node

import (
	"strings"

	"github.com/php-any/origami/data"
)

// InstanceOfExpression 表示 instanceof 表达式
type InstanceOfExpression struct {
	*Node
	Object    data.GetValue // 对象表达式
	ClassName data.GetValue // 类名
}

// NewInstanceOfExpression 创建一个新的 instanceof 表达式
func NewInstanceOfExpression(from data.From, object data.GetValue, className data.GetValue) *InstanceOfExpression {
	return &InstanceOfExpression{
		Node:      NewNode(from),
		Object:    object,
		ClassName: className,
	}
}

// GetValue 获取 instanceof 表达式的值
func (i *InstanceOfExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	objectValue, c := i.Object.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	className, acl := resolveInstanceofClassName(ctx, i.ClassName)
	if acl != nil {
		return nil, acl
	}
	if className == "" {
		return data.NewBoolValue(false), nil
	}
	return instanceof(ctx, className, objectValue)
}

// resolveInstanceofClassName 从 instanceof 右侧表达式解析待比较的类名
func resolveInstanceofClassName(ctx data.Context, classExpr data.GetValue) (string, data.Control) {
	var name string
	switch right := classExpr.(type) {
	case *StringLiteral:
		name = right.Value
	default:
		r, acl := classExpr.GetValue(ctx)
		if acl != nil {
			return "", acl
		}
		switch v := r.(type) {
		case data.AsString:
			name = v.AsString()
		case *data.ClassValue:
			name = v.Class.GetName()
		case data.GetName:
			name = v.GetName()
		case *data.ThisValue:
			if v.Class != nil {
				name = v.Class.GetName()
			}
		}
	}
	if name == "" {
		return "", nil
	}
	return resolveRuntimeClassName(ctx, name), nil
}

// resolveRuntimeClassName 按 PHP 命名空间规则解析类名，优先返回已加载的类/接口
func resolveRuntimeClassName(ctx data.Context, name string) string {
	if name == "" {
		return name
	}
	if strings.HasPrefix(name, "\\") {
		return strings.TrimPrefix(name, "\\")
	}

	vm := ctx.GetVM()
	ns := ctx.GetNamespace()

	if vm != nil {
		if _, ok := vm.GetClass(name); ok {
			return name
		}
		if _, ok := vm.GetInterface(name); ok {
			return name
		}
	}

	for _, candidate := range runtimeClassNameCandidates(ns, name) {
		if vm != nil {
			if _, ok := vm.GetClass(candidate); ok {
				return candidate
			}
			if _, ok := vm.GetInterface(candidate); ok {
				return candidate
			}
		}
	}
	candidates := runtimeClassNameCandidates(ns, name)
	if len(candidates) > 0 {
		return candidates[0]
	}
	return name
}

func runtimeClassNameCandidates(ns, name string) []string {
	if strings.Contains(name, "\\") {
		if ns != "" && !strings.HasPrefix(name, ns+"\\") {
			return []string{ns + "\\" + name, name}
		}
		return []string{name}
	}
	if ns != "" {
		return []string{ns + "\\" + name, name}
	}
	return []string{name}
}

func loadClassOrInterfaceForInstanceof(ctx data.Context, class string) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return nil, nil
	}
	if c, ok := vm.GetClass(class); ok {
		return c, nil
	}
	if inf, ok := vm.GetInterface(class); ok {
		return inf, nil
	}
	return vm.LoadPkg(class)
}

func instanceof(ctx data.Context, class string, objectValue data.GetValue) (data.GetValue, data.Control) {
	// 检查对象值是否为类实例
	if classValue, ok := objectValue.(*data.ClassValue); ok {
		c, acl := loadClassOrInterfaceForInstanceof(ctx, class)
		if acl != nil {
			return nil, acl
		}
		if c != nil {
			switch checkC := c.(type) {
			case data.ClassStmt:
				result, acl := checkClassIs(ctx, classValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			case data.InterfaceStmt:
				result, acl := checkClassIs(ctx, classValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			}
		}
	} else if thisValue, ok := objectValue.(*data.ThisValue); ok {
		// 处理 ThisValue（$this）
		c, acl := loadClassOrInterfaceForInstanceof(ctx, class)
		if acl != nil {
			return nil, acl
		}
		if c != nil {
			switch checkC := c.(type) {
			case data.ClassStmt:
				result, acl := checkClassIs(ctx, thisValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			case data.InterfaceStmt:
				result, acl := checkClassIs(ctx, thisValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			}
		}
	}

	switch class {
	case "object":
		switch objectValue.(type) {
		case *data.ClassValue:
			return data.NewBoolValue(true), nil
		case *data.ObjectValue:
			return data.NewBoolValue(true), nil
		}
	case "Closure", "closure":
		switch objectValue.(type) {
		case *data.FuncValue:
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}
