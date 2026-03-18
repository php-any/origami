package core

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveDirectoryIterator 内部状态属性名（存储在 ClassValue.ObjectValue.property）
// 类似 FilterIterator 模式，保证 PHP 子类继承时状态隔离
const (
	rdiPathKey  = "__rdi_path__"  // string: 根目录路径
	rdiFilesKey = "__rdi_files__" // *rdiFilesValue: 文件列表
	rdiPosKey   = "__rdi_pos__"   // int: 当前位置
	rdiFlagsKey = "__rdi_flags__" // int: flags
)

// rdiFilesValue 包装文件列表，实现 data.Value 接口
type rdiFilesValue struct {
	files []string
}

func (s *rdiFilesValue) GetValue(ctx data.Context) (data.GetValue, data.Control) { return s, nil }
func (s *rdiFilesValue) AsString() string                                        { return "rdiFiles" }
func (s *rdiFilesValue) Marshal(serializer data.Serializer) ([]byte, error)      { return nil, nil }
func (s *rdiFilesValue) Unmarshal(b []byte, serializer data.Serializer) error    { return nil }
func (s *rdiFilesValue) ToGoValue(serializer data.Serializer) (any, error)       { return nil, nil }

// RecursiveDirectoryIteratorClass 实现 PHP 的 RecursiveDirectoryIterator 类
// 状态通过 ClassValue.ObjectValue.property 存储，方法通过 ctx 访问（skill: php-class-state-sharing-pattern）
type RecursiveDirectoryIteratorClass struct {
	node.Node
}

// NewRecursiveDirectoryIteratorClass 创建新的 RecursiveDirectoryIterator 类
func NewRecursiveDirectoryIteratorClass() *RecursiveDirectoryIteratorClass {
	return &RecursiveDirectoryIteratorClass{}
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
	cv := data.NewClassValue(r, ctx.CreateBaseContext())
	// 初始化实例属性
	cv.SetProperty(rdiPathKey, data.NewStringValue(""))
	cv.SetProperty(rdiFilesKey, &rdiFilesValue{})
	cv.SetProperty(rdiPosKey, data.NewIntValue(0))
	cv.SetProperty(rdiFlagsKey, data.NewIntValue(4096)) // SKIP_DOTS
	return cv, nil
}

func (r *RecursiveDirectoryIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RDIConstruct{}, true
	case "rewind":
		return &RDIRewind{}, true
	case "current":
		return &RDICurrent{}, true
	case "key":
		return &RDIKey{}, true
	case "next":
		return &RDINext{}, true
	case "valid":
		return &RDIValid{}, true
	case "hasChildren":
		return &RDIHasChildren{}, true
	case "getChildren":
		return &RDIGetChildren{}, true
	case "getFilename":
		return &RDIGetFilename{}, true
	case "getPathname":
		return &RDIGetPathname{}, true
	case "getPath":
		return &RDIGetPath{}, true
	case "getRealPath":
		return &RDIGetRealPath{}, true
	case "getSubPath":
		return &RDIGetSubPath{}, true
	case "getSubPathname":
		return &RDIGetSubPathname{}, true
	case "isDir":
		return &RDIIsDir{}, true
	case "isFile":
		return &RDIIsFile{}, true
	}
	return nil, false
}

func (r *RecursiveDirectoryIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RDIConstruct{},
		&RDIRewind{},
		&RDICurrent{},
		&RDIKey{},
		&RDINext{},
		&RDIValid{},
		&RDIHasChildren{},
		&RDIGetChildren{},
		&RDIGetFilename{},
		&RDIGetPathname{},
		&RDIGetPath{},
		&RDIGetRealPath{},
		&RDIGetSubPath{},
		&RDIGetSubPathname{},
		&RDIIsDir{},
		&RDIIsFile{},
	}
}

func (r *RecursiveDirectoryIteratorClass) GetConstruct() data.Method {
	return &RDIConstruct{}
}

// ---- 辅助：从 ctx 获取 ClassValue ----

func rdiGetCV(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

// ---- 状态读写 ----

func rdiGetPath(cv *data.ClassValue) string {
	v, _ := cv.ObjectValue.GetProperty(rdiPathKey)
	if sv, ok := v.(*data.StringValue); ok {
		return sv.Value
	}
	return ""
}

func rdiSetPath(cv *data.ClassValue, path string) {
	cv.ObjectValue.SetProperty(rdiPathKey, data.NewStringValue(path))
}

func rdiGetFiles(cv *data.ClassValue) []string {
	v, _ := cv.ObjectValue.GetProperty(rdiFilesKey)
	if fv, ok := v.(*rdiFilesValue); ok {
		return fv.files
	}
	return nil
}

func rdiSetFiles(cv *data.ClassValue, files []string) {
	v, _ := cv.ObjectValue.GetProperty(rdiFilesKey)
	if fv, ok := v.(*rdiFilesValue); ok {
		fv.files = files
	} else {
		cv.ObjectValue.SetProperty(rdiFilesKey, &rdiFilesValue{files: files})
	}
}

func rdiGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(rdiPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return int(iv.Value)
	}
	return 0
}

func rdiSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(rdiPosKey, data.NewIntValue(pos))
}

func rdiGetFlags(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(rdiFlagsKey)
	if iv, ok := v.(*data.IntValue); ok {
		return int(iv.Value)
	}
	return 4096 // SKIP_DOTS
}

func rdiSetFlags(cv *data.ClassValue, flags int) {
	cv.ObjectValue.SetProperty(rdiFlagsKey, data.NewIntValue(flags))
}

// ---- 状态辅助方法 ----

func rdiCurrentFile(cv *data.ClassValue) string {
	files := rdiGetFiles(cv)
	pos := rdiGetPos(cv)
	if pos >= 0 && pos < len(files) {
		return files[pos]
	}
	return ""
}

func rdiLoadFiles(cv *data.ClassValue, path string, flags int) error {
	skipDots := (flags & 4096) != 0 // SKIP_DOTS = 4096

	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	var files []string
	if !skipDots {
		files = append(files, path+"/.")
		files = append(files, path+"/..")
	}

	sorted := make([]string, 0, len(entries))
	for _, name := range entries {
		if name != "." && name != ".." {
			sorted = append(sorted, name)
		}
	}
	sort.Strings(sorted)

	for _, name := range sorted {
		files = append(files, filepath.Join(path, name))
	}

	rdiSetFiles(cv, files)
	rdiSetPos(cv, 0)
	return nil
}

// ---- __construct ----

type RDIConstruct struct{}

func (m *RDIConstruct) GetName() string            { return "__construct" }
func (m *RDIConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RDIConstruct) GetIsStatic() bool          { return false }
func (m *RDIConstruct) GetReturnType() data.Types  { return nil }
func (m *RDIConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "flags", 1, nil, data.NewBaseType("int")),
	}
}
func (m *RDIConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	}
}
func (m *RDIConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}

	pathVal, _ := ctx.GetIndexValue(0)
	path := ""
	if ps, ok := pathVal.(data.AsString); ok {
		path = ps.AsString()
	}

	flagsVal, hasFlagsVal := ctx.GetIndexValue(1)
	flags := 4096 // 默认 SKIP_DOTS
	if hasFlagsVal {
		if fi, ok := flagsVal.(interface{ AsInt() int }); ok {
			flags = fi.AsInt()
		} else if fi64, ok := flagsVal.(interface{ AsInt64() int64 }); ok {
			flags = int(fi64.AsInt64())
		}
	}

	rdiSetPath(cv, path)
	rdiSetFlags(cv, flags)

	if err := rdiLoadFiles(cv, path, flags); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	return nil, nil
}

// ---- rewind ----

type RDIRewind struct{}

func (m *RDIRewind) GetName() string               { return "rewind" }
func (m *RDIRewind) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIRewind) GetIsStatic() bool             { return false }
func (m *RDIRewind) GetReturnType() data.Types     { return nil }
func (m *RDIRewind) GetParams() []data.GetValue    { return nil }
func (m *RDIRewind) GetVariables() []data.Variable { return nil }
func (m *RDIRewind) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	rdiSetPos(cv, 0)
	return nil, nil
}

// ---- current：返回 $this（ClassValue 本身，即迭代器对象作为 SplFileInfo） ----

type RDICurrent struct{}

func (m *RDICurrent) GetName() string               { return "current" }
func (m *RDICurrent) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDICurrent) GetIsStatic() bool             { return false }
func (m *RDICurrent) GetReturnType() data.Types     { return nil }
func (m *RDICurrent) GetParams() []data.GetValue    { return nil }
func (m *RDICurrent) GetVariables() []data.Variable { return nil }
func (m *RDICurrent) Call(ctx data.Context) (data.GetValue, data.Control) {
	// PHP RecursiveDirectoryIterator::current() 返回 $this
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue, nil
	}
	return data.NewNullValue(), nil
}

// ---- key ----

type RDIKey struct{}

func (m *RDIKey) GetName() string               { return "key" }
func (m *RDIKey) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIKey) GetIsStatic() bool             { return false }
func (m *RDIKey) GetReturnType() data.Types     { return data.Mixed{} }
func (m *RDIKey) GetParams() []data.GetValue    { return nil }
func (m *RDIKey) GetVariables() []data.Variable { return nil }
func (m *RDIKey) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(rdiCurrentFile(cv)), nil
}

// ---- next ----

type RDINext struct{}

func (m *RDINext) GetName() string               { return "next" }
func (m *RDINext) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDINext) GetIsStatic() bool             { return false }
func (m *RDINext) GetReturnType() data.Types     { return nil }
func (m *RDINext) GetParams() []data.GetValue    { return nil }
func (m *RDINext) GetVariables() []data.Variable { return nil }
func (m *RDINext) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	rdiSetPos(cv, rdiGetPos(cv)+1)
	return nil, nil
}

// ---- valid ----

type RDIValid struct{}

func (m *RDIValid) GetName() string               { return "valid" }
func (m *RDIValid) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIValid) GetIsStatic() bool             { return false }
func (m *RDIValid) GetReturnType() data.Types     { return data.Bool{} }
func (m *RDIValid) GetParams() []data.GetValue    { return nil }
func (m *RDIValid) GetVariables() []data.Variable { return nil }
func (m *RDIValid) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	files := rdiGetFiles(cv)
	pos := rdiGetPos(cv)
	return data.NewBoolValue(pos >= 0 && pos < len(files)), nil
}

// ---- hasChildren ----

type RDIHasChildren struct{}

func (m *RDIHasChildren) GetName() string               { return "hasChildren" }
func (m *RDIHasChildren) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIHasChildren) GetIsStatic() bool             { return false }
func (m *RDIHasChildren) GetReturnType() data.Types     { return data.Bool{} }
func (m *RDIHasChildren) GetParams() []data.GetValue    { return nil }
func (m *RDIHasChildren) GetVariables() []data.Variable { return nil }
func (m *RDIHasChildren) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(currentPath)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.IsDir()), nil
}

// ---- getChildren：返回子目录的 RecursiveDirectoryIterator ----

type RDIGetChildren struct{}

func (m *RDIGetChildren) GetName() string               { return "getChildren" }
func (m *RDIGetChildren) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetChildren) GetIsStatic() bool             { return false }
func (m *RDIGetChildren) GetReturnType() data.Types     { return nil }
func (m *RDIGetChildren) GetParams() []data.GetValue    { return nil }
func (m *RDIGetChildren) GetVariables() []data.Variable { return nil }
func (m *RDIGetChildren) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		childClass := NewRecursiveDirectoryIteratorClass()
		return data.NewClassValue(childClass, ctx.CreateBaseContext()), nil
	}
	currentPath := rdiCurrentFile(cv)
	flags := rdiGetFlags(cv)

	childClass := NewRecursiveDirectoryIteratorClass()
	childCV := data.NewClassValue(childClass, ctx.CreateBaseContext())
	// 初始化属性
	childCV.SetProperty(rdiPathKey, data.NewStringValue(currentPath))
	childCV.SetProperty(rdiFilesKey, &rdiFilesValue{})
	childCV.SetProperty(rdiPosKey, data.NewIntValue(0))
	childCV.SetProperty(rdiFlagsKey, data.NewIntValue(flags))

	if currentPath != "" {
		rdiLoadFiles(childCV, currentPath, flags)
	}

	return childCV, nil
}

// ---- getFilename ----

type RDIGetFilename struct{}

func (m *RDIGetFilename) GetName() string               { return "getFilename" }
func (m *RDIGetFilename) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetFilename) GetIsStatic() bool             { return false }
func (m *RDIGetFilename) GetReturnType() data.Types     { return data.String{} }
func (m *RDIGetFilename) GetParams() []data.GetValue    { return nil }
func (m *RDIGetFilename) GetVariables() []data.Variable { return nil }
func (m *RDIGetFilename) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	currentPath := rdiCurrentFile(cv)
	_, filename := filepath.Split(currentPath)
	return data.NewStringValue(filename), nil
}

// ---- getPathname ----

type RDIGetPathname struct{}

func (m *RDIGetPathname) GetName() string               { return "getPathname" }
func (m *RDIGetPathname) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetPathname) GetIsStatic() bool             { return false }
func (m *RDIGetPathname) GetReturnType() data.Types     { return data.String{} }
func (m *RDIGetPathname) GetParams() []data.GetValue    { return nil }
func (m *RDIGetPathname) GetVariables() []data.Variable { return nil }
func (m *RDIGetPathname) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(rdiCurrentFile(cv)), nil
}

// ---- getPath（当前文件所在目录） ----

type RDIGetPath struct{}

func (m *RDIGetPath) GetName() string               { return "getPath" }
func (m *RDIGetPath) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetPath) GetIsStatic() bool             { return false }
func (m *RDIGetPath) GetReturnType() data.Types     { return data.String{} }
func (m *RDIGetPath) GetParams() []data.GetValue    { return nil }
func (m *RDIGetPath) GetVariables() []data.Variable { return nil }
func (m *RDIGetPath) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewStringValue(rdiGetPath(cv)), nil
	}
	return data.NewStringValue(filepath.Dir(currentPath)), nil
}

// ---- getRealPath ----

type RDIGetRealPath struct{}

func (m *RDIGetRealPath) GetName() string               { return "getRealPath" }
func (m *RDIGetRealPath) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetRealPath) GetIsStatic() bool             { return false }
func (m *RDIGetRealPath) GetReturnType() data.Types     { return data.String{} }
func (m *RDIGetRealPath) GetParams() []data.GetValue    { return nil }
func (m *RDIGetRealPath) GetVariables() []data.Variable { return nil }
func (m *RDIGetRealPath) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewStringValue(""), nil
	}
	realPath, err := filepath.EvalSymlinks(currentPath)
	if err != nil {
		if absPath, err2 := filepath.Abs(currentPath); err2 == nil {
			return data.NewStringValue(absPath), nil
		}
		return data.NewStringValue(currentPath), nil
	}
	return data.NewStringValue(realPath), nil
}

// ---- getSubPath：相对于根路径的子目录路径 ----

type RDIGetSubPath struct{}

func (m *RDIGetSubPath) GetName() string               { return "getSubPath" }
func (m *RDIGetSubPath) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetSubPath) GetIsStatic() bool             { return false }
func (m *RDIGetSubPath) GetReturnType() data.Types     { return data.String{} }
func (m *RDIGetSubPath) GetParams() []data.GetValue    { return nil }
func (m *RDIGetSubPath) GetVariables() []data.Variable { return nil }
func (m *RDIGetSubPath) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewStringValue(""), nil
	}
	rootPath := rdiGetPath(cv)
	dir := filepath.Dir(currentPath)
	if dir == rootPath {
		return data.NewStringValue(""), nil
	}
	rel, err := filepath.Rel(rootPath, dir)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(rel), nil
}

// ---- getSubPathname：相对于根路径的完整相对路径 ----

type RDIGetSubPathname struct{}

func (m *RDIGetSubPathname) GetName() string               { return "getSubPathname" }
func (m *RDIGetSubPathname) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIGetSubPathname) GetIsStatic() bool             { return false }
func (m *RDIGetSubPathname) GetReturnType() data.Types     { return data.String{} }
func (m *RDIGetSubPathname) GetParams() []data.GetValue    { return nil }
func (m *RDIGetSubPathname) GetVariables() []data.Variable { return nil }
func (m *RDIGetSubPathname) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewStringValue(""), nil
	}
	rootPath := rdiGetPath(cv)
	rel, err := filepath.Rel(rootPath, currentPath)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(rel), nil
}

// ---- isDir ----

type RDIIsDir struct{}

func (m *RDIIsDir) GetName() string               { return "isDir" }
func (m *RDIIsDir) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIIsDir) GetIsStatic() bool             { return false }
func (m *RDIIsDir) GetReturnType() data.Types     { return data.Bool{} }
func (m *RDIIsDir) GetParams() []data.GetValue    { return nil }
func (m *RDIIsDir) GetVariables() []data.Variable { return nil }
func (m *RDIIsDir) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(currentPath)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.IsDir()), nil
}

// ---- isFile ----

type RDIIsFile struct{}

func (m *RDIIsFile) GetName() string               { return "isFile" }
func (m *RDIIsFile) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RDIIsFile) GetIsStatic() bool             { return false }
func (m *RDIIsFile) GetReturnType() data.Types     { return data.Bool{} }
func (m *RDIIsFile) GetParams() []data.GetValue    { return nil }
func (m *RDIIsFile) GetVariables() []data.Variable { return nil }
func (m *RDIIsFile) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := rdiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	currentPath := rdiCurrentFile(cv)
	if currentPath == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(currentPath)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(!info.IsDir()), nil
}
