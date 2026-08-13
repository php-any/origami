package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ServerBootMethod 加载 #[Application] 引导类并注册注解路由。
// 扫描范围由引导类上的 #[Application(scan: ...)] 声明。
// 用法: $server->boot(SpringApplication::class)
type ServerBootMethod struct {
	server *ServerClass
}

func (h *ServerBootMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, acl := resolveBootApplication(ctx, 0)
	if acl != nil {
		return nil, acl
	}

	vm := ctx.GetVM()
	if acl := loadApplicationClass(vm, className); acl != nil {
		return nil, acl
	}

	return mountAnnotationRoutes(h.server, vm, ctx, "boot")
}

func (h *ServerBootMethod) GetName() string            { return "boot" }
func (h *ServerBootMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerBootMethod) GetIsStatic() bool          { return false }
func (h *ServerBootMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "application", 0, nil, data.NewBaseType("string")),
	}
}
func (h *ServerBootMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "application", 0, data.NewBaseType("string")),
	}
}
func (h *ServerBootMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
