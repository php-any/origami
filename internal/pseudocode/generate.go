package pseudocode

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/context"
	netannotation "github.com/php-any/origami/std/net/annotation"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/system"
)

const pseudoCodeExt = ".php"
const defaultOutputDir = ".zy/std"

// phpReservedTypeNames 在 PHP 中大小写不敏感，不能作为类名（如 List -> list）。
var phpReservedTypeNames = map[string]struct{}{
	"list": {}, "array": {}, "callable": {}, "self": {}, "parent": {}, "static": {},
	"false": {}, "true": {}, "null": {}, "void": {}, "mixed": {},
}

type GenericMethodDocRule struct {
	TemplateDoc string
	ReturnDoc   string
	ParamDocs   []string
}

type GenericClassDocRule struct {
	TemplateDoc        string
	ReturnDocByPHPType map[string]string
	MethodRules        map[string]GenericMethodDocRule
}

// genericDocRules 用于配置泛型类的 PHPDoc 生成规则，便于扩展更多泛型容器。
var genericDocRules = map[string]GenericClassDocRule{
	"Database\\DB": {
		TemplateDoc: "T",
		ReturnDocByPHPType: map[string]string{
			"DB":    "DB<T>",
			"array": "array<T>",
		},
		MethodRules: map[string]GenericMethodDocRule{
			"model": {
				TemplateDoc: "T",
				ParamDocs: []string{
					"@param class-string<T> $className",
					"@param string|null $connectionName",
				},
			},
			"first":      {ReturnDoc: "T|null"},
			"groupBy":    {ReturnDoc: "DB<T>"},
			"join":       {ReturnDoc: "DB<T>"},
			"limit":      {ReturnDoc: "DB<T>"},
			"offset":     {ReturnDoc: "DB<T>"},
			"orderBy":    {ReturnDoc: "DB<T>"},
			"select":     {ReturnDoc: "DB<T>"},
			"table":      {ReturnDoc: "DB<T>"},
			"connection": {ReturnDoc: "DB<T>"},
			"get":        {ReturnDoc: "array<T>"},
			"query":      {ReturnDoc: "array<T>"},
			"execute":    {ReturnDoc: "object"},
			"where": {
				ReturnDoc: "DB<T>",
				ParamDocs: []string{
					"@param string $sql",
					"@param mixed ...$args",
				},
			},
		},
	},
}

// repeatableAnnotationClasses 声明可重复使用的注解类。
// 键为完整类名（如 Net\Annotation\Middleware）或短类名（如 Middleware，仅 annotation 命名空间）。
var repeatableAnnotationClasses = map[string]struct{}{}

// annotationTargetByClass 为未在 GetImplements 中声明 TypeTarget* 的注解类提供 TARGET_* 回退配置。
var annotationTargetByClass = map[string][]string{
	"Net\\Annotation\\Middleware":          {"TARGET_CLASS"},
	"Net\\Annotation\\Controller":          {"TARGET_CLASS"},
	"Net\\Annotation\\Route":               {"TARGET_CLASS"},
	"Net\\Annotation\\Application":         {"TARGET_FUNCTION"},
	"Net\\Annotation\\GetMapping":          {"TARGET_METHOD"},
	"Net\\Annotation\\PostMapping":         {"TARGET_METHOD"},
	"Net\\Annotation\\PutMapping":          {"TARGET_METHOD"},
	"Net\\Annotation\\DeleteMapping":       {"TARGET_METHOD"},
	"Annotation\\Inject":                   {"TARGET_PARAMETER", "TARGET_PROPERTY"},
	"Database\\Annotation\\Table":          {"TARGET_CLASS"},
	"Database\\Annotation\\Column":         {"TARGET_PROPERTY"},
	"Database\\Annotation\\Id":             {"TARGET_PROPERTY"},
	"Database\\Annotation\\GeneratedValue": {"TARGET_PROPERTY"},
}

var targetMarkerToFlag = map[string]string{
	node.TypeTargetClass:     "TARGET_CLASS",
	node.TypeTargetMethod:    "TARGET_METHOD",
	node.TypeTargetProperty:  "TARGET_PROPERTY",
	node.TypeTargetFunction:  "TARGET_FUNCTION",
	node.TypeTargetParameter: "TARGET_PARAMETER",
}

// Generate 加载 Go 实现的标准库并通过反射生成 PHP 伪代码。
func Generate(outputDir string) error {
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	loadStdLibraries(vm)

	ctx := vm.CreateContext(nil)
	modules := buildModules(ctx, vm.AllFuncs(), vm.AllClasses())

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	indexContent := generatePseudoCodeIndex(modules)
	if err := os.WriteFile(outputDir+"/pseudo_README.md", []byte(indexContent), 0o644); err != nil {
		return fmt.Errorf("写入索引文件失败: %w", err)
	}

	for _, module := range modules {
		if module.ModuleName == "" {
			continue
		}

		content := generatePHPPseudoCode(module)

		var filepath string
		if module.Namespace != "" {
			dirPath := strings.ReplaceAll(module.Namespace, "\\", "/")
			fullDirPath := fmt.Sprintf("%s/%s", outputDir, dirPath)
			if err := os.MkdirAll(fullDirPath, 0o755); err != nil {
				return fmt.Errorf("创建目录 %s 失败: %w", fullDirPath, err)
			}
			filepath = fmt.Sprintf("%s/%s%s", fullDirPath, strings.ToLower(module.ModuleName), pseudoCodeExt)
		} else {
			filepath = fmt.Sprintf("%s/%s%s", outputDir, strings.ToLower(module.ModuleName), pseudoCodeExt)
		}

		if err := os.WriteFile(filepath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("写入 %s 失败: %w", filepath, err)
		}
		removeLegacyZyStub(filepath)
	}

	fmt.Println("标准库伪代码生成完成！")
	fmt.Printf("输出目录: %s\n", outputDir)
	fmt.Printf("共生成 %d 个模块\n", len(modules))
	return nil
}

func removeLegacyZyStub(phpPath string) {
	zyPath := strings.TrimSuffix(phpPath, pseudoCodeExt) + ".zy"
	_ = os.Remove(zyPath)
}

func loadStdLibraries(vm data.VM) {
	std.Load(vm)
	http.Load(vm)
	websocket.Load(vm)
	netannotation.Load(vm)
	system.Load(vm)
	context.Load(vm)
}

func buildModules(ctx data.Context, functions []data.FuncStmt, classes []data.ClassStmt) []PseudoCode {
	var modules []PseudoCode

	if len(functions) > 0 {
		funcModulesByNamespace := make(map[string]*PseudoCode)
		for _, fn := range functions {
			fullName := fn.GetName()
			namespace := ""
			shortName := fullName

			if strings.Contains(fullName, "\\") {
				parts := strings.Split(fullName, "\\")
				namespace = strings.Join(parts[:len(parts)-1], "\\")
				shortName = parts[len(parts)-1]
			}

			module, ok := funcModulesByNamespace[namespace]
			if !ok {
				module = &PseudoCode{
					ModuleName:   "functions",
					Description:  "标准库函数",
					Namespace:    namespace,
					IsAnnotation: isAnnotationNamespace(namespace),
				}
				funcModulesByNamespace[namespace] = module
			}

			sig := analyzeFunction(fn)
			sig.Name = shortName
			module.Functions = append(module.Functions, sig)
		}

		for _, m := range funcModulesByNamespace {
			modules = append(modules, *m)
		}
	}

	for _, class := range classes {
		classSig := analyzeClass(ctx, class)

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

		var module *PseudoCode
		for i := range modules {
			if modules[i].ModuleName == moduleName && modules[i].Namespace == namespace {
				module = &modules[i]
				break
			}
		}

		if module == nil {
			modules = append(modules, PseudoCode{
				ModuleName:   moduleName,
				Description:  fmt.Sprintf("%s 类", className),
				Namespace:    namespace,
				IsAnnotation: isAnnotationNamespace(namespace),
			})
			module = &modules[len(modules)-1]
		}

		module.Classes = append(module.Classes, classSig)
	}

	return modules
}

// PseudoCode 表示伪代码结构
type PseudoCode struct {
	ModuleName   string
	Description  string
	Namespace    string
	IsAnnotation bool
	Functions    []FunctionSignature
	Classes      []ClassSignature
}

// FunctionSignature 表示函数签名
type FunctionSignature struct {
	Name       string
	Params     []Parameter
	ReturnType string
	ReturnDoc  string
	ParamDocs  []string
	FakeReturn string // 虚假 return 语句，便于 IDE 分析
}

// ClassSignature 表示类签名
type ClassSignature struct {
	Name                 string
	ClassName            string
	Description          string
	TemplateDoc          string // 如 "T"，用于输出 @template T
	AttributeFlags       string // Attribute 构造参数，如 \Attribute::TARGET_CLASS | \Attribute::IS_REPEATABLE
	Methods              []MethodSignature
	AnnotationCtorParams []Parameter
	Properties           []PropertySignature
}

// MethodSignature 表示方法签名
type MethodSignature struct {
	Name        string
	Params      []Parameter
	ReturnType  string
	ReturnDoc   string
	TemplateDoc string // 如 "T"，用于输出 @template T（常用于静态工厂）
	ParamDocs   []string
	FakeReturn  string
	Modifier    string
	IsStatic    bool
}

// PropertySignature 表示属性签名
type PropertySignature struct {
	Name     string
	Type     string
	Modifier string
	IsStatic bool
	Default  string
}

// Parameter 表示参数
type Parameter struct {
	Name       string
	Type       string
	IsVariadic bool
	Default    string // 含前导空格，如 ` = "App"`
}

func resolveClassStmt(class data.ClassStmt, ctx data.Context) data.ClassStmt {
	if cg, ok := class.(data.ClassGeneric); ok {
		if cloned, ok := cg.Clone(nil).(data.ClassStmt); ok {
			class = cloned
		}
	}

	val, ctl := class.GetValue(ctx)
	if ctl != nil {
		return class
	}
	if cv, ok := val.(*data.ClassValue); ok && cv.Class != nil {
		return cv.Class
	}
	if stmt, ok := val.(data.ClassStmt); ok {
		return stmt
	}
	return class
}

// collectClassMethods 收集类方法，并合并 GetConstruct / GetStaticMethod 暴露的方法。
func collectClassMethods(class data.ClassStmt) []data.Method {
	seen := make(map[string]bool)
	var collected []data.Method
	add := func(m data.Method) {
		if m == nil {
			return
		}
		name := m.GetName()
		if name == "" || seen[name] {
			return
		}
		seen[name] = true
		collected = append(collected, m)
	}

	for _, m := range class.GetMethods() {
		add(m)
	}
	add(class.GetConstruct())

	if gsm, ok := class.(data.GetStaticMethod); ok {
		for _, name := range probeStaticMethodNames(class) {
			if m, ok := gsm.GetStaticMethod(name); ok {
				add(m)
			}
		}
	}

	sort.Slice(collected, func(i, j int) bool {
		if collected[i].GetName() == "__construct" {
			return true
		}
		if collected[j].GetName() == "__construct" {
			return false
		}
		return collected[i].GetName() < collected[j].GetName()
	})

	return collected
}

// probeStaticMethodNames 探测可能存在的静态方法名（GetMethods 未包含时）。
func probeStaticMethodNames(class data.ClassStmt) []string {
	candidates := []string{
		"bind", "path", "debug", "error", "info", "warn", "notice", "trace", "fatal",
		"createFromFormat", "getLastErrors",
	}
	var names []string
	for _, name := range candidates {
		if _, ok := class.GetMethod(name); ok {
			continue
		}
		names = append(names, name)
	}
	return names
}

func analyzeClass(ctx data.Context, class data.ClassStmt) ClassSignature {
	class = resolveClassStmt(class, ctx)
	className := class.GetName()

	var shortClassName string
	if strings.Contains(className, "\\") {
		parts := strings.Split(className, "\\")
		shortClassName = parts[len(parts)-1]
	} else {
		shortClassName = className
	}

	sig := ClassSignature{
		Name:        className,
		ClassName:   safePHPClassName(shortClassName),
		Description: fmt.Sprintf("%s 类", shortClassName),
	}

	namespace := classNamespace(className)
	forAnnotation := isAnnotationNamespace(className)
	if forAnnotation {
		sig.AttributeFlags = buildAttributeFlags(class, className, namespace, shortClassName)
	}
	classRule, hasClassRule := genericDocRuleForClass(namespace, shortClassName)

	// 泛型容器类：从配置填充类级 template 文档。
	if hasClassRule {
		sig.TemplateDoc = classRule.TemplateDoc
	}
	if ctor := class.GetConstruct(); ctor != nil {
		sig.AnnotationCtorParams = analyzeMethodParams(ctor, forAnnotation, namespace)
	}

	for _, method := range collectClassMethods(class) {
		if method == nil {
			continue
		}

		methodSig := MethodSignature{
			Name:     method.GetName(),
			Modifier: "public",
			IsStatic: method.GetIsStatic(),
		}

		methodSig.Params = analyzeMethodParams(method, false, namespace)

		rawReturnType := methodReturnType(method)
		methodSig.ReturnType = formatPHPReturnType(rawReturnType, namespace, shortClassName)
		methodSig.ReturnDoc = formatPHPReturnDoc(rawReturnType, methodSig.ReturnType, namespace, shortClassName)
		methodSig.FakeReturn = fakeReturnStatement(methodSig.ReturnType, "        ")

		// 泛型文档规则：优先应用类级返回类型映射，再应用方法级覆盖。
		if hasClassRule {
			if mapped, ok := classRule.ReturnDocByPHPType[methodSig.ReturnType]; ok {
				methodSig.ReturnDoc = mapped
			}
			if mr, ok := classRule.MethodRules[methodSig.Name]; ok {
				if mr.TemplateDoc != "" {
					methodSig.TemplateDoc = mr.TemplateDoc
				}
				if mr.ReturnDoc != "" {
					methodSig.ReturnDoc = mr.ReturnDoc
				}
				if len(mr.ParamDocs) > 0 {
					methodSig.ParamDocs = mr.ParamDocs
				}
			}
		}

		if methodSig.Name == "__construct" {
			if len(sig.AnnotationCtorParams) == 0 {
				sig.AnnotationCtorParams = append([]Parameter(nil), methodSig.Params...)
			}
			if forAnnotation {
				continue
			}
		}

		sig.Methods = append(sig.Methods, methodSig)
	}

	return sig
}

func safePHPClassName(name string) string {
	if _, reserved := phpReservedTypeNames[strings.ToLower(name)]; reserved {
		return name + "Stub"
	}
	return name
}

func classNamespace(className string) string {
	if !strings.Contains(className, "\\") {
		return ""
	}
	parts := strings.Split(className, "\\")
	return strings.Join(parts[:len(parts)-1], "\\")
}

func genericDocRuleForClass(namespace, shortClassName string) (GenericClassDocRule, bool) {
	key := shortClassName
	if namespace != "" {
		key = namespace + `\` + shortClassName
	}
	rule, ok := genericDocRules[key]
	return rule, ok
}

func methodReturnType(method data.Method) data.Types {
	if returnTypeInterface, ok := method.(data.GetReturnType); ok {
		return returnTypeInterface.GetReturnType()
	}
	return nil
}

type methodParamsSource interface {
	GetParams() []data.GetValue
	GetVariables() []data.Variable
}

func analyzeMethodParams(method methodParamsSource, forAnnotation bool, namespace string) []Parameter {
	params := method.GetParams()
	vars := method.GetVariables()

	varByIndex := make(map[int]data.Variable)
	varByName := make(map[string]data.Variable)
	for _, v := range vars {
		varByIndex[v.GetIndex()] = v
		varByName[v.GetName()] = v
	}

	var result []Parameter
	for i, param := range params {
		if param == nil {
			continue
		}
		p := analyzeParam(param, forAnnotation, len(result), namespace)
		if p == nil {
			continue
		}
		if v, ok := varByIndex[i]; ok {
			enrichParamFromVariable(p, v, namespace)
		} else if v, ok := varByName[p.Name]; ok {
			enrichParamFromVariable(p, v, namespace)
		}
		result = append(result, *p)
	}
	return result
}

func enrichParamFromVariable(p *Parameter, v data.Variable, namespace string) {
	if ty := v.GetType(); ty != nil {
		p.Type = formatPHPType(ty, namespace)
	}
	applyNullableDefault(p)
}

func applyNullableDefault(p *Parameter) {
	if p.Default == " = null" && p.Type != "mixed" && !strings.HasPrefix(p.Type, "?") {
		p.Type = "?" + p.Type
	}
}

func analyzeParams(params []data.GetValue, forAnnotation bool) []Parameter {
	var result []Parameter
	for _, param := range params {
		if param == nil {
			continue
		}
		if p := analyzeParam(param, forAnnotation, len(result), ""); p != nil {
			result = append(result, *p)
		}
	}
	return result
}

func analyzeParam(param data.GetValue, forAnnotation bool, index int, namespace string) *Parameter {
	paramName := ""
	var paramTypes data.Types
	isVariadic := false
	var defaultValue data.GetValue

	switch p := param.(type) {
	case *node.AnnotationTargetParameter:
		paramName = p.GetName()
		paramTypes = p.GetType()
		defaultValue = p.GetDefaultValue()
	case *node.ParameterRawAST:
		paramName = p.GetName()
		paramTypes = p.GetType()
		defaultValue = p.GetDefaultValue()
	case *node.Parameter:
		paramName = p.GetName()
		paramTypes = p.GetType()
		defaultValue = p.GetDefaultValue()
	case *node.PromotedParameter:
		paramName = p.GetName()
		paramTypes = p.GetType()
		defaultValue = p.GetDefaultValue()
	case data.Parameter:
		paramName = p.GetName()
		paramTypes = p.GetType()
		defaultValue = p.GetDefaultValue()
	case data.Variable:
		paramName = p.GetName()
		paramTypes = p.GetType()
	default:
		if _, ok := param.(*node.Parameters); ok {
			isVariadic = true
			paramName = fmt.Sprintf("param%d", index)
		} else {
			return nil
		}
	}

	if paramName == "" {
		paramName = fmt.Sprintf("param%d", index)
	}

	if _, ok := param.(*node.Parameters); ok {
		isVariadic = true
	}
	if isVariadicParameter(param) {
		isVariadic = true
	}

	p := &Parameter{
		Name:       paramName,
		Type:       formatPHPType(paramTypes, namespace),
		IsVariadic: isVariadic,
		Default:    formatDefaultValue(defaultValue),
	}
	applyNullableDefault(p)
	return p
}

func isVariadicParameter(param data.GetValue) bool {
	if _, ok := param.(*node.Parameters); ok {
		return true
	}
	t := reflect.TypeOf(param)
	if t != nil && t.Kind() == reflect.Ptr && t.Elem().Name() == "ParametersTODO" {
		return true
	}
	return false
}

func formatPHPTypeFromString(s, namespace string) string {
	if s == "" || s == "mixed" {
		return "mixed"
	}
	return formatPHPType(&stubType{s: s}, namespace)
}

type stubType struct{ s string }

func (t *stubType) Is(data.Value) bool { return true }
func (t *stubType) String() string     { return t.s }

func formatPHPType(ty data.Types, namespace string) string {
	if ty == nil {
		return "mixed"
	}
	switch t := ty.(type) {
	case data.UnionType:
		parts := make([]string, len(t.Types))
		for i, pt := range t.Types {
			parts[i] = formatPHPType(pt, namespace)
		}
		return strings.Join(parts, "|")
	case data.MultipleReturnType:
		// 部分 Go 绑定误用多返回值类型表示联合类型，生成时用 | 连接。
		parts := make([]string, len(t.Types))
		for i, pt := range t.Types {
			parts[i] = formatPHPType(pt, namespace)
		}
		return strings.Join(parts, "|")
	}
	raw := strings.TrimSpace(ty.String())
	if raw == "" || raw == "void" || raw == "LspTypes" {
		return "mixed"
	}
	if raw == "<T>" || strings.HasPrefix(raw, "<") || strings.HasSuffix(raw, ">") {
		return "mixed"
	}
	if strings.Contains(raw, "|") {
		parts := strings.Split(raw, "|")
		for i, part := range parts {
			parts[i] = formatPHPTypeFromString(strings.TrimSpace(part), namespace)
		}
		return strings.Join(parts, "|")
	}
	if strings.Contains(raw, "\\") {
		if namespace != "" && strings.HasPrefix(raw, namespace+"\\") {
			short := strings.TrimPrefix(raw, namespace+"\\")
			if short != "" && isValidPHPIdentifier(short) {
				return short
			}
		}
		return "\\" + raw
	}
	if strings.HasPrefix(raw, "?") {
		inner := formatPHPTypeFromString(raw[1:], namespace)
		if inner == "mixed" {
			return "mixed"
		}
		return "?" + inner
	}
	return raw
}

func isValidPHPIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' {
				continue
			}
			return false
		}
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return false
	}
	return true
}

func formatPHPReturnType(ty data.Types, namespace, shortClassName string) string {
	if ty == nil {
		return ""
	}
	if g, ok := ty.(data.Generic); ok && g.Name == "" {
		if namespace == "Database" && shortClassName == "DB" {
			return "DB"
		}
		return ""
	}
	formatted := formatPHPType(ty, namespace)
	if formatted == "" || formatted == "mixed" {
		if g, ok := ty.(data.Generic); ok && (g.Name == "" || g.Name == "M") {
			if namespace == "Database" {
				return "DB"
			}
			return ""
		}
	}
	if formatted == "mixed" {
		return ""
	}
	return formatted
}

func formatPHPReturnDoc(ty data.Types, phpReturnType, namespace, _ string) string {
	if phpReturnType == "" || ty == nil {
		return ""
	}

	raw := strings.TrimSpace(ty.String())
	if raw == "" || raw == "void" {
		return ""
	}

	// 仅当反射类型含泛型信息时输出 @return，普通类型跳过。
	if strings.Contains(raw, "<") && strings.Contains(raw, ">") {
		return normalizePHPDocType(raw, namespace)
	}

	return ""
}

func normalizePHPDocType(raw, namespace string) string {
	docType := strings.TrimSpace(raw)
	if docType == "" {
		return ""
	}
	if namespace != "" {
		docType = strings.ReplaceAll(docType, namespace+"\\", "")
	}
	return docType
}

func formatDefaultValue(v data.GetValue) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case *data.StringValue:
		return " = " + strconv.Quote(val.AsString())
	case *data.IntValue:
		return fmt.Sprintf(" = %d", val.Value)
	case *data.FloatValue:
		return fmt.Sprintf(" = %g", val.Value)
	case *data.BoolValue:
		if val.Value {
			return " = true"
		}
		return " = false"
	case *data.NullValue:
		return " = null"
	default:
		return ""
	}
}

func isAnnotationNamespace(namespace string) bool {
	if namespace == "" {
		return false
	}
	return strings.Contains(strings.ToLower(strings.ReplaceAll(namespace, "\\", "/")), "annotation")
}

func isRepeatableAnnotation(class data.ClassStmt, className, namespace, shortClassName string) bool {
	if _, ok := repeatableAnnotationClasses[className]; ok {
		return true
	}
	if namespace != "" {
		if _, ok := repeatableAnnotationClasses[namespace+`\`+shortClassName]; ok {
			return true
		}
	}
	if _, ok := repeatableAnnotationClasses[shortClassName]; ok {
		return true
	}
	for _, impl := range class.GetImplements() {
		if impl == node.TypeRepeatable {
			return true
		}
	}
	return false
}

func resolveAnnotationTargetFlags(class data.ClassStmt, className, namespace, shortClassName string) []string {
	var flags []string
	seen := make(map[string]bool)
	add := func(flag string) {
		if flag == "" || seen[flag] {
			return
		}
		seen[flag] = true
		flags = append(flags, flag)
	}

	for _, impl := range class.GetImplements() {
		if flag, ok := targetMarkerToFlag[impl]; ok {
			add(flag)
		}
	}

	if len(flags) == 0 {
		keys := []string{className}
		if namespace != "" {
			keys = append(keys, namespace+`\`+shortClassName)
		}
		keys = append(keys, shortClassName)
		for _, key := range keys {
			if targets, ok := annotationTargetByClass[key]; ok {
				for _, target := range targets {
					add(target)
				}
				break
			}
		}
	}

	return flags
}

func buildAttributeFlags(class data.ClassStmt, className, namespace, shortClassName string) string {
	var parts []string
	for _, flag := range resolveAnnotationTargetFlags(class, className, namespace, shortClassName) {
		parts = append(parts, `\Attribute::`+flag)
	}
	if isRepeatableAnnotation(class, className, namespace, shortClassName) {
		parts = append(parts, `\Attribute::IS_REPEATABLE`)
	}
	return strings.Join(parts, " | ")
}

func analyzeFunction(fn data.FuncStmt) FunctionSignature {
	sig := FunctionSignature{
		Name: fn.GetName(),
	}

	sig.Params = analyzeMethodParams(fn, false, "")

	if returnTypeInterface, ok := fn.(data.GetReturnType); ok {
		rawReturnType := returnTypeInterface.GetReturnType()
		sig.ReturnType = formatPHPReturnType(rawReturnType, "", "")
		sig.ReturnDoc = formatPHPReturnDoc(rawReturnType, sig.ReturnType, "", "")
		sig.FakeReturn = fakeReturnStatement(sig.ReturnType, "    ")
	}

	return sig
}

// fakeReturnStatement 根据返回类型生成占位 return，使伪代码通过静态检查。
func fakeReturnStatement(returnType, indent string) string {
	if indent == "" {
		indent = "    "
	}
	rt := strings.TrimSpace(returnType)
	if rt == "" {
		return ""
	}

	ret := func(expr string) string {
		return indent + "return " + expr + ";\n"
	}

	if strings.HasPrefix(rt, "?") {
		return ret("null")
	}

	if strings.Contains(rt, "|") {
		parts := strings.Split(rt, "|")
		for _, p := range parts {
			if strings.TrimSpace(p) == "null" {
				return ret("null")
			}
		}
		return fakeReturnStatement(strings.TrimSpace(parts[0]), indent)
	}

	switch rt {
	case "void", "never":
		return ""
	case "string":
		return ret("''")
	case "int":
		return ret("0")
	case "float":
		return ret("0.0")
	case "bool":
		return ret("false")
	case "array":
		return ret("[]")
	case "object":
		return ret("new \\stdClass()")
	case "callable":
		return indent + "return static function () {};\n"
	case "static":
		return ret("new static()")
	case "self":
		return ret("new self()")
	case "mixed":
		return ret("null")
	default:
		if isValidPHPIdentifier(rt) || strings.HasPrefix(rt, "\\") {
			return ret("new " + rt + "()")
		}
		return ret("null")
	}
}

func generatePHPPseudoCode(module PseudoCode) string {
	tmpl := `<?php

{{if .Namespace}}namespace {{.Namespace}};

{{else}}namespace {

{{end}}{{if .Functions}}
{{range .Functions}}
{{if or .ReturnDoc .ParamDocs}}
/**
{{if .ParamDocs}}{{range .ParamDocs}}
 * {{.}}{{end}}{{end}}{{if .ReturnDoc}}
 * @return {{.ReturnDoc}}{{end}}
 */
{{end}}function {{.Name}}({{range $i, $param := .Params}}{{if $i}}, {{end}}{{if eq $param.IsVariadic true}}...${{$param.Name}}{{else}}{{if ne $param.Type "mixed"}}{{$param.Type}} {{end}}${{$param.Name}}{{$param.Default}}{{end}}{{end}}){{if .ReturnType}} : {{.ReturnType}}{{end}} {
{{if .FakeReturn}}{{.FakeReturn}}{{else}}    // 实现逻辑
{{end}}}
{{end}}
{{end}}
{{if .Classes}}
{{range .Classes}}
{{if .TemplateDoc}}
/**
 * @template {{.TemplateDoc}}
 */
{{end}}
{{if $.IsAnnotation}}{{if .AttributeFlags}}#[\Attribute({{.AttributeFlags}})]
{{else}}#[\Attribute]
{{end}}class {{.ClassName}} {
    public function __construct({{range $i, $param := .AnnotationCtorParams}}{{if $i}}, {{end}}{{if eq $param.IsVariadic true}}...${{$param.Name}}{{else}}{{if ne $param.Type "mixed"}}{{$param.Type}} {{end}}${{$param.Name}}{{$param.Default}}{{end}}{{end}}) {}
{{if .Methods}}{{range .Methods}}
{{if or .TemplateDoc .ReturnDoc .ParamDocs}}    /**
{{if .TemplateDoc}}     * @template {{.TemplateDoc}}
{{end}}{{if .ParamDocs}}{{range .ParamDocs}}     * {{.}}
{{end}}{{end}}{{if .ReturnDoc}}     * @return {{.ReturnDoc}}
{{end}}     */
{{end}}    {{.Modifier}} {{if .IsStatic}}static {{end}}function {{.Name}}({{range $i, $param := .Params}}{{if $i}}, {{end}}{{if eq $param.IsVariadic true}}...${{$param.Name}}{{else}}{{if ne $param.Type "mixed"}}{{$param.Type}} {{end}}${{$param.Name}}{{$param.Default}}{{end}}{{end}}){{if .ReturnType}} : {{.ReturnType}}{{end}} {
{{if .FakeReturn}}{{.FakeReturn}}{{else}}        // 实现逻辑
{{end}}    }
{{end}}{{end}}
}
{{else}}class {{.ClassName}} {
{{if .Properties}}{{range .Properties}}
    {{.Modifier}} {{if .IsStatic}}static {{end}}${{.Name}}{{if .Type}} : {{.Type}}{{end}}{{if .Default}} = {{.Default}}{{end}};
{{end}}{{end}}{{if .Methods}}{{range .Methods}}
{{if or .TemplateDoc .ReturnDoc .ParamDocs}}    /**
{{if .TemplateDoc}}     * @template {{.TemplateDoc}}
{{end}}{{if .ParamDocs}}{{range .ParamDocs}}     * {{.}}
{{end}}{{end}}{{if .ReturnDoc}}     * @return {{.ReturnDoc}}
{{end}}     */
{{end}}    {{.Modifier}} {{if .IsStatic}}static {{end}}function {{.Name}}({{range $i, $param := .Params}}{{if $i}}, {{end}}{{if eq $param.IsVariadic true}}...${{$param.Name}}{{else}}{{if ne $param.Type "mixed"}}{{$param.Type}} {{end}}${{$param.Name}}{{$param.Default}}{{end}}{{end}}){{if .ReturnType}} : {{.ReturnType}}{{end}} {
{{if .FakeReturn}}{{.FakeReturn}}{{else}}        // 实现逻辑
{{end}}    }
{{end}}{{end}}
}
{{end}}
{{end}}
{{end}}{{if not .Namespace}}
}
{{end}}`

	t, err := template.New("php_pseudocode").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, module); err != nil {
		panic(err)
	}

	return buf.String()
}

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
	if err := t.Execute(&buf, modules); err != nil {
		panic(err)
	}

	return buf.String()
}
