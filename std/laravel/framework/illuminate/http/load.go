package http

import (
	"github.com/php-any/origami/data"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

// Load：Request/Response + Json/Redirect/File/UploadedFile 常开。
func Load(vm data.VM) {
	httpfoundation.Load(vm)
	vm.AddClass(illuminateRequestClassStmt())
	vm.AddClass(NewIlluminateResponseClass())
	vm.AddClass(NewIlluminateJsonResponseClass())
	// RedirectResponse 不打 Go 桩，由 vendor 的 PHP 原文件解析：
	// redirect_response.go 那份桩只实现了 __construct/getTargetUrl/setTargetUrl/status/with，
	// 少了 setSession/getSession/setRequest/withInput/withErrors/withFragment/onlyInput/exceptInput/__call。
	// 桩一旦注册，ClassPathManager 就认为该类已存在、不再解析 vendor 文件，
	// 于是 Redirector::createRedirect() 里那句 $redirect->setSession($this->session)
	// 会以「类(Illuminate\Http\RedirectResponse)不存在对应函数(setSession)」抛错；
	// Handler::unauthenticated() 因此渲染失败走进 fallback，未登录访问 /admin 从 302 变成 500。
	// 实测：去掉这行注册后，vendor 原文件能解析出全部 16 个自身方法与全部 trait 方法。
	vm.AddClass(NewIlluminateFileClass())
	vm.AddClass(NewIlluminateUploadedFileClass())
}
