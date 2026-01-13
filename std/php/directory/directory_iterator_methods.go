package directory

import (
	"github.com/php-any/origami/data"
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
