package http

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ServerFlashMethod 启动时扫描注解路由并直接注册到 Server
// 支持多个独立应用：扫描目录下所有 main.php，每个 main.php 的 #[Application] 独立扫描各自的控制器目录
// 用法: $server->flash("./src") 或 $server->flash("./apps")；单应用推荐 $server->boot(ApplicationClass::class)
type ServerFlashMethod struct {
	server *ServerClass
}

func (h *ServerFlashMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	dir, acl := resolveAppDirPath(ctx, 0, "./src")
	if acl != nil {
		return nil, acl
	}

	vm := ctx.GetVM()

	mainFiles, err := findMainFiles(dir)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if len(mainFiles) == 0 {
		return nil, utils.NewThrowf("flash: 目录 %s 中未找到 main.php", dir)
	}

	for _, mainFile := range mainFiles {
		if _, acl := vm.LoadAndRun(mainFile); acl != nil {
			return nil, acl
		}
	}

	return mountAnnotationRoutes(h.server, vm, ctx, "flash")
}

// findMainFiles 递归扫描目录下所有 main.php 文件
func findMainFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.php" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (h *ServerFlashMethod) GetName() string            { return "flash" }
func (h *ServerFlashMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerFlashMethod) GetIsStatic() bool          { return false }
func (h *ServerFlashMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "dir", 0, data.NewStringValue("./src"), data.String{}),
	}
}
func (h *ServerFlashMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "dir", 0, nil),
	}
}
func (h *ServerFlashMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
