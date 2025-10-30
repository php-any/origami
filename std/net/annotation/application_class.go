package annotation

import (
	"errors"
	"github.com/php-any/origami/runtime"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ApplicationClass 应用入口注解（类似 Spring Boot 应用入口）
// 默认作为“特性注解”，但构造函数声明了 target，调度器会按需注入被注解的 AST 目标
// 用法示例：
// @Application(name: "DemoApp", port: 8081, scan: ["App\\Controller\\"])
// class Main {}
type ApplicationClass struct {
	node.Node
	process   data.Method
	register  data.Method
	construct data.Method
}

func (a *ApplicationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newApplication()

	return data.NewClassValue(&ApplicationClass{
		process:   &ApplicationProcessMethod{source},
		register:  &ApplicationRegisterMethod{source},
		construct: &ApplicationConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (a *ApplicationClass) GetName() string                            { return "Net\\Annotation\\Application" }
func (a *ApplicationClass) GetExtend() *string                         { return nil }
func (a *ApplicationClass) GetImplements() []string                    { return []string{node.TypeFeature} }
func (a *ApplicationClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (a *ApplicationClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (a *ApplicationClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return a.process, true
	case "register":
		return a.register, true
	case "__construct":
		return a.construct, true
	}
	return nil, false
}
func (a *ApplicationClass) GetMethods() []data.Method {
	return []data.Method{a.process, a.register, a.construct}
}
func (a *ApplicationClass) GetConstruct() data.Method { return a.construct }

// Application 应用入口元信息
type Application struct {
	name   string
	port   int64
	scan   string
	target any // 被注解的目标（通常是类）
}

func newApplication() *Application { return &Application{name: "App", port: 8080} }

// 构造函数：接收参数与 target（如存在）
type ApplicationConstructMethod struct{ app *Application }

func (m *ApplicationConstructMethod) GetName() string            { return "__construct" }
func (m *ApplicationConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ApplicationConstructMethod) GetIsStatic() bool          { return false }
func (m *ApplicationConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewStringValue("App"), data.NewBaseType("string")),
		node.NewParameter(nil, "port", 1, data.NewIntValue(8080), data.NewBaseType("int")),
		node.NewParameter(nil, "scan", 2, data.NewNullValue(), nil),
		node.NewParameter(nil, node.TargetName, 3, nil, nil),
	}
}
func (m *ApplicationConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "port", 1, nil),
		node.NewVariable(nil, "scan", 2, nil),
		node.NewVariable(nil, "target", 3, nil),
	}
}
func (m *ApplicationConstructMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *ApplicationConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, name: 0"))
	}
	port, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, port: 1"))
	}

	if v, ok := name.(*data.StringValue); ok {
		m.app.name = v.AsString()
	}
	if v, ok := port.(*data.IntValue); ok {
		m.app.port = int64(v.Value)
	}

	m.app.scan = "./src/controllers"
	if scan, ok := ctx.GetIndexValue(2); ok && scan != nil {
		if anyV, ok := scan.(*data.StringValue); ok {
			m.app.scan = anyV.AsString()
		}
	}

	if target, ok := ctx.GetIndexValue(3); ok {
		if anyT, ok := target.(*data.AnyValue); ok {
			m.app.target = anyT.Value
		}
	}

	// 若被标注的是函数 main，则在函数体首部注入一次性启动逻辑调用：__spring_bootstrap($request, $response)
	switch fn := m.app.target.(type) {
	case *node.FunctionStatement:
		fn.Body = append(m.BuildBoot(ctx), fn.Body...)
	}

	return nil, m.Scan(ctx)
}

func (m *ApplicationConstructMethod) BuildBoot(ctx data.Context) []data.GetValue {
	if temp, ok := ctx.GetVM().(*runtime.TempVM); ok {
		return []data.GetValue{
			&RegisterRoute{
				vm: temp,
			},
		}
	}
	return []data.GetValue{}
}

// 扫描目录下 *.zy 相关文件
func (m *ApplicationConstructMethod) Scan(ctx data.Context) data.Control {
	// 递归扫描 m.app.scan（默认 ./src）下的 .zy 文件（包含任意子目录），不做 http.zy 等特殊过滤
	cwd, err := os.Getwd()
	if err != nil {
		return data.NewErrorThrow(nil, err)
	}

	base := m.app.scan
	if !filepath.IsAbs(base) {
		base = filepath.Join(cwd, base)
	}

	var files []string
	err = filepath.WalkDir(base, func(path string, d os.DirEntry, wErr error) error {
		if wErr != nil {
			return wErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".zy" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return data.NewErrorThrow(nil, err)
	}

	vm := ctx.GetVM()
	for _, f := range files {
		if _, acl := vm.LoadAndRun(f); acl != nil {
			return acl
		}
	}
	return nil
}

// 处理方法（保留占位，供框架在启动时选择性触发）
type ApplicationProcessMethod struct{ app *Application }

func (m *ApplicationProcessMethod) GetName() string               { return "process" }
func (m *ApplicationProcessMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ApplicationProcessMethod) GetIsStatic() bool             { return false }
func (m *ApplicationProcessMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *ApplicationProcessMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *ApplicationProcessMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *ApplicationProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue("Application processed: " + m.app.name), nil
}

// 注册方法（保留占位，供框架在启动时选择性触发）
type ApplicationRegisterMethod struct{ app *Application }

func (m *ApplicationRegisterMethod) GetName() string               { return "register" }
func (m *ApplicationRegisterMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ApplicationRegisterMethod) GetIsStatic() bool             { return false }
func (m *ApplicationRegisterMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *ApplicationRegisterMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *ApplicationRegisterMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *ApplicationRegisterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue("Application registered on port " + data.NewIntValue(int(m.app.port)).AsString()), nil
}
