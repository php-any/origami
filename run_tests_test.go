package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/php-any/origami/cmd"
)

// TestRunTests 与主程序（zy.go 中的 cmd.RunScriptFile）完全一致的方式，
// 在进程内直接执行 ./tests/run_tests.php 并验收测试结果，无需重新构建二进制。
//
// 验收标准与 tests/run_tests.php 的失败机制一一对应：
//  1. 套件或任一用例失败时，run_tests.php 会调用 Log::fatal，而 Log::fatal 会
//     os.Exit(1) 直接终止当前进程 —— 因此任何失败都会让 go test 进程以非零码
//     退出，测试自然 FAIL；代码能执行到断言处，本身就说明套件全部通过；
//  2. 输出不得出现 "[FATAL]" 日志（Log::fatal 的输出前缀，双保险）；
//  3. 输出末尾必须包含完成标记 "接口功能测试完成"，确保整个套件完整跑完、
//     未被中途打断。
//
// Log 实例在创建 VM（NewLog）时捕获当时的 os.Stdout，因此这里先把 os.Stdout /
// os.Stderr 重定向到临时文件，再执行脚本，从而捕获测试套件的全部输出。
func TestRunTests(t *testing.T) {
	logFile, err := os.CreateTemp(t.TempDir(), "run_tests_*.log")
	if err != nil {
		t.Fatalf("创建日志文件失败: %v", err)
	}
	defer logFile.Close()

	oldStdout, oldStderr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = logFile, logFile
	defer func() {
		os.Stdout, os.Stderr = oldStdout, oldStderr
	}()

	// 与主程序相同的方式直接执行测试套件
	err = cmd.RunScriptFile(filepath.Join("tests", "run_tests.php"))
	if err != nil {
		t.Fatalf("tests/run_tests.php 执行失败: %v", err)
	}

	// 读取套件全部输出
	if _, err := logFile.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("读取日志失败: %v", err)
	}
	output, err := io.ReadAll(logFile)
	if err != nil {
		t.Fatalf("读取日志失败: %v", err)
	}
	outStr := string(output)

	if strings.Contains(outStr, "[FATAL]") {
		t.Fatalf("tests/run_tests.php 输出包含 FATAL 日志，存在失败用例:\n%s", outStr)
	}
	if !strings.Contains(outStr, "接口功能测试完成") {
		t.Fatalf("未检测到完成标记，测试套件可能提前中断:\n%s", outStr)
	}
	t.Log("tests/run_tests.php 全部用例通过")
}
