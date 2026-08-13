package php

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GlobFunction 实现 glob 函数。
// glob(string $pattern, int $flags = 0): array|false
// 基于 filepath.Glob 匹配通配符；GLOB_BRACE 支持 {a,b} 花括号展开。
// GLOB_NOESCAPE 下反斜杠转义语义与 PHP 存在差异，未做完整模拟。
type GlobFunction struct{}

func NewGlobFunction() data.FuncStmt {
	return &GlobFunction{}
}

func (f *GlobFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	flagsValue, _ := ctx.GetIndexValue(1)

	var pattern string
	switch p := patternValue.(type) {
	case data.AsString:
		pattern = p.AsString()
	default:
		if patternValue != nil {
			pattern = patternValue.AsString()
		}
	}

	flags := 0
	if flagsValue != nil {
		if asInt, ok := flagsValue.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				flags = v
			}
		}
	}

	if pattern == "" {
		return data.NewBoolValue(false), nil
	}

	patterns := []string{pattern}
	if flags&GLOB_BRACE != 0 {
		patterns = expandGlobBraces(pattern)
	}

	var matches []string
	for _, p := range patterns {
		found, err := filepath.Glob(p)
		if err != nil {
			if flags&GLOB_ERR != 0 {
				return data.NewBoolValue(false), nil
			}
			return data.NewBoolValue(false), nil
		}
		matches = append(matches, found...)
	}

	if flags&GLOB_NOSORT == 0 {
		sort.Strings(matches)
	} else {
		matches = dedupePreserveOrder(matches)
	}

	if flags&GLOB_ONLYDIR != 0 {
		filtered := matches[:0]
		for _, m := range matches {
			info, err := os.Stat(m)
			if err != nil {
				if flags&GLOB_ERR != 0 {
					return data.NewBoolValue(false), nil
				}
				continue
			}
			if info.IsDir() {
				filtered = append(filtered, m)
			}
		}
		matches = filtered
	}

	if len(matches) == 0 {
		if flags&GLOB_NOCHECK != 0 {
			matches = []string{pattern}
		} else {
			return data.NewArrayValue([]data.Value{}), nil
		}
	}

	if flags&GLOB_MARK != 0 {
		for i, m := range matches {
			info, err := os.Stat(m)
			if err != nil {
				if flags&GLOB_ERR != 0 {
					return data.NewBoolValue(false), nil
				}
				continue
			}
			if info.IsDir() && !strings.HasSuffix(m, "/") {
				matches[i] = m + "/"
			}
		}
	}

	values := make([]data.Value, len(matches))
	for i, m := range matches {
		values[i] = data.NewStringValue(m)
	}
	return data.NewArrayValue(values), nil
}

func expandGlobBraces(pattern string) []string {
	start := strings.Index(pattern, "{")
	if start < 0 {
		return []string{pattern}
	}
	depth := 0
	end := -1
	for i := start; i < len(pattern); i++ {
		switch pattern[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		return []string{pattern}
	}

	prefix := pattern[:start]
	suffixes := expandGlobBraces(pattern[end+1:])
	alternatives := strings.Split(pattern[start+1:end], ",")

	var results []string
	for _, alt := range alternatives {
		for _, suffix := range suffixes {
			results = append(results, prefix+alt+suffix)
		}
	}
	return results
}

func dedupePreserveOrder(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func (f *GlobFunction) GetName() string {
	return "glob"
}

func (f *GlobFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), nil),
	}
}

func (f *GlobFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	}
}
