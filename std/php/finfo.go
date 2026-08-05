package php

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// PHP fileinfo 常量（与 PHP 对齐的常用取值）
const (
	FILEINFO_NONE          = 0
	FILEINFO_SYMLINK       = 2
	FILEINFO_MIME_TYPE     = 16
	FILEINFO_MIME_ENCODING = 1024
	FILEINFO_MIME          = 1040 // FILEINFO_MIME_TYPE | FILEINFO_MIME_ENCODING
	FILEINFO_CONTINUE      = 32
	FILEINFO_RAW           = 256
	FILEINFO_EXTENSION     = 2097152
)

const finfoFlagsKey = "__finfo_flags__"

func InitFinfoConstants(vm data.VM) {
	vm.SetConstant("FILEINFO_NONE", data.NewIntValue(FILEINFO_NONE))
	vm.SetConstant("FILEINFO_SYMLINK", data.NewIntValue(FILEINFO_SYMLINK))
	vm.SetConstant("FILEINFO_MIME_TYPE", data.NewIntValue(FILEINFO_MIME_TYPE))
	vm.SetConstant("FILEINFO_MIME_ENCODING", data.NewIntValue(FILEINFO_MIME_ENCODING))
	vm.SetConstant("FILEINFO_MIME", data.NewIntValue(FILEINFO_MIME))
	vm.SetConstant("FILEINFO_CONTINUE", data.NewIntValue(FILEINFO_CONTINUE))
	vm.SetConstant("FILEINFO_RAW", data.NewIntValue(FILEINFO_RAW))
	vm.SetConstant("FILEINFO_EXTENSION", data.NewIntValue(FILEINFO_EXTENSION))
}

// FinfoClass 实现 PHP finfo 类（最小可用：构造 / file / buffer）
type FinfoClass struct {
	node.Node
}

func NewFinfoClass() data.ClassStmt {
	return &FinfoClass{}
}

func (c *FinfoClass) GetName() string                               { return "finfo" }
func (c *FinfoClass) GetExtend() *string                            { return nil }
func (c *FinfoClass) GetImplements() []string                       { return nil }
func (c *FinfoClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *FinfoClass) GetPropertyList() []data.Property              { return nil }
func (c *FinfoClass) GetConstruct() data.Method                     { return &FinfoConstructMethod{} }
func (c *FinfoClass) GetMethods() []data.Method {
	return []data.Method{
		&FinfoConstructMethod{},
		&FinfoFileMethod{},
		&FinfoBufferMethod{},
	}
}
func (c *FinfoClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &FinfoConstructMethod{}, true
	case "file":
		return &FinfoFileMethod{}, true
	case "buffer":
		return &FinfoBufferMethod{}, true
	}
	return nil, false
}
func (c *FinfoClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(finfoFlagsKey, data.NewIntValue(FILEINFO_NONE))
	return cv, nil
}

func finfoGetCV(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

func finfoGetFlags(cv *data.ClassValue) int {
	if cv == nil {
		return FILEINFO_NONE
	}
	v, _ := cv.ObjectValue.GetProperty(finfoFlagsKey)
	if asInt, ok := v.(data.AsInt); ok {
		if n, err := asInt.AsInt(); err == nil {
			return n
		}
	}
	return FILEINFO_NONE
}

func finfoSetFlags(cv *data.ClassValue, flags int) {
	if cv != nil {
		cv.SetProperty(finfoFlagsKey, data.NewIntValue(flags))
	}
}

type FinfoConstructMethod struct{}

func (m *FinfoConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	flags := FILEINFO_NONE
	if v, _ := ctx.GetIndexValue(0); v != nil {
		if asInt, ok := v.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				flags = n
			}
		}
	}
	finfoSetFlags(finfoGetCV(ctx), flags)
	return nil, nil
}
func (m *FinfoConstructMethod) GetName() string            { return "__construct" }
func (m *FinfoConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FinfoConstructMethod) GetIsStatic() bool          { return false }
func (m *FinfoConstructMethod) GetReturnType() data.Types  { return nil }
func (m *FinfoConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, data.NewIntValue(FILEINFO_NONE), nil),
		node.NewParameter(nil, "magic_database", 1, node.NewNullLiteral(nil), nil),
	}
}
func (m *FinfoConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, nil),
		node.NewVariable(nil, "magic_database", 1, nil),
	}
}

type FinfoFileMethod struct{}

func (m *FinfoFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	filename, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	flags := finfoGetFlags(finfoGetCV(ctx))
	if v, _ := ctx.GetIndexValue(1); v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			if asInt, ok := v.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil && n != 0 {
					flags = n
				}
			}
		}
	}
	mime, ok := detectMimeFromFile(filename)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(formatMimeByFlags(mime, flags)), nil
}
func (m *FinfoFileMethod) GetName() string            { return "file" }
func (m *FinfoFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FinfoFileMethod) GetIsStatic() bool          { return false }
func (m *FinfoFileMethod) GetReturnType() data.Types  { return nil }
func (m *FinfoFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "flags", 1, node.NewNullLiteral(nil), nil),
	}
}
func (m *FinfoFileMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, nil),
		node.NewVariable(nil, "flags", 1, nil),
	}
}

type FinfoBufferMethod struct{}

func (m *FinfoBufferMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	buf, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	flags := finfoGetFlags(finfoGetCV(ctx))
	if v, _ := ctx.GetIndexValue(1); v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			if asInt, ok := v.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil && n != 0 {
					flags = n
				}
			}
		}
	}
	mime := detectMimeFromBuffer([]byte(buf))
	return data.NewStringValue(formatMimeByFlags(mime, flags)), nil
}
func (m *FinfoBufferMethod) GetName() string            { return "buffer" }
func (m *FinfoBufferMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FinfoBufferMethod) GetIsStatic() bool          { return false }
func (m *FinfoBufferMethod) GetReturnType() data.Types  { return nil }
func (m *FinfoBufferMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "flags", 1, node.NewNullLiteral(nil), nil),
	}
}
func (m *FinfoBufferMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "flags", 1, nil),
	}
}

func detectMimeFromBuffer(b []byte) string {
	if len(b) == 0 {
		return "application/x-empty"
	}
	mime := http.DetectContentType(b)
	if i := strings.Index(mime, ";"); i >= 0 {
		mime = strings.TrimSpace(mime[:i])
	}
	return mime
}

func detectMimeFromFile(filename string) (string, bool) {
	f, err := os.Open(filename)
	if err != nil {
		return "", false
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		if info, statErr := f.Stat(); statErr == nil && info.Size() == 0 {
			return "application/x-empty", true
		}
		return "", false
	}
	mime := detectMimeFromBuffer(buf[:n])
	if mime == "application/octet-stream" || mime == "text/plain" {
		if extMime := mimeByExtension(filename); extMime != "" {
			return extMime, true
		}
	}
	return mime, true
}

func mimeByExtension(filename string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	switch ext {
	case "php":
		return "text/x-php"
	case "html", "htm":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "txt":
		return "text/plain"
	case "md":
		return "text/markdown"
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "svg":
		return "image/svg+xml"
	case "pdf":
		return "application/pdf"
	case "zip":
		return "application/zip"
	default:
		return ""
	}
}

func formatMimeByFlags(mime string, flags int) string {
	wantType := flags&FILEINFO_MIME_TYPE != 0 || flags&FILEINFO_MIME != 0 || flags == FILEINFO_NONE
	wantEnc := flags&FILEINFO_MIME_ENCODING != 0 || flags&FILEINFO_MIME == FILEINFO_MIME || (flags&FILEINFO_MIME != 0 && flags&FILEINFO_MIME_TYPE == 0 && flags&FILEINFO_MIME_ENCODING == 0)

	// FILEINFO_MIME = type | encoding
	if flags&FILEINFO_MIME == FILEINFO_MIME || flags == FILEINFO_MIME {
		return mime + "; charset=binary"
	}
	if flags&FILEINFO_MIME_ENCODING != 0 && flags&FILEINFO_MIME_TYPE == 0 {
		_ = wantEnc
		return "binary"
	}
	if wantType || flags&FILEINFO_MIME_TYPE != 0 {
		return mime
	}
	return mime
}

// ---- 过程式 API：finfo_open / finfo_file / finfo_buffer / finfo_close ----

type FinfoOpenFunction struct{}

func NewFinfoOpenFunction() data.FuncStmt { return &FinfoOpenFunction{} }

func (f *FinfoOpenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	flags := FILEINFO_NONE
	if v, _ := ctx.GetIndexValue(0); v != nil {
		if asInt, ok := v.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				flags = n
			}
		}
	}
	cls := NewFinfoClass()
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	finfoSetFlags(cv, flags)
	return cv, nil
}
func (f *FinfoOpenFunction) GetName() string { return "finfo_open" }
func (f *FinfoOpenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, data.NewIntValue(FILEINFO_NONE), nil),
		node.NewParameter(nil, "magic_database", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *FinfoOpenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, nil),
		node.NewVariable(nil, "magic_database", 1, nil),
	}
}

type FinfoFileFunction struct{}

func NewFinfoFileFunction() data.FuncStmt { return &FinfoFileFunction{} }

func (f *FinfoFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	finfoVal, _ := ctx.GetIndexValue(0)
	cv, ok := finfoVal.(*data.ClassValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	filename, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	flags := finfoGetFlags(cv)
	mime, ok := detectMimeFromFile(filename)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(formatMimeByFlags(mime, flags)), nil
}
func (f *FinfoFileFunction) GetName() string { return "finfo_file" }
func (f *FinfoFileFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "finfo", 0, nil, nil),
		node.NewParameter(nil, "filename", 1, nil, nil),
	}
}
func (f *FinfoFileFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "finfo", 0, nil),
		node.NewVariable(nil, "filename", 1, nil),
	}
}

type FinfoBufferFunction struct{}

func NewFinfoBufferFunction() data.FuncStmt { return &FinfoBufferFunction{} }

func (f *FinfoBufferFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	finfoVal, _ := ctx.GetIndexValue(0)
	cv, ok := finfoVal.(*data.ClassValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	buf, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	flags := finfoGetFlags(cv)
	mime := detectMimeFromBuffer([]byte(buf))
	return data.NewStringValue(formatMimeByFlags(mime, flags)), nil
}
func (f *FinfoBufferFunction) GetName() string { return "finfo_buffer" }
func (f *FinfoBufferFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "finfo", 0, nil, nil),
		node.NewParameter(nil, "string", 1, nil, nil),
	}
}
func (f *FinfoBufferFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "finfo", 0, nil),
		node.NewVariable(nil, "filename", 1, nil),
	}
}

type FinfoCloseFunction struct{}

func NewFinfoCloseFunction() data.FuncStmt { return &FinfoCloseFunction{} }

func (f *FinfoCloseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
func (f *FinfoCloseFunction) GetName() string { return "finfo_close" }
func (f *FinfoCloseFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "finfo", 0, nil, nil)}
}
func (f *FinfoCloseFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "finfo", 0, nil)}
}
