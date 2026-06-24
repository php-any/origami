package annotation

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// CliApplicationClass 命令行应用入口注解（类似 Symfony Console Application）
// 标注命令行引导类，在加载时扫描命令注册、调用 boot()，并将 exit() 注册为 shutdown 回调。
//
// 用法示例：
//
//	#[CliApplication(name: "MyCLI", version: "1.0.0")]
//	class App {
//	    public static function boot(): void { ... }
//	    public static function exit(): void { ... }
//	}
type CliApplicationClass struct {
	node.Node
	construct data.Method
}

func (a *CliApplicationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newCliApplication()

	return data.NewClassValue(&CliApplicationClass{
		construct: &CliApplicationConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (a *CliApplicationClass) GetName() string    { return "Cli\\Annotation\\CliApplication" }
func (a *CliApplicationClass) GetExtend() *string { return nil }
func (a *CliApplicationClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (a *CliApplicationClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (a *CliApplicationClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (a *CliApplicationClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return a.construct, true
	}
	return nil, false
}
func (a *CliApplicationClass) GetMethods() []data.Method {
	return []data.Method{a.construct}
}
func (a *CliApplicationClass) GetConstruct() data.Method { return a.construct }

// CliApplication 命令行应用元信息
type CliApplication struct {
	name    string
	version string
	scan    string
	target  any // 被注解的引导类
}

// cliScanningActive 防止递归扫描
var cliScanningActive = false

// cliScanningDirs 记录正在扫描中的目录，防止递归触发
var cliScanningDirs = make(map[string]bool)

// registeredCliExitClasses 记录已注册 exit 回调的引导类，避免重复注册。
var registeredCliExitClasses = make(map[string]bool)

func newCliApplication() *CliApplication { return &CliApplication{name: "CLI", version: "1.0.0"} }

type CliApplicationConstructMethod struct{ app *CliApplication }

func (m *CliApplicationConstructMethod) GetName() string            { return "__construct" }
func (m *CliApplicationConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CliApplicationConstructMethod) GetIsStatic() bool          { return false }
func (m *CliApplicationConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewStringValue("CLI"), data.NewBaseType("string")),
		node.NewParameter(nil, "version", 1, data.NewStringValue("1.0.0"), data.NewBaseType("string")),
		node.NewParameter(nil, "scan", 2, data.NewNullValue(), nil),
		node.NewAnnotationTargetParameter(nil, 3),
	}
}
func (m *CliApplicationConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "version", 1, nil),
		node.NewVariable(nil, "scan", 2, nil),
		node.NewAnnotationTargetVariable(nil, 3),
	}
}
func (m *CliApplicationConstructMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *CliApplicationConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, name: 0"))
	}
	version, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, version: 1"))
	}

	if v, ok := name.(*data.StringValue); ok {
		m.app.name = v.AsString()
	}
	if v, ok := version.(*data.StringValue); ok {
		m.app.version = v.AsString()
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
	}

	// 编译模式下只解析参数，跳过扫描目录和调用 boot()
	if data.CompileMode {
		return nil, nil
	}

	// 如果指定了 scan 参数，则扫描目录
	if m.app.scan != "" {
		scanDir := m.app.scan
		if !filepath.IsAbs(scanDir) {
			cwd, _ := os.Getwd()
			scanDir = filepath.Join(cwd, scanDir)
		}
		scanDir = filepath.Clean(scanDir)

		// 防止递归扫描
		if cliScanningDirs[scanDir] {
			m.registerCliExit(ctx)
			return nil, nil
		}
		cliScanningDirs[scanDir] = true
		defer func() { delete(cliScanningDirs, scanDir) }()

		if acl := m.Scan(ctx); acl != nil {
			return nil, acl
		}
	}
	if acl := m.invokeBoot(ctx); acl != nil {
		return nil, acl
	}
	m.registerCliExit(ctx)
	return nil, nil
}

func (m *CliApplicationConstructMethod) invokeBoot(ctx data.Context) data.Control {
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

func (m *CliApplicationConstructMethod) registerCliExit(ctx data.Context) {
	cls, ok := m.app.target.(*node.ClassStatement)
	if !ok {
		return
	}

	className := cls.GetName()
	if registeredCliExitClasses[className] {
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

	ctx.GetVM().AddShutdownCallback(fv)
	registeredCliExitClasses[className] = true
}

// GetAppName 获取应用名称
func (m *CliApplicationConstructMethod) GetAppName() string {
	return m.app.name
}

// GetAppVersion 获取应用版本
func (m *CliApplicationConstructMethod) GetAppVersion() string {
	return m.app.version
}

// ParseArgs 解析命令行参数
func (m *CliApplicationConstructMethod) ParseArgs() []string {
	return os.Args[1:]
}

// Scan 扫描目录中的文件并加载命令
func (m *CliApplicationConstructMethod) Scan(ctx data.Context) data.Control {
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

	return nil
}
