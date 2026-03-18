package directory

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplFileInfo 内部状态属性名
const sfiPathnameKey = "__sfi_pathname__"

// SplFileInfoClass 提供 PHP SplFileInfo 类定义
// 状态存储在 ClassValue.ObjectValue.property，支持 PHP 子类继承（skill: php-class-state-sharing-pattern）
type SplFileInfoClass struct {
	node.Node
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
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(sfiPathnameKey, data.NewStringValue(""))
	return cv, nil
}

// ---- 辅助：从 ctx 获取 ClassValue ----
func sfiGetCV(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

// ---- 状态读写 ----
func sfiGetPathname(cv *data.ClassValue) string {
	v, _ := cv.ObjectValue.GetProperty(sfiPathnameKey)
	if sv, ok := v.(*data.StringValue); ok {
		return sv.Value
	}
	return ""
}

func sfiSetPathname(cv *data.ClassValue, pathname string) {
	cv.ObjectValue.SetProperty(sfiPathnameKey, data.NewStringValue(pathname))
}

// ---- 辅助 ----
func sfiGetFilename(cv *data.ClassValue) string {
	return filepath.Base(sfiGetPathname(cv))
}

func sfiGetPath(cv *data.ClassValue) string {
	return filepath.Dir(sfiGetPathname(cv))
}

func (c *SplFileInfoClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &SplFileInfoConstructMethod{}, true
	case "getFilename":
		return &SplFileInfoGetFilenameMethod{}, true
	case "getBasename":
		return &SplFileInfoGetBasenameMethod{}, true
	case "getExtension":
		return &SplFileInfoGetExtensionMethod{}, true
	case "getPath":
		return &SplFileInfoGetPathMethod{}, true
	case "getPathname":
		return &SplFileInfoGetPathnameMethod{}, true
	case "getRealPath":
		return &SplFileInfoGetRealPathMethod{}, true
	case "isDir":
		return &SplFileInfoIsDirMethod{}, true
	case "isFile":
		return &SplFileInfoIsFileMethod{}, true
	case "isLink":
		return &SplFileInfoIsLinkMethod{}, true
	case "getSize":
		return &SplFileInfoGetSizeMethod{}, true
	case "getMTime":
		return &SplFileInfoGetMTimeMethod{}, true
	case "isReadable":
		return &SplFileInfoIsReadableMethod{}, true
	case "isWritable":
		return &SplFileInfoIsWritableMethod{}, true
	case "__toString":
		return &SplFileInfoToStringMethod{}, true
	case "getFileInfo":
		return &SplFileInfoGetFileInfoMethod{}, true
	case "getPathInfo":
		return &SplFileInfoGetPathInfoMethod{}, true
	}
	return nil, false
}

func (c *SplFileInfoClass) GetMethods() []data.Method {
	return []data.Method{
		&SplFileInfoConstructMethod{},
		&SplFileInfoGetFilenameMethod{},
		&SplFileInfoGetBasenameMethod{},
		&SplFileInfoGetExtensionMethod{},
		&SplFileInfoGetPathMethod{},
		&SplFileInfoGetPathnameMethod{},
		&SplFileInfoGetRealPathMethod{},
		&SplFileInfoIsDirMethod{},
		&SplFileInfoIsFileMethod{},
		&SplFileInfoIsLinkMethod{},
		&SplFileInfoGetSizeMethod{},
		&SplFileInfoGetMTimeMethod{},
		&SplFileInfoIsReadableMethod{},
		&SplFileInfoIsWritableMethod{},
		&SplFileInfoToStringMethod{},
		&SplFileInfoGetFileInfoMethod{},
		&SplFileInfoGetPathInfoMethod{},
	}
}

func (c *SplFileInfoClass) GetConstruct() data.Method {
	return &SplFileInfoConstructMethod{}
}

// ------- __construct -------

type SplFileInfoConstructMethod struct{}

func (m *SplFileInfoConstructMethod) GetName() string            { return "__construct" }
func (m *SplFileInfoConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFileInfoConstructMethod) GetIsStatic() bool          { return false }
func (m *SplFileInfoConstructMethod) GetReturnType() data.Types  { return nil }
func (m *SplFileInfoConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "filename", 0, nil, data.NewBaseType("string"))}
}
func (m *SplFileInfoConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "filename", 0, data.NewBaseType("string"))}
}
func (m *SplFileInfoConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			sfiSetPathname(cv, s.AsString())
		}
	}
	return nil, nil
}

// ------- getFilename -------

type SplFileInfoGetFilenameMethod struct{}

func (m *SplFileInfoGetFilenameMethod) GetName() string               { return "getFilename" }
func (m *SplFileInfoGetFilenameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetFilenameMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetFilenameMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetFilenameMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetFilenameMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetFilenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	filename := sfiGetFilename(cv)
	return data.NewStringValue(filename), nil
}

// ------- getBasename -------

type SplFileInfoGetBasenameMethod struct{}

func (m *SplFileInfoGetBasenameMethod) GetName() string            { return "getBasename" }
func (m *SplFileInfoGetBasenameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplFileInfoGetBasenameMethod) GetIsStatic() bool          { return false }
func (m *SplFileInfoGetBasenameMethod) GetReturnType() data.Types  { return data.String{} }
func (m *SplFileInfoGetBasenameMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "suffix", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *SplFileInfoGetBasenameMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "suffix", 0, data.NewBaseType("string"))}
}
func (m *SplFileInfoGetBasenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	suffix := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			suffix = s.AsString()
		}
	}
	filename := sfiGetFilename(cv)
	if suffix != "" && len(filename) >= len(suffix) && filename[len(filename)-len(suffix):] == suffix {
		filename = filename[:len(filename)-len(suffix)]
	}
	return data.NewStringValue(filename), nil
}

// ------- getExtension -------

type SplFileInfoGetExtensionMethod struct{}

func (m *SplFileInfoGetExtensionMethod) GetName() string               { return "getExtension" }
func (m *SplFileInfoGetExtensionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetExtensionMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetExtensionMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetExtensionMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetExtensionMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetExtensionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	filename := sfiGetFilename(cv)
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

type SplFileInfoGetPathMethod struct{}

func (m *SplFileInfoGetPathMethod) GetName() string               { return "getPath" }
func (m *SplFileInfoGetPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetPathMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetPathMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetPathMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetPathMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(sfiGetPath(cv)), nil
}

// ------- getPathname -------

type SplFileInfoGetPathnameMethod struct{}

func (m *SplFileInfoGetPathnameMethod) GetName() string               { return "getPathname" }
func (m *SplFileInfoGetPathnameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetPathnameMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetPathnameMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoGetPathnameMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetPathnameMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetPathnameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(sfiGetPathname(cv)), nil
}

// ------- getRealPath -------

type SplFileInfoGetRealPathMethod struct{}

func (m *SplFileInfoGetRealPathMethod) GetName() string               { return "getRealPath" }
func (m *SplFileInfoGetRealPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetRealPathMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetRealPathMethod) GetReturnType() data.Types     { return nil }
func (m *SplFileInfoGetRealPathMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetRealPathMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetRealPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	pathname := sfiGetPathname(cv)
	if real, err := filepath.EvalSymlinks(pathname); err == nil {
		return data.NewStringValue(real), nil
	}
	if abs, err := filepath.Abs(pathname); err == nil {
		return data.NewStringValue(abs), nil
	}
	return data.NewBoolValue(false), nil
}

// ------- isDir -------

type SplFileInfoIsDirMethod struct{}

func (m *SplFileInfoIsDirMethod) GetName() string               { return "isDir" }
func (m *SplFileInfoIsDirMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsDirMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsDirMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsDirMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsDirMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsDirMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(sfiGetPathname(cv))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.IsDir()), nil
}

// ------- isFile -------

type SplFileInfoIsFileMethod struct{}

func (m *SplFileInfoIsFileMethod) GetName() string               { return "isFile" }
func (m *SplFileInfoIsFileMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsFileMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsFileMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsFileMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsFileMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(sfiGetPathname(cv))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(!info.IsDir()), nil
}

// ------- isLink -------

type SplFileInfoIsLinkMethod struct{}

func (m *SplFileInfoIsLinkMethod) GetName() string               { return "isLink" }
func (m *SplFileInfoIsLinkMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsLinkMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsLinkMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsLinkMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsLinkMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsLinkMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Lstat(sfiGetPathname(cv))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.Mode()&os.ModeSymlink != 0), nil
}

// ------- getSize -------

type SplFileInfoGetSizeMethod struct{}

func (m *SplFileInfoGetSizeMethod) GetName() string               { return "getSize" }
func (m *SplFileInfoGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetSizeMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetSizeMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFileInfoGetSizeMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetSizeMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	info, err := os.Stat(sfiGetPathname(cv))
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.Size())), nil
}

// ------- getMTime -------

type SplFileInfoGetMTimeMethod struct{}

func (m *SplFileInfoGetMTimeMethod) GetName() string               { return "getMTime" }
func (m *SplFileInfoGetMTimeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetMTimeMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetMTimeMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplFileInfoGetMTimeMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetMTimeMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetMTimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	info, err := os.Stat(sfiGetPathname(cv))
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.ModTime().Unix())), nil
}

// ------- isReadable -------

type SplFileInfoIsReadableMethod struct{}

func (m *SplFileInfoIsReadableMethod) GetName() string               { return "isReadable" }
func (m *SplFileInfoIsReadableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsReadableMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsReadableMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsReadableMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsReadableMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsReadableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	f, err := os.Open(sfiGetPathname(cv))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	f.Close()
	return data.NewBoolValue(true), nil
}

// ------- isWritable -------

type SplFileInfoIsWritableMethod struct{}

func (m *SplFileInfoIsWritableMethod) GetName() string               { return "isWritable" }
func (m *SplFileInfoIsWritableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoIsWritableMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoIsWritableMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplFileInfoIsWritableMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoIsWritableMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoIsWritableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(sfiGetPathname(cv))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	mode := info.Mode()
	return data.NewBoolValue(mode&0200 != 0 || mode&0020 != 0 || mode&0002 != 0), nil
}

// ------- __toString -------

type SplFileInfoToStringMethod struct{}

func (m *SplFileInfoToStringMethod) GetName() string               { return "__toString" }
func (m *SplFileInfoToStringMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoToStringMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoToStringMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SplFileInfoToStringMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoToStringMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoToStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(sfiGetPathname(cv)), nil
}

// ------- getFileInfo -------

type SplFileInfoGetFileInfoMethod struct{}

func (m *SplFileInfoGetFileInfoMethod) GetName() string               { return "getFileInfo" }
func (m *SplFileInfoGetFileInfoMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetFileInfoMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetFileInfoMethod) GetReturnType() data.Types     { return nil }
func (m *SplFileInfoGetFileInfoMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetFileInfoMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetFileInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		clone := NewSplFileInfoClass()
		return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
	}
	pathname := sfiGetPathname(cv)
	clone := NewSplFileInfoClass()
	cloneCV := data.NewClassValue(clone, ctx.CreateBaseContext())
	cloneCV.SetProperty(sfiPathnameKey, data.NewStringValue(pathname))
	return cloneCV, nil
}

// ------- getPathInfo -------

type SplFileInfoGetPathInfoMethod struct{}

func (m *SplFileInfoGetPathInfoMethod) GetName() string               { return "getPathInfo" }
func (m *SplFileInfoGetPathInfoMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplFileInfoGetPathInfoMethod) GetIsStatic() bool             { return false }
func (m *SplFileInfoGetPathInfoMethod) GetReturnType() data.Types     { return nil }
func (m *SplFileInfoGetPathInfoMethod) GetParams() []data.GetValue    { return nil }
func (m *SplFileInfoGetPathInfoMethod) GetVariables() []data.Variable { return nil }
func (m *SplFileInfoGetPathInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		clone := NewSplFileInfoClass()
		return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
	}
	parentPath := filepath.Dir(sfiGetPathname(cv))
	clone := NewSplFileInfoClass()
	cloneCV := data.NewClassValue(clone, ctx.CreateBaseContext())
	cloneCV.SetProperty(sfiPathnameKey, data.NewStringValue(parentPath))
	return cloneCV, nil
}
