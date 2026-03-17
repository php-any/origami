package directory

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FilesystemIterator flags 常量（与 PHP 保持一致）
const (
	FSI_CURRENT_AS_PATHNAME = 32
	FSI_CURRENT_AS_FILEINFO = 0
	FSI_CURRENT_AS_SELF     = 16
	FSI_KEY_AS_PATHNAME     = 0
	FSI_KEY_AS_FILENAME     = 256
	FSI_FOLLOW_SYMLINKS     = 512
	FSI_SKIP_DOTS           = 4096
	FSI_UNIX_PATHS          = 8192
	FSI_NEW_CURRENT_AND_KEY = 256
	FSI_DEFAULT_FLAGS       = FSI_KEY_AS_PATHNAME | FSI_CURRENT_AS_FILEINFO | FSI_SKIP_DOTS
)

// FilesystemIteratorClass 提供 PHP FilesystemIterator 类定义
// 状态数据直接存储在类结构体中，方法通过 instance 字段访问（skill: php-class-state-sharing-pattern）
type FilesystemIteratorClass struct {
	node.Node
	path     string   // 目录路径（状态存储在类中）
	entries  []string // 过滤后的文件/目录名列表
	iterator int      // 当前迭代位置
	flags    int      // 当前 flags
}

// NewFilesystemIteratorClass 创建 FilesystemIteratorClass 实例
func NewFilesystemIteratorClass() *FilesystemIteratorClass {
	return &FilesystemIteratorClass{
		flags: FSI_DEFAULT_FLAGS,
	}
}

// GetValue 每次 new 时创建独立状态副本
func (c *FilesystemIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	clone := &FilesystemIteratorClass{
		path:     c.path,
		entries:  make([]string, len(c.entries)),
		iterator: c.iterator,
		flags:    c.flags,
	}
	copy(clone.entries, c.entries)
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "FilesystemIterator"
func (c *FilesystemIteratorClass) GetName() string { return "FilesystemIterator" }

// GetExtend FilesystemIterator 继承自 DirectoryIterator
func (c *FilesystemIteratorClass) GetExtend() *string {
	parent := "DirectoryIterator"
	return &parent
}

// GetImplements 实现 SeekableIterator 接口（继承自 DirectoryIterator 的 Iterator 已在继承链中）
func (c *FilesystemIteratorClass) GetImplements() []string {
	return []string{"SeekableIterator"}
}

// GetProperty 无额外属性
func (c *FilesystemIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetStaticProperty 返回 FilesystemIterator 类常量（与 PHP 保持一致）
func (c *FilesystemIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "CURRENT_AS_PATHNAME":
		return data.NewIntValue(FSI_CURRENT_AS_PATHNAME), true
	case "CURRENT_AS_FILEINFO":
		return data.NewIntValue(FSI_CURRENT_AS_FILEINFO), true
	case "CURRENT_AS_SELF":
		return data.NewIntValue(FSI_CURRENT_AS_SELF), true
	case "KEY_AS_PATHNAME":
		return data.NewIntValue(FSI_KEY_AS_PATHNAME), true
	case "KEY_AS_FILENAME":
		return data.NewIntValue(FSI_KEY_AS_FILENAME), true
	case "FOLLOW_SYMLINKS":
		return data.NewIntValue(FSI_FOLLOW_SYMLINKS), true
	case "SKIP_DOTS":
		return data.NewIntValue(FSI_SKIP_DOTS), true
	case "UNIX_PATHS":
		return data.NewIntValue(FSI_UNIX_PATHS), true
	case "NEW_CURRENT_AND_KEY":
		return data.NewIntValue(FSI_NEW_CURRENT_AND_KEY), true
	}
	return nil, false
}

// GetPropertyList 无额外属性
func (c *FilesystemIteratorClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 返回方法，传入当前实例引用（skill 关键：instance: c）
func (c *FilesystemIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &FilesystemIteratorConstructMethod{instance: c}, true
	case "current":
		return &FilesystemIteratorCurrentMethod{instance: c}, true
	case "key":
		return &FilesystemIteratorKeyMethod{instance: c}, true
	case "next":
		return &FilesystemIteratorNextMethod{instance: c}, true
	case "rewind":
		return &FilesystemIteratorRewindMethod{instance: c}, true
	case "valid":
		return &FilesystemIteratorValidMethod{instance: c}, true
	case "getFilename":
		return &FilesystemIteratorGetFilenameMethod{instance: c}, true
	case "getBasename":
		return &FilesystemIteratorGetBasenameMethod{instance: c}, true
	case "getExtension":
		return &FilesystemIteratorGetExtensionMethod{instance: c}, true
	case "getPath":
		return &FilesystemIteratorGetPathMethod{instance: c}, true
	case "getPathname":
		return &FilesystemIteratorGetPathnameMethod{instance: c}, true
	case "getRealPath":
		return &FilesystemIteratorGetRealPathMethod{instance: c}, true
	case "isDir":
		return &FilesystemIteratorIsDirMethod{instance: c}, true
	case "isFile":
		return &FilesystemIteratorIsFileMethod{instance: c}, true
	case "isDot":
		return &FilesystemIteratorIsDotMethod{instance: c}, true
	case "getSize":
		return &FilesystemIteratorGetSizeMethod{instance: c}, true
	case "getMTime":
		return &FilesystemIteratorGetMTimeMethod{instance: c}, true
	case "isReadable":
		return &FilesystemIteratorIsReadableMethod{instance: c}, true
	case "isWritable":
		return &FilesystemIteratorIsWritableMethod{instance: c}, true
	case "getFlags":
		return &FilesystemIteratorGetFlagsMethod{instance: c}, true
	case "setFlags":
		return &FilesystemIteratorSetFlagsMethod{instance: c}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *FilesystemIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&FilesystemIteratorConstructMethod{instance: c},
		&FilesystemIteratorCurrentMethod{instance: c},
		&FilesystemIteratorKeyMethod{instance: c},
		&FilesystemIteratorNextMethod{instance: c},
		&FilesystemIteratorRewindMethod{instance: c},
		&FilesystemIteratorValidMethod{instance: c},
		&FilesystemIteratorGetFilenameMethod{instance: c},
		&FilesystemIteratorGetBasenameMethod{instance: c},
		&FilesystemIteratorGetExtensionMethod{instance: c},
		&FilesystemIteratorGetPathMethod{instance: c},
		&FilesystemIteratorGetPathnameMethod{instance: c},
		&FilesystemIteratorGetRealPathMethod{instance: c},
		&FilesystemIteratorIsDirMethod{instance: c},
		&FilesystemIteratorIsFileMethod{instance: c},
		&FilesystemIteratorIsDotMethod{instance: c},
		&FilesystemIteratorGetSizeMethod{instance: c},
		&FilesystemIteratorGetMTimeMethod{instance: c},
		&FilesystemIteratorIsReadableMethod{instance: c},
		&FilesystemIteratorIsWritableMethod{instance: c},
		&FilesystemIteratorGetFlagsMethod{instance: c},
		&FilesystemIteratorSetFlagsMethod{instance: c},
	}
}

// GetConstruct 返回构造函数
func (c *FilesystemIteratorClass) GetConstruct() data.Method {
	return &FilesystemIteratorConstructMethod{instance: c}
}

// ------- 状态访问辅助方法（直接在类上操作，无需 ctx 属性传递）-------

func (c *FilesystemIteratorClass) currentEntry() string {
	if c.iterator >= 0 && c.iterator < len(c.entries) {
		return c.entries[c.iterator]
	}
	return ""
}

func (c *FilesystemIteratorClass) GetFilename() string { return c.currentEntry() }

func (c *FilesystemIteratorClass) GetPath() string { return c.path }

func (c *FilesystemIteratorClass) GetPathname() string {
	if c.iterator >= 0 && c.iterator < len(c.entries) {
		return filepath.Join(c.path, c.entries[c.iterator])
	}
	return c.path
}

func (c *FilesystemIteratorClass) GetRealPath() string {
	pathname := c.GetPathname()
	if realPath, err := filepath.EvalSymlinks(pathname); err == nil {
		return realPath
	}
	if absPath, err := filepath.Abs(pathname); err == nil {
		return absPath
	}
	return pathname
}

func (c *FilesystemIteratorClass) KeyStr() string {
	if c.flags&FSI_KEY_AS_FILENAME != 0 {
		return c.currentEntry()
	}
	return c.GetPathname()
}

func (c *FilesystemIteratorClass) IsDir() bool {
	fullPath := c.GetPathname()
	if fullPath == "" {
		return false
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (c *FilesystemIteratorClass) IsFile() bool {
	fullPath := c.GetPathname()
	if fullPath == "" {
		return false
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (c *FilesystemIteratorClass) IsDot() bool {
	name := c.currentEntry()
	return name == "." || name == ".."
}

func (c *FilesystemIteratorClass) GetBasename(suffix string) string {
	filename := c.GetFilename()
	if suffix != "" && len(filename) >= len(suffix) && filename[len(filename)-len(suffix):] == suffix {
		return filename[:len(filename)-len(suffix)]
	}
	return filename
}

func (c *FilesystemIteratorClass) GetExtension() string {
	filename := c.GetFilename()
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			if i > 0 {
				return filename[i+1:]
			}
			return ""
		}
		if filename[i] == '/' || filename[i] == '\\' {
			break
		}
	}
	return ""
}

func (c *FilesystemIteratorClass) GetSize() int64 {
	info, err := os.Stat(c.GetPathname())
	if err != nil {
		return 0
	}
	return info.Size()
}

func (c *FilesystemIteratorClass) GetMTime() int64 {
	info, err := os.Stat(c.GetPathname())
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}

func (c *FilesystemIteratorClass) IsReadable() bool {
	f, err := os.Open(c.GetPathname())
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func (c *FilesystemIteratorClass) IsWritable() bool {
	info, err := os.Stat(c.GetPathname())
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode&0200 != 0 || mode&0020 != 0 || mode&0002 != 0
}

// loadEntries 从磁盘读取目录条目并根据 flags 过滤
func (c *FilesystemIteratorClass) loadEntries() error {
	dir, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	skipDots := (c.flags & FSI_SKIP_DOTS) != 0
	filtered := make([]string, 0, len(names))
	for _, name := range names {
		if skipDots && (name == "." || name == "..") {
			continue
		}
		filtered = append(filtered, name)
	}
	sort.Strings(filtered)
	c.entries = filtered
	c.iterator = 0
	return nil
}
