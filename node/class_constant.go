package node

import (
	"fmt"
	"os"
	"strings"

	"github.com/php-any/origami/data"
)

// ClassConstant 表示 ::class 语法节点（含 $obj::class / $name::class）
type ClassConstant struct {
	From *TokenFrom
	Expr data.GetValue
}

// NewClassConstant 创建新的 ClassConstant 节点
func NewClassConstant(from *TokenFrom, expr data.GetValue) *ClassConstant {
	return &ClassConstant{
		From: from,
		Expr: expr,
	}
}

// GetValue 获取类常量值。PHP：$obj::class 返回实例的类名字符串，不触发类加载错误。
func (cc *ClassConstant) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	exprValue, control := cc.Expr.GetValue(ctx)
	if control != nil {
		return nil, control
	}

	if name, ok := classNameFromValue(exprValue); ok {
		return resolveClassNameString(ctx, name), nil
	}

	// 临时诊断：::class 无法解析出类名
	if cc.From != nil && strings.Contains(cc.From.GetSource(), "var-dumper") {
		start, end := cc.From.GetPosition()
		fmt.Fprintf(os.Stderr, "DIAG: ::class unresolvable expr=%T val=%v file=%q pos=%v-%v\n", exprValue, exprValue, cc.From.GetSource(), start, end)
	}

	if varExpr, ok := exprValue.(*VariableExpression); ok {
		varValue, ctl := varExpr.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if name, ok := classNameFromValue(varValue); ok {
			return resolveClassNameString(ctx, name), nil
		}
		return nil, data.NewErrorThrow(cc.From, fmt.Errorf("无法获取变量类型"))
	}

	return data.NewStringValue(""), nil
}

func classNameFromValue(v data.GetValue) (string, bool) {
	switch t := v.(type) {
	case *data.ClassValue:
		if t.Class != nil {
			return t.Class.GetName(), true
		}
	case *data.ThisValue:
		if t.Class != nil {
			return t.Class.GetName(), true
		}
	case *data.ThrowValue:
		// 异常对象也是对象：$e::class 应返回异常类名
		return t.GetName(), true
	case *StringLiteral:
		return t.Value, true
	case *data.StringValue:
		return t.AsString(), true
	}
	return "", false
}

func resolveClassNameString(ctx data.Context, className string) data.GetValue {
	className = strings.TrimPrefix(className, "\\")
	if className == "" {
		return data.NewStringValue("")
	}
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewStringValue(className)
	}
	if class, ok := vm.GetClass(className); ok && class != nil {
		return data.NewStringValue(class.GetName())
	}
	if class, acl := vm.GetOrLoadClass(className); acl == nil && class != nil {
		return data.NewStringValue(class.GetName())
	}
	// PHP：类不存在时 ::class 仍返回给定名称字符串
	return data.NewStringValue(className)
}
