package uid

import "github.com/php-any/origami/data"

var (
	_ data.ClassStmt         = (*UuidClass)(nil)
	_ data.GetStaticMethod   = (*UuidClass)(nil)
	_ data.GetStaticProperty = (*UuidClass)(nil)
	_ data.ClassStmt         = (*UlidClass)(nil)
	_ data.GetStaticMethod   = (*UlidClass)(nil)
	_ data.GetStaticProperty = (*UlidClass)(nil)
)

// Load 注册 Symfony Uid 原生类（Uuid / Ulid）。不注册 UuidV* / AbstractUid / Nil* / Max*，以便 PHP 子类继承原生 Uuid。
func Load(vm data.VM) {
	vm.AddClass(NewUuidClass())
	vm.AddClass(NewUlidClass())
}
