package http

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewUploadedFileClass() data.ClassStmt {
	return &UploadedFileClass{header: nil}
}

func NewUploadedFileClassFrom(header *multipart.FileHeader) data.ClassStmt {
	return &UploadedFileClass{header: header}
}

type UploadedFileClass struct {
	node.Node
	header *multipart.FileHeader
}

func (s *UploadedFileClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewUploadedFileClass(), ctx.CreateBaseContext()), nil
}

func (s *UploadedFileClass) GetName() string         { return uploadedFileTypeFQN }
func (s *UploadedFileClass) GetExtend() *string      { return nil }
func (s *UploadedFileClass) GetImplements() []string { return nil }
func (s *UploadedFileClass) AsString() string        { return "UploadedFile{}" }
func (s *UploadedFileClass) GetSource() any          { return s.header }
func (s *UploadedFileClass) GetConstruct() data.Method {
	return nil
}

func (s *UploadedFileClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "originalName":
		return &UploadedFileOriginalNameMethod{header: s.header}, true
	case "size":
		return &UploadedFileSizeMethod{header: s.header}, true
	case "mimeType":
		return &UploadedFileMimeTypeMethod{header: s.header}, true
	case "extension":
		return &UploadedFileExtensionMethod{header: s.header}, true
	case "isValid":
		return &UploadedFileIsValidMethod{header: s.header}, true
	case "getContent":
		return &UploadedFileGetContentMethod{header: s.header}, true
	case "store":
		return &UploadedFileStoreMethod{header: s.header}, true
	}
	return nil, false
}

func (s *UploadedFileClass) GetMethods() []data.Method {
	return []data.Method{
		&UploadedFileOriginalNameMethod{header: s.header},
		&UploadedFileSizeMethod{header: s.header},
		&UploadedFileMimeTypeMethod{header: s.header},
		&UploadedFileExtensionMethod{header: s.header},
		&UploadedFileIsValidMethod{header: s.header},
		&UploadedFileGetContentMethod{header: s.header},
		&UploadedFileStoreMethod{header: s.header},
	}
}

func (s *UploadedFileClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (s *UploadedFileClass) GetPropertyList() []data.Property {
	return nil
}

type UploadedFileOriginalNameMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileOriginalNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.header == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(h.header.Filename), nil
}

func (h *UploadedFileOriginalNameMethod) GetName() string               { return "originalName" }
func (h *UploadedFileOriginalNameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *UploadedFileOriginalNameMethod) GetIsStatic() bool             { return false }
func (h *UploadedFileOriginalNameMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *UploadedFileOriginalNameMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *UploadedFileOriginalNameMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

type UploadedFileSizeMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.header == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(h.header.Size)), nil
}

func (h *UploadedFileSizeMethod) GetName() string               { return "size" }
func (h *UploadedFileSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *UploadedFileSizeMethod) GetIsStatic() bool             { return false }
func (h *UploadedFileSizeMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *UploadedFileSizeMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *UploadedFileSizeMethod) GetReturnType() data.Types     { return data.NewBaseType("int") }

type UploadedFileMimeTypeMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileMimeTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.header == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(h.header.Header.Get("Content-Type")), nil
}

func (h *UploadedFileMimeTypeMethod) GetName() string               { return "mimeType" }
func (h *UploadedFileMimeTypeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *UploadedFileMimeTypeMethod) GetIsStatic() bool             { return false }
func (h *UploadedFileMimeTypeMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *UploadedFileMimeTypeMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *UploadedFileMimeTypeMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }

type UploadedFileExtensionMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileExtensionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.header == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(strings.TrimPrefix(filepath.Ext(h.header.Filename), ".")), nil
}

func (h *UploadedFileExtensionMethod) GetName() string               { return "extension" }
func (h *UploadedFileExtensionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *UploadedFileExtensionMethod) GetIsStatic() bool             { return false }
func (h *UploadedFileExtensionMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *UploadedFileExtensionMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *UploadedFileExtensionMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }

type UploadedFileIsValidMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileIsValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(h.header != nil && h.header.Size > 0 && h.header.Filename != ""), nil
}

func (h *UploadedFileIsValidMethod) GetName() string               { return "isValid" }
func (h *UploadedFileIsValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *UploadedFileIsValidMethod) GetIsStatic() bool             { return false }
func (h *UploadedFileIsValidMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *UploadedFileIsValidMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *UploadedFileIsValidMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }

type UploadedFileGetContentMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileGetContentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	content, err := readUploadedFileContent(h.header)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewStringValue(string(content)), nil
}

func (h *UploadedFileGetContentMethod) GetName() string            { return "getContent" }
func (h *UploadedFileGetContentMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *UploadedFileGetContentMethod) GetIsStatic() bool          { return false }
func (h *UploadedFileGetContentMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (h *UploadedFileGetContentMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *UploadedFileGetContentMethod) GetReturnType() data.Types { return data.NewBaseType("string") }

type UploadedFileStoreMethod struct {
	header *multipart.FileHeader
}

func (h *UploadedFileStoreMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	directory, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("store 缺少目录参数: %v", err)
	}
	name := ""
	if _, hasName := ctx.GetIndexValue(1); hasName {
		name, err = utils.ConvertFromIndex[string](ctx, 1)
		if err != nil {
			return nil, utils.NewThrowf("store 文件名参数无效: %v", err)
		}
	}
	path, err := storeUploadedFile(h.header, directory, name)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewStringValue(path), nil
}

func (h *UploadedFileStoreMethod) GetName() string            { return "store" }
func (h *UploadedFileStoreMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *UploadedFileStoreMethod) GetIsStatic() bool          { return false }
func (h *UploadedFileStoreMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, nil),
		node.NewParameter(nil, "name", 1, nil, nil),
	}
}
func (h *UploadedFileStoreMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, nil),
		node.NewVariable(nil, "name", 1, nil),
	}
}
func (h *UploadedFileStoreMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
