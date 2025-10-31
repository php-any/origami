package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// ControllerClass Controller注解类 - 特性注解
type ControllerClass struct {
	node.Node
	process   data.Method
	register  data.Method
	construct data.Method
}

func (c *ControllerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newController()

	return data.NewClassValue(&ControllerClass{
		process:   &ControllerProcessMethod{source},
		register:  &ControllerRegisterMethod{source},
		construct: &ControllerConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (c *ControllerClass) GetName() string { return "Net\\Annotation\\Controller" }

func (c *ControllerClass) GetExtend() *string {
	return nil
}

func (c *ControllerClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (c *ControllerClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *ControllerClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *ControllerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return c.process, true
	case "register":
		return c.register, true
	}
	return nil, false
}

func (c *ControllerClass) GetMethods() []data.Method {
	return []data.Method{
		c.process,
		c.register,
		c.construct,
	}
}

func (c *ControllerClass) GetConstruct() data.Method {
	return c.construct
}

// Controller 控制器实例
type Controller struct {
	name string
}

func newController() *Controller {
	return &Controller{}
}

// ControllerConstructMethod 构造函数 - 特性注解只接收注解参数
type ControllerConstructMethod struct {
	controller *Controller
}

func (m *ControllerConstructMethod) GetName() string {
	return "__construct"
}

func (m *ControllerConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ControllerConstructMethod) GetIsStatic() bool {
	return false
}

func (m *ControllerConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}

func (m *ControllerConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *ControllerConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ControllerConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var vm *runtime.TempVM
	if temp, ok := ctx.GetVM().(*runtime.TempVM); ok {
		vm = temp
	} else {
		return nil, utils.NewThrow(errors.New("@Controller 注解只能在 app() 内加载"))
	}
	// 读取 name
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, name: 0"))
	}
	if v, ok := a0.(*data.StringValue); ok {
		m.controller.name = v.AsString()
	}

	// 读取类节点并统一处理该类的方法路由：扫描 @Route 前缀 + @GetMapping/@PostMapping 路径 → router.routes
	if tv, ok := ctx.GetIndexValue(1); ok {
		if anyT, ok := tv.(*data.AnyValue); ok {
			if cls, ok := anyT.Value.(*node.ClassStatement); ok {
				// 读取类注解中的 @Route 前缀（若存在）
				prefix := ""
				for _, ann := range cls.Annotations {
					switch rc := ann.Class.(type) {
					case *RouteClass:
						prefix = rc.Prefix()
					}
				}

				// 规范化前缀
				normalize := func(p string) string {
					if p == "" {
						return ""
					}
					if p == "/" {
						return "/"
					}
					// 简单规范：确保单前导斜杠，去尾部斜杠
					if p[0] != '/' {
						p = "/" + p
					}
					for len(p) > 1 && p[len(p)-1] == '/' {
						p = p[:len(p)-1]
					}
					return p
				}
				join := func(prefix, path string) string {
					P := normalize(prefix)
					if path == "" {
						path = "/"
					}
					if path[0] != '/' {
						path = "/" + path
					}
					if P == "" || P == "/" {
						return path
					}
					if path == "/" {
						return P
					}
					return P + path
				}

				// 遍历类方法，读取方法注解并注册到 router.routes
				for _, method := range cls.GetMethods() {
					if cm, ok := method.(*node.ClassMethod); ok {
						for _, ann := range cm.Annotations {
							switch gc := ann.Class.(type) {
							case *GetMappingClass:
								full := join(prefix, gc.Path())
								vm.Cache = append(vm.Cache, runtime.Route{Method: "GET", Path: full, Target: method})
							case *PostMappingClass:
								full := join(prefix, gc.Path())
								vm.Cache = append(vm.Cache, runtime.Route{Method: "POST", Path: full, Target: method})
							case *PutMappingClass:
								full := join(prefix, gc.Path())
								vm.Cache = append(vm.Cache, runtime.Route{Method: "PUT", Path: full, Target: method})
							case *DeleteMappingClass:
								full := join(prefix, gc.Path())
								vm.Cache = append(vm.Cache, runtime.Route{Method: "DELETE", Path: full, Target: method})
							}
						}
					}
				}
			}
		}
	}

	return data.NewStringValue("Controller annotation constructed with name: " + m.controller.name), nil
}

// ControllerProcessMethod 处理注解的方法
type ControllerProcessMethod struct {
	controller *Controller
}

func (m *ControllerProcessMethod) GetName() string {
	return "process"
}

func (m *ControllerProcessMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ControllerProcessMethod) GetIsStatic() bool {
	return false
}

func (m *ControllerProcessMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ControllerProcessMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ControllerProcessMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ControllerProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解处理逻辑
	return data.NewStringValue("Controller processed with name: " + m.controller.name), nil
}

// ControllerRegisterMethod 注册控制器的方法
type ControllerRegisterMethod struct {
	controller *Controller
}

func (m *ControllerRegisterMethod) GetName() string {
	return "register"
}

func (m *ControllerRegisterMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ControllerRegisterMethod) GetIsStatic() bool {
	return false
}

func (m *ControllerRegisterMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ControllerRegisterMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ControllerRegisterMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ControllerRegisterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解注册逻辑
	return data.NewStringValue("Controller registered with name: " + m.controller.name), nil
}
