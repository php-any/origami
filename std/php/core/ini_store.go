package core

import "sync"

// iniStore 是全局 PHP ini 配置存储，模拟 php.ini 运行时可修改的配置项。
var iniStore sync.Map

// IniSet 设置配置项，返回旧值（若不存在则返回 false 字符串标识）。
func IniSet(key, value string) (string, bool) {
	old, loaded := iniStore.Swap(key, value)
	if !loaded {
		return "", false
	}
	if s, ok := old.(string); ok {
		return s, true
	}
	return "", false
}

// IniGet 获取配置项，若不存在返回 ("", false)。
func IniGet(key string) (string, bool) {
	v, ok := iniStore.Load(key)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
