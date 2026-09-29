// Package httpkernel 提供 Illuminate\Foundation\Http\Kernel 的 Go 实现。
package httpkernel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	fqnKernel         = "Illuminate\\Foundation\\Http\\Kernel"
	fqnKernelContract = "Illuminate\\Contracts\\Http\\Kernel"
)

// 默认 bootstrap 类，对齐 Foundation\Http\Kernel::$bootstrappers。
var defaultBootstrappers = []string{
	"Illuminate\\Foundation\\Bootstrap\\LoadEnvironmentVariables",
	"Illuminate\\Foundation\\Bootstrap\\LoadConfiguration",
	"Illuminate\\Foundation\\Bootstrap\\HandleExceptions",
	"Illuminate\\Foundation\\Bootstrap\\RegisterFacades",
	"Illuminate\\Foundation\\Bootstrap\\RegisterProviders",
	"Illuminate\\Foundation\\Bootstrap\\BootProviders",
}

// 默认中间件优先级，对齐 Foundation\Http\Kernel::$middlewarePriority。
var defaultMiddlewarePriority = []string{
	"Illuminate\\Foundation\\Http\\Middleware\\HandlePrecognitiveRequests",
	"Illuminate\\Cookie\\Middleware\\EncryptCookies",
	"Illuminate\\Cookie\\Middleware\\AddQueuedCookiesToResponse",
	"Illuminate\\Session\\Middleware\\StartSession",
	"Illuminate\\View\\Middleware\\ShareErrorsFromSession",
	"Illuminate\\Contracts\\Auth\\Middleware\\AuthenticatesRequests",
	"Illuminate\\Routing\\Middleware\\ThrottleRequests",
	"Illuminate\\Routing\\Middleware\\ThrottleRequestsWithRedis",
	"Illuminate\\Contracts\\Session\\Middleware\\AuthenticatesSessions",
	"Illuminate\\Routing\\Middleware\\SubstituteBindings",
	"Illuminate\\Auth\\Middleware\\Authorize",
}

// kernelState 保存 Kernel 实例可变状态。
type kernelState struct {
	app                data.Value
	router             data.Value
	bootstrappers      []string
	middleware         []string
	middlewareGroups   map[string][]string
	middlewareAliases  map[string]string
	middlewarePriority []string
	bootstrapped       bool
	// tel 由常驻 state 持有、请求沙箱共享指针（见 Sandbox）：
	// Telescope 的开关与匹配 pattern 都是启动期常量，没理由每请求重算一遍。
	tel *telescopeCache
}

func newKernelState() *kernelState {
	return &kernelState{
		bootstrappers:      append([]string(nil), defaultBootstrappers...),
		middleware:         nil,
		middlewareGroups:   make(map[string][]string),
		middlewareAliases:  make(map[string]string),
		middlewarePriority: append([]string(nil), defaultMiddlewarePriority...),
		tel:                &telescopeCache{},
	}
}

// KernelClass 实现 data.ClassStmt：Illuminate\Foundation\Http\Kernel。
type KernelClass struct {
	node.Node
	state      *kernelState
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

// NewClass 导出供 go-support.Load 注册。
func NewClass() data.ClassStmt {
	return newKernelClass(nil)
}

func newKernelClass(state *kernelState) *KernelClass {
	c := &KernelClass{
		state: state,
		properties: []data.Property{
			node.NewProperty(nil, "app", "protected", false, data.NewNullValue()),
			node.NewProperty(nil, "router", "protected", false, data.NewNullValue()),
			node.NewProperty(nil, "bootstrappers", "protected", false, stringsToArrayValue(defaultBootstrappers)),
			node.NewProperty(nil, "middleware", "protected", false, data.NewArrayValue(nil)),
			node.NewProperty(nil, "middlewareGroups", "protected", false, data.NewArrayValue(nil)),
			node.NewProperty(nil, "middlewareAliases", "protected", false, data.NewArrayValue(nil)),
			node.NewProperty(nil, "middlewarePriority", "protected", false, stringsToArrayValue(defaultMiddlewarePriority)),
		},
	}
	c.methods, c.methodList = kernelMethods()
	return c
}

func (c *KernelClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(newKernelClass(newKernelState()), ctx.CreateBaseContext()), nil
}

func (c *KernelClass) GetName() string         { return fqnKernel }
func (c *KernelClass) GetExtend() *string      { return nil }
func (c *KernelClass) GetImplements() []string { return []string{fqnKernelContract} }
func (c *KernelClass) GetSource() any          { return c.state }
func (c *KernelClass) GetConstruct() data.Method {
	return c.methods["__construct"]
}
func (c *KernelClass) GetPropertyList() []data.Property { return c.properties }
func (c *KernelClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *KernelClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *KernelClass) GetMethods() []data.Method { return c.methodList }
