package node

import (
	"github.com/php-any/origami/data"
)

// 注解类型常量
const (
	TypeFeature    = "TypeFeature"    // 特性注解
	TypeMacro      = "TypeMacro"      // 宏注解
	TypeRepeatable = "TypeRepeatable" // 可重复使用的注解（对应 Attribute::IS_REPEATABLE）

	// Attribute 目标类型标记（用于 gen-std 生成 TARGET_* 伪代码）
	TypeTargetClass     = "TypeTargetClass"
	TypeTargetMethod    = "TypeTargetMethod"
	TypeTargetProperty  = "TypeTargetProperty"
	TypeTargetFunction  = "TypeTargetFunction"
	TypeTargetParameter = "TypeTargetParameter"

	TargetName = "target"
)

// AnnotationTargetParameter 注解构造函数中的 target 参数，接收被注解的 AST 节点。
type AnnotationTargetParameter struct {
	*Parameter
}

// NewAnnotationTargetParameter 创建注解 target 参数：默认 null，类型为 AstNode。
func NewAnnotationTargetParameter(from data.From, index int) data.GetValue {
	return &AnnotationTargetParameter{
		Parameter: &Parameter{
			Node:         NewNode(from),
			Name:         TargetName,
			Index:        index,
			Type:         data.TypeAST,
			DefaultValue: data.NewNullValue(),
		},
	}
}

func (p *AnnotationTargetParameter) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return p.Parameter.GetValue(ctx)
}

// NewAnnotationTargetVariable 创建与 target 参数配套的变量声明。
func NewAnnotationTargetVariable(from data.From, index int) data.Variable {
	return NewVariable(from, TargetName, index, data.AST{})
}

func asAnnotationTargetParam(param data.GetValue) data.Variable {
	switch p := param.(type) {
	case *AnnotationTargetParameter:
		return p.Parameter
	case *Parameter:
		if p.Name == TargetName {
			return p
		}
	}
	return nil
}

// Annotation 表示注解节点
type Annotation struct {
	*Node
	Name      string          // 注解名称
	Arguments []data.GetValue // 注解参数
	Target    data.GetValue
	class     data.ClassStmt `pp:"-"` // 首次解析后缓存
}

// NewAnnotation 创建一个新的注解节点
func NewAnnotation(from data.From, name string, arguments []data.GetValue) *Annotation {
	return &Annotation{
		Node:      NewNode(from),
		Name:      name,
		Arguments: arguments,
	}
}

// GetName 返回注解名称
func (a *Annotation) GetName() string {
	return a.Name
}

// GetArguments 返回注解参数
func (a *Annotation) GetArguments() []data.GetValue {
	return a.Arguments
}

func (a *Annotation) resolveClass(ctx data.Context) (data.ClassStmt, data.Control) {
	return ctx.GetVM().GetOrLoadClass(a.Name)
}

// GetValue 获取注解节点的值
func (a *Annotation) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	stmt, acl := a.resolveClass(ctx)
	if acl != nil {
		return nil, acl
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	if object, ok := object.(*data.ClassValue); ok {
		if method := object.Class.GetConstruct(); method != nil {
			fnCtx, acl := bindAnnotationConstructorArgs(a, ctx, object)
			if acl != nil {
				return nil, acl
			}
			return object, &CallAnn{method: method, ctx: fnCtx}
		}
	}

	return object, acl
}

// unwrapAnnotationTarget 解开 data.GetValue 接口，避免 NewAnyValue 二次装箱
func unwrapAnnotationTarget(target data.GetValue) any {
	if target == nil {
		return nil
	}
	switch t := target.(type) {
	case *ClassStatement:
		return t
	case *ClassRegisterStmt:
		return t.Class
	case *AbstractClassStatement:
		return t.ClassStatement
	case *ClassGeneric:
		return t.ClassStatement
	default:
		return target
	}
}

type CallAnn struct {
	method data.Method
	ctx    data.Context
}

func (c *CallAnn) AsString() string {
	return "TODO"
}

func (c *CallAnn) GetValue(fnCtx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

func (c *CallAnn) InitAnnotation() data.Control {
	_, acl := c.method.Call(c.ctx)
	return acl
}
