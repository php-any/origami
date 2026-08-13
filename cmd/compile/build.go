package compile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// generateMainFile 生成 main.go 入口文件（模板驱动）
func generateMainFile(parsed []ParsedFile, entryPaths map[string]bool, sourceDir, outputDir, pkgName string) error {
	tmpl, err := loadTemplate(sourceDir, "main.go.tmpl")
	if err != nil {
		return err
	}
	data := buildTemplateData(parsed, entryPaths, pkgName)
	content, err := renderTemplate(tmpl, data)
	if err != nil {
		return err
	}
	return writeFormattedGoFile(filepath.Join(outputDir, "main.go"), content)
}

// buildBinary 调用 go build 编译为二进制。
// outputDir 为 Go 源码所在目录，binaryDir 为最终二进制输出目录。
func buildBinary(outputDir, binaryDir string) error {
	goCmd, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("未找到 go 命令: %w", err)
	}

	binaryName := "app"
	exeSuffix := ""
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}

	binaryPath := filepath.Join(binaryDir, binaryName+exeSuffix)
	fmt.Printf("正在解析依赖 ...\n")
	tidy := exec.Command(goCmd, "mod", "tidy")
	tidy.Dir = outputDir
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy 失败: %w", err)
	}

	fmt.Printf("正在编译 %s ...\n", binaryPath)
	absPath, err := filepath.Abs(binaryPath)
	if err != nil {
		absPath = binaryPath
	}
	cmd := exec.Command(goCmd, "build", "-o", absPath, ".")
	cmd.Dir = outputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build 失败: %w", err)
	}

	fmt.Printf("二进制文件已生成: %s\n", binaryPath)
	return nil
}
