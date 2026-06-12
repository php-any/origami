package container

import "github.com/php-any/origami/data"

// CompiledSingletonValue 构建预编译的 #[Singleton] 注解实例（供 compile 子命令使用）
func CompiledSingletonValue() *data.ClassValue {
	sc := &SingletonAnnotationClass{construct: &SingletonAnnotationConstructMethod{}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: sc}
}

// CompiledScopedValue 构建预编译的 #[Scoped] 注解实例（供 compile 子命令使用）
func CompiledScopedValue() *data.ClassValue {
	sc := &ScopedAnnotationClass{construct: &ScopedAnnotationConstructMethod{}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: sc}
}

// CompiledComponentValue 构建预编译的 #[Component] 注解实例（供 compile 子命令使用）
func CompiledComponentValue() *data.ClassValue {
	cc := &ComponentClass{construct: &ComponentConstructMethod{}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: cc}
}
