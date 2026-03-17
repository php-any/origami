package directory

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FilesystemIteratorGetFilenameMethod 实现 FilesystemIterator::getFilename
type FilesystemIteratorGetFilenameMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetFilenameMethod) GetName() string            { return "getFilename" }
func (m *FilesystemIteratorGetFilenameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilesystemIteratorGetFilenameMethod) GetIsStatic() bool          { return false }
func (m *FilesystemIteratorGetFilenameMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *FilesystemIteratorGetFilenameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *FilesystemIteratorGetFilenameMethod) GetReturnType() data.Types { return data.String{} }

func (m *FilesystemIteratorGetFilenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.GetFilename()), nil
}

// FilesystemIteratorGetBasenameMethod 实现 FilesystemIterator::getBasename
type FilesystemIteratorGetBasenameMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetBasenameMethod) GetName() string            { return "getBasename" }
func (m *FilesystemIteratorGetBasenameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilesystemIteratorGetBasenameMethod) GetIsStatic() bool          { return false }
func (m *FilesystemIteratorGetBasenameMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "suffix", 0, data.NewStringValue(""), data.String{}),
	}
}
func (m *FilesystemIteratorGetBasenameMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "suffix", 0, data.String{}),
	}
}
func (m *FilesystemIteratorGetBasenameMethod) GetReturnType() data.Types { return data.String{} }

func (m *FilesystemIteratorGetBasenameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	suffix := ""
	if sv, _ := ctx.GetIndexValue(0); sv != nil {
		suffix = sv.AsString()
	}
	return data.NewStringValue(m.instance.GetBasename(suffix)), nil
}

// FilesystemIteratorGetExtensionMethod 实现 FilesystemIterator::getExtension
type FilesystemIteratorGetExtensionMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetExtensionMethod) GetName() string { return "getExtension" }
func (m *FilesystemIteratorGetExtensionMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *FilesystemIteratorGetExtensionMethod) GetIsStatic() bool          { return false }
func (m *FilesystemIteratorGetExtensionMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *FilesystemIteratorGetExtensionMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *FilesystemIteratorGetExtensionMethod) GetReturnType() data.Types { return data.String{} }

func (m *FilesystemIteratorGetExtensionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.GetExtension()), nil
}

// FilesystemIteratorGetPathMethod 实现 FilesystemIterator::getPath
type FilesystemIteratorGetPathMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetPathMethod) GetName() string               { return "getPath" }
func (m *FilesystemIteratorGetPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorGetPathMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorGetPathMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorGetPathMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorGetPathMethod) GetReturnType() data.Types     { return data.String{} }

func (m *FilesystemIteratorGetPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.GetPath()), nil
}

// FilesystemIteratorGetPathnameMethod 实现 FilesystemIterator::getPathname
type FilesystemIteratorGetPathnameMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetPathnameMethod) GetName() string            { return "getPathname" }
func (m *FilesystemIteratorGetPathnameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilesystemIteratorGetPathnameMethod) GetIsStatic() bool          { return false }
func (m *FilesystemIteratorGetPathnameMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *FilesystemIteratorGetPathnameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *FilesystemIteratorGetPathnameMethod) GetReturnType() data.Types { return data.String{} }

func (m *FilesystemIteratorGetPathnameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.GetPathname()), nil
}

// FilesystemIteratorGetRealPathMethod 实现 FilesystemIterator::getRealPath
type FilesystemIteratorGetRealPathMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetRealPathMethod) GetName() string            { return "getRealPath" }
func (m *FilesystemIteratorGetRealPathMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilesystemIteratorGetRealPathMethod) GetIsStatic() bool          { return false }
func (m *FilesystemIteratorGetRealPathMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *FilesystemIteratorGetRealPathMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *FilesystemIteratorGetRealPathMethod) GetReturnType() data.Types { return data.String{} }

func (m *FilesystemIteratorGetRealPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.GetRealPath()), nil
}

// FilesystemIteratorIsDirMethod 实现 FilesystemIterator::isDir
type FilesystemIteratorIsDirMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorIsDirMethod) GetName() string               { return "isDir" }
func (m *FilesystemIteratorIsDirMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorIsDirMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorIsDirMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorIsDirMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorIsDirMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *FilesystemIteratorIsDirMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.IsDir()), nil
}

// FilesystemIteratorIsFileMethod 实现 FilesystemIterator::isFile
type FilesystemIteratorIsFileMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorIsFileMethod) GetName() string               { return "isFile" }
func (m *FilesystemIteratorIsFileMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorIsFileMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorIsFileMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorIsFileMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorIsFileMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *FilesystemIteratorIsFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.IsFile()), nil
}

// FilesystemIteratorIsDotMethod 实现 FilesystemIterator::isDot
type FilesystemIteratorIsDotMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorIsDotMethod) GetName() string               { return "isDot" }
func (m *FilesystemIteratorIsDotMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorIsDotMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorIsDotMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorIsDotMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorIsDotMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *FilesystemIteratorIsDotMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.IsDot()), nil
}

// FilesystemIteratorGetSizeMethod 实现 FilesystemIterator::getSize
type FilesystemIteratorGetSizeMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetSizeMethod) GetName() string               { return "getSize" }
func (m *FilesystemIteratorGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorGetSizeMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorGetSizeMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorGetSizeMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorGetSizeMethod) GetReturnType() data.Types     { return data.Int{} }

func (m *FilesystemIteratorGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(int(m.instance.GetSize())), nil
}

// FilesystemIteratorGetMTimeMethod 实现 FilesystemIterator::getMTime
type FilesystemIteratorGetMTimeMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetMTimeMethod) GetName() string               { return "getMTime" }
func (m *FilesystemIteratorGetMTimeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorGetMTimeMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorGetMTimeMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorGetMTimeMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorGetMTimeMethod) GetReturnType() data.Types     { return data.Int{} }

func (m *FilesystemIteratorGetMTimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(int(m.instance.GetMTime())), nil
}

// FilesystemIteratorIsReadableMethod 实现 FilesystemIterator::isReadable
type FilesystemIteratorIsReadableMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorIsReadableMethod) GetName() string               { return "isReadable" }
func (m *FilesystemIteratorIsReadableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorIsReadableMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorIsReadableMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorIsReadableMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorIsReadableMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *FilesystemIteratorIsReadableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.IsReadable()), nil
}

// FilesystemIteratorIsWritableMethod 实现 FilesystemIterator::isWritable
type FilesystemIteratorIsWritableMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorIsWritableMethod) GetName() string               { return "isWritable" }
func (m *FilesystemIteratorIsWritableMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorIsWritableMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorIsWritableMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorIsWritableMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorIsWritableMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *FilesystemIteratorIsWritableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.IsWritable()), nil
}

// FilesystemIteratorGetFlagsMethod 实现 FilesystemIterator::getFlags
type FilesystemIteratorGetFlagsMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorGetFlagsMethod) GetName() string               { return "getFlags" }
func (m *FilesystemIteratorGetFlagsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorGetFlagsMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorGetFlagsMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorGetFlagsMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorGetFlagsMethod) GetReturnType() data.Types     { return data.Int{} }

func (m *FilesystemIteratorGetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(m.instance.flags), nil
}

// FilesystemIteratorSetFlagsMethod 实现 FilesystemIterator::setFlags
type FilesystemIteratorSetFlagsMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorSetFlagsMethod) GetName() string            { return "setFlags" }
func (m *FilesystemIteratorSetFlagsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilesystemIteratorSetFlagsMethod) GetIsStatic() bool          { return false }
func (m *FilesystemIteratorSetFlagsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, nil, data.Int{}),
	}
}
func (m *FilesystemIteratorSetFlagsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, data.Int{}),
	}
}
func (m *FilesystemIteratorSetFlagsMethod) GetReturnType() data.Types { return nil }

func (m *FilesystemIteratorSetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if fv, _ := ctx.GetIndexValue(0); fv != nil {
		if iv, ok := fv.(interface{ AsInt() (int, error) }); ok {
			if n, err := iv.AsInt(); err == nil {
				m.instance.flags = n
			}
		}
	}
	return nil, nil
}
