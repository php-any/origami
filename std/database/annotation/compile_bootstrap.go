package annotation

import "github.com/php-any/origami/data"

// CompiledTableValue 构建预编译的 #[Table] 注解实例（供 compile 子命令使用）
func CompiledTableValue(name string) *data.ClassValue {
	tc := &TableClass{name: name}
	tc.construct = &TableConstructMethod{tableClass: tc}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: tc}
}
