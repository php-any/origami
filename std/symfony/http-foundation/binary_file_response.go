package httpfoundation

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnBinaryFileResponse = "Symfony\\Component\\HttpFoundation\\BinaryFileResponse"

// BinaryFileResponseClass 实现 Symfony\Component\HttpFoundation\BinaryFileResponse。
type BinaryFileResponseClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewBinaryFileResponseClass() data.ClassStmt {
	c := &BinaryFileResponseClass{}
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("file", 0, nil, nil),
				param("status", 1, data.NewIntValue(200), nil),
				param("headers", 2, data.NewArrayValue(nil), nil),
				param("public", 3, data.NewBoolValue(true), nil),
				param("contentDisposition", 4, data.NewNullValue(), nil),
				param("autoEtag", 5, data.NewBoolValue(false), nil),
				param("autoLastModified", 6, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				variable("file", 0, nil),
				variable("status", 1, nil),
				variable("headers", 2, nil),
				variable("public", 3, nil),
				variable("contentDisposition", 4, nil),
				variable("autoEtag", 5, nil),
				variable("autoLastModified", 6, nil),
			},
			nil, binaryConstruct),
		pubMethod("setFile",
			[]data.GetValue{
				param("file", 0, nil, nil),
				param("contentDisposition", 1, data.NewNullValue(), nil),
				param("autoEtag", 2, data.NewBoolValue(false), nil),
				param("autoLastModified", 3, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				variable("file", 0, nil),
				variable("contentDisposition", 1, nil),
				variable("autoEtag", 2, nil),
				variable("autoLastModified", 3, nil),
			},
			nil, binarySetFile),
		pubMethod("getFile", nil, nil, nil, binaryGetFile),
		pubMethod("setContentDisposition",
			[]data.GetValue{
				param("disposition", 0, nil, nil),
				param("filename", 1, data.NewStringValue(""), nil),
				param("filenameFallback", 2, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("disposition", 0, nil),
				variable("filename", 1, nil),
				variable("filenameFallback", 2, nil),
			},
			nil, binarySetContentDisposition),
		pubMethod("prepare",
			[]data.GetValue{param("request", 0, nil, nil)},
			[]data.Variable{variable("request", 0, nil)},
			nil, binaryPrepare),
		pubMethod("sendContent", nil, nil, nil, binarySendContent),
		pubMethod("send",
			[]data.GetValue{param("flush", 0, data.NewBoolValue(true), nil)},
			[]data.Variable{variable("flush", 0, nil)},
			nil, binarySend),
		pubMethod("setContent",
			[]data.GetValue{param("content", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("content", 0, nil)},
			nil, binarySetContent),
		pubMethod("getContent", nil, nil, nil, binaryGetContent),
		pubMethod("deleteFileAfterSend",
			[]data.GetValue{param("shouldDelete", 0, data.NewBoolValue(true), nil)},
			[]data.Variable{variable("shouldDelete", 0, nil)},
			nil, binaryDeleteFileAfterSend),
	}
	c.methods = indexMethods(list)
	c.methodList = list
	return c
}

func (c *BinaryFileResponseClass) GetName() string { return fqnBinaryFileResponse }
func (c *BinaryFileResponseClass) GetExtend() *string {
	parent := fqnSymfonyResponse
	return &parent
}
func (c *BinaryFileResponseClass) GetImplements() []string { return nil }
func (c *BinaryFileResponseClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.GetPropertyList() {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *BinaryFileResponseClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "file", "protected", false, data.NewNullValue()),
		node.NewProperty(nil, "offset", "protected", false, data.NewIntValue(0)),
		node.NewProperty(nil, "maxlen", "protected", false, data.NewIntValue(-1)),
		node.NewProperty(nil, "deleteFileAfterSend", "protected", false, data.NewBoolValue(false)),
		node.NewProperty(nil, "chunkSize", "protected", false, data.NewIntValue(16*1024)),
	}
}
func (c *BinaryFileResponseClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *BinaryFileResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *BinaryFileResponseClass) GetMethod(name string) (data.Method, bool) {
	return getIndexedMethod(c.methods, name)
}
func (c *BinaryFileResponseClass) GetMethods() []data.Method { return c.methodList }

func binaryConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	fileVal, _ := ctx.GetIndexValue(0)
	status := intParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)
	public := boolParam(ctx, 3, true)
	disposition, dispNull, dispOK := optionalStringParam(ctx, 4)
	autoEtag := boolParam(ctx, 5, false)
	autoLM := boolParam(ctx, 6, true)

	headers := createResponseHeaders(ctx, headersVal)
	_ = cv.SetProperty("headers", headers)
	_ = cv.SetProperty("version", data.NewStringValue("1.0"))
	_ = cv.SetProperty("content", data.NewStringValue(""))
	if ctl := applyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("offset", data.NewIntValue(0))
	_ = cv.SetProperty("maxlen", data.NewIntValue(-1))
	_ = cv.SetProperty("deleteFileAfterSend", data.NewBoolValue(false))
	_ = cv.SetProperty("chunkSize", data.NewIntValue(16*1024))

	disp := ""
	if dispOK && !dispNull {
		disp = disposition
	}
	if _, ctl := binarySetFileWith(ctx, cv, fileVal, disp, autoEtag, autoLM); ctl != nil {
		return nil, ctl
	}
	if public {
		if _, ctl := symfonyResponseSetPublic(ctx); ctl != nil {
			return nil, ctl
		}
	}
	return data.NewNullValue(), nil
}

func binarySetFile(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	fileVal, _ := ctx.GetIndexValue(0)
	disposition, dispNull, dispOK := optionalStringParam(ctx, 1)
	autoEtag := boolParam(ctx, 2, false)
	autoLM := boolParam(ctx, 3, true)
	disp := ""
	if dispOK && !dispNull {
		disp = disposition
	}
	return binarySetFileWith(ctx, cv, fileVal, disp, autoEtag, autoLM)
}

func isSymfonyFileObject(v data.Value) bool {
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil {
		return false
	}
	name := cv.GetName()
	return name == fqnFile || name == fqnUploadedFile ||
		stringsHasSuffix(name, "\\File") || stringsHasSuffix(name, "\\UploadedFile")
}

func binarySetFileWith(ctx data.Context, cv *data.ClassValue, fileVal data.Value, disposition string, autoEtag, autoLM bool) (data.GetValue, data.Control) {
	var file *data.ClassValue
	if isSymfonyFileObject(fileVal) {
		file = fileVal.(*data.ClassValue)
	} else {
		path := ""
		if obj, ok := fileVal.(*data.ClassValue); ok {
			path = fileGetPathname(obj)
			if path == "" {
				if ret, ctl := callObjMethod(obj, "getPathname"); ctl != nil {
					return nil, ctl
				} else if val, ok := ret.(data.Value); ok && val != nil {
					path = val.AsString()
				}
			}
		} else if fileVal != nil && !isNull(fileVal) {
			path = fileVal.AsString()
		}
		inst, ctl := constructNamedClass(ctx, fqnFile, data.NewStringValue(path), data.NewBoolValue(true))
		if ctl != nil {
			return nil, ctl
		}
		file = inst
	}

	readable, ctl := callObjMethod(file, "isReadable")
	if ctl != nil {
		return nil, ctl
	}
	if !valueIsTruthy(readable) {
		return nil, throwNamed(fqnFileException, "File must be readable.")
	}

	_ = cv.SetProperty("file", file)

	if autoEtag {
		binaryApplyAutoEtag(cv, file)
	}
	if autoLM {
		binaryApplyAutoLastModified(cv, file)
	}
	if disposition != "" {
		if _, ctl := binaryApplyDisposition(cv, file, disposition, "", ""); ctl != nil {
			return nil, ctl
		}
	}
	return cv, nil
}

func binaryGetFile(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := cv.GetProperty("file")
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func binaryFileOf(cv *data.ClassValue) *data.ClassValue {
	if cv == nil {
		return nil
	}
	v, _ := cv.GetProperty("file")
	if f, ok := v.(*data.ClassValue); ok {
		return f
	}
	return nil
}

func binaryApplyAutoLastModified(cv *data.ClassValue, file *data.ClassValue) {
	path := fileGetPathname(file)
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	respHeaderSet(responseHeaders(cv), "Last-Modified", []string{respFormatHTTPDate(info.ModTime())}, true)
}

func binaryApplyAutoEtag(cv *data.ClassValue, file *data.ClassValue) {
	path := fileGetPathname(file)
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s-%d-%d", path, info.Size(), info.ModTime().Unix())))
	etag := `"` + base64.StdEncoding.EncodeToString(sum[:]) + `"`
	respHeaderSet(responseHeaders(cv), "ETag", []string{etag}, true)
}

func binarySetContentDisposition(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	disposition := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		disposition = v.AsString()
	}
	filename := ""
	if v, ok := ctx.GetIndexValue(1); ok && v != nil && !isNull(v) {
		filename = v.AsString()
	}
	fallback := ""
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		fallback = v.AsString()
	}
	return binaryApplyDisposition(cv, binaryFileOf(cv), disposition, filename, fallback)
}

func asciiFilenameFallback(filename string) string {
	var b strings.Builder
	for _, r := range filename {
		if r == '%' || r < 32 || r > 126 {
			b.WriteByte('_')
			continue
		}
		if !unicode.IsPrint(r) {
			b.WriteByte('_')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func binaryApplyDisposition(cv *data.ClassValue, file *data.ClassValue, disposition, filename, fallback string) (data.GetValue, data.Control) {
	if filename == "" && file != nil {
		filename = fileBaseName(fileGetPathname(file))
	}
	if fallback == "" {
		need := false
		for _, r := range filename {
			if r < 0x20 || r > 0x7e || r == '%' {
				need = true
				break
			}
		}
		if need {
			fallback = asciiFilenameFallback(filename)
		}
	}
	s, err := makeDisposition(disposition, filename, fallback)
	if err != nil {
		return nil, throwNamed("InvalidArgumentException", "%s", err.Error())
	}
	respHeaderSet(responseHeaders(cv), "Content-Disposition", []string{s}, true)
	return cv, nil
}

func binaryPrepare(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headers := responseHeaders(cv)
	file := binaryFileOf(cv)
	if !respHeaderHas(headers, "Content-Type") {
		mime := "application/octet-stream"
		if file != nil {
			if ret, ctl := callObjMethod(file, "getMimeType"); ctl != nil {
				return nil, ctl
			} else if val, ok := ret.(data.Value); ok && val != nil && !isNull(val) && val.AsString() != "" {
				mime = val.AsString()
			}
		}
		respHeaderSet(headers, "Content-Type", []string{mime}, true)
	}
	if _, ctl := symfonyResponsePrepare(ctx); ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("offset", data.NewIntValue(0))
	_ = cv.SetProperty("maxlen", data.NewIntValue(-1))
	if file != nil {
		if info, err := os.Stat(fileGetPathname(file)); err == nil {
			respHeaderRemove(headers, "Transfer-Encoding")
			respHeaderSet(headers, "Content-Length", []string{strconv.Itoa(int(info.Size()))}, true)
		}
	}
	return responseSelf(ctx), nil
}

func binarySend(ctx data.Context) (data.GetValue, data.Control) {
	if _, ctl := symfonyResponseSendHeaders(ctx); ctl != nil {
		return nil, ctl
	}
	return binarySendContent(ctx)
}

func binarySendContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	code := responseStatusCode(cv)
	if code < 200 || code >= 300 {
		return responseSelf(ctx), nil
	}
	maxlen := responseGetIntProp(cv, "maxlen", -1)
	if maxlen == 0 {
		return responseSelf(ctx), nil
	}
	file := binaryFileOf(cv)
	path := fileGetPathname(file)
	b, err := os.ReadFile(path)
	if err != nil {
		return responseSelf(ctx), nil
	}
	offset := responseGetIntProp(cv, "offset", 0)
	if offset > 0 {
		if offset >= len(b) {
			b = nil
		} else {
			b = b[offset:]
		}
	}
	if maxlen > 0 && maxlen < len(b) {
		b = b[:maxlen]
	}
	if ctl := data.EmitOutput(ctx, string(b)); ctl != nil {
		return nil, ctl
	}
	if cookiePropBool(cv, "deleteFileAfterSend", false) && isRegularFile(path) {
		_ = os.Remove(path)
	}
	return responseSelf(ctx), nil
}

func binarySetContent(ctx data.Context) (data.GetValue, data.Control) {
	content, _ := ctx.GetIndexValue(0)
	if content != nil && !isNull(content) {
		return nil, throwNamed("LogicException", "The content cannot be set on a BinaryFileResponse instance.")
	}
	return responseSelf(ctx), nil
}

func binaryGetContent(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

func binaryDeleteFileAfterSend(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	_ = cv.SetProperty("deleteFileAfterSend", data.NewBoolValue(boolParam(ctx, 0, true)))
	return responseSelf(ctx), nil
}
