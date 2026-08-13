package tools

import (
	"strings"
	"text/template"
)

// 类定义模板
const classTemplate = `package {{.PackageName}}

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

{{$StructName := .StructName}}
{{$ClassName := .ClassName}}

func New{{$StructName}}Class() data.ClassStmt {
	source := &{{$StructName}}{}
	return &{{$StructName}}Class{
{{- range .Properties}}
		{{.Name | lower}}: node.NewProperty(nil, "{{.Name}}", "public", true, data.New{{.Type}}Value({{.Value}}), nil),
{{- end}}
{{- range .Methods}}
		{{.Name | camel}}: &{{$StructName}}{{.Name | title}}Method{source},
{{- end}}
	}
}

type {{$StructName}}Class struct {
	node.Node
{{- range .Properties}}
	{{.Name | lower}} data.Property
{{- end}}
{{- range .Methods}}
	{{.Name | camel}} data.Method
{{- end}}
}

func (s *{{$StructName}}Class) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *{{$StructName}}Class) GetName() string {
	return "{{$ClassName}}"
}

func (s *{{$StructName}}Class) GetExtend() *string {
	return nil
}

func (s *{{$StructName}}Class) GetImplements() []string {
	return nil
}

func (s *{{$StructName}}Class) GetProperty(name string) (data.Property, bool) {
	switch name {
{{- range .Properties}}
	case "{{.Name}}":
		return s.{{.Name | lower}}, true
{{- end}}
	}
	return nil, false
}

func (s *{{$StructName}}Class) GetProperties() map[string]data.Property {
	return map[string]data.Property{
{{- range .Properties}}
		"{{.Name}}": s.{{.Name | lower}},
{{- end}}
	}
}

func (s *{{$StructName}}Class) GetMethod(name string) (data.Method, bool) {
	switch name {
{{- range .Methods}}
	case "{{.Name | camel}}":
		return s.{{.Name | camel}}, true
{{- end}}
	}
	return nil, false
}

func (s *{{$StructName}}Class) GetMethods() []data.Method {
	return []data.Method{
{{- range .Methods}}
		s.{{.Name | camel}},
{{- end}}
	}
}

func (s *{{$StructName}}Class) GetConstruct() data.Method {
	return nil
}
`

// 方法包装器模板
const methodTemplate = `package {{.PackageName}}

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

{{$StructName := .StructName}}
type {{$StructName}}{{.MethodName | title}}Method struct {
	source *{{$StructName}}
}

func (h *{{$StructName}}{{.MethodName | title}}Method) Call(ctx data.Context) (data.GetValue, data.Control) {
{{- range .ParamChecks}}
	{{.}}
{{- end}}

	{{.SourceCall}}
	return {{.Return}}, nil
}

func (h *{{$StructName}}{{.MethodName | title}}Method) GetName() string {
	return "{{.MethodName | lower}}"
}

func (h *{{$StructName}}{{.MethodName | title}}Method) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *{{$StructName}}{{.MethodName | title}}Method) GetIsStatic() bool {
	return false
}

func (h *{{$StructName}}{{.MethodName | title}}Method) GetParams() []data.GetValue {
	return []data.GetValue{
{{- range .Params}}
		node.NewParameter(nil, "{{.Name}}", {{.Index}}, nil, nil),
{{- end}}
	}
}

func (h *{{$StructName}}{{.MethodName | title}}Method) GetVariables() []data.Variable {
	return []data.Variable{
{{- range .Params}}
		node.NewVariable(nil, "{{.Name}}", {{.Index}}, nil),
{{- end}}
	}
}
`

// 创建模板函数映射
func createTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":   strings.ToLower,
		"title":   strings.Title,
		"upper":   strings.ToUpper,
		"replace": strings.ReplaceAll,
		"camel":   toCamelCase,
	}
}

// toCamelCase 转换为小写开头的驼峰命名
func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// 将第一个字符转为小写
	result := strings.ToLower(string(s[0]))

	// 处理剩余字符，保持驼峰命名
	for i := 1; i < len(s); i++ {
		if i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z' && s[i] >= 'A' && s[i] <= 'Z' {
			// 如果前一个字符是大写且当前字符也是大写，则当前字符转小写
			result += strings.ToLower(string(s[i]))
		} else {
			// 否则保持原样
			result += string(s[i])
		}
	}

	return result
}

// GetClassTemplate 获取类定义模板
func GetClassTemplate() *template.Template {
	funcMap := createTemplateFuncMap()
	return template.Must(template.New("class").Funcs(funcMap).Parse(classTemplate))
}

// GetMethodTemplate 获取方法包装器模板
func GetMethodTemplate() *template.Template {
	funcMap := createTemplateFuncMap()
	return template.Must(template.New("method").Funcs(funcMap).Parse(methodTemplate))
}
