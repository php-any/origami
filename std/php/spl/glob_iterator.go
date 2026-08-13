package spl

import (
	"errors"
	"path/filepath"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// GlobIteratorClass 实现 PHP GlobIterator（extends FilesystemIterator�?
type GlobIteratorClass struct {
	FilesystemIteratorClass
}

func NewGlobIteratorClass() *GlobIteratorClass {
	return &GlobIteratorClass{
		FilesystemIteratorClass: FilesystemIteratorClass{
			flags: FSI_DEFAULT_FLAGS,
		},
	}
}

func (c *GlobIteratorClass) GetName() string { return "GlobIterator" }

func (c *GlobIteratorClass) GetExtend() *string {
	parent := "FilesystemIterator"
	return &parent
}

func (c *GlobIteratorClass) GetImplements() []string {
	return []string{"SeekableIterator"}
}

func (c *GlobIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	clone := &GlobIteratorClass{
		FilesystemIteratorClass: FilesystemIteratorClass{
			path:     c.path,
			entries:  make([]string, len(c.entries)),
			iterator: c.iterator,
			flags:    c.flags,
		},
	}
	copy(clone.entries, c.entries)
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

func (c *GlobIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &GlobIteratorConstructMethod{instance: &c.FilesystemIteratorClass}, true
	}
	return c.FilesystemIteratorClass.GetMethod(name)
}

func (c *GlobIteratorClass) GetMethods() []data.Method {
	methods := c.FilesystemIteratorClass.GetMethods()
	out := make([]data.Method, 0, len(methods)+1)
	out = append(out, &GlobIteratorConstructMethod{instance: &c.FilesystemIteratorClass})
	for _, m := range methods {
		if m.GetName() != "__construct" {
			out = append(out, m)
		}
	}
	return out
}

func (c *GlobIteratorClass) GetConstruct() data.Method {
	return &GlobIteratorConstructMethod{instance: &c.FilesystemIteratorClass}
}

func (c *GlobIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	return c.FilesystemIteratorClass.GetStaticProperty(name)
}

// GlobIteratorConstructMethod 实现 GlobIterator::__construct
type GlobIteratorConstructMethod struct {
	instance *FilesystemIteratorClass
}

func (m *GlobIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *GlobIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *GlobIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *GlobIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *GlobIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, data.String{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(FSI_DEFAULT_FLAGS), data.Int{}),
	}
}
func (m *GlobIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.String{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}
func (m *GlobIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	if patternValue == nil {
		return nil, utils.NewThrow(errors.New("GlobIterator::__construct() expects parameter 1 to be string"))
	}
	pattern := patternValue.AsString()
	if pattern == "" {
		return nil, utils.NewThrow(errors.New("GlobIterator::__construct(): Pattern must not be empty"))
	}

	flags := FSI_DEFAULT_FLAGS
	if flagsValue, _ := ctx.GetIndexValue(1); flagsValue != nil {
		if iv, ok := flagsValue.(interface{ AsInt() (int, error) }); ok {
			if n, err := iv.AsInt(); err == nil {
				flags = n
			}
		}
	}

	matches, err := filepath.Glob(filepath.FromSlash(pattern))
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if flags&FSI_SKIP_DOTS == 0 {
		// glob 结果不含 . / ..
	}
	sort.Strings(matches)

	// 提取目录路径，entries 只存储文件名（与 DirectoryIterator/FilesystemIterator 保持一致）
	dir := "."
	if len(matches) > 0 {
		first := matches[0]
		if abs, err := filepath.Abs(first); err == nil {
			first = abs
		}
		dir = filepath.Dir(first)
	}

	names := make([]string, len(matches))
	for i, match := range matches {
		match = filepath.Clean(match)
		if abs, err := filepath.Abs(match); err == nil {
			match = abs
		}
		names[i] = filepath.Base(match)
	}

	m.instance.path = dir
	m.instance.flags = flags
	m.instance.entries = names
	m.instance.iterator = 0
	return nil, nil
}
