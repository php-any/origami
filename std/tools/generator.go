package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Generate(instance interface{}, config *GeneratorConfig) error {
	generator := NewWrapperGenerator()
	files, err := generator.Generate(instance, config)
	if err != nil {
		return err
	}

	for _, file := range files {
		// 检查文件是否已存在
		if _, err := os.Stat(file.FilePath); err == nil {
			fmt.Printf("文件已存在，跳过: %s\n", file.FilePath)
			continue
		}

		err := writeFile(file)
		if err != nil {
			return err
		} else {
			fmt.Printf("已生成文件: %s\n", file.FilePath)
		}
	}
	return nil
}

func writeFile(file GeneratedFile) error {
	dir := filepath.Dir(file.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(file.FilePath, []byte(file.Content), 0644)
}

// WrapperGenerator 包装器生成器
type WrapperGenerator struct {
	analyzer *StructAnalyzer
}

// NewWrapperGenerator 创建新的包装器生成器
func NewWrapperGenerator() *WrapperGenerator {
	return &WrapperGenerator{
		analyzer: NewStructAnalyzer(),
	}
}

// Generate 生成所有文件
func (g *WrapperGenerator) Generate(instance interface{}, config *GeneratorConfig) ([]GeneratedFile, error) {
	// 分析结构体
	structInfo, err := g.analyzer.AnalyzeStruct(instance)
	if err != nil {
		return nil, fmt.Errorf("分析结构体失败: %v", err)
	}

	// 合并配置中的属性
	if config.Properties != nil {
		structInfo.Properties = config.Properties
	}

	var files []GeneratedFile

	// 生成类定义文件
	classFile, err := g.GenerateClass(structInfo, config)
	if err != nil {
		return nil, fmt.Errorf("生成类定义失败: %v", err)
	}
	files = append(files, *classFile)

	// 生成方法包装器文件
	methodFiles, err := g.GenerateMethods(structInfo, config)
	if err != nil {
		return nil, fmt.Errorf("生成方法包装器失败: %v", err)
	}
	files = append(files, methodFiles...)

	return files, nil
}

// GenerateClass 生成类定义文件
func (g *WrapperGenerator) GenerateClass(structInfo *StructInfo, config *GeneratorConfig) (*GeneratedFile, error) {
	tmpl := GetClassTemplate()

	var buf strings.Builder
	err := tmpl.Execute(&buf, map[string]interface{}{
		"PackageName": config.PackageName,
		"StructName":  config.StructName,
		"ClassName":   config.ClassName,
		"Methods":     structInfo.Methods,
		"Properties":  structInfo.Properties,
	})

	if err != nil {
		return nil, fmt.Errorf("模板执行失败: %v", err)
	}

	fileName := fmt.Sprintf("%s_class.go", strings.ToLower(config.StructName))
	filePath := fmt.Sprintf("%s/%s", config.OutputDir, fileName)

	return &GeneratedFile{
		FileName: fileName,
		Content:  buf.String(),
		FilePath: filePath,
		FileType: "class",
	}, nil
}

// GenerateMethods 生成所有方法包装器文件
func (g *WrapperGenerator) GenerateMethods(structInfo *StructInfo, config *GeneratorConfig) ([]GeneratedFile, error) {
	var files []GeneratedFile

	for _, method := range structInfo.Methods {
		methodFile, err := g.GenerateMethodWrapper(method, config)
		if err != nil {
			return nil, fmt.Errorf("生成方法 %s 包装器失败: %v", method.Name, err)
		}
		files = append(files, *methodFile)
	}

	return files, nil
}

// GenerateMethodWrapper 生成单个方法包装器文件
func (g *WrapperGenerator) GenerateMethodWrapper(method MethodInfo, config *GeneratorConfig) (*GeneratedFile, error) {
	// 生成参数检查代码
	paramChecks := g.generateParamChecks(method)

	// 生成原始方法调用代码
	sourceCall, ret := g.generateSourceCall(method, config)

	// 生成模板数据
	templateData := map[string]interface{}{
		"PackageName": config.PackageName,
		"StructName":  config.StructName,
		"MethodName":  method.Name,
		"Params":      method.Params,
		"ParamChecks": paramChecks,
		"SourceCall":  sourceCall,
		"Return":      ret,
	}

	// 执行模板
	tmpl := GetMethodTemplate()
	var buf strings.Builder
	err := tmpl.Execute(&buf, templateData)
	if err != nil {
		return nil, fmt.Errorf("模板执行失败: %v", err)
	}

	fileName := fmt.Sprintf("%s_%s_method.go",
		strings.ToLower(config.StructName),
		strings.ToLower(method.Name))
	filePath := fmt.Sprintf("%s/%s", config.OutputDir, fileName)

	return &GeneratedFile{
		FileName: fileName,
		Content:  buf.String(),
		FilePath: filePath,
		FileType: "method",
	}, nil
}

// generateParamChecks 生成参数检查代码
func (g *WrapperGenerator) generateParamChecks(method MethodInfo) []string {
	var checks []string

	for i, param := range method.Params {
		check := fmt.Sprintf(`a%d, ok := ctx.GetIndexValue(%d)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: %d"))
	}`, i, param.Index, param.Index)
		checks = append(checks, check)
	}

	return checks
}

// generateSourceCall 生成原始方法调用代码
func (g *WrapperGenerator) generateSourceCall(method MethodInfo, config *GeneratorConfig) (string, string) {
	var args []string

	for i, param := range method.Params {
		// 根据参数类型生成转换代码
		converter := g.analyzer.GetTypeConverter(param.Type)
		if converter == nil {
			if param.Type == "Context" {
				convertedArg := fmt.Sprintf("a%d.(%s)", i, "data.Context")
				args = append(args, convertedArg)
			} else {
				// 直接使用参数
				args = append(args, fmt.Sprintf("a%d", i))
			}
		} else {
			// 使用类型转换器
			convertedArg := fmt.Sprintf("a%d.(%s)%s", i, converter.DataType, converter.Converter)
			args = append(args, convertedArg)
		}
	}

	ret := "nil"
	switch method.ReturnType {
	case "string":
		if len(args) == 0 {
			return fmt.Sprintf("h.source.%s()", method.Name), "data.NewStringValue(\"\")"
		} else {
			return "", fmt.Sprintf("data.NewStringValue(h.source.%s(%s))", method.Name, strings.Join(args, ", "))
		}
	case "int":
		if len(args) == 0 {
			return fmt.Sprintf("h.source.%s()", method.Name), "data.NewIntValue(0)"
		} else {
			return "", fmt.Sprintf("data.NewIntValue(h.source.%s(%s))", method.Name, strings.Join(args, ", "))
		}
	case "bool":
		if len(args) == 0 {
			return fmt.Sprintf("h.source.%s()", method.Name), "data.NewBoolValue(false)"
		}
		return "", fmt.Sprintf("data.NewBoolValue(h.source.%s(%s))", method.Name, strings.Join(args, ", "))
	case "float":
		if len(args) == 0 {
			return fmt.Sprintf("h.source.%s()", method.Name), "data.NewFloatValue(0)"
		}
		return "", fmt.Sprintf("data.NewFloatValue(h.source.%s(%s))", method.Name, strings.Join(args, ", "))
	}

	// 构建方法调用
	if len(args) == 0 {
		return fmt.Sprintf("h.source.%s()", method.Name), ret
	} else {
		return fmt.Sprintf("h.source.%s(%s)", method.Name, strings.Join(args, ", ")), ret
	}
}
