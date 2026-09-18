package httpfoundation

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	fqnFile                  = "Symfony\\Component\\HttpFoundation\\File\\File"
	fqnFileNotFoundException = "Symfony\\Component\\HttpFoundation\\File\\Exception\\FileNotFoundException"
	fqnFileException         = "Symfony\\Component\\HttpFoundation\\File\\Exception\\FileException"
	sfiPathnameKey           = "__sfi_pathname__"
)

// 无 MimeTypes 组件时的扩展名简表，避免走 vendor PHP 抛 LogicException。
var mimeByExtension = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"svg":  "image/svg+xml",
	"bmp":  "image/bmp",
	"ico":  "image/x-icon",
	"pdf":  "application/pdf",
	"txt":  "text/plain",
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"zip":  "application/zip",
	"csv":  "text/csv",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"mp4":  "video/mp4",
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"bin":  "application/octet-stream",
}

var extensionByMime map[string]string

func init() {
	extensionByMime = make(map[string]string, len(mimeByExtension))
	for ext, mime := range mimeByExtension {
		if _, ok := extensionByMime[mime]; !ok {
			extensionByMime[mime] = ext
		}
	}
}

// FileClass 实现 Symfony\Component\HttpFoundation\File\File。
type FileClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewFileClass() data.ClassStmt {
	c := &FileClass{}
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("path", 0, nil, nil),
				param("checkPath", 1, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				variable("path", 0, nil),
				variable("checkPath", 1, nil),
			},
			nil, fileConstruct),
		pubMethod("guessExtension", nil, nil, nil, fileGuessExtension),
		pubMethod("getMimeType", nil, nil, nil, fileGetMimeType),
		pubMethod("move",
			[]data.GetValue{
				param("directory", 0, nil, nil),
				param("name", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("directory", 0, nil),
				variable("name", 1, nil),
			},
			nil, fileMove),
		pubMethod("getContent", nil, nil, data.NewBaseType("string"), fileGetContent),
		pubMethod("getPathname", nil, nil, data.NewBaseType("string"), fileGetPathnameMethod),
		pubMethod("getFilename", nil, nil, data.NewBaseType("string"), fileGetFilenameMethod),
		pubMethod("getSize", nil, nil, data.NewBaseType("int"), fileGetSizeMethod),
		pubMethod("isFile", nil, nil, data.NewBaseType("bool"), fileIsFileMethod),
	}
	c.methods = indexMethods(list)
	c.methodList = list
	return c
}

func (c *FileClass) GetName() string { return fqnFile }
func (c *FileClass) GetExtend() *string {
	parent := "SplFileInfo"
	return &parent
}
func (c *FileClass) GetImplements() []string { return nil }
func (c *FileClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.GetPropertyList() {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *FileClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "pathname", "protected", false, data.NewStringValue("")),
	}
}
func (c *FileClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *FileClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *FileClass) GetMethod(name string) (data.Method, bool) {
	return getIndexedMethod(c.methods, name)
}
func (c *FileClass) GetMethods() []data.Method { return c.methodList }

func fileSetPathname(cv *data.ClassValue, path string) {
	if cv == nil {
		return
	}
	val := data.NewStringValue(path)
	_ = cv.SetProperty(sfiPathnameKey, val)
	_ = cv.SetProperty("pathname", val)
}

func fileGetPathname(cv *data.ClassValue) string {
	if cv == nil {
		return ""
	}
	if v, _ := cv.GetProperty(sfiPathnameKey); v != nil {
		if _, ok := v.(*data.NullValue); !ok {
			if s := v.AsString(); s != "" {
				return s
			}
		}
	}
	if v, _ := cv.GetProperty("pathname"); v != nil {
		if _, ok := v.(*data.NullValue); !ok {
			return v.AsString()
		}
	}
	return ""
}

func fileBaseName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	if i := strings.LastIndex(name, "/"); i >= 0 {
		return name[i+1:]
	}
	return name
}

func fileExtension(name string) string {
	base := fileBaseName(name)
	i := strings.LastIndex(base, ".")
	if i < 0 || i == len(base)-1 {
		return ""
	}
	return base[i+1:]
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func fileApplyPath(cv *data.ClassValue, path string, checkPath bool) data.Control {
	if checkPath && !isRegularFile(path) {
		return throwNamed(fqnFileNotFoundException, `The file "%s" does not exist`, path)
	}
	fileSetPathname(cv, path)
	return nil
}

func fileConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := bagClassValue(ctx)
	path := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		path = v.AsString()
	}
	checkPath := boolParam(ctx, 1, true)
	if ctl := fileApplyPath(cv, path, checkPath); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func mimeFromPath(path string) (string, bool) {
	ext := strings.ToLower(fileExtension(path))
	if ext == "" {
		return "", false
	}
	mime, ok := mimeByExtension[ext]
	return mime, ok
}

func fileGuessExtension(ctx data.Context) (data.GetValue, data.Control) {
	ext := strings.ToLower(fileExtension(fileGetPathname(bagClassValue(ctx))))
	if ext == "" {
		return data.NewNullValue(), nil
	}
	if _, ok := mimeByExtension[ext]; ok {
		return data.NewStringValue(ext), nil
	}
	return data.NewStringValue(ext), nil
}

func fileGetMimeType(ctx data.Context) (data.GetValue, data.Control) {
	if mime, ok := mimeFromPath(fileGetPathname(bagClassValue(ctx))); ok {
		return data.NewStringValue(mime), nil
	}
	return data.NewNullValue(), nil
}

func fileMove(ctx data.Context) (data.GetValue, data.Control) {
	cv := bagClassValue(ctx)
	dir := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		dir = v.AsString()
	}
	name, isNull, hasName := optionalStringParam(ctx, 1)
	useName := hasName && !isNull && name != ""
	return fileMoveTo(ctx, cv, dir, name, useName)
}

func fileMoveTo(ctx data.Context, cv *data.ClassValue, directory, name string, useName bool) (data.GetValue, data.Control) {
	src := fileGetPathname(cv)
	if err := os.MkdirAll(directory, 0o777); err != nil {
		info, statErr := os.Stat(directory)
		if statErr == nil && info != nil && !info.IsDir() {
			return nil, throwNamed(fqnFileException, `Unable to create the "%s" directory: a similarly-named file exists.`, directory)
		}
		return nil, throwNamed(fqnFileException, `Unable to create the "%s" directory.`, directory)
	}
	base := fileBaseName(src)
	if useName {
		base = fileBaseName(name)
	}
	directory = strings.TrimRight(directory, `/\`)
	target := directory + string(os.PathSeparator) + base
	if err := os.Rename(src, target); err != nil {
		return nil, throwNamed(fqnFileException, `Could not move the file "%s" to "%s" (%s).`, src, target, err.Error())
	}
	_ = os.Chmod(target, 0o666)
	return constructNamedClass(ctx, fqnFile, data.NewStringValue(target), data.NewBoolValue(false))
}

func fileGetContent(ctx data.Context) (data.GetValue, data.Control) {
	path := fileGetPathname(bagClassValue(ctx))
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, throwNamed(fqnFileException, `Could not get the content of the file "%s".`, path)
	}
	return data.NewStringValue(string(b)), nil
}

func fileGetPathnameMethod(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(fileGetPathname(bagClassValue(ctx))), nil
}

func fileGetFilenameMethod(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(fileBaseName(fileGetPathname(bagClassValue(ctx)))), nil
}

func fileGetSizeMethod(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(fileGetPathname(bagClassValue(ctx)))
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.Size())), nil
}

func fileIsFileMethod(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(isRegularFile(fileGetPathname(bagClassValue(ctx)))), nil
}

func mimeFromTypeOrPath(mimeType, path string) string {
	if mimeType != "" && mimeType != "application/octet-stream" {
		return mimeType
	}
	if mime, ok := mimeFromPath(path); ok {
		return mime
	}
	if mimeType != "" {
		return mimeType
	}
	return ""
}
