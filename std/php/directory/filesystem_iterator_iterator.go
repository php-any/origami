package directory

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// FilesystemIteratorCurrentMethod 实现 FilesystemIterator::current
type FilesystemIteratorCurrentMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorCurrentMethod) GetName() string               { return "current" }
func (m *FilesystemIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorCurrentMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorCurrentMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }

func (m *FilesystemIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// CURRENT_AS_PATHNAME：返回路径名字符串
	if m.instance.flags&FSI_CURRENT_AS_PATHNAME != 0 {
		return data.NewStringValue(m.instance.GetPathname()), nil
	}
	// CURRENT_AS_SELF 或 CURRENT_AS_FILEINFO（默认）：返回 $this（ClassValue）
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return objCtx.ClassValue, nil
	}
	return data.NewNullValue(), nil
}

// FilesystemIteratorKeyMethod 实现 FilesystemIterator::key
type FilesystemIteratorKeyMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorKeyMethod) GetName() string               { return "key" }
func (m *FilesystemIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorKeyMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorKeyMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }

func (m *FilesystemIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.instance.KeyStr()), nil
}

// FilesystemIteratorNextMethod 实现 FilesystemIterator::next
type FilesystemIteratorNextMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorNextMethod) GetName() string               { return "next" }
func (m *FilesystemIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorNextMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorNextMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorNextMethod) GetReturnType() data.Types     { return nil }

func (m *FilesystemIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.path == "" {
		return nil, utils.NewThrow(errors.New("FilesystemIterator not initialized"))
	}
	m.instance.iterator++
	return nil, nil
}

// FilesystemIteratorRewindMethod 实现 FilesystemIterator::rewind
type FilesystemIteratorRewindMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *FilesystemIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorRewindMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorRewindMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorRewindMethod) GetReturnType() data.Types     { return nil }

func (m *FilesystemIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.instance.path == "" {
		return nil, utils.NewThrow(errors.New("FilesystemIterator not initialized"))
	}
	m.instance.iterator = 0
	return nil, nil
}

// FilesystemIteratorValidMethod 实现 FilesystemIterator::valid
type FilesystemIteratorValidMethod struct {
	instance *FilesystemIteratorClass
}

func (m *FilesystemIteratorValidMethod) GetName() string               { return "valid" }
func (m *FilesystemIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilesystemIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *FilesystemIteratorValidMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *FilesystemIteratorValidMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *FilesystemIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *FilesystemIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	valid := m.instance.iterator >= 0 && m.instance.iterator < len(m.instance.entries)
	return data.NewBoolValue(valid), nil
}
