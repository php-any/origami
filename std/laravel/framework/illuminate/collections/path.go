package collections

import "github.com/php-any/origami/data"

// PathGet 点号路径读取（供 config 等包复用）。
func PathGet(target data.Value, path string) (data.Value, bool) {
	return dataGetPath(target, path)
}

// PathSet 点号路径写入。
func PathSet(target *data.ArrayValue, path string, value data.Value) {
	dataSetPath(target, path, value)
}

// PathHas 点号路径是否存在。
func PathHas(target data.Value, path string) bool {
	_, ok := dataGetPath(target, path)
	return ok
}
