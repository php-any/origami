package laravel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework"
	"github.com/php-any/origami/std/laravel/httpkernel"
	"github.com/php-any/origami/std/laravel/prompts"
	"github.com/php-any/origami/std/laravel/sentinel"
	serializableclosure "github.com/php-any/origami/std/laravel/serializable-closure"
	"github.com/php-any/origami/std/laravel/telescope"
	"github.com/php-any/origami/std/laravel/tinker"
)

// Load proposes vendor acceleration registrations to the contract-gated VM.
// The HTTP Kernel is provided by Composer; the explicit serve host adapter is
// loaded separately by vendoraccel on the real VM.
func Load(vm data.VM) {
	framework.Load(vm)
	serializableclosure.Load(vm)
	sentinel.Load(vm)
	tinker.Load(vm)
	prompts.Load(vm)
	telescope.Load(vm)

	httpkernel.Load(vm)
}
