package filesystem

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
	sffinder "github.com/php-any/origami/std/symfony/finder"
)

const filesystemClassName = "Illuminate\\Filesystem\\Filesystem"

type FilesystemClass struct {
	node.Node
	methods map[string]data.Method
}

func NewFilesystemClass() data.ClassStmt {
	c := &FilesystemClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *FilesystemClass) GetName() string    { return filesystemClassName }
func (c *FilesystemClass) GetExtend() *string { return nil }
func (c *FilesystemClass) GetImplements() []string {
	return nil
}
func (c *FilesystemClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *FilesystemClass) GetPropertyList() []data.Property         { return nil }
func (c *FilesystemClass) GetConstruct() data.Method                  { return nil }
func (c *FilesystemClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *FilesystemClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *FilesystemClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *FilesystemClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *FilesystemClass) register() {
	c.methods["exists"] = kit.InstanceMethod("exists", []string{"path"}, fsExists)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"path", "lock"}, 1, fsGet)
	c.methods["put"] = kit.InstanceMethodOpt("put", []string{"path", "contents", "lock"}, 2, fsPut)
	c.methods["replace"] = kit.InstanceMethodOpt("replace", []string{"path", "content", "mode"}, 2, fsReplace)
	c.methods["replaceinfile"] = kit.InstanceMethod("replaceInFile", []string{"search", "replace", "path"}, fsReplaceInFile)
	c.methods["delete"] = kit.InstanceMethod("delete", []string{"paths"}, fsDelete)
	c.methods["files"] = kit.InstanceMethodOpt("files", []string{"directory", "hidden"}, 1, fsFiles)
	c.methods["allfiles"] = kit.InstanceMethodOpt("allFiles", []string{"directory", "hidden"}, 1, fsAllFiles)
	c.methods["directories"] = kit.InstanceMethod("directories", []string{"directory"}, fsDirectories)
	c.methods["alldirectories"] = kit.InstanceMethod("allDirectories", []string{"directory"}, fsAllDirectories)
	c.methods["isdirectory"] = kit.InstanceMethod("isDirectory", []string{"directory"}, fsIsDirectory)
	c.methods["isfile"] = kit.InstanceMethod("isFile", []string{"file"}, fsIsFile)
	c.methods["isreadable"] = kit.InstanceMethod("isReadable", []string{"path"}, fsIsReadable)
	c.methods["iswritable"] = kit.InstanceMethod("isWritable", []string{"path"}, fsIsWritable)
	c.methods["makedirectory"] = kit.InstanceMethodOpt("makeDirectory", []string{"path", "mode", "recursive", "force"}, 1, fsMakeDirectory)
	c.methods["ensuredirectoryexists"] = kit.InstanceMethodOpt("ensureDirectoryExists", []string{"path", "mode", "recursive"}, 1, fsEnsureDirectoryExists)
	c.methods["deletedirectory"] = kit.InstanceMethodOpt("deleteDirectory", []string{"directory", "preserve"}, 1, fsDeleteDirectory)
	c.methods["cleandirectory"] = kit.InstanceMethod("cleanDirectory", []string{"directory"}, fsCleanDirectory)
	c.methods["copy"] = kit.InstanceMethod("copy", []string{"path", "target"}, fsCopy)
	c.methods["move"] = kit.InstanceMethod("move", []string{"path", "target"}, fsMove)
	c.methods["size"] = kit.InstanceMethod("size", []string{"path"}, fsSize)
	c.methods["lastmodified"] = kit.InstanceMethod("lastModified", []string{"path"}, fsLastModified)
	c.methods["hash"] = kit.InstanceMethodOpt("hash", []string{"path", "algorithm"}, 1, fsHash)
	c.methods["sharedget"] = kit.InstanceMethod("sharedGet", []string{"path"}, fsSharedGet)
	c.methods["getrequire"] = kit.InstanceMethodOpt("getRequire", []string{"path", "data"}, 1, fsGetRequire)
	c.methods["requireonce"] = kit.InstanceMethodOpt("requireOnce", []string{"path", "data"}, 1, fsRequireOnce)
	// getRequire/requireOnce：第二参 $data 经 extract 注入独立作用域（对齐 PhpEngine）
	c.methods["missing"] = kit.InstanceMethod("missing", []string{"path"}, fsMissing)
	c.methods["json"] = kit.InstanceMethodOpt("json", []string{"path", "flags", "lock"}, 1, fsJSON)
	c.methods["append"] = kit.InstanceMethod("append", []string{"path", "data"}, fsAppend)
	c.methods["prepend"] = kit.InstanceMethod("prepend", []string{"path", "data"}, fsPrepend)
	c.methods["name"] = kit.InstanceMethod("name", []string{"path"}, fsName)
	c.methods["basename"] = kit.InstanceMethod("basename", []string{"path"}, fsBasename)
	c.methods["dirname"] = kit.InstanceMethod("dirname", []string{"path"}, fsDirname)
	c.methods["extension"] = kit.InstanceMethod("extension", []string{"path"}, fsExtension)
	c.methods["type"] = kit.InstanceMethod("type", []string{"path"}, fsType)
	c.methods["glob"] = kit.InstanceMethod("glob", []string{"pattern"}, fsGlob)
	kit.RegisterMacroable(c.methods, filesystemClassName)
}

func fsPath(ctx data.Context, i int) string {
	v := kit.Unwrap(kit.Arg(ctx, i))
	if v == nil {
		return ""
	}
	return v.AsString()
}

func fsBoolArg(ctx data.Context, i int, def bool) bool {
	v := kit.Unwrap(kit.Arg(ctx, i))
	if v == nil || kit.IsNull(v) {
		return def
	}
	if b, ok := v.(*data.BoolValue); ok {
		return b.Value
	}
	return kit.Truthy(v)
}

func fsIntArg(ctx data.Context, i int, def int) int {
	v := kit.Unwrap(kit.Arg(ctx, i))
	if v == nil || kit.IsNull(v) {
		return def
	}
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	if fv, ok := v.(*data.FloatValue); ok {
		return int(fv.Value)
	}
	s := strings.TrimSpace(v.AsString())
	if s == "" {
		return def
	}
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func fsFileNotFound(path string) data.Control {
	return data.NewErrorThrow(nil, fmt.Errorf("File does not exist at path %s.", path))
}

func fsExists(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	_, err := os.Stat(path)
	return data.NewBoolValue(err == nil), nil
}

func fsIsFile(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(path)
	return data.NewBoolValue(err == nil && !info.IsDir()), nil
}

func fsIsDirectory(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(path)
	return data.NewBoolValue(err == nil && info.IsDir()), nil
}

func fsGet(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	lock := fsBoolArg(ctx, 1, false)
	if !fileExists(path) {
		return nil, fsFileNotFound(path)
	}
	if lock {
		return fsSharedGet(ctx)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewStringValue(string(b)), nil
}

func fsSharedGet(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if !fileExists(path) {
		return data.NewStringValue(""), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(b)), nil
}

func fsPut(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	contents := fsPath(ctx, 1)
	lock := fsBoolArg(ctx, 2, false)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	_ = lock
	parent := filepath.Dir(path)
	if parent != "" && parent != "." {
		_ = os.MkdirAll(parent, 0o755)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(len(contents)), nil
}

func fsDeletePaths(ctx data.Context) []string {
	first := kit.Unwrap(kit.Arg(ctx, 0))
	if av, ok := first.(*data.ArrayValue); ok && av != nil {
		var paths []string
		for _, e := range av.List {
			if e == nil || e.Value == nil {
				continue
			}
			paths = append(paths, kit.Unwrap(e.Value).AsString())
		}
		return paths
	}
	paths := make([]string, 0, 4)
	for i := 0; ; i++ {
		v, _ := ctx.GetIndexValue(i)
		if v == nil {
			break
		}
		paths = append(paths, kit.Unwrap(v).AsString())
	}
	return paths
}

func fsDelete(ctx data.Context) (data.GetValue, data.Control) {
	success := true
	for _, path := range fsDeletePaths(ctx) {
		if path == "" {
			success = false
			continue
		}
		if err := os.Remove(path); err != nil {
			success = false
		}
	}
	return data.NewBoolValue(success), nil
}

func fsFiles(ctx data.Context) (data.GetValue, data.Control) {
	dir := fsPath(ctx, 0)
	hidden := fsBoolArg(ctx, 1, false)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if dir == "" || !dirExists(dir) {
		return out, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out, nil
	}
	i := 0
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !hidden && strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(dir, name)
		out.SetIntKey(i, fsSplFileInfo(ctx, dir, full))
		i++
	}
	return out, nil
}

func fsSplFileInfo(ctx data.Context, baseDir, fullPath string) data.Value {
	rel, err := filepath.Rel(baseDir, fullPath)
	if err != nil {
		rel = filepath.Base(fullPath)
	}
	relDir := filepath.Dir(rel)
	if relDir == "." {
		relDir = ""
	}
	return sffinder.NewSplFileInfoInstance(ctx, fullPath, relDir, rel)
}

func fsMakeDirectory(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	mode := fsIntArg(ctx, 1, 0o755)
	recursive := fsBoolArg(ctx, 2, false)
	force := fsBoolArg(ctx, 3, false)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	var err error
	switch {
	case recursive:
		err = os.MkdirAll(path, os.FileMode(mode))
	case force:
		err = os.MkdirAll(path, os.FileMode(mode))
	default:
		err = os.Mkdir(path, os.FileMode(mode))
	}
	return data.NewBoolValue(err == nil), nil
}

func fsCopy(ctx data.Context) (data.GetValue, data.Control) {
	src := fsPath(ctx, 0)
	dst := fsPath(ctx, 1)
	if src == "" || dst == "" {
		return data.NewBoolValue(false), nil
	}
	in, err := os.Open(src)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return data.NewBoolValue(false), nil
	}
	out, err := os.Create(dst)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return data.NewBoolValue(err == nil), nil
}

func fsMove(ctx data.Context) (data.GetValue, data.Control) {
	src := fsPath(ctx, 0)
	dst := fsPath(ctx, 1)
	if src == "" || dst == "" {
		return data.NewBoolValue(false), nil
	}
	err := os.Rename(src, dst)
	return data.NewBoolValue(err == nil), nil
}

func fsSize(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	info, err := os.Stat(path)
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.Size())), nil
}

func fsLastModified(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	info, err := os.Stat(path)
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(info.ModTime().Unix())), nil
}

func fsHash(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	algo := fsPath(ctx, 1)
	if algo == "" {
		algo = "md5"
	}
	if algo == "md5" {
		f, err := os.Open(path)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		defer f.Close()
		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(fmt.Sprintf("%x", h.Sum(nil))), nil
	}
	sum, err := hashFile(algo, path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(sum), nil
}

func hashFile(algo, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_ = algo
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// fsRequireOnce / fsGetRequire 对齐 PHP：
//
//	(static function () use ($path, $data) {
//	    extract($data, EXTR_SKIP);
//	    return require[_once] $path;
//	})();
//
// 必须把 $data（含视图 $layout/$content/__env）注入独立作用域，再用
// LoadInCallerContext 引入文件；否则 Blade/PhpEngine 读到 null 的 $layout。
func fsRequireOnce(ctx data.Context) (data.GetValue, data.Control) {
	return fsRequireWithData(ctx, true)
}

func fsGetRequire(ctx data.Context) (data.GetValue, data.Control) {
	return fsRequireWithData(ctx, false)
}

func fsRequireWithData(ctx data.Context, once bool) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if !fileExists(path) {
		return nil, fsFileNotFound(path)
	}
	scope := ctx.CreateContext(nil)
	if dataArg := kit.Unwrap(kit.Arg(ctx, 1)); dataArg != nil && !kit.IsNull(dataArg) {
		fsExtractSkip(scope, dataArg)
	}
	return node.IncludeCore(scope, data.NewStringValue(path), once, true, nil)
}

func fsExtractSkip(ctx data.Context, arr data.Value) {
	switch v := arr.(type) {
	case *data.ObjectValue:
		v.RangeProperties(func(key string, val data.Value) bool {
			if key != "" && !ctx.HasVariableByName(key) {
				ctx.SetVariableByName(key, val)
			}
			return true
		})
	case *data.ArrayValue:
		for _, zv := range v.List {
			if zv == nil || zv.Name == "" {
				continue
			}
			if !ctx.HasVariableByName(zv.Name) {
				ctx.SetVariableByName(zv.Name, zv.Value)
			}
		}
	}
}

func fsMissing(ctx data.Context) (data.GetValue, data.Control) {
	ex, ctl := fsExists(ctx)
	if ctl != nil {
		return nil, ctl
	}
	b, _ := ex.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func fsJSON(ctx data.Context) (data.GetValue, data.Control) {
	raw, ctl := fsGet(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s := ""
	if v, ok := raw.(data.Value); ok && v != nil {
		s = v.AsString()
	}
	// 轻量：交给 PHP json_decode
	fn, ok := ctx.GetVM().GetFunc("json_decode")
	if !ok || fn == nil {
		return data.NewNullValue(), nil
	}
	nctx := ctx.CreateContext(fn.GetVariables())
	data.BindDeclaredArgs(nctx, fn, []data.Value{data.NewStringValue(s), data.NewBoolValue(true)})
	return fn.Call(nctx)
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fsReplace(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	content := ""
	if v := kit.Unwrap(kit.Arg(ctx, 1)); v != nil {
		content = v.AsString()
	}
	mode := fsIntArg(ctx, 2, -1)
	if path == "" {
		return data.NewNullValue(), nil
	}
	if real, err := filepath.EvalSymlinks(path); err == nil {
		path = real
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".*")
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return nil, data.NewErrorThrow(nil, err)
	}
	_ = tmp.Close()
	if mode >= 0 {
		_ = os.Chmod(tmpName, os.FileMode(mode))
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewNullValue(), nil
}

func fsReplaceInFile(ctx data.Context) (data.GetValue, data.Control) {
	search := kit.Unwrap(kit.Arg(ctx, 0))
	replace := kit.Unwrap(kit.Arg(ctx, 1))
	path := fsPath(ctx, 2)
	if !fileExists(path) {
		return data.NewBoolValue(false), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	content := string(b)
	switch {
	case search == nil:
		return data.NewBoolValue(false), nil
	case replace == nil:
		replace = data.NewStringValue("")
	}
	if sav, ok := search.(*data.ArrayValue); ok {
		searches := sav.ToValueList()
		replaces := []data.Value{}
		if rav, ok := replace.(*data.ArrayValue); ok {
			replaces = rav.ToValueList()
		}
		for i, s := range searches {
			r := ""
			if i < len(replaces) && replaces[i] != nil {
				r = replaces[i].AsString()
			} else if _, isArr := replace.(*data.ArrayValue); !isArr && replace != nil {
				r = replace.AsString()
			}
			content = strings.ReplaceAll(content, s.AsString(), r)
		}
	} else {
		content = strings.ReplaceAll(content, search.AsString(), replace.AsString())
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func fsListDirs(dir string, recursive bool) *data.ArrayValue {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if dir == "" || !dirExists(dir) {
		return out
	}
	i := 0
	if !recursive {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return out
		}
		for _, ent := range entries {
			if ent.IsDir() {
				out.SetIntKey(i, data.NewStringValue(filepath.Join(dir, ent.Name())))
				i++
			}
		}
		return out
	}
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || path == dir {
			return nil
		}
		if d.IsDir() {
			out.SetIntKey(i, data.NewStringValue(path))
			i++
		}
		return nil
	})
	return out
}

func fsDirectories(ctx data.Context) (data.GetValue, data.Control) {
	return fsListDirs(fsPath(ctx, 0), false), nil
}

func fsAllDirectories(ctx data.Context) (data.GetValue, data.Control) {
	return fsListDirs(fsPath(ctx, 0), true), nil
}

func fsAllFiles(ctx data.Context) (data.GetValue, data.Control) {
	dir := fsPath(ctx, 0)
	hidden := fsBoolArg(ctx, 1, false)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if dir == "" || !dirExists(dir) {
		return out, nil
	}
	i := 0
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if !hidden && strings.HasPrefix(name, ".") {
			if d.IsDir() && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		out.SetIntKey(i, fsSplFileInfo(ctx, dir, path))
		i++
		return nil
	})
	return out, nil
}

func fsIsReadable(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	f, err := os.Open(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	_ = f.Close()
	return data.NewBoolValue(true), nil
}

func fsIsWritable(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if info.IsDir() {
		tmp := filepath.Join(path, ".origami_write_probe")
		if err := os.WriteFile(tmp, []byte{}, 0o644); err != nil {
			return data.NewBoolValue(false), nil
		}
		_ = os.Remove(tmp)
		return data.NewBoolValue(true), nil
	}
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	_ = f.Close()
	return data.NewBoolValue(true), nil
}

func fsEnsureDirectoryExists(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	mode := fsIntArg(ctx, 1, 0o755)
	recursive := fsBoolArg(ctx, 2, true)
	if path == "" || dirExists(path) {
		return data.NewNullValue(), nil
	}
	var err error
	if recursive {
		err = os.MkdirAll(path, os.FileMode(mode))
	} else {
		err = os.Mkdir(path, os.FileMode(mode))
	}
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewNullValue(), nil
}

func fsDeleteDirectory(ctx data.Context) (data.GetValue, data.Control) {
	dir := fsPath(ctx, 0)
	preserve := fsBoolArg(ctx, 1, false)
	if dir == "" || !dirExists(dir) {
		return data.NewBoolValue(false), nil
	}
	if preserve {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		for _, ent := range entries {
			p := filepath.Join(dir, ent.Name())
			if ent.IsDir() {
				_ = os.RemoveAll(p)
			} else {
				_ = os.Remove(p)
			}
		}
		return data.NewBoolValue(true), nil
	}
	if err := os.RemoveAll(dir); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func fsCleanDirectory(ctx data.Context) (data.GetValue, data.Control) {
	dir := fsPath(ctx, 0)
	if dir == "" || !dirExists(dir) {
		return data.NewBoolValue(false), nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	for _, ent := range entries {
		p := filepath.Join(dir, ent.Name())
		if ent.IsDir() {
			_ = os.RemoveAll(p)
		} else {
			_ = os.Remove(p)
		}
	}
	return data.NewBoolValue(true), nil
}

func fsAppend(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	dataStr := ""
	if v := kit.Unwrap(kit.Arg(ctx, 1)); v != nil {
		dataStr = v.AsString()
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	defer f.Close()
	n, err := f.WriteString(dataStr)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(n), nil
}

func fsPrepend(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	dataStr := ""
	if v := kit.Unwrap(kit.Arg(ctx, 1)); v != nil {
		dataStr = v.AsString()
	}
	existing := ""
	if b, err := os.ReadFile(path); err == nil {
		existing = string(b)
	}
	if err := os.WriteFile(path, []byte(dataStr+existing), 0o644); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(len(dataStr)), nil
}

func fsName(ctx data.Context) (data.GetValue, data.Control) {
	base := filepath.Base(fsPath(ctx, 0))
	ext := filepath.Ext(base)
	return data.NewStringValue(strings.TrimSuffix(base, ext)), nil
}

func fsBasename(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(filepath.Base(fsPath(ctx, 0))), nil
}

func fsDirname(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(filepath.Dir(fsPath(ctx, 0))), nil
}

func fsExtension(ctx data.Context) (data.GetValue, data.Control) {
	ext := filepath.Ext(fsPath(ctx, 0))
	return data.NewStringValue(strings.TrimPrefix(ext, ".")), nil
}

func fsType(ctx data.Context) (data.GetValue, data.Control) {
	path := fsPath(ctx, 0)
	info, err := os.Stat(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if info.IsDir() {
		return data.NewStringValue("dir"), nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return data.NewStringValue("link"), nil
	}
	return data.NewStringValue("file"), nil
}

func fsGlob(ctx data.Context) (data.GetValue, data.Control) {
	pattern := fsPath(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return out, nil
	}
	for i, m := range matches {
		out.SetIntKey(i, data.NewStringValue(m))
	}
	return out, nil
}
