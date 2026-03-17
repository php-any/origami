package directory

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplFileInfoClass 提供 PHP SplFileInfo 类定义
// SplFileInfo 是文件/目录信息的基类，DirectoryIterator 继承自它
type SplFileInfoClass struct {
	node.Node
	pathname string // 文件完整路径
}

// NewSplFileInfoClass 创建 SplFileInfoClass 实例
func NewSplFileInfoClass() *SplFileInfoClass {
	return &SplFileInfoClass{}
}

func (c *SplFileInfoClass) GetName() string { return "SplFileInfo" }

// GetExtend SplFileInfo 没有父类
func (c *SplFileInfoClass) GetExtend() *string { return nil }

// GetImplements SplFileInfo 实现 Stringable（PHP 8+）
func (c *SplFileInfoClass) GetImplements() []string {
	return []string{"Stringable"}
}

func (c *SplFileInfoClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplFileInfoClass) GetPropertyList() []data.Property              { return nil }

func (c *SplFileInfoClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	clone := &SplFileInfoClass{pathname: c.pathname}
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

func (c *SplFileInfoClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &SplFileInfoConstructMethod{instance: c}, true
	case "getFilename":
		return &SplFileInfoGetFilenameMethod{instance: c}, true
	case "getBasename":
		return &SplFileInfoGetBasenameMethod{instance: c}, true
	case "getExtension":
		return &SplFileInfoGetExtensionMethod{instance: c}, true
	case "getPath":
		return &SplFileInfoGetPathMethod{instance: c}, true
	case "getPathname":
		return &SplFileInfoGetPathnameMethod{instance: c}, true
	case "getRealPath":
		return &SplFileInfoGetRealPathMethod{instance: c}, true
	case "isDir":
		return &SplFileInfoIsDirMethod{instance: c}, true
	case "isFile":
		return &SplFileInfoIsFileMethod{instance: c}, true
	case "isLink":
		return &SplFileInfoIsLinkMethod{instance: c}, true
	case "getSize":
		return &SplFileInfoGetSizeMethod{instance: c}, true
	case "getMTime":
		return &SplFileInfoGetMTimeMethod{instance: c}, true
	case "isReadable":
		return &SplFileInfoIsReadableMethod{instance: c}, true
	case "isWritable":
		return &SplFileInfoIsWritableMethod{instance: c}, true
	case "__toString":
		return &SplFileInfoToStringMethod{instance: c}, true
	case "getFileInfo":
		return &SplFileInfoGetFileInfoMethod{instance: c}, true
	case "getPathInfo":
		return &SplFileInfoGetPathInfoMethod{instance: c}, true
	}
	return nil, false
}

func (c *SplFileInfoClass) GetMethods() []data.Method {
	return []data.Method{
		&SplFileInfoConstructMethod{instance: c},
		&SplFileInfoGetFilenameMethod{instance: c},
		&SplFileInfoGetBasenameMethod{instance: c},
		&SplFileInfoGetExtensionMethod{instance: c},
		&SplFileInfoGetPathMethod{instance: c},
		&SplFileInfoGetPathnameMethod{instance: c},
		&SplFileInfoGetRealPathMethod{instance: c},
		&SplFileInfoIsDirMethod{instance: c},
		&SplFileInfoIsFileMethod{instance: c},
		&SplFileInfoIsLinkMethod{instance: c},
		&SplFileInfoGetSizeMethod{instance: c},
		&SplFileInfoGetMTimeMethod{instance: c},
		&SplFileInfoIsReadableMethod{instance: c},
		&SplFileInfoIsWritableMethod{instance: c},
		&SplFileInfoToStringMethod{instance: c},
		&SplFileInfoGetFileInfoMethod{instance: c},
		&SplFileInfoGetPathInfoMethod{instance: c},
	}
}

func (c *SplFileInfoClass) GetConstruct() data.Method {
	return &SplFileInfoConstructMethod{instance: c}
}

// ------- 辅助方法 -------

func (c *SplFileInfoClass) getFilename() string {
	return filepath.Base(c.pathname)
}

func (c *SplFileInfoClass) getPath() string {
	return filepath.Dir(c.pathname)
}

// ------- __construct -------

type SplFileInfoConstructMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoConstructMethod) GetName() string            { return "__construct" }
func (m *SplFileInfoConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFileInfoConstructMethod) GetIsStatic() bool          { return false }
func (m *SplFileInfoConstructMethod) GetReturnType() data.Types  { return nil }
func (m *SplFileInfoConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "filename", 0, nil, data.NewBaseType("string"))}
}
func (m *SplFileInfoConstructMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			m.instance.pathname = s.AsString()
		}
	}
	return nil, nil
}

// ------- getFilename -------

type SplFileInfoGetFilenameMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetFilenameMethod) GetName() string               { return "getFilename" }
func (m *SplFileInfoGetFilenameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetFilenameMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetFilenameMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetFilenameMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetFilenameMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetFilenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.getFilename()), nil
}

// ------- getBasename -------

type SplFileInfoGetBasenameMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetBasenameMethod) GetName() string            { return "getBasename" }
func (m *SplFileInfoGetBasenameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFileInfoGetBasenameMethod) GetIsStatic() bool          { return false }
func (m *SplFileInfoGetBasenameMethod) GetReturnType() data.Types  { return data.String{} }
func (m *SplFileInfoGetBasenameMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "suffix", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *SplFileInfoGetBasenameMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetBasenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	suffix := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			suffix = s.AsString()
		}
	}
	filename := m.instance.getFilename()
	if suffix != "" && len(filename) >= len(suffix) && filename[len(filename)-len(suffix):] == suffix {
		filename = filename[:len(filename)-len(suffix)]
	}
	return data.NewStringValue(filename), nil
}

// ------- getExtension -------

type SplFileInfoGetExtensionMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetExtensionMethod) GetName() string               { return "getExtension" }
func (m *SplFileInfoGetExtensionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetExtensionMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetExtensionMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetExtensionMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetExtensionMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetExtensionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	filename := m.instance.getFilename()
	ext := ""
	for i := len(filename) - 1; i > 0; i-- {
		if filename[i] == '.' {
			ext = filename[i+1:]
			break
		}
	}
	return data.NewStringValue(ext), nil
}

// ------- getPath -------

type SplFileInfoGetPathMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetPathMethod) GetName() string               { return "getPath" }
func (m *SplFileInfoGetPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetPathMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetPathMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetPathMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetPathMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.getPath()), nil
}

// ------- getPathname -------

type SplFileInfoGetPathnameMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetPathnameMethod) GetName() string               { return "getPathname" }
func (m *SplFileInfoGetPathnameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetPathnameMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetPathnameMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetPathnameMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetPathnameMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetPathnameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.pathname), nil
}

// ------- getRealPath -------

type SplFileInfoGetRealPathMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetRealPathMethod) GetName() string               { return "getRealPath" }
func (m *SplFileInfoGetRealPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetRealPathMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetRealPathMethod) GetReturnType() data.Types     { return nil }
func (m *SplFileInfoGetRealPathMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetRealPathMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetRealPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if real, err := filepath.EvalSymlinks(m.instance.pathname); err == nil {
		return data.NewStringValue(real), nil
	}
	if abs, err := filepath.Abs(m.instance.pathname); err == nil {
		return data.NewStringValue(abs), nil
	}
	return data.NewBoolValue(false), nil
}

// ------- isDir -------

type SplFileInfoIsDirMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoIsDirMethod) GetName() string               { return "isDir" }
func (m *SplFileInfoIsDirMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsDirMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsDirMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsDirMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsDirMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsDirMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(m.instance.pathname)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.IsDir()), nil
}

// ------- isFile -------

type SplFileInfoIsFileMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoIsFileMethod) GetName() string               { return "isFile" }
func (m *SplFileInfoIsFileMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsFileMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsFileMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsFileMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsFileMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(m.instance.pathname)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(!info.IsDir()), nil
}

// ------- isLink -------

type SplFileInfoIsLinkMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoIsLinkMethod) GetName() string               { return "isLink" }
func (m *SplFileInfoIsLinkMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsLinkMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsLinkMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsLinkMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsLinkMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsLinkMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Lstat(m.instance.pathname)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.Mode()&os.ModeSymlink != 0), nil
}

// ------- getSize -------

type SplFileInfoGetSizeMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetSizeMethod) GetName() string               { return "getSize" }
func (m *SplFileInfoGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetSizeMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetSizeMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFileInfoGetSizeMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetSizeMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(m.instance.pathname)
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.Size())), nil
}

// ------- getMTime -------

type SplFileInfoGetMTimeMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetMTimeMethod) GetName() string               { return "getMTime" }
func (m *SplFileInfoGetMTimeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetMTimeMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetMTimeMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFileInfoGetMTimeMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetMTimeMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetMTimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(m.instance.pathname)
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.ModTime().Unix())), nil
}

// ------- isReadable -------

type SplFileInfoIsReadableMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoIsReadableMethod) GetName() string               { return "isReadable" }
func (m *SplFileInfoIsReadableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsReadableMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsReadableMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsReadableMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsReadableMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsReadableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	f, err := os.Open(m.instance.pathname)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	f.Close()
	return data.NewBoolValue(true), nil
}

// ------- isWritable -------

type SplFileInfoIsWritableMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoIsWritableMethod) GetName() string               { return "isWritable" }
func (m *SplFileInfoIsWritableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsWritableMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsWritableMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsWritableMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsWritableMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsWritableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(m.instance.pathname)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	mode := info.Mode()
	return data.NewBoolValue(mode&0200 != 0 || mode&0020 != 0 || mode&0002 != 0), nil
}

// ------- __toString -------

type SplFileInfoToStringMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoToStringMethod) GetName() string               { return "__toString" }
func (m *SplFileInfoToStringMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoToStringMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoToStringMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoToStringMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoToStringMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoToStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.pathname), nil
}

// ------- getFileInfo -------

type SplFileInfoGetFileInfoMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetFileInfoMethod) GetName() string               { return "getFileInfo" }
func (m *SplFileInfoGetFileInfoMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetFileInfoMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetFileInfoMethod) GetReturnType() data.Types     { return nil }
func (m *SplFileInfoGetFileInfoMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetFileInfoMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetFileInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	clone := &SplFileInfoClass{pathname: m.instance.pathname}
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

// ------- getPathInfo -------

type SplFileInfoGetPathInfoMethod struct{ instance *SplFileInfoClass }

func (m *SplFileInfoGetPathInfoMethod) GetName() string               { return "getPathInfo" }
func (m *SplFileInfoGetPathInfoMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetPathInfoMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetPathInfoMethod) GetReturnType() data.Types     { return nil }
func (m *SplFileInfoGetPathInfoMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetPathInfoMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetPathInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	parentPath := filepath.Dir(m.instance.pathname)
	clone := &SplFileInfoClass{pathname: parentPath}
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}
