package iconv

import "github.com/php-any/origami/data"

// Load 注册所有 iconv 相关函数。
// 未来如果新增 iconv_strlen、iconv_substr 等函数，都在此集中注册。
func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewIconvFunction(),
		NewIconvStrlenFunction(),
		NewIconvSubstrFunction(),
		NewIconvStrposFunction(),
		NewIconvStrrposFunction(),
	} {
		vm.AddFunc(fun)
	}
}
