package core

import (
	"strings"
	"sync"
)

// iniStore 是全局 PHP ini 配置存储，模拟 php.ini 运行时可修改的配置项。
var iniStore sync.Map

// IniSet 设置配置项，返回旧值（若不存在则返回 false 字符串标识）。
func IniSet(key, value string) (string, bool) {
	key = strings.ToLower(key)
	if old, loaded := iniStore.Load(key); loaded {
		if s, ok := old.(string); ok {
			iniStore.Store(key, value)
			return s, true
		}
	}
	prev, had := IniGet(key)
	iniStore.Store(key, value)
	if had {
		return prev, true
	}
	return "", false
}

// IniGet 获取配置项，若不存在返回 ("", false)。
func IniGet(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := iniStore.Load(key)
	if ok {
		s, _ := v.(string)
		return s, true
	}
	if d, ok := iniDefaults[key]; ok {
		return d, true
	}
	return "", false
}
