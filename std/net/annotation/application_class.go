package annotation

import (
	"errors"

	"github.com/php-any/origami/utils"

	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

// ApplicationClass 应用入口注解（类似 Spring Boot 应用入口）
// 标注引导类，在 flash / app 加载时扫描控制器、调用 boot()，并将 exit() 注册为 shutdown 回调。
//
// 用法示例：
//
//	#[Application(name: "DemoApp", scan: __DIR__)]
//	class DemoApplication {
//	    public static function boot(): void { ... }
//	    public static function exit(): void { ... }
//	}
type ApplicationClass struct {
	node.Node
	construct data.Method
}

func (a *ApplicationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newApplication()

	return data.NewClassValue(&ApplicationClass{
		construct: &ApplicationConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (a *ApplicationClass) GetName() string    { return "Net\\Annotation\\Application" }
func (a *ApplicationClass) GetExtend() *string { return nil }
func (a *ApplicationClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (a *ApplicationClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (a *ApplicationClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (a *ApplicationClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return a.construct, true
	}
	return nil, false
}
func (a *ApplicationClass) GetMethods() []data.Method {
	return []data.Method{a.construct}
}
func (a *ApplicationClass) GetConstruct() data.Method { return a.construct }

// Application 应用入口元信息
type Application struct {
	name   string
	port   int64
	scan   string
	target any // 被注解的引导类
}

// scanningDirs 记录正在扫描中的目录，防止 main.php 在 scan 目录内时递归触发 Scan
var scanningDirs = make(map[string]bool)

func newApplication() *Application { return &Application{name: "App", port: 8080} }

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
		return nil, utils.NewThrow(errors.New("缺少参数, name: 0"))
	}
	port, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, port: 1"))
	}

	if v, ok := name.(*data.StringValue); ok {
		m.app.name = v.AsString()
	}
	if v, ok := port.(*data.IntValue); ok {
		m.app.port = int64(v.Value)
	}

	scanSpecified := false
	if scan, ok := ctx.GetIndexValue(2); ok && scan != nil {
		if anyV, ok := scan.(*data.StringValue); ok {
			m.app.scan = anyV.AsString()
			scanSpecified = true
		}
	}

	if target, ok := ctx.GetIndexValue(3); ok {
		if anyT, ok := target.(*data.AnyValue); ok {
			m.app.target = anyT.Value
		}
	}

	if !scanSpecified {
		if cls, ok := m.app.target.(*node.ClassStatement); ok {
			if from := cls.GetFrom(); from != nil {
				filePath := from.GetSource()
				if filePath != "" {
					m.app.scan = filepath.Dir(filePath)
				}
			}
		}
		if m.app.scan == "" {
			m.app.scan = "./src"
		}
	}

	scanDir := m.app.scan
	if !filepath.IsAbs(scanDir) {
		cwd, _ := os.Getwd()
		scanDir = filepath.Join(cwd, scanDir)
	}
	scanDir = filepath.Clean(scanDir)

	// main.php 位于 scan 目录内时，Scan 会再次加载本文件；此处跳过后续扫描与 boot，保证生命周期只执行一次。
	if scanningDirs[scanDir] {
		return nil, nil
	}
	scanningDirs[scanDir] = true
	defer func() { delete(scanningDirs, scanDir) }()

	if acl := m.Scan(ctx); acl != nil {
		return nil, acl
	}
	if acl := m.invokeBoot(ctx); acl != nil {
		return nil, acl
	}
	m.registerExit(ctx)
	return nil, nil
}

func (m *ApplicationConstructMethod) invokeBoot(ctx data.Context) data.Control {
	cls, ok := m.app.target.(*node.ClassStatement)
	if !ok {
		return nil
	}

	baseCtx := ctx.CreateBaseContext()
	cv := data.NewClassValue(cls, baseCtx)

	method, has := cls.GetStaticMethod("boot")
	if !has {
		return nil
	}

	fnCtx := cv.CreateContext(method.GetVariables())
	_, acl := method.Call(fnCtx)
	return acl
}

func (m *ApplicationConstructMethod) registerExit(ctx data.Context) {
	cls, ok := m.app.target.(*node.ClassStatement)
	if !ok {
		return
	}

	method, has := cls.GetStaticMethod("exit")
	if !has {
		return
	}

	fn, acl := node.NewStaticMethodFuncValue(cls, method).GetValue(ctx)
	if acl != nil {
		return
	}
	fv, ok := fn.(*data.FuncValue)
	if !ok {
		return
	}

	vm, ok := ctx.GetVM().(*runtime.VM)
	if !ok {
		return
	}
	vm.AddShutdownCallback(fv)
}

func (m *ApplicationConstructMethod) Scan(ctx data.Context) data.Control {
	cwd, err := os.Getwd()
	if err != nil {
		return utils.NewThrow(err)
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
		ext := filepath.Ext(path)
		switch ext {
		case ".zy", ".php":
			files = append(files, path)
			return nil
		default:
			return nil
		}
	})
	if err != nil {
		return utils.NewThrow(err)
	}

	vm := ctx.GetVM()
	for _, f := range files {
		f = utils.NormalizePhpFilePath(f)
		if vm.GetPhpFileCache(f) {
			continue
		}
		if _, acl := vm.LoadAndRun(f); acl != nil {
			return acl
		}
	}

	RegisterPendingRoutes(vm)
	return nil
}
