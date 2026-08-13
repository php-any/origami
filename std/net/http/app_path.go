package http

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// resolveBootApplication 解析 boot 的应用引导类 FQCN。
func resolveBootApplication(ctx data.Context, argIndex int) (string, data.Control) {
	v, ok := ctx.GetIndexValue(argIndex)
	if !ok {
		return "", utils.NewThrow(errors.New("boot 需要应用引导类名"))
	}
	s, ok := v.(data.AsString)
	if !ok {
		return "", utils.NewThrow(errors.New("boot 参数必须是类名字符串"))
	}
	className := s.AsString()
	if className == "" {
		return "", utils.NewThrow(errors.New("boot 引导类名不能为空"))
	}
	return className, nil
}

// loadApplicationClass 加载带 #[Application] 的引导类（解析时触发扫描与生命周期）。
func loadApplicationClass(vm data.VM, className string) data.Control {
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return acl
	}
	if stmt == nil {
		return utils.NewThrowf("无法加载引导类: %s", className)
	}
	return nil
}

// resolveAppDirPath 解析 flash 的 dir 参数为绝对目录，并校验其中存在 main.php。
func resolveAppDirPath(ctx data.Context, argIndex int, defaultDir string) (string, data.Control) {
	var dir string
	if dirValue, ok := ctx.GetIndexValue(argIndex); ok {
		if pathStr, ok := dirValue.(data.AsString); ok {
			dir = pathStr.AsString()
		}
	}
	if dir == "" {
		dir = defaultDir
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", utils.NewThrowf("获取当前目录失败: %v", err)
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(currentDir, dir)
	}

	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", utils.NewThrowf("目录(%s)不存在", dir)
		}
		return "", utils.NewThrowf("无法访问目录: %s", dir)
	}
	if !info.IsDir() {
		return "", utils.NewThrowf("flash 需要目录路径，收到: %s", dir)
	}

	mainFile := filepath.Join(dir, "main.php")
	if _, err := os.Stat(mainFile); err != nil {
		if os.IsNotExist(err) {
			return "", utils.NewThrowf("应用目录(%s)中缺少 main.php", dir)
		}
		return "", utils.NewThrowf("无法访问文件: %s", mainFile)
	}
	return dir, nil
}
