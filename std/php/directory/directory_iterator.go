package directory

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DirectoryIteratorClass 提供 PHP DirectoryIterator 类定义
// DirectoryIterator 用于遍历目录中的文件和子目录
type DirectoryIteratorClass struct {
	node.Node
}

// GetValue 创建 DirectoryIterator 的实例
func (c *DirectoryIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "DirectoryIterator"
func (c *DirectoryIteratorClass) GetName() string { return "DirectoryIterator" }

// GetExtend 返回父类名，DirectoryIterator 没有父类
func (c *DirectoryIteratorClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，DirectoryIterator 实现 Iterator 接口
func (c *DirectoryIteratorClass) GetImplements() []string {
	return []string{"Iterator"}
}

// GetProperty 获取属性，DirectoryIterator 没有属性
func (c *DirectoryIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，DirectoryIterator 没有属性
func (c *DirectoryIteratorClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
func (c *DirectoryIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &DirectoryIteratorConstructMethod{}, true
	case "current":
		return &DirectoryIteratorCurrentMethod{}, true
	case "key":
		return &DirectoryIteratorKeyMethod{}, true
	case "next":
		return &DirectoryIteratorNextMethod{}, true
	case "rewind":
		return &DirectoryIteratorRewindMethod{}, true
	case "valid":
		return &DirectoryIteratorValidMethod{}, true
	case "getFilename":
		return &DirectoryIteratorGetFilenameMethod{}, true
	case "getPath":
		return &DirectoryIteratorGetPathMethod{}, true
	case "getPathname":
		return &DirectoryIteratorGetPathnameMethod{}, true
	case "isDir":
		return &DirectoryIteratorIsDirMethod{}, true
	case "isFile":
		return &DirectoryIteratorIsFileMethod{}, true
	case "isDot":
		return &DirectoryIteratorIsDotMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *DirectoryIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&DirectoryIteratorConstructMethod{},
		&DirectoryIteratorCurrentMethod{},
		&DirectoryIteratorKeyMethod{},
		&DirectoryIteratorNextMethod{},
		&DirectoryIteratorRewindMethod{},
		&DirectoryIteratorValidMethod{},
		&DirectoryIteratorGetFilenameMethod{},
		&DirectoryIteratorGetPathMethod{},
		&DirectoryIteratorGetPathnameMethod{},
		&DirectoryIteratorIsDirMethod{},
		&DirectoryIteratorIsFileMethod{},
		&DirectoryIteratorIsDotMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *DirectoryIteratorClass) GetConstruct() data.Method {
	return &DirectoryIteratorConstructMethod{}
}

// getDirectoryIteratorInfo 从上下文中获取 DirectoryIterator 的信息
func getDirectoryIteratorInfo(ctx data.Context) (*DirectoryIteratorData, bool) {
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		if objCtx.ObjectValue != nil {
			props := objCtx.ObjectValue.GetProperties()
			iterVal, hasIter := props["_iterator"]
			if hasIter {
				if iter, ok := iterVal.(*DirectoryIteratorData); ok {
					return iter, true
				}
			}
		}
	}
	return nil, false
}

// DirectoryIteratorData 存储 DirectoryIterator 的数据
type DirectoryIteratorData struct {
	path     string   // 目录路径
	entries  []string // 文件/目录名列表
	iterator int      // 当前迭代位置
}

// GetValue 实现 data.GetValue 接口
func (d *DirectoryIteratorData) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return d, nil
}

// AsString 实现 data.Value 接口
func (d *DirectoryIteratorData) AsString() string {
	return d.GetPathname()
}

// NewDirectoryIteratorData 创建新的 DirectoryIteratorData
func NewDirectoryIteratorData(path string) (*DirectoryIteratorData, error) {
	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// 读取目录内容
	entries, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// 对文件名进行排序（PHP DirectoryIterator 默认按字母顺序排序）
	sort.Strings(entries)

	return &DirectoryIteratorData{
		path:     path,
		entries:  entries,
		iterator: 0,
	}, nil
}

// Current 返回当前文件名
func (d *DirectoryIteratorData) Current() string {
	if d.iterator >= 0 && d.iterator < len(d.entries) {
		return d.entries[d.iterator]
	}
	return ""
}

// Key 返回当前索引
func (d *DirectoryIteratorData) Key() int {
	return d.iterator
}

// Next 移动到下一个
func (d *DirectoryIteratorData) Next() {
	d.iterator++
}

// Rewind 重置迭代器
func (d *DirectoryIteratorData) Rewind() {
	d.iterator = 0
}

// Valid 检查迭代器是否有效
func (d *DirectoryIteratorData) Valid() bool {
	return d.iterator >= 0 && d.iterator < len(d.entries)
}

// GetFilename 获取当前文件名
func (d *DirectoryIteratorData) GetFilename() string {
	return d.Current()
}

// GetPath 获取目录路径
func (d *DirectoryIteratorData) GetPath() string {
	return d.path
}

// GetPathname 获取完整路径名
func (d *DirectoryIteratorData) GetPathname() string {
	if d.iterator >= 0 && d.iterator < len(d.entries) {
		return filepath.Join(d.path, d.entries[d.iterator])
	}
	return d.path
}

// IsDir 检查当前项是否为目录
func (d *DirectoryIteratorData) IsDir() bool {
	if d.iterator >= 0 && d.iterator < len(d.entries) {
		fullPath := filepath.Join(d.path, d.entries[d.iterator])
		info, err := os.Stat(fullPath)
		if err != nil {
			return false
		}
		return info.IsDir()
	}
	return false
}

// IsFile 检查当前项是否为文件
func (d *DirectoryIteratorData) IsFile() bool {
	if d.iterator >= 0 && d.iterator < len(d.entries) {
		fullPath := filepath.Join(d.path, d.entries[d.iterator])
		info, err := os.Stat(fullPath)
		if err != nil {
			return false
		}
		return !info.IsDir()
	}
	return false
}

// IsDot 检查当前项是否为 "." 或 ".."
func (d *DirectoryIteratorData) IsDot() bool {
	if d.iterator >= 0 && d.iterator < len(d.entries) {
		filename := d.entries[d.iterator]
		return filename == "." || filename == ".."
	}
	return false
}
