package gosupport

import (
	"html"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Illuminate 相关库在 Origami 上运行时常用的 PHP 内置函数桩，仅注册在 Laravel 示例 VM。

type RestoreErrorHandlerFunction struct{}

func (f *RestoreErrorHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
func (f *RestoreErrorHandlerFunction) GetName() string            { return "restore_error_handler" }
func (f *RestoreErrorHandlerFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *RestoreErrorHandlerFunction) GetIsStatic() bool          { return false }
func (f *RestoreErrorHandlerFunction) GetParams() []data.GetValue { return nil }
func (f *RestoreErrorHandlerFunction) GetVariables() []data.Variable {
	return nil
}
func (f *RestoreErrorHandlerFunction) GetReturnType() data.Types { return data.NewBaseType("bool") }

type StreamSetChunkSizeFunction struct{}

func (f *StreamSetChunkSizeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	size := valueAsInt(mustIndex(ctx, 1), 8192)
	if size < 1 {
		size = 8192
	}
	return data.NewIntValue(size), nil
}
func (f *StreamSetChunkSizeFunction) GetName() string            { return "stream_set_chunk_size" }
func (f *StreamSetChunkSizeFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *StreamSetChunkSizeFunction) GetIsStatic() bool          { return false }
func (f *StreamSetChunkSizeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "size", 1, nil, nil),
	}
}
func (f *StreamSetChunkSizeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
		node.NewVariable(nil, "size", 1, data.NewBaseType("int")),
	}
}
func (f *StreamSetChunkSizeFunction) GetReturnType() data.Types { return data.NewBaseType("int") }

type FileinodeFunction struct{}

func (f *FileinodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	path := valueAsString(mustIndex(ctx, 0))
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return data.NewIntValue(int(st.Ino)), nil
	}
	return data.NewIntValue(1), nil
}
func (f *FileinodeFunction) GetName() string            { return "fileinode" }
func (f *FileinodeFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *FileinodeFunction) GetIsStatic() bool          { return false }
func (f *FileinodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}
}
func (f *FileinodeFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "filename", 0, data.NewBaseType("string"))}
}
func (f *FileinodeFunction) GetReturnType() data.Types { return nil }

type ParseStrFunction struct{}

func (f *ParseStrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	s := valueAsString(mustIndex(ctx, 0))
	values, err := url.ParseQuery(s)
	if err != nil {
		values = url.Values{}
	}
	list := make([]*data.ZVal, 0, len(values))
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		key = strings.TrimSuffix(key, "[]")
		list = append(list, data.NewNamedZVal(key, data.NewStringValue(vals[len(vals)-1])))
	}
	parsed := &data.ArrayValue{List: list}

	resultVal, _ := ctx.GetIndexValue(1)
	if dest, ok := resultVal.(*data.ArrayValue); ok {
		dest.List = parsed.List
	} else if dest, ok := resultVal.(*data.ObjectValue); ok {
		keys := make([]string, 0)
		dest.RangeProperties(func(k string, _ data.Value) bool {
			keys = append(keys, k)
			return true
		})
		for _, k := range keys {
			dest.UnsetProperty(k)
		}
		for _, zv := range parsed.List {
			if zv != nil && zv.Name != "" {
				dest.SetProperty(zv.Name, zv.Value)
			}
		}
	} else if zv := ctx.GetIndexZVal(1); zv != nil {
		zv.Value = parsed
	}
	return data.NewNullValue(), nil
}
func (f *ParseStrFunction) GetName() string            { return "parse_str" }
func (f *ParseStrFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *ParseStrFunction) GetIsStatic() bool          { return false }
func (f *ParseStrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.NewBaseType("string")),
		node.NewParameterReference(nil, "result", 1, nil, data.NewBaseType("array")),
	}
}
func (f *ParseStrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "result", 1, data.NewBaseType("array")),
	}
}
func (f *ParseStrFunction) GetReturnType() data.Types { return nil }

type HtmlEntityDecodeFunction struct{}

func (f *HtmlEntityDecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(html.UnescapeString(valueAsString(mustIndex(ctx, 0)))), nil
}
func (f *HtmlEntityDecodeFunction) GetName() string            { return "html_entity_decode" }
func (f *HtmlEntityDecodeFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HtmlEntityDecodeFunction) GetIsStatic() bool          { return false }
func (f *HtmlEntityDecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "flags", 1, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "encoding", 2, node.NewStringLiteralByAst(nil, "UTF-8"), nil),
	}
}
func (f *HtmlEntityDecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "encoding", 2, data.NewBaseType("string")),
	}
}
func (f *HtmlEntityDecodeFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

type GetcwdFunction struct{}

func (f *GetcwdFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	wd, err := os.Getwd()
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(wd), nil
}
func (f *GetcwdFunction) GetName() string            { return "getcwd" }
func (f *GetcwdFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *GetcwdFunction) GetIsStatic() bool          { return false }
func (f *GetcwdFunction) GetParams() []data.GetValue { return nil }
func (f *GetcwdFunction) GetVariables() []data.Variable {
	return nil
}
func (f *GetcwdFunction) GetReturnType() data.Types { return data.NewBaseType("string") }
