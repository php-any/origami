package intl

import "github.com/php-any/origami/data"

// Load 注册 intl 扩展相关函数（如 grapheme_strlen、grapheme_substr）。
func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewGraphemeStrlenFunction(),
		NewGraphemeSubstrFunction(),
	} {
		vm.AddFunc(fun)
	}
}
