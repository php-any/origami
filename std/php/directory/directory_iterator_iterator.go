package directory

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// DirectoryIteratorCurrentMethod 实现 DirectoryIterator::current (Iterator 接口)
type DirectoryIteratorCurrentMethod struct{}

func (m *DirectoryIteratorCurrentMethod) GetName() string               { return "current" }
func (m *DirectoryIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorCurrentMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorCurrentMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }

func (m *DirectoryIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// PHP 的 DirectoryIterator::current() 返回 $this，即 DirectoryIterator 对象本身
	// 这样在 foreach 循环中，$value 就是一个 DirectoryIterator 对象，可以调用 getBasename() 等方法
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 返回当前的 DirectoryIterator 对象本身（ClassValue，不是 ObjectValue）
		return objCtx.ClassValue, nil
	}
	return data.NewNullValue(), nil
}

// DirectoryIteratorKeyMethod 实现 DirectoryIterator::key (Iterator 接口)
type DirectoryIteratorKeyMethod struct{}

func (m *DirectoryIteratorKeyMethod) GetName() string               { return "key" }
func (m *DirectoryIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorKeyMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorKeyMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }

func (m *DirectoryIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(iterData.Key()), nil
}

// DirectoryIteratorNextMethod 实现 DirectoryIterator::next (Iterator 接口)
type DirectoryIteratorNextMethod struct{}

func (m *DirectoryIteratorNextMethod) GetName() string               { return "next" }
func (m *DirectoryIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorNextMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorNextMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorNextMethod) GetReturnType() data.Types     { return nil }

func (m *DirectoryIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return nil, utils.NewThrow(errors.New("DirectoryIterator not initialized"))
	}
	iterData.Next()
	return nil, nil
}

// DirectoryIteratorRewindMethod 实现 DirectoryIterator::rewind (Iterator 接口)
type DirectoryIteratorRewindMethod struct{}

func (m *DirectoryIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *DirectoryIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorRewindMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorRewindMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorRewindMethod) GetReturnType() data.Types     { return nil }

func (m *DirectoryIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return nil, utils.NewThrow(errors.New("DirectoryIterator not initialized"))
	}
	iterData.Rewind()
	return nil, nil
}

// DirectoryIteratorValidMethod 实现 DirectoryIterator::valid (Iterator 接口)
type DirectoryIteratorValidMethod struct{}

func (m *DirectoryIteratorValidMethod) GetName() string               { return "valid" }
func (m *DirectoryIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DirectoryIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *DirectoryIteratorValidMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DirectoryIteratorValidMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *DirectoryIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *DirectoryIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterData, ok := getDirectoryIteratorInfo(ctx)
	if !ok || iterData == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(iterData.Valid()), nil
}
