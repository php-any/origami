package core

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
)

// iniDefaults 模拟 PHP 内置 ini 默认值（未在 iniStore 中显式设置时生效）。
var iniDefaults = map[string]string{
	"default_charset":          "UTF-8",
	"input_encoding":           "",
	"internal_encoding":        "",
	"output_encoding":          "",
	"precision":                "14",
	"serialize_precision":      "-1",
	"memory_limit":             "128M",
	"max_memory_limit":         "-1",
	"post_max_size":            "8M",
	"upload_max_filesize":      "2M",
	"max_file_uploads":         "20",
	"file_uploads":             "1",
	"enable_post_data_reading": "1",
	"register_argc_argv":       "1",
}

// InitIniDefaults 在运行时启动时应用默认 ini，并合并 PHPT 传入的 ORIGAMI_PHPT_INI。
func InitIniDefaults() {
	for key, value := range iniDefaults {
		if _, ok := IniGet(key); !ok {
			iniStore.Store(key, value)
		}
	}
	applyIniFromEnv(os.Getenv("ORIGAMI_PHPT_INI"))
	applyMemoryLimitCap()
}

// applyMemoryLimitCap 在加载 ini 后根据 max_memory_limit 限制 memory_limit（GH-17951）。
func applyMemoryLimitCap() {
	maxMem, okMax := IniGet("max_memory_limit")
	mem, okMem := IniGet("memory_limit")
	if !okMax || !okMem || maxMem == "-1" {
		return
	}
	maxBytes, ok1 := ParseIniSizeBytes(maxMem)
	memBytes, ok2 := ParseIniSizeBytes(mem)
	if !ok1 {
		return
	}
	if mem == "-1" || (ok2 && memBytes > maxBytes) {
		iniStore.Store("memory_limit", maxMem)
	}
}

// ParseIniSizeBytes 解析 ini 尺寸字符串（如 128M）。
func ParseIniSizeBytes(raw string) (int, bool) {
	raw = strings.TrimSpace(strings.ToUpper(raw))
	if raw == "" || raw == "-1" {
		return 0, false
	}
	multiplier := 1
	switch raw[len(raw)-1] {
	case 'K':
		multiplier = 1024
		raw = raw[:len(raw)-1]
	case 'M':
		multiplier = 1024 * 1024
		raw = raw[:len(raw)-1]
	case 'G':
		multiplier = 1024 * 1024 * 1024
		raw = raw[:len(raw)-1]
	}
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n < 0 {
		return 0, false
	}
	return n * multiplier, true
}

func applyIniFromEnv(raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err == nil {
		raw = string(decoded)
	}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		iniStore.Store(key, value)
	}
}

// ApplyIniMap 批量设置 ini（PHPT --INI-- 等）。
func ApplyIniMap(values map[string]string) {
	for key, value := range values {
		iniStore.Store(strings.ToLower(key), value)
	}
}
