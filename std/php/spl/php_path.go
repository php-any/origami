package spl

import "strings"

// phpPathJoin 拼接目录与文件名，不 Clean「..」（对齐 PHP SplFileInfo / DirectoryIterator）。
// Go 的 filepath.Join / Dir 会解析「..」，导致 Flysystem PathPrefixer::stripPrefix
// 在根路径含 src/../ 时截断错误文件名。
func phpPathJoin(dir, name string) string {
	if dir == "" {
		return name
	}
	if name == "" {
		return strings.TrimRight(dir, `/\`)
	}
	dir = strings.TrimRight(dir, `/\`)
	name = strings.TrimLeft(name, `/\`)
	return dir + "/" + name
}

// phpPathDir 返回路径的目录部分，保留「..」段（对齐 PHP dirname / SplFileInfo::getPath）。
func phpPathDir(path string) string {
	if path == "" {
		return "."
	}
	// 去掉尾部分隔符（根 "/" 除外）
	for len(path) > 1 && (path[len(path)-1] == '/' || path[len(path)-1] == '\\') {
		path = path[:len(path)-1]
	}
	i := strings.LastIndexAny(path, `/\`)
	if i < 0 {
		return "."
	}
	if i == 0 {
		return path[:1]
	}
	return path[:i]
}

// phpPathBase 返回路径最后一段（文件名），不依赖 filepath.Clean。
func phpPathBase(path string) string {
	if path == "" {
		return ""
	}
	for len(path) > 1 && (path[len(path)-1] == '/' || path[len(path)-1] == '\\') {
		path = path[:len(path)-1]
	}
	i := strings.LastIndexAny(path, `/\`)
	if i < 0 {
		return path
	}
	return path[i+1:]
}

// phpPathRel 计算 full 相对 root 的路径；两端保持逻辑形式，优先字符串前缀剥离。
// 对齐 PHP RecursiveDirectoryIterator::getSubPathname：根层文件返回文件名，不返回 "."。
func phpPathRel(root, full string) string {
	root = strings.TrimRight(root, `/\`)
	full = strings.TrimRight(full, `/\`)
	if full == "" || full == root {
		return ""
	}
	prefix := root + "/"
	if strings.HasPrefix(full, prefix) {
		return full[len(prefix):]
	}
	return ""
}
