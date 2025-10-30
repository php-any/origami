package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// 注解类型常量
const (
	TypeFeature = "TypeFeature" // 特性注解
	TypeMacro   = "TypeMacro"   // 宏注解

	TargetName = "target"
)

// Annotation 表示注解节点
type Annotation struct {
	*Node
	Name      string          // 注解名称
	Arguments []data.GetValue // 注解参数
	Target    data.GetValue
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

// GetValue 获取注解节点的值
func (a *Annotation) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	stmt, ok := vm.GetClass(a.Name)
	if !ok {
		return nil, data.NewErrorThrow(a.from, errors.New(fmt.Sprintf("类 %s 不存在", a.Name)))
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	if object, ok := object.(*data.ClassValue); ok {
		if method := object.Class.GetConstruct(); method != nil {
			varies := method.GetVariables()
			fnCtx := object.CreateContext(varies)
			// 入参的值设置到上下文中
			for index, arg := range a.Arguments {
				switch argTV := arg.(type) {
				case *NamedArgument:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, argTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(a.from, err)
					}
					fnCtx.SetVariableValue(vari, tempV.(data.Value))
				default:
					tempV, acl := argTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}

					fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
				}
			}

			// 将被注解的 AST 目标按需注入构造函数：
			// 只要构造函数声明了名为 target 的参数，就注入，不再强依赖是否实现 TypeMacro
			if vari, err := findVariable(varies, TargetName); err == nil {
				fnCtx.SetVariableValue(vari, data.NewAnyValue(a.Target))
			}

			// 构造函数执行成功后，返回注解实例本身
			return object, &CallAnn{method: method, ctx: fnCtx}
		}
	}

	return object, acl
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
