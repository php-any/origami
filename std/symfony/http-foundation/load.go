package httpfoundation

import "github.com/php-any/origami/data"

// Load 注册本包 FQCN（Symfony\Component\HttpFoundation\*）。
func Load(vm data.VM) {
	vm.AddClass(NewParameterBagClass())
	vm.AddClass(NewHeaderBagClass())
	vm.AddClass(NewInputBagClass())
	vm.AddClass(NewServerBagClass())
	vm.AddClass(NewFileBagClass())
	vm.AddClass(NewResponseHeaderBagClass())
	vm.AddClass(NewCookieClass())
	vm.AddClass(NewSymfonyRequestClass())
	vm.AddClass(NewSymfonyResponseClass())
	vm.AddClass(NewJsonResponseClass())
	vm.AddClass(NewRedirectResponseClass())
	vm.AddClass(NewFileClass())
	vm.AddClass(NewUploadedFileClass())
	vm.AddClass(NewRequestStackClass())
	vm.AddClass(NewStreamedResponseClass())
	vm.AddClass(NewBinaryFileResponseClass())
}
