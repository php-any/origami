package directory

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DirectoryIteratorGetFilenameMethod 实现 DirectoryIterator::getFilename
type DirectoryIteratorGetFilenameMethod struct{}

func (m *DirectoryIteratorGetFilenameMethod) GetName() string               { return "getFilename" }
func (m *DirectoryIteratorGetFilenameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorGetFilenameMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorGetFilenameMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorGetFilenameMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorGetFilenameMethod) GetReturnType() data.Types     { return data.String{} }

func (m *DirectoryIteratorGetFilenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(iterData.GetFilename()), nil
}

// DirectoryIteratorGetBasenameMethod 实现 DirectoryIterator::getBasename
type DirectoryIteratorGetBasenameMethod struct{}

func (m *DirectoryIteratorGetBasenameMethod) GetName() string            { return "getBasename" }
func (m *DirectoryIteratorGetBasenameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DirectoryIteratorGetBasenameMethod) GetIsStatic() bool          { return false }
func (m *DirectoryIteratorGetBasenameMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "suffix", 0, data.NewStringValue(""), data.String{}),
	}
}
func (m *DirectoryIteratorGetBasenameMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "suffix", 0, data.String{}),
	}
}
func (m *DirectoryIteratorGetBasenameMethod) GetReturnType() data.Types { return data.String{} }

func (m *DirectoryIteratorGetBasenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewStringValue(""), nil
	}

	// 获取可选的 suffix 参数
	suffixValue, _ := ctx.GetIndexValue(0)
	suffix := ""
	if suffixValue != nil {
		suffix = suffixValue.AsString()
	}

	return data.NewStringValue(iterData.GetBasename(suffix)), nil
}

// DirectoryIteratorGetExtensionMethod 实现 DirectoryIterator::getExtension
type DirectoryIteratorGetExtensionMethod struct{}

func (m *DirectoryIteratorGetExtensionMethod) GetName() string            { return "getExtension" }
func (m *DirectoryIteratorGetExtensionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DirectoryIteratorGetExtensionMethod) GetIsStatic() bool          { return false }
func (m *DirectoryIteratorGetExtensionMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *DirectoryIteratorGetExtensionMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *DirectoryIteratorGetExtensionMethod) GetReturnType() data.Types { return data.String{} }

func (m *DirectoryIteratorGetExtensionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(iterData.GetExtension()), nil
}

// DirectoryIteratorGetPathMethod 实现 DirectoryIterator::getPath
type DirectoryIteratorGetPathMethod struct{}

func (m *DirectoryIteratorGetPathMethod) GetName() string               { return "getPath" }
func (m *DirectoryIteratorGetPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorGetPathMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorGetPathMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorGetPathMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorGetPathMethod) GetReturnType() data.Types     { return data.String{} }

func (m *DirectoryIteratorGetPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(iterData.GetPath()), nil
}

// DirectoryIteratorGetPathnameMethod 实现 DirectoryIterator::getPathname
type DirectoryIteratorGetPathnameMethod struct{}

func (m *DirectoryIteratorGetPathnameMethod) GetName() string               { return "getPathname" }
func (m *DirectoryIteratorGetPathnameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorGetPathnameMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorGetPathnameMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorGetPathnameMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorGetPathnameMethod) GetReturnType() data.Types     { return data.String{} }

func (m *DirectoryIteratorGetPathnameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(iterData.GetPathname()), nil
}

// DirectoryIteratorIsDirMethod 实现 DirectoryIterator::isDir
type DirectoryIteratorIsDirMethod struct{}

func (m *DirectoryIteratorIsDirMethod) GetName() string               { return "isDir" }
func (m *DirectoryIteratorIsDirMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorIsDirMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorIsDirMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorIsDirMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorIsDirMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *DirectoryIteratorIsDirMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iterData.IsDir()), nil
}

// DirectoryIteratorIsFileMethod 实现 DirectoryIterator::isFile
type DirectoryIteratorIsFileMethod struct{}

func (m *DirectoryIteratorIsFileMethod) GetName() string               { return "isFile" }
func (m *DirectoryIteratorIsFileMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorIsFileMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorIsFileMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorIsFileMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorIsFileMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *DirectoryIteratorIsFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iterData.IsFile()), nil
}

// DirectoryIteratorIsDotMethod 实现 DirectoryIterator::isDot
type DirectoryIteratorIsDotMethod struct{}

func (m *DirectoryIteratorIsDotMethod) GetName() string               { return "isDot" }
func (m *DirectoryIteratorIsDotMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorIsDotMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorIsDotMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorIsDotMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorIsDotMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *DirectoryIteratorIsDotMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iterData.IsDot()), nil
}

// DirectoryIteratorGetSizeMethod 实现 DirectoryIterator::getSize
type DirectoryIteratorGetSizeMethod struct{}

func (m *DirectoryIteratorGetSizeMethod) GetName() string               { return "getSize" }
func (m *DirectoryIteratorGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorGetSizeMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorGetSizeMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorGetSizeMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorGetSizeMethod) GetReturnType() data.Types     { return data.Int{} }

func (m *DirectoryIteratorGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(iterData.GetSize())), nil
}

// DirectoryIteratorGetMTimeMethod 实现 DirectoryIterator::getMTime
type DirectoryIteratorGetMTimeMethod struct{}

func (m *DirectoryIteratorGetMTimeMethod) GetName() string               { return "getMTime" }
func (m *DirectoryIteratorGetMTimeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorGetMTimeMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorGetMTimeMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorGetMTimeMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorGetMTimeMethod) GetReturnType() data.Types     { return data.Int{} }

func (m *DirectoryIteratorGetMTimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(iterData.GetMTime())), nil
}

// DirectoryIteratorIsReadableMethod 实现 DirectoryIterator::isReadable
type DirectoryIteratorIsReadableMethod struct{}

func (m *DirectoryIteratorIsReadableMethod) GetName() string               { return "isReadable" }
func (m *DirectoryIteratorIsReadableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorIsReadableMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorIsReadableMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorIsReadableMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorIsReadableMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *DirectoryIteratorIsReadableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iterData.IsReadable()), nil
}

// DirectoryIteratorIsWritableMethod 实现 DirectoryIterator::isWritable
type DirectoryIteratorIsWritableMethod struct{}

func (m *DirectoryIteratorIsWritableMethod) GetName() string               { return "isWritable" }
func (m *DirectoryIteratorIsWritableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorIsWritableMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorIsWritableMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorIsWritableMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorIsWritableMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *DirectoryIteratorIsWritableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iterData.IsWritable()), nil
}
