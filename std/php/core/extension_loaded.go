package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ExtensionLoadedFunction 实现 extension_loaded 函数
type ExtensionLoadedFunction struct{}

func NewExtensionLoadedFunction() data.FuncStmt {
	return &ExtensionLoadedFunction{}
}

func (f *ExtensionLoadedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	extension, _ := ctx.GetIndexValue(0)
	name := ""
	if extension != nil {
		name = strings.ToLower(extension.AsString())
	}

	return data.NewBoolValue(extensionAvailable(ctx.GetVM(), name)), nil
}

// Native optional extensions must opt in. Unknown names must never select
// unavailable fast paths (for example Xdebug's process restart machinery).
func extensionAvailable(vm data.VM, name string) bool {
	if name == "core" || name == "standard" {
		return true
	}
	functions := map[string]string{
		"date": "date", "json": "json_encode", "pcre": "preg_match",
		"hash": "hash", "filter": "filter_var", "iconv": "iconv",
		"mbstring": "mb_strlen", "ctype": "ctype_alpha", "tokenizer": "token_get_all",
		"session": "session_start", "zlib": "gzdecode", "openssl": "openssl_encrypt",
		"curl": "curl_init", "fileinfo": "finfo_open",
	}
	if fn, ok := functions[name]; ok {
		_, found := vm.GetFunc(fn)
		return found
	}
	classes := map[string]string{"phar": "Phar", "spl": "SplFileInfo", "reflection": "ReflectionClass", "pdo": "PDO", "intl": "NumberFormatter", "dom": "DOMDocument", "simplexml": "SimpleXMLElement", "zip": "ZipArchive"}
	if class, ok := classes[name]; ok {
		_, found := vm.GetClass(class)
		return found
	}
	return false
}

func (f *ExtensionLoadedFunction) GetName() string {
	return "extension_loaded"
}

var extensionLoadedFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "extension", 0, nil, data.String{}),
}

func (f *ExtensionLoadedFunction) GetParams() []data.GetValue {
	return extensionLoadedFunctionGetParams
}

var extensionLoadedFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "extension", 0, data.String{}),
}

func (f *ExtensionLoadedFunction) GetVariables() []data.Variable {
	return extensionLoadedFunctionGetVariables
}
