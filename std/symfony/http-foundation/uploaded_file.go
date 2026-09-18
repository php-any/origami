package httpfoundation

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnUploadedFile = "Symfony\\Component\\HttpFoundation\\File\\UploadedFile"

const (
	uploadErrOK        = 0
	uploadErrIniSize   = 1
	uploadErrFormSize  = 2
	uploadErrPartial   = 3
	uploadErrNoFile    = 4
	uploadErrNoTmpDir  = 6
	uploadErrCantWrite = 7
	uploadErrExtension = 8
)

// UploadedFileClass 实现 Symfony\Component\HttpFoundation\File\UploadedFile。
type UploadedFileClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewUploadedFileClass() data.ClassStmt {
	c := &UploadedFileClass{}
	maxFn := pubMethod("getMaxFilesize", nil, nil, nil, uploadedGetMaxFilesize).(*bagMethod)
	maxFn.static = true
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("path", 0, nil, nil),
				param("originalName", 1, nil, nil),
				param("mimeType", 2, data.NewNullValue(), nil),
				param("error", 3, data.NewNullValue(), nil),
				param("test", 4, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("path", 0, nil),
				variable("originalName", 1, nil),
				variable("mimeType", 2, nil),
				variable("error", 3, nil),
				variable("test", 4, nil),
			},
			nil, uploadedConstruct),
		pubMethod("getClientOriginalName", nil, nil, data.NewBaseType("string"), uploadedGetClientOriginalName),
		pubMethod("getClientOriginalPath", nil, nil, data.NewBaseType("string"), uploadedGetClientOriginalPath),
		pubMethod("getClientOriginalExtension", nil, nil, data.NewBaseType("string"), uploadedGetClientOriginalExtension),
		pubMethod("getClientMimeType", nil, nil, data.NewBaseType("string"), uploadedGetClientMimeType),
		pubMethod("guessClientExtension", nil, nil, nil, uploadedGuessClientExtension),
		pubMethod("getMimeType", nil, nil, nil, uploadedGetMimeType),
		pubMethod("getError", nil, nil, data.NewBaseType("int"), uploadedGetError),
		pubMethod("getErrorMessage", nil, nil, data.NewBaseType("string"), uploadedGetErrorMessage),
		pubMethod("isValid", nil, nil, data.NewBaseType("bool"), uploadedIsValid),
		pubMethod("move",
			[]data.GetValue{
				param("directory", 0, nil, nil),
				param("name", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("directory", 0, nil),
				variable("name", 1, nil),
			},
			nil, uploadedMove),
		pubMethod("getSize", nil, nil, data.NewBaseType("int"), fileGetSizeMethod),
		pubMethod("getPathname", nil, nil, data.NewBaseType("string"), fileGetPathnameMethod),
		pubMethod("getFilename", nil, nil, data.NewBaseType("string"), fileGetFilenameMethod),
		pubMethod("isFile", nil, nil, data.NewBaseType("bool"), fileIsFileMethod),
		maxFn,
	}
	c.methods = indexMethods(list)
	c.methodList = list
	return c
}

func (c *UploadedFileClass) GetName() string { return fqnUploadedFile }
func (c *UploadedFileClass) GetExtend() *string {
	parent := fqnFile
	return &parent
}
func (c *UploadedFileClass) GetImplements() []string { return nil }
func (c *UploadedFileClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.GetPropertyList() {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *UploadedFileClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "originalName", "private", false, data.NewStringValue("")),
		node.NewProperty(nil, "mimeType", "private", false, data.NewStringValue("application/octet-stream")),
		node.NewProperty(nil, "error", "private", false, data.NewIntValue(0)),
		node.NewProperty(nil, "originalPath", "private", false, data.NewStringValue("")),
		node.NewProperty(nil, "test", "private", false, data.NewBoolValue(false)),
		node.NewProperty(nil, "pathname", "protected", false, data.NewStringValue("")),
	}
}
func (c *UploadedFileClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *UploadedFileClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *UploadedFileClass) GetMethod(name string) (data.Method, bool) {
	return getIndexedMethod(c.methods, name)
}
func (c *UploadedFileClass) GetMethods() []data.Method { return c.methodList }
func (c *UploadedFileClass) GetStaticMethod(name string) (data.Method, bool) {
	if strings.EqualFold(name, "getMaxFilesize") {
		return getIndexedMethod(c.methods, "getMaxFilesize")
	}
	return nil, false
}

func uploadedConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := bagClassValue(ctx)
	path := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		path = v.AsString()
	}
	originalName := ""
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		originalName = v.AsString()
	}
	mimeType := "application/octet-stream"
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			if s := v.AsString(); s != "" {
				mimeType = s
			}
		}
	}
	errCode := 0
	if v, ok := ctx.GetIndexValue(3); ok && v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			errCode = valueAsInt(v, 0)
		}
	}
	test := boolParam(ctx, 4, false)

	_ = cv.SetProperty("originalName", data.NewStringValue(fileBaseName(originalName)))
	_ = cv.SetProperty("originalPath", data.NewStringValue(strings.ReplaceAll(originalName, "\\", "/")))
	_ = cv.SetProperty("mimeType", data.NewStringValue(mimeType))
	_ = cv.SetProperty("error", data.NewIntValue(errCode))
	_ = cv.SetProperty("test", data.NewBoolValue(test))

	if ctl := fileApplyPath(cv, path, errCode == uploadErrOK); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func uploadedStringProp(ctx data.Context, name, def string) string {
	cv := bagClassValue(ctx)
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	if _, ok := v.(*data.NullValue); ok {
		return def
	}
	return v.AsString()
}

func uploadedIntProp(ctx data.Context, name string, def int) int {
	cv := bagClassValue(ctx)
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	return valueAsInt(v, def)
}

func uploadedBoolProp(ctx data.Context, name string, def bool) bool {
	cv := bagClassValue(ctx)
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	if b, ok := v.(data.AsBool); ok {
		okv, err := b.AsBool()
		if err == nil {
			return okv
		}
	}
	return def
}

func uploadedGetClientOriginalName(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(uploadedStringProp(ctx, "originalName", "")), nil
}

func uploadedGetClientOriginalPath(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(uploadedStringProp(ctx, "originalPath", "")), nil
}

func uploadedGetClientOriginalExtension(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(fileExtension(uploadedStringProp(ctx, "originalName", ""))), nil
}

func uploadedGetClientMimeType(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(uploadedStringProp(ctx, "mimeType", "application/octet-stream")), nil
}

func uploadedGuessClientExtension(ctx data.Context) (data.GetValue, data.Control) {
	mime := strings.ToLower(uploadedStringProp(ctx, "mimeType", ""))
	if ext, ok := extensionByMime[mime]; ok {
		return data.NewStringValue(ext), nil
	}
	ext := strings.ToLower(fileExtension(uploadedStringProp(ctx, "originalName", "")))
	if ext == "" {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(ext), nil
}

func uploadedGetMimeType(ctx data.Context) (data.GetValue, data.Control) {
	path := fileGetPathname(bagClassValue(ctx))
	client := uploadedStringProp(ctx, "mimeType", "")
	if mime := mimeFromTypeOrPath(client, path); mime != "" {
		return data.NewStringValue(mime), nil
	}
	return data.NewNullValue(), nil
}

func uploadedGetError(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(uploadedIntProp(ctx, "error", 0)), nil
}

func uploadedIsValid(ctx data.Context) (data.GetValue, data.Control) {
	ok := uploadedIntProp(ctx, "error", 0) == uploadErrOK
	// test=true 时不要求 is_uploaded_file；Origami 无该内置时同样只看 error。
	return data.NewBoolValue(ok), nil
}

func uploadedMove(ctx data.Context) (data.GetValue, data.Control) {
	cv := bagClassValue(ctx)
	valid, ctl := uploadedIsValid(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if valueIsTruthy(valid) {
		dir := ""
		if v, ok := ctx.GetIndexValue(0); ok && v != nil {
			dir = v.AsString()
		}
		name, isNull, hasName := optionalStringParam(ctx, 1)
		useName := hasName && !isNull && name != ""
		return fileMoveTo(ctx, cv, dir, name, useName)
	}
	return nil, throwNamed(fqnFileException, "%s", uploadedExceptionMessage(ctx))
}

func uploadedGetErrorMessage(ctx data.Context) (data.GetValue, data.Control) {
	if uploadedIntProp(ctx, "error", 0) == uploadErrOK {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(uploadedExceptionMessage(ctx)), nil
}

func uploadedExceptionMessage(ctx data.Context) string {
	code := uploadedIntProp(ctx, "error", 0)
	name := uploadedStringProp(ctx, "originalName", "")
	switch code {
	case uploadErrIniSize:
		max := uploadedMaxFilesize(ctx) / 1024
		return fmt.Sprintf(`The file "%s" exceeds your upload_max_filesize ini directive (limit is %d KiB).`, name, max)
	case uploadErrFormSize:
		return fmt.Sprintf(`The file "%s" exceeds the upload limit defined in your form.`, name)
	case uploadErrPartial:
		return fmt.Sprintf(`The file "%s" was only partially uploaded.`, name)
	case uploadErrNoFile:
		return "No file was uploaded."
	case uploadErrCantWrite:
		return fmt.Sprintf(`The file "%s" could not be written on disk.`, name)
	case uploadErrNoTmpDir:
		return "File could not be uploaded: missing temporary directory."
	case uploadErrExtension:
		return "File upload was stopped by a PHP extension."
	default:
		return fmt.Sprintf(`The file "%s" was not uploaded due to an unknown error.`, name)
	}
}

func uploadedGetMaxFilesize(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(uploadedMaxFilesize(ctx)), nil
}

func uploadedMaxFilesize(ctx data.Context) int {
	postMax := parseIniFilesize(phpIniGetString(ctx, "post_max_size", "8M"))
	uploadMax := parseIniFilesize(phpIniGetString(ctx, "upload_max_filesize", "2M"))
	maxInt := int(^uint(0) >> 1)
	if postMax == 0 {
		postMax = maxInt
	}
	if uploadMax == 0 {
		uploadMax = maxInt
	}
	if postMax < uploadMax {
		return postMax
	}
	return uploadMax
}

func parseIniFilesize(size string) int {
	size = strings.TrimSpace(strings.ToLower(size))
	if size == "" {
		return 0
	}
	mult := 1
	switch size[len(size)-1] {
	case 'g':
		mult = 1024 * 1024 * 1024
		size = size[:len(size)-1]
	case 'm':
		mult = 1024 * 1024
		size = size[:len(size)-1]
	case 'k':
		mult = 1024
		size = size[:len(size)-1]
	}
	size = strings.TrimSpace(size)
	n := 0
	for _, r := range size {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n * mult
}

func phpIniGetString(ctx data.Context, key, def string) string {
	if ctx == nil || ctx.GetVM() == nil {
		return def
	}
	fn, ok := ctx.GetVM().GetFunc("ini_get")
	if !ok || fn == nil {
		return def
	}
	vars := fn.GetVariables()
	fnCtx := ctx.CreateContext(vars)
	if len(vars) > 0 {
		if ctl := fnCtx.SetVariableValue(vars[0], data.NewStringValue(key)); ctl != nil {
			return def
		}
	}
	ret, ctl := fn.Call(fnCtx)
	if ctl != nil || ret == nil {
		return def
	}
	val, ok := ret.(data.Value)
	if !ok || val == nil {
		return def
	}
	if _, isBool := val.(*data.BoolValue); isBool {
		return def
	}
	s := val.AsString()
	if s == "" {
		return def
	}
	return s
}
