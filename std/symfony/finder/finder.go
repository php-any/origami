package finder

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const finderClassName = "Symfony\\Component\\Finder\\Finder"
const splFileInfoClassName = "Symfony\\Component\\Finder\\SplFileInfo"

// FinderClass 实现 Symfony Finder 热路径子集。
type FinderClass struct {
	node.Node
	methods map[string]data.Method
}

func NewFinderClass() data.ClassStmt {
	c := &FinderClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = &finderMethod{name: "__construct", fn: finderConstruct}
	c.methods["create"] = &finderMethod{name: "create", static: true, fn: finderCreate}
	c.methods["in"] = &finderMethod{name: "in", params: []string{"dirs"}, fn: finderIn}
	c.methods["files"] = &finderMethod{name: "files", fn: finderFiles}
	c.methods["directories"] = &finderMethod{name: "directories", fn: finderDirectories}
	c.methods["name"] = &finderMethod{name: "name", params: []string{"patterns"}, fn: finderName}
	c.methods["notname"] = &finderMethod{name: "notName", params: []string{"patterns"}, fn: finderNotName}
	c.methods["path"] = &finderMethod{name: "path", params: []string{"patterns"}, fn: finderPath}
	c.methods["exclude"] = &finderMethod{name: "exclude", params: []string{"dirs"}, fn: finderExclude}
	c.methods["depth"] = &finderMethod{name: "depth", params: []string{"levels"}, fn: finderDepth}
	c.methods["ignoredotfiles"] = &finderMethod{name: "ignoreDotFiles", params: []string{"ignoreDotFiles"}, fn: finderIgnoreDotFiles}
	c.methods["sortbyname"] = &finderMethod{name: "sortByName", params: []string{"useNaturalSort"}, fn: finderSortByName}
	c.methods["getiterator"] = &finderMethod{name: "getIterator", fn: finderGetIterator}
	c.methods["count"] = &finderMethod{name: "count", fn: finderCount}
	return c
}

func (c *FinderClass) GetName() string                          { return finderClassName }
func (c *FinderClass) GetExtend() *string                       { return nil }
func (c *FinderClass) GetImplements() []string                  { return []string{"IteratorAggregate", "Countable"} }
func (c *FinderClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *FinderClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "dirs", "private", false, data.NewArrayValue(nil)),
		node.NewProperty(nil, "names", "private", false, data.NewArrayValue(nil)),
		node.NewProperty(nil, "notNames", "private", false, data.NewArrayValue(nil)),
		node.NewProperty(nil, "paths", "private", false, data.NewArrayValue(nil)),
		node.NewProperty(nil, "excludes", "private", false, data.NewArrayValue(nil)),
		node.NewProperty(nil, "mode", "private", false, data.NewStringValue("any")),
		node.NewProperty(nil, "maxDepth", "private", false, data.NewIntValue(-1)),
		node.NewProperty(nil, "ignoreDotFiles", "private", false, data.NewBoolValue(true)),
		node.NewProperty(nil, "sortByName", "private", false, data.NewBoolValue(false)),
	}
}
func (c *FinderClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *FinderClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *FinderClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *FinderClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *FinderClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

type finderMethod struct {
	name   string
	params []string
	static bool
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *finderMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *finderMethod) GetName() string                                     { return m.name }
func (m *finderMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *finderMethod) GetIsStatic() bool                                   { return m.static }
func (m *finderMethod) GetReturnType() data.Types                           { return nil }
func (m *finderMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p, i, nil, nil)
	}
	return out
}
func (m *finderMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}

func finderSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func finderConstruct(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func finderCreate(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(NewFinderClass(), ctx.CreateBaseContext()), nil
}

func appendStringProp(cv *data.ClassValue, prop string, val data.Value) {
	cur, _ := cv.GetProperty(prop)
	arr, _ := cur.(*data.ArrayValue)
	if arr == nil {
		arr = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	appendOne := func(p string) {
		if p == "" {
			return
		}
		arr.List = append(arr.List, data.NewZVal(data.NewStringValue(p)))
	}
	if av, ok := val.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil {
				appendOne(z.Value.AsString())
			}
		}
	} else if val != nil {
		appendOne(val.AsString())
	}
	_ = cv.SetProperty(prop, arr)
}

func finderIn(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	dirs, _ := ctx.GetIndexValue(0)
	appendStringProp(cv, "dirs", dirs)
	return cv, nil
}

func finderFiles(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	_ = cv.SetProperty("mode", data.NewStringValue("files"))
	return cv, nil
}

func finderDirectories(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	_ = cv.SetProperty("mode", data.NewStringValue("directories"))
	return cv, nil
}

func finderName(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	pat, _ := ctx.GetIndexValue(0)
	appendStringProp(cv, "names", pat)
	return cv, nil
}

func finderNotName(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	pat, _ := ctx.GetIndexValue(0)
	appendStringProp(cv, "notNames", pat)
	return cv, nil
}

func finderPath(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	pat, _ := ctx.GetIndexValue(0)
	appendStringProp(cv, "paths", pat)
	return cv, nil
}

func finderExclude(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	dirs, _ := ctx.GetIndexValue(0)
	appendStringProp(cv, "excludes", dirs)
	return cv, nil
}

func finderDepth(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	levels, _ := ctx.GetIndexValue(0)
	maxDepth := -1
	if levels != nil {
		s := strings.TrimSpace(levels.AsString())
		if n, err := strconv.Atoi(s); err == nil {
			maxDepth = n
		} else if strings.HasPrefix(s, "<=") {
			if n, err := strconv.Atoi(strings.TrimSpace(s[2:])); err == nil {
				maxDepth = n
			}
		} else if strings.HasPrefix(s, "<") {
			if n, err := strconv.Atoi(strings.TrimSpace(s[1:])); err == nil {
				maxDepth = n - 1
			}
		} else if strings.HasPrefix(s, ">=") || strings.HasPrefix(s, ">") {
			// 最小深度：热路径少见，忽略下限仅记 max=-1 全扫
			maxDepth = -1
		}
	}
	_ = cv.SetProperty("maxDepth", data.NewIntValue(maxDepth))
	return cv, nil
}

func finderIgnoreDotFiles(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	v, _ := ctx.GetIndexValue(0)
	ignore := true
	if v != nil {
		if b, ok := v.(data.AsBool); ok {
			ignore, _ = b.AsBool()
		}
	}
	_ = cv.SetProperty("ignoreDotFiles", data.NewBoolValue(ignore))
	return cv, nil
}

func finderSortByName(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	_ = cv.SetProperty("sortByName", data.NewBoolValue(true))
	return cv, nil
}

func stringListProp(cv *data.ClassValue, prop string) []string {
	v, _ := cv.GetProperty(prop)
	out := []string{}
	if av, ok := v.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil {
				out = append(out, z.Value.AsString())
			}
		}
	}
	return out
}

func matchAny(patterns []string, name string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, name); matched {
			return true
		}
		// Symfony 也接受正则 /.../
		if len(p) >= 2 && p[0] == '/' && p[len(p)-1] == '/' {
			// 简化：不引入 regexp 全量，跳过复杂正则
			continue
		}
	}
	return false
}

func matchNone(patterns []string, name string) bool {
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, name); matched {
			return false
		}
	}
	return true
}

type foundFile struct {
	path             string
	relativePath     string
	relativePathname string
}

func finderCollect(cv *data.ClassValue) []foundFile {
	dirs := stringListProp(cv, "dirs")
	names := stringListProp(cv, "names")
	notNames := stringListProp(cv, "notNames")
	paths := stringListProp(cv, "paths")
	excludes := stringListProp(cv, "excludes")
	modeV, _ := cv.GetProperty("mode")
	mode := "any"
	if modeV != nil {
		mode = modeV.AsString()
	}
	maxDepth := -1
	if md, _ := cv.GetProperty("maxDepth"); md != nil {
		if iv, ok := md.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				maxDepth = n
			}
		}
	}
	ignoreDot := true
	if id, _ := cv.GetProperty("ignoreDotFiles"); id != nil {
		if b, ok := id.(data.AsBool); ok {
			ignoreDot, _ = b.AsBool()
		}
	}
	doSort := false
	if sb, _ := cv.GetProperty("sortByName"); sb != nil {
		if b, ok := sb.(data.AsBool); ok {
			doSort, _ = b.AsBool()
		}
	}

	excludeSet := map[string]struct{}{}
	for _, e := range excludes {
		excludeSet[filepath.Clean(e)] = struct{}{}
	}

	var out []foundFile
	for _, root := range dirs {
		root = filepath.Clean(root)
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil {
				return nil
			}
			if path == root {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)
			depth := strings.Count(rel, "/")
			if info.IsDir() {
				// exclude 目录名
				base := info.Name()
				if _, ok := excludeSet[base]; ok {
					return filepath.SkipDir
				}
				if _, ok := excludeSet[filepath.Clean(rel)]; ok {
					return filepath.SkipDir
				}
				if ignoreDot && strings.HasPrefix(base, ".") {
					return filepath.SkipDir
				}
				if maxDepth >= 0 && depth > maxDepth {
					return filepath.SkipDir
				}
			}
			if maxDepth >= 0 && depth > maxDepth {
				return nil
			}
			if ignoreDot && strings.HasPrefix(info.Name(), ".") {
				return nil
			}
			if mode == "files" && info.IsDir() {
				return nil
			}
			if mode == "directories" && !info.IsDir() {
				return nil
			}
			base := info.Name()
			if len(names) > 0 && !matchAny(names, base) {
				return nil
			}
			if !matchNone(notNames, base) {
				return nil
			}
			if len(paths) > 0 {
				okPath := false
				for _, p := range paths {
					if matched, _ := filepath.Match(p, rel); matched {
						okPath = true
						break
					}
					if strings.Contains(rel, p) {
						okPath = true
						break
					}
				}
				if !okPath {
					return nil
				}
			}
			relDir := filepath.ToSlash(filepath.Dir(rel))
			if relDir == "." {
				relDir = ""
			}
			out = append(out, foundFile{
				path:             path,
				relativePath:     relDir,
				relativePathname: rel,
			})
			return nil
		})
	}
	if doSort {
		sort.Slice(out, func(i, j int) bool {
			return out[i].relativePathname < out[j].relativePathname
		})
	}
	return out
}

func finderResultsArray(ctx data.Context, cv *data.ClassValue) *data.ArrayValue {
	files := finderCollect(cv)
	vals := make([]data.Value, 0, len(files))
	for _, f := range files {
		vals = append(vals, newSplFileInfo(ctx, f.path, f.relativePath, f.relativePathname))
	}
	return data.NewArrayValue(vals).(*data.ArrayValue)
}

func finderGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	// ArrayValue 实现 data.Iterator，foreach / iterator_to_array 均可直接消费。
	return finderResultsArray(ctx, cv), nil
}

func finderCount(ctx data.Context) (data.GetValue, data.Control) {
	cv := finderSelf(ctx)
	return data.NewIntValue(len(finderCollect(cv))), nil
}

// --- SplFileInfo ---

type SplFileInfoClass struct {
	node.Node
	methods map[string]data.Method
}

func NewSplFileInfoClass() data.ClassStmt {
	c := &SplFileInfoClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = &finderMethod{name: "__construct", params: []string{"file", "relativePath", "relativePathname"}, fn: splConstruct}
	c.methods["getrelativepath"] = &finderMethod{name: "getRelativePath", fn: splGetRelativePath}
	c.methods["getrelativepathname"] = &finderMethod{name: "getRelativePathname", fn: splGetRelativePathname}
	c.methods["getpathname"] = &finderMethod{name: "getPathname", fn: splGetPathname}
	c.methods["getrealpath"] = &finderMethod{name: "getRealPath", fn: splGetRealPath}
	c.methods["getfilename"] = &finderMethod{name: "getFilename", fn: splGetFilename}
	c.methods["getpath"] = &finderMethod{name: "getPath", fn: splGetPath}
	c.methods["isfile"] = &finderMethod{name: "isFile", fn: splIsFile}
	c.methods["isdir"] = &finderMethod{name: "isDir", fn: splIsDir}
	c.methods["getcontents"] = &finderMethod{name: "getContents", fn: splGetContents}
	c.methods["getfilenamewithoutextension"] = &finderMethod{name: "getFilenameWithoutExtension", fn: splGetFilenameWithoutExtension}
	c.methods["__tostring"] = &finderMethod{name: "__toString", fn: splGetPathname}
	return c
}

func (c *SplFileInfoClass) GetName() string { return splFileInfoClassName }
func (c *SplFileInfoClass) GetExtend() *string {
	parent := "SplFileInfo"
	return &parent
}
func (c *SplFileInfoClass) GetImplements() []string                  { return nil }
func (c *SplFileInfoClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *SplFileInfoClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "pathname", "private", false, data.NewStringValue("")),
		node.NewProperty(nil, "relativePath", "private", false, data.NewStringValue("")),
		node.NewProperty(nil, "relativePathname", "private", false, data.NewStringValue("")),
	}
}
func (c *SplFileInfoClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *SplFileInfoClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *SplFileInfoClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *SplFileInfoClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}

func splSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func splConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := splSelf(ctx)
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		_ = cv.SetProperty("pathname", data.NewStringValue(v.AsString()))
	}
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		_ = cv.SetProperty("relativePath", data.NewStringValue(v.AsString()))
	}
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		_ = cv.SetProperty("relativePathname", data.NewStringValue(v.AsString()))
	}
	return data.NewNullValue(), nil
}

func newSplFileInfo(ctx data.Context, pathname, relPath, relPathname string) *data.ClassValue {
	cv := data.NewClassValue(NewSplFileInfoClass(), ctx.CreateBaseContext())
	_ = cv.SetProperty("pathname", data.NewStringValue(pathname))
	_ = cv.SetProperty("relativePath", data.NewStringValue(relPath))
	_ = cv.SetProperty("relativePathname", data.NewStringValue(relPathname))
	return cv
}

// NewSplFileInfoInstance 供 illuminate/filesystem 等构造 Finder SplFileInfo。
func NewSplFileInfoInstance(ctx data.Context, pathname, relPath, relPathname string) *data.ClassValue {
	return newSplFileInfo(ctx, pathname, relPath, relPathname)
}

func splProp(cv *data.ClassValue, name string) string {
	v, _ := cv.GetProperty(name)
	if v == nil {
		return ""
	}
	return v.AsString()
}

func splGetRelativePath(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(splProp(splSelf(ctx), "relativePath")), nil
}
func splGetRelativePathname(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(splProp(splSelf(ctx), "relativePathname")), nil
}
func splGetPathname(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(splProp(splSelf(ctx), "pathname")), nil
}
func splGetRealPath(ctx data.Context) (data.GetValue, data.Control) {
	p := splProp(splSelf(ctx), "pathname")
	if abs, err := filepath.Abs(p); err == nil {
		return data.NewStringValue(abs), nil
	}
	return data.NewStringValue(p), nil
}
func splGetFilename(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(filepath.Base(splProp(splSelf(ctx), "pathname"))), nil
}
func splGetPath(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(filepath.Dir(splProp(splSelf(ctx), "pathname"))), nil
}
func splIsFile(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(splProp(splSelf(ctx), "pathname"))
	return data.NewBoolValue(err == nil && !info.IsDir()), nil
}
func splIsDir(ctx data.Context) (data.GetValue, data.Control) {
	info, err := os.Stat(splProp(splSelf(ctx), "pathname"))
	return data.NewBoolValue(err == nil && info.IsDir()), nil
}
func splGetContents(ctx data.Context) (data.GetValue, data.Control) {
	b, err := os.ReadFile(splProp(splSelf(ctx), "pathname"))
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewStringValue(string(b)), nil
}
func splGetFilenameWithoutExtension(ctx data.Context) (data.GetValue, data.Control) {
	base := filepath.Base(splProp(splSelf(ctx), "pathname"))
	ext := filepath.Ext(base)
	return data.NewStringValue(strings.TrimSuffix(base, ext)), nil
}
