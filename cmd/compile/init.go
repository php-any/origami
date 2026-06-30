package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// presetTemplates 内置模板预设
var presetTemplates = map[string]map[string]string{
	"default": {
		"register.go.tmpl": defaultRegisterTmpl,
		"main.go.tmpl":     defaultMainTmpl,
		"go.mod.tmpl":      defaultModTmpl,
	},
	"minimal": {
		// 最小预设：仅 go.mod，register 和 main 由项目手写
		"go.mod.tmpl": defaultModTmpl,
	},
	"library": {
		// 库预设：仅 register + go.mod，main 通常不需要
		"register.go.tmpl": defaultRegisterTmpl,
		"go.mod.tmpl":      defaultModTmpl,
	},
	"fyne": {
		"main.go.tmpl": `package main

import (
	"fmt"
	"os"

	fyne "github.com/php-any/origami-fyne"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	system.Load(vm)
	fyne.Load(vm)

	Register(vm)
{{- if .HasEntry}}
	_, ctrl := vm.RunCompiledFile(EntryPath)
	if data.FlushAllBuffersFn != nil {
		data.FlushAllBuffersFn()
	}
	if ctrl != nil {
		fmt.Fprintf(os.Stderr, "run failed\n")
		p.ShowControl(ctrl)
		os.Exit(1)
	}
{{- end}}
	vm.RunShutdownCallbacks()
}
`,
		"go.mod.tmpl": `module {{.Pkg}}

go 1.25.0

require (
	github.com/php-any/origami v0.0.0
	github.com/php-any/origami-fyne v0.0.0
)

require fyne.io/fyne/v2 v2.5.5
`,
	},
	"web": {
		"main.go.tmpl": `package main

import (
	"fmt"
	"os"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/net/annotation"
	"github.com/php-any/origami/std/system"
)

func main() {
	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	http.Load(vm)
	websocket.Load(vm)
	annotation.Load(vm)
	system.Load(vm)

	Register(vm)
{{- if .HasEntry}}
	_, err := vm.RunCompiledFile(EntryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
{{- end}}
}
`,
		"go.mod.tmpl": defaultModTmpl,
	},
	"cli": {
		"main.go.tmpl": `package main

import (
	"fmt"
	"os"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	system.Load(vm)

	Register(vm)
{{- if .HasEntry}}
	_, err := vm.RunCompiledFile(EntryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
{{- end}}
}
`,
		"go.mod.tmpl": defaultModTmpl,
	},
}

// presetNames 返回排序后的预设名列表
func presetNames() []string {
	names := make([]string, 0, len(presetTemplates))
	for k := range presetTemplates {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// NewInitCommand 创建 init 子命令
func NewInitCommand() *cobra.Command {
	var (
		template string
		force    bool
		list     bool
	)

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "初始化项目的 .zy/ 编译模板",
		Long: `在目标目录下创建 .zy/ 目录并生成编译模板文件。

--template 支持：
  内置预设:  default, minimal, library, fyne, web, cli
  自定义目录:  --template=./my-templates（读取目录下所有 .tmpl 文件）

默认不会覆盖已有文件，使用 --force 可强制覆盖。
使用 --list 查看所有内置预设。

示例:
  zy init .                        # 默认模板
  zy init my-app --template=fyne   # Fyne 桌面应用
  zy init . --template=./templates # 自定义模板目录
  zy init . --force                # 强制覆盖`,
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if list {
				fmt.Println("内置模板预设:")
				for _, name := range presetNames() {
					files := presetTemplates[name]
					keys := make([]string, 0, len(files))
					for k := range files {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					fmt.Printf("  %-10s → %s\n", name, strings.Join(keys, ", "))
				}
				return nil
			}

			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("解析目录失败: %w", err)
			}

			if info, err := os.Stat(absDir); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("目录不存在: %s", absDir)
				}
				return fmt.Errorf("无法访问目录: %w", err)
			} else if !info.IsDir() {
				return fmt.Errorf("不是目录: %s", absDir)
			}

			zyDir := filepath.Join(absDir, ".zy")

			// 解析模板源
			files, source, err := resolveTemplateSource(template)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(zyDir, 0755); err != nil {
				return fmt.Errorf("创建 .zy/ 目录失败: %w", err)
			}

			fmt.Printf("模板源: %s\n", source)
			for name, content := range files {
				target := filepath.Join(zyDir, name)
				if !force {
					if _, err := os.Stat(target); err == nil {
						fmt.Printf("  跳过: .zy/%s\n", name)
						continue
					}
				}
				if err := os.WriteFile(target, []byte(content), 0644); err != nil {
					return fmt.Errorf("写入 .zy/%s 失败: %w", name, err)
				}
				fmt.Printf("  创建: .zy/%s\n", name)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&template, "template", "t", "default", "模板来源 (预设名或自定义目录路径)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制覆盖已有文件")
	cmd.Flags().BoolVarP(&list, "list", "l", false, "列出所有内置预设")

	return cmd
}

// resolveTemplateSource 解析模板源，返回 文件名→内容 映射和来源描述
func resolveTemplateSource(tmpl string) (map[string]string, string, error) {
	// 1. 先检查是否是内置预设
	if preset, ok := presetTemplates[tmpl]; ok {
		return preset, fmt.Sprintf("内置预设 (%s)", tmpl), nil
	}

	// 2. 检查是否自定义模板目录
	absPath, err := filepath.Abs(tmpl)
	if err != nil {
		return nil, "", fmt.Errorf("解析模板路径失败: %w", err)
	}
	if info, err := os.Stat(absPath); err == nil && info.IsDir() {
		files, err := loadTemplateDir(absPath)
		if err != nil {
			return nil, "", err
		}
		if len(files) == 0 {
			return nil, "", fmt.Errorf("模板目录 %s 中未找到 .tmpl 文件", absPath)
		}
		return files, fmt.Sprintf("自定义目录 (%s)", absPath), nil
	}

	return nil, "", fmt.Errorf(
		"未知模板: %s\n使用 --list 查看内置预设，或指定自定义模板目录路径",
		tmpl,
	)
}

// loadTemplateDir 从目录加载所有 .tmpl 文件
func loadTemplateDir(dir string) (map[string]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取模板目录失败: %w", err)
	}

	files := make(map[string]string)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tmpl") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("读取 %s 失败: %w", e.Name(), err)
		}
		files[e.Name()] = string(data)
	}
	return files, nil
}
