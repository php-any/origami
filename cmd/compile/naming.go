package compile

import (
	"path/filepath"
	"strings"
	"unicode"
)

// pathToFuncSuffix 将 PHP 路径转为 Go 标识符后缀（不含 AST_ 前缀）。
// 路径分隔符 / \ 划分的每一段单独 Pascal 化后用 _ 连接；
// 段内的 _ 转为驼峰（user_login → UserLogin），从而区分 user/login 与 user_login。
func pathToFuncSuffix(path string) string {
	path = strings.TrimSuffix(path, ".php")
	path = strings.TrimSuffix(path, ".zy")
	path = filepath.ToSlash(path)

	var segments []string
	for _, seg := range strings.Split(path, "/") {
		seg = sanitizePathSegment(seg)
		if seg == "" {
			continue
		}
		segments = append(segments, segmentToCamel(seg))
	}
	if len(segments) == 0 {
		return "File"
	}
	return strings.Join(segments, "_")
}

// funcNameFromPath 返回预编译 AST 函数名，如 AST_Examples_Spring_Src_Config_AppConfig
func funcNameFromPath(path string) string {
	suffix := pathToFuncSuffix(path)
	if suffix == "" {
		return "AST_File"
	}
	if suffix[0] >= '0' && suffix[0] <= '9' {
		suffix = "F_" + suffix
	}
	return "AST_" + suffix
}

// goFileNameFromPath 返回同包内的 Go 源文件名，如 ast_examples_spring_user_login.go
func goFileNameFromPath(path string) string {
	return strings.ToLower(funcNameFromPath(path)) + ".go"
}

func sanitizePathSegment(seg string) string {
	var b strings.Builder
	for _, c := range seg {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			b.WriteRune(c)
		} else if c == '_' || c == '-' {
			b.WriteByte('_')
		}
	}
	return b.String()
}

// segmentToCamel 将路径段（可含 _）转为 PascalCase 且无内部分隔符：user_login → UserLogin
func segmentToCamel(seg string) string {
	parts := strings.Split(seg, "_")
	var b strings.Builder
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		b.WriteString(titleWord(p))
	}
	if b.Len() == 0 {
		return "X"
	}
	return b.String()
}

func titleWord(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
