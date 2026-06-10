package compile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// generateMainFile 生成 main.go 入口文件
func generateMainFile(entryFile, outputDir, pkgName string) error {
	absEntry, err := filepath.Abs(entryFile)
	if err != nil {
		return fmt.Errorf("无法解析入口文件路径: %w", err)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n\n")
	b.WriteString("\t\"github.com/php-any/origami/parser\"\n")
	b.WriteString("\t\"github.com/php-any/origami/runtime\"\n")
	b.WriteString("\t\"github.com/php-any/origami/std\"\n")
	b.WriteString("\t\"github.com/php-any/origami/std/php\"\n")
	b.WriteString("\t\"github.com/php-any/origami/std/net/http\"\n")
	b.WriteString("\t\"github.com/php-any/origami/std/net/websocket\"\n")
	b.WriteString("\t\"github.com/php-any/origami/std/net/annotation\"\n")
	b.WriteString("\t\"github.com/php-any/origami/std/system\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func main() {\n")
	b.WriteString("\tp := parser.NewParser()\n")
	b.WriteString("\tvm := runtime.NewVM(p)\n")
	b.WriteString("\n")
	b.WriteString("\tstd.Load(vm)\n")
	b.WriteString("\tphp.Load(vm)\n")
	b.WriteString("\thttp.Load(vm)\n")
	b.WriteString("\twebsocket.Load(vm)\n")
	b.WriteString("\tannotation.Load(vm)\n")
	b.WriteString("\tsystem.Load(vm)\n")
	b.WriteString("\n")
	b.WriteString("\tRegister(vm)\n")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\t_, err := vm.RunCompiledFile(%q)\n", absEntry))
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tfmt.Fprintf(os.Stderr, \"错误: %v\\n\", err)\n")
	b.WriteString("\t\tos.Exit(1)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	return writeFormattedGoFile(filepath.Join(outputDir, "main.go"), []byte(b.String()))
}

// buildBinary 调用 go build 编译为二进制
func buildBinary(outputDir string) error {
	goCmd, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("未找到 go 命令: %w", err)
	}

	binaryName := "app"
	exeSuffix := ""
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}

	outputPath := filepath.Join(outputDir, binaryName+exeSuffix)
	fmt.Printf("正在解析依赖 ...\n")
	tidy := exec.Command(goCmd, "mod", "tidy")
	tidy.Dir = outputDir
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy 失败: %w", err)
	}

	fmt.Printf("正在编译 %s ...\n", outputPath)

	cmd := exec.Command(goCmd, "build", "-o", outputPath, ".")
	cmd.Dir = outputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build 失败: %w", err)
	}

	fmt.Printf("二进制文件已生成: %s\n", outputPath)
	return nil
}
