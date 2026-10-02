package node

import (
	"errors"
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
		// 解析期已通过 findFullClassNameByNamespace / self/parent 展开，
		// 不可再按当前命名空间重写，否则 \GlobalClass 会被误加成 Ns\GlobalClass。
		return strings.TrimPrefix(right.Value, "\\"), nil
	case *StaticClass:
		val, acl := right.GetValue(ctx)
		if acl != nil {
			return "", acl
		}
		if sv, ok := val.(*data.StringValue); ok {
			name = sv.AsString()
		}
	default:
		r, acl := classExpr.GetValue(ctx)
		if acl != nil {
			return "", acl
		}
		switch v := r.(type) {
		case *data.StringValue:
			name = v.AsString()
		case *data.ClassValue:
			name = v.Class.GetName()
		case *data.ThisValue:
			if v.Class != nil {
				name = v.Class.GetName()
			}
		case *data.FuncValue, *data.BoundFuncValue:
			name = "Closure"
		case data.Generator:
			name = "Generator"
		case *data.ThrowValue:
			name = v.GetName()
		default:
			return "", data.NewErrorThrowByName(nil, errors.New("Class name must be a valid object or a string"), "Error")
		}
	}
	if name == "" {
		return "", nil
	}
	return strings.TrimPrefix(name, "\\"), nil
}

func instanceof(ctx data.Context, class string, objectValue data.GetValue) (data.GetValue, data.Control) {
	value, ok := objectValue.(data.Value)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	// instanceof RHS is a class name; "object" and "iterable" are not classes.
	if data.TypeNameEqual(class, "object") || data.TypeNameEqual(class, "iterable") {
		return data.NewBoolValue(false), nil
	}
	switch object := objectValue.(type) {
	case *data.ClassValue:
		return data.NewBoolValue(data.NominalIsA(object.Class, class, ctx.GetVM())), nil
	case *data.ThisValue:
		return data.NewBoolValue(data.NominalIsA(object.Class, class, ctx.GetVM())), nil
	}
	return data.NewBoolValue(data.Class{Name: class}.Is(value)), nil
}
