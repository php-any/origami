package http

import (
	"github.com/php-any/origami/data"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

// Load：Request/Response + Json/Redirect/File/UploadedFile 常开。
func Load(vm data.VM) {
	httpfoundation.Load(vm)
	vm.AddClass(NewIlluminateRequestClass())
	vm.AddClass(NewIlluminateResponseClass())
	vm.AddClass(NewIlluminateJsonResponseClass())
	vm.AddClass(NewIlluminateRedirectResponseClass())
	vm.AddClass(NewIlluminateFileClass())
	vm.AddClass(NewIlluminateUploadedFileClass())
}
