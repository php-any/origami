package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"os"
	"path/filepath"
	"strings"
)

type autoloadExtensionHost interface {
	AutoloadExtensions() string
	SetAutoloadExtensions(string)
}

func autoloadExtensions(ctx data.Context) string {
	if host, ok := ctx.GetVM().(autoloadExtensionHost); ok {
		return host.AutoloadExtensions()
	}
	return ".inc,.php"
}

type SplAutoloadFunction struct{}

var splAutoloadParams = []data.GetValue{
	node.NewParameter(nil, "class", 0, nil, data.TypeString),
	node.NewParameter(nil, "file_extensions", 1, data.NewNullValue(), data.NewDeclaredUnionType([]data.Types{data.TypeString, data.TypeNull})),
}
var splAutoloadVariables = []data.Variable{node.NewVariable(nil, "class", 0, data.TypeString), node.NewVariable(nil, "file_extensions", 1, data.TypeMixed)}

func (*SplAutoloadFunction) GetName() string               { return "spl_autoload" }
func (*SplAutoloadFunction) GetParams() []data.GetValue    { return splAutoloadParams }
func (*SplAutoloadFunction) GetVariables() []data.Variable { return splAutoloadVariables }
func (*SplAutoloadFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	name := value.AsString()
	if _, found := ctx.GetVM().GetClass(name); found {
		return data.NewNullValue(), nil
	}
	if _, found := ctx.GetVM().GetInterface(name); found {
		return data.NewNullValue(), nil
	}
	extensions := autoloadExtensions(ctx)
	if value, _ := ctx.GetIndexValue(1); value != nil {
		if _, isNull := value.(*data.NullValue); !isNull {
			extensions = value.AsString()
		}
	}
	paths := []string{"."}
	if include, ok := core.IniGetInContext(ctx, "include_path"); ok {
		paths = filepath.SplitList(include)
	}
	base := strings.ReplaceAll(data.CanonicalTypeName(strings.TrimPrefix(name, "\\")), "\\", string(filepath.Separator))
	for _, extension := range strings.Split(extensions, ",") {
		for _, path := range paths {
			data.CheckRequest(ctx.GoContext())
			filename := filepath.Join(path, base+extension)
			if stat, err := os.Stat(filename); err != nil || stat.IsDir() {
				continue
			}
			if _, ctl := node.IncludeCore(ctx, data.NewStringValue(filename), false, false, nil); ctl != nil {
				return nil, ctl
			}
			if _, found := ctx.GetVM().GetClass(name); found {
				return data.NewNullValue(), nil
			}
			if _, found := ctx.GetVM().GetInterface(name); found {
				return data.NewNullValue(), nil
			}
		}
	}
	return data.NewNullValue(), nil
}

type SplAutoloadExtensionsFunction struct{}

var splAutoloadExtensionsParams = []data.GetValue{node.NewParameter(nil, "file_extensions", 0, data.NewNullValue(), data.NewDeclaredUnionType([]data.Types{data.TypeString, data.TypeNull}))}
var splAutoloadExtensionsVariables = []data.Variable{node.NewVariable(nil, "file_extensions", 0, data.TypeMixed)}

func (*SplAutoloadExtensionsFunction) GetName() string            { return "spl_autoload_extensions" }
func (*SplAutoloadExtensionsFunction) GetParams() []data.GetValue { return splAutoloadExtensionsParams }
func (*SplAutoloadExtensionsFunction) GetVariables() []data.Variable {
	return splAutoloadExtensionsVariables
}
func (*SplAutoloadExtensionsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if value, _ := ctx.GetIndexValue(0); value != nil {
		if _, isNull := value.(*data.NullValue); !isNull {
			if host, ok := ctx.GetVM().(autoloadExtensionHost); ok {
				host.SetAutoloadExtensions(value.AsString())
			}
		}
	}
	return data.NewStringValue(autoloadExtensions(ctx)), nil
}
