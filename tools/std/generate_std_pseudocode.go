package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/php-any/origami/std/context"

	"github.com/php-any/origami/std/net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/channel"
	"github.com/php-any/origami/std/exception"
	"github.com/php-any/origami/std/log"
	"github.com/php-any/origami/std/reflect"
	oslib "github.com/php-any/origami/std/system/os"
)

// 手动配置标准库
func getStdClasses() []data.ClassStmt {
	return []data.ClassStmt{
		log.NewLogClass(),
		exception.NewExceptionClass(),
		oslib.NewOSClass(),
		&reflect.ReflectClass{},
		http.NewServerClass(), // 暂时注释掉，因为存在初始化问题
		http.NewRequestClass(nil, nil),
		http.NewResponseClass(nil, nil),
		channel.NewChannelClass(),
		// sql
		//sql.NewConnClass(),
		//sql.NewDBClass(),
		//sql.NewRowClass(),
		//sql.NewRowsClass(),
		//sql.NewStmtClass(),
		//sql.NewTxClass(),
		//sql.NewTxOptionsClass(),
	}
}

func getStdFunctions() []data.FuncStmt {
	return []data.FuncStmt{
		std.NewDumpFunction(),
		std.NewIncludeFunction(),
		// 导出 go 的 context
		context.NewBackgroundFunction(),
		context.NewWithCancelFunction(),
		context.NewWithTimeoutFunction(),
		context.NewWithValueFunction(),
		// sql
		// sql.NewOpenFunction(),
	}
}

// PseudoCode 表示伪代码结构
type PseudoCode struct {
	ModuleName  string
	Description string
	Namespace   string
	Functions   []FunctionSignature
	Classes     []ClassSignature
}

// FunctionSignature 表示函数签名
type FunctionSignature struct {
	Name       string
	Params     []Parameter
	ReturnType string
	Comment    string
}

// ClassSignature 表示类签名
type ClassSignature struct {
	Name        string
	ClassName   string
	Description string
	Methods     []MethodSignature
	Properties  []PropertySignature
}

// MethodSignature 表示方法签名
type MethodSignature struct {
	Name       string
	Params     []Parameter
	ReturnType string
	Modifier   string
	IsStatic   bool
	Comment    string
}

// PropertySignature 表示属性签名
type PropertySignature struct {
	Name     string
	Type     string
	Modifier string
	IsStatic bool
	Default  string
	Comment  string
}

// Parameter 表示参数
type Parameter struct {
	Name       string
	Type       string
	IsVariadic bool
}

// 分析类并生成伪代码
func analyzeClass(class data.ClassStmt) ClassSignature {
	className := class.GetName()

	// 处理命名空间
	var shortClassName string

	if strings.Contains(className, "\\") {
		parts := strings.Split(className, "\\")
		shortClassName = parts[len(parts)-1]
	} else {
		shortClassName = className
	}

	sig := ClassSignature{
		Name:        className,
		ClassName:   shortClassName,
		Description: fmt.Sprintf("%s 类", shortClassName),
	}

	// 获取方法
	methods := class.GetMethods()
	for _, method := range methods {
		// 检查方法是否为 nil
		if method == nil {
			continue
		}

		methodSig := MethodSignature{
			Name:     method.GetName(),
			Modifier: "public",
			IsStatic: method.GetIsStatic(),
		}

		// 分析参数
		params := method.GetParams()
		for _, param := range params {
			// 检查参数是否为 nil
			if param == nil {
				continue
			}

			paramName := fmt.Sprintf("param%d", len(methodSig.Params))
			paramType := "mixed"

			// 尝试从参数中获取名称和类型
			isVariadic := false
			if variable, ok := param.(data.Variable); ok {
				paramName = variable.GetName()
				if variable.GetType() != nil {
					paramType = variable.GetType().String()
				} else {
					paramType = "mixed"
				}
			} else if parameter, ok := param.(data.Parameter); ok {
				paramName = parameter.GetName()
				if parameter.GetType() != nil {
					paramType = parameter.GetType().String()
				} else {
					paramType = "mixed"
				}
			}

			// 检查是否为可变参数 (node.Parameters)
			if _, ok := param.(*node.Parameters); ok {
				isVariadic = true
			}

			methodSig.Params = append(methodSig.Params, Parameter{
				Name:       paramName,
				Type:       paramType,
				IsVariadic: isVariadic,
			})
		}

		// 分析返回值 - 尝试获取真实的返回类型
		returnType := "void"
		if returnTypeInterface, ok := method.(data.GetReturnType); ok {
			if retType := returnTypeInterface.GetReturnType(); retType != nil {
				returnType = retType.String()
			}
		}
		methodSig.ReturnType = returnType

		sig.Methods = append(sig.Methods, methodSig)
	}

	return sig
}

// 分析函数并生成伪代码
func analyzeFunction(fn data.FuncStmt) FunctionSignature {
	sig := FunctionSignature{
		Name: fn.GetName(),
	}

	// 分析参数
	params := fn.GetParams()
	for _, param := range params {
		paramName := fmt.Sprintf("param%d", len(sig.Params))
		paramType := "mixed"

		// 尝试从参数中获取名称和类型
		isVariadic := false
		if variable, ok := param.(data.Variable); ok {
			paramName = variable.GetName()
			if variable.GetType() != nil {
				paramType = variable.GetType().String()
			} else {
				paramType = "mixed"
			}
		} else if parameter, ok := param.(data.Parameter); ok {
			paramName = parameter.GetName()
			if parameter.GetType() != nil {
				paramType = parameter.GetType().String()
			} else {
				paramType = "mixed"
			}
		}

		// 检查是否为可变参数 (node.Parameters)
		if _, ok := param.(*node.Parameters); ok {
			isVariadic = true
		}

		sig.Params = append(sig.Params, Parameter{
			Name:       paramName,
			Type:       paramType,
			IsVariadic: isVariadic,
		})
	}

	// 分析返回值 - 尝试获取真实的返回类型
	returnType := "void"
	if returnTypeInterface, ok := fn.(data.GetReturnType); ok {
		if retType := returnTypeInterface.GetReturnType(); retType != nil {
			returnType = retType.String()
		}
	}
	sig.ReturnType = returnType

	return sig
}

// 生成 PHP 伪代码
func generatePHPPseudoCode(module PseudoCode) string {
	tmpl := `<?php
{{if .Namespace}}namespace {{.Namespace}};

{{end}}
/**
 * {{.ModuleName}} - {{.Description}}
 * 
 * 此文件包含 {{.ModuleName}} 模块的伪代码接口定义
 * 这些是自动生成的接口，仅用于参考，不包含具体实现
 */
{{if .Functions}}
{{range .Functions}}
/**
 * {{.Name}} 函数
 * {{if .Comment}}{{.Comment}}{{end}}
 */
function {{.Name}}({{range $i, $param := .Params}}{{if $i}}, {{end}}{{if ne $param.Type "mixed"}}{{$param.Type}} {{end}}{{if eq $param.IsVariadic true}}...{{end}}${{$param.Name}}{{end}}){{if ne .ReturnType "void"}} : {{.ReturnType}}{{end}} {
    // 实现逻辑
}
{{end}}
{{end}}
{{if .Classes}}
{{range .Classes}}
/**
 * {{.Name}} 类
 * {{.Description}}
 */
class {{.ClassName}} {
{{if .Properties}}{{range .Properties}}
    /**
     * {{.Name}} 属性
     * {{if .Comment}}{{.Comment}}{{end}}
     */
    {{.Modifier}} {{if .IsStatic}}static {{end}}${{.Name}}{{if .Type}} : {{.Type}}{{end}}{{if .Default}} = {{.Default}}{{end}};
{{end}}{{end}}{{if .Methods}}{{range .Methods}}
    /**
     * {{.Name}} 方法
     * {{if .Comment}}{{.Comment}}{{end}}
     */
    {{.Modifier}} {{if .IsStatic}}static {{end}}function {{.Name}}({{range $i, $param := .Params}}{{if $i}}, {{end}}{{if ne $param.Type "mixed"}}{{$param.Type}} {{end}}{{if eq $param.IsVariadic true}}...{{end}}${{$param.Name}}{{end}}){{if ne .ReturnType "void"}} : {{.ReturnType}}{{end}} {
        // 实现逻辑
    }
{{end}}{{end}}
}
{{end}}
{{end}}`

	t, err := template.New("php_pseudocode").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var buf strings.Builder
	err = t.Execute(&buf, module)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// 生成索引文件
func generatePseudoCodeIndex(modules []PseudoCode) string {
	tmpl := `# 标准库伪代码参考

Origami 标准库的伪代码接口定义。

## 模块列表

{{range .}}
### [{{.ModuleName}}]({{if .Namespace}}./{{.Namespace}}/{{.ModuleName}}.php{{else}}./{{.ModuleName}}.php{{end}})

{{.Description}}

{{end}}

## 快速开始

` + "`" + `php
<?php
// 使用标准库函数
dump("Hello World");

// 使用标准库类
$log = new Log();
$log->info("Application started");

// 使用反射
$reflect = new Reflect();
$classInfo = $reflect->getClassInfo("MyClass");
` + "`" + `

## 模块说明

{{range .}}
### {{.ModuleName}}

{{.Description}}

**主要功能：**
{{if .Functions}}
- 函数：{{range .Functions}}{{.Name}}{{end}}
{{end}}
{{if .Classes}}
- 类：{{range .Classes}}{{.Name}}{{end}}
{{end}}

[查看伪代码]({{if .Namespace}}./{{.Namespace}}/{{.ModuleName}}.php{{else}}./{{.ModuleName}}.php{{end}})
{{end}}
`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var buf strings.Builder
	err = t.Execute(&buf, modules)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func main() {
	// 创建 docs/std 目录
	err := os.MkdirAll("docs/std", 0755)
	if err != nil {
		panic(err)
	}

	var modules []PseudoCode

	// 分析标准库函数
	stdFunctions := getStdFunctions()
	if len(stdFunctions) > 0 {
		// 按命名空间分组函数
		funcModulesByNamespace := make(map[string]*PseudoCode)

		for _, fn := range stdFunctions {
			fullName := fn.GetName()
			namespace := ""
			shortName := fullName

			if strings.Contains(fullName, "\\") {
				parts := strings.Split(fullName, "\\")
				namespace = strings.Join(parts[:len(parts)-1], "\\")
				shortName = parts[len(parts)-1]
			}

			// 获取或创建对应命名空间的模块
			module, ok := funcModulesByNamespace[namespace]
			if !ok {
				module = &PseudoCode{
					ModuleName:  "functions",
					Description: "标准库函数",
					Namespace:   namespace,
				}
				funcModulesByNamespace[namespace] = module
			}

			// 分析函数并使用短名称写入
			sig := analyzeFunction(fn)
			sig.Name = shortName
			module.Functions = append(module.Functions, sig)
		}

		// 汇总加入模块列表
		for _, m := range funcModulesByNamespace {
			modules = append(modules, *m)
		}
	}

	// 分析标准库类
	stdClasses := getStdClasses()
	for _, class := range stdClasses {
		classSig := analyzeClass(class)

		// 按类名分组，处理命名空间
		className := class.GetName()
		var moduleName string
		var namespace string

		if strings.Contains(className, "\\") {
			parts := strings.Split(className, "\\")
			namespace = strings.Join(parts[:len(parts)-1], "\\")
			moduleName = strings.ToLower(parts[len(parts)-1])
		} else {
			moduleName = strings.ToLower(className)
		}

		// 查找或创建模块
		var module *PseudoCode
		for i := range modules {
			if modules[i].ModuleName == moduleName {
				module = &modules[i]
				break
			}
		}

		if module == nil {
			modules = append(modules, PseudoCode{
				ModuleName:  moduleName,
				Description: fmt.Sprintf("%s 类", className),
				Namespace:   namespace,
			})
			module = &modules[len(modules)-1]
		}

		module.Classes = append(module.Classes, classSig)
	}

	// 生成索引文件
	indexContent := generatePseudoCodeIndex(modules)
	err = os.WriteFile("docs/std/pseudo_README.md", []byte(indexContent), 0644)
	if err != nil {
		panic(err)
	}

	// 生成每个模块的 PHP 伪代码
	for _, module := range modules {
		if module.ModuleName == "" {
			continue
		}

		content := generatePHPPseudoCode(module)

		// 根据命名空间创建目录结构
		var filepath string
		if module.Namespace != "" {
			// 将命名空间转换为目录路径
			dirPath := strings.ReplaceAll(module.Namespace, "\\", "/")
			fullDirPath := fmt.Sprintf("docs/std/%s", dirPath)

			// 创建目录
			err = os.MkdirAll(fullDirPath, 0755)
			if err != nil {
				panic(err)
			}

			filepath = fmt.Sprintf("%s/%s.php", fullDirPath, strings.ToLower(module.ModuleName))
		} else {
			filepath = fmt.Sprintf("docs/std/%s.php", strings.ToLower(module.ModuleName))
		}

		err = os.WriteFile(filepath, []byte(content), 0644)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("标准库伪代码生成完成！")
	fmt.Println("生成的伪代码位于 docs/std/ 目录")
	fmt.Printf("共分析了 %d 个模块\n", len(modules))
}
