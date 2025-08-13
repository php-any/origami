package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

// 处理补全请求
func handleTextDocumentCompletion(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/completion", true, req.Params)

	var params CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal completion params: %v", err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	logger.Info("请求代码补全：%s 位置 %d:%d", uri, position.Line, position.Character)

	doc, exists := documents[uri]
	if !exists {
		return CompletionList{IsIncomplete: false, Items: []CompletionItem{}}, nil
	}

	// 获取补全项
	items := getCompletionItems(doc.Content, position)

	result := CompletionList{
		IsIncomplete: false,
		Items:        items,
	}

	logLSPResponse("textDocument/completion", result, nil)
	return result, nil
}

// 获取补全项
func getCompletionItems(content string, position Position) []CompletionItem {
	items := []CompletionItem{}

	// 获取当前位置的上下文
	context := getCompletionContext(content, position)

	// 根据上下文提供相应的补全项
	switch context {
	case "keyword":
		items = append(items, getKeywordCompletions()...)
	case "type":
		items = append(items, getTypeCompletions()...)
	case "function":
		items = append(items, getFunctionCompletions()...)
	case "control":
		items = append(items, getControlStructureCompletions()...)
	case "operator":
		items = append(items, getOperatorCompletions()...)
	case "annotation":
		items = append(items, getAnnotationCompletions()...)
	case "builtin":
		items = append(items, getBuiltinCompletions()...)
	case "snippet":
		items = append(items, getSnippetCompletions()...)
	case "object_method":
		items = append(items, getObjectMethodCompletions()...)
		logger.Info("对象方法补全：%#v", items)
	default:
		// 默认提供所有补全项，但避免重复
		items = append(items, getSnippetCompletions()...)
		items = append(items, getFilteredKeywordCompletions()...)
		items = append(items, getTypeCompletions()...)
		items = append(items, getFunctionCompletions()...)
		items = append(items, getFilteredControlStructureCompletions()...)
		items = append(items, getOperatorCompletions()...)
		items = append(items, getAnnotationCompletions()...)
		items = append(items, getBuiltinCompletions()...)
	}

	return items
}

// 获取补全上下文
func getCompletionContext(content string, position Position) string {
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return "default"
	}

	line := lines[position.Line]
	if int(position.Character) > len(line) {
		return "default"
	}

	// 分析当前位置的上下文
	beforeCursor := line[:position.Character]

	// 检查是否在注解中
	if strings.Contains(beforeCursor, "@") {
		return "annotation"
	}

	// 检查是否在类型声明中
	if strings.Contains(beforeCursor, ":") && strings.Contains(beforeCursor, "$") {
		return "type"
	}

	// 检查是否在对象方法调用中 (如 $str-> 或 $str->le)
	if strings.Contains(beforeCursor, "->") {
		return "object_method"
	}

	// 检查是否在函数调用中
	if strings.Contains(beforeCursor, "(") {
		return "function"
	}

	// 检查是否在控制结构中
	if strings.Contains(beforeCursor, "if") || strings.Contains(beforeCursor, "for") ||
		strings.Contains(beforeCursor, "while") || strings.Contains(beforeCursor, "switch") {
		return "control"
	}

	// 检查是否在运算符附近
	if strings.Contains(beforeCursor, "+") || strings.Contains(beforeCursor, "-") ||
		strings.Contains(beforeCursor, "*") || strings.Contains(beforeCursor, "/") {
		return "operator"
	}

	// 检查是否在行首或独立位置，适合代码片段
	if strings.TrimSpace(beforeCursor) == "" ||
		strings.HasSuffix(strings.TrimSpace(beforeCursor), ";") ||
		strings.HasSuffix(strings.TrimSpace(beforeCursor), "}") ||
		strings.HasSuffix(strings.TrimSpace(beforeCursor), "\n") {
		return "snippet"
	}

	return "default"
}

// 代码片段补全 - 提供完整的代码结构
func getSnippetCompletions() []CompletionItem {
	snippets := []struct {
		label      string
		insertText string
		detail     string
		kind       CompletionItemKind
	}{
		// 控制结构片段
		{
			"for",
			"for (int $i = 0; $i < ${1:count}; $i++) {\n\t${2:// 循环体}\n}",
			"for 循环结构",
			CompletionItemKindSnippet,
		},
		{
			"fori",
			"for (int $i = 0; $i < ${1:count}; $i++) {\n\t${2:// 循环体}\n}",
			"for 循环结构 (带索引)",
			CompletionItemKindSnippet,
		},
		{
			"foreach",
			"foreach ($${1:array} as $${2:item}) {\n\t${3:// 循环体}\n}",
			"foreach 循环结构",
			CompletionItemKindSnippet,
		},
		{
			"foreachk",
			"foreach ($${1:array} as $${2:key} => $${3:value}) {\n\t${4:// 循环体}\n}",
			"foreach 循环结构 (带键值)",
			CompletionItemKindSnippet,
		},
		{
			"while",
			"while ($${1:condition}) {\n\t${2:// 循环体}\n}",
			"while 循环结构",
			CompletionItemKindSnippet,
		},
		{
			"do",
			"do {\n\t${1:// 循环体}\n} while ($${2:condition});",
			"do-while 循环结构",
			CompletionItemKindSnippet,
		},
		{
			"if",
			"if ($${1:condition}) {\n\t${2:// 条件为真时执行}\n}",
			"if 条件语句",
			CompletionItemKindSnippet,
		},
		{
			"ifelse",
			"if ($${1:condition}) {\n\t${2:// 条件为真时执行}\n} else {\n\t${3:// 条件为假时执行}\n}",
			"if-else 条件语句",
			CompletionItemKindSnippet,
		},
		{
			"ifelseif",
			"if ($${1:condition1}) {\n\t${2:// 条件1为真时执行}\n} elseif ($${3:condition2}) {\n\t${4:// 条件2为真时执行}\n} else {\n\t${5:// 所有条件为假时执行}\n}",
			"if-elseif-else 条件语句",
			CompletionItemKindSnippet,
		},
		{
			"switch",
			"switch ($${1:expression}) {\n\tcase ${2:value1}:\n\t\t${3:// 处理逻辑}\n\t\tbreak;\n\tcase ${4:value2}:\n\t\t${5:// 处理逻辑}\n\t\tbreak;\n\tdefault:\n\t\t${6:// 默认处理}\n\t\tbreak;\n}",
			"switch 分支语句",
			CompletionItemKindSnippet,
		},
		{
			"try",
			"try {\n\t${1:// 可能抛出异常的代码}\n} catch (Exception $${2:e}) {\n\t${3:// 异常处理}\n} finally {\n\t${4:// 总是执行的代码}\n}",
			"try-catch-finally 异常处理",
			CompletionItemKindSnippet,
		},

		// 函数和类片段
		{
			"function",
			"function ${1:functionName}($${2:parameters}) {\n\t${3:// 函数体}\n\treturn ${4:value};\n}",
			"函数定义",
			CompletionItemKindSnippet,
		},
		{
			"func",
			"function ${1:functionName}($${2:parameters}) {\n\t${3:// 函数体}\n\treturn ${4:value};\n}",
			"函数定义 (简写)",
			CompletionItemKindSnippet,
		},
		{
			"class",
			"class ${1:ClassName} {\n\t${2:// 类属性}\n\t\n\tpublic function __construct() {\n\t\t${3:// 构造函数}\n\t}\n\t\n\t${4:// 类方法}\n}",
			"类定义",
			CompletionItemKindSnippet,
		},
		{
			"interface",
			"interface ${1:InterfaceName} {\n\t${2:// 接口方法声明}\n}",
			"接口定义",
			CompletionItemKindSnippet,
		},
		{
			"trait",
			"trait ${1:TraitName} {\n\t${2:// trait 方法}\n}",
			"Trait 定义",
			CompletionItemKindSnippet,
		},
		{
			"constructor",
			"public function __construct($${1:parameters}) {\n\t${2:// 构造函数逻辑}\n}",
			"构造函数",
			CompletionItemKindSnippet,
		},
		{
			"destructor",
			"public function __destruct() {\n\t${1:// 析构函数逻辑}\n}",
			"析构函数",
			CompletionItemKindSnippet,
		},
		{
			"getter",
			"public function get${1:PropertyName}() {\n\treturn $this->${2:propertyName};\n}",
			"Getter 方法",
			CompletionItemKindSnippet,
		},
		{
			"setter",
			"public function set${1:PropertyName}($${2:value}) {\n\t$this->${3:propertyName} = $${2:value};\n}",
			"Setter 方法",
			CompletionItemKindSnippet,
		},

		// 注解片段
		{
			"@Controller",
			"@Controller(name: \"${1:ControllerName}\")",
			"Controller 注解",
			CompletionItemKindSnippet,
		},
		{
			"@Route",
			"@Route(prefix: \"${1:/api}\")",
			"Route 注解",
			CompletionItemKindSnippet,
		},
		{
			"@GetMapping",
			"@GetMapping(path: \"${1:/list}\")",
			"GetMapping 注解",
			CompletionItemKindSnippet,
		},
		{
			"@PostMapping",
			"@PostMapping(path: \"${1:/create}\")",
			"PostMapping 注解",
			CompletionItemKindSnippet,
		},
		{
			"@Autowired",
			"@Autowired\nprivate $${1:serviceName};",
			"Autowired 注解",
			CompletionItemKindSnippet,
		},
		{
			"@Inject",
			"@Inject(service: \"${1:ServiceName}\")\nprivate $${2:serviceName};",
			"Inject 注解",
			CompletionItemKindSnippet,
		},

		// 变量和类型声明片段
		{
			"var",
			"$${1:variableName} = ${2:value};",
			"变量声明",
			CompletionItemKindSnippet,
		},
		{
			"string",
			"string $${1:variableName} = \"${2:value}\";",
			"字符串变量声明",
			CompletionItemKindSnippet,
		},
		{
			"int",
			"int $${1:variableName} = ${2:0};",
			"整数变量声明",
			CompletionItemKindSnippet,
		},
		{
			"float",
			"float $${1:variableName} = ${2:0.0};",
			"浮点数变量声明",
			CompletionItemKindSnippet,
		},
		{
			"bool",
			"bool $${1:variableName} = ${2:true};",
			"布尔变量声明",
			CompletionItemKindSnippet,
		},
		{
			"array",
			"array $${1:variableName} = [${2:item1}, ${3:item2}];",
			"数组变量声明",
			CompletionItemKindSnippet,
		},
		{
			"const",
			"const ${1:CONSTANT_NAME} = ${2:value};",
			"常量声明",
			CompletionItemKindSnippet,
		},

		// 输出和调试片段
		{
			"echo",
			"echo $${1:variable};",
			"echo 输出",
			CompletionItemKindSnippet,
		},
		{
			"echo",
			"echo \"${1:message}\";",
			"echo 字符串输出",
			CompletionItemKindSnippet,
		},
		{
			"dump",
			"dump($${1:variable});",
			"dump 调试输出",
			CompletionItemKindSnippet,
		},
		{
			"log",
			"Log::${1:info}(\"${2:message}\");",
			"日志输出",
			CompletionItemKindSnippet,
		},

		// 命名空间和引用片段
		{
			"namespace",
			"namespace ${1:namespace\\path};",
			"命名空间声明",
			CompletionItemKindSnippet,
		},
		{
			"use",
			"use ${1:Namespace\\ClassName};",
			"use 语句",
			CompletionItemKindSnippet,
		},
		{
			"include",
			"include \"${1:file.php}\";",
			"include 文件",
			CompletionItemKindSnippet,
		},
		{
			"require",
			"require \"${1:file.php}\";",
			"require 文件",
			CompletionItemKindSnippet,
		},

		// 对象和数组操作片段
		{
			"new",
			"$${1:object} = new ${2:ClassName}(${3:parameters});",
			"创建对象",
			CompletionItemKindSnippet,
		},
		{
			"array_push",
			"$${1:array}->push($${2:item});",
			"数组添加元素",
			CompletionItemKindSnippet,
		},
		{
			"array_pop",
			"$${1:item} = $${2:array}->pop();",
			"数组移除末尾元素",
			CompletionItemKindSnippet,
		},
		{
			"array_map",
			"$${1:result} = $${2:array}->map(function($${3:item}) {\n\treturn ${4:// 处理逻辑};\n});",
			"数组映射",
			CompletionItemKindSnippet,
		},
		{
			"array_filter",
			"$${1:result} = $${2:array}->filter(function($${3:item}) {\n\treturn ${4:// 过滤条件};\n});",
			"数组过滤",
			CompletionItemKindSnippet,
		},
		{
			"array_reduce",
			"$${1:result} = $${2:array}->reduce(function($${3:acc}, $${4:item}) {\n\treturn ${5:// 归约逻辑};\n}, $${6:initial});",
			"数组归约",
			CompletionItemKindSnippet,
		},

		// 字符串操作片段
		{
			"str_length",
			"$${1:length} = $${2:string}->length;",
			"获取字符串长度",
			CompletionItemKindSnippet,
		},
		{
			"str_substring",
			"$${1:substring} = $${2:string}->substring($${3:start}, $${4:end});",
			"字符串子串",
			CompletionItemKindSnippet,
		},
		{
			"str_indexOf",
			"$${1:index} = $${2:string}->indexOf(\"${3:search}\");",
			"字符串查找",
			CompletionItemKindSnippet,
		},
		{
			"str_replace",
			"$${1:result} = $${2:string}->replace(\"${3:search}\", \"${4:replace}\");",
			"字符串替换",
			CompletionItemKindSnippet,
		},

		// 异常处理片段
		{
			"throw",
			"throw new Exception(\"${1:error message}\");",
			"抛出异常",
			CompletionItemKindSnippet,
		},
		{
			"catch",
			"catch (Exception $${1:e}) {\n\t${2:// 异常处理}\n}",
			"异常捕获",
			CompletionItemKindSnippet,
		},

		// 返回语句片段
		{
			"return",
			"return $${1:value};",
			"返回语句",
			CompletionItemKindSnippet,
		},
		{
			"return_null",
			"return null;",
			"返回空值",
			CompletionItemKindSnippet,
		},
		{
			"return_true",
			"return true;",
			"返回真值",
			CompletionItemKindSnippet,
		},
		{
			"return_false",
			"return false;",
			"返回假值",
			CompletionItemKindSnippet,
		},
	}

	var items []CompletionItem
	for _, snippet := range snippets {
		item := CompletionItem{
			Label:            snippet.label,
			Kind:             &[]CompletionItemKind{snippet.kind}[0],
			Detail:           &[]string{snippet.detail}[0],
			InsertText:       &[]string{snippet.insertText}[0],
			InsertTextFormat: &[]InsertTextFormat{InsertTextFormatSnippet}[0],
		}
		items = append(items, item)
	}

	return items
}

// 关键字补全
func getKeywordCompletions() []CompletionItem {
	keywords := []string{
		"namespace", "use", "class", "interface", "trait", "abstract", "final",
		"public", "private", "protected", "static", "const", "var", "function",
		"return", "if", "else", "elseif", "for", "foreach", "while", "do",
		"switch", "case", "default", "break", "continue", "try", "catch",
		"finally", "throw", "new", "this", "parent", "self", "instanceof",
		"echo", "print", "include", "require", "include_once", "require_once",
		"isset", "unset", "empty", "die", "exit", "eval", "clone", "yield",
		"global", "unset", "isset", "empty", "array", "string", "int", "float",
		"bool", "null", "void", "object", "callable", "iterable", "mixed",
		"resource", "scalar", "numeric", "true", "false", "null",
	}

	var items []CompletionItem
	for _, keyword := range keywords {
		item := CompletionItem{
			Label:  keyword,
			Kind:   &[]CompletionItemKind{CompletionItemKindKeyword}[0],
			Detail: &[]string{"Origami 关键字"}[0],
		}
		items = append(items, item)
	}

	return items
}

// 类型补全
func getTypeCompletions() []CompletionItem {
	types := []string{
		"int", "integer", "float", "double", "string", "bool", "boolean",
		"array", "object", "null", "void", "mixed", "callable", "iterable",
		"resource", "scalar", "numeric", "number", "text", "binary",
	}

	var items []CompletionItem
	for _, typeName := range types {
		item := CompletionItem{
			Label:  typeName,
			Kind:   &[]CompletionItemKind{CompletionItemKindTypeParameter}[0],
			Detail: &[]string{"数据类型"}[0],
		}
		items = append(items, item)
	}

	return items
}

// 函数补全
func getFunctionCompletions() []CompletionItem {
	functions := []string{
		"function", "func", "method", "constructor", "destructor",
		"getter", "setter", "callback", "closure", "lambda",
	}

	var items []CompletionItem
	for _, funcName := range functions {
		item := CompletionItem{
			Label:  funcName,
			Kind:   &[]CompletionItemKind{CompletionItemKindFunction}[0],
			Detail: &[]string{"函数相关"}[0],
		}
		items = append(items, item)
	}

	return items
}

// 控制结构补全
func getControlStructureCompletions() []CompletionItem {
	controls := []string{
		"if", "else", "elseif", "endif", "for", "foreach", "while", "do",
		"switch", "case", "default", "endswitch", "break", "continue",
		"try", "catch", "finally", "throw", "goto", "label",
	}

	var items []CompletionItem
	for _, control := range controls {
		item := CompletionItem{
			Label:  control,
			Kind:   &[]CompletionItemKind{CompletionItemKindKeyword}[0],
			Detail: &[]string{"控制结构"}[0],
		}
		items = append(items, item)
	}

	return items
}

// 运算符补全
func getOperatorCompletions() []CompletionItem {
	operators := []string{
		"+", "-", "*", "/", "%", "**", "++", "--", "=", "+=", "-=", "*=", "/=", "%=",
		"==", "===", "!=", "!==", "<", ">", "<=", ">=", "<=>", "&&", "||", "!", "&", "|", "^", "~",
		".", "->", "::", "?", ":", "??", "??=", "[]", "()", "{}", "@", "#",
	}

	var items []CompletionItem
	for _, op := range operators {
		item := CompletionItem{
			Label:  op,
			Kind:   &[]CompletionItemKind{CompletionItemKindOperator}[0],
			Detail: &[]string{"运算符"}[0],
		}
		items = append(items, item)
	}

	return items
}

// 注解补全
func getAnnotationCompletions() []CompletionItem {
	annotations := []string{
		"Controller", "Route", "GetMapping", "PostMapping", "PutMapping", "DeleteMapping",
		"Autowired", "Inject", "Service", "Repository", "Component", "Bean",
		"ResponseBody", "RequestBody", "PathVariable", "RequestParam", "RequestHeader",
		"Valid", "Validate", "Cache", "Transactional", "Async", "Scheduled",
		"Deprecated", "Override", "SuppressWarnings", "Target", "Retention",
		"Documented", "Inherited", "Repeatable", "SafeVarargs", "FunctionalInterface",
	}

	var items []CompletionItem
	for _, annotation := range annotations {
		item := CompletionItem{
			Label:  annotation,
			Kind:   &[]CompletionItemKind{CompletionItemKindClass}[0],
			Detail: &[]string{"注解"}[0],
		}
		items = append(items, item)
	}

	return items
}

// 内置函数和类补全
func getBuiltinCompletions() []CompletionItem {
	builtins := []struct {
		name   string
		kind   CompletionItemKind
		detail string
	}{
		// 标准库函数
		{"dump", CompletionItemKindFunction, "标准库函数 - 调试输出"},
		{"include", CompletionItemKindFunction, "标准库函数 - 包含文件"},

		// 标准库类
		{"Log", CompletionItemKindClass, "标准库类 - 日志记录"},
		{"Exception", CompletionItemKindClass, "标准库类 - 异常处理"},
		{"OS", CompletionItemKindClass, "标准库类 - 操作系统接口"},
		{"Reflect", CompletionItemKindClass, "标准库类 - 反射功能"},
		{"DateTime", CompletionItemKindClass, "标准库类 - 日期时间"},
		{"Channel", CompletionItemKindClass, "标准库类 - 通道通信"},

		// 网络相关
		{"Request", CompletionItemKindClass, "HTTP请求类"},
		{"Response", CompletionItemKindClass, "HTTP响应类"},
		{"Server", CompletionItemKindClass, "HTTP服务器类"},

		// 字符串方法
		{"length", CompletionItemKindProperty, "字符串属性 - 长度"},
		{"substring", CompletionItemKindMethod, "字符串方法 - 子字符串"},
		{"indexOf", CompletionItemKindMethod, "字符串方法 - 查找位置"},
		{"startsWith", CompletionItemKindMethod, "字符串方法 - 开始检查"},
		{"endsWith", CompletionItemKindMethod, "字符串方法 - 结束检查"},
		{"toLowerCase", CompletionItemKindMethod, "字符串方法 - 转小写"},
		{"toUpperCase", CompletionItemKindMethod, "字符串方法 - 转大写"},
		{"trim", CompletionItemKindMethod, "字符串方法 - 去除空白"},
		{"split", CompletionItemKindMethod, "字符串方法 - 分割"},
		{"replace", CompletionItemKindMethod, "字符串方法 - 替换"},

		// 数组方法
		{"push", CompletionItemKindMethod, "数组方法 - 添加元素"},
		{"pop", CompletionItemKindMethod, "数组方法 - 移除末尾元素"},
		{"shift", CompletionItemKindMethod, "数组方法 - 移除开头元素"},
		{"unshift", CompletionItemKindMethod, "数组方法 - 开头添加元素"},
		{"slice", CompletionItemKindMethod, "数组方法 - 切片"},
		{"splice", CompletionItemKindMethod, "数组方法 - 拼接"},
		{"concat", CompletionItemKindMethod, "数组方法 - 连接"},
		{"join", CompletionItemKindMethod, "数组方法 - 连接为字符串"},
		{"reverse", CompletionItemKindMethod, "数组方法 - 反转"},
		{"sort", CompletionItemKindMethod, "数组方法 - 排序"},
		{"map", CompletionItemKindMethod, "数组方法 - 映射"},
		{"filter", CompletionItemKindMethod, "数组方法 - 过滤"},
		{"reduce", CompletionItemKindMethod, "数组方法 - 归约"},
		{"forEach", CompletionItemKindMethod, "数组方法 - 遍历"},
		{"find", CompletionItemKindMethod, "数组方法 - 查找"},
		{"findIndex", CompletionItemKindMethod, "数组方法 - 查找索引"},
		{"includes", CompletionItemKindMethod, "数组方法 - 包含检查"},
		{"indexOf", CompletionItemKindMethod, "数组方法 - 查找位置"},
		{"some", CompletionItemKindMethod, "数组方法 - 部分满足"},
		{"every", CompletionItemKindMethod, "数组方法 - 全部满足"},
		{"flat", CompletionItemKindMethod, "数组方法 - 扁平化"},
		{"flatMap", CompletionItemKindMethod, "数组方法 - 扁平映射"},
	}

	var items []CompletionItem
	for _, builtin := range builtins {
		item := CompletionItem{
			Label:  builtin.name,
			Kind:   &[]CompletionItemKind{builtin.kind}[0],
			Detail: &[]string{builtin.detail}[0],
		}
		items = append(items, item)
	}

	return items
}

// 对象方法补全 - 专门用于 $str-> 等场景
func getObjectMethodCompletions() []CompletionItem {
	methods := []struct {
		name   string
		kind   CompletionItemKind
		detail string
	}{
		// 字符串方法和属性
		{"length", CompletionItemKindProperty, "字符串属性 - 长度"},
		{"substring", CompletionItemKindMethod, "字符串方法 - 子字符串 (start, end)"},
		{"indexOf", CompletionItemKindMethod, "字符串方法 - 查找位置 (search)"},
		{"startsWith", CompletionItemKindMethod, "字符串方法 - 开始检查 (search)"},
		{"endsWith", CompletionItemKindMethod, "字符串方法 - 结束检查 (search)"},
		{"toLowerCase", CompletionItemKindMethod, "字符串方法 - 转小写"},
		{"toUpperCase", CompletionItemKindMethod, "字符串方法 - 转大写"},
		{"trim", CompletionItemKindMethod, "字符串方法 - 去除空白"},
		{"split", CompletionItemKindMethod, "字符串方法 - 分割 (delimiter)"},
		{"replace", CompletionItemKindMethod, "字符串方法 - 替换 (search, replace)"},

		// 数组方法
		{"push", CompletionItemKindMethod, "数组方法 - 添加元素 (item)"},
		{"pop", CompletionItemKindMethod, "数组方法 - 移除末尾元素"},
		{"shift", CompletionItemKindMethod, "数组方法 - 移除开头元素"},
		{"unshift", CompletionItemKindMethod, "数组方法 - 开头添加元素 (item)"},
		{"slice", CompletionItemKindMethod, "数组方法 - 切片 (start, end)"},
		{"splice", CompletionItemKindMethod, "数组方法 - 拼接 (start, deleteCount, items...)"},
		{"concat", CompletionItemKindMethod, "数组方法 - 连接 (arrays...)"},
		{"join", CompletionItemKindMethod, "数组方法 - 连接为字符串 (separator)"},
		{"reverse", CompletionItemKindMethod, "数组方法 - 反转"},
		{"sort", CompletionItemKindMethod, "数组方法 - 排序 (compareFunction)"},
		{"map", CompletionItemKindMethod, "数组方法 - 映射 (callback)"},
		{"filter", CompletionItemKindMethod, "数组方法 - 过滤 (callback)"},
		{"reduce", CompletionItemKindMethod, "数组方法 - 归约 (callback, initialValue)"},
		{"forEach", CompletionItemKindMethod, "数组方法 - 遍历 (callback)"},
		{"find", CompletionItemKindMethod, "数组方法 - 查找 (callback)"},
		{"findIndex", CompletionItemKindMethod, "数组方法 - 查找索引 (callback)"},
		{"includes", CompletionItemKindMethod, "数组方法 - 包含检查 (item)"},
		{"indexOf", CompletionItemKindMethod, "数组方法 - 查找位置 (item)"},
		{"some", CompletionItemKindMethod, "数组方法 - 部分满足 (callback)"},
		{"every", CompletionItemKindMethod, "数组方法 - 全部满足 (callback)"},
		{"flat", CompletionItemKindMethod, "数组方法 - 扁平化 (depth)"},
		{"flatMap", CompletionItemKindMethod, "数组方法 - 扁平映射 (callback)"},

		// 通用对象方法
		{"toString", CompletionItemKindMethod, "转换为字符串"},
		{"valueOf", CompletionItemKindMethod, "获取原始值"},
		{"hasOwnProperty", CompletionItemKindMethod, "检查是否具有自身属性 (prop)"},
		{"isPrototypeOf", CompletionItemKindMethod, "检查是否为原型 (obj)"},
		{"propertyIsEnumerable", CompletionItemKindMethod, "检查属性是否可枚举 (prop)"},
	}

	var items []CompletionItem
	for _, method := range methods {
		item := CompletionItem{
			Label:  method.name,
			Kind:   &[]CompletionItemKind{method.kind}[0],
			Detail: &[]string{method.detail}[0],
		}
		items = append(items, item)
	}

	return items
}

// 过滤后的关键字补全 - 排除已有代码片段的关键字
func getFilteredKeywordCompletions() []CompletionItem {
	// 定义已有代码片段的关键字
	snippetKeywords := map[string]bool{
		"for": true, "foreach": true, "while": true, "do": true,
		"if": true, "else": true, "elseif": true, "endif": true,
		"switch": true, "case": true, "default": true, "endswitch": true,
		"try": true, "catch": true, "finally": true, "throw": true,
		"function": true, "func": true, "class": true, "interface": true,
		"trait": true, "constructor": true, "destructor": true,
		"getter": true, "setter": true, "namespace": true, "use": true,
		"include": true, "require": true, "echo": true, "dump": true,
		"log": true, "new": true, "return": true, "var": true,
		"string": true, "int": true, "float": true, "bool": true,
		"array": true, "const": true,
	}

	keywords := []string{
		"abstract", "final", "public", "private", "protected", "static",
		"break", "continue", "goto", "label", "instanceof", "this", "parent", "self",
		"print", "include_once", "require_once", "isset", "unset", "empty", "die", "exit",
		"eval", "clone", "yield", "global", "void", "object", "callable", "iterable",
		"mixed", "resource", "scalar", "numeric", "number", "text", "binary",
		"true", "false", "null",
	}

	var items []CompletionItem
	for _, keyword := range keywords {
		// 只添加不在代码片段中的关键字
		if !snippetKeywords[keyword] {
			item := CompletionItem{
				Label:  keyword,
				Kind:   &[]CompletionItemKind{CompletionItemKindKeyword}[0],
				Detail: &[]string{"Origami 关键字"}[0],
			}
			items = append(items, item)
		}
	}

	return items
}

// 过滤后的控制结构补全 - 排除已有代码片段的控制结构
func getFilteredControlStructureCompletions() []CompletionItem {
	// 定义已有代码片段的控制结构
	snippetControls := map[string]bool{
		"if": true, "else": true, "elseif": true, "endif": true,
		"for": true, "foreach": true, "while": true, "do": true,
		"switch": true, "case": true, "default": true, "endswitch": true,
		"try": true, "catch": true, "finally": true, "throw": true,
	}

	controls := []string{
		"break", "continue", "goto", "label",
	}

	var items []CompletionItem
	for _, control := range controls {
		// 只添加不在代码片段中的控制结构
		if !snippetControls[control] {
			item := CompletionItem{
				Label:  control,
				Kind:   &[]CompletionItemKind{CompletionItemKindKeyword}[0],
				Detail: &[]string{"控制结构"}[0],
			}
			items = append(items, item)
		}
	}

	return items
}
