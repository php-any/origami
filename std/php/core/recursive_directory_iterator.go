package core

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveDirectoryIteratorClass 实现 PHP 的 RecursiveDirectoryIterator 类
type RecursiveDirectoryIteratorClass struct {
	node.Node
	files []string // 文件列表
	pos   int      // 当前位置
	path  string   // 目录路径
}

// NewRecursiveDirectoryIteratorClass 创建新的 RecursiveDirectoryIterator 类
func NewRecursiveDirectoryIteratorClass() *RecursiveDirectoryIteratorClass {
	return &RecursiveDirectoryIteratorClass{
		files: []string{},
		pos:   0,
		path:  "",
	}
}

func (r *RecursiveDirectoryIteratorClass) GetName() string {
	return "RecursiveDirectoryIterator"
}

func (r *RecursiveDirectoryIteratorClass) GetExtend() *string {
	// RecursiveDirectoryIterator extends FilesystemIterator
	parent := "FilesystemIterator"
	return &parent
}

func (r *RecursiveDirectoryIteratorClass) GetImplements() []string {
	// RecursiveDirectoryIterator implements RecursiveIterator (on top of inherited SeekableIterator)
	return []string{"RecursiveIterator"}
}

func (r *RecursiveDirectoryIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetStaticProperty 返回 FilesystemIterator / RecursiveDirectoryIterator 常量（与 PHP 保持一致）
func (r *RecursiveDirectoryIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "CURRENT_AS_PATHNAME":
		return data.NewIntValue(32), true
	case "CURRENT_AS_FILEINFO":
		return data.NewIntValue(0), true
	case "CURRENT_AS_SELF":
		return data.NewIntValue(16), true
	case "KEY_AS_PATHNAME":
		return data.NewIntValue(0), true
	case "KEY_AS_FILENAME":
		return data.NewIntValue(256), true
	case "FOLLOW_SYMLINKS":
		return data.NewIntValue(512), true
	case "SKIP_DOTS":
		return data.NewIntValue(4096), true
	case "UNIX_PATHS":
		return data.NewIntValue(8192), true
	case "NEW_CURRENT_AND_KEY":
		return data.NewIntValue(256), true
	}
	return nil, false
}

func (r *RecursiveDirectoryIteratorClass) GetPropertyList() []data.Property {
	return nil
}

func (r *RecursiveDirectoryIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建新的实例，每个实例有自己的状态
	clone := &RecursiveDirectoryIteratorClass{
		files: make([]string, len(r.files)),
		pos:   r.pos,
		path:  r.path,
	}
	return data.NewClassValue(clone, ctx), nil
}

func (r *RecursiveDirectoryIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RecursiveDirectoryIteratorConstruct{instance: r}, true
	case "rewind":
		return &RecursiveDirectoryIteratorRewind{instance: r}, true
	case "current":
		return &RecursiveDirectoryIteratorCurrent{instance: r}, true
	case "key":
		return &RecursiveDirectoryIteratorKey{instance: r}, true
	case "next":
		return &RecursiveDirectoryIteratorNext{instance: r}, true
	case "valid":
		return &RecursiveDirectoryIteratorValid{instance: r}, true
	case "hasChildren":
		return &RecursiveDirectoryIteratorHasChildren{instance: r}, true
	case "getChildren":
		return &RecursiveDirectoryIteratorGetChildren{instance: r}, true
	case "getFilename":
		return &RecursiveDirectoryIteratorGetFilename{instance: r}, true
	case "getPathname":
		return &RecursiveDirectoryIteratorGetPathname{instance: r}, true
	case "isDir":
		return &RecursiveDirectoryIteratorIsDir{instance: r}, true
	case "isFile":
		return &RecursiveDirectoryIteratorIsFile{instance: r}, true
	}
	return nil, false
}

func (r *RecursiveDirectoryIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RecursiveDirectoryIteratorConstruct{instance: r},
		&RecursiveDirectoryIteratorRewind{instance: r},
		&RecursiveDirectoryIteratorCurrent{instance: r},
		&RecursiveDirectoryIteratorKey{instance: r},
		&RecursiveDirectoryIteratorNext{instance: r},
		&RecursiveDirectoryIteratorValid{instance: r},
		&RecursiveDirectoryIteratorHasChildren{instance: r},
		&RecursiveDirectoryIteratorGetChildren{instance: r},
		&RecursiveDirectoryIteratorGetFilename{instance: r},
		&RecursiveDirectoryIteratorGetPathname{instance: r},
		&RecursiveDirectoryIteratorIsDir{instance: r},
		&RecursiveDirectoryIteratorIsFile{instance: r},
	}
}

func (r *RecursiveDirectoryIteratorClass) GetConstruct() data.Method {
	return &RecursiveDirectoryIteratorConstruct{instance: r}
}

// RecursiveDirectoryIteratorConstruct 构造函数
type RecursiveDirectoryIteratorConstruct struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorConstruct) GetName() string            { return "__construct" }
func (m *RecursiveDirectoryIteratorConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorConstruct) GetIsStatic() bool          { return false }
func (m *RecursiveDirectoryIteratorConstruct) GetReturnType() data.Types  { return nil }
func (m *RecursiveDirectoryIteratorConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "flags", 1, nil, data.NewBaseType("int")),
	}
}
func (m *RecursiveDirectoryIteratorConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	}
}
func (m *RecursiveDirectoryIteratorConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathVal, _ := ctx.GetIndexValue(0)
	var path string
	if pathStr, ok := pathVal.(data.AsString); ok {
		path = pathStr.AsString()
	}

	// 直接访问实例状态
	m.instance.path = path
	m.instance.pos = 0
	m.instance.files = []string{}

	// 按照PHP RecursiveDirectoryIterator的行为：
	// 1. 首先添加 . 和 .. （特殊目录）
	// 2. 然后添加实际的文件和目录
	m.instance.files = append(m.instance.files, path+"/.")
	m.instance.files = append(m.instance.files, path+"/..")

	// 添加实际的文件和目录
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if name != "." && name != ".." {
			fullPath := filepath.Join(path, name)
			m.instance.files = append(m.instance.files, fullPath)
		}
	}

	return nil, nil
}

// RecursiveDirectoryIteratorRewind 重置迭代器
type RecursiveDirectoryIteratorRewind struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorRewind) GetName() string               { return "rewind" }
func (m *RecursiveDirectoryIteratorRewind) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorRewind) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorRewind) GetReturnType() data.Types     { return nil }
func (m *RecursiveDirectoryIteratorRewind) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorRewind) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorRewind) Call(ctx data.Context) (data.GetValue, data.Control) {
	m.instance.pos = 0
	return nil, nil
}

// RecursiveDirectoryIteratorCurrent 返回当前元素
type RecursiveDirectoryIteratorCurrent struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorCurrent) GetName() string               { return "current" }
func (m *RecursiveDirectoryIteratorCurrent) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorCurrent) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorCurrent) GetReturnType() data.Types     { return nil }
func (m *RecursiveDirectoryIteratorCurrent) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorCurrent) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorCurrent) Call(ctx data.Context) (data.GetValue, data.Control) {
	// PHP 的 RecursiveDirectoryIterator::current() 返回 $this，即迭代器对象本身
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return objCtx.ClassValue, nil
	}
	return data.NewNullValue(), nil
}

// RecursiveDirectoryIteratorKey 返回当前键
type RecursiveDirectoryIteratorKey struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorKey) GetName() string               { return "key" }
func (m *RecursiveDirectoryIteratorKey) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorKey) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorKey) GetReturnType() data.Types     { return data.Int{} }
func (m *RecursiveDirectoryIteratorKey) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorKey) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorKey) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		return data.NewStringValue(""), nil
	}
	// 返回完整路径作为key，与PHP原生行为一致
	return data.NewStringValue(m.instance.files[m.instance.pos]), nil
}

// RecursiveDirectoryIteratorNext 移动到下一个元素
type RecursiveDirectoryIteratorNext struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorNext) GetName() string               { return "next" }
func (m *RecursiveDirectoryIteratorNext) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorNext) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorNext) GetReturnType() data.Types     { return nil }
func (m *RecursiveDirectoryIteratorNext) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorNext) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorNext) Call(ctx data.Context) (data.GetValue, data.Control) {
	m.instance.pos++
	return nil, nil
}

// RecursiveDirectoryIteratorValid 检查当前位置是否有效
type RecursiveDirectoryIteratorValid struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorValid) GetName() string               { return "valid" }
func (m *RecursiveDirectoryIteratorValid) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorValid) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorValid) GetReturnType() data.Types     { return data.Bool{} }
func (m *RecursiveDirectoryIteratorValid) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorValid) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorValid) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.pos >= 0 && m.instance.pos < len(m.instance.files)), nil
}

// RecursiveDirectoryIteratorHasChildren 检查当前元素是否有子元素
type RecursiveDirectoryIteratorHasChildren struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorHasChildren) GetName() string { return "hasChildren" }
func (m *RecursiveDirectoryIteratorHasChildren) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveDirectoryIteratorHasChildren) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorHasChildren) GetReturnType() data.Types     { return data.Bool{} }
func (m *RecursiveDirectoryIteratorHasChildren) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorHasChildren) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorHasChildren) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		return data.NewBoolValue(false), nil
	}

	path := m.instance.files[m.instance.pos]
	info, err := os.Stat(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.IsDir()), nil
}

// RecursiveDirectoryIteratorGetChildren 返回子迭代器
type RecursiveDirectoryIteratorGetChildren struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorGetChildren) GetName() string { return "getChildren" }
func (m *RecursiveDirectoryIteratorGetChildren) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveDirectoryIteratorGetChildren) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorGetChildren) GetReturnType() data.Types     { return nil }
func (m *RecursiveDirectoryIteratorGetChildren) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorGetChildren) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorGetChildren) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		// 返回空的迭代器
		emptyClass := NewRecursiveDirectoryIteratorClass()
		return data.NewClassValue(emptyClass, ctx.CreateBaseContext()), nil
	}

	path := m.instance.files[m.instance.pos]

	// 创建子迭代器实例
	childClass := NewRecursiveDirectoryIteratorClass()
	childClass.path = path

	// 加载子目录内容
	entries, err := os.ReadDir(path)
	if err != nil {
		return data.NewClassValue(childClass, ctx.CreateBaseContext()), nil
	}

	for _, entry := range entries {
		name := entry.Name()
		if name != "." && name != ".." {
			childClass.files = append(childClass.files, filepath.Join(path, name))
		}
	}

	return data.NewClassValue(childClass, ctx.CreateBaseContext()), nil
}

// RecursiveDirectoryIteratorGetFilename 获取文件名
type RecursiveDirectoryIteratorGetFilename struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorGetFilename) GetName() string { return "getFilename" }
func (m *RecursiveDirectoryIteratorGetFilename) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveDirectoryIteratorGetFilename) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorGetFilename) GetReturnType() data.Types     { return data.String{} }
func (m *RecursiveDirectoryIteratorGetFilename) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorGetFilename) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorGetFilename) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		return data.NewStringValue(""), nil
	}

	currentPath := m.instance.files[m.instance.pos]
	_, filename := filepath.Split(currentPath)

	// 特殊处理 . 和 .. 目录
	if filename == "." || filename == ".." {
		return data.NewStringValue(filename), nil
	}

	return data.NewStringValue(filename), nil
}

// RecursiveDirectoryIteratorGetPathname 获取完整路径名
type RecursiveDirectoryIteratorGetPathname struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorGetPathname) GetName() string { return "getPathname" }
func (m *RecursiveDirectoryIteratorGetPathname) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveDirectoryIteratorGetPathname) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorGetPathname) GetReturnType() data.Types     { return data.String{} }
func (m *RecursiveDirectoryIteratorGetPathname) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorGetPathname) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorGetPathname) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(m.instance.files[m.instance.pos]), nil
}

// RecursiveDirectoryIteratorIsDir 检查是否为目录
type RecursiveDirectoryIteratorIsDir struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorIsDir) GetName() string               { return "isDir" }
func (m *RecursiveDirectoryIteratorIsDir) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorIsDir) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorIsDir) GetReturnType() data.Types     { return data.Bool{} }
func (m *RecursiveDirectoryIteratorIsDir) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorIsDir) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorIsDir) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		return data.NewBoolValue(false), nil
	}

	currentPath := m.instance.files[m.instance.pos]
	info, err := os.Stat(currentPath)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.IsDir()), nil
}

// RecursiveDirectoryIteratorIsFile 检查是否为文件
type RecursiveDirectoryIteratorIsFile struct {
	instance *RecursiveDirectoryIteratorClass
}

func (m *RecursiveDirectoryIteratorIsFile) GetName() string               { return "isFile" }
func (m *RecursiveDirectoryIteratorIsFile) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RecursiveDirectoryIteratorIsFile) GetIsStatic() bool             { return false }
func (m *RecursiveDirectoryIteratorIsFile) GetReturnType() data.Types     { return data.Bool{} }
func (m *RecursiveDirectoryIteratorIsFile) GetParams() []data.GetValue    { return nil }
func (m *RecursiveDirectoryIteratorIsFile) GetVariables() []data.Variable { return nil }
func (m *RecursiveDirectoryIteratorIsFile) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.pos < 0 || m.instance.pos >= len(m.instance.files) {
		return data.NewBoolValue(false), nil
	}

	currentPath := m.instance.files[m.instance.pos]
	info, err := os.Stat(currentPath)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(!info.IsDir()), nil
}
