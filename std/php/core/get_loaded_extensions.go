package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type GetLoadedExtensionsFunction struct{}

func (*GetLoadedExtensionsFunction) GetName() string { return "get_loaded_extensions" }
func (*GetLoadedExtensionsFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "zend_extensions", 0, data.NewBoolValue(false), nil)}
}
func (*GetLoadedExtensionsFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "zend_extensions", 0, nil)}
}
func (*GetLoadedExtensionsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if v, ok := ctx.GetIndexValue(0); ok {
		if b, ok := v.(data.AsBool); ok {
			yes, _ := b.AsBool()
			if yes {
				return data.NewArrayValue(nil), nil
			}
		}
	}
	return loadedExtensions(ctx), nil
}

func loadedExtensions(ctx data.Context) *data.ArrayValue {
	var extensions []data.Value
	for _, name := range []string{"core", "standard", "date", "json", "pcre", "hash", "filter", "iconv", "mbstring", "ctype", "tokenizer", "session", "zlib", "openssl", "curl", "fileinfo", "phar", "spl", "reflection", "pdo", "intl", "dom", "simplexml", "zip"} {
		if extensionAvailable(ctx.GetVM(), name) {
			extensions = append(extensions, data.NewStringValue(name))
		}
	}
	return data.NewArrayValue(extensions).(*data.ArrayValue)
}
